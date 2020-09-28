package candidate

// ApplicantSort sort which is used in `GetJobApplicants`
type ApplicantSort int

const (
	// AppicantLastname sort by last name
	AppicantLastname ApplicantSort = iota
	// AppicantFirstname sort by first name
	AppicantFirstname
	// AppicantPostedDate sort by posted date
	AppicantPostedDate
	// AppicantExperienceYears sort by posted date
	AppicantExperienceYears
)
