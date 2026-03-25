package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

// ListInterfacesTool 定义 yapi_list_interfaces 工具
func ListInterfacesTool() mcp.Tool {
	return mcp.NewTool("yapi_list_interfaces",
		mcp.WithDescription("获取指定项目或分类下的接口列表，支持分页"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("项目 ID")),
		mcp.WithString("cat_id", mcp.Description("分类 ID，不填则返回项目下全部接口")),
		mcp.WithNumber("page", mcp.Description("页码，默认 1")),
		mcp.WithNumber("limit", mcp.Description("每页数量，默认 20")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

// ListInterfacesHandler 处理 yapi_list_interfaces 调用
func (d *Deps) ListInterfacesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := ParseProjectID(req.GetString("project_id", ""))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	page := req.GetInt("page", 1)
	limit := req.GetInt("limit", 20)

	client, err := d.GetClient(projectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	catIDStr := req.GetString("cat_id", "")
	if catIDStr != "" {
		catID, err := strconv.Atoi(catIDStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cat_id 必须是数字: %s", catIDStr)), nil
		}
		result, err := client.GetInterfaceListByCat(catID, page, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取接口列表失败: %v", err)), nil
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	}

	result, err := client.GetInterfaceList(page, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取接口列表失败: %v", err)), nil
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}
