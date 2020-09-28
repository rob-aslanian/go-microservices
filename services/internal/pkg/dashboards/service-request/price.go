package servicerequest

// Price which servicer/requester will get/get according to agreemant. it can be fixed(total 1000$), hourly(100$ for each hour) or negotiable between the seller-buyer
type Price string

const (
	// PriceFixed Price = "Fixed"
	PriceFixed Price = "Fixed"

	// PriceHourly Price = "Hourly"
	PriceHourly Price = "Hourly"

	// PriceNegotiable Price = "Negotiable"
	PriceNegotiable Price = "Negotiable"
)

var price = map[string]int{
	"Hour": 1,
}

// GetHours returns amount of hours of given interval
func (p Price) GetHours() int {
	return price[string(p)]
}
