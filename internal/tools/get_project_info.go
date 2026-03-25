package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetProjectInfoTool 定义 yapi_get_project_info 工具
func GetProjectInfoTool() mcp.Tool {
	return mcp.NewTool("yapi_get_project_info",
		mcp.WithDescription("获取 YAPI 项目的详细信息，包括名称、描述、基础路径等"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("项目 ID")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

// GetProjectInfoHandler 处理 yapi_get_project_info 调用
func (d *Deps) GetProjectInfoHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := ParseProjectID(req.GetString("project_id", ""))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cacheKey := fmt.Sprintf("project_info:%d", projectID)
	if cached, ok := d.Cache.Get(cacheKey); ok {
		data, _ := json.MarshalIndent(cached, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	}

	client, err := d.GetClient(projectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	info, err := client.GetProject()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取项目信息失败: %v", err)), nil
	}

	d.Cache.Set(cacheKey, info)
	data, _ := json.MarshalIndent(info, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}
