package storage

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/MessageDream/goby/module/infrastructure"
)

func NewLocal(config *Config) *Local {
	if err := os.MkdirAll(config.LocalStoragePath, os.ModePerm); err != nil {
		log.Fatal(4, "Fail to create local storage directory '%s': %v", path.Dir(config.LocalStoragePath), err)
	}
	return &Local{config}
}

type Local struct {
	Config *Config
}

func (self *Local) Upload(key, filePath string) (string, error) {
	fileName := key + path.Ext(filePath)
	savePath := path.Join(self.Config.LocalStoragePath, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	w, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer w.Close()
	_, err = io.Copy(w, file)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (self *Local) Download(key string) (string, error) {
	savedPath := path.Join(self.Config.LocalStoragePath, key)
	file, err := os.Open(savedPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return infrastructure.SaveFileToTemp(key, file)
}

func (self *Local) GenerateDownloadURL(key string) string {
	return path.Join(self.Config.DownloadURL, key)
}
