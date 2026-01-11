package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// APIResponse md2wechat.cn API 响应
type APIResponse struct {
	Code int    `json:"code"` // 0 表示成功
	Msg  string `json:"msg"`  // 错误信息
	Data struct {
		HTML string `json:"html"` // 转换后的 HTML
	} `json:"data"`
}

// APIRequest md2wechat.cn API 请求
type APIRequest struct {
	Markdown string `json:"markdown"`
	Theme    string `json:"theme"`
	FontSize string `json:"fontSize,omitempty"`
}

// apiConverter API 模式转换器
type apiConverter struct {
	log     *zap.Logger
	baseURL string
	timeout time.Duration
}

// NewAPIConverter 创建 API 转换器
func NewAPIConverter(log *zap.Logger) *apiConverter {
	return &apiConverter{
		log:     log,
		baseURL: "https://www.md2wechat.cn/api/convert",
		timeout: 30 * time.Second,
	}
}

// convertViaAPI 通过 API 执行转换
func (c *converter) convertViaAPI(req *ConvertRequest) *ConvertResult {
	result := &ConvertResult{
		Mode:    ModeAPI,
		Theme:   req.Theme,
		Success: false,
	}

	// 获取 API 主题名
	apiTheme, err := c.theme.GetAPITheme(req.Theme)
	if err != nil {
		// 如果不是预定义的 API 主题，直接使用传入的主题名
		apiTheme = req.Theme
	}

	// 创建 API 转换器
	apiConv := NewAPIConverter(c.log)

	// 调用 API
	html, err := apiConv.Convert(&APIRequest{
		Markdown: req.Markdown,
		Theme:    apiTheme,
		FontSize: req.FontSize,
	}, req.APIKey)

	if err != nil {
		result.Error = fmt.Sprintf("API call failed: %s", err.Error())
		c.log.Error("API conversion failed",
			zap.String("theme", req.Theme),
			zap.Error(err))
		return result
	}

	// 提取图片引用
	images := c.ExtractImages(req.Markdown)

	result.HTML = html
	result.Images = images
	result.Success = true

	c.log.Info("API conversion succeeded",
		zap.String("theme", req.Theme),
		zap.Int("image_count", len(images)))

	return result
}

// Convert 调用 md2wechat.cn API 进行转换
func (a *apiConverter) Convert(req *APIRequest, apiKey string) (string, error) {
	// 序列化请求
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", a.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", apiKey)

	// 创建客户端
	client := &http.Client{
		Timeout: a.timeout,
	}

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	// 解析响应
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("parse response: %w (body: %s)", err, string(body))
	}

	// 检查响应状态
	if apiResp.Code != 0 {
		return "", &ConvertError{
			Code:    "API_ERROR",
			Message: fmt.Sprintf("API returned error code %d: %s", apiResp.Code, apiResp.Msg),
		}
	}

	// 返回 HTML
	return apiResp.Data.HTML, nil
}

// SetBaseURL 设置 API 基础 URL（用于测试）
func (a *apiConverter) SetBaseURL(url string) {
	a.baseURL = url
}

// SetTimeout 设置请求超时
func (a *apiConverter) SetTimeout(timeout time.Duration) {
	a.timeout = timeout
}
