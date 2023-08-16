package delayed

// ConsumerOptions 代表消费者的配置选项
type ConsumerOptions struct {
	Sleep      int // 当就绪队列为空时, 休眠时间, 单位为秒
	ClosedWait int // 关闭队列时, 等待队列中的任务执行完成的时间, 单位为秒
	MaxRetries int // 任务失败后的重试次数
	WorkerNum  int // 任务处理器的数量
}

func WithClosedWait(closedWait int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ClosedWait = closedWait
	}
}

func WithSleep(sleep int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.Sleep = sleep
	}
}

func WithMaxRetries(maxRetries int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.MaxRetries = maxRetries
	}
}

func WithWorkerNum(workerNum int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.WorkerNum = workerNum
	}
}
