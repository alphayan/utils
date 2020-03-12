package config

import (
	"flag"
	"fmt"
)

var (
	//配置文件路径
	configPath = flag.String("config", "", "config file path")
)

func init() {
	flag.Parse()
	_, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Println("load config err")
		panic(err)
	}
}
