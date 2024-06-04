package utils

import (
	"strconv"
	"time"
)

type Token struct {
	Value string
}

func NewToken(id string) *Token {
	now := time.Now().Unix()
	str := id + strconv.FormatInt(now, 10) + GlobalObject.TokenSeed
	md5Str := Md5(str)
	token := Token{Value: md5Str}
	return &token
}
