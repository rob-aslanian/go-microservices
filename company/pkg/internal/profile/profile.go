package profile

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	careercenter "gitlab.lan/Rightnao-site/microservices/company/pkg/internal/career-center"
)

// Profile ...
type Profile struct {
	account.Account `bson:",inline"`
	Avatar          string     `bson:"avatar"`
	Cover           string     `bson:"cover"`
	Description     string     `bson:"description"`
	Mission         string     `bson:"mission"`
	Founders        []*Founder `bson:"founders"`
	Awards          []*Award   `bson:"awards"`
	// Gallery         []*Gallery   `bson:"gallery"`
	Milestones []*Milestone `bson:"milestones"`
	Benefits   []Benefit    `bson:"benefits"`
	Services   []*Service   `bson:"services"` // where to store it?
	Products   []*Product   `bson:"products"` // where to store it?
	// Reports []*Report `bson:"reports"` // should be seperate

	Translation           map[string]*Translation `bson:"translation"`
	CurrentTranslation    string                  `bson:"-"`
	AvailableTranslations []string                `bson:"-"`

	AmountOfFollowings int32   `bson:"-"`
	AmountOfFollowers  int32   `bson:"-"`
	AmountOfEmployees  int32   `bson:"-"`
	AmountOfJobs       int32   `bson:"-"`
	AvarageRating      float32 `bson:"-"`

	IsAboutUsSet bool `bson:"changed_about_us"`

	CareerCenter *careercenter.CareerCenter `bson:"career_center"`

	IsFollow   bool `bson:"-"`
	IsFavorite bool `bson:"-"`
	IsOnline   bool `bson:"-"`
	IsBlocked  bool `bson:"-"`
}

// Translation ...
type Translation struct {
	Name        string `bson:"name"`
	Mission     string `bson:"mission"`
	Description string `bson:"description"`
}

// // GetID returns id of company profile
// func (a Profile) GetID() string {
// 	return a.ID.Hex()
// }
//
// // SetID saves id of company profile. If id has a wrong format returns usersErrors.ErrWrongID error.
// func (a *Profile) SetID(id string) error {
// 	if bson.IsObjectIdHex(id) {
// 		a.ID = bson.ObjectIdHex(id)
// 		return nil
// 	}
// 	return companyErrors.ErrWrongID
// }
//
// // GenerateID creates new random id for company profile
// func (a *Profile) GenerateID() string {
// 	id := bson.NewObjectId()
// 	a.ID = id
// 	return id.Hex()
// }

// Translate ...
func (p *Profile) Translate(ctx context.Context, lang string) string {
	if p == nil || lang == "" {
		return "en"
	}

	if tr, isExists := p.Translation[lang]; isExists {
		p.Name = tr.Name
		p.Mission = tr.Mission
		p.Description = tr.Description
	} else {
		return "en"
	}

	// milestone
	if lang != "" {
		for i := range p.Milestones {
			if p.Milestones[i] != nil {
				p.Milestones[i].Translate(ctx, lang)
			}
		}
	}

	// awards
	if lang != "" {
		for i := range p.Awards {
			if p.Awards[i] != nil {
				p.Awards[i].Translate(ctx, lang)
			}
		}
	}

	return lang
}
