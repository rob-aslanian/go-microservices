package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/shop/internal/file"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/shop"
)

// Service define functions inside Service
type Service interface {
	CreateShop(ctx context.Context, data *shop.Shop) (string, error)
	AddProduct(ctx context.Context, companyID string, data *shop.Product) (string, error)
	GetProduct(ctx context.Context, id string) (*shop.Product, error)
	GetShop(ctx context.Context, id string) (*shop.Shop, error)
	ChangeShowcase(ctx context.Context, companyID string, shopID string, showcase string) error
	MakeOrder(ctx context.Context, productIDs []string, address *shop.Address) error
	AddToWishlist(ctx context.Context, productID string) error
	GetMyWishlist(ctx context.Context, first, after uint) ([]*shop.Product, error)
	GetOrdersForBuyer(ctx context.Context, first, after uint) ([]*shop.Order, error)
	GetOrdersForSeller(ctx context.Context, companyID, shopID string, first, after uint) ([]*shop.Order, error)
	ChangeOrderStatus(ctx context.Context, companyID string, orderID string, orderStatus string) error
	GetProductsWithDeleted(ctx context.Context, ids []string) ([]*shop.Product, error)
	ChangeProduct(ctx context.Context, companyID string, prod *shop.Product) error

	ChangeLogo(ctx context.Context, companyID string, shopID string, url string) error
	ChangeCover(ctx context.Context, companyID string, shopID string, url string) error
	RemoveLogo(ctx context.Context, companyID string, shopID string) error
	RemoveCover(ctx context.Context, companyID string, shopID string) error
	ChangeImagesInProduct(ctx context.Context, companyID, productID string, images []*file.File) error
	RemoveProduct(ctx context.Context, companyID, productID string) error
	HideProduct(ctx context.Context, companyID string, productID string, value bool) error
	GetMyShops(ctx context.Context, companyID string) ([]*shop.Shop, error)
	FindProducts(ctx context.Context, shopID string, filter *shop.SearchFilter, first, after uint) ([]*shop.Product, error)
}
