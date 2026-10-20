package jwt

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
)

// ============================== JWT 令牌 ==============================

// TokenType 令牌类型
const (
	TokenTypeAccess  = "access"  // 访问令牌
	TokenTypeRefresh = "refresh" // 刷新令牌
)

// Claims JWT 载荷
//
// 包含身份信息、租户信息与令牌语义
type Claims struct {
	UserID   int    `json:"user_id"`   // 用户ID
	TenantID int    `json:"tenant_id"` // 租户ID
	Username string `json:"username"`  // 用户名
	Realname string `json:"realname"`  // 真实姓名
	Type     string `json:"typ"`       // 令牌类型：access 或 refresh
	jwt.RegisteredClaims
}

// cfg 读取字符串配置
//
// 配置缺失或为空时返回默认值
func cfg(ctx context.Context, key string, def ...string) string {
	v := g.Cfg().MustGet(ctx, key)
	if v.IsNil() || v.String() == "" {
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}
	return v.String()
}

// cfgDuration 读取时长配置
//
// 配置值为秒数，缺失或非正时返回默认值
func cfgDuration(ctx context.Context, key string, def time.Duration) time.Duration {
	v := g.Cfg().MustGet(ctx, key)
	if v.IsNil() || v.Int() <= 0 {
		return def
	}
	return time.Duration(v.Int()) * time.Second
}

// Sign 签发令牌
//
// 返回签名后的令牌与载荷，载荷中的 ID 即 jti，可作为黑名单键
func Sign(ctx context.Context, tokenType string, userID, tenantID int, username, realname string) (string, *Claims, error) {
	now := time.Now()
	expireSec := cfgDuration(ctx, "jwt.accessTokenExpire", 2*time.Hour)
	if tokenType == TokenTypeRefresh {
		expireSec = cfgDuration(ctx, "jwt.refreshTokenExpire", 168*time.Hour)
	}
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
		Realname: realname,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        newJTI(),
			Issuer:    cfg(ctx, "jwt.issuer", "easygoadmin"),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expireSec)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg(ctx, "jwt.secret", "easygoadmin")))
	if err != nil {
		return "", nil, gerror.Wrap(err, "令牌签发失败")
	}
	return signed, claims, nil
}

// Parse 解析并校验令牌
//
// 校验签名算法、签名、过期时间与签发者
func Parse(ctx context.Context, tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errInvalidToken()
			}
			return []byte(cfg(ctx, "jwt.secret", "easygoadmin")), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(cfg(ctx, "jwt.issuer", "easygoadmin")),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, errInvalidToken()
	}
	return claims, nil
}

// errInvalidToken 构造无效令牌错误
//
// 状态码保持 401，供前端识别并跳转登录页
func errInvalidToken() error {
	return gerror.NewCode(gcode.New(401, "Token无效或已失效", nil))
}

// ExpireDuration 返回令牌有效期
//
// 用于黑名单 TTL 计算
func ExpireDuration(ctx context.Context, tokenType string) time.Duration {
	if tokenType == TokenTypeRefresh {
		return cfgDuration(ctx, "jwt.refreshTokenExpire", 168*time.Hour)
	}
	return cfgDuration(ctx, "jwt.accessTokenExpire", 2*time.Hour)
}
