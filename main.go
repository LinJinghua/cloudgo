package main

import (
	"fmt"
	"os"
	"log"

	"github.com/LinJinghua/cloudgo/service"
	flag "github.com/spf13/pflag"
)

const (
	// PORT : default port
	PORT string = "8080"
)

var port string

func init()  {
	// 设置命令行参数
	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	// 解析命令行参数
	flag.Parse()

	// 解析监听端口
	// 用户设置第一，系统环境次之，默认最后
	if len(*pPort) != 0 {
		port = *pPort
	} else {
		if port = os.Getenv("PORT"); len(port) == 0 {
			port = PORT
		}
	}
}

func main() {
	// 开始服务
	if err := service.Run(port); err != nil {
        log.Fatal("start fasthttp fail:", err)
    } else {
		fmt.Println("listening on :" + port)
	}
}
