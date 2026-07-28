package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lsjhtang/kid"
)

func main() {
	engin := kid.New(kid.WithConfigs("./config.yaml.yaml"))
	engin.GET("", func(ctx *kid.Context) interface{} {
		fmt.Println(1)
		return gin.H{"hello": "world"}
	}, func(ctx *kid.Context) {
		fmt.Println(2)
		ctx.Next()
	}, func(ctx *kid.Context) {
		fmt.Println(3)
		ctx.Next()
	})
	engin.Launch()
}
