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
	"encoding/json"
	"net"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"easygoadmin/utility/ip2region"
)

// ============================== 中间件通用工具 ==============================

// cfgBool 读取布尔配置
// 配置缺失时返回默认值
func cfgBool(ctx context.Context, key string, def bool) bool {
	v := g.Cfg().MustGet(ctx, key)
	if v.IsNil() {
		return def
	}
	return v.Bool()
}

// cfgInt 读取整数配置
// 配置缺失时返回默认值（显式配置 0 表示业务含义上的关闭，原样返回）
func cfgInt(ctx context.Context, key string, def int) int {
	v := g.Cfg().MustGet(ctx, key)
	if v.IsNil() {
		return def
	}
	return v.Int()
}

// gjsonEncode JSON 序列化
// 序列化失败时返回空串
func gjsonEncode(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// trunc 按字节截断且保证 UTF-8 完整
//
// 截断点回退到完整的字符边界，避免 MySQL 严格模式拒绝非法字节
func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && !utf8.RuneStart(s[n]) {
		n--
	}
	return s[:n] + "..."
}

// maskJSONField 脱敏 JSON 文本中的指定字段
// 将该字段的值替换为 ***
var maskReCache = map[string]*regexp.Regexp{}

func maskJSONField(s, field string) string {
	re, ok := maskReCache[field]
	if !ok {
		re = regexp.MustCompile(`("` + regexp.QuoteMeta(field) + `"\s*:\s*)"?[^",}]*"?`)
		maskReCache[field] = re
	}
	return re.ReplaceAllString(s, `${1}***`)
}

// realClientIP 获取客户端真实IP
//
// 读取配置 api.proxyIps，直连来源为可信代理时，
// 依次取 X-Forwarded-For 首段或 X-Real-IP
func realClientIP(r *ghttp.Request) string {
	remote := r.GetRemoteIp()
	if !isTrustedProxy(r.Context(), remote) {
		return remote
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip := strings.TrimSpace(strings.Split(xff, ",")[0]); ip != "" {
			return ip
		}
	}
	if real := strings.TrimSpace(r.Header.Get("X-Real-IP")); real != "" {
		return real
	}
	return remote
}

// isTrustedProxy 判断来源是否为可信代理
// 支持精确 IP 与 CIDR 两种配置形式
func isTrustedProxy(ctx context.Context, ip string) bool {
	v := g.Cfg().MustGet(ctx, "api.proxyIps")
	if v.IsNil() || v.String() == "" {
		return false
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, item := range strings.Split(v.String(), ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if item == ip {
			return true
		}
		if _, cidr, err := net.ParseCIDR(item); err == nil && cidr.Contains(parsed) {
			return true
		}
	}
	return false
}

// requestLocation 按请求IP解析归属地
//
// 基于 ip2region 实现，解析失败时返回 "-"
func requestLocation(ctx context.Context, ip string) string {
	_ = ctx
	return ip2region.Location(ip)
}
