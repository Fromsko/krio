# 性能优化文档

本文档介绍了 Krio 项目的性能优化功能,包括并发处理和缓存机制。

## 目录

- [缓存机制](#缓存机制)
- [并发处理](#并发处理)
- [配置选项](#配置选项)
- [使用示例](#使用示例)
- [性能监控](#性能监控)

## 缓存机制

### 功能特性

- **自动缓存**: 所有成功抓取的网页内容都会被自动缓存
- **TTL 过期**: 缓存条目会在指定时间后自动过期 (默认 1 小时)
- **线程安全**: 缓存使用读写锁,支持并发访问
- **内存管理**: 可手动清空缓存,释放内存

### 缓存流程

```
请求 URL → 检查缓存 → 命中? → 是: 返回缓存数据
                          ↓
                         否: 抓取网页 → 存入缓存 → 返回数据
```

### 缓存优势

1. **减少网络请求**: 重复访问相同 URL 时直接从缓存读取
2. **提升响应速度**: 缓存命中时无需网络请求,响应速度提升 100 倍以上
3. **降低 API 调用**: 减少 AI 总结的重复计算
4. **节省带宽**: 避免重复下载相同网页内容

## 并发处理

### 功能特性

- **并发控制**: 可配置最大并发数,避免资源耗尽
- **批量处理**: 支持一次性处理多个 URL
- **上下文支持**: 支持取消操作
- **错误隔离**: 单个 URL 失败不影响其他 URL 的处理

### 并发流程

```
批量 URL 请求
    ↓
创建 goroutine 池 (限制并发数)
    ↓
并发抓取网页 (使用缓存)
    ↓
AI 总结处理 (可并行)
    ↓
保存到 Obsidian (串行,避免冲突)
    ↓
返回所有结果
```

### 并发优势

1. **提升吞吐量**: 5 个并发 URL 的处理时间接近单个 URL 的时间
2. **资源利用**: 充分利用网络带宽和 CPU 资源
3. **用户体验**: 批量操作时大幅减少等待时间

## 配置选项

在 [config.yaml](../config.yaml) 中配置性能选项:

```yaml
scraper:
  # ... 基础配置 ...

  # 是否启用缓存 (推荐: true)
  enable_cache: true

  # 缓存过期时间 (推荐: 1h)
  # 可用单位: s (秒), m (分), h (小时)
  cache_ttl: 1h

  # 最大并发数 (推荐: 5)
  # 根据网络带宽和 CPU 性能调整
  # 过高可能导致资源耗尽或被封禁
  max_concurrency: 5
```

### 参数建议

| 环境 | 缓存 TTL | 并发数 | 说明 |
|------|----------|--------|------|
| 开发 | 30m | 3 | 频繁更新内容,较低并发避免影响调试 |
| 生产 | 1h | 5 | 平衡性能和资源使用 |
| 高性能 | 2h | 10 | 需要良好网络和服务器配置 |

## 使用示例

### 单个 URL (自动使用缓存)

```go
import "github.com/fromsko/krio/internal/tool"

// 创建工具 (自动启用缓存)
webNoteTool, _ := tool.NewSaveWebNoteTool(ctx, cfg)

// 第一次访问: 网络抓取 + 存入缓存
resp1, _ := webNoteTool.SaveWebNote(ctx, tool.SaveWebNoteRequest{
    URL: "https://example.com",
})

// 第二次访问: 直接从缓存返回 (速度快 100 倍)
resp2, _ := webNoteTool.SaveWebNote(ctx, tool.SaveWebNoteRequest{
    URL: "https://example.com",
})
```

### 批量 URL (并发处理)

```go
urls := []string{
    "https://example.com",
    "https://www.iana.org/domains/reserved",
    "https://www.wikipedia.org/",
    // ... 更多 URL
}

// 并发处理所有 URL
responses := webNoteTool.SaveWebNoteBatch(ctx, urls, tags, "Inbox")

// 处理结果
for i, resp := range responses {
    if resp.Success {
        fmt.Printf("✅ %s: %s\n", urls[i], resp.FilePath)
    } else {
        fmt.Printf("❌ %s: %s\n", urls[i], resp.Message)
    }
}
```

### 运行批量演示

```bash
# 编译批量演示程序
go build -o batch_demo.exe ./cmd/batch_demo

# 运行演示
./batch_demo.exe
```

## 性能监控

### 获取缓存统计

```go
// 获取缓存状态
stats := webNoteTool.GetCacheStats()

// stats 示例:
// {
//   "enabled": true,
//   "cache_size": 42,        // 当前缓存条目数
//   "cache_ttl": "1h0m0s",   // 缓存过期时间
//   "max_concurrency": 5     // 最大并发数
// }
```

### 清空缓存

```go
// 清空所有缓存
webNoteTool.ClearCache()
```

### 性能指标

#### 单个 URL 性能

| 场景 | 首次访问 | 缓存命中 | 提升 |
|------|----------|----------|------|
| 网络抓取 | 3-5 秒 | N/A | - |
| 缓存读取 | 3-5 秒 | <50ms | **100x** |

#### 批量 URL 性能 (5 个 URL, 并发数 5)

| 模式 | 总耗时 | 说明 |
|------|--------|------|
| 串行 | ~15-25 秒 | 5 个 URL 依次处理 |
| 并发 | ~3-5 秒 | 5 个 URL 同时处理 |
| 提升 | **5x** | 接近单个 URL 的时间 |

#### 批量 URL 性能 (10 个 URL, 并发数 5)

| 模式 | 总耗时 | 说明 |
|------|--------|------|
| 串行 | ~30-50 秒 | 10 个 URL 依次处理 |
| 并发 | ~6-10 秒 | 分 2 批并发处理 |
| 提升 | **5x** | 受限于并发数配置 |

## 最佳实践

### 1. 合理设置缓存 TTL

```yaml
# ✅ 推荐: 根据内容更新频率设置
cache_ttl: 1h    # 适用于内容更新较慢的网站
cache_ttl: 10m   # 适用于新闻、社交媒体等实时内容
cache_ttl: 24h   # 适用于文档、静态页面
```

```yaml
# ❌ 不推荐: 过短或过长
cache_ttl: 1m    # 太短,缓存效果不明显
cache_ttl: 720h  # 太长,可能返回过期内容
```

### 2. 控制并发数

```yaml
# ✅ 推荐: 根据实际情况调整
max_concurrency: 3   # 开发环境,避免影响调试
max_concurrency: 5   # 生产环境,平衡性能
max_concurrency: 10  # 高性能环境,需要良好网络
```

```yaml
# ❌ 不推荐: 过高或过低
max_concurrency: 1   # 太低,失去并发优势
max_concurrency: 100 # 太高,可能被封禁或资源耗尽
```

### 3. 使用批量 API

```go
// ✅ 推荐: 使用批量 API
responses := webNoteTool.SaveWebNoteBatch(ctx, urls, tags, folder)

// ❌ 不推荐: 循环调用单次 API
for _, url := range urls {
    webNoteTool.SaveWebNote(ctx, tool.SaveWebNoteRequest{URL: url})
}
```

### 4. 监控缓存使用

```go
// 定期检查缓存大小
stats := webNoteTool.GetCacheStats()
if stats["cache_size"].(int) > 1000 {
    // 缓存过大,考虑清空
    webNoteTool.ClearCache()
}
```

### 5. 错误处理

```go
// 批量处理时检查每个结果
responses := webNoteTool.SaveWebNoteBatch(ctx, urls, tags, folder)

successCount := 0
for _, resp := range responses {
    if resp.Success {
        successCount++
    }
}

log.Info("批量处理完成",
    zap.Int("total", len(responses)),
    zap.Int("success", successCount),
    zap.Int("failed", len(responses) - successCount),
)
```

## 技术实现

### 缓存结构

```go
type Cache struct {
    mu    sync.RWMutex           // 读写锁
    items map[string]*cacheEntry // 缓存条目
}

type cacheEntry struct {
    page      *WebPage    // 网页数据
    expiresAt time.Time   // 过期时间
}
```

### 并发控制

使用 buffered channel 作为信号量:

```go
type CachedFetcher struct {
    semaphore chan struct{} // 并发控制
}

// 获取信号量
f.semaphore <- struct{}{}
defer func() { <-f.semaphore }()
```

### goroutine 池

```go
var wg sync.WaitGroup
for _, url := range urls {
    wg.Add(1)
    go func(urlStr string) {
        defer wg.Done()
        // 处理逻辑...
    }(url)
}
wg.Wait()
```

## 故障排查

### 问题: 缓存未生效

**检查配置:**
```yaml
scraper:
  enable_cache: true  # 确保启用
```

**检查日志:**
```
DEBUG 缓存命中 url="https://example.com"  # 缓存命中
DEBUG 缓存未命中,开始抓取 url="https://example.com"  # 缓存未命中
```

### 问题: 并发性能不明显

**可能原因:**
1. 并发数设置过低
2. 网络带宽不足
3. 目标网站限速

**解决方案:**
```yaml
max_concurrency: 10  # 增加并发数
```

### 问题: 内存占用过高

**解决方案:**
```go
// 定期清空缓存
webNoteTool.ClearCache()

// 或减少 TTL
cache_ttl: 30m
```

## 未来优化

- [ ] 支持 Redis 缓存 (分布式)
- [ ] 支持 LRU 缓存淘汰策略
- [ ] 支持持久化缓存 (重启后保留)
- [ ] 添加缓存命中率统计
- [ ] 支持动态调整并发数
- [ ] 添加请求队列管理
