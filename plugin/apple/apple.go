package apple

import (
	"context"
	"fmt"
	"os"

	gopay "github.com/go-pay/gopay/apple"
	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/utils"
)

type AppleStore struct {
	client  *gopay.Client
	Options *Options
}

type Option func(*Options)

func (apple *AppleStore) Provide(ctx context.Context) any {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("apple.pay%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件为找到 [%s.*]，请检查配置文件", confPrefix))
	}
	return New(
		WithIsProduct(config.GetBool(utils.GetConfigurationItem(confPrefix, "is_product"))),
		WithKeyId(config.GetString(utils.GetConfigurationItem(confPrefix, "key_id"))),
		WithIssuerId(config.GetString(utils.GetConfigurationItem(confPrefix, "issuer_id"))),
		WithBid(config.GetString(utils.GetConfigurationItem(confPrefix, "bundle_id"))),
		WithPrivateKey(config.GetString(utils.GetConfigurationItem(confPrefix, "private_key"))),
		WithPrivateKeyFile(config.GetString(utils.GetConfigurationItem(confPrefix, "private_key_file"))),
	)
}

func (apple *AppleStore) init() error {
	if apple.Options.PrivateKeyFile == "" && apple.Options.PrivateKey == "" {
		return fmt.Errorf("apple: 私钥文件和私钥不能同时为空")
	}
	privateKey := apple.Options.PrivateKey
	if len(privateKey) <= 0 {
		privateKeyBytes, err := os.ReadFile(apple.Options.PrivateKeyFile)
		if err != nil {
			return fmt.Errorf("apple: 读取私钥文件失败: %w", err)
		}
		privateKey = string(privateKeyBytes)
	}

	client, err := gopay.NewClient(
		apple.Options.IssuerId,
		apple.Options.Bid,
		apple.Options.KeyId,
		privateKey,
		apple.Options.IsProduct,
	)
	if err != nil {
		return fmt.Errorf("apple: 创建客户端失败: %w", err)
	}
	apple.client = client
	return nil
}

func New(opts ...Option) *AppleStore {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	var apple AppleStore
	apple.Options = options
	if err := apple.init(); err != nil {
		panic(err)
	}
	return &apple
}
