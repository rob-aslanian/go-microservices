package feedback

type FeedBackCompliment struct {
	FavoriteFeatures  string `bson:"favorite_features"`
	ImproveExperience string `bson:"improve_experience"`
	ServicesToHave    string `bson:"services_to_have"`
}
