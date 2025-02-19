package users

import (
	"gitee.com/unitedrhino/share/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Custom claims structure
type LoginClaims struct {
	UserID  int64 `json:",string"`
	AppCode string
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserID     int64  `json:",string"`
	Account    string //账号
	RoleIDs    []int64
	RoleCodes  []string
	TenantCode string `json:",string"`
	IsAdmin    int64
	IsAllData  int64
}

func GetLoginJwtToken(secretKey string, userID int64, appCode string, id string) (string, error) {
	claims := LoginClaims{
		UserID:  userID,
		AppCode: appCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: id,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 更新token
func RefreshLoginToken(tokenString string, secretKey string, AccessExpire int64) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*LoginClaims); ok && token.Valid {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(AccessExpire) * time.Second))
		return CreateToken(secretKey, *claims)
	}
	return "", errors.TokenInvalid
}
