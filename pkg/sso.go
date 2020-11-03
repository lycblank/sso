package sso

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var defaultSSO SingleSignOn

func Init(gdb *gorm.DB, rdb *redis.Client) {
	defaultSSO = NewMysqlRedisSSO(gdb, rdb)
}


type Loginer interface{
	Login(ctx context.Context, openid string, passwd string) (token string, refreshToken string, err error)
}

type Register interface {
	Registe(ctx context.Context, openid string, passwd string) (token string, refreshToken string, err error)
}

type Checker interface {
	Check(ctx context.Context, token string) (openid string, refreshToken string, err error)
}

type Refersher interface {
	Refresh(ctx context.Context, refreshToken string) (openid string, token string, newRefreshToken string, err error)
}

type SingleSignOn interface {
	Loginer
	Register
	Checker
	Refersher
}

func Login(ctx context.Context, openid string, passwd string) (token string, refreshToken string, err error) {
	return defaultSSO.Login(ctx, openid, passwd)
}

func Registe(ctx context.Context, openid string, passwd string) (token string, refreshToken string, err error) {
	return defaultSSO.Registe(ctx, openid, passwd)
}

func Check(ctx context.Context, token string) (openid string, refreshToken string, err error) {
	return defaultSSO.Check(ctx, token)
}

func Refresh(ctx context.Context, refreshToken string) (openid string, token string, newRefreshToken string, err error) {
	return defaultSSO.Refresh(ctx, refreshToken)
}
