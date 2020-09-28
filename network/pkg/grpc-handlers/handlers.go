package grpc_handlers

import (
	"github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
	"golang.org/x/net/context"
)

func (h *GrpcHandlers) SendFriendRequest(ctx context.Context, data *networkRPC.FriendshipRequest) (f *networkRPC.Friendship, err error) {
	defer recoverHandler(&err)

	friendship, e := h.ns.SendFriendRequest(ctx, data.FriendId, data.Description)
	panicIf(e)

	return friendship.ToRPC(), nil
}

func (h *GrpcHandlers) ApproveFriendRequest(ctx context.Context, data *networkRPC.Friendship) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.ApproveFriendRequest(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) DenyFriendRequest(ctx context.Context, data *networkRPC.Friendship) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.ApproveFriendRequest(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) IgnoreFriendRequest(ctx context.Context, data *networkRPC.Friendship) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.IgnoreFriendRequest(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetFriendRequests(ctx context.Context, data *networkRPC.FriendRequestFilter) (arr *networkRPC.FriendshipArr, err error) {
	defer recoverHandler(&err)

	filter := model.NewFriendshipRequestFilterFromRPC(data)
	friendships, e := h.ns.GetFriendshipRequests(ctx, filter)
	panicIf(e)

	return model.FriendshipArr(friendships).ToRPC(), nil
}

func (h *GrpcHandlers) GetAllFriendships(ctx context.Context, data *networkRPC.FriendshipFilter) (arr *networkRPC.FriendshipArr, err error) {
	defer recoverHandler(&err)

	filter := model.NewFriendshipFilterFromRPC(data)
	friendships, e := h.ns.GetAllFriendships(ctx, filter)
	panicIf(e)

	return model.FriendshipArr(friendships).ToRPC(), nil
}

func (h *GrpcHandlers) GetAllFriendshipsID(ctx context.Context, data *networkRPC.Empty) (arr *networkRPC.StringArr, err error) {
	defer recoverHandler(&err)

	friendships, e := h.ns.GetAllFriendshipID(ctx)
	panicIf(e)

	return &networkRPC.StringArr{
		List: friendships,
	}, nil
}

func (h *GrpcHandlers) Unfriend(ctx context.Context, data *networkRPC.User) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.Unfriend(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) IsFriend(ctx context.Context, data *networkRPC.User) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsFriend(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) Follow(ctx context.Context, data *networkRPC.User) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.Follow(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) Unfollow(ctx context.Context, data *networkRPC.User) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.Unfollow(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetFollowers(ctx context.Context, data *networkRPC.UserFilter) (response *networkRPC.FollowsArr, err error) {
	defer recoverHandler(&err)

	follows, e := h.ns.GetFollowers(ctx, model.NewUserFilterFromRPC(data))
	panicIf(e)

	return model.FollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) GetFollowings(ctx context.Context, data *networkRPC.UserFilter) (response *networkRPC.FollowsArr, err error) {
	defer recoverHandler(&err)

	follows, e := h.ns.GetFollowings(ctx, model.NewUserFilterFromRPC(data))
	panicIf(e)

	return model.FollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) AddToFavourites(ctx context.Context, data *networkRPC.User) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddToFavourites(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) AddToFollowingsFavourites(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddToFollowingsFavourites(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) AddToFollowingsFavouritesForCompany(ctx context.Context, data *networkRPC.CompanyCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddToFollowingsFavouritesForCompany(ctx, data.Company1Id, data.Company2Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveFromFavourites(ctx context.Context, data *networkRPC.User) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveFromFavourites(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveFromFollowingsFavourites(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveFromFollowingsFavourites(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveFromFollowingsFavouritesForCompany(ctx context.Context, data *networkRPC.CompanyCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveFromFollowingsFavouritesForCompany(ctx, data.Company1Id, data.Company2Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) AddToCategory(ctx context.Context, data *networkRPC.CategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddToCategory(ctx, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) AddToFollowingsCategory(ctx context.Context, data *networkRPC.CategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddToFollowingsCategory(ctx, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) AddToFollowingsCategoryForCompany(ctx context.Context, data *networkRPC.CategoryRequestWithCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddToFollowingsCategoryForCompany(ctx, data.CompanyId, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveFromCategory(ctx context.Context, data *networkRPC.CategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveFromCategory(ctx, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveFromFollowingsCategory(ctx context.Context, data *networkRPC.CategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveFromFollowingsCategory(ctx, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveFromFollowingsCategoryForCompany(ctx context.Context, data *networkRPC.CategoryRequestWithCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveFromFollowingsCategoryForCompany(ctx, data.CompanyId, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) BatchRemoveFromCategory(ctx context.Context, data *networkRPC.BatchRemoveFromCategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.BatchRemoveFromCategory(ctx, data.Ids, data.Category, data.All)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) BatchRemoveFromFollowingsCategory(ctx context.Context, data *networkRPC.BatchRemoveFromCategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.BatchRemoveFromFollowingsCategory(ctx, data.Ids, data.Category, data.All)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) BatchRemoveFromFollowingsCategoryForCompany(ctx context.Context, data *networkRPC.BatchRemoveFromCategoryRequestWithCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.BatchRemoveFromFollowingsCategoryForCompany(ctx, data.CompanyId, data.Request.Ids, data.Request.Category, data.Request.All)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetFriendSuggestions(ctx context.Context, data *networkRPC.Pagination) (response *networkRPC.UserSuggestionArr, err error) {
	defer recoverHandler(&err)

	suggestions, e := h.ns.GetFriendSuggestions(ctx, model.NewPagination(data))
	panicIf(e)

	return model.UserSuggestionArr(suggestions).ToRPC(), nil

}

func (h *GrpcHandlers) MakeCompanyOwner(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.MakeCompanyOwner(ctx, model.NewUserCompanyFromRPC(data))
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) IsCompanyOwner(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, e := h.ns.IsCompanyOwner(ctx, model.NewUserCompanyFromRPC(data))
	panicIf(e)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) MakeCompanyAdmin(ctx context.Context, data *networkRPC.MakeCompanyAdminRequest) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.MakeCompanyAdmin(ctx, model.NewAdminEdgeFromRPC(data))
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetAdminObject(ctx context.Context, data *networkRPC.Company) (response *networkRPC.AdminObject, err error) {
	defer recoverHandler(&err)

	admin, e := h.ns.GetAdminObject(ctx, data.Id)
	panicIf(e)

	return admin.ToRPC(), nil
}

func (h *GrpcHandlers) GetCompanyAdmins(ctx context.Context, data *networkRPC.Company) (response *networkRPC.AdminObjectArr, err error) {
	defer recoverHandler(&err)

	admins, e := h.ns.GetCompanyAdmins(ctx, data.Id)
	panicIf(e)

	return model.AdminArr(admins).ToRPC(), nil
}

func (h *GrpcHandlers) ChangeAdminLevel(ctx context.Context, data *networkRPC.MakeCompanyAdminRequest) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.ChangeAdminLevel(ctx, model.NewAdminEdgeFromRPC(data))
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) DeleteCompanyAdmin(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.DeleteCompanyAdmin(ctx, model.NewUserCompanyFromRPC(data))
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetUserCompanies(ctx context.Context, data *networkRPC.User) (response *networkRPC.AdminObjectArr, err error) {
	defer recoverHandler(&err)

	admins, e := h.ns.GetUserCompanies(ctx, data.Id)
	panicIf(e)

	return model.AdminArr(admins).ToRPC(), nil
}

func (h *GrpcHandlers) FollowCompany(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.FollowCompany(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) UnfollowCompany(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.UnfollowCompany(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetFollowerCompanies(ctx context.Context, data *networkRPC.CompanyFilter) (response *networkRPC.CompanyFollowsArr, err error) {
	defer recoverHandler(&err)

	follows, e := h.ns.GetFollowerCompanies(ctx, model.NewCompanyFilterFromRPC(data))
	panicIf(e)

	return model.CompanyFollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) GetFilteredFollowingCompanies(ctx context.Context, data *networkRPC.CompanyFilter) (response *networkRPC.CompanyFollowsArr, err error) {
	defer recoverHandler(&err)

	follows, e := h.ns.GetFilteredFollowingCompanies(ctx, model.NewCompanyFilterFromRPC(data))
	panicIf(e)

	return model.CompanyFollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) AddCompanyToFavourites(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddCompanyToFavourites(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveCompanyFromFavourites(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveCompanyFromFavourites(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) AddCompanyToCategory(ctx context.Context, data *networkRPC.CategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddCompanyToCategory(ctx, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) RemoveCompanyFromCategory(ctx context.Context, data *networkRPC.CategoryRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.RemoveCompanyFromCategory(ctx, data.Id, data.CategoryName)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetSuggestedCompanies(ctx context.Context, data *networkRPC.Pagination) (response *networkRPC.CompanySuggestionArr, err error) {
	defer recoverHandler(&err)

	companies, e := h.ns.GetSuggestedCompanies(ctx, model.NewPagination(data))
	panicIf(e)

	return model.CompanySuggestionArr(companies).ToRPC(), nil
}

func (h *GrpcHandlers) AddExperience(ctx context.Context, data *networkRPC.AddExperienceRequest) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.AddExperience(ctx, model.NewAddExperienceRequestFromRPC(data))
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetAllExperience(ctx context.Context, data *networkRPC.Empty) (response *networkRPC.ExperienceArr, err error) {
	defer recoverHandler(&err)

	experiences, e := h.ns.GetAllExperience(ctx)
	panicIf(e)

	return model.ExperienceArr(experiences).ToRPC(), nil
}

func (h *GrpcHandlers) AskRecommendation(ctx context.Context, data *networkRPC.RecommendationParams) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.AskRecommendation(ctx, model.NewRecommendationRequestModel(data))
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) IgnoreRecommendationRequest(ctx context.Context, data *networkRPC.RecommendationRequest) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.IgnoreRecommendationRequest(ctx, data.Id)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetRequestedRecommendations(ctx context.Context, data *networkRPC.Pagination) (response *networkRPC.RecommendationRequestArr, err error) {
	defer recoverHandler(&err)

	requests, err := h.ns.GetRequestedRecommendations(ctx, model.NewPagination(data))
	panicIf(err)

	return model.RecommendationRequestArr(requests).ToRPC(), nil
}

func (h *GrpcHandlers) GetReceivedRecommendationRequests(ctx context.Context, data *networkRPC.Pagination) (response *networkRPC.RecommendationRequestArr, err error) {
	defer recoverHandler(&err)

	requests, err := h.ns.GetReceivedRecommendationRequests(ctx, model.NewPagination(data))
	panicIf(err)

	return model.RecommendationRequestArr(requests).ToRPC(), nil
}

func (h *GrpcHandlers) WriteRecommendation(ctx context.Context, data *networkRPC.RecommendationParams) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.WriteRecommendation(ctx, model.NewRecommendationModel(data))
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) SetRecommendationVisibility(ctx context.Context, data *networkRPC.RecommendationVisibility) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.SetRecommendationVisibility(ctx, data.RecommendationId, data.Visible)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetReceivedRecommendations(ctx context.Context, data *networkRPC.Pagination) (response *networkRPC.RecommendationArr, err error) {
	defer recoverHandler(&err)

	recommendations, err := h.ns.GetReceivedRecommendations(ctx, model.NewPagination(data))
	panicIf(err)

	return model.RecommendationArr(recommendations).ToRPC(), nil
}

func (h *GrpcHandlers) GetGivenRecommendations(ctx context.Context, data *networkRPC.Pagination) (response *networkRPC.RecommendationArr, err error) {
	defer recoverHandler(&err)

	recommendations, err := h.ns.GetGivenRecommendations(ctx, model.NewPagination(data))
	panicIf(err)

	return model.RecommendationArr(recommendations).ToRPC(), nil
}

func (h *GrpcHandlers) GetReceivedRecommendationById(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.RecommendationArr, err error) {
	defer recoverHandler(&err)

	recommendations, err := h.ns.GetReceivedRecommendationsById(ctx, data.Id, model.NewPagination(data.Pagination))
	panicIf(err)

	return model.RecommendationArr(recommendations).ToRPC(), nil
}

func (h *GrpcHandlers) GetGivenRecommendationsById(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.RecommendationArr, err error) {
	defer recoverHandler(&err)

	recommendations, err := h.ns.GetGivenRecommendationsById(ctx, data.Id, model.NewPagination(data.Pagination))
	panicIf(err)

	return model.RecommendationArr(recommendations).ToRPC(), nil
}

func (h *GrpcHandlers) GetHiddenRecommendationByID(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.RecommendationArr, err error) {
	defer recoverHandler(&err)

	recommendations, err := h.ns.GetHiddenRecommendationsById(ctx, data.Id, model.NewPagination(data.Pagination))
	panicIf(err)

	return model.RecommendationArr(recommendations).ToRPC(), nil
}

func (h *GrpcHandlers) BlockUser(ctx context.Context, data *networkRPC.User) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.BlockUser(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) UnblockUser(ctx context.Context, data *networkRPC.User) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.UnblockUser(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) BlockCompany(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.BlockCompany(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) UnblockCompany(ctx context.Context, data *networkRPC.Company) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.UnblockCompany(ctx, data.Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetblockedUsersOrCompanies(ctx context.Context, data *networkRPC.Empty) (response *networkRPC.BlockedUserOrCompanyArr, err error) {
	defer recoverHandler(&err)

	list, e := h.ns.GetBlockedUsersOrCompanies(ctx)
	panicIf(e)

	return model.BlockerUserOrCompanyArr(list).ToRPC(), nil
}

func (h *GrpcHandlers) BlockUserForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.BlockUserForCompany(ctx, data.GetCompanyId(), data.GetUserId())
	if err != nil {
		return nil, err
	}

	return EMPTY, nil
}

func (h *GrpcHandlers) UnblockUserForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.UnblockUserForCompany(ctx, data.GetCompanyId(), data.GetUserId())
	if err != nil {
		return nil, err
	}

	return EMPTY, nil
}

func (h *GrpcHandlers) GetBlockedUsersForCompany(ctx context.Context, data *networkRPC.Company) (response *networkRPC.UserArr, err error) {
	defer recoverHandler(&err)

	users, e := h.ns.GetBlockedUsersForCompany(ctx, data.GetId())
	panicIf(e)

	return model.UserArr(users).ToRPC(), nil
}

//func (h *GrpcHandlers) GetBlockedUsers(ctx context.Context, data *networkRPC.Empty) (response *networkRPC.UserArr, err error) {
//	defer recoverHandler(&err)
//
//	users, e := h.ns.GetBlockedUsers(ctx)
//	panicIf(e)
//
//	return model.UserArr(users).ToRPC(), nil
//}
//
//func (h *GrpcHandlers) GetBlockedCompanies(ctx context.Context, data *networkRPC.Empty) (response *networkRPC.CompanyArr, err error) {
//	defer recoverHandler(&err)
//
//	companies, e := h.ns.GetBlockedCompanies(ctx)
//	panicIf(e)
//
//	return model.CompanyArr(companies).ToRPC(), nil
//}

func (h *GrpcHandlers) GetFollowingsForCompany(ctx context.Context, data *networkRPC.IdWithUserFilter) (response *networkRPC.FollowsArr, err error) {
	defer recoverHandler(&err)

	filter := model.NewIdWithUserFilterFromRPC(data)
	follows, e := h.ns.GetFollowingsForCompany(ctx, filter)
	panicIf(e)

	return model.FollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) GetFollowersForCompany(ctx context.Context, data *networkRPC.IdWithUserFilter) (response *networkRPC.FollowsArr, err error) {
	defer recoverHandler(&err)

	filter := model.NewIdWithUserFilterFromRPC(data)
	follows, e := h.ns.GetFollowersForCompany(ctx, filter)
	panicIf(e)

	return model.FollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) GetFollowerCompaniesForCompany(ctx context.Context, data *networkRPC.IdWithCompanyFilter) (response *networkRPC.CompanyFollowsArr, err error) {
	defer recoverHandler(&err)

	filter := model.NewIdWithCompanyFilterFromRPC(data)
	follows, e := h.ns.GetFollowerCompaniesForCompany(ctx, filter)
	panicIf(e)

	return model.CompanyFollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) GetFollowingCompaniesForCompany(ctx context.Context, data *networkRPC.IdWithCompanyFilter) (response *networkRPC.CompanyFollowsArr, err error) {
	defer recoverHandler(&err)

	filter := model.NewIdWithCompanyFilterFromRPC(data)
	follows, e := h.ns.GetFollowingCompaniesForCompany(ctx, filter)
	panicIf(e)

	return model.CompanyFollowArr(follows).ToRPC(), nil
}

func (h *GrpcHandlers) FollowForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.FollowForCompany(ctx, data.CompanyId, data.UserId)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) UnfollowForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.UnfollowForCompany(ctx, data.CompanyId, data.UserId)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) FollowCompanyForCompany(ctx context.Context, data *networkRPC.CompanyCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.FollowCompanyForCompany(ctx, data.Company1Id, data.Company2Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) UnfollowCompanyForCompany(ctx context.Context, data *networkRPC.CompanyCompanyId) (empty *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	e := h.ns.UnfollowCompanyForCompany(ctx, data.Company1Id, data.Company2Id)
	panicIf(e)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetSuggestedPeopleForCompany(ctx context.Context, data *networkRPC.Company) (response *networkRPC.UserSuggestionArr, err error) {
	defer recoverHandler(&err)

	res, e := h.ns.GetSuggestedPeopleForCompany(ctx, data.Id)
	panicIf(e)

	return model.UserSuggestionArr(res).ToRPC(), nil
}

func (h *GrpcHandlers) GetSuggestedCompaniesForCompany(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.CompanySuggestionArr, err error) {
	defer recoverHandler(&err)
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("companyId", data.Id)
	}
	res, e := h.ns.GetSuggestedCompaniesForCompany(ctx, data.GetId(), model.NewPagination(data.GetPagination()))
	panicIf(e)

	return model.CompanySuggestionArr(res).ToRPC(), nil
}

func (h *GrpcHandlers) GetNumberOfFollowersForCompany(ctx context.Context, data *networkRPC.Company) (response *networkRPC.IntValue, err error) {
	defer recoverHandler(&err)

	value, err := h.ns.GetNumberOfFollowersForCompany(ctx, data.Id)
	panicIf(err)

	return &networkRPC.IntValue{Value: int32(value)}, nil
}

func (h *GrpcHandlers) IsBlocked(ctx context.Context, data *networkRPC.User) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlocked(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}
func (h *GrpcHandlers) IsBlockedCompany(ctx context.Context, data *networkRPC.Company) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlockedCompany(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsBlockedForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlockedForCompany(ctx, data.GetUserId(), data.GetCompanyId())
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsBlockedCompanyForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlockedCompanyForCompany(ctx, data.GetUserId(), data.GetCompanyId())
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsBlockedByUser(ctx context.Context, data *networkRPC.User) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlockedByUser(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}
func (h *GrpcHandlers) IsBlockedCompanyByUser(ctx context.Context, data *networkRPC.Company) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlockedCompanyByUser(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsBlockedByCompany(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlockedByCompany(ctx, data.GetUserId(), data.GetCompanyId())
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsBlockedCompanyByCompany(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsBlockedCompanyByCompany(ctx, data.GetUserId(), data.GetCompanyId())
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsFollowing(ctx context.Context, data *networkRPC.User) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsFollowing(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}
func (h *GrpcHandlers) IsFollowingCompany(ctx context.Context, data *networkRPC.Company) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsFollowingCompany(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsFollowingForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsFollowingForCompany(ctx, data.GetUserId(), data.GetCompanyId())
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsFollowingCompanyForCompany(ctx context.Context, data *networkRPC.UserCompanyId) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsFollowingCompanyForCompany(ctx, data.GetUserId(), data.GetCompanyId())
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (h *GrpcHandlers) IsFavourite(ctx context.Context, data *networkRPC.User) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsFavourite(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}
func (h *GrpcHandlers) IsFavouriteCompany(ctx context.Context, data *networkRPC.Company) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	val, err := h.ns.IsFavouriteCompany(ctx, data.Id)
	panicIf(err)

	return &networkRPC.BooleanValue{Value: val}, nil
}

func (this *GrpcHandlers) GetFriendIdsOf(ctx context.Context, data *networkRPC.User) (response *networkRPC.StringArr, err error) {
	defer recoverHandler(&err)

	arr, err := this.ns.GetFriendIdsOf(ctx, data.Id)
	panicIf(err)

	return &networkRPC.StringArr{List: arr}, nil
}

func (this *GrpcHandlers) GetUserCountings(ctx context.Context, data *networkRPC.User) (response *networkRPC.UserCountings, err error) {
	defer recoverHandler(&err)

	countings, err := this.ns.GetUserCountings(ctx, data.Id)
	panicIf(err)

	return countings.ToRPC(), nil
}

func (this *GrpcHandlers) GetCompanyCountings(ctx context.Context, data *networkRPC.Company) (response *networkRPC.CompanyCountings, err error) {
	defer recoverHandler(&err)

	countings, err := this.ns.GetCompanyCountings(ctx, data.Id)
	panicIf(err)

	return &networkRPC.CompanyCountings{
		Employees:  countings.AmountOfEmployees,
		Followers:  countings.AmountOfFollowers,
		Followings: countings.AmountOfFollowings,
	}, nil
}

func (this *GrpcHandlers) GetCategoryTree(ctx context.Context, data *networkRPC.Empty) (response *networkRPC.CategoryTree, err error) {
	defer recoverHandler(&err)

	tree, err := this.ns.GetCategoryTree(ctx)
	panicIf(err)

	return tree.ToRPC(), nil
}

func (this *GrpcHandlers) GetCategoryTreeForFollowings(ctx context.Context, data *networkRPC.Empty) (response *networkRPC.CategoryTree, err error) {
	defer recoverHandler(&err)

	tree, err := this.ns.GetCategoryTreeForFollowings(ctx)
	panicIf(err)

	return tree.ToRPC(), nil
}

func (this *GrpcHandlers) GetCategoryTreeForFollowingsForCompany(ctx context.Context, data *networkRPC.Company) (response *networkRPC.CategoryTree, err error) {
	defer recoverHandler(&err)

	tree, err := this.ns.GetCategoryTreeForFollowingsForCompany(ctx, data.Id)
	panicIf(err)

	return tree.ToRPC(), nil
}

func (this *GrpcHandlers) CreateCategory(ctx context.Context, data *networkRPC.CategoryPath) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.CreateCategory(ctx, data.Parent, data.Name)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) CreateCategoryForFollowings(ctx context.Context, data *networkRPC.CategoryPath) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.CreateCategoryForFollowings(ctx, data.Parent, data.Name)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) CreateCategoryForFollowingsForCompany(ctx context.Context, data *networkRPC.CategoryPathWithCompanyId) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.CreateCategoryForFollowingsForCompany(ctx, data.CompanyId, data.CategoryPath.Parent, data.CategoryPath.Name)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) RemoveCategory(ctx context.Context, data *networkRPC.CategoryPath) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.RemoveCategory(ctx, data.Parent, data.Name)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) RemoveCategoryForFollowings(ctx context.Context, data *networkRPC.CategoryPath) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.RemoveCategoryForFollowings(ctx, data.Parent, data.Name)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) RemoveCategoryForFollowingsForCompany(ctx context.Context, data *networkRPC.CategoryPathWithCompanyId) (response *networkRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.RemoveCategoryForFollowingsForCompany(ctx, data.CompanyId, data.CategoryPath.Parent, data.CategoryPath.Name)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) IsFriendRequestSend(ctx context.Context, data *networkRPC.User) (response *networkRPC.BooleanValue, err error) {
	defer recoverHandler(&err)

	value, err := this.ns.IsFriendRequestSend(ctx, data.GetId())
	panicIf(err)

	return &networkRPC.BooleanValue{Value: value}, nil
}

func (this *GrpcHandlers) IsFriendRequestRecieved(ctx context.Context, data *networkRPC.User) (response *networkRPC.FriendRequestID, err error) {
	defer recoverHandler(&err)

	value, friendshipID, err := this.ns.IsFriendRequestRecieved(ctx, data.GetId())
	panicIf(err)

	return &networkRPC.FriendRequestID{
		FriendshipID: friendshipID,
		Recivied:     value,
	}, nil
}

func (this *GrpcHandlers) GetFriendshipID(ctx context.Context, data *networkRPC.User) (response *networkRPC.Friendship, err error) {
	defer recoverHandler(&err)
	friendshipID, err := this.ns.GetFriendshipID(ctx, data.GetId())
	panicIf(err)

	return &networkRPC.Friendship{
		Id: friendshipID,
	}, nil
}

func (this *GrpcHandlers) GetMutualFriendsOfUser(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.FriendList, err error) {
	defer recoverHandler(&err)

	mutualInterface, amount, err := this.ns.GetMutualFriendsOfUser(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var mutual []*userRPC.Profile

	if mu, ok := mutualInterface.([]*userRPC.Profile); ok {
		mutual = make([]*userRPC.Profile, 0, len(mu))
		for i := range mu {
			mutual = append(mutual, mu[i])
		}
	}

	return &networkRPC.FriendList{
		Friends: mutual,
		Amount:  amount,
	}, nil
}

func (this *GrpcHandlers) GetFriendsOfUser(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.FriendList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFriendsOfUser(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var friends []*userRPC.Profile

	if mu, ok := friendsInterface.([]*userRPC.Profile); ok {
		friends = make([]*userRPC.Profile, 0, len(mu))
		for i := range mu {
			friends = append(friends, mu[i])
		}
	}

	return &networkRPC.FriendList{
		Friends: friends,
		Amount:  amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowsOfUser(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.FriendList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowsOfUser(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var friends []*userRPC.Profile

	if mu, ok := friendsInterface.([]*userRPC.Profile); ok {
		friends = make([]*userRPC.Profile, 0, len(mu))
		for i := range mu {
			friends = append(friends, mu[i])
		}
	}

	return &networkRPC.FriendList{
		Friends: friends,
		Amount:  amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowersOfUser(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.FriendList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowersOfUser(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var friends []*userRPC.Profile

	if mu, ok := friendsInterface.([]*userRPC.Profile); ok {
		friends = make([]*userRPC.Profile, 0, len(mu))
		for i := range mu {
			friends = append(friends, mu[i])
		}
	}

	return &networkRPC.FriendList{
		Friends: friends,
		Amount:  amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowsCompaniesOfUser(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.CompanyList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowsCompaniesOfUser(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var company []*companyRPC.Profile

	if mu, ok := friendsInterface.([]*companyRPC.Profile); ok {
		company = make([]*companyRPC.Profile, 0, len(mu))
		for i := range mu {
			company = append(company, mu[i])
		}
	}

	return &networkRPC.CompanyList{
		Companies: company,
		Amount:    amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowersCompaniesOfUser(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.CompanyList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowersCompaniesOfUser(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var company []*companyRPC.Profile

	if mu, ok := friendsInterface.([]*companyRPC.Profile); ok {
		company = make([]*companyRPC.Profile, 0, len(mu))
		for i := range mu {
			company = append(company, mu[i])
		}
	}

	return &networkRPC.CompanyList{
		Companies: company,
		Amount:    amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowsOfCompany(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.FriendList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowsOfCompany(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var friends []*userRPC.Profile

	if mu, ok := friendsInterface.([]*userRPC.Profile); ok {
		friends = make([]*userRPC.Profile, 0, len(mu))
		for i := range mu {
			friends = append(friends, mu[i])
		}
	}

	return &networkRPC.FriendList{
		Friends: friends,
		Amount:  amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowersOfCompany(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.FriendList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowersOfCompany(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var friends []*userRPC.Profile

	if mu, ok := friendsInterface.([]*userRPC.Profile); ok {
		friends = make([]*userRPC.Profile, 0, len(mu))
		for i := range mu {
			friends = append(friends, mu[i])
		}
	}

	return &networkRPC.FriendList{
		Friends: friends,
		Amount:  amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowsCompaniesOfCompany(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.CompanyList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowsCompaniesOfCompany(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var company []*companyRPC.Profile

	if mu, ok := friendsInterface.([]*companyRPC.Profile); ok {
		company = make([]*companyRPC.Profile, 0, len(mu))
		for i := range mu {
			company = append(company, mu[i])
		}
	}

	return &networkRPC.CompanyList{
		Companies: company,
		Amount:    amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowersCompaniesOfCompany(ctx context.Context, data *networkRPC.IdWithPagination) (response *networkRPC.CompanyList, err error) {
	defer recoverHandler(&err)

	friendsInterface, amount, err := this.ns.GetFollowersCompaniesOfCompany(
		ctx,
		data.GetId(),
		uint32(data.GetPagination().GetAmount()),
		uint32(data.GetPagination().GetAfter()),
	)
	panicIf(err)

	var company []*companyRPC.Profile

	if mu, ok := friendsInterface.([]*companyRPC.Profile); ok {
		company = make([]*companyRPC.Profile, 0, len(mu))
		for i := range mu {
			company = append(company, mu[i])
		}
	}

	return &networkRPC.CompanyList{
		Companies: company,
		Amount:    amount,
	}, nil
}

func (this *GrpcHandlers) GetAmountOfMutualFriends(ctx context.Context, data *networkRPC.ID) (response *networkRPC.Amount, err error) {
	defer recoverHandler(&err)

	amount, err := this.ns.GetAmountOfMutualFriends(
		ctx,
		data.GetID(),
	)
	panicIf(err)

	return &networkRPC.Amount{
		Amount: amount,
	}, nil
}

func (this *GrpcHandlers) GetFollowersIDs(ctx context.Context, data *networkRPC.GetFollowersIDsRequest) (response *networkRPC.IDs, err error) {
	defer recoverHandler(&err)

	ids, err := this.ns.GetFollowersIDs(ctx, data.GetID(), data.GetIsCompany())
	if err != nil {
		return nil, err
	}

	return &networkRPC.IDs{
		IDs: ids,
	}, nil
}
