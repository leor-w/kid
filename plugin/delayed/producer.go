package delayed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Producer struct {
	queue Queue `inject:""`
}

func (p *Producer) Provide(context.Context) interface{} {
	return p
}

// Dispatch 向就绪队列投递任务
func (p *Producer) Dispatch(topic string, payload []byte) (string, error) {
	jobId, err := p.getJobId(topic)
	if err != nil {
		return "", fmt.Errorf("delayed.Producer.Dispatch: %w", err)
	}
	return jobId, p.queue.Push(&Job{
		Id:        jobId,
		Topic:     topic,
		Delay:     0,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
	})
}

// DispatchDelay 向延迟队列投递任务 delay 为延迟时间, delayUnit 为延迟时间单位
func (p *Producer) DispatchDelay(topic string, payload []byte, delay int, delayUnit DelayUnit) (string, error) {
	delayTime := DelayTime(delay, delayUnit)
	jobId, err := p.getJobId(topic)
	if err != nil {
		return "", fmt.Errorf("delayed.Producer.DispatchDelay: %w", err)
	}
	return jobId, p.queue.PushDelay(&Job{
		Id:        jobId,
		Topic:     topic,
		Delay:     delayTime,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
	})
}

// getJobId 获取一个可用的 JobId
func (p *Producer) getJobId(topic string) (string, error) {
	var (
		retries int
		jobId   string
	)
	for {
		jobId = strings.ReplaceAll(uuid.New().String(), "-", "")
		exist, err := p.queue.Exists(topic, jobId)
		if err != nil {
			retries++
			continue
		}
		if exist {
			retries++
			if retries >= 5 {
				return "", ErrGetJobIdRetries
			}
		} else {
			break
		}
	}
	return jobId, nil
}
