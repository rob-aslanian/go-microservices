package resolver

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/notificationsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func (_ *Resolver) GetProfile(ctx context.Context, input GetProfileRequest) (*ProfileResolver, error) {
	profile, err := user.GetProfile(ctx, &userRPC.ProfileRequest{
		URL:      input.Url,
		Language: NullToString(input.Lang),
	})
	if err != nil {
		return nil, err
	}

	p := ToProfile(ctx, profile)

	// countings, _ := network.GetUserCountings(ctx, &networkRPC.User{
	// 	Id: profile.GetID(),
	// })
	//
	// p.Network_info = &NetworkInfoInUserProfile{
	// 	Connections:     countings.GetNumOfConnections(),
	// 	Followers:       countings.GetNumOfFollowers(),
	// 	Followings:      countings.GetNumOfFollowings(),
	// 	Recommendations: countings.GetNumOfReceivedRecommendations(),
	// 	Reviews:         0,
	// }

	return &ProfileResolver{
		R:        &p,
		language: NullToString(input.Lang),
	}, nil
}

func (_ *Resolver) GetProfileByID(ctx context.Context, input GetProfileByIDRequest) (*ProfileResolver, error) {
	profile, err := user.GetProfileByID(ctx, &userRPC.ID{
		ID: input.User_id,
	})
	if err != nil {
		return nil, err
	}

	p := ToProfile(ctx, profile)

	countings, _ := network.GetUserCountings(ctx, &networkRPC.User{
		Id: profile.GetID(),
	})

	p.Network_info = &NetworkInfoInUserProfile{
		Connections:     countings.GetNumOfConnections(),
		Followers:       countings.GetNumOfFollowers(),
		Followings:      countings.GetNumOfFollowings(),
		Recommendations: countings.GetNumOfReceivedRecommendations(),
		Reviews:         0,
	}

	return &ProfileResolver{
		R:        &p,
		language: NullToString(input.Lang),
	}, nil
}

func (_ *Resolver) GetMyCompanies(ctx context.Context) (*[]CompanyProfileResolver, error) {
	companies, err := user.GetMyCompanies(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}

	comp := make([]CompanyProfileResolver, len(companies.GetCompanies()))
	for i, v := range companies.GetCompanies() {
		c := toCompanyProfile(ctx, *v)
		comp[i] = CompanyProfileResolver{
			R: &c,
		}
	}

	return &comp, nil
}

func (_ *Resolver) GetOriginAvatar(ctx context.Context) (*string, error) {
	var a string
	avatar, err := user.GetOriginAvatar(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, nil
	}
	a = avatar.GetURL()

	return &a, nil
}

func (_ *Resolver) RemoveAvatar(ctx context.Context) (*SuccessResolver, error) {
	_, err := user.RemoveAvatar(ctx, &userRPC.Empty{})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeHeadline(ctx context.Context, input ChangeHeadlineRequest) (*SuccessResolver, error) {
	_, err := user.ChangeHeadline(ctx, &userRPC.Headline{Headline: input.Input})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeStory(ctx context.Context, input ChangeStoryRequest) (*SuccessResolver, error) {
	_, err := user.ChangeStory(ctx, &userRPC.Story{Story: input.Input})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) GetUserPortfolios(ctx context.Context, input GetUserPortfoliosRequest) (*UserPortfoliosResolver, error) {

	var first uint32 = 2
	var after string = "0"

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	port, err := user.GetPortfolios(ctx, &userRPC.RequestPortfolios{
		CompanyID:   NullToString(input.Company_id),
		UserID:      input.User_id,
		First:       first,
		After:       after,
		ContentType: contentTypeToRPC(input.Content_type),
	})
	if err != nil {
		return nil, err
	}

	portfolios := make([]UserPortfolio, 0, len(port.GetPortfolios()))

	for i, p := range port.GetPortfolios() {
		portfolios = append(portfolios, UserPortfolio{
			ID:         p.GetID(),
			Title:      p.GetTittle(),
			View_count: p.GetViewsCount(),
			Like_count: p.GetLikesCount(),
			Files:      make([]File, len(p.GetFiles())),
			Created_at: p.GetCreatedAt(),
			Has_liked:  p.GetHasLiked(),
		})

		/// Files
		for j, f := range p.GetFiles() {
			portfolios[i].Files[j] = File{
				Name:      f.GetName(),
				Address:   f.GetURL(),
				ID:        f.GetID(),
				Mime_type: f.GetMimeType(),
			}
		}
	}

	return &UserPortfoliosResolver{
		R: &UserPortfolios{
			Portfolio_amount: port.GetPortfolioAmount(),
			Portfolios:       portfolios,
		},
	}, nil

}

func (_ *Resolver) GetUserPortfolioInfo(ctx context.Context, input GetUserPortfolioInfoRequest) (*PortfolioInfoResolver, error) {

	res, err := user.GetUserPortfolioInfo(ctx, &userRPC.UserId{
		Id: input.User_id,
	})

	if err != nil {
		return nil, err
	}

	return &PortfolioInfoResolver{
		R: portoflioInfoRPCToPortfolioInfo(res),
	}, nil
}
func (_ *Resolver) GetUserPortfolioByID(ctx context.Context, input GetUserPortfolioByIDRequest) (*UserPortfolioResolver, error) {

	respsonse, err := user.GetPortfolioByID(ctx, &userRPC.RequestPortfolio{
		CompanyID:   NullToString(input.Company_id),
		UserID:      input.User_id,
		PortfolioID: input.Portfolio_id,
	})

	if err != nil {
		return nil, err
	}

	res := &UserPortfolio{
		ID:                   respsonse.GetID(),
		Title:                respsonse.GetTittle(),
		Description:          respsonse.GetDescription(),
		View_count:           respsonse.GetViewsCount(),
		Like_count:           respsonse.GetLikesCount(),
		Has_liked:            respsonse.GetHasLiked(),
		Saved_count:          respsonse.GetSavedCount(),
		Is_comments_disabled: respsonse.GetIsCommentDisabled(),
		Share_count:          respsonse.GetSharedCount(),
		Content_type:         userRPCContentTypeEnumToString(respsonse.GetContentType()),
		Created_at:           respsonse.GetCreatedAt(),
		Tools:                make([]string, 0, len(respsonse.GetTools())),
		Files:                make([]File, 0, len(respsonse.GetFiles())),
	}

	// Tools
	for i := range respsonse.GetTools() {
		res.Tools = append(res.Tools, respsonse.GetTools()[i])
	}

	// Files
	if len(respsonse.GetFiles()) > 0 {
		for _, f := range respsonse.GetFiles() {
			res.Files = append(res.Files, File{
				Name:      f.GetName(),
				Address:   f.GetURL(),
				ID:        f.GetID(),
				Mime_type: f.GetMimeType(),
			})
		}

	}

	return &UserPortfolioResolver{
		R: res,
	}, nil
}

func (_ *Resolver) GetAllUsersForAdmin(ctx context.Context, input GetAllUsersForAdminRequest) (*UsersResolver, error) {

	respsonse, err := user.GetAllUsersForAdmin(ctx, &userRPC.Pagination{
		First: Nullint32ToUint32(input.Pagination.First),
		After: NullToString(input.Pagination.After),
	})

	if err != nil {
		return nil, err
	}

	return &UsersResolver{
		R: &Users{
			User_amount: respsonse.GetUserAmount(),
			Users:       usersFomAdimRPCToUsers(respsonse.GetUsers()),
		},
	}, nil
}

func (_ *Resolver) ChangeUserStatus(ctx context.Context, input ChangeUserStatusRequest) (*SuccessResolver, error) {

	_, err := user.ChangeUserStatus(ctx, &userRPC.ChangeUserStatusRequest{
		UserID: input.User_id,
		Status: statusToRPC(input.Status),
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

func usersFomAdimRPCToUsers(data []*userRPC.UserForAdmin) []UserForAdmin {
	if data == nil {
		return nil
	}

	users := make([]UserForAdmin, 0, len(data))

	for _, u := range data {
		users = append(users, userFomAdimRPCToUser(u))
	}

	return users
}

func userFomAdimRPCToUser(data *userRPC.UserForAdmin) UserForAdmin {
	if data == nil {
		return UserForAdmin{}
	}

	profile := UserForAdmin{
		ID:                       data.GetID(),
		Firstname:                data.GetFirstname(),
		Lastname:                 data.GetLasttname(),
		Avatar:                   data.GetAvatar(),
		Url:                      data.GetURL(),
		Status:                   statusRPCToStatus(data.GetStatus()),
		Email:                    data.GetEmail(),
		Phone:                    data.GetPhoneNumber(),
		Birthday:                 data.GetBirthday(),
		Gender:                   userGenderRPCToGender(data.GetGender()),
		Date_of_registration:     data.GetDateOfActivation(),
		Profile_complete_percent: data.GetProfileCompletePercent(),
	}

	if data.GetLocation() == nil {
		profile.Location = &LocationProfile{} // it should ne nil
	} else {
		profile.Location = &LocationProfile{
			City:    data.GetLocation().GetCity(),
			Country: data.GetLocation().GetCountryID(),
		}
	}

	return profile

}

func statusToRPC(data string) userRPC.Status {
	if data == "" {
		return userRPC.Status_BLOCKED
	}

	switch data {
	case "ACTIVATED":
		return userRPC.Status_ACTIVATED
	case "NOT_ACTIVATED":
		return userRPC.Status_NOT_ACTIVATED
	case "DEACTIVATED":
		return userRPC.Status_DISABLED
	}

	return userRPC.Status_BLOCKED
}

func statusRPCToStatus(data userRPC.Status) string {

	switch data {
	case userRPC.Status_ACTIVATED:
		return "ACTIVATED"
	case userRPC.Status_NOT_ACTIVATED:
		return "NOT_ACTIVATED"
	case userRPC.Status_DISABLED:
		return "DEACTIVATED"
	}

	return "UNKNOWN"
}

func userGenderRPCToGender(data userRPC.GenderValue) string {
	if data.String() == "" {
		return ""
	}

	if data == userRPC.GenderValue_FEMALE {
		return "female"
	}

	return "male"
}

// Wallet

func (_ *Resolver) ContactInvitationForWallet(ctx context.Context, input ContactInvitationForWalletRequest) (*WalletResponseResolver, error) {

	resp, err := user.ContactInvitationForWallet(ctx, &userRPC.InvitationWalletRequest{
		Email:       input.Wallet_input.Email,
		Name:        input.Wallet_input.Name,
		Message:     NullToString(input.Wallet_input.Message),
		SilverCoins: input.Wallet_input.Silver_coins,
	})

	if err != nil {
		return nil, err
	}

	return &WalletResponseResolver{
		R: walletRPCToWallet(resp),
	}, nil

}

func walletRPCToWallet(data *userRPC.WalletResponse) *WalletResponse {
	if data == nil {
		return nil
	}

	return &WalletResponse{
		Amount: walletAmountRPCToWalletAmount(data.GetAmount()),
		Status: walletStatusRPCToWalletStatus(data.GetStatus()),
	}
}

func walletAmountRPCToWalletAmount(data *userRPC.WalletAmountResponse) WalletAmount {
	if data == nil {
		return WalletAmount{}
	}

	return WalletAmount{
		Gold_coins:     data.GetGoldCoins(),
		Silver_coins:   data.GetSilverCoins(),
		Pending_amount: data.GetPendingAmount(),
	}
}

func walletStatusRPCToWalletStatus(data userRPC.WalletStatusEnum) string {
	switch data {
	case userRPC.WalletStatusEnum_DONE:
		return "done"
	case userRPC.WalletStatusEnum_PENDING:
		return "pending"
	}

	return "rejected"
}

//@@@ New Portfolio @@@//

func (_ *Resolver) AddPortfolio(ctx context.Context, input AddPortfolioRequest) (*SuccessResolver, error) {

	response, err := user.AddPortfolio(ctx, &userRPC.Portfolio{
		Tittle:            input.Portfolio.Title,
		Description:       input.Portfolio.Description,
		IsCommentDisabled: input.Portfolio.Is_comment_disabled,
		Tools:             toolsToRPC(input.Portfolio.Tools),
		ContentType:       contentTypeToRPC(input.Portfolio.Content_type),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      response.GetID(),
			Success: true,
		},
	}, nil
}

// LikeUserPortfolio ...
func (_ *Resolver) LikeUserPortfolio(ctx context.Context, input LikeUserPortfolioRequest) (*SuccessResolver, error) {

	_, err := user.LikeUserPortfolio(ctx, &userRPC.PortfolioAction{
		CompanyID:   NullToString(input.Company_id),
		OwnerID:     input.Owner_id,
		PortfolioID: input.Portfolio_id,
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

// UnLikeUserPortfolio ...
func (_ *Resolver) UnLikeUserPortfolio(ctx context.Context, input UnLikeUserPortfolioRequest) (*SuccessResolver, error) {

	_, err := user.UnLikeUserPortfolio(ctx, &userRPC.PortfolioAction{
		CompanyID:   NullToString(input.Company_id),
		OwnerID:     input.Owner_id,
		PortfolioID: input.Portfolio_id,
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

// AddViewCountToPortfolio ...
func (_ *Resolver) AddViewCountToPortfolio(ctx context.Context, input AddViewCountToPortfolioRequest) (*SuccessResolver, error) {

	_, err := user.AddViewCountToPortfolio(ctx, &userRPC.PortfolioAction{
		CompanyID:   NullToString(input.Company_id),
		OwnerID:     input.Owner_id,
		PortfolioID: input.Portfolio_id,
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

// AddSavedCountToPortfolio ...
func (_ *Resolver) AddSavedCountToPortfolio(ctx context.Context, input AddSavedCountToPortfolioRequest) (*SuccessResolver, error) {

	_, err := user.AddSavedCountToPortfolio(ctx, &userRPC.PortfolioAction{
		OwnerID:     input.Owner_id,
		PortfolioID: input.Portfolio_id,
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

// GetUserPortfolioComments ...
func (r *Resolver) GetUserPortfolioComments(ctx context.Context, input GetUserPortfolioCommentsRequest) (*UserPortfolioCommentResponseResolver, error) {

	var first uint32 = 2
	var after string = "0"

	if input.Pagination.First != nil {
		first = uint32(*input.Pagination.First)
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	response, err := user.GetUserPortfolioComments(ctx, &userRPC.GetPortfolioComment{
		PortfolioID: input.Portfolio_id,
		First:       first,
		After:       after,
	})

	if err != nil {
		return nil, err
	}

	comments := make([]UserPortfolioComment, 0, len(response.GetPortfolioComment()))

	// log.Printf("Comments %+v" , comments)
	// log.Printf("Comments from service %+v" , response.GetPortfolioComment())

	for _, comment := range response.GetPortfolioComment() {

		if c := portfolioCommentToRPC(comment); c != nil {
			log.Printf("Comment %+v", *c)

			c.User_profile = &Profile{}
			c.Company_profile = &CompanyProfile{}

			comments = append(comments, *c)
		}
	}

	userIDs := make([]string, 0, len(response.GetPortfolioComment()))
	companiesIDs := make([]string, 0, len(response.GetPortfolioComment()))

	for _, c := range response.GetPortfolioComment() {
		if c.GetCompanyID() != "" {
			companiesIDs = append(companiesIDs, c.GetCompanyID())
		} else {
			userIDs = append(userIDs, c.GetUserID())
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

	for _, p := range response.GetPortfolioComment() {
		for i := range comments {

			if comments[i].ID == p.GetID() {

				// company profile
				if p.GetCompanyID() != "" {
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == p.GetCompanyID() {
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
					profile := ToProfile(ctx, userResp.GetProfiles()[p.GetUserID()])
					comments[i].User_profile = &profile

				}

			}

		}
	}

	return &UserPortfolioCommentResponseResolver{
		R: &UserPortfolioCommentResponse{
			Comments_amount: response.GetCommentAmount(),
			Comments:        comments,
		},
	}, nil

}

func portfolioCommentToRPC(data *userRPC.PortfolioComment) *UserPortfolioComment {
	if data == nil {
		return nil
	}

	return &UserPortfolioComment{
		ID:         data.GetID(),
		Comment:    data.GetComment(),
		Created_at: data.GetCreatedAt(),
	}
}

// AddCommentToPortfolio ...
func (_ *Resolver) AddCommentToPortfolio(ctx context.Context, input AddCommentToPortfolioRequest) (*SuccessResolver, error) {

	id, err := user.AddCommentToPortfolio(ctx, &userRPC.PortfolioComment{
		CompanyID:   NullToString(input.Comment.Company_id),
		OwnerID:     input.Comment.Owner_id,
		PortfolioID: input.Comment.Portfolio_id,
		Comment:     input.Comment.Comment,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      id.ID,
			Success: true,
		},
	}, nil
}

// RemoveCommentInPortfolio
func (_ *Resolver) RemoveCommentInPortfolio(ctx context.Context, input RemoveCommentInPortfolioRequest) (*SuccessResolver, error) {

	_, err := user.RemoveCommentInPortfolio(ctx, &userRPC.RemovePortfolioComment{
		CommentID:   input.Comment_id,
		PortfolioID: input.Portfolio_id,
		CompanyID:   NullToString(input.Company_id),
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

func userRPCContentTypeEnumToString(data userRPC.ContentTypeEnum) string {
	content := "Photo"

	switch data {
	case userRPC.ContentTypeEnum_Content_Type_Article:
		content = "Article"
	case userRPC.ContentTypeEnum_Content_Type_Video:
		content = "Video"
	case userRPC.ContentTypeEnum_Content_Type_Audio:
		content = "Audio"
	}

	return content
}

func contentTypeToRPC(ctp string) userRPC.ContentTypeEnum {
	switch ctp {
	case "Video":
		return userRPC.ContentTypeEnum_Content_Type_Video
	case "Article":
		return userRPC.ContentTypeEnum_Content_Type_Article
	case "Audio":
		return userRPC.ContentTypeEnum_Content_Type_Audio
	}

	return userRPC.ContentTypeEnum_Content_Type_Photo
}

func toolsToRPC(t *[]string) []string {
	if t == nil {
		return nil
	}

	tools := make([]string, 0, len(*t))

	for _, tool := range *t {
		tools = append(tools, tool)
	}

	return tools

}

func filesToRPC(ids *[]string) []*userRPC.File {
	if ids == nil {
		return nil
	}
	files := make([]*userRPC.File, 0, len(*ids))

	for i := range *ids {
		files = append(files, &userRPC.File{ID: (*ids)[i]})
	}

	return files
}

func (_ *Resolver) RemovePortfolio(ctx context.Context, input RemovePortfolioRequest) (*SuccessResolver, error) {
	_, err := user.RemovePortfolio(ctx, &userRPC.Portfolio{
		ID: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeOrderFilesInPortfolio(ctx context.Context, input ChangeOrderFilesInPortfolioRequest) (*SuccessResolver, error) {
	_, err := user.ChangeOrderFilesInPortfolio(ctx, &userRPC.PortfolioFile{
		ID:       input.File.ID,
		FileID:   input.File.FileID,
		Position: uint32(input.File.Position),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangePortfolio(ctx context.Context, input ChangePortfolioRequest) (*SuccessResolver, error) {
	exp := userRPC.Portfolio{
		ID:                input.ID,
		Tittle:            input.Portfolio.Title,
		Description:       input.Portfolio.Description,
		IsCommentDisabled: input.Portfolio.Is_comment_disabled,
	}

	// Tools
	if input.Portfolio.Tools != nil {
		exp.Tools = make([]string, 0, len(*input.Portfolio.Tools))
		for _, tool := range *input.Portfolio.Tools {
			exp.Tools = append(exp.Tools, tool)
		}
	}

	_, err := user.ChangePortfolio(ctx, &exp)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddLinksInPortfolio(ctx context.Context, input AddLinksInPortfolioRequest) (*SuccessResolver, error) {
	_, err := user.AddLinksInPortfolio(ctx, &userRPC.AddLinksRequest{
		ID: input.ID,
		Links: func(in []LinkInput) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					URL: in[r].Url,
				}
			}
			return ar
		}(input.Input),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeLinkInPortfolio(ctx context.Context, input ChangeLinkInPortfolioRequest) (*SuccessResolver, error) {
	_, err := user.ChangeLinkInPortfolio(ctx, &userRPC.ChangeLinkRequest{
		ID: input.ID,
		Link: &userRPC.Link{
			ID:  input.Link_id,
			URL: input.Url,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveLinksInPortfolio(ctx context.Context, input RemoveLinksInPortfolioRequest) (*SuccessResolver, error) {
	_, err := user.RemoveLinksInPortfolio(ctx, &userRPC.RemoveLinksRequest{
		ID: input.ID,

		Links: func(in []string) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					ID: in[r],
				}
			}
			return ar
		}(input.Links_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveFilesInPortfolio(ctx context.Context, input RemoveFilesInPortfolioRequest) (*SuccessResolver, error) {
	_, err := user.RemoveFilesInPortfolio(ctx, &userRPC.RemoveFilesRequest{
		ID: input.ID,

		Files: func(in []string) []*userRPC.File {
			ar := make([]*userRPC.File, len(in))
			for r := range in {
				ar[r] = &userRPC.File{
					ID: in[r],
				}
			}

			return ar
		}(input.Files_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

// AddToolTechnology ...
func (_ *Resolver) AddToolTechnology(ctx context.Context, input AddToolTechnologyRequest) (IDsResolver, error) {
	response, err := user.AddToolTechnology(ctx, &userRPC.ToolTechnologyList{
		ToolsTechnologies: func(in []ToolTechnologyInput) []*userRPC.ToolTechnology {
			ar := make([]*userRPC.ToolTechnology, len(in))
			for r := range in {
				ar[r] = &userRPC.ToolTechnology{
					ToolTechnology: in[r].Tool_Technology,
					Rank:           stringToUserRPCToolTechnologyLevel(in[r].Rank),
				}
			}
			return ar
		}(input.Tools_technologies),
	})
	if err != nil {
		return IDsResolver{}, err
	}

	return IDsResolver{
		R: &IDs{
			Ids: response.GetIDs(),
		},
	}, nil
}

// ChangeToolTechnology ...
func (_ *Resolver) ChangeToolTechnology(ctx context.Context, input ChangeToolTechnologyRequest) (*SuccessResolver, error) {
	tools := userRPC.ToolTechnologyList{}

	for _, tool := range input.Tools_technologies {
		tools.ToolsTechnologies = append(tools.ToolsTechnologies, &userRPC.ToolTechnology{
			Rank:           stringToUserRPCToolTechnologyLevel(tool.Rank),
			ToolTechnology: tool.Tool_Technology,
			ID:             tool.ID,
		})
	}
	_, err := user.ChangeToolTechnology(ctx, &tools)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}
func (_ *Resolver) RemoveToolTechnology(ctx context.Context, input RemoveToolTechnologyRequest) (*SuccessResolver, error) {
	_, err := user.RemoveToolTechnology(ctx, &userRPC.ToolTechnologyList{
		ToolsTechnologies: func(in []string) []*userRPC.ToolTechnology {
			ar := make([]*userRPC.ToolTechnology, len(in))
			for r := range in {
				ar[r] = &userRPC.ToolTechnology{
					ID: in[r],
				}
			}
			return ar
		}(input.ID),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddExperience(ctx context.Context, input AddExperienceRequest) (*SuccessResolver, error) {
	response, err := user.AddExperience(ctx, &userRPC.Experience{
		Position: input.Experience.Position,
		Location: &userRPC.Location{
			City: &userRPC.City{
				Id: NullIDToInt32(input.Experience.City_id),
			},
		},
		Company:       NullToString(input.Experience.Company),
		CurrentlyWork: NullBoolToBool(input.Experience.Currently_work),
		Description:   NullToString(input.Experience.Description),
		FinishDate:    NullToString(input.Experience.Finish_date),
		StartDate:     NullToString(input.Experience.Start_date),

		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Experience.Files_id),

		Links: func(in *[]LinkInput) []*userRPC.Link {
			if in != nil {
				ar := make([]*userRPC.Link, len(*in))
				for r := range *in {
					ar[r] = &userRPC.Link{
						URL: (*in)[r].Url,
					}
				}

				return ar
			}
			return []*userRPC.Link{}
		}(input.Experience.Links),
	})

	if err != nil {
		return nil, err
	}

	success := &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}

	if response.GetID() != "" {
		// It should be a pointer!!!
		success.R.ID = response.GetID()
	}

	return success, nil
}

func (_ *Resolver) RemoveExperience(ctx context.Context, input RemoveExperienceRequest) (*SuccessResolver, error) {
	_, err := user.RemoveExperience(ctx, &userRPC.Experience{
		ID: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeExperience(ctx context.Context, input ChangeExperienceRequest) (*SuccessResolver, error) {
	exp := userRPC.Experience{
		ID:       input.ID,
		Position: NullToString(input.Experience.Position),
		Location: &userRPC.Location{
			City: &userRPC.City{
				Id: NullIDToInt32(input.Experience.City_id),
			},
		},
		Company:       NullToString(input.Experience.Company),
		CurrentlyWork: NullBoolToBool(input.Experience.Currently_work),
		Description:   NullToString(input.Experience.Description),
		FinishDate:    NullToString(input.Experience.Finish_date),
		StartDate:     NullToString(input.Experience.Start_date),
		IsCurrentlyWorkNull: func() bool {
			if input.Experience.Currently_work == nil {
				return true
			}
			return false
		}(),
	}

	if input.Experience.City_id == nil {
		exp.IsLocationNull = true
	}

	if input.Experience.Description == nil {
		exp.IsDescriptionNull = true
	}

	if input.Experience.Links != nil {
		exp.Links = make([]*userRPC.Link, 0, len(*input.Experience.Links))
		for _, link := range *input.Experience.Links {
			exp.Links = append(exp.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	_, err := user.ChangeExperience(ctx, &exp)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeLinkInExperience(ctx context.Context, input ChangeLinkInExperienceRequest) (*SuccessResolver, error) {
	_, err := user.ChangeLinkInExperience(ctx, &userRPC.ChangeLinkRequest{
		ID: input.ID,
		Link: &userRPC.Link{
			ID:  input.Link_id,
			URL: input.Url,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveLinksInExperience(ctx context.Context, input RemoveLinksInExperienceRequest) (*SuccessResolver, error) {
	_, err := user.RemoveLinksInExperience(ctx, &userRPC.RemoveLinksRequest{
		ID: input.ID,

		Links: func(in []string) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					ID: in[r],
				}
			}
			return ar
		}(input.Links_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveFilesInExperience(ctx context.Context, input RemoveFilesInExperienceRequest) (*SuccessResolver, error) {
	_, err := user.RemoveFilesInExperience(ctx, &userRPC.RemoveFilesRequest{
		ID: input.ID,

		Files: func(in []string) []*userRPC.File {
			ar := make([]*userRPC.File, len(in))
			for r := range in {
				ar[r] = &userRPC.File{
					ID: in[r],
				}
			}

			return ar
		}(input.Files_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddLinksInExperience(ctx context.Context, input AddLinksInExperienceRequest) (*SuccessResolver, error) {
	_, err := user.AddLinksInExperience(ctx, &userRPC.AddLinksRequest{
		ID: input.ID,
		Links: func(in []LinkInput) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					URL: in[r].Url,
				}
			}
			return ar
		}(input.Input),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) GetListOfUnuploadFilesInExperience(ctx context.Context) (*[]FileResolver, error) {
	files, err := user.GetUploadedFilesInExperience(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}

	if len(files.GetFiles()) == 0 {
		return nil, nil
	}

	f := make([]FileResolver, len(files.GetFiles()))

	for i := range files.GetFiles() {
		f[i].R = &File{
			ID:        files.GetFiles()[i].GetID(),
			Name:      files.GetFiles()[i].GetName(),
			Address:   files.GetFiles()[i].GetURL(),
			Mime_type: files.GetFiles()[i].GetMimeType(),
		}
	}

	return &f, nil
}

func (_ *Resolver) GetListOfUnuploadFilesInEducation(ctx context.Context) (*[]FileResolver, error) {
	files, err := user.GetUploadedFilesInEducation(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}

	if len(files.GetFiles()) == 0 {
		return nil, nil
	}

	f := make([]FileResolver, len(files.GetFiles()))

	for i := range files.GetFiles() {
		f[i].R = &File{
			ID:        files.GetFiles()[i].GetID(),
			Name:      files.GetFiles()[i].GetName(),
			Address:   files.GetFiles()[i].GetURL(),
			Mime_type: files.GetFiles()[i].GetMimeType(),
		}
	}

	return &f, nil
}

func (_ *Resolver) GetUnuploadImageInInterest(ctx context.Context) (*FileResolver, error) {
	file, err := user.GetUnuploadImageInInterest(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}

	if file.GetURL() == "" {
		return nil, nil
	}

	var f FileResolver
	f.R = &File{
		ID:        file.GetID(),
		Name:      file.GetName(),
		Address:   file.GetURL(),
		Mime_type: file.GetMimeType(),
	}

	return &f, nil
}

func (_ *Resolver) AddEducation(ctx context.Context, input AddEducationRequest) (*SuccessResolver, error) {
	id, err := user.AddEducation(ctx, &userRPC.Education{
		School:           input.Education.School,
		Degree:           NullToString(input.Education.Degree),
		FieldStudy:       input.Education.Field_study,
		Grade:            NullToString(input.Education.Grade),
		Description:      NullToString(input.Education.Description),
		FinishDate:       NullToString(input.Education.Finish_date),
		StartDate:        input.Education.Start_date,
		IsCurrentlyStudy: input.Education.Currently_study,
		Location: &userRPC.Location{
			City: &userRPC.City{
				Id: NullIDToInt32(input.Education.City_id),
			},
		},
		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Education.Files_id),

		Links: func(in *[]LinkInput) []*userRPC.Link {
			if in != nil {
				ar := make([]*userRPC.Link, len(*in))
				for r := range *in {
					ar[r] = &userRPC.Link{
						URL: (*in)[r].Url,
					}
				}
				return ar
			}
			return []*userRPC.Link{}
		}(input.Education.Links),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) RemoveEducation(ctx context.Context, input RemoveEducationRequest) (*SuccessResolver, error) {
	_, err := user.RemoveEducation(ctx, &userRPC.Education{
		ID: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeEducation(ctx context.Context, input ChangeEducationRequest) (*SuccessResolver, error) {
	edu := userRPC.Education{
		ID: input.ID,
		Location: &userRPC.Location{
			City: &userRPC.City{
				Id: NullIDToInt32(input.Education.City_id),
			},
		},
		School:           NullToString(input.Education.School),
		Degree:           NullToString(input.Education.Degree),
		FieldStudy:       NullToString(input.Education.Field_study),
		Grade:            NullToString(input.Education.Grade),
		FinishDate:       NullToString(input.Education.Finish_date),
		StartDate:        NullToString(input.Education.Start_date),
		IsCurrentlyStudy: input.Education.Currently_study,
		Description:      NullToString(input.Education.Description),
		// Links:            make([]*userRPC.Link, 0, len(input.Education.Links)),
	}

	if input.Education.Degree == nil {
		edu.IsDegreeNull = true
	}

	if input.Education.Grade == nil {
		edu.IsGradeNull = true
	}

	if input.Education.Description == nil {
		edu.IsDescriptionNull = true
	}

	if input.Education.Links != nil {
		edu.Links = make([]*userRPC.Link, 0, len(*input.Education.Links))
		for _, link := range *input.Education.Links {
			edu.Links = append(edu.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	_, err := user.ChangeEducation(ctx, &edu)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeLinkInEducation(ctx context.Context, input ChangeLinkInEducationRequest) (*SuccessResolver, error) {
	_, err := user.ChangeLinkInEducation(ctx, &userRPC.ChangeLinkRequest{
		ID: input.ID,
		Link: &userRPC.Link{
			ID:  input.Link_id,
			URL: input.Url,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveLinksInEducation(ctx context.Context, input RemoveLinksInEducationRequest) (*SuccessResolver, error) {
	_, err := user.RemoveLinksInEducation(ctx, &userRPC.RemoveLinksRequest{
		ID: input.ID,

		Links: func(in []string) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					ID: in[r],
				}
			}
			return ar
		}(input.Links_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveFilesInEducation(ctx context.Context, input RemoveFilesInEducationRequest) (*SuccessResolver, error) {
	_, err := user.RemoveFilesInEducation(ctx, &userRPC.RemoveFilesRequest{
		ID: input.ID,

		Files: func(in []string) []*userRPC.File {
			ar := make([]*userRPC.File, len(in))
			for r := range in {
				ar[r] = &userRPC.File{
					ID: in[r],
				}
			}
			return ar
		}(input.Files_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddLinksInEducation(ctx context.Context, input AddLinksInEducationRequest) (*SuccessResolver, error) {
	response, err := user.AddLinksInEducation(ctx, &userRPC.AddLinksRequest{
		ID: input.ID,
		Links: func(in []LinkInput) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					URL: in[r].Url,
				}
			}
			return ar
		}(input.Input),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
			ID:      response.GetID(),
		}}, nil
}

func (_ *Resolver) AddSkills(ctx context.Context, input AddSkillsRequest) (*SuccessResolver, error) {
	result, err := user.AddSkills(ctx, &userRPC.SkillList{
		Skills: func(in []SkillInput) []*userRPC.Skill {
			ar := make([]*userRPC.Skill, len(in))
			for r := range in {
				ar[r] = &userRPC.Skill{
					Skill: in[r].Skill,
				}
			}
			return ar
		}(input.Skills),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      result.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) ChangeOrderOfSkill(ctx context.Context, input ChangeOrderOfSkillRequest) (*SuccessResolver, error) {
	_, err := user.ChangeOrderOfSkill(ctx, &userRPC.Skill{
		ID:       input.Skill.ID,
		Position: uint32(input.Skill.Position),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveSkills(ctx context.Context, input RemoveSkillsRequest) (*SuccessResolver, error) {
	_, err := user.RemoveSkills(ctx, &userRPC.SkillList{
		Skills: func(in []string) []*userRPC.Skill {
			ar := make([]*userRPC.Skill, len(in))
			for r := range in {
				ar[r] = &userRPC.Skill{
					ID: in[r],
				}
			}
			return ar
		}(input.Skill_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) VerifySkill(ctx context.Context, input VerifySkillRequest) (*SuccessResolver, error) {
	_, err := user.VerifySkill(ctx, &userRPC.VerifySkillRequest{
		UserID: input.User_id,
		Skill: &userRPC.Skill{
			ID: input.Skill_id,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) UnverifySkill(ctx context.Context, input VerifySkillRequest) (*SuccessResolver, error) {
	_, err := user.UnverifySkill(ctx, &userRPC.VerifySkillRequest{
		UserID: input.User_id,
		Skill: &userRPC.Skill{
			ID: input.Skill_id,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddKnownLanguage(ctx context.Context, input AddKnownLanguageRequest) (*SuccessResolver, error) {
	id, err := user.AddKnownLanguage(ctx, &userRPC.KnownLanguage{
		Language: input.Language.Language_id,
		Rank:     uint32(input.Language.Rank),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) ChangeKnownLanguage(ctx context.Context, input ChangeKnownLanguageRequest) (*SuccessResolver, error) {
	_, err := user.ChangeKnownLanguage(ctx, &userRPC.KnownLanguage{
		ID:       input.ID,
		Language: NullToString(input.Language.Language_id),
		Rank:     uint32(NullToInt32(input.Language.Rank)),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveKnownLanguage(ctx context.Context, input RemoveKnownLanguageRequest) (*SuccessResolver, error) {
	_, err := user.RemoveKnownLanguage(ctx, &userRPC.KnownLanguage{
		ID: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

// AddLinksInAccomplishment ...
func (_ *Resolver) AddLinksInAccomplishment(ctx context.Context, input AddLinksInAccomplishmentRequest) (*SuccessResolver, error) {
	_, err := user.AddLinksInAccomplishment(ctx, &userRPC.AddLinksRequest{
		ID: input.ID,
		Links: func(in []LinkInput) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					URL: in[r].Url,
				}
			}
			return ar
		}(input.Input),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveLinksInAccomplishment(ctx context.Context, input RemoveLinksInAccomplishmentRequest) (*SuccessResolver, error) {
	_, err := user.RemoveLinksInAccomplishment(ctx, &userRPC.RemoveLinksRequest{
		ID: input.ID,

		Links: func(in []string) []*userRPC.Link {
			ar := make([]*userRPC.Link, len(in))
			for r := range in {
				ar[r] = &userRPC.Link{
					ID: in[r],
				}
			}
			return ar
		}(input.Links_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

// RemoveFilesInAccomplishment ...
func (_ *Resolver) RemoveFilesInAccomplishment(ctx context.Context, input RemoveFilesInAccomplishmentRequest) (*SuccessResolver, error) {
	_, err := user.RemoveFilesInAccomplishment(ctx, &userRPC.RemoveFilesRequest{
		ID: input.ID,

		Files: func(in []string) []*userRPC.File {
			ar := make([]*userRPC.File, len(in))
			for r := range in {
				ar[r] = &userRPC.File{
					ID: in[r],
				}
			}

			return ar
		}(input.Files_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddAccomplishmentPublication(ctx context.Context, input AddAccomplishmentPublicationRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		Name: input.Input.Title,
		Type: userRPC.Accomplishment_Publication,
		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Input.Files_id),
	}

	if input.Input.Publisher == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Input.Publisher
	}

	if input.Input.Date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Input.Date
	}

	if input.Input.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Input.Description
	}

	if input.Input.Url == nil {
		acc.IsURLNull = true
	} else {
		acc.URL = *input.Input.Url
	}

	id, err := user.AddAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) AddAccomplishmentCertification(ctx context.Context, input AddAccomplishmentCertificationRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		Name:     input.Input.Name,
		Type:     userRPC.Accomplishment_Certificate,
		IsExpire: input.Input.Is_expire,
		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Input.Files_id),

		Links: func(in *[]LinkInput) []*userRPC.Link {
			if in != nil {
				ar := make([]*userRPC.Link, len(*in))
				for r := range *in {
					ar[r] = &userRPC.Link{
						URL: (*in)[r].Url,
					}
				}

				return ar
			}
			return []*userRPC.Link{}
		}(input.Input.Link),
	}

	if input.Input.Certification_authority == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Input.Certification_authority
	}

	if input.Input.Start_date == "" {
		acc.IsStartDateNull = true
	} else {
		acc.StartDate = input.Input.Start_date
	}

	if input.Input.Finish_date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Input.Finish_date
	}

	if input.Input.Is_expire == false {
		acc.IsIsExpireNull = true
	} else {
		acc.IsExpire = input.Input.Is_expire
	}

	if input.Input.License_number == nil {
		acc.IsLicenseNumberNull = true
	} else {
		acc.LicenseNumber = *input.Input.License_number
	}

	if input.Input.Url == nil {
		acc.IsURLNull = true
	} else {
		acc.URL = *input.Input.Url
	}

	id, err := user.AddAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) AddAccomplishmentLicense(ctx context.Context, input AddAccomplishmentLicenseRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		Name:     input.Input.Name,
		Type:     userRPC.Accomplishment_License,
		IsExpire: input.Input.Is_expire,
		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Input.Files_id),

		Links: func(in *[]LinkInput) []*userRPC.Link {
			if in != nil {
				ar := make([]*userRPC.Link, len(*in))
				for r := range *in {
					ar[r] = &userRPC.Link{
						URL: (*in)[r].Url,
					}
				}

				return ar
			}
			return []*userRPC.Link{}
		}(input.Input.Link),
	}

	if input.Input.Start_date == "" {
		acc.IsStartDateNull = true
	} else {
		acc.StartDate = input.Input.Start_date
	}

	if input.Input.Finish_date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Input.Finish_date
	}

	if input.Input.Is_expire == false {
		acc.IsIsExpireNull = true
	} else {
		acc.IsExpire = input.Input.Is_expire
	}

	if input.Input.Issuer == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Input.Issuer
	}

	if input.Input.License_number == nil {
		acc.IsLicenseNumberNull = true
	} else {
		acc.LicenseNumber = *input.Input.License_number
	}

	id, err := user.AddAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) AddAccomplishmentAward(ctx context.Context, input AddAccomplishmentAwardRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		Name: input.Input.Title,
		Type: userRPC.Accomplishment_Award,
		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Input.Files_id),

		Links: func(in *[]LinkInput) []*userRPC.Link {
			if in != nil {
				ar := make([]*userRPC.Link, len(*in))
				for r := range *in {
					ar[r] = &userRPC.Link{
						URL: (*in)[r].Url,
					}
				}

				return ar
			}
			return []*userRPC.Link{}
		}(input.Input.Link),
	}

	if input.Input.Issuer == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Input.Issuer
	}

	if input.Input.Date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Input.Date
	}

	if input.Input.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Input.Description
	}

	id, err := user.AddAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) AddAccomplishmentProject(ctx context.Context, input AddAccomplishmentProjectRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		Name: input.Input.Name,
		Type: userRPC.Accomplishment_Project,
		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Input.Files_id),

		Links: func(in *[]LinkInput) []*userRPC.Link {
			if in != nil {
				ar := make([]*userRPC.Link, len(*in))
				for r := range *in {
					ar[r] = &userRPC.Link{
						URL: (*in)[r].Url,
					}
				}

				return ar
			}
			return []*userRPC.Link{}
		}(input.Input.Link),
	}

	if input.Input.Start_date == "" {
		acc.IsStartDateNull = true
	} else {
		acc.StartDate = input.Input.Start_date
	}

	if input.Input.Finish_date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Input.Finish_date
	}

	if input.Input.Url == nil {
		acc.IsURLNull = true
	} else {
		acc.URL = *input.Input.Url
	}

	if input.Input.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Input.Description
	}

	id, err := user.AddAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) AddAccomplishmentTest(ctx context.Context, input AddAccomplishmentTestRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		Name:  input.Input.Title,
		Type:  userRPC.Accomplishment_Test,
		Score: float32(input.Input.Score),
		Files: func(ids *[]string) []*userRPC.File {
			if ids == nil {
				return nil
			}
			files := make([]*userRPC.File, 0, len(*ids))
			for i := range *ids {
				files = append(files, &userRPC.File{ID: (*ids)[i]})
			}

			return files
		}(input.Input.Files_id),

		Links: func(in *[]LinkInput) []*userRPC.Link {
			if in != nil {
				ar := make([]*userRPC.Link, len(*in))
				for r := range *in {
					ar[r] = &userRPC.Link{
						URL: (*in)[r].Url,
					}
				}

				return ar
			}
			return []*userRPC.Link{}
		}(input.Input.Link),
	}

	if input.Input.Date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Input.Date
	}

	if input.Input.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Input.Description
	}

	id, err := user.AddAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) RemoveAccomplishment(ctx context.Context, input RemoveAccomplishmentRequest) (*SuccessResolver, error) {
	_, err := user.RemoveAccomplishment(ctx, &userRPC.Accomplishment{
		ID: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeAccomplishmentCertification(ctx context.Context, input ChangeAccomplishmentCertificationRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		ID:       input.ID,
		IsExpire: input.Accomplishment.Is_expire,
		Type:     userRPC.Accomplishment_Certificate,
	}

	if input.Accomplishment.Link != nil {
		acc.Links = make([]*userRPC.Link, 0, len(*input.Accomplishment.Link))
		for _, link := range *input.Accomplishment.Link {
			acc.Links = append(acc.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	if input.Accomplishment.Name == nil {
		acc.IsNameNull = true
	} else {
		acc.Name = *input.Accomplishment.Name
	}

	// if input.Accomplishment.Is_expire == false {
	// 	acc.IsIsExpireNull = true
	// } else {
	// 	acc.IsExpire = input.Accomplishment.Is_expire
	// }

	// if input.Accomplishment.Finish_date == nil {
	// 	acc.IsFinishDateNull = true
	// } else {
	// 	acc.FinishDate = *input.Accomplishment.Finish_date
	// }

	if input.Accomplishment.Is_expire == true {
		acc.IsExpire = input.Accomplishment.Is_expire
		if input.Accomplishment.Is_expire == false {
			acc.FinishDate = *input.Accomplishment.Finish_date
			acc.IsIsExpireNull = true
		} else {
			acc.IsExpire = input.Accomplishment.Is_expire
			acc.IsFinishDateNull = true
		}
	}

	if input.Accomplishment.Start_date == "" {
		acc.IsStartDateNull = true
	} else {
		acc.StartDate = input.Accomplishment.Start_date
	}

	if input.Accomplishment.Certification_authority == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Accomplishment.Certification_authority
	}

	if input.Accomplishment.License_number == nil {
		acc.IsLicenseNumberNull = true
	} else {
		acc.LicenseNumber = *input.Accomplishment.License_number
	}

	if input.Accomplishment.Url == nil {
		acc.IsURLNull = true
	} else {
		acc.URL = *input.Accomplishment.Url
	}

	_, err := user.ChangeAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeAccomplishmentLicense(ctx context.Context, input ChangeAccomplishmentLicenseRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		ID:       input.ID,
		IsExpire: input.Accomplishment.Is_expire,
		Type:     userRPC.Accomplishment_License,
	}

	if input.Accomplishment.Link != nil {
		acc.Links = make([]*userRPC.Link, 0, len(*input.Accomplishment.Link))
		for _, link := range *input.Accomplishment.Link {
			acc.Links = append(acc.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	if input.Accomplishment.Name == nil {
		acc.IsNameNull = true
	} else {
		acc.Name = *input.Accomplishment.Name
	}

	// if input.Accomplishment.Is_expire == false {
	// 	acc.IsIsExpireNull = true
	// } else {
	// 	acc.IsExpire = input.Accomplishment.Is_expire
	// }

	// if input.Accomplishment.Finish_date == nil {
	// 	acc.IsFinishDateNull = true
	// } else {
	// 	acc.FinishDate = *input.Accomplishment.Finish_date
	// }

	if input.Accomplishment.Is_expire == true {
		acc.IsExpire = input.Accomplishment.Is_expire
		if input.Accomplishment.Is_expire == false {
			acc.FinishDate = *input.Accomplishment.Finish_date
			acc.IsIsExpireNull = true
		} else {
			acc.IsExpire = input.Accomplishment.Is_expire
			acc.IsFinishDateNull = true
		}
	}

	if input.Accomplishment.Start_date == "" {
		acc.IsStartDateNull = true
	} else {
		acc.StartDate = input.Accomplishment.Start_date
	}

	if input.Accomplishment.Finish_date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Accomplishment.Finish_date
	}

	if input.Accomplishment.Issuer == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Accomplishment.Issuer
	}

	if input.Accomplishment.License_number == nil {
		acc.IsLicenseNumberNull = true
	} else {
		acc.LicenseNumber = *input.Accomplishment.License_number
	}

	_, err := user.ChangeAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeAccomplishmentAward(ctx context.Context, input ChangeAccomplishmentAwardRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		ID:   input.ID,
		Type: userRPC.Accomplishment_Award,
	}

	if input.Accomplishment.Link != nil {
		acc.Links = make([]*userRPC.Link, 0, len(*input.Accomplishment.Link))
		for _, link := range *input.Accomplishment.Link {
			acc.Links = append(acc.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	if input.Accomplishment.Title == nil {
		acc.IsNameNull = true
	} else {
		acc.Name = *input.Accomplishment.Title
	}

	if input.Accomplishment.Issuer == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Accomplishment.Issuer
	}

	if input.Accomplishment.Date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Accomplishment.Date
	}

	if input.Accomplishment.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Accomplishment.Description
	}

	_, err := user.ChangeAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeAccomplishmentProject(ctx context.Context, input ChangeAccomplishmentProjectRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		ID:   input.ID,
		Type: userRPC.Accomplishment_Project,
	}

	if input.Accomplishment.Link != nil {
		acc.Links = make([]*userRPC.Link, 0, len(*input.Accomplishment.Link))
		for _, link := range *input.Accomplishment.Link {
			acc.Links = append(acc.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	if input.Accomplishment.Name == nil {
		acc.IsNameNull = true
	} else {
		acc.Name = *input.Accomplishment.Name
	}

	if input.Accomplishment.Start_date == "" {
		acc.IsStartDateNull = true
	} else {
		acc.StartDate = input.Accomplishment.Start_date
	}

	if input.Accomplishment.Finish_date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Accomplishment.Finish_date
	}

	if input.Accomplishment.Url == nil {
		acc.IsURLNull = true
	} else {
		acc.URL = *input.Accomplishment.Url
	}

	if input.Accomplishment.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Accomplishment.Description
	}

	if input.Accomplishment.Start_date == "" {
		acc.IsStartDateNull = true
	} else {
		acc.StartDate = input.Accomplishment.Start_date
	}

	_, err := user.ChangeAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeAccomplishmentPublication(ctx context.Context, input ChangeAccomplishmentPublicationRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		ID:   input.ID,
		Type: userRPC.Accomplishment_Publication,
	}

	if input.Accomplishment.Link != nil {
		acc.Links = make([]*userRPC.Link, 0, len(*input.Accomplishment.Link))
		for _, link := range *input.Accomplishment.Link {
			acc.Links = append(acc.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	if input.Accomplishment.Title == nil {
		acc.IsNameNull = true
	} else {
		acc.Name = *input.Accomplishment.Title
	}

	if input.Accomplishment.Publisher == nil {
		acc.IsIssuerNull = true
	} else {
		acc.Issuer = *input.Accomplishment.Publisher
	}

	if input.Accomplishment.Date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Accomplishment.Date
	}

	if input.Accomplishment.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Accomplishment.Description
	}

	if input.Accomplishment.Url == nil {
		acc.IsURLNull = true
	} else {
		acc.URL = *input.Accomplishment.Url
	}

	_, err := user.AddAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeAccomplishmentTest(ctx context.Context, input ChangeAccomplishmentTestRequest) (*SuccessResolver, error) {
	acc := userRPC.Accomplishment{
		ID:   input.ID,
		Type: userRPC.Accomplishment_Test,
	}

	if input.Accomplishment.Link != nil {
		acc.Links = make([]*userRPC.Link, 0, len(*input.Accomplishment.Link))
		for _, link := range *input.Accomplishment.Link {
			acc.Links = append(acc.Links, &userRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	if input.Accomplishment.Title == nil {
		acc.IsNameNull = true
	} else {
		acc.Name = *input.Accomplishment.Title
	}

	if input.Accomplishment.Date == nil {
		acc.IsFinishDateNull = true
	} else {
		acc.FinishDate = *input.Accomplishment.Date
	}

	if input.Accomplishment.Description == nil {
		acc.IsDescriptionNull = true
	} else {
		acc.Description = *input.Accomplishment.Description
	}

	if input.Accomplishment.Score == nil {
		acc.IsScoreNull = true
	} else {
		acc.Score = float32(*input.Accomplishment.Score)
	}

	_, err := user.ChangeAccomplishment(ctx, &acc)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddInterest(ctx context.Context, input AddInterestRequest) (*SuccessResolver, error) {
	interest := userRPC.Interest{
		Interest: input.Input.Interest,
	}

	if input.Input.Description == nil {
		interest.IsDescriptionNull = true
	} else {
		interest.Description = *input.Input.Description
	}

	id, err := user.AddInterest(ctx, &interest)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) ChangeInterest(ctx context.Context, input ChangeInterestRequest) (*SuccessResolver, error) {
	interest := userRPC.Interest{
		ID: input.ID,
	}

	if input.Interest.Interest == nil {
		interest.IsInterestNull = true
	} else {
		interest.Interest = *input.Interest.Interest
	}

	if input.Interest.Description == nil {
		interest.IsDescriptionNull = true
	} else {
		interest.Description = *input.Interest.Description
	}

	_, err := user.ChangeInterest(ctx, &interest)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveInterest(ctx context.Context, input RemoveInterestRequest) (*SuccessResolver, error) {
	interest := userRPC.Interest{
		ID: input.ID,
	}

	_, err := user.RemoveInterest(ctx, &interest)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveImageInInterest(ctx context.Context, input RemoveImageInInterestRequest) (*SuccessResolver, error) {
	interest := userRPC.Interest{
		ID: input.ID,
	}

	_, err := user.RemoveImageInInterest(ctx, &interest)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) GetOriginImageInInterest(ctx context.Context, input GetOriginImageInInterestRequest) (*string, error) {
	var a string
	avatar, err := user.GetOriginImageInInterest(ctx, &userRPC.Interest{
		ID: input.Interest_id,
	})
	if err != nil {
		return nil, nil
	}
	a = avatar.GetURL()

	return &a, nil
}

func (_ *Resolver) AskRecommendation(ctx context.Context, input AskRecommendationRequest) (*SuccessResolver, error) {
	_, err := network.AskRecommendation(ctx, &networkRPC.RecommendationParams{
		UserId:    input.User_id,
		Text:      input.Text,
		Title:     NullToString(input.Title),
		Relations: stringToRelationRPC(input.Relation),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) IgnoreRecommendationRequest(ctx context.Context, input IgnoreRecommendationRequestRequest) (*SuccessResolver, error) {
	_, err := network.IgnoreRecommendationRequest(ctx, &networkRPC.RecommendationRequest{
		Id: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SetVisibilityRecommendation(ctx context.Context, input SetVisibilityRecommendationRequest) (*SuccessResolver, error) {
	_, err := network.SetRecommendationVisibility(ctx, &networkRPC.RecommendationVisibility{
		RecommendationId: input.ID,
		Visible:          input.Is_visible,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) WriteRecommendation(ctx context.Context, input WriteRecommendationRequest) (*SuccessResolver, error) {
	_, err := network.WriteRecommendation(ctx, &networkRPC.RecommendationParams{
		UserId:    input.User_id,
		Text:      input.Text,
		Title:     NullToString(input.Title),
		Relations: stringToRelationRPC(input.Relation),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveProfileTranslation(ctx context.Context, input SaveProfileTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserProfileTranslation(ctx, &userRPC.ProfileTranslation{
		Firstname: input.Translations.Firstname,
		Headline:  input.Translations.Headline,
		Language:  input.LanguageID,
		Lastname:  input.Translations.Lastname,
		Nickname:  input.Translations.Nickname,
		Story:     input.Translations.Story,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveUserExperienceTranslation(ctx context.Context, input SaveUserExperienceTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserExperienceTranslation(ctx, &userRPC.ExperienceTranslation{
		Language:     input.LanguageID,
		ExperienceID: input.Translations.Experience_id,
		Company:      input.Translations.Company,
		Position:     input.Translations.Position,
		Description:  input.Translations.Description,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveUserEducationTranslation(ctx context.Context, input SaveUserEducationTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserEducationTranslation(ctx, &userRPC.EducationTranslation{
		Language:    input.LanguageID,
		EducationID: input.Translations.Education_id,
		Degree:      input.Translations.Degree,
		Description: input.Translations.Description,
		FieldStudy:  input.Translations.Field_of_study,
		Grade:       input.Translations.Grade,
		School:      input.Translations.School,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveUserInterestTranslation(ctx context.Context, input SaveUserInterestTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserInterestTranslation(ctx, &userRPC.InterestTranslation{
		Language:    input.LanguageID,
		Description: input.Translations.Description,
		InterestID:  input.Translations.Interest_id,
		Interest:    input.Translations.Interest,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveUserPortfolioTranslation(ctx context.Context, input SaveUserPortfolioTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserPortfolioTranslation(ctx, &userRPC.PortfolioTranslation{
		Language:    input.LanguageID,
		Description: input.Translations.Description,
		PortfolioID: input.Translations.Portfolio_id,
		Title:       input.Translations.Tittle,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveUserToolTechnologyTranslation(ctx context.Context, input SaveUserToolTechnologyTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserToolTechnologyTranslation(ctx, &userRPC.ToolTechnologyTranslation{
		Language:         input.LanguageID,
		ToolTechnology:   input.Translations.Tool_technology,
		TooltechnologyID: input.Translations.Tool_technology_id,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveUserSkillTranslation(ctx context.Context, input SaveUserSkillTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserSkillTranslation(ctx, &userRPC.SkillTranslation{
		Language: input.LanguageID,
		Skill:    input.Translations.Skill,
		SkillID:  input.Translations.Skill_id,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SaveUserAccomplishmentTranslation(ctx context.Context, input SaveUserAccomplishmentTranslationRequest) (*SuccessResolver, error) {
	_, err := user.SaveUserAccomplishmentTranslation(ctx, &userRPC.AccomplishmentTranslation{
		Language:         input.LanguageID,
		AccomplishmentID: input.Translations.Accomplishment_id,
		Description:      input.Translations.Description,
		Issuer:           input.Translations.Issuer,
		Name:             input.Translations.Name,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveTranslation(ctx context.Context, data RemoveTranslationRequest) (*SuccessResolver, error) {
	_, err := user.RemoveTranslation(ctx, &userRPC.Language{Language: data.LanguageID})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveNotification(ctx context.Context, input RemoveNotificationRequest) (*SuccessResolver, error) {
	_, err := notifications.RemoveNotification(ctx, &notificationsRPC.IDs{
		ID: input.Ids,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) MarkNotificationAsSeen(ctx context.Context, input MarkNotificationAsSeenRequest) (*SuccessResolver, error) {
	_, err := notifications.MarkAsSeen(ctx, &notificationsRPC.IDs{
		ID: input.Ids,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ReportUser(ctx context.Context, input ReportUserRequest) (*SuccessResolver, error) {
	_, err := user.ReportUser(ctx, &userRPC.ReportUserRequest{
		Description: NullToString(input.Input.Text),
		UserID:      input.User_id,
		Type:        stringToUserRPCReportType(input.Input.Reason),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func stringToRelationRPC(data *string) userRPC.RecommendationRelationEnum {
	if data == nil {
		return userRPC.RecommendationRelationEnum_NO_RELATION
	}

	switch *data {
	case "experience":
		return userRPC.RecommendationRelationEnum_EXPERIENCE
	case "education":
		return userRPC.RecommendationRelationEnum_EDUCATION
	case "accomplishment":
		return userRPC.RecommendationRelationEnum_ACCOMPLISHMENT
	}

	return userRPC.RecommendationRelationEnum_NO_RELATION
}

func relationRPCToString(data userRPC.RecommendationRelationEnum) string {

	switch data {
	case userRPC.RecommendationRelationEnum_EXPERIENCE:
		return "experience"
	case userRPC.RecommendationRelationEnum_EDUCATION:
		return "education"
	case userRPC.RecommendationRelationEnum_ACCOMPLISHMENT:
		return "accomplishment"
	}

	return "no_relation"
}

func stringToUserRPCReportType(s string) userRPC.ReportUserRequest_ReportType {
	switch s {
	case "report_violates_terms_of_use":
		return userRPC.ReportUserRequest_VolatationTermsOfUse
	case "not_real_individual":
		return userRPC.ReportUserRequest_NotRealIndividual
	case "pretending_to_be_someone":
		return userRPC.ReportUserRequest_PretendingToBeSomone
	case "may_be_hacked":
		return userRPC.ReportUserRequest_MayBeHacked
	case "avatar_is_not_person":
		return userRPC.ReportUserRequest_PictureIsNotPerson
	case "avatar_is_offensive":
		return userRPC.ReportUserRequest_PictureIsOffensive
	}

	return userRPC.ReportUserRequest_Other
}

func ToProfile(ctx context.Context, profile *userRPC.Profile) Profile {
	if profile == nil {
		return Profile{
			Location: &LocationProfile{},
		}
	}

	p := Profile{
		ID:                   profile.GetID(),
		Firstname:            profile.GetFirstname(),
		Lastname:             profile.GetLastname(),
		Url:                  profile.GetURL(),
		Date_of_registration: profile.GetDateOfActivation(),
		Emails:               profile.GetEmails(),
		Phones:               profile.GetPhones(),
		Network_info:         &NetworkInfoInUserProfile{},
	}

	p.Avatar = /*NullableStringClause(*/ profile.GetAvatar()          //, profile.GetIsAvatarNull())
	p.Middlename = /* NullableStringClause(*/ profile.GetMiddlename() //, profile.GetIsMiddlenameNull())
	p.Patronymic = /*NullableStringClause(*/ profile.GetPatronymic()  //, profile.GetIsPatronymicNull())
	p.Nickname = /*NullableStringClause(*/ profile.GetNickname()      //, profile.GetIsNicknameNull())

	if profile.GetIsNativeNameNull() && profile.GetNativeName() == nil {
		p.Native_name = nil
	} else {
		p.Native_name = &NativeName{
			Name:     profile.GetNativeName().GetName(),
			Language: profile.GetNativeName().GetLanguage(),
		}
	}

	p.Headline = profile.GetHeadline()
	p.Email = /*NullableStringClause(*/ profile.GetEmail()       //, profile.GetIsEmailNull())
	p.Phone = /*NullableStringClause(*/ profile.GetPhoneNumber() //, profile.GetPhoneNumberNull())
	p.Birthday = /*NullableStringClause(*/ profile.GetBirthday() //, profile.GetBirthdayNull())

	p.Location = &LocationProfile{} // it should ne nil
	if profile.GetIsLocationNull() || profile.GetLocation() == nil {
		p.Location = &LocationProfile{} // it should ne nil
	} else {

		p.Location = &LocationProfile{
			City:    profile.GetLocation().GetCity(),
			Country: profile.GetLocation().GetCountryID(),
		}
	}

	p.Gender = func(s string) string {
		if s == "FEMALE" {
			return "female"
		}
		return "male"
	}(profile.GetGender().String())
	p.Story = profile.GetStory()

	p.Online = profile.GetIsOnline()
	p.Me = profile.GetIsMe()
	p.Friend = profile.GetIsFriend()
	p.Follow = profile.GetIsFollow()
	p.Favorite = profile.GetIsFavorite()
	p.Blocked = profile.GetIsBlocked()
	p.Friend_request = profile.GetIsFriendRequestSend()
	p.Recieved_friend_request = profile.GetIsFriendRequestRecieved()
	p.Friendship_id = profile.GetFriendshipID()
	p.Profile_complete_percent = profile.GetProfileCompletePercent()
	p.Current_translation = profile.GetCurrentTranslation()
	p.Available_translations = profile.GetAvailableTranslations()

	countings, _ := network.GetUserCountings(ctx, &networkRPC.User{
		Id: profile.GetID(),
	})

	reviewAmount, _ := company.GetAmountOfReviewsOfUser(ctx, &companyRPC.ID{
		ID: profile.GetID(),
	})

	p.Network_info = &NetworkInfoInUserProfile{
		Connections:               countings.GetNumOfConnections(),
		Followers:                 countings.GetNumOfFollowers(),
		Followings:                countings.GetNumOfFollowings(),
		Recommendations:           countings.GetNumOfReceivedRecommendations(),
		Reviews:                   reviewAmount.GetAmount(),
		Mutual_connections_amount: profile.GetMutualConnectionsAmount(),
	}

	return p
}

func portoflioInfoRPCToPortfolioInfo(data *userRPC.PortfolioInfo) *PortfolioInfo {
	if data == nil {
		return nil
	}

	return &PortfolioInfo{
		Portfolio_statistic: PortfolioInfoStatistic{
			Like_count:    data.GetLikeCount(),
			View_count:    data.GetViewCount(),
			Comment_count: data.GetCommentCount(),
		},

		Portfolio_amount: PortfolioInfoAmount{
			Has_photo:   data.GetHasPhoto(),
			Has_video:   data.GetHasVideo(),
			Has_article: data.GetHasArticle(),
			Has_audio:   data.GetHasAudio(),
		},
	}
}
