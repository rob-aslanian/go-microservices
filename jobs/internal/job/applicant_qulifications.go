package job

import "gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"

type ApplicantQualification struct {
	Experience     candidate.ExperienceEnum
	ToolTechnology []ToolTechnology `bson:"tools_technology"`
	Language       []Language       `bson:"languages"`
	Skills         []string         `bson:"skills"`
	Education      []string         `bson:"education"`
	License        string           `bson:"license"`
	Work           string           `bson:"work"`
}
