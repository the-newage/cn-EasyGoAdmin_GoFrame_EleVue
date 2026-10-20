package cache

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

// 有序集合（ZSET）门面：在线用户等需要按分数排序与区间查询的场景。
// Redis 模式走 ZADD/ZREM/ZRANGEBYSCORE；
// 内存降级用进程内 map，读取时由调用方自行排序（单机语义）。

var (
	memZMu  sync.RWMutex
	memZSet = map[string]map[string]float64{} // key → member → score
)

// ZAdd 添加有序集合成员
func ZAdd(ctx context.Context, key, member string, score float64) error {
	if !memFallback.Load() {
		_, err := g.Redis().Do(ctx, "ZADD", key, score, member)
		if err == nil {
			return nil
		}
		memFallback.Store(true)
		g.Log().Warningf(ctx, "Redis 写入失败(%v)，缓存降级为内存模式", err)
	}
	memZMu.Lock()
	defer memZMu.Unlock()
	if memZSet[key] == nil {
		memZSet[key] = map[string]float64{}
	}
	memZSet[key][member] = score
	return nil
}

// ZRem 移除有序集合成员
func ZRem(ctx context.Context, key, member string) error {
	if !memFallback.Load() {
		if _, err := g.Redis().Do(ctx, "ZREM", key, member); err == nil {
			return nil
		}
	}
	memZMu.Lock()
	defer memZMu.Unlock()
	delete(memZSet[key], member)
	return nil
}

// ZRangeByScore 按分数区间返回成员（闭区间；max 传 math.Inf(1) 表示无上限）
func ZRangeByScore(ctx context.Context, key string, min, max float64) ([]string, error) {
	if !memFallback.Load() {
		v, err := g.Redis().Do(ctx, "ZRANGEBYSCORE", key, min, max)
		if err == nil {
			return v.Strings(), nil
		}
	}
	memZMu.RLock()
	defer memZMu.RUnlock()
	m := memZSet[key]
	out := make([]string, 0, len(m))
	for member, score := range m {
		if score >= min && score <= max {
			out = append(out, member)
		}
	}
	return out, nil
}
