package notmes

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Type represents what kind of notification is
type Type string

const (
	// TypeNewOrder ...
	TypeNewOrder Type = "new_order"
	// TypeNewProposal ...
	TypeNewProposal Type = "new_proposal"
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

// NewOrder represent notification for job new order
type NewOrder struct {
	Notification `json:",inline"`

	CompanyID string `json:"company_id"`
	UserSenderID  string `json:"user_sender_id"` 
	ServiceID string `json:"service_id"`
}

// NewProposal ...
type NewProposal struct {
	Notification `json:",inline"`

	CompanyID string `json:"company_id"`
	UserSenderID  string `json:"user_sender_id"` 
	RequestID string `json:"request_id"`
}
