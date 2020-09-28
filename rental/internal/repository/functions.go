package servicesrepo

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/rental"
	"go.mongodb.org/mongo-driver/bson"
)

// docker-compose exec mongo_db mongo --username developer --authenticationDatabase admin -p

// AddHouseRentalAppartament ...
func (r *Repository) AddHouseRentalAppartament(ctx context.Context, data rental.Appartament) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil

}

// AddRealEstateStorageRooms ...
func (r *Repository) AddRealEstateStorageRooms(ctx context.Context, data rental.StorageRooms) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateBuildings ...
func (r *Repository) AddRealEstateBuildings(ctx context.Context, data rental.Buildings) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil

}

// AddRealEstateCommercial ...
func (r *Repository) AddRealEstateCommercial(ctx context.Context, data rental.Commercial) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateGarage ...
func (r *Repository) AddRealEstateGarage(ctx context.Context, data rental.Garage) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateHotelRooms ...
func (r *Repository) AddRealEstateHotelRooms(ctx context.Context, data rental.HotelRooms) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateLand ...
func (r *Repository) AddRealEstateLand(ctx context.Context, data rental.Land) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateOffice ...
func (r *Repository) AddRealEstateOffice(ctx context.Context, data rental.Office) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateRuralFarm ...
func (r *Repository) AddRealEstateRuralFarm(ctx context.Context, data rental.RuralFarm) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateRenovation ...
func (r *Repository) AddRealEstateRenovation(ctx context.Context, data rental.Renovation) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateMaterials ...
func (r *Repository) AddRealEstateMaterials(ctx context.Context, data rental.Materials) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AddRealEstateMove ...
func (r *Repository) AddRealEstateMove(ctx context.Context, data rental.Move) error {
	_, err := r.housesCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// GetRealEstates ...
func (r *Repository) GetRealEstates(ctx context.Context, dealType rental.DealType, first int, after int) (rental.GetRental, error) {

	cursor, err := r.housesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"post_status": "active",
					"deal_type":   dealType,
				},
			},
			{
				"$sort": bson.M{
					"is_urgent":  -1,
					"created_at": -1,
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"estates": bson.M{
						"$push": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$size": "$estates"},
							0,
						},
					},
					"materials": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.deal_type", "materials"}},
							}},
							0,
						},
					},
					"move": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.deal_type", "move"}},
							}},
							0,
						},
					},
					"renovation": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.deal_type", "renovation"}},
							}},
							0,
						},
					},
					"land": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "land"}},
							}},
							0,
						},
					},
					"buildings": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "buildings"}},
							}},
							0,
						},
					},
					"commercial_properties": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "commercial_properties"}},
							}},
							0,
						},
					},
					"offices": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "offices"}},
							}},
							0,
						},
					},
					"garages": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "garages"}},
							}},
							0,
						},
					},
					"appartments": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "appartments"}},
							}},
							0,
						},
					},
					"houses": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "houses"}},
							}},
							0,
						},
					},
					"homes": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "homes"}},
							}},
							0,
						},
					},
					"new_homes": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "new_homes"}},
							}},
							0,
						},
					},
					"summer_cottage": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "summer_cottage"}},
							}},
							0,
						},
					},
					"rural_farm": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "rural_farm"}},
							}},
							0,
						},
					},
					"storage_rooms": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$estates"},
							bson.M{"$filter": bson.M{
								"input": "$estates",
								"as":    "e",
								"cond":  bson.M{"$eq": []interface{}{"$$e.property_type", "storage_rooms"}},
							}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"_id": 0,
					"storage_rooms": bson.M{
						"$slice": []interface{}{
							"$storage_rooms",
							after,
							first,
						},
					},
					"new_homes": bson.M{
						"$slice": []interface{}{
							"$new_homes",
							after,
							first,
						},
					},
					"homes": bson.M{
						"$slice": []interface{}{
							"$homes",
							after,
							first,
						},
					},
					"houses": bson.M{
						"$slice": []interface{}{
							"$houses",
							after,
							first,
						},
					},
					"appartments": bson.M{
						"$slice": []interface{}{
							"$appartments",
							after,
							first,
						},
					},
					"garages": bson.M{
						"$slice": []interface{}{
							"$garages",
							after,
							first,
						},
					},
					"offices": bson.M{
						"$slice": []interface{}{
							"$offices",
							after,
							first,
						},
					},
					"commercial_properties": bson.M{
						"$slice": []interface{}{
							"$commercial_properties",
							after,
							first,
						},
					},
					"buildings": bson.M{
						"$slice": []interface{}{
							"$buildings",
							after,
							first,
						},
					},
					"land": bson.M{
						"$slice": []interface{}{
							"$land",
							after,
							first,
						},
					},
					"rural_farm": bson.M{
						"$slice": []interface{}{
							"$rural_farm",
							after,
							first,
						},
					},
					"materials": bson.M{
						"$slice": []interface{}{
							"$materials",
							after,
							first,
						},
					},
					"renovation": bson.M{
						"$slice": []interface{}{
							"$renovation",
							after,
							first,
						},
					},
					"move": bson.M{
						"$slice": []interface{}{
							"$move",
							after,
							first,
						},
					},
					"amount": 1,
				},
			},
		},
	)

	if err != nil {
		return rental.GetRental{}, err
	}
	defer cursor.Close(ctx)

	res := rental.GetRental{}
	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return rental.GetRental{}, err
		}
	}

	return res, nil
}
