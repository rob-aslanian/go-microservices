package model

import "gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"

type Pagination struct {
	Amount int32
	After  int32
}

func NewPagination(p *networkRPC.Pagination) *Pagination {
	return &Pagination{
		Amount: p.Amount,
		After:  p.After,
	}
}
