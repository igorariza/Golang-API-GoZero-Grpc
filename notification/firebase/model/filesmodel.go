package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
)

var _ FilesModel = (*customFilesModel)(nil)

type (
	// FilesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFilesModel.
	FilesModel interface {
		filesModel

		// add more functions to model for use in logic files

		// ListByPage ... list files with pagination
		ListByPage(page, size int32) ([]*Files, error)
	}

	customFilesModel struct {
		*defaultFilesModel
	}
)

// NewFilesModel returns a model for the mongo.
func NewFilesModel(url, db, collection string, c cache.CacheConf) FilesModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customFilesModel{
		defaultFilesModel: newDefaultFilesModel(conn),
	}
}
