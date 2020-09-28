package job

import "go.mongodb.org/mongo-driver/bson/primitive"

// ViewJobWithSeenStat ...
type ViewJobWithSeenStat struct {
	ID           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	Status       string             `bson:"status"`
	TotalAmount  int32              `bson:"total_amount"`
	UnseenAmount int32              `bson:"unseen_amount"`
}

// GetID returns id
func (p ViewJobWithSeenStat) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *ViewJobWithSeenStat) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}
