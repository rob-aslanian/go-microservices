package job

// Meta ...
type Meta struct {
	AdvertisementCountries []string  `bson:"advertisement_countries"`
	Highlight              Highlight `bson:"highlight"`
	Renewal                int32
	// JobPlan                Plan `bson:"job_plan"`
	AmountOfDays   int32 `bson:"amount_of_days"`
	Anonymous      bool
	NumOfLanguages int32 `bson:"num_of_languages"`
	Currency       string
}
