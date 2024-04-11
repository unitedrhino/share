package users

import (
	"gitee.com/i-Things/share/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Custom claims structure
type OpenClaims struct {
	Account    string //账号
	TenantCode string `json:",string"`
	jwt.RegisteredClaims
}

func GetOpenJwtToken(secretKey string, t time.Time, seconds int64, account string, tenantCode string) (string, error) {
	IssuedAt := jwt.NewNumericDate(t)
	claims := OpenClaims{
		TenantCode: tenantCode,
		Account:    account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.Add(time.Duration(seconds) * time.Second)),
			IssuedAt:  IssuedAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 更新token
func RefreshOpenToken(tokenString string, secretKey string, AccessExpire int64) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OpenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*OpenClaims); ok && token.Valid {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
		return CreateToken(secretKey, *claims)
	}
	return "", errors.TokenInvalid
}
