package delayed

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lsjhtang/kid/config"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/logger"
)

const (
	queueReady   = iota // 队列就绪
	queueClosed         // 队列关闭
	queueStarted        // 队列启动
)

type RedisQueue struct {
	rds      *redis.Client `inject:""`
	topics   []string
	stopCh   chan struct{}
	status   uint32
	canClose bool // 是否可以关闭

	// 用于操作队列的 Redis-lua脚本
	push, pop, pushDelay, cancelDelay, migrate *redisv8.Script

	sync.RWMutex // 读写锁

	slot string // Redis 集群的槽位
}

// pushScript 用于将任务推送到就绪队列
var pushScript = `
-- 获取KEYS
-- local readyQueueKey = KEYS[1]
-- local delayJobKey = KEYS[2]

-- 获取ARGV内的参数
local jobId = ARGV[1]
local jobBytes = ARGV[2]

-- 判断任务是否存在
local exists = redis.call("HEXISTS", KEYS[2], jobId)
if exists == 1 then
    return redis.error_reply("job_id already exists")
end

-- 保存任务详情
local saveDetail = redis.call("HSETNX", KEYS[2], jobId, jobBytes)
if not saveDetail then
    return redis.error_reply("save job detail failed")
end

-- 将任务添加到就绪队列
local pushReady = redis.call("RPUSH", KEYS[1], jobId)
if not pushReady then
    return redis.error_reply("push ready queue failed")
end
return 0
`

// popScript 用于从就绪队列中弹出任务
var popScript = `
-- 获取KEYS
-- local readyQueueKey = KEYS[1]
-- local delayJobKey = KEYS[2]

-- 从就绪队列中弹出
local jobId = redis.call("LPOP", KEYS[1])
if not jobId then
	return 0
end

-- 从任务详情中获取
local jobBytes = redis.call("HGET", KEYS[2], jobId)
if not jobBytes then
	return redis.error_reply("get job detail failed")
end

-- 从任务详情中删除
local delJob = redis.call("HDEL", KEYS[2], jobId)
return jobBytes
`

// pushDelayScript 用于将任务推送到延迟队列
var pushDelayScript = `
-- 获取KEYS
-- local delayQueueKey = KEYS[1]
-- local delayJobKey = KEYS[2]

-- 获取ARGV内的参数
local jobId = ARGV[1]
local jobBytes = ARGV[2]
local delay = ARGV[3]

-- 判断任务是否存在
local exists = redis.call("HEXISTS", KEYS[2], jobId)
if exists == 1 then
	return redis.error_reply("job_id already exists")
end

-- 保存任务详情
local saveDetail = redis.call("HSETNX", KEYS[2], jobId, jobBytes)
if not saveDetail then
	return redis.error_reply("save job detail failed")
end

-- 将任务添加到延迟队列
local pushDelay = redis.call("ZADD", KEYS[1], delay, jobId)
if not pushDelay then
	return redis.error_reply("push delay queue failed")
end
return 0
`

// cancelDelayScript 用于取消延迟任务
var cancelDelayScript = `
-- 获取KEYS
-- local readyQueueKey = KEYS[1]
-- local delayQueueKey = KEYS[2]
-- local delayJobKey = KEYS[3]

-- 获取ARGV内的参数
local jobId = ARGV[1]

-- 判断任务是否存在
local exists = redis.call("HEXISTS", KEYS[3], jobId)
if exists == 0 then
	return redis.error_reply("job_id not exists")
end

-- 从任务详情中删除
local delJob = redis.call("HDEL", KEYS[3], jobId)
if not delJob then
	return redis.error_reply("delete job detail failed")
end

-- 从就绪队列中删除
local delReady = redis.call("LREM", KEYS[1], 0, jobId)
if not delReady then
	return redis.error_reply("delete ready queue failed")
end

-- 从延迟队列中删除
local delDelay = redis.call("ZREM", KEYS[2], jobId)
if not delDelay then
	return redis.error_reply("delete delay queue failed")
end
return 0
`

// migrateScript 用于迁移延迟队列中的任务到就绪队列
var migrateScript = `
-- 获取KEYS
-- local readyQueueKey = KEYS[1]
-- local delayQueueKey = KEYS[2]

-- 获取ARGV内的参数
local now = ARGV[1]

-- 获取延迟队列中的任务
local jobIds = redis.call("ZRANGEBYSCORE", KEYS[2], 0, now)
if not jobIds then
	return redis.error_reply("get job ids failed")
end

for _, jobId in ipairs(jobIds) do
    -- 从任务详情中获取
    local pushReady = redis.call("RPUSH", KEYS[1], jobId)
	if not pushReady then
		return redis.error_reply("push ready queue failed")
	end

    -- 从任务详情中删除
	local delDelay = redis.call("ZREM", KEYS[2], jobId)
	if not delDelay then
		return redis.error_reply("delete delay queue failed")
	end
end
return 0
`

func (rq *RedisQueue) Provide(context.Context) interface{} {
	return NewRedisQueue()
}

func NewRedisQueue() Queue {
	return &RedisQueue{
		stopCh:      make(chan struct{}),
		status:      queueReady,
		push:        redisv8.NewScript(pushScript),
		pop:         redisv8.NewScript(popScript),
		pushDelay:   redisv8.NewScript(pushDelayScript),
		cancelDelay: redisv8.NewScript(cancelDelayScript),
		migrate:     redisv8.NewScript(migrateScript),
		slot:        config.GetString("delayed.redisPrefix"),
		RWMutex:     sync.RWMutex{},
	}
}

func (rq *RedisQueue) Run() error {
	if rq.status == queueStarted {
		return fmt.Errorf("定时任务队列: 队列已经启动")
	}
	go func() {
		atomic.StoreUint32(&rq.status, queueStarted)
		go rq.watchSystemSignal()
		go func() {
			err := rq.migrationDueJobs()
			if err != nil {
				logger.Errorf("定时任务队列: 迁移延迟队列中的任务到就绪队列失败: %v", err)
			}
		}()
		<-rq.stopCh
		rq.Close()
	}()
	return nil
}

func (rq *RedisQueue) AddTopic(topic string) error {
	if rq.status == queueClosed {
		return fmt.Errorf("定时任务队列: 队列已经关闭")
	}
	rq.Lock()
	defer rq.Unlock()
	rq.topics = append(rq.topics, topic)
	return nil
}

// migration 用于迁移延迟队列中的任务到就绪队列
func (rq *RedisQueue) migrationDueJobs() error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if rq.status == queueClosed {
				return nil
			}
			now := time.Now().Unix()
			for _, topic := range rq.topics {
				go func(topic string) {
					if err := rq.migrate.Run(context.Background(), rq.rds.Client,
						[]string{GetReadyQueueKey(rq.slot, topic), GetDelayQueueKey(rq.slot, topic)}, now).Err(); err != nil {
						logger.Errorf("定时任务队列: 迁移延迟队列中的任务到就绪队列失败, topic: %s, err: %s", topic, err.Error())
					}
				}(topic)
			}
		}
	}
}

// watchSystemSignal 监听系统退出信号
func (rq *RedisQueue) watchSystemSignal() {
	signals := make(chan os.Signal, 1)
	// 监听所有 Linux 常见的系统退出信号
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT,
		syscall.SIGBUS, syscall.SIGFPE, syscall.SIGSEGV, syscall.SIGPIPE,
		syscall.SIGALRM, syscall.SIGTERM)
	for {
		select {
		case <-signals: // 收到系统信号
			logger.Infof("定时任务队列: 收到系统退出信号, 服务即将关闭")
			rq.stopCh <- struct{}{}
			return
		}
	}
}

// Push 将任务推送到就绪队列
func (rq *RedisQueue) Push(job *Job) error {
	jobId := job.Id
	jobBytes, err := msgpack.Marshal(job)
	if err != nil {
		return fmt.Errorf("delayed.RedisQueue.Push: Job json 序列化失败, err: %s", err.Error())
	}
	if err := rq.push.Run(context.Background(), rq.rds.Client,
		[]string{GetReadyQueueKey(rq.slot, job.Topic), GetDelayJobKey(rq.slot, job.Topic)}, jobId, string(jobBytes)).Err(); err != nil {
		return fmt.Errorf("delayed.RedisQueue.Push: 执行 Push 脚本错误: %s", err.Error())
	}
	return nil
}

func (rq *RedisQueue) Pop(topic string) (*Job, error) {
	rest, err := rq.pop.Run(context.Background(), rq.rds.Client, []string{GetReadyQueueKey(rq.slot, topic), GetDelayJobKey(rq.slot, topic)}).Result()
	if err != nil {
		return nil, fmt.Errorf("delayed.RedisQueue.Pop: 执行 Pop 脚本错误: %s", err.Error())
	}
	if cast.ToString(rest) == "0" {
		return nil, nil
	}
	var job Job
	if err := msgpack.Unmarshal([]byte(cast.ToString(rest)), &job); err != nil {
		return nil, fmt.Errorf("delayed.RedisQueue.Pop: Job json 反序列化失败, err: %s", err.Error())
	}
	return &job, nil
}

func (rq *RedisQueue) PushDelay(job *Job) error {
	jobBytes, err := msgpack.Marshal(job)
	if err != nil {
		return fmt.Errorf("delayed.RedisQueue.PushDelay: Job json 序列化失败, err: %s", err.Error())
	}
	if err := rq.pushDelay.Run(context.Background(), rq.rds.Client,
		[]string{GetDelayQueueKey(rq.slot, job.Topic), GetDelayJobKey(rq.slot, job.Topic)}, job.Id, string(jobBytes), job.Delay).Err(); err != nil {
		return fmt.Errorf("delayed.RedisQueue.PushDelay: 执行 PushDelay 脚本错误: %s", err.Error())
	}
	return nil
}

// Cancel 取消延迟任务
func (rq *RedisQueue) Cancel(topic, jobId string) error {
	if err := rq.cancelDelay.Run(context.Background(), rq.rds.Client, []string{GetReadyQueueKey(rq.slot, topic),
		GetDelayQueueKey(rq.slot, topic), GetDelayJobKey(rq.slot, topic)}, jobId).Err(); err != nil {
		return fmt.Errorf("delayed.RedisQueue.Cancel: 执行 Cancel 脚本错误: %s", err.Error())
	}
	return nil
}

// FetchReadyJob 获取就绪队列中的任务
func (rq *RedisQueue) FetchReadyJob(topic string) ([]*Job, error) {
	readyJobs, err := rq.rds.LRange(GetReadyQueueKey(rq.slot, topic), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("delayed.RedisQueue.FetchReadyJob: 获取就绪队列失败, err: %s", err.Error())
	}
	if len(readyJobs) == 0 {
		return nil, nil
	}
	rets, err := rq.rds.HMGet(GetDelayJobKey(rq.slot, topic), readyJobs...).Result()
	if err != nil {
		return nil, fmt.Errorf("delayed.RedisQueue.FetchReadyJob: 获取就绪队列详情失败, err: %s", err.Error())
	}
	var jobs []*Job
	for _, ret := range rets {
		var job Job
		if err := json.Unmarshal([]byte(cast.ToString(ret)), &job); err != nil {
			return nil, fmt.Errorf("delayed.RedisQueue.FetchReadyJob: Job json 反序列化失败, err: %s", err.Error())
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// FetchDelayJob 获取延迟队列中的任务
func (rq *RedisQueue) FetchDelayJob(topic string) ([]*Job, error) {
	delayJobs, err := rq.rds.ZRange(GetDelayQueueKey(rq.slot, topic), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("delayed.RedisQueue.FetchDelayJob: 获取延迟队列失败, err: %s", err.Error())
	}
	if len(delayJobs) == 0 {
		return nil, nil
	}
	rets, err := rq.rds.HMGet(GetDelayJobKey(rq.slot, topic), delayJobs...).Result()
	if err != nil {
		return nil, fmt.Errorf("delayed.RedisQueue.FetchDelayJob: 获取延迟队列详情失败, err: %s", err.Error())
	}
	var jobs []*Job
	for _, ret := range rets {
		var job Job
		if err := json.Unmarshal([]byte(cast.ToString(ret)), &job); err != nil {
			return nil, fmt.Errorf("delayed.RedisQueue.FetchDelayJob: Job json 反序列化失败, err: %s", err.Error())
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// ReadyLen 获取就绪队列长度
func (rq *RedisQueue) ReadyLen(topic string) (int, error) {
	count, err := rq.rds.LLen(GetReadyQueueKey(rq.slot, topic)).Result()
	if err != nil {
		return 0, fmt.Errorf("delayed.RedisQueue.ReadyLen: 获取就绪队列长度失败, err: %s", err.Error())
	}
	return int(count), nil
}

// DelayLen 获取延迟队列长度
func (rq *RedisQueue) DelayLen(topic string) (int, error) {
	count, err := rq.rds.ZCard(GetDelayQueueKey(rq.slot, topic)).Result()
	if err != nil {
		return 0, fmt.Errorf("delayed.RedisQueue.DelayLen: 获取延迟队列长度失败, err: %s", err.Error())
	}
	return int(count), nil
}

// Exists 判断任务是否存在
func (rq *RedisQueue) Exists(topic, jobId string) (bool, error) {
	exist, err := rq.rds.HExists(GetDelayJobKey(rq.slot, topic), jobId).Result()
	if err != nil {
		return false, fmt.Errorf("delayed.RedisQueue.Exists: 判断任务是否存在失败, err: %s", err.Error())
	}
	return exist, nil
}

func (rq *RedisQueue) Close() error {
	atomic.StoreUint32(&rq.status, queueClosed)
	return nil
}
