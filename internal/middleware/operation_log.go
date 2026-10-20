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

package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"easygoadmin/internal/consts"
	"easygoadmin/internal/router/meta"
	"easygoadmin/utility/response"
	"easygoadmin/utility/ua"
)

// ============================== 操作日志中间件 ==============================
//
// 请求结束后根据路由元数据（gmeta log/logType 标签）组装记录，
// 经有界通道异步批量落库，不阻塞业务请求。
//
// 开关：operationLog.enable 为总开关；operationLog.logQuery 控制查询类是否记录。

const (
	logChanSize   = 2048        // 日志通道容量
	logBatchSize  = 100         // 单批落库条数
	logFlushEvery = time.Second // 定时刷新间隔
	logMaxLen     = 8000        // param/result 截断长度
)

var (
	logOnce sync.Once                                    // 保证日志工作协程仅启动一次
	logChan = make(chan *OperationLogEntry, logChanSize) // 日志缓冲通道
	logDrop int64                                        // 累计丢弃条数
	logMu   sync.Mutex                                   // 丢弃计数互斥锁
)

// OperationLog 操作日志采集
//
// 依据路由元数据判断是否记录，组装日志记录后投入异步通道
func OperationLog(r *ghttp.Request) {
	if !cfgBool(r.Context(), "operationLog.enable", true) {
		r.Middleware.Next()
		return
	}
	start := time.Now()
	r.Middleware.Next()

	rm, ok := r.Context().Value(meta.CtxKey).(meta.RouteMeta)
	if !ok || rm.Log == "" {
		return
	}
	logType := consts.LogTypeMap[rm.LogType]
	if logType == consts.LogTypeQuery && !cfgBool(r.Context(), "operationLog.logQuery", true) {
		return
	}

	code, _ := response.GetCode(r)
	entry := &OperationLogEntry{
		Title:         rm.Log,
		Type:          logType,
		Source:        0,
		Method:        rm.Method + " " + rm.Path,
		RequestMethod: r.Method,
		URL:           r.URL.String(),
		Param:         trunc(collectParam(r), logMaxLen),
		Result:        trunc(collectResult(r), logMaxLen),
		IP:            realClientIP(r),
		Location:      requestLocation(r.Context(), realClientIP(r)),
		OS:            ua.ParseOS(r.Header.Get("User-Agent")),
		Browser:       ua.ParseBrowser(r.Header.Get("User-Agent")),
		ConsumeTime:   time.Since(start).Milliseconds(),
		Status:        0,
		Error:         "",
		CreateUser:    consts.GetUsername(r.Context()),
	}
	if code != 0 {
		entry.Status = 1
		if body, ok := response.GetBody(r); ok && body.Msg != "" {
			entry.Error = trunc(body.Msg, 500)
		}
	}

	logOnce.Do(startLogWorker)
	select {
	case logChan <- entry:
	default:
		logMu.Lock()
		logDrop++
		n := logDrop
		logMu.Unlock()
		if n%100 == 1 {
			g.Log().Warningf(r.Context(), "操作日志通道已满，累计丢弃 %d 条", n)
		}
	}
}

// collectParam 组装脱敏后的请求参数
//
// GET 与 DELETE 请求取 query 参数，写操作优先取请求体
func collectParam(r *ghttp.Request) string {
	if r.Method == "GET" || r.Method == "DELETE" {
		m := r.GetMap()
		if len(m) == 0 {
			return ""
		}
		return trunc(gjsonEncode(m), logMaxLen)
	}
	body := r.GetBodyString()
	if body != "" {
		return maskSecret(body)
	}
	return trunc(gjsonEncode(r.GetMap()), logMaxLen)
}

// collectResult 组装响应信封内容
func collectResult(r *ghttp.Request) string {
	body, ok := response.GetBody(r)
	if !ok {
		return ""
	}
	return trunc(gjsonEncode(body), logMaxLen)
}

// maskSecret 脱敏敏感字段
//
// 覆盖密码、盐与令牌等字段
func maskSecret(s string) string {
	for _, k := range []string{"password", "Password", "oldPassword", "newPassword", "salt", "access_token", "refresh_token"} {
		s = maskJSONField(s, k)
	}
	return s
}

// startLogWorker 启动异步落库工作协程
//
// 按批量大小或定时触发刷新，将缓冲日志写入 LogSink
func startLogWorker() {
	go func() {
		batch := make([]*OperationLogEntry, 0, logBatchSize)
		flush := func() {
			if len(batch) == 0 || LogSink == nil {
				batch = batch[:0]
				return
			}
			ctx := context.Background()
			for _, e := range batch {
				if err := LogSink(ctx, e); err != nil {
					g.Log().Errorf(ctx, "操作日志落库失败: %v", err)
				}
			}
			batch = batch[:0]
		}
		ticker := time.NewTicker(logFlushEvery)
		defer ticker.Stop()
		for {
			select {
			case e := <-logChan:
				batch = append(batch, e)
				if len(batch) >= logBatchSize {
					flush()
				}
			case <-ticker.C:
				flush()
			}
		}
	}()
}
