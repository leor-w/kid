package awss3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/injector"
	"github.com/lsjhtang/kid/utils"
)

type AwsS3 struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	options       *Options
}

type Option func(*Options)

func (awsS3 *AwsS3) Provide(ctx context.Context) any {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
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
		WithCacheControl(config.GetString(utils.GetConfigurationItem(confPrefix, "cache_control"))),
	)
}

func (awsS3 *AwsS3) init() {
	awsS3.Client()
	awsS3.PreSignClient()
}

func (awsS3 *AwsS3) PreSignClient() *s3.PresignClient {

	if awsS3.presignClient == nil {
		awsS3.presignClient = s3.NewPresignClient(awsS3.s3Client)
	}
	return awsS3.presignClient
}

func (awsS3 *AwsS3) Client() *s3.Client {
	if awsS3.s3Client == nil {
		cfg, err := awsConf.LoadDefaultConfig(context.TODO(),
			awsConf.WithRegion(awsS3.options.Region),
			awsConf.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					awsS3.options.AccessKey,
					awsS3.options.SecretKey,
					"",
				),
			),
		)
		if err != nil {
			panic(fmt.Sprintf("创建S3客户端失败: %v", err))
		}
		awsS3.s3Client = s3.NewFromConfig(cfg)
	}
	return awsS3.s3Client
}

// GetPreSignUploadURL 获取预签名上传链接
func (awsS3 *AwsS3) GetPreSignUploadURL(conf *GetPreSignURLConf) (string, error) {
	cacheControl := awsS3.options.CacheControl
	if conf.CacheControl != "" {
		cacheControl = conf.CacheControl
	}
	req, err := awsS3.PreSignClient().PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:       &conf.Bucket,
		Key:          &conf.ObjectKey,
		CacheControl: aws.String(cacheControl),
	}, s3.WithPresignExpires(conf.Expires))
	if err != nil {
		return "", fmt.Errorf("获取预签名上传链接失败: %v", err)
	}
	return req.URL, nil
}

// GetPreSignDownloadURL 获取预签名下载链接
func (awsS3 *AwsS3) GetPreSignDownloadURL(bucket, key string, expires time.Duration) (string, error) {
	req, err := awsS3.PreSignClient().PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", fmt.Errorf("获取预签名下载链接失败: %v", err)
	}
	return req.URL, nil
}

// CreateMultipartUpload 创建分片上传
func (awsS3 *AwsS3) CreateMultipartUpload(conf *GetPreSignURLConf) (*MultipartResponse, error) {
	cacheControl := awsS3.options.CacheControl
	if conf.CacheControl != "" {
		cacheControl = conf.CacheControl
	}
	expiredAt := time.Now().Add(conf.Expires)
	input := &s3.CreateMultipartUploadInput{
		Bucket:       &conf.Bucket,
		Key:          &conf.ObjectKey,
		ContentType:  &conf.ContentType,
		CacheControl: &cacheControl,
		Expires:      &expiredAt,
	}
	resp, err := awsS3.Client().CreateMultipartUpload(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("创建分片上传失败: %v", err)
	}
	return &MultipartResponse{
		UploadId: *resp.UploadId,
	}, nil
}

// GetMultipartUploadPreSignURL 获取预签名分片上传链接
// bucket: 存储桶
// uploadId: 上传ID
// objectKey: 对象键
// partNumber: 分片编号
// expires: 预签 URL 过期时间
func (awsS3 *AwsS3) GetMultipartUploadPreSignURL(conf *MultipartUploadPreSignConfig) (string, error) {
	req, err := awsS3.PreSignClient().PresignUploadPart(context.TODO(), &s3.UploadPartInput{
		Bucket:     aws.String(conf.Bucket),
		Key:        aws.String(conf.ObjectKey),
		PartNumber: aws.Int32(conf.PartNumber),
		UploadId:   aws.String(conf.UploadId),
	}, s3.WithPresignExpires(conf.Expires))
	if err != nil {
		return "", fmt.Errorf("获取预签名分片上传链接失败: %v", err)
	}
	return req.URL, nil
}

// AbortMultipartUpload 取消分片上传
func (awsS3 *AwsS3) AbortMultipartUpload(conf *CancelMultipartUploadConfig) error {
	var owner *string
	if len(conf.ExpectedBucketOwner) > 0 {
		owner = &conf.ExpectedBucketOwner
	}
	_, err := awsS3.Client().AbortMultipartUpload(context.TODO(), &s3.AbortMultipartUploadInput{
		Bucket:              aws.String(conf.Bucket),
		Key:                 aws.String(conf.ObjectKey),
		UploadId:            aws.String(conf.UploadId),
		ExpectedBucketOwner: owner,
	})
	if err != nil {
		return fmt.Errorf("取消分片上传失败: %v", err)
	}
	return nil
}

// CompleteMultipartUpload 完成分片上传, 将会合并所有分片
// conf 包含了完成分片上传所需的参数
func (awsS3 *AwsS3) CompleteMultipartUpload(conf *CompleteMultipartUploadConfig) error {
	completedParts := make([]types.CompletedPart, len(conf.Parts))
	for i, part := range conf.Parts {
		completedParts[i] = types.CompletedPart{
			ETag:       aws.String(part.ETag),
			PartNumber: aws.Int32(part.PartNumber),
		}
	}
	input := &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(conf.Bucket),
		Key:             aws.String(conf.ObjectKey),
		UploadId:        aws.String(conf.UploadId),
		MultipartUpload: &types.CompletedMultipartUpload{Parts: completedParts},
	}
	if _, err := awsS3.Client().CompleteMultipartUpload(context.TODO(), input); err != nil {
		return fmt.Errorf("完成分片上传失败: %v", err)
	}
	return nil
}

// ListFileByPrefix 获取指定前缀的文件列表
func (awsS3 *AwsS3) ListFileByPrefix(bucket, prefix string) ([]types.Object, error) {
	objects, err := awsS3.Client().ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %s", err.Error())
	}
	return objects.Contents, nil
}

func (awsS3 *AwsS3) DeleteFile(conf *DeleteObjectConfig) error {
	var owner *string
	if len(conf.ExpectedBucketOwner) > 0 {
		owner = &conf.ExpectedBucketOwner
	}
	var versionId *string
	if len(conf.VersionId) > 0 {
		versionId = &conf.VersionId
	}
	if _, err := awsS3.Client().DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket:              aws.String(conf.Bucket),
		Key:                 aws.String(conf.ObjectKey),
		ExpectedBucketOwner: owner,
		VersionId:           versionId,
	}); err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}

func NewAwsS3(options ...Option) *AwsS3 {
	var opts Options
	for _, o := range options {
		o(&opts)
	}
	awsS3 := &AwsS3{
		options: &opts,
	}
	awsS3.init()
	return awsS3
}
