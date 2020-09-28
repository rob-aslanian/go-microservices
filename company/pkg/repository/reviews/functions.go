package reviewsRepository

import (
	"context"
	"errors"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
)

const (
	companyReviewsCollection        = "company_reviews"
	companyReviewsReportsCollection = "company_reviews_reports"
)

// // SaveNewCompanyAccount ...
// func (r Repository) SaveNewCompanyAccount(ctx context.Context, data *account.Account) error {
// 	err := r.collections[companyReviewsCollection].Insert(data)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// AddReview ...
func (r Repository) AddReview(ctx context.Context, review *profile.Review) error {
	if !bson.IsObjectIdHex(review.GetID()) || !bson.IsObjectIdHex(review.GetAuthorID()) ||
		!bson.IsObjectIdHex(review.GetCompanyID()) {
		return errors.New("wrong_id")
	}

	err := r.collections[companyReviewsCollection].Insert(review)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyReview ...
func (r Repository) DeleteCompanyReview(ctx context.Context, companyID, userID string, reviewID string) error {
	if !bson.IsObjectIdHex(companyID) || !bson.IsObjectIdHex(userID) ||
		!bson.IsObjectIdHex(reviewID) {
		return errors.New("wrong_id")
	}

	err := r.collections[companyReviewsCollection].Remove(
		bson.M{
			"_id":        bson.ObjectIdHex(reviewID),
			"company_id": bson.ObjectIdHex(companyID),
			"author_id":  bson.ObjectIdHex(userID),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetCompanyReviews ...
func (r Repository) GetCompanyReviews(ctx context.Context, companyID string, first uint32, after uint32) ([]*profile.Review, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]*profile.Review, 0)

	res := r.collections[companyReviewsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": bson.ObjectIdHex(companyID),
				},
			},
			{
				"$sort": bson.M{
					"reviews.created_at": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := res.All(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetUsersRevies ...
func (r Repository) GetUsersRevies(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Review, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]*profile.Review, 0)

	res := r.collections[companyReviewsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"author_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$sort": bson.M{
					"reviews.created_at": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := res.All(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// AddCompanyReviewReport ...
func (r Repository) AddCompanyReviewReport(ctx context.Context, reviewReport *profile.ReviewReport) error {
	err := r.collections[companyReviewsReportsCollection].Insert(reviewReport)
	if err != nil {
		return err
	}

	return nil
}

// GetAvarageRateOfCompany ...
func (r Repository) GetAvarageRateOfCompany(ctx context.Context, companyID string) (float32, uint32, error) {
	if !bson.IsObjectIdHex(companyID) {
		return 0, 0, errors.New("wrong_id")
	}

	result := struct {
		AmountOfRate uint32  `bson:"amount"`
		AvarageRate  float32 `bson:"avarage_rate"`
	}{}

	log.Println("companyID:", companyID)

	aggregationPipeline := []bson.M{
		{
			"$match": bson.M{
				"company_id": bson.ObjectIdHex(companyID),
			},
		},
		{
			"$group": bson.M{
				"_id": "company_id",
				"avarage_rate": bson.M{
					"$avg": "$rate",
				},
				"amount": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	pipe := r.collections[companyReviewsCollection].Pipe(aggregationPipeline)

	err := pipe.One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	return result.AvarageRate, result.AmountOfRate, nil
}

// GetAmountOfEachRate ...
func (r Repository) GetAmountOfEachRate(ctx context.Context, companyID string) (map[uint32]uint32, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	result := []struct {
		Rate   uint32 `bson:"_id"`
		Amount uint32 `bson:"amount"`
	}{}

	aggregationPipeline := []bson.M{
		{
			"$match": bson.M{
				"company_id": bson.ObjectIdHex(companyID),
			},
		},
		{
			"$group": bson.M{
				"_id": "$rate",
				"amount": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	pipe := r.collections[companyReviewsCollection].Pipe(aggregationPipeline)

	err := pipe.All(&result)
	if err != nil {
		return nil, err
	}

	m := make(map[uint32]uint32)

	for i := range result {
		m[result[i].Rate] = result[i].Amount
	}

	return m, nil
}

// GetAmountOfReviewsOfUser ...
func (r Repository) GetAmountOfReviewsOfUser(ctx context.Context, userID string) (int32, error) {
	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	result := struct {
		Amount int32 `bson:"amount"`
	}{}

	aggregationPipeline := []bson.M{
		{
			"$match": bson.M{
				"author_id": bson.ObjectIdHex(userID),
			},
		},
		{
			"$group": bson.M{
				"_id": "$author_id",
				"amount": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	pipe := r.collections[companyReviewsCollection].Pipe(aggregationPipeline)

	err := pipe.One(&result)
	if err != nil {
		return 0, err
	}

	return result.Amount, nil
}
