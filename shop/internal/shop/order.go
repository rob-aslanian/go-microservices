package shop

import (
	"time"

	"gitlab.lan/Rightnao-site/microservices/shop/internal/price"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order ...
type Order struct {
	ID           primitive.ObjectID  `bson:"_id"`
	ShopID       primitive.ObjectID  `bson:"shop_id"`
	UserID       *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID    *primitive.ObjectID `bson:"company_id,omitempty"`
	Address      *Address            `bson:"address"`
	ProductID    primitive.ObjectID  `bson:"product_id"`
	OrderStatus  string              `bson:"status"`
	CreatedAt    time.Time           `bson:"created_at"`
	DeliveryTime *time.Time          `bson:"delivery_at"`
	Quantity     uint32              `bson:"quantity"`
	IsPaid       bool                `bson:"is_paid"`
	Price        price.Price         `bson:"price"`
}

// GetID returns id of ad
func (g Order) GetID() string {
	return g.ID.Hex()
}

// SetID ...
func (g *Order) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.ID = objID
	return nil
}

// GenerateID creates new id
func (g *Order) GenerateID() string {
	g.ID = primitive.NewObjectID()
	return g.ID.Hex()
}

// GetUserID ...
func (g Order) GetUserID() string {
	if g.UserID == nil {
		return ""
	}

	return g.UserID.Hex()
}

// SetUserID ...
func (g *Order) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.UserID = &objID
	return nil
}

// GetCompanyID ...
func (g Order) GetCompanyID() string {
	if g.CompanyID == nil {
		return ""
	}

	return g.CompanyID.Hex()
}

// SetCompanyID ...
func (g *Order) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.CompanyID = &objID
	return nil
}

// GetProductID ...
func (g Order) GetProductID() string {

	return g.ProductID.Hex()
}

// SetProductID ...
func (g *Order) SetProductID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.ProductID = objID
	return nil
}

// GetShopID ...
func (g Order) GetShopID() string {

	return g.ShopID.Hex()
}

// SetShopID ...
func (g *Order) SetShopID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.ShopID = objID
	return nil
}
