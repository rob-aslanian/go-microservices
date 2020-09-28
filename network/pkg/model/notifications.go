package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Type represents what kind of notification is
type Type string

const (
	// TypeNewFollow ...
	TypeNewFollow Type = "new_follow"
	// TypeNewConnection ...
	TypeNewConnection Type = "new_connection"
	// TypeApproveConnectionRequest ...
	TypeApproveConnectionRequest Type = "approved_connection"
	// TypeNewRecommendationRequest ...
	TypeNewRecommendationRequest Type = "recommendation_request"
	// TypeNewRecommendation ...
	TypeNewRecommendation Type = "new_recommendation"
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

// NewFollow represent notification for a new follow
type NewFollow struct {
	Notification `json:",inline"`

	UserSenderID string `json:"user_sender_id"` // id of user who send follow request
}

// NewConnectionRequest represent notification when someone send connection request
type NewConnectionRequest struct {
	Notification `json:",inline"`

	UserSenderID string `json:"user_sender_id"` // id of user who send connection request
	FriendshipID string `json:"friendship_id"`
}

// NewApproveConnectionRequest ...
type NewApproveConnectionRequest struct {
	Notification `json:",inline"`

	UserSenderID string `json:"user_sender_id"` // id of user who approved connection request
}

// NewRecommendationRequest ...
type NewRecommendationRequest struct {
	Notification `json:",inline"`

	UserSenderID     string `json:"user_sender_id"` // id of user who sent recommendation
	Text             string `json:"text"`
	RecommendationID string `json:"recommendation_id"`
}

// NewRecommendation ...
type NewRecommendation struct {
	Notification `json:",inline"`

	UserSenderID string `json:"user_sender_id"` // id of user who sent recommendation
	Text         string `json:"text"`
}
