package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"minigo/internal/infrastructure/config"
)

// OSSService OSS服务
type OSSService struct {
	endpoint        string
	accessKeyID     string
	accessKeySecret string
	bucketName      string
	region          string
	expireSeconds   int
}

// NewOSSService 创建OSS服务实例
func NewOSSService() *OSSService {
	return &OSSService{
		endpoint:        config.GetOSSEndpoint(),
		accessKeyID:     config.GetOSSAccessKeyID(),
		accessKeySecret: config.GetOSSAccessKeySecret(),
		bucketName:      config.GetOSSBucketName(),
		region:          config.GetOSSRegion(),
		expireSeconds:   config.GetOSSTokenExpireSeconds(),
	}
}

// UploadToken 上传token结构
type UploadToken struct {
	AccessKeyID string `json:"accessKeyId"`
	Policy      string `json:"policy"`
	Signature   string `json:"signature"`
	Dir         string `json:"dir"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Callback    string `json:"callback,omitempty"`
	CallbackVar string `json:"callbackVar,omitempty"`
}

// PolicyCondition 策略条件
type PolicyCondition struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

// normalizeDir 标准化目录路径
func (s *OSSService) normalizeDir(dir string) string {
	if dir == "" {
		return fmt.Sprintf("uploads/%s/", time.Now().Format("2006/01/02"))
	}
	// 确保目录以/结尾
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	return dir
}

// GetUploadToken 获取上传token
func (s *OSSService) GetUploadToken(dir string, maxSize int64) (*UploadToken, error) {
	if s.accessKeyID == "" || s.accessKeySecret == "" || s.bucketName == "" {
		return nil, fmt.Errorf("OSS配置不完整")
	}

	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(s.expireSeconds) * time.Second)
	expireTimeStr := expireTime.UTC().Format("2006-01-02T15:04:05.000Z")

	// 标准化目录路径
	dir = s.normalizeDir(dir)

	// 如果没有指定最大文件大小，默认10MB
	if maxSize <= 0 {
		maxSize = 10 * 1024 * 1024 // 10MB
	}

	// 构建策略
	policy := PolicyCondition{
		Expiration: expireTimeStr,
		Conditions: []interface{}{
			map[string]string{"bucket": s.bucketName},
			[]interface{}{"starts-with", "$key", dir},
			[]interface{}{"content-length-range", 0, maxSize},
			[]interface{}{"in", "$content-type", []string{"image/jpg", "image/png", "image/jpeg"}},
		},
	}

	// 序列化策略
	policyBytes, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("序列化策略失败: %v", err)
	}

	// Base64编码策略
	policyBase64 := base64.StdEncoding.EncodeToString(policyBytes)

	// 生成签名
	signature := s.generateSignature(policyBase64, false)

	// 构建host
	host := s.getHost()

	return &UploadToken{
		AccessKeyID: s.accessKeyID,
		Policy:      policyBase64,
		Signature:   signature,
		Dir:         dir,
		Host:        host,
		Expire:      expireTime.Unix(),
	}, nil
}

// generateSignature 生成HMAC-SHA1签名
func (s *OSSService) generateSignature(data string, urlEncode bool) string {
	h := hmac.New(sha1.New, []byte(s.accessKeySecret))
	h.Write([]byte(data))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if urlEncode {
		// URL编码签名中的特殊字符（用于URL签名）
		return url.QueryEscape(signature)
	}

	// 返回纯base64编码（用于policy签名）
	return signature
}

// getHost 获取host地址
func (s *OSSService) getHost() string {
	if s.endpoint != "" {
		return fmt.Sprintf("https://%s.%s", s.bucketName, s.endpoint)
	}
	// 默认阿里云OSS格式
	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", s.bucketName, s.region)
}

// GetFileURL 获取文件访问URL（带签名）
func (s *OSSService) GetFileURL(objectKey string) string {
	// 如果为空，直接返回空字符串
	if objectKey == "" {
		return ""
	}

	// 如果是完整URL，直接返回
	if strings.HasPrefix(objectKey, "http") {
		return objectKey
	}

	// 如果是开发环境，返回随机图片URL
	if config.IsDevEnv() {
		return fmt.Sprintf("https://loremflickr.com/400/400?lock=%d", time.Now().UnixNano())
	}

	return s.GetSignedURL(objectKey, 3600) // 默认1小时过期
}

// GetSignedURL 获取带签名的文件访问URL
func (s *OSSService) GetSignedURL(objectKey string, expireSeconds int) string {
	if s.accessKeyID == "" || s.accessKeySecret == "" || s.bucketName == "" {
		// 如果配置不完整，返回简单URL
		host := s.getHost()
		return fmt.Sprintf("%s/%s", host, objectKey)
	}

	// 计算过期时间
	expires := time.Now().Unix() + int64(expireSeconds)

	// 构建要签名的字符串
	// 格式: GET\n\n\n{expires}\n/{bucket}/{object}
	stringToSign := fmt.Sprintf("GET\n\n\n%d\n/%s/%s", expires, s.bucketName, objectKey)

	// 生成签名
	signature := s.generateSignature(stringToSign, true)

	// 构建带签名的URL
	host := s.getHost()
	return fmt.Sprintf("%s/%s?OSSAccessKeyId=%s&Expires=%d&Signature=%s",
		host, objectKey, s.accessKeyID, expires, signature)
}
