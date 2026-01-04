package scraper

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/pkg/logger"
	"go.uber.org/zap"
)

// CachedFetcher 带缓存和并发的抓取器
type CachedFetcher struct {
	fetcher    *Fetcher
	cache      *Cache
	semaphore  chan struct{} // 并发控制
}

// Cache 缓存结构
type Cache struct {
	mu    sync.RWMutex
	items map[string]*cacheEntry
}

// cacheEntry 缓存条目
type cacheEntry struct {
	page      *WebPage
	expiresAt time.Time
}

// NewCache 创建缓存
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]*cacheEntry),
	}
}

// Get 获取缓存
func (c *Cache) Get(key string) (*WebPage, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.page, true
}

// Set 设置缓存
func (c *Cache) Set(key string, page *WebPage, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheEntry{
		page:      page,
		expiresAt: time.Now().Add(ttl),
	}
}

// Clear 清空缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheEntry)
}

// Size 返回缓存大小
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// NewCachedFetcher 创建带缓存的抓取器
func NewCachedFetcher(cfg *config.ScraperConfig, maxConcurrency int, cacheTTL time.Duration) *CachedFetcher {
	return &CachedFetcher{
		fetcher:   NewFetcher(cfg),
		cache:     NewCache(),
		semaphore: make(chan struct{}, maxConcurrency),
	}
}

// Fetch 抓取单个网页 (带缓存)
func (f *CachedFetcher) Fetch(url string) (*WebPage, error) {
	// 尝试从缓存获取
	if page, found := f.cache.Get(url); found {
		logger.Get().Debug("缓存命中", zap.String("url", url))
		return page, nil
	}

	// 缓存未命中,执行抓取
	logger.Get().Debug("缓存未命中,开始抓取", zap.String("url", url))
	page, err := f.fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}

	// 存入缓存 (默认 TTL 1小时)
	f.cache.Set(url, page, 1*time.Hour)
	return page, nil
}

// FetchBatch 批量并发抓取网页
func (f *CachedFetcher) FetchBatch(ctx context.Context, urls []string) map[string]*FetchResult {
	log := logger.Get()
	log.Info("开始批量抓取", zap.Int("total_urls", len(urls)))

	results := make(map[string]*FetchResult)
	var wg sync.WaitGroup
	var successCount, failCount int32

	// 为每个 URL 启动一个 goroutine
	for _, url := range urls {
		wg.Add(1)

		go func(urlStr string) {
			defer wg.Done()

			// 获取信号量(限制并发数)
			select {
			case f.semaphore <- struct{}{}:
				defer func() { <-f.semaphore }()
			case <-ctx.Done():
				results[urlStr] = &FetchResult{
					URL:  urlStr,
					Err:  fmt.Errorf("操作已取消"),
					Page: nil,
				}
				return
			}

			// 抓取网页
			page, err := f.Fetch(urlStr)

			// 记录结果
			if err != nil {
				atomic.AddInt32(&failCount, 1)
				log.Warn("抓取失败", zap.String("url", urlStr), zap.Error(err))
			} else {
				atomic.AddInt32(&successCount, 1)
				log.Debug("抓取成功", zap.String("url", urlStr), zap.Int("content_length", len(page.Content)))
			}

			results[urlStr] = &FetchResult{
				URL:  urlStr,
				Err:  err,
				Page: page,
			}
		}(url)
	}

	// 等待所有任务完成
	wg.Wait()

	log.Info("批量抓取完成",
		zap.Int("total", len(urls)),
		zap.Int32("success", successCount),
		zap.Int32("failed", failCount),
		zap.Int("cache_size", f.cache.Size()),
	)

	return results
}

// FetchResult 抓取结果
type FetchResult struct {
	URL  string
	Page *WebPage
	Err  error
}

// ClearCache 清空缓存
func (f *CachedFetcher) ClearCache() {
	f.cache.Clear()
	logger.Get().Info("缓存已清空")
}

// GetCacheSize 获取缓存大小
func (f *CachedFetcher) GetCacheSize() int {
	return f.cache.Size()
}
