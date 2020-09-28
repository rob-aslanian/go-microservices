package fileList

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type FileInfoModel struct {
	ID           bson.ObjectId `bson:"_id"`
	URL          string        `bson:"url"`
	InternalName string        `bson:"internal_name"`
	Name         string        `bson:"name"`
	MimeType     string        `bson:"mime-type"`
	Size         int64         `bson:"size"`
	OwnerID      bson.ObjectId `bson:"owner"`
	CreatedAt    time.Time     `bson:"created_at"`
}

func (f *FileInfoModel) GetID() string {
	return f.ID.Hex()
}

func (f *FileInfoModel) GetURL() string {
	return f.URL
}

func (f *FileInfoModel) GetName() string {
	return f.Name
}

func (f *FileInfoModel) GetInternalName() string {
	return f.InternalName
}

func (f *FileInfoModel) GetMimeType() string {
	return f.MimeType
}

func (f *FileInfoModel) GetSize() int64 {
	return f.Size
}

func (f *FileInfoModel) GetOwnerID() string {
	return f.OwnerID.Hex()
}

func (f *FileInfoModel) GetCreatedAt() time.Time {
	return f.CreatedAt
}
