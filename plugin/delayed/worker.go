package delayed

import "golang.org/x/sync/semaphore"

// Worker 代表一个 topic 的 worker
type Worker struct {
	Topic       string              // 任务 topic
	Handler     JobHandler          // 任务处理器
	WorkerCount int                 // 并行任务数
	WorkerPool  *semaphore.Weighted // 通过信号量控制并发协程数
}

// NewWorker 用于创建一个 Worker
func NewWorker(topic string, handler JobHandler, workerCount int) *Worker {
	return &Worker{
		Topic:       topic,
		Handler:     handler,
		WorkerCount: workerCount,
		WorkerPool:  semaphore.NewWeighted(int64(workerCount)),
	}
}
