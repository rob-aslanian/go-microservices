package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/repo/network"
	"golang.org/x/net/context"
)

func (s *NetworkService) SendFriendRequest(ctx context.Context, friendId, description string) (*model.Friendship, error) {
	err := s.validator.ValidateUserId(friendId)
	if err != nil {
		return nil, err
	}
	senderId := s.authenticateUser(ctx)

	if senderId == friendId {
		return nil, errors.New("You can not send friend request to yourself")
	}

	// TODO: check if request from friendId wans't sent
	frindshipRequests, err := s.GetFriendshipRequests(ctx, &model.FriendshipRequestFilter{
		Status: "Requested",
	})
	if err != nil {
		return nil, err
	}

	for _, fr := range frindshipRequests {
		if fr.Friend.Id == friendId {
			return fr, nil
		}
	}

	friendship := &model.NewFriendship{
		SenderId:    network.ToUserId(senderId),
		ReceiverId:  network.ToUserId(friendId),
		Description: description,
		CreatedAt:   time.Now(),
		Status:      model.FriendshipStatus_Requested,
	}
	insertedKey, err := s.repo.InsertFriendship(ctx, friendship)
	if err != nil {
		return nil, err
	}

	fr, err := s.repo.GetFriendship(ctx, insertedKey)
	if err != nil {
		return nil, err
	}

	// send notification
	not := model.NewConnectionRequest{
		UserSenderID: senderId,
		FriendshipID: fr.Id,
	}
	not.GenerateID()

	err = s.mq.SendNewConnection(friendId, &not)
	if err != nil {
		log.Println("error while sending notification:", err)
		// s.tracer
	}

	return fr, nil
}

func (s *NetworkService) ApproveFriendRequest(ctx context.Context, key string) error {
	// validate input
	//err := s.validator.ValidateUserId(key) // This is not user id, this is friend request key
	//if err != nil {
	//	return err
	//}

	senderId := s.authenticateUser(ctx)

	err := s.repo.ChangeFriendshipStatus(ctx, key, model.FriendshipStatus_Approved)
	if err != nil {
		return err
	}

	friendship, err := s.repo.GetFriendship(ctx, key)
	if err != nil {
		return err
	}

	myId := network.ToUserId(senderId)
	friendId := network.ToUserId(friendship.Friend.Id)
	now := time.Now()
	s.repo.Follow(ctx, &model.FollowRequest{FollowerId: myId, FollowingId: friendId, CreatedAt: now})
	s.repo.Follow(ctx, &model.FollowRequest{FollowerId: friendId, FollowingId: myId, CreatedAt: now})

	// send notification
	not := model.NewApproveConnectionRequest{
		UserSenderID: senderId,
	}
	not.GenerateID()

	err = s.mq.ApproveConnectionRequest(friendship.Friend.Id, &not)
	if err != nil {
		log.Println("error while sending notification:", err)
		// s.tracer
	}

	return nil
}

func (s *NetworkService) DenyFriendRequest(ctx context.Context, key string) error {
	// validate input
	//err := s.validator.ValidateUserId(key) // This is not user id, this is friend request key
	//if err != nil {
	//	return err
	//}

	s.authenticateUser(ctx)

	return s.repo.ChangeFriendshipStatus(ctx, key, model.FriendshipStatus_Denied)
}

func (s *NetworkService) IgnoreFriendRequest(ctx context.Context, key string) error {
	// validate input
	//err := s.validator.ValidateUserId(key) // This is not user id, this is friend request key
	//if err != nil {
	//	return err
	//}

	s.authenticateUser(ctx)

	// return s.repo.ChangeFriendshipStatus(ctx, key, model.FriendshipStatus_Ignored)
	return s.repo.RemoveFriendshipByID(ctx, key)
}

func (s *NetworkService) GetFriendshipRequests(ctx context.Context, filter *model.FriendshipRequestFilter) ([]*model.Friendship, error) {
	// validate input
	err := s.validator.ValidateStruct(filter)
	if err != nil {
		return nil, err
	}

	senderId := s.authenticateUser(ctx)

	friendships, err := s.repo.GetFriendshipRequests(ctx, senderId, filter)
	if err != nil {
		return nil, err
	}
	return friendships, nil
}

func (s *NetworkService) GetAllFriendships(ctx context.Context, filter *model.FriendshipFilter) ([]*model.Friendship, error) {
	// validate input
	err := s.validator.ValidateStruct(filter)
	if err != nil {
		return nil, err
	}

	senderId := s.authenticateUser(ctx)

	friendships, err := s.repo.GetFriendsOf(ctx, senderId, filter)
	if err != nil {
		return nil, err
	}
	return friendships, nil
}

func (s *NetworkService) GetAllFriendshipID(ctx context.Context) ([]string, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.GetAllFriendshipID(ctx, senderId)
}

func (s *NetworkService) Unfriend(ctx context.Context, friendId string) error {
	err := s.validator.ValidateUserId(friendId)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)
	err = s.repo.RemoveFriendship(ctx, senderId, friendId)
	return err
}

func (s *NetworkService) Follow(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	_, err = s.repo.Follow(ctx, &model.FollowRequest{
		FollowerId:  network.ToUserId(senderId),
		FollowingId: network.ToUserId(id),
		CreatedAt:   time.Now(),
	})

	// send notification
	not := model.NewFollow{
		UserSenderID: senderId,
	}
	not.GenerateID()

	err = s.mq.SendNewFollow(id, &not)
	if err != nil {
		log.Println("error while sending notification:", err)
		// s.tracer
	}

	return err
}

func (s *NetworkService) IsFriend(ctx context.Context, userKey string) (bool, error) {
	err := s.validator.ValidateKey(userKey)
	if err != nil {
		return false, err
	}

	senderId := s.authenticateUser(ctx)

	return s.repo.IsFriend(ctx, senderId, userKey)
}

func (s *NetworkService) Unfollow(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)
	err = s.repo.Unfollow(ctx, network.ToUserId(senderId), network.ToUserId(id))
	return err
}

func (s *NetworkService) GetFollowers(ctx context.Context, filter *model.UserFilter) ([]*model.Follow, error) {
	senderId := s.authenticateUser(ctx)
	followers, err := s.repo.GetFollowers(ctx, senderId, filter)
	return followers, err
}

func (s *NetworkService) GetFollowings(ctx context.Context, filter *model.UserFilter) ([]*model.Follow, error) {
	senderId := s.authenticateUser(ctx)
	followers, err := s.repo.GetFollowings(ctx, senderId, filter)
	return followers, err
}

func (s *NetworkService) AddToFavourites(ctx context.Context, userId string) error {
	return s.AddToCategory(ctx, userId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) AddToFollowingsFavourites(ctx context.Context, companyId string) error {
	return s.AddToFollowingsCategory(ctx, companyId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) AddToFollowingsFavouritesForCompany(ctx context.Context, companyId, refCompanyId string) error {
	return s.AddToFollowingsCategoryForCompany(ctx, companyId, refCompanyId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) RemoveFromFavourites(ctx context.Context, userId string) error {
	return s.RemoveFromCategory(ctx, userId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) RemoveFromFollowingsFavourites(ctx context.Context, companyId string) error {
	return s.RemoveFromFollowingsCategory(ctx, companyId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) RemoveFromFollowingsFavouritesForCompany(ctx context.Context, companyId, refCompanyId string) error {
	return s.RemoveFromFollowingsCategoryForCompany(ctx, companyId, refCompanyId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) AddToCategory(ctx context.Context, userId, categoryName string) error {
	err := s.validator.ValidateKey(userId)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	cat := &model.CategoryRelation{
		OwnerId:      network.ToUserId(senderId),
		ReferralId:   network.ToUserId(userId),
		CategoryName: categoryName,
		CreatedAt:    time.Now(),
	}

	_, err = s.repo.AddToCategory(ctx, cat)
	return err
}

func (s *NetworkService) AddToFollowingsCategory(ctx context.Context, companyId, categoryName string) error {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	cat := &model.CategoryRelation{
		OwnerId:      network.ToUserId(senderId),
		ReferralId:   network.ToCompanyId(companyId),
		CategoryName: categoryName,
		CreatedAt:    time.Now(),
	}

	_, err = s.repo.AddToFollowingsCategory(ctx, cat)
	return err
}

func (s *NetworkService) AddToFollowingsCategoryForCompany(ctx context.Context, companyId, refCompanyId, categoryName string) error {
	err := s.validator.ValidateKey(refCompanyId)
	if err != nil {
		return err
	}

	s.requireAdminLevelForCompany(ctx, companyId, "Admin")

	cat := &model.CategoryRelation{
		OwnerId:      network.ToCompanyId(companyId),
		ReferralId:   network.ToCompanyId(refCompanyId),
		CategoryName: categoryName,
		CreatedAt:    time.Now(),
	}

	_, err = s.repo.AddToFollowingsCategory(ctx, cat)
	return err
}

func (s *NetworkService) RemoveFromCategory(ctx context.Context, userId, categoryName string) error {
	err := s.validator.ValidateKey(userId)
	if err != nil {
		return err
	}
	senderId := s.authenticateUser(ctx)

	cat := &model.CategoryRelation{
		OwnerId:      network.ToUserId(senderId),
		ReferralId:   network.ToUserId(userId),
		CategoryName: categoryName,
	}

	err = s.repo.RemoveFromCategory(ctx, cat)
	return err
}

func (s *NetworkService) RemoveFromFollowingsCategory(ctx context.Context, companyId, categoryName string) error {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return err
	}
	senderId := s.authenticateUser(ctx)

	cat := &model.CategoryRelation{
		OwnerId:      network.ToUserId(senderId),
		ReferralId:   network.ToCompanyId(companyId),
		CategoryName: categoryName,
	}

	err = s.repo.RemoveFromFollowingsCategory(ctx, cat)
	return err
}

func (s *NetworkService) RemoveFromFollowingsCategoryForCompany(ctx context.Context, companyId, refCompanyId, categoryName string) error {
	err := s.validator.ValidateKey(refCompanyId)
	if err != nil {
		return err
	}
	s.requireAdminLevelForCompany(ctx, companyId, "Admin")

	cat := &model.CategoryRelation{
		OwnerId:      network.ToCompanyId(companyId),
		ReferralId:   network.ToCompanyId(refCompanyId),
		CategoryName: categoryName,
	}

	err = s.repo.RemoveFromFollowingsCategory(ctx, cat)
	return err
}

func (s *NetworkService) BatchRemoveFromCategory(ctx context.Context, userIds []string, categoryName string, all bool) error {
	var err error
	for i, id := range userIds {
		err = s.validator.ValidateKey(id)
		if err != nil {
			return err
		}
		userIds[i] = network.ToUserId(id)
	}
	senderId := s.authenticateUser(ctx)

	return s.repo.BatchRemoveFromCategory(ctx, senderId, userIds, categoryName, all)
}

func (s *NetworkService) BatchRemoveFromFollowingsCategory(ctx context.Context, companyIds []string, categoryName string, all bool) error {
	var err error
	for i, id := range companyIds {
		err = s.validator.ValidateKey(id)
		if err != nil {
			return err
		}
		companyIds[i] = network.ToCompanyId(id)
	}
	senderId := s.authenticateUser(ctx)

	return s.repo.BatchRemoveFromCategoryForFollowings(ctx, senderId, companyIds, categoryName, all)
}

func (s *NetworkService) BatchRemoveFromFollowingsCategoryForCompany(ctx context.Context, companyId string, companyIds []string, categoryName string, all bool) error {
	var err error
	for _, id := range companyIds {
		err = s.validator.ValidateKey(id)
		if err != nil {
			return err
		}
	}
	s.requireAdminLevelForCompany(ctx, companyId, "Admin")

	return s.repo.BatchRemoveFromCategoryForFollowingsForCompany(ctx, companyId, companyIds, categoryName, all)
}

func (s *NetworkService) GetFriendSuggestions(ctx context.Context, pagination *model.Pagination) ([]*model.UserSuggestion, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetFriendSuggestions(ctx, senderId, pagination)
}

func (s *NetworkService) MakeCompanyOwner(ctx context.Context, request *model.UserCompanyId) error {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return err
	}

	request.CompanyId = network.ToCompanyId(request.CompanyId)
	request.UserId = network.ToUserId(request.UserId)

	return s.repo.MakeCompanyOwner(ctx, request)
}

func (s *NetworkService) IsCompanyOwner(ctx context.Context, request *model.UserCompanyId) (bool, error) {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return false, err
	}

	request.CompanyId = network.ToCompanyId(request.CompanyId)
	request.UserId = network.ToUserId(request.UserId)

	return s.repo.IsCompanyOwner(ctx, request)
}

func (s *NetworkService) MakeCompanyAdmin(ctx context.Context, request *model.AdminEdge) error {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return err
	}
	senderId := s.authenticateUser(ctx)

	isOwner, err := s.repo.IsCompanyOwner(ctx, &model.UserCompanyId{UserId: network.ToUserId(senderId), CompanyId: network.ToCompanyId(request.CompanyId)})
	if err != nil {
		return err
	}
	if !isOwner {
		admin, err := s.repo.GetAdminObject(ctx, &model.UserCompanyId{UserId: network.ToUserId(senderId), CompanyId: network.ToCompanyId(request.CompanyId)})
		if err != nil {
			return err
		}
		if admin.Level != model.AdminLevel_Admin {
			return errors.New("You don't have access to this operation")
		}
	}

	request.CompanyId = network.ToCompanyId(request.CompanyId)
	request.UserId = network.ToUserId(request.UserId)
	request.CreatedBy = network.ToUserId(senderId)
	request.CreatedAt = time.Now()

	return s.repo.MakeCompanyAdmin(ctx, request)
}

func (s *NetworkService) GetAdminObject(ctx context.Context, companyId string) (*model.Admin, error) {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return nil, err
	}

	senderId := s.authenticateUser(ctx)

	return s.repo.GetAdminObject(ctx, &model.UserCompanyId{UserId: network.ToUserId(senderId), CompanyId: network.ToCompanyId(companyId)})
}

func (s *NetworkService) GetCompanyAdmins(ctx context.Context, companyId string) ([]*model.Admin, error) {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return nil, err
	}

	return s.repo.GetCompanyAdmins(ctx, companyId)
}

func (s *NetworkService) ChangeAdminLevel(ctx context.Context, request *model.AdminEdge) error {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return err
	}
	senderId := s.authenticateUser(ctx)

	isOwner, err := s.repo.IsCompanyOwner(ctx, &model.UserCompanyId{UserId: network.ToUserId(senderId), CompanyId: network.ToCompanyId(request.CompanyId)})
	if err != nil {
		return err
	}
	if !isOwner {
		admin, err := s.repo.GetAdminObject(ctx, &model.UserCompanyId{UserId: network.ToUserId(senderId), CompanyId: network.ToCompanyId(request.CompanyId)})
		if err != nil {
			return err
		}
		if admin.Level != model.AdminLevel_Admin {
			return errors.New("You don't have access to this operation")
		}
	}

	userAdmin, err := s.repo.GetAdminObject(ctx, &model.UserCompanyId{UserId: network.ToUserId(request.UserId), CompanyId: network.ToCompanyId(request.CompanyId)})
	if err != nil {
		return err
	}

	return s.repo.ChangeAdminLevel(ctx, userAdmin.Id, request.Level)
}

func (s *NetworkService) DeleteCompanyAdmin(ctx context.Context, request *model.UserCompanyId) error {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return err
	}
	senderId := s.authenticateUser(ctx)

	isOwner, err := s.repo.IsCompanyOwner(ctx, &model.UserCompanyId{UserId: network.ToUserId(senderId), CompanyId: network.ToCompanyId(request.CompanyId)})
	if err != nil {
		return err
	}
	if !isOwner {
		admin, err := s.repo.GetAdminObject(ctx, &model.UserCompanyId{UserId: network.ToUserId(senderId), CompanyId: network.ToCompanyId(request.CompanyId)})
		if err != nil {
			return err
		}
		if admin.Level != model.AdminLevel_Admin {
			return errors.New("You don't have access to this operation")
		}
	}

	userAdmin, err := s.repo.GetAdminObject(ctx, &model.UserCompanyId{UserId: network.ToUserId(request.UserId), CompanyId: network.ToCompanyId(request.CompanyId)})
	if err != nil {
		return err
	}

	return s.repo.DeleteCompanyAdmin(ctx, userAdmin.Id)
}

func (s *NetworkService) GetUserCompanies(ctx context.Context, userId string) ([]*model.Admin, error) {
	err := s.validator.ValidateKey(userId)
	if err != nil {
		return nil, err
	}

	return s.repo.GetUserCompanies(ctx, userId)
}

func (s *NetworkService) FollowCompany(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	_, err = s.repo.Follow(ctx, &model.FollowRequest{
		FollowerId:  network.ToUserId(senderId),
		FollowingId: network.ToCompanyId(id),
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return err
	}

	// send notification
	not := model.NewFollow{
		UserSenderID: senderId,
	}
	not.GenerateID()

	err = s.mq.SendNewFollow(id, &not)
	if err != nil {
		log.Println("error while sending notification:", err)
		// s.tracer
	}

	return nil
}

func (s *NetworkService) UnfollowCompany(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)
	err = s.repo.Unfollow(ctx, network.ToUserId(senderId), network.ToCompanyId(id))
	return err
}

func (s *NetworkService) GetFollowerCompanies(ctx context.Context, filter *model.CompanyFilter) ([]*model.CompanyFollow, error) {
	senderId := s.authenticateUser(ctx)
	followers, err := s.repo.GetFollowerCompanies(ctx, senderId, filter)
	return followers, err
}

func (s *NetworkService) GetFilteredFollowingCompanies(ctx context.Context, filter *model.CompanyFilter) ([]*model.CompanyFollow, error) {
	senderId := s.authenticateUser(ctx)
	followers, err := s.repo.GetFilteredFollowingCompanies(ctx, senderId, filter)
	return followers, err
}

func (s *NetworkService) AddCompanyToFavourites(ctx context.Context, companyId string) error {
	return s.AddCompanyToCategory(ctx, companyId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) RemoveCompanyFromFavourites(ctx context.Context, companyId string) error {
	return s.RemoveCompanyFromCategory(ctx, companyId, FAVOURITES_CATEGORY_NAME)
}

func (s *NetworkService) AddCompanyToCategory(ctx context.Context, companyId, categoryName string) error {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	cat := &model.CategoryRelation{
		OwnerId:      network.ToUserId(senderId),
		ReferralId:   network.ToCompanyId(companyId),
		CategoryName: categoryName,
		CreatedAt:    time.Now(),
	}

	_, err = s.repo.AddToCategory(ctx, cat)
	return err
}

func (s *NetworkService) RemoveCompanyFromCategory(ctx context.Context, companyId, categoryName string) error {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return err
	}
	senderId := s.authenticateUser(ctx)

	cat := &model.CategoryRelation{
		OwnerId:      network.ToUserId(senderId),
		ReferralId:   network.ToCompanyId(companyId),
		CategoryName: categoryName,
	}

	err = s.repo.RemoveFromCategory(ctx, cat)
	return err
}

func (s *NetworkService) GetSuggestedCompanies(ctx context.Context, pagination *model.Pagination) ([]*model.CompanySuggestion, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetSuggestedCompanies(ctx, senderId, pagination)
}

func (s *NetworkService) AddExperience(ctx context.Context, experience *model.AddExperienceRequest) error {
	err := s.validator.ValidateAddExperienceRequest(experience)
	if err != nil {
		return err
	}
	senderId := s.authenticateUser(ctx)
	experience.UserId = network.ToUserId(senderId)
	experience.CompanyId = network.ToCompanyId(experience.CompanyId)

	err = s.repo.AddExperience(ctx, experience)
	return err
}

func (s *NetworkService) GetAllExperience(ctx context.Context) ([]*model.Experience, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetAllExperience(ctx, senderId)
}

func (s *NetworkService) AskRecommendation(ctx context.Context, request *model.RecommendationRequestModel) error {
	senderId := s.authenticateUser(ctx)

	if senderId == request.To {
		return errors.New("can_not_send_to yourself")
	}

	recivier := request.To

	request.From = network.ToUserId(senderId)
	request.To = network.ToUserId(recivier)
	request.CreatedAt = time.Now()

	key, err := s.repo.SaveRecommendationRequest(ctx, request)
	if err != nil {
		return err
	}

	// send notification
	not := model.NewRecommendationRequest{
		RecommendationID: key,
		UserSenderID:     senderId,
		Text:             request.Text,
	}
	not.GenerateID()

	err = s.mq.SendNewRecommendationRequest(recivier, &not)
	if err != nil {
		// s.tracer
		log.Println("error while sending notification:", err)
	}

	return nil
}

func (s *NetworkService) IgnoreRecommendationRequest(ctx context.Context, key string) error {
	senderId := s.authenticateUser(ctx)

	recommendationRequest, err := s.repo.GetRecommendationRequestModel(ctx, key)
	if err != nil {
		return err
	}
	if recommendationRequest.To != network.ToUserId(senderId) {
		return errors.New("You can not ignore this recommendation request")
	}
	recommendationRequest.IsIgnored = true
	return s.repo.UpdateRecommendationRequest(ctx, key, recommendationRequest)
}

func (s *NetworkService) GetRequestedRecommendations(ctx context.Context, pagination *model.Pagination) ([]*model.RecommendationRequest, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetRequestedRecommendations(ctx, senderId, pagination)

}
func (s *NetworkService) GetReceivedRecommendationRequests(ctx context.Context, pagination *model.Pagination) ([]*model.RecommendationRequest, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetReceivedRecommendationRequests(ctx, senderId, pagination)
}

func (s *NetworkService) WriteRecommendation(ctx context.Context, recommendation *model.RecommendationModel) error {
	senderId := s.authenticateUser(ctx)

	if senderId == recommendation.To {
		return errors.New("can_not_send_to yourself")
	}

	recivier := recommendation.To

	recommendation.From = network.ToUserId(senderId)
	recommendation.To = network.ToUserId(recivier)
	recommendation.CreatedAt = time.Now()

	err := s.repo.WriteRecommendation(ctx, recommendation)
	if err != nil {
		return err
	}

	// TODO: remove request

	// send notification
	not := model.NewRecommendation{
		UserSenderID: senderId,
		Text:         recommendation.Text,
	}
	not.GenerateID()

	err = s.mq.SendNewRecommendation(recivier, &not)
	if err != nil {
		log.Println("error while sending notification:", err)
		// s.tracer
	}

	return nil
}

func (s *NetworkService) SetRecommendationVisibility(ctx context.Context, key string, visible bool) error {
	senderId := s.authenticateUser(ctx)

	rec, err := s.repo.GetRecommendationModel(ctx, key)
	if err != nil {
		return err
	}

	if rec.To != network.ToUserId(senderId) {
		return errors.New("You can not change this recommendation")
	}

	return s.repo.SetRecommendationVisibility(ctx, key, visible)
}

func (s *NetworkService) GetReceivedRecommendations(ctx context.Context, pagination *model.Pagination) ([]*model.Recommendation, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetReceivedRecommendations(ctx, senderId, pagination /*, true*/)
}

func (s *NetworkService) GetGivenRecommendations(ctx context.Context, pagination *model.Pagination) ([]*model.Recommendation, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetGivenRecommendations(ctx, senderId, pagination)
}

func (s *NetworkService) GetReceivedRecommendationsById(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.Recommendation, error) {

	// senderId := s.authenticateUser(ctx)
	//
	// isMe := false
	//
	// if senderId == userId {
	// 	isMe = true
	// }

	return s.repo.GetReceivedRecommendations(ctx, userId, pagination /*, isMe*/)
}

func (s *NetworkService) GetGivenRecommendationsById(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.Recommendation, error) {
	return s.repo.GetGivenRecommendations(ctx, userId, pagination)
}

func (s *NetworkService) GetHiddenRecommendationsById(ctx context.Context, userId string, pagination *model.Pagination) ([]*model.Recommendation, error) {
	// senderId := s.authenticateUser(ctx)

	return s.repo.GetHiddenRecommendations(ctx, userId, pagination)
}

func (s *NetworkService) BlockUser(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	err = s.chat.GetProfilesByIDs(ctx, senderId, id, true)
	if err != nil {
		return err
	}

	err = s.repo.Block(ctx, network.ToUserId(senderId), network.ToUserId(id))
	return err
}

func (s *NetworkService) UnblockUser(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	err = s.chat.GetProfilesByIDs(ctx, senderId, id, false)
	if err != nil {
		return err
	}

	err = s.repo.Unblock(ctx, network.ToUserId(senderId), network.ToUserId(id))
	return err
}

func (s *NetworkService) BlockCompany(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	err = s.chat.GetProfilesByIDs(ctx, senderId, id, true)
	if err != nil {
		return err
	}

	err = s.repo.Block(ctx, network.ToUserId(senderId), network.ToCompanyId(id))
	return err
}

func (s *NetworkService) UnblockCompany(ctx context.Context, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	senderId := s.authenticateUser(ctx)

	err = s.chat.GetProfilesByIDs(ctx, senderId, id, false)
	if err != nil {
		return err
	}

	err = s.repo.Unblock(ctx, network.ToUserId(senderId), network.ToCompanyId(id))
	return err
}

func (s *NetworkService) GetBlockedUsersOrCompanies(ctx context.Context) ([]*model.BlockedUserOrCompany, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.GetBlockedUsersOrCompanies(ctx, senderId)
}

// TODO:
func (s *NetworkService) BlockUserForCompany(ctx context.Context, companyID string, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	s.requireAdminLevelForCompany(ctx, companyID, model.AdminLevel_Admin)

	err = s.chat.GetProfilesByIDs(ctx, companyID, id, true)
	if err != nil {
		return err
	}

	err = s.repo.Block(ctx, network.ToCompanyId(companyID), network.ToUserId(id))
	return err
}

// TODO:
func (s *NetworkService) UnblockUserForCompany(ctx context.Context, companyID string, id string) error {
	err := s.validator.ValidateKey(id)
	if err != nil {
		return err
	}

	s.requireAdminLevelForCompany(ctx, companyID, model.AdminLevel_Admin)

	err = s.chat.GetProfilesByIDs(ctx, companyID, id, false)
	if err != nil {
		return err
	}

	err = s.repo.Unblock(ctx, network.ToCompanyId(companyID), network.ToUserId(id))
	return err
}

func (s *NetworkService) GetBlockedUsersForCompany(ctx context.Context, companyID string) ([]*model.User, error) {
	s.requireAdminLevelForCompany(ctx, companyID, model.AdminLevel_Admin)
	return s.repo.GetBlockedUsersForCompany(ctx, companyID)
}

//func (s *NetworkService) GetBlockedUsers(ctx context.Context) ([]*model.User, error) {
//	senderId := s.authenticateUser(ctx)
//	return s.repo.GetBlockedUsers(ctx, senderId)
//}

//func (s *NetworkService) GetBlockedCompanies(ctx context.Context) ([]*model.Company, error) {
//	senderId := s.authenticateUser(ctx)
//	return s.repo.GetBlockedCompanies(ctx, senderId)
//}

func (s *NetworkService) GetFollowingsForCompany(ctx context.Context, request *model.IdWithUserFilter) ([]*model.Follow, error) {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return nil, err
	}

	s.requireAdminLevelForCompany(ctx, request.Id, model.AdminLevel_Admin)

	return s.repo.GetFollowingsForCompany(ctx, request.Id, request.Filter)
}

func (s *NetworkService) GetFollowersForCompany(ctx context.Context, request *model.IdWithUserFilter) ([]*model.Follow, error) {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return nil, err
	}

	s.requireAdminLevelForCompany(ctx, request.Id, model.AdminLevel_Admin)

	return s.repo.GetFollowersForCompany(ctx, request.Id, request.Filter)
}

func (s *NetworkService) GetFollowingCompaniesForCompany(ctx context.Context, request *model.IdWithCompanyFilter) ([]*model.CompanyFollow, error) {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return nil, err
	}

	s.requireAdminLevelForCompany(ctx, request.Id, model.AdminLevel_Admin)

	return s.repo.GetFollowingCompaniesForCompany(ctx, request.Id, request.Filter)
}

func (s *NetworkService) GetFollowerCompaniesForCompany(ctx context.Context, request *model.IdWithCompanyFilter) ([]*model.CompanyFollow, error) {
	err := s.validator.ValidateStruct(request)
	if err != nil {
		return nil, err
	}

	s.requireAdminLevelForCompany(ctx, request.Id, model.AdminLevel_Admin)

	return s.repo.GetFollowerCompaniesForCompany(ctx, request.Id, request.Filter)
}

func (s *NetworkService) FollowForCompany(ctx context.Context, companyId, userId string) error {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return err
	}
	err = s.validator.ValidateKey(userId)
	if err != nil {
		return err
	}

	s.requireAdminLevelForCompany(ctx, companyId, model.AdminLevel_Admin)

	_, err = s.repo.Follow(ctx, &model.FollowRequest{
		FollowerId:  network.ToCompanyId(companyId),
		FollowingId: network.ToUserId(userId),
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *NetworkService) UnfollowForCompany(ctx context.Context, companyId, userId string) error {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return err
	}
	err = s.validator.ValidateKey(userId)
	if err != nil {
		return err
	}

	s.requireAdminLevelForCompany(ctx, companyId, model.AdminLevel_Admin)

	err = s.repo.Unfollow(ctx, network.ToCompanyId(companyId), network.ToUserId(userId))
	return err
}

func (s *NetworkService) FollowCompanyForCompany(ctx context.Context, followerCompanyId, followingCompanyId string) error {
	err := s.validator.ValidateKey(followerCompanyId)
	if err != nil {
		return err
	}
	err = s.validator.ValidateKey(followingCompanyId)
	if err != nil {
		return err
	}

	s.requireAdminLevelForCompany(ctx, followerCompanyId, model.AdminLevel_Admin)

	_, err = s.repo.Follow(ctx, &model.FollowRequest{
		FollowerId:  network.ToCompanyId(followerCompanyId),
		FollowingId: network.ToCompanyId(followingCompanyId),
		CreatedAt:   time.Now(),
	})

	return err
}

func (s *NetworkService) UnfollowCompanyForCompany(ctx context.Context, followerCompanyId, followingCompanyId string) error {
	err := s.validator.ValidateKey(followerCompanyId)
	if err != nil {
		return err
	}
	err = s.validator.ValidateKey(followingCompanyId)
	if err != nil {
		return err
	}

	s.requireAdminLevelForCompany(ctx, followerCompanyId, model.AdminLevel_Admin)

	err = s.repo.Unfollow(ctx, network.ToCompanyId(followerCompanyId), network.ToCompanyId(followingCompanyId))
	return err
}

func (s *NetworkService) GetSuggestedPeopleForCompany(ctx context.Context, companyId string) ([]*model.UserSuggestion, error) {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return nil, err
	}

	s.requireAdminLevelForCompany(ctx, companyId, model.AdminLevel_Admin)

	return s.repo.GetSuggestedPeopleForCompany(ctx, companyId)
}

func (s *NetworkService) GetSuggestedCompaniesForCompany(ctx context.Context, companyId string, pagination *model.Pagination) ([]*model.CompanySuggestion, error) {
	err := s.validator.ValidateKey(companyId)
	if err != nil {
		return nil, err
	}

	s.requireAdminLevelForCompany(ctx, companyId, model.AdminLevel_Admin)

	return s.repo.GetSuggestedCompaniesForCompany(ctx, companyId, pagination)
}

func (s *NetworkService) GetNumberOfFollowersForCompany(ctx context.Context, companyId string) (int, error) {
	return s.repo.GetNumberOfFollowers(ctx, network.ToCompanyId(companyId))
}

func (s *NetworkService) IsBlocked(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.IsBlocked(ctx, network.ToUserId(senderId), network.ToUserId(id))
}

func (s *NetworkService) IsBlockedForCompany(ctx context.Context, id string, companyID string) (bool, error) {
	_ = s.authenticateUser(ctx)
	v, err := s.repo.IsBlocked(ctx, network.ToCompanyId(companyID), network.ToUserId(id))
	if err != nil {
		return false, err
	}

	return v, nil
}

func (s *NetworkService) IsBlockedCompany(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	isBLocked, err := s.repo.IsBlocked(ctx, network.ToUserId(senderId), network.ToCompanyId(id))
	if err != nil {
		return false, nil
	}

	return isBLocked, nil
}

func (s *NetworkService) IsBlockedCompanyForCompany(ctx context.Context, id string, companyID string) (bool, error) {
	_ = s.authenticateUser(ctx)
	isBLocked, err := s.repo.IsBlocked(ctx, network.ToCompanyId(companyID), network.ToCompanyId(id))
	if err != nil {
		return false, nil
	}

	return isBLocked, nil
}

func (s *NetworkService) IsBlockedByUser(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.IsBlockedByUser(ctx, network.ToUserId(senderId), network.ToUserId(id))
}

func (s *NetworkService) IsBlockedByCompany(ctx context.Context, id string, companyID string) (bool, error) {
	_ = s.authenticateUser(ctx)
	v, err := s.repo.IsBlockedByUser(ctx, network.ToCompanyId(companyID), network.ToUserId(id))
	if err != nil {
		return false, err
	}

	return v, nil
}

func (s *NetworkService) IsBlockedCompanyByUser(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	isBLocked, err := s.repo.IsBlockedByUser(ctx, network.ToCompanyId(id), network.ToUserId(senderId))
	if err != nil {
		return false, nil
	}

	return isBLocked, nil
}

func (s *NetworkService) IsBlockedCompanyByCompany(ctx context.Context, id string, companyID string) (bool, error) {
	_ = s.authenticateUser(ctx)
	isBLocked, err := s.repo.IsBlockedByUser(ctx, network.ToCompanyId(companyID), network.ToCompanyId(id))
	if err != nil {
		return false, nil
	}

	return isBLocked, nil
}

func (s *NetworkService) IsFollowing(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.IsFollowing(ctx, network.ToUserId(senderId), network.ToUserId(id))
}

func (s *NetworkService) IsFollowingForCompany(ctx context.Context, id string, companyID string) (bool, error) {
	_ = s.authenticateUser(ctx)
	v, err := s.repo.IsFollowing(ctx, network.ToCompanyId(companyID), network.ToUserId(id))
	if err != nil {
		log.Println("error:", err)
		return false, err
	}

	return v, nil
}

func (s *NetworkService) IsFollowingCompany(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.IsFollowing(ctx, network.ToUserId(senderId), network.ToCompanyId(id))
}

func (s *NetworkService) IsFollowingCompanyForCompany(ctx context.Context, id string, companyID string) (bool, error) {
	_ = s.authenticateUser(ctx)
	return s.repo.IsFollowing(ctx, network.ToCompanyId(companyID), network.ToCompanyId(id))
}

func (s *NetworkService) IsFavourite(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.IsFavourite(ctx, network.ToUserId(senderId), network.ToUserId(id))
}
func (s *NetworkService) IsFavouriteCompany(ctx context.Context, id string) (bool, error) {
	senderId := s.authenticateUser(ctx)
	return s.repo.IsFavourite(ctx, network.ToUserId(senderId), network.ToCompanyId(id))
}

func (s *NetworkService) GetFriendIdsOf(ctx context.Context, userId string) ([]string, error) {
	friendships, err := s.repo.GetFriendsOf(ctx, userId, &model.FriendshipFilter{})
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(friendships))
	for i, friendship := range friendships {
		ids[i] = friendship.Friend.Id
	}

	return ids, nil
}

func (s *NetworkService) GetUserCountings(ctx context.Context, userId string) (*model.UserCountings, error) {
	senderId := s.authenticateUser(ctx)

	countings, err := s.repo.GetUserCountings(ctx, userId)
	if err != nil {
		return nil, err
	}

	isAllowed, err := s.checkConnectionPrivacy(ctx, userId, senderId)
	if err != nil {
		return nil, err
	}
	if !isAllowed {
		if countings != nil {
			countings.NumOfConnections = 0
		}
	}

	return countings, nil
}

func (s *NetworkService) GetCategoryTree(ctx context.Context) (*model.CategoryTree, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetCategoryTree(ctx, senderId)
}

func (s *NetworkService) GetCategoryTreeForFollowings(ctx context.Context) (*model.CategoryTree, error) {
	senderId := s.authenticateUser(ctx)

	return s.repo.GetCategoryTreeForFollowings(ctx, senderId)
}

func (s *NetworkService) GetCategoryTreeForFollowingsForCompany(ctx context.Context, companyKey string) (*model.CategoryTree, error) {
	s.requireAdminLevelForCompany(ctx, companyKey, "Admin")

	return s.repo.GetCategoryTreeForFollowings(ctx, companyKey)
}

func (s *NetworkService) CreateCategory(ctx context.Context, parent, name string) error {
	senderId := s.authenticateUser(ctx)

	tree, err := s.repo.GetCategoryTree(ctx, senderId)
	if err != nil {
		return err
	}
	err = handleCreateCategory(tree, parent, name)
	if err != nil {
		return err
	}

	return s.repo.SetCategoryTree(ctx, tree)
}

func (s *NetworkService) CreateCategoryForFollowings(ctx context.Context, parent, name string) error {
	senderId := s.authenticateUser(ctx)

	tree, err := s.repo.GetCategoryTreeForFollowings(ctx, senderId)
	if err != nil {
		return err
	}
	err = handleCreateCategory(tree, parent, name)
	if err != nil {
		return err
	}

	return s.repo.SetCategoryTreeForFollowings(ctx, tree)
}

func (s *NetworkService) CreateCategoryForFollowingsForCompany(ctx context.Context, companyId, parent, name string) error {
	s.requireAdminLevelForCompany(ctx, companyId, "Admin")

	tree, err := s.repo.GetCategoryTreeForFollowings(ctx, companyId)
	if err != nil {
		return err
	}
	err = handleCreateCategory(tree, parent, name)
	if err != nil {
		return err
	}

	return s.repo.SetCategoryTreeForFollowings(ctx, tree)
}

func handleCreateCategory(tree *model.CategoryTree, parent, name string) error {
	added := false
	for i, item := range tree.Categories {
		if item.UniqueName == parent && item.HasChildren {
			for _, subitem := range item.Children {
				if subitem.Name == name {
					return errors.New("This category already exists")
				}
			}
			tree.Categories[i].Children = append(item.Children, model.CategoryItem{
				UniqueName:  fmt.Sprint(item.UniqueName, "__", name),
				Name:        name,
				HasChildren: false,
			})
			added = true
			break
		}
	}
	if added {
		return nil
	} else {
		return errors.New("Invalid category path")
	}
}

func (s *NetworkService) RemoveCategory(ctx context.Context, parent, name string) error {
	senderId := s.authenticateUser(ctx)
	tree, err := s.repo.GetCategoryTree(ctx, senderId)
	if err != nil {
		return err
	}
	uniqueName, err := handleRemoveCategory(tree, parent, name)
	if err != nil {
		return err
	}
	err = s.repo.SetCategoryTree(ctx, tree)
	if err != nil {
		return err
	}
	return s.repo.ClearCategory(ctx, senderId, uniqueName)
}

func (s *NetworkService) RemoveCategoryForFollowings(ctx context.Context, parent, name string) error {
	senderId := s.authenticateUser(ctx)
	tree, err := s.repo.GetCategoryTreeForFollowings(ctx, senderId)
	if err != nil {
		return err
	}
	uniqueName, err := handleRemoveCategory(tree, parent, name)
	if err != nil {
		return err
	}
	err = s.repo.SetCategoryTreeForFollowings(ctx, tree)
	if err != nil {
		return err
	}
	return s.repo.ClearCategoryForFollowings(ctx, senderId, uniqueName)
}

func (s *NetworkService) RemoveCategoryForFollowingsForCompany(ctx context.Context, companyId, parent, name string) error {
	s.requireAdminLevelForCompany(ctx, companyId, "Admin")

	tree, err := s.repo.GetCategoryTreeForFollowings(ctx, companyId)
	if err != nil {
		return err
	}
	uniqueName, err := handleRemoveCategory(tree, parent, name)
	if err != nil {
		return err
	}
	err = s.repo.SetCategoryTreeForFollowings(ctx, tree)
	if err != nil {
		return err
	}
	return s.repo.ClearCategoryForFollowingsForCompany(ctx, companyId, uniqueName)
}

func handleRemoveCategory(tree *model.CategoryTree, parent, name string) (string, error) {
	for i, item := range tree.Categories {
		if item.UniqueName == parent && item.HasChildren {
			for j, subitem := range item.Children {
				if subitem.Name == name {
					tree.Categories[i].Children = append(item.Children[:j], item.Children[j+1:]...)
					return subitem.UniqueName, nil
				}
			}
		}
	}
	return "", errors.New("Category path not found")
}

func (s *NetworkService) IsFriendRequestSend(ctx context.Context, userID string) (bool, error) {
	senderId := s.authenticateUser(ctx)

	friendships, err := s.repo.GetFriendshipRequests(ctx, senderId, &model.FriendshipRequestFilter{
		Sent:   true,
		Status: string(model.FriendshipStatus_Requested),
	})
	if err != nil {
		return false, err
	}

	var friendshipID string

	for i := range friendships {
		if friendships[i].Friend.Id == userID {
			friendshipID = friendships[i].Id
			break
		}
	}

	if friendshipID == "" {
		return false, nil
	}

	friendship, err := s.repo.GetFriendship(ctx, friendshipID)
	if err != nil {
		return false, err
	}

	if friendship.Status == model.FriendshipStatus_Requested {
		return true, nil
	}

	return false, nil
}

func (s *NetworkService) IsFriendRequestRecieved(ctx context.Context, userID string) (bool, string, error) {
	senderId := s.authenticateUser(ctx)

	friendships, err := s.repo.GetFriendshipRequests(ctx, userID, &model.FriendshipRequestFilter{
		Sent:   true,
		Status: string(model.FriendshipStatus_Requested),
	})
	if err != nil {
		return false, "", err
	}

	var friendshipID string

	for i := range friendships {
		if friendships[i].Friend.Id == senderId {
			friendshipID = friendships[i].Id
			break
		}
	}

	if friendshipID == "" {
		return false, "", nil
	}

	friendship, err := s.repo.GetFriendship(ctx, friendshipID)
	if err != nil {
		return false, "", err
	}

	if friendship.Status == model.FriendshipStatus_Requested {
		return true, friendshipID, nil
	}

	return false, "", nil
}

func (s *NetworkService) GetFriendshipID(ctx context.Context, userID string) (string, error) {
	senderId := s.authenticateUser(ctx)

	if senderId == userID {
		return "", nil
	}

	friendships, err := s.repo.GetAllFriendship(ctx, userID, &model.FriendshipRequestFilter{})
	if err != nil {
		return "", err
	}

	var friendshipID string

	for i := range friendships {
		if friendships[i].Friend.Id == senderId {
			friendshipID = friendships[i].Id
			break
		}
	}

	if friendshipID == "" {
		return "", nil
	}

	return friendshipID, nil
}

func (s *NetworkService) GetFriendsOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	senderId := s.authenticateUser(ctx)

	isAllowed, err := s.checkConnectionPrivacy(ctx, userID, senderId)
	if err != nil {
		return nil, 0, err
	}
	if !isAllowed {
		return nil, 0, nil
	}

	ids, amount, err := s.repo.GetFriendsOfUser(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.user.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from user service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowsOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowsOfUser(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.user.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from user service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowersOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowersOfUser(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.user.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from user service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowsCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowsCompaniesOfUser(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.company.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from company service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowersCompaniesOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowersCompaniesOfUser(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.company.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from company service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowsOfCompany(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowsOfCompany(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.user.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from user service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowersOfCompany(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowersOfCompany(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.user.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from user service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowsCompaniesOfCompany(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowsCompaniesOfCompany(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.company.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from company service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetFollowersCompaniesOfCompany(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	// senderId := s.authenticateUser(ctx)

	ids, amount, err := s.repo.GetFollowersCompaniesOfCompany(ctx, userID, first, after)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, 0, err
	}

	profiles, err := s.company.GetProfilesByIDs(ctx, ids)
	if err != nil {
		log.Println("error get profile from company service:", err)
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetMutualFriendsOfUser(ctx context.Context, userID string, first uint32, after uint32) (friends interface{}, amount int64, err error) {
	senderId := s.authenticateUser(ctx)

	isAllowed, err := s.checkConnectionPrivacy(ctx, userID, senderId)
	if err != nil {
		return nil, 0, err
	}
	if !isAllowed {
		return nil, 0, nil
	}

	ids, amount, err := s.repo.GetMutualFriendsOfUser(ctx, senderId, userID, first, after)
	if err != nil {
		return nil, 0, err
	}

	profiles, err := s.user.GetProfilesByIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	return profiles, amount, nil
}

func (s *NetworkService) GetCompanyCountings(ctx context.Context, companyID string) (*model.CompanyCountings, error) {
	return s.repo.GetCompanyCountings(ctx, companyID)
}

func (s *NetworkService) checkConnectionPrivacy(ctx context.Context, userID string, senderID string) (bool, error) {
	privacy, err := s.user.GetConectionsPrivacy(ctx, userID)
	if err != nil {
		log.Println("error: get privacy:", err)
		return false, err
	}
	if privacy != "NONE" {
		switch privacy {
		case "ME":
			if senderID != userID {
				return false, nil
			}
		case "MEMBERS":
			if senderID == "" {
				return false, nil
			}
		case "MY_CONNECTIONS":
			isFriend, err := s.IsFriend(ctx, userID)
			if err != nil {
				return false, err
			}
			if !isFriend {
				return false, nil
			}
		}
	}
	return true, nil
}

func (s *NetworkService) GetAmountOfMutualFriends(ctx context.Context, userID string) (int32, error) {
	senderId := s.authenticateUser(ctx)

	isAllowed, err := s.checkConnectionPrivacy(ctx, userID, senderId)
	if err != nil {
		return 0, err
	}
	if !isAllowed {
		return 0, nil
	}

	amount, err := s.repo.GetAmountOfMutualFriends(ctx, senderId, userID)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (s *NetworkService) GetFollowersIDs(ctx context.Context, id string, isCompanyID bool) ([]string, error) {
	log.Println(id, isCompanyID)
	ids, err := s.repo.GetFollowersIDs(ctx, id, isCompanyID)
	if err != nil {
		log.Println("error request DB:", err)
		return nil, err
	}

	return ids, nil
}
