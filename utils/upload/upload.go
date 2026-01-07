/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package upload

import (
	"github.com/google/wire"
	"mime/multipart"
	"github.com/sufo/bailu-admin/app/config"
)

type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	//仅支持local存储
	UploadFileToDir(file *multipart.FileHeader, subDir string) (string, string, error)
	DeleteFile(key ...string) error
	//仅支持local存储
	DeleteFileInDir(key string, subDir string) error
}

// NewOss OSS的实例化方法
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [ccfish86](https://github.com/ccfish86)
func NewOSS() OSS {
	switch config.Conf.Upload.Type {
	case "local":
		wire.Bind(new(OSS), new(Local))
		return &Local{}
	//case "qiniu":
	//	return &Qiniu{}
	//case "tencent-cos":
	//	return &TencentCOS{}
	//case "aliyun-oss":
	//	return &AliyunOSS{}
	//case "huawei-obs":
	//	return HuaWeiObs
	//case "aws-s3":
	//	return &AwsS3{}
	default:
		return &Local{}
	}
}
