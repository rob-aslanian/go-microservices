package houses

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/rental"

	companyadmin "gitlab.lan/Rightnao-site/microservices/rental/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/location"
)

// AddHouseRentalAppartament ...
func (h HouseRental) AddHouseRentalAppartament(ctx context.Context, companyID string, data rental.Appartament) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddHouseRentalAppartament")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddHouseRentalAppartament(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateMove ...
func (h HouseRental) AddRealEstateMove(ctx context.Context, companyID string, data rental.Move) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateMove")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateMove(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateRenovation ...
func (h HouseRental) AddRealEstateRenovation(ctx context.Context, companyID string, data rental.Renovation) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateRenovation")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateRenovation(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateMaterials ...
func (h HouseRental) AddRealEstateMaterials(ctx context.Context, companyID string, data rental.Materials) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateMaterials")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateMaterials(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateRuralFarm ...
func (h HouseRental) AddRealEstateRuralFarm(ctx context.Context, companyID string, data rental.RuralFarm) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateRuralFarm")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateRuralFarm(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateStorageRooms ...
func (h HouseRental) AddRealEstateStorageRooms(ctx context.Context, companyID string, data rental.StorageRooms) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateStorageRooms")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateStorageRooms(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateBuildings ...
func (h HouseRental) AddRealEstateBuildings(ctx context.Context, companyID string, data rental.Buildings) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateBuildings")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateBuildings(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateCommercial ...
func (h HouseRental) AddRealEstateCommercial(ctx context.Context, companyID string, data rental.Commercial) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateCommercial")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateCommercial(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateGarage ...
func (h HouseRental) AddRealEstateGarage(ctx context.Context, companyID string, data rental.Garage) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateGarage")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateGarage(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateHotelRooms ...
func (h HouseRental) AddRealEstateHotelRooms(ctx context.Context, companyID string, data rental.HotelRooms) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateHotelRooms")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateHotelRooms(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateLand ...
func (h HouseRental) AddRealEstateLand(ctx context.Context, companyID string, data rental.Land) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateLand")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateLand(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstateOffice ...
func (h HouseRental) AddRealEstateOffice(ctx context.Context, companyID string, data rental.Office) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstateOffice")
	defer span.Finish()

	// retrieve id of user
	profileID, err := h.authRPC.GetUserID(ctx)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := h.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
		profileID = companyID
		data.RentalInfo.IsCompany = true
	}

	rentalID, err := h.setRentalInfo(&data.RentalInfo, profileID)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}

	err = h.repository.Houses.AddRealEstateOffice(ctx, data)
	if err != nil {
		h.tracer.LogError(span, err)
		return "", err
	}
	return rentalID, nil
}

// AddRealEstate ...
func (h HouseRental) AddRealEstate(ctx context.Context, companyID string, data interface{}) (string, error) {
	span := h.tracer.MakeSpan(ctx, "AddRealEstate")
	defer span.Finish()

	switch reflect.TypeOf(data).String() {
	case "rental.Appartament":
		return h.AddHouseRentalAppartament(ctx, companyID, data.(rental.Appartament))
	case "rental.Garage":
		return h.AddRealEstateGarage(ctx, companyID, data.(rental.Garage))
	case "rental.StorageRooms":
		return h.AddRealEstateStorageRooms(ctx, companyID, data.(rental.StorageRooms))
	case "rental.Office":
		return h.AddRealEstateOffice(ctx, companyID, data.(rental.Office))
	case "rental.Commercial":
		return h.AddRealEstateCommercial(ctx, companyID, data.(rental.Commercial))
	case "rental.Buildings":
		return h.AddRealEstateBuildings(ctx, companyID, data.(rental.Buildings))
	case "rental.Land":
		return h.AddRealEstateLand(ctx, companyID, data.(rental.Land))
	case "rental.Renovation":
		return h.AddRealEstateRenovation(ctx, companyID, data.(rental.Renovation))
	case "rental.Materials":
		return h.AddRealEstateMaterials(ctx, companyID, data.(rental.Materials))
	case "rental.Move":
		return h.AddRealEstateMove(ctx, companyID, data.(rental.Move))
	}

	return "", nil

}

// GetRealEstates ...
func (h HouseRental) GetRealEstates(ctx context.Context, dealType rental.DealType, first uint32, after string) (r rental.GetRental, e error) {
	span := h.tracer.MakeSpan(ctx, "GetRealEstates")
	defer span.Finish()

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return
	}
	if afterNumber < 0 {
		return
	}

	res, err := h.repository.Houses.GetRealEstates(ctx, dealType, int(first), int(afterNumber))
	if err != nil {
		h.tracer.LogError(span, err)
		return r, err
	}
	return res, nil
}

func (h HouseRental) setRentalInfo(info *rental.CommonRental, profileID string) (string, error) {

	if info == nil {
		return "", errors.New("Info can`t be nil")
	}

	id := info.GenerateID()
	info.SetOwnerID(profileID)
	info.CreatedAt = time.Now()
	info.PostStatus = "active"

	cityID, _ := strconv.Atoi(info.Location.City.ID)
	lang := "en"
	/// Location
	cityName, subdivision, countryID, err := h.infoRPC.GetCityInformationByID(context.TODO(), int32(cityID), &lang)

	if err != nil {
		return "", err
	}

	info.Location = location.Location{
		City: &location.City{
			Name:        cityName,
			Subdivision: subdivision,
		},
		Country: &location.Country{
			ID: countryID,
		},
	}

	return id, nil
}

// return false if level doesn't much
func (h HouseRental) checkAdminLevel(ctx context.Context, companyID string, requiredLevels ...companyadmin.AdminLevel) bool {
	span := h.tracer.MakeSpan(ctx, "checkAdminLevel")
	defer span.Finish()

	actualLevel, err := h.networkRPC.GetAdminLevel(ctx, companyID)
	if err != nil {
		h.tracer.LogError(span, err)
		log.Println("Error: checkAdminLevel:", err)
		return false
	}

	for _, lvl := range requiredLevels {
		if lvl == actualLevel {
			return true
		}
	}

	return false
}
