package notmes

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Type represents what kind of notification is
type Type string

const (
	// TypeNewInvitation ...
	TypeNewInvitation Type = "new_job_invitation"
	// TypeNewApplicant ...
	TypeNewApplicant Type = "new_job_applicant"
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

// NewInvitation represent notification for job invitation
type NewInvitation struct {
	Notification `json:",inline"`

	CompanyID string `json:"company_id"`
	JobID     string `json:"job_id"`
}

// NewJobApplicant ...
type NewJobApplicant struct {
	Notification `json:",inline"`

	CandidateID string `json:"candidate_id"`
	JobID       string `json:"job_id"`
	CompanyID   string `json:"company_id"`
}
