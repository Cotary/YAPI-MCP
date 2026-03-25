package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/Cotary/YAPI-MCP/internal/yapi"
)

// SaveInterfaceTool 定义 yapi_save_interface 工具
func SaveInterfaceTool() mcp.Tool {
	return mcp.NewTool("yapi_save_interface",
		mcp.WithDescription("新增或更新 YAPI 接口。提供 id 参数则更新已有接口，不提供则新增"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("项目 ID")),
		mcp.WithString("cat_id", mcp.Required(), mcp.Description("分类 ID")),
		mcp.WithString("title", mcp.Required(), mcp.Description("接口标题")),
		mcp.WithString("path", mcp.Required(), mcp.Description("接口路径，如 /api/users")),
		mcp.WithString("method", mcp.Required(), mcp.Description("请求方法"),
			mcp.Enum("GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS")),
		mcp.WithString("id", mcp.Description("接口 ID，提供则更新该接口，不提供则新增")),
		mcp.WithString("status", mcp.Description("接口状态: done / undone")),
		mcp.WithString("desc", mcp.Description("接口描述")),
		mcp.WithString("markdown", mcp.Description("Markdown 格式文档")),
		mcp.WithString("req_body_type", mcp.Description("请求体类型: json / form / raw")),
		mcp.WithString("req_body_other", mcp.Description("请求体内容（JSON 字符串，用于 json/raw 类型）")),
		mcp.WithString("res_body_type", mcp.Description("响应体类型: json / raw")),
		mcp.WithString("res_body", mcp.Description("响应体内容（JSON 字符串）")),
		mcp.WithString("req_query", mcp.Description("查询参数 JSON 数组，如 [{\"name\":\"page\",\"desc\":\"页码\"}]")),
		mcp.WithString("req_headers", mcp.Description("请求头 JSON 数组")),
		mcp.WithString("req_body_form", mcp.Description("表单请求体 JSON 数组")),
		// 标记为非只读、破坏性操作
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

// SaveInterfaceHandler 处理 yapi_save_interface 调用
func (d *Deps) SaveInterfaceHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := ParseProjectID(req.GetString("project_id", ""))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	catIDStr := req.GetString("cat_id", "")
	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cat_id 必须是数字: %s", catIDStr)), nil
	}

	reqBodyType := req.GetString("req_body_type", "")
	reqBodyOther := req.GetString("req_body_other", "")
	resBodyType := req.GetString("res_body_type", "")
	resBody := req.GetString("res_body", "")

	params := yapi.SaveInterfaceParams{
		CatID:               catID,
		Title:               req.GetString("title", ""),
		Path:                req.GetString("path", ""),
		Method:              req.GetString("method", "GET"),
		Status:              req.GetString("status", "undone"),
		Description:         req.GetString("desc", ""),
		Markdown:            req.GetString("markdown", ""),
		ReqBodyType:         reqBodyType,
		ReqBodyOther:        reqBodyOther,
		ReqBodyIsJSONSchema: reqBodyType == "json" && reqBodyOther != "",
		ResBodyType:         resBodyType,
		ResBody:             resBody,
		ResBodyIsJSONSchema: resBodyType == "json" && resBody != "",
	}

	if idStr := req.GetString("id", ""); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("id 必须是数字: %s", idStr)), nil
		}
		params.ID = id
	}

	if q := req.GetString("req_query", ""); q != "" {
		var items []yapi.ReqQuery
		if err := json.Unmarshal([]byte(q), &items); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("req_query JSON 解析失败: %v", err)), nil
		}
		params.ReqQuery = items
	}
	if h := req.GetString("req_headers", ""); h != "" {
		var items []yapi.ReqHeader
		if err := json.Unmarshal([]byte(h), &items); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("req_headers JSON 解析失败: %v", err)), nil
		}
		params.ReqHeaders = items
	}
	if f := req.GetString("req_body_form", ""); f != "" {
		var items []yapi.ReqBodyFormItem
		if err := json.Unmarshal([]byte(f), &items); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("req_body_form JSON 解析失败: %v", err)), nil
		}
		params.ReqBodyForm = items
	}

	client, err := d.GetClient(projectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, err := client.SaveInterface(params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("保存接口失败: %v", err)), nil
	}

	// 清除相关缓存
	d.Cache.Delete(fmt.Sprintf("menu:%d", projectID))

	action := "新增"
	if params.ID > 0 {
		action = "更新"
		d.Cache.Delete(fmt.Sprintf("iface_detail:%d:%d", projectID, params.ID))
	}

	return mcp.NewToolResultText(fmt.Sprintf("接口%s成功，ID: %s", action, result.ID)), nil
}
