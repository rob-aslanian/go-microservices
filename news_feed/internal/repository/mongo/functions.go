package repository

import (
	"context"
	"errors"

	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/file"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/post"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SavePost ...
func (r Repository) SavePost(ctx context.Context, p *post.Post) error {
	_, err := r.postsCollection.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

// GetPostByID changes only text, hashtags and changed_at
func (r Repository) GetPostByID(ctx context.Context, id string) (*post.Post, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id": objID,
					"deleted": bson.M{
						"$ne": true,
					},
				},
			},
			{
				"$project": bson.M{
					"comments": 0,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	p := post.Post{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&p)
		if err != nil {
			return nil, err
		}
	}

	return &p, nil
}

// ChangePost ...
func (r Repository) ChangePost(ctx context.Context, p *post.Post) error {
	_, err := r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": p.ID,
			"deleted": bson.M{
				"$ne": true,
			},
		},
		bson.M{
			"$set": bson.M{
				"text":       p.Text,
				"hashtags":   p.Hashtags,
				"changed_at": p.ChangedAt,
				"is_pinned":  p.IsPinned,
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// RemovePost ...
func (r Repository) RemovePost(ctx context.Context, postID string) error {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"deleted": true,
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// GetNewsfeed ...
func (r Repository) GetNewsfeed(ctx context.Context, requestorID, id string, pinned bool, first, after uint32) (*post.Newsfeed, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objRequestorID, err := primitive.ObjectIDFromHex(requestorID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	match := bson.M{
		"$or": []bson.M{
			{"newsfeed_user_id": objID},
			{"newsfeed_company_id": objID},
		},
		"deleted": bson.M{
			"$ne": true,
		},
	}

	if pinned {
		match["is_pinned"] = true
	} else {
		match["is_pinned"] = false
	}

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$addFields": bson.M{
					"liked": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{
								"$filter": bson.M{
									"input": "$likes",
									"cond":  bson.M{"$eq": []interface{}{"$$this._id", objRequestorID}},
								},
							},
							nil,
						},
					},
					"comments_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$comments"},
							bson.M{"$size": bson.M{
								"$filter": bson.M{
									"input": "$comments",
									"cond":  bson.M{"$eq": []interface{}{"$$this.parent_id", nil}},
								},
							},
							},
							0,
						},
					},
					"likes_amount.like": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "like"}},
							}}},
							0,
						},
					},
					"likes_amount.heart": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "heart"}},
							}}},
							0,
						},
					},
					"likes_amount.stop": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "stop"}},
							}}},
							0,
						},
					},
					"likes_amount.hmm": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "hmm"}},
							}}},
							0,
						},
					},
					"likes_amount.clap": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "clap"}},
							}}},
							0,
						},
					},
					"likes_amount.rocket": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "rocket"}},
							}}},
							0,
						},
					},
					"likes_amount.shit": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "shit"}},
							}}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"likes":    0,
					"comments": 0,
				},
			},
			{
				"$facet": bson.M{
					"posts": []bson.M{
						{
							"$sort": bson.M{
								"created_at": -1,
							},
						},
						{"$skip": after},
						{"$limit": first},
						{"$addFields": bson.M{
							"liked": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$liked"},
									bson.M{"$arrayElemAt": []interface{}{"$liked", 0}},
									nil,
								},
							},
						}},
					},
					"amount": []bson.M{
						{"$count": "count"},
					},
				},
			},
			{
				"$project": bson.M{
					"posts": 1,
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$amount"},
							bson.M{"$arrayElemAt": []interface{}{"$amount", 0}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"posts":  1,
					"amount": "$amount.count",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	feed := post.Newsfeed{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&feed)
		if err != nil {
			return nil, err
		}
	}

	return &feed, nil
}

// AddComment ...
func (r Repository) AddComment(ctx context.Context, com *post.Comment) error {
	match := bson.M{
		"_id":              com.PostID,
		"comment_disabled": false,
		"deleted": bson.M{
			"$ne": true,
		},
	}

	_, err := r.postsCollection.UpdateOne(
		ctx,
		match,
		bson.M{
			"$push": bson.M{
				"comments": com,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeComment ...
func (r Repository) ChangeComment(ctx context.Context, com *post.Comment) error {
	postObjID, err := primitive.ObjectIDFromHex(com.GetPostID())
	if err != nil {
		return errors.New(`wrong_id`)
	}

	commentObjID, err := primitive.ObjectIDFromHex(com.GetID())
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": postObjID,
			"deleted": bson.M{
				"$ne": true,
			},
			"comments._id": commentObjID,
		},
		bson.M{
			"$set": bson.M{
				"comments.$.text":       com.Text,
				"comments.$.tags":       com.Tags,
				"comments.$.changed_at": com.ChangedAt,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveComment ...
func (r Repository) RemoveComment(ctx context.Context, postID, commentID string) error {
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": postObjID,
			"deleted": bson.M{
				"$ne": true,
			},
			"comments._id": commentObjID,
		},
		bson.M{
			"$set": bson.M{
				"comments.$.deleted": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetComments ...
func (r Repository) GetComments(ctx context.Context, myID, postID string, sort post.CommentSort, first, after uint32) (*post.Comments, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objMyID, err := primitive.ObjectIDFromHex(myID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	var sortAgg bson.D

	if sort == post.CommentSortTop {
		sortAgg = bson.D{
			{
				Key:   "total_likes",
				Value: -1,
			},
		}
	}

	sortAgg = append(sortAgg, bson.E{
		Key:   "created_at",
		Value: -1,
	})

	agg := make([]bson.M, 0, 7)
	agg = append(agg,
		bson.M{
			"$match": bson.M{
				"_id": objID,
				"deleted": bson.M{
					"$ne": true,
				},
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$comments",
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": "$comments",
			},
		},
	)

	if sort == post.CommentSortTop {
		agg = append(agg, bson.M{
			"$addFields": bson.M{
				"total_likes": bson.M{
					"$cond": []interface{}{
						bson.M{"$isArray": "$likes"},
						bson.M{"$size": "$likes"},
						0,
					},
				},
			},
		},
		)
	}

	agg = append(agg,

		bson.M{
			"$facet": bson.M{
				"comments": []bson.M{
					{
						"$match": bson.M{
							"parent_id": nil,
							"deleted": bson.M{
								"$ne": true,
							},
						},
					},
					{
						"$sort": sortAgg,
					},
					{
						"$skip": after,
					},
					{
						"$limit": first,
					},
					{
						"$addFields": bson.M{
							"liked": bson.M{
								"$cond": []bson.M{
									{"$and": []bson.M{
										{"$isArray": "$likes"},
										{"$gt": []interface{}{bson.M{"$size": "$likes"}, 0}},
									},
									},
									{
										"$filter": bson.M{
											"input": "$likes",
											"cond":  bson.M{"$eq": []interface{}{"$$this._id", objMyID}},
										},
									},
									nil,
								},
							},
							"likes_amount.like": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "like"}},
									}}},
									0,
								},
							},
							"likes_amount.heart": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "heart"}},
									}}},
									0,
								},
							},
							"likes_amount.stop": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "stop"}},
									}}},
									0,
								},
							},
							"likes_amount.hmm": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "hmm"}},
									}}},
									0,
								},
							},
							"likes_amount.clap": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "clap"}},
									}}},
									0,
								},
							},
							"likes_amount.rocket": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "rocket"}},
									}}},
									0,
								},
							},
							"likes_amount.shit": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "shit"}},
									}}},
									0,
								},
							},
						},
					},
					{
						"$addFields": bson.M{
							"liked": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$liked"},
									bson.M{"$arrayElemAt": []interface{}{"$liked", 0}},
									nil,
								},
							},
						},
					},
					{
						"$project": bson.M{
							"replies": 0,
						},
					},
				},
				"amount": []bson.M{
					{
						"$match": bson.M{
							"deleted": bson.M{
								"$ne": true,
							},
							"parent_id": nil,
						},
					},
					{
						"$count": "count",
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"comments": 1,
				"amount": bson.M{
					"$cond": []interface{}{
						bson.M{
							"$isArray": "$amount",
						},
						bson.M{
							"$arrayElemAt": []interface{}{
								"$amount",
								0,
							},
						},
						0,
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"comments": 1,
				"amount":   "$amount.count",
			},
		})

	cursor, err := r.postsCollection.Aggregate(ctx, agg)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	com := post.Comments{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&com)
		if err != nil {
			return nil, err
		}
	}
	return &com, nil
}

// GetCommentReplies ...
func (r Repository) GetCommentReplies(ctx context.Context, myID, postID string, commentID string, first, after uint32) (*post.Comments, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objCommentID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objMyID, err := primitive.ObjectIDFromHex(myID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	agg := []bson.M{
		{
			"$match": bson.M{
				"_id": objID,
				"deleted": bson.M{
					"$ne": true,
				},
			},
		},
		{
			"$unwind": bson.M{
				"path": "$comments",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$comments",
			},
		},
		{
			"$facet": bson.M{
				"comments": []bson.M{
					{
						"$match": bson.M{
							"parent_id": objCommentID,
							"deleted": bson.M{
								"$ne": true,
							},
						},
					},
					{
						"$sort": bson.M{
							"comments.created_at": -1,
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
							"liked": bson.M{
								"$cond": []bson.M{
									{"$and": []bson.M{
										{"$isArray": "$likes"},
										{"$gt": []interface{}{bson.M{"$size": "$likes"}, 0}},
									},
									},
									{
										"$filter": bson.M{
											"input": "$likes",
											"cond":  bson.M{"$eq": []interface{}{"$$this._id", objMyID}},
										},
									},
									nil,
								},
							},
							"likes_amount.like": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "like"}},
									}}},
									0,
								},
							},
							"likes_amount.heart": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "heart"}},
									}}},
									0,
								},
							},
							"likes_amount.stop": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "stop"}},
									}}},
									0,
								},
							},
							"likes_amount.hmm": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "hmm"}},
									}}},
									0,
								},
							},
							"likes_amount.clap": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "clap"}},
									}}},
									0,
								},
							},
							"likes_amount.rocket": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "rocket"}},
									}}},
									0,
								},
							},
							"likes_amount.shit": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$likes"},
									bson.M{"$size": bson.M{"$filter": bson.M{
										"input": "$likes",
										"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "shit"}},
									}}},
									0,
								},
							},
						},
					},
					{
						"$addFields": bson.M{
							"liked": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$liked"},
									bson.M{"$arrayElemAt": []interface{}{"$liked", 0}},
									nil,
								},
							},
						},
					},
					{
						"$project": bson.M{
							"replies": 0,
						},
					},
				},
				"amount": []bson.M{
					{
						"$match": bson.M{
							"parent_id": objCommentID,
							"deleted": bson.M{
								"$ne": true,
							},
						},
					},
					{
						"$count": "count",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"comments": 1,
				"amount": bson.M{
					"$cond": []interface{}{
						bson.M{
							"$isArray": "$amount",
						},
						bson.M{
							"$arrayElemAt": []interface{}{
								"$amount",
								0,
							},
						},
						0,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"comments": 1,
				"amount":   "$amount.count",
			},
		},
	}

	cursor, err := r.postsCollection.Aggregate(ctx, agg)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	com := post.Comments{}

	if cursor.Next(ctx) {
		err = cursor.Decode(&com)
		if err != nil {
			return nil, err
		}
	}
	return &com, nil
}

// GetSharedPosts ...
func (r Repository) GetSharedPosts(ctx context.Context, requestorID, id string, first, after uint32) (*post.Newsfeed, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objRequestorID, err := primitive.ObjectIDFromHex(requestorID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	match := bson.M{
		"shared_post_id": objID,
		"deleted": bson.M{
			"$ne": true,
		},
	}

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$addFields": bson.M{
					"liked": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{
								"$filter": bson.M{
									"input": "$likes",
									"cond":  bson.M{"$eq": []interface{}{"$$this._id", objRequestorID}},
								},
							},
							nil,
						},
					},
					"comments_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$comments"},
							bson.M{"$size": "$comments"},
							0,
						},
					},
					"likes_amount.like": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "like"}},
							}}},
							0,
						},
					},
					"likes_amount.heart": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "heart"}},
							}}},
							0,
						},
					},
					"likes_amount.stop": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "stop"}},
							}}},
							0,
						},
					},
					"likes_amount.hmm": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "hmm"}},
							}}},
							0,
						},
					},
					"likes_amount.clap": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "clap"}},
							}}},
							0,
						},
					},
					"likes_amount.rocket": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "rocket"}},
							}}},
							0,
						},
					},
					"likes_amount.shit": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "shit"}},
							}}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"comments": 0,
					"likes":    0,
				},
			},
			{
				"$facet": bson.M{
					"posts": []bson.M{
						{
							"$sort": bson.M{
								"created_at": -1,
							},
						},
						{"$skip": after},
						{"$limit": first},
						{"$addFields": bson.M{
							"liked": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$liked"},
									bson.M{"$arrayElemAt": []interface{}{"$liked", 0}},
									nil,
								},
							},
						}},
					},
					"amount": []bson.M{
						{"$count": "count"},
					},
				},
			},
			{
				"$project": bson.M{
					"posts": 1,
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$amount"},
							bson.M{"$arrayElemAt": []interface{}{"$amount", 0}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"posts":  1,
					"amount": "$amount.count",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	feed := post.Newsfeed{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&feed)
		if err != nil {
			return nil, err
		}
	}

	return &feed, nil
}

// GetAmountOfSharedPosts ...
func (r Repository) GetAmountOfSharedPosts(ctx context.Context, postIDs []string) (map[string]uint32, error) {
	objIDs := make([]primitive.ObjectID, 0, len(postIDs))
	for _, postID := range postIDs {
		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	cursor, err := r.postsCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"shared_post_id": bson.M{
					"$in": objIDs,
				},
				"deleted": bson.M{
					"$ne": true,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$shared_post_id",
				"amount_shares": bson.M{
					"$sum": 1,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	m := make(map[string]uint32)

	s := struct {
		ID     primitive.ObjectID `bson:"_id"`
		Amount uint32             `bson:"amount_shares"`
	}{}

	for cursor.Next(ctx) {
		err = cursor.Decode(&s)
		if err != nil {
			return nil, err
		}
		m[s.ID.Hex()] = s.Amount
	}

	return m, nil
}

// AddFile ...
func (r Repository) AddFile(ctx context.Context, postID string, f *file.File) error {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}
	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
			"deleted": bson.M{
				"$ne": true,
			},
		},
		bson.M{
			"$push": bson.M{
				"files": f,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddFileInComment ...
func (r Repository) AddFileInComment(ctx context.Context, postID, commentID string, f *file.File) error {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	commentObjID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
			"deleted": bson.M{
				"$ne": true,
			},
			"comments._id": commentObjID,
		},
		bson.M{
			"$push": bson.M{
				"comments.$.files": f,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFile ...
func (r Repository) RemoveFile(ctx context.Context, postID string, fileID string) error {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objFileID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objPostID,
			"deleted": bson.M{
				"$ne": true,
			},
		},
		bson.M{
			"$pull": bson.M{
				"files": bson.M{
					"id": objFileID,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFileInComment ...
func (r Repository) RemoveFileInComment(ctx context.Context, postID, commentID string, fileID string) error {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objFileID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCommentID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objPostID,
			"deleted": bson.M{
				"$ne": true,
			},
			"comments._id": objCommentID,
		},
		bson.M{
			"$pull": bson.M{
				"comments.$.files": bson.M{
					"id": objFileID,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetNewsfeedOfFollowings ...
func (r Repository) GetNewsfeedOfFollowings(ctx context.Context, requestorID string, ids []string, first, after uint32) (*post.Newsfeed, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	objRequestorID, err := primitive.ObjectIDFromHex(requestorID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	match := bson.M{
		"$or": []bson.M{
			{
				"user_id": bson.M{
					"$in": objIDs,
				},
				"newsfeed_user_id": bson.M{
					"$in": objIDs,
				},
			},
			{
				"company_id": bson.M{
					"$in": objIDs,
				},
				"newsfeed_company_id": bson.M{
					"$in": objIDs,
				},
			},
		},
		"deleted": bson.M{
			"$ne": true,
		},
	}

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$addFields": bson.M{
					"liked": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{
								"$filter": bson.M{
									"input": "$likes",
									"cond":  bson.M{"$eq": []interface{}{"$$this._id", objRequestorID}},
								},
							},
							nil,
						},
					},
					"comments_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$comments"},
							bson.M{"$size": "$comments"},
							0,
						},
					},
					"likes_amount.like": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "like"}},
							}}},
							0,
						},
					},
					"likes_amount.heart": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "heart"}},
							}}},
							0,
						},
					},
					"likes_amount.stop": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "stop"}},
							}}},
							0,
						},
					},
					"likes_amount.hmm": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "hmm"}},
							}}},
							0,
						},
					},
					"likes_amount.clap": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "clap"}},
							}}},
							0,
						},
					},
					"likes_amount.rocket": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "rocket"}},
							}}},
							0,
						},
					},
					"likes_amount.shit": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "shit"}},
							}}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"likes":    0,
					"comments": 0,
				},
			},
			{
				"$facet": bson.M{
					"posts": []bson.M{
						{
							"$sort": bson.M{
								"created_at": -1,
							},
						},
						{"$skip": after},
						{"$limit": first},
						{"$addFields": bson.M{
							"liked": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$liked"},
									bson.M{"$arrayElemAt": []interface{}{"$liked", 0}},
									nil,
								},
							},
						}},
					},
					"amount": []bson.M{
						{"$count": "count"},
					},
				},
			},
			{
				"$project": bson.M{
					"posts": 1,
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$amount"},
							bson.M{"$arrayElemAt": []interface{}{"$amount", 0}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"posts":  1,
					"amount": "$amount.count",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	feed := post.Newsfeed{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&feed)
		if err != nil {
			return nil, err
		}
	}

	return &feed, nil
}

// GetAmountOfReplies ...
func (r Repository) GetAmountOfReplies(ctx context.Context, postID string, commentsIDs []string) (map[string]uint32, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objIDs := make([]primitive.ObjectID, 0, len(commentsIDs))
	for _, postID := range commentsIDs {
		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
		objIDs = append(objIDs, objID)
	}

	cursor, err := r.postsCollection.Aggregate(ctx, []bson.M{

		{
			"$match": bson.M{
				"_id": objID,
			},
		},
		{
			"$unwind": bson.M{
				"path": "$comments",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$comments",
			},
		},
		{
			"$match": bson.M{
				"parent_id": bson.M{
					"$in": objIDs,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$parent_id",
				"amount": bson.M{
					"$sum": 1,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	m := make(map[string]uint32)

	s := struct {
		ID     primitive.ObjectID `bson:"_id"`
		Amount uint32             `bson:"amount"`
	}{}

	for cursor.Next(ctx) {
		err = cursor.Decode(&s)
		if err != nil {
			return nil, err
		}
		m[s.ID.Hex()] = s.Amount
	}

	return m, nil
}

// GetCommentByID ...
func (r Repository) GetCommentByID(ctx context.Context, postID string, commentID string) (*post.Comment, error) {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objCommentID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	res := r.postsCollection.FindOne(
		ctx,
		bson.M{
			"_id": objPostID,
			"deleted": bson.M{
				"$ne": true,
			},
			"comments._id": objCommentID,
		},
	)

	if res.Err() != nil {
		return nil, err
	}

	c := post.Comment{}

	err = res.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// LikePost ...
func (r Repository) LikePost(ctx context.Context, postID string, like *post.Like) error {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_ = r.UnlikePost(ctx, postID, like.GetID())

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objPostID,
		},
		bson.M{
			"$push": bson.M{
				"likes": like,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// UnlikePost ...
func (r Repository) UnlikePost(ctx context.Context, postID, id string) error {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objPostID,
		},
		bson.M{
			"$pull": bson.M{
				"likes": bson.M{
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

// LikeComment ...
func (r Repository) LikeComment(ctx context.Context, postID, commentID string, like *post.Like) error {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCommentID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_ = r.UnlikeComment(ctx, postID, commentID, like.GetID())

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":          objPostID,
			"comments._id": objCommentID,
		},
		bson.M{
			"$push": bson.M{
				"comments.$.likes": like,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// UnlikeComment ...
func (r Repository) UnlikeComment(ctx context.Context, postID, commentID string, id string) error {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objCommentID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.postsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":          objPostID,
			"comments._id": objCommentID,
		},
		bson.M{
			"$pull": bson.M{
				"comments.$.likes": bson.M{
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

// GetLikeInPostByID ...
func (r Repository) GetLikeInPostByID(ctx context.Context, postID, id string) (*post.Like, error) {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id": objPostID,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$likes",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$likes",
				},
			},
			{
				"$match": bson.M{
					"_id": objID,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	like := new(post.Like)

	if cursor.Next(ctx) {
		err = cursor.Decode(like)
		if err != nil {
			return nil, err
		}
	}

	return like, nil
}

// GetLikeInCommentByID ...
func (r Repository) GetLikeInCommentByID(ctx context.Context, postID, commentID, id string) (*post.Like, error) {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objCommentID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id":          objPostID,
					"comments._id": objCommentID,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$comments",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$comments",
				},
			},
			{
				"$match": bson.M{
					"_id": objCommentID,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$likes",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$likes",
				},
			},
			{
				"$match": bson.M{
					"_id": objID,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	like := new(post.Like)

	if cursor.Next(ctx) {
		err = cursor.Decode(like)
		if err != nil {
			return nil, err
		}
	}

	return like, nil
}

// GetLikedListInPost ...
func (r Repository) GetLikedListInPost(ctx context.Context, postID string, emoji *post.EmojiType, first, after uint32) ([]*post.Like, error) {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	agg := make([]bson.M, 0, 5)

	agg = append(agg,
		bson.M{
			"$match": bson.M{
				"_id": objPostID,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$likes",
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": "$likes",
			},
		},
	)

	if emoji != nil {
		agg = append(agg,
			bson.M{
				"$match": bson.M{
					"emoji": *emoji,
				},
			},
		)
	}

	agg = append(agg,
		bson.M{
			"$skip": after,
		},
		bson.M{
			"$limit": first,
		},
	)

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		agg,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	likes := make([]*post.Like, 0)

	for cursor.Next(ctx) {
		l := new(post.Like)
		err = cursor.Decode(&l)
		if err != nil {
			return nil, err
		}
		likes = append(likes, l)
	}

	return likes, nil
}

// GetLikedListInComment ...
func (r Repository) GetLikedListInComment(ctx context.Context, postID, commentID string, emoji *post.EmojiType, first, after uint32) ([]*post.Like, error) {
	objPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	objCommentID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	agg := make([]bson.M, 0, 8)
	agg = append(agg,
		bson.M{
			"$match": bson.M{
				"_id": objPostID,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$comments",
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": "$comments",
			},
		},
		bson.M{
			"$match": bson.M{
				"_id": objCommentID,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$likes",
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": "$likes",
			},
		},
	)

	if emoji != nil {
		agg = append(agg,
			bson.M{
				"$match": bson.M{
					"emoji": *emoji,
				},
			},
		)
	}

	agg = append(agg,
		bson.M{
			"$skip": after,
		},
		bson.M{
			"$limit": first,
		},
	)

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		agg,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	likes := make([]*post.Like, 0)

	for cursor.Next(ctx) {
		l := new(post.Like)
		err = cursor.Decode(&l)
		if err != nil {
			return nil, err
		}
		likes = append(likes, l)
	}

	return likes, nil
}

// Search ...
func (r Repository) Search(ctx context.Context, requestorID, newsfeedID string, keywords string, hashtags []string, first, after uint32) (*post.Newsfeed, error) {

	var objID primitive.ObjectID
	var err error

	if newsfeedID != "" {
		objID, err = primitive.ObjectIDFromHex(newsfeedID)
		if err != nil {
			return nil, errors.New(`wrong_id`)
		}
	}

	objRequestorID, err := primitive.ObjectIDFromHex(requestorID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	match := bson.M{
		"deleted": bson.M{
			"$ne": true,
		},
	}

	if newsfeedID != "" {
		match["$or"] = []bson.M{
			{"newsfeed_user_id": objID},
			{"newsfeed_company_id": objID},
		}
	}

	if len(keywords) > 0 {
		match["$text"] = bson.M{
			"$search": keywords,
		}
	}

	if len(hashtags) > 0 {
		match["hashtags"] = bson.M{
			"$in": hashtags,
		}
	}

	cursor, err := r.postsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$addFields": bson.M{
					"liked": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{
								"$filter": bson.M{
									"input": "$likes",
									"cond":  bson.M{"$eq": []interface{}{"$$this._id", objRequestorID}},
								},
							},
							nil,
						},
					},
					"comments_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$comments"},
							bson.M{"$size": bson.M{
								"$filter": bson.M{
									"input": "$comments",
									"cond":  bson.M{"$eq": []interface{}{"$$this.parent_id", nil}},
								},
							},
							},
							0,
						},
					},
					"likes_amount.like": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "like"}},
							}}},
							0,
						},
					},
					"likes_amount.heart": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "heart"}},
							}}},
							0,
						},
					},
					"likes_amount.stop": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "stop"}},
							}}},
							0,
						},
					},
					"likes_amount.hmm": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "hmm"}},
							}}},
							0,
						},
					},
					"likes_amount.clap": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "clap"}},
							}}},
							0,
						},
					},
					"likes_amount.rocket": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "rocket"}},
							}}},
							0,
						},
					},
					"likes_amount.shit": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": bson.M{"$filter": bson.M{
								"input": "$likes",
								"cond":  bson.M{"$eq": []interface{}{"$$this.emoji", "shit"}},
							}}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"likes":    0,
					"comments": 0,
				},
			},
			{
				"$facet": bson.M{
					"posts": []bson.M{
						{
							"$sort": bson.M{
								"created_at": -1,
							},
						},
						{"$skip": after},
						{"$limit": first},
						{"$addFields": bson.M{
							"liked": bson.M{
								"$cond": []interface{}{
									bson.M{"$isArray": "$liked"},
									bson.M{"$arrayElemAt": []interface{}{"$liked", 0}},
									nil,
								},
							},
						}},
					},
					"amount": []bson.M{
						{"$count": "count"},
					},
				},
			},
			{
				"$project": bson.M{
					"posts": 1,
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$amount"},
							bson.M{"$arrayElemAt": []interface{}{"$amount", 0}},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"posts":  1,
					"amount": "$amount.count",
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	feed := post.Newsfeed{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&feed)
		if err != nil {
			return nil, err
		}
	}

	return &feed, nil
}
