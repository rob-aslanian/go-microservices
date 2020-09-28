package serverRPC

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/invitation"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/status"
	userReport "gitlab.lan/Rightnao-site/microservices/user/pkg/user_report"
)

// IdentifyCountry ...
func (s Server) IdentifyCountry(ctx context.Context, data *userRPC.Empty) (*userRPC.CountryID, error) {
	id, err := s.service.IdentifyCountry(ctx)
	if err != nil {
		return nil, err
	}

	return &userRPC.CountryID{
		ID: id,
	}, nil
}

// CheckToken ...
func (s Server) CheckToken(ctx context.Context, data *userRPC.Empty) (*userRPC.BooleanValue, error) {
	value, err := s.service.CheckToken(ctx)
	if err != nil {
		return nil, err
	}
	return &userRPC.BooleanValue{Value: value}, nil
}

// IsUsernameBusy ...
func (s Server) IsUsernameBusy(ctx context.Context, data *userRPC.Username) (*userRPC.BooleanValue, error) {
	value, err := s.service.IsUsernameBusy(ctx, data.GetUsername())
	if err != nil {
		return nil, err
	}

	return &userRPC.BooleanValue{Value: value}, nil
}

// RegisterUser creates new user account
func (s Server) RegisterUser(ctx context.Context, data *userRPC.RegisterRequest) (*userRPC.LoginResponse, error) {
	userID, url, tmpToken, err := s.service.CreateNewAccount(ctx, registerRequestToAccount(data), data.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userRPC.LoginResponse{
		Status: userRPC.Status_NOT_ACTIVATED,
		Token:  tmpToken,
		URL:    url,
		UserId: userID,
	}, nil
}

// ActivateUser ...
func (s Server) ActivateUser(ctx context.Context, data *userRPC.ActivateUserRequest) (*userRPC.ActivateUserResponse, error) {
	res, err := s.service.ActivateUser(ctx, data.GetCode(), data.GetUserID())
	if err != nil {
		return nil, err
	}

	return &userRPC.ActivateUserResponse{
		UserId:    res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Avatar:    res.Avatar,
		Token:     res.Token,
		URL:       res.URL,
	}, nil
}

// SendRecover creates temprorary code for recovery
func (s Server) SendRecover(ctx context.Context, data *userRPC.SendRecoverRequest) (*userRPC.Empty, error) {
	// make in future option how to recover
	err := s.service.Recover(ctx, data.GetLogin(), data.GetSendUsername(), data.GetResetPassword())
	return &userRPC.Empty{}, err
}

// RecoverPassword ...
func (s Server) RecoverPassword(ctx context.Context, data *userRPC.RecoverPasswordRequest) (*userRPC.Empty, error) {
	err := s.service.RecoverPassword(ctx, data.GetCode(), data.GetUserId(), data.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// Login ...
func (s Server) Login(ctx context.Context, data *userRPC.Credentials) (*userRPC.LoginResponse, error) {
	res, err := s.service.Login(ctx, data.GetLogin(), data.GetPassword(), data.GetTwoFACode())
	if err != nil {
		return nil, err
	}

	return &userRPC.LoginResponse{
		Avatar:        res.Avatar,
		UserId:        res.ID,
		URL:           res.URL,
		TwoFARequired: res.Is2FAEnabled,
		Token:         res.Token,
		FirstName:     res.FirstName,
		LastName:      res.LastName,
		Status:        userStatusToStatusRPC(res.Status),
		Gender:        res.Gender,
		Email:         res.Email,
	}, nil
}

// SignOut ...
func (s Server) SignOut(ctx context.Context, data *userRPC.Empty) (*userRPC.Empty, error) {
	err := s.service.SignOut(ctx)
	if err != nil {
		return nil, err
	}
	return &userRPC.Empty{}, nil
}

// SignOutSession ...
func (s Server) SignOutSession(ctx context.Context, sessionID *userRPC.SessionID) (*userRPC.Empty, error) {
	err := s.service.SignOutSession(ctx, sessionID.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SignOutFromAll ...
func (s Server) SignOutFromAll(ctx context.Context, data *userRPC.Empty) (*userRPC.Empty, error) {
	err := s.service.SignOutFromAll(ctx)
	if err != nil {
		return nil, err
	}
	return &userRPC.Empty{}, nil
}

// GetAccount ...
func (s Server) GetAccount(ctx context.Context, data *userRPC.Empty) (*userRPC.Account, error) {
	acc, err := s.service.GetAccount(ctx)
	if err != nil {
		return nil, err
	}

	return accountToRPC(acc), nil
}

// ChangeFirstName changes first name of user.
func (s Server) ChangeFirstName(ctx context.Context, data *userRPC.FirstName) (*userRPC.Empty, error) {
	err := s.service.ChangeFirstName(ctx, data.GetFirstName())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeLastname changes last name of user.
func (s Server) ChangeLastname(ctx context.Context, data *userRPC.Lastname) (*userRPC.Empty, error) {
	err := s.service.ChangeLastName(ctx, data.GetLastname())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangePatronymic changes patronomyc and its permission
func (s Server) ChangePatronymic(ctx context.Context, data *userRPC.Patronymic) (*userRPC.Empty, error) {
	var patr *string

	if !data.GetIsPatronymicNull() {
		p := data.GetPatronymic()
		patr = &p
	}

	err := s.service.ChangePatronymic(
		ctx,
		patr,
		permissionRPCToAccountPermission(data.GetPermission()),
	)

	return &userRPC.Empty{}, err
}

// ChangeNickname changes nickname and its permission
func (s Server) ChangeNickname(ctx context.Context, data *userRPC.Nickname) (*userRPC.Empty, error) {
	var nick *string

	if !data.GetIsNicknameNull() {
		p := data.GetNickname()
		nick = &p
	}

	err := s.service.ChangeNickname(
		ctx,
		nick,
		permissionRPCToAccountPermission(data.GetPermission()),
	)

	return &userRPC.Empty{}, err
}

// ChangeMiddleName changes nickname and its permission
func (s Server) ChangeMiddleName(ctx context.Context, data *userRPC.Middlename) (*userRPC.Empty, error) {
	var name *string

	if !data.GetIsMiddlenameNull() {
		p := data.GetMiddlename()
		name = &p
	}

	err := s.service.ChangeMiddleName(
		ctx,
		name,
		permissionRPCToAccountPermission(data.GetPermission()),
	)

	return &userRPC.Empty{}, err
}

// ChangeNameOnNativeLanguage ...
func (s Server) ChangeNameOnNativeLanguage(ctx context.Context, data *userRPC.NativeName) (*userRPC.Empty, error) {
	var name *string
	var lang *string

	if !data.GetIsNameNull() {
		p := data.GetName()
		name = &p
	}

	if data.GetLanguageID() != "" {
		p := data.GetLanguageID()
		lang = &p
	}

	err := s.service.ChangeNameOnNativeLanguage(
		ctx,
		name,
		lang,
		permissionRPCToAccountPermission(data.GetPermission()),
	)

	return &userRPC.Empty{}, err
}

// ChangeBirthday ...
func (s Server) ChangeBirthday(ctx context.Context, data *userRPC.Birthday) (*userRPC.Empty, error) {
	var bd *time.Time
	if data.GetBirthday() != "" {
		t := stringDateToTime(data.GetBirthday())
		bd = &t
	}

	err := s.service.ChangeBirthday(
		ctx,
		bd,
		permissionRPCToAccountPermission(data.GetPermission()),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeGender ...
func (s Server) ChangeGender(ctx context.Context, data *userRPC.Gender) (*userRPC.Empty, error) {
	var gen *string

	if !data.GetIsGenderNull() {
		t := data.GetGender().String()
		gen = &t
	}

	err := s.service.ChangeGender(
		ctx,
		gen,
		permissionRPCToAccountPermission(data.GetPermission()),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddEmail ...
func (s Server) AddEmail(ctx context.Context, data *userRPC.Email) (*userRPC.ID, error) {
	id, err := s.service.AddEmail(
		ctx,
		data.GetEmail(),
		permissionRPCToAccountPermission(data.GetPermission()),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// RemoveEmail ...
func (s Server) RemoveEmail(ctx context.Context, data *userRPC.Email) (*userRPC.Empty, error) {
	err := s.service.RemoveEmail(ctx, data.GetId())
	if err != nil {
		return nil, err
	}
	return &userRPC.Empty{}, nil
}

// ChangeEmail ...
func (s Server) ChangeEmail(ctx context.Context, data *userRPC.Email) (*userRPC.Empty, error) {
	err := s.service.ChangeEmail(
		ctx,
		data.GetId(),
		permissionRPCToAccountPermission(data.GetPermission()),
		data.GetIsPrimary(),
	)
	if err != nil {
		return nil, err
	}
	return &userRPC.Empty{}, nil
}

// AddPhone ...
func (s Server) AddPhone(ctx context.Context, data *userRPC.Phone) (*userRPC.ID, error) {
	id, err := s.service.AddPhone(
		ctx,
		countryCodeRPCtoAccountCountryCode(data.GetCountryCode()),
		data.GetNumber(),
		permissionRPCToAccountPermission(data.GetPermission()),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// RemovePhone ...
func (s Server) RemovePhone(ctx context.Context, data *userRPC.Phone) (*userRPC.Empty, error) {
	err := s.service.RemovePhone(ctx, data.GetId())
	if err != nil {
		return nil, err
	}
	return &userRPC.Empty{}, nil
}

// ChangePhone ...
func (s Server) ChangePhone(ctx context.Context, data *userRPC.Phone) (*userRPC.Empty, error) {
	err := s.service.ChangePhone(
		ctx,
		data.GetId(),
		permissionRPCToAccountPermission(data.GetPermission()),
		data.GetIsPrimary(),
	)
	if err != nil {
		return nil, err
	}
	return &userRPC.Empty{}, nil
}

// AddMyAddress ...
func (s Server) AddMyAddress(ctx context.Context, data *userRPC.MyAddress) (*userRPC.ID, error) {
	id, err := s.service.AddMyAddress(
		ctx,
		myAddressRPCToAccountMyAddress(data),
	)

	return &userRPC.ID{ID: id}, err
}

// RemoveMyAddress ...
func (s Server) RemoveMyAddress(ctx context.Context, data *userRPC.MyAddress) (*userRPC.Empty, error) {
	err := s.service.RemoveMyAddress(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeMyAddress ...
func (s Server) ChangeMyAddress(ctx context.Context, data *userRPC.MyAddress) (*userRPC.Empty, error) {
	err := s.service.ChangeMyAddress(
		ctx,
		myAddressRPCToAccountMyAddress(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddOtherAddress ...
// func (s Server) AddOtherAddress(ctx context.Context, data *userRPC.OtherAddress) (*userRPC.ID, error) {
func (s Server) AddOtherAddress(ctx context.Context, data *userRPC.OtherAddress) (*userRPC.ID, error) {
	id, err := s.service.AddOtherAddress(
		ctx,
		otherAddressRPCToAccountMyAddress(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// RemoveOtherAddress ...
func (s Server) RemoveOtherAddress(ctx context.Context, data *userRPC.OtherAddress) (*userRPC.Empty, error) {
	err := s.service.RemoveOtherAddress(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeOtherAddress ...
func (s Server) ChangeOtherAddress(ctx context.Context, data *userRPC.OtherAddress) (*userRPC.Empty, error) {
	err := s.service.ChangeOtherAddress(
		ctx,
		otherAddressRPCToAccountMyAddress(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeUILanguage ...
func (s Server) ChangeUILanguage(ctx context.Context, data *userRPC.Language) (*userRPC.Empty, error) {
	err := s.service.ChangeUILanguage(ctx, data.GetLanguage())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangePrivacy ...
func (s Server) ChangePrivacy(ctx context.Context, data *userRPC.ChangePrivacyRequest) (*userRPC.Empty, error) {
	priv := data.GetPrivacy()
	permType := data.GetPermission()

	err := s.service.ChangePrivacy(
		ctx,
		privacySettingsRPCToAccountPrivacyItem(&priv),
		permissionTypeRPCToAccountPermissionType(&permType),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangePassword ...
func (s Server) ChangePassword(ctx context.Context, data *userRPC.ChangePasswordRequest) (*userRPC.Empty, error) {
	err := s.service.ChangePassword(ctx, data.GetOldPassword(), data.GetNewPassword())

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// Init2FA ...
func (s Server) Init2FA(ctx context.Context, data *userRPC.Empty) (*userRPC.TwoFAResponse, error) {
	qr, url, key, err := s.service.Init2FA(ctx)
	if err != nil {
		return nil, err
	}

	return &userRPC.TwoFAResponse{
		QR:  qr,
		URL: url,
		Key: key,
	}, nil
}

// Enable2FA ...
func (s Server) Enable2FA(ctx context.Context, data *userRPC.TwoFACode) (*userRPC.Empty, error) {
	err := s.service.Enable2FA(ctx, data.GetCode())
	if err != nil {
		return nil, err
	}
	return &userRPC.Empty{}, nil
}

// Disable2FA ...
func (s Server) Disable2FA(ctx context.Context, data *userRPC.TwoFACode) (*userRPC.Empty, error) {
	err := s.service.Disable2FA(ctx, data.GetCode())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// DeactivateAccount ...
func (s Server) DeactivateAccount(ctx context.Context, data *userRPC.CheckPasswordRequest) (*userRPC.Empty, error) {
	err := s.service.DeactivateAccount(ctx, data.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// Profile

// GetProfile ...
func (s Server) GetProfile(ctx context.Context, data *userRPC.ProfileRequest) (*userRPC.Profile, error) {
	profile, err := s.service.GetProfile(ctx, data.GetURL(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	return profileToProfileRPC(profile), nil
}

// GetProfileByID ...
func (s Server) GetProfileByID(ctx context.Context, data *userRPC.ID) (*userRPC.Profile, error) {
	profile, err := s.service.GetProfileByID(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return profileToProfileRPC(profile), nil
}

// GetProfilesByID ...
func (s Server) GetProfilesByID(ctx context.Context, data *userRPC.UserIDs) (*userRPC.ProfileList, error) {
	profiles, err := s.service.GetProfilesByID(ctx, data.GetID(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	prof := make([]*userRPC.Profile, 0, len(profiles))

	for i := range profiles {
		prof = append(prof, profileToProfileRPC(profiles[i]))
	}

	return &userRPC.ProfileList{
		Profiles: prof,
	}, nil
}

// GetMapProfilesByID ...
func (s Server) GetMapProfilesByID(ctx context.Context, data *userRPC.UserIDs) (*userRPC.MapProfiles, error) {
	profiles, err := s.service.GetMapProfilesByID(ctx, data.GetID(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	mapProfiles := make(map[string]*userRPC.Profile, len(profiles))

	for _, pr := range profiles {
		mapProfiles[pr.GetID()] = profileToProfileRPC(pr)
	}

	return &userRPC.MapProfiles{
		Profiles: mapProfiles,
	}, nil
}

// GetMyCompanies ...
func (s Server) GetMyCompanies(ctx context.Context, data *userRPC.Empty) (*userRPC.Companies, error) {
	prof, err := s.service.GetMyCompanies(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: refactor it
	if profiles, ok := prof.([]*companyRPC.Profile); ok {
		return &userRPC.Companies{
			Companies: profiles,
		}, nil
	}

	return nil, errors.New("internal_error")
}

// GetPortfolios ...
func (s Server) GetPortfolios(ctx context.Context, data *userRPC.RequestPortfolios) (*userRPC.PortfolioList, error) {
	result, err := s.service.GetPortfolios(
		ctx,
		data.GetCompanyID(),
		data.GetUserID(),
		data.GetFirst(),
		data.GetAfter(),
		rpcContentTypeToContentType(data.GetContentType()),
	)

	if err != nil {
		return nil, err
	}

	var amount int32 = 0
	ports := profilePortfoliosToPortfoliosRPC(result)

	if result != nil {
		amount = result.PortfolioAmount
	}

	return &userRPC.PortfolioList{
		PortfolioAmount: amount,
		Portfolios:      ports,
	}, nil
}

func profilePortfoliosToPortfoliosRPC(data *profile.Portfolios) []*userRPC.Portfolio {

	if data == nil {
		return nil
	}

	if data.Portfolios == nil {
		return nil
	}

	ports := make([]*userRPC.Portfolio, 0, len(data.Portfolios))

	for i := range data.Portfolios {
		ports = append(ports, profilePortfolioToPortfolioRPC(data.Portfolios[i]))
	}

	return ports

}

func (s Server) GetPortfolioByID(ctx context.Context, data *userRPC.RequestPortfolio) (*userRPC.Portfolio, error) {

	result, err := s.service.GetPortfolioByID(
		ctx,
		data.GetCompanyID(),
		data.GetUserID(),
		data.GetPortfolioID(),
	)

	if err != nil {
		return nil, err
	}

	res := &userRPC.Portfolio{
		ID:                result.GetID(),
		Tittle:            result.Title,
		Description:       result.Description,
		LikesCount:        result.LikesCount,
		Tools:             make([]string, 0, len(result.Tools)),
		Files:             make([]*userRPC.File, 0, len(result.Files)),
		ViewsCount:        result.ViewsCount,
		SharedCount:       result.SharedCount,
		SavedCount:        result.SavedCount,
		IsCommentDisabled: result.IsCommentClosed,
		CreatedAt:         result.CreatedAt.String(),
		ContentType:       contentTypeToRPC(result.ContentType),
		HasLiked:          result.HasLiked,
	}

	// Files
	for _, f := range result.Files {
		if f != nil {
			res.Files = append(res.Files, profileFiletoFileRPC(f))
		}
	}

	// Tools
	for i := range result.Tools {
		if result.Tools[i] != nil {
			res.Tools = append(res.Tools, *result.Tools[i])
		}
	}

	return res, nil

}

func contentTypeToRPC(ctp profile.ContentType) userRPC.ContentTypeEnum {
	var content userRPC.ContentTypeEnum = userRPC.ContentTypeEnum_Content_Type_Photo

	switch ctp {
	case profile.ContentTypeArticle:
		content = userRPC.ContentTypeEnum_Content_Type_Article
	case profile.ContentTypeVideo:
		content = userRPC.ContentTypeEnum_Content_Type_Video
	case profile.ContentTypeAudio:
		content = userRPC.ContentTypeEnum_Content_Type_Audio
	}

	return content
}

func rpcContentTypeToContentType(ctp userRPC.ContentTypeEnum) string {
	var content string = "photo"

	switch ctp {
	case userRPC.ContentTypeEnum_Content_Type_Article:
		content = "article"
	case userRPC.ContentTypeEnum_Content_Type_Video:
		content = "video"
	case userRPC.ContentTypeEnum_Content_Type_Audio:
		content = "audio"
	}

	return content
}

// GetAllUsersForAdmin ...
func (s Server) GetAllUsersForAdmin(ctx context.Context, data *userRPC.Pagination) (*userRPC.GetAllUsersForAdminResponse, error) {
	res, err := s.service.GetAllUsersForAdmin(
		ctx,
		data.GetFirst(),
		data.GetAfter(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.GetAllUsersForAdminResponse{
		UserAmount: res.UsersAmount,
		Users:      usersToRPC(res.Users),
	}, nil
}

func (s Server) ChangeUserStatus(ctx context.Context, data *userRPC.ChangeUserStatusRequest) (*userRPC.Empty, error) {
	err := s.service.ChangeUserStatus(
		ctx,
		data.GetUserID(),
		statusRPCToStatus(data.GetStatus()),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

func usersToRPC(data []*account.User) []*userRPC.UserForAdmin {
	if data == nil {
		return nil
	}

	users := make([]*userRPC.UserForAdmin, 0, len(data))

	for _, u := range data {
		users = append(users, userToRPC(u))
	}

	return users
}

func userToRPC(data *account.User) *userRPC.UserForAdmin {
	if data == nil {
		return nil
	}

	profile := &userRPC.UserForAdmin{
		ID:                     data.GetUserID(),
		Avatar:                 data.Avatar,
		Firstname:              data.FirstName,
		Lasttname:              data.Lastname,
		URL:                    data.URL,
		Status:                 userStatusToStatusRPC(data.Status),
		ProfileCompletePercent: int32(data.CompletePercent),
		Location:               accountLocationToProfileLocationRPC(data.Location),
		DateOfActivation:       timeToStringDayMonthAndYear(data.CreatedAt),
	}

	if data.Gender != nil {
		if *data.Gender == "FEMALE" {
			profile.Gender = userRPC.GenderValue_FEMALE
		} else {
			profile.Gender = userRPC.GenderValue_MALE
		}
	}

	// BirthDay
	if data.Birthday != nil {
		profile.Birthday = timeToStringDayMonthAndYear(data.Birthday.Birthday)
	}

	// Email
	if data.Email != nil {
		profile.Email = *data.Email
	}

	// Phones
	for _, e := range data.Phones {
		if e.Primary {
			profile.PhoneNumber = e.CountryCode.Code + e.Number
		}
	}

	return profile
}

func statusRPCToStatus(data userRPC.Status) status.UserStatus {

	s := status.UserStatusUnKnown

	switch data {
	case userRPC.Status_ACTIVATED:
		s = status.UserStatusActivated
	case userRPC.Status_NOT_ACTIVATED:
		s = status.UserStatusNotActivated
	case userRPC.Status_DISABLED:
		s = status.UserStatusDeactivated
	}

	return s
}

// Wallet

// ContactInvitationForWallet ...
func (s Server) ContactInvitationForWallet(ctx context.Context, wallet *userRPC.InvitationWalletRequest) (*userRPC.WalletResponse, error) {

	err := s.service.ContactInvitationForWallet(
		ctx,
		wallet.GetName(),
		wallet.GetEmail(),
		wallet.GetMessage(),
		wallet.GetSilverCoins(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.WalletResponse{
		Amount: &userRPC.WalletAmountResponse{
			GoldCoins:     0,
			SilverCoins:   wallet.GetSilverCoins(),
			PendingAmount: 0,
		},
		Status: userRPC.WalletStatusEnum_DONE,
	}, nil
}

// AddGoldCoinsToWallet ...
func (s Server) AddGoldCoinsToWallet(ctx context.Context, data *userRPC.WalletAddGoldCoins) (*userRPC.Empty, error) {
	err := s.service.AddGoldCoinsToWallet(
		ctx,
		data.GetUserID(),
		data.GetCoins(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// CreateWalletAccount ...
func (s Server) CreateWalletAccount(ctx context.Context, data *userRPC.UserId) (*userRPC.Empty, error) {
	err := s.service.CreateWalletAccount(
		ctx,
		data.GetId(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetUserByInvitedID ...
func (s Server) GetUserByInvitedID(ctx context.Context, data *userRPC.UserId) (*userRPC.WalletInvitedByCount, error) {
	count, err := s.service.GetUserByInvitedID(
		ctx,
		data.GetId(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.WalletInvitedByCount{
		Count: count,
	}, nil
}

//@@@  NEW_PORTOFLIO @@@//

// GetUserPortfolioInfo ...
func (s Server) GetUserPortfolioInfo(ctx context.Context, data *userRPC.UserId) (*userRPC.PortfolioInfo, error) {

	res, err := s.service.GetUserPortfolioInfo(
		ctx,
		data.GetId(),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.PortfolioInfo{
		CommentCount: res.CommentCount,
		LikeCount:    res.LikesCount,
		ViewCount:    res.ViewsCount,
		HasPhoto:     res.HasPhoto,
		HasVideo:     res.HasVideo,
		HasArticle:   res.HasArticle,
		HasAudio:     res.HasAudio,
	}, nil
}

// AddPortfolio ...
func (s Server) AddPortfolio(ctx context.Context, portoflio *userRPC.Portfolio) (*userRPC.ID, error) {

	id, err := s.service.AddPortfolio(
		ctx,
		portfolioRPCToProfilePortfolio(portoflio),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// LikeUserPortfolio ...
func (s Server) LikeUserPortfolio(ctx context.Context, data *userRPC.PortfolioAction) (*userRPC.Empty, error) {

	err := s.service.LikeUserPortfolio(
		ctx,
		data.GetOwnerID(),
		data.GetPortfolioID(),
		data.GetCompanyID(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// UnLikeUserPortfolio ...
func (s Server) UnLikeUserPortfolio(ctx context.Context, data *userRPC.PortfolioAction) (*userRPC.Empty, error) {

	err := s.service.UnLikeUserPortfolio(
		ctx,
		data.GetOwnerID(),
		data.GetPortfolioID(),
		data.GetCompanyID(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddViewCountToPortfolio ...
func (s Server) AddViewCountToPortfolio(ctx context.Context, data *userRPC.PortfolioAction) (*userRPC.Empty, error) {

	err := s.service.AddViewCountToPortfolio(
		ctx,
		data.GetOwnerID(),
		data.GetPortfolioID(),
		data.GetCompanyID(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddSavedCountToPortfolio ...
func (s Server) AddSavedCountToPortfolio(ctx context.Context, data *userRPC.PortfolioAction) (*userRPC.Empty, error) {

	err := s.service.AddSavedCountToPortfolio(
		ctx,
		data.GetOwnerID(),
		data.GetPortfolioID(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetUserPortfolioComments ...
func (s Server) GetUserPortfolioComments(ctx context.Context, data *userRPC.GetPortfolioComment) (*userRPC.PortfolioCommentResponse, error) {

	result, err := s.service.GetPortfolioComments(
		ctx,
		data.GetPortfolioID(),
		data.GetFirst(),
		data.GetAfter(),
	)

	if err != nil {
		return nil, err
	}

	comments := make([]*userRPC.PortfolioComment, 0, len(result.Comments))

	for _, c := range result.Comments {

		var userID = ""
		var companyID = ""

		if c.IsCompany {
			companyID = c.ProfileID.Hex()
		} else {
			userID = c.ProfileID.Hex()
		}

		comments = append(comments, &userRPC.PortfolioComment{
			ID:      c.ID.Hex(),
			Comment: c.Comment,
			// OwnerID: c.OwnerID,
			CompanyID: companyID,
			UserID:    userID,
			CreatedAt: c.CreatedAt.String(),
		})

	}

	return &userRPC.PortfolioCommentResponse{
		CommentAmount:    result.CommentAmount,
		PortfolioComment: comments,
	}, nil
}

// AddCommentToPortfolio ..
func (s Server) AddCommentToPortfolio(ctx context.Context, data *userRPC.PortfolioComment) (*userRPC.ID, error) {

	id, err := s.service.AddCommentToPortfolio(
		ctx,
		portfolioCommentRPCToComment(data),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// RemoveCommentInPortfolio ...
func (s Server) RemoveCommentInPortfolio(ctx context.Context, data *userRPC.RemovePortfolioComment) (*userRPC.Empty, error) {

	err := s.service.RemoveCommentInPortfolio(
		ctx,
		data.GetPortfolioID(),
		data.GetCommentID(),
		data.GetCompanyID(),
	)

	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeOrderFilesInPortfolio ...
func (s Server) ChangeOrderFilesInPortfolio(ctx context.Context, data *userRPC.PortfolioFile) (*userRPC.Empty, error) {

	err := s.service.ChangeOrderFilesInPortfolio(ctx, data.GetID(), data.GetFileID(), data.GetPosition())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangePortfolio ...
func (s Server) ChangePortfolio(ctx context.Context, data *userRPC.Portfolio) (*userRPC.Empty, error) {
	err := s.service.ChangePortfolio(ctx, portfolioRPCToProfilePortfolio(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemovePortfolio ...
func (s Server) RemovePortfolio(ctx context.Context, data *userRPC.Portfolio) (*userRPC.Empty, error) {
	err := s.service.RemovePortfolio(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddLinksInPortfolio ...
func (s Server) AddLinksInPortfolio(ctx context.Context, data *userRPC.AddLinksRequest) (*userRPC.Empty, error) {
	urls := make([]*profile.Link, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		urls = append(urls, linkRPCToProfileLink(data.GetLinks()[i]))
	}

	_ /*ids*/, err := s.service.AddLinksInPortfolio(ctx, data.GetID(), urls)
	if err != nil {
		return nil, err
	}

	// TODO: return ids

	return &userRPC.Empty{}, nil
}

// ChangeLinkInPortfolio ...
func (s Server) ChangeLinkInPortfolio(ctx context.Context, data *userRPC.ChangeLinkRequest) (*userRPC.Empty, error) {
	// ids := make([]string, 0, len(data.GetLinks()))
	// for i := range data.GetLinks() {
	// 	ids = append(ids, data.GetLinks()[i].GetID())
	// }

	err := s.service.ChangeLinkInPortfolio(ctx, data.GetID(), data.GetLink().GetID(), data.GetLink().GetURL())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveLinksInPortfolio ...
func (s Server) RemoveLinksInPortfolio(ctx context.Context, data *userRPC.RemoveLinksRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		ids = append(ids, data.GetLinks()[i].GetID())
	}

	err := s.service.RemoveLinksInPortfolio(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddFileInPortfolio ...
func (s Server) AddFileInPortfolio(ctx context.Context, data *userRPC.File) (*userRPC.ID, error) {

	id, err := s.service.AddFileInPortfolio(ctx, data.GetUserID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{
		ID: id,
	}, nil
}

// RemoveFilesInPortfolio ...
func (s Server) RemoveFilesInPortfolio(ctx context.Context, data *userRPC.RemoveFilesRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInPortfolio(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetToolsTechnologies ...
func (s Server) GetToolsTechnologies(ctx context.Context, data *userRPC.RequestToolTechnology) (*userRPC.ToolsTechnologiesList, error) {
	result, err := s.service.GetToolsTechnologies(
		ctx,
		data.GetUserID(),
		data.GetFirst(),
		data.GetAfter(),
		data.GetLanguage(),
	)
	if err != nil {
		return nil, err
	}

	tool := make([]*userRPC.ToolTechnology, 0, len(result))

	for i := range result {
		tool = append(tool, profileToolTechnologyToUserRPCToolTechnology(result[i]))
	}

	return &userRPC.ToolsTechnologiesList{
		ToolsTechnologies: tool,
	}, nil
}

// AddToolTechnology ...
func (s Server) AddToolTechnology(ctx context.Context, data *userRPC.ToolTechnologyList) (*userRPC.IDs, error) {

	tl := make([]*profile.ToolTechnology, 0, len(data.GetToolsTechnologies()))

	for i := range data.GetToolsTechnologies() {
		if data.GetToolsTechnologies() != nil {
			tl = append(tl, userRPCToolsToProfileToolTechnology(data.GetToolsTechnologies()[i]))
		}
	}

	stringIDs, err := s.service.AddToolTechnology(ctx, tl)
	if err != nil {
		return nil, err
	}

	return &userRPC.IDs{
		IDs: stringIDs,
	}, nil
}

// ChangeToolTechnology ...
func (s Server) ChangeToolTechnology(ctx context.Context, data *userRPC.ToolTechnologyList) (*userRPC.Empty, error) {

	tl := make([]*profile.ToolTechnology, 0, len(data.GetToolsTechnologies()))

	for i := range data.GetToolsTechnologies() {
		if data.GetToolsTechnologies() != nil {
			tl = append(tl, userRPCToolsToProfileToolTechnology(data.GetToolsTechnologies()[i]))
		}
	}

	err := s.service.ChangeToolTechnology(ctx, tl)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveToolTechnology ...
func (s *Server) RemoveToolTechnology(ctx context.Context, data *userRPC.ToolTechnologyList) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetToolsTechnologies()))
	for i := range data.GetToolsTechnologies() {
		ids = append(ids, data.GetToolsTechnologies()[i].GetID())
	}

	err := s.service.RemoveToolTechnology(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetExperiences ...
func (s Server) GetExperiences(ctx context.Context, data *userRPC.RequestExperiences) (*userRPC.ExperienceList, error) {
	result, err := s.service.GetExperiences(
		ctx,
		data.GetUserID(),
		data.GetFirst(),
		data.GetAfter(),
		data.GetLanguage(),
	)
	if err != nil {
		return nil, err
	}

	exps := make([]*userRPC.Experience, 0, len(result))

	for i := range result {
		exps = append(exps, profileExperienceToExperienceRPC(result[i]))
	}

	return &userRPC.ExperienceList{
		Experiences: exps,
	}, nil
}

// AddExperience ...
func (s Server) AddExperience(ctx context.Context, data *userRPC.Experience) (*userRPC.ID, error) {
	id, err := s.service.AddExperience(
		ctx,
		experienceRPCToProfileExperience(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// ChangeExperience ...
func (s Server) ChangeExperience(ctx context.Context, data *userRPC.Experience) (*userRPC.Empty, error) {
	err := s.service.ChangeExperience(ctx, experienceRPCToProfileExperience(data), !data.GetIsCurrentlyWorkNull())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveExperience ...
func (s Server) RemoveExperience(ctx context.Context, data *userRPC.Experience) (*userRPC.Empty, error) {
	err := s.service.RemoveExperience(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddLinksInExperience ...
func (s Server) AddLinksInExperience(ctx context.Context, data *userRPC.AddLinksRequest) (*userRPC.Empty, error) {
	urls := make([]*profile.Link, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		urls = append(urls, linkRPCToProfileLink(data.GetLinks()[i]))
	}

	_ /*ids*/, err := s.service.AddLinksInExperience(ctx, data.GetID(), urls)
	if err != nil {
		return nil, err
	}

	// TODO: return ids

	return &userRPC.Empty{}, nil
}

// AddFileInExperience ...
func (s Server) AddFileInExperience(ctx context.Context, data *userRPC.File) (*userRPC.ID, error) {
	id, err := s.service.AddFileInExperience(ctx, data.GetUserID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{
		ID: id,
	}, nil
}

// RemoveFilesInExperience ...
func (s Server) RemoveFilesInExperience(ctx context.Context, data *userRPC.RemoveFilesRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInExperience(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeLinkInExperience ...
func (s Server) ChangeLinkInExperience(ctx context.Context, data *userRPC.ChangeLinkRequest) (*userRPC.Empty, error) {
	// ids := make([]string, 0, len(data.GetLinks()))
	// for i := range data.GetLinks() {
	// 	ids = append(ids, data.GetLinks()[i].GetID())
	// }

	err := s.service.ChangeLinkInExperience(ctx, data.GetID(), data.GetLink().GetID(), data.GetLink().GetURL())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveLinksInExperience ...
func (s Server) RemoveLinksInExperience(ctx context.Context, data *userRPC.RemoveLinksRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		ids = append(ids, data.GetLinks()[i].GetID())
	}

	err := s.service.RemoveLinksInExperience(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetUploadedFilesInExperience ...
func (s Server) GetUploadedFilesInExperience(ctx context.Context, data *userRPC.Empty) (*userRPC.Files, error) {
	files, err := s.service.GetUploadedFilesInExperience(ctx)
	if err != nil {
		return nil, err
	}

	fileRPC := make([]*userRPC.File, 0, len(files))
	for i := range files {
		fileRPC = append(fileRPC, profileFiletoFileRPC(files[i]))
	}

	return &userRPC.Files{
		Files: fileRPC,
	}, nil
}

// GetEducations ...
func (s Server) GetEducations(ctx context.Context, data *userRPC.RequestEducations) (*userRPC.EducationList, error) {
	result, err := s.service.GetEducations(ctx, data.GetUserID(), data.GetFirst(), data.GetAfter(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	edus := make([]*userRPC.Education, 0, len(result))

	for i := range result {
		edus = append(edus, profileEducationToEducationRPC(result[i]))
	}

	return &userRPC.EducationList{
		Educations: edus,
	}, nil
}

// AddEducation ...
func (s Server) AddEducation(ctx context.Context, data *userRPC.Education) (*userRPC.ID, error) {
	id, err := s.service.AddEducation(
		ctx,
		educationRPCToProfileEducation(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// ChangeEducation ...
func (s Server) ChangeEducation(ctx context.Context, data *userRPC.Education) (*userRPC.Empty, error) {
	err := s.service.ChangeEducation(ctx, educationRPCToProfileEducation(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveEducation ...
func (s Server) RemoveEducation(ctx context.Context, data *userRPC.Education) (*userRPC.Empty, error) {
	err := s.service.RemoveEducation(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddLinksInEducation ...
func (s Server) AddLinksInEducation(ctx context.Context, data *userRPC.AddLinksRequest) (*userRPC.ID, error) {
	urls := make([]*profile.Link, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		urls = append(urls, linkRPCToProfileLink(data.GetLinks()[i]))
	}

	ids, err := s.service.AddLinksInEducation(ctx, data.GetID(), urls)
	if err != nil {
		return nil, err
	}

	// TODO: return several ids
	id := ""

	for _, i := range ids {
		id = i
		break
	}

	return &userRPC.ID{
		ID: id,
	}, nil
}

// AddFileInEducation ...
func (s Server) AddFileInEducation(ctx context.Context, data *userRPC.File) (*userRPC.ID, error) {
	id, err := s.service.AddFileInEducation(ctx, data.GetUserID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{
		ID: id,
	}, nil
}

// RemoveFilesInEducation ...
func (s Server) RemoveFilesInEducation(ctx context.Context, data *userRPC.RemoveFilesRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInEducation(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeLinkInEducation ...
func (s Server) ChangeLinkInEducation(ctx context.Context, data *userRPC.ChangeLinkRequest) (*userRPC.Empty, error) {

	err := s.service.ChangeLinkInEducation(ctx, data.GetID(), data.GetLink().GetID(), data.GetLink().GetURL())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveLinksInEducation ...
func (s Server) RemoveLinksInEducation(ctx context.Context, data *userRPC.RemoveLinksRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		ids = append(ids, data.GetLinks()[i].GetID())
	}

	err := s.service.RemoveLinksInEducation(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetUploadedFilesInEducation ...
func (s Server) GetUploadedFilesInEducation(ctx context.Context, data *userRPC.Empty) (*userRPC.Files, error) {
	files, err := s.service.GetUploadedFilesInEducation(ctx)
	if err != nil {
		return nil, err
	}

	fileRPC := make([]*userRPC.File, 0, len(files))
	for i := range files {
		fileRPC = append(fileRPC, profileFiletoFileRPC(files[i]))
	}

	return &userRPC.Files{
		Files: fileRPC,
	}, nil
}

// GetSkills ...
func (s *Server) GetSkills(ctx context.Context, data *userRPC.RequestSkills) (*userRPC.SkillList, error) {
	skills, err := s.service.GetSkills(ctx, data.GetUserID(), data.GetFirst(), data.GetAfter(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	sk := make([]*userRPC.Skill, 0, len(skills))
	for i := range skills {
		sk = append(sk, profileSkillToSkillRPC(skills[i]))
	}

	return &userRPC.SkillList{
		Skills: sk,
	}, nil
}

// GetEndorsements ...
func (s *Server) GetEndorsements(ctx context.Context, data *userRPC.RequestEndorsements) (*userRPC.ProfileList, error) {
	profiles, err := s.service.GetEndorsements(ctx, data.GetSkillID(), data.GetFirst(), data.GetAfter(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	pr := make([]*userRPC.Profile, 0, len(profiles))
	for i := range profiles {
		pr = append(pr, profileToProfileRPC(profiles[i]))
	}

	return &userRPC.ProfileList{
		Profiles: pr,
	}, nil
}

// AddSkills ...
func (s *Server) AddSkills(ctx context.Context, data *userRPC.SkillList) (*userRPC.ID, error) {

	sk := make([]*profile.Skill, 0, len(data.GetSkills()))

	for i := range data.GetSkills() {
		if data.GetSkills() != nil {
			sk = append(sk, skillRPCToProfileSkill(data.GetSkills()[i]))
		}
	}

	stringIDs, err := s.service.AddSkills(ctx, sk)
	if err != nil {
		return nil, err
	}

	if len(stringIDs) > 0 {
		return &userRPC.ID{
			ID: stringIDs[0],
		}, nil
	}

	return &userRPC.ID{}, nil
}

// ChangeOrderOfSkill ...
func (s *Server) ChangeOrderOfSkill(ctx context.Context, data *userRPC.Skill) (*userRPC.Empty, error) {
	err := s.service.ChangeOrderOfSkill(ctx, data.GetID(), data.GetPosition())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveSkills ...
func (s *Server) RemoveSkills(ctx context.Context, data *userRPC.SkillList) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetSkills()))
	for i := range data.GetSkills() {
		ids = append(ids, data.GetSkills()[i].GetID())
	}

	err := s.service.RemoveSkills(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// VerifySkill ...
func (s *Server) VerifySkill(ctx context.Context, data *userRPC.VerifySkillRequest) (*userRPC.Empty, error) {
	err := s.service.VerifySkill(ctx, data.GetUserID(), data.GetSkill().GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// UnverifySkill ...
func (s *Server) UnverifySkill(ctx context.Context, data *userRPC.VerifySkillRequest) (*userRPC.Empty, error) {
	err := s.service.UnverifySkill(ctx, data.GetUserID(), data.GetSkill().GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetInterests ...
func (s *Server) GetInterests(ctx context.Context, data *userRPC.RequestInterests) (*userRPC.InterestList, error) {
	interests, err := s.service.GetInterests(ctx, data.GetUserID(), data.GetFirst(), data.GetAfter(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	inter := make([]*userRPC.Interest, 0, len(interests))
	for i := range interests {
		inter = append(inter, profileInterestToInterestRPC(interests[i]))
	}

	return &userRPC.InterestList{
		Interests: inter,
	}, nil
}

// AddInterest ...
func (s *Server) AddInterest(ctx context.Context, data *userRPC.Interest) (*userRPC.ID, error) {
	id, err := s.service.AddInterest(
		ctx,
		interestRPCToProfileInterest(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// ChangeInterest ...
func (s *Server) ChangeInterest(ctx context.Context, data *userRPC.Interest) (*userRPC.Empty, error) {
	err := s.service.ChangeInterest(ctx, interestRPCToProfileInterest(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveInterest ...
func (s *Server) RemoveInterest(ctx context.Context, data *userRPC.Interest) (*userRPC.Empty, error) {
	err := s.service.RemoveInterest(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeImageInterest ...
func (s *Server) ChangeImageInterest(ctx context.Context, data *userRPC.File) (*userRPC.ID, error) {
	id, err := s.service.ChangeImageInterest(ctx, data.GetUserID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}
	return &userRPC.ID{ID: id}, nil
}

// RemoveImageInInterest ...
func (s *Server) RemoveImageInInterest(ctx context.Context, data *userRPC.Interest) (*userRPC.Empty, error) {
	err := s.service.RemoveImageInInterest(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetUnuploadImageInInterest ...
func (s *Server) GetUnuploadImageInInterest(ctx context.Context, data *userRPC.Empty) (*userRPC.File, error) {
	file, err := s.service.GetUnuploadImageInInterest(ctx)
	if err != nil {
		return nil, err
	}

	fileRPC := profileFiletoFileRPC(file)

	if fileRPC == nil {
		return &userRPC.File{}, nil
	}

	return fileRPC, nil
}

// GetOriginImageInInterest ...
func (s *Server) GetOriginImageInInterest(ctx context.Context, data *userRPC.Interest) (*userRPC.File, error) {
	image, err := s.service.GetOriginImageInInterest(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.File{
		URL: image,
	}, nil
}

// ChangeOriginImageInInterest ...
func (s *Server) ChangeOriginImageInInterest(ctx context.Context, data *userRPC.File) (*userRPC.Empty, error) {
	err := s.service.ChangeOriginImageInInterest(ctx, data.GetUserID(), data.GetTargetID(), data.GetURL())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetAccomplishments ...
func (s *Server) GetAccomplishments(ctx context.Context, data *userRPC.RequestAccomplshments) (*userRPC.AccomplishmentList, error) {
	accs, err := s.service.GetAccomplishments(ctx, data.GetUserID(), data.GetFirst(), data.GetAfter(), data.GetLanguage())
	if err != nil {
		return nil, err
	}

	accomplishments := make([]*userRPC.Accomplishment, 0, len(accs))
	for i := range accs {
		accomplishments = append(accomplishments, profileAccomplishmentToAccomplishmentRPC(accs[i]))
	}

	return &userRPC.AccomplishmentList{
		Accomplishments: accomplishments,
	}, nil
}

// AddFileInAccomplishment ...
func (s Server) AddFileInAccomplishment(ctx context.Context, data *userRPC.File) (*userRPC.ID, error) {
	id, err := s.service.AddFileInAccomplishment(ctx, data.GetUserID(), data.GetTargetID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{
		ID: id,
	}, nil
}

// RemoveFilesInAccomplishment ...
func (s Server) RemoveFilesInAccomplishment(ctx context.Context, data *userRPC.RemoveFilesRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInAccomplishment(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// AddLinksInAccomplishment ...
func (s Server) AddLinksInAccomplishment(ctx context.Context, data *userRPC.AddLinksRequest) (*userRPC.Empty, error) {
	urls := make([]*profile.Link, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		urls = append(urls, linkRPCToProfileLink(data.GetLinks()[i]))
	}

	_ /*ids*/, err := s.service.AddLinksInAccomplishment(ctx, data.GetID(), urls)
	if err != nil {
		return nil, err
	}

	// TODO: return ids

	return &userRPC.Empty{}, nil
}

// RemoveLinksInAccomplishment ...
func (s Server) RemoveLinksInAccomplishment(ctx context.Context, data *userRPC.RemoveLinksRequest) (*userRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		ids = append(ids, data.GetLinks()[i].GetID())
	}

	err := s.service.RemoveLinksInAccomplishment(ctx, data.GetID(), ids)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetUploadedFilesInAccomplishment ...
func (s Server) GetUploadedFilesInAccomplishment(ctx context.Context, data *userRPC.Empty) (*userRPC.Files, error) {
	files, err := s.service.GetUploadedFilesInAccomplishment(ctx)
	if err != nil {
		return nil, err
	}

	fileRPC := make([]*userRPC.File, 0, len(files))
	for i := range files {
		fileRPC = append(fileRPC, profileFiletoFileRPC(files[i]))
	}

	return &userRPC.Files{
		Files: fileRPC,
	}, nil
}

// AddAccomplishment ...
func (s *Server) AddAccomplishment(ctx context.Context, data *userRPC.Accomplishment) (*userRPC.ID, error) {
	id, err := s.service.AddAccomplishment(
		ctx,
		accomplishmentRPCToProfileAccomplishment(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// ChangeAccomplishment ...
func (s *Server) ChangeAccomplishment(ctx context.Context, data *userRPC.Accomplishment) (*userRPC.Empty, error) {
	err := s.service.ChangeAccomplishment(ctx, accomplishmentRPCToProfileAccomplishment(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveAccomplishment ...
func (s *Server) RemoveAccomplishment(ctx context.Context, data *userRPC.Accomplishment) (*userRPC.Empty, error) {
	err := s.service.RemoveAccomplishment(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetReceivedRecommendations ...
func (s *Server) GetReceivedRecommendations(ctx context.Context, data *userRPC.IDWithPagination) (*userRPC.Recommendations, error) {
	recoms, err := s.service.GetReceivedRecommendations(
		ctx,
		data.GetID(),
		data.GetPagination().GetFirst(),
		data.GetPagination().GetAfter(),
		data.GetLanguage(),
	)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*userRPC.Recommendation, 0, len(recoms))
	for i := range recoms {
		recommendations = append(recommendations, profileRecommendationToRecommendationRPC(recoms[i]))
	}

	return &userRPC.Recommendations{
		Recommendations: recommendations,
	}, nil
}

// GetGivenRecommendations ...
func (s *Server) GetGivenRecommendations(ctx context.Context, data *userRPC.IDWithPagination) (*userRPC.Recommendations, error) {
	recoms, err := s.service.GetGivenRecommendations(
		ctx,
		data.GetID(),
		data.GetPagination().GetFirst(),
		data.GetPagination().GetAfter(),
		data.GetLanguage(),
	)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*userRPC.Recommendation, 0, len(recoms))
	for i := range recoms {
		recommendations = append(recommendations, profileRecommendationToRecommendationRPC(recoms[i]))
	}

	return &userRPC.Recommendations{
		Recommendations: recommendations,
	}, nil
}

// GetReceivedRecommendationRequests ...
func (s *Server) GetReceivedRecommendationRequests(ctx context.Context, data *userRPC.IDWithPagination) (*userRPC.RecommendationRequests, error) {
	recomsRequest, err := s.service.GetReceivedRecommendationRequests(
		ctx,
		data.GetID(),
		data.GetPagination().GetFirst(),
		data.GetPagination().GetAfter(),
		data.GetLanguage(),
	)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*userRPC.RecommendationRequest, 0, len(recomsRequest))
	for i := range recomsRequest {
		recommendations = append(recommendations, profileRecommendationRequestToRecommendationRequestRPC(recomsRequest[i]))
	}

	return &userRPC.RecommendationRequests{
		RecommendationRequests: recommendations,
	}, nil
}

// GetRequestedRecommendationRequests ...
func (s *Server) GetRequestedRecommendationRequests(ctx context.Context, data *userRPC.IDWithPagination) (*userRPC.RecommendationRequests, error) {
	recomsRequest, err := s.service.GetRequestedRecommendationRequests(
		ctx,
		data.GetID(),
		data.GetPagination().GetFirst(),
		data.GetPagination().GetAfter(),
		data.GetLanguage(),
	)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*userRPC.RecommendationRequest, 0, len(recomsRequest))
	for i := range recomsRequest {
		recommendations = append(recommendations, profileRecommendationRequestToRecommendationRequestRPC(recomsRequest[i]))
	}

	return &userRPC.RecommendationRequests{
		RecommendationRequests: recommendations,
	}, nil
}

// GetHiddenRecommendations ...
func (s *Server) GetHiddenRecommendations(ctx context.Context, data *userRPC.IDWithPagination) (*userRPC.Recommendations, error) {
	recoms, err := s.service.GetHiddenRecommendations(
		ctx,
		data.GetID(),
		data.GetPagination().GetFirst(),
		data.GetPagination().GetAfter(),
		data.GetLanguage(),
	)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*userRPC.Recommendation, 0, len(recoms))
	for i := range recoms {
		recommendations = append(recommendations, profileRecommendationToRecommendationRPC(recoms[i]))
	}

	return &userRPC.Recommendations{
		Recommendations: recommendations,
	}, nil
}

// GetKnownLanguages ...
func (s *Server) GetKnownLanguages(ctx context.Context, data *userRPC.RequestKnownLanguages) (*userRPC.KnownLanguageList, error) {
	result, err := s.service.GetKnownLanguages(ctx, data.GetUserID(), data.GetFirst(), data.GetAfter())
	if err != nil {
		return nil, err
	}

	langs := make([]*userRPC.KnownLanguage, 0, len(result))

	for i := range result {
		langs = append(langs, profileKnownLanguageToKnownLanguageRPC(result[i]))
	}

	return &userRPC.KnownLanguageList{
		KnownLanguages: langs,
	}, nil
}

// AddKnownLanguage ...
func (s *Server) AddKnownLanguage(ctx context.Context, data *userRPC.KnownLanguage) (*userRPC.ID, error) {
	id, err := s.service.AddKnownLanguage(
		ctx,
		knownLanguageRPCToProfileKnownLanguage(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.ID{ID: id}, nil
}

// ChangeKnownLanguage ...
func (s *Server) ChangeKnownLanguage(ctx context.Context, data *userRPC.KnownLanguage) (*userRPC.Empty, error) {
	err := s.service.ChangeKnownLanguage(ctx, knownLanguageRPCToProfileKnownLanguage(data))
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveKnownLanguage ...
func (s *Server) RemoveKnownLanguage(ctx context.Context, data *userRPC.KnownLanguage) (*userRPC.Empty, error) {
	err := s.service.RemoveKnownLanguage(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeHeadline ...
func (s *Server) ChangeHeadline(ctx context.Context, data *userRPC.Headline) (*userRPC.Empty, error) {
	err := s.service.ChangeHeadline(ctx, data.GetHeadline())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeStory ...
func (s *Server) ChangeStory(ctx context.Context, data *userRPC.Story) (*userRPC.Empty, error) {
	err := s.service.ChangeStory(ctx, data.GetStory())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetOriginAvatar ...
func (s *Server) GetOriginAvatar(ctx context.Context, data *userRPC.Empty) (*userRPC.File, error) {
	avatar, err := s.service.GetOriginAvatar(ctx)
	if err != nil {
		return nil, err
	}

	return &userRPC.File{
		URL: avatar,
	}, nil
}

// ChangeOriginAvatar ...
func (s *Server) ChangeOriginAvatar(ctx context.Context, data *userRPC.File) (*userRPC.Empty, error) {
	err := s.service.ChangeOriginAvatar(ctx, data.GetUserID(), data.GetURL())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ChangeAvatar ...
func (s *Server) ChangeAvatar(ctx context.Context, data *userRPC.File) (*userRPC.Empty, error) {
	err := s.service.ChangeAvatar(ctx, data.GetUserID(), data.GetURL())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveAvatar ...
func (s *Server) RemoveAvatar(ctx context.Context, data *userRPC.Empty) (*userRPC.Empty, error) {
	err := s.service.RemoveAvatar(ctx)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// CheckPassword ...
func (s *Server) CheckPassword(ctx context.Context, data *userRPC.CheckPasswordRequest) (*userRPC.Empty, error) {
	err := s.service.CheckPassword(ctx, data.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// ReportUser ...
func (s *Server) ReportUser(ctx context.Context, data *userRPC.ReportUserRequest) (*userRPC.Empty, error) {

	err := s.service.ReportUser(
		ctx,
		reportRPCToUserReport(data),
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil

}

// SaveUserProfileTranslation ...
func (s *Server) SaveUserProfileTranslation(ctx context.Context, data *userRPC.ProfileTranslation) (*userRPC.Empty, error) {
	err := s.service.SaveUserProfileTranslation(
		ctx,
		data.GetLanguage(),
		&profile.Translation{
			FirstName: data.GetFirstname(),
			Lastname:  data.GetLastname(),
			Headline:  data.GetHeadline(),
			Nickname:  data.GetNickname(),
			Story:     data.GetStory(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SaveUserExperienceTranslation ...
func (s *Server) SaveUserExperienceTranslation(ctx context.Context, data *userRPC.ExperienceTranslation) (*userRPC.Empty, error) {
	desc := data.GetDescription()
	err := s.service.SaveUserExperienceTranslation(
		ctx,
		data.GetExperienceID(),
		data.GetLanguage(),
		&profile.ExperienceTranslation{
			Company:     data.GetCompany(),
			Description: &desc,
			Position:    data.GetPosition(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SaveUserEducationTranslation ...
func (s *Server) SaveUserEducationTranslation(ctx context.Context, data *userRPC.EducationTranslation) (*userRPC.Empty, error) {
	desc := data.GetDescription()
	degree := data.GetDegree()
	grade := data.GetGrade()

	err := s.service.SaveUserEducationTranslation(
		ctx,
		data.GetEducationID(),
		data.GetLanguage(),
		&profile.EducationTranslation{
			Description: &desc,
			Degree:      &degree,
			FieldStudy:  data.GetFieldStudy(),
			School:      data.GetSchool(),
			Grade:       &grade,
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SaveUserInterestTranslation ...
func (s *Server) SaveUserInterestTranslation(ctx context.Context, data *userRPC.InterestTranslation) (*userRPC.Empty, error) {
	desc := data.GetDescription()

	err := s.service.SaveUserInterestTranslation(
		ctx,
		data.GetInterestID(),
		data.GetLanguage(),
		&profile.InterestTranslation{
			Interest:    data.GetInterest(),
			Description: &desc,
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SaveUserSkillTranslation ...
func (s *Server) SaveUserSkillTranslation(ctx context.Context, data *userRPC.SkillTranslation) (*userRPC.Empty, error) {
	err := s.service.SaveUserSkillTranslation(
		ctx,
		data.GetSkillID(),
		data.GetLanguage(),
		&profile.SkillTranslation{
			Skill: data.GetSkill(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SaveUserAccomplishmentTranslation ...
func (s *Server) SaveUserAccomplishmentTranslation(ctx context.Context, data *userRPC.AccomplishmentTranslation) (*userRPC.Empty, error) {
	desc := data.GetDescription()
	issuer := data.GetIssuer()

	err := s.service.SaveUserAccomplishmentTranslation(
		ctx,
		data.GetAccomplishmentID(),
		data.GetLanguage(),
		&profile.AccomplishmentTranslation{
			Name:        data.GetName(),
			Description: &desc,
			Issuer:      &issuer,
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SentEmailInvitation ...
func (s *Server) SentEmailInvitation(ctx context.Context, data *userRPC.EmailInvitation) (*userRPC.Empty, error) {
	err := s.service.SentEmailInvitation(ctx, data.GetAddress(), data.GetName(), data.GetCompanyID())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SaveUserPortfolioTranslation ...
func (s *Server) SaveUserPortfolioTranslation(ctx context.Context, data *userRPC.PortfolioTranslation) (*userRPC.Empty, error) {

	err := s.service.SaveUserPortfolioTranslation(
		ctx,
		data.GetPortfolioID(),
		data.GetLanguage(),
		&profile.PortfolioTranslation{
			Description: data.GetDescription(),
			Tittle:      data.GetTitle(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// SaveUserToolTechnologyTranslation ...
func (s *Server) SaveUserToolTechnologyTranslation(ctx context.Context, data *userRPC.ToolTechnologyTranslation) (*userRPC.Empty, error) {

	err := s.service.SaveUserToolTechnologyTranslation(
		ctx,
		data.GetTooltechnologyID(),
		data.GetLanguage(),
		&profile.ToolTechnologyTranslation{
			ToolTechnology: data.GetToolTechnology(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// RemoveTranslation ...
func (s *Server) RemoveTranslation(ctx context.Context, data *userRPC.Language) (*userRPC.Empty, error) {
	err := s.service.RemoveTransaltion(ctx, data.GetLanguage())
	if err != nil {
		return nil, err
	}

	return &userRPC.Empty{}, nil
}

// GetInvitation ...
func (s *Server) GetInvitation(ctx context.Context, data *userRPC.Empty) (*userRPC.Invitations, error) {
	res, amount, err := s.service.GetInvitation(ctx)
	if err != nil {
		return nil, err
	}

	invs := make([]*userRPC.Invitation, 0, len(res))

	for _, inv := range res {
		invs = append(invs, invitaionInvitationToInvitaionRPC(&inv))
	}

	return &userRPC.Invitations{
		Amount:      amount,
		Invitations: invs,
	}, nil
}

// GetInvitationForCompany ...
func (s *Server) GetInvitationForCompany(ctx context.Context, data *userRPC.ID) (*userRPC.Invitations, error) {
	res, amount, err := s.service.GetInvitationForCompany(
		ctx,
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}

	invs := make([]*userRPC.Invitation, 0, len(res))

	for _, inv := range res {
		invs = append(invs, invitaionInvitationToInvitaionRPC(&inv))
	}

	return &userRPC.Invitations{
		Amount:      amount,
		Invitations: invs,
	}, nil
}

// GetConectionsPrivacy ...
func (s *Server) GetConectionsPrivacy(ctx context.Context, data *userRPC.ID) (*userRPC.Permission, error) {
	perm, err := s.service.GetConectionsPrivacy(ctx, data.GetID())
	if err != nil {
		return nil, err
	}
	p := accountPermissionTypeToPermissionTypeRPC(&perm)

	pe := userRPC.PermissionType_NONE

	if p != nil {
		pe = *p
	}

	return &userRPC.Permission{
		Type: pe,
	}, nil
}

// GetUsersForAdvert ...
func (s *Server) GetUsersForAdvert(ctx context.Context, data *userRPC.GetUsersForAdvertRequest) (*userRPC.IDs, error) {

	ids, err := s.service.GetUsersForAdvert(ctx, account.UserForAdvert{
		OwnerID:   data.GetOwnerID(),
		AgeFrom:   data.GetAgeFrom(),
		AgeTo:     data.GetAgeTo(),
		Gender:    data.GetGender(),
		Locations: data.GetLocations(),
		Languages: data.GetLanguages(),
	})
	if err != nil {
		return nil, err
	}

	return &userRPC.IDs{
		IDs: ids,
	}, nil

}

// Converters to service types

func registerRequestToAccount(data *userRPC.RegisterRequest) *account.Account {
	acc := account.Account{
		FirstName:  data.GetFirstName(),
		Username:   data.GetUsername(),
		Lastname:   data.GetLastName(),
		LanguageID: data.GetLanguageId(),
		Gender:     account.Gender{Gender: data.Gender.String()},
		Emails:     make([]*account.Email, 0),
		// Phones:     make([]*account.Phone, 0),
		Birthday: &account.Birthday{},
	}

	if data.GetInvitedBy() != "" {
		acc.SetInvitedByID(data.GetInvitedBy())
	}

	acc.Emails = append(acc.Emails, &account.Email{
		Email: data.GetEmail(),
	})

	// acc.Phones = append(acc.Phones, &account.Phone{
	// 	CountryCode: account.CountryCode{
	// 		ID: uint32(data.GetCountryPrefixCode()),
	// 	},
	// 	Number: data.GetPhoneNumber(),
	// })

	acc.Birthday.Birthday = stringDateToTime(data.GetBirthday())

	acc.Location = &account.UserLocation{
		Location: location.Location{
			City: &location.City{
				ID: data.GetCityID(),
			},
			Country: &location.Country{
				ID: data.GetCountryId(),
			},
		},
	}

	// data.GetPassword()

	return &acc
}

func genderRPCToAccountGender(data *userRPC.Gender) account.Gender {
	perm := account.Permission{}
	if data.GetPermission() != nil {
		t := permissionRPCToAccountPermission(data.GetPermission())
		perm = *t
	}
	return account.Gender{
		Gender:     data.GetGender().String(),
		Permission: perm,
	}
}

func permissionRPCToAccountPermission(data *userRPC.Permission) *account.Permission {
	if data == nil {
		return nil
	}

	perm := account.Permission{}

	tmp := data.GetType()

	permType := permissionTypeRPCToAccountPermissionType(&tmp)
	if permType != nil {
		perm.Type = *permType
	}

	// switch data.GetType() {
	// case userRPC.PermissionType_NONE:
	// 	log.Println("Experimental: got PermissionType_NONE")
	// 	// perm.Type = account.PermissionTypeNone
	// 	return nil
	// case userRPC.PermissionType_ME:
	// 	perm.Type = account.PermissionTypeMe
	// case userRPC.PermissionType_MEMBERS:
	// 	perm.Type = account.PermissionTypeMembers
	// case userRPC.PermissionType_MY_CONNECTIONS:
	// 	perm.Type = account.PermissionTypeMyConnections
	// default:
	// 	return nil
	// }

	return &perm
}

func permissionTypeRPCToAccountPermissionType(data *userRPC.PermissionType) *account.PermissionType {
	if data == nil {
		return nil
	}

	switch *data {
	case userRPC.PermissionType_NONE:
		// perm.Type = account.PermissionTypeNone
		return nil
	case userRPC.PermissionType_ME:
		p := account.PermissionTypeMe
		return &p
	case userRPC.PermissionType_MEMBERS:
		p := account.PermissionTypeMembers
		return &p
	case userRPC.PermissionType_MY_CONNECTIONS:
		p := account.PermissionTypeMyConnections
		return &p
	}

	return nil
}

func countryCodeRPCtoAccountCountryCode(data *userRPC.CountryCode) *account.CountryCode {
	if data == nil {
		return nil
	}

	return &account.CountryCode{
		ID:   data.GetID(),
		Code: data.GetCode(),
		// CountryID: data.
	}
}

func myAddressRPCToAccountMyAddress(data *userRPC.MyAddress) *account.MyAddress {
	if data == nil {
		return nil
	}

	addr := account.MyAddress{
		Name:      data.GetName(),
		Firstname: data.GetFirstname(),
		Lastname:  data.GetLastname(),
		Apartment: data.GetApartment(),
		Street:    data.GetStreet(),
		ZIP:       data.GetZIP(),
		IsPrimary: data.GetIsPrimary(),
	}

	err := addr.SetID(data.ID)
	if err != nil {
		log.Println("error: myAddressRPCToAccountMyAddress:", err)
	}

	if loc := locationRPCToLocationLocation(data.GetLocation()); loc != nil {
		addr.Location = *loc
	}

	return &addr
}

func locationRPCToLocationLocation(data *userRPC.Location) *location.Location {
	if data == nil {
		return nil
	}

	loc := location.Location{}

	if city := cityRPCToLocationCity(data.GetCity()); city != nil {
		loc.City = city
	}

	if country := countryRPCToLocationCountry(data.GetCountry()); country != nil {
		loc.Country = country
	}

	return &loc
}

func cityRPCToLocationCity(data *userRPC.City) *location.City {
	if data == nil {
		return nil
	}

	city := location.City{
		ID:          data.GetId(),
		Name:        data.GetTitle(),
		Subdivision: data.GetSubdivision(),
	}

	return &city
}

func countryRPCToLocationCountry(data *userRPC.Country) *location.Country {
	if data == nil {
		return nil
	}

	country := location.Country{
		ID: data.GetId(),
		// Name: data.
	}

	return &country
}

func otherAddressRPCToAccountMyAddress(data *userRPC.OtherAddress) *account.OtherAddress {
	if data == nil {
		return nil
	}

	addr := account.OtherAddress{
		Name:      data.GetName(),
		Firstname: data.GetFirstname(),
		Lastname:  data.GetLastname(),
		Apartment: data.GetApartment(),
		Street:    data.GetStreet(),
		ZIP:       data.GetZIP(),
	}

	err := addr.SetID(data.ID)
	if err != nil {
		log.Println("error: myAddressRPCToAccountMyAddress:", err)
	}

	if loc := locationRPCToLocationLocation(data.GetLocation()); loc != nil {
		addr.Location = *loc
	}

	return &addr
}

func privacySettingsRPCToAccountPrivacyItem(data *userRPC.PrivacySettings) *account.PrivacyItem {
	if data == nil {
		return nil
	}

	switch *data {
	case userRPC.PrivacySettings_sharing_edits:
		p := account.PrivacyItemShareEdits
		return &p
	case userRPC.PrivacySettings_active_status:
		p := account.PrivacyItemActiveStatus
		return &p
	case userRPC.PrivacySettings_find_by_email:
		p := account.PrivacyItemFindByEmail
		return &p
	case userRPC.PrivacySettings_find_by_phone:
		p := account.PrivacyItemFindByPhone
		return &p
	case userRPC.PrivacySettings_my_connections:
		p := account.PrivacyItemMyConnections
		return &p
		// case userRPC.PrivacySettings_profile_pictures:
		// 	itm =
	}

	p := account.PrivacyItemProfilePicture
	return &p
}

func experienceRPCToProfileExperience(data *userRPC.Experience) *profile.Experience {
	if data == nil {
		return nil
	}

	exp := profile.Experience{
		Company:   data.GetCompany(),
		Position:  data.Position,
		StartDate: stringDayMonthAndYearToTime(data.GetStartDate()),
		Links:     make([]*profile.Link, 0, len(data.GetLinks())),
		Files:     make([]*profile.File, 0, len(data.GetFiles())),
	}

	if !data.IsLocationNull {
		// c := data.GetCityID()
		// exp.CityID = &c
		if loc := locationRPCToLocationLocation(data.GetLocation()); loc != nil {
			exp.Location = *loc
		}
	}

	if !data.IsDescriptionNull {
		d := data.GetDescription()
		exp.Description = &d
	}

	if data.GetCurrentlyWork() {
		exp.CurrentlyWork = true
		t := time.Time{}
		exp.FinishDate = &t
	} else {
		t := stringDateToTime(data.GetFinishDate())
		exp.FinishDate = &t
	}

	if !data.IsCurrentlyWorkNull || data.GetFinishDate() != "" {
		f := stringDayMonthAndYearToTime(data.GetFinishDate())
		exp.FinishDate = &f
	}

	_ = exp.SetID(data.GetID())

	for i := range data.GetFiles() {
		exp.Files = append(exp.Files, fileRPCToProfileFile(data.Files[i]))
	}

	for i := range data.GetLinks() {
		exp.Links = append(exp.Links, linkRPCToProfileLink(data.Links[i]))
	}

	return &exp
}

func portfolioRPCToProfilePortfolio(data *userRPC.Portfolio) *profile.Portfolio {
	if data == nil {
		return nil
	}

	port := profile.Portfolio{
		Title:           data.GetTittle(),
		Description:     data.GetDescription(),
		IsCommentClosed: data.GetIsCommentDisabled(),
		Tools:           make([]*string, 0, len(data.GetTools())),
		ContentType:     userRPCContentTypeToProfileContentType(data.GetContentType()),
	}

	_ = port.SetID(data.GetID())

	for i := range data.GetTools() {
		port.Tools = append(port.Tools, &data.Tools[i])
	}

	log.Printf("Tools %v", data.GetTools())

	if len(data.GetFiles()) > 0 {
		port.Files = make([]*profile.File, 0, len(data.GetFiles()))
		for i := range data.GetFiles() {
			port.Files = append(port.Files, fileRPCToProfileFile(data.Files[i]))
		}
	}

	return &port
}

func portfolioCommentRPCToComment(data *userRPC.PortfolioComment) *profile.PortfolioComment {
	if data == nil {
		return nil
	}

	comment := profile.PortfolioComment{
		Comment:     data.GetComment(),
		OwnerID:     data.GetOwnerID(),
		PortfolioID: data.GetPortfolioID(),
	}

	_ = comment.SetCommentID(data.GetID())

	comment.IsCompany = false

	if data.GetCompanyID() != "" {
		comment.ProfileID = data.GetCompanyID()
		comment.IsCompany = true
	}

	return &comment
}

func portfolioPhotoRPCToPortfolioPhoto(data *userRPC.Portfolio) *profile.Portfolio {
	if data == nil {
		return nil
	}

	port := profile.Portfolio{
		Description: data.GetDescription(),
		Files:       make([]*profile.File, 0, len(data.GetFiles())),
	}

	_ = port.SetID(data.GetID())

	for i := range data.GetFiles() {
		port.Files = append(port.Files, fileRPCToProfileFile(data.Files[i]))
	}

	return &port
}

func userRPCToolTechnologuToProfileToolTechnology(data *userRPC.ToolTechnology) *profile.ToolTechnology {
	if data == nil {
		return nil
	}

	tool := profile.ToolTechnology{
		ToolTechnology: data.ToolTechnology,
		Rank:           userRPCToolTechnologyToProfileToolTechnology(data.Rank),
	}

	_ = tool.SetID(data.GetID())

	if data.ToolTechnology != "" {
		tool.ToolTechnology = data.GetToolTechnology()
	}
	tool.Rank = userRPCToolTechnologyToProfileToolTechnology(data.GetRank())

	return &tool
}

func fileRPCToProfileFile(data *userRPC.File) *profile.File {
	if data == nil {
		return nil
	}

	f := profile.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		URL:      data.GetURL(),
	}

	_ = f.SetID(data.GetID())

	return &f
}

func linkRPCToProfileLink(data *userRPC.Link) *profile.Link {
	if data == nil {
		return nil
	}

	l := profile.Link{
		URL: data.GetURL(),
	}

	_ = l.SetID(data.GetID())

	return &l
}

func educationRPCToProfileEducation(data *userRPC.Education) *profile.Education {
	if data == nil {
		return nil
	}

	edu := profile.Education{
		FieldStudy:       data.GetFieldStudy(),
		School:           data.GetSchool(),
		StartDate:        stringDayMonthAndYearToTime(data.GetStartDate()),
		FinishDate:       stringDayMonthAndYearToTime(data.GetFinishDate()),
		IsCurrentlyStudy: data.GetIsCurrentlyStudy(),
		Links:            make([]*profile.Link, 0, len(data.GetLinks())),
		Files:            make([]*profile.File, 0, len(data.GetFiles())),
	}

	if !data.IsDegreeNull {
		d := data.GetDegree()
		edu.Degree = &d
	}

	if !data.IsGradeNull {
		g := data.GetGrade()
		edu.Grade = &g
	}

	if !data.IsDescriptionNull {
		d := data.GetDescription()
		edu.Description = &d
	}

	_ = edu.SetID(data.GetID())

	for i := range data.GetFiles() {
		edu.Files = append(edu.Files, fileRPCToProfileFile(data.Files[i]))
	}

	for i := range data.GetLinks() {
		edu.Links = append(edu.Links, linkRPCToProfileLink(data.Links[i]))
	}

	if !data.IsLocationNull {
		// c := data.GetCityID()
		// exp.CityID = &c
		if loc := locationRPCToLocationLocation(data.GetLocation()); loc != nil {
			edu.Location = *loc
		}
	}

	return &edu
}

func skillRPCToProfileSkill(data *userRPC.Skill) *profile.Skill {
	if data == nil {
		return nil
	}

	sk := profile.Skill{
		Position: data.GetPosition(),
		Skill:    data.GetSkill(),
	}

	_ = sk.SetID(data.GetID())

	return &sk
}

func interestRPCToProfileInterest(data *userRPC.Interest) *profile.Interest {
	if data == nil {
		return nil
	}

	in := profile.Interest{
		// Image: data.GetImage(),
		Interest: data.GetInterest(),
		// Description: data.GetDescription(),
	}

	if !data.IsDescriptionNull {
		s := data.GetDescription()
		in.Description = &s
	}

	_ = in.SetID(data.GetID())

	return &in
}

func accomplishmentRPCToProfileAccomplishment(data *userRPC.Accomplishment) *profile.Accomplishment {
	if data == nil {
		return nil
	}

	acc := profile.Accomplishment{
		Name: data.GetName(),
	}

	_ = acc.SetID(data.GetID())

	t := data.GetType()
	acc.Type = accomplishmentTypeRPCToProfileAccomplishmentType(t)

	if !data.IsURLNull {
		url := data.GetURL()
		acc.URL = &url
	}

	if !data.IsIssuerNull {
		iss := data.GetIssuer()
		acc.Issuer = &iss
	}

	if !data.IsIsExpireNull {
		exp := data.IsExpire
		acc.IsExpire = &exp
	}

	if !data.IsFinishDateNull {
		finishDate := stringDayMonthAndYearToTime(data.GetFinishDate())
		acc.FinishDate = &finishDate
	}

	if !data.IsDescriptionNull {
		desc := data.Description
		acc.Description = &desc
	}

	if !data.IsLicenseNumberNull {
		licNum := data.LicenseNumber
		acc.LicenseNumber = &licNum
	}

	if !data.IsScoreNull {
		score := data.Score
		acc.Score = &score
	}

	if !data.IsStartDateNull {
		startDate := stringDayMonthAndYearToTime(data.GetStartDate())
		acc.StartDate = &startDate
	}

	for i := range data.GetFiles() {
		acc.Files = append(acc.Files, fileRPCToProfileFile(data.Files[i]))
	}

	for i := range data.GetLinks() {
		acc.Links = append(acc.Links, linkRPCToProfileLink(data.Links[i]))
	}

	log.Println("HANDLER:", acc.StartDate, acc.FinishDate)

	return &acc
}

func accomplishmentTypeRPCToProfileAccomplishmentType(data userRPC.Accomplishment_AccomplishmentType) profile.AccomplishmentType {

	var accomp profile.AccomplishmentType

	switch data {
	case userRPC.Accomplishment_Certificate:
		accomp = profile.AccomplishmentTypeCertificate
	case userRPC.Accomplishment_Test:
		accomp = profile.AccomplishmentTypeTest
	case userRPC.Accomplishment_Award:
		accomp = profile.AccomplishmentTypeAward
	case userRPC.Accomplishment_License:
		accomp = profile.AccomplishmentTypeLicense
	case userRPC.Accomplishment_Project:
		accomp = profile.AccomplishmentTypeProject
	case userRPC.Accomplishment_Publication:
		accomp = profile.AccomplishmentTypePublication
	}

	return accomp
}

func knownLanguageRPCToProfileKnownLanguage(data *userRPC.KnownLanguage) *profile.KnownLanguage {
	if data == nil {
		return nil
	}

	lang := profile.KnownLanguage{
		Language: data.GetLanguage(),
		Rank:     data.GetRank(),
	}

	lang.SetID(data.GetID())

	return &lang
}

func reportRPCToUserReport(data *userRPC.ReportUserRequest) *userReport.Report {
	if data == nil {
		return nil
	}

	rep := userReport.Report{
		Description: data.GetDescription(),
	}

	rep.SetUserID(data.GetUserID())

	t := data.GetType()
	rep.Type = *reportTypeToUserReportType(&t)

	return &rep
}

func reportTypeToUserReportType(data *userRPC.ReportUserRequest_ReportType) *userReport.Type {
	if data == nil {
		return nil
	}

	var t userReport.Type

	switch *data {
	case userRPC.ReportUserRequest_MayBeHacked:
		t = userReport.ReportUserRequestMayBeHacked
	case userRPC.ReportUserRequest_NotRealIndividual:
		t = userReport.ReportUserRequestNotRealIndividual
	case userRPC.ReportUserRequest_PictureIsNotPerson:
		t = userReport.ReportUserRequestPictureIsNotPerson
	case userRPC.ReportUserRequest_PictureIsOffensive:
		t = userReport.ReportUserRequestPictureIsOffensive
	case userRPC.ReportUserRequest_PretendingToBeSomone:
		t = userReport.ReportUserRequestPretendingToBeSomone
	case userRPC.ReportUserRequest_VolatationTermsOfUse:
		t = userReport.ReportUserRequestVolatationTermsOfUse
	default:
		t = userReport.ReportUserRequestOther
	}

	return &t
}

// Converters to gRPC types

// TODO:
func accountToRPC(data *account.Account) *userRPC.Account {
	if data == nil {
		return nil
	}

	acc := &userRPC.Account{
		FirstName:          data.FirstName,
		Lastname:           data.Lastname,
		LanguageID:         data.LanguageID,
		Emails:             make([]*userRPC.Email, 0, len(data.Emails)),
		Phones:             make([]*userRPC.Phone, 0, len(data.Phones)),
		Patronymic:         accountPatronymicToRPC(data.Patronymic),
		NickName:           accountNicknameToRPC(data.Nickname),
		NativeName:         accountNativeNameToRPC(data.NativeName),
		MiddleName:         accountMiddleNameToRPC(data.MiddleName),
		MyAddresses:        make([]*userRPC.MyAddress, 0, len(data.MyAddresses)),
		OtherAddresses:     make([]*userRPC.OtherAddress, 0, len(data.OtherAddresses)),
		Location:           accountUserLocationToRPC(data.Location),
		Gender:             genderToRPC(data.Gender),
		Birthday:           accountBirthdayToRPC(data.Birthday),
		Privacies:          accountPrivaciesToRPC(data.Privacy),
		Notifications:      accountNotificationsToRPC(data.Notification),
		LastChangePassword: timeToStringDayMonthAndYear(data.LastChangePassword),
		AmountOfSessions:   data.AmountOfSessions,

		// DateOfActivation: timeToStringDayMonthAndYear(data.CreatedAt),
		// TODO:
		// :data.Status,
		// IsEditable: ,
	}

	for _, email := range data.Emails {
		acc.Emails = append(acc.Emails, accountEmailToRPC(email))
	}

	for _, phone := range data.Phones {
		acc.Phones = append(acc.Phones, accountPhoneToRPC(phone))
	}

	for _, address := range data.MyAddresses {
		acc.MyAddresses = append(acc.MyAddresses, accountMyAddressToRPC(address))
	}

	for _, address := range data.OtherAddresses {
		acc.OtherAddresses = append(acc.OtherAddresses, accountOtherAddressToRPC(address))
	}

	return acc
}

func accountBirthdayToRPC(data *account.Birthday) *userRPC.Birthday {
	y, m, d := data.Birthday.UTC().Date()
	bd := strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)

	return &userRPC.Birthday{
		Birthday:   bd,
		Permission: accountPermissionToRPC(&data.Permission),
	}
}

func accountPermissionToRPC(data *account.Permission) *userRPC.Permission {
	perm := &userRPC.Permission{}
	if data == nil {
		perm.Type = userRPC.PermissionType_MEMBERS
		return perm
	}
	perm.Type = accountPermissionTypeToRPC(data.Type)

	return perm
}

func accountPermissionTypeToRPC(data account.PermissionType) userRPC.PermissionType {
	var perm userRPC.PermissionType
	switch data {
	case account.PermissionTypeMe:
		perm = userRPC.PermissionType_ME
	case account.PermissionTypeMembers:
		perm = userRPC.PermissionType_MEMBERS
	case account.PermissionTypeMyConnections:
		perm = userRPC.PermissionType_MY_CONNECTIONS
	default:
		perm = userRPC.PermissionType_MEMBERS
	}

	return perm
}

func accountPatronymicToRPC(data *account.Patronymic) *userRPC.Patronymic {
	if data == nil {
		return nil
	}

	return &userRPC.Patronymic{
		Patronymic: data.Patronymic,
		Permission: accountPermissionToRPC(data.Permission),
	}
}

func accountEmailToRPC(data *account.Email) *userRPC.Email {
	if data == nil {
		return nil
	}

	return &userRPC.Email{
		Id:          data.GetID(),
		Email:       data.Email,
		IsActivated: data.Activated,
		IsPrimary:   data.Primary,
		Permission:  accountPermissionToRPC(&data.Permission),
	}
}

func accountPhoneToRPC(data *account.Phone) *userRPC.Phone {
	if data == nil {
		return nil
	}

	return &userRPC.Phone{
		// CountryAbbreviation: data.CountryAbbreviation,
		Id:          data.GetID(),
		IsActivated: data.Activated,
		IsPrimary:   data.Primary,
		Number:      data.Number,
		Permission:  accountPermissionToRPC(&data.Permission),
		CountryCode: accountCountryCodeToRPC(&data.CountryCode),
	}
}

func accountCountryCodeToRPC(data *account.CountryCode) *userRPC.CountryCode {
	if data == nil {
		return nil
	}

	return &userRPC.CountryCode{
		Code: data.Code,
		ID:   data.ID,
		// : data.CountryID,
	}
}

func accountMiddleNameToRPC(data *account.MiddleName) *userRPC.Middlename {
	if data == nil {
		return nil
	}

	return &userRPC.Middlename{
		Middlename: data.Middlename,
		Permission: accountPermissionToRPC(data.Permission),
	}
}

func accountNicknameToRPC(data *account.Nickname) *userRPC.Nickname {
	if data == nil {
		return nil
	}

	return &userRPC.Nickname{
		Nickname:   data.Nickname,
		Permission: accountPermissionToRPC(data.Permission),
	}
}

func accountNativeNameToRPC(data *account.NativeName) *userRPC.NativeName {
	if data == nil {
		return nil
	}

	return &userRPC.NativeName{
		LanguageID: data.LanguageID,
		Name:       data.Name,
		Permission: accountPermissionToRPC(data.Permission),
	}
}

func accountUserLocationToRPC(data *account.UserLocation) *userRPC.Location {
	if data == nil {
		return nil
	}

	return &userRPC.Location{
		Permission: accountPermissionToRPC(&data.Permission),
		City:       locationCityToRPC(data.City),
		Country:    locationCountryToRPC(data.Country),
	}
}

func locationCityToRPC(data *location.City) *userRPC.City {
	if data == nil {
		return nil
	}

	return &userRPC.City{
		Id:          data.ID,
		Subdivision: data.Subdivision,
		Title:       data.Name,
	}
}

func locationCountryToRPC(data *location.Country) *userRPC.Country {
	if data == nil {
		return nil
	}

	return &userRPC.Country{
		Id: data.ID,
		// : data.Name,
	}
}

func locationToRPC(data *location.Location) *userRPC.Location {
	if data == nil {
		return nil
	}

	return &userRPC.Location{
		City:    locationCityToRPC(data.City),
		Country: locationCountryToRPC(data.Country),
		// Permission
	}
}

func accountLocationToProfileLocationRPC(data *account.UserLocation) *userRPC.LocationProfile {
	if data == nil {
		return nil
	}

	l := userRPC.LocationProfile{}

	if data.Country != nil {
		l.CountryID = data.Country.ID
	}

	if data.City != nil {
		l.City = data.City.Name
	}

	return &l
}

func accountMyAddressToRPC(data *account.MyAddress) *userRPC.MyAddress {
	if data == nil {
		return nil
	}

	return &userRPC.MyAddress{
		Apartment: data.Apartment,
		Firstname: data.Firstname,
		ID:        data.GetID(),
		IsPrimary: data.IsPrimary,
		Lastname:  data.Lastname,
		Location:  locationToRPC(&data.Location),
		Name:      data.Name,
		Street:    data.Street,
		ZIP:       data.ZIP,
	}
}

func accountOtherAddressToRPC(data *account.OtherAddress) *userRPC.OtherAddress {
	if data == nil {
		return nil
	}

	return &userRPC.OtherAddress{
		Apartment: data.Apartment,
		Firstname: data.Firstname,
		ID:        data.GetID(),
		Lastname:  data.Lastname,
		Location:  locationToRPC(&data.Location),
		Name:      data.Name,
		Street:    data.Street,
		ZIP:       data.ZIP,
	}
}

func genderToRPC(data account.Gender) *userRPC.Gender {
	gender := userRPC.Gender{}

	if data.Gender == "FEMALE" {
		gender.Gender = userRPC.GenderValue_FEMALE
	} else {
		gender.Gender = userRPC.GenderValue_MALE
	}

	perm := data.Permission
	gender.Permission = accountPermissionToRPC(&perm)

	return &gender
}

func accountPrivaciesToRPC(data *account.Privacy) *userRPC.Privacies {
	if data == nil {
		return nil
	}

	return &userRPC.Privacies{
		ActiveStatus:   accountPermissionTypeToRPC(data.ActiveStatus),
		FindByEmail:    accountPermissionTypeToRPC(data.FindByEmail),
		FindByPhone:    accountPermissionTypeToRPC(data.FindByPhone),
		MyConnections:  accountPermissionTypeToRPC(data.MyConnections),
		ProfilePicture: accountPermissionTypeToRPC(data.ProfilePicture),
		ShareEdits:     accountPermissionTypeToRPC(data.ShareEdits),
	}
}

func accountNotificationsToRPC(data *account.Notifications) *userRPC.Notifications {
	if data == nil {
		return nil
	}

	return &userRPC.Notifications{
		AcceptInvitation:     data.AcceptInvitation,
		Birthdays:            data.Birthdays,
		ConnectionRequest:    data.ConnectionRequest,
		EmailUpdates:         data.Endorsements,
		Endorsements:         data.Endorsements,
		ImportContactsJoined: data.ImportContactsJoined,
		JobChangesInNetwork:  data.JobChangesInNetwork,
		JobRecommendations:   data.JobRecommendations,
		NewChatMessage:       data.NewChatMessage,
		NewFollowers:         data.NewFollowers,
	}
}

func profileToProfileRPC(data *profile.Profile) *userRPC.Profile {
	if data == nil {
		return nil
	}

	prof := &userRPC.Profile{
		ID:     data.GetID(),
		Avatar: data.Avatar,
		URL:    data.URL,

		Location: accountLocationToProfileLocationRPC(data.Location),
		// NativeName: "",

		Headline:  data.Headline,
		Firstname: data.FirstName,
		Lastname:  data.Lastname,
		Story:     data.Story,

		IsMe:                    data.IsMe,
		IsFollow:                data.IsFollow,
		IsFriend:                data.IsFriend,
		IsOnline:                data.IsOnline,
		IsBlocked:               data.IsBlocked,
		IsFavorite:              data.IsFavorite,
		IsFriendRequestSend:     data.IsFriendRequestSend,
		IsFriendRequestRecieved: data.IsFriendRequestRecieved,
		FriendshipID:            data.FriendshipID,
		MutualConnectionsAmount: data.MutualConectionsAmount,
		ProfileCompletePercent:  int32(data.CompletePercent),
		CurrentTranslation:      data.CurrentTranslation,
		DateOfActivation:        timeToStringDayMonthAndYear(data.CreatedAt),
		Emails:                  make([]string, 0, len(data.Emails)),
		Phones:                  make([]string, 0, len(data.Phones)),
	}

	if data.MiddleName != nil {
		prof.Middlename = data.MiddleName.Middlename
	} else {
		prof.IsMiddlenameNull = true
	}

	if data.Nickname != nil {
		prof.Nickname = data.Nickname.Nickname
	} else {
		prof.IsNicknameNull = true
	}

	if data.Patronymic != nil {
		prof.Patronymic = data.Patronymic.Patronymic
	} else {
		prof.IsPatronymicNull = true
	}

	if nm := accountNativeNameToRPC(data.NativeName); nm != nil {
		prof.NativeName = &userRPC.NativeNameProfile{
			Language: nm.GetName(),
			Name:     nm.GetLanguageID(),
		}
	}

	if gender := genderToRPC(data.Gender); gender != nil {
		prof.Gender = gender.GetGender()
	}

	if data.Birthday != nil {
		prof.Birthday = timeToStringDayMonthAndYear(data.Birthday.Birthday)
	} else {
		prof.IsBirthdayNull = true
	}

	for _, e := range data.Emails {
		if e.Activated {
			if e.Primary {
				prof.Email = e.Email
			}

			prof.Emails = append(prof.Emails, e.Email)
		}
	}

	for _, e := range data.Phones {
		// TODO: uncomment when we will have SMS gateway
		// if e.Activated {
		if e.Primary {
			prof.PhoneNumber = e.CountryCode.Code + e.Number
		}

		prof.Phones = append(prof.Phones, e.CountryCode.Code+e.Number)
		// }
	}

	if len(data.AvailableTranslations) > 0 {
		prof.AvailableTranslations = data.AvailableTranslations
	}

	return prof
}

func profileExperienceToExperienceRPC(data *profile.Experience) *userRPC.Experience {
	if data == nil {
		return nil
	}

	exp := userRPC.Experience{
		ID:            data.GetID(),
		Position:      data.Position,
		StartDate:     timeToStringMonthAndYear(data.StartDate),
		CurrentlyWork: data.CurrentlyWork,
		Company:       data.Company, // remake
		Links:         make([]*userRPC.Link, 0, len(data.Links)),
		Files:         make([]*userRPC.File, 0, len(data.Files)),
		// IsCurrentlyWorkNull
	}

	// if data.CityID == nil {
	// 	exp.IsLocationNull = true
	// } else {
	// 	exp.CityID = *data.CityID
	// }

	exp.Location = locationToRPC(&data.Location)

	if data.Description == nil {
		exp.IsDescriptionNull = true
	} else {
		exp.Description = *data.Description
	}

	if data.FinishDate != nil {
		exp.FinishDate = timeToStringMonthAndYear(*(data.FinishDate))
	}

	for i := range data.Links {
		if data.Links[i] != nil { // TODO: findout why it why is saves nulls
			exp.Links = append(exp.Links, profileLinktoLinkRPC(data.Links[i]))
		}
	}

	for i := range data.Files {
		exp.Files = append(exp.Files, profileFiletoFileRPC(data.Files[i]))
	}

	return &exp
}

func profilePortfolioToPortfolioRPC(data *profile.Portfolio) *userRPC.Portfolio {
	if data == nil {
		return nil
	}

	port := userRPC.Portfolio{
		ID:         data.GetID(),
		Tittle:     data.Title,
		Files:      make([]*userRPC.File, 0, len(data.Files)),
		ViewsCount: data.ViewsCount,
		CreatedAt:  data.CreatedAt.String(),
		HasLiked:   data.HasLiked,
		LikesCount: data.LikesCount,
		// ContentType: profileContentTypeToUserRPCContentTypeEnum(data.ContentType),
	}

	if port.Description != "" {
		port.Description = data.Description
	}
	if port.Tittle != "" {
		port.Description = data.Description
	}

	for i := range data.Files {
		port.Files = append(port.Files, profileFiletoFileRPC(data.Files[i]))
	}

	return &port
}

func profileContentTypeToUserRPCContentTypeEnum(data profile.ContentType) userRPC.ContentTypeEnum {
	var content userRPC.ContentTypeEnum

	switch data {
	case profile.ContentTypePhoto:
		content = userRPC.ContentTypeEnum_Content_Type_Photo
	case profile.ContentTypeArticle:
		content = userRPC.ContentTypeEnum_Content_Type_Article
	case profile.ContentTypeVideo:
		content = userRPC.ContentTypeEnum_Content_Type_Video
	case profile.ContentTypeAudio:
		content = userRPC.ContentTypeEnum_Content_Type_Audio
	}

	return content
}

func userRPCContentTypeToProfileContentType(data userRPC.ContentTypeEnum) profile.ContentType {

	var content = profile.ContentTypePhoto

	switch data {
	case userRPC.ContentTypeEnum_Content_Type_Article:
		content = profile.ContentTypeArticle
	case userRPC.ContentTypeEnum_Content_Type_Video:
		content = profile.ContentTypeVideo
	case userRPC.ContentTypeEnum_Content_Type_Audio:
		content = profile.ContentTypeAudio
	}

	return content
}

func profileToolTechnologyToUserRPCToolTechnology(data *profile.ToolTechnology) *userRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	tool := userRPC.ToolTechnology{
		ID:             data.GetID(),
		ToolTechnology: data.ToolTechnology,
		Rank:           profileLevelToUserRPCLevel(data.Rank),
	}

	if tool.ToolTechnology != "" {
		tool.ToolTechnology = data.ToolTechnology
	}

	return &tool
}

func profileLevelToUserRPCLevel(data profile.Level) userRPC.ToolTechnology_Level {
	var tool userRPC.ToolTechnology_Level

	switch data {
	case profile.LevelBeginner:
		tool = userRPC.ToolTechnology_Level_Beginner
	case profile.LevelIntermediate:
		tool = userRPC.ToolTechnology_Level_Intermediate
	case profile.LevelAdvanced:
		tool = userRPC.ToolTechnology_Level_Advanced
	case profile.LevelMaster:
		tool = userRPC.ToolTechnology_Level_Master
	}

	return tool
}

func userRPCToolTechnologyToProfileToolTechnology(data userRPC.ToolTechnology_Level) profile.Level {
	var tool profile.Level

	switch data {
	case userRPC.ToolTechnology_Level_Beginner:
		tool = profile.LevelBeginner
	case userRPC.ToolTechnology_Level_Intermediate:
		tool = profile.LevelIntermediate
	case userRPC.ToolTechnology_Level_Advanced:
		tool = profile.LevelAdvanced
	case userRPC.ToolTechnology_Level_Master:
		tool = profile.LevelMaster
	}

	return tool
}

func profileLinktoLinkRPC(data *profile.Link) *userRPC.Link {
	if data == nil {
		return nil
	}

	return &userRPC.Link{
		ID:  data.GetID(),
		URL: data.URL,
	}
}

func profileFiletoFileRPC(data *profile.File) *userRPC.File {
	if data == nil {
		return nil
	}

	return &userRPC.File{
		ID:       data.GetID(),
		URL:      data.URL,
		MimeType: data.MimeType,
		Name:     data.Name,
	}
}

func userStatusToStatusRPC(stat status.UserStatus) userRPC.Status {
	var st userRPC.Status

	switch stat {
	case status.UserStatusActivated:
		st = userRPC.Status_ACTIVATED
	case status.UserStatusNotActivated:
		st = userRPC.Status_NOT_ACTIVATED
	}

	return st
}

func profileEducationToEducationRPC(data *profile.Education) *userRPC.Education {
	if data == nil {
		return nil
	}

	edu := userRPC.Education{
		ID:               data.GetID(),
		School:           data.School,
		FieldStudy:       data.FieldStudy,
		StartDate:        timeToStringMonthAndYear(data.StartDate),
		FinishDate:       timeToStringMonthAndYear(data.FinishDate),
		IsCurrentlyStudy: data.IsCurrentlyStudy,
		Links:            make([]*userRPC.Link, 0, len(data.Links)),
		Files:            make([]*userRPC.File, 0, len(data.Files)),
	}

	if data.Degree == nil {
		edu.IsDegreeNull = true
	} else {
		edu.Degree = *data.Degree
	}

	if data.Grade == nil {
		edu.IsGradeNull = true
	} else {
		edu.Grade = *data.Grade
	}

	if data.Description == nil {
		edu.IsDescriptionNull = true
	} else {
		edu.Description = *data.Description
	}

	for i := range data.Links {
		if data.Links[i] != nil { // TODO: findout why it why is saves nulls
			edu.Links = append(edu.Links, profileLinktoLinkRPC(data.Links[i]))
		}
	}

	for i := range data.Files {
		edu.Files = append(edu.Files, profileFiletoFileRPC(data.Files[i]))
	}

	edu.Location = locationToRPC(&data.Location)

	return &edu
}

func profileSkillToSkillRPC(data *profile.Skill) *userRPC.Skill {
	if data == nil {
		return nil
	}

	return &userRPC.Skill{
		ID:       data.GetID(),
		Position: data.Position,
		Skill:    data.Skill,
	}
}

func profileInterestToInterestRPC(data *profile.Interest) *userRPC.Interest {
	if data == nil {
		return nil
	}

	in := userRPC.Interest{
		ID:       data.GetID(),
		Interest: data.Interest,
	}

	if data.Image != nil {
		in.Image = data.Image.URL
	} else {
		in.IsImageNull = true
	}

	if data.Description != nil {
		in.Description = *data.Description
	} else {
		in.IsDescriptionNull = true
	}

	return &in
}

func profileAccomplishmentToAccomplishmentRPC(data *profile.Accomplishment) *userRPC.Accomplishment {
	if data == nil {
		return nil
	}

	acc := userRPC.Accomplishment{
		ID:   data.GetID(),
		Name: data.Name,
	}

	t := profileAccomplishmentTypeToAccomplishmentTypeRPC(data.Type)
	acc.Type = t

	if data.Description != nil {
		acc.Description = *data.Description
	} else {
		acc.IsDescriptionNull = true
	}

	if data.FinishDate != nil {
		acc.FinishDate = timeToStringMonthAndYear(*data.FinishDate)
	} else {
		acc.IsFinishDateNull = true
	}

	if data.IsExpire != nil {
		acc.IsExpire = *data.IsExpire
	} else {
		acc.IsIsExpireNull = true
	}

	if data.Issuer != nil {
		acc.Issuer = *data.Issuer
	} else {
		acc.IsIssuerNull = true
	}

	if data.LicenseNumber != nil {
		acc.LicenseNumber = *data.LicenseNumber
	} else {
		acc.IsLicenseNumberNull = true
	}

	if data.Score != nil {
		acc.Score = *data.Score
	} else {
		acc.IsScoreNull = true
	}

	if data.StartDate != nil {
		acc.StartDate = timeToStringMonthAndYear(*data.StartDate)
	} else {
		acc.IsStartDateNull = true
	}

	if data.URL != nil {
		acc.URL = *data.URL
	} else {
		acc.IsURLNull = true
	}

	for i := range data.Links {
		if data.Links[i] != nil { // TODO: findout why it why is saves nulls
			acc.Links = append(acc.Links, profileLinktoLinkRPC(data.Links[i]))
		}
	}

	for i := range data.Files {
		acc.Files = append(acc.Files, profileFiletoFileRPC(data.Files[i]))
	}

	return &acc
}

func profileAccomplishmentTypeToAccomplishmentTypeRPC(data profile.AccomplishmentType) userRPC.Accomplishment_AccomplishmentType {

	var accompl userRPC.Accomplishment_AccomplishmentType

	switch data {
	case profile.AccomplishmentTypeCertificate:
		accompl = userRPC.Accomplishment_Certificate
	case profile.AccomplishmentTypeTest:
		accompl = userRPC.Accomplishment_Test
	case profile.AccomplishmentTypeAward:
		accompl = userRPC.Accomplishment_Award
	case profile.AccomplishmentTypeLicense:
		accompl = userRPC.Accomplishment_License
	case profile.AccomplishmentTypeProject:
		accompl = userRPC.Accomplishment_Project
	case profile.AccomplishmentTypePublication:
		accompl = userRPC.Accomplishment_Publication
	}

	return accompl
}

func profileRecommendationToRecommendationRPC(data *profile.Recommendation) *userRPC.Recommendation {
	if data == nil {
		return nil
	}

	rec := userRPC.Recommendation{
		CreatedAt:     data.CreatedAt,
		ID:            data.ID,
		Text:          data.Text,
		Receiver:      profileToProfileRPC(data.Receiver),
		Recommendator: profileToProfileRPC(data.Recommendator),
		Title:         data.Title,
		Relation:      stringToRelationRPC(data.Relation),
	}
	if data.IsHidden != nil {
		rec.IsHidden = *data.IsHidden
	} else {
		rec.IsIsHiddenNull = true
	}

	return &rec
}

func profileRecommendationRequestToRecommendationRequestRPC(data *profile.RecommendationRequest) *userRPC.RecommendationRequest {
	if data == nil {
		return nil
	}

	r := userRPC.RecommendationRequest{
		ID:        data.ID,
		Text:      data.Text,
		Requestor: profileToProfileRPC(data.Requestor),
		Requested: profileToProfileRPC(data.Requested),
		CreatedAt: data.CreatedAt,
		Title:     data.Title,
		Relation:  stringToRelationRPC(data.Relation),
	}

	return &r
}
func stringToRelationRPC(data string) userRPC.RecommendationRelationEnum {

	switch data {
	case "EXPERIENCE":
		return userRPC.RecommendationRelationEnum_EXPERIENCE
	case "EDUCATION":
		return userRPC.RecommendationRelationEnum_EDUCATION
	case "ACCOMPLISHMENT":
		return userRPC.RecommendationRelationEnum_ACCOMPLISHMENT
	}

	return userRPC.RecommendationRelationEnum_NO_RELATION
}

func profileKnownLanguageToKnownLanguageRPC(data *profile.KnownLanguage) *userRPC.KnownLanguage {
	if data == nil {
		return nil
	}

	lang := userRPC.KnownLanguage{
		ID:       data.GetID(),
		Language: data.Language,
		Rank:     data.Rank,
	}

	return &lang
}

func invitaionInvitationToInvitaionRPC(data *invitation.Invitation) *userRPC.Invitation {
	if data == nil {
		return nil
	}

	return &userRPC.Invitation{
		Email: data.Email,
		Name:  data.Name,
	}
}

func accountPermissionTypeToPermissionTypeRPC(data *account.PermissionType) *userRPC.PermissionType {
	if data == nil {
		return nil
	}

	switch *data {
	case account.PermissionTypeMe:
		p := userRPC.PermissionType_ME
		return &p
	case account.PermissionTypeMembers:
		p := userRPC.PermissionType_MEMBERS
		return &p
	case account.PermissionTypeMyConnections:
		p := userRPC.PermissionType_MY_CONNECTIONS
		return &p
	}

	p := userRPC.PermissionType_NONE
	return &p
}

// ---------

func stringDateToTime(s string) time.Time {
	if date, err := time.Parse("2-1-2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func stringDayMonthAndYearToTime(s string) time.Time {
	if date, err := time.Parse("1-2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func timeToStringMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, _ := t.UTC().Date()
	return strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func timeToStringDayMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, d := t.UTC().Date()
	return strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func userRPCToolsToProfileToolTechnology(data *userRPC.ToolTechnology) *profile.ToolTechnology {
	if data == nil {
		return nil
	}

	tl := profile.ToolTechnology{
		ToolTechnology: data.GetToolTechnology(),
		Rank:           userRPCToolTechnologyToProfileToolTechnology(data.GetRank()),
	}

	_ = tl.SetID(data.GetID())

	return &tl
}
