package middleware

import (
	"context"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewHooks 创建 MCP 请求生命周期钩子
//
// Hooks 可以拦截所有 MCP 请求的生命周期事件，适用于：
//   - 日志记录：追踪所有请求的来源和结果
//   - 监控指标：统计请求量、错误率
//   - 审计追踪：记录谁在什么时候做了什么
func NewHooks() *server.Hooks {
	hooks := &server.Hooks{}

	// BeforeAny：每个请求进入时触发
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, msg any) {
		log.Printf("[MCP] --> %s (id=%v)", method, id)
	})

	// OnSuccess：请求成功完成时触发
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, msg any, result any) {
		log.Printf("[MCP] <-- %s OK", method)
	})

	// OnError：请求处理出错时触发
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, msg any, err error) {
		log.Printf("[MCP] <-- %s ERROR: %v", method, err)
	})

	// BeforeInitialize：客户端连接初始化时触发，可用于记录客户端信息
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, msg *mcp.InitializeRequest) {
		log.Printf("[MCP] 客户端初始化: %s %s",
			msg.Params.ClientInfo.Name,
			msg.Params.ClientInfo.Version,
		)
	})

	// BeforeCallTool：工具调用前触发，可用于鉴权、限流
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, msg *mcp.CallToolRequest) {
		log.Printf("[MCP] 调用工具: %s", msg.Params.Name)
	})

	return hooks
}

// ToolTimingMiddleware 工具调用耗时统计中间件
//
// Middleware 是一种链式处理模式：每个中间件包裹下一个处理函数，
// 可以在调用前后添加逻辑。这种模式常见于 HTTP 框架（如 Go 的 net/http、Express.js）。
func ToolTimingMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()
		result, err := next(ctx, req)
		elapsed := time.Since(start)

		if err != nil {
			log.Printf("[Tool] %s 失败 (耗时 %v): %v", req.Params.Name, elapsed, err)
		} else {
			log.Printf("[Tool] %s 完成 (耗时 %v)", req.Params.Name, elapsed)
		}

		return result, err
	}
}
