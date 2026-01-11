package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"github.com/geekjourneyx/md2wechat-skill/internal/wechat"
	"go.uber.org/zap"
)

// Processor 图片处理器
type Processor struct {
	cfg        *config.Config
	log        *zap.Logger
	ws         *wechat.Service
	compressor *Compressor
}

// NewProcessor 创建图片处理器
func NewProcessor(cfg *config.Config, log *zap.Logger) *Processor {
	return &Processor{
		cfg:        cfg,
		log:        log,
		ws:         wechat.NewService(cfg, log),
		compressor: NewCompressor(log, cfg.MaxImageWidth, cfg.MaxImageSize),
	}
}

// UploadResult 上传结果
type UploadResult struct {
	MediaID   string `json:"media_id"`
	WechatURL string `json:"wechat_url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// UploadLocalImage 上传本地图片
func (p *Processor) UploadLocalImage(filePath string) (*UploadResult, error) {
	p.log.Info("uploading local image", zap.String("path", filePath))

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// 检查图片格式
	if !IsValidImageFormat(filePath) {
		return nil, fmt.Errorf("unsupported image format: %s", filePath)
	}

	// 如果需要压缩，先处理
	processedPath := filePath
	if p.cfg.CompressImages {
		compressedPath, compressed, err := p.compressor.CompressImage(filePath)
		if err != nil {
			p.log.Warn("compress failed, using original", zap.Error(err))
		} else if compressed {
			processedPath = compressedPath
			defer os.Remove(compressedPath)
			p.log.Info("using compressed image", zap.String("path", processedPath))
		}
	}

	// 上传到微信
	result, err := p.ws.UploadMaterialWithRetry(processedPath, 3)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		MediaID:   result.MediaID,
		WechatURL: result.WechatURL,
	}, nil
}

// DownloadAndUpload 下载在线图片并上传
func (p *Processor) DownloadAndUpload(url string) (*UploadResult, error) {
	p.log.Info("downloading and uploading image", zap.String("url", url))

	// 下载图片
	tmpPath, err := wechat.DownloadFile(url)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpPath)

	// 检查格式
	if !IsValidImageFormat(tmpPath) {
		return nil, fmt.Errorf("downloaded file is not a valid image")
	}

	// 压缩（如果需要）
	processedPath := tmpPath
	if p.cfg.CompressImages {
		compressedPath, compressed, err := p.compressor.CompressImage(tmpPath)
		if err != nil {
			p.log.Warn("compress failed, using original", zap.Error(err))
		} else if compressed {
			processedPath = compressedPath
			defer os.Remove(compressedPath)
			p.log.Info("using compressed image", zap.String("path", processedPath))
		}
	}

	// 上传到微信
	result, err := p.ws.UploadMaterialWithRetry(processedPath, 3)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		MediaID:   result.MediaID,
		WechatURL: result.WechatURL,
	}, nil
}

// GenerateAndUploadResult AI 生成图片结果
type GenerateAndUploadResult struct {
	Prompt      string `json:"prompt"`
	OriginalURL string `json:"original_url"`
	MediaID     string `json:"media_id"`
	WechatURL   string `json:"wechat_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

// GenerateAndUpload AI 生成图片并上传
func (p *Processor) GenerateAndUpload(prompt string) (*GenerateAndUploadResult, error) {
	p.log.Info("generating image via AI", zap.String("prompt", prompt))

	// 验证配置
	if err := p.cfg.ValidateForImageGeneration(); err != nil {
		return nil, err
	}

	// 调用图片生成 API
	imageURL, err := p.callImageAPI(prompt)
	if err != nil {
		return nil, fmt.Errorf("generate image: %w", err)
	}
	p.log.Info("image generated", zap.String("url", imageURL))

	// 下载生成的图片
	tmpPath, err := wechat.DownloadFile(imageURL)
	if err != nil {
		return nil, fmt.Errorf("download generated image: %w", err)
	}
	defer os.Remove(tmpPath)

	// 压缩（如果需要）
	processedPath := tmpPath
	if p.cfg.CompressImages {
		compressedPath, compressed, err := p.compressor.CompressImage(tmpPath)
		if err != nil {
			p.log.Warn("compress failed, using original", zap.Error(err))
		} else if compressed {
			processedPath = compressedPath
			defer os.Remove(compressedPath)
			p.log.Info("using compressed image", zap.String("path", processedPath))
		}
	}

	// 上传到微信
	result, err := p.ws.UploadMaterialWithRetry(processedPath, 3)
	if err != nil {
		return nil, err
	}

	return &GenerateAndUploadResult{
		Prompt:      prompt,
		OriginalURL: imageURL,
		MediaID:     result.MediaID,
		WechatURL:   result.WechatURL,
	}, nil
}

// callImageAPI 调用图片生成 API（兼容 OpenAI DALL-E）
func (p *Processor) callImageAPI(prompt string) (string, error) {
	// 构造请求
	reqBody := map[string]any{
		"model":  "dall-e-3",
		"prompt": prompt,
		"n":      1,
		"size":   "1024x1024",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 创建请求
	url := p.cfg.ImageAPIBase + "/images/generations"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.cfg.ImageAPIKey)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result struct {
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Data) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	return result.Data[0].URL, nil
}

// GetImageInfo 获取图片信息
func (p *Processor) GetImageInfo(filePath string) (*ImageInfo, error) {
	return GetImageInfo(filePath)
}

// CompressImage 压缩图片（公开方法）
func (p *Processor) CompressImage(filePath string) (string, bool, error) {
	return p.compressor.CompressImage(filePath)
}

// SetCompressQuality 设置压缩质量
func (p *Processor) SetCompressQuality(quality int) {
	p.compressor.SetQuality(quality)
}
