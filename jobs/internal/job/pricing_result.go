package job

// PricingResult ...
type PricingResult struct {
	Total       float64                  `bson:"total"`
	Currency    string                   `bson:"currency"`
	ByCountries []PricingResultByCountry `bson:"by_countries"`
}

// PricingResultByCountry ...
type PricingResultByCountry struct {
	Country            string             `bson:"country"`
	PlanPrice          float64            `bson:"plan_price"`
	RenewalPrice       float64            `bson:"renewal_price"`
	AdditionalFeatures map[string]float64 `bson:"features"`
	TotalPrice         float64            `bson:"total_price"`
}

// PlanPrices ...
type PlanPrices struct {
	Country   string        `bson:"country"`
	Currency  string        `bson:"currency"`
	Prices    PricesPerPlan `bson:"plans"`
	Features  Features      `bson:"features"`
	Discounts Discounts     `bson:"discounts"`
}

// PricesPerPlan ...
type PricesPerPlan struct {
	Basic            float64 `bson:"basic"`
	Start            float64 `bson:"start"`
	Standard         float64 `bson:"standard"`
	Professional     float64 `bson:"professional"`
	ProfessionalPlus float64 `bson:"professionalPlus"`
	Exclusive        float64 `bson:"exclusive"`
	Premium          float64 `bson:"premium"`
}

// Features ...
type Features struct {
	Anonymously float64 `bson:"publish_anonymously"`
	Language    float64 `bson:"language"`
}

// Discounts ...
type Discounts struct {
	Renewal []float64 `bson:"renewal"`
}
