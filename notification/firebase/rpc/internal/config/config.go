package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

// Config ... struct to load configurations store in etc folder
// you can append more fields that you need.
type Config struct {
	zrpc.RpcServerConf

	DB struct { // database config
		DataSource     string `json:"datasource"`
		DatabaseName   string `json:"databaseName"`
		CollectionName string `json:"collectionName"`
	}

	StorageProvider struct {
		DigitalOcean struct {
			SpaceName     string
			SpaceEndpoint string
			SpaceKey      string
			SpaceSecret   string
			SpaceRegion   string
		}
	}
}
