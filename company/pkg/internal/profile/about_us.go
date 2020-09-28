package profile

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
)

// AboutUs ...
type AboutUs struct {
	// ID             bson.ObjectId           `bson:"id"`
	Size           account.Size            `bson:"size"`
	Type           account.Type            `bson:"type"`
	Mission        string                  `bson:"mission"`
	Parking        account.Parking         `bson:"parking"`
	Industry       account.Industry        `bson:"industry"`
	Description    string                  `bson:"description"`
	BusinessHours  []*account.BusinessHour `bson:"business_hours"`
	FoundationDate time.Time               `bson:"foundation_date"`
	IsAboutUsSet   bool                    `bson:"changed_about_us"`

	Translations map[string]*AboutUsTranslation `bson:"translations"`

	IsDescriptionNull bool `bson:"-"`
	IsMissionNull     bool `bson:"-"`
	IsTypeNull        bool `bson:"-"`
	IsSizeNull        bool `bson:"-"`
	IsParkingNull     bool `bson:"-"`
	IsSubindustryNull bool `bson:"-"`
}

// AboutUsTranslation ...
type AboutUsTranslation struct {
	Mission     string `bson:"mission"`
	Description string `bson:"description"`
}

// Translate ...
func (a *AboutUs) Translate(ctx context.Context, lang string) string {
	if a == nil || lang == "" {
		return "en"
	}

	if tr, isExists := a.Translations[lang]; isExists {
		a.Mission = tr.Mission
		a.Description = tr.Description
	} else {
		return "en"
	}

	return lang
}

// // GetID returns id of company profile
// func (a AboutUs) GetID() string {
// 	return a.ID.Hex()
// }
//
// // SetID saves id of company profile. If id has a wrong format returns usersErrors.ErrWrongID error.
// func (a *AboutUs) SetID(id string) error {
// 	if bson.IsObjectIdHex(id) {
// 		a.ID = bson.ObjectIdHex(id)
// 		return nil
// 	}
// 	return companyErrors.ErrWrongID
// }
