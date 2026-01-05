package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	resp "minigo/internal/interfaces/response"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
	Reset(key string)
}

// TokenBucket 令牌桶限流器
type TokenBucket struct {
	capacity int           // 桶容量
	tokens   int           // 当前令牌数
	rate     time.Duration // 令牌生成速率
	lastTime time.Time     // 上次更新时间
	mutex    sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity int, rate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

// Allow 检查是否允许请求
func (tb *TokenBucket) Allow(key string) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	// 计算应该添加的令牌数
	elapsed := now.Sub(tb.lastTime)
	tokensToAdd := int(elapsed / tb.rate)

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastTime = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// Reset 重置令牌桶
func (tb *TokenBucket) Reset(key string) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.tokens = tb.capacity
	tb.lastTime = time.Now()
}

// RateLimiterManager 限流器管理器
type RateLimiterManager struct {
	limiters map[string]*TokenBucket
	mutex    sync.RWMutex
	capacity int
	rate     time.Duration
}

// NewRateLimiterManager 创建限流器管理器
func NewRateLimiterManager(capacity int, rate time.Duration) *RateLimiterManager {
	return &RateLimiterManager{
		limiters: make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}

// GetLimiter 获取或创建限流器
func (rlm *RateLimiterManager) GetLimiter(key string) *TokenBucket {
	rlm.mutex.RLock()
	limiter, exists := rlm.limiters[key]
	rlm.mutex.RUnlock()

	if exists {
		return limiter
	}

	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()

	// 双重检查
	if limiter, exists = rlm.limiters[key]; exists {
		return limiter
	}

	limiter = NewTokenBucket(rlm.capacity, rlm.rate)
	rlm.limiters[key] = limiter
	return limiter
}

// Allow 检查是否允许请求
func (rlm *RateLimiterManager) Allow(key string) bool {
	limiter := rlm.GetLimiter(key)
	return limiter.Allow(key)
}

// Reset 重置限流器
func (rlm *RateLimiterManager) Reset(key string) {
	rlm.mutex.RLock()
	limiter, exists := rlm.limiters[key]
	rlm.mutex.RUnlock()

	if exists {
		limiter.Reset(key)
	}
}

// CleanupExpired 清理过期的限流器
func (rlm *RateLimiterManager) CleanupExpired() {
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()

	now := time.Now()
	for key, limiter := range rlm.limiters {
		// 如果限流器超过1小时没有使用，则删除
		if now.Sub(limiter.lastTime) > time.Hour {
			delete(rlm.limiters, key)
		}
	}
}

// 全局限流器管理器
var (
	globalRateLimiter *RateLimiterManager
	once              sync.Once
)

// GetGlobalRateLimiter 获取全局限流器
func GetGlobalRateLimiter() *RateLimiterManager {
	once.Do(func() {
		// 默认配置：每秒100个请求，桶容量200
		globalRateLimiter = NewRateLimiterManager(200, 10*time.Millisecond)

		// 启动清理协程
		go func() {
			ticker := time.NewTicker(10 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				globalRateLimiter.CleanupExpired()
			}
		}()
	})
	return globalRateLimiter
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	limiter := GetGlobalRateLimiter()

	return func(c *gin.Context) {
		key := keyFunc(c)
		if !limiter.Allow(key) {
			resp.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPBasedRateLimitMiddleware 基于IP的限流中间件
func IPBasedRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(func(c *gin.Context) string {
		return c.ClientIP()
	})
}

// UserBasedRateLimitMiddleware 基于用户的限流中间件
func UserBasedRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(func(c *gin.Context) string {
		userID, exists := c.Get(ContextUserIDKey)
		if !exists {
			return c.ClientIP() // 如果没有用户信息，使用IP
		}
		return fmt.Sprintf("user_%d", userID.(int64))
	})
}

// APIKeyBasedRateLimitMiddleware 基于API Key的限流中间件
func APIKeyBasedRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(func(c *gin.Context) string {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			return c.ClientIP()
		}
		return fmt.Sprintf("api_%s", apiKey)
	})
}

// CustomRateLimitMiddleware 自定义限流中间件
func CustomRateLimitMiddleware(capacity int, rate time.Duration, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	limiter := NewRateLimiterManager(capacity, rate)

	return func(c *gin.Context) {
		key := keyFunc(c)
		if !limiter.Allow(key) {
			resp.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
