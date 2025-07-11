package util

import (
	"Douyin/config"
	"bytes"
	"context"
	"fmt"
	"strings"

	// TODO 到时候再修改
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func UploadVideo(data []byte) (VideoUrl string, err error) {
	config.InitConfig()
	size := int64(len(data))
	key := fmt.Sprintf("%s.mp4", GenerateUUID())
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", config.Config.AliyunConfig.Bucket, key),
	}
	mac := qbox.NewMac(config.Config.AliyunConfig.AccessKey, config.Config.AliyunConfig.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	uploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err = uploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), size, &putExtra)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", config.Config.AliyunConfig.Domain, ret.Key), nil
}

func UploadJPG(imgPath string, videoUrl string) string {
	config.InitConfig()

	videoName := strings.Split(strings.Replace(videoUrl, config.Config.AliyunConfig.Domain+"/", "", -1), ".")[0]
	key := fmt.Sprintf("%s.%s", videoName+"_cover", "jpg")

	putPolicy := storage.PutPolicy{
		Scope: config.Config.AliyunConfig.Bucket,
	}
	mac := qbox.NewMac(config.Config.AliyunConfig.AccessKey, config.Config.AliyunConfig.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, imgPath, &putExtra)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return fmt.Sprintf("%s/%s", config.Config.AliyunConfig.Domain, ret.Key)
}
