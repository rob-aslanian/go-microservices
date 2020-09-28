package shop

import (
	"errors"
	"strings"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Shop ...
type Shop struct {
	ID           primitive.ObjectID  `bson:"_id"`
	UserID       *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID    *primitive.ObjectID `bson:"company_id,omitempty"`
	Logo         *string             `bson:"logo,omitempty"`
	Cover        *string             `bson:"cover,omitempty"`
	Title        string              `bson:"title"`
	Category     Category            `bson:"category"`
	SellerType   string              `bson:"seller_type"`
	ProductsType []string            `bson:"products_type"`
	Description  string              `bson:"description"`
	Showcase     *string             `bson:"showcase,omitempty"`
}

// GetID returns id of ad
func (g Shop) GetID() string {
	return g.ID.Hex()
}

// SetID ...
func (g *Shop) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.ID = objID
	return nil
}

// GenerateID creates new id
func (g *Shop) GenerateID() string {
	g.ID = primitive.NewObjectID()
	return g.ID.Hex()
}

// GetUserID ...
func (g Shop) GetUserID() string {
	if g.UserID == nil {
		return ""
	}

	return g.UserID.Hex()
}

// SetUserID ...
func (g *Shop) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.UserID = &objID
	return nil
}

// GetCompanyID ...
func (g Shop) GetCompanyID() string {
	if g.CompanyID == nil {
		return ""
	}

	return g.CompanyID.Hex()
}

// SetCompanyID ...
func (g *Shop) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.CompanyID = &objID
	return nil
}

// ValidateTitle checks if not ampty and no more then 60 characters
func (g Shop) ValidateTitle() error {
	err := ValidateName(g.Title)
	if err != nil {
		return err
	}
	return nil
}

// ValidateName ...
func ValidateName(name string) error {
	if len(name) == 0 {
		return errors.New("empty_name")
	}
	if utf8.RuneCountInString(name) > 60 {
		return errors.New("empty_name")
	}
	return nil
}

// Trim trims title, description.
func (g *Shop) Trim() {
	strings.TrimSpace(g.Title)
	strings.TrimSpace(g.Description)
}
