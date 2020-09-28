package resolver

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/newsfeedRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

// AddPostInNewsfeed ...
func (r *Resolver) AddPostInNewsfeed(ctx context.Context, input AddPostInNewsfeedRequest) (*SuccessResolver, error) {
	resp, err := newsfeed.AddPost(
		ctx,
		postNewsfeedInputToNewsfeedRPCPost(&input.Post),
	)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      resp.GetID(),
			Success: true,
		},
	}, nil
}

// ChangePostInNewsfeed ...
func (r *Resolver) ChangePostInNewsfeed(ctx context.Context, input ChangePostInNewsfeedRequest) (*SuccessResolver, error) {
	p := newsfeedRPC.Post{
		ID:       input.Post.Post_id,
		Text:     input.Post.Text,
		IsPinned: input.Post.Is_pinned,
	}

	if input.Post.Company_id != nil {
		p.CompanyID = *input.Post.Company_id
	}

	_, err := newsfeed.ChangePost(
		ctx,
		&p,
	)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// RemovePostInNewsfeed ...
func (r *Resolver) RemovePostInNewsfeed(ctx context.Context, input RemovePostInNewsfeedRequest) (*SuccessResolver, error) {
	var companyID string

	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	_, err := newsfeed.RemovePost(
		ctx,
		&newsfeedRPC.RemovePostRequest{
			PostID:    input.Post_id,
			CompanyID: companyID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// GetNewsfeed ...
func (r *Resolver) GetNewsfeed(ctx context.Context, input GetNewsfeedRequest) (*NewsfeedResolverCustom, error) {
	var pinned bool

	if input.Pinned != nil {
		pinned = *input.Pinned
	}

	id := ""
	companyID := ""

	if input.ID != nil {
		id = *input.ID
	}

	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	resp, err := newsfeed.GetNewsfeed(ctx, &newsfeedRPC.GetNewsfeedRequest{
		ID:        id,
		CompanyID: companyID,
		Pinned:    pinned,
		Pagination: &newsfeedRPC.Pagination{
			First: input.First,
			After: input.After,
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.GetAmount() == 0 {
		return nil, nil
	}

	posts := make([]NewsfeedPostCustom, 0, len(resp.GetPosts()))

	for _, post := range resp.GetPosts() {
		if p := newsfeedPostToNewsfeedRPCPost(post); p != nil {
			posts = append(posts, *p)
		}
	}

	// --------

	userIDs := make([]string, 0, len(resp.GetPosts()))
	companiesIDs := make([]string, 0, len(resp.GetPosts()))

	for _, post := range resp.GetPosts() {
		if post.GetCompanyID() != "" {
			companiesIDs = append(companiesIDs, post.GetCompanyID())
		} else {
			userIDs = append(userIDs, post.GetUserID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, postRPC := range resp.GetPosts() {
		for i := range posts {
			if posts[i].ID == postRPC.GetID() {
				if postRPC.GetCompanyID() != "" {
					// company profile
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == postRPC.GetCompanyID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						posts[i].Company_profile = &pr
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[postRPC.GetUserID()])
					posts[i].User_profile = &profile
				}

				break
			}
		}
	}

	// --------

	return &NewsfeedResolverCustom{
		R: &NewsfeedCustom{
			Post_amount: int32(resp.GetAmount()),
			Posts:       posts,
		},
	}, nil
}

// AddCommentInPostInNewsfeed ...
func (r *Resolver) AddCommentInPostInNewsfeed(ctx context.Context, input AddCommentInPostInNewsfeedRequest) (*SuccessResolver, error) {
	resp, err := newsfeed.AddComment(
		ctx,
		commentNewsfeedInputToNewsfeedRPCComment(&input.Comment),
	)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      resp.GetID(),
			Success: true,
		},
	}, nil
}

// ChangeCommentInPostInNewsfeed ...
func (r *Resolver) ChangeCommentInPostInNewsfeed(ctx context.Context, input ChangeCommentInPostInNewsfeedRequest) (*SuccessResolver, error) {
	var companyID string

	if input.Comment.Company_id != nil {
		companyID = *input.Comment.Company_id
	}

	com := newsfeedRPC.Comment{
		ID:        input.Comment.Comment_id,
		CompanyID: companyID,
		Text:      input.Comment.Text,
		PostID:    input.Comment.Post_id,
	}

	if input.Comment.Tags != nil {
		com.Tags = make([]*newsfeedRPC.Tag, 0, len(*input.Comment.Tags))

		for _, t := range *input.Comment.Tags {
			com.Tags = append(com.Tags, newsfeedTagInputToNewsfeedRPCTag(&t))
		}
	}

	_, err := newsfeed.ChangeComment(
		ctx,
		&com,
	)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// RemoveCommentInPostInNewsfeed ...
func (r *Resolver) RemoveCommentInPostInNewsfeed(ctx context.Context, input RemoveCommentInPostInNewsfeedRequest) (*SuccessResolver, error) {
	req := newsfeedRPC.RemoveCommentRequest{
		PostID:    input.Post_id,
		CommentID: input.Comment_id,
	}

	if input.Company_id != nil {
		req.CompanyID = *input.Company_id
	}

	_, err := newsfeed.RemoveComment(
		ctx,
		&req,
	)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// GetCommentsOfNewsfeedPost ...
func (r *Resolver) GetCommentsOfNewsfeedPost(ctx context.Context, input GetCommentsOfNewsfeedPostRequest) (*NewsfeedPostCommentsResolverCustom, error) {
	var companyID string
	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	var sort newsfeedRPC.GetCommentsRequest_SortOption
	if input.Sort != nil {
		sort = stringToGetCommentsRequestSortOption(*input.Sort)
	}

	resp, err := newsfeed.GetComments(ctx, &newsfeedRPC.GetCommentsRequest{
		ID:        input.Post_id,
		CompanyID: companyID,
		Sort:      sort,
		Pagination: &newsfeedRPC.Pagination{
			First: input.First,
			After: input.After,
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.GetAmount() == 0 {
		return nil, nil
	}

	comments := make([]NewsfeedPostCommentCustom, 0, len(resp.GetComments()))

	for _, com := range resp.GetComments() {
		if p := newsfeedPostCommentToNewsfeedRPCComment(com); p != nil {
			comments = append(comments, *p)
		}
	}

	// --------

	userIDs := make([]string, 0, len(resp.GetComments()))
	companiesIDs := make([]string, 0, len(resp.GetComments()))

	for _, post := range resp.GetComments() {
		if post.GetCompanyID() != "" {
			companiesIDs = append(companiesIDs, post.GetCompanyID())
		} else {
			userIDs = append(userIDs, post.GetUserID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, postRPC := range resp.GetComments() {
		for i := range comments {
			if comments[i].ID == postRPC.GetID() {
				if postRPC.GetCompanyID() != "" {
					// company profile
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == postRPC.GetCompanyID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						comments[i].Company_profile = &pr
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[postRPC.GetUserID()])
					comments[i].User_profile = &profile
				}

				break
			}
		}
	}

	// --------

	return &NewsfeedPostCommentsResolverCustom{
		R: &NewsfeedPostCommentsCustom{
			Amount:   int32(resp.GetAmount()),
			Comments: comments,
		},
	}, nil
}

// GetCommentRepliesOfNewsfeedPost ...
func (r *Resolver) GetCommentRepliesOfNewsfeedPost(ctx context.Context, input GetCommentRepliesOfNewsfeedPostRequest) (*NewsfeedPostCommentsResolverCustom, error) {
	var companyID string
	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	resp, err := newsfeed.GetCommentReplies(ctx, &newsfeedRPC.GetCommentRepliesRequest{
		PostID:    input.Post_id,
		CommentID: input.CommentID,
		CompanyID: companyID,
		Pagination: &newsfeedRPC.Pagination{
			First: input.First,
			After: input.After,
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.GetAmount() == 0 {
		return nil, nil
	}

	comments := make([]NewsfeedPostCommentCustom, 0, len(resp.GetComments()))

	for _, com := range resp.GetComments() {
		if p := newsfeedPostCommentToNewsfeedRPCComment(com); p != nil {
			comments = append(comments, *p)
		}
	}

	// --------

	userIDs := make([]string, 0, len(resp.GetComments()))
	companiesIDs := make([]string, 0, len(resp.GetComments()))

	for _, post := range resp.GetComments() {
		if post.GetCompanyID() != "" {
			companiesIDs = append(companiesIDs, post.GetCompanyID())
		} else {
			userIDs = append(userIDs, post.GetUserID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, postRPC := range resp.GetComments() {
		for i := range comments {
			if comments[i].ID == postRPC.GetID() {
				if postRPC.GetCompanyID() != "" {
					// company profile
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == postRPC.GetCompanyID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						comments[i].Company_profile = &pr
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[postRPC.GetUserID()])
					comments[i].User_profile = &profile
				}

				break
			}
		}
	}

	// --------

	return &NewsfeedPostCommentsResolverCustom{
		R: &NewsfeedPostCommentsCustom{
			Amount:   int32(resp.GetAmount()),
			Comments: comments,
		},
	}, nil
}

// GetSharedPost ...
func (r *Resolver) GetSharedPost(ctx context.Context, input GetSharedPostRequest) (*NewsfeedResolverCustom, error) {
	companyID := ""
	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	resp, err := newsfeed.GetSharedPosts(ctx, &newsfeedRPC.GetSharedPostsRequest{
		ID:        input.ID,
		CompanyID: companyID,
		Pagination: &newsfeedRPC.Pagination{
			First: input.First,
			After: input.After,
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.GetAmount() == 0 {
		return nil, nil
	}

	posts := make([]NewsfeedPostCustom, 0, len(resp.GetPosts()))

	for _, post := range resp.GetPosts() {
		if p := newsfeedPostToNewsfeedRPCPost(post); p != nil {
			posts = append(posts, *p)
		}
	}

	// --------

	userIDs := make([]string, 0, len(resp.GetPosts()))
	companiesIDs := make([]string, 0, len(resp.GetPosts()))

	for _, post := range resp.GetPosts() {
		if post.GetCompanyID() != "" {
			companiesIDs = append(companiesIDs, post.GetCompanyID())
		} else {
			userIDs = append(userIDs, post.GetUserID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, postRPC := range resp.GetPosts() {
		for i := range posts {
			if posts[i].ID == postRPC.GetID() {
				if postRPC.GetCompanyID() != "" {
					// company profile
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == postRPC.GetCompanyID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						posts[i].Company_profile = &pr
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[postRPC.GetUserID()])
					posts[i].User_profile = &profile
				}

				break
			}
		}
	}

	// --------

	return &NewsfeedResolverCustom{
		R: &NewsfeedCustom{
			Post_amount: int32(resp.GetAmount()),
			Posts:       posts,
		},
	}, nil
}

// GetNewsfeedPost ...
func (r *Resolver) GetNewsfeedPost(ctx context.Context, input GetNewsfeedPostRequest) (*NewsfeedPostResolverCustom, error) {
	resp, err := newsfeed.GetPostByID(ctx, &newsfeedRPC.ID{
		ID: input.ID,
	})
	if err != nil {
		return nil, err
	}

	p := newsfeedPostToNewsfeedRPCPost(resp)

	// --------

	// userIDs := make([]string, 0, len(resp.GetPosts()))
	// companiesIDs := make([]string, 0, len(resp.GetPosts()))
	//
	// for _, post := range resp.GetPosts() {
	// 	if post.GetCompanyID() != "" {
	// 		companiesIDs = append(companiesIDs, post.GetCompanyID())
	// 	} else {
	// 		userIDs = append(userIDs, post.GetUserID())
	// 	}
	// }

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: []string{resp.GetCompanyID()},
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: []string{resp.GetUserID()},
	})
	if err != nil {
		return nil, err
	}

	// for _, postRPC := range resp.GetPosts() {
	// 	for i := range posts {
	// 		if posts[i].ID == postRPC.GetID() {
	if resp.GetCompanyID() != "" {
		// company profile
		var profile *companyRPC.Profile

		for _, company := range companyResp.GetProfiles() {
			if company.GetId() == resp.GetCompanyID() {
				profile = company
				break
			}
		}

		if profile != nil {
			pr := toCompanyProfile(ctx, *profile)
			p.Company_profile = &pr
		}
	} else {
		// user profile
		profile := ToProfile(ctx, userResp.GetProfiles()[resp.GetUserID()])
		p.User_profile = &profile
	}
	//
	// 			break
	// 		}
	// 	}
	// }

	// --------

	return &NewsfeedPostResolverCustom{
		R: p,
	}, nil
}

// RemoveFileInNewsfeed ...
func (r *Resolver) RemoveFileInNewsfeed(ctx context.Context, input RemoveFileInNewsfeedRequest) (*SuccessResolver, error) {
	var companyID string

	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	var commentID string

	if input.Comment_id != nil {
		commentID = *input.Comment_id
	}

	_, err := newsfeed.RemoveFileInPost(
		ctx,
		&newsfeedRPC.RemoveFileInPostRequest{
			PostID:    input.Post_id,
			FileID:    input.File_id,
			CompanyID: companyID,
			CommentID: commentID,
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

// LikePostInNewsfeed ...
func (r *Resolver) LikePostInNewsfeed(ctx context.Context, input LikePostInNewsfeedRequest) (*SuccessResolver, error) {
	var commentID string

	if input.Comment_id != nil {
		commentID = *input.Comment_id
	}

	_, err := newsfeed.Like(
		ctx,
		&newsfeedRPC.LikeRequest{
			PostID:    input.Post_id,
			CommentID: commentID,
			Like:      likeInputToNewsfeedRPCLike(&input.Like),
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

// UnlikePostInNewsfeed ...
func (r *Resolver) UnlikePostInNewsfeed(ctx context.Context, input UnlikePostInNewsfeedRequest) (*SuccessResolver, error) {
	var commentID string

	if input.Comment_id != nil {
		commentID = *input.Comment_id
	}

	_, err := newsfeed.Unlike(
		ctx,
		&newsfeedRPC.UnlikeRequest{
			PostID:    input.Post_id,
			CommentID: commentID,
			ID:        input.ID,
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

// GetListOfLikedInNewsfeed ...
func (r *Resolver) GetListOfLikedInNewsfeed(ctx context.Context, input GetListOfLikedInNewsfeedRequest) (*[]*LikableEntityUnionResolver, error) {
	var commentID string
	if input.Comment_id != nil {
		commentID = *input.Comment_id
	}

	var companyID string
	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	response, err := newsfeed.GetLikedList(
		ctx,
		&newsfeedRPC.GetLikedListRequest{
			PostID:    input.Post_id,
			CommentID: commentID,
			CompanyID: companyID,
			Emoji:     NullToString(input.Emoji),
			Pagination: &newsfeedRPC.Pagination{
				First: input.First,
				After: input.After,
			},
		})
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(response.GetLikes()))
	companiesIDs := make([]string, 0, len(response.GetLikes()))

	for _, l := range response.GetLikes() {
		if l != nil {
			switch l.GetEntity().String() {
			case "User":
				userIDs = append(userIDs, l.GetID())
			case "Company":
				companiesIDs = append(companiesIDs, l.GetID())
			}
		}
	}

	// ---

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	userProfiles := make(map[string]Profile, len(userResp.GetProfiles()))
	companyProfiles := make(map[string]CompanyProfile, len(companyResp.GetProfiles()))

	for key, value := range userResp.GetProfiles() {
		profile := ToProfile(ctx, value)
		userProfiles[key] = profile
	}

	for _, pr := range companyResp.GetProfiles() {
		if pr != nil {
			companyProfiles[pr.GetId()] = toCompanyProfile(ctx, *pr)
		}
	}
	// ---

	likes := make([]*LikableEntityUnionResolver, 0, len(response.GetLikes()))

	for _, l := range response.GetLikes() {
		if l != nil {
			switch l.GetEntity().String() {

			case "User":
				likes = append(likes, &LikableEntityUnionResolver{
					Result: &UserLikedItemResolver{
						R: &UserLikedItem{
							Profile: userProfiles[l.GetID()],
							Emoji:   l.GetEmoji(),
						},
					},
				})

			case "Company":
				likes = append(likes, &LikableEntityUnionResolver{
					Result: &CompanyLikedItemResolver{
						R: &CompanyLikedItem{
							Profile: companyProfiles[l.GetID()],
							Emoji:   l.GetEmoji(),
						},
					},
				})
			}
		}
	}

	return &likes, nil
}

// SearchAmongNewsfeedPosts ...
func (r *Resolver) SearchAmongNewsfeedPosts(ctx context.Context, input SearchAmongNewsfeedPostsRequest) (*NewsfeedResolverCustom, error) {
	companyID := ""
	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	newsfeedID := ""
	if input.Newsfeed_id != nil {
		newsfeedID = *input.Newsfeed_id
	}

	resp, err := newsfeed.SearchAmongPosts(ctx, &newsfeedRPC.SearchAmongPostsRequest{
		CompanyID:  companyID,
		NewsfeedID: newsfeedID,
		Keyword:    input.Keyword,
		Pagination: &newsfeedRPC.Pagination{
			First: input.First,
			After: input.After,
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.GetAmount() == 0 {
		return nil, nil
	}

	posts := make([]NewsfeedPostCustom, 0, len(resp.GetPosts()))

	for _, post := range resp.GetPosts() {
		if p := newsfeedPostToNewsfeedRPCPost(post); p != nil {
			posts = append(posts, *p)
		}
	}

	// --------

	userIDs := make([]string, 0, len(resp.GetPosts()))
	companiesIDs := make([]string, 0, len(resp.GetPosts()))

	for _, post := range resp.GetPosts() {
		if post.GetCompanyID() != "" {
			companiesIDs = append(companiesIDs, post.GetCompanyID())
		} else {
			userIDs = append(userIDs, post.GetUserID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, postRPC := range resp.GetPosts() {
		for i := range posts {
			if posts[i].ID == postRPC.GetID() {
				if postRPC.GetCompanyID() != "" {
					// company profile
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == postRPC.GetCompanyID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						posts[i].Company_profile = &pr
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[postRPC.GetUserID()])
					posts[i].User_profile = &profile
				}

				break
			}
		}
	}

	// --------

	return &NewsfeedResolverCustom{
		R: &NewsfeedCustom{
			Post_amount: int32(resp.GetAmount()),
			Posts:       posts,
		},
	}, nil
}

func postNewsfeedInputToNewsfeedRPCPost(data *PostNewsfeedInput) *newsfeedRPC.Post {
	if data == nil {
		return nil
	}

	p := newsfeedRPC.Post{
		Text: data.Text,
	}

	if data.Company_id != nil {
		p.CompanyID = *data.Company_id
	}

	if data.Newsfeed_company_id != nil {
		p.NewsFeedCompanyID = *data.Newsfeed_company_id
	}

	if data.Newsfeed_user_id != nil {
		p.NewsFeedUserID = *data.Newsfeed_user_id
	}

	if data.Shared_post_id != nil {
		p.SharedPostID = *data.Shared_post_id
	}

	if data.Tags != nil {
		p.Tags = make([]*newsfeedRPC.Tag, 0, len(*data.Tags))

		for _, t := range *data.Tags {
			p.Tags = append(p.Tags, newsfeedTagInputToNewsfeedRPCTag(&t))
		}
	}

	return &p
}

func newsfeedTagInputToNewsfeedRPCTag(data *NewsfeedTagInput) *newsfeedRPC.Tag {
	if data == nil {
		return nil
	}

	tag := newsfeedRPC.Tag{
		ID:     data.ID,
		Entity: stringToNewsfeedRPCEntityType(data.Type),
	}

	return &tag
}

func stringToNewsfeedRPCEntityType(data string) newsfeedRPC.EntityType {
	switch data {
	case "group":
		return newsfeedRPC.EntityType_Group
	case "company":
		return newsfeedRPC.EntityType_Company
	case "community":
		return newsfeedRPC.EntityType_Community
	case "organization":
		return newsfeedRPC.EntityType_Organization
	}

	return newsfeedRPC.EntityType_User
}

func commentNewsfeedInputToNewsfeedRPCComment(data *CommentPostNewsfeedInput) *newsfeedRPC.Comment {
	if data == nil {
		return nil
	}

	com := newsfeedRPC.Comment{
		Text:   data.Text,
		PostID: data.Post_id,
	}

	if data.Company_id != nil {
		com.CompanyID = *data.Company_id
	}

	if data.Parent_id != nil {
		com.ParentID = *data.Parent_id
	}

	if data.Tags != nil {
		com.Tags = make([]*newsfeedRPC.Tag, 0, len(*data.Tags))

		for _, t := range *data.Tags {
			com.Tags = append(com.Tags, newsfeedTagInputToNewsfeedRPCTag(&t))
		}
	}

	return &com
}

func likeInputToNewsfeedRPCLike(data *LikeInput) *newsfeedRPC.Like {
	if data == nil {
		return nil
	}

	like := newsfeedRPC.Like{
		ID:     data.ID,
		Emoji:  data.Emoji,
		Entity: stringToNewsfeedRPCEntityType(data.Type),
	}

	return &like
}

// ---------

func newsfeedPostToNewsfeedRPCPost(data *newsfeedRPC.Post) *NewsfeedPostCustom {
	if data == nil {
		return nil
	}

	post := NewsfeedPostCustom{
		ID:                       data.GetID(),
		Changed_at:               data.GetChangedAt(),
		Created_at:               data.GetCreatedAt(),
		Is_comments_disabled:     data.GetIsCommentedDisabled(),
		Is_notification_disabled: data.GetIsNotificationDisabled(),
		Comments_amount:          int32(data.GetCommentsAmount()),
		Shares_amount:            int32(data.GetSharesAmount()),
		Text:                     data.GetText(),
		Hashtags:                 data.GetHashtags(),
		Liked:                    data.GetLiked(),
		Tags:                     make([]NewsfeedTag, 0, len(data.GetTags())),
		Shared_post_id:           data.GetSharedPostID(),
		Files:                    make([]File, 0, len(data.GetFiles())),
	}

	if lk := newsfeedRPCLikesAmountToLikesAmount(data.GetLikesAmount()); lk != nil {
		post.Likes_amount = *lk
	}

	for _, f := range data.GetFiles() {
		fi := newsfeedRPCFileToFile(f)
		if fi != nil {
			post.Files = append(post.Files, *fi)
		}
	}

	for _, t := range data.GetTags() {
		if ta := newsfeedRPCTagToNewsfeedTag(t); ta != nil {
			post.Tags = append(post.Tags, *ta)
		}
	}

	return &post
}

func newsfeedRPCLikesAmountToLikesAmount(data *newsfeedRPC.LikesAmount) *LikesAmount {
	if data == nil {
		return nil
	}

	likesAmount := LikesAmount{
		Clap:   int32(data.GetClap()),
		Heart:  int32(data.GetHeart()),
		Hmm:    int32(data.GetHmm()),
		Like:   int32(data.GetLike()),
		Rocket: int32(data.GetRocket()),
		Shit:   int32(data.GetShit()),
		Stop:   int32(data.GetStop()),
	}

	return &likesAmount
}

func newsfeedRPCTagToNewsfeedTag(data *newsfeedRPC.Tag) *NewsfeedTag {
	if data == nil {
		return nil
	}

	tag := NewsfeedTag{
		ID:   data.GetID(),
		Type: newsfeedRPCEntityTypeToString(data.GetEntity()),
	}

	return &tag
}

func newsfeedRPCEntityTypeToString(data newsfeedRPC.EntityType) string {
	switch data {
	case newsfeedRPC.EntityType_Group:
		return "group"
	case newsfeedRPC.EntityType_Company:
		return "company"
	case newsfeedRPC.EntityType_Community:
		return "community"
	case newsfeedRPC.EntityType_Organization:
		return "organization"
	}
	return "user"
}

func newsfeedRPCFileToFile(data *newsfeedRPC.File) *File {
	if data == nil {
		return nil
	}

	f := File{
		Address:   data.GetURL(),
		ID:        data.GetID(),
		Mime_type: data.GetMimeType(),
		Name:      data.GetName(),
	}

	return &f
}

func newsfeedPostCommentToNewsfeedRPCComment(data *newsfeedRPC.Comment) *NewsfeedPostCommentCustom {
	if data == nil {
		return nil
	}

	post := NewsfeedPostCommentCustom{
		ID:             data.GetID(),
		Created_at:     data.GetCreatedAt(),
		Changed_at:     data.GetChangedAt(),
		Replies_amount: int32(data.GetRepliesAmount()),
		Text:           data.GetText(),
		Liked:          data.GetLiked(),
		Tags:           make([]NewsfeedTag, 0, len(data.Tags)),
		Files:          make([]File, 0, len(data.GetFiles())),
	}

	if lk := newsfeedRPCLikesAmountToLikesAmount(data.GetLikesAmount()); lk != nil {
		post.Likes_amount = *lk
	}

	for _, f := range data.GetFiles() {
		if fi := newsfeedRPCFileToFile(f); fi != nil {
			post.Files = append(post.Files, *fi)
		}
	}

	for _, t := range data.Tags {
		if ta := newsfeedRPCTagToNewsfeedTag(t); ta != nil {
			post.Tags = append(post.Tags, *ta)
		}
	}

	return &post
}

func stringToGetCommentsRequestSortOption(s string) newsfeedRPC.GetCommentsRequest_SortOption {
	switch s {
	case "amount_of_like":
		return newsfeedRPC.GetCommentsRequest_ByTopLiked
	}

	return newsfeedRPC.GetCommentsRequest_ByCreationTime
}
