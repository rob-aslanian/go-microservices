package profile

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

type PortfolioSharedStruct struct {
	ID              bson.ObjectId `bson:"id"`
	Title           string        `bson:"title"`
	Description     string        `bson:"description"`
	IsCommentClosed bool          `bson:"is_comment_closed"`
	ViewsCount      int           `bson:"views_count"`
	LikesCount      int           `bson:"likes_count"`
	SharedCount     int           `bson:"shared_count"`
	SavedCount      int           `bson:"saved_count"`
	CreatedAt       time.Time     `bson:"create_at"`
}

// Portfolio . it contains file, title and info about the office
type Portfolio struct {
	ID              bson.ObjectId `bson:"id"`
	Title           string        `bson:"title"`
	Description     string        `bson:"description"`
	IsCommentClosed bool          `bson:"is_comment_closed"`
	ViewsCount      int32         `bson:"views_count"`
	LikesCount      int32         `bson:"likes_count"`
	SharedCount     int32         `bson:"shared_count"`
	SavedCount      int32         `bson:"saved_count"`
	Tools           []*string     `bson:"tools"`
	Files           []*File       `bson:"files"`
	ContentType     ContentType   `bson:"content_type"`
	CreatedAt       time.Time     `bson:"create_at"`
	HasLiked 	    bool

	Translations map[string]*PortfolioTranslation `bson:"translations"`
}

type Portfolios struct {
	PortfolioAmount int32 `bson:"portfolio_amount"`
	Portfolios []*Portfolio `bson:"portfolios"`
}

type PortfolioAction struct {
	PortfolioID string 
	OwnerID 	string 			`bson:"owner_id"`
	ProfileID   string 			`bson:"profile_id"`
	IsCompany 	bool 			`bson:"is_company"`
	CreatedAt   time.Time     	`bson:"create_at"`
}

// PortfolioLikes ... 
type PortfolioLikes struct {
	Likes    int32 	`bson:"likes"`
	HasLiked bool 	`bson:"has_liked"`
}

// PortfolioComment ...
type PortfolioComment struct {
	ID          bson.ObjectId 	`bson:"id"`
	PortfolioID string 
	Comment 	string 			`bson:"comment"`
	OwnerID 	string 			`bson:"owner_id"`
	ProfileID   string 			`bson:"profile_id"`
	IsCompany 	bool 			`bson:"is_company"`
	CreatedAt   time.Time     	`bson:"create_at"`
}

type GetPortfolioComment struct {
	ID          bson.ObjectId 	`bson:"id"`
	Comment 	string 			`bson:"comment"`
	OwnerID 	bson.ObjectId 	`bson:"owner_id"`
	ProfileID   bson.ObjectId 	`bson:"profile_id"`
	IsCompany 	bool 			`bson:"is_company"`
	CreatedAt   time.Time     	`bson:"created_at"`
}


// PortfolioInfo ... 
type PortfolioInfo struct{
	ViewsCount  	int32 			
	LikesCount  	int32 			
	CommentCount  	int32 		
	HasPhoto 		bool  `bson:"has_photo"`
	HasVideo	 	bool  `bson:"has_video"`  
	HasArticle 		bool  `bson:"has_article"`  
	HasAudio	 	bool  `bson:"has_audio"`  


}

type GetPortfolioComments struct {
	CommentAmount int32   			`bson:"comments_amount"`
	Comments []*GetPortfolioComment `bson:"comments"`
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
	ContentTypePhoto ContentType = "photo"
	// ContentTypeArticle ...
	ContentTypeArticle ContentType = "article"
	// ContentTypeVideo ...
	ContentTypeVideo ContentType = "video"
	// ContentTypeAudio ...
	ContentTypeAudio ContentType = "audio"
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

// GetCommentID ...
func (p *PortfolioComment) GetCommentID() string {
	return p.ID.Hex()
}

// SetCommentID set id
func (p *PortfolioComment) SetCommentID(id string) error {
	if ok := bson.IsObjectIdHex(id); !ok {
		return errors.New(`wrong_id`)
	}
	objID := bson.ObjectIdHex(id)

	p.ID = objID
	return nil

}

// GenerateCommentID creates new id
func (p *PortfolioComment) GenerateCommentID() string {
	p.ID = bson.NewObjectId()
	return p.ID.Hex()
}

// Translate ...
// func (i *Portfolio) Translate(lang string) {
// 	if i == nil || lang == "" {
// 		return
// 	}

// 	if tr, isExists := i.Translations[lang]; isExists {
// 		i.Tittle = tr.Tittle
// 		i.Description = tr.Description
// 	}
// }
