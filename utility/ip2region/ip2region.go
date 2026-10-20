// Package ip2region 基于 ip2region xdb 的 IP 归属地查询
//
// 供登录日志与操作日志的 location 字段使用。
// 数据文件为 manifest/data/ip2region.xdb（IPv4），路径可用配置 ip2region.path 覆盖。
package ip2region

import (
	"context"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	once     sync.Once     // 保证查询器仅初始化一次
	searcher *xdb.Searcher // 归属地查询器
	mu       sync.RWMutex  // Searcher 底层共享 ReadSeeker，读操作加读锁
	loadErr  error         // 初始化错误
)

// initSearcher 懒加载查询器
//
// 失败时仅告警一次，后续查询返回空
func initSearcher() {
	once.Do(func() {
		path := "manifest/data/ip2region.xdb"
		if v, err := g.Cfg().Get(context.Background(), "ip2region.path"); err == nil && v.String() != "" {
			path = v.String()
		}
		searcher, loadErr = xdb.NewWithFileOnly(xdb.IPv4, path)
		if loadErr != nil {
			g.Log().Warningf(context.Background(), "ip2region 加载失败(%s): %v（归属地将显示为 -）", path, loadErr)
		}
	})
}

// City 根据 IP 返回城市信息
//
// 优先返回市，其次返回省，失败或内网地址返回空串。
// 区域串段位兼容两种格式，各段可能为空：
//
//	新版：国家|省|市|ISP|国家码
//	经典：国家|区域|省|市|ISP
func City(ip string) string {
	ip = strings.TrimSpace(ip)
	if ip == "" || ip == "-" {
		return ""
	}
	// 内网/回环直接归空（xdb 返回 Reserved/LAN 等无意义值）
	if isPrivateIP(ip) {
		return ""
	}
	initSearcher()
	if searcher == nil {
		return ""
	}
	mu.RLock()
	region, err := searcher.Search(ip)
	mu.RUnlock()
	if err != nil || region == "" {
		return ""
	}
	// 典型内网占位值
	for _, kw := range []string{"Reserved", "LAN", "None", "广播"} {
		if strings.Contains(region, kw) {
			return ""
		}
	}
	parts := strings.Split(region, "|")
	var prov, city string
	if len(parts) >= 4 {
		p1 := strings.TrimSpace(parts[1])
		if isISP(parts[3]) || hasAdminSuffix(p1) {
			prov, city = p1, strings.TrimSpace(parts[2]) // 新版：第二段即省份
		} else {
			prov, city = strings.TrimSpace(parts[2]), strings.TrimSpace(parts[3]) // 经典：第三段起为省/市
		}
		switch {
		case city != "" && city != "0" && prov != "" && prov != "0":
			return prov + city
		case city != "" && city != "0":
			return city
		case prov != "" && prov != "0":
			return prov
		}
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" && p != "0" && !isISP(p) {
			return p
		}
	}
	return ""
}

// Location 返回归属地字段值
//
// 供登录日志与操作日志使用，无归属地时返回 "-"
func Location(ip string) string {
	if c := City(ip); c != "" {
		return c
	}
	return "-"
}

// isPrivateIP 判断是否为内网或回环地址
func isPrivateIP(ip string) bool {
	return strings.HasPrefix(ip, "127.") ||
		strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "172.16.") ||
		strings.HasPrefix(ip, "172.17.") || strings.HasPrefix(ip, "172.18.") ||
		strings.HasPrefix(ip, "172.19.") || strings.HasPrefix(ip, "172.2") ||
		strings.HasPrefix(ip, "172.30.") || strings.HasPrefix(ip, "172.31.") ||
		ip == "::1" || ip == "0.0.0.0"
}

// isISP 粗略识别 ISP 文本
//
// 避免将运营商名称误判为城市
func isISP(s string) bool {
	for _, kw := range []string{"电信", "联通", "移动", "铁通", "教育", "科研", "广电"} {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}

// hasAdminSuffix 判断是否为行政区划名称
//
// 识别省、市、自治区、特别行政区等后缀，用于识别新版段位的省份段
func hasAdminSuffix(s string) bool {
	for _, kw := range []string{"省", "市", "自治区", "特别行政区"} {
		if strings.HasSuffix(s, kw) {
			return true
		}
	}
	return false
}
