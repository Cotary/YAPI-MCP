package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProjectConfig 单个 YAPI 项目的配置
type ProjectConfig struct {
	ProjectID   int    `yaml:"project_id"`
	Token       string `yaml:"token"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	BaseURL     string `yaml:"base_url"` // 可选，覆盖全局 base_url
}

// Config 应用全局配置
type Config struct {
	YapiBaseURL     string          `yaml:"yapi_base_url"`
	CacheTTLMinutes int             `yaml:"cache_ttl_minutes"`
	LogLevel        string          `yaml:"log_level"`
	SkipTLSVerify   bool            `yaml:"skip_tls_verify"`
	Projects        []ProjectConfig `yaml:"projects"`
}

// GetBaseURL 获取项目的 base URL，优先使用项目级别配置
func (c *Config) GetBaseURL(project *ProjectConfig) string {
	if project.BaseURL != "" {
		return strings.TrimRight(project.BaseURL, "/")
	}
	return strings.TrimRight(c.YapiBaseURL, "/")
}

// FindProject 按 project_id 查找项目配置
func (c *Config) FindProject(projectID int) *ProjectConfig {
	for i := range c.Projects {
		if c.Projects[i].ProjectID == projectID {
			return &c.Projects[i]
		}
	}
	return nil
}

// Load 加载配置文件，优先级：configPath 参数 > YAPI_MCP_CONFIG 环境变量 > 当前目录 yapi-mcp.config.yaml
func Load(configPath string) *Config {
	path := resolveConfigPath(configPath)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("无法读取配置文件 %s: %v", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
	}

	applyEnvOverrides(&cfg)
	if err := validate(&cfg); err != nil {
		log.Fatalf("配置校验失败: %v", err)
	}

	return &cfg
}

func resolveConfigPath(flagPath string) string {
	if flagPath != "" {
		return flagPath
	}
	if envPath := os.Getenv("YAPI_MCP_CONFIG"); envPath != "" {
		return envPath
	}
	candidates := []string{"yapi-mcp.config.yaml", "config.yaml"}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	log.Fatal("未找到配置文件，请通过 -config 参数指定或在当前目录放置 yapi-mcp.config.yaml")
	return ""
}

// applyEnvOverrides 环境变量覆盖（优先级最高）
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("YAPI_BASE_URL"); v != "" {
		cfg.YapiBaseURL = v
	}
	if v := os.Getenv("YAPI_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("YAPI_SKIP_TLS_VERIFY"); v == "true" || v == "1" {
		cfg.SkipTLSVerify = true
	}
}

func validate(cfg *Config) error {
	if cfg.YapiBaseURL == "" {
		return fmt.Errorf("yapi_base_url 不能为空")
	}
	if len(cfg.Projects) == 0 {
		return fmt.Errorf("至少需要配置一个项目")
	}
	for i, p := range cfg.Projects {
		if p.ProjectID == 0 {
			return fmt.Errorf("projects[%d].project_id 不能为 0", i)
		}
		if p.Token == "" {
			return fmt.Errorf("projects[%d].token 不能为空", i)
		}
	}
	if cfg.CacheTTLMinutes < 0 {
		cfg.CacheTTLMinutes = 0
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	return nil
}
