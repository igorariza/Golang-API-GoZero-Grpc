package storageprovider

import (
	"errors"
	"fmt"

	bylogs "github.com/cuemby/by-go-utils/pkg/bylogger"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/pkg/file"
)

type IStorageProvider interface {
	Connect() error
	ProviderName() string
	UploadFile(genericFile *file.GenericFile) error
	ListFiles() ([]*file.GenericFile, error)
	ListFilesPaginated(page int64, size int64) ([]*file.GenericFile, error)
	DeleteFile(filename string) (bool, error)
	GetFile(filename string) (*file.GenericFile, error)
}

type StorageProvider string

const (
	DigitalOceanSpace StorageProvider = "DO_SPACES"
	AwsS3             StorageProvider = "AWS_S3"
	GoogleStorage     StorageProvider = "GOOGLE_STORAGE"
)

func NewStorageProvider(provider StorageProvider) (IStorageProvider, error) {
	switch provider {
	case DigitalOceanSpace:
		client := &DigitalOceanSpaceProvider{}
		err := client.Connect()
		return client, err
	default:
		bylogs.LogErr("unsupported storage provider", provider)
		return nil, errors.New(fmt.Sprint("unsupported storage provider:", provider))
	}
}
