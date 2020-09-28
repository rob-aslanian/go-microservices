package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

// ServiceSearch ...
type ServiceSearch struct {
	Keyword       []string
	CityID        []string
	CountryID     []string
	Remote        bool
	OnSiteWork    bool
	DeliveryTime  DeliveryTime
	PriceType     string
	FixedPrice    int32
	MinPrice      int32
	MaxPrice      int32
	CurrencyPrice string
	Skills        []string
	Languages     []string
	LocationType  LocationType
	Price         Price
	IsAlwaysOpen  bool
	WeekDays      []WeekDay
	HourFrom      string
	HourTo        string
	ServiceOwner  ServiceOwner
	First         uint32 `bson:"-"`
	After         string `bson:"-"`
}

// ServiceSearchFilter holds the fields by which the JobSearchFilter will be saved
type ServiceSearchFilter struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	Type      FilterType          `bson:"type"`
	Name      string              `bson:"name"`
	ServiceSearch
}

// GetID returns id
func (p ServiceSearchFilter) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *ServiceSearchFilter) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *ServiceSearchFilter) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns id
func (p ServiceSearchFilter) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return p.UserID.Hex()
}

// SetUserID set id
func (p *ServiceSearchFilter) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = &objID
	return nil
}

// GetCompanyID returns id
func (p ServiceSearchFilter) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}
	return p.CompanyID.Hex()
}

// SetCompanyID set id
func (p *ServiceSearchFilter) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// WeekDay ...
type WeekDay string

const (
	// WeekDayMonday ...
	WeekDayMonday WeekDay = "monday"
	// WeekDayTuesday ...
	WeekDayTuesday WeekDay = "tuesday"
	// WeekDayWednesday ...
	WeekDayWednesday WeekDay = "wednesday"
	// WeekDayThursday ...
	WeekDayThursday WeekDay = "thursday"
	// WeekDayFriday ...
	WeekDayFriday WeekDay = "friday"
	// WeekDaySaturday ...
	WeekDaySaturday WeekDay = "saturday"
	// WeekDaySunday ...
	WeekDaySunday WeekDay = "sunday"
)

// Price ...
type Price string

const (
	// PriceAny Price = "Any"
	PriceAny Price = "Any"

	// PriceFixed Price = "Fixed"
	PriceFixed Price = "Fixed"

	// PriceHourly Price = "Hourly"
	PriceHourly Price = "Hourly"

	// PriceNegotiable Price = "Negotiable"
	PriceNegotiable Price = "Negotiable"
)

// LocationType ...
type LocationType string

const (
	// LocationTypeAny LocationType = "RemoTeOnly"
	LocationTypeAny = "Any"
	// LocationTypeRemoTeOnly LocationType = "RemoTeOnly"
	LocationTypeRemoTeOnly LocationType = "RemoteOnly"

	// LocationTypeOnSiteWork OnSiteWork = "OnSiteWork"
	LocationTypeOnSiteWork LocationType = "OnSiteWork"
)

// DeliveryTime ...
type DeliveryTime string

const (
	// DeliveryAny DeliveryTime = "UpTo24Hours"
	DeliveryAny DeliveryTime = "any"

	// DeliveryUpTo24Hours DeliveryTime = "UpTo24Hours"
	DeliveryUpTo24Hours DeliveryTime = "UpTo24Hours"

	// DeliveryUpTo3Days DeliveryTime = "UpTo3Days"
	DeliveryUpTo3Days DeliveryTime = "UpTo3Days"

	// DeliveryUpTo7Days DeliveryTime = "UpTo7Days"
	DeliveryUpTo7Days DeliveryTime = "UpTo7Days"

	// Delivery12Weeks DeliveryTime = "1-2Weeks"
	Delivery12Weeks DeliveryTime = "1-2Weeks"

	// Delivery2Weeks DeliveryTime = "2-4Weeks"
	Delivery2Weeks DeliveryTime = "2-4Weeks"

	// DeliveryMonthAndMore DeliveryTime = "MonthAndMore"
	DeliveryMonthAndMore DeliveryTime = "MonthAndMore"

	// DeliveryCustom DeliveryTime = "Custom"
	DeliveryCustom DeliveryTime = "Custom"
)

// ServiceOwner ...
type ServiceOwner string

const (
	// ServiceOwnerAny ...
	ServiceOwnerAny ServiceOwner = "Any_Owner"
	// ServiceOwnerUser ...
	ServiceOwnerUser ServiceOwner = "Owner_User"
	// ServiceOwnerCompany ...
	ServiceOwnerCompany ServiceOwner = "Owner_Company"
)
