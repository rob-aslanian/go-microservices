package account

import (
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
)

// UserLocation ...
type UserLocation struct {
	location.Location `bson:",inline"`
	Permission        Permission `bson:"permission,omitempty"`
}
