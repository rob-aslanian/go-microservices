package resolver

import (
	"context"
	"errors"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func (_ *Resolver) GetCategoryTree(ctx context.Context) ([]CategoryItemResolver, error) {
	res, err := network.GetCategoryTree(ctx, &networkRPC.Empty{})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	items := make([]CategoryItemResolver, len(res.Categories))
	for i, cat := range res.Categories {
		it := categoryItemToGql(cat)
		items[i] = CategoryItemResolver{
			R: &it,
		}
	}

	return items, nil
}

func (_ *Resolver) GetFollowingsCategoryTree(ctx context.Context) ([]CategoryItemResolver, error) {
	res, err := network.GetCategoryTreeForFollowings(ctx, &networkRPC.Empty{})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	items := make([]CategoryItemResolver, len(res.Categories))
	for i, cat := range res.Categories {
		it := categoryItemToGql(cat)
		items[i] = CategoryItemResolver{
			R: &it,
		}
	}

	return items, nil
}

func (_ *Resolver) GetFollowingsCategoryTreeForCompany(ctx context.Context, input GetFollowingsCategoryTreeForCompanyRequest) ([]CategoryItemResolver, error) {
	res, err := network.GetCategoryTreeForFollowingsForCompany(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	items := make([]CategoryItemResolver, len(res.Categories))
	for i, cat := range res.Categories {
		it := categoryItemToGql(cat)
		items[i] = CategoryItemResolver{
			R: &it,
		}
	}

	return items, nil
}

func (_ *Resolver) CreateCategory(ctx context.Context, input CreateCategoryRequest) (*bool, error) {
	_, err := network.CreateCategory(ctx, &networkRPC.CategoryPath{
		Name:   input.Name,
		Parent: input.Parent,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) CreateFollowingsCategory(ctx context.Context, input CreateFollowingsCategoryRequest) (*bool, error) {
	_, err := network.CreateCategoryForFollowings(ctx, &networkRPC.CategoryPath{
		Name:   input.Name,
		Parent: input.Parent,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) CreateFollowingsCategoryForCompany(ctx context.Context, input CreateFollowingsCategoryForCompanyRequest) (*bool, error) {
	_, err := network.CreateCategoryForFollowingsForCompany(ctx, &networkRPC.CategoryPathWithCompanyId{
		CompanyId: input.CompanyId,
		CategoryPath: &networkRPC.CategoryPath{
			Name:   input.Name,
			Parent: input.Parent,
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RemoveCategory(ctx context.Context, input RemoveCategoryRequest) (*bool, error) {
	_, err := network.RemoveCategory(ctx, &networkRPC.CategoryPath{
		Name:   input.Name,
		Parent: input.Parent,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RemoveFollowingsCategory(ctx context.Context, input RemoveFollowingsCategoryRequest) (*bool, error) {
	_, err := network.RemoveCategoryForFollowings(ctx, &networkRPC.CategoryPath{
		Name:   input.Name,
		Parent: input.Parent,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RemoveFollowingsCategoryForCompany(ctx context.Context, input RemoveFollowingsCategoryForCompanyRequest) (*bool, error) {
	_, err := network.RemoveCategoryForFollowingsForCompany(ctx, &networkRPC.CategoryPathWithCompanyId{
		CompanyId: input.CompanyId,
		CategoryPath: &networkRPC.CategoryPath{
			Name:   input.Name,
			Parent: input.Parent,
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetFriendRequests(ctx context.Context, input GetFriendRequestsRequest) ([]FriendshipResolver, error) {
	res, err := network.GetFriendRequests(ctx, friendshipRequestToRPC(input))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	requests := make([]FriendshipResolver, 0, len(res.Friendships))
	for _, fr := range res.Friendships {
		requests = append(requests, friendshipToResolver(fr))
	}
	return requests, nil
}

func (_ *Resolver) GetFriendships(ctx context.Context, input GetFriendshipsRequest) ([]FriendshipWithProfileResolver, error) {
	res, err := network.GetAllFriendships(ctx, friendshipFilterToRPC(input))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	requests := make([]FriendshipWithProfileResolver, 0, len(res.Friendships))
	for _, fr := range res.Friendships {
		requests = append(requests, friendshipWithProfileToResolver(fr))
	}

	ids := make([]string, 0, len(requests))
	for _, r := range requests {
		ids = append(ids, string(r.Friend().ID()))
	}

	if len(ids) > 0 {
		profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
			ID: ids,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range requests {
			r.R.Friend_profile = ToProfile(ctx, profiles.GetProfiles()[string(r.Friend().ID())])
			r.R.Friend_profile.ID = string(r.Friend().ID())
		}
	}
	return requests, nil
}

func (_ *Resolver) SendFriendRequest(ctx context.Context, input SendFriendRequestRequest) (*FriendshipResolver, error) {
	if input.Description == nil {
		descValue := ""
		input.Description = &descValue
	}
	res, err := network.SendFriendRequest(ctx, &networkRPC.FriendshipRequest{FriendId: input.UserId, Description: *input.Description})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	resolver := friendshipToResolver(res)
	return &resolver, nil
}

func (_ *Resolver) ApproveFriendRequest(ctx context.Context, input ApproveFriendRequestRequest) (*bool, error) {
	_, err := network.ApproveFriendRequest(ctx, &networkRPC.Friendship{Id: input.RequestId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) DenyFriendRequest(ctx context.Context, input DenyFriendRequestRequest) (*bool, error) {
	_, err := network.DenyFriendRequest(ctx, &networkRPC.Friendship{Id: input.RequestId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) IgnoreFriendRequest(ctx context.Context, input IgnoreFriendRequestRequest) (*bool, error) {
	_, err := network.IgnoreFriendRequest(ctx, &networkRPC.Friendship{Id: input.RequestId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) Unfriend(ctx context.Context, input UnfriendRequest) (*bool, error) {
	_, err := network.Unfriend(ctx, &networkRPC.User{Id: input.UserId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetFollowers(ctx context.Context, input GetFollowersRequest) ([]FollowInfoWithProfileResolver, error) {
	res, err := network.GetFollowers(ctx, userFilterToRPC(input))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	followInfos := make([]FollowInfoWithProfileResolver, 0, len(res.Follows))
	for _, fol := range res.Follows {
		followInfos = append(followInfos, followInfoWithProfileToResolver(fol))
	}

	ids := make([]string, 0, len(followInfos))
	for _, r := range followInfos {
		ids = append(ids, string(r.User().ID()))
	}

	if len(ids) > 0 {
		profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
			ID: ids,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range followInfos {
			r.R.User_profile = ToProfile(ctx, profiles.GetProfiles()[string(r.User().ID())])
			r.R.User_profile.ID = string(r.User().ID())
		}
	}

	return followInfos, nil
}

func (_ *Resolver) GetFollowings(ctx context.Context, input GetFollowingsRequest) ([]FollowInfoWithProfileResolver, error) {
	res, err := network.GetFollowings(ctx, userFilterToRPC(input))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	followInfos := make([]FollowInfoWithProfileResolver, 0, len(res.Follows))
	for _, fol := range res.Follows {
		followInfos = append(followInfos, followInfoWithProfileToResolver(fol))
	}

	ids := make([]string, 0, len(followInfos))
	for _, r := range followInfos {
		ids = append(ids, string(r.User().ID()))
	}

	if len(ids) > 0 {
		profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
			ID: ids,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range followInfos {
			r.R.User_profile = ToProfile(ctx, profiles.GetProfiles()[string(r.User().ID())])
			r.R.User_profile.ID = string(r.User().ID())
		}
	}

	return followInfos, nil
}

func (_ *Resolver) Follow(ctx context.Context, input FollowRequest) (*bool, error) {
	_, err := network.Follow(ctx, &networkRPC.User{Id: input.UserId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) Unfollow(ctx context.Context, input UnfollowRequest) (*bool, error) {
	_, err := network.Unfollow(ctx, &networkRPC.User{Id: input.UserId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) AddToFavourites(ctx context.Context, input AddToFavouritesRequest) (CategoryNameWithUserIdResolver, error) {
	_, err := network.AddToFavourites(ctx, &networkRPC.User{Id: input.UserId})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithUserIdResolver{}, e
	}
	return CategoryNameWithUserIdResolver{R: &CategoryNameWithUserId{Unique_name: "favourites", User_id: input.UserId}}, nil // "favourites" taken from network constant
}

func (_ *Resolver) AddToFollowingsFavourites(ctx context.Context, input AddToFollowingsFavouritesRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.AddToFollowingsFavourites(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: "favourites", Company_id: input.CompanyId}}, nil // "favourites" taken from network constant
}

func (_ *Resolver) AddToFollowingsFavouritesForCompany(ctx context.Context, input AddToFollowingsFavouritesForCompanyRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.AddToFollowingsFavouritesForCompany(ctx, &networkRPC.CompanyCompanyId{Company1Id: input.CompanyId, Company2Id: input.RefCompanyId})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: "favourites", Company_id: input.CompanyId}}, nil // "favourites" taken from network constant
}

func (_ *Resolver) RemoveFromFavourites(ctx context.Context, input RemoveFromFavouritesRequest) (CategoryNameWithUserIdResolver, error) {
	_, err := network.RemoveFromFavourites(ctx, &networkRPC.User{Id: input.UserId})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithUserIdResolver{}, e
	}
	return CategoryNameWithUserIdResolver{R: &CategoryNameWithUserId{Unique_name: "favourites", User_id: input.UserId}}, nil
}

func (_ *Resolver) RemoveFromFollowingsFavourites(ctx context.Context, input RemoveFromFollowingsFavouritesRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.RemoveFromFollowingsFavourites(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: "favourites", Company_id: input.CompanyId}}, nil
}

func (_ *Resolver) RemoveFromFollowingsFavouritesForCompany(ctx context.Context, input RemoveFromFollowingsFavouritesForCompanyRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.RemoveFromFollowingsFavouritesForCompany(ctx, &networkRPC.CompanyCompanyId{Company1Id: input.CompanyId, Company2Id: input.RefCompanyId})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: "favourites", Company_id: input.CompanyId}}, nil
}

func (_ *Resolver) AddToCategory(ctx context.Context, input AddToCategoryRequest) (CategoryNameWithUserIdResolver, error) {
	_, err := network.AddToCategory(ctx, &networkRPC.CategoryRequest{Id: input.UserId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithUserIdResolver{}, e
	}
	return CategoryNameWithUserIdResolver{R: &CategoryNameWithUserId{Unique_name: input.Category_name, User_id: input.UserId}}, nil
}

func (_ *Resolver) AddToFollowingsCategory(ctx context.Context, input AddToFollowingsCategoryRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.AddToFollowingsCategory(ctx, &networkRPC.CategoryRequest{Id: input.CompanyId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: input.Category_name, Company_id: input.CompanyId}}, nil
}

func (_ *Resolver) AddToFollowingsCategoryForCompany(ctx context.Context, input AddToFollowingsCategoryForCompanyRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.AddToFollowingsCategoryForCompany(ctx, &networkRPC.CategoryRequestWithCompanyId{CompanyId: input.CompanyId, Id: input.RefCompanyId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: input.Category_name, Company_id: input.RefCompanyId}}, nil // Tornike (frontend) asked  to return id of ref company
}

func (_ *Resolver) RemoveFromCategory(ctx context.Context, input RemoveFromCategoryRequest) (CategoryNameWithUserIdResolver, error) {
	_, err := network.RemoveFromCategory(ctx, &networkRPC.CategoryRequest{Id: input.UserId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithUserIdResolver{}, e
	}
	return CategoryNameWithUserIdResolver{R: &CategoryNameWithUserId{Unique_name: input.Category_name, User_id: input.UserId}}, nil
}

func (_ *Resolver) RemoveFromFollowingsCategory(ctx context.Context, input RemoveFromFollowingsCategoryRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.RemoveFromFollowingsCategory(ctx, &networkRPC.CategoryRequest{Id: input.CompanyId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: input.Category_name, Company_id: input.CompanyId}}, nil
}

func (_ *Resolver) RemoveFromFollowingsCategoryForCompany(ctx context.Context, input RemoveFromFollowingsCategoryForCompanyRequest) (CategoryNameWithCompanyIdResolver, error) {
	_, err := network.RemoveFromFollowingsCategoryForCompany(ctx, &networkRPC.CategoryRequestWithCompanyId{CompanyId: input.CompanyId, Id: input.RefCompanyId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return CategoryNameWithCompanyIdResolver{}, e
	}
	return CategoryNameWithCompanyIdResolver{R: &CategoryNameWithCompanyId{Unique_name: input.Category_name, Company_id: input.RefCompanyId}}, nil // Tornike (frontend) asked  to return id of ref company
}

func (_ *Resolver) BatchRemoveFromCategory(ctx context.Context, input BatchRemoveFromCategoryRequest) (*bool, error) {
	_, err := network.BatchRemoveFromCategory(ctx, &networkRPC.BatchRemoveFromCategoryRequest{
		Ids:      input.UserIds,
		Category: input.Category_name,
		All:      input.All,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) BatchRemoveFromFollowingsCategory(ctx context.Context, input BatchRemoveFromFollowingsCategoryRequest) (*bool, error) {
	_, err := network.BatchRemoveFromFollowingsCategory(ctx, &networkRPC.BatchRemoveFromCategoryRequest{
		Ids:      input.CompanyIds,
		Category: input.Category_name,
		All:      input.All,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) BatchRemoveFromFollowingsCategoryForCompany(ctx context.Context, input BatchRemoveFromFollowingsCategoryForCompanyRequest) (*bool, error) {
	_, err := network.BatchRemoveFromFollowingsCategoryForCompany(ctx, &networkRPC.BatchRemoveFromCategoryRequestWithCompanyId{
		CompanyId: input.CompanyId,
		Request: &networkRPC.BatchRemoveFromCategoryRequest{
			Ids:      input.CompanyIds,
			Category: input.Category_name,
			All:      input.All,
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetFriendSuggestions(ctx context.Context, input GetFriendSuggestionsRequest) ([]UserSuggestionResolver, error) {
	res, err := network.GetFriendSuggestions(ctx, &networkRPC.Pagination{
		After:  NullIDToInt32(input.Pagination.After),
		Amount: NullToInt32(input.Pagination.First),
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	suggestions := make([]UserSuggestionResolver, 0, len(res.Suggestions))
	for _, suggestion := range res.Suggestions {
		res := suggestionToResolver(suggestion)
		suggestions = append(suggestions, res)
	}

	ids := make([]string, 0, len(suggestions))
	for _, r := range suggestions {
		ids = append(ids, string(r.User().ID()))
	}

	profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: ids,
	})
	if err != nil {
		return nil, err
	}

	for _, r := range suggestions {
		r.R.User_profile = ToProfile(ctx, profiles.GetProfiles()[string(r.User().ID())])
	}

	return suggestions, nil
}

func (_ *Resolver) FollowCompany(ctx context.Context, input FollowCompanyRequest) (*bool, error) {
	_, err := network.FollowCompany(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnfollowCompany(ctx context.Context, input UnfollowCompanyRequest) (*bool, error) {
	_, err := network.UnfollowCompany(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetFollowerCompanies(ctx context.Context, input GetFollowerCompaniesRequest) ([]CompanyFollowInfoResolver, error) {
	res, err := network.GetFollowerCompanies(ctx, companyFilterToRPC(input))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	follows := make([]CompanyFollowInfoResolver, 0, len(res.Follows))
	for _, follow := range res.Follows {
		res := companyFollowToResolver(follow)
		follows = append(follows, res)
	}

	ids := make([]string, 0, len(res.GetFollows()))
	for _, v := range res.GetFollows() {
		ids = append(ids, v.GetCompany().GetId())
	}

	companies := make(map[string]CompanyProfile, len(res.GetFollows()))

	comps, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	for _, pr := range comps.GetProfiles() {
		if pr != nil {
			companies[pr.GetId()] = toCompanyProfile(ctx, *pr)
		}
	}

	for i := range follows {
		if follows[i].R != nil {
			follows[i].R.Company_profile = companies[follows[i].R.Company.ID]
		}
	}

	return follows, nil
}

func (_ *Resolver) GetFollowingCompanies(ctx context.Context, input GetFollowingCompaniesRequest) ([]CompanyFollowInfoResolver, error) {
	res, err := network.GetFilteredFollowingCompanies(ctx, companyFilterToRPC(input))
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	follows := make([]CompanyFollowInfoResolver, 0, len(res.Follows))
	for _, follow := range res.Follows {
		res := companyFollowToResolver(follow)
		follows = append(follows, res)
	}

	ids := make([]string, 0, len(res.GetFollows()))
	for _, v := range res.GetFollows() {
		ids = append(ids, v.GetCompany().GetId())
	}

	companies := make(map[string]CompanyProfile, len(res.GetFollows()))

	comps, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	for _, pr := range comps.GetProfiles() {
		if pr != nil {
			companies[pr.GetId()] = toCompanyProfile(ctx, *pr)
		}
	}

	for i := range follows {
		if follows[i].R != nil {
			follows[i].R.Company_profile = companies[follows[i].R.Company.ID]
		}
	}

	return follows, nil
}

func (_ *Resolver) AddCompanyToFavourites(ctx context.Context, input AddCompanyToFavouritesRequest) (*bool, error) {
	_, err := network.AddCompanyToFavourites(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RemoveCompanyFromFavourites(ctx context.Context, input RemoveCompanyFromFavouritesRequest) (*bool, error) {
	_, err := network.RemoveCompanyFromFavourites(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) AddCompanyToCategory(ctx context.Context, input AddCompanyToCategoryRequest) (*bool, error) {
	_, err := network.AddCompanyToCategory(ctx, &networkRPC.CategoryRequest{Id: input.CompanyId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RemoveCompanyFromCategory(ctx context.Context, input RemoveCompanyFromCategoryRequest) (*bool, error) {
	_, err := network.RemoveCompanyFromCategory(ctx, &networkRPC.CategoryRequest{Id: input.CompanyId, CategoryName: input.Category_name})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetSuggestedCompanies(ctx context.Context, input GetSuggestedCompaniesRequest) ([]CompanySuggestionResolver, error) {
	res, err := network.GetSuggestedCompanies(ctx, &networkRPC.Pagination{
		After:  NullIDToInt32(input.Pagination.After),
		Amount: NullToInt32(input.Pagination.First),
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	suggestions := make([]CompanySuggestionResolver, 0, len(res.Suggestions))
	for _, suggestion := range res.Suggestions {
		res := companySuggestionToResolver(ctx, suggestion)
		suggestions = append(suggestions, res)
	}

	return suggestions, nil
}

func (_ *Resolver) BlockUser(ctx context.Context, input BlockUserRequest) (*bool, error) {
	_, err := network.BlockUser(ctx, &networkRPC.User{Id: input.UserId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnblockUser(ctx context.Context, input UnblockUserRequest) (*bool, error) {
	_, err := network.UnblockUser(ctx, &networkRPC.User{Id: input.UserId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) BlockCompany(ctx context.Context, input BlockCompanyRequest) (*bool, error) {
	_, err := network.BlockCompany(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnblockCompany(ctx context.Context, input UnblockCompanyRequest) (*bool, error) {
	_, err := network.UnblockCompany(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetBlockedUsersOrCompanies(ctx context.Context) ([]BlockedUserOrCompanyResolver, error) {
	res, err := network.GetblockedUsersOrCompanies(ctx, &networkRPC.Empty{})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	blockedUserResolvers := make([]BlockedUserOrCompanyResolver, 0, len(res.List))
	for _, user := range res.List {
		res := BlockedUserOrCompany{ID: user.Id, Name: user.Name, Avatar: user.Avatar, Is_company: user.IsCompany}
		blockedUserResolvers = append(blockedUserResolvers, BlockedUserOrCompanyResolver{R: &res})
	}
	return blockedUserResolvers, nil
}

func (_ *Resolver) BlockUserForCompany(ctx context.Context, input BlockUserForCompanyRequest) (*bool, error) {
	_, err := network.BlockUserForCompany(ctx, &networkRPC.UserCompanyId{
		UserId:    input.User_id,
		CompanyId: input.Company_id,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnblockUserForCompany(ctx context.Context, input UnblockUserForCompanyRequest) (*bool, error) {
	_, err := network.UnblockUserForCompany(ctx, &networkRPC.UserCompanyId{
		UserId:    input.User_id,
		CompanyId: input.Company_id,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetBlockedUsersForCompany(ctx context.Context, input GetBlockedUsersForCompanyRequest) (*[]BlockedUserOrCompanyResolver, error) {
	results, err := network.GetBlockedUsersForCompany(ctx, &networkRPC.Company{Id: input.Company_id})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	blocked := make([]BlockedUserOrCompanyResolver, 0, len(results.GetUsers()))

	for _, user := range results.GetUsers() {
		res := BlockedUserOrCompany{
			ID:         user.GetId(),
			Name:       user.GetFirstName() + user.GetLastName(),
			Avatar:     user.GetAvatar(),
			Is_company: false,
		}

		blocked = append(blocked, BlockedUserOrCompanyResolver{R: &res})
	}

	// ids := make([]string, 0, len(results.GetUsers()))
	// for _, p := range results.GetUsers() {
	// 	ids = append(ids, p.GetId())
	// }

	// profs := make([]ProfileResolver, 0, len(ids))
	//
	// if len(ids) > 0 {
	// 	profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
	// 		ID: ids,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	for _, r := range results.GetUsers() {
	// 		pr := Profile{}
	//
	// 		pr = ToProfile(profiles.GetProfiles()[string(r.GetId())])
	// 		pr.ID = string(r.GetId())
	//
	// 		profs = append(profs, ProfileResolver{
	// 			R: &pr,
	// 		})
	//
	// 	}
	// }

	// return &profs, nil

	return &blocked, nil
}

//func (_ *Resolver) GetBlockedUsers(ctx context.Context) ([]BlockedUserResolver, error) {
//	res, err := network.GetBlockedUsers(ctx, &networkRPC.Empty{})
//	if e, isErr := handleError(err); isErr {
//		return nil, e
//	}
//
//	blockedUserResolvers := make([]BlockedUserResolver, 0, len(res.Users))
//	for _, user := range res.Users {
//		res := BlockedUser{ID: user.Id, First_name: user.FirstName, Last_name: user.LastName}
//		blockedUserResolvers = append(blockedUserResolvers, BlockedUserResolver{R: &res})
//	}
//	return blockedUserResolvers, nil
//}
//
//func (_ *Resolver) GetBlockedCompanies(ctx context.Context) ([]BlockedCompanyResolver, error) {
//	res, err := network.GetBlockedCompanies(ctx, &networkRPC.Empty{})
//	if e, isErr := handleError(err); isErr {
//		return nil, e
//	}
//
//	blockedCompaniesResolvers := make([]BlockedCompanyResolver, 0, len(res.Companies))
//	for _, comp := range res.Companies {
//		res := BlockedCompany{ID: comp.Id, Name: comp.Name}
//		blockedCompaniesResolvers = append(blockedCompaniesResolvers, BlockedCompanyResolver{R: &res})
//	}
//	return blockedCompaniesResolvers, nil
//}

func (_ *Resolver) FollowForCompany(ctx context.Context, input FollowForCompanyRequest) (*bool, error) {
	_, err := network.FollowForCompany(ctx, &networkRPC.UserCompanyId{CompanyId: input.CompanyId, UserId: input.UserId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnfollowForCompany(ctx context.Context, input UnfollowForCompanyRequest) (*bool, error) {
	_, err := network.UnfollowForCompany(ctx, &networkRPC.UserCompanyId{CompanyId: input.CompanyId, UserId: input.UserId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) FollowCompanyForCompany(ctx context.Context, input FollowCompanyForCompanyRequest) (*bool, error) {
	_, err := network.FollowCompanyForCompany(ctx, &networkRPC.CompanyCompanyId{Company1Id: input.CompanyId, Company2Id: input.FollowId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) UnfollowCompanyForCompany(ctx context.Context, input UnfollowCompanyForCompanyRequest) (*bool, error) {
	_, err := network.UnfollowCompanyForCompany(ctx, &networkRPC.CompanyCompanyId{Company1Id: input.CompanyId, Company2Id: input.FollowId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetFollowingsForCompany(ctx context.Context, input GetFollowingsForCompanyRequest) ([]FollowInfoWithProfileResolver, error) {
	res, err := network.GetFollowingsForCompany(ctx, &networkRPC.IdWithUserFilter{Id: input.CompanyId, Filter: userFilterToRPC(input)})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	followInfos := make([]FollowInfoWithProfileResolver, 0, len(res.Follows))
	for _, fol := range res.Follows {
		followInfos = append(followInfos, followInfoWithProfileToResolver(fol))
	}
	ids := make([]string, 0, len(followInfos))
	for _, r := range followInfos {
		ids = append(ids, string(r.User().ID()))
	}

	profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: ids,
	})
	if err != nil {
		return nil, err
	}

	for _, r := range followInfos {
		r.R.User_profile = ToProfile(ctx, profiles.GetProfiles()[string(r.User().ID())])
	}

	return followInfos, nil
}

func (_ *Resolver) GetFollowersForCompany(ctx context.Context, input GetFollowersForCompanyRequest) ([]FollowInfoWithProfileResolver, error) {
	res, err := network.GetFollowersForCompany(ctx, &networkRPC.IdWithUserFilter{Id: input.CompanyId, Filter: userFilterToRPC(input)})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	followInfos := make([]FollowInfoWithProfileResolver, 0, len(res.Follows))
	for _, fol := range res.Follows {
		followInfos = append(followInfos, followInfoWithProfileToResolver(fol))
	}

	ids := make([]string, 0, len(followInfos))
	for _, r := range followInfos {
		ids = append(ids, string(r.User().ID()))
	}

	profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: ids,
	})
	if err != nil {
		return nil, err
	}

	for _, r := range followInfos {
		r.R.User_profile = ToProfile(ctx, profiles.GetProfiles()[string(r.User().ID())])
	}

	return followInfos, nil
}

func (_ *Resolver) GetFollowingCompaniesForCompany(ctx context.Context, input GetFollowingCompaniesForCompanyRequest) ([]CompanyFollowInfoResolver, error) {
	res, err := network.GetFollowingCompaniesForCompany(ctx, &networkRPC.IdWithCompanyFilter{Id: input.CompanyId, Filter: companyFilterToRPC(input)})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	follows := make([]CompanyFollowInfoResolver, 0, len(res.Follows))
	for _, follow := range res.Follows {
		res := companyFollowToResolver(follow)
		follows = append(follows, res)
	}

	ids := make([]string, 0, len(res.GetFollows()))
	for _, v := range res.GetFollows() {
		ids = append(ids, v.GetCompany().GetId())
	}

	companies := make(map[string]CompanyProfile, len(res.GetFollows()))

	comps, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	for _, pr := range comps.GetProfiles() {
		if pr != nil {
			companies[pr.GetId()] = toCompanyProfile(ctx, *pr)
		}
	}

	for i := range follows {
		if follows[i].R != nil {
			follows[i].R.Company_profile = companies[follows[i].R.Company.ID]
		}
	}

	return follows, nil
}

func (_ *Resolver) GetFollowerCompaniesForCompany(ctx context.Context, input GetFollowerCompaniesForCompanyRequest) ([]CompanyFollowInfoResolver, error) {
	res, err := network.GetFollowerCompaniesForCompany(ctx, &networkRPC.IdWithCompanyFilter{Id: input.CompanyId, Filter: companyFilterToRPC(input)})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	follows := make([]CompanyFollowInfoResolver, 0, len(res.Follows))
	for _, follow := range res.Follows {
		res := companyFollowToResolver(follow)
		follows = append(follows, res)
	}

	ids := make([]string, 0, len(res.GetFollows()))
	for _, v := range res.GetFollows() {
		ids = append(ids, v.GetCompany().GetId())
	}

	companies := make(map[string]CompanyProfile, len(res.GetFollows()))

	comps, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	for _, pr := range comps.GetProfiles() {
		if pr != nil {
			companies[pr.GetId()] = toCompanyProfile(ctx, *pr)
		}
	}

	for i := range follows {
		if follows[i].R != nil {
			follows[i].R.Company_profile = companies[follows[i].R.Company.ID]
		}
	}

	return follows, nil
}

func (_ *Resolver) GetSuggestedPeopleForCompany(ctx context.Context, input GetSuggestedPeopleForCompanyRequest) ([]UserSuggestionResolver, error) {
	res, err := network.GetSuggestedPeopleForCompany(ctx, &networkRPC.Company{Id: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	suggestions := make([]UserSuggestionResolver, 0, len(res.Suggestions))
	for _, suggestion := range res.Suggestions {
		res := suggestionToResolver(suggestion)
		suggestions = append(suggestions, res)
	}
	return suggestions, nil
}

func (_ *Resolver) GetSuggestedCompaniesForCompany(ctx context.Context, input GetSuggestedCompaniesForCompanyRequest) ([]CompanySuggestionResolver, error) {
	res, err := network.GetSuggestedCompaniesForCompany(ctx, &networkRPC.IdWithPagination{
		Id: input.CompanyId,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	},
	)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	suggestions := make([]CompanySuggestionResolver, 0, len(res.Suggestions))
	for _, suggestion := range res.Suggestions {
		res := companySuggestionToResolver(ctx, suggestion)
		suggestions = append(suggestions, res)
	}
	return suggestions, nil
}

func (_ *Resolver) GetConnectionsOfUser(ctx context.Context, input GetConnectionsOfUserRequest) (*ListOfFriendsResolver, error) {
	res, err := network.GetFriendsOfUser(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]Profile, 0, len(res.GetFriends()))

	for i := range res.GetFriends() {
		prof = append(prof, ToProfile(ctx, res.GetFriends()[i]))
	}

	return &ListOfFriendsResolver{
		R: &ListOfFriends{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowsOfUser(ctx context.Context, input GetConnectionsOfUserRequest) (*ListOfFriendsResolver, error) {
	res, err := network.GetFollowsOfUser(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]Profile, 0, len(res.GetFriends()))

	for i := range res.GetFriends() {
		prof = append(prof, ToProfile(ctx, res.GetFriends()[i]))
	}

	return &ListOfFriendsResolver{
		R: &ListOfFriends{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowersOfUser(ctx context.Context, input GetConnectionsOfUserRequest) (*ListOfFriendsResolver, error) {
	res, err := network.GetFollowersOfUser(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]Profile, 0, len(res.GetFriends()))

	for i := range res.GetFriends() {
		prof = append(prof, ToProfile(ctx, res.GetFriends()[i]))
	}

	return &ListOfFriendsResolver{
		R: &ListOfFriends{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowsCompaniesOfUser(ctx context.Context, input GetFollowsCompaniesOfUserRequest) (*ListOfCompaniesResolver, error) {
	res, err := network.GetFollowsCompaniesOfUser(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	log.Println("length:", len(res.GetCompanies()))

	prof := make([]CompanyProfile, 0, len(res.GetCompanies()))

	for i := range res.GetCompanies() {
		if res.GetCompanies()[i] == nil {
			return nil, errors.New("internal_error")
		}
		prof = append(prof, toCompanyProfile(ctx, *res.GetCompanies()[i]))
	}

	return &ListOfCompaniesResolver{
		R: &ListOfCompanies{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowersCompaniesOfUser(ctx context.Context, input GetFollowersCompaniesOfUserRequest) (*ListOfCompaniesResolver, error) {
	res, err := network.GetFollowersCompaniesOfUser(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]CompanyProfile, 0, len(res.GetCompanies()))

	for i := range res.GetCompanies() {
		if res.GetCompanies()[i] == nil {
			return nil, errors.New("internal_error")
		}
		prof = append(prof, toCompanyProfile(ctx, *res.GetCompanies()[i]))
	}

	return &ListOfCompaniesResolver{
		R: &ListOfCompanies{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowsOfCompany(ctx context.Context, input GetConnectionsOfUserRequest) (*ListOfFriendsResolver, error) {
	res, err := network.GetFollowsOfCompany(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]Profile, 0, len(res.GetFriends()))

	for i := range res.GetFriends() {
		prof = append(prof, ToProfile(ctx, res.GetFriends()[i]))
	}

	return &ListOfFriendsResolver{
		R: &ListOfFriends{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowersOfCompany(ctx context.Context, input GetConnectionsOfUserRequest) (*ListOfFriendsResolver, error) {
	res, err := network.GetFollowersOfCompany(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]Profile, 0, len(res.GetFriends()))

	for i := range res.GetFriends() {
		prof = append(prof, ToProfile(ctx, res.GetFriends()[i]))
	}

	return &ListOfFriendsResolver{
		R: &ListOfFriends{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowsCompaniesOfCompany(ctx context.Context, input GetFollowsCompaniesOfUserRequest) (*ListOfCompaniesResolver, error) {
	res, err := network.GetFollowsCompaniesOfCompany(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	log.Println("length:", len(res.GetCompanies()))

	prof := make([]CompanyProfile, 0, len(res.GetCompanies()))

	for i := range res.GetCompanies() {
		if res.GetCompanies()[i] == nil {
			return nil, errors.New("internal_error")
		}
		prof = append(prof, toCompanyProfile(ctx, *res.GetCompanies()[i]))
	}

	return &ListOfCompaniesResolver{
		R: &ListOfCompanies{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetFollowersCompaniesOfCompany(ctx context.Context, input GetFollowersCompaniesOfUserRequest) (*ListOfCompaniesResolver, error) {
	res, err := network.GetFollowersCompaniesOfCompany(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]CompanyProfile, 0, len(res.GetCompanies()))

	for i := range res.GetCompanies() {
		if res.GetCompanies()[i] == nil {
			return nil, errors.New("internal_error")
		}
		prof = append(prof, toCompanyProfile(ctx, *res.GetCompanies()[i]))
	}

	return &ListOfCompaniesResolver{
		R: &ListOfCompanies{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

// func (_ *Resolver) GetFollowOfUser(ctx context.Context, input GetConnectionsOfUserRequest) (*ListOfFriendsResolver, error) {
// 	res, err := network.GetFollowOfUser(ctx, &networkRPC.IdWithPagination{
// 		Id: input.User_id,
// 		Pagination: &networkRPC.Pagination{
// 			After:  NullIDToInt32(input.Pagination.After),
// 			Amount: NullToInt32(input.Pagination.First),
// 		},
// 	})
// 	if e, isErr := handleError(err); isErr {
// 		return nil, e
// 	}
//
// 	prof := make([]Profile, 0, len(res.GetFriends()))
//
// 	for i := range res.GetFriends() {
// 		prof = append(prof, ToProfile(res.GetFriends()[i]))
// 	}
//
// 	return &ListOfFriendsResolver{
// 		R: &ListOfFriends{
// 			Profiles: prof,
// 			Amount:   int32(res.GetAmount()),
// 		},
// 	}, nil
// }

func (_ *Resolver) GetMutualConnectionsOfUser(ctx context.Context, input GetConnectionsOfUserRequest) (*ListOfFriendsResolver, error) {
	res, err := network.GetMutualFriendsOfUser(ctx, &networkRPC.IdWithPagination{
		Id: input.User_id,
		Pagination: &networkRPC.Pagination{
			After:  NullIDToInt32(input.Pagination.After),
			Amount: NullToInt32(input.Pagination.First),
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}

	prof := make([]Profile, 0, len(res.GetFriends()))

	for i := range res.GetFriends() {
		prof = append(prof, ToProfile(ctx, res.GetFriends()[i]))
	}

	return &ListOfFriendsResolver{
		R: &ListOfFriends{
			Profiles: prof,
			Amount:   int32(res.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) SentEmailInvitation(ctx context.Context, input SentEmailInvitationRequest) (*SuccessResolver, error) {
	inv := userRPC.EmailInvitation{
		Address: input.Email,
		Name:    input.Name,
	}

	if input.Company_id != nil {
		inv.CompanyID = *input.Company_id
	}

	_, err := user.SentEmailInvitation(ctx, &inv)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{Success: true},
	}, nil
}

func (_ *Resolver) GetInvitation(ctx context.Context) (ListOfInvitationsResolver, error) {
	res, err := user.GetInvitation(ctx, &userRPC.Empty{})
	if err != nil {
		return ListOfInvitationsResolver{}, err
	}

	invitations := make([]Invitation, 0, len(res.GetInvitations()))

	for _, inv := range res.GetInvitations() {
		if inv != nil {
			invitations = append(invitations, Invitation{
				Email: inv.GetEmail(),
				Name:  inv.GetName(),
			})
		}
	}

	return ListOfInvitationsResolver{
		R: &ListOfInvitations{
			Amount:      res.GetAmount(),
			Invitations: invitations,
		},
	}, nil
}

func (_ *Resolver) GetInvitationForCompany(ctx context.Context, input GetInvitationForCompanyRequest) (ListOfInvitationsResolver, error) {
	res, err := user.GetInvitationForCompany(ctx, &userRPC.ID{
		ID: input.Company_id,
	})
	if err != nil {
		return ListOfInvitationsResolver{}, err
	}

	invitations := make([]Invitation, 0, len(res.GetInvitations()))

	for _, inv := range res.GetInvitations() {
		if inv != nil {
			invitations = append(invitations, Invitation{
				Email: inv.GetEmail(),
				Name:  inv.GetName(),
			})
		}
	}

	return ListOfInvitationsResolver{
		R: &ListOfInvitations{
			Amount:      res.GetAmount(),
			Invitations: invitations,
		},
	}, nil
}
