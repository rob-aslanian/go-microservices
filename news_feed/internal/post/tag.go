package post

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tag ...
type Tag struct {
	ID   primitive.ObjectID `bson:"id" json:"id"`
	Type string             `bson:"type" json:"type"`
}

// GetID returns id of ad
func (t Tag) GetID() string {
	return t.ID.Hex()
}

// SetID ...
func (t *Tag) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	t.ID = objID
	return nil
}

// MarshalJSON ...
func (t *Tag) MarshalJSON() ([]byte, error) {
	type TagAlias Tag

	m, err := json.Marshal(&struct {
		ID string `json:"id"`
		*TagAlias
	}{
		ID:       t.GetID(),
		TagAlias: (*TagAlias)(t),
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}
