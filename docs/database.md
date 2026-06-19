# 数据库与缓存

Kid 提供 MySQL、Redis 和仓储接口封装。MySQL 基于 Gorm，Redis 基于 `go-redis/redis/v8`。

## 1. MySQL 初始化

直接创建连接：

```go
db := mysql.New(
	mysql.WithHost("127.0.0.1"),
	mysql.WithPort(3306),
	mysql.WithUser("root"),
	mysql.WithPassword("secret"),
	mysql.WithDb("app"),
	mysql.WithMaxIdle(20),
	mysql.WithMaxOpen(100),
	mysql.WithMaxLife(12*time.Hour),
	mysql.WithLogLevel(3),
)
```

`mysql.New` 会：

1. 使用 Gorm 打开 MySQL 连接。
2. 配置 Kid Gorm logger。
3. 设置最大空闲连接数、最大连接数和连接最大存活时间。
4. 使用 `PingContext` 测试连接状态。

连接失败会 `panic` 并输出明确错误，例如开启连接失败、创建连接对象失败或测试连接状态失败。

## 2. MySQL 配置注入

配置示例：

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

注册到 Kid 容器：

```go
app := kid.New(kid.WithConfigs("./config.yaml"))

app.Registry(&mysql.MySQL{})

if err := app.Populate(); err != nil {
	panic(err)
}
```

`mysql.MySQL.Provide` 会读取 `mysql` 配置前缀。多实例可以配合别名使用，例如 `mysql.read`、`mysql.write`，注册时传入容器别名。

## 3. Gorm Scope

分页：

```go
db.Scopes(mysql.Paginate(1, 20)).Find(&users)
```

条件：

```go
db.Scopes(mysql.Where(mysql.Wrappers(
	mysql.NewWhereWrapper("status", mysql.ConditionTypeEq, 1),
	mysql.NewWhereWrapper("id", mysql.ConditionTypeIn, []int64{1, 2, 3}),
))).Find(&users)
```

模糊搜索：

```go
db.Scopes(mysql.Search(mysql.Wrappers(
	mysql.NewWhereWrapper("name", mysql.ConditionTypeEq, "kid"),
))).Find(&users)
```

软删除过滤：

```go
db.Scopes(mysql.FilterDeleted()).Find(&users)
```

`Paginate` 会保护分页参数：页码小于等于 0 时使用 1，页大小小于等于 0 时使用 10，页大小大于等于 100 时限制为 100。

## 4. 自动迁移

实现模型提供者：

```go
type MigrateModels struct{}

func (MigrateModels) Models() []interface{} {
	return []interface{}{
		&User{},
	}
}
```

启动入口传入：

```go
app := kid.New(
	kid.WithConfigs("./config.yaml"),
	kid.WithDbMigrate(MigrateModels{}),
)
```

执行迁移：

```bash
go run . database --automigrate -c ./config.yaml
```

如需关闭外键检查：

```bash
go run . database --automigrate -c ./config.yaml --closefkcheck
```

如果没有配置有效的模型提供者，迁移命令会输出错误并退出。

## 5. MySQL Repository

`database/mysql.Repository` 通过容器注入 `*mysql.MySQL`，实现了 `database/repos` 中的基础仓储接口。

主要能力：

| 方法 | 说明 |
| --- | --- |
| `Transaction` | 开启事务并在函数返回错误时回滚。 |
| `ExecWithTx` | 如果上下文已有事务则复用，否则新开事务。 |
| `Exist` | 根据 Finder 判断记录是否存在。 |
| `GetOne`、`Find`、`GetById`、`GetByKV` | 查询单条或多条记录。 |
| `Create`、`Save` | 创建或保存模型。 |
| `Update`、`UpdateByKV` | 更新模型。 |
| `Delete` | 删除模型。 |
| `Count`、`Sum` | 聚合查询。 |
| `CreateTable`、`TableHasExist` | 表结构辅助方法。 |

注册示例：

```go
app.Registry(&mysql.MySQL{})
app.Registry(&mysql.Repository{})

if err := app.Populate(); err != nil {
	panic(err)
}
```

## 6. Redis 初始化

直接创建连接：

```go
rdb := redis.New(
	redis.WithAddr("127.0.0.1:6379"),
	redis.WithDb(0),
	redis.WithPassword("secret"),
	redis.WithPoolSize(50),
	redis.WithMinIdle(10),
)
```

`redis.New` 会创建 `go-redis` Client，并在初始化阶段执行 `Ping`。连接失败会 `panic`，错误信息包含 `ping redis failed`。

配置示例：

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

注册到容器：

```go
app.Registry(&redis.Client{})

if err := app.Populate(); err != nil {
	panic(err)
}
```

## 7. Redis Repository

`database/redis.RedisRepository` 通过容器注入 `*redis.Client`，提供常用 key-value 操作。

| 方法 | 说明 |
| --- | --- |
| `Set` | 设置 key，过期时间单位为秒。 |
| `Get` | 读取字符串值。 |
| `Del` | 删除一个或多个 key。 |
| `Expire` | 更新过期时间，单位为秒。 |
| `Exists` | 判断 key 是否存在。 |

注册示例：

```go
app.Registry(&redis.Client{})
app.Registry(&redis.RedisRepository{})

if err := app.Populate(); err != nil {
	panic(err)
}
```

## 8. 性能建议

1. MySQL 根据业务并发设置 `maxIdle`、`maxOpen` 和 `maxLife`，避免连接池过小导致请求堆积，也避免连接过多压垮数据库。
2. Redis 根据 QPS 和命令耗时设置 `poolSize`、`minIdleConn` 和 timeout。高并发场景优先使用 pipeline 或批量命令。
3. 批量迁移、批量写入和外部 API 同步任务不要直接在请求 goroutine 中执行，优先使用队列、worker pool 或异步任务。
4. 生产环境不要裸跑服务，应由进程管理器接管重启、日志和告警。
