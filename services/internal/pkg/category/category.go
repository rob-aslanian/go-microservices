package category

// Category is the list of which services can the seller/buyer offer/get. In v-office he can choose 3 main categories,
//  but in service/request - only one.
type Category struct {
	Main string   `bson:"main"`
	Sub  []string `bson:"sub"`
}

// VOfficeCategory is used only for Voffice when registering
type VOfficeCategory struct {
	Main []string `bson:"office_main"`
}
