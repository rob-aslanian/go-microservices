package careercenter

// CVOptions ...
type CVOptions struct {
	YoungProfessionals       bool `bson:"young_professionals"`
	ExpierencedProfessionals bool `bson:"expierenced_professionals"`
	NewJobSeekers            bool `bson:"new_job_seekers"`
}
