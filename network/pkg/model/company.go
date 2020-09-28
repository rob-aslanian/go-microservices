package model

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"time"
)

type UserCompanyId struct {
	UserId    string `json:"_from" validate:"len=24,alphanum"`
	CompanyId string `json:"_to" validate:"len=24,alphanum"`
}

func NewUserCompanyFromRPC(req *networkRPC.UserCompanyId) *UserCompanyId {
	return &UserCompanyId{
		UserId:    req.UserId,
		CompanyId: req.CompanyId,
	}
}

type Company struct {
	Id             string    `json:"_key"`
	OwnerId        string    `json:"owner_id"`
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	Industry       string    //`json:"industry"`
	Type           string    `json:"type"`
	Address        string    `json:"address"`
	Apartment      string    `json:"apartment"`
	Zip            string    `json:"zip"`
	FoundationYear uint      `json:"foundation_year"`
	Email          string    `json:"email"`
	CreatedAt      time.Time //`json:"created_at"`
	Status         string    `json:"status"`
}

func (c *Company) ToRPC() *networkRPC.Company {
	return &networkRPC.Company{
		Id:             c.Id,
		Name:           c.Name,
		Url:            c.URL,
		Industry:       c.Industry,
		Type:           c.Type,
		Address:        c.Address,
		Email:          c.Email,
		FoundationYear: int32(c.FoundationYear),
	}
}

type CompanyArr []*Company

func (companies CompanyArr) ToRPC() *networkRPC.CompanyArr {
	networkComps := make([]*networkRPC.Company, len(companies))
	for i, company := range companies {
		networkComps[i] = company.ToRPC()
	}
	return &networkRPC.CompanyArr{Companies: networkComps}
}

type CompanyFollow struct {
	Company    Company  `json:"company"`
	Following  bool     `json:"following"`
	Followers  int      `json:"followers"`
	Rating     int      `json:"rating"`
	Size       int      `json:"size"`
	Categories []string `json:"categories"`
}

func (f *CompanyFollow) ToRPC() *networkRPC.CompanyFollowInfo {
	return &networkRPC.CompanyFollowInfo{
		Company:    f.Company.ToRPC(),
		Followers:  int32(f.Followers),
		Following:  f.Following,
		Rating:     int32(f.Rating),
		Size:       int32(f.Size),
		Categories: f.Categories,
	}
}

type CompanyFollowArr []*CompanyFollow

func (follows CompanyFollowArr) ToRPC() *networkRPC.CompanyFollowsArr {
	networkFollows := make([]*networkRPC.CompanyFollowInfo, len(follows))
	for i, follow := range follows {
		networkFollows[i] = follow.ToRPC()
	}
	return &networkRPC.CompanyFollowsArr{Follows: networkFollows}
}
