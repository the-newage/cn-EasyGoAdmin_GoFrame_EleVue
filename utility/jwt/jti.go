package jwt

import (
	"crypto/rand"
	"encoding/hex"
)

// ============================== 令牌标识 ==============================

// newJTI 生成令牌唯一标识
//
// 生成 16 字节随机数并转为 32 位十六进制字符串，用作黑名单键；
// 随机源读取失败时返回固定兜底值
func newJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "jti-fallback"
	}
	return hex.EncodeToString(b)
}
