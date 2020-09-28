package rental

// DealType ...
type DealType string

const (
	// DealTypeAny ...
	DealTypeAny = "any"
	// DealTypeSell ...
	DealTypeSell DealType = "sell"
	// DealTypeRent ...
	DealTypeRent DealType = "rent"
	// DealTypeLease ...
	DealTypeLease DealType = "lease"
	// DealTypeShare ...
	DealTypeShare DealType = "share"
	// DealTypeBuild ...
	DealTypeBuild DealType = "build"
	// DealTypeMaterials ...
	DealTypeMaterials DealType = "materials"
	// DealTypeRenovation ...
	DealTypeRenovation DealType = "renovation"
	// DealTypeMove ...
	DealTypeMove DealType = "move"
)
