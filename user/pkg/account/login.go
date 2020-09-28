package account

import (
	"gitlab.lan/Rightnao-site/microservices/user/pkg/status"
)

// LoginResponse ... 
type LoginResponse struct {
	ID           string
	Status       status.UserStatus
	URL          string
	Token        string
	FirstName    string
	LastName     string
	Is2FAEnabled bool
	TwoFASecret  []byte
	Avatar       string
	Password     string
	Gender       string
	Email 		 string
}
