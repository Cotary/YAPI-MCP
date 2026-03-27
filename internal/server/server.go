package server

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/Cotary/YAPI-MCP/internal/config"
	"github.com/Cotary/YAPI-MCP/internal/middleware"
	"github.com/Cotary/YAPI-MCP/internal/prompts"
	"github.com/Cotary/YAPI-MCP/internal/resources"
	"github.com/Cotary/YAPI-MCP/internal/tools"
)

// NewMCPServer 创建并配置完整的 MCP Server 实例
//
// 这是整个 MCP 服务的装配中心，依次完成：
//  1. 创建 MCPServer 并启用所需能力（Tools/Resources/Prompts/Logging）
//  2. 配置 Hooks（请求生命周期钩子）和 Middleware（工具调用中间件）
//  3. 注册所有 Tools（9 个 YAPI 操作工具）
//  4. 注册所有 Resources（项目列表 + 接口文档模板）
//  5. 注册所有 Prompts（API 审查、文档生成、测试生成）
func NewMCPServer(cfg *config.Config) *server.MCPServer {
	hooks := middleware.NewHooks()

	s := server.NewMCPServer(
		"yapi-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
		server.WithHooks(hooks),
		server.WithToolHandlerMiddleware(middleware.ToolTimingMiddleware),
		server.WithInstructions(
			"这是一个 YAPI API 文档管理的 MCP 服务。"+
				"你可以通过 Tools 查询和管理 YAPI 接口，"+
				"通过 Resources 读取项目信息，"+
				"通过 Prompts 获取 API 审查和文档生成的引导。"+
				"建议先调用 yapi_list_projects 了解可用项目。",
		),
	)

	// 注册 Tools
	deps := tools.NewDeps(cfg)
	s.AddTool(tools.ListProjectsTool(), deps.ListProjectsHandler)
	s.AddTool(tools.GetProjectInfoTool(), deps.GetProjectInfoHandler)
	s.AddTool(tools.GetCategoriesTool(), deps.GetCategoriesHandler)
	s.AddTool(tools.ListInterfacesTool(), deps.ListInterfacesHandler)
	s.AddTool(tools.GetInterfaceDetailTool(), deps.GetInterfaceDetailHandler)
	s.AddTool(tools.SearchInterfacesTool(), deps.SearchInterfacesHandler)
	s.AddTool(tools.SaveInterfaceTool(), deps.SaveInterfaceHandler)
	s.AddTool(tools.AddCategoryTool(), deps.AddCategoryHandler)
	s.AddTool(tools.ImportDataTool(), deps.ImportDataHandler)

	// 注册 Resources
	resources.Register(s, cfg)

	// 注册 Prompts
	prompts.Register(s)

	return s
}
