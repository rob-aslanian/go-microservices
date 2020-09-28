package service

import (
	"context"
	"errors"
	"log"
	"time"

	companyadmin "gitlab.lan/Rightnao-site/microservices/shop/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/file"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/shop"
	"google.golang.org/grpc/metadata"
)

// CreateShop ...
func (s Service) CreateShop(ctx context.Context, sh *shop.Shop) (string, error) {
	span := s.tracer.MakeSpan(ctx, "CreateShop")
	defer span.Finish()

	ownerID := ""

	if sh.GetCompanyID() != "" {
		if ok := s.checkAdminLevel(ctx, sh.GetCompanyID(), companyadmin.AdminLevelAdmin); !ok {
			return "", errors.New("not_allowed")
		}
		ownerID = sh.GetCompanyID()
	} else {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}
		sh.SetUserID(userID)
		ownerID = userID
	}

	// check amount of shops. Maximum 3.
	shopsAmount, err := s.repository.GetAmountOfShops(ctx, ownerID)
	if err != nil {
		return "", err
	}
	if shopsAmount >= 3 {
		return "", errors.New("too_many_shops")
	}

	id := sh.GenerateID()

	err = s.repository.CreateShop(ctx, sh)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// AddProduct ...
func (s Service) AddProduct(ctx context.Context, companyID string, pr *shop.Product) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddProduct")
	defer span.Finish()

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return "", errors.New("not_allowed")
		}
		pr.SetCompanyID(companyID)
	} else {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}
		pr.SetUserID(userID)
	}

	id := pr.GenerateID()
	err := pr.Validate()
	if err != nil {
		return "", err
	}
	pr.CreatedAt = time.Now()

	err = s.repository.AddProduct(ctx, pr)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// GetProduct ...
func (s Service) GetProduct(ctx context.Context, id string) (*shop.Product, error) {
	span := s.tracer.MakeSpan(ctx, "GetProduct")
	defer span.Finish()

	pr, err := s.repository.GetProduct(ctx, id, false)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return pr, nil
}

// GetShop ...
func (s Service) GetShop(ctx context.Context, id string) (*shop.Shop, error) {
	span := s.tracer.MakeSpan(ctx, "GetShop")
	defer span.Finish()

	pr, err := s.repository.GetShop(ctx, id)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return pr, nil
}

// ChangeShowcase ...
func (s Service) ChangeShowcase(ctx context.Context, companyID string, shopID string, showcase string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeShowcase")
	defer span.Finish()

	// TODO: check owner of shop

	err := s.repository.ChangeShowcase(ctx, shopID, showcase)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// MakeOrder ...
func (s Service) MakeOrder(ctx context.Context, productIDs []string, address *shop.Address) error {
	span := s.tracer.MakeSpan(ctx, "MakeOrder")
	defer span.Finish()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	orders := make([]shop.Order, len(productIDs))

	products, err := s.repository.GetProducts(ctx, productIDs, uint(len(productIDs)), 0)
	if err != nil {
		return err
	}

	if len(products) == 0 {
		return errors.New("nothing_found")
	}

	for i := range productIDs {
		orders[i] = shop.Order{
			Address: address,
		}

		for j := range products {
			if products[j].GetID() == productIDs[i] {
				orders[i].SetShopID(products[j].GetShopID())
				orders[i].Price = products[j].Price
				// TODO: apply discount
				break
			}
		}

		orders[i].GenerateID()
		orders[i].SetUserID(userID)
		orders[i].SetProductID(productIDs[i])
		orders[i].CreatedAt = time.Now()
		orders[i].OrderStatus = "Pending"
	}

	err = s.repository.MakeOrder(ctx, orders)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddToWishlist ...
func (s Service) AddToWishlist(ctx context.Context, productID string) error {
	span := s.tracer.MakeSpan(ctx, "AddToWishlist")
	defer span.Finish()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.AddToWishlist(ctx, userID, productID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetMyWishlist ...
func (s Service) GetMyWishlist(ctx context.Context, first, after uint) ([]*shop.Product, error) {
	span := s.tracer.MakeSpan(ctx, "GetMyWishlist")
	defer span.Finish()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids, err := s.repository.GetProductsIDFromWishlist(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	prs, err := s.repository.GetProducts(ctx, ids, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return prs, nil
}

// GetOrdersForBuyer ...
func (s Service) GetOrdersForBuyer(ctx context.Context, first, after uint) ([]*shop.Order, error) {
	span := s.tracer.MakeSpan(ctx, "GetOrdersForBuyer")
	defer span.Finish()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids, err := s.repository.GetOrdersIDForBuyer(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	orders, err := s.repository.GetOrders(ctx, ids, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return orders, nil
}

// GetOrdersForSeller ...
func (s Service) GetOrdersForSeller(ctx context.Context, companyID, shopID string, first, after uint) ([]*shop.Order, error) {
	span := s.tracer.MakeSpan(ctx, "GetOrdersForSeller")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return nil, errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return nil, errors.New("not_allowed")
	}

	ids, err := s.repository.GetOrdersIDForSeller(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	orders, err := s.repository.GetOrders(ctx, ids, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return orders, nil
}

// ChangeOrderStatus ...
func (s Service) ChangeOrderStatus(ctx context.Context, companyID string, orderID string, orderStatus string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOrderStatus")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	orders, err := s.repository.GetOrders(ctx, []string{orderID}, 1, 0)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if len(orders) < 1 {
		return errors.New("order_not_found")
	}

	product, err := s.repository.GetProduct(ctx, orders[0].GetProductID(), true)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if product.GetShopID() != userID &&
		product.GetShopID() != companyID &&
		product.GetUserID() != userID &&
		product.GetCompanyID() != companyID {
		return errors.New(`not_allowed`)
	}

	err = s.repository.ChangeOrderStatus(ctx, orderID, orderStatus)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetProductsWithDeleted ...
func (s Service) GetProductsWithDeleted(ctx context.Context, ids []string) ([]*shop.Product, error) {
	span := s.tracer.MakeSpan(ctx, "GetProductsWithDeleted")
	defer span.Finish()

	products, err := s.repository.GetProductsWithDeleted(ctx, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return products, nil
}

// ChangeProduct ...
func (s Service) ChangeProduct(ctx context.Context, companyID string, prod *shop.Product) error {
	span := s.tracer.MakeSpan(ctx, "ChangeProduct")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	prod2, err := s.repository.GetProduct(ctx, prod.GetID(), false)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, prod2.GetShopID())
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.ChangeProduct(ctx, prod.GetID(), prod)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeLogo ...
func (s Service) ChangeLogo(ctx context.Context, companyID string, shopID string, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeLogo")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.ChangeLogo(ctx, shopID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCover ...
func (s Service) ChangeCover(ctx context.Context, companyID string, shopID string, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCover")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.ChangeCover(ctx, shopID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveLogo ...
func (s Service) RemoveLogo(ctx context.Context, companyID string, shopID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveLogo")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.RemoveLogo(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveCover ...
func (s Service) RemoveCover(ctx context.Context, companyID string, shopID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveCover")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.RemoveCover(ctx, shopID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeImagesInProduct ...
func (s Service) ChangeImagesInProduct(ctx context.Context, companyID, productID string, images []*file.File) error {
	span := s.tracer.MakeSpan(ctx, "ChangeImagesInProduct")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	prod, err := s.repository.GetProduct(ctx, productID, false)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, prod.GetShopID())
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.ChangeImagesInProduct(ctx, productID, images)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveProduct ...
func (s Service) RemoveProduct(ctx context.Context, companyID, productID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveProduct")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	prod2, err := s.repository.GetProduct(ctx, productID, false)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, prod2.GetShopID())
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.RemoveProduct(ctx, productID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// HideProduct ...
func (s Service) HideProduct(ctx context.Context, companyID string, productID string, value bool) error {
	span := s.tracer.MakeSpan(ctx, "HideProduct")
	defer span.Finish()

	var userID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return errors.New("not_allowed")
		}
	} else {
		token := s.retriveToken(ctx)
		var err error
		userID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	prod2, err := s.repository.GetProduct(ctx, productID, false)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if owner
	sh, err := s.repository.GetShop(ctx, prod2.GetShopID())
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if sh.GetCompanyID() != companyID && sh.GetUserID() != userID {
		return errors.New("not_allowed")
	}

	err = s.repository.HideProduct(ctx, productID, value)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetMyShops ...
func (s Service) GetMyShops(ctx context.Context, companyID string) ([]*shop.Shop, error) {
	span := s.tracer.MakeSpan(ctx, "GetMyShops")
	defer span.Finish()

	var ownerID string

	if companyID != "" {
		if ok := s.checkAdminLevel(ctx, companyID, companyadmin.AdminLevelAdmin); !ok {
			return nil, errors.New("not_allowed")
		}
		ownerID = companyID
	} else {
		token := s.retriveToken(ctx)
		var err error
		ownerID, err = s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
	}

	shops, err := s.repository.GetMyShops(ctx, ownerID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return shops, nil
}

// FindProducts ...
func (s Service) FindProducts(ctx context.Context, shopID string, filter *shop.SearchFilter, first, after uint) ([]*shop.Product, error) {
	span := s.tracer.MakeSpan(ctx, "FindProducts")
	defer span.Finish()

	prod, err := s.repository.FindProducts(ctx, shopID, filter, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return prod, nil
}

// ---

func (s Service) retriveToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}

// checkAdminLevel return false if level doesn't much
func (s Service) checkAdminLevel(ctx context.Context, companyID string, requiredLevels ...companyadmin.AdminLevel) bool {
	span := s.tracer.MakeSpan(ctx, "checkAdminLevel")
	defer span.Finish()

	actualLevel, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		log.Println("Error: checkAdminLevel:", err)
		return false
	}

	for _, lvl := range requiredLevels {
		if lvl == actualLevel {
			return true
		}
	}

	return false
}
