package servicerequest

// DeliveryTime tells us how soon the service is expected to be delivered
type DeliveryTime string

const (
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
