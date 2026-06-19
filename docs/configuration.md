# 配置说明

Kid 的配置能力基于 Viper 封装，默认使用本地 YAML 文件。Web 服务启动时通常通过 `kid.WithConfigs` 指定配置文件路径。

## 1. 服务启动配置

```go
app := kid.New(
	kid.WithAppName("demo"),
	kid.WithConfigs("./config.yaml"),
	kid.WithRunMode("app.run_mode"),
)
```

`kid.New` 内部会调用 `config.New(config.WithProviders(kid.Options.Configs), config.WithDefault(true))`，并把该配置对象设置为默认配置。

`kid.WithRunMode("app.run_mode")` 表示从配置中读取 `app.run_mode`。示例：

```yaml
app:
  run_mode: release
```

## 2. 多配置文件

可以传入多个配置文件：

```go
app := kid.New(
	kid.WithConfigs("./config.yaml", "./config.local.yaml"),
)
```

`config.Config` 会按 provider 顺序查找配置。`GetString`、`GetInt` 等类型化读取方法会返回第一个可转换为非空字符串的值；`Get` 会返回第一个非 `nil` 值。

命令行启动也可以覆盖配置文件：

```bash
go run . server -c ./config.yaml
```

多个文件可以用逗号分隔，符合 Cobra `StringSlice` 参数行为：

```bash
go run . server -c ./config.yaml,./config.local.yaml
```

## 3. 直接使用配置包

```go
conf := config.New(
	config.WithProviders([]string{"./config.yaml"}),
)

name := conf.GetString("app.name")
port := conf.GetInt("server.port")
enabled := conf.GetBool("feature.enabled")
```

使用默认配置：

```go
name := config.GetString("app.name")
port := config.GetInt("server.port")
enabled := config.GetBool("feature.enabled")
```

常用读取方法：

| 方法 | 返回类型 |
| --- | --- |
| `config.Get` | `interface{}` |
| `config.GetString` | `string` |
| `config.DefaultString` | `string` |
| `config.GetStringSlice` | `[]string` |
| `config.GetStringMap` | `map[string]string` |
| `config.GetInt` | `int` |
| `config.GetInt64` | `int64` |
| `config.GetBool` | `bool` |
| `config.GetFloat` | `float64` |
| `config.GetTime` | `time.Time` |
| `config.GetDuration` | `time.Duration` |
| `config.Exist` | `bool` |

## 4. 常用配置片段

### MySQL

```yaml
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "secret"
  database: "app"
  maxLife: 12h
  maxIdle: 20
  maxOpen: 100
  logLevel: 3
```

`database/mysql` 的 `maxLife` 由 `config.GetDuration` 读取后直接传入 Gorm 连接池，建议写成 `12h` 这类 duration 字符串，并在启动前验证实际值。

### Redis

```yaml
redis:
  addr: "127.0.0.1:6379"
  db: 0
  password: "secret"
  dialTimeout: 3
  readTimeout: 3
  writeTimeout: 3
  maxConnAge: 8
  poolTimeout: 3
  idleTimeout: 3
  poolSize: 50
  minIdleConn: 10
  ctxTimeout: 3
  tls: false
```

Redis 配置样例中的 timeout 字段按当前 options 以数字换算：dial/read/write/pool/context timeout 为秒，idle timeout 为分钟，max connection age 为小时。

### JWT

```yaml
jwt:
  issuer: kid
  refresh_expire: 7
  expire: 1
  signing_method: sha256
  secret: kid_secret
```

JWT 的 `refresh_expire` 和 `expire` 按天配置，源码中会乘以 `24h`。

### Token

```yaml
token:
  expire: 1
  store: redis
  secret: kid_secret
```

`guard/token` 的 `expire` 按天配置。当前 `Provide` 读取 `secret` 和 `expire`，`store` 用于表达业务配置意图，实际 Store 需要通过容器注册。

### 日志

```yaml
logger:
  level: 3
  name: kid
  reportCaller: false
  path: "./logs"
  link: "link/dir/file"
  rotation: 24
  maxAge: 14
```

默认 `kid.New` 会初始化 Logrus 默认 logger。需要日志切割、Hook 或自定义 Formatter 时，可按 `plugin/logrus` 下的 options 和 rotate/hook 包扩展。

## 5. 敏感配置建议

1. 不要把真实密钥、数据库密码、支付证书、云存储密钥提交到 Git。
2. 本地开发可以用 `config.local.yaml` 覆盖基础配置，生产环境建议使用环境注入或配置中心生成最终配置文件。
3. 启动时显式传入配置路径，避免因为工作目录变化读取到错误配置。
4. 对 MySQL、Redis、JWT、支付、短信等外部资源，在部署流水线里增加配置存在性检查。
5. 生产配置变更需要保留审计记录，尤其是 JWT secret、支付证书、短信签名和对象存储权限。

## 6. 配置缺失时的行为

部分基础设施插件在 `Provide` 阶段会检查配置前缀是否存在。例如 MySQL 会检查 `mysql`，Redis 会检查 `redis`，JWT 会检查 `jwt`。缺失时会 `panic` 并输出明确错误，便于在启动阶段暴露问题。

生产环境建议将启动日志接入统一日志系统，并由进程管理器负责失败重启和告警。
