package model

import (
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

type RecommendationRequestModel struct {
	From      string    `json:"_from"`
	To        string    `json:"_to"`
	Text      string    `json:"text"`
	IsIgnored bool      `json:"is_ignored"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRecommendationRequestModel(request *networkRPC.RecommendationParams) *RecommendationRequestModel {
	return &RecommendationRequestModel{
		To:   request.UserId,
		Text: request.Text,
	}
}

type RecommendationRequest struct {
	Id        string    `json:"_key"`
	Requestor User      `json:"requestor"`
	Requested User      `json:"requested"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *RecommendationRequest) ToRPC() *networkRPC.RecommendationRequest {
	return &networkRPC.RecommendationRequest{
		Id:        r.Id,
		Requestor: r.Requestor.ToRPC(),
		Requested: r.Requested.ToRPC(),
		Text:      r.Text,
		CreatedAt: r.CreatedAt.Unix(),
	}
}

type RecommendationRequestArr []*RecommendationRequest

func (requests RecommendationRequestArr) ToRPC() *networkRPC.RecommendationRequestArr {
	arr := make([]*networkRPC.RecommendationRequest, len(requests))
	for i, r := range requests {
		arr[i] = r.ToRPC()
	}
	return &networkRPC.RecommendationRequestArr{Requests: arr}
}

type RecommendationModel struct {
	From      string    `json:"_from"`
	To        string    `json:"_to"`
	Text      string    `json:"text"`
	Relation  string    `json:"relation"`
	Title     string    `json:"title"`
	Hidden    *bool     `json:"hidden"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRecommendationModel(params *networkRPC.RecommendationParams) *RecommendationModel {
	return &RecommendationModel{
		To:       params.UserId,
		Text:     params.Text,
		Relation: params.GetRelations().String(),
		Title:    params.GetTitle(),
	}
}

type Recommendation struct {
	Id            string    `json:"_key"`
	Recommendator User      `json:"recommendator"`
	Receiver      User      `json:"receiver"`
	Text          string    `json:"text"`
	Hidden        *bool     `json:"hidden"`
	Relation      string    `json:"relation"`
	Title         string    `json:"title"`
	CreatedAt     time.Time `json:"created_at"`
}

func (r *Recommendation) ToRPC() *networkRPC.Recommendation {
	rec := networkRPC.Recommendation{
		Id:            r.Id,
		Text:          r.Text,
		Recommendator: r.Recommendator.ToRPC(),
		Receiver:      r.Receiver.ToRPC(),
		CreatedAt:     r.CreatedAt.Unix(),
		Title:         r.Title,
		Relation:      stringToRelationRPC(r.Relation),
	}

	if r.Hidden != nil {
		rec.Hidden = *r.Hidden
	} else {
		rec.IsHiddenNull = true
	}

	return &rec
}

func stringToRelationRPC(data string) userRPC.RecommendationRelationEnum {

	switch data {
	case "EXPERIENCE":
		return userRPC.RecommendationRelationEnum_EXPERIENCE
	case "EDUCATION":
		return userRPC.RecommendationRelationEnum_EDUCATION
	case "ACCOMPLISHMENT":
		return userRPC.RecommendationRelationEnum_ACCOMPLISHMENT
	}

	return userRPC.RecommendationRelationEnum_NO_RELATION
}

type RecommendationArr []*Recommendation

func (recommendations RecommendationArr) ToRPC() *networkRPC.RecommendationArr {
	arr := make([]*networkRPC.Recommendation, len(recommendations))
	for i, r := range recommendations {
		arr[i] = r.ToRPC()
	}
	return &networkRPC.RecommendationArr{Recommendations: arr}
}
