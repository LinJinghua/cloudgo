package service

import (
	"encoding/json"
	"fmt"

    "github.com/valyala/fasthttp"
    "github.com/buaazp/fasthttprouter"
)

func httpHandle(ctx *fasthttp.RequestCtx) {
	// fasthttprouter.RequestCtx.UserValue() 获得路由匹配得到的参数，如规则 /hello/:id 中的 :id
	data, _ := json.Marshal(struct{ Test string }{"Hello " + ctx.UserValue("id").(string)})
	fmt.Fprintln(ctx, string(data[:]))
}

// Run : run service at the port
func Run(port string) error {
    // 使用 fasthttprouter 创建路由
    router := fasthttprouter.New()

	// 设置访问的路由，处理逻辑
	router.GET("/hello/:id", httpHandle)

	// 设置监听的端口
	return fasthttp.ListenAndServe("0.0.0.0:" + port, router.Handler)
}
