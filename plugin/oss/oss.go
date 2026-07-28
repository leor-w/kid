package oss

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/injector"
	"github.com/lsjhtang/kid/utils"
)

type OSS struct {
	options *Options
}

func (alioss *OSS) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("oss%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithEndpoint(config.GetString(utils.GetConfigurationItem(confPrefix, "endpoint"))),
		WithBucketName(config.GetString(utils.GetConfigurationItem(confPrefix, "bucketName"))),
		WithAccessKey(config.GetString(utils.GetConfigurationItem(confPrefix, "accessKey"))),
		WithSecretKey(config.GetString(utils.GetConfigurationItem(confPrefix, "secretKey"))),
	)
}

type Option func(*Options)

func New(opts ...Option) *OSS {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &OSS{
		options: options,
	}
}

func (alioss *OSS) Upload(_ context.Context, dir string, file *multipart.FileHeader) (interface{}, error) {
	var url string
	// 打开文件
	fileHandle, err := file.Open()
	if err != nil {
		return "", err
	}
	fmt.Println("空间名:", alioss.options.bucketName)
	defer func(fileHandle multipart.File) {
		_ = fileHandle.Close()
	}(fileHandle)
	// 计算文件的MD5值
	// 检查数据库中是否已存在该文件
	// 配置OSS客户端
	// 创建OSS服务实例
	//client, err := oss.New(entity.Endpoint, entity.AccessKeyId, entity.AccessKeySecret)
	client, err := oss.New(alioss.options.endpoint, alioss.options.accessKey, alioss.options.secretKey)
	if err != nil {
		return "", err
	}
	// 获取存储空间
	bucket, err := client.Bucket(alioss.options.bucketName)
	if err != nil {
		return "", err
	}
	// 上传文件
	path := strconv.FormatInt(time.Now().Unix(), 10)
	if len(dir) > 0 {
		err = bucket.PutObject(dir+"/"+path, fileHandle)
		if err != nil {
			return "", err
		}
		fmt.Println("文件相对路径:", dir+"/"+path)
		url = "https://" + alioss.options.bucketName + "." + alioss.options.endpoint + "/" + dir + "/" + path
	} else {
		err = bucket.PutObject(file.Filename, fileHandle)
		if err != nil {
			return "", err
		}
		url = "https://" + alioss.options.bucketName + "." + alioss.options.endpoint + "/" + path
	}

	return url, nil
}

func (alioss *OSS) UploadDirectoryToOSS(ctx context.Context, localPath, ossDir string) ([]string, error) {
	var xzUrls []string // 存储所有.xz文件的URL

	// 获取OSS Bucket实例
	client, err := oss.New(alioss.options.endpoint, alioss.options.accessKey, alioss.options.secretKey)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(alioss.options.bucketName)
	if err != nil {
		return nil, err
	}

	// 生成一个统一的时间戳
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// 遍历目录，上传每个文件
	err = filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 可能遇到的错误
		}

		if !info.IsDir() {
			relativePath, err := filepath.Rel(filepath.Dir(localPath), path)
			if err != nil {
				return err
			}

			// 转换为使用正斜杠的路径，适用于OSS
			relativePath = timestamp + "/" + filepath.ToSlash(relativePath)
			if ossDir != "" {
				relativePath = ossDir + "/" + relativePath
			}
			// 上传文件至OSS
			uploadErr := alioss.uploadFileToOSS(ctx, bucket, path, relativePath)
			if uploadErr != nil {
				return uploadErr
			}
			if strings.HasSuffix(info.Name(), ".gltf.xz") {
				// 如果是.xz文件，构建URL
				xzUrl := fmt.Sprintf("https://%s.%s/%s", alioss.options.bucketName, alioss.options.endpoint, relativePath)
				xzUrls = append(xzUrls, xzUrl)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return xzUrls, nil
}

// uploadFileToOSS 将单个文件上传到指定的OSS地址
func (alioss *OSS) uploadFileToOSS(_ context.Context, bucket *oss.Bucket, localFilePath, ossFilePath string) error {
	file, err := os.Open(localFilePath)
	fmt.Println("文件为:", file)
	if err != nil {
		return err
	}
	defer file.Close()

	// 上传文件到OSS
	err = bucket.PutObject(ossFilePath, file)
	if err != nil {
		return err
	}
	fmt.Printf("文件上传成功: %s\n", ossFilePath)
	return nil
}

func (alioss *OSS) UnzipFile(zipPath, extractPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func(r *zip.ReadCloser) {
		_ = r.Close()
	}(r)

	for _, f := range r.File {
		fPath := filepath.Join(extractPath, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func(rc io.ReadCloser) {
			err := rc.Close()
			if err != nil {

			}
		}(rc)

		_, err = io.Copy(outFile, rc)

		if err != nil {
			return err
		}
	}
	return nil
}

func (alioss *OSS) ClearDirectory(dir string) error {
	// 读取目录中的所有文件和子目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// 遍历所有文件和子目录
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		// 如果是目录，则递归调用删除
		if entry.IsDir() {
			if err := alioss.ClearDirectory(path); err != nil {
				return err
			}
			// 删除空目录
			if err := os.Remove(path); err != nil {
				return err
			}
		} else {
			// 如果是文件，则直接删除
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}
