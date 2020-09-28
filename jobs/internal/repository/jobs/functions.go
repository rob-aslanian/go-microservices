package jobsrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	careercenter "gitlab.lan/Rightnao-site/microservices/jobs/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetCandidateProfile recieves userID as string, transforms it into mongodb ObjectID, searches by that ID and returns that user's profile.
func (r *Repository) GetCandidateProfile(ctx context.Context, userID string) (*candidate.Profile, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	var profile candidate.Profile

	err = r.profileCollection.FindOne(
		ctx,
		bson.M{
			"user_id": objID,
		},
	).Decode(&profile)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}

	return &profile, nil
}

// UpsertCandidateProfile changes ``is_open``, ``career_interests`` values.
func (r *Repository) UpsertCandidateProfile(ctx context.Context, profile *candidate.Profile) error {
	objID, err := primitive.ObjectIDFromHex(profile.GetUserID())
	if err != nil {
		return errors.New(`wrong_id`)
	}

	r.profileCollection.UpdateOne(
		ctx,
		bson.M{
			"user_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"is_open":          profile.IsOpen,
				"career_interests": profile.CareerInterests,
			},
		},
		options.Update().SetUpsert(true),
	)

	return nil
}

// SetOpenFlag changes ``is_open`` value
func (r *Repository) SetOpenFlag(ctx context.Context, userID string, flag bool) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.profileCollection.UpdateOne(
		ctx,
		bson.M{
			"user_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"is_open": flag,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// PostJob adds another job post
func (r *Repository) PostJob(ctx context.Context, post *job.Posting) error {
	_, err := r.jobsCollection.InsertOne(ctx, post)
	if err != nil {
		return err
	}

	return nil
}

// UpdateJobPosting ...
func (r *Repository) UpdateJobPosting(ctx context.Context, post *job.Posting) error {
	objID, err := primitive.ObjectIDFromHex(post.GetID())
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": post,
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// UpdateJobDetails ...
func (r *Repository) UpdateJobDetails(ctx context.Context, post *job.Posting) error {
	objID, err := primitive.ObjectIDFromHex(post.GetID())
	if err != nil {
		return errors.New(`wrong_id`)
	}

	log.Println("details:", objID, post.JobDetails)

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"job_details": post.JobDetails,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// DeleteExpiredPost ...
func (r Repository) DeleteExpiredPost(ctx context.Context, postID, companyID string) error {
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return err
	}
	_, err = r.jobsCollection.DeleteOne(
		ctx,
		bson.M{
			"_id":        postObjID,
			"company_id": companyObjID,
			"status":     "Expired",
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetJobPosting ...
func (r *Repository) GetJobPosting(ctx context.Context, jobID string, revertApplicant bool) (*job.Posting, error) {
	objID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	var result job.Posting

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": objID,
			},
		},
	}

	if revertApplicant {
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"applications": bson.M{
					"$reverseArray": "$applications",
				},
			},
		})
	}

	cursor, err := r.jobsCollection.Aggregate(ctx, pipeline)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

// GetDraft ...
func (r *Repository) GetDraft(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error) {
	objID, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	var draft job.ViewForCompany

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id":        objID,
					"company_id": objCompanyID,
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		cursor.Decode(&draft)
	}

	return &draft, nil
}

// GetPost ...
func (r *Repository) GetPost(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error) {
	objID, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	var draft job.ViewForCompany

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id":        objID,
					"company_id": objCompanyID,
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		cursor.Decode(&draft)
	}

	return &draft, nil
}

// GetIDsOfApplicants ...
func (r *Repository) GetIDsOfApplicants(ctx context.Context, jobID string, reverse bool) ([]string, error) {
	objID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": objID,
			},
		},
	}

	if reverse {
		pipeline = append(pipeline, bson.M{
			"$project": bson.M{
				"applications": bson.M{
					"$reverseArray": "$applications",
				},
			},
		})
	}

	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"ids": "$applications.user_id",
		},
	})

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		pipeline,
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	objIDs := struct {
		IDs []primitive.ObjectID `bson:"ids"`
	}{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&objIDs)
		if err != nil {
			return nil, err
		}
	}

	ids := make([]string, 0, len(objIDs.IDs))
	for _, id := range objIDs.IDs {
		ids = append(ids, id.Hex())
	}

	return ids, nil
}

// GetJobForCompany ...
func (r *Repository) GetJobForCompany(ctx context.Context, companyID, jobID string) (*job.ViewForCompany, error) {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	var post job.ViewForCompany

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id":        objJobID,
					"company_id": objCompanyID,
				},
			},
			{
				"$addFields": bson.M{
					"num_of_applications": mongoCondition("$applications", "$isArray", "$size", 0),
					"num_of_views":        mongoCondition("$views", "$isArray", "$size", 0),
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		cursor.Decode(&post)
	}

	return &post, nil
}

// GetListOfJobsWithSeenStat ...
func (r *Repository) GetListOfJobsWithSeenStat(ctx context.Context, companyID string, first, after int32) ([]*job.ViewJobWithSeenStat, error) {
	objID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objID,
				},
			},
			{
				"$sort": bson.M{
					"activation_date": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$project": bson.M{
					"title": "$job_details.title",
					"total_amount": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$isArray": "$applications",
							},
							bson.M{
								"$size": "$applications",
							},
							0,
						},
					},
					"unseen_amount": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$isArray": "$applications",
							},
							bson.M{
								"$size": bson.M{
									"$filter": bson.M{
										"input": "$applications",
										"as":    "tmp",
										"cond": bson.M{
											"$eq": []interface{}{"$$tmp.metadata.seen", false},
										},
									},
								},
							},
							0,
						},
					},
					"status": "$status",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	jobs := make([]*job.ViewJobWithSeenStat, 0)
	for cursor.Next(ctx) {
		var val job.ViewJobWithSeenStat
		err = cursor.Decode(&val)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &val)
	}

	return jobs, nil
}

// IsJobApplied ...
func (r *Repository) IsJobApplied(ctx context.Context, userID, jobID string) (bool, error) {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	var isExist bool

	cursor, err := r.jobsCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"_id":                  objJobID,
				"applications.user_id": objUserID,
			},
		},
	})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		isExist = true
	}

	return isExist, nil
}

// IgnoreInvitation ...
func (r *Repository) IgnoreInvitation(ctx context.Context, userID, jobID string) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objJobID,
		},
		bson.M{
			"$pull": bson.M{
				"invited_candidates.user_id": objUserID,
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// IsJobSaved ...
func (r *Repository) IsJobSaved(ctx context.Context, userID, jobID string) (bool, error) {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	var isExist bool

	cursor, err := r.profileCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"user_id":    objUserID,
				"saved_jobs": objJobID,
			},
		},
	})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		isExist = true
	}

	return isExist, nil
}

// GetCompanyIDByJobID ...
func (r *Repository) GetCompanyIDByJobID(ctx context.Context, jobID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return "", errors.New(`wrong_id`)
	}

	cursor, err := r.jobsCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"_id": objID,
			},
		},
		{
			"$project": bson.M{
				"company_id": 1,
			},
		},
	})
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	id := struct {
		ID primitive.ObjectID `bson:"company_id"`
	}{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&id)
		if err != nil {
			return "", err
		}
	}

	return id.ID.Hex(), nil
}

// ApplyJob ...
func (r *Repository) ApplyJob(ctx context.Context, userID, jobID string, application *job.Application, fileIDs []string) error {
	objID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	files, err := r.getUnuploadedFilesForApplication(ctx, userID, jobID, fileIDs)
	if err != nil {
		return err
	}
	application.Documents = files

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$push": bson.M{
				"applications": application,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	err = r.removeUnuploadedFilesForApplication(ctx, userID, jobID)
	if err != nil {
		log.Println("error: removeUnuploadedFilesForApplication", err)
	}

	return nil
}

// AddJobView ...
func (r *Repository) AddJobView(ctx context.Context, userID, jobID string) error {
	objID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$addToSet": bson.M{
				"views": objID,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// GetRecommendedJobs ...
func (r *Repository) GetRecommendedJobs(ctx context.Context, userID, countryID string, first, after int32) ([]*job.ViewForUser, error) {
	rec := make([]*job.ViewForUser, 0)

	var ids, idsSaved []string

	if userID != "" {
		// skipped jobs
		var err error

		ids, err = r.GetIDOfSkippedJobs(ctx, userID)
		if err != nil {
			return nil, err
		}

		// saved jobs
		idsSaved, err = r.GetIDOfSavedJobs(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	objSavedIDs := make([]primitive.ObjectID, 0, len(idsSaved))

	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, objID)
		} else {
			log.Println("can't convert id into objectID:", err)
		}
	}

	for _, id := range idsSaved {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objSavedIDs = append(objSavedIDs, objID)
		} else {
			log.Println("can't convert id into objectID:", err)
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"status":                               "Active",
				"job_metadata.advertisement_countries": countryID,
				"_id": bson.M{
					"$nin": objIDs,
				},
			},
		},
		// {
		// 	"$skip": after,
		// },
		// {
		// 	"$limit": first,
		// },
		{
			"$sample": bson.M{
				"size": first,
			},
		},
		{
			"$addFields": bson.M{
				"is_saved": bson.M{
					"$cond": []interface{}{
						bson.M{
							"$in": []interface{}{
								"$_id",
								objSavedIDs,
							},
						},
						true,
						false,
					},
				},
			},
		},
	}

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var r job.ViewForUser

		err = cursor.Decode(&r)
		if err != nil {
			return nil, err
		}

		rec = append(rec, &r)
	}

	log.Println("GetRecommendedJobs: len res:", len(rec), "nin: ", objIDs)

	return rec, nil
}

// SaveJob ...
func (r *Repository) SaveJob(ctx context.Context, userID, jobID string) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.profileCollection.UpdateOne(
		ctx,
		bson.M{
			"user_id": objUserID,
		},
		bson.M{
			"$addToSet": bson.M{
				"saved_jobs": objJobID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// UnsaveJob ...
func (r *Repository) UnsaveJob(ctx context.Context, userID, jobID string) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.profileCollection.UpdateOne(
		ctx,
		bson.M{
			"user_id": objUserID,
		},
		bson.M{
			"$pull": bson.M{
				"saved_jobs": objJobID,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SkipJob ...
func (r *Repository) SkipJob(ctx context.Context, userID, jobID string) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.profileCollection.UpdateOne(
		ctx,
		bson.M{
			"user_id": objUserID,
		},
		bson.M{
			"$addToSet": bson.M{
				"skipped_jobs": objJobID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// UnskipJob ...
func (r *Repository) UnskipJob(ctx context.Context, userID, jobID string) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.profileCollection.UpdateOne(
		ctx,
		bson.M{
			"user_id": objUserID,
		},
		bson.M{
			"$pull": bson.M{
				"skipped_jobs": objJobID,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedJobs ...
func (r *Repository) GetSavedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	jobIDs := struct {
		IDs []primitive.ObjectID `bson:"saved_jobs"`
	}{}

	cursor, err := r.profileCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"user_id": objUserID,
			},
		},
		{
			"$skip": after,
		},
		{
			"$limit": first,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&jobIDs)
		if err != nil {
			return nil, err
		}
	}

	if len(jobIDs.IDs) == 0 {
		return []*job.ViewForUser{}, nil
	}

	// getting jobs
	cursor2, err := r.jobsCollection.Find(
		ctx,
		bson.M{
			"_id": bson.M{
				"$in": jobIDs.IDs,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor2.Close(ctx)

	views := make([]*job.ViewForUser, 0, len(jobIDs.IDs))

	for cursor2.Next(ctx) {
		j := job.ViewForUser{}

		err = cursor2.Decode(&j)
		if err != nil {
			return nil, err
		}

		views = append(views, &j)
	}

	return views, nil
}

// GetSkippedJobs ...
func (r *Repository) GetSkippedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	jobIDs := struct {
		IDs []primitive.ObjectID `bson:"skipped_jobs"`
	}{}

	cursor, err := r.profileCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"user_id": objUserID,
			},
		},
		{
			"$skip": after,
		},
		{
			"$limit": first,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err = cursor.Decode(&jobIDs)
		if err != nil {
			return nil, err
		}
	}

	if len(jobIDs.IDs) == 0 {
		return []*job.ViewForUser{}, nil
	}

	// getting jobs
	cursor2, err := r.jobsCollection.Find(
		ctx,
		bson.M{
			"_id": bson.M{
				"$in": jobIDs.IDs,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor2.Close(ctx)

	views := make([]*job.ViewForUser, 0, len(jobIDs.IDs))

	for cursor2.Next(ctx) {
		j := job.ViewForUser{}

		err = cursor2.Decode(&j)
		if err != nil {
			return nil, err
		}

		views = append(views, &j)
	}

	return views, nil
}

// GetAppliedJobs ...
func (r *Repository) GetAppliedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	// saved jobs
	idsSaved, err := r.GetIDOfSavedJobs(ctx, userID)
	if err != nil {
		return nil, err
	}

	objSavedIDs := make([]primitive.ObjectID, 0, len(idsSaved))

	for _, id := range idsSaved {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objSavedIDs = append(objSavedIDs, objID)
		} else {
			log.Println("can't convert id into objectID:", err)
		}
	}

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"applications.user_id": objID,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$addFields": bson.M{
					"application": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$isArray": "$applications",
							},
							bson.M{
								"$arrayElemAt": []interface{}{
									bson.M{
										"$filter": bson.M{
											"input": "$applications",
											"as":    "a",
											"cond": bson.M{
												"$eq": []interface{}{
													"$$a.user_id", objID,
												},
											},
										},
									},
									0,
								},
							},
							0,
						},
					},
					"is_saved": bson.M{
						"$in": []interface{}{
							"$_id",
							objSavedIDs,
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

	views := make([]*job.ViewForUser, 0)
	for cursor.Next(ctx) {
		j := job.ViewForUser{}

		err = cursor.Decode(&j)
		if err != nil {
			return nil, err
		}

		views = append(views, &j)
	}

	return views, nil
}

// SetJobApplicationSeen ...
func (r *Repository) SetJobApplicationSeen(ctx context.Context, jobID, applicantID string, seen bool) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objApplicantID, err := primitive.ObjectIDFromHex(applicantID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":                  objJobID,
			"applications.user_id": objApplicantID,
		},
		bson.M{
			"$set": bson.M{
				"applications.$.metadata.seen": seen,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// SetJobApplicationCategory ...
func (r *Repository) SetJobApplicationCategory(ctx context.Context, jobID, applicantID string, category job.ApplicantCategory) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objApplicantID, err := primitive.ObjectIDFromHex(applicantID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":                  objJobID,
			"applications.user_id": objApplicantID,
		},
		bson.M{
			"$set": bson.M{
				"applications.$.metadata.category": category,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// GetAmountOfApplicantsPerCategory ...
func (r *Repository) GetAmountOfApplicantsPerCategory(ctx context.Context, companyID string) (total, unseen, favorite, inReview, disqualified int32, err error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		err = errors.New(`wrong_id`)
		return
	}

	counts := struct {
		Total        int32 `bson:"total_amount"`
		Unseen       int32 `bson:"unseen_amount"`
		Favorite     int32 `bson:"favorite_amount"`
		InReview     int32 `bson:"in_review_amount"`
		Disqualified int32 `bson:"disqualified_amount"`
	}{}

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objCompanyID,
				},
			},
			{
				"$project": bson.M{
					"total_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$applications"},
							bson.M{"$size": "$applications"},
							0,
						},
					},

					"unseen_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$applications"},
							bson.M{
								"$size": bson.M{
									"$filter": bson.M{
										"input": "$applications",
										"as":    "tmp",
										"cond":  bson.M{"$eq": []interface{}{"$$tmp.metadata.seen", false}},
									},
								},
							},
							0,
						},
					},

					"favorite_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$applications"},
							bson.M{
								"$size": bson.M{
									"$filter": bson.M{
										"input": "$applications",
										"as":    "tmp",
										"cond":  bson.M{"$eq": []interface{}{"$$tmp.metadata.category", "favorite"}},
									},
								},
							},
							0,
						},
					},

					"in_review_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$applications"},
							bson.M{
								"$size": bson.M{
									"$filter": bson.M{
										"input": "$applications",
										"as":    "tmp",
										"cond":  bson.M{"$eq": []interface{}{"$$tmp.metadata.category", "in_review"}},
									},
								},
							},
							0,
						},
					},

					"disqualified_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$applications"},
							bson.M{
								"$size": bson.M{
									"$filter": bson.M{
										"input": "$applications",
										"as":    "tmp",
										"cond":  bson.M{"$eq": []interface{}{"$$tmp.metadata.category", "disqualified"}},
									},
								},
							},
							0,
						},
					},
				},
			},
		},
	)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&counts)
		if err != nil {
			return
		}
	}

	return counts.Total, counts.Unseen, counts.Favorite, counts.InReview, counts.Disqualified, nil
}

// GetCandidates ...
func (r *Repository) GetCandidates(ctx context.Context, companyID string, country string, first, after int32) ([]*candidate.ViewForCompany, error) {

	// skipped
	ids, err := r.GetIDOfSkippedCandidates(ctx, companyID)
	if err != nil {
		return nil, err
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))

	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, objID)
		} else {
			log.Println("can't convert id into objectID:", err)
		}
	}

	// saved
	idsSaved, err := r.GetIDOfSavedCandidates(ctx, companyID)
	if err != nil {
		return nil, err
	}

	objSavedIDs := make([]primitive.ObjectID, 0, len(idsSaved))

	for _, id := range idsSaved {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objSavedIDs = append(objSavedIDs, objID)
		} else {
			log.Println("can't convert id into objectID:", err)
		}
	}

	cursor, err := r.profileCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"is_open":                    true,
					"career_interests.countries": country,
					"user_id": bson.M{
						"$nin": objIDs,
					},
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$addFields": bson.M{
					"is_saved": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$in": []interface{}{
									"$user_id",
									objSavedIDs,
								},
							},
							true,
							false,
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

	views := make([]*candidate.ViewForCompany, 0)
	for cursor.Next(ctx) {
		j := candidate.ViewForCompany{}

		err = cursor.Decode(&j)
		if err != nil {
			return nil, err
		}

		views = append(views, &j)
	}

	return views, nil
}

// SaveCandidate ...
func (r *Repository) SaveCandidate(ctx context.Context, companyID, candidateID string) error {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCandidateID, err := primitive.ObjectIDFromHex(candidateID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.companiesCollection.UpdateOne(
		ctx,
		bson.M{
			"company_id": objCompanyID,
		},
		bson.M{
			"$addToSet": bson.M{
				"saved_candidates": objCandidateID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// UnsaveCandidate ...
func (r *Repository) UnsaveCandidate(ctx context.Context, companyID, candidateID string) error {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCandidateID, err := primitive.ObjectIDFromHex(candidateID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.companiesCollection.UpdateOne(
		ctx,
		bson.M{
			"company_id": objCompanyID,
		},
		bson.M{
			"$pull": bson.M{
				"saved_candidates": objCandidateID,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// GetIDOfSavedCandidates ...
func (r *Repository) GetIDOfSavedCandidates(ctx context.Context, companyID string) ([]string, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objIDs := make([]primitive.ObjectID, 0)

	// get ids
	cursor, err := r.companiesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objCompanyID,
				},
			},
			{
				"$project": bson.M{
					"saved_candidates": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$saved_candidates",
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		i := struct {
			ID primitive.ObjectID `bson:"saved_candidates"`
		}{}

		cursor.Decode(&i)
		objIDs = append(objIDs, i.ID)
	}

	ids := make([]string, 0, len(objIDs))
	for _, o := range objIDs {
		ids = append(ids, o.Hex())
	}

	return ids, nil
}

// SkipCandidate ...
func (r *Repository) SkipCandidate(ctx context.Context, companyID, candidateID string) error {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCandidateID, err := primitive.ObjectIDFromHex(candidateID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.companiesCollection.UpdateOne(
		ctx,
		bson.M{
			"company_id": objCompanyID,
		},
		bson.M{
			"$addToSet": bson.M{
				"skipped_candidates": objCandidateID,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// UnskipCandidate ...
func (r *Repository) UnskipCandidate(ctx context.Context, companyID, candidateID string) error {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCandidateID, err := primitive.ObjectIDFromHex(candidateID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.companiesCollection.UpdateOne(
		ctx,
		bson.M{
			"company_id": objCompanyID,
		},
		bson.M{
			"$pull": bson.M{
				"skipped_candidates": objCandidateID,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(`not_found`)
		}
		return err
	}

	return nil
}

// GetAmountsOfManageCandidates returns amount of saved candidates, amount of skipped candidates
func (r *Repository) GetAmountsOfManageCandidates(ctx context.Context, companyID string) (int32, int32, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return 0, 0, errors.New(`wrong_id`)
	}

	amount := struct {
		Saved   int32 `bson:"saved_candidates"`
		Skipped int32 `bson:"skipped_candidates"`
		// Alerts  int32 `bson:"alerts"`
	}{}

	cursor, err := r.companiesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objCompanyID,
				},
			},
			{
				"$addFields": bson.M{
					"saved_candidates": bson.M{
						"$size": "$saved_candidates",
					},
					"skipped_candidates": bson.M{
						"$size": "$skipped_candidates",
					},
					// "alerts": bson.M{
					// 	"$size": "$alerts", // TODO:
					// },
				},
			},
			{
				"$project": bson.M{
					"saved_candidates":   1,
					"skipped_candidates": 1,
					// "alerts":             1,
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, 0, errors.New(`not_found`)
		}
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&amount)
		if err != nil {
			return 0, 0, err
		}
	}

	return amount.Saved, amount.Skipped, nil
}

// GetSavedCandidates ...
func (r *Repository) GetSavedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	ids := make([]primitive.ObjectID, 0)
	candidates := make([]*candidate.ViewForCompany, 0)

	// get ids
	cursor, err := r.companiesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objCompanyID,
				},
			},
			{
				"$project": bson.M{
					"saved_candidates": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$saved_candidates",
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		i := struct {
			ID primitive.ObjectID `bson:"saved_candidates"`
		}{}
		cursor.Decode(&i)
		ids = append(ids, i.ID)
	}

	cursor2, err := r.profileCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"user_id": bson.M{
						"$in": ids,
					},
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor2.Close(ctx)

	for cursor2.Next(ctx) {
		c := candidate.ViewForCompany{}
		cursor2.Decode(&c)
		candidates = append(candidates, &c)
	}

	return candidates, nil
}

// GetIDOfSkippedCandidates ...
func (r *Repository) GetIDOfSkippedCandidates(ctx context.Context, companyID string) ([]string, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objIDs := make([]primitive.ObjectID, 0)

	// get ids
	cursor, err := r.companiesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objCompanyID,
				},
			},
			{
				"$project": bson.M{
					"skipped_candidates": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$skipped_candidates",
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		i := struct {
			ID primitive.ObjectID `bson:"skipped_candidates"`
		}{}

		cursor.Decode(&i)
		objIDs = append(objIDs, i.ID)
	}

	ids := make([]string, 0, len(objIDs))
	for _, o := range objIDs {
		ids = append(ids, o.Hex())
	}

	return ids, nil
}

// GetSkippedCandidates ...
func (r *Repository) GetSkippedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	ids := make([]primitive.ObjectID, 0)
	candidates := make([]*candidate.ViewForCompany, 0)

	// get ids
	cursor, err := r.companiesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objCompanyID,
				},
			},
			{
				"$project": bson.M{
					"skipped_candidates": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$skipped_candidates",
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		i := struct {
			ID primitive.ObjectID `bson:"skipped_candidates"`
		}{}
		cursor.Decode(&i)
		ids = append(ids, i.ID)
	}

	cursor2, err := r.profileCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"user_id": bson.M{
						"$in": ids,
					},
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor2.Close(ctx)

	for cursor2.Next(ctx) {
		c := candidate.ViewForCompany{}
		cursor2.Decode(&c)
		candidates = append(candidates, &c)
	}

	return candidates, nil
}

// GetPostedJobs ...
func (r *Repository) GetPostedJobs(ctx context.Context, companyID string) ([]*job.ViewForCompany, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	rec := make([]*job.ViewForCompany, 0)

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{"company_id": objCompanyID},
			},
			{
				"$addFields": bson.M{
					"num_of_applications": mongoCondition("$applications", "$isArray", "$size", 0),
					"num_of_views":        mongoCondition("$views", "$isArray", "$size", 0),
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var r job.ViewForCompany

		err = cursor.Decode(&r)
		if err != nil {
			return nil, err
		}

		rec = append(rec, &r)
	}

	return rec, nil
}

// InviteUserToApply ...
func (r *Repository) InviteUserToApply(ctx context.Context, jobID string, invitation job.Invitation) error {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = primitive.ObjectIDFromHex(invitation.GetUserID())
	if err != nil {
		return errors.New(`wrong_id`)
	}

	result, err := r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objJobID,
		},
		bson.M{
			"$addToSet": bson.M{
				"invited_candidates": invitation,
			},
		},
	)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("already_invited")
	}

	return nil
}

// GetInvitedJobs ...
func (r *Repository) GetInvitedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error) {

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	// // saved jobs
	// idsSaved, err := r.GetIDOfSavedJobs(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// objSavedIDs := make([]primitive.ObjectID, 0, len(idsSaved))
	//
	// for _, id := range idsSaved {
	// 	objID, err := primitive.ObjectIDFromHex(id)
	// 	if err == nil {
	// 		objSavedIDs = append(objSavedIDs, objID)
	// 	} else {
	// 		log.Println("can't convert id into objectID:", err)
	// 	}
	// }

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"invited_candidates.user_id": objUserID,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
			{
				"$addFields": bson.M{
					"application": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$isArray": "$applications",
							},
							bson.M{
								"$arrayElemAt": []interface{}{
									bson.M{
										"$filter": bson.M{
											"input": "$applications",
											"as":    "a",
											"cond": bson.M{
												"$eq": []interface{}{
													"$$a.user_id", objUserID,
												},
											},
										},
									},
									0,
								},
							},
							bson.M{},
						},
					},
					"invitation_text": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$isArray": "$invited_candidates",
							},
							bson.M{
								"$arrayElemAt": []interface{}{
									bson.M{
										"$filter": bson.M{
											"input": "$invited_candidates",
											"as":    "a",
											"cond": bson.M{
												"$eq": []interface{}{
													"$$a.user_id", objUserID,
												},
											},
										},
									},
									0,
								},
							},
							"",
						},
					},
					// "is_saved": bson.M{
					// 	"$in": []interface{}{
					// 		"$_id",
					// 		objSavedIDs,
					// 	},
					// },
				},
			},
			{
				"$addFields": bson.M{
					"invitation_text": "$invitation_text.text",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	views := make([]*job.ViewForUser, 0)
	for cursor.Next(ctx) {
		j := job.ViewForUser{}

		err = cursor.Decode(&j)
		if err != nil {
			return nil, err
		}

		views = append(views, &j)
	}

	return views, nil
}

// IsCandidateInvited ...
func (r *Repository) IsCandidateInvited(ctx context.Context, jobID, candidateID string) (bool, error) {
	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}
	objUserID, err := primitive.ObjectIDFromHex(candidateID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	cursor, err := r.jobsCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"_id":                        objJobID,
				"invited_candidates.user_id": objUserID,
			},
		},
	})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	var isInvited bool
	if cursor.Next(ctx) {
		isInvited = true
	}

	return isInvited, nil
}

// ReportJob ...
func (r *Repository) ReportJob(ctx context.Context, report *job.Report) error {
	_, err := r.jobReportsCollection.InsertOne(ctx, report)
	if err != nil {
		return err
	}

	return nil
}

// ReportCandidate ...
func (r *Repository) ReportCandidate(ctx context.Context, report *candidate.Report) error {
	_, err := r.candidateReportsCollection.InsertOne(ctx, report)
	if err != nil {
		return err
	}

	return nil
}

// SaveJobSearchFilter ...
func (r *Repository) SaveJobSearchFilter(ctx context.Context, filter *job.NamedSearchFilter) error {
	_, err := r.jobFiltersCollection.InsertOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedJobSearchFilters ...
func (r *Repository) GetSavedJobSearchFilters(ctx context.Context, userID string) ([]*job.NamedSearchFilter, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.candidateFiltersCollection.Find(
		ctx,
		bson.M{
			"user_id": objID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	filters := make([]*job.NamedSearchFilter, 0)
	for cursor.Next(ctx) {
		j := job.NamedSearchFilter{}

		err = cursor.Decode(&j)
		if err != nil {
			return nil, err
		}

		filters = append(filters, &j)
	}

	return filters, nil
}

// SaveCandidateSearchFilter ...
func (r *Repository) SaveCandidateSearchFilter(ctx context.Context, companyID string, filter *candidate.NamedSearchFilter) error {
	_, err := r.candidateFiltersCollection.InsertOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedCandidateSearchFilters ...
func (r *Repository) GetSavedCandidateSearchFilters(ctx context.Context, companyID string) ([]*candidate.NamedSearchFilter, error) {
	objID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.candidateFiltersCollection.Find(
		ctx,
		bson.M{
			"company_id": objID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	filters := make([]*candidate.NamedSearchFilter, 0)
	for cursor.Next(ctx) {
		j := candidate.NamedSearchFilter{}

		err = cursor.Decode(&j)
		if err != nil {
			return nil, err
		}

		filters = append(filters, &j)
	}

	return filters, nil
}

// GetPlanPrices ...
func (r *Repository) GetPlanPrices(ctx context.Context, countriesID []string, currency string) ([]job.PlanPrices, error) {
	cursor, err := r.pricesCollection.Find(ctx, bson.M{
		"country": bson.M{
			"$in": countriesID,
		},
		"currency": currency,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	plans := make([]job.PlanPrices, 0)

	for cursor.Next(ctx) {
		val := job.PlanPrices{}
		err = cursor.Decode(&val)
		if err != nil {
			return nil, err
		}

		plans = append(plans, val)
	}

	return plans, nil
}

// GetPricingFor ...
func (r *Repository) GetPricingFor(ctx context.Context, meta *job.Meta) (*job.PricingResult, error) {
	anonymCoef := 0
	if meta.Anonymous {
		anonymCoef = 1
	}

	projection := bson.M{
		"country":                      1,
		"currency":                     1,
		"plan_price":                   fmt.Sprint("$amount_of_days.", meta.AmountOfDays),
		"renewal_discount":             bson.M{"$arrayElemAt": []interface{}{"$discounts.renewal", meta.Renewal}},
		"features.language":            bson.M{"$multiply": []interface{}{"$features.language", meta.Renewal + 1}},
		"features.publish_anonymously": bson.M{"$multiply": []interface{}{anonymCoef, "$features.publish_anonymously", meta.Renewal + 1}},
	}

	result := job.PricingResult{}

	cursor, err := r.pricesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"country":  bson.M{"$in": meta.AdvertisementCountries},
					"currency": meta.Currency,
				},
			},
			{
				"$project": projection,
			},
			{
				"$addFields": bson.M{
					"renewal_price": bson.M{
						"$divide": []interface{}{
							bson.M{
								"$multiply": []interface{}{
									"$plan_price", bson.M{
										"$subtract": []interface{}{
											100, "$renewal_discount",
										},
									},
									meta.Renewal,
								},
							}, 100,
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"total_price": bson.M{
						"$add": []interface{}{
							"$plan_price", "$features.publish_anonymously", "$features.language", "$renewal_price",
						},
					},
				},
			},
			{
				"$group": bson.M{
					"_id":          "$currency",
					"currency":     bson.M{"$first": "$currency"},
					"total":        bson.M{"$sum": "$total_price"},
					"by_countries": bson.M{"$push": "$$ROOT"},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

// GetAmountOfActiveJobsOfCompany ...
func (r *Repository) GetAmountOfActiveJobsOfCompany(ctx context.Context, companyID string) (int32, error) {
	objID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return 0, errors.New(`wrong_id`)
	}

	result := struct {
		Amount int32 `bson:"amount"`
	}{}

	cursor, err := r.jobsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objID,
					"status":     "Active",
				},
			},
			{
				"$count": "amount",
			},
		},
	)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		cursor.Decode(&result)
	}

	return result.Amount, nil
}

// GetCareerInterestsByIds ...
func (r *Repository) GetCareerInterestsByIds(ctx context.Context, ids []string, first, after int32) (map[string]*candidate.CareerInterests, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	cursor, err := r.profileCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"user_id": bson.M{
						"$in": objIDs,
					},
				},
			},
			{
				"$project": bson.M{
					"user_id":          1,
					"career_interests": 1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	profiles := make(map[string]*candidate.CareerInterests, len(ids))

	for cursor.Next(ctx) {
		if userID, ok := cursor.Current.Lookup("user_id").ObjectIDOK(); ok {
			var profile candidate.Profile
			err = cursor.Decode(&profile)
			if err != nil {
				return nil, err
			}
			profiles[userID.Hex()] = profile.CareerInterests
		} else {
			log.Println("can't retrive user_id")
			return nil, errors.New("internal_error")
		}
	}

	return profiles, nil
}

// GetSortedListOfIDsOfApplicant ...
func (r *Repository) GetSortedListOfIDsOfApplicant(ctx context.Context, userIDs []string, sort candidate.ApplicantSort, first, after int32) ([]string, error) {
	objIDs := make([]primitive.ObjectID, 0, len(userIDs))
	for _, id := range userIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	IDs := make([]string, 0)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id": bson.M{
					"$in": objIDs,
				},
			},
		},
	}

	if sort == candidate.AppicantExperienceYears {
		pipeline = append(pipeline, bson.M{
			"$sort": bson.M{
				"career_interests.experience": -1,
			},
		})
	}

	pipeline = append(pipeline, []bson.M{
		// {
		// 	"$skip": after,
		// },
		// {
		// 	"$limit": first,
		// },
		{
			"$project": bson.M{
				"_id": 1,
			},
		},
	}...)

	cursor, err := r.profileCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		val := struct {
			ID primitive.ObjectID `bson:"_id"`
		}{}
		err = cursor.Decode(&val)
		if err != nil {
			return nil, err
		}

		IDs = append(IDs, val.ID.Hex())
	}

	return IDs, nil
}

// GetIDOfSkippedJobs ...
func (r *Repository) GetIDOfSkippedJobs(ctx context.Context, userID string) ([]string, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objIDs := make([]primitive.ObjectID, 0)

	// get ids
	cursor, err := r.profileCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"user_id": objUserID,
				},
			},
			{
				"$project": bson.M{
					"skipped_jobs": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$skipped_jobs",
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		i := struct {
			ID primitive.ObjectID `bson:"skipped_jobs"`
		}{}

		cursor.Decode(&i)
		objIDs = append(objIDs, i.ID)
	}

	ids := make([]string, 0, len(objIDs))
	for _, o := range objIDs {
		ids = append(ids, o.Hex())
	}

	return ids, nil
}

// GetIDOfSavedJobs ...
func (r *Repository) GetIDOfSavedJobs(ctx context.Context, userID string) ([]string, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objIDs := make([]primitive.ObjectID, 0)

	// get ids
	cursor, err := r.profileCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"user_id": objUserID,
				},
			},
			{
				"$project": bson.M{
					"saved_jobs": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$saved_jobs",
				},
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(`not_found`)
		}

		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		i := struct {
			ID primitive.ObjectID `bson:"saved_jobs"`
		}{}

		cursor.Decode(&i)
		objIDs = append(objIDs, i.ID)
	}

	ids := make([]string, 0, len(objIDs))
	for _, o := range objIDs {
		ids = append(ids, o.Hex())
	}

	return ids, nil
}

// UploadFileForJob ...
func (r *Repository) UploadFileForJob(ctx context.Context, companyID string, jobID string, file *job.File) error {

	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":        objJobID,
			"company_id": objCompanyID,
		},
		bson.M{
			"$push": bson.M{
				"job_details.files": file,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// UploadFileForApplication ...
func (r *Repository) UploadFileForApplication(ctx context.Context, userID string, jobID string, file *job.File) error {
	// TODO:
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objJobID,
		},
		bson.M{
			"$push": bson.M{
				"unuploaded_files." + userID: file,
			},
		},
		// options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// IsCandidateSaved ...
func (r *Repository) IsCandidateSaved(ctx context.Context, companyID, candidateID string) (bool, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}
	objUserID, err := primitive.ObjectIDFromHex(candidateID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	cursor, err := r.companiesCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"company_id":       objCompanyID,
				"saved_candidates": objUserID,
			},
		},
	})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	var isSaved bool
	if cursor.Next(ctx) {
		isSaved = true
	}

	return isSaved, nil
}

// AddCVInCareerCenter ...
func (r *Repository) AddCVInCareerCenter(ctx context.Context, userID, companyID string, opt careercenter.CVOptions) error {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.companiesCollection.UpdateOne(
		ctx,
		bson.M{
			"company_id": objCompanyID,
		},
		bson.M{
			"$push": bson.M{
				"cvs": bson.M{
					"user_id":                   objUserID,
					"young_professionals":       opt.YoungProfessionals,
					"expierenced_professionals": opt.ExpierencedProfessionals,
					"new_job_seekers":           opt.NewJobSeekers,
					"created_at":                time.Now(),
				},
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedCVs ...
func (r *Repository) GetSavedCVs(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error) {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	type CVs struct {
		IDs     primitive.ObjectID `bson:"user_id"`
		IsSaved bool               `bson:"is_saved"`
	}

	// get sent CVs
	cursor, err := r.companiesCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": objCompanyID,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$cvs",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$cvs",
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ids := make([]primitive.ObjectID, 0)
	mm := make([]CVs, 0)

	for cursor.Next(ctx) {
		var m CVs

		err = cursor.Decode(&m)
		if err != nil {
			return nil, err
		}

		ids = append(ids, m.IDs)
		mm = append(mm, m)
	}

	// get career interests
	cursor, err = r.profileCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"is_open": true,
					"user_id": bson.M{
						"$in": ids,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	views := make([]*candidate.ViewForCompany, 0)
	for cursor.Next(ctx) {
		j := candidate.ViewForCompany{}

		err = cursor.Decode(&j)
		if err != nil {
			return nil, err
		}

		views = append(views, &j)
	}

	for i := range views {
		for j := range mm {
			if views[i].UserID == mm[j].IDs {
				views[i].IsSaved = &mm[j].IsSaved
			}
		}
	}

	return views, nil
}

// RemoveCVs ...
func (r *Repository) RemoveCVs(ctx context.Context, companyID string, ids []string) error {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	_, err = r.companiesCollection.UpdateOne(
		ctx,
		bson.M{
			"company_id": objCompanyID,
		},
		bson.M{
			"$pull": bson.M{
				"cvs": bson.M{
					"user_id": bson.M{
						"$in": objIDs,
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

// MakeFavoriteCVs ...
func (r *Repository) MakeFavoriteCVs(ctx context.Context, companyID string, ids []string, isFavourite bool) error {
	objCompanyID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	_, err = r.companiesCollection.UpdateMany(ctx,
		bson.M{
			"company_id": objCompanyID,
		},
		bson.M{
			"$set": bson.M{
				"cvs.$[elem].is_saved": isFavourite,
			},
		},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{
					"elem.user_id": bson.M{
						"$in": objIDs,
					},
				},
			},
		}),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) getUnuploadedFilesForApplication(ctx context.Context, userID, jobID string, fileIDs []string) ([]*job.File, error) {
	objFileIDs := make([]primitive.ObjectID, 0, len(fileIDs))
	for _, id := range fileIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objFileIDs = append(objFileIDs, objID)
	}

	/*objUserID*/
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.jobsCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"_id": objJobID,
			},
		},
		{
			"$project": bson.M{
				"unuploaded_files." + userID: 1,
			},
		},
		{
			"$unwind": bson.M{
				"path": "$unuploaded_files." + userID,
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$unuploaded_files." + userID,
			},
		},
		{
			"$match": bson.M{
				"id": bson.M{
					"$in": objFileIDs,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	files := make([]*job.File, 0)

	for cursor.Next(ctx) {
		f := job.File{}
		err = cursor.Decode(&f)
		if err != nil {
			return nil, err
		}

		files = append(files, &f)
	}

	return files, nil
}

func (r *Repository) removeUnuploadedFilesForApplication(ctx context.Context, userID, jobID string) error {
	// objFileIDs := make([]primitive.ObjectID, 0, len(fileIDs))
	// for _, id := range fileIDs {
	// 	objID, err := primitive.ObjectIDFromHex(id)
	// 	if err != nil {
	// 		return errors.New(`wrong_id`)
	// 	}
	// 	objFileIDs = append(objFileIDs, objID)
	// }

	/*objUserID*/
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objJobID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	// cursor, err := r.jobsCollection.Aggregate(ctx, []bson.M{
	// 	{
	// 		"$match": bson.M{
	// 			"_id": objJobID,
	// 		},
	// 	},
	// 	{
	// 		"$project": bson.M{
	// 			"unuploaded_files." + userID: 1,
	// 		},
	// 	},
	// 	{
	// 		"$unwind": bson.M{
	// 			"path": "$unuploaded_files." + userID,
	// 		},
	// 	},
	// 	{
	// 		"$replaceRoot": bson.M{
	// 			"newRoot": "$unuploaded_files." + userID,
	// 		},
	// 	},
	// 	{
	// 		"$match": bson.M{
	// 			"id": bson.M{
	// 				"$in": objFileIDs,
	// 			},
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// defer cursor.Close(ctx)

	_, err = r.jobsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objJobID,
		},
		bson.M{
			"$unset": bson.M{
				"unuploaded_files." + userID: "",
			},
		},
	)
	if err != nil {
		return err
	}

	// files := make([]*job.File, 0)
	//
	// for cursor.Next(ctx) {
	// 	f := job.File{}
	// 	err = cursor.Decode(&f)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	files = append(files, &f)
	// }

	return nil
}

// -----

func mongoCondition(field, condition, thenValue string, elseValue interface{}) bson.M {
	return bson.M{
		"$cond": bson.M{
			"if": bson.M{
				condition: field,
			},
			"then": bson.M{
				thenValue: field,
			},
			"else": elseValue,
		},
	}
}
