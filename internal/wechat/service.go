package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"github.com/silenceper/wechat/v2"
	wechatcache "github.com/silenceper/wechat/v2/cache"
	wechatconfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"go.uber.org/zap"
)

// Service 微信服务
type Service struct {
	cfg *config.Config
	log *zap.Logger
	wc  *wechat.Wechat
}

// NewService 创建微信服务
func NewService(cfg *config.Config, log *zap.Logger) *Service {
	return &Service{
		cfg: cfg,
		log: log,
		wc:  wechat.NewWechat(),
	}
}

// getOfficialAccount 获取公众号实例
func (s *Service) getOfficialAccount() any {
	memory := wechatcache.NewMemory()
	wechatCfg := &wechatconfig.Config{
		AppID:     s.cfg.WechatAppID,
		AppSecret: s.cfg.WechatSecret,
		Cache:     memory,
	}
	return s.wc.GetOfficialAccount(wechatCfg)
}

// UploadMaterialResult 上传素材结果
type UploadMaterialResult struct {
	MediaID   string `json:"media_id"`
	WechatURL string `json:"wechat_url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// UploadMaterial 上传素材到微信
func (s *Service) UploadMaterial(filePath string) (*UploadMaterialResult, error) {
	startTime := time.Now()
	oa := s.getOfficialAccount()

	// 通过反射或接口获取 Material 实例
	// 这里使用类型断言获取实际的官方账号接口
	mat := getMaterialFromOA(oa)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		s.log.Error("open file failed",
			zap.String("path", filePath),
			zap.Error(err))
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	// 调用微信 API 上传
	mediaID, url, err := mat.AddMaterial(material.MediaTypeImage, file)
	if err != nil {
		s.log.Error("upload material failed",
			zap.String("path", filePath),
			zap.Error(err))
		return nil, fmt.Errorf("upload material: %w", err)
	}

	duration := time.Since(startTime)
	s.log.Info("material uploaded",
		zap.String("path", filePath),
		zap.String("media_id", maskMediaID(mediaID)),
		zap.Duration("duration", duration))

	return &UploadMaterialResult{
		MediaID:   mediaID,
		WechatURL: url,
	}, nil
}

// CreateDraftResult 创建草稿结果
type CreateDraftResult struct {
	MediaID  string `json:"media_id"`
	DraftURL string `json:"draft_url,omitempty"`
}

// CreateDraft 创建草稿
func (s *Service) CreateDraft(articles []*draft.Article) (*CreateDraftResult, error) {
	startTime := time.Now()
	oa := s.getOfficialAccount()

	// 通过反射或接口获取 Draft 实例
	dm := getDraftFromOA(oa)

	// 转换为 SDK 格式
	var draftArticles []draft.Article
	for _, a := range articles {
		draftArticles = append(draftArticles, *a)
	}

	mediaID, err := dm.AddDraft(draftArticles)
	if err != nil {
		s.log.Error("create draft failed", zap.Error(err))
		return nil, fmt.Errorf("create draft: %w", err)
	}

	duration := time.Since(startTime)
	s.log.Info("draft created",
		zap.String("media_id", maskMediaID(mediaID)),
		zap.Duration("duration", duration))

	// 构造草稿 URL
	draftURL := fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/appmsg?t=media/appmsg_edit_v2&action=edit&createType=0&token=")

	return &CreateDraftResult{
		MediaID:  mediaID,
		DraftURL: draftURL,
	}, nil
}

// UploadMaterialFromBytes 从字节数据上传素材
func (s *Service) UploadMaterialFromBytes(data []byte, filename string) (*UploadMaterialResult, error) {
	// 创建临时文件
	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, "md2wechat_"+filename)
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return nil, fmt.Errorf("write temp file: %w", err)
	}
	defer os.Remove(tmpPath)

	return s.UploadMaterial(tmpPath)
}

// AccessTokenResult 获取 access_token 结果（用于调试）
type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken 获取 access_token（调试用）
func (s *Service) GetAccessToken() (*AccessTokenResult, error) {
	oa := s.getOfficialAccount()
	at := getAccessTokenFromOA(oa)
	accessToken, err := at.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	return &AccessTokenResult{
		AccessToken: accessToken.AccessToken,
		ExpiresIn:   accessToken.ExpiresIn,
	}, nil
}

// maskMediaID 遮蔽 media_id 用于日志
func maskMediaID(id string) string {
	if len(id) < 8 {
		return "***"
	}
	return id[:4] + "***" + id[len(id)-4:]
}

// UploadMaterialWithRetry 带重试的上传
func (s *Service) UploadMaterialWithRetry(filePath string, maxRetries int) (*UploadMaterialResult, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		result, err := s.UploadMaterial(filePath)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(time.Second)
		}
	}
	return nil, lastErr
}

// DownloadFile 下载文件到临时目录
func DownloadFile(url string) (string, error) {
	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// 发起请求
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// 创建临时文件
	tmpDir := os.TempDir()
	ext := filepath.Ext(url)
	if ext == "" {
		ext = ".jpg"
	}
	tmpPath := filepath.Join(tmpDir, "md2wechat_download_"+ext)
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmpFile.Close()

	// 写入文件
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("write file: %w", err)
	}

	return tmpPath, nil
}

// CreateMultipartFormData 创建 multipart 表单数据
func CreateMultipartFormData(fieldName, filename string, data []byte) (string, *bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		writer.Close()
		return "", nil, ""
	}

	if _, err := part.Write(data); err != nil {
		writer.Close()
		return "", nil, ""
	}

	contentType := writer.FormDataContentType()
	writer.Close()

	return contentType, body, filename
}

// JSONMarshal 自定义 JSON 序列化
func JSONMarshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// 以下辅助函数用于从 OfficialAccount 接口获取具体实例
// 由于 silenceper/wechat SDK 的接口设计，需要通过接口方法获取

type materialInterface interface {
	AddMaterial(mediaType material.MediaType, file *os.File) (mediaID string, url string, err error)
}

type draftInterface interface {
	AddDraft(articles []draft.Article) (string, error)
}

type accessTokenInterface interface {
	GetAccessToken() (*AccessToken, error)
}

// AccessToken access_token 结构
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// getMaterialFromOA 从 OfficialAccount 获取 Material 实例
func getMaterialFromOA(oa any) materialInterface {
	// 这里使用类型断言，实际类型由 SDK 提供
	type materialGetter interface {
		GetMaterial() materialInterface
	}
	if getter, ok := oa.(materialGetter); ok {
		return getter.GetMaterial()
	}
	return nil
}

// getDraftFromOA 从 OfficialAccount 获取 Draft 实例
func getDraftFromOA(oa any) draftInterface {
	type draftGetter interface {
		GetDraft() draftInterface
	}
	if getter, ok := oa.(draftGetter); ok {
		return getter.GetDraft()
	}
	return nil
}

// getAccessTokenFromOA 从 OfficialAccount 获取 AccessToken 实例
func getAccessTokenFromOA(oa any) accessTokenInterface {
	type accessTokenGetter interface {
		GetAccessToken() accessTokenInterface
	}
	if getter, ok := oa.(accessTokenGetter); ok {
		return getter.GetAccessToken()
	}
	return nil
}
