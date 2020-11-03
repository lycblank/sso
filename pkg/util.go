package sso

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

func AddSalt(oldPassword string, salt string) string {
	hash := md5.New()
	io.WriteString(hash, oldPassword + salt)
	return fmt.Sprintf("%X", hash.Sum(nil))
}

func GetTimeUnix() int64 {
	return time.Now().Unix()
}


func CreateToken(val string) string {
	hash := md5.New()
	bw := bufio.NewWriter(hash)
	bw.WriteString(val)
	bw.WriteString(fmt.Sprintf("%d", time.Now().UnixNano()))
	buf := make([]byte, 8)
	io.ReadFull(rand.Reader, buf)
	bw.Write(buf)
	bw.Flush()
	return fmt.Sprintf("%X", hash.Sum(nil))
}

func GetTokenKey(token string) string {
	return fmt.Sprintf("sso:token:%s", token)
}

func GetRefreshTokenKey(refreshToken string) string {
	return fmt.Sprintf("sso:refreshtoken:%s", refreshToken)
}

func GetOpenidKey(openid string) string {
	return fmt.Sprintf("sso:openid:%s", openid)
}
