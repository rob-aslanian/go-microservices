package filters

import (
	"context"
	"errors"
	"log"

	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SaveUserSearchFilter ...
func (r Repository) SaveUserSearchFilter(ctx context.Context, data *requests.UserSearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// SaveJobSearchFilter ...
func (r Repository) SaveJobSearchFilter(ctx context.Context, data *requests.JobSearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// // SaveCandidateSearchFilters ...
// func (r Repository) SaveCandidateSearchFilter(ctx context.Context, data *requests.CandidateSearchFilter) error {
// 	_, err := r.filtersCollection.InsertOne(ctx, data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// SaveCompanySearchFilter ...
func (r Repository) SaveCompanySearchFilter(ctx context.Context, data *requests.CompanySearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// SaveServiceSearchFilter ...
func (r Repository) SaveServiceSearchFilter(ctx context.Context, data *requests.ServiceSearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// SaveServiceRequestSearchFilter ...
func (r Repository) SaveServiceRequestSearchFilter(ctx context.Context, data *requests.ServiceRequestSearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// 	return filters, nil
// }

// GetAllFilters ...
func (r Repository) GetAllFilters(ctx context.Context, userID string) ([]interface{}, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.filtersCollection.Find(ctx, bson.M{
		"user_id": userObjID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	filters := make([]interface{}, 0)

	for cursor.Next(ctx) {
		f := bson.M{}
		err := cursor.Decode(&f)
		if err != nil {
			return nil, err
		}
		log.Printf("got: %+v\n", f)

		if value, isEsxists := f["type"]; isEsxists {
			switch value {
			case requests.TypeCandidateFilterType:
				candidateFilter := requests.CandidateSearchFilter{}
				err = cursor.Decode(&candidateFilter)
				if err != nil {
					return nil, err
				}
				filters = append(filters, candidateFilter)

			case requests.TypeCompanyFilterType:
				companyFilter := requests.CompanySearchFilter{}
				err = cursor.Decode(&companyFilter)
				if err != nil {
					return nil, err
				}
				filters = append(filters, companyFilter)

			case requests.TypeJobFilterType:
				jobFilter := requests.JobSearchFilter{}
				err = cursor.Decode(&jobFilter)
				if err != nil {
					return nil, err
				}
				filters = append(filters, jobFilter)

			case requests.TypeUserFilterType:
				userFilter := requests.UserSearchFilter{}
				err = cursor.Decode(&userFilter)
				if err != nil {
					return nil, err
				}
				filters = append(filters, userFilter)
			case requests.TypeServiceFilterType:
				serviceFilter := requests.ServiceSearchFilter{}
				err = cursor.Decode(&serviceFilter)
				if err != nil {
					return nil, err
				}
				filters = append(filters, serviceFilter)
			case requests.TypeServiceRequestFilterType:
				serviceFilter := requests.ServiceRequestSearchFilter{}
				err = cursor.Decode(&serviceFilter)
				if err != nil {
					return nil, err
				}
				filters = append(filters, serviceFilter)
			}
		}
	}

	return filters, nil
}

// RemoveFilter ...
func (r Repository) RemoveFilter(ctx context.Context, filterID string) error {
	//TODO check if ID is valid objectID
	if filterID == "" {
		return errors.New("empty_id")
	}
	filterObjID, err := primitive.ObjectIDFromHex(filterID)
	if err != nil {
		return errors.New("wrong_id")
	}

	_, err = r.filtersCollection.DeleteOne(
		ctx,
		bson.M{
			"_id": filterObjID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserSearchFilterForCompany ...
func (r Repository) SaveUserSearchFilterForCompany(ctx context.Context, data *requests.UserSearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// SaveJobSearchFilterForCompany ...
func (r Repository) SaveJobSearchFilterForCompany(ctx context.Context, data *requests.JobSearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// SaveCandidateSearchFilterForCompany ...
func (r Repository) SaveCandidateSearchFilterForCompany(ctx context.Context, data *requests.CandidateSearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// SaveCompanySearchFilterForCompany ...
func (r Repository) SaveCompanySearchFilterForCompany(ctx context.Context, data *requests.CompanySearchFilter) error {
	_, err := r.filtersCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// GetAllFiltersForCompany ...
func (r Repository) GetAllFiltersForCompany(ctx context.Context, companyID string) ([]interface{}, error) {
	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.filtersCollection.Find(ctx, bson.M{
		"company_id": companyObjID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	filters := make([]interface{}, 0)

	for cursor.Next(ctx) {
		f := bson.M{}
		cursor.Decode(&f)

		if value, isEsxists := f["type"]; isEsxists {
			switch value {
			case requests.TypeCandidateFilterType:
				candidateFilter := requests.CandidateSearchFilter{}
				cursor.Decode(&candidateFilter)
				filters = append(filters, candidateFilter)

			case requests.TypeCompanyFilterType:
				companyFilter := requests.CompanySearchFilter{}
				cursor.Decode(&companyFilter)
				filters = append(filters, companyFilter)

			case requests.TypeJobFilterType:
				jobFilter := requests.JobSearchFilter{}
				cursor.Decode(&jobFilter)
				filters = append(filters, jobFilter)

			case requests.TypeUserFilterType:
				userFilter := requests.UserSearchFilter{}
				cursor.Decode(&userFilter)
				filters = append(filters, userFilter)
			case requests.TypeServiceFilterType:
				serviceFilter := requests.ServiceSearchFilter{}
				cursor.Decode(&serviceFilter)
				filters = append(filters, serviceFilter)
			case requests.TypeServiceRequestFilterType:
				serviceFilter := requests.ServiceRequestSearchFilter{}
				cursor.Decode(&serviceFilter)
				filters = append(filters, serviceFilter)
			}
		}
	}

	return filters, nil
}

// RemoveFilterForCompany ...
func (r Repository) RemoveFilterForCompany(ctx context.Context, filterID, companyID string) error {
	filterObjID, err := primitive.ObjectIDFromHex(filterID)
	if err != nil {
		return errors.New("wrong_id")
	}

	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New("wrong_id")
	}

	_, err = r.filtersCollection.DeleteOne(
		ctx,
		bson.M{
			"_id":        filterObjID,
			"company_id": companyObjID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
