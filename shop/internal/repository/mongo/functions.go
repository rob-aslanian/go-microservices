package repository

import (
	"context"
	"errors"

	"gitlab.lan/Rightnao-site/microservices/shop/internal/file"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/shop"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateShop ...
func (r Repository) CreateShop(ctx context.Context, sh *shop.Shop) error {
	_, err := r.shopsCollection.InsertOne(ctx, sh)
	if err != nil {
		return err
	}

	return nil
}

// AddProduct ...
func (r Repository) AddProduct(ctx context.Context, pr *shop.Product) error {
	_, err := r.productsCollection.InsertOne(ctx, pr)
	if err != nil {
		return err
	}

	return nil
}

// GetProduct ...
func (r Repository) GetProduct(ctx context.Context, id string, withDeleted bool) (*shop.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	filter := bson.M{
		"_id": objID,
	}

	if !withDeleted {
		filter["is_deleted"] = false
		filter["is_hidden"] = false
	}

	res := r.productsCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	pr := new(shop.Product)
	err = res.Decode(pr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return pr, nil
}

// GetProducts ...
func (r Repository) GetProducts(ctx context.Context, ids []string, first, after uint) ([]*shop.Product, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	opt := options.Find()
	opt.SetSkip(int64(after))
	opt.SetLimit(int64(first))

	cursor, err := r.productsCollection.Find(ctx,
		bson.M{
			"_id": bson.M{
				"$in": objIDs,
			},
		},
		opt,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	products := make([]*shop.Product, 0)
	for cursor.Next(ctx) {
		pr := new(shop.Product)
		err = cursor.Decode(&pr)
		if err != nil {
			return nil, err
		}
		products = append(products, pr)
	}

	return products, nil
}

// GetShop ...
func (r Repository) GetShop(ctx context.Context, id string) (*shop.Shop, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	res := r.shopsCollection.FindOne(ctx, bson.M{
		"_id": objID,
	})

	v := new(shop.Shop)
	err = res.Decode(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// ChangeShowcase ...
func (r Repository) ChangeShowcase(ctx context.Context, shopID string, showcase string) error {
	objID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.shopsCollection.UpdateOne(ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"showcase": showcase,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// MakeOrder ...
func (r Repository) MakeOrder(ctx context.Context, orders []shop.Order) error {

	ords := make([]interface{}, len(orders))

	for i := range orders {
		ords[i] = orders[i]
	}

	_, err := r.ordersCollection.InsertMany(ctx,
		ords,
	)
	if err != nil {
		return err
	}

	return nil
}

// AddToWishlist ...
func (r Repository) AddToWishlist(ctx context.Context, userID string, productID string) error {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objProductID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.wishlistCollection.UpdateOne(ctx,
		bson.M{
			"_id": objUserID,
		},
		bson.M{
			"$addToSet": bson.M{
				"products": objProductID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// GetProductsIDFromWishlist ...
func (r Repository) GetProductsIDFromWishlist(ctx context.Context, userID string) ([]string, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.wishlistCollection.Find(ctx, bson.M{
		"_id": objUserID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	products := struct {
		ProductIDs []primitive.ObjectID `bson:"products"`
	}{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&products)
		if err != nil {
			return nil, err
		}
	}

	ids := make([]string, len(products.ProductIDs))

	for i := range products.ProductIDs {
		ids = append(ids, products.ProductIDs[i].Hex())
	}

	return ids, nil
}

// GetOrdersIDForBuyer ...
func (r Repository) GetOrdersIDForBuyer(ctx context.Context, userID string) ([]string, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.ordersCollection.Find(ctx,
		bson.M{
			"user_id": objUserID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ids := make([]string, 0)

	for cursor.Next(ctx) {
		id := struct {
			ID primitive.ObjectID `bson:"_id"`
		}{}
		err = cursor.Decode(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id.ID.Hex())
	}

	return ids, nil
}

// GetOrders ...
func (r Repository) GetOrders(ctx context.Context, ids []string, first, after uint) ([]*shop.Order, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	opt := options.Find()
	opt.SetSkip(int64(after))
	opt.SetLimit(int64(first))

	cursor, err := r.ordersCollection.Find(ctx,
		bson.M{
			"_id": bson.M{
				"$in": objIDs,
			},
		},
		opt,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	orders := make([]*shop.Order, 0)
	for cursor.Next(ctx) {
		order := new(shop.Order)
		err = cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrdersIDForSeller ...
func (r Repository) GetOrdersIDForSeller(ctx context.Context, shopID string) ([]string, error) {
	objShopID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.ordersCollection.Find(ctx,
		bson.M{
			"shop_id": objShopID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	orderIDs := make([]string, 0)

	for cursor.Next(ctx) {
		id := struct {
			ID primitive.ObjectID `bson:"_id"`
		}{}
		err = cursor.Decode(&id)
		if err != nil {
			return nil, err
		}
		orderIDs = append(orderIDs, id.ID.Hex())
	}

	return orderIDs, nil
}

// ChangeOrderStatus ...
func (r Repository) ChangeOrderStatus(ctx context.Context, orderID string, orderStatus string) error {
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.ordersCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"status": orderStatus,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetProductsWithDeleted ...
func (r Repository) GetProductsWithDeleted(ctx context.Context, ids []string) ([]*shop.Product, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	cursor, err := r.productsCollection.Find(ctx,
		bson.M{
			"_id": bson.M{
				"$in": objIDs,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	products := make([]*shop.Product, 0)
	for cursor.Next(ctx) {
		product := new(shop.Product)
		err = cursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// ChangeLogo ...
func (r Repository) ChangeLogo(ctx context.Context, shopID, url string) error {
	objID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.shopsCollection.UpdateOne(ctx, bson.M{
		"_id": objID,
	},
		bson.M{
			"$set": bson.M{
				"logo": url,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCover ...
func (r Repository) ChangeCover(ctx context.Context, shopID, url string) error {
	objID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.shopsCollection.UpdateOne(ctx, bson.M{
		"_id": objID,
	},
		bson.M{
			"$set": bson.M{
				"cover": url,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveLogo ...
func (r Repository) RemoveLogo(ctx context.Context, shopID string) error {
	objID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.shopsCollection.UpdateOne(ctx, bson.M{
		"_id": objID,
	},
		bson.M{
			"$set": bson.M{
				"logo": nil,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveCover ...
func (r Repository) RemoveCover(ctx context.Context, shopID string) error {
	objID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.shopsCollection.UpdateOne(ctx, bson.M{
		"_id": objID,
	},
		bson.M{
			"$set": bson.M{
				"cover": nil,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeImagesInProduct ...
func (r Repository) ChangeImagesInProduct(ctx context.Context, productID string, images []*file.File) error {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.productsCollection.UpdateOne(ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"images": images,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeProduct ...
func (r Repository) ChangeProduct(ctx context.Context, productID string, prod *shop.Product) error {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.productsCollection.UpdateOne(ctx, bson.M{
		"_id": objID,
	},
		bson.M{
			"$set": bson.M{
				"in_stock":      prod.InStock,
				"is_used":       prod.IsUsed,
				"quantity":      prod.Quantity,
				"sku":           prod.SKU,
				"price":         prod.Price,
				"category":      prod.Category,
				"specification": prod.Specification,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveProduct ...
func (r Repository) RemoveProduct(ctx context.Context, productID string) error {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.productsCollection.UpdateOne(ctx, bson.M{
		"_id": objID,
	},
		bson.M{
			"$set": bson.M{
				"is_deleted": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAmountOfShops ...
func (r Repository) GetAmountOfShops(ctx context.Context, ownerID string) (uint8, error) {
	objID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return 0, errors.New(`wrong_id`)
	}

	cursor, err := r.shopsCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{
						"user_id": objID,
					},
					{
						"company_id": objID,
					},
				},
			},
		},
		{"$count": "amount"},
	})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	result := struct {
		Amount uint8 `bson:"amount"`
	}{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return 0, err
		}
	}

	return result.Amount, nil
}

// HideProduct ...
func (r Repository) HideProduct(ctx context.Context, productID string, value bool) error {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.productsCollection.UpdateOne(ctx, bson.M{
		"_id": objID,
	},
		bson.M{
			"$set": bson.M{
				"is_hidden": value,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetMyShops ...
func (r Repository) GetMyShops(ctx context.Context, ownerID string) ([]*shop.Shop, error) {
	objOwnerID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.shopsCollection.Find(ctx, bson.M{
		"$or": []bson.M{
			{
				"user_id": objOwnerID,
			},
			{
				"company_id": objOwnerID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	shops := make([]*shop.Shop, 0)

	for cursor.Next(ctx) {
		sh := new(shop.Shop)
		err = cursor.Decode(&sh)
		if err != nil {
			return nil, err
		}
		shops = append(shops, sh)
	}

	return shops, nil
}

// FindProducts ...
func (r Repository) FindProducts(ctx context.Context, shopID string, filter *shop.SearchFilter, first, after uint) ([]*shop.Product, error) {
	pipe := make([]bson.M, 0)

	match := bson.M{}

	if shopID != "" {
		objID, err := primitive.ObjectIDFromHex(shopID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}

		match["shop_id"] = objID
	}

	if filter != nil {
		// keyword
		if filter.Keyword != "" {
			match["$text"] = bson.M{
				"$search": filter.Keyword,
			}
		}
		// category
		if len(filter.Category) > 0 {
			match["category.main"] = bson.M{
				"$in": filter.Category,
			}
		}
		// price
		if filter.PriceMax != 0 && filter.PriceMin != 0 {
			match["$and"] = []bson.M{
				{
					"price.amount": bson.M{
						"$lte": filter.PriceMax,
					},
				},
				{
					"price.amount": bson.M{
						"$gte": filter.PriceMin,
					},
				}}
		} else {
			if filter.PriceMax != 0 {
				match["price.amount"] = bson.M{
					"$lte": filter.PriceMax,
				}
			}
			if filter.PriceMin != 0 {
				match["price.amount"] = bson.M{
					"$gte": filter.PriceMin,
				}
			}
		}
		// in stock
		if filter.InStock != nil {
			match["in_stock"] = *filter.InStock
		}
		// is used
		if filter.IsUsed != nil {
			match["is_used"] = *filter.IsUsed
		}
		// sort
		//

	}
	if len(match) > 0 {
		pipe = append(pipe, bson.M{"$match": match})
	}

	pipe = append(pipe, bson.M{
		"$skip": after,
	},
		bson.M{
			"$limit": first,
		},
	)

	cursor, err := r.productsCollection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	products := make([]*shop.Product, 0)
	for cursor.Next(ctx) {
		pr := new(shop.Product)
		err = cursor.Decode(&pr)
		if err != nil {
			return nil, err
		}
		products = append(products, pr)
	}

	return products, nil
}
