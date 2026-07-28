package workerpool

import (
	"context"

	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"
)

type Job interface {
	Do()
}

type Worker struct {
	JobQueue chan Job
	Quit     chan struct{}
}

func NewWorker() *Worker {
	return &Worker{
		JobQueue: make(chan Job),
		Quit:     make(chan struct{}),
	}
}

func (worker *Worker) Run(wq chan chan Job) {
	go func() {
		for {
			wq <- worker.JobQueue
			select {
			case job := <-worker.JobQueue:
				if err := utils.NoPanic(func() {
					job.Do()
				}); err != nil {
					logger.Errorf("任务执行失败: %s", err.Error())
				}
			case <-worker.Quit:
				return
			}
		}
	}()
}

type Pool struct {
	workerLen   int
	jobQueue    chan Job
	workerQueue chan chan Job
}

func (pool *Pool) Provide(ctx context.Context) interface{} {
	poolWorkerLen := ctx.Value("workerLen").(int)
	if poolWorkerLen <= 0 {
		poolWorkerLen = 20
	}
	p := NewPool(poolWorkerLen)
	p.Run()
	return p
}

func NewPool(workerLen int) *Pool {
	return &Pool{
		workerLen:   workerLen,
		jobQueue:    make(chan Job),
		workerQueue: make(chan chan Job),
	}
}

func (pool *Pool) Join(job Job) {
	pool.jobQueue <- job
}

func (pool *Pool) Run() {
	for i := 0; i < pool.workerLen; i++ {
		worker := NewWorker()
		worker.Run(pool.workerQueue)
	}

	go func() {
		for {
			select {
			case job := <-pool.jobQueue:
				worker := <-pool.workerQueue
				worker <- job
			}
		}
	}()
}
