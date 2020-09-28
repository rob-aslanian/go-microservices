package model

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"time"
)

type AdminEdge struct {
	UserId    string     `json:"_from"`
	CompanyId string     `json:"_to"`
	Level     AdminLevel `json:"level"`
	CreatedBy string     `json:"created_by_id"`
	CreatedAt time.Time  `json:"created_at"`
}

func NewAdminEdgeFromRPC(r *networkRPC.MakeCompanyAdminRequest) *AdminEdge {
	return &AdminEdge{
		UserId:    r.UserId,
		CompanyId: r.CompanyId,
		Level:     AdminLevelFromRPC(r.Level),
	}
}

type AdminLevel string

const AdminLevel_Admin = AdminLevel("Admin")
const AdminLevel_JobAdmin = AdminLevel("JobAdmin")
const AdminLevel_CommercialAdmin = AdminLevel("CommercialAdmin")
const AdminLevel_VShopAdmin = AdminLevel("VShopAdmin")
const AdminLevel_VServiceAdmin = AdminLevel("VServiceAdmin")

func (l AdminLevel) ToRPC() networkRPC.AdminLevel {
	return networkRPC.AdminLevel(networkRPC.AdminLevel_value[string(l)])
}

func AdminLevelFromRPC(l networkRPC.AdminLevel) AdminLevel {
	return AdminLevel(l.String())
}

type Admin struct {
	Id        string     `json:"_key"`
	User      User       `json:"user"`
	Company   Company    `json:"company"`
	Level     AdminLevel `json:"level"`
	CreatedBy User       `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
}

func (a *Admin) ToRPC() *networkRPC.AdminObject {
	return &networkRPC.AdminObject{
		User:      a.User.ToRPC(),
		Company:   a.Company.ToRPC(),
		Level:     a.Level.ToRPC(),
		CreatedBy: a.CreatedBy.ToRPC(),
		CreatedAt: a.CreatedAt.Unix(),
	}
}

type AdminArr []*Admin

func (admins AdminArr) ToRPC() *networkRPC.AdminObjectArr {
	networkAdmins := make([]*networkRPC.AdminObject, len(admins))
	for i, admin := range admins {
		networkAdmins[i] = admin.ToRPC()
	}
	return &networkRPC.AdminObjectArr{List: networkAdmins}
}
