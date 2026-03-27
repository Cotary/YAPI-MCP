package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/Cotary/YAPI-MCP/internal/yapi"
)

// ImportDataTool 定义 yapi_import_data 工具
func ImportDataTool() mcp.Tool {
	return mcp.NewTool("yapi_import_data",
		mcp.WithDescription("向 YAPI 项目导入接口数据，支持 Swagger、JSON、HAR、Postman 格式。可通过 json 参数传入数据或通过 url 参数从远程拉取"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("项目 ID")),
		mcp.WithString("type", mcp.Required(), mcp.Description("导入数据格式"),
			mcp.Enum("swagger", "json", "har", "postman")),
		mcp.WithString("merge", mcp.Required(), mcp.Description("数据合并策略: normal(普通模式)、good(智能合并)、merge(完全覆盖)"),
			mcp.Enum("normal", "good", "merge")),
		mcp.WithString("json", mcp.Description("导入的 JSON 数据（序列化后的字符串），与 url 参数二选一")),
		mcp.WithString("url", mcp.Description("数据来源 URL，服务端将通过此 URL 获取导入数据，与 json 参数二选一")),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(true),
	)
}

// ImportDataHandler 处理 yapi_import_data 调用
func (d *Deps) ImportDataHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := ParseProjectID(req.GetString("project_id", ""))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	importType := req.GetString("type", "")
	merge := req.GetString("merge", "normal")
	jsonData := req.GetString("json", "")
	dataURL := req.GetString("url", "")

	if jsonData == "" && dataURL == "" {
		return mcp.NewToolResultError("json 和 url 参数至少提供一个"), nil
	}

	client, err := d.GetClient(projectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = client.ImportData(yapi.ImportDataParams{
		Type:  importType,
		JSON:  jsonData,
		Merge: merge,
		URL:   dataURL,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("数据导入失败: %v", err)), nil
	}

	d.Cache.Delete(fmt.Sprintf("categories:%d", projectID))
	d.Cache.Delete(fmt.Sprintf("menu:%d", projectID))

	mergeDesc := map[string]string{
		"normal": "普通模式",
		"good":   "智能合并",
		"merge":  "完全覆盖",
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"数据导入成功\n  项目: %d\n  格式: %s\n  合并策略: %s",
		projectID, importType, mergeDesc[merge],
	)), nil
}
