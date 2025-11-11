package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrTokenNotValidYet = errors.New("token is not valid yet")
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // user, admin, superadmin
	jwt.RegisteredClaims
}

// Config JWT 配置
type Config struct {
	Secret     string        // 密钥
	Issuer     string        // 签发者
	ExpireTime time.Duration // 过期时间
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username, role string, config Config) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(config.ExpireTime)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新 Token
func RefreshToken(tokenString string, secret string, config Config) (string, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return "", err
	}

	// 即使 token 过期，如果在允许刷新的时间范围内，也可以刷新
	// 这里简化处理，直接生成新 token
	return GenerateToken(claims.UserID, claims.Username, claims.Role, config)
}

// ValidateToken 验证 Token 是否有效
func ValidateToken(tokenString string, secret string) error {
	_, err := ParseToken(tokenString, secret)
	return err
}

