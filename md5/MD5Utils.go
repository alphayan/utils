package MD5Utils

import (
	"crypto/md5"
	"encoding/hex"
)

/**
md5加密
*/
func MD5Encrypt(value string) string {
	md := md5.New()
	md.Write([]byte(value))
	cipherStr := md.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
