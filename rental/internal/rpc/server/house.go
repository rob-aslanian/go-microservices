package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/rental"
)

// House define functions inside house
type House interface {
	AddHouseRentalAppartament(ctx context.Context, companyID string, data rental.Appartament) (string, error)
	AddRealEstateBuildings(ctx context.Context, companyID string, data rental.Buildings) (string, error)
	AddRealEstateCommercial(ctx context.Context, companyID string, data rental.Commercial) (string, error)
	AddRealEstateGarage(ctx context.Context, companyID string, data rental.Garage) (string, error)
	AddRealEstateHotelRooms(ctx context.Context, companyID string, data rental.HotelRooms) (string, error)
	AddRealEstateLand(ctx context.Context, companyID string, data rental.Land) (string, error)
	AddRealEstateOffice(ctx context.Context, companyID string, data rental.Office) (string, error)
	AddRealEstate(ctx context.Context, companyID string, data interface{}) (string, error)
	GetRealEstates(ctx context.Context, dealType rental.DealType, first uint32, after string) (rental.GetRental, error)
}
