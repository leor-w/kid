# 插件系统

Kid 的插件系统基于一个很小的合同：插件实现 `Provide(context.Context) interface{}`，然后注册到容器中，由 `Populate` 完成依赖注入。

## 1. 插件合同

```go
type Plugin interface {
	Provide(context.Context) interface{}
}
```

注册和注入流程：

```go
app := kid.New(kid.WithConfigs("./config.yaml"))

app.Registry(&redis.Client{})
app.Registry(&mysql.MySQL{})

if err := app.Populate(); err != nil {
	panic(err)
}
```

`Registry` 只把插件提供的实例放入容器。`Populate` 会扫描结构体字段上的 `inject` tag，并把依赖写入对应字段。

## 2. 自定义插件

```go
type DemoPlugin struct {
	Message string
}

func (p *DemoPlugin) Provide(ctx context.Context) interface{} {
	return &DemoPlugin{Message: "hello"}
}
```

注册：

```go
app.Registry(&DemoPlugin{})
```

在其他插件中注入：

```go
type DemoService struct {
	Demo *DemoPlugin `inject:""`
}

func (s *DemoService) Provide(ctx context.Context) interface{} {
	return s
}
```

## 3. 多实例和别名

注册别名：

```go
app.Registry(&mysql.MySQL{}, container.WithAlias("read"))
app.Registry(&mysql.MySQL{}, container.WithAlias("write"))
```

配置可以按别名拆分：

```yaml
mysql:
  read:
    host: "127.0.0.1"
    port: 3306
    user: "reader"
    password: "secret"
    database: "app"
  write:
    host: "127.0.0.1"
    port: 3306
    user: "writer"
    password: "secret"
    database: "app"
```

注入指定别名：

```go
type UserRepository struct {
	ReadDB  *mysql.MySQL `inject:"alias:read"`
	WriteDB *mysql.MySQL `inject:"alias:write"`
}
```

可选依赖：

```go
type OptionalService struct {
	Cache *redis.Client `inject:"opts:NR"`
}
```

`opts:NR` 表示依赖不存在时不返回错误。

## 4. 获取插件实例

```go
raw, err := app.Get(&mysql.MySQL{})
if err != nil {
	panic(err)
}

db := raw.(*mysql.MySQL)
_ = db
```

获取别名实例：

```go
raw, err := app.Get(&mysql.MySQL{}, "read")
if err != nil {
	panic(err)
}

readDB := raw.(*mysql.MySQL)
_ = readDB
```

## 5. 常用插件索引

| 分类 | 包路径 | 说明 | 配置样例 |
| --- | --- | --- | --- |
| 日志 | `plugin/logrus` | Logrus logger、Formatter、日志切割、CloudWatch Hook。 | `plugin/logrus/config.yaml` |
| MySQL | `database/mysql` | Gorm 连接、连接池、Scope、Repository。 | `database/mysql/config.yaml` |
| Redis | `database/redis` | Redis Client、Repository、连接池、TLS。 | `database/redis/config.yaml` |
| JWT | `guard/jwt` | JWT 发行、验证、刷新、吊销、黑名单。 | `guard/jwt/config.yaml` |
| Token | `guard/token` | HMAC 签名 Token。 | `guard/token/config.yaml` |
| 验证码 | `plugin/captcha` | 数字、字符、音频等验证码。 | `plugin/captcha/config.yaml` |
| 支付 | `plugin/alipay`、`plugin/wechat/payments`、`plugin/stripe`、`plugin/apple`、`plugin/google/pay`、`plugin/douyin` | 支付、订阅、回调和订单相关能力。 | 各插件目录下 `config.yaml` |
| OAuth | `plugin/oauth2/*`、`plugin/qq` | Apple、Facebook、Google、Instagram、TikTok、QQ 登录。 | 各插件目录下 `config.yaml` 或 `config.yml` |
| 短信 | `plugin/sms/alicloud`、`plugin/sms/twilio` | 阿里云短信、Twilio 验证。 | 各插件目录下 `config.yaml` |
| 对象存储 | `plugin/oss`、`plugin/qiniu`、`plugin/aws/awss3` | 阿里云 OSS、七牛、AWS S3。 | 各插件目录下 `config.yaml` |
| 消息队列 | `plugin/mq/rabbitmq` | RabbitMQ 生产者和消费者。 | `plugin/mq/rabbitmq/config.yaml` |
| 延迟队列 | `plugin/delayed` | Redis 延迟任务、生产者、消费者、worker。 | `plugin/delayed/config.yaml` |
| 分布式锁 | `plugin/lock` | Redis 分布式锁和 redsync。 | `plugin/lock/options.go` |
| WebSocket | `plugin/websocket` | WebSocket Server、Connection、Emitter。 | `plugin/websocket/config.yaml` |
| 定时任务 | `plugin/cron` | Cron 任务。 | 源码 options |
| 限流 | `plugin/rate_limit` | 请求限流配置。 | `plugin/rate_limit/config.yaml` |
| RBAC | `plugin/rbac` | Casbin 权限控制。 | `plugin/rbac/config.yaml` |
| 机器人 | `plugin/robot/dingtalk` | 钉钉机器人和 webhook。 | `plugin/robot/dingtalk/config.yaml` |
| 加密 | `plugin/crypt` | RSA、Ed25519。 | `plugin/crypt/config.yaml` |
| 搜索 | `plugin/elasticsearch` | Elasticsearch Client。 | `plugin/elasticsearch/elasticsearch.yaml` |
| 工作者池 | `plugin/workerpool` | 固定 worker 数任务池。 | 源码 options |

## 6. 插件配置入口

插件配置样例分散在插件目录内。接入某个插件前，优先查看同目录的 `config.yaml` 或 `config.yml`。

示例：

```text
plugin/stripe/config.yaml
plugin/wechat/config.yaml
plugin/sms/alicloud/config.yaml
plugin/oauth2/google/config.yaml
plugin/websocket/config.yaml
```

生产环境不要直接使用样例密钥。支付证书、OAuth secret、短信密钥和对象存储密钥应通过安全配置系统注入。

## 7. 接入建议

1. 只注册业务实际使用的插件，避免启动时初始化不需要的外部资源。
2. 注册顺序建议先注册基础资源，再注册依赖基础资源的插件，例如先 Redis，再 Redis Store，再 JWT。
3. 每次接入插件时先确认配置前缀是否存在。多个插件在 `Provide` 阶段会因为配置缺失直接 `panic`。
4. 支付、OAuth、短信、云存储等平台类插件要按平台官方规则补充业务级验签、回调幂等和审计日志。
5. 多实例插件必须使用别名，避免容器无法区分同类型实例。
