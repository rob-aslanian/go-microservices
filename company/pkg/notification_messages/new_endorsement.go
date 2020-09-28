package notmes

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Type represents what kind of notification is
type Type string

const (
	// TypeNewReview ...
	TypeNewReview Type = "new_review"
	// TypeNewFounderRequest ...
	TypeNewFounderRequest Type = "new_founder_request"
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

// NewCompanyReview ...
type NewCompanyReview struct {
	Notification `json:",inline"`

	ReviewerUserID string `json:"reviewer_user_id"`
	ReviewID       string `json:"rewiew_id"`
}
