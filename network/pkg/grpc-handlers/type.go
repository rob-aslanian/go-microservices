package grpc_handlers

import (
	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
	"golang.org/x/net/context"
)

type NetworkService interface {
	SendFriendRequest(context.Context, string, string) (*model.Friendship, error)
	ApproveFriendRequest(ctx context.Context, key string) error
	DenyFriendRequest(ctx context.Context, key string) error
	IgnoreFriendRequest(ctx context.Context, key string) error
	GetFriendshipRequests(ctx context.Context, filter *model.FriendshipRequestFilter) ([]*model.Friendship, error)
	GetAllFriendships(context.Context, *model.FriendshipFilter) ([]*model.Friendship, error)
	GetAllFriendshipID(ctx context.Context) ([]string, error)
	Unfriend(ctx context.Context, friendId string) error
	IsFriend(ctx context.Context, userId string) (bool, error)
	Follow(ctx context.Context, userId string) error
	Unfollow(ctx context.Context, userId string) error
	GetFollowers(ctx context.Context, filter *model.UserFilter) ([]*model.Follow, error)
	GetFollowings(ctx context.Context, filter *model.UserFilter) ([]*model.Follow, error)
	AddToFavourites(ctx context.Context, userId string) error
	AddToFollowingsFavourites(ctx context.Context, companyId string) error
	AddToFollowingsFavouritesForCompany(ctx context.Context, companyId, refCompanyId string) error
	RemoveFromFavourites(ctx context.Context, userId string) error
	RemoveFromFollowingsFavourites(ctx context.Context, companyId string) error
	RemoveFromFollowingsFavouritesForCompany(ctx context.Context, companyId, refCompanyId string) error
	AddToCategory(ctx context.Context, userId, categoryName string) error
	AddToFollowingsCategory(ctx context.Context, companyId, categoryName string) error
	AddToFollowingsCategoryForCompany(ctx context.Context, companyId, refCompanyId, categoryName string) error
	RemoveFromCategory(ctx context.Context, userId, categoryName string) error
	RemoveFromFollowingsCategory(ctx context.Context, companyId, categoryName string) error
	RemoveFromFollowingsCategoryForCompany(ctx context.Context, companyId, refCompanyId, categoryName string) error
	BatchRemoveFromCategory(ctx context.Context, userIds []string, categoryName string, all bool) error
	BatchRemoveFromFollowingsCategory(ctx context.Context, companyIds []string, categoryName string, all bool) error
	BatchRemoveFromFollowingsCategoryForCompany(ctx context.Context, companyId string, companyIds []string, categoryName string, all bool) error
	GetFriendSuggestions(ctx context.Context, pagination *model.Pagination) ([]*model.UserSuggestion, error)

	MakeCompanyOwner(ctx context.Context, request *model.UserCompanyId) error
	IsCompanyOwner(ctx context.Context, request *model.UserCompanyId) (bool, error)
	MakeCompanyAdmin(ctx context.Context, request *model.AdminEdge) error
	GetAdminObject(ctx context.Context, companyId string) (*model.Admin, error)
	GetCompanyAdmins(ctx context.Context, companyId string) ([]*model.Admin, error)
	ChangeAdminLevel(ctx context.Context, request *model.AdminEdge) error
	DeleteCompanyAdmin(ctx context.Context, request *model.UserCompanyId) error
	GetUserCompanies(ctx context.Context, userId string) ([]*model.Admin, error)
	FollowCompany(ctx context.Context, companyId string) error
	UnfollowCompany(ctx context.Context, companyId string) error
	GetFollowerCompanies(ctx context.Context, filter *model.CompanyFilter) ([]*model.CompanyFollow, error)
	GetFilteredFollowingCompanies(ctx context.Context, filter *model.CompanyFilter) ([]*model.CompanyFollow, error)
	AddCompanyToFavourites(ctx context.Context, companyId string) error
	RemoveCompanyFromFavourites(ctx context.Context, companyId string) error
	AddCompanyToCategory(ctx context.Context, companyId, categoryName string) error
	RemoveCompanyFromCategory(ctx context.Context, companyId, categoryName string) error
	GetSuggestedCompanies(ctx context.Context, pagination *model.Pagination) ([]*model.CompanySuggestion, error)

	AddExperience(ctx context.Context, experience *model.AddExperienceRequest) error
	GetAllExperience(ctx context.Context) ([]*model.Experience, error)
	AskRecommendation(ctx context.Context, request *model.RecommendationRequestModel) error
	IgnoreRecommendationRequest(ctx context.Context, key string) error
	GetRequestedRecommendations(ctx context.Context, pagination *model.Pagination) ([]*model.RecommendationRequest, error)
	GetReceivedRecommendationRequests(ctx context.Context, pagination *model.Pagination) ([]*model.RecommendationRequest, error)
	WriteRecommendation(ctx context.Context, recommendation *model.RecommendationModel) error
	SetRecommendationVisibility(ctx context.Context, key string, visible bool) error
	GetReceivedRecommendations(ctx context.Context, pagination *model.Pagination) ([]*model.Recommendation, error)
	GetGivenRecommendations(ctx context.Context, pagination *model.Pagination) ([]*model.Recommendation, error)
	GetReceivedRecommendationsById(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.Recommendation, error)
	GetGivenRecommendationsById(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.Recommendation, error)
	GetHiddenRecommendationsById(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.Recommendation, error)

	BlockUser(ctx context.Context, id string) error
	UnblockUser(ctx context.Context, id string) error
	BlockCompany(ctx context.Context, id string) error
	UnblockCompany(ctx context.Context, id string) error
	GetBlockedUsersOrCompanies(ctx context.Context) ([]*model.BlockedUserOrCompany, error)
	// GetBlockedUsers(ctx context.Context) ([]*model.User, error)
	// GetBlockedCompanies(ctx context.Context) ([]*model.Company, error)

	BlockUserForCompany(ctx context.Context, companyID string, id string) error
	UnblockUserForCompany(ctx context.Context, companyID string, id string) error
	GetBlockedUsersForCompany(ctx context.Context, companyID string) ([]*model.User, error)

	GetFollowingsForCompany(ctx context.Context, request *model.IdWithUserFilter) ([]*model.Follow, error)
	GetFollowersForCompany(ctx context.Context, request *model.IdWithUserFilter) ([]*model.Follow, error)
	GetFollowingCompaniesForCompany(ctx context.Context, request *model.IdWithCompanyFilter) ([]*model.CompanyFollow, error)
	GetFollowerCompaniesForCompany(ctx context.Context, request *model.IdWithCompanyFilter) ([]*model.CompanyFollow, error)
	FollowForCompany(ctx context.Context, companyId, userId string) error
	UnfollowForCompany(ctx context.Context, companyId, userId string) error
	FollowCompanyForCompany(ctx context.Context, followerCompanyId, followingCompanyId string) error
	UnfollowCompanyForCompany(ctx context.Context, followerCompanyId, followingCompanyId string) error
	GetSuggestedPeopleForCompany(ctx context.Context, companyId string) ([]*model.UserSuggestion, error)
	GetSuggestedCompaniesForCompany(ctx context.Context, companyId string, pagination *model.Pagination) ([]*model.CompanySuggestion, error)

	IsBlocked(ctx context.Context, id string) (bool, error)
	IsBlockedCompany(ctx context.Context, id string) (bool, error)
	IsBlockedByUser(ctx context.Context, id string) (bool, error)
	IsBlockedCompanyByUser(ctx context.Context, id string) (bool, error)
	IsFollowing(ctx context.Context, id string) (bool, error)
	IsFollowingCompany(ctx context.Context, id string) (bool, error)
	IsFavourite(ctx context.Context, id string) (bool, error)
	IsFavouriteCompany(ctx context.Context, id string) (bool, error)

	IsBlockedForCompany(ctx context.Context, id string, companyID string) (bool, error)
	IsBlockedCompanyForCompany(ctx context.Context, id string, companyID string) (bool, error)
	IsBlockedByCompany(ctx context.Context, id string, companyID string) (bool, error)
	IsBlockedCompanyByCompany(ctx context.Context, id string, companyID string) (bool, error)
	IsFollowingForCompany(ctx context.Context, id string, companyID string) (bool, error)
	IsFollowingCompanyForCompany(ctx context.Context, id string, companyID string) (bool, error)

	GetNumberOfFollowersForCompany(ctx context.Context, companyId string) (int, error)
	GetFriendIdsOf(ctx context.Context, userId string) ([]string, error)
	GetUserCountings(ctx context.Context, userId string) (*model.UserCountings, error)
	GetCompanyCountings(ctx context.Context, companyID string) (*model.CompanyCountings, error)
	GetCategoryTree(ctx context.Context) (*model.CategoryTree, error)
	GetCategoryTreeForFollowings(ctx context.Context) (*model.CategoryTree, error)
	GetCategoryTreeForFollowingsForCompany(ctx context.Context, companyKey string) (*model.CategoryTree, error)
	CreateCategory(ctx context.Context, parent, name string) error
	CreateCategoryForFollowings(ctx context.Context, parent, name string) error
	CreateCategoryForFollowingsForCompany(ctx context.Context, companyId, parent, name string) error
	RemoveCategory(ctx context.Context, parent, name string) error
	RemoveCategoryForFollowings(ctx context.Context, parent, name string) error
	RemoveCategoryForFollowingsForCompany(ctx context.Context, companyId, parent, name string) error

	IsFriendRequestSend(ctx context.Context, userID string) (bool, error)
	IsFriendRequestRecieved(ctx context.Context, userID string) (bool, string, error)
	GetFriendshipID(ctx context.Context, userID string) (string, error)

	GetFriendsOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetMutualFriendsOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowsOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowersOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowsCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowersCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowsOfCompany(ctx context.Context, companyID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowersOfCompany(ctx context.Context, companyID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowsCompaniesOfCompany(ctx context.Context, companyID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetFollowersCompaniesOfCompany(ctx context.Context, companyID string, first uint32, after uint32) (friends interface{}, amount int64, err error)
	GetAmountOfMutualFriends(ctx context.Context, userID string) (int32, error)

	GetFollowersIDs(ctx context.Context, id string, isCompanyID bool) ([]string, error)
}

type GrpcHandlers struct {
	ns NetworkService
}

func NewGrpcHandlers(networkService NetworkService) *GrpcHandlers {
	return &GrpcHandlers{ns: networkService}
}
