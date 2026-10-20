package cache

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

// 统一缓存门面：Redis 优先（g.Redis，配置读取 manifest/config 的 redis 节），
// Redis 不可用时自动降级为进程内存（gcache），保证单机开发可用。
var memFallback atomic.Bool

// Init 探测 Redis 连通性，失败则切换内存模式
func Init(ctx context.Context) {
	if _, err := g.Redis().Do(ctx, "PING"); err != nil {
		memFallback.Store(true)
		g.Log().Warningf(ctx, "Redis 不可用(%v)，缓存降级为进程内存模式", err)
		return
	}
	g.Log().Info(ctx, "Redis 缓存已连接")
}

// UsingMemory 是否处于内存降级模式
func UsingMemory() bool { return memFallback.Load() }

// Set 写入缓存（ttlSeconds<=0 表示不过期）
func Set(ctx context.Context, key, value string, ttlSeconds int) error {
	if !memFallback.Load() {
		var err error
		if ttlSeconds > 0 {
			_, err = g.Redis().Do(ctx, "SETEX", key, ttlSeconds, value)
		} else {
			_, err = g.Redis().Do(ctx, "SET", key, value)
		}
		if err == nil {
			return nil
		}
		memFallback.Store(true)
		g.Log().Warningf(ctx, "Redis 写入失败(%v)，缓存降级为内存模式", err)
	}
	if ttlSeconds > 0 {
		return gcache.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second)
	}
	return gcache.Set(ctx, key, value, -1)
}

// Get 读取缓存；第二个返回值表示是否存在
func Get(ctx context.Context, key string) (string, bool) {
	if !memFallback.Load() {
		v, err := g.Redis().Do(ctx, "GET", key)
		if err == nil {
			if v.IsNil() {
				return "", false
			}
			return v.String(), true
		}
		memFallback.Store(true)
		g.Log().Warningf(ctx, "Redis 读取失败(%v)，缓存降级为内存模式", err)
	}
	v, err := gcache.Get(ctx, key)
	if err != nil || v.IsNil() {
		return "", false
	}
	return v.String(), true
}

// Del 删除缓存
func Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if !memFallback.Load() {
		args := make([]any, 0, len(keys))
		for _, k := range keys {
			args = append(args, k)
		}
		if _, err := g.Redis().Do(ctx, "DEL", args...); err == nil {
			return nil
		}
	}
	for _, k := range keys {
		_, _ = gcache.Remove(ctx, k)
	}
	return nil
}

// Exists 判断键是否存在
func Exists(ctx context.Context, key string) bool {
	if !memFallback.Load() {
		v, err := g.Redis().Do(ctx, "EXISTS", key)
		if err == nil {
			return v.Int() > 0
		}
	}
	ok, err := gcache.Contains(ctx, key)
	if err != nil {
		return false
	}
	return ok
}

// SetNX 不存在时写入（防重复提交），成功返回 true
func SetNX(ctx context.Context, key, value string, ttlSeconds int) bool {
	if !memFallback.Load() {
		v, err := g.Redis().Do(ctx, "SET", key, value, "NX", "EX", ttlSeconds)
		if err == nil {
			return !v.IsNil() && v.String() == "OK"
		}
	}
	if Exists(ctx, key) {
		return false
	}
	_ = Set(ctx, key, value, ttlSeconds)
	return true
}
