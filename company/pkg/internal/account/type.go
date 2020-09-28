package account

// Type ...
type Type string

const (
	// TypeSelfEmployed ...
	TypeSelfEmployed Type = "self_employed"
	// TypeEducationalInstitution ...
	TypeEducationalInstitution Type = "educational_institution"
	// TypeGovernmentAgency ...
	TypeGovernmentAgency Type = "government_agency"
	// TypeSoleProprietorship ...
	TypeSoleProprietorship Type = "sole_proprietorship"
	// TypePrivatelyHeld ...
	TypePrivatelyHeld Type = "privately_held"
	// TypePartnership ...
	TypePartnership Type = "partnership"
	// TypePublicCompany ...
	TypePublicCompany Type = "public_company"
	// TypeUnknown ...
	TypeUnknown Type = "unknown"
)
