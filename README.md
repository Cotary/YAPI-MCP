# YAPI-MCP

[English](#english) | [中文](#中文)

---

<a id="中文"></a>

## 中文

基于 Go 和 [mcp-go](https://github.com/mark3labs/mcp-go) 构建的 YAPI MCP Server，让 AI 助手（Cursor、Claude Desktop 等）能够直接查询和管理 YAPI 接口文档。

> 本项目同时附带 MCP 协议教学文档，详见 [MCP 教学（中文）](docs/mcp-tutorial.md)

### 功能

**Tools（7 个工具）**

| 工具 | 描述 | 类型 |
|------|------|------|
| `yapi_list_projects` | 列出所有已配置的 YAPI 项目 | 只读 |
| `yapi_get_project_info` | 获取项目详细信息 | 只读 |
| `yapi_get_categories` | 获取接口分类及分类下的接口 | 只读 |
| `yapi_list_interfaces` | 获取接口列表（分页、按分类筛选） | 只读 |
| `yapi_get_interface_detail` | 获取接口完整详情 | 只读 |
| `yapi_search_interfaces` | 按关键词搜索接口（支持跨项目） | 只读 |
| `yapi_save_interface` | 新增或更新接口 | 写入 |

**Resources（2 个资源）**

| URI | 描述 |
|-----|------|
| `yapi://projects` | 已配置项目列表（静态） |
| `yapi://projects/{projectId}/interfaces/{interfaceId}` | 接口文档（动态模板） |

**Prompts（3 个提示模板）**

| 名称 | 描述 |
|------|------|
| `api_design_review` | API 设计审查 |
| `generate_api_docs` | 生成 API 文档 |
| `generate_api_tests` | 生成接口测试用例 |

**其他特性**：双模式传输（stdio + Streamable HTTP）、多项目配置、内存缓存、Hooks + Middleware

### 快速开始

**1. 构建**

```bash
git clone https://github.com/Cotary/YAPI-MCP.git
cd YAPI-MCP
go build -o yapi-mcp ./cmd/yapi-mcp
```

**2. 配置**

```bash
cp config.example.yaml config.yaml
# 编辑 config.yaml，填入 YAPI 地址和项目 Token
```

```yaml
yapi_base_url: "https://yapi.your-company.com"
cache_ttl_minutes: 10
log_level: "info"

projects:
  - project_id: 1026
    token: "your-project-token"
    name: "用户服务"
    description: "用户相关的 API 接口"
```

**3. 在 Cursor 中使用**

创建 `.cursor/mcp.json`：

stdio 模式（个人使用）：
```json
{
  "mcpServers": {
    "yapi": {
      "command": "/path/to/yapi-mcp",
      "args": ["-transport", "stdio", "-config", "/path/to/config.yaml"]
    }
  }
}
```

HTTP 模式（团队共享）：
```bash
# 先在服务器启动
./yapi-mcp -transport http -config config.yaml -port 8080
```
```json
{
  "mcpServers": {
    "yapi": {
      "url": "http://your-server:8080/mcp"
    }
  }
}
```

### 配置说明

**配置加载优先级**：`-config` 参数 > `YAPI_MCP_CONFIG` 环境变量 > 当前目录 `config.yaml`

**环境变量覆盖**：

| 变量 | 说明 |
|------|------|
| `YAPI_BASE_URL` | 覆盖全局 base_url |
| `YAPI_LOG_LEVEL` | 覆盖日志级别 |
| `YAPI_MCP_CONFIG` | 指定配置文件路径 |

**命令行参数**：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-transport` | `stdio` | 传输模式：`stdio` 或 `http` |
| `-config` | - | 配置文件路径 |
| `-port` | `8080` | HTTP 模式监听端口 |

### 项目结构

```
yapi-mcp/
├── cmd/yapi-mcp/main.go             # 入口
├── internal/
│   ├── config/config.go             # 配置加载
│   ├── yapi/client.go, types.go     # YAPI HTTP 客户端
│   ├── cache/cache.go               # 内存 TTL 缓存
│   ├── server/server.go             # MCP Server 装配
│   ├── tools/*.go                   # 7 个 MCP 工具
│   ├── resources/resources.go       # MCP 资源
│   ├── prompts/prompts.go           # MCP 提示模板
│   └── middleware/middleware.go      # Hooks + Middleware
├── docs/
│   └── mcp-tutorial.md              # MCP 教学（中文）
├── config.example.yaml              # 配置示例
├── Makefile                         # 构建命令
└── go.mod / go.sum
```

---

<a id="english"></a>

## English

A YAPI MCP Server built with Go and [mcp-go](https://github.com/mark3labs/mcp-go), enabling AI assistants (Cursor, Claude Desktop, etc.) to directly query and manage YAPI API documentation.

> This project also includes an MCP protocol tutorial: [MCP 教学（中文）](docs/mcp-tutorial.md)

### Features

**Tools (7 tools)**

| Tool | Description | Type |
|------|-------------|------|
| `yapi_list_projects` | List all configured YAPI projects | Read-only |
| `yapi_get_project_info` | Get project details | Read-only |
| `yapi_get_categories` | Get categories with interfaces | Read-only |
| `yapi_list_interfaces` | List interfaces (pagination, filter by category) | Read-only |
| `yapi_get_interface_detail` | Get full interface details | Read-only |
| `yapi_search_interfaces` | Search interfaces by keyword (cross-project) | Read-only |
| `yapi_save_interface` | Create or update an interface | Write |

**Resources (2 resources)**

| URI | Description |
|-----|-------------|
| `yapi://projects` | Configured project list (static) |
| `yapi://projects/{projectId}/interfaces/{interfaceId}` | Interface document (dynamic template) |

**Prompts (3 prompt templates)**

| Name | Description |
|------|-------------|
| `api_design_review` | API design review |
| `generate_api_docs` | Generate API documentation |
| `generate_api_tests` | Generate interface test cases |

**Other features**: Dual transport (stdio + Streamable HTTP), multi-project config, in-memory caching, Hooks + Middleware

### Quick Start

**1. Build**

```bash
git clone https://github.com/Cotary/YAPI-MCP.git
cd YAPI-MCP
go build -o yapi-mcp ./cmd/yapi-mcp
```

**2. Configure**

```bash
cp config.example.yaml config.yaml
# Edit config.yaml with your YAPI URL and project tokens
```

```yaml
yapi_base_url: "https://yapi.your-company.com"
cache_ttl_minutes: 10
log_level: "info"

projects:
  - project_id: 1026
    token: "your-project-token"
    name: "User Service"
    description: "User-related APIs"
```

**3. Use with Cursor**

Create `.cursor/mcp.json`:

stdio mode (personal use):
```json
{
  "mcpServers": {
    "yapi": {
      "command": "/path/to/yapi-mcp",
      "args": ["-transport", "stdio", "-config", "/path/to/config.yaml"]
    }
  }
}
```

HTTP mode (team sharing):
```bash
# Start on your server first
./yapi-mcp -transport http -config config.yaml -port 8080
```
```json
{
  "mcpServers": {
    "yapi": {
      "url": "http://your-server:8080/mcp"
    }
  }
}
```

### Configuration

**Config loading priority**: `-config` flag > `YAPI_MCP_CONFIG` env > `config.yaml` in current directory

**Environment variable overrides**:

| Variable | Description |
|----------|-------------|
| `YAPI_BASE_URL` | Override global base_url |
| `YAPI_LOG_LEVEL` | Override log level |
| `YAPI_MCP_CONFIG` | Config file path |

**CLI flags**:

| Flag | Default | Description |
|------|---------|-------------|
| `-transport` | `stdio` | Transport mode: `stdio` or `http` |
| `-config` | - | Config file path |
| `-port` | `8080` | HTTP mode listen port |

### Project Structure

```
yapi-mcp/
├── cmd/yapi-mcp/main.go             # Entry point
├── internal/
│   ├── config/config.go             # Configuration loader
│   ├── yapi/client.go, types.go     # YAPI HTTP client
│   ├── cache/cache.go               # In-memory TTL cache
│   ├── server/server.go             # MCP Server assembly
│   ├── tools/*.go                   # 7 MCP tools
│   ├── resources/resources.go       # MCP resources
│   ├── prompts/prompts.go           # MCP prompt templates
│   └── middleware/middleware.go      # Hooks + Middleware
├── docs/
│   └── mcp-tutorial.md              # MCP Tutorial (Chinese)
├── config.example.yaml              # Config example
├── Makefile                         # Build commands
└── go.mod / go.sum
```

---

## License

MIT
