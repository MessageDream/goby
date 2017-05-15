package storage

type Config struct {
	DownloadURL string

	LocalStoragePath string

	AccessKey string
	SecretKey string
	Bucket    string
	Prefix    string

	QiNiuZone int

	OSSEndpoint string
}

type Storage interface {
	Upload(key, filePath string) (string, error)

	Download(key string) (string, error)

	GenerateDownloadURL(key string) string
}

var (
	storage Storage
)

func InitStorage(storageType string, config *Config) {

	switch storageType {
	case "qiniu":
		storage = NewQiNiu(config)
	case "oss":
		storage = NewOSS(config)
	default:
		storage = NewLocal(config)
	}
}

func Upload(key, filePath string) (string, error) {
	return storage.Upload(key, filePath)
}

func Download(key string) (string, error) {
	return storage.Download(key)
}

func GenerateDownloadURL(key string) string {
	return storage.GenerateDownloadURL(key)
}
