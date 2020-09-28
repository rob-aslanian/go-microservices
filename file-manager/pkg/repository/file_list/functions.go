package fileList

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/globalsign/mgo/bson"
)

func (r *repository) SaveFileInfo(header multipart.FileHeader, userID string, url string, internalName string) error {
	return r.collection.Insert(
		bson.M{
			"_id":           bson.NewObjectId(),
			"url":           url,
			"internal_name": internalName,
			"name":          header.Filename,
			"mime-type":     header.Header.Get("Content-Type"),
			"size":          header.Size,
			"owner":         bson.ObjectIdHex(userID),
			"created_at":    time.Now().UTC(),
		},
	)
}

func (r *repository) GetFileInfo(ctx context.Context, url string) (FileInfo, error) {
	var result FileInfoModel

	err := r.collection.Find(
		bson.M{
			"url": url,
		},
	).One(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
