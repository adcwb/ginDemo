package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// StringToMD5 字符串转MD5
func StringToMD5(data string) string {
	signByte := []byte(data)
	hash := md5.New()
	hash.Write(signByte)
	return strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}
