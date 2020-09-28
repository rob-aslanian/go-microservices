package services

import (
	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
	"golang.org/x/net/context"
)

type NetworkRepo interface {
	InsertFriendship(context.Context, *model.NewFriendship) (string, error)
	ChangeFriendshipStatus(ctx context.Context, friendshipKey string, newStatus model.FriendshipStatus) error
	GetFriendship(context.Context, string) (*model.Friendship, error)
	GetFriendshipRequests(ctx context.Context, userKey string, filter *model.FriendshipRequestFilter) ([]*model.Friendship, error)
	GetAllFriendship(ctx context.Context, userKey string, filter *model.FriendshipRequestFilter) ([]*model.Friendship, error)
	GetAllFriendshipID(ctx context.Context, userKey string) ([]string, error)
	GetFriendsOf(ctx context.Context, userKey string, filter *model.FriendshipFilter) ([]*model.Friendship, error)
	RemoveFriendship(ctx context.Context, user1Key, user2Key string) error
	Follow(ctx context.Context, follow *model.FollowRequest) (string, error)
	Unfollow(ctx context.Context, user1Key, user2Key string) error
	IsFriend(ctx context.Context, user1Key, user2Key string) (bool, error)
	GetFollowers(ctx context.Context, userKey string, filter *model.UserFilter) ([]*model.Follow, error)
	GetFollowings(ctx context.Context, userKey string, filter *model.UserFilter) ([]*model.Follow, error)
	AddToCategory(ctx context.Context, categoryRelation *model.CategoryRelation) (string, error)
	AddToFollowingsCategory(ctx context.Context, categoryRelation *model.CategoryRelation) (string, error)
	RemoveFromCategory(ctx context.Context, categoryRelation *model.CategoryRelation) error
	RemoveFromFollowingsCategory(ctx context.Context, categoryRelation *model.CategoryRelation) error
	BatchRemoveFromCategory(ctx context.Context, userKey string, userIds []string, categoryName string, all bool) error
	BatchRemoveFromCategoryForFollowings(ctx context.Context, ownerKey string, referalKeys []string, categoryName string, all bool) error
	GetFriendSuggestions(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.UserSuggestion, error)

	MakeCompanyOwner(ctx context.Context, request *model.UserCompanyId) error
	IsCompanyOwner(ctx context.Context, request *model.UserCompanyId) (bool, error)
	MakeCompanyAdmin(ctx context.Context, request *model.AdminEdge) error
	GetAdminObject(ctx context.Context, request *model.UserCompanyId) (*model.Admin, error)
	GetCompanyAdmins(ctx context.Context, companyKey string) ([]*model.Admin, error)
	ChangeAdminLevel(ctx context.Context, adminKey string, level model.AdminLevel) error
	DeleteCompanyAdmin(ctx context.Context, adminKey string) error
	GetUserCompanies(ctx context.Context, userKey string) ([]*model.Admin, error)
	GetFollowerCompanies(ctx context.Context, userKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error)
	GetFilteredFollowingCompanies(ctx context.Context, userKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error)
	GetSuggestedCompanies(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.CompanySuggestion, error)

	AddExperience(ctx context.Context, experience *model.AddExperienceRequest) error
	GetAllExperience(ctx context.Context, userKey string) ([]*model.Experience, error)
	SaveRecommendationRequest(ctx context.Context, request *model.RecommendationRequestModel) (string, error)
	GetRecommendationRequestModel(ctx context.Context, key string) (*model.RecommendationRequestModel, error)
	UpdateRecommendationRequest(ctx context.Context, key string, request *model.RecommendationRequestModel) error
	GetRequestedRecommendations(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.RecommendationRequest, error)
	GetReceivedRecommendationRequests(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.RecommendationRequest, error)
	WriteRecommendation(ctx context.Context, recommendation *model.RecommendationModel) error
	GetRecommendationModel(ctx context.Context, key string) (*model.RecommendationModel, error)
	SetRecommendationVisibility(ctx context.Context, key string, visible bool) error
	GetReceivedRecommendations(ctx context.Context, userKey string, pagination *model.Pagination /*, isMe bool*/) ([]*model.Recommendation, error)
	GetGivenRecommendations(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.Recommendation, error)
	GetHiddenRecommendations(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.Recommendation, error)

	Block(ctx context.Context, blockerId, blockedId string) error
	Unblock(ctx context.Context, blockerId, blockedId string) error
	GetBlockedUsersOrCompanies(ctx context.Context, userKey string) ([]*model.BlockedUserOrCompany, error)
	GetBlockedUsersForCompany(ctx context.Context, userKey string) ([]*model.User, error)
	//GetBlockedUsers(ctx context.Context, userKey string) ([]*model.User, error)
	//GetBlockedCompanies(ctx context.Context, userKey string) ([]*model.Company, error)

	GetFollowingsForCompany(ctx context.Context, companyKey string, filter *model.UserFilter) ([]*model.Follow, error)
	GetFollowersForCompany(ctx context.Context, companyKey string, filter *model.UserFilter) ([]*model.Follow, error)
	GetFollowingCompaniesForCompany(ctx context.Context, companyKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error)
	GetFollowerCompaniesForCompany(ctx context.Context, companyKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error)
	GetSuggestedPeopleForCompany(ctx context.Context, companyKey string) ([]*model.UserSuggestion, error)
	GetSuggestedCompaniesForCompany(ctx context.Context, companyId string, pagination *model.Pagination) ([]*model.CompanySuggestion, error)
	BatchRemoveFromCategoryForFollowingsForCompany(ctx context.Context, companyId string, companyIds []string, categoryName string, all bool) error

	IsBlocked(ctx context.Context, id1, id2 string) (bool, error)
	IsBlockedByUser(ctx context.Context, id1, id2 string) (bool, error)
	IsFollowing(ctx context.Context, id1, id2 string) (bool, error)
	IsFavourite(ctx context.Context, id1, id2 string) (bool, error)
	GetNumberOfFollowers(ctx context.Context, id string) (int, error)
	GetUserCountings(ctx context.Context, userId string) (*model.UserCountings, error)
	GetCompanyCountings(ctx context.Context, companyID string) (*model.CompanyCountings, error)
	SetCategoryTree(ctx context.Context, tree *model.CategoryTree) error
	SetCategoryTreeForFollowings(ctx context.Context, tree *model.CategoryTree) error
	GetCategoryTree(ctx context.Context, userKey string) (*model.CategoryTree, error)
	GetCategoryTreeForFollowings(ctx context.Context, key string) (*model.CategoryTree, error)
	ClearCategory(ctx context.Context, userKey, category string) error
	ClearCategoryForFollowings(ctx context.Context, ownerKey, category string) error
	ClearCategoryForFollowingsForCompany(ctx context.Context, companyKey, category string) error

	GetFriendsOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowsOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowersOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowsCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowersCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowsOfCompany(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowersOfCompany(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowsCompaniesOfCompany(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetFollowersCompaniesOfCompany(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error)
	GetMutualFriendsOfUser(ctx context.Context, senderID string, userID string, first uint32, after uint32) ([]string, int64, error)
	RemoveFriendshipByID(ctx context.Context, key string) error
	GetAmountOfMutualFriends(ctx context.Context, senderID string, userID string) (int32, error)
	GetFollowersIDs(ctx context.Context, id string, isCompany bool) ([]string, error)
}

type AuthClient interface {
	GetUserId(context.Context, string) (string, error)
}

type UserClient interface {
	GetProfilesByIDs(ctx context.Context, ids []string) (interface{}, error)
	GetConectionsPrivacy(ctx context.Context, userID string) (string, error)
}

type CompanyClient interface {
	GetProfilesByIDs(ctx context.Context, ids []string) (interface{}, error)
}

type ChatClient interface {
	GetProfilesByIDs(ctx context.Context, senderID, targetID string, value bool) error
}

type MQClient interface {
	SendNewFollow(userID string, message *model.NewFollow) error
	SendNewConnection(userID string, message *model.NewConnectionRequest) error
	ApproveConnectionRequest(userID string, message *model.NewApproveConnectionRequest) error
	SendNewRecommendationRequest(userID string, message *model.NewRecommendationRequest) error
	SendNewRecommendation(userID string, message *model.NewRecommendation) error
}

type NetworkValidator interface {
	ValidateStruct(value interface{}) error
	ValidateUserId(id string) error
	ValidateId(id string) error
	ValidateKey(key string) error
	ValidateAddExperienceRequest(req *model.AddExperienceRequest) error
}

type NetworkService struct {
	repo      NetworkRepo
	auth      AuthClient
	user      UserClient
	company   CompanyClient
	chat      ChatClient
	validator NetworkValidator
	mq        MQClient
}

func NewNetworkService(repo NetworkRepo, validator NetworkValidator, auth AuthClient, user UserClient, company CompanyClient, chat ChatClient, mq MQClient) *NetworkService {
	return &NetworkService{
		repo:      repo,
		validator: validator,
		auth:      auth,
		user:      user,
		company:   company,
		chat:      chat,
		mq:        mq,
	}
}
