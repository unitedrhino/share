package users

import (
	"gitee.com/i-Things/share/errors"
	"github.com/golang-jwt/jwt/v5"
)

// 创建一个token
func CreateToken(secretKey string, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseTokenWithFunc(claim jwt.Claims, tokenString string, f jwt.Keyfunc) error {
	token, err := jwt.ParseWithClaims(tokenString, claim, f)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return errors.TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return errors.TokenMalformed
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return errors.TokenNotValidYet
		default:
			return errors.TokenInvalid
		}
	}
	if token != nil {
		if token.Valid {
			return nil
		}
		return errors.TokenInvalid

	} else {
		return errors.TokenInvalid
	}
}

// 解析 token
func ParseToken(claim jwt.Claims, tokenString string, secretKey string) error {
	return ParseTokenWithFunc(claim, tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
}
