package notmes

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Type represents what kind of notification is
type Type string

const (
	// TypeNewEndorsement ...
	TypeNewEndorsement Type = "new_endorsement"
)

// Notification ...
type Notification struct {
	ID         string    `json:"id"`
	Type       Type      `json:"type"`
	ReceiverID string    `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
	// IsSeen     bool      `json:"seen"`
}

// GenerateID ...
func (n *Notification) GenerateID() {
	n.ID = bson.NewObjectId().Hex()
}

// NewEndorsement represent notification for verifying skill
type NewEndorsement struct {
	Notification `json:",inline"`

	UserSenderID  string `json:"user_sender_id"` // id of user who send connection request
	EndorsementID string `json:"endorsement"`
	SkillID       string `json:"skill_id"`
}
