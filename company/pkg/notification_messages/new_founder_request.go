package notmes

// NewFounderRequest ...
type NewFounderRequest struct {
	Notification `json:",inline"`

	RequestID string `json:"request_id"`
	Founder   string `json:"founder"`
	CompanyID string `json:"company_id"`
}
