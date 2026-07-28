package qiniu

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lsjhtang/kid/injector"

	"github.com/qiniu/go-sdk/v7/auth"

	"github.com/qiniu/go-sdk/v7/cdn"

	"github.com/jinzhu/copier"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/utils"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Qiniu struct {
	options          *Options
	mac              *qbox.Mac
	bucketManager    *storage.BucketManager
	operationManager *storage.OperationManager
	cdnManger        *cdn.CdnManager
}

func (qiniu *Qiniu) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("qiniu%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithDomain(config.GetString(utils.GetConfigurationItem(confPrefix, "domain"))),
		WithBucket(config.GetString(utils.GetConfigurationItem(confPrefix, "bucket"))),
		WithAccess(config.GetString(utils.GetConfigurationItem(confPrefix, "access"))),
		WithSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "secret"))),
		WithPrivate(config.GetBool(utils.GetConfigurationItem(confPrefix, "private"))),
		WithTTL(config.GetDuration(utils.GetConfigurationItem(confPrefix, "ttl"))),
		WithStorageConfig(&Config{
			UseHttps:      config.GetBool(utils.GetConfigurationItem(confPrefix, "bucketConf.useHttps")),
			UseCdnDomains: config.GetBool(utils.GetConfigurationItem(confPrefix, "bucketConf.useDomains")),
			CentralRsHost: config.GetString(utils.GetConfigurationItem(confPrefix, "bucketConf.centralRsHost")),
		}),
		WithOperationConf(&Config{
			UseHttps:      config.GetBool(utils.GetConfigurationItem(confPrefix, "operation.useHttps")),
			UseCdnDomains: config.GetBool(utils.GetConfigurationItem(confPrefix, "operation.useDomains")),
			CentralRsHost: config.GetString(utils.GetConfigurationItem(confPrefix, "operation.centralRsHost")),
			NotifyUrl:     config.GetString(utils.GetConfigurationItem(confPrefix, "operation.notifyUrl")),
			Pipeline:      config.GetString(utils.GetConfigurationItem(confPrefix, "operation.pipeline")),
		}),
	)
}

type Option func(*Options)

func (qiniu *Qiniu) UploadToken(policies ...*Policy) string {
	var (
		policy      *Policy
		storePolicy storage.PutPolicy
	)
	if len(policies) > 0 {
		policy = policies[0]
		if len(policy.Scope) <= 0 {
			policy.Scope = qiniu.options.bucket
		}
		if err := copier.Copy(&storePolicy, &policy); err != nil {
			return ""
		}
	} else {
		storePolicy = storage.PutPolicy{Scope: qiniu.options.bucket}
	}
	return storePolicy.UploadToken(qiniu.mac)
}

// MakePrivateURL 获取私有文件的下载链接地址
func (qiniu *Qiniu) MakePrivateURL(domain, key string, deadline int64) string {
	return storage.MakePrivateURLv2(qiniu.mac, domain, key, deadline)
}

func (qiniu *Qiniu) MakePrivateURLWithQuery(domain, key string, query url.Values, deadline int64) string {
	return storage.MakePrivateURLv2WithQuery(qiniu.mac, domain, key, query, deadline)
}

func (qiniu *Qiniu) MakePrivateURLWithQueryString(domain, key, query string, deadline int64) string {
	return storage.MakePrivateURLv2WithQueryString(qiniu.mac, domain, key, query, deadline)
}

type FopConfig struct {
	Fops    []string
	Force   bool
	FileKey string
}

func (qiniu *Qiniu) Pfop(config *FopConfig) (string, error) {
	fops := strings.Join(config.Fops, ";")
	persistentId, err := qiniu.operationManager.Pfop(qiniu.options.bucket, config.FileKey,
		fops, qiniu.options.operationConf.Pipeline, qiniu.options.operationConf.NotifyUrl, config.Force)
	if err != nil {
		return "", err
	}
	return persistentId, nil
}

func (qiniu *Qiniu) Prefop(persistentId string) (storage.PrefopRet, error) {
	return qiniu.operationManager.Prefop(persistentId)
}

func (qiniu *Qiniu) Move(srcBucket, srcFile, destBucket, destFile string, force bool) error {
	return qiniu.bucketManager.Move(srcBucket, srcFile, destBucket, destFile, force)
}

func (qiniu *Qiniu) Delete(file string) error {
	return qiniu.bucketManager.Delete(qiniu.options.bucket, file)
}

func (qiniu *Qiniu) Deletes(files ...string) error {
	for _, file := range files {
		if err := qiniu.Delete(file); err != nil {
			return err
		}
	}
	return nil
}

func (qiniu *Qiniu) Stat(file string) (storage.FileInfo, error) {
	return qiniu.bucketManager.Stat(qiniu.options.bucket, file)
}

// SmallFilesMkZip 少量文件压缩
func (qiniu *Qiniu) SmallFilesMkZip(conf *ZipConfig) string {
	if len(conf.SaveAs) <= 0 {
		return ""
	}
	fop := strings.Builder{}
	fop.WriteString(fmt.Sprintf("mkzip/2/encoding/%s", base64.URLEncoding.EncodeToString([]byte("utf-8"))))
	for _, file := range conf.ZipFiles {
		f := storage.MakePrivateURLv2(qiniu.mac, qiniu.options.domain, file.Source, time.Now().Add(time.Minute*30).Unix())
		fop.WriteString(fmt.Sprintf("/url/%s", base64.URLEncoding.EncodeToString([]byte(f))))
		if len(file.Alias) > 0 {
			fop.WriteString(fmt.Sprintf("/alias/%s", base64.URLEncoding.EncodeToString([]byte(file.Alias))))
		}
	}
	return fmt.Sprintf("%s|saveas/%s", fop.String(),
		base64.URLEncoding.EncodeToString([]byte(qiniu.options.bucket+":"+conf.SaveAs)))
}

func (qiniu *Qiniu) EncodeEntry(fileName string) string {
	return storage.EncodedEntry(qiniu.options.bucket, fileName)
}

func (qiniu *Qiniu) CdnManager() *cdn.CdnManager {
	if qiniu.cdnManger == nil {
		qiniu.cdnManger = cdn.NewCdnManager(auth.New(qiniu.options.access, qiniu.options.secret))
	}
	return qiniu.cdnManger
}

func (qiniu *Qiniu) BucketManager() *storage.BucketManager {
	if qiniu.bucketManager == nil {
		qiniu.bucketManager = storage.NewBucketManager(qiniu.mac, &storage.Config{
			UseHTTPS:      qiniu.options.bucketConf.UseHttps,
			UseCdnDomains: qiniu.options.bucketConf.UseCdnDomains,
			CentralRsHost: qiniu.options.bucketConf.CentralRsHost,
		})
	}
	return qiniu.bucketManager
}

func (qiniu *Qiniu) OperationManager() *storage.OperationManager {
	if qiniu.operationManager == nil {
		qiniu.operationManager = storage.NewOperationManager(qiniu.mac, &storage.Config{
			UseHTTPS:      qiniu.options.operationConf.UseHttps,
			UseCdnDomains: qiniu.options.operationConf.UseCdnDomains,
			CentralRsHost: qiniu.options.operationConf.CentralRsHost,
		})
	}
	return qiniu.operationManager
}

// GetRegion 获取存储区域
func (qiniu *Qiniu) GetRegion() (*storage.Region, error) {
	return storage.GetRegion(qiniu.mac.AccessKey, qiniu.options.bucket)
}

// NewUploader 获取上传对象
func (qiniu *Qiniu) NewUploader(cfg *storage.Config) *storage.FormUploader {
	region, err := qiniu.GetRegion()
	if err != nil {
		return nil
	}
	if cfg.Region == nil {
		cfg.Region = region
	}
	return storage.NewFormUploader(cfg)
}

func New(opts ...Option) *Qiniu {
	options := &Options{
		private: false,
		ttl:     5,
	}
	for _, opt := range opts {
		opt(options)
	}
	mac := qbox.NewMac(options.access, options.secret)
	return &Qiniu{
		options: options,
		mac:     mac,
		bucketManager: storage.NewBucketManager(mac, &storage.Config{
			UseHTTPS:      options.bucketConf.UseHttps,
			UseCdnDomains: options.bucketConf.UseCdnDomains,
			CentralRsHost: options.bucketConf.CentralRsHost,
		}),
		operationManager: storage.NewOperationManager(mac, &storage.Config{
			UseHTTPS:      options.operationConf.UseHttps,
			UseCdnDomains: options.operationConf.UseCdnDomains,
			CentralRsHost: options.operationConf.CentralRsHost,
		}),
		cdnManger: cdn.NewCdnManager(auth.New(options.access, options.secret)),
	}
}
