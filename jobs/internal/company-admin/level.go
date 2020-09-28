/*Package companyadmin containg possible values of company admin level*/
package companyadmin

// AdminLevel ...
type AdminLevel string

const (
	// AdminLevelAdmin main admin level. Can everything.
	AdminLevelAdmin AdminLevel = "Admin"
	// AdminLevelJob only reponse for jobs.
	AdminLevelJob AdminLevel = "JobAdmin"
	// AdminLevelCommercial only reponse for adverts.
	AdminLevelCommercial AdminLevel = "CommercialAdmin"
	// AdminLevelVShop only reponse for shop.
	AdminLevelVShop AdminLevel = "VShopAdmin"
	// AdminLevelVService only reponse for services.
	AdminLevelVService AdminLevel = "VServiceAdmin"
	// AdminLevelUnknown not an admin.
	AdminLevelUnknown AdminLevel = "Unknown"
)
