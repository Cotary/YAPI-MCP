package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/Cotary/YAPI-MCP/internal/yapi"
)

// GetInterfaceDetailTool 定义 yapi_get_interface_detail 工具
func GetInterfaceDetailTool() mcp.Tool {
	return mcp.NewTool("yapi_get_interface_detail",
		mcp.WithDescription("获取接口的完整详细信息，包括请求参数、请求头、请求体、响应体等"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("项目 ID")),
		mcp.WithString("interface_id", mcp.Required(), mcp.Description("接口 ID")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

// GetInterfaceDetailHandler 处理 yapi_get_interface_detail 调用
func (d *Deps) GetInterfaceDetailHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := ParseProjectID(req.GetString("project_id", ""))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	ifaceIDStr := req.GetString("interface_id", "")
	ifaceID, err := strconv.Atoi(ifaceIDStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("interface_id 必须是数字: %s", ifaceIDStr)), nil
	}

	cacheKey := fmt.Sprintf("iface_detail:%d:%d", projectID, ifaceID)
	if cached, ok := d.Cache.Get(cacheKey); ok {
		return mcp.NewToolResultText(cached.(string)), nil
	}

	client, err := d.GetClient(projectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	detail, err := client.GetInterfaceDetail(ifaceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取接口详情失败: %v", err)), nil
	}

	formatted := formatInterfaceDetail(detail)
	d.Cache.Set(cacheKey, formatted)
	return mcp.NewToolResultText(formatted), nil
}

// formatInterfaceDetail 将接口详情格式化为结构清晰的 JSON
func formatInterfaceDetail(d *yapi.InterfaceDetail) string {
	result := map[string]any{
		"基本信息": map[string]any{
			"ID":     d.ID,
			"标题":     d.Title,
			"路径":     d.Path,
			"方法":     d.Method,
			"状态":     d.Status,
			"描述":     d.Description,
			"标签":     d.Tag,
		},
	}

	if len(d.ReqParams) > 0 {
		result["路径参数"] = d.ReqParams
	}
	if len(d.ReqQuery) > 0 {
		result["查询参数"] = d.ReqQuery
	}
	if len(d.ReqHeaders) > 0 {
		result["请求头"] = d.ReqHeaders
	}

	reqBody := map[string]any{
		"类型": d.ReqBodyType,
	}
	if d.ReqBodyType == "form" && len(d.ReqBodyForm) > 0 {
		reqBody["表单参数"] = d.ReqBodyForm
	}
	if d.ReqBodyType == "json" && d.ReqBodyOther != "" {
		var jsonSchema any
		if err := json.Unmarshal([]byte(d.ReqBodyOther), &jsonSchema); err == nil {
			reqBody["JSON Schema"] = jsonSchema
		} else {
			reqBody["原始内容"] = d.ReqBodyOther
		}
	}
	if d.ReqBodyType == "raw" && d.ReqBodyOther != "" {
		reqBody["原始内容"] = d.ReqBodyOther
	}
	result["请求体"] = reqBody

	resBody := map[string]any{
		"类型": d.ResBodyType,
	}
	if d.ResBody != "" {
		var jsonBody any
		if err := json.Unmarshal([]byte(d.ResBody), &jsonBody); err == nil {
			resBody["内容"] = jsonBody
		} else {
			resBody["原始内容"] = d.ResBody
		}
	}
	result["响应体"] = resBody

	if d.Markdown != "" {
		result["Markdown文档"] = d.Markdown
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data)
}
