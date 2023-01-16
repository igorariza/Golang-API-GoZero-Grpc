package model

import (
	"context"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
)

// ListByPage ... List files with pagination
func (m *customFilesModel) ListByPage(page, size int32) (files []*Files, err error) {
	// records to skip in query
	var skip int64
	// parse size page to int64
	sizeInt64 := int64(size)

	if page > 0 {
		skip = int64(page-1) * sizeInt64
	}

	options := options2.FindOptions{
		Skip:  &skip,
		Limit: &sizeInt64,
	}

	err = m.conn.Find(context.Background(), &files, options)
	if err != nil {
		return files, err
	}

	return nil, nil
}
