package utils

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/lsjhtang/kid/logger"
	"github.com/sirupsen/logrus"
)

// SafeGo 提供安全启动协程的能力，捕获 panic 并记录日志, 避免程序崩溃
func SafeGo(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("安全协程捕获异常: %v\n", err)
				if debug.Stack() != nil {
					stack := debug.Stack()
					logger.Errorf("堆栈信息:\n%s", stack)
				}
			}
		}()
		fn()
	}()
}

// JobFunc 定义执行的任务函数类型，接受上下文并返回错误
type JobFunc func(ctx context.Context) error

// SafeExecutor 提供安全执行任务的能力
type SafeExecutor struct {
	logger      *logrus.Logger
	maxWorkers  int
	jobQueue    chan JobFunc
	errorChan   chan error
	wg          sync.WaitGroup
	shutdownCtx context.Context
	cancel      context.CancelFunc
}

// NewSafeExecutor 创建一个新的 SafeExecutor 实例
func NewSafeExecutor(logger *logrus.Logger, maxWorkers int) *SafeExecutor {
	ctx, cancel := context.WithCancel(context.Background())
	executor := &SafeExecutor{
		logger:      logger,
		maxWorkers:  maxWorkers,
		jobQueue:    make(chan JobFunc),
		errorChan:   make(chan error),
		shutdownCtx: ctx,
		cancel:      cancel,
	}
	executor.start()
	return executor
}

// start 启动工作池中的 worker
func (e *SafeExecutor) start() {
	for i := 0; i < e.maxWorkers; i++ {
		e.wg.Add(1)
		go e.worker(i)
	}
}

// worker 是工作池中的单个 worker
func (e *SafeExecutor) worker(id int) {
	defer e.wg.Done()
	e.logger.Infof("Worker %d 启动", id)
	for {
		select {
		case <-e.shutdownCtx.Done():
			e.logger.Infof("Worker %d 接收到关闭信号", id)
			return
		case job, ok := <-e.jobQueue:
			if !ok {
				e.logger.Infof("Worker %d 任务队列关闭", id)
				return
			}
			e.logger.Infof("Worker %d 开始执行任务", id)
			if err := e.executeJob(job); err != nil {
				e.logger.Errorf("Worker %d 任务执行出错: %v", id, err)
				select {
				case e.errorChan <- err:
				default:
					e.logger.Warn("错误通道已满，丢弃错误")
				}
			} else {
				e.logger.Infof("Worker %d 任务执行成功", id)
			}
		}
	}
}

// executeJob 安全执行一个任务，捕获 panic 并返回错误
func (e *SafeExecutor) executeJob(job JobFunc) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("任务 panic: %v", r)
			e.logger.Errorf("任务 panic: %v", r)
		}
	}()
	err = job(e.shutdownCtx)
	return
}

// Submit 提交一个任务到执行器
func (e *SafeExecutor) Submit(job JobFunc) error {
	select {
	case <-e.shutdownCtx.Done():
		return errors.New("执行器已关闭，无法提交新任务")
	default:
	}

	select {
	case e.jobQueue <- job:
		return nil
	default:
		return errors.New("任务队列已满")
	}
}

// Errors 返回一个只读的错误通道，用于接收任务执行中的错误
func (e *SafeExecutor) Errors() <-chan error {
	return e.errorChan
}

// Shutdown 优雅关闭执行器，等待所有 worker 完成
func (e *SafeExecutor) Shutdown(ctx context.Context) error {
	e.cancel()
	close(e.jobQueue)
	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		close(e.errorChan)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

//func main() {
//	// 初始化日志
//	logger := logrus.New()
//	logger.SetFormatter(&logrus.TextFormatter{
//		FullTimestamp: true,
//	})
//
//	// 创建 SafeExecutor，限制最大并发数为5
//	executor := NewSafeExecutor(logger, 5)
//
//	// 监听错误
//	go func() {
//		for err := range executor.Errors() {
//			logger.Errorf("接收到任务错误: %v", err)
//		}
//	}()
//
//	// 提交一些任务
//	for i := 0; i < 10; i++ {
//		index := i
//		err := executor.Submit(func(ctx context.Context) error {
//			logger.Infof("执行任务 %d", index)
//			if index == 5 {
//				panic("模拟 panic")
//			}
//			time.Sleep(time.Millisecond * 500)
//			logger.Infof("任务 %d 完成", index)
//			return nil
//		})
//		if err != nil {
//			logger.Errorf("无法提交任务 %d: %v", index, err)
//		}
//	}
//
//	// 等待一段时间后关闭执行器
//	time.Sleep(time.Second * 3)
//	if err := executor.Shutdown(context.Background()); err != nil {
//		logger.Errorf("关闭执行器时出错: %v", err)
//	} else {
//		logger.Info("执行器已成功关闭")
//	}
//}
