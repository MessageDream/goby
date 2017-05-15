package storage

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/MessageDream/goby/module/infrastructure"

	"github.com/denverdino/aliyungo/oss"
)

func NewOSS(config *Config) *OSS {
	region := oss.Region(config.OSSEndpoint)
	client := oss.NewOSSClient(region, false, config.AccessKey, config.SecretKey, true)
	resp, err := client.GetService()
	if err != nil {
		log.Fatal(4, "Fail to get oss service'%v': %v", config, err)
	}

	bucket := client.Bucket(config.Bucket)

	bucketExist := false
	for _, b := range resp.Buckets {
		if b.Name == bucket.Name {
			bucketExist = true
		}
	}

	if !bucketExist {
		if err := bucket.PutBucket(oss.PublicRead); err != nil {
			log.Fatal(4, "Fail to create bucket:'%v': %v", config, err)
		}
	}

	return &OSS{
		Config: config,
		Bucket: bucket,
	}
}

type OSS struct {
	Config *Config
	Bucket *oss.Bucket
}

func (self *OSS) Upload(key, filePath string) (string, error) {
	fileName := key + path.Ext(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileRemotePath := path.Join(self.Config.Prefix, fileName)
	if err := self.Bucket.PutFile(fileRemotePath, file, oss.PublicRead, oss.Options{}); err != nil {
		return "", nil
	}
	return fileName, nil
}

func (self *OSS) Download(key string) (string, error) {
	headers := http.Header{}
	fileRemotePath := path.Join(self.Config.Prefix, key)
	resp, err := self.Bucket.GetResponseWithHeaders(fileRemotePath, headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return infrastructure.SaveFileToTemp(key, resp.Body)
}

func (self *OSS) GenerateDownloadURL(key string) string {
	fileRemotePath := path.Join(self.Config.Prefix, key)
	return self.Bucket.URL(fileRemotePath)
}
