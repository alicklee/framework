package utils

import (
	"crypto/md5"
	"encoding/hex"
)

/*
MD5方法简单封装
*/
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
