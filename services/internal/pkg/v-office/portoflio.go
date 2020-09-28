package office

import (
	"time"

	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Portfolio is for v-office only. it contains file, title and info about the office
type Portfolio struct {
	ID          primitive.ObjectID  `bson:"_id"`
	UserID      primitive.ObjectID  `bson:"user_id"`
	CompanyID   *primitive.ObjectID `bson:"company_id,omitempty"`
	OfficeID    primitive.ObjectID  `bson:"office_id"`
	ContentType ContentType         `bson:"content_type"`
	Tittle      string              `bson:"tittle"`
	Description string              `bson:"description"`
	Files       []*file.File        `bson:"files"`
	Link        []*file.Link        `bson:"links"`
	CreatedAt   time.Time           `bson:"create_at"`
}

// GetID returns id
func (p Portfolio) GetID() string {
	if !p.ID.IsZero() {
		return p.ID.Hex()
	}
	return ""
}

// SetID set id
func (p *Portfolio) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Portfolio) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// func (p *Portfolio) GenerateLinkID() string {
// 	for _, link := range p.Link {
// 		link.ID = primitive.NewObjectID()
// 	}
// }

// GetUserID returns user id
func (p Portfolio) GetUserID() string {
	if !p.UserID.IsZero() {
		return p.UserID.Hex()
	}
	return ""
}

// SetUserID set user id
func (p *Portfolio) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// GetCompanyID returns user id
func (p Portfolio) GetCompanyID() string {
	if p.CompanyID != nil {
		return p.CompanyID.Hex()
	}
	return ""
}

// SetCompanyID set user id
func (p *Portfolio) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// GetOfficeID returns user id
func (p Portfolio) GetOfficeID() string {
	return p.OfficeID.Hex()
}

// SetOfficeID set user id
func (p *Portfolio) SetOfficeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}
