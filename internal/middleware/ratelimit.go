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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"easygoadmin/internal/consts"
	"easygoadmin/utility/cache"
	"easygoadmin/utility/response"
)

// ============================== API 限流中间件 ==============================
//
// 固定窗口计数：键按窗口起点对齐（ratelimit:{ip}:{path}:{窗口起点}），
// 计数存 utility/cache（Redis 优先、内存降级），多实例部署时共享计数。
//
// 规则：
//   - 仅限 /api/ 前缀接口；文档（/api.json）、静态资源、健康检查不受限
//   - 认证入口（登录/验证码/令牌兑换）按 rateLimit.authLimit 单独更严阈值
//   - 计数存储异常时放行（限流不阻断业务）

// RateLimit API 限流
func RateLimit(r *ghttp.Request) {
	if !cfgBool(r.Context(), "rateLimit.enable", true) {
		r.Middleware.Next()
		return
	}
	path := r.URL.Path
	if !isRateLimitedPath(path) {
		r.Middleware.Next()
		return
	}
	window := cfgInt(r.Context(), "rateLimit.windowSeconds", 60)
	limit := cfgInt(r.Context(), "rateLimit.limit", 120)
	if isAuthPath(path) {
		limit = cfgInt(r.Context(), "rateLimit.authLimit", 10)
	}
	if limit <= 0 || window <= 0 {
		r.Middleware.Next()
		return
	}

	// 窗口起点对齐：同一窗口内共享计数键
	windowStart := time.Now().Unix() / int64(window) * int64(window)
	key := fmt.Sprintf("ratelimit:%s:%s:%d", realClientIP(r), path, windowStart)
	n, err := incrCounter(r, key, window)
	if err != nil {
		r.Middleware.Next()
		return
	}
	if n > int64(limit) {
		response.Write(r, response.Fail(gerror.NewCode(consts.CodeTooManyRequests)))
		return
	}
	r.Middleware.Next()
}

// isRateLimitedPath 判断是否为参与限流的接口路径
// /api.json 与 /swagger 不带 /api/ 前缀，自然排除；健康检查显式豁免
func isRateLimitedPath(path string) bool {
	if path == "/api/health" {
		return false
	}
	return path == "/api" || strings.HasPrefix(path, "/api/")
}

// isAuthPath 判断是否为认证入口（适用更严格的独立阈值）
func isAuthPath(path string) bool {
	switch path {
	case "/api/login", "/api/captcha", "/api/oauth2/token":
		return true
	}
	return false
}

// incrCounter 固定窗口计数自增
// 键已含窗口起点，TTL 仅用于回收存储，续期不影响窗口归属
func incrCounter(r *ghttp.Request, key string, windowSeconds int) (int64, error) {
	ttl := windowSeconds + 5
	if raw, ok := cache.Get(r.Context(), key); ok {
		n, _ := strconv.ParseInt(raw, 10, 64)
		n++
		_ = cache.Set(r.Context(), key, strconv.FormatInt(n, 10), ttl)
		return n, nil
	}
	if cache.SetNX(r.Context(), key, "1", ttl) {
		return 1, nil
	}
	// 并发下 SetNX 未抢到：读现值再加一
	if raw, ok := cache.Get(r.Context(), key); ok {
		n, _ := strconv.ParseInt(raw, 10, 64)
		n++
		_ = cache.Set(r.Context(), key, strconv.FormatInt(n, 10), ttl)
		return n, nil
	}
	return 1, nil
}
