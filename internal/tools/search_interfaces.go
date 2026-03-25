package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/Cotary/YAPI-MCP/internal/config"
	"github.com/Cotary/YAPI-MCP/internal/yapi"
)

// SearchInterfacesTool 定义 yapi_search_interfaces 工具
func SearchInterfacesTool() mcp.Tool {
	return mcp.NewTool("yapi_search_interfaces",
		mcp.WithDescription("按关键词搜索接口，支持按名称和路径模糊匹配，可跨项目搜索"),
		mcp.WithString("keyword", mcp.Required(), mcp.Description("搜索关键词，按名称和路径匹配")),
		mcp.WithString("project_id", mcp.Description("限定项目 ID，不填则搜索所有已配置项目")),
		mcp.WithNumber("limit", mcp.Description("返回数量限制，默认 20")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

// SearchInterfacesHandler 处理 yapi_search_interfaces 调用
func (d *Deps) SearchInterfacesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keyword := req.GetString("keyword", "")
	if keyword == "" {
		return mcp.NewToolResultError("keyword 不能为空"), nil
	}

	limit := req.GetInt("limit", 20)
	if limit <= 0 {
		limit = 20
	}

	projectIDStr := req.GetString("project_id", "")
	var projects []config.ProjectConfig
	if projectIDStr != "" {
		pid, err := ParseProjectID(projectIDStr)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		p := d.Config.FindProject(pid)
		if p == nil {
			return mcp.NewToolResultError(fmt.Sprintf("未找到 project_id=%d 的配置", pid)), nil
		}
		projects = []config.ProjectConfig{*p}
	} else {
		projects = d.Config.Projects
	}

	type searchResult struct {
		ProjectID   int    `json:"project_id"`
		ProjectName string `json:"project_name"`
		InterfaceID int    `json:"interface_id"`
		Title       string `json:"title"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		CatID       int    `json:"cat_id"`
		Status      string `json:"status"`
	}

	var results []searchResult
	kw := strings.ToLower(keyword)

	for _, p := range projects {
		menu, err := d.getMenuCached(p)
		if err != nil {
			continue
		}
		for _, cat := range menu {
			for _, iface := range cat.List {
				if matchKeyword(iface, kw) {
					results = append(results, searchResult{
						ProjectID:   p.ProjectID,
						ProjectName: p.Name,
						InterfaceID: iface.ID,
						Title:       iface.Title,
						Path:        iface.Path,
						Method:      iface.Method,
						CatID:       iface.CatID,
						Status:      iface.Status,
					})
					if len(results) >= limit {
						break
					}
				}
			}
			if len(results) >= limit {
				break
			}
		}
		if len(results) >= limit {
			break
		}
	}

	if len(results) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("未找到匹配 \"%s\" 的接口", keyword)), nil
	}

	data, _ := json.MarshalIndent(results, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("找到 %d 个匹配接口:\n%s", len(results), string(data))), nil
}

func (d *Deps) getMenuCached(p config.ProjectConfig) ([]yapi.MenuCategory, error) {
	cacheKey := fmt.Sprintf("menu:%d", p.ProjectID)
	if cached, ok := d.Cache.Get(cacheKey); ok {
		return cached.([]yapi.MenuCategory), nil
	}

	baseURL := d.Config.GetBaseURL(&p)
	client := yapi.NewClient(baseURL, p.Token, p.ProjectID, d.Config.SkipTLSVerify)
	menu, err := client.GetMenuWithInterfaces()
	if err != nil {
		return nil, err
	}
	d.Cache.Set(cacheKey, menu)
	return menu, nil
}

func matchKeyword(iface yapi.InterfaceSummary, kw string) bool {
	return strings.Contains(strings.ToLower(iface.Title), kw) ||
		strings.Contains(strings.ToLower(iface.Path), kw)
}
