# 运维排障

本文说明 Kid 服务的启动流程、日志、Recovery、常见故障和性能建议。

## 1. 启动流程

`kid.New` 会执行以下流程：

1. 处理 Cobra 命令行。
2. 加载配置文件。
3. 初始化日志。
4. 从配置 key 读取 Gin 运行模式。
5. 创建 Gin Engine。
6. 注册默认请求日志中间件。
7. 注册 Recovery。
8. 初始化 RouterGroup。

`Launch` 会打印 Kid 标识并启动 HTTP 服务：

```go
app.Launch(":8080")
```

如果不传 host，默认监听 `:8080`。

## 2. 生产启动建议

生产环境禁止裸跑服务。建议使用 Docker、systemd、Kubernetes 或其他进程管理方式托管进程。

systemd 示例：

```ini
[Unit]
Description=Kid Demo Service
After=network.target

[Service]
WorkingDirectory=/opt/kid-demo
ExecStart=/opt/kid-demo/server server -c /etc/kid-demo/config.yaml
Restart=always
RestartSec=3
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

容器启动示例：

```bash
docker run --rm \
  -p 8080:8080 \
  -v "$PWD/config.yaml:/app/config.yaml:ro" \
  kid-demo:latest \
  server -c /app/config.yaml
```

## 3. 日志

默认启动会初始化 Logrus logger。请求日志由 `middleware.Logger` 输出，格式包含：

- HTTP 方法
- 客户端 IP
- 路径
- 协议
- 状态码
- 请求耗时
- User-Agent
- Gin 私有错误

状态码为 200 时输出 info 日志，否则输出 error 日志。

日志配置示例：

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

`plugin/logrus` 支持 Formatter、日志切割和 CloudWatch Hook。日志目录必须确保服务进程有写权限。

## 4. Recovery

默认 Recovery 会捕获 panic，记录错误和堆栈，并返回：

```json
{"code":50000,"message":"内部错误"}
```

排障时优先搜索日志中的：

```text
panic:
堆栈:
```

如果 panic 发生在插件初始化阶段，服务可能无法启动；如果发生在请求处理阶段，Recovery 会返回 500 响应并记录堆栈。

## 5. 常见故障

| 现象 | 常见原因 | 排查动作 |
| --- | --- | --- |
| 启动后找不到配置 | `server -c` 路径错误，或工作目录不符合预期。 | 使用绝对路径启动，确认配置文件存在。 |
| `config.yaml file not found configuration item` | MySQL、Redis、Token 等配置前缀缺失。 | 检查 `mysql`、`redis`、`token`、`jwt` 等配置块。 |
| `jwt 注入失败` | 缺少 `jwt` 配置前缀。 | 补齐 `jwt.issuer`、`jwt.secret`、`jwt.expire` 等字段。 |
| `开启MySQL连接错误` | 地址、账号、密码、库名或网络错误。 | 使用 MySQL 客户端验证连接，检查安全组和 DNS。 |
| `测试MySQL连接状态错误` | Gorm 打开连接后 Ping 失败。 | 检查数据库可用性、账号权限和连接数限制。 |
| `ping redis failed` | Redis 地址、密码、DB、TLS 或网络错误。 | 使用 `redis-cli` 验证连接，检查 Redis TLS 配置。 |
| `container.populate` 未找到实例 | 插件未注册或别名不一致。 | 检查 `Registry` 顺序、`inject` tag 和 `container.WithAlias`。 |
| 请求返回 500 且日志有 `panic` | Handler 或插件发生 panic。 | 根据堆栈定位代码行，补充错误处理和输入校验。 |

## 6. 健康检查

建议业务服务显式增加健康检查接口：

```go
app.GET("/healthz", func(ctx *kid.Context) any {
	return gin.H{"status": "ok"}
})
```

如果服务依赖 MySQL、Redis 或外部支付网关，建议另建 readiness 接口检查关键依赖，避免把依赖未就绪的实例加入流量。

## 7. 性能建议

1. MySQL 和 Redis 连接池必须按真实并发配置，不要使用未经评估的默认值直接上线。
2. 批量任务、延迟任务、Webhook 重试、短信批发和外部 API 同步优先使用队列、`plugin/workerpool` 或 `utils.SafeExecutor` 控制并发。
3. 需要请求级限流时优先评估 `plugin/rate_limit`，避免单实例被突发流量压垮。
4. Redis 高频操作优先使用 pipeline 或批量命令，减少网络往返。
5. 外部平台回调必须做幂等处理，避免重试导致重复支付、重复发货或重复发券。
6. 大文件上传、转码、压缩和云存储同步不要在 HTTP 请求链路内长时间阻塞。

## 8. 安全建议

1. JWT secret、支付证书、OAuth secret、短信密钥、对象存储密钥不要入库。
2. 生产日志不要输出完整 token、密码、私钥、支付证书或用户敏感信息。
3. 对管理接口、迁移命令和内部回调接口增加网络隔离或鉴权。
4. 命令行迁移前确认目标数据库，避免误连生产库。
5. panic 日志包含堆栈，日志系统需要有访问控制。

## 9. 发布前检查

发布前建议至少检查：

```bash
go test ./config ./container ./guard/token ./utils
git diff --check
```

全仓库测试可能触发大量第三方插件依赖下载或外部服务初始化。若要把全仓库测试作为门禁，需要先为依赖外部服务的包补隔离测试或 mock。
