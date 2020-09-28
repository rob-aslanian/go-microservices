package advert

import (
	"errors"
	"strings"
	"time"

	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	durationDays = 30
)

// Advert represents statistic for Ads Manager
type Advert struct {
	ID             primitive.ObjectID `bson:"_id"`
	Type           Type               `bson:"type"`
	Format         Format             `bson:"format"`
	Name           string             `bson:"name"`
	Button         string             `bson:"button"`
	Status         Status             `bson:"status"`
	Budget         float32            `bson:"-"`
	Currency       string             `bson:"currency"`
	StartDate      time.Time          `bson:"start_date"`
	FinishDate     *time.Time         `bson:"finish_date,omitempty"`
	PausedTime     *time.Duration     `bson:"paused_time"`
	LastPausedTime *time.Time         `bson:"last_paused_time,omitempty"`
	CreatorID      primitive.ObjectID `bson:"creator_id"`

	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	Locations []location.Location `bson:"location"`
}

// GetID returns id of ad
func (ad Advert) GetID() string {
	return ad.ID.Hex()
}

// SetID ...
func (ad *Advert) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.ID = objID
	return nil
}

// GenerateID creates new id
func (ad *Advert) GenerateID() string {
	ad.ID = primitive.NewObjectID()
	return ad.ID.Hex()
}

// SetCreatorID ...
func (ad *Advert) SetCreatorID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.CreatorID = objID
	return nil
}

// GetCreatorID ...
func (ad Advert) GetCreatorID() string {
	return ad.CreatorID.Hex()
}

// SetCompanyID ...
func (ad *Advert) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ad.CompanyID = &objID
	return nil
}

// GetCompanyID ...
func (ad Advert) GetCompanyID() string {
	if ad.CompanyID == nil {
		return ""
	}
	return ad.CompanyID.Hex()
}

// CalculateDates set statuts according start date and calculate finish date
func (ad *Advert) CalculateDates() {
	if isToday(ad.StartDate) {
		ad.Status = StatusActive
		finishDate := ad.StartDate.AddDate(0, 0, durationDays)
		ad.FinishDate = &finishDate
	} else {
		ad.Status = StatusNotRunning
	}
}

func isToday(startDate time.Time) bool {
	year, month, day := startDate.Date()
	yearToday, monthToday, dayToday := time.Now().Date()

	if year != yearToday ||
		month != monthToday ||
		day != dayToday {
		return false
	}

	return true
}

// Validate checks budget, currency, country id, ....
func (ad Advert) Validate() error {
	if ad.Currency == "" {
		return errors.New("wrong_currency")
	}
	if len(ad.Locations) == 0 {
		return errors.New("wrong_country")
	}
	for _, l := range ad.Locations {
		if l.CountryID == "" {
			return errors.New("wrong_country")
		}
	}
	if ad.Name == "" {
		return errors.New("empty_name")
	}

	// check if date not in the past
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	if ad.StartDate.Before(today) {
		return errors.New("wrong_start_date")
	}

	return nil
}

// Trim removes spaces from Name
func (ad *Advert) Trim() {
	ad.Name = strings.TrimSpace(ad.Name)
}
