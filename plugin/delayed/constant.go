package delayed

import (
	"fmt"
	"time"
)

var (
	ErrGetJobIdRetries = fmt.Errorf("获取 JobId 重试次数超过 5 次, 请检查 Redis 连接是否正常")
)

const (
	ReadyQueueKey = "delayed:ready:topic.%s" // 就绪队列的 key
	DelayQueueKey = "delayed:delay:topic.%s" // 延迟SortedSet的 key
	DelayJob      = "delayed:job:topic.%s"   // 延迟任务详情存放的 key
)

// GetReadyQueueKey 获取就绪队列的 key
func GetReadyQueueKey(slot, topic string) string {
	if slot != "" {
		return fmt.Sprintf("{%s}%s", slot, fmt.Sprintf(ReadyQueueKey, topic))
	}
	return fmt.Sprintf(ReadyQueueKey, topic)
}

// GetDelayQueueKey 获取延迟SortedSet的 key
func GetDelayQueueKey(slot, topic string) string {
	if slot != "" {
		return fmt.Sprintf("{%s}%s", slot, fmt.Sprintf(DelayQueueKey, topic))
	}
	return fmt.Sprintf(DelayQueueKey, topic)
}

// GetDelayJobKey 获取延迟任务详情存放的 key
func GetDelayJobKey(slot, topic string) string {
	if slot != "" {
		return fmt.Sprintf("{%s}%s", slot, fmt.Sprintf(DelayJob, topic))
	}
	return fmt.Sprintf(DelayJob, topic)
}

// DelayUnit 延时时间计算类型
type DelayUnit uint8

const (
	DelayUnitSecond DelayUnit = iota + 1 // 延迟计算单位：秒
	DelayUnitMinute                      // 延迟计算单位：分
	DelayUnitHour                        // 延迟计算单位：时
	DelayUnitDay                         // 延迟计算单位：天
)

// DelayTime 计算延迟时间
func DelayTime(delay int, delayUnit DelayUnit) int64 {
	var delayTime int64
	switch delayUnit {
	case DelayUnitSecond:
		delayTime = time.Now().Add(time.Duration(delay) * time.Second).Unix()
	case DelayUnitMinute:
		delayTime = time.Now().Add(time.Duration(delay) * time.Minute).Unix()
	case DelayUnitHour:
		delayTime = time.Now().Add(time.Duration(delay) * time.Hour).Unix()
	case DelayUnitDay:
		delayTime = time.Now().Add(time.Duration(delay) * time.Hour * 24).Unix()
	default:
		delayTime = time.Now().Add(time.Duration(delay) * time.Second).Unix()
	}
	return delayTime
}
