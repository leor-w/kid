package injector

// injectTag 依赖注入标签
// 值应为 `inject:"别名,选项"`
// 例如 DB *mysql.MysqlConn `inject:"alias:db1,opts:unessential|xxx"`
// 上面的例子中, DB 为实体的字段名, mysql.MysqlConn 为实体的类型, db1 为实体的别名,
// unessential 为标注该字段在执行注入时, 如果没有找到对应的实例, 则忽略该字段
// 如果不标注则在执行注入时, 没有找到对应的实例, 会返回错误
const injectTag = "inject" // 依赖注入标签

const (
	tagAliasKey    = "alias"    // 别名
	tagScopeKey    = "scope"    // 分组
	tagOptionalKey = "optional" // 选项
)

type NameKey struct{}

// TagOption 标签选项

const (
	TagOptionNotRequired = "NR" // Not Required, 标注该字段在执行注入时, 如果没有找到对应的实例, 则忽略该字段
)
