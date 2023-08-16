package dingtalk

import "time"

// Options 基础配置
type Options struct {
	AccessKey       string        // 验证 key
	AccessKeySecret string        // 验证秘钥
	RegionId        string        // 区域 例：cn-hangzhou
	Network         string        // 需访问服务器的网络类型。例：inner，vpc
	Endpoint        string        // 服务器地址，例：ecs-cn-hangzhou.aliyuncs.com
	ReadTimeout     time.Duration // 读超时，单位：秒
	ConnectTimeout  time.Duration // 连接超时，单位：秒
	MaxIdleConns    int           // 最大连接数
	Robots          []*RobotOptions
}

func WithAccessKey(accessKey string) Option {
	return func(o *Options) {
		o.AccessKey = accessKey
	}
}

func WithAccessKeySecret(accessKeySecret string) Option {
	return func(o *Options) {
		o.AccessKeySecret = accessKeySecret
	}
}

func WithRegionId(regionId string) Option {
	return func(o *Options) {
		o.RegionId = regionId
	}
}

func WithNetwork(network string) Option {
	return func(o *Options) {
		o.Network = network
	}
}

func WithEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.Endpoint = endpoint
	}
}

func WithReadTimeout(readTimeout time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = readTimeout
	}
}

func WithConnectTimeout(connectTimeout time.Duration) Option {
	return func(o *Options) {
		o.ConnectTimeout = connectTimeout
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(o *Options) {
		o.MaxIdleConns = maxIdleConns
	}
}

func WithRobots(robots []*RobotOptions) Option {
	return func(o *Options) {
		o.Robots = robots
	}
}

// RobotOptions 机器人配置
type RobotOptions struct {
	Name           string
	RobotCode      string
	ConversationId string
	Token          string
	CoolAppCode    string
}

func WithName(name string) RobotOption {
	return func(o *RobotOptions) {
		o.Name = name
	}
}

func WithRobotCode(robotCode string) RobotOption {
	return func(o *RobotOptions) {
		o.RobotCode = robotCode
	}
}

func WithConversationId(conversationId string) RobotOption {
	return func(o *RobotOptions) {
		o.ConversationId = conversationId
	}
}

func WithToken(token string) RobotOption {
	return func(o *RobotOptions) {
		o.Token = token
	}
}

func WithCoolAppCode(coolAppCode string) RobotOption {
	return func(o *RobotOptions) {
		o.CoolAppCode = coolAppCode
	}
}
