package dir

import (
	"log"
	"os"
)

//文件或目录是否存在
func IsExist(dirPath string) bool {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	log.Fatal("dir is exist err:", err)
	return false
}

// 创建文件或者目录
func MkdirAll(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("make dir or file error:", err)
	}
}
