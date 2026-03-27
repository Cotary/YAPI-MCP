package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/Cotary/YAPI-MCP/internal/yapi"
)

// AddCategoryTool 定义 yapi_add_category 工具
func AddCategoryTool() mcp.Tool {
	return mcp.NewTool("yapi_add_category",
		mcp.WithDescription("新增 YAPI 接口分类，在指定项目下创建一个新的接口分组"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("项目 ID")),
		mcp.WithString("name", mcp.Required(), mcp.Description("分类名称")),
		mcp.WithString("desc", mcp.Description("分类描述")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(true),
	)
}

// AddCategoryHandler 处理 yapi_add_category 调用
func (d *Deps) AddCategoryHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := ParseProjectID(req.GetString("project_id", ""))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	name := req.GetString("name", "")
	if name == "" {
		return mcp.NewToolResultError("分类名称不能为空"), nil
	}

	client, err := d.GetClient(projectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := client.AddCategory(yapi.AddCategoryParams{
		Name:        name,
		Description: req.GetString("desc", ""),
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("新增分类失败: %v", err)), nil
	}

	d.Cache.Delete(fmt.Sprintf("categories:%d", projectID))
	d.Cache.Delete(fmt.Sprintf("menu:%d", projectID))

	return mcp.NewToolResultText(fmt.Sprintf(
		"分类创建成功\n  ID: %d\n  名称: %s\n  项目: %d",
		result.ID, result.Name, result.ProjectID,
	)), nil
}
