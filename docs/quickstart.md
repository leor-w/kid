# 快速开始

本文说明如何在一个新 Go 服务中接入 Kid，并使用它完成最小 Web 服务、路由、中间件、参数读取、响应和校验。

## 1. 创建项目

```bash
mkdir demo-kid-service
cd demo-kid-service
go mod init example.com/demo-kid-service
go get github.com/leor-w/kid
```

Kid 当前模块声明为 `github.com/leor-w/kid`，Go 版本为 1.22。

## 2. 编写入口

创建 `main.go`：

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid"
)

func main() {
	app := kid.New(
		kid.WithAppName("demo"),
		kid.WithConfigs("./config.yaml"),
		kid.WithRunMode("app.run_mode"),
	)

	app.GET("/ping", func(ctx *kid.Context) any {
		return gin.H{"message": "pong"}
	})

	app.Launch(":8080")
}
```

`kid.New` 会执行初始化流程：处理命令行、加载配置、初始化日志、设置 Gin 运行模式、创建 Gin Engine、注册默认请求日志中间件和 Recovery。

`kid.WithRunMode("app.run_mode")` 传入的是配置 key，不是直接传入 `debug`、`release` 或 `test`。如果配置 key 没有读到值，源码会回退到 `gin.TestMode`。

## 3. 编写配置

创建 `config.yaml`：

```yaml
app:
  run_mode: debug
```

生产环境建议显式传入配置文件，不依赖默认路径。

## 4. 启动和访问

```bash
go run . server -c ./config.yaml
```

访问测试接口：

```bash
curl http://127.0.0.1:8080/ping
```

预期返回：

```json
{"message":"pong"}
```

## 5. 路由和分组

Kid 的路由封装在 `RouterGroup` 上，支持 `GET`、`POST`、`PUT`、`DELETE` 和 `Group`。

```go
api := app.Group("/api")

api.GET("/users/:id", func(ctx *kid.Context) any {
	id := ctx.FindInt64("id")
	return gin.H{"id": id}
})

api.POST("/users", func(ctx *kid.Context) any {
	name := ctx.FindString("name")
	return gin.H{"name": name}
})
```

Handler 签名是：

```go
func(ctx *kid.Context) any
```

返回值不为 `nil` 时，Kid 默认使用 JSON 响应。返回 `nil` 时，封装层会 `Abort` 当前请求，适合中间件或 Handler 已经自行写出响应的场景。

## 6. 中间件

Kid 中间件签名是：

```go
func(ctx *kid.Context)
```

示例：

```go
func RequestID() kid.Middleware {
	return func(ctx *kid.Context) {
		ctx.Set("request_id", "req-001")
		ctx.Next()
	}
}

app.GET("/trace", func(ctx *kid.Context) any {
	return gin.H{"request_id": ctx.GetString("request_id")}
}, RequestID())
```

也可以对分组使用中间件：

```go
api := app.Group("/api", RequestID())
```

已有的 Gin Handler 可以用 `kid.WrapHandler` 包装。

## 7. 参数读取

`kid.Context` 提供一组 `Find*` 方法，会按以下顺序查找参数：

1. 路由参数 `Param`
2. Query 参数
3. Gin Context 中的 key
4. Header
5. PostForm

示例：

```go
app.GET("/users/:id", func(ctx *kid.Context) any {
	return gin.H{
		"id":     ctx.FindInt64("id"),
		"source": ctx.FindDefaultString("source", "web"),
	}
})
```

常用方法包括 `FindString`、`FindInt`、`FindInt64`、`FindBool`、`FindFloat64` 以及对应的 `FindDefault*`。

## 8. 响应类型

默认响应是 JSON。可以通过 `ctx.SetResponseType` 指定响应格式：

```go
app.GET("/text", func(ctx *kid.Context) any {
	ctx.SetResponseType(kid.ResponseTypeString)
	return "hello kid"
})

app.GET("/redirect", func(ctx *kid.Context) any {
	ctx.SetResponseType(kid.ResponseTypeRedirect)
	return "https://example.com"
})
```

已定义的响应类型包括 JSON、XML、YAML、HTML、IndentedJSON、SecureJSON、JSONP、AsciiJSON、PureJSON、TOML、ProtoBuf、String、Redirect、File 等。

## 9. 参数校验

字段规则校验：

```go
app.GET("/search", func(ctx *kid.Context) any {
	if err := ctx.ValidFiled(kid.Rules{
		"page": "omitempty,numeric",
		"size": "omitempty,numeric",
	}); err != nil {
		return gin.H{"error": err.Error()}
	}

	return gin.H{
		"page": ctx.FindDefaultInt("page", 1),
		"size": ctx.FindDefaultInt("size", 10),
	}
})
```

结构体校验可以使用 `kid.Struct`：

```go
type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
}

app.POST("/users", func(ctx *kid.Context) any {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return gin.H{"error": err.Error()}
	}
	if err := kid.Struct(req); err != nil {
		return gin.H{"error": err.Error()}
	}
	return gin.H{"name": req.Name}
})
```

框架也提供 `ctx.Valid`，其内部会先绑定再执行 `kid.Struct`。复杂请求建议显式绑定后再调用 `kid.Struct`，排障更直接。

## 10. 下一步

- 需要多配置文件或配置读取方式，见 [配置说明](configuration.md)。
- 需要 MySQL、Redis 和仓储接口，见 [数据库与缓存](database.md)。
- 需要 JWT、Token 和 Redis Store，见 [认证与授权](auth.md)。
- 需要支付、短信、OAuth、对象存储等能力，见 [插件系统](plugins.md)。
