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
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"

	"easygoadmin/api/v1"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
)

// ============================== 服务器监控业务实现 ==============================

// sServerMonitor 服务器监控业务实现
// 汇集进程运行时、磁盘、数据库与缓存连通信息
type sServerMonitor struct{}

// procStart 进程启动时间（包加载时刻）
var procStart = time.Now()

func init() {
	service.RegisterServerMonitor(&sServerMonitor{})
}

// Info 查询服务器监控信息
func (s *sServerMonitor) Info(ctx context.Context, req *v1.ServerInfoReq) (res *v1.ServerInfoRes, err error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	hostname, _ := os.Hostname()

	res = &v1.ServerInfoRes{
		Server: v1.ServerOsInfo{
			Hostname:   hostname,
			GoVersion:  runtime.Version(),
			Arch:       runtime.GOARCH,
			NumCPU:     runtime.NumCPU(),
			Goroutines: runtime.NumGoroutine(),
			Uptime:     time.Since(procStart).Round(time.Second).String(),
		},
		Mem: v1.ServerMemInfo{
			AllocMB:      roundMB(m.Alloc),
			TotalAllocMB: roundMB(m.TotalAlloc),
			SysMB:        roundMB(m.Sys),
			HeapAllocMB:  roundMB(m.HeapAlloc),
			HeapSysMB:    roundMB(m.HeapSys),
			HeapObjects:  m.HeapObjects,
			GCCycles:     m.NumGC,
		},
		Disk: diskUsage(),
	}

	// 数据库连通与版本
	start := time.Now()
	_, err = g.DB().Ctx(ctx).Exec(ctx, "SELECT VERSION()")
	pingMs := time.Since(start).Milliseconds()
	version := ""
	if err == nil {
		if v, verr := g.DB().Ctx(ctx).GetValue(ctx, "SELECT VERSION()"); verr == nil {
			version = v.String()
		}
		res.Db = v1.ServerDbInfo{Version: version, PingMs: pingMs}
	} else {
		res.Db = v1.ServerDbInfo{PingMs: pingMs}
	}

	// 缓存连通
	if cache.UsingMemory() {
		res.Cache = v1.ServerCacheBrief{Mode: "memory", Connected: false}
	} else if _, cerr := g.Redis().Do(ctx, "PING"); cerr == nil {
		res.Cache = v1.ServerCacheBrief{Mode: "redis", Connected: true}
	} else {
		res.Cache = v1.ServerCacheBrief{Mode: "redis", Connected: false}
	}

	// 定时任务调度器状态
	// 挂载数统计 gcron 中 job: 前缀条目（排除上传清理等非任务条目），
	// 应挂载数取任务表 status=1（运行中）的记录数，挂载失败或数据被外部改动会出现差值
	mounted := 0
	for _, e := range gcron.Entries() {
		if strings.HasPrefix(e.Name, "job:") {
			mounted++
		}
	}
	expected, _ := dao.Job.Ctx(ctx).Where("status", 1).Count()
	res.Scheduler = v1.ServerSchedulerInfo{
		Enabled:  !daemonStart.IsZero(),
		Healthy:  !daemonStart.IsZero() && mounted >= int(expected),
		Mounted:  mounted,
		Expected: int(expected),
		StartTime: func() string {
			if daemonStart.IsZero() {
				return ""
			}
			return daemonStart.Format("2006-01-02 15:04:05")
		}(),
	}

	if err != nil {
		// 数据库不通不视为整体失败，仅版本留空；返回详情供页面展示
		err = nil
	}
	return res, nil
}

// roundMB 字节转 MB（保留两位小数）
func roundMB(v uint64) float64 {
	return math.Round(float64(v)/1048576*100) / 100
}
