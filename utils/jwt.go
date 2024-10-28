package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const TokenExpireDuration = time.Hour * 2

var Secret = []byte("W2JzDLCNpdO8Xq94ZlAc6Hge7fshbvtymkGKRn3YQruaM1BToIjVwU0SxEF5PiGY")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(username string) (string, error) {
	expirationTime := time.Now().Add(TokenExpireDuration)

	claims := MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (any, error) {
		// 确保签名方法是我们使用的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return Secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
