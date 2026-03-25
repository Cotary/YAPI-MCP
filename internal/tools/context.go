package tools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Cotary/YAPI-MCP/internal/cache"
	"github.com/Cotary/YAPI-MCP/internal/config"
	"github.com/Cotary/YAPI-MCP/internal/yapi"
)

// Deps 所有工具共享的依赖
type Deps struct {
	Config *config.Config
	Cache  *cache.Cache
}

// NewDeps 创建工具依赖
func NewDeps(cfg *config.Config) *Deps {
	ttl := time.Duration(cfg.CacheTTLMinutes) * time.Minute
	if ttl <= 0 {
		ttl = 1 * time.Second // 最低缓存兜底
	}
	return &Deps{
		Config: cfg,
		Cache:  cache.New(ttl),
	}
}

// GetClient 根据 project_id 创建 YAPI 客户端
func (d *Deps) GetClient(projectID int) (*yapi.Client, error) {
	p := d.Config.FindProject(projectID)
	if p == nil {
		return nil, fmt.Errorf("未找到 project_id=%d 的配置，请先通过 yapi_list_projects 查看可用项目", projectID)
	}
	baseURL := d.Config.GetBaseURL(p)
	return yapi.NewClient(baseURL, p.Token, p.ProjectID, d.Config.SkipTLSVerify), nil
}

// ParseProjectID 从字符串参数解析 project_id
func ParseProjectID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("project_id 必须是数字: %s", s)
	}
	return id, nil
}
