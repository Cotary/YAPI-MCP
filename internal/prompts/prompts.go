package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register 注册所有 MCP Prompts 到服务器
//
// MCP Prompts 是预定义的交互模板，帮助 LLM 以一致的方式完成特定任务。
// 与 Tools 不同，Prompts 不执行操作，而是返回构造好的消息序列供 LLM 参考。
func Register(s *server.MCPServer) {
	registerAPIReviewPrompt(s)
	registerAPIDocGenPrompt(s)
	registerAPITestGenPrompt(s)
}

// registerAPIReviewPrompt API 设计审查提示模板
func registerAPIReviewPrompt(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("api_design_review",
		mcp.WithPromptDescription("审查 YAPI 中的 API 接口设计，检查是否符合 RESTful 规范、命名约定、参数设计等最佳实践"),
		mcp.WithArgument("project_id",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("YAPI 项目 ID"),
		),
		mcp.WithArgument("interface_id",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("要审查的接口 ID"),
		),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		projectID := req.Params.Arguments["project_id"]
		interfaceID := req.Params.Arguments["interface_id"]

		return mcp.NewGetPromptResult(
			"API 设计审查",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(fmt.Sprintf(
						`请审查 YAPI 项目 %s 中接口 %s 的设计，请按以下维度进行分析：

1. **RESTful 规范**：HTTP 方法是否正确使用、URI 设计是否合理
2. **命名规范**：路径和参数命名是否一致、清晰
3. **参数设计**：必填/选填是否合理、参数类型是否正确
4. **请求体**：结构是否清晰、字段是否完整
5. **响应体**：是否包含必要的状态码和错误信息格式
6. **安全性**：是否需要认证、是否有敏感数据暴露风险

请先调用 yapi_get_interface_detail 工具获取接口详情，然后给出具体的审查意见和改进建议。`,
						projectID, interfaceID,
					)),
				),
			},
		), nil
	})
}

// registerAPIDocGenPrompt API 文档生成提示模板
func registerAPIDocGenPrompt(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("generate_api_docs",
		mcp.WithPromptDescription("基于 YAPI 项目中的接口定义，生成开发者友好的 API 文档"),
		mcp.WithArgument("project_id",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("YAPI 项目 ID"),
		),
		mcp.WithArgument("format",
			mcp.ArgumentDescription("文档格式：markdown（默认）或 openapi"),
		),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		projectID := req.Params.Arguments["project_id"]
		format := req.Params.Arguments["format"]
		if format == "" {
			format = "markdown"
		}

		return mcp.NewGetPromptResult(
			"API 文档生成",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(fmt.Sprintf(
						`请为 YAPI 项目 %s 生成 %s 格式的 API 文档。

操作步骤：
1. 先调用 yapi_get_categories 获取项目的分类和接口列表
2. 对每个接口调用 yapi_get_interface_detail 获取详情
3. 按分类组织，生成包含以下内容的文档：
   - 接口概述（路径、方法、描述）
   - 请求参数说明（路径参数、查询参数、请求头）
   - 请求体示例
   - 响应体示例
   - 错误码说明（如有）`,
						projectID, format,
					)),
				),
			},
		), nil
	})
}

// registerAPITestGenPrompt API 测试用例生成提示模板
func registerAPITestGenPrompt(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("generate_api_tests",
		mcp.WithPromptDescription("基于 YAPI 接口定义生成自动化测试用例代码"),
		mcp.WithArgument("project_id",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("YAPI 项目 ID"),
		),
		mcp.WithArgument("interface_id",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("接口 ID"),
		),
		mcp.WithArgument("language",
			mcp.ArgumentDescription("测试代码语言：go（默认）、typescript、python"),
		),
	), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		projectID := req.Params.Arguments["project_id"]
		interfaceID := req.Params.Arguments["interface_id"]
		language := req.Params.Arguments["language"]
		if language == "" {
			language = "go"
		}

		return mcp.NewGetPromptResult(
			"API 测试用例生成",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(fmt.Sprintf(
						`请为 YAPI 项目 %s 中的接口 %s 生成 %s 语言的自动化测试用例。

步骤：
1. 调用 yapi_get_interface_detail 获取接口详情
2. 根据接口定义生成测试代码，包含：
   - 正常请求测试（Happy Path）
   - 参数缺失测试
   - 参数类型错误测试
   - 边界值测试
   - 权限/认证测试（如适用）
3. 使用该语言常用的测试框架`,
						projectID, interfaceID, language,
					)),
				),
			},
		), nil
	})
}
