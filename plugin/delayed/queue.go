package delayed

import (
	"context"
)

type IProduct interface {
	Dispatch(topic string, payload interface{}) (string, error)                                      // 添加任务
	DispatchDelay(topic string, payload interface{}, delay int, delayUnit DelayUnit) (string, error) // 添加延迟任务
}

type IConsumer interface {
	Run() error                  // 启动消费者
	JoinWorker(worker *Worker)   // 加入一个 worker
	Close(context.Context) error // 关闭消费者
}

// Queue 代表一个队列, 用于存储 Job, 并提供 Job 的增删改查等操作, 以及 Job 的投递和消费, 以及 Job 的迁移, 以及队列的关闭等操作
// 延迟 Job 会被存储在延迟队列中, 就绪 Job 会被存储在就绪队列中, 就绪 Job 会被消费, 延迟 Job 在延迟时间到达后会被迁移到就绪队列中
type Queue interface {
	Run() error                           // 启动队列
	AddTopic(string) error                // 添加一个 topic
	Push(*Job) error                      // 添加 Job 到队列中, 投递的方向为队列尾部
	Pop(string) (*Job, error)             // 从队列中获取 Job
	PushDelay(*Job) error                 // 将 Job 添加到延迟队列中
	Cancel(string, string) error          // 通过指定 topic 和 jobId 取消延迟任务
	FetchReadyJob(string) ([]*Job, error) // 获取队列中已就绪的 Job
	FetchDelayJob(string) ([]*Job, error) // 获取延迟队列中的所有 Job
	ReadyLen(string) (int, error)         // 获取队列中已就绪的 Job 数量
	DelayLen(string) (int, error)         // 获取延迟队列中的 Job 数量
	Exists(string, string) (bool, error)  // 判断 JobID 是否存在
	Close() error                         // 关闭队列
}

type JobHandler interface {
	Topic() string                          // 获取 Job 的 topic
	OnHandle(context.Context, []byte) error // 处理 Job
	OnError(context.Context, []byte) error  // 处理 Job 处理失败的情况
}
