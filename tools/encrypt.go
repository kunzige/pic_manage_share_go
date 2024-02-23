package tools

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func Md5(s string) string {
	str := []byte(s)
	hash := md5.Sum(str)
	encoded := fmt.Sprintf("%x", hash) //encoded就是加密后的字符串
	return strings.ToUpper(encoded)
}
