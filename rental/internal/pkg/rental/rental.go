package rental

import (
	"time"

	file "gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/files"

	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/location"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CommonRental ...
type CommonRental struct {
	ID            primitive.ObjectID `bson:"_id"`
	OwnerID       primitive.ObjectID `bson:"owner_id"`
	IsCompany     bool               `bson:"is_company"`
	PostStatus    PostStatus         `bson:"post_status"`
	PostCurrency  string             `bson:"post_currency"`
	DealType      DealType           `bson:"deal_type"`
	PropertyType  PropertyType       `bson:"property_type"`
	Location      location.Location  `bson:"location"`
	Files         []*file.File       `bson:"files"`
	CreatedAt     time.Time          `bson:"created_at"`
	ExpiredDays   int32              `bson:"expired_days"`
	IsUrgent      bool               `bson:"is_urgent"`
	HasRepossesed bool               `bson:"has_repossesed"`
	Views         int32              `bson:"-"`
	Shares        int32              `bson:"-"`
	Alerts        int32              `bson:"-"`
	Offers        int32              `bson:"-"`
}

// GetRental ...
type GetRental struct {
	Amount               int32          `bson:"amount"`
	Appartaments         []Appartament  `bson:"appartments"`
	Homes                []Appartament  `bson:"homes"`
	NewHomes             []Appartament  `bson:"new_homes"`
	SummerCotages        []Appartament  `bson:"summer_cottage"`
	Houses               []Appartament  `bson:"houses"`
	Garages              []Garage       `bson:"garages"`
	StorageRooms         []StorageRooms `bson:"storage_rooms"`
	Offices              []Office       `bson:"offices"`
	CommercialProperties []Commercial   `bson:"commercial_properties"`
	Buildings            []Buildings    `bson:"buildings"`
	Lands                []Land         `bson:"land"`
	RuralFarms           []RuralFarm    `bson:"rural_farm"`
	Materials            []Materials    `bson:"materials"`
	Move                 []Move         `bson:"move"`
	Renovation           []Renovation   `bson:"renovation"`
}

// Appartament ...
type Appartament struct {
	RentalInfo      CommonRental      `bson:",inline"`
	TypeOfProperty  []TypeOfProperty  `bson:"type_of_property"`
	Status          Status            `bson:"status"`
	BadRooms        int32             `bson:"badrooms"`
	BathRooms       int32             `bson:"bathrooms"`
	TotalArea       int32             `bson:"total_area"`
	Metric          PriceType         `bson:"metric"`
	Floor           int32             `bson:"floor"`
	Floors          int32             `bson:"floors"`
	CarSpaces       int32             `bson:"car_spaces"`
	OutdoorFeatures []OutdoorFeatures `bson:"outdoor_features,omitempty"`
	IndoorFeatures  []IndoorFeatures  `bson:"indoor_features,omitempty"`
	ClimatControl   []ClimatControl   `bson:"climate_control,omitempty"`
	AvailibatiFrom  string            `bson:"avilable_from"`
	AvailibatiTo    string            `bson:"avilable_to"`
	Details         []Detail          `bson:"detail"`
	Price           Price             `bson:"price"`
	Phones          []Phone           `bson:"phones"`
	IsAgent         bool              `bson:"is_agent"`
	WhoLive         []WhoLive         `bson:"who_live,omitempty"`
}

// Garage ...
type Garage struct {
	RentalInfo        CommonRental        `bson:",inline"`
	AdditionalFilters []AdditionalFilters `bson:"additional_filters,omitempty"`
	TotalArea         int32               `bson:"total_area"`
	Metric            PriceType           `bson:"metric"`
	Details           []Detail            `bson:"detail"`
	IsAgent           bool                `bson:"is_agent"`
	Price             Price               `bson:"price"`
	Phones            []Phone             `bson:"phones"`
}

// StorageRooms ...
type StorageRooms struct {
	RentalInfo     CommonRental `bson:",inline"`
	Status         Status       `bson:"status"`
	TotalArea      int32        `bson:"total_area"`
	Metric         PriceType    `bson:"metric"`
	Details        []Detail     `bson:"detail"`
	IsAgent        bool         `bson:"is_agent"`
	Price          Price        `bson:"price"`
	Phones         []Phone      `bson:"phones"`
	AvailibatiFrom string       `bson:"avilable_from"`
	AvailibatiTo   string       `bson:"avilable_to"`
}

// Buildings ...
type Buildings struct {
	RentalInfo     CommonRental `bson:",inline"`
	Status         Status       `bson:"status,omitempty"`
	TotalArea      int32        `bson:"total_area"`
	Metric         PriceType    `bson:"metric"`
	AvailibatiFrom string       `bson:"avilable_from,omitempty"`
	AvailibatiTo   string       `bson:"avilable_to,omitempty"`
	Details        []Detail     `bson:"detail"`
	Price          Price        `bson:"price"`
	Phones         []Phone      `bson:"phones"`
	IsAgent        bool         `bson:"is_agent"`
}

// HotelRooms ...
type HotelRooms struct {
	RentalInfo     CommonRental `bson:",inline"`
	Status         Status       `bson:"status,omitempty"`
	Rooms          int32        `bson:"rooms"`
	TotalArea      int32        `bson:"total_area"`
	Metric         PriceType    `bson:"metric"`
	AvailibatiFrom string       `bson:"avilable_from"`
	AvailibatiTo   string       `bson:"avilable_to"`
	Details        []Detail     `bson:"detail"`
	Price          Price        `bson:"price"`
	IsAgent        bool         `bson:"is_agent"`
	Phones         []Phone      `bson:"phones"`
}

// Commercial ...
type Commercial struct {
	RentalInfo           CommonRental                 `bson:",inline"`
	CommercialProperties []CommericalProperty         `bson:"commercial_property,omitempty"`
	CommericalLocation   []CommericalPropertyLocation `bson:"commerical_location,omitempty"`
	AdditionalFilters    []AdditionalFilters          `bson:"additional_filters,omitempty"`
	Status               Status                       `bson:"status"`
	AvailibatiFrom       string                       `bson:"avilable_from"`
	AvailibatiTo         string                       `bson:"avilable_to"`
	TotalArea            int32                        `bson:"total_area"`
	Metric               PriceType                    `bson:"metric"`
	Details              []Detail                     `bson:"detail"`
	Price                Price                        `bson:"price"`
	IsAgent              bool                         `bson:"is_agent"`
	Phones               []Phone                      `bson:"phones"`
}

// Land ...
type Land struct {
	RentalInfo     CommonRental        `bson:",inline"`
	TypeOfLand     []Status            `bson:"type_of_land"`
	More           []AdditionalFilters `bson:"more"`
	TotalArea      int32               `bson:"total_area"`
	Metric         PriceType           `bson:"metric"`
	AvailibatiFrom string              `bson:"avilable_from"`
	AvailibatiTo   string              `bson:"avilable_to"`
	Details        []Detail            `bson:"detail"`
	Price          Price               `bson:"price"`
	Phones         []Phone             `bson:"phones"`
	IsAgent        bool                `bson:"is_agent"`
}

// Office ...
type Office struct {
	RentalInfo     CommonRental `bson:",inline"`
	Layout         Layout       `bson:"layout"`
	BuildingUse    BuildingUse  `bson:"building_use"`
	Status         Status       `bson:"status"`
	TotalArea      int32        `bson:"total_area"`
	Metric         PriceType    `bson:"metric"`
	AvailibatiFrom string       `bson:"avilable_from"`
	AvailibatiTo   string       `bson:"avilable_to"`
	Details        []Detail     `bson:"detail"`
	Price          Price        `bson:"price"`
	IsAgent        bool         `bson:"is_agent"`
	Phones         []Phone      `bson:"phones"`
}

// RuralFarm ...
type RuralFarm struct {
	RentalInfo     CommonRental      `bson:",inline"`
	PropertyType   []PropertyType    `bson:"property_types,omitempty"`
	Additional     []OutdoorFeatures `bson:"additional,omitempty"`
	Status         Status            `bson:"status"`
	TotalArea      int32             `bson:"total_area"`
	Metric         PriceType         `bson:"metric"`
	AvailibatiFrom string            `bson:"avilable_from"`
	AvailibatiTo   string            `bson:"avilable_to"`
	Details        []Detail          `bson:"detail"`
	IsAgent        bool              `bson:"is_agent"`
	Price          Price             `bson:"price"`
	Phones         []Phone           `bson:"phones"`
}

// Build
type Build struct {
	RentalInfo    CommonRental `bson:",inline"`
	PurchasePrice Price        `bson:"purchase_price"`
	TotalArea     int32        `bson:"total_area"`
	Metric        PriceType    `bson:"metric"`
	Details       []Detail     `bson:"detail"`
	Timing        Timing       `bson:"timing"`
	IsAgent       bool         `bson:"is_agent"`
}

// Move ...
type Move struct {
	RentalInfo   CommonRental `bson:",inline"`
	LocationType []string     `bson:"location_type"`
	CountryIDs   []string     `bson:"country_ids,omitempty"`
	Services     []Service    `bson:"services,omitempty"`
	Details      []Detail     `bson:"detail"`
}

// Materials ...
type Materials struct {
	RentalInfo CommonRental `bson:",inline"`
	Materials  []Material   `bson:"materials,omitempty"`
	Details    []Detail     `bson:"detail"`
}

// Renovation ...
type Renovation struct {
	RentalInfo          CommonRental `bson:",inline"`
	CountryIDs          []string     `bson:"country_ids,omitempty"`
	CityIDs             []string     `bson:"city_ids,omitempty"`
	Exetior             Price        `bson:"exetior"`
	Interior            Price        `bson:"interior"`
	InteriorAndExterior Price        `bson:"interior_exterior"`
	Timing              Timing       `bson:"timing"`
	Details             []Detail     `bson:"detail"`
}

// GetID returns id
func (p CommonRental) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *CommonRental) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *CommonRental) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetOwnerID returns owner id
func (p CommonRental) GetOwnerID() string {
	return p.OwnerID.Hex()
}

// SetOwnerID set owner id
func (p *CommonRental) SetOwnerID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OwnerID = objID
	return nil
}
