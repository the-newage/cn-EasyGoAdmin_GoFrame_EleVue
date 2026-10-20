package ua

import "strings"

// ============================== User-Agent 解析 ==============================

// ParseOS 从 User-Agent 解析操作系统
//
// 供操作日志与登录日志记录使用，无法识别时返回 Unknown
func ParseOS(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "windows nt 10"):
		return "Windows 10/11"
	case strings.Contains(ua, "windows nt"):
		return "Windows"
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"):
		return "iOS"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "mac os x"), strings.Contains(ua, "macintosh"):
		return "macOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	}
	return "Unknown"
}

// ParseBrowser 从 User-Agent 解析浏览器
//
// 供操作日志与登录日志记录使用，无法识别时返回 Unknown
func ParseBrowser(userAgent string) string {
	ua := userAgent
	switch {
	case strings.Contains(ua, "Edg/"), strings.Contains(ua, "Edge/"):
		return "Edge"
	case strings.Contains(ua, "OPR/"), strings.Contains(ua, "Opera"):
		return "Opera"
	case strings.Contains(ua, "Firefox/"):
		return "Firefox"
	case strings.Contains(ua, "Chrome/"), strings.Contains(ua, "CriOS"):
		return "Chrome"
	case strings.Contains(ua, "Safari/"):
		return "Safari"
	case strings.Contains(ua, "curl/"):
		return "curl"
	}
	return "Unknown"
}
