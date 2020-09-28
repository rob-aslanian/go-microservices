package advert

import (
	"time"

	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/file"

	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Campaing ...
type Campaing struct {
	ID            primitive.ObjectID   `bson:"_id"`
	Type          Type                 `bson:"type"`
	Formats       []Format             `bson:"format"`
	Name          string               `bson:"name"`
	Status        Status               `bson:"status"`
	Gender        *string              `bson:"gender,omitempty"`
	AgeFrom       *int32               `bson:"age_from,omitempty"`
	AgeTo         *int32               `bson:"age_to,omitempty"`
	Currency      string               `bson:"currency"`
	Impressions   int32                `bson:"impressions"`
	Clicks        int32                `bson:"clicks"`
	Forwarding    int32                `bson:"forwarding"`
	Referals      int32                `bson:"referals"`
	StartDate     time.Time            `bson:"start_date"`
	FinishDate    *time.Time           `bson:"finish_date,omitempty"`
	CreatorID     primitive.ObjectID   `bson:"creator_id"`
	IsComapny     bool                 `bson:"is_company"`
	Locations     []location.Location  `bson:"location"`
	Languages     []string             `bson:"languages"`
	Adverts       []CampaingAdvert     `bson:"adverts"`
	RelevantUsers []primitive.ObjectID `bson:"relevant_users,omitempty"`
	CtrAVG        float64              `bson:"ctr_avg,omitempty"`
}

// Campaings ...
type Campaings struct {
	Amount    int32      `bson:"amount"`
	Campaings []Campaing `bson:"campaigns"`
}

// Adverts ...
type Adverts struct {
	Amount   int32            `bson:"amount"`
	Campaing Campaing         `bson:"campaign"`
	Adverts  []CampaingAdvert `bson:"adverts"`
}

// GetAdvert ...
type GetAdvert struct {
	Impressions int32          `bson:"impressions"`
	Clicks      int32          `bson:"clicks"`
	Adverts     CampaingAdvert `bson:"adverts"`
}

// CampaingAdvert ...
type CampaingAdvert struct {
	ID          primitive.ObjectID `bson:"_id"`
	TypeID      primitive.ObjectID `bson:"type_id,omitempty"`
	Status      Status             `bson:"status"`
	AdvertType  Type               `bson:"type"`
	Name        string             `bson:"name"`
	URL         string             `bson:"url,omitempty"`
	Impressions int32              `bson:"impressions"`
	Clicks      int32              `bson:"clicks"`
	Forwarding  int32              `bson:"forwarding"`
	CreatedAt   time.Time          `bson:"created_at"`
	Contents    *[]Content         `bson:"contents,omitempty"`
	Formats     []Format           `bson:"formats,omitempty"`
	Files       []file.File        `bson:"files,omitempty"`
}

// GetAdvertID returns id of ad
func (ad CampaingAdvert) GetAdvertID() string {
	return ad.ID.Hex()
}

// SetAdvertID ...
func (ad *CampaingAdvert) SetAdvertID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.ID = objID
	return nil
}

// GetTypeID returns id of ad
func (ad CampaingAdvert) GetTypeID() string {
	return ad.TypeID.Hex()
}

// SetTypeID ...
func (ad *CampaingAdvert) SetTypeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.TypeID = objID
	return nil
}

// GenerateAdvertID creates new id
func (ad *CampaingAdvert) GenerateAdvertID() string {
	ad.ID = primitive.NewObjectID()
	return ad.ID.Hex()
}

// GetID returns id of ad
func (ad Campaing) GetID() string {
	return ad.ID.Hex()
}

// SetID ...
func (ad *Campaing) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.ID = objID
	return nil
}

// GenerateID creates new id
func (ad *Campaing) GenerateID() string {
	ad.ID = primitive.NewObjectID()
	return ad.ID.Hex()
}

// SetCreatorID ...
func (ad *Campaing) SetCreatorID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.CreatorID = objID
	return nil
}

// GetCreatorID ...
func (ad Campaing) GetCreatorID() string {
	return ad.CreatorID.Hex()
}
