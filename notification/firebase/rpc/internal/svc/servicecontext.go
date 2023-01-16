package svc

import (
	"os"

	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/model"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/pkg/storageprovider"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type ServiceContext struct {
	Config                config.Config
	FileModel             model.FilesModel
	StorageProviderClient storageprovider.IStorageProvider
}

func NewServiceContext(c config.Config) *ServiceContext {
	var err error
	//configure cache with redis
	cacheConfig := cache.CacheConf{
		cache.NodeConf{
			RedisConf: c.Redis.RedisConf,
			Weight:    100,
		},
	}
	// instance model
	fileModel := model.NewFilesModel(c.DB.DataSource, c.DB.DatabaseName, c.DB.CollectionName, cacheConfig)

	// instance storage provider
	var storageProviderClient storageprovider.IStorageProvider
	if c.StorageProvider.DigitalOcean.SpaceName != "" {
		if err = os.Setenv("DO_SPACES_NAME", c.StorageProvider.DigitalOcean.SpaceName); err != nil {
			panic(err)
		}
		if err = os.Setenv("DO_SPACES_ENDPOINT", c.StorageProvider.DigitalOcean.SpaceEndpoint); err != nil {
			panic(err)
		}
		if err = os.Setenv("DO_SPACES_KEY", c.StorageProvider.DigitalOcean.SpaceKey); err != nil {
			panic(err)
		}
		if err = os.Setenv("DO_SPACES_SECRET", c.StorageProvider.DigitalOcean.SpaceSecret); err != nil {
			panic(err)
		}
		if err = os.Setenv("DO_SPACES_REGION", c.StorageProvider.DigitalOcean.SpaceRegion); err != nil {
			panic(err)
		}

		storageProviderClient, err = storageprovider.NewStorageProvider(storageprovider.DigitalOceanSpace)
		if err != nil {
			panic(err)
		}
	}

	//configure service context
	return &ServiceContext{
		Config:                c,
		FileModel:             fileModel,
		StorageProviderClient: storageProviderClient,
	}
}
