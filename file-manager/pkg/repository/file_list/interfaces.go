package fileList

import "time"

type (
	Configuration interface {
		GetAddress() []string
		GetUser() string
		GetPassword() string
		GetDatabase() string
		GetCollectionFiles() string
	}

	FileInfo = interface {
		GetID() string
		GetURL() string
		GetName() string
		GetInternalName() string
		GetMimeType() string
		GetSize() int64
		GetOwnerID() string
		GetCreatedAt() time.Time
	}
)
