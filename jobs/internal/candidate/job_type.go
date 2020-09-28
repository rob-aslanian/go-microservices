package candidate

// JobType type of job.
type JobType string

const (
	// JobTypeUnknown unknown type
	JobTypeUnknown JobType = "Unknown"
	// JobTypeFullTime full time type
	JobTypeFullTime JobType = "FullTime"
	// JobTypePartTime part time type
	JobTypePartTime JobType = "PartTime"
	// JobTypePartner partner type
	JobTypePartner JobType = "Partner"
	// JobTypeContractual contractual type
	JobTypeContractual JobType = "Contractual"
	// JobTypeVolunteer volunteer type
	JobTypeVolunteer JobType = "Volunteer"
	// JobTypeTemporary temporary type
	JobTypeTemporary JobType = "Temporary"
	// JobTypeConsultancy consultancy type
	JobTypeConsultancy JobType = "Consultancy"
	// JobTypeInternship internship type
	JobTypeInternship JobType = "Internship"
)
