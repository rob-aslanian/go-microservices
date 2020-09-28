package shop

import (
	"errors"
	"strings"
	"time"

	"gitlab.lan/Rightnao-site/microservices/shop/internal/file"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/price"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product ...
type Product struct {
	ID            primitive.ObjectID  `bson:"_id"`
	ShopID        primitive.ObjectID  `bson:"shop_id"`
	UserID        *primitive.ObjectID `bson:"user_id"`
	CompanyID     *primitive.ObjectID `bson:"company_id"`
	Title         string              `bson:"title"`
	Category      Category            `bson:"category"`
	Brand         *string             `bson:"brand"`
	IsUsed        bool                `bson:"is_used"`
	ProductType   string              `bson:"product_type"`
	Price         price.Price         `bson:"price"`
	SKU           string              `bson:"sku"`
	InStock       bool                `bson:"in_stock"`
	Quantity      uint32              `bson:"quantity"`
	Images        []file.File         `bson:"images,omitempty"`
	Description   string              `bson:"description"`
	Specification Specification       `bson:"specification"`
	Discount      Discount            `bson:"discount"`
	IsDeleted     bool                `bson:"is_deleted"`
	IsHidden      bool                `bson:"is_hidden"`
	CreatedAt     time.Time           `bson:"created_at"`
}

// GetID returns id of ad
func (g Product) GetID() string {
	return g.ID.Hex()
}

// SetID ...
func (g *Product) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.ID = objID
	return nil
}

// GenerateID creates new id
func (g *Product) GenerateID() string {
	g.ID = primitive.NewObjectID()
	return g.ID.Hex()
}

// GetShopID ...
func (g Product) GetShopID() string {
	return g.ShopID.Hex()
}

// SetShopID ...
func (g *Product) SetShopID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.ShopID = objID
	return nil
}

// GetUserID ...
func (g Product) GetUserID() string {
	if g.UserID == nil {
		return ""
	}

	return g.UserID.Hex()
}

// SetUserID ...
func (g *Product) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.UserID = &objID
	return nil
}

// GetCompanyID ...
func (g Product) GetCompanyID() string {
	if g.CompanyID == nil {
		return ""
	}

	return g.CompanyID.Hex()
}

// SetCompanyID ...
func (g *Product) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	g.CompanyID = &objID
	return nil
}

// Validate ...
func (g *Product) Validate() error {
	if strings.TrimSpace(g.Title) == "" {
		return errors.New("empty_title")
	}
	// if g.Category.Main == "" || g.Category.Sub == "" {
	// 	return errors.New("empty_category")
	// }
	if g.Price.Amount == 0 {
		return errors.New("empty_price")
	}
	if strings.TrimSpace(g.Price.Currency) == "" {
		return errors.New("empty_currency")
	}
	return nil
}

// Specification ...
type Specification struct {
	Size       string      `bson:"size"`
	Color      string      `bson:"color"`
	Material   string      `bson:"material"`
	Variations []Variation `bson:"variations,omitempty"`
}

// Variation ...
type Variation struct {
	Price    price.Price `bson:"price"`
	SKU      string      `bson:"sku"`
	InStock  bool        `bson:"in_stock"`
	Quantity uint32      `bson:"quantity"`
}

// Discount ...
type Discount struct {
	AmountOfProducts uint32    `bson:"amount_of_products"`
	DiscountType     string    `bson:"type"`
	DicountValue     string    `bson:"value"`
	StartDate        time.Time `bson:"start_date"`
	EndDate          time.Time `bson:"end_date"`
}
