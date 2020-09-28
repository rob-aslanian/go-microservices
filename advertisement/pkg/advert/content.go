package advert

import (
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/file"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Content ...
type Content struct {
	ID             primitive.ObjectID `bson:"_id"`
	ImageURL       string             `bson:"image_url"`
	Title          string             `bson:"title"`
	Content        string             `bson:"content,omitempty"`
	CustomButton   string             `bson:"custom_button,omitempty"`
	DestinationURL string             `bson:"destination_url,omitempty"`
	Files          []file.File        `bson:"files,omitempty"`
}

// GetID returns id of ad
func (ad Content) GetID() string {
	return ad.ID.Hex()
}

// SetID ...
func (ad *Content) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.ID = objID
	return nil
}

// GenerateID creates new id
func (ad *Content) GenerateID() string {
	ad.ID = primitive.NewObjectID()
	return ad.ID.Hex()
}

// // GetFileID returns id of file
// func (c Content) GetFileID() string {
// 	return c.FileID.Hex()
// }

// // SetFileID ...
// func (c *Content) SetFileID(id string) error {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	c.FileID = objID
// 	return nil
// }
