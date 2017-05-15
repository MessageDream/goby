package storage

import (
	"net/http"
	"path"

	"github.com/MessageDream/goby/module/infrastructure"

	"golang.org/x/net/context"
	"qiniupkg.com/api.v7/kodo"
)

func NewQiNiu(config *Config) *QiNiu {

	kodo.SetMac(config.AccessKey, config.SecretKey)

	client := kodo.New(config.QiNiuZone, nil)

	return &QiNiu{
		Config: config,
		Bucket: client.Bucket(config.Bucket),
	}
}

type QiNiu struct {
	Config *Config
	Bucket kodo.Bucket
}

func (self *QiNiu) Upload(key, filePath string) (string, error) {
	fileName := key + path.Ext(filePath)
	fileRemotePath := path.Join(self.Config.Prefix, fileName)
	ctx := context.Background()
	if err := self.Bucket.PutFile(ctx, nil, fileRemotePath, filePath, nil); err != nil {
		return "", err
	}
	return fileName, nil
}

func (self *QiNiu) Download(key string) (string, error) {
	fileRemotePath := path.Join(self.Config.Prefix, key)
	url := kodo.MakeBaseUrl(self.Config.DownloadURL, fileRemotePath)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return infrastructure.SaveFileToTemp(key, resp.Body)

}

func (self *QiNiu) GenerateDownloadURL(key string) string {
	fileRemotePath := path.Join(self.Config.Prefix, key)
	return kodo.MakeBaseUrl(self.Config.DownloadURL, fileRemotePath)
}
