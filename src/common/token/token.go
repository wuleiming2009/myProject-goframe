package token

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"

	"myProject/cache"
	"myProject/common/response"
	"myProject/conf"
)

const (
	RedisUserHashKey = "user_agent"
)

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
}

func verifyAction(userToken string) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(userToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		cfg, err := conf.GetAppConfig()
		if err != nil {
			return nil, err
		}
		return []byte(cfg.JwtSecret), nil
	})
	if err != nil {
		return nil, response.ErrRelogin
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, response.ErrRelogin
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, response.ErrRelogin
	}
	return claims, nil
}

func getToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	cfg, err := conf.GetAppConfig()
	if err != nil {
		return "", err
	}
	signedToken, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return "", response.ErrInternal
	}
	return signedToken, nil
}

func AddToken(ctx context.Context, userId uint64) (string, error) {
	cfg, err := conf.GetAppConfig()
	if err != nil {
		return "", err
	}
	claims := &JWTClaims{}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * cfg.TokenExpireTime).Unix()
	signedToken, err := getToken(claims)
	if err != nil {
		return "", err
	}

	// redis中插入用户hash
	rds, err := cache.RedisClient()
	if err != nil {
		return "", err
	}

	strCmd := rds.HSet(ctx, RedisUserHashKey, signedToken, userId)
	if strCmd.Err() != nil {
		return "", err
	}

	return signedToken, nil
}

func GetUserIdByToken(ctx context.Context, userToken string) (uint64, error) {

	_, err := verifyAction(userToken)
	if err != nil {
		return 0, err
	}
	// redis中查用户hash
	rds, err := cache.RedisClient()
	if err != nil {
		return 0, err
	}
	strCmd := rds.HGet(ctx, RedisUserHashKey, userToken)
	if strCmd.Err() != nil {
		return 0, err
	}
	userId, err := strCmd.Int64()
	if err != nil {
		return 0, err
	}
	return uint64(userId), nil
}
