package aliyun

import (
	"context"
	"fmt"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/utils"

	afs "github.com/alibabacloud-go/afs-20180112/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/lsjhtang/kid/config"
)

const endpoint = "afs.aliyuncs.com"

type NoTraceBehavioral struct {
	config  *rpc.Config
	client  *afs.Client
	options *Options
}

func (b *NoTraceBehavioral) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("aliyun.afs%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件未找到配置项 [%s]", confPrefix))
	}
	return New(
		WithAccessKeyId(config.GetString(utils.GetConfigurationItem(confPrefix, "accessKeyId"))),
		WithAccessKeySecret(config.GetString(utils.GetConfigurationItem(confPrefix, "accessKeySecret"))),
		WithRegionId(config.GetString(utils.GetConfigurationItem(confPrefix, "regionId"))),
	)
}

func (b *NoTraceBehavioral) Verify(ip, data string) (string, error) {
	req := new(afs.AnalyzeNvcRequest)
	req.SetSourceIp(ip)
	req.SetScoreJsonStr("{\"200\":\"PASS\",\"400\":\"NC\",\"800\":\"BLOCK\"}")
	req.SetData(data)
	resp, err := b.client.AnalyzeNvc(req)
	if err != nil {
		return "", err
	}
	return *resp.BizCode, nil
}

func (b *NoTraceBehavioral) Init() error {
	b.config = new(rpc.Config)
	b.config.SetAccessKeyId(b.options.AccessKeyId).
		SetAccessKeySecret(b.options.AccessKeySecret).
		SetRegionId(b.options.RegionId).
		SetEndpoint(endpoint)
	cli, err := afs.NewClient(b.config)
	if err != nil {
		return err
	}
	b.client = cli
	return nil
}

func New(opts ...Option) *NoTraceBehavioral {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	var behavioral = &NoTraceBehavioral{
		options: &options,
	}
	if err := behavioral.Init(); err != nil {
		panic(fmt.Sprintf("初始化阿里云行为式验证码失败: %s", err.Error()))
	}
	return behavioral
}
