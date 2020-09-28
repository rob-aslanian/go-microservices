package serverRPC

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/rentalRPC"
	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/location"
	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/rental"
)

func realEstateRPCToStruct(data *rentalRPC.AddRentalRequest) interface{} {
	if data == nil {
		return nil
	}

	propertyType := data.GetRental().GetPropertyType()
	dealType := data.GetRental().GetDealType()

	if dealType == rentalRPC.DealTypeEnum_DealType_Renovation ||
		dealType == rentalRPC.DealTypeEnum_DealType_Materials ||
		dealType == rentalRPC.DealTypeEnum_DealType_Move {

		switch dealType {
		case rentalRPC.DealTypeEnum_DealType_Renovation:
			return rentalRenovationRPCToStruct(data)
		case rentalRPC.DealTypeEnum_DealType_Materials:
			return rentalMaterialsRPCToStruct(data)
		case rentalRPC.DealTypeEnum_DealType_Move:
			return rentalMoveRPCToStruct(data)
		}

	} else {
		switch propertyType {
		case rentalRPC.PropertyTypeEnum_PropertyType_Appartments:
			return rentalAppartamentRPCToAppartament(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_NewHomes:
			return rentalAppartamentRPCToAppartament(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_Homes:
			return rentalAppartamentRPCToAppartament(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_Houses:
			return rentalAppartamentRPCToAppartament(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_SummerCottage:
			return rentalAppartamentRPCToAppartament(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_Garages:
			return rentalGarageRPCToStruct(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_StorageRooms:
			return rentalStorageRoomsRPCToStruct(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_Offices:
			return rentalOfficeRPCToStruct(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_CommercialProperties:
			return rentalCommercialRPCToStruct(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_Buildings:
			return rentalBuildingRPCToStruct(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_Land:
			return rentalLandRPCToStruct(data)
		case rentalRPC.PropertyTypeEnum_PropertyType_RuralFarm:
			return rentalRuralFarmRPCToStruct(data)
		}

	}

	return nil
}

/* TODO
   - Add Land type
   - Add Building type
   - Add Summer cotage type
   - Add Rural farm type
   - Add house type ( is`s like appartments )
   - Add homes type ( is`s like appartments )
   - Add new home type ( is`s like appartments )
   - Add Garage type
*/

func rentalEstatesToRPC(data rental.GetRental) []*rentalRPC.Estate {

	estates := make([]*rentalRPC.Estate, 0, data.Amount)

	if data.Amount > 0 {

		// Appartaments
		if len(data.Appartaments) > 0 {
			for _, a := range data.Appartaments {
				estate := rentalRPC.Estate{}
				estate.Estates = &rentalRPC.Estate_Appartaments{
					Appartaments: rentalAppartanentsToRPC(a),
				}
				estates = append(estates, &estate)
			}

		}
		// Storage Rooms
		if len(data.StorageRooms) > 0 {
			for _, sr := range data.StorageRooms {
				estate := rentalRPC.Estate{}
				estate.Estates = &rentalRPC.Estate_StorageRoom{
					StorageRoom: rentalStorageRoomToRPC(sr),
				}
				estates = append(estates, &estate)
			}

		}
		// Office
		if len(data.Offices) > 0 {
			for _, of := range data.Offices {
				estate := rentalRPC.Estate{}
				estate.Estates = &rentalRPC.Estate_Office{
					Office: rentalOfficeToRPC(of),
				}
				estates = append(estates, &estate)
			}

		}
		// CommercialProperties
		if len(data.CommercialProperties) > 0 {
			for _, cp := range data.CommercialProperties {
				estate := rentalRPC.Estate{}
				estate.Estates = &rentalRPC.Estate_CommercialAndRuralFarm{
					CommercialAndRuralFarm: rentalCommercialToRPC(cp),
				}
				estates = append(estates, &estate)
			}

		}

	}

	return estates
}

func rentalAppartanentsToRPC(data rental.Appartament) *rentalRPC.Appartaments {
	return &rentalRPC.Appartaments{
		HasRepossesed:   data.RentalInfo.HasRepossesed,
		AvailibatiFrom:  data.AvailibatiFrom,
		AvailibatiTo:    data.AvailibatiTo,
		BadRooms:        data.BadRooms,
		BathRooms:       data.BathRooms,
		CarSpaces:       data.CarSpaces,
		Floor:           data.Floor,
		Floors:          data.Floors,
		IsAgent:         data.IsAgent,
		TotalArea:       data.TotalArea,
		Rental:          rentalInfoToRPC(data.RentalInfo),
		Details:         rentalDetailsToRPC(data.Details),
		ClimatControl:   rentalClimatControlsToRPC(data.ClimatControl),
		IndoorFeatures:  rentalIndoorFeaturesToRPC(data.IndoorFeatures),
		OutdoorFeatures: rentalOutdoorFeaturesToRPC(data.OutdoorFeatures),
		MetricType:      rentalPriceTypeToRPC(data.Metric),
		Phones:          rentalPhonesToRPC(data.Phones),
		Price:           rentalPriceToRPC(data.Price),
		TypeOfProperty:  rentalOfPropertiesToRPC(data.TypeOfProperty),
		Status:          rentalStatusToRPC(data.Status),
	}
}

func rentalCommercialToRPC(data rental.Commercial) *rentalRPC.CommercialAndRuralFarm {
	return &rentalRPC.CommercialAndRuralFarm{
		HasRepossesed:              data.RentalInfo.HasRepossesed,
		AvailibatiFrom:             data.AvailibatiFrom,
		AvailibatiTo:               data.AvailibatiTo,
		IsAgent:                    data.IsAgent,
		TotalArea:                  data.TotalArea,
		CommercialProperties:       rentalCommercilaPropertiesToRPC(data.CommercialProperties),
		CommericalPropertyLocation: rentalCommercialLocationsToRPC(data.CommericalLocation),
		AdditionalFilters:          rentalAddtionalFiltersToRPC(data.AdditionalFilters),
		Rental:                     rentalInfoToRPC(data.RentalInfo),
		Details:                    rentalDetailsToRPC(data.Details),
		MetricType:                 rentalPriceTypeToRPC(data.Metric),
		Phones:                     rentalPhonesToRPC(data.Phones),
		Price:                      rentalPriceToRPC(data.Price),
		Status:                     rentalStatusToRPC(data.Status),
	}
}

func rentalOfficeToRPC(data rental.Office) *rentalRPC.Office {
	return &rentalRPC.Office{
		HasRepossesed:  data.RentalInfo.HasRepossesed,
		AvailibatiFrom: data.AvailibatiFrom,
		AvailibatiTo:   data.AvailibatiTo,
		IsAgent:        data.IsAgent,
		TotalArea:      data.TotalArea,
		Layout:         rentalLayoutToRPC(data.Layout),
		BuildingUse:    rentalBuildinUseToRPC(data.BuildingUse),
		Rental:         rentalInfoToRPC(data.RentalInfo),
		Details:        rentalDetailsToRPC(data.Details),
		MetricType:     rentalPriceTypeToRPC(data.Metric),
		Phones:         rentalPhonesToRPC(data.Phones),
		Price:          rentalPriceToRPC(data.Price),
		Status:         rentalStatusToRPC(data.Status),
	}
}

func rentalStorageRoomToRPC(data rental.StorageRooms) *rentalRPC.StorageRoom {
	return &rentalRPC.StorageRoom{
		HasRepossesed:  data.RentalInfo.HasRepossesed,
		AvailibatiFrom: data.AvailibatiFrom,
		AvailibatiTo:   data.AvailibatiTo,
		IsAgent:        data.IsAgent,
		TotalArea:      data.TotalArea,
		Rental:         rentalInfoToRPC(data.RentalInfo),
		Details:        rentalDetailsToRPC(data.Details),
		MetricType:     rentalPriceTypeToRPC(data.Metric),
		Phones:         rentalPhonesToRPC(data.Phones),
		Price:          rentalPriceToRPC(data.Price),
		Status:         rentalStatusToRPC(data.Status),
	}
}

func rentalAddtionalFiltersToRPC(data []rental.AdditionalFilters) []rentalRPC.AdditionalFiltersEnum {
	if len(data) <= 0 {
		return nil
	}

	filters := make([]rentalRPC.AdditionalFiltersEnum, 0, len(data))

	for _, f := range data {
		filters = append(filters, rentalAddtionaFilterToRPC(f))
	}

	return filters
}

func rentalCommercialLocationsToRPC(data []rental.CommericalPropertyLocation) []rentalRPC.CommericalPropertyLocationEnum {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]rentalRPC.CommericalPropertyLocationEnum, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalCommercialLocationToRPC(p))
	}

	return properties
}

func rentalCommercilaPropertiesToRPC(data []rental.CommericalProperty) []rentalRPC.CommercialPropertyEnum {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]rentalRPC.CommercialPropertyEnum, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalCommercialPropertyToRPC(p))
	}

	return properties
}

func rentalOfPropertiesToRPC(data []rental.TypeOfProperty) []rentalRPC.TypeOfPropertyEnum {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]rentalRPC.TypeOfPropertyEnum, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalOfPropertyToRPC(p))
	}

	return properties
}

func rentalPriceToRPC(data rental.Price) *rentalRPC.Price {
	return &rentalRPC.Price{
		Currency:  data.Currency,
		MaxPrice:  data.MaxPrice,
		MinPrice:  data.MinPrice,
		PriceType: rentalPriceTypeToRPC(data.PriceType),
	}
}

func rentalPhonesToRPC(data []rental.Phone) []*rentalRPC.Phone {
	if len(data) <= 0 {
		return nil
	}

	phones := make([]*rentalRPC.Phone, 0, len(data))

	for _, p := range data {
		phones = append(phones, &rentalRPC.Phone{
			Id:          p.GetID(),
			CountryCode: p.CountryCode,
			Number:      p.Number,
		})
	}

	return phones
}

func rentalOutdoorFeaturesToRPC(data []rental.OutdoorFeatures) []rentalRPC.OutdoorFeaturesEnum {
	if len(data) <= 0 {
		return nil
	}

	features := make([]rentalRPC.OutdoorFeaturesEnum, 0, len(data))

	for _, f := range data {
		features = append(features, rentalOutdoorFeatureToRPC(f))
	}

	return features
}

func rentalIndoorFeaturesToRPC(data []rental.IndoorFeatures) []rentalRPC.IndoorFeaturesEnum {
	if len(data) <= 0 {
		return nil
	}

	features := make([]rentalRPC.IndoorFeaturesEnum, 0, len(data))

	for _, f := range data {
		features = append(features, rentalIndoorFeatureToRPC(f))
	}

	return features
}

func rentalClimatControlsToRPC(data []rental.ClimatControl) []rentalRPC.ClimatControlEnum {
	if len(data) <= 0 {
		return nil
	}

	controls := make([]rentalRPC.ClimatControlEnum, 0, len(data))

	for _, c := range data {
		controls = append(controls, rentalClimatControlToRPC(c))
	}

	return controls
}

func rentalDetailsToRPC(data []rental.Detail) []*rentalRPC.Detail {
	if len(data) <= 0 {
		return nil
	}

	details := make([]*rentalRPC.Detail, 0, len(data))

	for _, d := range data {
		details = append(details, &rentalRPC.Detail{
			ID:          d.GetID(),
			Title:       d.Title,
			HouseRules:  d.HouseRules,
			Description: d.Description,
		})
	}

	return details
}

func rentalLocationToRPC(data location.Location) *rentalRPC.Location {
	return &rentalRPC.Location{
		Address: nullToString(data.Address),
		Street:  nullToString(data.Street),
		City: &rentalRPC.City{
			Id:          data.City.ID,
			City:        data.City.Name,
			Subdivision: data.City.Subdivision,
		},
		Country: &rentalRPC.Country{
			Id: data.Country.ID,
		},
	}
}

func rentalInfoToRPC(data rental.CommonRental) *rentalRPC.Rental {
	return &rentalRPC.Rental{
		ID:           data.GetID(),
		OwnerID:      data.GetOwnerID(),
		ExpiredDays:  data.ExpiredDays,
		IsCompany:    data.IsCompany,
		Offers:       data.Offers,
		IsUrgent:     data.IsUrgent,
		Shares:       data.Shares,
		Views:        data.Views,
		Alerts:       data.Alerts,
		PostCurrency: data.PostCurrency,
		CreatedAt:    data.CreatedAt.String(),
		DealType:     rentalDealTypeToRPC(data.DealType),
		Location:     rentalLocationToRPC(data.Location),
		PropertyType: rentalPropertyTypeToRPC(data.PropertyType),
	}
}

func rentalMoveRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Move {
	if data == nil {
		return rental.Move{}
	}

	return rental.Move{
		CountryIDs:   data.GetCountryIDs(),
		LocationType: rentalLocationTypesRPCToLocationTypes(data.GetLocationType()),
		Services:     rentalServicesRPCToServices(data.GetServices()),
		RentalInfo:   rentalInfoRPCToRentalInfo(data.GetRental()),
		Details:      rentalDetailsRPCToDetails(data.GetDetails()),
	}
}

func rentalMaterialsRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Materials {
	if data == nil {
		return rental.Materials{}
	}

	return rental.Materials{
		Materials:  rentalMaterialsRPCToMaterials(data.GetMaterials()),
		RentalInfo: rentalInfoRPCToRentalInfo(data.GetRental()),
		Details:    rentalDetailsRPCToDetails(data.GetDetails()),
	}
}

func rentalRenovationRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Renovation {
	if data == nil {
		return rental.Renovation{}
	}

	return rental.Renovation{
		CityIDs:             data.GetCityIDs(),
		CountryIDs:          data.GetCountryIDs(),
		Exetior:             rentalPriceRPCToPrice(data.GetExterior()),
		Interior:            rentalPriceRPCToPrice(data.GetInterior()),
		InteriorAndExterior: rentalPriceRPCToPrice(data.GetInteriorAndExterior()),
		Timing:              rentalTimingRPCToStruct(data.GetTiming()),
		RentalInfo:          rentalInfoRPCToRentalInfo(data.GetRental()),
		Details:             rentalDetailsRPCToDetails(data.GetDetails()),
	}
}

func rentalStorageRoomsRPCToStruct(data *rentalRPC.AddRentalRequest) rental.StorageRooms {
	if data == nil {
		return rental.StorageRooms{}
	}

	return rental.StorageRooms{
		IsAgent:    data.GetIsAgent(),
		RentalInfo: rentalInfoRPCToRentalInfo(data.GetRental()),
		Status:     rentalStatusRPCToStatus(data.GetStatus()),
		TotalArea:  data.GetTotalArea(),
		Metric:     rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		Phones:     rentalPhonesRPCToPhones(data.GetPhones()),
		Price:      rentalPriceRPCToPrice(data.GetPrice()),
		Details:    rentalDetailsRPCToDetails(data.GetDetails()),
	}
}

func rentalCommercialRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Commercial {
	if data == nil {
		return rental.Commercial{}
	}

	return rental.Commercial{
		AvailibatiFrom:       data.GetAvailibatiFrom(),
		AvailibatiTo:         data.GetAvailibatiTo(),
		IsAgent:              data.GetIsAgent(),
		TotalArea:            data.GetTotalArea(),
		Metric:               rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		RentalInfo:           rentalInfoRPCToRentalInfo(data.GetRental()),
		Phones:               rentalPhonesRPCToPhones(data.GetPhones()),
		Price:                rentalPriceRPCToPrice(data.GetPrice()),
		Status:               rentalStatusRPCToStatus(data.GetStatus()),
		Details:              rentalDetailsRPCToDetails(data.GetDetails()),
		AdditionalFilters:    rentalAddtionaFiltersRPCToAddiotionalFilters(data.GetAdditionalFilters()),
		CommercialProperties: rentalCommercialPropertiesRPCCommercialProperties(data.GetCommercialProperties()),
		CommericalLocation:   rentalCommercialLocationsRPCLocations(data.GetCommericalPropertyLocation()),
	}
}

func rentalGarageRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Garage {
	if data == nil {
		return rental.Garage{}
	}
	return rental.Garage{
		TotalArea:         data.GetTotalArea(),
		RentalInfo:        rentalInfoRPCToRentalInfo(data.GetRental()),
		Phones:            rentalPhonesRPCToPhones(data.GetPhones()),
		Price:             rentalPriceRPCToPrice(data.GetPrice()),
		Details:           rentalDetailsRPCToDetails(data.GetDetails()),
		Metric:            rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		AdditionalFilters: rentalAddtionaFiltersRPCToAddiotionalFilters(data.GetAdditionalFilters()),
	}
}

func rentalHotelRoomsRPCToStruct(data *rentalRPC.HotelRooms) rental.HotelRooms {
	if data == nil {
		return rental.HotelRooms{}
	}
	return rental.HotelRooms{
		TotalArea:  data.GetTotalArea(),
		Rooms:      data.GetRooms(),
		Status:     rentalStatusRPCToStatus(data.GetStatus()),
		RentalInfo: rentalInfoRPCToRentalInfo(data.GetRental()),
		Phones:     rentalPhonesRPCToPhones(data.GetPhones()),
		Price:      rentalPriceRPCToPrice(data.GetPrice()),
		Details:    rentalDetailsRPCToDetails(data.GetDetails()),
	}
}

func rentalBuildingRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Buildings {
	if data == nil {
		return rental.Buildings{}
	}

	return rental.Buildings{
		IsAgent:        data.GetIsAgent(),
		TotalArea:      data.GetTotalArea(),
		AvailibatiFrom: data.GetAvailibatiFrom(),
		AvailibatiTo:   data.GetAvailibatiTo(),
		Metric:         rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		RentalInfo:     rentalInfoRPCToRentalInfo(data.GetRental()),
		Phones:         rentalPhonesRPCToPhones(data.GetPhones()),
		Price:          rentalPriceRPCToPrice(data.GetPrice()),
		Status:         rentalStatusRPCToStatus(data.GetStatus()),
		Details:        rentalDetailsRPCToDetails(data.GetDetails()),
	}
}

func rentalAppartamentRPCToAppartament(data *rentalRPC.AddRentalRequest) rental.Appartament {
	if data == nil {
		return rental.Appartament{}
	}

	return rental.Appartament{
		AvailibatiFrom:  data.GetAvailibatiFrom(),
		AvailibatiTo:    data.GetAvailibatiTo(),
		BadRooms:        data.GetBadRooms(),
		BathRooms:       data.GetBathRooms(),
		CarSpaces:       data.GetCarSpaces(),
		Floor:           data.GetFloor(),
		Floors:          data.GetFloors(),
		IsAgent:         data.GetIsAgent(),
		TotalArea:       data.GetTotalArea(),
		Metric:          rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		RentalInfo:      rentalInfoRPCToRentalInfo(data.GetRental()),
		IndoorFeatures:  rentalIndoorFeaturesRPCtoIndoorFeatures(data.GetIndoorFeatures()),
		OutdoorFeatures: rentalOutdoorFeaturesRPCtoOutdoorFeatures(data.GetOutdoorFeatures()),
		Phones:          rentalPhonesRPCToPhones(data.GetPhones()),
		Price:           rentalPriceRPCToPrice(data.GetPrice()),
		Status:          rentalStatusRPCToStatus(data.GetStatus()),
		ClimatControl:   rentalClimatControlsRPCToClimatControls(data.GetClimatControl()),
		Details:         rentalDetailsRPCToDetails(data.GetDetails()),
		TypeOfProperty:  rentalOfPropertiesRPCToTypeOfProperties(data.GetTypeOfProperty()),
		WhoLive:         rentalWhoLivesRPCToWhoLive(data.GetWhoLive()),
	}
}

func rentalOfficeRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Office {
	if data == nil {
		return rental.Office{}
	}
	return rental.Office{
		TotalArea:      data.GetTotalArea(),
		Metric:         rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		Status:         rentalStatusRPCToStatus(data.GetStatus()),
		Layout:         rentalLayoutRPCToLayout(data.GetLayout()),
		BuildingUse:    rentalBuildinUseRPCToBuildinUse(data.GetBuildingUse()),
		AvailibatiFrom: data.GetAvailibatiFrom(),
		AvailibatiTo:   data.GetAvailibatiTo(),
		RentalInfo:     rentalInfoRPCToRentalInfo(data.GetRental()),
		Phones:         rentalPhonesRPCToPhones(data.GetPhones()),
		Price:          rentalPriceRPCToPrice(data.GetPrice()),
		Details:        rentalDetailsRPCToDetails(data.GetDetails()),
		IsAgent:        data.GetIsAgent(),
	}
}

func rentalRuralFarmRPCToStruct(data *rentalRPC.AddRentalRequest) rental.RuralFarm {
	if data == nil {
		return rental.RuralFarm{}
	}
	return rental.RuralFarm{
		AvailibatiFrom: data.GetAvailibatiFrom(),
		AvailibatiTo:   data.GetAvailibatiTo(),
		IsAgent:        data.GetIsAgent(),
		TotalArea:      data.GetTotalArea(),
		PropertyType:   rentalPropertyTypesRPCToPropertyTypes(data.GetPropertyType()),
		Additional:     rentalOutdoorFeaturesRPCtoOutdoorFeatures(data.GetOutdoorFeatures()),
		Metric:         rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		RentalInfo:     rentalInfoRPCToRentalInfo(data.GetRental()),
		Phones:         rentalPhonesRPCToPhones(data.GetPhones()),
		Price:          rentalPriceRPCToPrice(data.GetPrice()),
		Details:        rentalDetailsRPCToDetails(data.GetDetails()),
	}
}
func rentalLandRPCToStruct(data *rentalRPC.AddRentalRequest) rental.Land {
	if data == nil {
		return rental.Land{}
	}
	return rental.Land{
		AvailibatiFrom: data.GetAvailibatiFrom(),
		AvailibatiTo:   data.GetAvailibatiTo(),
		IsAgent:        data.GetIsAgent(),
		TotalArea:      data.GetTotalArea(),
		Metric:         rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		TypeOfLand:     rentalStatusesRPCToStatuses(data.GetTypeOfLand()),
		More:           rentalAddtionaFiltersRPCToAddiotionalFilters(data.GetAdditionalFilters()),
		RentalInfo:     rentalInfoRPCToRentalInfo(data.GetRental()),
		Phones:         rentalPhonesRPCToPhones(data.GetPhones()),
		Price:          rentalPriceRPCToPrice(data.GetPrice()),
		Details:        rentalDetailsRPCToDetails(data.GetDetails()),
	}
}

func rentalServicesRPCToServices(data []rentalRPC.ServiceEnum) []rental.Service {
	if len(data) <= 0 {
		return nil
	}

	services := make([]rental.Service, 0, len(data))

	for _, s := range data {
		services = append(services, rentalServiceRPCToService(s))
	}

	return services
}

func rentalLocationTypesRPCToLocationTypes(data []rentalRPC.LocationEnum) []string {
	if len(data) <= 0 {
		return nil
	}

	locations := make([]string, 0, len(data))

	for _, l := range data {
		locations = append(locations, l.String())
	}

	return locations
}

func rentalMaterialsRPCToMaterials(data []rentalRPC.MaterialEnum) []rental.Material {
	if len(data) <= 0 {
		return nil
	}

	materials := make([]rental.Material, 0, len(data))

	for _, m := range data {
		materials = append(materials, rentalMaterialRPCToMaterial(m))
	}

	return materials
}

func rentalPropertyTypesRPCToPropertyTypes(data []rentalRPC.PropertyTypeEnum) []rental.PropertyType {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]rental.PropertyType, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalPropertyTypeRPCToPropertyType(p))
	}

	return properties
}

func rentalStatusesRPCToStatuses(data []rentalRPC.StatusEnum) []rental.Status {
	if len(data) <= 0 {
		return nil
	}

	statuses := make([]rental.Status, 0, len(data))

	for _, s := range data {
		statuses = append(statuses, rentalStatusRPCToStatus(s))
	}

	return statuses
}

func rentalCommercialLocationsRPCLocations(data []rentalRPC.CommericalPropertyLocationEnum) []rental.CommericalPropertyLocation {
	if len(data) <= 0 {
		return nil
	}

	locations := make([]rental.CommericalPropertyLocation, 0, len(data))

	for _, p := range data {
		locations = append(locations, rentalCommercialLocationRPCLocation(p))
	}

	return locations
}

func rentalCommercialPropertiesRPCCommercialProperties(data []rentalRPC.CommercialPropertyEnum) []rental.CommericalProperty {

	if len(data) <= 0 {
		return nil
	}

	properties := make([]rental.CommericalProperty, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalCommercialPropertyRPCCommercialProperty(p))
	}

	return properties
}

func rentalAddtionaFiltersRPCToAddiotionalFilters(data []rentalRPC.AdditionalFiltersEnum) []rental.AdditionalFilters {
	if len(data) <= 0 {
		return nil
	}

	filters := make([]rental.AdditionalFilters, 0, len(data))

	for _, f := range data {
		filters = append(filters, rentalAddtionaFilterRPCToAddiotionalFilter(f))
	}

	return filters
}
func rentalWhoLivesRPCToWhoLive(data []rentalRPC.WhoLiveEnum) []rental.WhoLive {
	if len(data) <= 0 {
		return nil
	}

	lives := make([]rental.WhoLive, 0, len(data))

	for _, l := range data {
		lives = append(lives, rentalWhoLiveRPCToWhoLive(l))
	}

	return lives
}

func rentalOfPropertiesRPCToTypeOfProperties(data []rentalRPC.TypeOfPropertyEnum) []rental.TypeOfProperty {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]rental.TypeOfProperty, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalOfPropertyRPCToTypeOfProperty(p))
	}

	return properties
}

func rentalDetailsRPCToDetails(data []*rentalRPC.Detail) []rental.Detail {
	if len(data) <= 0 {
		return nil
	}

	details := make([]rental.Detail, 0, len(data))

	for i, d := range data {
		details = append(details, rental.Detail{
			Title:       d.GetTitle(),
			HouseRules:  d.GetHouseRules(),
			Description: d.GetDescription(),
		})

		id := details[i].GenerateID()
		details[i].SetID(id)
	}

	return details
}

func rentalClimatControlsRPCToClimatControls(data []rentalRPC.ClimatControlEnum) []rental.ClimatControl {
	if len(data) <= 0 {
		return nil
	}

	controls := make([]rental.ClimatControl, 0, len(data))

	for _, c := range data {
		controls = append(controls, rentalClimatControlRPCToClimatControl(c))
	}

	return controls
}

func rentalPriceRPCToPrice(data *rentalRPC.Price) (r rental.Price) {
	if data == nil {
		return
	}

	return rental.Price{
		MinPrice:  data.GetMinPrice(),
		MaxPrice:  data.GetMaxPrice(),
		Currency:  data.GetCurrency(),
		PriceType: rentalPriceTypeRPCToPriceType(data.GetPriceType()),
	}
}

func rentalPhonesRPCToPhones(data []*rentalRPC.Phone) []rental.Phone {
	if data == nil || len(data) <= 0 {
		return nil
	}

	phones := make([]rental.Phone, 0, len(data))

	for i, p := range data {
		phones = append(phones, rental.Phone{
			CountryCode: p.GetCountryCode(),
			Number:      p.GetNumber(),
		})

		id := phones[i].GenerateID()
		phones[i].SetID(id)
	}

	return phones
}

func rentalOutdoorFeaturesRPCtoOutdoorFeatures(data []rentalRPC.OutdoorFeaturesEnum) []rental.OutdoorFeatures {
	if len(data) <= 0 {
		return nil
	}

	features := make([]rental.OutdoorFeatures, 0, len(data))

	for _, f := range data {
		features = append(features, rentalOutdoorFeatureRPCtoOutdoorFeature(f))
	}

	return features
}
func rentalIndoorFeaturesRPCtoIndoorFeatures(data []rentalRPC.IndoorFeaturesEnum) []rental.IndoorFeatures {
	if len(data) <= 0 {
		return nil
	}

	features := make([]rental.IndoorFeatures, 0, len(data))

	for _, f := range data {
		features = append(features, rentalIndoorFeatureRPCtoIndoorFeature(f))
	}

	return features
}

func rentalInfoRPCToRentalInfo(data *rentalRPC.Rental) rental.CommonRental {
	if data == nil {
		return rental.CommonRental{}
	}

	return rental.CommonRental{
		ExpiredDays:  data.GetExpiredDays(),
		IsUrgent:     data.GetIsUrgent(),
		PostCurrency: data.GetPostCurrency(),
		PropertyType: rentalPropertyTypeRPCToPropertyType(data.GetPropertyType()),
		DealType:     rentalDealTypeRPCToDealType(data.GetDealType()),
		Location:     rentalLocationRPCToLocation(data.GetLocation()),
	}
}

func rentalLocationRPCToLocation(data *rentalRPC.Location) location.Location {
	if data == nil {
		return location.Location{}
	}

	return location.Location{
		Country: &location.Country{},
		City: &location.City{
			ID: data.GetCity().GetId(),
		},
		Street:  stringToNull(data.GetStreet()),
		Address: stringToNull(data.GetAddress()),
	}
}

func rentalWhoLiveRPCToWhoLive(data rentalRPC.WhoLiveEnum) rental.WhoLive {
	switch data {
	case rentalRPC.WhoLiveEnum_WhoLive_Mortgagor:
		return rental.WhoLiveMortgagor
	case rentalRPC.WhoLiveEnum_WhoLiveOwner:
		return rental.WhoLiveOwner
	}

	return rental.WhoLiveAny
}

func rentalWhoLiveToRPC(data rental.WhoLive) rentalRPC.WhoLiveEnum {
	switch data {
	case rental.WhoLiveMortgagor:
		return rentalRPC.WhoLiveEnum_WhoLive_Mortgagor
	case rental.WhoLiveOwner:
		return rentalRPC.WhoLiveEnum_WhoLiveOwner
	}

	return rentalRPC.WhoLiveEnum_WhoLive_Any
}

func rentalOfPropertyRPCToTypeOfProperty(data rentalRPC.TypeOfPropertyEnum) rental.TypeOfProperty {
	switch data {
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Appartaments:
		return rental.TypeOfPropertyAppartaments
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Houses:
		return rental.TypeOfPropertyHouses
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_CountryHomes:
		return rental.TypeOfPropertyCountryHomes
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Duplex:
		return rental.TypeOfPropertyDuplex
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Penthouses:
		return rental.TypeOfPropertyPenthouses
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Bungalow:
		return rental.TypeOfPropertyBungalow
	}
	return rental.TypeOfPropertyAny
}

func rentalOfPropertyToRPC(data rental.TypeOfProperty) rentalRPC.TypeOfPropertyEnum {
	switch data {
	case rental.TypeOfPropertyAppartaments:
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Appartaments
	case rental.TypeOfPropertyHouses:
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Houses
	case rental.TypeOfPropertyCountryHomes:
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_CountryHomes
	case rental.TypeOfPropertyDuplex:
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Duplex
	case rental.TypeOfPropertyPenthouses:
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Penthouses
	case rental.TypeOfPropertyBungalow:
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Bungalow
	}
	return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Any
}

func rentalClimatControlRPCToClimatControl(data rentalRPC.ClimatControlEnum) rental.ClimatControl {
	switch data {
	case rentalRPC.ClimatControlEnum_ClimatControl_AirConditioning:
		return rental.ClimatControlAirConditioning
	case rentalRPC.ClimatControlEnum_ClimatControl_Hearting:
		return rental.ClimatControlHearting
	case rentalRPC.ClimatControlEnum_ClimatControl_WaterTank:
		return rental.ClimatControlWaterTank
	case rentalRPC.ClimatControlEnum_ClimatControl_SolarPanels:
		return rental.ClimatControlSolarPanels
	case rentalRPC.ClimatControlEnum_ClimatControl_HighEnergyEfficiency:
		return rental.ClimatControlHighEnergyEfficiency
	case rentalRPC.ClimatControlEnum_ClimatControlSolarHotWater:
		return rental.ClimatControlSolarHotWater
	case rentalRPC.ClimatControlEnum_ClimatControl_ZonalHeating:
		return rental.ClimatControlZonalHeating
	case rentalRPC.ClimatControlEnum_ClimatControl_HeatPumps:
		return rental.ClimatControlHeatPumps
	}
	return rental.ClimatControlAny
}

func rentalClimatControlToRPC(data rental.ClimatControl) rentalRPC.ClimatControlEnum {
	switch data {
	case rental.ClimatControlAirConditioning:
		return rentalRPC.ClimatControlEnum_ClimatControl_AirConditioning
	case rental.ClimatControlHearting:
		return rentalRPC.ClimatControlEnum_ClimatControl_Hearting
	case rental.ClimatControlWaterTank:
		return rentalRPC.ClimatControlEnum_ClimatControl_WaterTank
	case rental.ClimatControlSolarPanels:
		return rentalRPC.ClimatControlEnum_ClimatControl_SolarPanels
	case rental.ClimatControlHighEnergyEfficiency:
		return rentalRPC.ClimatControlEnum_ClimatControl_HighEnergyEfficiency
	case rental.ClimatControlSolarHotWater:
		return rentalRPC.ClimatControlEnum_ClimatControlSolarHotWater
	case rental.ClimatControlZonalHeating:
		return rentalRPC.ClimatControlEnum_ClimatControl_ZonalHeating
	case rental.ClimatControlHeatPumps:
		return rentalRPC.ClimatControlEnum_ClimatControl_HeatPumps
	}
	return rentalRPC.ClimatControlEnum_ClimatControl_Any
}

func rentalStatusRPCToStatus(data rentalRPC.StatusEnum) rental.Status {
	switch data {
	case rentalRPC.StatusEnum_Status_OldBuild:
		return rental.StatusOldBuild
	case rentalRPC.StatusEnum_Status_NewBuilding:
		return rental.StatusNewBuilding
	case rentalRPC.StatusEnum_Status_UnderConstruction:
		return rental.StatusUnderConstruction
	case rentalRPC.StatusEnum_StatusDeveloped:
		return rental.StatusDeveloped
	case rentalRPC.StatusEnum_Status_Buildable:
		return rental.StatusBuildable
	case rentalRPC.StatusEnum_Status_NonBuilding:
		return rental.StatusNonBuilding
	}
	return rental.StatusAny
}

func rentalStatusToRPC(data rental.Status) rentalRPC.StatusEnum {
	switch data {
	case rental.StatusOldBuild:
		return rentalRPC.StatusEnum_Status_OldBuild
	case rental.StatusNewBuilding:
		return rentalRPC.StatusEnum_Status_NewBuilding
	case rental.StatusUnderConstruction:
		return rentalRPC.StatusEnum_Status_UnderConstruction
	case rental.StatusDeveloped:
		return rentalRPC.StatusEnum_StatusDeveloped
	case rental.StatusBuildable:
		return rentalRPC.StatusEnum_Status_Buildable
	case rental.StatusNonBuilding:
		return rentalRPC.StatusEnum_Status_NonBuilding
	}
	return rentalRPC.StatusEnum_Status_Any
}

func rentalPriceTypeRPCToPriceType(data rentalRPC.PriceTypeEnum) rental.PriceType {

	switch data {
	case rentalRPC.PriceTypeEnum_PriceType_MetreSquare:
		return rental.PriceTypeMetreSquare
	case rentalRPC.PriceTypeEnum_PriceType_FeetSquare:
		return rental.PriceTypeFeetSquare
	case rentalRPC.PriceTypeEnum_PriceType_Total:
		return rental.PriceTypeTotal

	}
	return rental.PriceTypeAny
}

func rentalPriceTypeToRPC(data rental.PriceType) rentalRPC.PriceTypeEnum {

	switch data {
	case rental.PriceTypeMetreSquare:
		return rentalRPC.PriceTypeEnum_PriceType_MetreSquare
	case rental.PriceTypeTotal:
		return rentalRPC.PriceTypeEnum_PriceType_Total
	}
	return rentalRPC.PriceTypeEnum_PriceType_Any
}
func rentalOutdoorFeatureRPCtoOutdoorFeature(data rentalRPC.OutdoorFeaturesEnum) rental.OutdoorFeatures {
	switch data {
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_SwimmingPool:
		return rental.OutdoorFeaturesSwimmingPool
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Balcony:
		return rental.OutdoorFeaturesBalcony
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_UndercoverParking:
		return rental.OutdoorFeaturesUndercoverParking
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_FullyFenced:
		return rental.OutdoorFeaturesFullyFenced
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_TennisCourt:
		return rental.OutdoorFeaturesTennisCourt
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garden:
		return rental.OutdoorFeaturesGarden
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garage:
		return rental.OutdoorFeaturesGarage
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorArea:
		return rental.OutdoorFeaturesOutdoorArea
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Shed:
		return rental.OutdoorFeaturesShed
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorSpa:
		return rental.OutdoorFeaturesOutdoorSpa
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Outbuildings:
		return rental.OutdoorFeaturesOutbuildings
	}
	return rental.OutdoorFeaturesAny
}

func rentalOutdoorFeatureToRPC(data rental.OutdoorFeatures) rentalRPC.OutdoorFeaturesEnum {
	switch data {
	case rental.OutdoorFeaturesSwimmingPool:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_SwimmingPool
	case rental.OutdoorFeaturesBalcony:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Balcony
	case rental.OutdoorFeaturesUndercoverParking:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_UndercoverParking
	case rental.OutdoorFeaturesFullyFenced:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_FullyFenced
	case rental.OutdoorFeaturesTennisCourt:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_TennisCourt
	case rental.OutdoorFeaturesGarden:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garden
	case rental.OutdoorFeaturesGarage:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garage
	case rental.OutdoorFeaturesOutdoorArea:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorArea
	case rental.OutdoorFeaturesShed:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Shed
	case rental.OutdoorFeaturesOutdoorSpa:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorSpa
	case rental.OutdoorFeaturesOutbuildings:
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Outbuildings
	}
	return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Any
}

func rentalIndoorFeatureRPCtoIndoorFeature(data rentalRPC.IndoorFeaturesEnum) rental.IndoorFeatures {

	switch data {
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Ensuit:
		return rental.IndoorFeaturesEnsuit
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Study:
		return rental.IndoorFeaturesStudy
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_AlarmSystem:
		return rental.IndoorFeaturesAlarmSystem
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Floorboards:
		return rental.IndoorFeaturesFloorboards
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_RumpusRoom:
		return rental.IndoorFeaturesRumpusRoom
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_StorageRoom:
		return rental.IndoorFeaturesStorageRoom
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Dishwasher:
		return rental.IndoorFeaturesDishwasher
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Lift:
		return rental.IndoorFeaturesLift
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_BuiltInRobes:
		return rental.IndoorFeaturesBuiltInRobes
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Broadband:
		return rental.IndoorFeaturesBroadband
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Gym:
		return rental.IndoorFeaturesGym
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Workshop:
		return rental.IndoorFeaturesWorkshop
	}
	return rental.IndoorFeaturesAny
}

func rentalIndoorFeatureToRPC(data rental.IndoorFeatures) rentalRPC.IndoorFeaturesEnum {

	switch data {
	case rental.IndoorFeaturesEnsuit:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Ensuit
	case rental.IndoorFeaturesStudy:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Study
	case rental.IndoorFeaturesAlarmSystem:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_AlarmSystem
	case rental.IndoorFeaturesFloorboards:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Floorboards
	case rental.IndoorFeaturesRumpusRoom:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_RumpusRoom
	case rental.IndoorFeaturesStorageRoom:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_StorageRoom
	case rental.IndoorFeaturesDishwasher:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Dishwasher
	case rental.IndoorFeaturesLift:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Lift
	case rental.IndoorFeaturesBuiltInRobes:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_BuiltInRobes
	case rental.IndoorFeaturesBroadband:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Broadband
	case rental.IndoorFeaturesGym:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Gym
	case rental.IndoorFeaturesWorkshop:
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Workshop
	}
	return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Any
}

func rentalPropertyTypeRPCToPropertyType(data rentalRPC.PropertyTypeEnum) rental.PropertyType {
	switch data {
	case rentalRPC.PropertyTypeEnum_PropertyType_All:
		return rental.PropertyTypeAll
	case rentalRPC.PropertyTypeEnum_PropertyType_NewHomes:
		return rental.PropertyTypeNewHomes
	case rentalRPC.PropertyTypeEnum_PropertyType_Homes:
		return rental.PropertyTypeHomes
	case rentalRPC.PropertyTypeEnum_PropertyType_Houses:
		return rental.PropertyTypeHouses
	case rentalRPC.PropertyTypeEnum_PropertyType_Appartments:
		return rental.PropertyTypeAppartments
	case rentalRPC.PropertyTypeEnum_PropertyType_Garages:
		return rental.PropertyTypeGarages
	case rentalRPC.PropertyTypeEnum_PropertyType_StorageRooms:
		return rental.PropertyTypeStorageRooms
	case rentalRPC.PropertyTypeEnum_PropertyType_Offices:
		return rental.PropertyTypeOffices
	case rentalRPC.PropertyTypeEnum_PropertyType_CommercialProperties:
		return rental.PropertyTypeCommercialProperties
	case rentalRPC.PropertyTypeEnum_PropertyType_Buildings:
		return rental.PropertyTypeBuildings
	case rentalRPC.PropertyTypeEnum_PropertyType_Land:
		return rental.PropertyTypeLand
	case rentalRPC.PropertyTypeEnum_PropertyType_BareLand:
		return rental.PropertyTypeBareLand
	case rentalRPC.PropertyTypeEnum_PropertyType_Barn:
		return rental.PropertyTypeBarn
	case rentalRPC.PropertyTypeEnum_PropertyType_SummerCottage:
		return rental.PropertyTypeSummerCottage
	case rentalRPC.PropertyTypeEnum_PropertyType_RuralFarm:
		return rental.PropertyTypeRuralFarm
	case rentalRPC.PropertyTypeEnum_PropertyType_HotelRoom:
		return rental.PropertyTypeHotelRoom
	}

	return rental.PropertyTypeAny
}

func rentalPropertyTypeToRPC(data rental.PropertyType) rentalRPC.PropertyTypeEnum {
	switch data {
	case rental.PropertyTypeAll:
		return rentalRPC.PropertyTypeEnum_PropertyType_All
	case rental.PropertyTypeNewHomes:
		return rentalRPC.PropertyTypeEnum_PropertyType_NewHomes
	case rental.PropertyTypeHomes:
		return rentalRPC.PropertyTypeEnum_PropertyType_Homes
	case rental.PropertyTypeHouses:
		return rentalRPC.PropertyTypeEnum_PropertyType_Houses
	case rental.PropertyTypeAppartments:
		return rentalRPC.PropertyTypeEnum_PropertyType_Appartments
	case rental.PropertyTypeGarages:
		return rentalRPC.PropertyTypeEnum_PropertyType_Garages
	case rental.PropertyTypeStorageRooms:
		return rentalRPC.PropertyTypeEnum_PropertyType_StorageRooms
	case rental.PropertyTypeOffices:
		return rentalRPC.PropertyTypeEnum_PropertyType_Offices
	case rental.PropertyTypeCommercialProperties:
		return rentalRPC.PropertyTypeEnum_PropertyType_CommercialProperties
	case rental.PropertyTypeBuildings:
		return rentalRPC.PropertyTypeEnum_PropertyType_Buildings
	case rental.PropertyTypeLand:
		return rentalRPC.PropertyTypeEnum_PropertyType_Land
	case rental.PropertyTypeBareLand:
		return rentalRPC.PropertyTypeEnum_PropertyType_BareLand
	case rental.PropertyTypeBarn:
		return rentalRPC.PropertyTypeEnum_PropertyType_Barn
	case rental.PropertyTypeSummerCottage:
		return rentalRPC.PropertyTypeEnum_PropertyType_SummerCottage
	case rental.PropertyTypeRuralFarm:
		return rentalRPC.PropertyTypeEnum_PropertyType_RuralFarm
	case rental.PropertyTypeHotelRoom:
		return rentalRPC.PropertyTypeEnum_PropertyType_HotelRoom
	}

	return rentalRPC.PropertyTypeEnum_PropertyType_Any
}

func rentalDealTypeRPCToDealType(data rentalRPC.DealTypeEnum) rental.DealType {

	switch data {
	case rentalRPC.DealTypeEnum_DealType_Lease:
		return rental.DealTypeLease
	case rentalRPC.DealTypeEnum_DealType_Rent:
		return rental.DealTypeRent
	case rentalRPC.DealTypeEnum_DealType_Sell:
		return rental.DealTypeSell
	case rentalRPC.DealTypeEnum_DealType_Share:
		return rental.DealTypeShare
	case rentalRPC.DealTypeEnum_DealType_Build:
		return rental.DealTypeBuild
	case rentalRPC.DealTypeEnum_DealType_Materials:
		return rental.DealTypeMaterials
	case rentalRPC.DealTypeEnum_DealType_Renovation:
		return rental.DealTypeRenovation
	case rentalRPC.DealTypeEnum_DealType_Move:
		return rental.DealTypeMove
	}

	return rental.DealTypeAny
}

func rentalDealTypeToRPC(data rental.DealType) rentalRPC.DealTypeEnum {

	switch data {
	case rental.DealTypeLease:
		return rentalRPC.DealTypeEnum_DealType_Lease
	case rental.DealTypeRent:
		return rentalRPC.DealTypeEnum_DealType_Rent
	case rental.DealTypeSell:
		return rentalRPC.DealTypeEnum_DealType_Sell
	case rental.DealTypeShare:
		return rentalRPC.DealTypeEnum_DealType_Share
	case rental.DealTypeBuild:
		return rentalRPC.DealTypeEnum_DealType_Build
	case rental.DealTypeMaterials:
		return rentalRPC.DealTypeEnum_DealType_Materials
	case rental.DealTypeRenovation:
		return rentalRPC.DealTypeEnum_DealType_Renovation
	case rental.DealTypeMove:
		return rentalRPC.DealTypeEnum_DealType_Move
	}

	return rentalRPC.DealTypeEnum_DealType_Any
}

func rentalAddtionaFilterRPCToAddiotionalFilter(data rentalRPC.AdditionalFiltersEnum) rental.AdditionalFilters {

	switch data {
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_AireConditioning:
		return rental.AdditionalFiltersAireConditioning
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Electricity:
		return rental.AdditionalFiltersElectricity
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_AutomaticDoor:
		return rental.AdditionalFiltersAutomaticDoor
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Heating:
		return rental.AdditionalFiltersHeating
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_MotoBikeGarage:
		return rental.AdditionalFiltersMotoBikeGarage
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_NaturalGas:
		return rental.AdditionalFiltersNaturalGas
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_OnCorner:
		return rental.AdditionalFiltersOnCorner
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_SecuritySystem:
		return rental.AdditionalFiltersSecuritySystem
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Sewage:
		return rental.AdditionalFiltersSewage
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_SmokeExtractor:
		return rental.AdditionalFiltersSmokeExtractor
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Water:
		return rental.AdditionalFiltersWater
	}

	return rental.AdditionalFiltersAny
}

func rentalAddtionaFilterToRPC(data rental.AdditionalFilters) rentalRPC.AdditionalFiltersEnum {

	switch data {
	case rental.AdditionalFiltersAireConditioning:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_AireConditioning
	case rental.AdditionalFiltersElectricity:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Electricity
	case rental.AdditionalFiltersAutomaticDoor:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_AutomaticDoor
	case rental.AdditionalFiltersHeating:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Heating
	case rental.AdditionalFiltersMotoBikeGarage:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_MotoBikeGarage
	case rental.AdditionalFiltersNaturalGas:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_NaturalGas
	case rental.AdditionalFiltersOnCorner:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_OnCorner
	case rental.AdditionalFiltersSecuritySystem:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_SecuritySystem
	case rental.AdditionalFiltersSewage:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Sewage
	case rental.AdditionalFiltersSmokeExtractor:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_SmokeExtractor
	case rental.AdditionalFiltersWater:
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Water
	}

	return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Any
}

func rentalCommercialPropertyRPCCommercialProperty(data rentalRPC.CommercialPropertyEnum) rental.CommericalProperty {

	switch data {
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_Basement:
		return rental.CommericalPropertyBasement
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_CommercialPremises:
		return rental.CommericalPropertyCommercialPremises
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_FoodFacility:
		return rental.CommericalPropertyFoodFacility
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_Garage:
		return rental.CommericalPropertyGarage
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_IndustrialBuilding:
		return rental.CommericalPropertyIndustrialBuilding
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_OfficeSpace:
		return rental.CommericalPropertyOfficeSpace
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_TradingPlace:
		return rental.CommericalPropertyTradingPlace
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_Warehouse:
		return rental.CommericalPropertyWarehouse
	}
	return rental.CommericalPropertyAny
}

func rentalCommercialPropertyToRPC(data rental.CommericalProperty) rentalRPC.CommercialPropertyEnum {

	switch data {
	case rental.CommericalPropertyBasement:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_Basement
	case rental.CommericalPropertyCommercialPremises:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_CommercialPremises
	case rental.CommericalPropertyFoodFacility:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_FoodFacility
	case rental.CommericalPropertyGarage:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_Garage
	case rental.CommericalPropertyIndustrialBuilding:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_IndustrialBuilding
	case rental.CommericalPropertyOfficeSpace:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_OfficeSpace
	case rental.CommericalPropertyTradingPlace:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_TradingPlace
	case rental.CommericalPropertyWarehouse:
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_Warehouse
	}
	return rentalRPC.CommercialPropertyEnum_CommericalProperty_Any
}

func rentalCommercialLocationRPCLocation(data rentalRPC.CommericalPropertyLocationEnum) rental.CommericalPropertyLocation {
	switch data {
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_BelowGround:
		return rental.CommericalPropertyLocationBelowGround
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_InShoppingCentre:
		return rental.CommericalPropertyLocationInShoppingCentre
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Indifferent:
		return rental.CommericalPropertyLocationIndifferent
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Mezzanine:
		return rental.CommericalPropertyLocationMezzanine
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Other:
		return rental.CommericalPropertyLocationOther
	}
	return rental.CommericalPropertyLocationAny
}

func rentalCommercialLocationToRPC(data rental.CommericalPropertyLocation) rentalRPC.CommericalPropertyLocationEnum {
	switch data {
	case rental.CommericalPropertyLocationBelowGround:
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_BelowGround
	case rental.CommericalPropertyLocationInShoppingCentre:
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_InShoppingCentre
	case rental.CommericalPropertyLocationIndifferent:
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Indifferent
	case rental.CommericalPropertyLocationMezzanine:
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Mezzanine
	case rental.CommericalPropertyLocationOther:
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Other
	}
	return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Any
}

func rentalLayoutRPCToLayout(data rentalRPC.LayoutEnum) rental.Layout {
	switch data {
	case rentalRPC.LayoutEnum_Layout_Indifferent:
		return rental.LayoutIndifferent
	case rentalRPC.LayoutEnum_Layout_OpenPlan:
		return rental.LayoutOpenPlan
	case rentalRPC.LayoutEnum_Layout_Walls:
		return rental.LayoutWalls
	}
	return rental.LayoutAny
}

func rentalLayoutToRPC(data rental.Layout) rentalRPC.LayoutEnum {
	switch data {
	case rental.LayoutIndifferent:
		return rentalRPC.LayoutEnum_Layout_Indifferent
	case rental.LayoutOpenPlan:
		return rentalRPC.LayoutEnum_Layout_OpenPlan
	case rental.LayoutWalls:
		return rentalRPC.LayoutEnum_Layout_Walls
	}
	return rentalRPC.LayoutEnum_Layout_Any
}

func rentalBuildinUseRPCToBuildinUse(data rentalRPC.BuildingUseEnum) rental.BuildingUse {

	switch data {
	case rentalRPC.BuildingUseEnum_Building_Use_Indifferent:
		return rental.BuildingUseIndifferent
	case rentalRPC.BuildingUseEnum_Building_Use_OnlyOffice:
		return rental.BuildingUseOnlyOffice
	case rentalRPC.BuildingUseEnum_Building_Use_Mixed:
		return rental.BuildingUseMixed
	}
	return rental.BuildingUseAny
}

func rentalBuildinUseToRPC(data rental.BuildingUse) rentalRPC.BuildingUseEnum {

	switch data {
	case rental.BuildingUseIndifferent:
		return rentalRPC.BuildingUseEnum_Building_Use_Indifferent
	case rental.BuildingUseOnlyOffice:
		return rentalRPC.BuildingUseEnum_Building_Use_OnlyOffice
	case rental.BuildingUseMixed:
		return rentalRPC.BuildingUseEnum_Building_Use_Mixed
	}
	return rentalRPC.BuildingUseEnum_Building_Use_Any
}

func rentalTimingRPCToStruct(data rentalRPC.TimingEnum) rental.Timing {
	switch data {
	case rentalRPC.TimingEnum_Timing_Flexible:
		return rental.TimingFlexible
	case rentalRPC.TimingEnum_Timing_6Months:
		return rental.Timing6Months
	case rentalRPC.TimingEnum_Timing_Year:
		return rental.TimingYear
	}

	return rental.AnyTiming
}

func rentalTimingToRPC(data rental.Timing) rentalRPC.TimingEnum {
	switch data {
	case rental.TimingFlexible:
		return rentalRPC.TimingEnum_Timing_Flexible
	case rental.Timing6Months:
		return rentalRPC.TimingEnum_Timing_6Months
	case rental.TimingYear:
		return rentalRPC.TimingEnum_Timing_Year
	}

	return rentalRPC.TimingEnum_Any_Timing
}

func nullToString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func stringToNull(str string) *string {
	if str == "" {
		return nil
	}

	return &str
}

func rentalMaterialRPCToMaterial(data rentalRPC.MaterialEnum) rental.Material {

	switch data {
	case rentalRPC.MaterialEnum_Material_Lumber_Composites:
		return rental.MaterialLumberComposites
	case rentalRPC.MaterialEnum_Material_Fencing:
		return rental.MaterialFencing
	case rentalRPC.MaterialEnum_Material_Decking:
		return rental.MaterialDecking
	case rentalRPC.MaterialEnum_Material_Fastners:
		return rental.MaterialFastners
	case rentalRPC.MaterialEnum_Material_Moulding_Millwork:
		return rental.MaterialMouldingMillwork
	case rentalRPC.MaterialEnum_Material_Paint:
		return rental.MaterialPaint
	case rentalRPC.MaterialEnum_Material_Drywall:
		return rental.MaterialDrywall
	case rentalRPC.MaterialEnum_Material_Doors_Windows:
		return rental.MaterialDoorsWindows
	case rentalRPC.MaterialEnum_Material_Roofing_Gutters:
		return rental.MaterialRoofingGutters
	case rentalRPC.MaterialEnum_Material_Ladders:
		return rental.MaterialLadders
	case rentalRPC.MaterialEnum_Material_Scaffolding:
		return rental.MaterialScaffolding
	case rentalRPC.MaterialEnum_Material_Plumbing:
		return rental.MaterialPlumbing
	case rentalRPC.MaterialEnum_Material_Siding:
		return rental.MaterialSiding
	case rentalRPC.MaterialEnum_Material_Insulation:
		return rental.MaterialInsulation
	case rentalRPC.MaterialEnum_Material_Ceilings:
		return rental.MaterialCeilings
	case rentalRPC.MaterialEnum_Material_Wall_Paneling:
		return rental.MaterialWallPaneling
	case rentalRPC.MaterialEnum_Material_Flooring:
		return rental.MaterialFlooring
	case rentalRPC.MaterialEnum_Material_Concrete_Cement_Masonry:
		return rental.MaterialConcreteCementMasonry
	case rentalRPC.MaterialEnum_Material_Material_Handling_Equipment:
		return rental.MaterialMaterialHandlingEquipment
	case rentalRPC.MaterialEnum_Material_Glass_and_Plastic_Sheets:
		return rental.MaterialGlassandPlasticSheets
	case rentalRPC.MaterialEnum_Material_Building_Hardware:
		return rental.MaterialBuildingHardware
	case rentalRPC.MaterialEnum_Material_Heating_venting_Cooling:
		return rental.MaterialHeatingventingCooling
	case rentalRPC.MaterialEnum_Material_Other:
		return rental.MaterialOther
	}

	return rental.AnyMaterial
}

func rentalMaterialToRPC(data rental.Material) rentalRPC.MaterialEnum {

	switch data {
	case rental.MaterialLumberComposites:
		return rentalRPC.MaterialEnum_Material_Lumber_Composites
	case rental.MaterialFencing:
		return rentalRPC.MaterialEnum_Material_Fencing
	case rental.MaterialDecking:
		return rentalRPC.MaterialEnum_Material_Decking
	case rental.MaterialFastners:
		return rentalRPC.MaterialEnum_Material_Fastners
	case rental.MaterialMouldingMillwork:
		return rentalRPC.MaterialEnum_Material_Moulding_Millwork
	case rental.MaterialPaint:
		return rentalRPC.MaterialEnum_Material_Paint
	case rental.MaterialDrywall:
		return rentalRPC.MaterialEnum_Material_Drywall
	case rental.MaterialDoorsWindows:
		return rentalRPC.MaterialEnum_Material_Doors_Windows
	case rental.MaterialRoofingGutters:
		return rentalRPC.MaterialEnum_Material_Roofing_Gutters
	case rental.MaterialLadders:
		return rentalRPC.MaterialEnum_Material_Ladders
	case rental.MaterialScaffolding:
		return rentalRPC.MaterialEnum_Material_Scaffolding
	case rental.MaterialPlumbing:
		return rentalRPC.MaterialEnum_Material_Plumbing
	case rental.MaterialSiding:
		return rentalRPC.MaterialEnum_Material_Siding
	case rental.MaterialInsulation:
		return rentalRPC.MaterialEnum_Material_Insulation
	case rental.MaterialCeilings:
		return rentalRPC.MaterialEnum_Material_Ceilings
	case rental.MaterialWallPaneling:
		return rentalRPC.MaterialEnum_Material_Wall_Paneling
	case rental.MaterialFlooring:
		return rentalRPC.MaterialEnum_Material_Flooring
	case rental.MaterialConcreteCementMasonry:
		return rentalRPC.MaterialEnum_Material_Concrete_Cement_Masonry
	case rental.MaterialMaterialHandlingEquipment:
		return rentalRPC.MaterialEnum_Material_Material_Handling_Equipment
	case rental.MaterialGlassandPlasticSheets:
		return rentalRPC.MaterialEnum_Material_Glass_and_Plastic_Sheets
	case rental.MaterialBuildingHardware:
		return rentalRPC.MaterialEnum_Material_Building_Hardware
	case rental.MaterialHeatingventingCooling:
		return rentalRPC.MaterialEnum_Material_Heating_venting_Cooling
	case rental.MaterialOther:
		return rentalRPC.MaterialEnum_Material_Other
	}

	return rentalRPC.MaterialEnum_Any_Material
}

func rentalServiceRPCToService(data rentalRPC.ServiceEnum) rental.Service {
	switch data {
	case rentalRPC.ServiceEnum_Service_Auto_Transport:
		return rental.ServiceAutoTransport
	case rentalRPC.ServiceEnum_Service_Storage:
		return rental.ServiceStorage
	case rentalRPC.ServiceEnum_Service_Moving_Supplies:
		return rental.ServiceMovingSupplies
	case rentalRPC.ServiceEnum_Service_Furniture_Movers:
		return rental.ServiceFurnitureMovers
	}
	return rental.AnyService
}

func rentalServiceToRPC(data rental.Service) rentalRPC.ServiceEnum {
	switch data {
	case rental.ServiceAutoTransport:
		return rentalRPC.ServiceEnum_Service_Auto_Transport
	case rental.ServiceStorage:
		return rentalRPC.ServiceEnum_Service_Storage
	case rental.ServiceMovingSupplies:
		return rentalRPC.ServiceEnum_Service_Moving_Supplies
	case rental.ServiceFurnitureMovers:
		return rentalRPC.ServiceEnum_Service_Furniture_Movers
	}
	return rentalRPC.ServiceEnum_Any_Service
}
