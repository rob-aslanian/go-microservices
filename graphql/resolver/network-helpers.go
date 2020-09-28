package resolver

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func categoryItemArrayToGql(input []*networkRPC.CategoryItem) []CategoryItem {
	items := make([]CategoryItem, len(input))
	for i, cat := range input {
		items[i] = categoryItemToGql(cat)
	}
	return items
}

func categoryItemToGql(input *networkRPC.CategoryItem) CategoryItem {
	item := CategoryItem{}
	if input != nil {
		item.Name = input.Name
		item.Unique_name = input.UniqueName
		item.Has_children = input.HasChildren
		item.Children = categoryItemArrayToGql(input.Children)
	}
	return item
}

func friendshipRequestToRPC(input GetFriendRequestsRequest) *networkRPC.FriendRequestFilter {
	if input.Status == nil {
		defStatus := "Requested"
		input.Status = &defStatus
	}
	if input.Sent == nil {
		defSent := false
		input.Sent = &defSent
	}
	return &networkRPC.FriendRequestFilter{
		Status: *input.Status,
		Sent:   *input.Sent,
	}
}

func friendshipFilterToRPC(input GetFriendshipsRequest) *networkRPC.FriendshipFilter {
	emptyString := ""
	if input.Category == nil {
		input.Category = &emptyString
	}
	if input.Letter == nil {
		input.Letter = &emptyString
	}
	if input.Query == nil {
		input.Query = &emptyString
	}
	if input.Sort_by == nil {
		input.Sort_by = &emptyString
	}
	if input.Companies == nil {
		input.Companies = &[]string{}
	}
	return &networkRPC.FriendshipFilter{
		Query:     *input.Query,
		Category:  *input.Category,
		Letter:    *input.Letter,
		SortBy:    *input.Sort_by,
		Companies: *input.Companies,
	}
}

func userToNetworkUser(u *networkRPC.User) NetworkUser {
	exp, err := user.GetExperiences(context.Background(), &userRPC.RequestExperiences{
		First:  1,
		UserID: u.Id,
	})
	var lastExp LastExperience
	if err == nil && len(exp.Experiences) > 0 {
		last := exp.Experiences[0]
		lastExp.Company = last.Company
		lastExp.Position = last.Position
	}
	return NetworkUser{
		ID:              u.Id,
		Primary_email:   u.PrimaryEmail,
		Primary_phone:   u.PrimaryPhone,
		First_name:      u.FirstName,
		Last_name:       u.LastName,
		Url:             u.Url,
		Avatar:          u.Avatar,
		Gender:          u.Gender,
		Last_experience: lastExp,
	}
}

func friendshipToResolver(fr *networkRPC.Friendship) FriendshipResolver {
	return FriendshipResolver{
		R: &Friendship{
			ID:           fr.Id,
			Status:       fr.Status,
			Description:  fr.Description,
			Following:    fr.Following,
			Categories:   fr.Categories,
			My_request:   fr.MyRequest,
			Friend:       userToNetworkUser(fr.Friend),
			Created_at:   float64(fr.CreatedAt),
			Responded_at: float64(fr.RespondedAt),
		},
	}
}

func friendshipWithProfileToResolver(fr *networkRPC.Friendship) FriendshipWithProfileResolver {
	return FriendshipWithProfileResolver{
		R: &FriendshipWithProfile{
			ID:           fr.Id,
			Status:       fr.Status,
			Description:  fr.Description,
			Following:    fr.Following,
			Categories:   fr.Categories,
			My_request:   fr.MyRequest,
			Friend:       userToNetworkUser(fr.Friend),
			Created_at:   float64(fr.CreatedAt),
			Responded_at: float64(fr.RespondedAt),
		},
	}
}

func followInfoToResolver(fol *networkRPC.FollowInfo) FollowInfoResolver {
	return FollowInfoResolver{
		R: &FollowInfo{
			User:      userToNetworkUser(fol.User),
			Followers: fol.Followers,
			Following: fol.Following,
			Is_friend: fol.IsFriend,
		},
	}
}

func followInfoWithProfileToResolver(fol *networkRPC.FollowInfo) FollowInfoWithProfileResolver {
	return FollowInfoWithProfileResolver{
		R: &FollowInfoWithProfile{
			User: userToNetworkUser(fol.User),
			// User_profile: ToProfile(nil),
			Followers: fol.Followers,
			Following: fol.Following,
			Is_friend: fol.IsFriend,
		},
	}
}

func suggestionToResolver(sug *networkRPC.UserSuggestion) UserSuggestionResolver {
	return UserSuggestionResolver{
		R: &UserSuggestion{
			User:      userToNetworkUser(sug.User),
			Following: sug.Following,
			Followers: sug.Followers,
		},
	}
}

func companySuggestionToResolver(ctx context.Context, sug *networkRPC.CompanySuggestion) CompanySuggestionResolver {
	comp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: []string{sug.GetCompany().GetId()},
	})
	if err != nil {
		log.Println(err)
	}

	sg := CompanySuggestionResolver{
		R: &CompanySuggestion{
			Company:   companyToNetworkCompany(sug.Company),
			Followers: sug.Followers,
		},
	}

	if len(comp.GetProfiles()) > 0 {
		sg.R.Company_profile = toCompanyProfile(ctx, *(comp.GetProfiles()[0]))
	}

	return sg
}

func companyToNetworkCompany(company *networkRPC.Company) NetworkCompany {
	return NetworkCompany{
		ID:              company.Id,
		Name:            company.Name,
		Url:             company.Url,
		Avatar:          company.Avatar,
		Industry:        company.Industry,
		Type:            companyTypeToEnum(company.Type),
		Email:           company.Email,
		Address:         company.Address,
		Foundation_year: company.FoundationYear,
	}
}

func companyTypeToEnum(s string) string {
	t := "type_unknown"

	switch s {
	case "TYPE_SELF_EMPLOYED":
		t = "type_self_emplyed"
	case "TYPE_EDUCATIONAL_INSTITUTION":
		t = "type_educational_institution"
	case "TYPE_GOVERNMENT_AGENSY":
		t = "type_government_agency"
	case "TYPE_SOLE_PROPRIETORSHIP":
		t = "type_sole_proprietorship"
	case "TYPE_PRIVATELY_HELD":
		t = "type_privately_held"
	case "TYPE_PARTNERSHIP":
		t = "type_partnership"
	case "TYPE_PUBLIC_COMPANY":
		t = "type_public_company"
	}

	return t
}

func companyFollowToResolver(fol *networkRPC.CompanyFollowInfo) CompanyFollowInfoResolver {
	return CompanyFollowInfoResolver{
		R: &CompanyFollowInfo{
			Company:    companyToNetworkCompany(fol.Company),
			Following:  fol.Following,
			Followers:  fol.Followers,
			Rating:     fol.Rating,
			Size:       fol.Size,
			Categories: fol.Categories,
		},
	}
}

func companyFilterToRPC(inp interface{}) *networkRPC.CompanyFilter {
	switch input := inp.(type) {
	case GetFollowingCompaniesRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		return &networkRPC.CompanyFilter{
			Query:    *input.Query,
			Category: *input.Category,
			Letter:   *input.Letter,
			SortBy:   *input.Sort_by,
		}
	case GetFollowerCompaniesRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		return &networkRPC.CompanyFilter{
			Query:    *input.Query,
			Category: *input.Category,
			Letter:   *input.Letter,
			SortBy:   *input.Sort_by,
		}
	case GetFollowingCompaniesForCompanyRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		return &networkRPC.CompanyFilter{
			Query:    *input.Query,
			Category: *input.Category,
			Letter:   *input.Letter,
			SortBy:   *input.Sort_by,
		}
	case GetFollowerCompaniesForCompanyRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		return &networkRPC.CompanyFilter{
			Query:    *input.Query,
			Category: *input.Category,
			Letter:   *input.Letter,
			SortBy:   *input.Sort_by,
		}
	default:
		return &networkRPC.CompanyFilter{}
	}
}

func userFilterToRPC(inp interface{}) *networkRPC.UserFilter {
	switch input := inp.(type) {
	case GetFollowingsRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		if input.Companies == nil {
			input.Companies = &[]string{}
		}
		return &networkRPC.UserFilter{
			Query:     *input.Query,
			Category:  *input.Category,
			Letter:    *input.Letter,
			SortBy:    *input.Sort_by,
			Companies: *input.Companies,
		}
	case GetFollowersRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		if input.Companies == nil {
			input.Companies = &[]string{}
		}
		return &networkRPC.UserFilter{
			Query:     *input.Query,
			Category:  *input.Category,
			Letter:    *input.Letter,
			SortBy:    *input.Sort_by,
			Companies: *input.Companies,
		}
	case GetFollowingsForCompanyRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		if input.Companies == nil {
			input.Companies = &[]string{}
		}
		return &networkRPC.UserFilter{
			Query:     *input.Query,
			Category:  *input.Category,
			Letter:    *input.Letter,
			SortBy:    *input.Sort_by,
			Companies: *input.Companies,
		}
	case GetFollowersForCompanyRequest:
		emptyString := ""
		if input.Category == nil {
			input.Category = &emptyString
		}
		if input.Letter == nil {
			input.Letter = &emptyString
		}
		if input.Query == nil {
			input.Query = &emptyString
		}
		if input.Sort_by == nil {
			input.Sort_by = &emptyString
		}
		if input.Companies == nil {
			input.Companies = &[]string{}
		}
		return &networkRPC.UserFilter{
			Query:     *input.Query,
			Category:  *input.Category,
			Letter:    *input.Letter,
			SortBy:    *input.Sort_by,
			Companies: *input.Companies,
		}
	default:
		return &networkRPC.UserFilter{}
	}

}
