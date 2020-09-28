package resolver

import (
	"context"
	"log"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/shopRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

// GetShop ...
func (_ *Resolver) GetShop(ctx context.Context, data GetShopRequest) (*ShopResolver, error) {
	res, err := shop.GetShop(ctx, &shopRPC.ID{
		ID: data.ID,
	})
	if err != nil {
		return nil, err
	}

	sh := shopRPCToShop(res)

	// ---

	sh.User = new(Profile)
	sh.Company = new(CompanyProfile)

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: []string{res.GetCompanyID()},
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: []string{res.GetUserID()},
	})
	if err != nil {
		return nil, err
	}

	if res.GetCompanyID() != "" {
		// company profile
		var profile *companyRPC.Profile

		for _, company := range companyResp.GetProfiles() {
			if company.GetId() == res.GetCompanyID() {
				profile = company
				break
			}
		}

		if profile != nil {
			pr := toCompanyProfile(ctx, *profile)
			sh.Company = &pr
		}
	} else {
		// user profile
		profile := ToProfile(ctx, userResp.GetProfiles()[res.GetUserID()])
		sh.User = &profile
	}

	// ---

	return &ShopResolver{
		R: sh,
	}, nil
}

// GetProduct ...
func (_ *Resolver) GetProduct(ctx context.Context, data GetProductRequest) (*ProductResolver, error) {
	res, err := shop.GetProduct(ctx, &shopRPC.ID{
		ID: data.ID,
	})
	if err != nil {
		return nil, err
	}

	return &ProductResolver{
		R: shopRPCProductToProduct(res),
	}, nil
}

// GetMyWishlist ...
func (_ *Resolver) GetMyWishlist(ctx context.Context, data GetMyWishlistRequest) (*[]ProductResolver, error) {
	var after, first string

	if data.Pagination.After != nil {
		after = *data.Pagination.After
	}
	if data.Pagination.First != nil {
		first = strconv.Itoa(int(*data.Pagination.First))
	}

	res, err := shop.GetMyWishlist(ctx, &shopRPC.IDWithPagination{
		// ID: data.
		After: after,
		First: first,
	})
	if err != nil {
		return nil, err
	}

	products := make([]ProductResolver, 0, len(res.GetProducts()))

	for i := range res.GetProducts() {
		products = append(products, ProductResolver{
			R: shopRPCProductToProduct(res.GetProducts()[i]),
		})
	}

	return &products, nil
}

// GetOrdersForBuyer ...
func (_ *Resolver) GetOrdersForBuyer(ctx context.Context, data GetOrdersForBuyerRequest) (*[]OrderResolver, error) {
	var after, first string

	if data.Pagination.After != nil {
		after = *data.Pagination.After
	}
	if data.Pagination.First != nil {
		first = strconv.Itoa(int(*data.Pagination.First))
	}

	res, err := shop.GetOrdersForBuyer(ctx, &shopRPC.IDWithPagination{
		After: after,
		First: first,
	})
	if err != nil {
		return nil, err
	}

	productIDs := make([]string, len(res.GetOrders()))
	for i := range res.GetOrders() {
		productIDs[i] = res.GetOrders()[i].GetProductID()
	}

	products, err := shop.GetProductsWithDeleted(ctx, &shopRPC.IDs{
		ID: productIDs,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]OrderResolver, 0, len(res.GetOrders()))

	log.Printf("orders: len: %d\n%+v\n", len(res.GetOrders()), res.GetOrders())

	for i := range res.GetOrders() {
		order := shopRPCOrderToOrder(ctx, res.GetOrders()[i])

		for j := range products.GetProducts() {
			if products.GetProducts()[j].GetID() == order.ID {
				if pr := shopRPCProductToProduct(products.GetProducts()[j]); pr != nil {
					order.Product = *pr
				}
				break
			}
		}

		orders = append(orders, OrderResolver{
			R: order,
		})
	}

	return &orders, nil
}

// GetOrdersForSeller ...
func (_ *Resolver) GetOrdersForSeller(ctx context.Context, data GetOrdersForSellerRequest) (*[]OrderResolver, error) {
	var after, first string

	if data.Pagination.After != nil {
		after = *data.Pagination.After
	}
	if data.Pagination.First != nil {
		first = strconv.Itoa(int(*data.Pagination.First))
	}

	res, err := shop.GetOrdersForSeller(ctx, &shopRPC.GetOrdersForSellerRequest{
		CompanyID: NullToString(data.Company_id),
		ShopID:    data.Shop_id,
		After:     after,
		First:     first,
	})
	if err != nil {
		return nil, err
	}

	productIDs := make([]string, len(res.GetOrders()))
	for i := range res.GetOrders() {
		productIDs[i] = res.GetOrders()[i].GetProductID()
	}

	products, err := shop.GetProductsWithDeleted(ctx, &shopRPC.IDs{
		ID: productIDs,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]OrderResolver, 0, len(res.GetOrders()))

	for i := range res.GetOrders() {
		order := shopRPCOrderToOrder(ctx, res.GetOrders()[i])

		for j := range products.GetProducts() {
			if products.GetProducts()[j] != nil && products.GetProducts()[j].GetID() == order.ID {
				if pr := shopRPCProductToProduct(products.GetProducts()[j]); pr != nil {
					order.Product = *pr
				}
				break
			}
		}

		orders = append(orders, OrderResolver{
			R: order,
		})
	}

	return &orders, nil
}

// CreateShop ...
func (_ *Resolver) CreateShop(ctx context.Context, data CreateShopRequest) (*SuccessResolver, error) {
	productTypes := make([]shopRPC.ProductsType, len(data.Input.ProductsType))

	for i := range data.Input.ProductsType {
		productTypes[i] = shopRPC.ProductsType(shopRPC.ProductsType_value[data.Input.ProductsType[i]])
	}

	res, err := shop.CreateShop(ctx, &shopRPC.CreateShopRequest{
		CompanyID:    NullToString(data.Company_id),
		Category:     data.Input.Category,
		Description:  data.Input.Description,
		ProductsType: productTypes,
		Title:        data.Input.Title,
		// SellerType
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      res.GetID(),
		},
	}, nil
}

// AddProduct ...
func (_ *Resolver) AddProduct(ctx context.Context, data AddProductRequest) (*SuccessResolver, error) {
	res, err := shop.AddProduct(ctx, &shopRPC.AddProductRequest{
		ShopID:        data.Input.ShopID,
		Brand:         NullToString(data.Input.Brand),
		Category:      categoryInputToShopRPCCategory(&data.Input.Category),
		CompanyID:     NullToString(data.Company_id),
		Description:   data.Input.Description,
		Discount:      discountInputToShopRPCDiscount(&data.Input.Discount),
		InStock:       data.Input.In_stock,
		IsUsed:        data.Input.Is_used,
		Price:         priceInputToShopRPCPrice(&data.Input.Price),
		ProductType:   shopRPC.ProductsType(shopRPC.ProductsType_value[data.Input.ProductType]),
		Quantity:      uint32(data.Input.Quantity),
		SKU:           data.Input.Sku,
		Specification: specificationInputToShopRPCSpecification(&data.Input.Specification),
		Title:         data.Input.Title,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

// ChangeShowcase ...
func (_ *Resolver) ChangeShowcase(ctx context.Context, data ChangeShowcaseRequest) (*SuccessResolver, error) {
	_, err := shop.ChangeShowcase(ctx, &shopRPC.ChangeShowcaseRequest{
		CompanyID: NullToString(data.Company_id),
		ShopID:    data.Input.ShopID,
		Showcase:  data.Input.Showcase,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// MakeOrder ...
func (_ *Resolver) MakeOrder(ctx context.Context, data MakeOrderRequest) (*SuccessResolver, error) {
	_, err := shop.MakeOrder(ctx, &shopRPC.MakeOrderRequest{
		ProductIDs: data.Input.Product_ids,
		Address:    addressInputToShopRPCAddress(&data.Input.Address),
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// AddToWishlist ...
func (_ *Resolver) AddToWishlist(ctx context.Context, data AddToWishlistRequest) (*SuccessResolver, error) {
	_, err := shop.AddToWishlist(ctx, &shopRPC.ID{
		ID: data.Product_id,
		// data.Company_id
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// ChangeOrderStatus ...
func (_ *Resolver) ChangeOrderStatus(ctx context.Context, data ChangeOrderStatusRequest) (*SuccessResolver, error) {
	_, err := shop.ChangeOrderStatus(ctx, &shopRPC.ChangeOrderStatusRequest{
		CompanyID:   NullToString(data.Company_id),
		OrderID:     data.Input.Order_id,
		OrderStatus: shopRPC.OrderStatus(shopRPC.OrderStatus_value[data.Input.Order_status]),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// RemoveShopLogo ...
func (_ *Resolver) RemoveShopLogo(ctx context.Context, data RemoveShopLogoRequest) (*SuccessResolver, error) {
	_, err := shop.RemoveLogo(ctx, &shopRPC.RemoveLogoRequest{
		CompanyID: NullToString(data.Company_id),
		ShopID:    data.Shop_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// RemoveShopCover ...
func (_ *Resolver) RemoveShopCover(ctx context.Context, data RemoveShopCoverRequest) (*SuccessResolver, error) {
	_, err := shop.RemoveCover(ctx, &shopRPC.RemoveLogoRequest{
		CompanyID: NullToString(data.Company_id),
		ShopID:    data.Shop_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// ChangeProduct ...
func (_ *Resolver) ChangeProduct(ctx context.Context, data ChangeProductRequest) (*SuccessResolver, error) {
	_, err := shop.ChangeProduct(ctx, &shopRPC.Product{
		ID: data.Product_id,
		// Brand:         NullToString(data.Input.Brand),
		Category:  categoryInputToShopRPCCategory(&data.Input.Category),
		CompanyID: NullToString(data.Company_id),
		// Description:   data.Input.Description,
		// Discount:      discountInputToShopRPCDiscount(&data.Input.Discount),
		InStock: data.Input.In_stock,
		IsUsed:  data.Input.Is_used,
		Price:   priceInputToShopRPCPrice(&data.Input.Price),
		// ProductType:   shopRPC.ProductsType(shopRPC.ProductsType_value[data.Input.ProductType]),
		Quantity:      uint32(data.Input.Quantity),
		SKU:           data.Input.Sku,
		Specification: specificationInputToShopRPCSpecification(&data.Input.Specification),
		// Title:         data.Input.Title,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// RemoveProduct ...
func (_ *Resolver) RemoveProduct(ctx context.Context, data RemoveProductRequest) (*SuccessResolver, error) {
	_, err := shop.RemoveProduct(ctx, &shopRPC.RemoveProductRequest{
		CompanyID: NullToString(data.Company_id),
		ProductID: data.Product_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// HideProduct ...
func (_ *Resolver) HideProduct(ctx context.Context, data HideProductRequest) (*SuccessResolver, error) {
	_, err := shop.HideProduct(ctx, &shopRPC.HideProductRequest{
		CompanyID: NullToString(data.Company_id),
		ProductID: data.Product_id,
		Value:     data.Value,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// GetMyShops ...
func (_ *Resolver) GetMyShops(ctx context.Context, data GetMyShopsRequest) (*[]ShopResolver, error) {
	shops, err := shop.GetMyShops(ctx, &shopRPC.ID{
		ID: NullToString(data.CompanyID),
	})
	if err != nil {
		return nil, err
	}

	sh := make([]ShopResolver, len(shops.GetShops()))

	for i := range shops.GetShops() {
		sh[i] = ShopResolver{
			R: shopRPCToShop(shops.GetShops()[i]),
		}
	}

	return &sh, nil
}

// SearchProduct ...
func (_ *Resolver) SearchProduct(ctx context.Context, data SearchProductRequest) (*[]ProductResolver, error) {
	var after, first string

	if data.Pagination.After != nil {
		after = *data.Pagination.After
	}
	if data.Pagination.First != nil {
		first = strconv.Itoa(int(*data.Pagination.First))
	}

	res, err := shop.FindProducts(ctx, &shopRPC.FindProductsRequest{
		ShopID:       NullToString(data.Shop_id),
		SearchFilter: searchProductFilterInputToSearchFilterRPC(data.Filter),
		After:        after,
		First:        first,
	})
	if err != nil {
		return nil, err
	}

	products := make([]ProductResolver, 0, len(res.GetProducts()))

	for i := range res.GetProducts() {
		products = append(products, ProductResolver{
			R: shopRPCProductToProduct(res.GetProducts()[i]),
		})
	}

	return &products, nil
}

// ---

func shopRPCToShop(data *shopRPC.Shop) *Shop {
	if data == nil {
		return nil
	}

	sh := Shop{
		ID:           data.GetID(),
		Description:  data.GetDescription(),
		ProductTypes: data.GetProductsType(),
		SellerType:   data.GetSellerType().String(),
		Showcase:     data.GetShowcase(),
		Title:        data.GetTitle(),
		Logo:         data.GetLogo(),
		Cover:        data.GetCover(),
	}

	if cat := shopRPCCategoryToCategory(data.GetCategory()); cat != nil {
		sh.Category = *cat
	}

	return &sh
}

func shopRPCCategoryToCategory(data *shopRPC.Category) *Category {
	if data == nil {
		return nil
	}

	return &Category{
		Main: data.GetMain(),
		// TODO:
		// Sub_Category: data.GetSub(),
	}
}

func shopRPCProductToProduct(data *shopRPC.Product) *Product {
	if data == nil {
		return nil
	}

	pr := Product{
		Brand:       data.GetBrand(),
		ID:          data.GetID(),
		In_stock:    data.GetInStock(),
		Is_used:     data.GetIsUsed(),
		ProductType: data.GetProductType(),
		Quantity:    int32(data.GetQuantity()),
		Sku:         data.GetSKU(),
		Title:       data.GetTitle(),
		Images:      make([]File, 0, len(data.GetImages())),
	}

	if cat := shopRPCCategoryToCategory(data.GetCategory()); cat != nil {
		pr.Category = *cat
	}

	if dis := shopRPCDiscountToDiscount(data.GetDiscount()); dis != nil {
		pr.Discount = *dis
	}

	if price := shopRPCPriceToPrice(data.GetPrice()); price != nil {
		pr.Price = *price
	}

	if spec := shopRPCSpecificationToSpecification(data.GetSpecification()); spec != nil {
		pr.Specification = *spec
	}

	for i := range data.GetImages() {
		if f := shopRPCFileToFile(data.GetImages()[i]); f != nil {
			pr.Images[i] = *f
		}
	}

	return &pr
}

func shopRPCDiscountToDiscount(data *shopRPC.Discount) *Discount {
	if data == nil {
		return nil
	}

	return &Discount{
		AmountOfProducts: int32(data.GetAmountOfProducts()),
		DiscountType:     data.GetDiscountType(),
		DiscountValue:    data.GetDiscountValue(),
		EndDate:          data.GetEndDate(),
		StartDate:        data.GetStartDate(),
	}
}

func shopRPCPriceToPrice(data *shopRPC.Price) *Price {
	if data == nil {
		return nil
	}

	return &Price{
		Amount:   float64(data.GetAmount()) / 100.0,
		Currency: data.GetCurrency(),
	}
}

func shopRPCSpecificationToSpecification(data *shopRPC.Specification) *Specification {
	if data == nil {
		return nil
	}

	spec := Specification{
		Color:      data.GetColor(),
		Material:   data.GetMaterial(),
		Size:       data.GetSize(),
		Variations: make([]Variation, 0, len(data.GetVariations())),
	}

	for i := range data.GetVariations() {
		if variation := shopRPCVariationToVariation(data.GetVariations()[i]); variation != nil {
			spec.Variations = append(spec.Variations, *variation)
		}
	}

	return &spec
}

func shopRPCVariationToVariation(data *shopRPC.Variation) *Variation {
	if data == nil {
		return nil
	}

	v := Variation{
		In_stock: data.GetInStock(),
		Quantity: int32(data.GetQuantity()),
		Sku:      data.GetSKU(),
	}

	if price := shopRPCPriceToPrice(data.GetPrice()); price != nil {
		v.Price = *price
	}

	return &v
}

func shopRPCOrderToOrder(ctx context.Context, data *shopRPC.Order) *Order {
	if data == nil {
		return nil
	}

	or := Order{
		Created_at:    data.GetCreatedAt(),
		Delivery_time: data.GetDeliverTime(),
		ID:            data.GetID(),
		Order_status:  data.GetOrderStatus().String(),
		Quantity:      int32(data.GetQuantity()),
	}

	if ad := shopRPCAddressToAddress(data.GetAddress()); ad != nil {
		or.Address = *ad
	}

	if price := shopRPCPriceToPrice(data.GetPrice()); price != nil {
		or.Price = *price
	}

	if data.GetProductID() != "" {
		prod, err := shop.GetProductsWithDeleted(ctx, &shopRPC.IDs{
			ID: []string{data.GetProductID()},
		})
		if err != nil {
			log.Println(err)
			return nil
		}

		if len(prod.GetProducts()) != 0 {
			if p := shopRPCProductToProduct(prod.GetProducts()[0]); p != nil {
				or.Product = *p
			}
		}
	}

	return &or
}

func shopRPCAddressToAddress(data *shopRPC.Address) *ShopAddress {
	if data == nil {
		return nil
	}

	return &ShopAddress{
		Address:   data.GetAddress(),
		CityID:    data.GetCityID(),
		Comments:  data.GetComments(),
		Firstname: data.GetFirstName(),
		Lastname:  data.GetLastName(),
		Phone:     data.GetPhone(),
		Zip_code:  data.GetZIPCode(),
		// data.GetApartment(),
	}
}

func shopRPCFileToFile(data *shopRPC.File) *File {
	if data == nil {
		return nil
	}

	return &File{
		Address:   data.GetURL(),
		ID:        data.GetID(),
		Mime_type: data.GetMimeType(),
		Name:      data.GetName(),
	}
}

// ---

func specificationInputToShopRPCSpecification(data *SpecificationInput) *shopRPC.Specification {
	if data == nil {
		return nil
	}

	spec := shopRPC.Specification{
		Color:    data.Color,
		Material: data.Material,
		Size:     data.Size,
	}

	if data.Variations != nil {
		spec.Variations = make([]*shopRPC.Variation, len(*data.Variations))

		for i := range *data.Variations {
			spec.Variations[i] = variationInputToShopRPCVariation(&(*data.Variations)[i])
		}
	}

	return &spec
}

func variationInputToShopRPCVariation(data *VariationInput) *shopRPC.Variation {
	if data == nil {
		return nil
	}

	return &shopRPC.Variation{
		InStock:  data.In_stock,
		Price:    priceInputToShopRPCPrice(&data.Price),
		Quantity: uint32(data.Quantity),
		SKU:      data.Sku,
	}
}

func priceInputToShopRPCPrice(data *PriceInput) *shopRPC.Price {
	if data == nil {
		return nil
	}

	return &shopRPC.Price{
		Amount:   uint32(int(data.Amount) * 100),
		Currency: data.Currency,
	}
}

func categoryInputToShopRPCCategory(data *CategoryInput) *shopRPC.Category {
	if data == nil {
		return nil
	}

	return &shopRPC.Category{
		Main: data.Main,
		// Sub:  data.Sub_category,
	}
}

func discountInputToShopRPCDiscount(data *DiscountInput) *shopRPC.Discount {
	if data == nil {
		return nil
	}

	return &shopRPC.Discount{
		AmountOfProducts: uint32(data.AmountOfProducts),
		DiscountType:     data.DiscountType,
		DiscountValue:    data.DiscountValue,
		EndDate:          data.EndDate,
		StartDate:        data.StartDate,
	}
}

func addressInputToShopRPCAddress(data *ShopAddressInput) *shopRPC.Address {
	if data == nil {
		return nil
	}

	return &shopRPC.Address{
		Address:   data.Address,
		Apartment: data.Apartment,
		CityID:    data.CityID,
		Comments:  NullToString(data.Comments),
		FirstName: data.Firstname,
		LastName:  data.Lastname,
		Phone:     data.Phone,
		ZIPCode:   data.Zip_code,
	}
}

func searchProductFilterInputToSearchFilterRPC(data *SearchProductInput) *shopRPC.SearchFilter {
	if data == nil {
		return nil
	}

	sf := shopRPC.SearchFilter{
		Keyword:  NullToString(data.Keyword),
		PriceMax: Nullint32ToUint32(data.PriceMax) * 100,
		PriceMin: Nullint32ToUint32(data.PriceMin) * 100,
	}

	if data.CategoryMain != nil {
		sf.CategoryMain = *data.CategoryMain
	}

	if data.CategorySub != nil {
		sf.CategorySub = *data.CategorySub
	}

	if data.IsInStock == nil {
		sf.IsInStockNull = true
	} else {
		sf.InStock = *data.IsInStock
	}

	if data.IsUsed == nil {
		sf.IsIsUsedNull = true
	} else {
		sf.IsUsed = *data.IsUsed
	}

	return &sf
}
