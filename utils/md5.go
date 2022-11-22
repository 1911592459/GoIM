package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

//大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

//加密,给md5后面加一个随机数salt
func MakePassword(pwd, salt string) string {
	return Md5Encode(pwd + salt)
}

//解密
func ValidPassword(pwd, salt string, password string) bool {
	return Md5Encode(pwd+salt) == password
}
