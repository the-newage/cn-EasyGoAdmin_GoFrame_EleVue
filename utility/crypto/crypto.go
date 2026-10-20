package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// ============================== 密码加密与校验 ==============================

// GenerateSalt 生成随机外部盐值
//
// 取 8 字节随机数经 base64url 编码后截取 10 个字符，去除填充
func GenerateSalt() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:10]
}

// HashPassword 双重加盐加密
//
// 使用 bcrypt 对 password 与 salt 拼接后的内容加密，返回哈希值；
// 两值分别写入 password 与 salt 列，salt 为空时等价于直接加密密码，兼容旧数据
func HashPassword(password, salt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword 校验密码
//
// 使用 bcrypt 比对哈希值与 password 加 salt 的组合，salt 为空时直接比对，兼容旧数据
func VerifyPassword(hashedPassword, password, salt string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt)) == nil
}

// IsBcrypt 判断是否为 bcrypt 格式
//
// 旧版 md5 加 salt 的数据返回 false
func IsBcrypt(hashed string) bool {
	return len(hashed) > 4 && (hashed[:4] == "$2a$" || hashed[:4] == "$2b$" || hashed[:4] == "$2y$")
}
