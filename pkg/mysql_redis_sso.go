package sso

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lycblank/sso/pkg/db"
	"gorm.io/gorm"
	"io"
	"time"
)

type MysqlRedisSSO struct {
	gdb *gorm.DB
	rd *redis.Client
}

func NewMysqlRedisSSO(gdb *gorm.DB, rd *redis.Client) *MysqlRedisSSO {
	return &MysqlRedisSSO{
		gdb:gdb,
		rd:rd,
	}
}

func (mr *MysqlRedisSSO) Login(ctx context.Context, openid string, password string) (token string, refreshToken string, err error) {
	login := &db.Login{
		OpenId:openid,
	}
	err = mr.gdb.First(login, "open_id = ?", openid).Error
	if err != nil {
		return
	}

	if !login.OK() {
		err =UserNotExists
		return
	}

	newPassword := AddSalt(password, login.Salt)
	if login.Password != newPassword {
		err = UserOrPasswordIncorrect
		return
	}
	login.LastLoginTime = GetTimeUnix()

	err = login.Sync(ctx, mr.gdb)
	if err != nil {
		return
	}

	token = CreateToken(openid+password)
	refreshToken = CreateToken(token)

	pipe := mr.rd.Pipeline()
	defer pipe.Close()
	mr.saveToken(ctx, openid, token, refreshToken, pipe)
	_, _ = pipe.Exec(ctx)
	return
}


func (mr *MysqlRedisSSO) Registe(ctx context.Context, openid string, password string) (token string, refreshToken string, err error) {
	buf := make([]byte, 8)
	io.ReadFull(rand.Reader, buf)
	salt := fmt.Sprintf("%X", buf)
	newPassword := AddSalt(password, salt)
	err = mr.gdb.Create(&db.Login{
		OpenId:        openid,
		Password:      newPassword,
		Salt:          salt,
		CreateTime:    GetTimeUnix(),
	}).Error
	if err != nil {
		return
	}
	return mr.Login(ctx, openid, password)
}


func (mr *MysqlRedisSSO) Check(ctx context.Context, token string) (openid string, refreshToken string, err error) {
	tokenKey := GetTokenKey(token)
	var vals []interface{}
	vals, err = mr.rd.HMGet(ctx, tokenKey, "open_id", "refresh_token").Result()
	if err != nil {
		return
	}
	openid, _ = vals[0].(string)
	refreshToken, _ = vals[1].(string)
	if openid == "" || refreshToken == "" {
		err = TokenIncorrect
		return
	}
	return
}

func (mr *MysqlRedisSSO) Refresh(ctx context.Context, refreshToken string) (openid string, token string, newRefreshToken string, err error) {
	refreshTokenKey := GetRefreshTokenKey(refreshToken)
	var vals []interface{}
	vals, err = mr.rd.HMGet(ctx, refreshTokenKey, "open_id", "token").Result()
	if err != nil {
		return
	}
	token, _ = vals[0].(string)
	openid, _ = vals[1].(string)
	if openid == "" || token == "" {
		err = RefreshTokenIncorrect
		return
	}

	token = CreateToken(openid)
	newRefreshToken = CreateToken(token)
	pipe := mr.rd.Pipeline()
	defer pipe.Close()
	mr.saveToken(ctx, openid, token, newRefreshToken, pipe)
	_, _ = pipe.Exec(ctx)
	return
}


func (mr *MysqlRedisSSO) saveToken(ctx context.Context, openid string, token string, refreshToken string, pipe redis.Pipeliner) {
	openidKey := GetOpenidKey(openid)
	if vals, err := mr.rd.HMGet(ctx, openidKey, "token", "refresh_token").Result(); err == nil {
		oldToken := vals[0].(string)
		oldRefreshToken := vals[1].(string)
		pipe.Del(ctx, GetOpenidKey(openid), GetTokenKey(oldToken), GetRefreshTokenKey(oldRefreshToken))
	}

	tokenKey := GetTokenKey(token)
	_ = pipe.HMSet(ctx, tokenKey, "open_id", openid, "refresh_token", refreshToken)
	_ =  pipe.Expire(ctx, token, time.Hour*48)

	refreshTokenKey := GetRefreshTokenKey(refreshToken)
	_ = pipe.HMSet(ctx, refreshTokenKey, "open_id", openid, "token", token)
	_ = pipe.Expire(ctx, token, time.Hour*96)

	openidTokenKey := GetOpenidKey(openid)
	_ =  pipe.HMSet(ctx, openidTokenKey, "refresh_token", refreshToken, "token", token)
}
