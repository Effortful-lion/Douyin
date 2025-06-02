package main

import (
	"Douyin/config"
	"fmt"
)

func main() {
	config.InitConfig()
	fmt.Println(config.Config)
}
