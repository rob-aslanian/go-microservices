package network

import (
	"context"
	"log"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/utils"
	"gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
)

func (r *NetworkRepo) GetFriendship(ctx context.Context, friendshipKey string) (*model.Friendship, error) {
	md := utils.ExtractMetadata(ctx, "user-id")

	res, err := r.db.Query(ctx, GET_FRIENDSHIP_QUERY, map[string]interface{}{
		"id":   ToFriendshipId(friendshipKey),
		"myId": ToUserId(md["user-id"]),
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var fr model.Friendship
	_, err = res.ReadDocument(ctx, &fr)
	if err != nil {
		return nil, err
	}
	return &fr, nil
}

func (r *NetworkRepo) InsertFriendship(ctx context.Context, friendship *model.NewFriendship) (string, error) {
	meta, err := r.friendships.CreateDocument(ctx, friendship)
	if err != nil {
		return "", err
	}
	return meta.Key, nil
}

func (r *NetworkRepo) ChangeFriendshipStatus(ctx context.Context, friendshipKey string, newStatus model.FriendshipStatus) error {
	var friendship model.NewFriendship
	_, err := r.friendships.ReadDocument(ctx, friendshipKey, &friendship)
	if err != nil {
		return err
	}
	md := utils.ExtractMetadata(ctx, "user-id")
	if friendship.ReceiverId != ToUserId(md["user-id"]) {
		return errors.New("You can not respond to this friend request")
	}

	if friendship.Status != model.FriendshipStatus_Requested && friendship.Status != model.FriendshipStatus_Ignored {
		return errors.New("You can not change status of this friend request")
	}
	_, err = r.friendships.UpdateDocument(ctx, friendshipKey, map[string]interface{}{
		"status":       newStatus,
		"responded_at": time.Now(),
	})

	return err
}

func (r *NetworkRepo) GetFriendshipRequests(ctx context.Context, userKey string, filter *model.FriendshipRequestFilter) ([]*model.Friendship, error) {
	query, params := GET_FILTERED_FRIENDSHIP_REQUEST_QUERY(ToUserId(userKey), filter)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFriendships := make([]*model.Friendship, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var friendship model.Friendship
		cursor.ReadDocument(ctx, &friendship)
		allFriendships[i] = &friendship
	}

	return allFriendships, nil
}

func (r *NetworkRepo) GetAllFriendship(ctx context.Context, userKey string, filter *model.FriendshipRequestFilter) ([]*model.Friendship, error) {
	query, params := GET_ALL_FRIENDSHIP_QUERY(ToUserId(userKey), filter)
	// ---

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFriendships := make([]*model.Friendship, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var friendship model.Friendship
		cursor.ReadDocument(ctx, &friendship)
		allFriendships[i] = &friendship
	}

	return allFriendships, nil
}

func (r *NetworkRepo) GetAllFriendshipID(ctx context.Context, userKey string) ([]string, error) {
	query, params := GET_ALL_FRIENDSHIPID_QUERY(ToUserId(userKey))

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	ids := make([]string, 0, cursor.Count())
	for cursor.HasMore() {
		result := struct {
			IDs string `json:"ids"`
		}{}
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, err
		}
		ids = append(ids, result.IDs)
	}

	return ids, nil
}

func (r *NetworkRepo) GetFriendsOf(ctx context.Context, userKey string, filter *model.FriendshipFilter) ([]*model.Friendship, error) {
	for i := 0; i < len(filter.Companies); i++ {
		filter.Companies[i] = ToCompanyId(filter.Companies[i])
	}

	query, params := GET_FILTERED_FRIENDSHIP_QUERY(ToUserId(userKey), filter.Query, filter.Category, filter.Letter, filter.SortBy, filter.Companies)

	log.Printf("GetFriendsOf: \n %s \n %s \n", query, params)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFriendships := make([]*model.Friendship, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var friendship model.Friendship
		cursor.ReadDocument(ctx, &friendship)
		allFriendships[i] = &friendship
	}

	return allFriendships, nil

}

func (r *NetworkRepo) RemoveFriendship(ctx context.Context, user1Key, user2Key string) error {

	_, err := r.db.Query(ctx, REMOVE_UNIDIRECTIONAL_RELATION_QUERY, map[string]interface{}{
		"user1":     ToUserId(user1Key),
		"user2":     ToUserId(user2Key),
		"@relation": FriendshipName,
	})
	if err != nil {
		return err
	}

	_, err = r.db.Query(ctx, REMOVE_UNIDIRECTIONAL_RELATION_QUERY, map[string]interface{}{
		"user1":     ToUserId(user1Key),
		"user2":     ToUserId(user2Key),
		"@relation": CategoriesName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *NetworkRepo) IsFriend(ctx context.Context, user1Key, user2Key string) (bool, error) {

	return r.executeBoolQuery(IS_FRIEND_QUERY, map[string]interface{}{
		"user1": ToUserId(user1Key),
		"user2": ToUserId(user2Key),
	})
}

func (r *NetworkRepo) Follow(ctx context.Context, follow *model.FollowRequest) (string, error) {

	meta, err := r.follows.CreateDocument(ctx, follow)
	if err != nil {
		return "", err
	}
	return meta.Key, nil
}

func (r *NetworkRepo) Unfollow(ctx context.Context, user1Id, user2Id string) error {

	_, err := r.db.Query(ctx, REMOVE_DIRECTIONAL_RELATION_QUERY, map[string]interface{}{
		"from":      user1Id,
		"to":        user2Id,
		"@relation": FollowName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) GetFollowers(ctx context.Context, userKey string, filter *model.UserFilter) ([]*model.Follow, error) {
	query, params := GET_FILTERED_USER_FOLLOWERS(userKey, filter)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFollows := make([]*model.Follow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var follow model.Follow
		cursor.ReadDocument(ctx, &follow)
		allFollows[i] = &follow
	}

	return allFollows, nil
}

func (r *NetworkRepo) GetFollowings(ctx context.Context, userKey string, filter *model.UserFilter) ([]*model.Follow, error) {
	query, params := GET_FILTERED_USER_FOLLOWINGS(userKey, filter)

	log.Println("GetFollowings", "\n", query, "\n", params, "----")

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFollows := make([]*model.Follow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var friendship model.Follow
		cursor.ReadDocument(ctx, &friendship)
		allFollows[i] = &friendship
	}

	return allFollows, nil
}

func (r *NetworkRepo) AddToCategory(ctx context.Context, categoryRelation *model.CategoryRelation) (string, error) {

	meta, err := r.categories.CreateDocument(ctx, categoryRelation)
	if err != nil {
		return "", err
	}
	return meta.Key, nil
}

func (r *NetworkRepo) AddToFollowingsCategory(ctx context.Context, categoryRelation *model.CategoryRelation) (string, error) {

	meta, err := r.categoriesForFollowings.CreateDocument(ctx, categoryRelation)
	if err != nil {
		return "", err
	}
	return meta.Key, nil
}

func (r *NetworkRepo) RemoveFromCategory(ctx context.Context, categoryRelation *model.CategoryRelation) error {

	_, err := r.db.Query(ctx, REMOVE_FROM_CATEGORY_QUERY, map[string]interface{}{
		"from":     categoryRelation.OwnerId,
		"to":       categoryRelation.ReferralId,
		"category": categoryRelation.CategoryName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) RemoveFromFollowingsCategory(ctx context.Context, categoryRelation *model.CategoryRelation) error {

	_, err := r.db.Query(ctx, REMOVE_FROM_FOLLOWINGS_CATEGORY_QUERY, map[string]interface{}{
		"from":     categoryRelation.OwnerId,
		"to":       categoryRelation.ReferralId,
		"category": categoryRelation.CategoryName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) BatchRemoveFromCategory(ctx context.Context, userKey string, userIds []string, categoryName string, all bool) error {
	query, params := BATCH_REMOVE_FROM_CATEGORY_QUERY(userKey, userIds, categoryName, all)

	_, err := r.db.Query(ctx, query, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) BatchRemoveFromCategoryForFollowings(ctx context.Context, ownerKey string, referalKeys []string, categoryName string, all bool) error {
	query, params := BATCH_REMOVE_FROM_CATEGORY_FOR_FOLLOWINGS_QUERY(ownerKey, referalKeys, categoryName, all)

	_, err := r.db.Query(ctx, query, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) GetFriendSuggestions(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.UserSuggestion, error) {
	query, params := GET_FRIEND_SUGGESTIONS(userId, pagination)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allSuggestions := make([]*model.UserSuggestion, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var suggestion model.UserSuggestion
		cursor.ReadDocument(ctx, &suggestion)
		allSuggestions[i] = &suggestion
	}
	return allSuggestions, nil
}

func (r *NetworkRepo) MakeCompanyOwner(ctx context.Context, request *model.UserCompanyId) error {
	_, err := r.owns_company.CreateDocument(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) IsCompanyOwner(ctx context.Context, request *model.UserCompanyId) (bool, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_COMPANY_OWNERSHIP, map[string]interface{}{
		"userId":    request.UserId,
		"companyId": request.CompanyId,
	})
	if err != nil {
		return false, err
	}
	defer cursor.Close()

	return cursor.Count() > 0, nil
}

func (r *NetworkRepo) MakeCompanyAdmin(ctx context.Context, request *model.AdminEdge) error {
	_, err := r.admins.CreateDocument(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) GetAdminObject(ctx context.Context, request *model.UserCompanyId) (*model.Admin, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_COMPANY_ADMIN, map[string]interface{}{
		"id":        request.UserId,
		"companyId": request.CompanyId,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	if cursor.Count() > 0 {
		var admin model.Admin
		cursor.ReadDocument(ctx, &admin)
		return &admin, nil
	}

	return nil, hc_errors.JsonError{Type: hc_errors.NOT_FOUNT_ERROR_TYPE, Description: "User is not admin of the company"}
}

func (r *NetworkRepo) GetCompanyAdmins(ctx context.Context, companyKey string) ([]*model.Admin, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_ALL_COMPANY_ADMIN, map[string]interface{}{
		"id": ToCompanyId(companyKey),
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allAdmins := make([]*model.Admin, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var admin model.Admin
		cursor.ReadDocument(ctx, &admin)
		allAdmins[i] = &admin
	}
	return allAdmins, nil
}

func (r *NetworkRepo) ChangeAdminLevel(ctx context.Context, adminKey string, level model.AdminLevel) error {
	_, err := r.admins.UpdateDocument(ctx, adminKey, map[string]interface{}{
		"level": level,
	})
	return err
}

func (r *NetworkRepo) DeleteCompanyAdmin(ctx context.Context, adminKey string) error {
	_, err := r.admins.RemoveDocument(ctx, adminKey)
	return err
}

func (r *NetworkRepo) GetUserCompanies(ctx context.Context, userKey string) ([]*model.Admin, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_ALL_ADMINED_COMPANIES, map[string]interface{}{
		"id": ToUserId(userKey),
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allComps := make([]*model.Admin, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var admin model.Admin
		cursor.ReadDocument(ctx, &admin)
		allComps[i] = &admin
	}
	return allComps, nil
}

func (r *NetworkRepo) GetFollowerCompanies(ctx context.Context, userKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error) {
	query, params := GET_FILTERED_COMPANY_FOLLOWERS(userKey, filter)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFollows := make([]*model.CompanyFollow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var follow model.CompanyFollow
		cursor.ReadDocument(ctx, &follow)
		allFollows[i] = &follow
	}

	return allFollows, nil
}

func (r *NetworkRepo) GetFilteredFollowingCompanies(ctx context.Context, userKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error) {
	query, params := GET_FILTERED_COMPANY_FOLLOWINGS(userKey, filter)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allCompanies := make([]*model.CompanyFollow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var follow model.CompanyFollow
		cursor.ReadDocument(ctx, &follow)
		allCompanies[i] = &follow
	}

	return allCompanies, nil
}

func (r *NetworkRepo) SaveRecommendationRequest(ctx context.Context, request *model.RecommendationRequestModel) (string, error) {
	meta, err := r.recommendationRequest.CreateDocument(ctx, request)
	return meta.Key, err
}

func (r *NetworkRepo) GetRecommendationRequestModel(ctx context.Context, key string) (*model.RecommendationRequestModel, error) {
	var req model.RecommendationRequestModel
	_, err := r.recommendationRequest.ReadDocument(ctx, key, &req)
	return &req, err
}

func (r *NetworkRepo) UpdateRecommendationRequest(ctx context.Context, key string, request *model.RecommendationRequestModel) error {
	_, err := r.recommendationRequest.UpdateDocument(ctx, key, request)
	return err
}

func (r *NetworkRepo) GetRequestedRecommendations(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.RecommendationRequest, error) {
	query, params := GET_REQUESTED_RECOMMENDATIONS(userKey, pagination)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	recommendationRequests := make([]*model.RecommendationRequest, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var request model.RecommendationRequest
		cursor.ReadDocument(ctx, &request)
		recommendationRequests[i] = &request
	}

	return recommendationRequests, nil
}

func (r *NetworkRepo) GetReceivedRecommendationRequests(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.RecommendationRequest, error) {
	query, params := GET_RECEIVED_RECOMMENDATION_REQUESTS(userKey, pagination)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	recommendationRequests := make([]*model.RecommendationRequest, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var request model.RecommendationRequest
		cursor.ReadDocument(ctx, &request)
		recommendationRequests[i] = &request
	}

	return recommendationRequests, nil
}

func (r *NetworkRepo) WriteRecommendation(ctx context.Context, recommendation *model.RecommendationModel) error {
	_, err := r.recommendations.CreateDocument(ctx, recommendation)
	return err
}

func (r *NetworkRepo) GetRecommendationModel(ctx context.Context, key string) (*model.RecommendationModel, error) {
	var rec model.RecommendationModel
	_, err := r.recommendations.ReadDocument(ctx, key, &rec)
	return &rec, err
}
func (r *NetworkRepo) SetRecommendationVisibility(ctx context.Context, key string, visible bool) error {
	_, err := r.recommendations.UpdateDocument(ctx, key, map[string]interface{}{
		"hidden": !visible,
	})
	return err
}

func (r *NetworkRepo) GetReceivedRecommendations(ctx context.Context, userKey string, pagination *model.Pagination /*, isMe bool*/) ([]*model.Recommendation, error) {
	query, params := GET_RECEIVED_RECOMMENDATIONS(userKey, pagination /*, isMe*/)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	recommendations := make([]*model.Recommendation, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var recommendation model.Recommendation
		cursor.ReadDocument(ctx, &recommendation)
		recommendations[i] = &recommendation
		// log.Printf("RECEIVED user: %q\n%+v\n", userKey, recommendation)
	}

	return recommendations, nil
}

func (r *NetworkRepo) GetGivenRecommendations(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.Recommendation, error) {
	query, params := GET_GIVEN_RECOMMENDATIONS(userKey, pagination)

	// log.Println("GET_GIVEN_RECOMMENDATIONS:", query, "\n", params)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	recommendations := make([]*model.Recommendation, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var recommendation model.Recommendation
		cursor.ReadDocument(ctx, &recommendation)
		recommendations[i] = &recommendation
		// log.Printf("GIVEN user: %q\n%+v\n", userKey, recommendation)
	}

	return recommendations, nil
}

func (r *NetworkRepo) GetHiddenRecommendations(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.Recommendation, error) {
	query, params := GET_HIDDEN_RECOMMENDATIONS(userKey, pagination)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	recommendations := make([]*model.Recommendation, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var recommendation model.Recommendation
		cursor.ReadDocument(ctx, &recommendation)
		recommendations[i] = &recommendation
	}

	return recommendations, nil
}

func (r *NetworkRepo) AddExperience(ctx context.Context, experience *model.AddExperienceRequest) error {
	_, err := r.works_at.CreateDocument(ctx, experience)
	return err
}

func (r *NetworkRepo) GetAllExperience(ctx context.Context, userKey string) ([]*model.Experience, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_ALL_EXPERIENCE, map[string]interface{}{
		"id": ToUserId(userKey),
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allExp := make([]*model.Experience, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var exp model.Experience
		cursor.ReadDocument(ctx, &exp)
		allExp[i] = &exp
	}

	return allExp, nil
}

func (r *NetworkRepo) GetSuggestedCompanies(ctx context.Context, userKey string, pagination *model.Pagination) ([]*model.CompanySuggestion, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_SUGGESTED_COMPANIES(pagination), map[string]interface{}{
		"id": ToUserId(userKey),
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	companies := make([]*model.CompanySuggestion, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var exp model.CompanySuggestion
		cursor.ReadDocument(ctx, &exp)
		companies[i] = &exp
	}

	return companies, nil
}

func (r *NetworkRepo) BatchRemoveFromCategoryForFollowingsForCompany(ctx context.Context, ownerKey string, referalKeys []string, categoryName string, all bool) error {
	query, params := BATCH_REMOVE_FROM_CATEGORY_FOR_FOLLOWINGS_QUERY_FOR_COMPANY(ownerKey, referalKeys, categoryName, all)
	_, err := r.db.Query(ctx, query, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) Block(ctx context.Context, blockerId, blockedId string) error {
	// _, err := r.blocks.CreateDocument(ctx, model.BlockRequest{
	// 	Blocker:   blockerId,
	// 	Blocked:   blockedId,
	// 	CreatedAt: time.Now(),
	// })
	// return err

	cursor, err := r.db.Query(
		driver.WithQueryCount(ctx), `
		INSERT {
    "_from": @blocker,
    "_to": @blocked,
    "created_at":@time
}
INTO blocks
		`, map[string]interface{}{
			"blocker": blockerId,
			"blocked": blockedId,
			"time":    time.Now(),
		})
	if err != nil {
		return err
	}
	defer cursor.Close()

	return nil
}

func (r *NetworkRepo) Unblock(ctx context.Context, blockerId, blockedId string) error {
	_, err := r.db.Query(ctx, REMOVE_DIRECTIONAL_RELATION_QUERY, map[string]interface{}{
		"from":      blockerId,
		"to":        blockedId,
		"@relation": BlocksName,
	})
	return err
}

func (r *NetworkRepo) GetBlockedUsersOrCompanies(ctx context.Context, userKey string) ([]*model.BlockedUserOrCompany, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_BLOCKED_USERS_OR_COMPANIES, map[string]interface{}{
		"id": ToUserId(userKey),
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	userOrCompanies := make([]*model.BlockedUserOrCompany, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var userOrCompany model.BlockedUserOrCompany
		cursor.ReadDocument(ctx, &userOrCompany)
		userOrCompanies[i] = &userOrCompany
	}

	return userOrCompanies, nil
}

func (r *NetworkRepo) GetBlockedUsersForCompany(ctx context.Context, companyID string) ([]*model.User, error) {
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_BLOCKED_USERS, map[string]interface{}{
		"id": ToCompanyId(companyID),
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	users := make([]*model.User, 0, cursor.Count())
	for cursor.HasMore() {
		user := model.User{}

		_, err := cursor.ReadDocument(ctx, &user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

//func (r *NetworkRepo) GetBlockedUsers(ctx context.Context, userKey string) ([]*model.User, error) {
//	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_BLOCKED_USERS, map[string]interface{}{
//		"id": ToUserId(userKey),
//	})
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close()
//
//	users := make([]*model.User, cursor.Count())
//	for i := 0; cursor.HasMore(); i++ {
//		var user model.User
//		cursor.ReadDocument(ctx, &user)
//		users[i] = &user
//	}
//
//	return users, nil
//}
//
//func (r *NetworkRepo) GetBlockedCompanies(ctx context.Context, userKey string) ([]*model.Company, error) {
//	cursor, err := r.db.Query(driver.WithQueryCount(ctx), GET_BLOCKED_COMPANIES, map[string]interface{}{
//		"id": ToUserId(userKey),
//	})
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close()
//
//	companies := make([]*model.Company, cursor.Count())
//	for i := 0; cursor.HasMore(); i++ {
//		var company model.Company
//		cursor.ReadDocument(ctx, &company)
//		companies[i] = &company
//	}
//
//	return companies, nil
//}

func (r *NetworkRepo) GetFollowingsForCompany(ctx context.Context, companyKey string, filter *model.UserFilter) ([]*model.Follow, error) {
	query, params := GET_FILTERED_USER_FOLLOWINGS_FOR_COMPANY(companyKey, filter)

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFollows := make([]*model.Follow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var follow model.Follow
		cursor.ReadDocument(ctx, &follow)
		allFollows[i] = &follow
	}

	return allFollows, nil
}

func (r *NetworkRepo) GetFollowersForCompany(ctx context.Context, companyKey string, filter *model.UserFilter) ([]*model.Follow, error) {
	query, params := GET_FILTERED_USER_FOLLOWERS_FOR_COMPANY(companyKey, filter)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFollows := make([]*model.Follow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var follow model.Follow
		cursor.ReadDocument(ctx, &follow)
		allFollows[i] = &follow
	}

	return allFollows, nil
}

func (r *NetworkRepo) GetFollowingCompaniesForCompany(ctx context.Context, companyKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error) {
	query, params := GET_FILTERED_COMPANY_FOLLOWINGS_FOR_COMPANY(companyKey, filter)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	companyFollows := make([]*model.CompanyFollow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var follow model.CompanyFollow
		cursor.ReadDocument(ctx, &follow)
		companyFollows[i] = &follow
	}

	return companyFollows, nil
}

func (r *NetworkRepo) GetFollowerCompaniesForCompany(ctx context.Context, companyKey string, filter *model.CompanyFilter) ([]*model.CompanyFollow, error) {
	query, params := GET_FILTERED_COMPANY_FOLLOWERS_FOR_COMPANY(companyKey, filter)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	companyFollows := make([]*model.CompanyFollow, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var follow model.CompanyFollow
		cursor.ReadDocument(ctx, &follow)
		companyFollows[i] = &follow
	}

	return companyFollows, nil
}

func (r *NetworkRepo) GetSuggestedPeopleForCompany(ctx context.Context, companyKey string) ([]*model.UserSuggestion, error) {
	query, params := GET_SUGGESTED_PEOPLE_FOR_COMPANY(companyKey, 20)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	users := make([]*model.UserSuggestion, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var suggestion model.UserSuggestion
		cursor.ReadDocument(ctx, &suggestion)
		users[i] = &suggestion
	}

	return users, nil
}

func (r *NetworkRepo) GetSuggestedCompaniesForCompany(ctx context.Context, companyKey string, pagination *model.Pagination) ([]*model.CompanySuggestion, error) {
	query, params := GET_SUGGESTED_COMPANIES_FOR_COMPANY(companyKey, pagination)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	companies := make([]*model.CompanySuggestion, cursor.Count())
	for i := 0; cursor.HasMore(); i++ {
		var suggestion model.CompanySuggestion
		cursor.ReadDocument(ctx, &suggestion)
		companies[i] = &suggestion
	}

	return companies, nil
}

func (r *NetworkRepo) GetNumberOfFollowers(ctx context.Context, id string) (int, error) {
	query, params := GET_NUMBER_OF_FOLLOWERS(id)
	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	var number int
	_, err = cursor.ReadDocument(ctx, &number)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (r *NetworkRepo) IsBlocked(ctx context.Context, id1, id2 string) (bool, error) {
	query, params := IS_BLOCKED(id1, id2)

	log.Printf("IsBlocked: \n %v\n%v\n", query, params)

	return r.executeBoolQuery(query, params)
}

func (r *NetworkRepo) IsBlockedByUser(ctx context.Context, id1, id2 string) (bool, error) {
	query, params := IS_BLOCKED(id2, id1)

	log.Printf("IsBlockedByCompany:\n%v\n%v\n", query, params)

	res, err := r.executeBoolQuery(query, params)
	if err != nil {
		return false, err
	}

	return res, nil
}

// func (r *NetworkRepo) IsBlockedCompany(ctx context.Context, id1, id2 string) (bool, error) {
// 	query, params := IS_BLOCKED_COMPANY(id1, id2)
// 	return r.executeBoolQuery(query, params)
// }

func (r *NetworkRepo) IsFollowing(ctx context.Context, id1, id2 string) (bool, error) {
	query, params := IS_FOLLOWING(id1, id2)

	log.Printf("IsFollowingForCompany: \n%v\n%v\n", query, params)

	return r.executeBoolQuery(query, params)
}

func (r *NetworkRepo) IsFavourite(ctx context.Context, id1, id2 string) (bool, error) {
	query, params := IS_FAVOURITE(id1, id2)
	return r.executeBoolQuery(query, params)
}

func (r *NetworkRepo) GetUserCountings(ctx context.Context, userKey string) (*model.UserCountings, error) {
	query, params := GET_USER_COUNTINGS(ToUserId(userKey))

	log.Println("GET_USER_COUNTINGS", "\n", query, "\n", params)

	cursor, err := r.db.Query(ctx, query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var countings model.UserCountings
	_, err = cursor.ReadDocument(ctx, &countings)
	if err != nil {
		return nil, err
	}
	return &countings, nil
}

func (r *NetworkRepo) SetCategoryTree(ctx context.Context, tree *model.CategoryTree) error {
	exists, err := r.category_tree.DocumentExists(ctx, tree.OwnerId)
	if err != nil {
		return err
	}
	if exists {
		_, err = r.category_tree.UpdateDocument(ctx, tree.OwnerId, tree)
		return err
	} else {
		_, err = r.category_tree.CreateDocument(ctx, tree)
		return err
	}
}

func (r *NetworkRepo) GetCategoryTree(ctx context.Context, userKey string) (*model.CategoryTree, error) {
	var tree model.CategoryTree
	_, err := r.category_tree.ReadDocument(ctx, userKey, &tree)

	if driver.IsNotFound(err) {
		tree = model.CategoryTree{
			OwnerId: userKey,
			Categories: []model.CategoryItem{
				{
					UniqueName:  "favorite",
					Name:        "Favorite",
					HasChildren: false,
				}, {
					UniqueName:  "friends_and_family",
					Name:        "Friends & Family",
					HasChildren: true,
					Children:    []model.CategoryItem{},
				}, {
					UniqueName:  "work",
					Name:        "Work",
					HasChildren: true,
					Children:    []model.CategoryItem{},
				}, {
					UniqueName:  "business",
					Name:        "Business",
					HasChildren: true,
					Children:    []model.CategoryItem{},
				}, {
					UniqueName:  "other",
					Name:        "Other",
					HasChildren: true,
					Children:    []model.CategoryItem{},
				},
			},
		}
		_, err = r.category_tree.CreateDocument(ctx, tree)
	}
	return &tree, err
}

func (r *NetworkRepo) SetCategoryTreeForFollowings(ctx context.Context, tree *model.CategoryTree) error {
	exists, err := r.category_tree_followings.DocumentExists(ctx, tree.OwnerId)
	if err != nil {
		return err
	}
	if exists {
		_, err = r.category_tree_followings.UpdateDocument(ctx, tree.OwnerId, tree)
		return err
	} else {
		_, err = r.category_tree_followings.CreateDocument(ctx, tree)
		return err
	}
}

func (r *NetworkRepo) GetCategoryTreeForFollowings(ctx context.Context, key string) (*model.CategoryTree, error) {
	var tree model.CategoryTree
	_, err := r.category_tree_followings.ReadDocument(ctx, key, &tree)

	if driver.IsNotFound(err) {
		tree = model.CategoryTree{
			OwnerId: key,
			Categories: []model.CategoryItem{
				{
					UniqueName:  "favorite",
					Name:        "Favorite",
					HasChildren: false,
				}, {
					UniqueName:  "business",
					Name:        "Business",
					HasChildren: true,
					Children:    []model.CategoryItem{},
				}, {
					UniqueName:  "other",
					Name:        "Other",
					HasChildren: true,
					Children:    []model.CategoryItem{},
				},
			},
		}
		_, err = r.category_tree_followings.CreateDocument(ctx, tree)
	}
	return &tree, err
}

func (r *NetworkRepo) ClearCategory(ctx context.Context, userKey, category string) error {
	query, params := CLEAR_CATEGORIES(userKey, category)
	_, err := r.db.Query(ctx, query, params)
	return err
}

func (r *NetworkRepo) ClearCategoryForFollowings(ctx context.Context, userKey, category string) error {
	query, params := CLEAR_CATEGORIES_FOR_FOLLOWINGS(ToUserId(userKey), category)
	_, err := r.db.Query(ctx, query, params)
	return err
}

func (r *NetworkRepo) ClearCategoryForFollowingsForCompany(ctx context.Context, companyKey, category string) error {
	query, params := CLEAR_CATEGORIES_FOR_FOLLOWINGS(ToCompanyId(companyKey), category)
	_, err := r.db.Query(ctx, query, params)
	return err
}

func (r *NetworkRepo) GetFriendsOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users
		LET amount = (
    FOR v, e IN 1..1 ANY DOCUMENT("users", @target_id) friendship
        FILTER v.status == "ACTIVATED"
        FILTER e.status == "Approved"
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 ANY DOCUMENT("users", @target_id) friendship
    FILTER v.status == "ACTIVATED"
    FILTER e.status == "Approved"
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": userID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowsOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 OUTBOUND DOCUMENT("users", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("users", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 OUTBOUND DOCUMENT("users", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("users", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": userID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowersOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 INBOUND DOCUMENT("users", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("users", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 INBOUND DOCUMENT("users", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("users", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": userID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowsCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 OUTBOUND DOCUMENT("users", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("companies", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 OUTBOUND DOCUMENT("users", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("companies", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": userID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowersCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 INBOUND DOCUMENT("users", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("companies", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 INBOUND DOCUMENT("users", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("companies", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": userID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowsOfCompany(ctx context.Context, companyID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 OUTBOUND DOCUMENT("companies", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("users", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 OUTBOUND DOCUMENT("companies", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("users", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": companyID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowersOfCompany(ctx context.Context, companyID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 INBOUND DOCUMENT("companies", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("users", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 INBOUND DOCUMENT("companies", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("users", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": companyID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowsCompaniesOfCompany(ctx context.Context, companyID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 OUTBOUND DOCUMENT("companies", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("companies", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 OUTBOUND DOCUMENT("companies", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("companies", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": companyID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetFollowersCompaniesOfCompany(ctx context.Context, companyID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx,
		`WITH users, companies
		LET amount = (
    FOR v, e IN 1..1 INBOUND DOCUMENT("companies", @target_id) follow
        FILTER v.status == "ACTIVATED"
				FILTER IS_SAME_COLLECTION("companies", v._id)
				COLLECT WITH COUNT INTO cnt
				RETURN cnt
    )

LET ids = (FOR v, e IN 1..1 INBOUND DOCUMENT("companies", @target_id) follow
    FILTER v.status == "ACTIVATED"
		FILTER IS_SAME_COLLECTION("companies", v._id)
    LIMIT @after, @first
    RETURN v._key)

RETURN {"ids":ids, "amount": FIRST(amount)}`,
		map[string]interface{}{
			"target_id": companyID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetMutualFriendsOfUser(ctx context.Context, senderID string, userID string, first uint32, after uint32) ([]string, int64, error) {
	cursor, err := r.db.Query(ctx, `
		WITH users
		LET my_friends = (
		FOR v, e IN 1..1 ANY DOCUMENT("users", @my_id) friendship
				FILTER v.status == "ACTIVATED"
				FILTER e.status == "Approved"
				RETURN v._key
		)

LET amount = (
FOR v, e IN 1..1 ANY DOCUMENT("users", @target_id) friendship
		FILTER v.status == "ACTIVATED"
		FILTER e.status == "Approved"
		FILTER e._to == CONCAT("users/", @my_id) OR e._from == CONCAT("users/", @my_id)
		COLLECT WITH COUNT INTO cnt
		RETURN cnt
)


LET friends = (
	FOR v, e IN 1..1 ANY DOCUMENT("users", @target_id) friendship
			FILTER v.status == "ACTIVATED"
			FILTER e.status == "Approved"
			LIMIT  @after, @first
			RETURN v._key
	)

RETURN {
    "ids": INTERSECTION(friends, my_friends),
    "amount": FIRST(amount)
}`,
		map[string]interface{}{
			"my_id":     senderID,
			"target_id": userID,
			"first":     first,
			"after":     after,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()

	result := struct {
		IDs    []string `json:"ids"`
		Amount int64    `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return nil, 0, err
		}
	}

	return result.IDs, result.Amount, nil
}

func (r *NetworkRepo) GetCompanyCountings(ctx context.Context, companyID string) (*model.CompanyCountings, error) {
	query, params := GET_COMPANY_COUNTINGS(ToCompanyId(companyID))

	cursor, err := r.db.Query(ctx, query, params)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var countings model.CompanyCountings
	_, err = cursor.ReadDocument(ctx, &countings)
	if err != nil {
		return nil, err
	}
	return &countings, nil
}

func (r *NetworkRepo) RemoveFriendshipByID(ctx context.Context, key string) error {

	_, err := r.db.Query(ctx, REMOVE_FRIENDSHIP_BY_ID, map[string]interface{}{
		"friendship_id": key,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *NetworkRepo) GetAmountOfMutualFriends(ctx context.Context, senderID string, userID string) (int32, error) {
	cursor, err := r.db.Query(ctx, `
		WITH users

		LET amount = (
		FOR v, e IN 1..1 ANY DOCUMENT("users", @target_id) friendship
				FILTER v.status == "ACTIVATED"
				FILTER e.status == "Approved"
				FILTER COUNT(for b in blocks filter (b._from == @target_id AND b._to == @my_id) OR (b._to == @target_id AND b._from == @my_id) return 1) == 0
				FOR in_v, in_e IN 1..1 ANY DOCUMENT("users", @my_id) friendship
		    		FILTER in_v.status == "ACTIVATED"
		    		FILTER in_e.status == "Approved"
		    		FILTER in_v._key != @my_id
		    		FILTER in_v._key != @target_id
				    FILTER v._key == in_v._key
				    //RETURN {in_v, in_e}
		    		COLLECT WITH COUNT INTO cnt
		    		RETURN cnt
		)

		RETURN {
		    "amount": FIRST(amount)
		}`,
		map[string]interface{}{
			"my_id":     senderID,
			"target_id": userID,
		},
	)
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	result := struct {
		Amount int32 `json:"amount"`
	}{}

	for cursor.HasMore() {
		_, err = cursor.ReadDocument(ctx, &result)
		if err != nil {
			return 0, err
		}
	}

	return result.Amount, nil
}

func (r *NetworkRepo) GetFollowersIDs(ctx context.Context, id string, isCompany bool) ([]string, error) {
	key := ""
	if isCompany {
		key = "companies/" + id
	} else {
		key = "users/" + id
	}

	query, params := GET_ALL_FOLLOWERS_IDS(key)

	log.Println("GetFollowersIDs", "\n", query, "\n", params, "----")

	cursor, err := r.db.Query(driver.WithQueryCount(ctx), query, params)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	allFollows := make([]string, 0, cursor.Count())
	for cursor.HasMore() {

		var id string
		_, err := cursor.ReadDocument(ctx, &id)
		if err != nil {
			return nil, err
		}

		allFollows = append(allFollows, id)
	}

	return allFollows, nil
}
