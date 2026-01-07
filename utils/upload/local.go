package upload

import (
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils"
	"errors"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type Local struct{}

func (local *Local) UploadFileToDir(file *multipart.FileHeader, subDir string) (string, string, error) {
	//fmt.Printf("%v", file)
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	//读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5(name)

	// 拼接新文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	//建议可以默认按日期分割
	mkdirErr := os.MkdirAll(config.Conf.Local.Dir+"/"+subDir, os.ModePerm)

	if mkdirErr != nil {
		log.L.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	// 拼接路径和文件名
	if subDir != "" && !strings.HasSuffix(subDir, "/") {
		subDir = subDir + "/"
	}
	filePath := config.Conf.Local.Dir + "/" + subDir + filename
	//filepath := global.GVA_CONFIG.Local.Path + "/" + filename

	f, openError := file.Open() // 读取文件
	if openError != nil {
		log.L.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(filePath)
	if createErr != nil {
		log.L.Error("function os.Create() Filed", zap.Any("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		log.L.Error("function io.Copy() Filed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	if subDir != "" && !strings.HasSuffix(subDir, "/") {
		subDir = subDir + "/"
	}
	return subDir + filename, filename, nil
}

func (local *Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	return local.UploadFileToDir(file, "")
}

//@function: DeleteFile
//@description: 删除文件
//@param: key string
//@return: error

func (*Local) DeleteFile(keys ...string) error {
	for _, key := range keys {
		p := config.Conf.Local.Dir + "/" + key
		if strings.Contains(p, config.Conf.Local.Dir) {
			if err := os.Remove(p); err != nil {
				return errors.New("本地文件删除失败, err:" + err.Error())
			}
		}
	}
	return nil
}

func (*Local) DeleteFileInDir(key string, subDir string) error {
	if subDir != "" && !strings.HasSuffix(subDir, "/") {
		subDir = subDir + "/"
	}

	p := config.Conf.Local.Dir + "/" + subDir + key
	if strings.Contains(p, config.Conf.Local.Dir+"/"+subDir) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
