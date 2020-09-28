package resolver

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/rentalRPC"
)

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

func realEstatesRPCToEstates(data []*rentalRPC.Estate) []RealEstateInterface {
	if len(data) <= 0 {
		return nil
	}

	estates := make([]RealEstateInterface, 0, len(data))

	for _, e := range data {
		// Appartament
		if k := e.GetAppartaments(); k != nil {
			estates = append(estates, AppartmentResolver{
				R: rentalAppartamentRPCToStruct(k),
			})
		}
		// Storage room
		if s := e.GetStorageRoom(); s != nil {
			estates = append(estates, StorageRoomResolver{
				R: rentalStorageRoomRPCToStruct(s),
			})
		}
		// Office
		if o := e.GetOffice(); o != nil {
			estates = append(estates, OfficeResolver{
				R: rentalOfficeRPCToStruct(o),
			})
		}
		// CommercialProperties
		if cp := e.GetCommercialAndRuralFarm(); cp != nil {
			estates = append(estates, CommericalPropertiesResolver{
				R: rentalCommercilaPropertiesRPCToStruct(cp),
			})
		}

	}

	return estates
}

func realEstateToRPC(data AddRealEstateInput) *rentalRPC.AddRentalRequest {
	res := &rentalRPC.AddRentalRequest{
		CompanyID:           NullToString(data.Company_id),
		AvailibatiFrom:      NullToString(data.Availability_from),
		AvailibatiTo:        NullToString(data.Availability_to),
		TotalArea:           NullToInt32(data.Total_area),
		BadRooms:            NullToInt32(data.Badrooms),
		BathRooms:           NullToInt32(data.Bathrooms),
		CarSpaces:           NullToInt32(data.Car_spaces),
		Rooms:               NullToInt32(data.Rooms),
		Floor:               NullToInt32(data.Floor),
		Floors:              NullToInt32(data.Floors),
		Rental:              rentalInfoToRPC(data.Rental_info),
		Details:             rentalDetailsToRPC(data.Rental_detail),
		Phones:              rentalPhonesToRPC(data.Phones),
		Price:               rentalPriceToRPC(data.Price),
		BuildingUse:         rentalBuilingUseToRPC(data.Building_use),
		CityIDs:             NullStringArrayToStringArray(data.City_ids),
		CountryIDs:          NullStringArrayToStringArray(data.Country_ids),
		Status:              rentalStatusToRPC(data.Status),
		WhoLive:             rentalWhoLivesToRPC(data.Who_live),
		Exterior:            rentalPriceToRPC(data.Exterior),
		Interior:            rentalPriceToRPC(data.Interior),
		Purchase:            rentalPriceToRPC(data.Purchase),
		InteriorAndExterior: rentalPriceToRPC(data.Interior_and_exterior),
		LocationType:        rentalLocationTypesToRPC(data.Location_type),
		IsAgent:             data.Is_agent,
		HasRepossesed:       data.Has_repossesed,
	}

	// Property types
	if data.Property_types != nil {
		res.PropertyType = rentalPropertyTypesToRPC(*data.Property_types)
	}

	// Type of land
	if data.Type_of_land != nil {
		res.TypeOfLand = rentalStatusesToRPC(*data.Type_of_land)
	}

	// Metric
	if data.Metrict_type != nil {
		res.MetricType = rentalPriceTypeToRPC(*data.Metrict_type)
	}

	// Timing
	if data.Timing != nil {
		res.Timing = rentalTimingToRPC(*data.Timing)
	}

	// Commercial properties
	if data.Commercial_properties != nil {
		res.CommercialProperties = rentalCommercialPropertiesToRPC(*data.Commercial_properties)

	}
	// Commercial location
	if data.Commercial_locations != nil {
		res.CommericalPropertyLocation = rentalCommercialLocationsToRPC(data.Commercial_locations)
	}
	// Services
	if data.Services != nil {
		res.Services = rentalServicesToRPC(data.Services)
	}

	// Materials
	if data.Materials != nil {
		res.Materials = rentalMaterilasToRPC(data.Materials)
	}

	// Type of property
	if data.Type_of_property != nil {
		res.TypeOfProperty = rentalTypeOfPropertiesToRPC(*data.Type_of_property)
	}

	// Additional filters
	if data.Additional_filters != nil {
		res.AdditionalFilters = rentalAddiotanlFiltersToRPC(*data.Additional_filters)
	}

	// Layout
	if data.Layout != nil {
		res.Layout = rentalLayoutToRPC(*data.Layout)
	}
	// Climat Control
	if data.Climat_control != nil {
		res.ClimatControl = rentalClimatControlsToRPC(*data.Climat_control)
	}

	// Indoor Features
	if data.Indoor_features != nil {
		res.IndoorFeatures = rentalIndoorFeaturesToRPC(*data.Indoor_features)
	}

	// Outdoor Features
	if data.Outdoor_features != nil {
		res.OutdoorFeatures = rentalOutdoorFeaturesToRPC(*data.Outdoor_features)
	}

	return res
}

func rentalOfficeToRPC(data RealEstateOfficeInput) *rentalRPC.Office {
	return &rentalRPC.Office{
		Status:         rentalStatusToRPC(data.Status),
		AvailibatiFrom: NullToString(data.Availability_from),
		AvailibatiTo:   NullToString(data.Availability_to),
		TotalArea:      data.Total_area,
		Rental:         rentalInfoToRPC(data.Rental_info),
		Details:        rentalDetailsToRPC(data.Rental_detail),
		Layout:         rentalLayoutToRPC(data.Layout),
	}
}

func rentalLandToRPC(data RealEstateLandInput) *rentalRPC.Land {
	return &rentalRPC.Land{
		TypeOfLand:        rentalStatusesToRPC(data.Type_of_land),
		AvailibatiFrom:    NullToString(data.Availability_from),
		AvailibatiTo:      NullToString(data.Availability_to),
		TotalArea:         data.Total_area,
		Rental:            rentalInfoToRPC(data.Rental_info),
		Details:           rentalDetailsToRPC(data.Rental_detail),
		AdditionalFilters: rentalAddiotanlFiltersToRPC(data.More),
	}
}

func rentalCommerciaToRPC(data RealEstateCommercialInput) *rentalRPC.CommercialAndRuralFarm {
	return &rentalRPC.CommercialAndRuralFarm{
		AvailibatiFrom:       NullToString(data.Availability_from),
		AvailibatiTo:         NullToString(data.Availability_to),
		TotalArea:            data.Total_area,
		CommercialProperties: rentalCommercialPropertiesToRPC(data.Commercial_properties),
		AdditionalFilters:    rentalAddiotanlFiltersToRPC(data.Additional_filters),
		Rental:               rentalInfoToRPC(data.Rental_info),
		Details:              rentalDetailsToRPC(data.Rental_detail),
	}
}

func rentalBuildingsToRPC(data RealEstateBuildingsInput) *rentalRPC.BuilldingsAndGarage {
	return &rentalRPC.BuilldingsAndGarage{
		AvailibatiFrom: NullToString(data.Availability_from),
		AvailibatiTo:   NullToString(data.Availability_to),
		TotalArea:      data.Total_area,
		Rental:         rentalInfoToRPC(data.Rental_info),
		Details:        rentalDetailsToRPC(data.Rental_detail),
	}
}

func rentalGarageToRPC(data RealEstateGarageInput) *rentalRPC.BuilldingsAndGarage {
	return &rentalRPC.BuilldingsAndGarage{
		TotalArea:         data.Total_area,
		Rental:            rentalInfoToRPC(data.Rental_info),
		Details:           rentalDetailsToRPC(data.Rental_detail),
		AdditionalFilters: rentalAddiotanlFiltersToRPC(data.Additional_filters),
	}
}

func rentalHotelRoomsToRPC(data RealEstateHotelRoomsInput) *rentalRPC.HotelRooms {
	return &rentalRPC.HotelRooms{
		Status:         rentalStatusToRPC(data.Status),
		AvailibatiFrom: NullToString(data.Availability_from),
		AvailibatiTo:   NullToString(data.Availability_to),
		TotalArea:      data.Total_area,
		Rooms:          data.Rooms,
		Rental:         rentalInfoToRPC(data.Rental_info),
		Details:        rentalDetailsToRPC(data.Rental_detail),
	}
}

func rentalAppartamentToRPC(data HouseRentalAppartamentInput) *rentalRPC.Appartaments {
	return &rentalRPC.Appartaments{
		AvailibatiFrom:  NullToString(data.Availability_from),
		AvailibatiTo:    NullToString(data.Availability_to),
		BadRooms:        data.Badrooms,
		BathRooms:       data.Bathrooms,
		TotalArea:       data.Total_area,
		CarSpaces:       data.Car_spaces,
		Floor:           data.Floor,
		Floors:          data.Floors,
		HasRepossesed:   data.Has_repossesed,
		IsAgent:         data.Is_agent,
		Details:         rentalDetailsToRPC(data.Rental_detail),
		ClimatControl:   rentalClimatControlsToRPC(data.Climat_control),
		IndoorFeatures:  rentalIndoorFeaturesToRPC(data.Indoor_features),
		OutdoorFeatures: rentalOutdoorFeaturesToRPC(data.Outdoor_features),
		Rental:          rentalInfoToRPC(data.Rental_info),
		TypeOfProperty:  rentalTypeOfPropertiesToRPC(data.Type_of_home),
		Status:          rentalStatusToRPC(data.Status),
		WhoLive:         rentalWhoLivesToRPC(data.Who_live),
	}
}

func rentalCommercilaPropertiesRPCToStruct(data *rentalRPC.CommercialAndRuralFarm) *CommericalProperties {
	if data == nil {
		return &CommericalProperties{}
	}

	rental := data.GetRental()

	res := &CommericalProperties{
		ID:                    rental.GetID(),
		Alerts:                rental.GetAlerts(),
		Offers:                rental.GetOffers(),
		Is_urgent:             rental.GetIsUrgent(),
		Has_repossessed:       data.GetHasRepossesed(),
		Is_agent:              data.GetIsAgent(),
		Total_area:            data.GetTotalArea(),
		Availability_from:     data.GetAvailibatiFrom(),
		Availability_to:       data.GetAvailibatiTo(),
		Commercial_properties: rentalCommercialPropertiesRPCToStruct(data.GetCommercialProperties()),
		Commercial_locations:  rentalCommercialLocationsRPCToStruct(data.GetCommericalPropertyLocation()),
		Additional_filters:    rentalAditionalFiltersRPCToStruct(data.GetAdditionalFilters()),
		Status:                rentalStatusRPCToString(data.GetStatus()),
		Files:                 rentalFilesRPCToFiles(rental.GetFiles().GetFiles()),
		Metrict_type:          rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		Phones:                rentalPhoneRPCToPhone(data.GetPhones()),
		Price:                 rentalPriceRPCToPrice(data.GetPrice()),
		Rental_details:        rentalDetailsRPCToDetails(data.GetDetails()),
		Rental_info:           rentalInfoRPCToInfo(data.GetRental()),
	}

	if rental.GetIsCompany() {
		res.Company_id = rental.GetOwnerID()
	} else {
		res.User_id = rental.GetOwnerID()
	}

	return res
}

func rentalOfficeRPCToStruct(data *rentalRPC.Office) *Office {
	if data == nil {
		return &Office{}
	}

	rental := data.GetRental()

	res := &Office{
		ID:                rental.GetID(),
		Alerts:            rental.GetAlerts(),
		Offers:            rental.GetOffers(),
		Is_urgent:         rental.GetIsUrgent(),
		Has_repossessed:   data.GetHasRepossesed(),
		Is_agent:          data.GetIsAgent(),
		Total_area:        data.GetTotalArea(),
		Availability_from: data.GetAvailibatiFrom(),
		Availability_to:   data.GetAvailibatiTo(),
		Layout:            rentalLayoutRPCToString(data.GetLayout()),
		Building_use:      rentalBuilingUseRPCToString(data.GetBuildingUse()),
		Status:            rentalStatusRPCToString(data.GetStatus()),
		Files:             rentalFilesRPCToFiles(rental.GetFiles().GetFiles()),
		Metrict_type:      rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		Phones:            rentalPhoneRPCToPhone(data.GetPhones()),
		Price:             rentalPriceRPCToPrice(data.GetPrice()),
		Rental_details:    rentalDetailsRPCToDetails(data.GetDetails()),
		Rental_info:       rentalInfoRPCToInfo(data.GetRental()),
	}

	if rental.GetIsCompany() {
		res.Company_id = rental.GetOwnerID()
	} else {
		res.User_id = rental.GetOwnerID()
	}

	return res
}

func rentalStorageRoomRPCToStruct(data *rentalRPC.StorageRoom) *StorageRoom {
	if data == nil {
		return &StorageRoom{}
	}

	rental := data.GetRental()

	res := &StorageRoom{
		ID:                rental.GetID(),
		Alerts:            rental.GetAlerts(),
		Offers:            rental.GetOffers(),
		Is_urgent:         rental.GetIsUrgent(),
		Has_repossessed:   data.GetHasRepossesed(),
		Is_agent:          data.GetIsAgent(),
		Total_area:        data.GetTotalArea(),
		Availability_from: data.GetAvailibatiFrom(),
		Availability_to:   data.GetAvailibatiTo(),
		Status:            rentalStatusRPCToString(data.GetStatus()),
		Files:             rentalFilesRPCToFiles(rental.GetFiles().GetFiles()),
		Metrict_type:      rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		Phones:            rentalPhoneRPCToPhone(data.GetPhones()),
		Price:             rentalPriceRPCToPrice(data.GetPrice()),
		Rental_details:    rentalDetailsRPCToDetails(data.GetDetails()),
		Rental_info:       rentalInfoRPCToInfo(data.GetRental()),
	}

	if rental.GetIsCompany() {
		res.Company_id = rental.GetOwnerID()
	} else {
		res.User_id = rental.GetOwnerID()
	}

	return res
}

func rentalAppartamentRPCToStruct(data *rentalRPC.Appartaments) *Appartment {
	if data == nil {
		return &Appartment{}
	}

	rental := data.GetRental()

	res := &Appartment{
		ID:                rental.GetID(),
		Alerts:            rental.GetAlerts(),
		Offers:            rental.GetOffers(),
		Is_urgent:         rental.GetIsUrgent(),
		Badrooms:          data.GetBadRooms(),
		Bathrooms:         data.GetBathRooms(),
		Car_spaces:        data.GetCarSpaces(),
		Floor:             data.GetFloor(),
		Floors:            data.GetFloors(),
		Has_repossessed:   data.GetHasRepossesed(),
		Is_agent:          data.GetIsAgent(),
		Total_area:        data.GetTotalArea(),
		Availability_from: data.GetAvailibatiFrom(),
		Availability_to:   data.GetAvailibatiTo(),
		Status:            rentalStatusRPCToString(data.GetStatus()),
		Climat_control:    rentalClimatControlsRPCToArray(data.GetClimatControl()),
		Files:             rentalFilesRPCToFiles(rental.GetFiles().GetFiles()),
		Indoor_features:   rentalIndoorFeaturesRPCToArray(data.GetIndoorFeatures()),
		Outdoor_features:  rentalOutdoorFeaturesRPCToArray(data.GetOutdoorFeatures()),
		Metrict_type:      rentalPriceTypeRPCToPriceType(data.GetMetricType()),
		Phones:            rentalPhoneRPCToPhone(data.GetPhones()),
		Price:             rentalPriceRPCToPrice(data.GetPrice()),
		Rental_details:    rentalDetailsRPCToDetails(data.GetDetails()),
		Rental_info:       rentalInfoRPCToInfo(data.GetRental()),
	}

	if rental.GetIsCompany() {
		res.Company_id = rental.GetOwnerID()
	} else {
		res.User_id = rental.GetOwnerID()
	}

	return res

}

func rentalInfoRPCToInfo(data *rentalRPC.Rental) RentalInfo {
	if data == nil {
		return RentalInfo{}
	}

	return RentalInfo{
		Property_type: rentalPropertyTypeRPCToString(data.GetPropertyType()),
		Deal_type:     rentalDealTypeRPCToString(data.GetDealType()),
		Location:      rentalLocationRPCToLocation(data.GetLocation()),
		Expired_days:  data.GetExpiredDays(),
		Post_currency: data.GetPostCurrency(),
		Created_at:    data.GetCreatedAt(),
	}
}

func rentalLocationRPCToLocation(data *rentalRPC.Location) RentalLocation {
	if data == nil {
		return RentalLocation{}
	}

	return RentalLocation{
		City:    data.GetCity().GetCity(),
		Country: data.GetCountry().GetId(),
		Address: data.GetAddress(),
		Street:  data.GetStreet(),
	}
}

func rentalDetailsRPCToDetails(data []*rentalRPC.Detail) []RentalDetail {
	if len(data) <= 0 {
		return nil
	}

	details := make([]RentalDetail, 0, len(data))

	for _, d := range data {
		details = append(details, RentalDetail{
			ID:          d.GetID(),
			Title:       d.GetTitle(),
			House_rules: d.GetHouseRules(),
			Description: d.GetDescription(),
		})
	}

	return details
}

func rentalPriceRPCToPrice(data *rentalRPC.Price) RentalPrice {
	if data == nil {
		return RentalPrice{}
	}

	return RentalPrice{
		Price_type: rentalPriceTypeRPCToPriceType(data.GetPriceType()),
		Currency:   data.GetCurrency(),
		Max_price:  data.GetMaxPrice(),
		Min_price:  data.GetMinPrice(),
	}
}

func rentalPhoneRPCToPhone(data []*rentalRPC.Phone) []RentalPhone {
	if len(data) <= 0 {
		return nil
	}

	phones := make([]RentalPhone, 0, len(data))

	for _, p := range data {
		phones = append(phones, RentalPhone{
			Country_code_id: p.GetCountryCode(),
			Number:          p.GetNumber(),
		})
	}

	return phones
}
func rentalFilesRPCToFiles(data []*rentalRPC.File) []File {
	if len(data) <= 0 {
		return nil
	}

	files := make([]File, 0, len(data))

	for _, f := range data {
		files = append(files, File{
			ID:        f.GetID(),
			Address:   f.GetURL(),
			Mime_type: f.GetMimeType(),
			Name:      f.GetName(),
		})
	}

	return files
}

func rentalOutdoorFeaturesRPCToArray(data []rentalRPC.OutdoorFeaturesEnum) []string {
	if len(data) <= 0 {
		return nil
	}

	features := make([]string, 0, len(data))

	for _, c := range data {
		features = append(features, rentalOutdoorFeatureRPCToString(c))
	}

	return features

}

func rentalIndoorFeaturesRPCToArray(data []rentalRPC.IndoorFeaturesEnum) []string {
	if len(data) <= 0 {
		return nil
	}

	features := make([]string, 0, len(data))

	for _, c := range data {
		features = append(features, rentalIndoorFeatureRPCToString(c))
	}

	return features
}

func rentalCommercialPropertiesRPCToStruct(data []rentalRPC.CommercialPropertyEnum) []string {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]string, 0, len(data))

	for _, c := range data {
		properties = append(properties, rentalCommercialPropertyRPCToString(c))
	}

	return properties
}

func rentalAditionalFiltersRPCToStruct(data []rentalRPC.AdditionalFiltersEnum) []string {
	if len(data) <= 0 {
		return nil
	}

	filters := make([]string, 0, len(data))

	for _, f := range data {
		filters = append(filters, rentalAddiotanlFilterRPCToString(f))
	}

	return filters
}

func rentalCommercialLocationsRPCToStruct(data []rentalRPC.CommericalPropertyLocationEnum) []string {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]string, 0, len(data))

	for _, c := range data {
		properties = append(properties, rentalCommercialLocationRPCToString(c))
	}

	return properties
}

func rentalClimatControlsRPCToArray(data []rentalRPC.ClimatControlEnum) []string {
	if len(data) <= 0 {
		return nil
	}

	constrols := make([]string, 0, len(data))

	for _, c := range data {
		constrols = append(constrols, rentalClimatControlRPCToString(c))
	}

	return constrols
}

func rentalPropertyTypesToRPC(data []string) []rentalRPC.PropertyTypeEnum {
	if data == nil {
		return nil
	}

	properties := make([]rentalRPC.PropertyTypeEnum, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalPropertyTypeToRPC(p))
	}

	return properties
}

func rentalCommercialLocationsToRPC(data *[]string) []rentalRPC.CommericalPropertyLocationEnum {
	if data == nil {
		return nil
	}

	locations := make([]rentalRPC.CommericalPropertyLocationEnum, 0, len(*data))

	for _, s := range *data {
		locations = append(locations, rentalCommercialLocationToRPC(s))
	}

	return locations
}

func rentalServicesToRPC(data *[]string) []rentalRPC.ServiceEnum {
	if data == nil {
		return nil
	}

	services := make([]rentalRPC.ServiceEnum, 0, len(*data))

	for _, s := range *data {
		services = append(services, rentalServiceToRPC(s))
	}

	return services

}

func rentalLocationTypesToRPC(data *[]string) []rentalRPC.LocationEnum {
	if data == nil {
		return nil
	}

	locations := make([]rentalRPC.LocationEnum, 0, len(*data))

	for _, l := range *data {
		locations = append(locations, rentalLocationTypeToRPC(l))
	}

	return locations
}

func rentalMaterilasToRPC(data *[]string) []rentalRPC.MaterialEnum {
	if data == nil {
		return nil
	}

	materials := make([]rentalRPC.MaterialEnum, 0, len(*data))

	for _, m := range *data {
		materials = append(materials, rentalMaterialToRPC(m))
	}

	return materials

}

func rentalStatusesToRPC(data []string) []rentalRPC.StatusEnum {
	if len(data) <= 0 {
		return nil
	}

	statuses := make([]rentalRPC.StatusEnum, 0, len(data))

	for _, s := range data {
		statuses = append(statuses, rentalStatusToRPC(s))
	}

	return statuses
}

func rentalAddiotanlFiltersToRPC(data []string) []rentalRPC.AdditionalFiltersEnum {
	if len(data) <= 0 {
		return nil
	}

	filters := make([]rentalRPC.AdditionalFiltersEnum, 0, len(data))

	for _, p := range data {
		filters = append(filters, rentalAddiotanlFilterToRPC(p))
	}

	return filters
}

func rentalCommercialPropertiesToRPC(data []string) []rentalRPC.CommercialPropertyEnum {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]rentalRPC.CommercialPropertyEnum, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalCommercialPropertyToRPC(p))
	}

	return properties
}

func rentalWhoLivesToRPC(data *[]string) []rentalRPC.WhoLiveEnum {
	if data == nil {
		return nil
	}

	lives := make([]rentalRPC.WhoLiveEnum, 0, len(*data))

	for _, l := range *data {
		lives = append(lives, rentalWhoLiveToRPC(l))
	}

	return lives
}

func rentalTypeOfPropertiesToRPC(data []string) []rentalRPC.TypeOfPropertyEnum {
	if len(data) <= 0 {
		return nil
	}

	properties := make([]rentalRPC.TypeOfPropertyEnum, 0, len(data))

	for _, p := range data {
		properties = append(properties, rentalTypeOfPropertyToRPC(p))
	}

	return properties
}

func rentalInfoToRPC(data RentalInput) *rentalRPC.Rental {
	return &rentalRPC.Rental{
		DealType:     rentalDealTypeToRPC(data.Deal_type),
		PropertyType: rentalPropertyTypeToRPC(data.Property_type),
		Location:     rentalLocationToRPC(data.Location),
		ExpiredDays:  data.Expired_days,
		PostCurrency: data.Post_currency,
	}
}

func rentalLocationToRPC(data RentalLocationInput) *rentalRPC.Location {
	return &rentalRPC.Location{
		Country: &rentalRPC.Country{
			Id: data.Country_id,
		},
		City:    rentalCityToRPC(data.City),
		Street:  NullToString(data.Street),
		Address: NullToString(data.Address),
	}
}

func rentalCityToRPC(data CityInput) *rentalRPC.City {
	return &rentalRPC.City{
		Id:          NullToString(data.ID),
		City:        NullToString(data.City),
		Subdivision: NullToString(data.Subdivision),
	}
}

func rentalPriceToRPC(data *RentalPriceInput) *rentalRPC.Price {
	if data == nil {
		return nil
	}

	return &rentalRPC.Price{
		PriceType: rentalPriceTypeToRPC(data.Price_type),
		MinPrice:  NullToInt32(data.Min_price),
		MaxPrice:  NullToInt32(data.Max_price),
		FixPrice:  NullToInt32(data.Fix_price),
		Currency:  data.Currency,
	}
}

func rentalPhonesToRPC(data *[]RentalPhoneInput) []*rentalRPC.Phone {
	if data == nil {
		return nil
	}

	phones := make([]*rentalRPC.Phone, 0, len(*data))

	for _, p := range *data {
		phones = append(phones, rentalPhoneToRPC(p))
	}

	return phones
}

func rentalPhoneToRPC(data RentalPhoneInput) *rentalRPC.Phone {

	return &rentalRPC.Phone{
		Number:      data.Number,
		CountryCode: data.Country_code_id,
	}
}
func rentalClimatControlsToRPC(data []string) []rentalRPC.ClimatControlEnum {
	if len(data) <= 0 {
		return nil
	}

	controls := make([]rentalRPC.ClimatControlEnum, 0, len(data))

	for _, c := range data {
		controls = append(controls, rentalClimatControlToRPC(c))
	}

	return controls
}

func rentalDetailsToRPC(data []RentalDetailInput) []*rentalRPC.Detail {
	if len(data) <= 0 {
		return nil
	}

	details := make([]*rentalRPC.Detail, 0, len(data))

	for _, d := range data {
		details = append(details, rentalDetailToRPC(d))
	}

	return details
}

func rentalDetailToRPC(data RentalDetailInput) *rentalRPC.Detail {
	return &rentalRPC.Detail{
		Title:       data.Title,
		HouseRules:  NullToString(data.House_rules),
		Description: NullToString(data.Description),
	}
}

func rentalIndoorFeaturesToRPC(data []string) []rentalRPC.IndoorFeaturesEnum {
	if len(data) <= 0 {
		return nil
	}

	features := make([]rentalRPC.IndoorFeaturesEnum, 0, len(data))

	for _, c := range data {
		features = append(features, rentalIndoorFeatureToRPC(c))
	}

	return features
}

func rentalOutdoorFeaturesToRPC(data []string) []rentalRPC.OutdoorFeaturesEnum {
	if len(data) <= 0 {
		return nil
	}

	features := make([]rentalRPC.OutdoorFeaturesEnum, 0, len(data))

	for _, c := range data {
		features = append(features, rentalOutdoorFeatureToRPC(c))
	}

	return features
}

func rentalWhoLiveToRPC(data string) rentalRPC.WhoLiveEnum {

	switch data {
	case "WhoLive_Mortgagor":
		return rentalRPC.WhoLiveEnum_WhoLive_Mortgagor
	case "WhoLive_Owner":
		return rentalRPC.WhoLiveEnum_WhoLiveOwner
	}

	return rentalRPC.WhoLiveEnum_WhoLive_Any
}

func rentalWhoLiveRPCToString(data rentalRPC.WhoLiveEnum) string {

	switch data {
	case rentalRPC.WhoLiveEnum_WhoLive_Mortgagor:
		return "WhoLive_Mortgagor"
	case rentalRPC.WhoLiveEnum_WhoLiveOwner:
		return "WhoLive_Owner"
	}

	return "WhoLive_Any"
}

func rentalPriceTypeToRPC(data string) rentalRPC.PriceTypeEnum {

	switch data {
	case "PriceType_Total":
		return rentalRPC.PriceTypeEnum_PriceType_Total
	case "PriceType_MetreSquare":
		return rentalRPC.PriceTypeEnum_PriceType_MetreSquare
	case "PriceType_FeetSquare":
		return rentalRPC.PriceTypeEnum_PriceType_FeetSquare

	}

	return rentalRPC.PriceTypeEnum_PriceType_Any
}

func rentalPriceTypeRPCToPriceType(data rentalRPC.PriceTypeEnum) string {

	switch data {
	case rentalRPC.PriceTypeEnum_PriceType_Total:
		return "PriceType_Total"
	case rentalRPC.PriceTypeEnum_PriceType_MetreSquare:
		return "PriceType_MetreSquare"
	case rentalRPC.PriceTypeEnum_PriceType_FeetSquare:
		return "PriceType_FeetSquare"
	}

	return "PriceType_Any"
}

func rentalClimatControlToRPC(data string) rentalRPC.ClimatControlEnum {
	switch data {
	case "ClimatControl_AirConditioning":
		return rentalRPC.ClimatControlEnum_ClimatControl_AirConditioning
	case "ClimatControl_Hearting":
		return rentalRPC.ClimatControlEnum_ClimatControl_Hearting
	case "ClimatControl_WaterTank":
		return rentalRPC.ClimatControlEnum_ClimatControl_WaterTank
	case "ClimatControl_SolarPanels":
		return rentalRPC.ClimatControlEnum_ClimatControl_SolarPanels
	case "ClimatControl_HighEnergyEfficiency":
		return rentalRPC.ClimatControlEnum_ClimatControl_HighEnergyEfficiency
	case "ClimatControlSolarHotWater":
		return rentalRPC.ClimatControlEnum_ClimatControlSolarHotWater
	case "ClimatControl_ZonalHeating":
		return rentalRPC.ClimatControlEnum_ClimatControl_ZonalHeating
	case "ClimatControl_HeatPumps":
		return rentalRPC.ClimatControlEnum_ClimatControl_HeatPumps
	}
	return rentalRPC.ClimatControlEnum_ClimatControl_Any
}

func rentalClimatControlRPCToString(data rentalRPC.ClimatControlEnum) string {
	switch data {
	case rentalRPC.ClimatControlEnum_ClimatControl_AirConditioning:
		return "ClimatControl_AirConditioning"
	case rentalRPC.ClimatControlEnum_ClimatControl_Hearting:
		return "ClimatControl_Hearting"
	case rentalRPC.ClimatControlEnum_ClimatControl_WaterTank:
		return "ClimatControl_WaterTank"
	case rentalRPC.ClimatControlEnum_ClimatControl_SolarPanels:
		return "ClimatControl_SolarPanels"
	case rentalRPC.ClimatControlEnum_ClimatControl_HighEnergyEfficiency:
		return "ClimatControl_HighEnergyEfficiency"
	case rentalRPC.ClimatControlEnum_ClimatControlSolarHotWater:
		return "ClimatControlSolarHotWater"
	case rentalRPC.ClimatControlEnum_ClimatControl_ZonalHeating:
		return "ClimatControl_ZonalHeating"
	case rentalRPC.ClimatControlEnum_ClimatControl_HeatPumps:
		return "ClimatControl_HeatPumps"
	}
	return "ClimatControl_Any"
}

func rentalIndoorFeatureToRPC(data string) rentalRPC.IndoorFeaturesEnum {

	switch data {
	case "IndoorFeatures_Ensuit":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Ensuit
	case "IndoorFeatures_Study":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Study
	case "IndoorFeatures_AlarmSystem":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_AlarmSystem
	case "IndoorFeatures_Floorboards":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Floorboards
	case "IndoorFeatures_RumpusRoom":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_RumpusRoom
	case "IndoorFeatures_StorageRoom":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_StorageRoom
	case "IndoorFeatures_Dishwasher":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Dishwasher
	case "IndoorFeatures_Lift":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Lift
	case "IndoorFeatures_BuiltInRobes":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_BuiltInRobes
	case "IndoorFeatures_Broadband":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Broadband
	case "IndoorFeatures_Gym":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Gym
	case "IndoorFeatures_Workshop":
		return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Workshop
	}

	return rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Any
}

func rentalIndoorFeatureRPCToString(data rentalRPC.IndoorFeaturesEnum) string {

	switch data {
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Ensuit:
		return "IndoorFeatures_Ensuit"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Study:
		return "IndoorFeatures_Study"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_AlarmSystem:
		return "IndoorFeatures_AlarmSystem"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Floorboards:
		return "IndoorFeatures_Floorboards"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_RumpusRoom:
		return "IndoorFeatures_RumpusRoom"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_StorageRoom:
		return "IndoorFeatures_StorageRoom"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Dishwasher:
		return "IndoorFeatures_Dishwasher"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Lift:
		return "IndoorFeatures_Lift"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_BuiltInRobes:
		return "IndoorFeatures_BuiltInRobes"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Broadband:
		return "IndoorFeatures_Broadband"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Gym:
		return "IndoorFeatures_Gym"
	case rentalRPC.IndoorFeaturesEnum_IndoorFeatures_Workshop:
		return "IndoorFeatures_Workshop"
	}

	return "IndoorFeatures_Any"
}

func rentalOutdoorFeatureToRPC(data string) rentalRPC.OutdoorFeaturesEnum {

	switch data {
	case "OutdoorFeatures_SwimmingPool":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_SwimmingPool
	case "OutdoorFeatures_Balcony":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Balcony
	case "OutdoorFeatures_UndercoverParking":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_UndercoverParking
	case "OutdoorFeatures_FullyFenced":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_FullyFenced
	case "OutdoorFeatures_TennisCourt":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_TennisCourt
	case "OutdoorFeatures_Garden":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garden
	case "OutdoorFeatures_Garage":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garage
	case "OutdoorFeatures_OutdoorArea":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorArea
	case "OutdoorFeatures_Shed":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Shed
	case "OutdoorFeatures_OutdoorSpa":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorSpa
	case "OutdoorFeatures_Outbuildings":
		return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Outbuildings
	}
	return rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Any
}

func rentalOutdoorFeatureRPCToString(data rentalRPC.OutdoorFeaturesEnum) string {

	switch data {
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_SwimmingPool:
		return "OutdoorFeatures_SwimmingPool"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Balcony:
		return "OutdoorFeatures_Balcony"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_UndercoverParking:
		return "OutdoorFeatures_UndercoverParking"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_FullyFenced:
		return "OutdoorFeatures_FullyFenced"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_TennisCourt:
		return "OutdoorFeatures_TennisCourt"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garden:
		return "OutdoorFeatures_Garden"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Garage:
		return "OutdoorFeatures_Garage"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorArea:
		return "OutdoorFeatures_OutdoorArea"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Shed:
		return "OutdoorFeatures_Shed"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_OutdoorSpa:
		return "OutdoorFeatures_OutdoorSpa"
	case rentalRPC.OutdoorFeaturesEnum_OutdoorFeatures_Outbuildings:
		return "OutdoorFeatures_Outbuildings"

	}
	return "rentalRPC.OutdoorFeaturesEnum"
}

func rentalDealTypeToRPC(data string) rentalRPC.DealTypeEnum {

	switch data {
	case "DealType_Rent":
		return rentalRPC.DealTypeEnum_DealType_Rent
	case "DealType_Lease":
		return rentalRPC.DealTypeEnum_DealType_Lease
	case "DealType_Share":
		return rentalRPC.DealTypeEnum_DealType_Share
	case "DealType_Sell":
		return rentalRPC.DealTypeEnum_DealType_Sell
	case "DealType_Build":
		return rentalRPC.DealTypeEnum_DealType_Build
	case "DealType_Materials":
		return rentalRPC.DealTypeEnum_DealType_Materials
	case "DealType_Renovation":
		return rentalRPC.DealTypeEnum_DealType_Renovation
	case "DealType_Move":
		return rentalRPC.DealTypeEnum_DealType_Move
	}
	return rentalRPC.DealTypeEnum_DealType_Any
}

func rentalDealTypeRPCToString(data rentalRPC.DealTypeEnum) string {

	switch data {
	case rentalRPC.DealTypeEnum_DealType_Rent:
		return "DealType_Rent"
	case rentalRPC.DealTypeEnum_DealType_Lease:
		return "DealType_Lease"
	case rentalRPC.DealTypeEnum_DealType_Share:
		return "DealType_Share"
	case rentalRPC.DealTypeEnum_DealType_Sell:
		return "DealType_Sell"
	case rentalRPC.DealTypeEnum_DealType_Build:
		return "DealType_Build"
	case rentalRPC.DealTypeEnum_DealType_Materials:
		return "DealType_Materials"
	case rentalRPC.DealTypeEnum_DealType_Renovation:
		return "DealType_Renovation"
	case rentalRPC.DealTypeEnum_DealType_Move:
		return "DealType_Move"
	}
	return "rentalRPC.DealTypeEnum"
}

func rentalPropertyTypeToRPC(data string) rentalRPC.PropertyTypeEnum {

	switch data {
	case "PropertyType_All":
		return rentalRPC.PropertyTypeEnum_PropertyType_All
	case "PropertyType_NewHomes":
		return rentalRPC.PropertyTypeEnum_PropertyType_NewHomes
	case "PropertyType_Homes":
		return rentalRPC.PropertyTypeEnum_PropertyType_Homes
	case "PropertyType_Houses":
		return rentalRPC.PropertyTypeEnum_PropertyType_Houses
	case "PropertyType_Appartments":
		return rentalRPC.PropertyTypeEnum_PropertyType_Appartments
	case "PropertyType_Garages":
		return rentalRPC.PropertyTypeEnum_PropertyType_Garages
	case "PropertyType_StorageRooms":
		return rentalRPC.PropertyTypeEnum_PropertyType_StorageRooms
	case "PropertyType_Offices":
		return rentalRPC.PropertyTypeEnum_PropertyType_Offices
	case "PropertyType_CommercialProperties":
		return rentalRPC.PropertyTypeEnum_PropertyType_CommercialProperties
	case "PropertyType_Buildings":
		return rentalRPC.PropertyTypeEnum_PropertyType_Buildings
	case "PropertyType_Land":
		return rentalRPC.PropertyTypeEnum_PropertyType_Land
	case "PropertyType_BareLand":
		return rentalRPC.PropertyTypeEnum_PropertyType_BareLand
	case "PropertyType_Barn":
		return rentalRPC.PropertyTypeEnum_PropertyType_Barn
	case "PropertyType_SummerCottage":
		return rentalRPC.PropertyTypeEnum_PropertyType_SummerCottage
	case "PropertyType_RuralFarm":
		return rentalRPC.PropertyTypeEnum_PropertyType_RuralFarm
	case "PropertyType_HotelRoom":
		return rentalRPC.PropertyTypeEnum_PropertyType_HotelRoom
	}
	return rentalRPC.PropertyTypeEnum_PropertyType_Any
}

func rentalPropertyTypeRPCToString(data rentalRPC.PropertyTypeEnum) string {

	switch data {
	case rentalRPC.PropertyTypeEnum_PropertyType_All:
		return "PropertyType_All"
	case rentalRPC.PropertyTypeEnum_PropertyType_NewHomes:
		return "PropertyType_NewHomes"
	case rentalRPC.PropertyTypeEnum_PropertyType_Homes:
		return "PropertyType_Homes"
	case rentalRPC.PropertyTypeEnum_PropertyType_Houses:
		return "PropertyType_Houses"
	case rentalRPC.PropertyTypeEnum_PropertyType_Appartments:
		return "PropertyType_Appartments"
	case rentalRPC.PropertyTypeEnum_PropertyType_Garages:
		return "PropertyType_Garages"
	case rentalRPC.PropertyTypeEnum_PropertyType_StorageRooms:
		return "PropertyType_StorageRooms"
	case rentalRPC.PropertyTypeEnum_PropertyType_Offices:
		return "PropertyType_Offices"
	case rentalRPC.PropertyTypeEnum_PropertyType_CommercialProperties:
		return "PropertyType_CommercialProperties"
	case rentalRPC.PropertyTypeEnum_PropertyType_Buildings:
		return "PropertyType_Buildings"
	case rentalRPC.PropertyTypeEnum_PropertyType_Land:
		return "PropertyType_Land"
	case rentalRPC.PropertyTypeEnum_PropertyType_BareLand:
		return "PropertyType_BareLand"
	case rentalRPC.PropertyTypeEnum_PropertyType_Barn:
		return "PropertyType_Barn"
	case rentalRPC.PropertyTypeEnum_PropertyType_SummerCottage:
		return "PropertyType_SummerCottage"
	case rentalRPC.PropertyTypeEnum_PropertyType_RuralFarm:
		return "PropertyType_RuralFarm"
	}
	return "PropertyType_Any"
}

func rentalTypeOfPropertyToRPC(data string) rentalRPC.TypeOfPropertyEnum {

	switch data {
	case "TypeOfProperty_Appartaments":
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Appartaments
	case "TypeOfProperty_Houses":
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Houses
	case "TypeOfProperty_CountryHomes":
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_CountryHomes
	case "TypeOfProperty_Duplex":
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Duplex
	case "TypeOfProperty_Penthouses":
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Penthouses
	case "TypeOfProperty_Bungalow":
		return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Bungalow
	}
	return rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Any
}

func rentalTypeOfPropertyRPCToString(data rentalRPC.TypeOfPropertyEnum) string {

	switch data {
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Appartaments:
		return "TypeOfProperty_Appartaments"
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Houses:
		return "TypeOfProperty_Houses"
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_CountryHomes:
		return "TypeOfProperty_CountryHomes"
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Duplex:
		return "TypeOfProperty_Duplex"
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Penthouses:
		return "TypeOfProperty_Penthouses"
	case rentalRPC.TypeOfPropertyEnum_TypeOfProperty_Bungalow:
		return "TypeOfProperty_Bungalow"
	}
	return "TypeOfProperty_Any"
}

func rentalStatusToRPC(data string) rentalRPC.StatusEnum {

	switch data {
	case "Status_OldBuild":
		return rentalRPC.StatusEnum_Status_OldBuild
	case "Status_NewBuilding":
		return rentalRPC.StatusEnum_Status_NewBuilding
	case "Status_UnderConstruction":
		return rentalRPC.StatusEnum_Status_UnderConstruction
	case "StatusDeveloped":
		return rentalRPC.StatusEnum_StatusDeveloped
	case "Status_Buildable":
		return rentalRPC.StatusEnum_Status_Buildable
	case "Status_NonBuilding":
		return rentalRPC.StatusEnum_Status_NonBuilding
	}
	return rentalRPC.StatusEnum_Status_Any
}

func rentalStatusRPCToString(data rentalRPC.StatusEnum) string {

	switch data {
	case rentalRPC.StatusEnum_Status_OldBuild:
		return "Status_OldBuild"
	case rentalRPC.StatusEnum_Status_NewBuilding:
		return "Status_NewBuilding"
	case rentalRPC.StatusEnum_Status_UnderConstruction:
		return "Status_UnderConstruction"
	case rentalRPC.StatusEnum_StatusDeveloped:
		return "StatusDeveloped"
	case rentalRPC.StatusEnum_Status_Buildable:
		return "Status_Buildable"
	case rentalRPC.StatusEnum_Status_NonBuilding:
		return "Status_NonBuilding"
	}
	return "Status_Any"
}

func rentalCommercialPropertyToRPC(data string) rentalRPC.CommercialPropertyEnum {

	switch data {
	case "CommericalProperty_OfficeSpace":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_OfficeSpace
	case "CommericalProperty_CommercialPremises":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_CommercialPremises
	case "CommericalProperty_IndustrialBuilding":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_IndustrialBuilding
	case "CommericalProperty_Warehouse":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_Warehouse
	case "CommericalProperty_FoodFacility":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_FoodFacility
	case "CommericalProperty_Garage":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_Garage
	case "CommericalProperty_Basement":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_Basement
	case "CommericalProperty_TradingPlace":
		return rentalRPC.CommercialPropertyEnum_CommericalProperty_TradingPlace
	}

	return rentalRPC.CommercialPropertyEnum_CommericalProperty_Any
}

func rentalCommercialPropertyRPCToString(data rentalRPC.CommercialPropertyEnum) string {

	switch data {
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_OfficeSpace:
		return "CommericalProperty_OfficeSpace"
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_CommercialPremises:
		return "CommericalProperty_CommercialPremises"
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_IndustrialBuilding:
		return "CommericalProperty_IndustrialBuilding"
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_Warehouse:
		return "CommericalProperty_Warehouse"
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_FoodFacility:
		return "CommericalProperty_FoodFacility"
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_Garage:
		return "CommericalProperty_Garage"
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_Basement:
		return "CommericalProperty_Basement"
	case rentalRPC.CommercialPropertyEnum_CommericalProperty_TradingPlace:
		return "CommericalProperty_TradingPlace"
	}

	return "CommericalProperty_Any"
}

func rentalCommercialLocationToRPC(data string) rentalRPC.CommericalPropertyLocationEnum {

	switch data {
	case "CommericalPropertyLocation_Indifferent":
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Indifferent
	case "CommericalPropertyLocation_InShoppingCentre":
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_InShoppingCentre
	case "CommericalPropertyLocation_Mezzanine":
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Mezzanine
	case "CommericalPropertyLocation_BelowGround":
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_BelowGround
	case "CommericalPropertyLocation_Other":
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Other
	case "CommericalPropertyLocation_Garage":
		return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Garage
	}

	return rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Any
}

func rentalCommercialLocationRPCToString(data rentalRPC.CommericalPropertyLocationEnum) string {

	switch data {
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Indifferent:
		return "CommericalPropertyLocation_Indifferent"
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_InShoppingCentre:
		return "CommericalPropertyLocation_InShoppingCentre"
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Mezzanine:
		return "CommericalPropertyLocation_Mezzanine"
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_BelowGround:
		return "CommericalPropertyLocation_BelowGround"
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Other:
		return "CommericalPropertyLocation_Other"
	case rentalRPC.CommericalPropertyLocationEnum_CommericalPropertyLocation_Garage:
		return "CommericalPropertyLocation_Garage"
	}

	return "CommericalPropertyLocation_Any"
}

func rentalAddiotanlFilterToRPC(data string) rentalRPC.AdditionalFiltersEnum {

	switch data {
	case "Additional_Filter_Electricity":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Electricity
	case "Additional_Filter_Water":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Water
	case "Additional_Filter_NaturalGas":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_NaturalGas
	case "Additional_Filter_Sewage":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Sewage
	case "Additional_Filter_AireConditioning":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_AireConditioning
	case "Additional_Filter_Heating":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Heating
	case "Additional_Filter_OnCorner":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_OnCorner
	case "Additional_Filter_SmokeExtractor":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_SmokeExtractor
	case "Additional_Filter_MotoBikeGarage":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_MotoBikeGarage
	case "Additional_Filter_AutomaticDoor":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_AutomaticDoor
	case "Additional_Filter_SecuritySystem":
		return rentalRPC.AdditionalFiltersEnum_Additional_Filter_SecuritySystem
	}

	return rentalRPC.AdditionalFiltersEnum_Additional_Filter_Any
}

func rentalAddiotanlFilterRPCToString(data rentalRPC.AdditionalFiltersEnum) string {

	switch data {
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Electricity:
		return "Additional_Filter_Electricity"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Water:
		return "Additional_Filter_Water"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_NaturalGas:
		return "Additional_Filter_NaturalGas"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Sewage:
		return "Additional_Filter_Sewage"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_AireConditioning:
		return "Additional_Filter_AireConditioning"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_Heating:
		return "Additional_Filter_Heating"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_OnCorner:
		return "Additional_Filter_OnCorner"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_SmokeExtractor:
		return "Additional_Filter_SmokeExtractor"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_MotoBikeGarage:
		return "Additional_Filter_MotoBikeGarage"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_AutomaticDoor:
		return "Additional_Filter_AutomaticDoor"
	case rentalRPC.AdditionalFiltersEnum_Additional_Filter_SecuritySystem:
		return "Additional_Filter_SecuritySystem"
	}

	return "Additional_Filter_Any"
}

func rentalLayoutToRPC(data string) rentalRPC.LayoutEnum {
	switch data {
	case "Layout_Indifferent":
		return rentalRPC.LayoutEnum_Layout_Indifferent
	case "Layout_OpenPlan":
		return rentalRPC.LayoutEnum_Layout_OpenPlan
	case "Layout_Walls":
		return rentalRPC.LayoutEnum_Layout_Walls
	}
	return rentalRPC.LayoutEnum_Layout_Any
}

func rentalLayoutRPCToString(data rentalRPC.LayoutEnum) string {
	switch data {
	case rentalRPC.LayoutEnum_Layout_Indifferent:
		return "Layout_Indifferent"
	case rentalRPC.LayoutEnum_Layout_OpenPlan:
		return "Layout_OpenPlan"
	case rentalRPC.LayoutEnum_Layout_Walls:
		return "Layout_Walls"
	}
	return "Layout_Any"
}

func rentalBuilingUseToRPC(data *string) rentalRPC.BuildingUseEnum {
	if data == nil {
		return rentalRPC.BuildingUseEnum_Building_Use_Any
	}

	switch *data {
	case "Building_Use_Indifferent":
		return rentalRPC.BuildingUseEnum_Building_Use_Indifferent
	case "Building_Use_OnlyOffice":
		return rentalRPC.BuildingUseEnum_Building_Use_OnlyOffice
	case "Building_Use_Mixed":
		return rentalRPC.BuildingUseEnum_Building_Use_Mixed
	}
	return rentalRPC.BuildingUseEnum_Building_Use_Any
}

func rentalBuilingUseRPCToString(data rentalRPC.BuildingUseEnum) string {
	switch data {
	case rentalRPC.BuildingUseEnum_Building_Use_Indifferent:
		return "Building_Use_Indifferent"
	case rentalRPC.BuildingUseEnum_Building_Use_OnlyOffice:
		return "Building_Use_OnlyOffice"
	case rentalRPC.BuildingUseEnum_Building_Use_Mixed:
		return "Building_Use_Mixed"
	}
	return "Building_Use_Any"
}

func rentalMaterialToRPC(data string) rentalRPC.MaterialEnum {

	switch data {
	case "Material_Lumber_Composites":
		return rentalRPC.MaterialEnum_Material_Lumber_Composites
	case "Material_Fencing":
		return rentalRPC.MaterialEnum_Material_Fencing
	case "Material_Decking":
		return rentalRPC.MaterialEnum_Material_Decking
	case "Material_Fastners":
		return rentalRPC.MaterialEnum_Material_Fastners
	case "Material_Moulding_Millwork":
		return rentalRPC.MaterialEnum_Material_Moulding_Millwork
	case "Material_Drywall":
		return rentalRPC.MaterialEnum_Material_Drywall
	case "Material_Doors_Windows":
		return rentalRPC.MaterialEnum_Material_Doors_Windows
	case "Material_Roofing_Gutters":
		return rentalRPC.MaterialEnum_Material_Roofing_Gutters
	case "Material_Ladders":
		return rentalRPC.MaterialEnum_Material_Ladders
	case "Material_Scaffolding":
		return rentalRPC.MaterialEnum_Material_Scaffolding
	case "Material_Plumbing":
		return rentalRPC.MaterialEnum_Material_Plumbing
	case "Material_Siding":
		return rentalRPC.MaterialEnum_Material_Siding
	case "Material_Insulation":
		return rentalRPC.MaterialEnum_Material_Insulation
	case "Material_Ceilings":
		return rentalRPC.MaterialEnum_Material_Ceilings
	case "Material_Wall_Paneling":
		return rentalRPC.MaterialEnum_Material_Wall_Paneling
	case "Material_Flooring":
		return rentalRPC.MaterialEnum_Material_Flooring
	case "Material_Concrete_Cement_Masonry":
		return rentalRPC.MaterialEnum_Material_Concrete_Cement_Masonry
	case "Material_Material_Handling_Equipment":
		return rentalRPC.MaterialEnum_Material_Material_Handling_Equipment
	case "Material_Building_Hardware":
		return rentalRPC.MaterialEnum_Material_Building_Hardware
	case "Material_Glass_and_Plastic_Sheets":
		return rentalRPC.MaterialEnum_Material_Glass_and_Plastic_Sheets
	case "Material_Heating_venting_Cooling":
		return rentalRPC.MaterialEnum_Material_Heating_venting_Cooling
	case "Material_Other":
		return rentalRPC.MaterialEnum_Material_Other

	}

	return rentalRPC.MaterialEnum_Any_Material

}

func rentalServiceToRPC(data string) rentalRPC.ServiceEnum {

	switch data {
	case "Service_Auto_Transport":
		return rentalRPC.ServiceEnum_Service_Auto_Transport
	case "Service_Storage":
		return rentalRPC.ServiceEnum_Service_Storage
	case "Service_Moving_Supplies":
		return rentalRPC.ServiceEnum_Service_Moving_Supplies
	case "Service_Furniture_Movers":
		return rentalRPC.ServiceEnum_Service_Furniture_Movers
	}

	return rentalRPC.ServiceEnum_Any_Service
}

func rentalServiceRPCToString(data rentalRPC.ServiceEnum) string {

	switch data {
	case rentalRPC.ServiceEnum_Service_Auto_Transport:
		return "Service_Auto_Transport"
	case rentalRPC.ServiceEnum_Service_Storage:
		return "Service_Storage"
	case rentalRPC.ServiceEnum_Service_Moving_Supplies:
		return "Service_Moving_Supplies"
	case rentalRPC.ServiceEnum_Service_Furniture_Movers:
		return "Service_Furniture_Movers"
	}

	return "Any_Service"
}

func rentalTimingToRPC(data string) rentalRPC.TimingEnum {
	switch data {
	case "Timing_Flexible":
		return rentalRPC.TimingEnum_Timing_Flexible
	case "Timing_6Months":
		return rentalRPC.TimingEnum_Timing_6Months
	case "Timing_Year":
		return rentalRPC.TimingEnum_Timing_Year
	}

	return rentalRPC.TimingEnum_Any_Timing
}

func rentalTimingRPCToString(data rentalRPC.TimingEnum) string {
	switch data {
	case rentalRPC.TimingEnum_Timing_Flexible:
		return "Timing_Flexible"
	case rentalRPC.TimingEnum_Timing_6Months:
		return "Timing_6Months"
	case rentalRPC.TimingEnum_Timing_Year:
		return "Timing_Year"
	}

	return "Any_Timing"
}

func rentalLocationTypeToRPC(data string) rentalRPC.LocationEnum {
	switch data {
	case "Location_Local":
		return rentalRPC.LocationEnum_Location_Local
	case "Location_International":
		return rentalRPC.LocationEnum_Location_International
	}
	return rentalRPC.LocationEnum_Any_Location
}

func rentalLocationRPCToString(data rentalRPC.LocationEnum) string {
	switch data {
	case rentalRPC.LocationEnum_Location_Local:
		return "Location_Local"
	case rentalRPC.LocationEnum_Location_International:
		return "Location_International"
	}
	return "Any_Location"
}
