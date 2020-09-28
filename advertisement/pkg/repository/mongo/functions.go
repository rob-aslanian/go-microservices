package repository

import (
	"context"
	"errors"
	"log"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/advert"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/file"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SaveBanner ...
func (r Repository) SaveBanner(ctx context.Context, ad *advert.Banner) error {
	_, err := r.advertCollection.InsertOne(ctx, ad)
	if err != nil {
		return err
	}

	return nil
}

// ChangeBanner ...
func (r Repository) ChangeBanner(ctx context.Context, ad *advert.Banner) error {
	_, err := r.advertCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": ad.ID,
		},
		bson.M{
			"$set": ad,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetBannerByID ...
func (r Repository) GetBannerByID(ctx context.Context, bannerID string) (*advert.Banner, error) {
	objID, err := primitive.ObjectIDFromHex(bannerID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	result := r.advertCollection.FindOne(
		ctx,
		bson.M{
			"_id": objID,
		},
	)

	if result.Err() != nil {
		return nil, result.Err()
	}

	ad := advert.Banner{}

	err = result.Decode(&ad)
	if err != nil {
		return nil, err
	}

	return &ad, nil
}

// AddImageToGallery ...
func (r Repository) AddImageToGallery(ctx context.Context, id string, f *file.File) error {

	objCampaingID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("wrong_officeID")
	}

	_, dbErr := r.advertCollection.UpdateOne(
		ctx,
		bson.M{
			"adverts._id": objCampaingID,
		},
		bson.M{
			"$push": bson.M{
				"adverts.$.files": f,
			},
		},
	)

	if dbErr != nil {
		return dbErr
	}
	return nil
}

// GetImageURLByID ...
func (r Repository) GetImageURLByID(ctx context.Context, fileID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return "", errors.New(`wrong_id`)
	}

	result := r.galleryCollection.FindOne(
		ctx,
		bson.M{
			"_id": objID,
		},
	)

	if result.Err() != nil {
		return "", result.Err()
	}

	file := struct {
		URL string `bson:"url"`
	}{}

	err = result.Decode(&file)
	if err != nil {
		return "", err
	}

	return file.URL, nil
}

// GetGallery ...
func (r Repository) GetGallery(ctx context.Context, userID, companyID string, first, after uint32) ([]*file.File, uint32, error) {

	match := bson.M{}

	if companyID != "" {
		objCompanyID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return nil, 0, errors.New(`wrong_id`)
		}

		match["company_id"] = objCompanyID
	} else {
		objUserID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, 0, errors.New(`wrong_id`)
		}

		match["user_id"] = objUserID
		match["company_id"] = bson.M{
			"$exists": false,
		}
	}

	cursor, err := r.galleryCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$facet": bson.M{
					"files": []bson.M{
						{"$skip": after},
						{"$limit": first},
					},
					"totalCount": []bson.M{
						{"$count": "count"},
					},
				},
			},
			{
				"$project": bson.M{
					"files": 1,
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$totalCount"},
							bson.M{"$arrayElemAt": []interface{}{"$totalCount", 0}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"files":  1,
					"amount": "$amount.count",
				},
			},
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	result := struct {
		Files  []*file.File `bson:"files"`
		Amount uint32       `bson:"amount"`
	}{}

	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.Files, result.Amount, nil
}

// GetMyAdverts ...
func (r Repository) GetMyAdverts(ctx context.Context, userID, companyID string, first, after uint32) ([]*advert.Advert, uint32, error) {
	match := bson.M{
		"deleted": bson.M{
			"$ne": true,
		},
	}

	if companyID != "" {
		objCompanyID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return nil, 0, errors.New(`wrong_id`)
		}

		match["company_id"] = objCompanyID
	} else {
		objUserID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, 0, errors.New(`wrong_id`)
		}

		match["creator_id"] = objUserID
		match["company_id"] = bson.M{
			"$exists": false,
		}
	}

	cursor, err := r.advertCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$facet": bson.M{
					"ads": []bson.M{
						{"$skip": after},
						{"$limit": first},
					},
					"totalCount": []bson.M{
						{"$count": "count"},
					},
				},
			},
			{
				"$project": bson.M{
					"ads": 1,
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$totalCount"},
							bson.M{"$arrayElemAt": []interface{}{"$totalCount", 0}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"ads":    1,
					"amount": "$amount.count",
				},
			},
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	result := struct {
		Adverts []*advert.Advert `bson:"ads"`
		Amount  uint32           `bson:"amount"`
	}{}

	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.Adverts, result.Amount, nil
}

// SaveJob ...
func (r Repository) SaveJob(ctx context.Context, ad *advert.Job) error {
	_, err := r.advertCollection.InsertOne(ctx, ad)
	if err != nil {
		return err
	}

	return nil
}

// ChangeAdvert ...
func (r Repository) ChangeAdvert(ctx context.Context, ad *advert.Advert) error {
	_, err := r.advertCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": ad.ID,
		},
		bson.M{
			"$set": ad,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAdvertByID ...
func (r Repository) GetAdvertByID(ctx context.Context, advertID string) (*advert.Advert, error) {
	objID, err := primitive.ObjectIDFromHex(advertID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	result := r.advertCollection.FindOne(
		ctx,
		bson.M{
			"_id": objID,
		},
	)

	if result.Err() != nil {
		return nil, result.Err()
	}

	ad := advert.Advert{}

	err = result.Decode(&ad)
	if err != nil {
		return nil, err
	}

	return &ad, nil
}

// SaveCandidate ...
func (r Repository) SaveCandidate(ctx context.Context, ad *advert.Candidate) error {
	_, err := r.advertCollection.InsertOne(ctx, ad)
	if err != nil {
		return err
	}

	return nil
}

// GetBanners ...
func (r Repository) GetBanners(ctx context.Context, countryID string /*, place advert.Place*/, format advert.Format, amount uint32) ([]*advert.Banner, error) {
	match := bson.M{
		"status": advert.StatusActive,
		"type":   advert.TypeBanner,
		"format": format,
		// "places": place,
	}

	if countryID != "" {
		match["location.country_id"] = countryID
	}

	cursor, err := r.advertCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$sample": bson.M{
					"size": amount,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	adverts := make([]*advert.Banner, 0)

	for cursor.Next(ctx) {
		var ad advert.Banner
		err = cursor.Decode(&ad)
		if err != nil {
			return nil, err
		}
		adverts = append(adverts, &ad)
	}

	return adverts, nil
}

// GetCandidates ...
func (r Repository) GetCandidates(ctx context.Context, countryID string, format advert.Format, amount uint32) ([]string, error) {
	match := bson.M{
		"status": advert.StatusActive,
		"type":   advert.TypeCandidate,
		"format": format,
	}

	if countryID != "" {
		match["location.country_id"] = countryID
	}

	cursor, err := r.advertCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$sample": bson.M{
					"size": amount,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ids := make([]string, 0)

	for cursor.Next(ctx) {
		var ad advert.Banner
		err = cursor.Decode(&ad)
		if err != nil {
			return nil, err
		}
		ids = append(ids, ad.GetCreatorID())
	}

	return ids, nil
}

// GetJobs ...
func (r Repository) GetJobs(ctx context.Context, countryID string, amount uint32) ([]string, error) {
	match := bson.M{
		"status": advert.StatusActive,
		"type":   advert.TypeJob,
	}

	if countryID != "" {
		match["location.country_id"] = countryID
	}

	cursor, err := r.advertCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$sample": bson.M{
					"size": amount,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ids := make([]string, 0)

	for cursor.Next(ctx) {
		var ad advert.Banner
		err = cursor.Decode(&ad)
		if err != nil {
			return nil, err
		}
		ids = append(ids, ad.GetCompanyID())
	}

	log.Println("ids:", ids)

	return ids, nil
}

// RemoveAdvert ...
func (r Repository) RemoveAdvert(ctx context.Context, campaignID, advertID string) error {

	objCampaingID, err := primitive.ObjectIDFromHex(campaignID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objID, err := primitive.ObjectIDFromHex(advertID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.advertCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objCampaingID,
		},
		bson.M{
			"$pull": bson.M{
				"adverts": bson.M{
					"_id": objID,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAdvertCampaign ...
func (r Repository) RemoveAdvertCampaign(ctx context.Context, campaignID string) error {

	objCampaingID, err := primitive.ObjectIDFromHex(campaignID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.advertCollection.DeleteOne(
		ctx,
		bson.M{
			"_id": objCampaingID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// CreateAdvertByCampaign ...
func (r Repository) CreateAdvertByCampaign(ctx context.Context, campaignID string, data advert.CampaingAdvert) error {
	objCampaignID, err := primitive.ObjectIDFromHex(campaignID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.advertCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objCampaignID,
		},
		bson.M{
			"$push": bson.M{
				"adverts": data,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil

}

// CreateAdvertCampaign ...
func (r Repository) CreateAdvertCampaign(ctx context.Context, data advert.Campaing) error {

	_, err := r.advertCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// GetAdvertCampaigns ...
func (r Repository) GetAdvertCampaigns(ctx context.Context, profileID string, first int, after int) (*advert.Campaings, error) {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.advertCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"creator_id": objProfileID,
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"campaigns": bson.M{
						"$push": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$campaigns"},
							bson.M{"$size": "$campaigns"},
							0,
						},
					},
				},
			},
			{
				"$sort": bson.M{
					"campaigns.start_date": -1,
				},
			},
			{
				"$project": bson.M{
					"_id": 0,
					"campaigns": bson.M{
						"$slice": []interface{}{
							"$campaigns",
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
		return nil, err
	}
	defer cursor.Close(ctx)

	res := new(advert.Campaings)

	if cursor.Next(ctx) {
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}
	}

	if res != nil {
		for i := range res.Campaings {
			for _, ad := range res.Campaings[i].Adverts {
				if ad.Clicks > 0 && ad.Impressions > 0 {
					res.Campaings[i].CtrAVG = float64(ad.Clicks) / float64(ad.Impressions)
				}
			}
		}
	}

	return res, nil

}

// GetAdvertsByCampaignID ...
func (r Repository) GetAdvertsByCampaignID(ctx context.Context, campaignID, profileID string, first int, after int) (*advert.Adverts, error) {

	objCampaignID, err := primitive.ObjectIDFromHex(campaignID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.advertCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id":        objCampaignID,
					"creator_id": objProfileID,
				},
			},
			{
				"$addFields": bson.M{
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$adverts"},
							bson.M{"$size": "$adverts"},
							0,
						},
					},
					"adverts": "$adverts",
				},
			},
			{
				"$project": bson.M{
					"_id":      0,
					"campaign": "$$ROOT",
					"adverts": bson.M{
						"$slice": []interface{}{
							"$adverts",
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
		return nil, err
	}
	defer cursor.Close(ctx)

	res := new(advert.Adverts)

	if cursor.Next(ctx) {
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// GetAdvert ...
func (r Repository) GetAdvert(ctx context.Context, profileID string, advertType advert.Type) (*advert.GetAdvert, error) {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	match := bson.M{
		"status": "active",
		"impressions": bson.M{
			"$gt": 0,
		},
		"relevant_users": bson.M{
			"$elemMatch": bson.M{
				"$in": []interface{}{objProfileID},
			},
		},
	}

	if advertType != advert.TypeAny {
		match["type"] = advertType
	}

	cursor, err := r.advertCollection.Aggregate(
		ctx,
		[]bson.M{

			{
				"$match": match,
			},
			{
				"$sort": bson.M{
					"impressions": -1,
				},
			},
			{
				"$facet": bson.M{
					"adverts": []interface{}{
						bson.M{"$unwind": "$adverts"},
						bson.M{"$match": bson.M{
							"adverts.status": "active",
						}},
						bson.M{
							"$addFields": bson.M{
								"adverts.formats": "$$ROOT.format",
							},
						},
						bson.M{"$sort": bson.M{
							"adverts.impressions": 1,
						},
						},
					},
				},
			},
			{
				"$project": bson.M{
					"adverts": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$adverts.adverts"},
							bson.M{"$arrayElemAt": []interface{}{
								"$adverts.adverts", 0,
							}},
							0,
						},
					},
					"impressions": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$adverts"},
							bson.M{"$arrayElemAt": []interface{}{
								"$adverts.impressions", 0,
							}},
							0,
						},
					},
					"clicks": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$adverts"},
							bson.M{"$arrayElemAt": []interface{}{
								"$adverts.clicks", 0,
							}},
							0,
						},
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	res := advert.GetAdvert{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

// ChangeAdvertActions ...
func (r Repository) ChangeAdvertActions(ctx context.Context, advertID string, data advert.ActionType) error {
	objAdvertID, err := primitive.ObjectIDFromHex(advertID)
	if err != nil {
		return err
	}

	advertField := bson.M{}

	switch data {
	case advert.ActionTypeImpressions:
		advertField = bson.M{
			"impressions":           -1,
			"adverts.$.impressions": 1,
		}
	case advert.ActionTypeClicks:
		advertField = bson.M{
			"clicks":           -1,
			"adverts.$.clicks": 1,
		}
	}

	_, err = r.advertCollection.UpdateOne(
		ctx,
		bson.M{
			"adverts._id": objAdvertID,
		},
		bson.M{
			"$inc": advertField,
		},
	)
	if err != nil {
		return err
	}

	return nil

}

// ChangeStatus ...
func (r Repository) ChangeStatus(ctx context.Context, campaignID string, advertID string, data advert.Status) error {

	match := bson.M{}
	statusKey := "status"

	// Advert
	if advertID != "" {
		objAdvertID, err := primitive.ObjectIDFromHex(advertID)
		if err != nil {
			return err
		}
		statusKey = "adverts.$.status"
		match["adverts._id"] = objAdvertID
	}

	// Campaign
	if campaignID != "" {
		objCampaingID, err := primitive.ObjectIDFromHex(campaignID)
		if err != nil {
			return err
		}

		match["_id"] = objCampaingID
	}

	statusField := bson.M{}

	switch data {
	case advert.StatusActive:
		statusField = bson.M{statusKey: "active"}
	case advert.StatusCompleted:
		statusField = bson.M{statusKey: "completed"}
	case advert.StatusInActive:
		statusField = bson.M{statusKey: "in_active"}
	case advert.StatusPaused:
		statusField = bson.M{statusKey: "paused"}
	}

	_, err := r.advertCollection.UpdateOne(
		ctx,
		match,
		bson.M{
			"$set": statusField,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
