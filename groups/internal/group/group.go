package group

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"gitlab.lan/Rightnao-site/microservices/groups/internal/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group ...
type Group struct {
	ID                  primitive.ObjectID `bson:"_id"`
	OwnerID             primitive.ObjectID `bson:"owner_id"`
	Name                string             `bson:"name"`
	Type                string             `bson:"type"`
	Privacy             PrivacyType        `bson:"privacy_type"`
	Members             []Member           `bson:"members"`
	Tagline             string             `bson:"tagline"`
	Description         string             `bson:"description"`
	Rules               string             `bson:"rules"`
	Location            *location.Location `bson:"location"`
	Cover               string             `bson:"cover"`
	CoverOriginal       string             `bson:"cover_original"`
	URL                 string             `bson:"url"`
	PostApprovalByAdmin bool               `bson:"post_approval_by_admin"`
	// AdminActivity []Log
	// MemberRequests []Member
	// PendingPosts []primitive.ObjectID
	// ReportedContent
	// Admins []primitive.ObjectID `bson:"admins"`
	// Blocked []primitive.ObjectID
	// Invited []primitive.ObjectID
	CreatedAt       time.Time `bson:"created_at"`
	AmountOfMembers uint32    `bson:"amount_of_members"`
}

// GetID returns id of ad
func (g Group) GetID() string {
	return g.ID.Hex()
}

// SetID ...
func (g *Group) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.ID = objID
	return nil
}

// GenerateID creates new id
func (g *Group) GenerateID() string {
	g.ID = primitive.NewObjectID()
	return g.ID.Hex()
}

// GetOwnerID ...
func (g Group) GetOwnerID() string {
	return g.OwnerID.Hex()
}

// SetOwnerID ...
func (g *Group) SetOwnerID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.OwnerID = objID
	return nil
}

// ValidateName checks if not ampty and no more then 60 characters
func (g Group) ValidateName() error {
	err := ValidateName(g.Name)
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

// Trim trims name, tagline, description, rules.
func (g *Group) Trim() {
	strings.TrimSpace(g.Name)
	strings.TrimSpace(g.Tagline)
	strings.TrimSpace(g.Description)
	strings.TrimSpace(g.Rules)
}

// IsAdmin ...
func (g Group) IsAdmin(id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	for i := range g.Members {
		if g.Members[i].ID == objID &&
			g.Members[i].IsAdmin != nil &&
			*g.Members[i].IsAdmin == true {
			return true, nil
		}
	}
	return false, nil
}

// // AddAdmin ...
// func (g *Group) AddAdmin(id string) error {
// 	if g.Admins == nil {
// 		g.Admins = make([]primitive.ObjectID, 0, 1)
// 	}
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}
//
// 	g.Admins = append(g.Admins, objID)
// 	return nil
// }

// IsMember ...
func (g Group) IsMember(id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	for i := range g.Members {
		if g.Members[i].ID == objID {
			return true, nil
		}
	}

	return false, nil
}

// AddMember ...
func (g *Group) AddMember(id string, isAdmin bool) error {
	if g.Members == nil {
		g.Members = make([]Member, 0, 1)
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	m := Member{
		ID:        objID,
		CreatedAt: time.Now(),
	}

	if isAdmin {
		m.IsAdmin = &isAdmin
	}

	g.Members = append(g.Members, m)
	return nil
}
