package servicesrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/review"

	offer "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/offers"
	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/qualifications"
	serviceorder "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/service-order"
	office "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/v-office"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// docker-compose exec mongo_db mongo --username developer --authenticationDatabase admin -p

// CreateVOffice adds a new vOffice to the db
func (r *Repository) CreateVOffice(ctx context.Context, vOffice *office.Office) error {
	_, err := r.officeCollection.InsertOne(ctx, vOffice)
	if err != nil {
		return err
	}

	return nil
}

// ChangeVOffice ...
func (r *Repository) ChangeVOffice(ctx context.Context, vOffice *office.Office) error {

	if vOffice.ID.IsZero() {
		return errors.New(`wrong_id`)
	}

	match := bson.M{
		"_id": vOffice.ID,
	}

	if vOffice.CompanyID != nil {
		match["company_id"] = vOffice.CompanyID
	} else {
		match["user_id"] = vOffice.UserID
	}

	_, err := r.officeCollection.UpdateOne(ctx,
		match,
		bson.M{
			"$set": bson.M{
				"name":        vOffice.Name,
				"description": vOffice.Description,
				"category":    vOffice.Category,
				"location":    vOffice.Location,
				"languages":   vOffice.Languages,
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// GetVOffice returns us the vOffice
func (r *Repository) GetVOffice(ctx context.Context, companyID, userID string) ([]*office.Office, error) {

	match := bson.M{}

	var userObjID, companyObjID primitive.ObjectID

	if userID != "" {
		var err error
		userObjID, err = primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
	}

	if companyID != "" {
		var err error
		companyObjID, err = primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}

	}

	if !companyObjID.IsZero() {
		match["company_id"] = companyObjID
	} else {
		match["user_id"] = userObjID
	}

	cursor, err := r.officeCollection.Aggregate(ctx, []bson.M{
		{
			"$match": match,
		},
	})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not_found")
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*office.Office, 0)

	for cursor.Next(ctx) {
		r := new(office.Office)
		err := cursor.Decode(r)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

// GetVOfficeByID ...
func (r *Repository) GetVOfficeByID(ctx context.Context, officeID string) (*office.Office, error) {

	objID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return nil, errors.New("wrong_ID")
	}

	match := bson.M{
		"_id": objID,
	}
	result := r.officeCollection.FindOne(
		ctx,
		match,
	)

	if result.Err() != nil {
		return nil, result.Err()
	}

	office := office.Office{}

	err = result.Decode(&office)
	if err != nil {
		return nil, err
	}

	return &office, nil
}

// RemoveVOffice ...
func (r *Repository) RemoveVOffice(ctx context.Context, officeID string) error {
	objID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_ID")
	}

	r.officeCollection.DeleteOne(
		ctx,
		bson.M{
			"_id": objID,
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	return nil
}

// GetServicesRequest ...
func (r *Repository) GetServicesRequest(ctx context.Context, userID, companyID string) ([]*servicerequest.Request, error) {

	match := []bson.M{}

	var userObjID, companyObjID primitive.ObjectID

	if userID != "" {
		var err error
		userObjID, err = primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
	}

	if companyID != "" {
		var err error
		companyObjID, err = primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}

	}

	if !companyObjID.IsZero() {
		match = []bson.M{
			{

				"$match": bson.M{
					"company_id": companyObjID,
				},
			},
		}
	} else {
		match = []bson.M{
			{
				"$match": bson.M{
					"user_id": userObjID,
					"company_id": bson.M{
						"$exists": false,
					},
				},
			},
		}
	}

	cursor, err := r.requestCollection.Aggregate(ctx, match)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	services := make([]*servicerequest.Request, 0)

	for cursor.Next(ctx) {
		r := new(servicerequest.Request)
		err := cursor.Decode(r)
		if err != nil {
			return nil, err
		}

		services = append(services, r)
	}

	return services, nil

}

// GetServiceRequest ...
func (r *Repository) GetServiceRequest(ctx context.Context, serviceID string) (*servicerequest.Request, error) {

	serviceObjID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return nil, errors.New("wrong_ID")
	}

	service := servicerequest.Request{}

	result := r.requestCollection.FindOne(ctx, bson.M{
		"_id": serviceObjID,
	})
	err = result.Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}

	return &service, nil

}

// ChangeServicesRequestStatus ...
func (r *Repository) ChangeServicesRequestStatus(ctx context.Context, serviceID string, serviceStatus servicerequest.ServiceReqestStatus) error {

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_serviceID")
	}

	_, err = r.requestCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objServiceID,
		},
		bson.M{
			"$set": serviceStatus,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsOutOfOffice sets the office IsOut property to the selected value.
func (r *Repository) IsOutOfOffice(ctx context.Context, officeID string, isOut bool, returnDate *time.Time) error {
	objID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_ID")
	}

	var date *time.Time = nil

	if !isOut {
		date = nil
	} else {
		if returnDate != nil {
			date = returnDate
		} else {
			date = nil
		}

	}

	r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"is_out":      isOut,
				"return_date": date,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	return nil
}

// IsURLBusy checks if url is busy
func (r *Repository) IsURLBusy(ctx context.Context, url string) (bool, error) {
	count, err := r.officeCollection.CountDocuments(
		ctx,
		bson.M{
			"url": url,
		})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, err
	}

	return true, nil
}

// ChangeVOfficeName changes the voffice name
func (r *Repository) ChangeVOfficeName(ctx context.Context, officeID, name string) error {
	objID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_ID")
	}
	r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},

		bson.M{
			"$set": bson.M{
				"name": name,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	return nil
}

// AddChangeVOfficeDescription adds or changes the voffice description
func (r *Repository) AddChangeVOfficeDescription(ctx context.Context, officeID, description string) error {
	objID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_ID")
	}

	r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},

		bson.M{
			"$set": bson.M{
				"description": description,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	return nil
}

// AddVOfficeLanguages adds or changes the voffice qualifications
func (r *Repository) AddVOfficeLanguages(ctx context.Context, userID, companyID, officeID string, langs []*qualifications.Language) error {

	objID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_ID")
	}
	match := bson.M{
		"_id": objID,
	}

	var userObjID, companyObjID primitive.ObjectID

	if userID != "" {
		var err error
		userObjID, err = primitive.ObjectIDFromHex(userID)
		if err != nil {
			return errors.New(`wrong_id`)
		}
	}

	if companyID != "" {
		var err error
		companyObjID, err = primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return errors.New(`wrong_id`)
		}

	}

	if !companyObjID.IsZero() {
		match["company_id"] = companyObjID
	} else {
		match["user_id"] = userObjID
	}

	r.officeCollection.UpdateOne(
		ctx,
		match,
		bson.M{
			"$push": bson.M{
				"languages": bson.M{
					"$each": langs,
				},
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	return nil
}

// ChangeVOfficeLanguage adds or changes the voffice qualifications
func (r *Repository) ChangeVOfficeLanguage(ctx context.Context, officeID string, langs []*qualifications.Language) error {

	objID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_ID")
	}
	r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},

		bson.M{
			"$set": bson.M{
				"languages": langs,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	return nil
}

// RemoveVOfficeLanguages removes qualifications languages from db
func (r *Repository) RemoveVOfficeLanguages(ctx context.Context, officeID string, languageIDs []string) error {
	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_officeID")
	}

	objLanguages := make([]primitive.ObjectID, 0, len(languageIDs))
	for _, objLanguage := range languageIDs {
		obj, err := primitive.ObjectIDFromHex(objLanguage)
		if err != nil {
			return err
		}
		objLanguages = append(objLanguages, obj)
	}

	_, err = r.officeCollection.UpdateMany(
		ctx,
		bson.M{
			"_id": objOfficeID,
		},
		bson.M{
			"$pull": bson.M{
				"languages": bson.M{
					"_id": bson.M{
						"$in": objLanguages,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddVOfficePortfolio add portfolio to the vOffice
// TODO: needs to add also other than just files
func (r *Repository) AddVOfficePortfolio(ctx context.Context, officeID string, portfolio *office.Portfolio,
	userID string, companyID string) error {
	if officeID == "" {
		return errors.New("officeID empty")
	}

	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	fileIDs := make([]primitive.ObjectID, 0, len(portfolio.Files))
	for _, files := range portfolio.Files {
		fileIDs = append(fileIDs, files.ID)
	}
	// append files
	files := []struct {
		File *file.File `bson:"unused_files"`
	}{}
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"_id": objOfficeID,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":          0,
				"unused_files": 1,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$unused_files",
			},
		},
		bson.M{
			"$match": bson.M{
				"unused_files.type": "Portfolio",
				"unused_files.id": bson.M{
					"$in": fileIDs,
				},
			},
		},
	}

	cursor, err := r.officeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	portfolio.Files = make([]*file.File, len(files))
	for i := range files {
		portfolio.Files[i] = files[i].File
	}

	if companyID != "" {
		portfolio.SetCompanyID(companyID)
	} else {
		portfolio.SetUserID(userID)
	}

	_, err = r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objOfficeID,
		},

		bson.M{
			"$push": bson.M{
				"portfolio": portfolio,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	// delete old files
	_, err = r.officeCollection.UpdateMany(
		ctx,
		bson.M{
			"_id": objOfficeID,
		},
		bson.M{
			"$pull": bson.M{
				"unused_files": bson.M{
					"type": "Portfolio",
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeVOfficePortfolio add portfolio to the vOffice
// TODO: needs to add also other than just files
func (r *Repository) ChangeVOfficePortfolio(ctx context.Context, officeID, portfolioID string, portfolio *office.Portfolio) error {
	if officeID == "" {
		return errors.New("officeID empty")
	}
	if portfolioID == "" {
		return errors.New("portfolioID empty")
	}

	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objPortfolioID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	// linkIntest := bson.M{}

	_, err = r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":           objOfficeID,
			"portfolio._id": objPortfolioID,
		},

		bson.M{
			"$set": bson.M{
				"portfolio.$": portfolio,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	// delete old files
	_, err = r.officeCollection.UpdateMany(
		ctx,
		bson.M{
			"_id": objOfficeID,
		},
		bson.M{
			"$pull": bson.M{
				"unused_files": bson.M{
					"type": "Portfolio",
				},
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// RemoveVOfficePortfolio ...
func (r *Repository) RemoveVOfficePortfolio(ctx context.Context, officeID, portfolioID string) error {
	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_officeID")
	}
	objPortfolioID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return errors.New("wrong_portfolioID")
	}

	_, err = r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objOfficeID,
		},
		bson.M{
			"$pull": bson.M{
				"portfolio": bson.M{
					"_id": objPortfolioID,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddFileInVOfficePortfolio adds files
func (r *Repository) AddFileInVOfficePortfolio(ctx context.Context, officeID, portfolioID, companyID string, file *file.File) error {
	obOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_officeID")
	}
	objPortfolioID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return errors.New("wrong_portfolioID")
	}

	match := bson.M{
		"_id":           obOfficeID,
		"portfolio._id": objPortfolioID,
	}

	if companyID != "" {
		objCompanyID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return errors.New("wrong_objServiceID")
		}
		match["company_id"] = objCompanyID
	}

	_, dbErr := r.officeCollection.UpdateOne(
		ctx,
		match,
		bson.M{
			"$push": bson.M{
				"portfolio.$.files": file,
			},
		},
	)

	if dbErr != nil {
		return dbErr
	}

	return nil
}

// AddFileInServiceRequest adds files
func (r *Repository) AddFileInServiceRequest(ctx context.Context, serviceID, companyID string, file *file.File) error {

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_officeID")
	}

	match := bson.M{
		"_id": objServiceID,
	}

	if companyID != "" {
		objCompanyID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return errors.New("wrong_objServiceID")
		}
		match["company_id"] = objCompanyID
	}

	_, dbErr := r.requestCollection.UpdateOne(
		ctx,
		match,
		bson.M{
			"$push": bson.M{
				"files": file,
			},
		},
	)

	if dbErr != nil {
		return dbErr
	}

	return nil
}

// AddFileInOrderService adds files
func (r *Repository) AddFileInOrderService(ctx context.Context, orderID string, file *file.File) error {

	objOrderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("wrong_officeID")
	}

	_, err = r.orderCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objOrderID,
		},
		bson.M{
			"$push": bson.M{
				"order_detail.files": file,
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// ChangeVofficeCover ...
func (r *Repository) ChangeVofficeCover(ctx context.Context, officeID string, companyID string, image string) error {

	obOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return err
	}

	r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": obOfficeID,
		},
		bson.M{
			"$set": bson.M{
				"cover_image": image,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveVofficeCover ...
func (r *Repository) RemoveVofficeCover(ctx context.Context, officeID, companyID string) error {

	obOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return err
	}

	_, err = r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": obOfficeID,
		},
		bson.M{
			"$set": bson.M{
				"cover_origin_image": "",
				"cover_image":        "",
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddFileInVofficeService ...
func (r *Repository) AddFileInVofficeService(ctx context.Context, officeID, serviceID, companyID string, file *file.File) error {
	obOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_officeID")
	}

	obServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_portfolioID")
	}

	match := bson.M{
		"_id":       obServiceID,
		"office_id": obOfficeID,
	}

	if companyID != "" {
		objCompanyID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return errors.New("wrong_objServiceID")
		}
		match["company_id"] = objCompanyID
	}

	_, dbErr := r.servicesCollection.UpdateOne(
		ctx,
		match,
		bson.M{
			"$push": bson.M{
				"files": file,
			},
		},
	)

	if dbErr != nil {
		return dbErr
	}

	return nil
}

// ChangeVofficeOriginCover ...
func (r *Repository) ChangeVofficeOriginCover(ctx context.Context, officeID string, companyID string, image string) error {

	obOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return err
	}

	r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": obOfficeID,
		},
		bson.M{
			"$set": bson.M{
				"cover_origin_image": image,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFilesInVOfficePortfolio ...
func (r *Repository) RemoveFilesInVOfficePortfolio(ctx context.Context, officeID, portfolioID string, ids []string) error {
	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_officeID")
	}
	objPortfolioID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return errors.New("wrong_portfolioID")
	}

	isObjectIDs := make([]primitive.ObjectID, 0, len(ids))
	for i := range ids {
		objIDs, err := primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			return errors.New("wrong_fileID")
		}
		isObjectIDs = append(isObjectIDs, objIDs)
	}

	removeIDs := make([]primitive.ObjectID, 0, len(isObjectIDs))
	for i := range isObjectIDs {
		removeIDs = append(removeIDs, isObjectIDs[i])
	}

	_, err = r.officeCollection.DeleteMany(
		ctx,
		bson.M{
			"_id":                    objOfficeID,
			"portfolios._id":         objPortfolioID,
			"portfolios.$.files._id": removeIDs,
		},
	)
	if err != nil {
		return err
	}

	return nil

}

// RemoveFilesInServiceRequest ...
func (r *Repository) RemoveFilesInServiceRequest(ctx context.Context, serviceID string, ids []string) error {

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_officeID")
	}

	isObjectIDs := make([]primitive.ObjectID, 0, len(ids))
	for i := range ids {
		objIDs, err := primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			return errors.New("wrong_fileID")
		}
		isObjectIDs = append(isObjectIDs, objIDs)
	}

	removeIDs := make([]primitive.ObjectID, 0, len(isObjectIDs))
	for i := range isObjectIDs {
		removeIDs = append(removeIDs, isObjectIDs[i])
	}

	_, err = r.requestCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objServiceID,
		},
		bson.M{
			"$pull": bson.M{
				"files": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil

}

// RemoveLinksInVOfficePortfolio ...
func (r *Repository) RemoveLinksInVOfficePortfolio(ctx context.Context, officeID, portfolioID string, ids []string) error {
	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_officeID")
	}
	objPortfolioID, err := primitive.ObjectIDFromHex(portfolioID)
	if err != nil {
		return errors.New("wrong_portfolioID")
	}

	isObjectIDs := make([]primitive.ObjectID, 0, len(ids))
	for i := range ids {
		objIDs, err := primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			return errors.New("wrong_fileID")
		}
		isObjectIDs = append(isObjectIDs, objIDs)
	}

	_, err = r.officeCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objOfficeID,
			//
			"portfolio._id": objPortfolioID,
		},
		bson.M{
			"$pull": bson.M{
				"portfolio.$.links": bson.M{
					"id": bson.M{"$in": isObjectIDs}},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetVOfficeServices retrieves all the services where officeID is the ID passed
func (r *Repository) GetVOfficeServices(ctx context.Context, officeID string) ([]*servicerequest.Service, error) {
	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"office_id": objOfficeID,
			},
		},
	}

	cursor, err := r.servicesCollection.Aggregate(ctx, pipeline)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not_found")
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*servicerequest.Service, 0)

	for cursor.Next(ctx) {
		r := new(servicerequest.Service)
		err := cursor.Decode(r)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

// GetAllServices retrieves all the services where officeID is the ID passed
func (r *Repository) GetAllServices(ctx context.Context, profileID string, isCompany bool) ([]*servicerequest.Service, error) {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, err
	}

	match := bson.M{}

	if isCompany {
		match["company_id"] = objProfileID
	} else {
		match["user_id"] = objProfileID
	}

	cursor, err := r.servicesCollection.Aggregate(ctx, []bson.M{
		{
			"$match": match,
		},
	})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not_found")
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*servicerequest.Service, 0)

	for cursor.Next(ctx) {
		r := new(servicerequest.Service)
		err := cursor.Decode(r)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

// GetVOfficeService ...
func (r *Repository) GetVOfficeService(ctx context.Context, officeID, serviceID string) (*servicerequest.Service, error) {

	match := bson.M{}

	// Get By service ID
	if serviceID != "" {
		objServiceID, err := primitive.ObjectIDFromHex(serviceID)
		if err != nil {
			return nil, err
		}

		match["_id"] = objServiceID
	}

	// Get By office ID
	if officeID != "" {
		objOfficeID, err := primitive.ObjectIDFromHex(officeID)
		if err != nil {
			return nil, err
		}

		match["office_id"] = objOfficeID

	}

	result := r.servicesCollection.FindOne(ctx, match)

	service := &servicerequest.Service{}

	err := result.Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not_found")
		}
		return nil, err
	}

	return service, nil
}

// GetVOfficeServiceByID ...
func (r *Repository) GetVOfficeServiceByID(ctx context.Context, serviceID string) (*servicerequest.Service, error) {
	objServiceID, err := primitive.ObjectIDFromHex(serviceID)

	if err != nil {
		return nil, err
	}

	match := bson.M{
		"_id": objServiceID,
	}

	result := r.servicesCollection.FindOne(ctx, match)

	service := &servicerequest.Service{}

	err = result.Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not_found")
		}
		return nil, err
	}

	return service, nil
}

// ChangeVOfficeServiceStatus ...
func (r *Repository) ChangeVOfficeServiceStatus(ctx context.Context, officeID, serviceID string, serviceStatus servicerequest.ServiceStatus) error {

	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return errors.New("wrong_officeID")
	}
	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_serviceID")
	}

	log.Printf("Check service status in REPOO %+v", serviceStatus)

	_, err = r.servicesCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":       objServiceID,
			"office_id": objOfficeID,
		},
		bson.M{
			"$set": serviceStatus,
		},
	)
	if err != nil {
		return err
	}

	return nil

}

// GetVOfficePortfolio ...
func (r *Repository) GetVOfficePortfolio(ctx context.Context, officeID string) ([]*servicerequest.Portfolio, error) {
	objOfficeID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": objOfficeID,
			},
		},
		{
			"$project": bson.M{
				"portfolio": 1,
			},
		},
	}

	cursor, err := r.officeCollection.Aggregate(ctx, pipeline)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not_found")
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*servicerequest.Portfolio, 0)

	for cursor.Next(ctx) {
		r := new(servicerequest.Portfolio)
		err := cursor.Decode(r)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

// AddVOfficeService adds or changes a service in the service collection
func (r *Repository) AddVOfficeService(ctx context.Context, service *servicerequest.Service) error {

	_, err := r.servicesCollection.InsertOne(
		ctx,
		service,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	return nil
}

// ChangeVOfficeService adds or changes a service in the service collection
func (r *Repository) ChangeVOfficeService(ctx context.Context, serviceID, offiecID string, service *servicerequest.Service) error {
	log.Println("Repository ChangeVOfficeService ServiceID:", serviceID)
	if serviceID == "" {
		return errors.New("serviceID empty")
	}

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_ID")
	}
	objOfficeID, err := primitive.ObjectIDFromHex(offiecID)
	if err != nil {
		return errors.New("wrong_ID")
	}

	_, err = r.servicesCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":       objServiceID,
			"office_id": objOfficeID,
		},

		bson.M{
			"$set": bson.M{
				"additional_details": service.AdditionalDetails,
				"title":              service.Title,
				"description":        service.Description,
				"status":             service.Status,
				"category":           service.Category,
				"currency":           service.Currency,
				"price":              service.Price,
				"delivery_time":      service.DeliveryTime,
				"fixed_price_amount": service.FixedPriceAmmount,
				"min_price_amount":   service.MinPriceAmmout,
				"max_price_amount":   service.MaxPriceAmmout,
				"location_type":      service.LocationType,
				"location":           service.Location,
				"is_draft":           service.IsDraft,
				"is_remote":          service.IsRemote,
			},
		},
		// options.Update().SetUpsert(true),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}

	// delete old files
	// r.servicesCollection.UpdateMany(
	// 	ctx,
	// 	bson.M{
	// 		"_id": objServiceID,
	// 	},
	// 	bson.M{
	// 		"$pull": bson.M{
	// 			"unused_files": bson.M{
	// 				"type": "Service",
	// 			},
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	return nil
}

// RemoveVOfficeService ...
func (r *Repository) RemoveVOfficeService(ctx context.Context, serviceID string) error {
	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_officeID")
	}
	_, err = r.servicesCollection.DeleteOne(
		ctx,
		bson.M{
			"_id": objServiceID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFilesInVOfficeService ...
func (r *Repository) RemoveFilesInVOfficeService(ctx context.Context, serviceID string, ids []string) error {
	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return errors.New("wrong_officeID")
	}

	isObjectIDs := make([]primitive.ObjectID, 0, len(ids))
	for i := range ids {
		objIDs, err := primitive.ObjectIDFromHex(ids[i])
		if err != nil {
			return errors.New("wrong_fileID")
		}
		isObjectIDs = append(isObjectIDs, objIDs)
	}

	removeIDs := make([]primitive.ObjectID, 0, len(isObjectIDs))
	for i := range isObjectIDs {
		removeIDs = append(removeIDs, isObjectIDs[i])
	}

	_, err = r.servicesCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objServiceID,
		},
		bson.M{
			"$pull": bson.M{
				"files": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// func (r *Repository) RemoveVOfficeService(ctx context.Context, officeID string,)

// AddServicesRequest adds or changes a request in the request collection
func (r *Repository) AddServicesRequest(ctx context.Context, companyID string, request *servicerequest.Request) error {

	var objCompanyID primitive.ObjectID

	if companyID != "" {
		var err error
		objCompanyID, err = primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return errors.New("wrong_ID")
		}
		request.CompanyID = &objCompanyID
	}

	_, err := r.requestCollection.InsertOne(ctx, request)
	if err != nil {
		return err
	}

	return nil

}

// ChangeServicesRequest ...
func (r *Repository) ChangeServicesRequest(ctx context.Context, request *servicerequest.Request) error {

	objServicetID, err := primitive.ObjectIDFromHex(request.GetID())
	if err != nil {
		return errors.New("wrong_serviceID")
	}

	_, err = r.requestCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objServicetID,
		},
		bson.M{
			"$set": request,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveServicesRequest ...
func (r *Repository) RemoveServicesRequest(ctx context.Context, requestID string) error {
	objRequestID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		return errors.New("wrong_serviceID")
	}

	_, err = r.requestCollection.DeleteOne(
		ctx,
		bson.M{
			"_id": objRequestID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SendProposalForServiceRequest ...
func (r *Repository) SendProposalForServiceRequest(ctx context.Context, data *offer.Proposal) error {
	_, err := r.offerCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// IgnoreProposalForServiceRequest ...
func (r *Repository) IgnoreProposalForServiceRequest(ctx context.Context, profileID, proposalID string) error {

	// Profile ID
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return errors.New("wrong_serviceID")
	}

	// Proposal ID
	objProposalID, err := primitive.ObjectIDFromHex(proposalID)
	if err != nil {
		return errors.New("wrong_serviceID")
	}

	r.offerCollection.DeleteOne(
		ctx,
		bson.M{
			"_id":      objProposalID,
			"owner_id": objProfileID,
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("not_found")
		}
		return err
	}
	return nil
}

// OrderService ...
func (r *Repository) OrderService(ctx context.Context, data *serviceorder.Order) error {
	_, err := r.orderCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// AcceptOrderService ...
func (r *Repository) AcceptOrderService(ctx context.Context, data *serviceorder.Order, oldID string) error {
	objOldID, err := primitive.ObjectIDFromHex(oldID)
	if err != nil {
		return errors.New("wrong_serviceID")
	}

	// Set new order (buyer)
	_, err = r.orderCollection.InsertOne(ctx, data)
	// Update old order  (seller)
	_, err = r.orderCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objOldID,
		},
		bson.M{
			"$set": bson.M{
				"order_detail.status": "in_progress",
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderByReferalID ...
func (r *Repository) GetOrderByReferalID(ctx context.Context, referalID string) (*serviceorder.Order, error) {
	objReferalID, err := primitive.ObjectIDFromHex(referalID)

	if err != nil {
		return nil, err
	}

	result := r.orderCollection.FindOne(ctx,
		bson.M{
			"referal_id": objReferalID,
		},
	)

	order := &serviceorder.Order{}

	err = result.Decode(&order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetProposalByID ...
func (r *Repository) GetProposalByID(ctx context.Context, profileID, proposalID string) (*offer.Proposal, error) {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	objProposalID, err := primitive.ObjectIDFromHex(proposalID)
	if err != nil {
		return nil, err
	}

	result := r.offerCollection.FindOne(ctx,
		bson.M{
			"_id":      objProposalID,
			"owner_id": objProfileID,
		},
	)

	proposal := &offer.Proposal{}

	err = result.Decode(&proposal)
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

// CancelServiceOrder ...
func (r *Repository) CancelServiceOrder(ctx context.Context, profileID, orderID string) error {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return errors.New("wrong_profileID")
	}

	objOrderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("wrong_orderID")
	}

	_, err = r.orderCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":      objOrderID,
			"owner_id": objProfileID,
		},
		bson.M{
			"$set": bson.M{
				"order_detail.status": "canceled",
			},
		},
	)

	if err != nil {
		return err
	}

	seller, _ := r.GetOrderByReferalID(ctx, orderID)

	if seller != nil {
		return r.CancelServiceOrder(ctx, seller.GetOwnerID(), seller.GetID())
	}

	return nil
}

// DeclineServiceOrder ...
func (r *Repository) DeclineServiceOrder(ctx context.Context, profileID, orderID string) error {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return errors.New("wrong_profileID")
	}

	objOrderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("wrong_orderID")
	}

	_, err = r.orderCollection.DeleteOne(
		ctx,
		bson.M{
			"_id":      objOrderID,
			"owner_id": objProfileID,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// DeliverServiceOrder ...
func (r *Repository) DeliverServiceOrder(ctx context.Context, profileID, orderID string) error {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return errors.New("wrong_profileID")
	}

	objOrderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("wrong_orderID")
	}

	_, err = r.orderCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":      objOrderID,
			"owner_id": objProfileID,
		},
		bson.M{
			"$set": bson.M{
				"order_detail.status": "delivered",
			},
		},
	)

	if err != nil {
		return err
	}

	seller, _ := r.GetOrderByReferalID(ctx, orderID)

	if seller != nil {
		return r.DeliverServiceOrder(ctx, seller.GetOwnerID(), seller.GetID())
	}

	return nil
}

// AcceptDeliverdServiceOrder ...
func (r *Repository) AcceptDeliverdServiceOrder(ctx context.Context, profileID, orderID string) error {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return errors.New("wrong_profileID")
	}

	objOrderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("wrong_orderID")
	}

	res := r.orderCollection.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id":      objOrderID,
			"owner_id": objProfileID,
		},
		bson.M{
			"$set": bson.M{
				"order_detail.status": "completed",
			},
		},
	)

	if res.Err() != nil {
		return res.Err()
	}

	order := new(serviceorder.Order)

	err = res.Decode(&order)
	if err != nil {
		return err
	}

	buyer, _ := r.GetVOfficerServiceOrderByID(ctx, order.GetReferalID())

	if buyer != nil {
		return r.AcceptDeliverdServiceOrder(ctx, buyer.GetOwnerID(), buyer.GetID())
	}

	return nil
}

// CancelDeliverdServiceOrder ...
func (r *Repository) CancelDeliverdServiceOrder(ctx context.Context, profileID, orderID string) error {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return errors.New("wrong_profileID")
	}

	objOrderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("wrong_orderID")
	}

	res := r.orderCollection.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id":      objOrderID,
			"owner_id": objProfileID,
		},
		bson.M{
			"$set": bson.M{
				"order_detail.status": "in_progress",
			},
		},
	)

	if res.Err() != nil {
		return res.Err()
	}

	order := new(serviceorder.Order)

	err = res.Decode(&order)
	if err != nil {
		return err
	}

	buyer, _ := r.GetVOfficerServiceOrderByID(ctx, order.GetReferalID())

	if buyer != nil {
		return r.CancelDeliverdServiceOrder(ctx, buyer.GetOwnerID(), buyer.GetID())
	}

	return nil
}

// GetVOfficerServiceOrders ...
func (r *Repository) GetVOfficerServiceOrders(ctx context.Context, ownerID string, officeID string,
	orderType serviceorder.OrderType, orderstatus serviceorder.OrderStatus, first int, after int) (*serviceorder.GetOrder, error) {

	objOwnerID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, err
	}

	match := bson.M{
		"owner_id":   objOwnerID,
		"order_type": orderType,
	}

	if officeID != "" {
		objOfficeID, err := primitive.ObjectIDFromHex(officeID)
		if err != nil {
			return nil, err
		}

		match["office_id"] = objOfficeID
	}

	if orderstatus != "any" {
		match["order_detail.status"] = orderstatus
	}

	cursor, err := r.orderCollection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"orders": bson.M{
						"$push": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"order_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$orders"},
							bson.M{"$size": "$orders"},
							0,
						},
					},
				},
			},
			{

				"$unwind": bson.M{
					"path": "$orders",
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"orders": bson.M{
						"$push": "$orders",
					},
					"order_amount": bson.M{
						"$first": "$order_amount",
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	order := &serviceorder.GetOrder{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&order)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("not_found")
			}
			return nil, err
		}
	}

	return order, nil
}

// GetVOfficerServiceOrderByID ...
func (r *Repository) GetVOfficerServiceOrderByID(ctx context.Context, orderID string) (*serviceorder.Order, error) {

	objOrderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}

	result := r.orderCollection.FindOne(
		ctx,
		bson.M{
			"_id": objOrderID,
		},
	)

	if result.Err() != nil {
		return nil, result.Err()
	}

	order := new(serviceorder.Order)

	err = result.Decode(&order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetReceivedProposals ...
func (r *Repository) GetReceivedProposals(ctx context.Context, profileID, requestID string, first int, after int) (*offer.GetProposal, error) {
	objOwnerID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, err
	}

	match := bson.M{
		"owner_id": objOwnerID,
	}

	if requestID != "" {
		objRequestID, err := primitive.ObjectIDFromHex(requestID)
		if err != nil {
			return nil, err
		}

		match["request_id"] = objRequestID
	}

	cursor, err := r.offerCollection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"proposals": bson.M{
						"$push": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"proposal_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$proposals"},
							bson.M{"$size": "$proposals"},
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

	proposal := &offer.GetProposal{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&proposal)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("not_found")
			}
			return nil, err
		}
	}

	return proposal, nil
}

// GetSendedProposals ...
func (r *Repository) GetSendedProposals(ctx context.Context, profileID string, first int, after int) (*offer.GetProposal, error) {

	objOwnerID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.offerCollection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"proposal_detail.profile_id": objOwnerID,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"proposals": bson.M{
						"$push": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"proposal_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$proposals"},
							bson.M{"$size": "$proposals"},
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

	proposal := &offer.GetProposal{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&proposal)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("not_found")
			}
			return nil, err
		}
	}

	return proposal, nil
}

// GetProposalAmount ...
func (r *Repository) GetProposalAmount(ctx context.Context, requestID string) (int32, error) {
	objRequestID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		return 0, err
	}

	cursor, err := r.offerCollection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"request_id": objRequestID,
				},
			},
			{
				"$count": "count",
			},
		},
	)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	amount := struct {
		Count int32 `bson:"count"`
	}{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&amount)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return 0, errors.New("not_found")
			}
			return 0, err
		}
	}

	return amount.Count, nil
}

// SaveVOfficeService ...
func (r *Repository) SaveVOfficeService(ctx context.Context, profileID, serviceID string) error {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return err
	}

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return err
	}

	_, err = r.savedServicesCollection.UpdateOne(ctx,
		bson.M{
			"_id": objProfileID,
		},
		bson.M{
			"$addToSet": bson.M{
				"saved_services": objServiceID,
			},
		},

		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// UnSaveVOfficeService ...
func (r *Repository) UnSaveVOfficeService(ctx context.Context, profileID, serviceID string) error {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return err
	}

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return err
	}

	_, err = r.savedServicesCollection.UpdateOne(ctx,
		bson.M{
			"_id": objProfileID,
		},
		bson.M{
			"$pull": bson.M{
				"saved_services": objServiceID,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveServiceRequest ...
func (r *Repository) SaveServiceRequest(ctx context.Context, profileID, serviceID string) error {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return err
	}

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return err
	}

	_, err = r.savedServicesCollection.UpdateOne(ctx,
		bson.M{
			"_id": objProfileID,
		},
		bson.M{
			"$addToSet": bson.M{
				"saved_services_request": objServiceID,
			},
		},

		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// UnSaveServiceRequest ...
func (r *Repository) UnSaveServiceRequest(ctx context.Context, profileID, serviceID string) error {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return err
	}

	objServiceID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return err
	}

	_, err = r.savedServicesCollection.UpdateOne(ctx,
		bson.M{
			"_id": objProfileID,
		},
		bson.M{
			"$pull": bson.M{
				"saved_services_request": objServiceID,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedVOfficeServices ...
func (r *Repository) GetSavedVOfficeServices(ctx context.Context, profileID string, first int, after int) ([]string, int32, error) {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.savedServicesCollection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": bson.M{"_id": objProfileID},
			},
			{
				"$addFields": bson.M{
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$saved_services"},
							bson.M{"$size": "$saved_services"},
							0,
						},
					},
				},
			},
			{
				"$unwind": bson.M{
					"path": "$saved_services",
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"services": bson.M{
						"$push": "$saved_services",
					},
					"amount": bson.M{
						"$first": "$amount",
					},
				},
			},
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	services := struct {
		SavedServices []primitive.ObjectID `bson:"services"`
		Amount        int32                `bson:"amount"`
	}{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&services)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, 0, errors.New("not_found")
			}
			return nil, 0, err
		}
	}

	res := make([]string, 0, len(services.SavedServices))

	for _, s := range services.SavedServices {
		res = append(res, s.Hex())
	}

	return res, services.Amount, nil
}

// GetSavedServicesRequest ...
func (r *Repository) GetSavedServicesRequest(ctx context.Context, profileID string, first int, after int) ([]string, int32, error) {
	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.savedServicesCollection.Aggregate(ctx,
		[]bson.M{
			{
				"$match": bson.M{"_id": objProfileID},
			},
			{
				"$addFields": bson.M{
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$saved_services_request"},
							bson.M{"$size": "$saved_services_request"},
							0,
						},
					},
				},
			},
			{
				"$unwind": bson.M{
					"path": "$saved_services_request",
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"services": bson.M{
						"$push": "$saved_services_request",
					},
					"amount": bson.M{
						"$first": "$amount",
					},
				},
			},
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	services := struct {
		SavedServices []primitive.ObjectID `bson:"services"`
		Amount        int32                `bson:"amount"`
	}{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&services)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, 0, errors.New("not_found")
			}
			return nil, 0, err
		}
	}

	res := make([]string, 0, len(services.SavedServices))

	for _, s := range services.SavedServices {
		res = append(res, s.Hex())
	}

	return res, services.Amount, nil
}

// HasLikedService ...
func (r *Repository) HasLikedService(ctx context.Context, profileID string, serviceID string, requestID string) (bool, error) {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return false, err
	}

	match := bson.M{
		"_id": objProfileID,
	}

	// Get Service ID
	if serviceID != "" {
		objServiceID, err := primitive.ObjectIDFromHex(serviceID)
		if err != nil {
			return false, err
		}

		match["saved_services"] = objServiceID
	}

	// Get Request ID
	if requestID != "" {
		objRequestID, err := primitive.ObjectIDFromHex(requestID)
		if err != nil {
			return false, err
		}

		match["saved_services_request"] = objRequestID
	}

	result := r.savedServicesCollection.FindOne(
		ctx,
		match,
	)

	if result.Err() != nil {
		return false, result.Err()
	}

	res := struct {
		ID primitive.ObjectID `bson:"_id"`
	}{}

	err = result.Decode(&res)
	if err != nil {
		return false, err
	}

	return true, nil

}

// HasOrderedService ...
func (r *Repository) HasOrderedService(ctx context.Context, profileID string, serviceID string, requestID string) (bool, error) {

	objProfileID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return false, err
	}

	match := bson.M{
		"order_detail.profile_id": objProfileID,
	}

	// Get Service ID
	if serviceID != "" {
		objServiceID, err := primitive.ObjectIDFromHex(serviceID)
		if err != nil {
			return false, err
		}

		match["service_id"] = objServiceID
	}

	// Get Request ID
	if requestID != "" {
		objRequestID, err := primitive.ObjectIDFromHex(requestID)
		if err != nil {
			return false, err
		}

		match["request_id"] = objRequestID
	}

	result := r.orderCollection.FindOne(
		ctx,
		match,
	)

	if result.Err() != nil {
		return false, result.Err()
	}

	res := struct {
		ID primitive.ObjectID `bson:"_id"`
	}{}

	err = result.Decode(&res)
	if err != nil {
		return false, err
	}

	return true, nil

}

// AddNoteForOrderService ...
func (r *Repository) AddNoteForOrderService(ctx context.Context, orderID, profileID, text string) error {
	objOrdertID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("wrong_orderID")
	}

	objProfleID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return errors.New("wrong_orderID")
	}

	_, err = r.orderCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":      objOrdertID,
			"owner_id": objProfleID,
		},
		bson.M{
			"$set": bson.M{
				"order_detail.note": text,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// WriteReviewForService ...
func (r *Repository) WriteReviewForService(ctx context.Context, review review.Review) error {

	_, err := r.reviewCollections.InsertOne(ctx, review)

	if err != nil {
		return err
	}

	return nil

}

// WriteReviewForServiceRequest ...
func (r *Repository) WriteReviewForServiceRequest(ctx context.Context, review review.Review) error {

	_, err := r.reviewCollections.InsertOne(ctx, review)

	if err != nil {
		return err
	}

	return nil

}

// GetServicesReview ...
func (r *Repository) GetServicesReview(ctx context.Context, profileID, officeID string, first int, after int) (*review.GetReview, error) {

	match := bson.M{}

	if profileID != "" {
		objOwnerID, err := primitive.ObjectIDFromHex(profileID)
		if err != nil {
			return nil, err
		}
		match["owner_id"] = objOwnerID
	}

	if officeID != "" {
		objOfficeID, err := primitive.ObjectIDFromHex(officeID)
		if err != nil {
			return nil, err
		}

		match["office_id"] = objOfficeID
	}

	cursor, err := r.reviewCollections.Aggregate(ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$addFields": bson.M{
					"review_avg": bson.M{
						"$avg": []interface{}{
							"$review_detail.clarity",
							"$review_detail.communication",
							"$review_detail.payment",
						},
					},
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"reviews": bson.M{
						"$push": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"review_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$reviews"},
							bson.M{"$size": "$reviews"},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"reviews": bson.M{
						"$slice": []interface{}{
							"$reviews",
							after,
							first,
						},
					},
					"review_amount": 1,
					"clartity_avg": bson.M{
						"$avg": []interface{}{
							"$reviews.review_detail.clarity",
						},
					},
					"communication_avg": bson.M{
						"$avg": []interface{}{
							"$reviews.review_detail.communication",
						},
					},
					"payment_avg": bson.M{
						"$avg": []interface{}{
							"$reviews.review_detail.payment",
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

	review := &review.GetReview{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&review)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("not_found")
			}
			return nil, err
		}
	}

	return review, nil
}
