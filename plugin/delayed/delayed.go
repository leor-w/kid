package delayed

import (
	"context"
)

type Delayed struct {
	product  *Producer `inject:""`        // 任务生产者
	consumer *Consumer `inject:"opts:NR"` // 任务消费者
}

func (d *Delayed) Provide(ctx context.Context) interface{} {
	return d
}

func (d *Delayed) Run() error {
	if err := d.consumer.Run(); err != nil {
		return err
	}
	return nil
}

func (d *Delayed) Register(worker *Worker) {
	d.consumer.Register(worker)
}

func (d *Delayed) Close() error {
	if err := d.consumer.Close(); err != nil {
		return err
	}
	return nil
}

// Dispatch 用于添加任务
func (d *Delayed) Dispatch(topic string, payload []byte) (string, error) {
	return d.product.Dispatch(topic, payload)
}

// DispatchDelay 用于添加延迟任务
func (d *Delayed) DispatchDelay(topic string, payload []byte, delay int, delayUnit DelayUnit) (string, error) {
	return d.product.DispatchDelay(topic, payload, delay, delayUnit)
}

// Cancel 通过指定 topic 和 jobId 取消延迟任务
func (d *Delayed) Cancel(topic, jobId string) error {
	return d.consumer.Cancel(topic, jobId)
}
