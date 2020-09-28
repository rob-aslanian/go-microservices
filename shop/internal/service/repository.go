package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/shop/internal/file"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/shop"
)

// Repository ...
type Repository interface {
	CreateShop(ctx context.Context, data *shop.Shop) error
	AddProduct(ctx context.Context, data *shop.Product) error
	GetProduct(ctx context.Context, id string, withDeleted bool) (*shop.Product, error)
	GetShop(ctx context.Context, id string) (*shop.Shop, error)
	ChangeShowcase(ctx context.Context, shopID string, showcase string) error
	MakeOrder(ctx context.Context, orders []shop.Order) error
	AddToWishlist(ctx context.Context, userID string, productID string) error
	GetProductsIDFromWishlist(ctx context.Context, userID string) ([]string, error)
	GetProducts(ctx context.Context, ids []string, first, after uint) ([]*shop.Product, error)
	GetOrdersIDForBuyer(ctx context.Context, userID string) ([]string, error)
	GetOrders(ctx context.Context, ids []string, first, after uint) ([]*shop.Order, error)
	GetOrdersIDForSeller(ctx context.Context, shopID string) ([]string, error)
	ChangeOrderStatus(ctx context.Context, orderID string, orderStatus string) error
	GetProductsWithDeleted(ctx context.Context, ids []string) ([]*shop.Product, error)
	ChangeLogo(ctx context.Context, shopID, url string) error
	ChangeCover(ctx context.Context, shopID, url string) error
	RemoveLogo(ctx context.Context, shopID string) error
	RemoveCover(ctx context.Context, shopID string) error
	ChangeImagesInProduct(ctx context.Context, productID string, images []*file.File) error
	ChangeProduct(ctx context.Context, productID string, prod *shop.Product) error
	RemoveProduct(ctx context.Context, productID string) error
	GetAmountOfShops(ctx context.Context, ownerID string) (uint8, error)
	HideProduct(ctx context.Context, productID string, value bool) error
	GetMyShops(ctx context.Context, ownerID string) ([]*shop.Shop, error)
	FindProducts(ctx context.Context, shopID string, filter *shop.SearchFilter, first, after uint) ([]*shop.Product, error)
}
