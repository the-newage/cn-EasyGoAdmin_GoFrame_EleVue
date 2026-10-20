// +======================================================================
// | EasyGoAdmin 企业数字化开发平台 [ 赋能开发者，助力企业发展 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2019-2026 半城风雨团队 版权所有
// +----------------------------------------------------------------------
// | 开源协议: Apache-2.0
// | 重要声明: 未经书面许可，不得移除或修改本框架相关版权标识。
// +----------------------------------------------------------------------
// | 官方网站: https://www.easygoadmin.vip
// +----------------------------------------------------------------------
// | 研发团队: 小蚂蚁云团队 荣誉出品
// +======================================================================
// | 版权和免责声明:
// | 一、知识产权声明
// | 本团队对该软件框架产品（以下简称“本产品”）拥有完整知识产权，包括但不限于商标权、专利权、
// | 著作权、商业秘密等。前述权利均受中华人民共和国及相关国际法律法规保护。未经本团队书面授权，
// | 任何个人或单位不得基于本产品本身申请任何形式的知识产权。
// | 二、使用限制与禁止条款
// | 您理解并同意，禁止将本产品用于以下任何行为：
// | 1. 任何违反中华人民共和国法律法规的用途；
// | 2. 任何侵害他人合法权益（包括但不限于人身权、财产权、知识产权等）的行为；
// | 3. 任何危害国家安全、公共利益及社会道德风尚的行为。
// | 三、责任豁免
// | 因使用本产品而产生的任何损失或法律纠纷（包括但不限于直接损失、间接损失、附带损失、衍生损失
// | 以及因违约、侵权等引发的责任），本团队均不承担任何法律责任。若因违规使用本产品导致任何第三
// | 方提起索赔或行政机关进行调查，您应独立承担全部责任。
// | 四、违约与法律责任
// | 如发现任何违反本声明的行为，本团队有权：
// | 1. 无条件配合公安机关及其他有权机关的调查取证工作；
// | 2. 依法追究相关方的民事责任，并要求赔偿全部损失；
// | 3. 保留进一步采取刑事控告等法律措施的权利。
// | 五、许可范围说明
// | 1、被授权方应视本产品源代码为核心机密信息，负有最高标准的保密与保管义务。禁止以任何形式（包
// | 括但不限于直接或间接的分享、复制、传播、分发、转售、出租、开源、反向工程、反编译、反汇编或通
// | 过网络服务提供给任何第三方）泄露、披露或处置源代码的全部或部分内容。任何未经授权的复制、使用、
// | 传播行为均构成对本团队知识产权的严重侵权，本团队有权依法追究其全部法律责任，并要求赔偿因此造
// | 成的一切经济损失及商誉损害。
// | 2、本产品仅限于公司或个人在法律允许范围内，用于内部且合规的软件产品研发。具体权利义务以《软件
// | 许可协议》附件为准。
// +======================================================================

package logic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"easygoadmin/api/v1"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
)

// ============================== 缓存监控业务实现 ==============================

// sCacheMonitor 缓存监控业务实现
// Redis 模式解析 INFO 各节；内存降级模式仅返回模式标识
type sCacheMonitor struct{}

func init() {
	service.RegisterCacheMonitor(&sCacheMonitor{})
}

// Info 查询缓存监控信息
func (s *sCacheMonitor) Info(ctx context.Context, req *v1.CacheInfoReq) (res *v1.CacheInfoRes, err error) {
	if cache.UsingMemory() {
		return &v1.CacheInfoRes{Mode: "memory", Connected: false}, nil
	}
	res = &v1.CacheInfoRes{Mode: "redis", Connected: true}
	sections := map[string]map[string]string{}
	for _, name := range []string{"server", "memory", "stats", "clients", "keyspace"} {
		v, err := g.Redis().Do(ctx, "INFO", name)
		if err != nil {
			res.Connected = false
			return res, nil
		}
		sections[name] = parseRedisInfo(v.String())
	}

	res.Version = sections["server"]["redis_version"]
	if sec, e := strconv.ParseInt(sections["server"]["uptime_in_seconds"], 10, 64); e == nil {
		res.Uptime = (time.Duration(sec) * time.Second).String()
	}
	res.UsedMemory = sections["memory"]["used_memory_human"]
	res.UsedMemoryPeak = sections["memory"]["used_memory_peak_human"]
	res.MaxMemory = sections["memory"]["maxmemory_human"]
	res.FragRatio = sections["memory"]["mem_fragmentation_ratio"]
	res.ConnectedClients, _ = strconv.ParseInt(sections["clients"]["connected_clients"], 10, 64)
	res.ExpiredKeys, _ = strconv.ParseInt(sections["stats"]["expired_keys"], 10, 64)
	res.EvictedKeys, _ = strconv.ParseInt(sections["stats"]["evicted_keys"], 10, 64)
	res.TotalCommands, _ = strconv.ParseInt(sections["stats"]["total_commands_processed"], 10, 64)
	res.TotalConns, _ = strconv.ParseInt(sections["stats"]["total_connections_received"], 10, 64)
	res.Hits, _ = strconv.ParseInt(sections["stats"]["keyspace_hits"], 10, 64)
	res.Misses, _ = strconv.ParseInt(sections["stats"]["keyspace_misses"], 10, 64)
	if total := res.Hits + res.Misses; total > 0 {
		res.HitRate = strconv.FormatFloat(float64(res.Hits)/float64(total)*100, 'f', 2, 64)
	} else {
		res.HitRate = "0.00"
	}
	if v, e := g.Redis().Do(ctx, "DBSIZE"); e == nil {
		res.Keys = v.Int64()
	}
	for k, line := range sections["keyspace"] {
		if !strings.HasPrefix(k, "db") {
			continue
		}
		entry := v1.CacheKeySpaceDb{Db: k}
		for _, kv := range strings.Split(line, ",") {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				continue
			}
			switch parts[0] {
			case "keys":
				entry.Keys, _ = strconv.ParseInt(parts[1], 10, 64)
			case "expires":
				entry.Expires, _ = strconv.ParseInt(parts[1], 10, 64)
			}
		}
		res.KeySpace = append(res.KeySpace, entry)
	}
	return res, nil
}

// parseRedisInfo 解析 INFO 输出为 键→值 映射
func parseRedisInfo(s string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			out[line[:idx]] = strings.TrimSpace(line[idx+1:])
		}
	}
	return out
}
