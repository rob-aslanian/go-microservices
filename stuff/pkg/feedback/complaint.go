package feedback

type FeedBackComplaint struct {
	MissingOrWrong    string `bson:"missing_or_wrong"`
	ImproveExperience string `bson:"improve_experience"`
	TellUsMore        string `bson:"tell_us_more"`
}
