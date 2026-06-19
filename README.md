# Kid

Kid 是一个基于 [Gin](https://github.com/gin-gonic/gin) 二次封装的 Go Web 框架。它保留 Gin 的路由和中间件模型，同时补充了统一上下文、返回值自动响应、配置加载、日志、依赖注入、数据库、缓存、认证守卫和常用业务插件。

这个仓库更适合作为业务服务的基础框架，而不是一个只包含 HTTP 路由的轻量库。项目内置了 MySQL、Redis、JWT、Token、OAuth、短信、支付、对象存储、验证码、WebSocket、定时任务、消息队列、分布式锁等能力。

## 环境要求

- Go 1.22 或以上。
- 使用方项目需要能访问 `github.com/leor-w/kid` 及其依赖。
- MySQL、Redis、支付、短信、对象存储等能力按需启用；未使用的插件不需要在业务代码中注册。

## 快速开始

安装依赖：

```bash
go get github.com/leor-w/kid
```

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

创建 `config.yaml`：

```yaml
app:
  run_mode: debug
```

启动服务：

```bash
go run . server -c ./config.yaml
```

访问接口：

```bash
curl http://127.0.0.1:8080/ping
```

更多接入细节见 [快速开始](docs/quickstart.md)。

## 核心能力

| 能力 | 入口 | 说明 |
| --- | --- | --- |
| Web 服务 | `kid.New`、`Launch` | 初始化配置、日志、Gin Engine、默认日志中间件和 Recovery。 |
| 路由封装 | `GET`、`POST`、`PUT`、`DELETE`、`Group` | Handler 使用 `func(*kid.Context) any`，返回值默认序列化为 JSON。 |
| 响应类型 | `SetResponseType` | 支持 JSON、XML、YAML、HTML、String、Redirect、File 等响应模式。 |
| 上下文增强 | `kid.Context` | 支持参数查找、默认值、用户绑定、结构体校验、字段规则校验。 |
| 配置 | `config.New`、`kid.WithConfigs` | 基于 Viper 读取本地 YAML，可加载多个配置文件。 |
| 日志 | `logger`、`plugin/logrus` | 默认初始化 Logrus，可扩展 Hook、Formatter、日志切割。 |
| 依赖注入 | `Registry`、`Populate`、`inject:""` | 插件通过 `Provide(context.Context)` 进入容器，支持别名和接口注入。 |
| 数据库 | `database/mysql` | 封装 Gorm 连接、连接池、分页、条件、搜索、自动迁移。 |
| Redis | `database/redis` | 封装 Redis Client、连接池、TLS、Repository 接口。 |
| 认证 | `guard/jwt`、`guard/token` | 提供统一 `guard.Guard` 接口，支持发行、验证、刷新、吊销。 |
| 插件 | `plugin/*` | 覆盖支付、OAuth、短信、对象存储、验证码、队列、锁、WebSocket 等场景。 |

## 推荐目录结构

```text
your-service/
├── cmd/
│   └── server/
│       └── main.go
├── config.yaml
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
└── go.mod
```

Kid 不强制业务目录结构。建议把 HTTP Handler、业务服务、仓储和模型分层管理，插件注册和路由注册集中放在启动入口附近。

## 常用命令

启动服务并指定配置：

```bash
go run . server -c ./config.yaml
```

数据库自动迁移：

```bash
go run . database --automigrate -c ./config.yaml
```

使用命令行参数指定数据库连接并执行迁移：

```bash
go run . database --automigrate \
  --database app \
  --address 127.0.0.1 \
  --port 3306 \
  --user root \
  --password secret
```

命令行由 Cobra 提供，`help` 子命令会输出可用命令和参数。

## 使用文档

- [快速开始](docs/quickstart.md)：最小服务、路由、参数、中间件、响应和校验。
- [配置说明](docs/configuration.md)：配置文件加载、读取 API、常用配置片段和敏感配置建议。
- [数据库与缓存](docs/database.md)：MySQL、Redis、Gorm Scope、自动迁移和仓储接口。
- [认证与授权](docs/auth.md)：JWT、签名 Token、Redis Store、黑名单和安全建议。
- [插件系统](docs/plugins.md)：插件注册、依赖注入、别名、多实例和插件能力索引。
- [运维排障](docs/operations.md)：启动流程、日志、Recovery、常见错误、性能和并发建议。

## 生产使用注意事项

1. 不要裸跑服务。生产环境建议使用 Docker、systemd、Kubernetes 或其他进程管理方式托管进程、日志和重启策略。
2. 显式传入配置文件，例如 `kid.WithConfigs("./config.yaml")`，避免依赖默认路径。
3. MySQL、Redis、JWT、支付、短信、对象存储等敏感配置不要提交真实密钥到仓库。
4. 对外部资源初始化失败会直接 `panic` 或返回明确错误，应在部署前通过配置检查和健康检查提前发现。
5. 根据流量设置 MySQL 与 Redis 连接池参数；批量任务、延迟任务和高并发外部调用优先评估协程池、队列或限流插件。
6. 默认 Recovery 会返回 `{"code":50000,"message":"内部错误"}` 并记录堆栈日志，线上排障应优先检索服务日志。

## 当前文档边界

本次文档只描述源码中已经存在的公开 API、配置样例和可确认行为。支付、OAuth、短信、云存储等第三方插件的完整业务流程依赖平台侧规则，当前只提供配置入口和接入索引。
