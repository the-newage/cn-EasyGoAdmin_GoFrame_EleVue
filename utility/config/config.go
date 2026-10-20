// Package config 提供配置读取的环境变量覆盖接线
//
// 覆盖规则遵循 12-Factor 标准优先级：GF_ 前缀环境变量 > 配置文件。
// 键名规则：配置路径大写、点换下划线并加 GF_ 前缀
// （database.default.link → GF_DATABASE_DEFAULT_LINK）。
// 实现方式为包装默认文件适配器（gcfg.Adapter 是框架官方扩展点），
// 注册后全部 Get/MustGet 读取点零改动地获得覆盖能力。
package config

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/genv"
)

// adapterEnvFirst 环境变量优先的配置适配器
//
// 仅覆盖配置文件中已存在的键：环境变量名由配置路径正向推导（无歧义），
// 而反向拆分下划线无法还原 uploadDir 这类驼峰键名，故不做补键注入。
// 文件返回的 map 会被就地叠加环境变量值，进程内环境恒定，重复叠加幂等。
type adapterEnvFirst struct {
	base gcfg.Adapter // 默认文件适配器，环境未命中时回落
}

// RegisterEnvOverride 注册环境变量优先的配置适配器（包装并替换默认文件适配器）
//
// 须在首次读取配置之前调用（服务启动入口最前，含缓存/数据库/Redis 初始化）；
// 未设置任何 GF_ 环境变量时行为与默认文件适配器完全一致。
func RegisterEnvOverride() {
	gcfg.Instance().SetAdapter(&adapterEnvFirst{base: gcfg.Instance().GetAdapter()})
}

// Available 委托文件适配器
func (a *adapterEnvFirst) Available(ctx context.Context, resource ...string) bool {
	return a.base.Available(ctx, resource...)
}

// Get 先查 GF_ 前缀环境变量，未命中回落文件适配器；
// 文件返回 map（如 database 整节点）时对已有键逐叶叠加环境变量覆盖
func (a *adapterEnvFirst) Get(ctx context.Context, pattern string) (any, error) {
	if v := envValue(pattern); v != nil {
		return v, nil
	}
	value, err := a.base.Get(ctx, pattern)
	if err != nil || value == nil {
		return value, err
	}
	if m, ok := value.(map[string]any); ok {
		overlayMap(m, pattern)
	}
	return value, nil
}

// Data 返回全量配置并叠加环境变量覆盖
func (a *adapterEnvFirst) Data(ctx context.Context) (map[string]any, error) {
	data, err := a.base.Data(ctx)
	if err != nil || data == nil {
		return data, err
	}
	overlayMap(data, "")
	return data, nil
}

// envValue 按键路径取 GF_ 前缀环境变量；未设置或空值返回 nil（空串视为未设置，
// 避免容器编排注入的空占位把配置文件中的有效值清空）
func envValue(pattern string) any {
	if pattern == "" || pattern == "." {
		return nil
	}
	v := genv.Get("GF_" + strings.ToUpper(strings.ReplaceAll(pattern, ".", "_")))
	if v == nil || v.String() == "" {
		return nil
	}
	return v.Val()
}

// overlayMap 对 map 中已存在的键递归叠加环境变量覆盖
//
// prefix 为当前层级的配置路径（顶层为空）；仅处理标量覆盖，
// 切片/数组不支持环境变量替换（本项目可覆盖键均为标量）
func overlayMap(m map[string]any, prefix string) {
	for k, val := range m {
		path := k
		if prefix != "" && prefix != "." {
			path = prefix + "." + k
		}
		if env := envValue(path); env != nil {
			m[k] = env
			continue
		}
		if child, ok := val.(map[string]any); ok {
			overlayMap(child, path)
		}
	}
}
