package account

// AdminLevel ...
type AdminLevel string

const (
	// AdminLevelAdmin ...
	AdminLevelAdmin AdminLevel = "Admin"
	// AdminLevelJob ...
	AdminLevelJob AdminLevel = "JobAdmin"
	// AdminLevelCommercial ...
	AdminLevelCommercial AdminLevel = "CommercialAdmin"
	// AdminLevelVShop ...
	AdminLevelVShop AdminLevel = "VShopAdmin"
	// AdminLevelVService ...
	AdminLevelVService AdminLevel = "VServiceAdmin"
	// AdminLevelUnknown ...
	AdminLevelUnknown AdminLevel = "Unknown"
)
