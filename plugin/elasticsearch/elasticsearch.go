package elasticsearch

import (
	"context"
	"fmt"

	"github.com/lsjhtang/kid/injector"

	"github.com/elastic/go-elasticsearch/v7"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/utils"
)

type Elasticsearch struct {
	*elasticsearch.Client
	options *Options
}

func (es *Elasticsearch) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("elasticsearch%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithAddresses(config.GetStringSlice(utils.GetConfigurationItem(confPrefix, "addresses"))),
		WithUsername(config.GetString(utils.GetConfigurationItem(confPrefix, "username"))),
		WithPassword(config.GetString(utils.GetConfigurationItem(confPrefix, "password"))),
		WithCloudID(config.GetString(utils.GetConfigurationItem(confPrefix, "cloudId"))),
		WithApiKey(config.GetString(utils.GetConfigurationItem(confPrefix, "apiKey"))),
		WithServiceToken(config.GetString(utils.GetConfigurationItem(confPrefix, "serviceToken"))),
		WithCertificateFingerprint(config.GetString(utils.GetConfigurationItem(confPrefix, "certificateFingerprint"))),
		WithHeader(config.GetStringMapStringSlice(utils.GetConfigurationItem(confPrefix, "header"))),
		WithCACert(config.GetString(utils.GetConfigurationItem(confPrefix, "caCert"))),
		WithRetryOnStatus(config.GetIntSlice(utils.GetConfigurationItem(confPrefix, "retryOnStatus"))),
		WithMaxRetries(config.GetInt(utils.GetConfigurationItem(confPrefix, "maxRetries"))),
		WithCompressRequestBody(config.GetBool(utils.GetConfigurationItem(confPrefix, "compressRequestBody"))),
		WithDiscoverNodesOnStart(config.GetBool(utils.GetConfigurationItem(confPrefix, "discoverNodesOnStart"))),
		WithDiscoverNodesInterval(config.GetDuration(utils.GetConfigurationItem(confPrefix, "discoverNodesInterval"))),
		WithEnableMetrics(config.GetBool(utils.GetConfigurationItem(confPrefix, "enableMetrics"))),
		WithEnableDebugLogger(config.GetBool(utils.GetConfigurationItem(confPrefix, "enableDebugLogger"))),
		WithEnableCompatibilityMode(config.GetBool(utils.GetConfigurationItem(confPrefix, "enableCompatibilityMode"))),
		WithDisableMetaHeader(config.GetBool(utils.GetConfigurationItem(confPrefix, "disableMetaHeader"))),
	)
}

func (es *Elasticsearch) Init(opts ...Option) error {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:             options.Addresses,
		Username:              options.Username,
		Password:              options.Password,
		CloudID:               options.CloudID,
		APIKey:                options.APIKey,
		Header:                options.Header,
		CACert:                nil,
		RetryOnStatus:         options.RetryOnStatus,
		DisableRetry:          options.DisableRetry,
		EnableRetryOnTimeout:  false,
		MaxRetries:            options.MaxRetries,
		DiscoverNodesOnStart:  options.DiscoverNodesOnStart,
		DiscoverNodesInterval: options.DiscoverNodesInterval,
		EnableMetrics:         options.EnableMetrics,
		EnableDebugLogger:     options.EnableDebugLogger,
		RetryBackoff:          nil,
		Transport:             nil,
		Logger:                nil,
		Selector:              nil,
		ConnectionPoolFunc:    nil,
	})
	if err != nil {
		return fmt.Errorf("init Elasticsearch client error: %s", err.Error())
	}
	es.Client = client
	es.options = &options
	return nil
}

func New(opts ...Option) *Elasticsearch {
	es := &Elasticsearch{}
	if err := es.Init(opts...); err != nil {
		panic(err.Error())
	}
	return es
}
