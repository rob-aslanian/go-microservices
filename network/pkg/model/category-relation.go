package model

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"time"
)

type CategoryTree struct {
	OwnerId    string         `json:"_key"`
	Categories []CategoryItem `json:"categories"`
}

func (this *CategoryTree) ToRPC() *networkRPC.CategoryTree {
	arr := make([]*networkRPC.CategoryItem, len(this.Categories))
	for i, item := range this.Categories {
		arr[i] = item.ToRPC()
	}
	return &networkRPC.CategoryTree{Categories: arr}
}

type CategoryItem struct {
	UniqueName  string         `json:"unique_name"`
	Name        string         `json:"name" validate:"min=3,max=50"`
	HasChildren bool           `json:"has_children"`
	Children    []CategoryItem `json:"children"`
}

func (this CategoryItem) ToRPC() *networkRPC.CategoryItem {
	arr := make([]*networkRPC.CategoryItem, len(this.Children))
	for i, item := range this.Children {
		arr[i] = item.ToRPC()
	}
	return &networkRPC.CategoryItem{
		UniqueName:  this.UniqueName,
		Name:        this.Name,
		HasChildren: this.HasChildren,
		Children:    arr,
	}
}

type CategoryRelation struct {
	OwnerId    string `json:"_from"`
	ReferralId string `json:"_to"`

	CategoryName string `json:"name" validate:"min=3,max=50"`

	CreatedAt time.Time `json:"created_at"`
}
