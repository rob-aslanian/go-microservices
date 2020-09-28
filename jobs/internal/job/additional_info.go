package job

import suitable "gitlab.lan/Rightnao-site/microservices/jobs/internal/suitablefor"

type AdditionalInfo struct {
	SuitableFor       []suitable.SuitableFor `bson:"suitable_for"`
	TravelRequirement TravelRequirement      `bson:"travle_requirement"`
}
