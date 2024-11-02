package awss3

import "time"

type MultipartResponse struct {
	UploadId string
}

type MultipartUploadPreSignConfig struct {
	Bucket     string        // 存储桶名称
	UploadId   string        // 上传ID
	ObjectKey  string        // 对象名称
	PartNumber int32         // 分片编号
	Expires    time.Duration // 预签名URL过期时间
}

type CompleteMultipartUploadConfig struct {
	bucket      string
	UploadId    string         // 上传ID
	ObjectKey   string         // 对象名称
	ContentType string         // 文件类型
	Parts       []CompletePart // 分片列表
}

type CompletePart struct {
	PartNumber int32  // 分片编号
	ETag       string // 分片的ETag
}
