package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/rentalRPC"
)

// AddRealEstate ...
func (*Resolver) AddRealEstate(ctx context.Context, in AddRealEstateRequest) (*SuccessResolver, error) {
	id, err := rental.AddRealEstate(ctx, realEstateToRPC(in.Input))

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// GetRealEstates ...
func (*Resolver) GetRealEstates(ctx context.Context, in GetRealEstatesRequest) (*RealEstatesResolver, error) {
	res, err := rental.GetRealEstates(ctx, &rentalRPC.GetRealEstateRequest{
		DealType: rentalDealTypeToRPC(in.Deal_type),
		First:    Nullint32ToUint32(in.Pagination.First),
		After:    NullToString(in.Pagination.After),
	})

	if err != nil {
		return nil, err
	}

	return &RealEstatesResolver{
		R: &RealEstates{
			Amount:  res.GetAmount(),
			Estates: realEstatesRPCToEstates(res.GetEstates()),
		},
	}, nil
}

// AddRealEstateAppartamentOrHouse ...
func (*Resolver) AddRealEstateAppartamentOrHouse(ctx context.Context, in AddRealEstateAppartamentOrHouseRequest) (*SuccessResolver, error) {

	id, err := rental.AddHouseRentalAppartament(ctx, &rentalRPC.AddRentalRequest{
		CompanyID: NullToString(in.Input.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// AddRealEstateBuildings ...
func (*Resolver) AddRealEstateBuildings(ctx context.Context, in AddRealEstateBuildingsRequest) (*SuccessResolver, error) {
	id, err := rental.AddRealEstateBuildings(ctx, &rentalRPC.AddRentalRequest{
		CompanyID: NullToString(in.Input.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// AddRealEstateCommercial ...
func (*Resolver) AddRealEstateCommercial(ctx context.Context, in AddRealEstateCommercialRequest) (*SuccessResolver, error) {
	id, err := rental.AddRealEstateCommercial(ctx, &rentalRPC.AddRentalRequest{
		CompanyID: NullToString(in.Input.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// AddRealEstateGarage ...
func (*Resolver) AddRealEstateGarage(ctx context.Context, in AddRealEstateGarageRequest) (*SuccessResolver, error) {
	id, err := rental.AddRealEstateGarage(ctx, &rentalRPC.AddRentalRequest{
		CompanyID: NullToString(in.Input.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// AddRealEstateHotelRooms ...
func (*Resolver) AddRealEstateHotelRooms(ctx context.Context, in AddRealEstateHotelRoomsRequest) (*SuccessResolver, error) {
	id, err := rental.AddRealEstateHotelRooms(ctx, &rentalRPC.AddRentalRequest{
		CompanyID: NullToString(in.Input.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// AddRealEstateLand ...
func (*Resolver) AddRealEstateLand(ctx context.Context, in AddRealEstateLandRequest) (*SuccessResolver, error) {
	id, err := rental.AddRealEstateLand(ctx, &rentalRPC.AddRentalRequest{
		CompanyID: NullToString(in.Input.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// AddRealEstateOffice ...
func (*Resolver) AddRealEstateOffice(ctx context.Context, in AddRealEstateOfficeRequest) (*SuccessResolver, error) {
	id, err := rental.AddRealEstateOffice(ctx, &rentalRPC.AddRentalRequest{
		CompanyID: NullToString(in.Input.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      id.GetID(),
		},
	}, nil
}

// ID ...
func (r RealEstateInterfaceResolver) ID() graphql.ID {
	return (*r.r).ID()
}

// UserID ...
func (r RealEstateInterfaceResolver) UserID() graphql.ID {
	return (*r.r).User_id()
}

// CompanyID ...
func (r RealEstateInterfaceResolver) CompanyID() graphql.ID {
	return (*r.r).Company_id()
}

// RentalInfo ...
func (r RealEstateInterfaceResolver) RentalInfo() RentalInfoResolver {
	return (*r.r).Rental_info()
}

// RentalDetails ...
func (r RealEstateInterfaceResolver) RentalDetails() []RentalDetailResolver {
	return (*r.r).Rental_details()
}

// Files ...
func (r RealEstateInterfaceResolver) Files() []FileResolver {
	return (*r.r).Files()
}

// HasRepossessed ...
func (r RealEstateInterfaceResolver) HasRepossessed() bool {
	return (*r.r).Has_repossessed()
}

// IsUrgent ...
func (r RealEstateInterfaceResolver) IsUrgent() bool {
	return (*r.r).Is_urgent()
}

// Alerts ...
func (r RealEstateInterfaceResolver) Alerts() int32 {
	return (*r.r).Alerts()
}

// Offers ...
func (r RealEstateInterfaceResolver) Offers() int32 {
	return (*r.r).Offers()
}

// Likes ...
func (r RealEstateInterfaceResolver) Likes() int32 {
	return (*r.r).Likes()
}

// MetrictType ...
func (r RealEstateInterfaceResolver) MetrictType() string {
	return (*r.r).Metrict_type()
}

// AvailabilityFrom ...
func (r RealEstateInterfaceResolver) AvailabilityFrom() string {
	return (*r.r).Availability_from()
}

// AvailabilityTo ...
func (r RealEstateInterfaceResolver) AvailabilityTo() string {
	return (*r.r).Availability_to()
}

// ToAppartment ...
func (r RealEstateInterfaceResolver) ToAppartment() (*AppartmentResolver, bool) {
	res, ok := (*r.r).(AppartmentResolver)
	return &res, ok
}

// ToStorageRoom ...
func (r RealEstateInterfaceResolver) ToStorageRoom() (*StorageRoomResolver, bool) {
	res, ok := (*r.r).(StorageRoomResolver)
	return &res, ok
}

// ToOffice ...
func (r RealEstateInterfaceResolver) ToOffice() (*OfficeResolver, bool) {
	res, ok := (*r.r).(OfficeResolver)
	return &res, ok
}

// ToCommericalProperties ...
func (r RealEstateInterfaceResolver) ToCommericalProperties() (*CommericalPropertiesResolver, bool) {
	res, ok := (*r.r).(CommericalPropertiesResolver)
	return &res, ok
}
