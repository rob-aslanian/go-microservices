package company

import "go.mongodb.org/mongo-driver/bson/primitive"

// Details ...
type Details struct {
	CompanyID     primitive.ObjectID `bson:"company_id"`
	CompanyName   string             `bson:"company_name"`
	CompanyAvatar string
	CompanyURL    string
	CompanySize   string `bson:"company_size"`
	Industry      string
	Subindustry   string
}

// GetCompanyID returns user id
func (p Details) GetCompanyID() string {
	return p.CompanyID.Hex()
}

// SetCompanyID set user id
func (p *Details) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = objID
	return nil
}
