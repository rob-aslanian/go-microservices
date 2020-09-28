package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/rentalRPC"
)

// AddRealEstate ...
func (s Server) AddRealEstate(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	id, err := s.house.AddRealEstate(ctx,
		data.GetCompanyID(),
		realEstateRPCToStruct(data),
	)

	if err != nil {
		return &rentalRPC.ID{}, err
	}

	return &rentalRPC.ID{
		ID: id,
	}, nil
}

// GetRealEstates ...
func (s Server) GetRealEstates(ctx context.Context, data *rentalRPC.GetRealEstateRequest) (*rentalRPC.Estates, error) {
	res, err := s.house.GetRealEstates(ctx,
		rentalDealTypeRPCToDealType(data.GetDealType()),
		data.GetFirst(),
		data.GetAfter(),
	)

	if err != nil {
		return nil, err
	}

	return &rentalRPC.Estates{
		Amount:  res.Amount,
		Estates: rentalEstatesToRPC(res),
	}, nil
}

// AddHouseRentalAppartament ...
func (s Server) AddHouseRentalAppartament(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	return nil, nil
}

// AddRealEstateBuildings ...
func (s Server) AddRealEstateBuildings(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	return nil, nil
}

// AddRealEstateCommercial ...
func (s Server) AddRealEstateCommercial(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	return nil, nil
}

// AddRealEstateGarage ...
func (s Server) AddRealEstateGarage(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	return nil, nil
}

// AddRealEstateHotelRooms ...
func (s Server) AddRealEstateHotelRooms(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	return nil, nil
}

// AddRealEstateLand ...
func (s Server) AddRealEstateLand(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	return nil, nil
}

// AddRealEstateOffice ...
func (s Server) AddRealEstateOffice(ctx context.Context, data *rentalRPC.AddRentalRequest) (*rentalRPC.ID, error) {
	return nil, nil
}
