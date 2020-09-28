package status

// CompanyStatus ...
type CompanyStatus string

const (
	// CompanyStatusNotActivated ...
	CompanyStatusNotActivated CompanyStatus = "NOT_ACTIVATED"

	// CompanyStatusActivated ...
	CompanyStatusActivated CompanyStatus = "ACTIVATED"

	// CompanyStatusDeactivated ...
	CompanyStatusDeactivated CompanyStatus = "DEACTIVATED"
)
