package delayed

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/utils"

	redisv8 "github.com/go-redis/redis/v8"

	"github.com/lsjhtang/kid/logger"
)

// Consumer 消费者
type Consumer struct {
	options       *ConsumerOptions
	topicWorkers  []*Worker
	queue         Queue `inject:""`
	stopCh        chan struct{}
	processingNum uint32
	status        uint32
	ctx           context.Context

	sync.Mutex // 互斥锁
}

type ConsumerOption func(*ConsumerOptions)

func (c *Consumer) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("delayed.consumer%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return NewConsumer(
		WithSleep(config.GetInt(utils.GetConfigurationItem(confPrefix, "sleep"))),
		WithClosedWait(config.GetInt(utils.GetConfigurationItem(confPrefix, "closeWait"))),
		WithMaxRetries(config.GetInt(utils.GetConfigurationItem(confPrefix, "maxRetires"))),
		WithWorkerNum(config.GetInt(utils.GetConfigurationItem(confPrefix, "workerNum"))),
	)
}

const (
	consumerStatusReady   = iota // 服务准备就绪
	consumerStatusRunning        // 服务运行中
	consumerStatusClosed         // 服务关闭
)

// Run 启动队列
func (c *Consumer) Run() error {
	if c.status != consumerStatusReady {
		return errors.New("定时任务消费者已经启动")
	}
	go func() {
		// 启动队列 Worker 处理任务
		for _, worker := range c.topicWorkers {
			go func(worker *Worker) {
				if err := c.processJob(c.ctx, worker); err != nil {
					logger.Errorf("定时任务队列 JoinWorker: 定时任务执行失败, topic: %s, err: %s", worker.Topic, err.Error())
				}
			}(worker)
		}
		// 监听系统退出信号
		go c.watchSystemSignal()
		if err := c.queue.Run(); err != nil {
			panic(fmt.Errorf("定时任务队列 Run: %w", err))
		}
		atomic.StoreUint32(&c.status, consumerStatusRunning)
		<-c.stopCh
		c.Close()
	}()
	return nil
}

func (c *Consumer) Register(worker *Worker) {
	c.Lock()
	defer c.Unlock()
	c.queue.AddTopic(worker.Topic)
	c.topicWorkers = append(c.topicWorkers, worker)
}

// processJob 处理 job
func (c *Consumer) processJob(ctx context.Context, worker *Worker) error {
	for {
		// 如果关闭则退出循环
		if atomic.LoadUint32(&c.status) == consumerStatusClosed {
			break
		}
		// 通过 workerPool 控制并发数
		if err := worker.WorkerPool.Acquire(ctx, 1); err != nil {
			return err
		}
		// 尝试从队列中获取 job
		job, err := c.queue.Pop(worker.Topic)
		if err != nil && !errors.Is(err, redisv8.Nil) {
			logger.Errorf("定时任务队列: 定时任务获取失败, topic: %s, err: %s", worker.Topic, err.Error())
			// 获取失败, 释放 workerPool
			worker.WorkerPool.Release(1)
			time.Sleep(time.Second * time.Duration(c.options.Sleep))
			continue
		}
		if job == nil {
			logger.Debug("定时任务队列: 定时任务队列为空")
			// 队列为空, 释放 workerPool
			worker.WorkerPool.Release(1)
			time.Sleep(time.Second * time.Duration(c.options.Sleep))
			continue
		}
		go func() {
			c.Lock()
			atomic.AddUint32(&c.processingNum, 1)
			c.Unlock()
			defer func() {
				c.Lock()
				atomic.AddUint32(&c.processingNum, ^uint32(0))
				c.Unlock()
			}()
			// 处理 job
			if err := worker.Handler.OnHandle(ctx, job.Payload); err != nil {
				// 处理失败
				// 判断是否需要重试
				job.Retries++
				if job.Retries >= c.options.MaxRetries {
					// 重试次数超过最大重试次数, 丢弃
					if err := worker.Handler.OnError(ctx, job.Payload); err != nil {
						// 丢弃失败, 打印日志
						logger.Errorf("定时任务队列: 定时任务执行失败, topic: %s, payload: %v, err: %s",
							job.Topic, job.Payload, err.Error())
					}
					worker.WorkerPool.Release(1)
					return
				}
				// 重新入队
				if err := c.queue.Push(job); err != nil {
					// 入队失败, 丢弃
					if err := worker.Handler.OnError(ctx, job.Payload); err != nil {
						// 丢弃失败, 打印日志
						logger.Errorf("定时任务队列: 定时任务执行失败, topic: %s, payload: %v, err: %s",
							job.Topic, job.Payload, err.Error())
					}
				}
			}
			worker.WorkerPool.Release(1)
		}()
	}
	// 服务关闭, 等待所有 worker 退出
	if err := worker.WorkerPool.Acquire(ctx, 1); err != nil {
		return err
	}
	return nil
}

// watchSystemSignal 监听系统退出信号
func (c *Consumer) watchSystemSignal() {
	signals := make(chan os.Signal, 1)
	// 监听所有 Linux 常见的系统退出信号
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT,
		syscall.SIGBUS, syscall.SIGFPE, syscall.SIGSEGV, syscall.SIGPIPE,
		syscall.SIGALRM, syscall.SIGTERM)
	for {
		select {
		case <-signals: // 收到系统信号
			logger.Infof("定时任务消费者: 收到系统退出信号, 服务即将关闭")
			c.stopCh <- struct{}{}
			break
		}
	}
}

// Close 关闭队列
func (c *Consumer) Close() error {
	atomic.StoreUint32(&c.status, consumerStatusClosed)
	for c.processingNum > 0 {
		fmt.Println("定时任务队列: 等待 %d 个任务处理完成\n", c.processingNum)
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}

func (c *Consumer) Cancel(topic, jobId string) error {
	if err := c.queue.Cancel(topic, jobId); err != nil {
		return fmt.Errorf("定时任务队列 Cancel: %w", err)
	}
	return nil
}

func NewConsumer(opts ...ConsumerOption) *Consumer {
	options := &ConsumerOptions{
		Sleep:      3,
		ClosedWait: 5,
		MaxRetries: 5,
		WorkerNum:  20,
	}
	for _, opt := range opts {
		opt(options)
	}
	return &Consumer{
		options: options,
		stopCh:  make(chan struct{}, 1),
		status:  consumerStatusReady,
		Mutex:   sync.Mutex{},
	}
}
