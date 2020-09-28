package houses

import (
	"context"

	notmes "gitlab.lan/Rightnao-site/microservices/rental/internal/notification_messages"
	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/rental"

	companyadmin "gitlab.lan/Rightnao-site/microservices/rental/internal/company-admin"
)

// AuthRPC represents auth service
type AuthRPC interface {
	GetUserID(ctx context.Context) (string, error)
}

// NetworkRPC represents network service
type NetworkRPC interface {
	GetAdminLevel(ctx context.Context, companyID string) (companyadmin.AdminLevel, error)
}

// InfoRPC represents info service
type InfoRPC interface {
	GetUserCountry(ctx context.Context) (string, error)
	GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error)
}

// UserRPC represents a Info gRPC client
type UserRPC interface {
	CheckPassword(ctx context.Context, password string) error
}

// ChatRPC ...
type ChatRPC interface {
	IsLive(ctx context.Context, id string) (bool, error)
}

// HousesRepository is ll about offices. creating it, changing it, posting services from it etc
type HousesRepository interface {
	AddHouseRentalAppartament(ctx context.Context, data rental.Appartament) error
	AddRealEstateBuildings(ctx context.Context, data rental.Buildings) error
	AddRealEstateCommercial(ctx context.Context, data rental.Commercial) error
	AddRealEstateGarage(ctx context.Context, data rental.Garage) error
	AddRealEstateHotelRooms(ctx context.Context, data rental.HotelRooms) error
	AddRealEstateLand(ctx context.Context, data rental.Land) error
	AddRealEstateOffice(ctx context.Context, data rental.Office) error
	AddRealEstateStorageRooms(ctx context.Context, data rental.StorageRooms) error
	AddRealEstateRuralFarm(ctx context.Context, data rental.RuralFarm) error
	AddRealEstateRenovation(ctx context.Context, data rental.Renovation) error
	AddRealEstateMaterials(ctx context.Context, data rental.Materials) error
	AddRealEstateMove(ctx context.Context, data rental.Move) error
	GetRealEstates(ctx context.Context, dealType rental.DealType, first int, after int) (rental.GetRental, error)
}

// CacheRepository contains functions which have to be in cache repository
type CacheRepository interface {
	CreateTemporaryCodeForEmailActivation(ctx context.Context, companyID string, email string) (string, error)
	CheckTemporaryCodeForEmailActivation(ctx context.Context, companyID string, code string) (bool, string, error)
	Remove(ctx context.Context, key string) error
}

// MQ ...
type MQ interface {
	OrderService(targetID string, not *notmes.NewOrder) error
	SendProposalForServiceRequest(targetID string, not *notmes.NewProposal) error
}
