package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/Cotary/YAPI-MCP/internal/config"
	"github.com/Cotary/YAPI-MCP/internal/yapi"
)

// Register 注册所有 MCP Resources 到服务器
//
// MCP Resources 是只读数据源，让 LLM 可以按需读取结构化信息。
// 这里演示两种资源类型：
//   - 静态资源 (Resource)：固定 URI，如项目列表
//   - 动态资源 (ResourceTemplate)：URI 含参数模板，如具体接口文档
func Register(s *server.MCPServer, cfg *config.Config) {
	registerProjectListResource(s, cfg)
	registerInterfaceDocTemplate(s, cfg)
}

// registerProjectListResource 注册静态资源：已配置项目列表
func registerProjectListResource(s *server.MCPServer, cfg *config.Config) {
	resource := mcp.NewResource(
		"yapi://projects",
		"YAPI 项目列表",
		mcp.WithResourceDescription("所有已配置的 YAPI 项目的概览信息"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(resource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		type projectItem struct {
			ProjectID   int    `json:"project_id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			BaseURL     string `json:"base_url"`
		}

		items := make([]projectItem, 0, len(cfg.Projects))
		for _, p := range cfg.Projects {
			items = append(items, projectItem{
				ProjectID:   p.ProjectID,
				Name:        p.Name,
				Description: p.Description,
				BaseURL:     cfg.GetBaseURL(&p),
			})
		}

		data, _ := json.MarshalIndent(items, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "yapi://projects",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

// registerInterfaceDocTemplate 注册动态资源模板：接口文档
// URI 格式: yapi://projects/{projectId}/interfaces/{interfaceId}
func registerInterfaceDocTemplate(s *server.MCPServer, cfg *config.Config) {
	template := mcp.NewResourceTemplate(
		"yapi://projects/{projectId}/interfaces/{interfaceId}",
		"YAPI 接口文档",
		mcp.WithTemplateDescription("获取指定项目中某个接口的完整文档"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(template, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		projectID, interfaceID, err := parseResourceURI(req.Params.URI)
		if err != nil {
			return nil, err
		}

		p := cfg.FindProject(projectID)
		if p == nil {
			return nil, fmt.Errorf("未找到 project_id=%d 的配置", projectID)
		}

		client := yapi.NewClient(cfg.GetBaseURL(p), p.Token, p.ProjectID, cfg.SkipTLSVerify)
		detail, err := client.GetInterfaceDetail(interfaceID)
		if err != nil {
			return nil, fmt.Errorf("获取接口详情失败: %w", err)
		}

		data, _ := json.MarshalIndent(detail, "", "  ")
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

func parseResourceURI(uri string) (projectID, interfaceID int, err error) {
	// yapi://projects/123/interfaces/456
	trimmed := strings.TrimPrefix(uri, "yapi://projects/")
	parts := strings.Split(trimmed, "/interfaces/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("URI 格式错误: %s", uri)
	}
	projectID, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("projectId 无效: %s", parts[0])
	}
	interfaceID, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("interfaceId 无效: %s", parts[1])
	}
	return
}
