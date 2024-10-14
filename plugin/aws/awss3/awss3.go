package awss3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials"

	config2 "github.com/aws/aws-sdk-go-v2/config"

	"github.com/leor-w/kid/utils"

	"github.com/leor-w/injector"
	"github.com/leor-w/kid/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsS3 struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	options       *Options
}

type Option func(*Options)

func (aws *AwsS3) Provide(ctx context.Context) any {
	var confName string
	name, ok := ctx.Value(new(injector.NameKey)).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("aws.s3%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件为找到 [%s.*]，请检查配置文件", confPrefix))
	}
	return NewAwsS3(
		WithRegion(config.GetString(utils.GetConfigurationItem(confPrefix, "region"))),
		WithAccessKey(config.GetString(utils.GetConfigurationItem(confPrefix, "access_key"))),
		WithSecretKey(config.GetString(utils.GetConfigurationItem(confPrefix, "secret_key"))),
	)
}

func (aws *AwsS3) GeneratePresignedURL(bucketName, key string, expireTime time.Duration) (string, error) {
	req, err := aws.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}
	return req.URL, nil
}

func NewAwsS3(options ...Option) *AwsS3 {
	var opts Options
	for _, o := range options {
		o(&opts)
	}
	cfg, err := config2.LoadDefaultConfig(context.TODO(),
		config2.WithRegion(opts.Region),
		config2.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, "")),
	)
	if err != nil {
		panic(fmt.Sprintf("创建S3客户端失败: %v", err))
	}
	s3Client := s3.NewFromConfig(cfg)

	return &AwsS3{
		s3Client:      s3Client,
		presignClient: s3.NewPresignClient(s3Client),
		options:       &opts,
	}
}
