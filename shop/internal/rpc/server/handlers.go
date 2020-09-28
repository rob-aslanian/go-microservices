package serverRPC

import (
	"context"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/shopRPC"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/file"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/price"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/shop"
)

// CreateShop ...
func (s Server) CreateShop(ctx context.Context, data *shopRPC.CreateShopRequest) (*shopRPC.ID, error) {
	id, err := s.service.CreateShop(
		ctx,
		createShopRequestRPCToShop(data),
	)
	if err != nil {
		return nil, err
	}

	return &shopRPC.ID{
		ID: id,
	}, nil
}

// AddProduct ...
func (s Server) AddProduct(ctx context.Context, data *shopRPC.AddProductRequest) (*shopRPC.ID, error) {
	id, err := s.service.AddProduct(
		ctx,
		data.GetCompanyID(),
		shopRPCAddProductRequestToShopProduct(data),
	)
	if err != nil {
		return nil, err
	}

	return &shopRPC.ID{
		ID: id,
	}, nil
}

// GetProduct ...
func (s Server) GetProduct(ctx context.Context, data *shopRPC.ID) (*shopRPC.Product, error) {
	pr, err := s.service.GetProduct(
		ctx,
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}

	return shopProductToShopRPCProduct(pr), nil
}

// GetShop ...
func (s Server) GetShop(ctx context.Context, data *shopRPC.ID) (*shopRPC.Shop, error) {
	sh, err := s.service.GetShop(
		ctx,
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}

	return shopToShopRPCShop(sh), nil
}

// ChangeShowcase ...
func (s Server) ChangeShowcase(ctx context.Context, data *shopRPC.ChangeShowcaseRequest) (*shopRPC.Empty, error) {
	err := s.service.ChangeShowcase(ctx, data.GetCompanyID(), data.GetShopID(), data.GetShowcase())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// MakeOrder ...
func (s Server) MakeOrder(ctx context.Context, data *shopRPC.MakeOrderRequest) (*shopRPC.Empty, error) {
	err := s.service.MakeOrder(ctx,
		data.GetProductIDs(),
		shopRPCAddressToShopAddress(data.GetAddress()),
	)
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// AddToWishlist ...
func (s Server) AddToWishlist(ctx context.Context, data *shopRPC.ID) (*shopRPC.Empty, error) {
	err := s.service.AddToWishlist(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// GetMyWishlist ...
func (s Server) GetMyWishlist(ctx context.Context, data *shopRPC.IDWithPagination) (*shopRPC.Products, error) {
	var first uint
	var after uint

	if data.GetFirst() != "" {
		firstInt, _ := strconv.Atoi(data.GetFirst())
		first = uint(firstInt)
	}

	if data.GetAfter() != "" {
		afterInt, _ := strconv.Atoi(data.GetAfter())
		after = uint(afterInt)
	}

	products, err := s.service.GetMyWishlist(ctx, first, after)
	if err != nil {
		return nil, err
	}

	prs := make([]*shopRPC.Product, len(products))

	for i := range products {
		prs[i] = shopProductToShopRPCProduct(products[i])
	}

	return &shopRPC.Products{
		Products: prs,
	}, nil
}

// GetOrdersForBuyer ...
func (s Server) GetOrdersForBuyer(ctx context.Context, data *shopRPC.IDWithPagination) (*shopRPC.Orders, error) {
	var first uint
	var after uint

	if data.GetFirst() != "" {
		firstInt, _ := strconv.Atoi(data.GetFirst())
		first = uint(firstInt)
	}

	if data.GetAfter() != "" {
		afterInt, _ := strconv.Atoi(data.GetAfter())
		after = uint(afterInt)
	}

	products, err := s.service.GetOrdersForBuyer(ctx, first, after)
	if err != nil {
		return nil, err
	}

	orders := make([]*shopRPC.Order, len(products))

	for i := range products {
		orders[i] = shopOrderToOrderRPC(products[i])
	}

	return &shopRPC.Orders{
		Orders: orders,
	}, nil
}

// GetOrdersForSeller ...
func (s Server) GetOrdersForSeller(ctx context.Context, data *shopRPC.GetOrdersForSellerRequest) (*shopRPC.Orders, error) {
	var first uint
	var after uint

	if data.GetFirst() != "" {
		firstInt, _ := strconv.Atoi(data.GetFirst())
		first = uint(firstInt)
	}

	if data.GetAfter() != "" {
		afterInt, _ := strconv.Atoi(data.GetAfter())
		after = uint(afterInt)
	}

	products, err := s.service.GetOrdersForSeller(ctx, data.GetCompanyID(), data.GetShopID(), first, after)
	if err != nil {
		return nil, err
	}

	orders := make([]*shopRPC.Order, len(products))

	for i := range products {
		orders[i] = shopOrderToOrderRPC(products[i])
	}

	return &shopRPC.Orders{
		Orders: orders,
	}, nil
}

// ChangeOrderStatus ...
func (s Server) ChangeOrderStatus(ctx context.Context, data *shopRPC.ChangeOrderStatusRequest) (*shopRPC.Empty, error) {
	err := s.service.ChangeOrderStatus(ctx, data.GetCompanyID(), data.GetOrderID(), data.GetOrderStatus().String())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// GetProductsWithDeleted ...
func (s Server) GetProductsWithDeleted(ctx context.Context, data *shopRPC.IDs) (*shopRPC.Products, error) {
	products, err := s.service.GetProductsWithDeleted(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	prod := make([]*shopRPC.Product, len(products))

	for i := range products {
		prod[i] = shopProductToShopRPCProduct(products[i])
	}

	return &shopRPC.Products{
		Products: prod,
	}, nil
}

// ChangeLogo ...
func (s Server) ChangeLogo(ctx context.Context, data *shopRPC.File) (*shopRPC.Empty, error) {
	err := s.service.ChangeLogo(ctx, data.GetCompanyID(), data.GetTargetID(), data.GetURL())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// ChangeCover ...
func (s Server) ChangeCover(ctx context.Context, data *shopRPC.File) (*shopRPC.Empty, error) {
	err := s.service.ChangeCover(ctx, data.GetCompanyID(), data.GetTargetID(), data.GetURL())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// RemoveLogo ...
func (s Server) RemoveLogo(ctx context.Context, data *shopRPC.RemoveLogoRequest) (*shopRPC.Empty, error) {
	err := s.service.RemoveLogo(ctx, data.GetCompanyID(), data.GetShopID())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// RemoveCover ...
func (s Server) RemoveCover(ctx context.Context, data *shopRPC.RemoveLogoRequest) (*shopRPC.Empty, error) {
	err := s.service.RemoveCover(ctx, data.GetCompanyID(), data.GetShopID())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// ChangeImagesInProduct ...
func (s Server) ChangeImagesInProduct(ctx context.Context, data *shopRPC.ChangeImagesInProductRequest) (*shopRPC.Empty, error) {
	files := make([]*file.File, len(data.GetFiles()))

	for i := range data.GetFiles() {
		files[i] = shopRPCFileToFile(data.GetFiles()[i])
	}

	err := s.service.ChangeImagesInProduct(ctx, data.GetCompanyID(), data.GetProductID(), files)
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// ChangeProduct ...
func (s Server) ChangeProduct(ctx context.Context, data *shopRPC.Product) (*shopRPC.Empty, error) {
	err := s.service.ChangeProduct(ctx, data.GetCompanyID(), shopRPCProductToProduct(data))
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// RemoveProduct ...
func (s Server) RemoveProduct(ctx context.Context, data *shopRPC.RemoveProductRequest) (*shopRPC.Empty, error) {
	err := s.service.RemoveProduct(ctx, data.GetCompanyID(), data.GetProductID())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// HideProduct ...
func (s Server) HideProduct(ctx context.Context, data *shopRPC.HideProductRequest) (*shopRPC.Empty, error) {
	err := s.service.HideProduct(ctx, data.GetCompanyID(), data.GetProductID(), data.GetValue())
	if err != nil {
		return nil, err
	}

	return &shopRPC.Empty{}, nil
}

// GetMyShops ...
func (s Server) GetMyShops(ctx context.Context, data *shopRPC.ID) (*shopRPC.Shops, error) {
	shops, err := s.service.GetMyShops(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	shopsPRC := make([]*shopRPC.Shop, len(shops))

	for i := range shops {
		shopsPRC[i] = shopToShopRPCShop(shops[i])
	}

	return &shopRPC.Shops{
		Shops: shopsPRC,
	}, nil
}

// FindProducts ...
func (s Server) FindProducts(ctx context.Context, data *shopRPC.FindProductsRequest) (*shopRPC.Products, error) {
	var first uint
	var after uint

	if data.GetFirst() != "" {
		firstInt, _ := strconv.Atoi(data.GetFirst())
		first = uint(firstInt)
	}

	if data.GetAfter() != "" {
		afterInt, _ := strconv.Atoi(data.GetAfter())
		after = uint(afterInt)
	}

	products, err := s.service.FindProducts(ctx, data.GetShopID(), shopRPCSearchFilterToSearchFilter(data.GetSearchFilter()), first, after)
	if err != nil {
		return nil, err
	}

	prod := make([]*shopRPC.Product, len(products))

	for i := range products {
		prod[i] = shopProductToShopRPCProduct(products[i])
	}

	return &shopRPC.Products{
		Products: prod,
	}, nil
}

// ---

func createShopRequestRPCToShop(data *shopRPC.CreateShopRequest) *shop.Shop {
	if data == nil {
		return nil
	}

	productTypes := make([]string, len(data.GetProductsType()))

	for i := range data.GetProductsType() {
		productTypes[i] = data.GetProductsType()[i].String()
	}

	sh := shop.Shop{
		Title: data.GetTitle(),
		Category: shop.Category{
			Main: data.GetCategory(),
		},
		SellerType:   data.GetSellerType().String(),
		Description:  data.GetDescription(),
		ProductsType: productTypes,
	}

	sh.SetCompanyID(data.GetCompanyID())

	return &sh
}

func shopRPCAddProductRequestToShopProduct(data *shopRPC.AddProductRequest) *shop.Product {
	if data == nil {
		return nil
	}

	sh := shop.Product{
		Description: data.GetDescription(),
		InStock:     data.GetInStock(),
		IsUsed:      data.GetIsUsed(),
		ProductType: data.GetProductType().String(),
		Quantity:    data.GetQuantity(),
		SKU:         data.GetSKU(),
		Title:       data.GetTitle(),
	}

	sh.SetShopID(data.GetShopID())

	if data.GetBrand() != "" {
		brand := data.GetBrand()
		sh.Brand = &brand
	}

	if cat := shopRPCCategory(data.GetCategory()); cat != nil {
		sh.Category = *cat
	}

	if dis := shopRPCDiscountToShopDiscount(data.GetDiscount()); dis != nil {
		sh.Discount = *dis
	}

	if pr := shopRPCPriceToPrice(data.GetPrice()); pr != nil {
		sh.Price = *pr
	}

	if sp := shopRPCSpecificationToShopSpecification(data.GetSpecification()); sp != nil {
		sh.Specification = *sp
	}

	return &sh
}

func shopRPCCategory(data *shopRPC.Category) *shop.Category {
	if data == nil {
		return nil
	}

	return &shop.Category{
		Main: data.GetMain(),
		Sub:  data.GetSub(),
	}
}

func shopRPCDiscountToShopDiscount(data *shopRPC.Discount) *shop.Discount {
	if data == nil {
		return nil
	}

	return &shop.Discount{
		AmountOfProducts: data.GetAmountOfProducts(),
		DicountValue:     data.GetDiscountValue(),
		DiscountType:     data.GetDiscountType(),
		EndDate:          stringToTime(data.GetEndDate()),
		StartDate:        stringToTime(data.GetStartDate()),
	}
}

func shopRPCPriceToPrice(data *shopRPC.Price) *price.Price {
	if data == nil {
		return nil
	}

	return &price.Price{
		Amount:   data.GetAmount(),
		Currency: data.GetCurrency(),
	}
}

func shopRPCSpecificationToShopSpecification(data *shopRPC.Specification) *shop.Specification {
	if data == nil {
		return nil
	}

	sp := shop.Specification{
		Color:      data.GetColor(),
		Material:   data.GetMaterial(),
		Size:       data.GetSize(),
		Variations: make([]shop.Variation, len(data.GetVariations())),
	}

	for i := range data.GetVariations() {
		if v := shopRPCVariationsToShopVariation(data.GetVariations()[i]); v != nil {
			sp.Variations[i] = *v
		}
	}

	return &sp
}

func shopRPCVariationsToShopVariation(data *shopRPC.Variation) *shop.Variation {
	if data == nil {
		return nil
	}

	v := shop.Variation{
		InStock:  data.GetInStock(),
		Quantity: data.GetQuantity(),
		SKU:      data.GetSKU(),
	}

	if pr := shopRPCPriceToPrice(data.GetPrice()); pr != nil {
		v.Price = *pr
	}

	return &v
}

func shopRPCAddressToShopAddress(data *shopRPC.Address) *shop.Address {
	if data == nil {
		return nil
	}

	a := shop.Address{
		Address:   data.GetAddress(),
		Apartment: data.GetApartment(),
		Comments:  data.GetComments(),
		FirstName: data.GetFirstName(),
		LastName:  data.GetLastName(),
		Phone:     data.GetPhone(),
		ZIPCode:   data.GetZIPCode(),
	}

	if id, err := strconv.Atoi(data.GetCityID()); err == nil {
		a.CityID = uint(id)
	}

	return &a
}

func shopRPCFileToFile(data *shopRPC.File) *file.File {
	if data == nil {
		return nil
	}

	f := file.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		Position: data.GetPosition(),
		URL:      data.GetURL(),
	}

	f.SetID(data.GetID())

	return &f
}

func shopRPCProductToProduct(data *shopRPC.Product) *shop.Product {
	if data == nil {
		return nil
	}

	pr := shop.Product{
		InStock:  data.GetInStock(),
		IsUsed:   data.GetIsUsed(),
		Quantity: data.GetQuantity(),
		SKU:      data.GetSKU(),
		// ProductType
		// IsDeleted: data.GetIsDeleted(),
		// Discount
	}

	pr.SetID(data.GetID())

	if price := shopRPCPriceToPrice(data.GetPrice()); price != nil {
		pr.Price = *price
	}

	if cat := shopRPCCategory(data.GetCategory()); cat != nil {
		pr.Category = *cat
	}

	if spec := shopRPCSpecificationToShopSpecification(data.GetSpecification()); spec != nil {
		pr.Specification = *spec
	}

	return &pr
}

func shopRPCSearchFilterToSearchFilter(data *shopRPC.SearchFilter) *shop.SearchFilter {
	if data == nil {
		return nil
	}

	sf := shop.SearchFilter{
		Category:    data.GetCategoryMain(),
		Subcategory: data.GetCategorySub(),
		PriceMax:    data.GetPriceMax(),
		PriceMin:    data.GetPriceMin(),
		Keyword:     data.GetKeyword(),
	}

	if !data.GetIsIsUsedNull() {
		sf.IsUsed = &data.IsUsed
	}

	if !data.GetIsInStockNull() {
		sf.InStock = &data.InStock
	}

	return &sf
}

// ---

func shopProductToShopRPCProduct(data *shop.Product) *shopRPC.Product {
	if data == nil {
		return nil
	}

	sh := shopRPC.Product{
		Category:      shopCategoryToShopRPCCategory(&data.Category),
		CompanyID:     data.GetCompanyID(),
		Description:   data.Description,
		Discount:      shopDiscountToShopRPCDiscount(&data.Discount),
		ID:            data.GetID(),
		InStock:       data.InStock,
		IsUsed:        data.IsUsed,
		Price:         priceToShopRPCPrice(&data.Price),
		ProductType:   data.ProductType,
		Quantity:      data.Quantity,
		SKU:           data.SKU,
		Title:         data.Title,
		Specification: shopSpecificationToShopRPCSpecification(&data.Specification),
		Images:        make([]*shopRPC.File, 0, len(data.Images)),
	}

	if data.Brand != nil {
		sh.Brand = *data.Brand
	}

	for i := range data.Images {
		if f := fileToShopRPCFile(&data.Images[i]); f != nil {
			sh.Images = append(sh.Images, f)
		}
	}

	return &sh
}

func shopCategoryToShopRPCCategory(data *shop.Category) *shopRPC.Category {
	if data == nil {
		return nil
	}

	return &shopRPC.Category{
		Main: data.Main,
		Sub:  data.Sub,
	}
}

func shopDiscountToShopRPCDiscount(data *shop.Discount) *shopRPC.Discount {
	if data == nil {
		return nil
	}

	return &shopRPC.Discount{
		AmountOfProducts: data.AmountOfProducts,
		DiscountType:     data.DiscountType,
		DiscountValue:    data.DicountValue,
		EndDate:          timeToString(data.EndDate),
		StartDate:        timeToString(data.StartDate),
	}
}

func priceToShopRPCPrice(data *price.Price) *shopRPC.Price {
	if data == nil {
		return nil
	}

	return &shopRPC.Price{
		Amount:   data.Amount,
		Currency: data.Currency,
	}
}

func shopSpecificationToShopRPCSpecification(data *shop.Specification) *shopRPC.Specification {
	if data == nil {
		return nil
	}

	sp := shopRPC.Specification{
		Color:      data.Color,
		Material:   data.Material,
		Size:       data.Size,
		Variations: make([]*shopRPC.Variation, len(data.Variations)),
	}

	for i := range data.Variations {
		sp.Variations[i] = shopSpecificationToShopRPCVariations(&data.Variations[i])
	}

	return &sp
}

func shopSpecificationToShopRPCVariations(data *shop.Variation) *shopRPC.Variation {
	if data == nil {
		return nil
	}

	return &shopRPC.Variation{
		InStock:  data.InStock,
		Price:    priceToShopRPCPrice(&data.Price),
		Quantity: data.Quantity,
		SKU:      data.SKU,
	}
}

func shopToShopRPCShop(data *shop.Shop) *shopRPC.Shop {
	if data == nil {
		return nil
	}

	sh := shopRPC.Shop{
		Category:     shopCategoryToShopRPCCategory(&data.Category),
		CompanyID:    data.GetCompanyID(),
		Description:  data.Description,
		ID:           data.GetID(),
		ProductsType: data.ProductsType,
		SellerType:   shopRPC.SellerType(shopRPC.SellerType_value[data.SellerType]),
		Title:        data.Title,
		UserID:       data.GetUserID(),
	}

	if data.Logo != nil {
		sh.Logo = *data.Logo
	}

	if data.Cover != nil {
		sh.Cover = *data.Cover
	}

	if data.Showcase != nil {
		sh.Showcase = *data.Showcase
	}

	return &sh
}

func shopOrderToOrderRPC(data *shop.Order) *shopRPC.Order {
	if data == nil {
		return nil
	}

	order := shopRPC.Order{
		ID:          data.GetID(),
		Address:     shopAddressToShopRPCAddress(data.Address),
		CompanyID:   data.GetCompanyID(),
		CreatedAt:   timeToString(data.CreatedAt),
		OrderStatus: shopRPC.OrderStatus(shopRPC.OrderStatus_value[data.OrderStatus]),
		Price:       priceToShopRPCPrice(&data.Price),
		ProductID:   data.GetProductID(),
		Quantity:    data.Quantity,
		UserID:      data.GetUserID(),
	}

	if data.DeliveryTime != nil {
		order.DeliverTime = timeToString(*data.DeliveryTime)
	}

	return &order
}

func shopAddressToShopRPCAddress(data *shop.Address) *shopRPC.Address {
	if data == nil {
		return nil
	}

	a := shopRPC.Address{
		Address:   data.Address,
		Apartment: data.Apartment,
		Comments:  data.Comments,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Phone:     data.Phone,
		ZIPCode:   data.ZIPCode,
		CityID:    strconv.Itoa(int(data.CityID)),
	}

	return &a
}

func fileToShopRPCFile(data *file.File) *shopRPC.File {
	if data == nil {
		return nil
	}

	f := shopRPC.File{
		ID:       data.GetID(),
		MimeType: data.MimeType,
		Name:     data.Name,
		Position: data.Position,
		URL:      data.URL,
	}

	return &f
}

// ---

func timeToString(s time.Time) string {
	return s.Format(time.RFC3339)
}

func stringToTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
