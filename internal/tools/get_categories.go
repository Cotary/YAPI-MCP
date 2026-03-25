package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetCategoriesTool 定义 yapi_get_categories 工具
func GetCategoriesTool() mcp.Tool {
	return mcp.NewTool("yapi_get_categories",
		mcp.WithDescription("获取项目下的接口分类列表，每个分类包含其下的接口摘要（ID、标题、路径、方法）"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("项目 ID")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
}

// GetCategoriesHandler 处理 yapi_get_categories 调用
func (d *Deps) GetCategoriesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := ParseProjectID(req.GetString("project_id", ""))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cacheKey := fmt.Sprintf("menu:%d", projectID)
	if cached, ok := d.Cache.Get(cacheKey); ok {
		data, _ := json.MarshalIndent(cached, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	}

	client, err := d.GetClient(projectID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	menu, err := client.GetMenuWithInterfaces()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取分类列表失败: %v", err)), nil
	}

	d.Cache.Set(cacheKey, menu)

	type compactInterface struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Path   string `json:"path"`
		Method string `json:"method"`
	}
	type compactCategory struct {
		ID          int                `json:"id"`
		Name        string             `json:"name"`
		Description string             `json:"description,omitempty"`
		Interfaces  []compactInterface `json:"interfaces"`
	}

	result := make([]compactCategory, 0, len(menu))
	for _, cat := range menu {
		interfaces := make([]compactInterface, 0, len(cat.List))
		for _, iface := range cat.List {
			interfaces = append(interfaces, compactInterface{
				ID:     iface.ID,
				Title:  iface.Title,
				Path:   iface.Path,
				Method: iface.Method,
			})
		}
		result = append(result, compactCategory{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Interfaces:  interfaces,
		})
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}
