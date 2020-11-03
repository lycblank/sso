
//Note that it is automatically generated, please do not modify

package db

import (
	"gorm.io/gorm"
	"bytes"
	"fmt"
	"context"
	"github.com/go-redis/redis/v8"
)

type Login struct {
	Uid int32 `gorm:"column:uid;primaryKey;comment:用户uid"`
	OpenId string `gorm:"column:open_id;comment:用户标识"`
	Password string `gorm:"column:password;comment:密码"`
	Salt string `gorm:"column:salt;comment:盐"`
	LastLoginTime int64 `gorm:"column:LastLoginTime;comment:最后一次登录时间"`
	Version int32 `gorm:"column:version"`
	CreateTime int64 `gorm:"column:create_time"`
	Deleted int32 `gorm:"column:deleted;comment:1:表示已删除，0:表示未删除"`
	DeleteTime int64 `gorm:"column:delete_time"`
}

func (l *Login) SyncScheme(ctx context.Context, gdb *gorm.DB) error {
	return gdb.AutoMigrate(l)
}

func (Login) TableName() string {
	return "login"
}

func (l *Login) DataKey() string {
	var buf bytes.Buffer
	buf.WriteString("login")
	buf.WriteString(":")
	buf.WriteString(fmt.Sprint(l.Uid))
	return buf.String()
}

func (l *Login) OK() bool {
	return l.Deleted == 0
}

func (l *Login) Sync(ctx context.Context, gdb *gorm.DB) error {
	return gdb.Save(l).Error
}

func (l *Login) Find(ctx context.Context, rdb *redis.Client, gdb *gorm.DB) error {
	if err := l.FindByCache(ctx, rdb); err == nil {
		return nil
	}
	if err := l.FindByDB(ctx, gdb); err != nil {
		return err
	}
	_ = l.SaveCache(ctx, rdb)
	return nil
}

func (l *Login) FindByDB(ctx context.Context, gdb *gorm.DB) error {
	err := gdb.First(l, l.Uid).Error
	return err
}

func (l *Login) FindByCache(ctx context.Context, rdb *redis.Client) error {
	dataKey := l.DataKey()
	pipe := rdb.Pipeline()
	cmd1 := pipe.HGet(ctx, dataKey, "uid")
	cmd2 := pipe.HGet(ctx, dataKey, "open_id")
	cmd3 := pipe.HGet(ctx, dataKey, "password")
	cmd4 := pipe.HGet(ctx, dataKey, "salt")
	cmd5 := pipe.HGet(ctx, dataKey, "last_login_time")
	if _, err := pipe.Exec(ctx); err != nil {
		 return err
	}
	val, _ := cmd1.Int64()
	l.Uid = int32(val)
	l.OpenId = cmd2.Val()
	l.Password = cmd3.Val()
	l.Salt = cmd4.Val()
	l.LastLoginTime, _ = cmd5.Int64()
	return nil
}

func (l *Login) SaveCache(ctx context.Context, rdb *redis.Client) error {
	dataKey := l.DataKey()
	pipe := rdb.Pipeline()
	cmd1 := pipe.HSet(ctx, dataKey, "uid", l.Uid)
	cmd2 := pipe.HSet(ctx, dataKey, "open_id", l.OpenId)
	cmd3 := pipe.HSet(ctx, dataKey, "password", l.Password)
	cmd4 := pipe.HSet(ctx, dataKey, "salt", l.Salt)
	cmd5 := pipe.HSet(ctx, dataKey, "last_login_time", l.LastLoginTime)
	if _, err := pipe.Exec(ctx); err != nil {
		 return err
	}
	return nil
}

func (l *Login) SetReadRedisCmd(ctx context.Context, pipe *redis.Pipeline) error {
	dataKey := l.DataKey()
	cmd1 := pipe.HGet(ctx, dataKey, "uid")
	cmd2 := pipe.HGet(ctx, dataKey, "open_id")
	cmd3 := pipe.HGet(ctx, dataKey, "password")
	cmd4 := pipe.HGet(ctx, dataKey, "salt")
	cmd5 := pipe.HGet(ctx, dataKey, "last_login_time")
	return nil
}

func (l *Login) ParseRedisCmd(ctx context.Context, cmds []redis.Cmder) (cs []redis.Cmder, err error) {
	if len(cmds) > 0 {
		if terr := cmds[0].Err(); terr != nil {
			err = terr
		}
		val, _ := cmds[0].(*redis.StringCmd).Int64()
		l.Uid = int32(val)
		cmds = cmds[1:]
	}
	if len(cmds) > 0 {
		if terr := cmds[0].Err(); terr != nil {
			err = terr
		}
		l.OpenId = cmds[0].(*redis.StringCmd).Val()
		cmds = cmds[1:]
	}
	if len(cmds) > 0 {
		if terr := cmds[0].Err(); terr != nil {
			err = terr
		}
		l.Password = cmds[0].(*redis.StringCmd).Val()
		cmds = cmds[1:]
	}
	if len(cmds) > 0 {
		if terr := cmds[0].Err(); terr != nil {
			err = terr
		}
		l.Salt = cmds[0].(*redis.StringCmd).Val()
		cmds = cmds[1:]
	}
	if len(cmds) > 0 {
		if terr := cmds[0].Err(); terr != nil {
			err = terr
		}
		l.LastLoginTime, _ = cmds[0].(*redis.StringCmd).Int64()
		cmds = cmds[1:]
	}
	return cmds, err
}
