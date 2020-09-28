package group

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Member ...
type Member struct {
	ID        primitive.ObjectID `bson:"id"`
	CreatedAt time.Time          `bson:"created_at"`
	IsAdmin   *bool              `bson:"is_admin,omitempty"`
}

// GetID returns id of ad
func (m Member) GetID() string {
	return m.ID.Hex()
}

// SetID ...
func (m *Member) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	m.ID = objID
	return nil
}

// InvitedMember ...
type InvitedMember struct {
	Member    `bson:",inline"`
	InvitedBy primitive.ObjectID `bson:"invited_by"`
}

// GetInvitedByID ...
func (m InvitedMember) GetInvitedByID() string {
	return m.InvitedBy.Hex()
}

// SetInvitedByID ...
func (m *InvitedMember) SetInvitedByID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	m.InvitedBy = objID
	return nil
}
