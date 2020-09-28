package servicerequest

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
)

// Portfolio . it contains file, title and info about the office
type Portfolio struct {
	ID          bson.ObjectId `bson:"id"`
	ContentType ContentType   `bson:"content_type"`
	Tittle      string        `bson:"tittle"`
	Description string        `bson:"description"`
	Links       []*file.Link  `bson:"links"`
	Files       []*file.File  `bson:"files"`
	CreatedAt   time.Time     `bson:"create_at"`

	Translations map[string]*PortfolioTranslation `bson:"translations"`
}

// PortfolioTranslation ...
type PortfolioTranslation struct {
	ContentType ContentType `bson:"content_type"`
	Tittle      string      `bson:"tittle"`
	Description string      `bson:"description"`
}

// ContentType to know which type of file user is uploading in v-office
type ContentType string

const (
	// ContentTypeImage ...
	ContentTypeImage ContentType = "image"
	// ContentTypeArticle ...
	ContentTypeArticle ContentType = "article"
	// ContentTypeCode ...
	ContentTypeCode ContentType = "code"
	// ContentTypeVideo ...
	ContentTypeVideo ContentType = "video"
	// ContentTypeAudio ...
	ContentTypeAudio ContentType = "audio"
	//ContentTypeOther ...
	ContentTypeOther ContentType = "other"
)

// GetID returns id
func (p Portfolio) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Portfolio) SetID(id string) error {
	if ok := bson.IsObjectIdHex(id); !ok {
		return errors.New(`wrong_id`)
	}
	objID := bson.ObjectIdHex(id)

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Portfolio) GenerateID() string {
	p.ID = bson.NewObjectId()
	return p.ID.Hex()
}

// Translate ...
func (i *Portfolio) Translate(lang string) {
	if i == nil || lang == "" {
		return
	}

	if tr, isExists := i.Translations[lang]; isExists {
		i.Tittle = tr.Tittle
		i.Description = tr.Description
	}
}
