package model

type Report struct {
	UserId         string `bson:"user_id"`
	CompanyId      string `bson:"company_id"`
	ConversationId string `bson:"conversation_id"`
	Text           string `bson:"text"`
}
