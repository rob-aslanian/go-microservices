package office

import (
	"time"
	
	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/location"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/qualifications"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/review"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Office is the v-office for seller to sell his services.
type Office struct {
	ID                    primitive.ObjectID            `bson:"_id"`
	UserID                *primitive.ObjectID           `bson:"user_id"`
	CompanyID             *primitive.ObjectID           `bson:"company_id"`
	SellerDashBoardID     primitive.ObjectID            `bson:"seller_dashboard_id"`
	Name                  string                        `bson:"name"`
	IsOut                 bool                          `bson:"is_out"`
	ReturnDate            *time.Time                    `bson:"return_date"`
	Category              string      					`bson:"category"`
	Location              location.Location             `bson:"location"`
	Portfolio             []*Portfolio                  `bson:"portfolio,omitempty"`
	Description           string                        `bson:"description"`
	Reviews               *[]review.Review              `bson:"reviews"`
	Files                 []*file.File                  `bson:"files"`
	CoverImage            string                        `bson:"cover_image,omitempty"`
	CoverOriginImage      string                        `bson:"cover_origin_image,omitempty"`
	Languages             []*qualifications.Language    `bson:"languages"`
	CreatedAt             time.Time                     `bson:"created_at"`
	Services              []*servicerequest.Service
	Translation           map[string]*Translation `bson:"translation"`
	IsMe                  bool                    `bson:"-"`
	CurrentTranslation    string                  `bson:"-"`
	AvailableTranslations []string                `bson:"-"`
}

// Translation ...
type Translation struct {
	Name        string `bson:"name"`
	Description string `bson:"description"`
}

// GetID returns id
func (p Office) GetID() string {
	if !p.ID.IsZero() {
		return p.ID.Hex()
	}
	return ""
}

// SetID set id
func (p *Office) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Office) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns user id
func (p Office) GetUserID() string {
	if p.UserID == nil {
		return ""
	}

	if !p.UserID.IsZero() {
		return p.UserID.Hex()
	}
	return ""
}

// SetUserID set user id
func (p *Office) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = &objID
	return nil
}

// GetCompanyID returns user id
func (p Office) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}

	if !p.CompanyID.IsZero() {
		return p.CompanyID.Hex()
	}
	return ""
}

// SetCompanyID set user id
func (p *Office) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// GetSellerDashBoardID returns user id
func (p Office) GetSellerDashBoardID() string {
	return p.SellerDashBoardID.Hex()
}

// SetSellerDashBoardID set user id
func (p *Office) SetSellerDashBoardID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}
