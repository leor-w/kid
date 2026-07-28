package injector

var (
	NotFoundBucketError        = "未找到对应的 bucket: [%v]"
	NonPointerError            = "无法为非指针类型注入: [%v]"
	NotFoundEntityError        = "[%v] 未找到对应的实例: [%v]; 注入范围: %s; 别名: [%s]"
	InjectionUnfinishedError   = "注入未完成: [%v]"
	InvalidInjectionFiledError = "类型 [%v], 别名: [%s] 值无效"
)
