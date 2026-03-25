package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// ListProjectsTool 定义 yapi_list_projects 工具
func ListProjectsTool() mcp.Tool {
	return mcp.NewTool("yapi_list_projects",
		mcp.WithDescription("列出所有已配置的 YAPI 项目及其基本信息，这是使用其他工具的起点"),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

// ListProjectsHandler 处理 yapi_list_projects 调用
func (d *Deps) ListProjectsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	type projectItem struct {
		ProjectID   int    `json:"project_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	items := make([]projectItem, 0, len(d.Config.Projects))
	for _, p := range d.Config.Projects {
		items = append(items, projectItem{
			ProjectID:   p.ProjectID,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	data, _ := json.MarshalIndent(items, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("已配置 %d 个项目:\n%s", len(items), string(data))), nil
}
