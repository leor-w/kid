# 认证与授权

Kid 的认证入口是 `guard.Guard`。当前仓库提供两类实现：JWT 和签名 Token。两者都依赖 Store 保存 token 详情和用户 token 列表。

## 1. 统一 Guard 接口

```go
type Guard interface {
	License(user *guard.User, newRefresh bool) (string, string, error)
	RefreshLicense(refreshLicense string) (string, error)
	GetLicense(guard.UserType, int64) ([]string, error)
	Verify(string) (*guard.User, error)
	Cancellation(string) error
	CancellationAll(guard.UserType, int64) error
	ExpiresAt(string) int64
	IssuerAt(string) int64
}
```

用户信息结构：

```go
type User struct {
	Uid       int64
	Type      guard.UserType
	TokenType guard.TokenType
}
```

内置用户类型：

| 常量 | 说明 |
| --- | --- |
| `guard.Admin` | 后台管理员。 |
| `guard.General` | 普通用户。 |

## 2. JWT

JWT 配置：

```yaml
jwt:
  issuer: kid
  refresh_expire: 7
  expire: 1
  signing_method: sha256
  secret: kid_secret
```

JWT 特点：

1. 支持 access token 和 refresh token。
2. 支持 `RefreshLicense` 刷新 access token。
3. 支持 Redis Store 保存 token 详情。
4. 支持 Redis blacklist 吊销 JWT。
5. `Verify` 会先检查黑名单和 Store，再解析 token。

`jwt.New` 可以创建 JWT 实例，但 `store` 和 `blacklist` 字段需要通过容器注入后才能安全使用，因此业务服务推荐使用 `Registry` 和 `Populate` 完成装配。

依赖注册：

```go
app := kid.New(kid.WithConfigs("./config.yaml"))

app.Registry(&redis.Client{})
app.Registry(&store.RedisStore{})
app.Registry(&jwt.RedisBlacklist{})
app.Registry(&jwt.Jwt{})

if err := app.Populate(); err != nil {
	panic(err)
}
```

获取 JWT 实例：

```go
rawGuard, err := app.Get(&jwt.Jwt{})
if err != nil {
	panic(err)
}

guarder := rawGuard.(*jwt.Jwt)
```

发行凭证：

```go
accessToken, refreshToken, err := guarder.License(&guard.User{
	Uid:  10001,
	Type: guard.General,
}, true)
if err != nil {
	panic(err)
}

_ = accessToken
_ = refreshToken
```

验证凭证：

```go
user, err := guarder.Verify(accessToken)
if err != nil {
	return gin.H{"error": err.Error()}
}

return gin.H{"uid": user.Uid}
```

刷新凭证：

```go
newAccessToken, err := guarder.RefreshLicense(refreshToken)
if err != nil {
	return gin.H{"error": err.Error()}
}

return gin.H{"access_token": newAccessToken}
```

吊销凭证：

```go
if err := guarder.Cancellation(accessToken); err != nil {
	return gin.H{"error": err.Error()}
}

return gin.H{"revoked": true}
```

## 3. 签名 Token

Token 配置：

```yaml
token:
  expire: 1
  store: redis
  secret: kid_secret
```

签名 Token 特点：

1. 通过 HMAC 签名生成 token。
2. token 详情仍保存在 Store 中。
3. 当前 `RefreshLicense` 方法未实现刷新逻辑，会返回空字符串和 `nil` 错误。
4. 适合不需要 JWT 标准 claims 和 refresh token 的内部系统。

`token.New` 可以创建 Token 实例，但 `store` 字段需要通过容器注入后才能安全使用，因此业务服务推荐使用 `Registry` 和 `Populate` 完成装配。

依赖注册：

```go
app.Registry(&redis.Client{})
app.Registry(&store.RedisStore{})
app.Registry(&token.Token{})

if err := app.Populate(); err != nil {
	panic(err)
}
```

获取 Token 实例：

```go
rawGuard, err := app.Get(&token.Token{})
if err != nil {
	panic(err)
}

guarder := rawGuard.(*token.Token)
```

发行凭证：

```go
license, _, err := guarder.License(&guard.User{
	Uid:  10001,
	Type: guard.General,
})
if err != nil {
	panic(err)
}
```

验证凭证：

```go
user, err := guarder.Verify(license)
if err != nil {
	return gin.H{"error": err.Error()}
}

return gin.H{"uid": user.Uid}
```

## 4. Redis Store

`guard/store.RedisStore` 负责保存：

1. token 详情。
2. 用户 token 列表。
3. token 过期和批量过期。

Redis key 由 `guard/constant` 生成，按用户类型、用户 ID 和 token 组合索引。Store 使用 Redis pipeline 同时写入用户 token 索引和 token 详情。

## 5. HTTP 中间件示例

```go
func Auth(guarder guard.Guard) kid.Middleware {
	return func(ctx *kid.Context) {
		license := ctx.GetHeader("Authorization")
		if license == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "missing authorization"})
			return
		}

		user, err := guarder.Verify(license)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		ctx.BindUser(user)
		ctx.Next()
	}
}
```

读取绑定用户：

```go
app.GET("/me", func(ctx *kid.Context) any {
	var user guard.User
	if err := ctx.GetBindUser(&user); err != nil {
		return gin.H{"error": err.Error()}
	}

	return gin.H{"uid": user.Uid, "type": user.Type}
}, Auth(guarder))
```

## 6. 安全建议

1. JWT secret 和 Token secret 必须使用生产级随机字符串，不要使用示例值。
2. access token 生命周期建议短，refresh token 生命周期可以稍长，但要支持吊销。
3. Redis Store 是 JWT/Token 有效性判断的一部分，Redis 不可用会影响认证。
4. `CancellationAll` 可用于用户修改密码、账号冻结或后台强制下线。
5. Header 中的 token 格式由业务决定，若使用 `Bearer <token>`，中间件需要自行剥离前缀。
6. 登录、刷新、吊销失败必须记录清晰日志，便于定位 Redis、签名、过期和黑名单问题。
