package dingtalk

import (
	"context"
	"fmt"

	"github.com/lsjhtang/kid/plugin/robot"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin"
	"github.com/lsjhtang/kid/utils"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	client "github.com/alibabacloud-go/dingtalk/robot_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
)

var (
	ErrNotFoundRobot = fmt.Errorf("未找到对应机器人")
)

type Dingtalk struct {
	robots  map[string]*Robot
	options *Options
}

type Option func(*Options)

type RobotOption func(*RobotOptions)

func (ding *Dingtalk) Provide(ctx context.Context) interface{} {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("dingtalk%s", configName)
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("配置文件中未找到对应配置项 [%s], 请检查", confPrefix))
	}
	robotConfs, ok := config.Get(utils.GetConfigurationItem(confPrefix, "robots")).([]interface{})
	if !ok {
		panic(fmt.Sprintf("钉钉机器人配置错误, 请检查"))
	}
	var robots []*RobotOptions
	for _, conf := range robotConfs {
		robotConf, ok := conf.(map[string]interface{})
		if !ok {
			panic(fmt.Sprintf("钉钉机器人配置错误, 请检查"))
		}
		robots = append(robots, &RobotOptions{
			Name:           cast.ToString(robotConf["name"]),
			RobotCode:      cast.ToString(robotConf["robot_code"]),
			ConversationId: cast.ToString(robotConf["conversation_id"]),
			Token:          cast.ToString(robotConf["token"]),
			CoolAppCode:    cast.ToString(robotConf["cool_app_code"]),
		})
	}
	return New(WithRobots(robots))
}

func (ding *Dingtalk) SendMessage(params interface{}) error {
	reqConf, ok := params.(*SendMessageRequest)
	if !ok {
		return robot.ErrWrongParamType
	}
	r, exist := ding.robots[reqConf.RobotName]
	if !exist {
		return ErrNotFoundRobot
	}
	_, err := r.client.OrgGroupSendWithOptions(
		&client.OrgGroupSendRequest{
			CoolAppCode:        tea.String(r.options.CoolAppCode),
			MsgKey:             tea.String(reqConf.MsgKey),
			MsgParam:           tea.String(reqConf.MsgParam),
			OpenConversationId: tea.String(r.options.ConversationId),
			RobotCode:          tea.String(r.options.RobotCode),
			Token:              tea.String(r.options.Token),
		},
		&client.OrgGroupSendHeaders{
			XAcsDingtalkAccessToken: tea.String(r.options.Token),
		},
		&util.RuntimeOptions{})
	return err
}

func (ding *Dingtalk) WithdrawMessage(params interface{}) error {
	reqConf, ok := params.(*WithdrawMessageRequest)
	if !ok {
		return robot.ErrWrongParamType
	}
	r, exist := ding.robots[reqConf.RobotName]
	if !exist {
		return ErrNotFoundRobot
	}
	_, err := r.client.OrgGroupRecallWithOptions(
		&client.OrgGroupRecallRequest{
			OpenConversationId: tea.String(r.options.ConversationId),
			ProcessQueryKeys:   tea.StringSlice(reqConf.ProcessQueryKeys),
			RobotCode:          tea.String(r.options.RobotCode),
		},
		&client.OrgGroupRecallHeaders{
			XAcsDingtalkAccessToken: tea.String(r.options.Token),
		},
		&util.RuntimeOptions{})
	return err
}

func New(options ...Option) *Dingtalk {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	var robots map[string]*Robot
	for _, r := range opts.Robots {
		robots[r.Name] = NewRobot(r)
	}
	return &Dingtalk{
		options: opts,
		robots:  robots,
	}
}

type Robot struct {
	client  *client.Client
	options *RobotOptions
}

func NewRobot(options *RobotOptions) *Robot {
	conf := &openapi.Config{
		Protocol: tea.String("https"),
		RegionId: tea.String("central"),
	}
	cli, err := client.NewClient(conf)
	if err != nil {
		panic(fmt.Sprintf("钉钉机器人创建客户端错误: %s", err.Error()))
	}
	return &Robot{
		client:  cli,
		options: options,
	}
}
