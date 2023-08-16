package delayed

// Job 代表一个定时任务
type Job struct {
	Id        string `msgpack:"1"` // 任务ID
	Topic     string `msgpack:"2"` // 任务主题
	Delay     int64  `msgpack:"3"` // 任务延迟执行的时间戳
	Payload   []byte `msgpack:"4"` // 任务内容
	Timestamp int64  `msgpack:"5"` // 任务投递时间
	Retries   int    `msgpack:"6"` // 任务重试次数
}
