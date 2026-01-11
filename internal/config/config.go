package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	// 微信公众号配置
	WechatAppID  string `json:"wechat_appid" yaml:"wechat_appid" env:"WECHAT_APPID"`
	WechatSecret string `json:"wechat_secret" yaml:"wechat_secret" env:"WECHAT_SECRET"`

	// md2wechat.cn API 配置
	MD2WechatAPIKey    string `json:"md2wechat_api_key" yaml:"md2wechat_api_key" env:"MD2WECHAT_API_KEY"`
	DefaultConvertMode string `json:"default_convert_mode" yaml:"default_convert_mode" env:"CONVERT_MODE"`
	DefaultTheme       string `json:"default_theme" yaml:"default_theme" env:"DEFAULT_THEME"`

	// 图片生成 API 配置
	ImageAPIKey  string `json:"image_api_key" yaml:"image_api_key" env:"IMAGE_API_KEY"`
	ImageAPIBase string `json:"image_api_base" yaml:"image_api_base" env:"IMAGE_API_BASE"`

	// 图片处理配置
	CompressImages bool  `json:"compress_images" yaml:"compress_images" env:"COMPRESS_IMAGES"`
	MaxImageWidth  int   `json:"max_image_width" yaml:"max_image_width" env:"MAX_IMAGE_WIDTH"`
	MaxImageSize   int64 `json:"max_image_size" yaml:"max_image_size" env:"MAX_IMAGE_SIZE"`

	// 超时配置
	HTTPTimeout int `json:"http_timeout" yaml:"http_timeout" env:"HTTP_TIMEOUT"`

	// 配置文件路径（用于追踪）
	configFile string
}

// ConfigFile 配置文件结构（YAML/JSON）
type configFile struct {
	Wechat struct {
		AppID  string `json:"appid" yaml:"appid"`
		Secret string `json:"secret" yaml:"secret"`
	} `json:"wechat" yaml:"wechat"`

	API struct {
		MD2WechatKey string `json:"md2wechat_key" yaml:"md2wechat_key"`
		ImageKey     string `json:"image_key" yaml:"image_key"`
		ImageBaseURL string `json:"image_base_url" yaml:"image_base_url"`
		ConvertMode  string `json:"convert_mode" yaml:"convert_mode"`
		DefaultTheme string `json:"default_theme" yaml:"default_theme"`
		HTTPTimeout  int    `json:"http_timeout" yaml:"http_timeout"`
	} `json:"api" yaml:"api"`

	Image struct {
		Compress bool `json:"compress" yaml:"compress"`
		MaxWidth int  `json:"max_width" yaml:"max_width"`
		MaxSize  int  `json:"max_size_mb" yaml:"max_size_mb"`
	} `json:"image" yaml:"image"`
}

// Load 从配置文件和环境变量加载配置
// 优先级：环境变量 > 配置文件 > 默认值
func Load() (*Config, error) {
	return LoadWithDefaults("")
}

// LoadWithDefaults 使用指定配置文件路径加载配置
func LoadWithDefaults(configPath string) (*Config, error) {
	cfg := &Config{
		DefaultConvertMode: "api",
		DefaultTheme:       "default",
		CompressImages:     true,
		MaxImageWidth:      1920,
		MaxImageSize:       5 * 1024 * 1024, // 5MB
		HTTPTimeout:        30,
		ImageAPIBase:       "https://api.openai.com/v1",
	}

	// 1. 尝试从配置文件加载
	if configPath == "" {
		configPath = findConfigFile()
	}
	if configPath != "" {
		if err := loadFromFile(cfg, configPath); err != nil {
			// 配置文件加载失败不是致命错误，继续使用环境变量和默认值
			fmt.Fprintf(os.Stderr, "Warning: failed to load config file: %v\n", err)
		} else {
			cfg.configFile = configPath
		}
	}

	// 2. 环境变量覆盖配置文件
	loadFromEnv(cfg)

	// 3. 验证必需配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 4. 处理 MaxImageSize (配置文件中是 MB)
	if cfg.configFile != "" && cfg.MaxImageSize < 1024*1024 {
		// 如果值小于 1MB，可能是配置文件使用了 MB 单位
		cfg.MaxImageSize = cfg.MaxImageSize * 1024 * 1024
	}

	return cfg, nil
}

// findConfigFile 查找配置文件
func findConfigFile() string {
	// 按优先级查找
	paths := []string{
		"md2wechat.yaml",
		"md2wechat.yml",
		"md2wechat.json",
		".md2wechat.yaml",
		".md2wechat.yml",
		".md2wechat.json",
		filepath.Join(os.Getenv("HOME"), ".md2wechat.yaml"),
		filepath.Join(os.Getenv("HOME"), ".config", "md2wechat", "config.yaml"),
	}

	for _, path := range paths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}

	return ""
}

// loadFromFile 从文件加载配置
func loadFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".json" {
		return loadFromJSON(cfg, data)
	}
	// 默认使用 YAML
	return loadFromYAML(cfg, data)
}

// loadFromYAML 从 YAML 加载
func loadFromYAML(cfg *Config, data []byte) error {
	var cf configFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	// 映射到 Config
	if cf.Wechat.AppID != "" {
		cfg.WechatAppID = cf.Wechat.AppID
	}
	if cf.Wechat.Secret != "" {
		cfg.WechatSecret = cf.Wechat.Secret
	}
	if cf.API.MD2WechatKey != "" {
		cfg.MD2WechatAPIKey = cf.API.MD2WechatKey
	}
	if cf.API.ImageKey != "" {
		cfg.ImageAPIKey = cf.API.ImageKey
	}
	if cf.API.ImageBaseURL != "" {
		cfg.ImageAPIBase = cf.API.ImageBaseURL
	}
	if cf.API.ConvertMode != "" {
		cfg.DefaultConvertMode = cf.API.ConvertMode
	}
	if cf.API.DefaultTheme != "" {
		cfg.DefaultTheme = cf.API.DefaultTheme
	}
	if cf.API.HTTPTimeout > 0 {
		cfg.HTTPTimeout = cf.API.HTTPTimeout
	}
	cfg.CompressImages = cf.Image.Compress
	if cf.Image.MaxWidth > 0 {
		cfg.MaxImageWidth = cf.Image.MaxWidth
	}
	if cf.Image.MaxSize > 0 {
		cfg.MaxImageSize = int64(cf.Image.MaxSize) * 1024 * 1024
	}

	return nil
}

// loadFromJSON 从 JSON 加载
func loadFromJSON(cfg *Config, data []byte) error {
	var cf configFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}

	// 映射到 Config（与 loadFromYAML 相同的逻辑）
	if cf.Wechat.AppID != "" {
		cfg.WechatAppID = cf.Wechat.AppID
	}
	if cf.Wechat.Secret != "" {
		cfg.WechatSecret = cf.Wechat.Secret
	}
	if cf.API.MD2WechatKey != "" {
		cfg.MD2WechatAPIKey = cf.API.MD2WechatKey
	}
	if cf.API.ImageKey != "" {
		cfg.ImageAPIKey = cf.API.ImageKey
	}
	if cf.API.ImageBaseURL != "" {
		cfg.ImageAPIBase = cf.API.ImageBaseURL
	}
	if cf.API.ConvertMode != "" {
		cfg.DefaultConvertMode = cf.API.ConvertMode
	}
	if cf.API.DefaultTheme != "" {
		cfg.DefaultTheme = cf.API.DefaultTheme
	}
	if cf.API.HTTPTimeout > 0 {
		cfg.HTTPTimeout = cf.API.HTTPTimeout
	}
	cfg.CompressImages = cf.Image.Compress
	if cf.Image.MaxWidth > 0 {
		cfg.MaxImageWidth = cf.Image.MaxWidth
	}
	if cf.Image.MaxSize > 0 {
		cfg.MaxImageSize = int64(cf.Image.MaxSize) * 1024 * 1024
	}

	return nil
}

// loadFromEnv 从环境变量加载
func loadFromEnv(cfg *Config) {
	if v := os.Getenv("WECHAT_APPID"); v != "" {
		cfg.WechatAppID = v
	}
	if v := os.Getenv("WECHAT_SECRET"); v != "" {
		cfg.WechatSecret = v
	}
	if v := os.Getenv("MD2WECHAT_API_KEY"); v != "" {
		cfg.MD2WechatAPIKey = v
	}
	if v := os.Getenv("CONVERT_MODE"); v != "" {
		cfg.DefaultConvertMode = v
	}
	if v := os.Getenv("DEFAULT_THEME"); v != "" {
		cfg.DefaultTheme = v
	}
	if v := os.Getenv("IMAGE_API_KEY"); v != "" {
		cfg.ImageAPIKey = v
	}
	if v := os.Getenv("IMAGE_API_BASE"); v != "" {
		cfg.ImageAPIBase = v
	}
	if v := os.Getenv("COMPRESS_IMAGES"); v != "" {
		cfg.CompressImages = getEnvBool("COMPRESS_IMAGES", true)
	}
	if v := os.Getenv("MAX_IMAGE_WIDTH"); v != "" {
		cfg.MaxImageWidth = getEnvInt("MAX_IMAGE_WIDTH", cfg.MaxImageWidth)
	}
	if v := os.Getenv("MAX_IMAGE_SIZE"); v != "" {
		cfg.MaxImageSize = int64(getEnvInt("MAX_IMAGE_SIZE", int(cfg.MaxImageSize)))
	}
	if v := os.Getenv("HTTP_TIMEOUT"); v != "" {
		cfg.HTTPTimeout = getEnvInt("HTTP_TIMEOUT", cfg.HTTPTimeout)
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.WechatAppID == "" {
		return &ConfigError{Field: "WechatAppID", Message: "WECHAT_APPID is required"}
	}
	if c.WechatSecret == "" {
		return &ConfigError{Field: "WechatSecret", Message: "WECHAT_SECRET is required"}
	}

	// 验证转换模式
	if c.DefaultConvertMode != "api" && c.DefaultConvertMode != "ai" {
		return &ConfigError{Field: "ConvertMode", Message: "must be 'api' or 'ai'"}
	}

	// 验证数值范围
	if c.MaxImageWidth < 100 || c.MaxImageWidth > 10000 {
		return &ConfigError{Field: "MaxImageWidth", Message: "must be between 100 and 10000"}
	}
	if c.MaxImageSize < 1024*100 { // 最小 100KB
		return &ConfigError{Field: "MaxImageSize", Message: "must be at least 100KB"}
	}
	if c.HTTPTimeout < 1 || c.HTTPTimeout > 300 {
		return &ConfigError{Field: "HTTPTimeout", Message: "must be between 1 and 300 seconds"}
	}

	return nil
}

// ValidateForImageGeneration 验证图片生成配置
func (c *Config) ValidateForImageGeneration() error {
	if c.ImageAPIKey == "" {
		return &ConfigError{Field: "ImageAPIKey", Message: "IMAGE_API_KEY is required for image generation"}
	}
	return nil
}

// ValidateForAPIConversion 验证 API 转换配置
func (c *Config) ValidateForAPIConversion() error {
	if c.MD2WechatAPIKey == "" && c.DefaultConvertMode == "api" {
		return &ConfigError{Field: "MD2WechatAPIKey", Message: "MD2WECHAT_API_KEY is required for API mode"}
	}
	return nil
}

// GetConfigFile 获取配置文件路径
func (c *Config) GetConfigFile() string {
	return c.configFile
}

// ToMap 转换为 map 用于显示
func (c *Config) ToMap(maskSecret bool) map[string]any {
	result := map[string]any{
		"wechat_appid":         c.WechatAppID,
		"wechat_secret":        maskIf(c.WechatSecret, maskSecret),
		"default_convert_mode": c.DefaultConvertMode,
		"default_theme":        c.DefaultTheme,
		"md2wechat_api_key":    maskIf(c.MD2WechatAPIKey, maskSecret),
		"image_api_key":        maskIf(c.ImageAPIKey, maskSecret),
		"image_api_base":       c.ImageAPIBase,
		"compress_images":      c.CompressImages,
		"max_image_width":      c.MaxImageWidth,
		"max_image_size_mb":    c.MaxImageSize / 1024 / 1024,
		"http_timeout":         c.HTTPTimeout,
		"config_file":          c.configFile,
	}
	return result
}

// SaveConfig 保存配置到文件
func SaveConfig(path string, cfg *Config) error {
	ext := strings.ToLower(filepath.Ext(path))

	cf := configFile{}
	cf.Wechat.AppID = cfg.WechatAppID
	cf.Wechat.Secret = cfg.WechatSecret
	cf.API.MD2WechatKey = cfg.MD2WechatAPIKey
	cf.API.ImageKey = cfg.ImageAPIKey
	cf.API.ImageBaseURL = cfg.ImageAPIBase
	cf.API.ConvertMode = cfg.DefaultConvertMode
	cf.API.DefaultTheme = cfg.DefaultTheme
	cf.API.HTTPTimeout = cfg.HTTPTimeout
	cf.Image.Compress = cfg.CompressImages
	cf.Image.MaxWidth = cfg.MaxImageWidth
	cf.Image.MaxSize = int(cfg.MaxImageSize / 1024 / 1024)

	var data []byte
	var err error

	if ext == ".json" {
		data, err = json.MarshalIndent(cf, "", "  ")
	} else {
		data, err = yaml.Marshal(cf)
	}

	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

// ConfigError 配置错误
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error [%s]: %s", e.Field, e.Message)
}

// getEnvBool 获取布尔型环境变量
func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val == "true" || val == "1"
}

// getEnvInt 获取整型环境变量
func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

// getEnvString 获取字符串型环境变量
func getEnvString(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// maskIf 掩码处理
func maskIf(value string, mask bool) string {
	if !mask || value == "" {
		return value
	}
	if len(value) <= 4 {
		return "***"
	}
	return value[:2] + "***" + value[len(value)-2:]
}
