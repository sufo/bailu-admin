package upload

type Qiniu struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Qiniu
//@function: UploadFile
//@description: 上传文件
//@param: file *multipart.FileHeader
//@return: string, string, error

//func (*Qiniu) UploadFile(file *multipart.FileHeader) (string, string, error) {
//	putPolicy := storage.PutPolicy{Scope: config.Conf.Qiniu.Bucket}
//	mac := qbox.NewMac(config.Conf.Qiniu.AccessKey, config.Conf.Qiniu.SecretKey)
//	upToken := putPolicy.UploadToken(mac)
//	cfg := qiniuConfig()
//	formUploader := storage.NewFormUploader(cfg)
//	ret := storage.PutRet{}
//	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}
//
//	f, openError := file.Open()
//	if openError != nil {
//		log.L.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
//
//		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
//	}
//	defer f.Close()                                                  // 创建文件 defer 关闭
//	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
//	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
//	if putErr != nil {
//		log.L.Error("function formUploader.Put() Filed", zap.Any("err", putErr.Error()))
//		return "", "", errors.New("function formUploader.Put() Filed, err:" + putErr.Error())
//	}
//	return config.Conf.Qiniu.ImgPath + "/" + ret.Key, ret.Key, nil
//}
//
//// 不支持subDir
//func (q *Qiniu) UploadFileToDir(file *multipart.FileHeader, subDir string) (string, string, error) {
//	log.L.Warn("do not support subDir")
//	return q.UploadFile(file)
//}
//
////@author: [piexlmax](https://github.com/piexlmax)
////@author: [ccfish86](https://github.com/ccfish86)
////@author: [SliverHorn](https://github.com/SliverHorn)
////@object: *Qiniu
////@function: DeleteFile
////@description: 删除文件
////@param: key string
////@return: error
//
//func (*Qiniu) DeleteFile(key string) error {
//	mac := qbox.NewMac(config.Conf.Qiniu.AccessKey, config.Conf.Qiniu.SecretKey)
//	cfg := qiniuConfig()
//	bucketManager := storage.NewBucketManager(mac, cfg)
//	if err := bucketManager.Delete(config.Conf.Qiniu.Bucket, key); err != nil {
//		global.GVA_LOG.Error("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
//		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
//	}
//	return nil
//}
//
//func (q *Qiniu) DeleteFileInDir(key string, subDir string) error {
//	log.L.Warn("do not support subDir")
//	return q.DeleteFile(key)
//}
//
////@author: [SliverHorn](https://github.com/SliverHorn)
////@object: *Qiniu
////@function: qiniuConfig
////@description: 根据配置文件进行返回七牛云的配置
////@return: *storage.Config
//
//func qiniuConfig() *storage.Config {
//	cfg := storage.Config{
//		UseHTTPS:      config.Conf.Qiniu.UseHTTPS,
//		UseCdnDomains: config.Conf.Qiniu.UseCdnDomains,
//	}
//	switch config.Conf.Qiniu.Region { // 根据配置文件进行初始化空间对应的机房
//	case "ZoneHuadong":
//		cfg.Region = &storage.ZoneHuadong
//	case "ZoneHuabei":
//		cfg.Region = &storage.ZoneHuabei
//	case "ZoneHuanan":
//		cfg.Region = &storage.ZoneHuanan
//	case "ZoneBeimei":
//		cfg.Region = &storage.ZoneBeimei
//	case "ZoneXinjiapo":
//		cfg.Region = &storage.ZoneXinjiapo
//	}
//	return &cfg
//}
