package notification

// Settings ...
type Settings struct {
	NewEndorsement        *bool `bson:"new_endorsement"`
	NewFollow             *bool `bson:"new_follow"`
	NewConnection         *bool `bson:"new_connection"`
	ApprovedConnection    *bool `bson:"approved_connection"`
	RecommendationRequest *bool `bson:"recommendation_request"`
	NewRecommendation     *bool `bson:"new_recommendation"`
	NewJobInvitation      *bool `bson:"new_job_invitation"`
}

// CompanySettings ...
type CompanySettings struct {
	NewFollow    *bool `bson:"new_follow"`
	NewReview    *bool `bson:"new_review"`
	NewApplicant *bool `bson:"new_applicant"`
}

// ParameterSetting ...
type ParameterSetting string

const (
	// ParameterUnknonwn ...
	ParameterUnknonwn ParameterSetting = "unknown" // user and company
	// NewEndorsement ...
	NewEndorsement ParameterSetting = "new_endorsement" // user
	// NewFollow ...
	NewFollow ParameterSetting = "new_follow" // user and company
	// NewConnection ...
	NewConnection ParameterSetting = "new_connection" // user
	// ApprovedConnection ...
	ApprovedConnection ParameterSetting = "approved_connection" // user
	// RecommendationRequest ...
	RecommendationRequest ParameterSetting = "recommendation_request" // user
	// NewRecommendation ...
	NewRecommendation ParameterSetting = "new_recommendation" // user
	// NewReview ...
	NewReview ParameterSetting = "new_review" // company
	// NewJobInvitation ...
	NewJobInvitation ParameterSetting = "new_job_invitation" // user
	// NewApplicant ...
	NewApplicant ParameterSetting = "new_applicant" // company
)
