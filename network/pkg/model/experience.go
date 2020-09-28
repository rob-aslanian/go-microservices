package model

import "gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"

type Experience struct {
	Company          Company `json:"company"`
	Title            string  `json:"title"`
	FromMonth        int     `json:"from_month"`
	FromYear         int     `json:"from_year"`
	ToMonth          int     `json:"to_month"`
	ToYear           int     `json:"to_year"`
	WorkingCurrently bool    `json:"working_currently"`
	Description      string  `json:"description"`
}

func (e *Experience) ToRPC() *networkRPC.Experience {
	return &networkRPC.Experience{
		Company:          e.Company.ToRPC(),
		Title:            e.Title,
		StartMonth:       int32(e.FromMonth),
		StartYear:        int32(e.FromYear),
		EndMonth:         int32(e.ToMonth),
		EndYear:          int32(e.ToYear),
		WorkingCurrently: e.WorkingCurrently,
		Description:      e.Description,
	}
}

type ExperienceArr []*Experience

func (experiences ExperienceArr) ToRPC() *networkRPC.ExperienceArr {
	networkExps := make([]*networkRPC.Experience, len(experiences))
	for i, exp := range experiences {
		networkExps[i] = exp.ToRPC()
	}
	return &networkRPC.ExperienceArr{Experiences: networkExps}
}

type AddExperienceRequest struct {
	UserId           string `json:"_from"`
	CompanyId        string `json:"_to"`
	Title            string `json:"title" validate:"min=3,max=50"`
	FromMonth        int    `json:"from_month" validate:"min=1,max=12"`
	FromYear         int    `json:"from_year" validate:"min=1900,max=2100"`
	ToMonth          int    `json:"to_month"`
	ToYear           int    `json:"to_year"`
	WorkingCurrently bool   `json:"working_currently"`
	Description      string `json:"description" validate:"omitempty,max=500"`
}

func NewAddExperienceRequestFromRPC(r *networkRPC.AddExperienceRequest) *AddExperienceRequest {
	return &AddExperienceRequest{
		CompanyId:        r.CompanyId,
		Title:            r.Title,
		FromMonth:        int(r.StartMonth),
		FromYear:         int(r.StartYear),
		ToMonth:          int(r.EndMonth),
		ToYear:           int(r.EndYear),
		WorkingCurrently: r.WorkingCurrently,
		Description:      r.Description,
	}
}
