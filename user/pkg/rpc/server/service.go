package serverRPC

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/status"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/invitation"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
	userReport "gitlab.lan/Rightnao-site/microservices/user/pkg/user_report"
)

// Service define functions inside Service
type Service interface {
	IsUsernameBusy(ctx context.Context, username string) (bool, error)
	CreateNewAccount(ctx context.Context, acc *account.Account, password string) (id, url, tmpToken string, err error)
	ActivateUser(ctx context.Context, code string, userID string) (res *account.LoginResponse, err error)

	Recover(ctx context.Context, email string, remindUsername, resetPasswrod bool) error
	GenerateRecoveryCode(ctx context.Context, login string /* methodOfRecovery */) error
	RecoverPassword(ctx context.Context, code, userID, password string) error

	IdentifyCountry(ctx context.Context) (string, error)
	CheckToken(ctx context.Context) (bool, error)
	Login(ctx context.Context, login, password, twoFACode string) (result account.LoginResponse, err error)
	SignOut(ctx context.Context) error
	SignOutSession(ctx context.Context, sessionID string) error
	SignOutFromAll(ctx context.Context) error

	GetAccount(ctx context.Context) (*account.Account, error)
	ChangeFirstName(ctx context.Context, firstname string) error
	ChangeLastName(ctx context.Context, lastname string) error
	ChangePatronymic(ctx context.Context, patronymic *string, permission *account.Permission) error
	ChangeNickname(ctx context.Context, nickname *string, permission *account.Permission) error
	ChangeMiddleName(ctx context.Context, middlename *string, permission *account.Permission) error
	ChangeNameOnNativeLanguage(ctx context.Context, name *string, lang *string, permission *account.Permission) error
	ChangeBirthday(ctx context.Context, birthday *time.Time, permission *account.Permission) error
	ChangeGender(ctx context.Context, gender *string, permission *account.Permission) error

	AddEmail(ctx context.Context, email string, permission *account.Permission) (id string, err error) // is it right?
	RemoveEmail(ctx context.Context, emailID string) error
	ChangeEmail(ctx context.Context, emailID string, permission *account.Permission, isPrimary bool) error

	AddPhone(ctx context.Context, countryCode *account.CountryCode, number string, permission *account.Permission) (string, error)
	RemovePhone(ctx context.Context, phoneID string) error
	ChangePhone(ctx context.Context, phoneID string, permission *account.Permission, isPrimary bool) error

	AddMyAddress(ctx context.Context, address *account.MyAddress) (string, error)
	RemoveMyAddress(ctx context.Context, addressID string) error
	ChangeMyAddress(ctx context.Context, address *account.MyAddress) error

	AddOtherAddress(ctx context.Context, address *account.OtherAddress) (string, error)
	RemoveOtherAddress(ctx context.Context, addressID string) error
	ChangeOtherAddress(ctx context.Context, address *account.OtherAddress) error

	ChangeUILanguage(ctx context.Context, lang string) error
	ChangePrivacy(ctx context.Context, priv *account.PrivacyItem, value *account.PermissionType) error
	ChangePassword(ctx context.Context, oldPass string, newPass string) error

	Init2FA(ctx context.Context) (qr string, url string, key string, err error)
	Enable2FA(ctx context.Context, code string) error
	Disable2FA(ctx context.Context, code string) error

	DeactivateAccount(ctx context.Context, password string) error

	// Profile
	GetProfile(ctx context.Context, url string, lang string) (*profile.Profile, error)
	GetProfileByID(ctx context.Context, userID string) (*profile.Profile, error)
	GetProfilesByID(ctx context.Context, ids []string, lang string) ([]*profile.Profile, error)
	GetMapProfilesByID(ctx context.Context, ids []string, lang string) (map[string]*profile.Profile, error)
	GetMyCompanies(ctx context.Context) (interface{}, error)

	GetUserPortfolioInfo(ctx context.Context, userID string) (*profile.PortfolioInfo, error)
	GetPortfolios(ctx context.Context, companyID string, userID string, first uint32, after, contentType string) (*profile.Portfolios, error)
	GetPortfolioByID(ctx context.Context, companyID string, userID string, portfolioID string) (*profile.Portfolio, error)

	AddPortfolio(ctx context.Context, port *profile.Portfolio) (id string, err error)
	AddSavedCountToPortfolio(ctx context.Context, ownerID string, porfolioID string) error
	AddViewCountToPortfolio(ctx context.Context, ownerID string, porfolioID string, companyID string) error

	GetPortfolioComments(ctx context.Context, porfolioID string, first uint32, after string) (*profile.GetPortfolioComments, error)
	AddCommentToPortfolio(ctx context.Context, comment *profile.PortfolioComment) (id string, err error)
	RemoveCommentInPortfolio(ctx context.Context, portfolioID, commentID, companyID string) error

	LikeUserPortfolio(ctx context.Context, ownerID string, porfolioID string, companyID string) error
	UnLikeUserPortfolio(ctx context.Context, ownerID string, porfolioID string, companyID string) error

	ChangeOrderFilesInPortfolio(ctx context.Context, portfolioID, fileID string, position uint32) error
	ChangePortfolio(ctx context.Context, port *profile.Portfolio) error
	RemovePortfolio(ctx context.Context, portID string) error
	AddLinksInPortfolio(ctx context.Context, expID string, links []*profile.Link) ([]string, error)
	ChangeLinkInPortfolio(ctx context.Context, eduID string, linkID, url string) error
	RemoveLinksInPortfolio(ctx context.Context, expID string, ids []string) error
	AddFileInPortfolio(ctx context.Context, userID string, portID string, file *profile.File) (string, error)
	RemoveFilesInPortfolio(ctx context.Context, expID string, ids []string) error

	GetToolsTechnologies(ctx context.Context, userID string, first uint32, after, lang string) ([]*profile.ToolTechnology, error)
	AddToolTechnology(ctx context.Context, exp []*profile.ToolTechnology) (ids []string, err error)
	ChangeToolTechnology(ctx context.Context, tool []*profile.ToolTechnology) error
	RemoveToolTechnology(ctx context.Context, toolID []string) error

	GetExperiences(ctx context.Context, userID string, first uint32, after string, lang string) ([]*profile.Experience, error)
	AddExperience(ctx context.Context, exp *profile.Experience) (id string, err error)
	ChangeExperience(ctx context.Context, exp *profile.Experience, changeIsCurrentlyWorking bool) error
	RemoveExperience(ctx context.Context, expID string) error
	AddLinksInExperience(ctx context.Context, expID string, links []*profile.Link) ([]string, error)
	AddFileInExperience(ctx context.Context, userID string, expID string, file *profile.File) (string, error)
	RemoveFilesInExperience(ctx context.Context, expID string, ids []string) error
	ChangeLinkInExperience(ctx context.Context, expID string, linkID, url string) error
	RemoveLinksInExperience(ctx context.Context, expID string, ids []string) error
	GetUploadedFilesInExperience(ctx context.Context) ([]*profile.File, error)

	GetEducations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Education, error)
	AddEducation(ctx context.Context, edu *profile.Education) (id string, err error)
	ChangeEducation(ctx context.Context, edu *profile.Education) error
	RemoveEducation(ctx context.Context, eduID string) error
	AddLinksInEducation(ctx context.Context, expID string, links []*profile.Link) ([]string, error)
	AddFileInEducation(ctx context.Context, userID string, expID string, file *profile.File) (string, error)
	RemoveFilesInEducation(ctx context.Context, eduID string, ids []string) error
	ChangeLinkInEducation(ctx context.Context, eduID string, linkID, url string) error
	RemoveLinksInEducation(ctx context.Context, eduID string, ids []string) error
	GetUploadedFilesInEducation(ctx context.Context) ([]*profile.File, error)

	GetSkills(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Skill, error)
	GetEndorsements(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Profile, error)
	AddSkills(ctx context.Context, skills []*profile.Skill) ([]string, error)
	ChangeOrderOfSkill(ctx context.Context, skillID string, position uint32) error
	RemoveSkills(ctx context.Context, ids []string) error
	VerifySkill(ctx context.Context, targetID string, skillID string) error
	UnverifySkill(ctx context.Context, targetID string, skillID string) error

	GetInterests(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Interest, error)
	AddInterest(ctx context.Context, interest *profile.Interest) (id string, err error)
	ChangeInterest(ctx context.Context, interest *profile.Interest) error
	RemoveInterest(ctx context.Context, interestID string) error
	ChangeImageInterest(ctx context.Context, targetID string, interestID string, image *profile.File) (id string, err error)
	RemoveImageInInterest(ctx context.Context, interestID string) error
	GetUnuploadImageInInterest(ctx context.Context) (*profile.File, error)
	GetOriginImageInInterest(ctx context.Context, interestID string) (string, error)
	ChangeOriginImageInInterest(ctx context.Context, userID string, interestID string, url string) error

	GetAccomplishments(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Accomplishment, error)
	AddAccomplishment(ctx context.Context, accoplishment *profile.Accomplishment) (id string, err error)
	ChangeAccomplishment(ctx context.Context, accoplishment *profile.Accomplishment) error
	RemoveAccomplishment(ctx context.Context, accoplishmentID string) error
	AddFileInAccomplishment(ctx context.Context, userID string, expID string, file *profile.File) (string, error)
	RemoveFilesInAccomplishment(ctx context.Context, expID string, ids []string) error
	AddLinksInAccomplishment(ctx context.Context, expID string, links []*profile.Link) ([]string, error)
	RemoveLinksInAccomplishment(ctx context.Context, eduID string, ids []string) error
	GetUploadedFilesInAccomplishment(ctx context.Context) ([]*profile.File, error)

	GetReceivedRecommendations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Recommendation, error)
	GetGivenRecommendations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Recommendation, error)
	GetReceivedRecommendationRequests(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.RecommendationRequest, error)
	GetRequestedRecommendationRequests(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.RecommendationRequest, error)
	GetHiddenRecommendations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Recommendation, error)

	GetKnownLanguages(ctx context.Context, targetID string, first uint32, after string) ([]*profile.KnownLanguage, error)
	AddKnownLanguage(ctx context.Context, lang *profile.KnownLanguage) (string, error)
	ChangeKnownLanguage(ctx context.Context, lang *profile.KnownLanguage) error
	RemoveKnownLanguage(ctx context.Context, langID string) error

	ChangeHeadline(ctx context.Context, headline string) error
	ChangeStory(ctx context.Context, story string) error

	GetOriginAvatar(ctx context.Context) (string, error)
	ChangeOriginAvatar(ctx context.Context, userID string, url string) error
	ChangeAvatar(ctx context.Context, userID string, url string) error
	RemoveAvatar(ctx context.Context) error

	CheckPassword(ctx context.Context, password string) error

	SaveUserProfileTranslation(ctx context.Context, lang string, tr *profile.Translation) error
	SaveUserExperienceTranslation(ctx context.Context, expID string, lang string, tr *profile.ExperienceTranslation) error
	SaveUserEducationTranslation(ctx context.Context, educationID string, lang string, tr *profile.EducationTranslation) error
	SaveUserInterestTranslation(ctx context.Context, interestID string, lang string, tr *profile.InterestTranslation) error
	SaveUserSkillTranslation(ctx context.Context, skillID string, lang string, tr *profile.SkillTranslation) error
	SaveUserAccomplishmentTranslation(ctx context.Context, accomplishmentID string, lang string, tr *profile.AccomplishmentTranslation) error
	SaveUserPortfolioTranslation(ctx context.Context, portfolioID string, lang string, tr *profile.PortfolioTranslation) error
	SaveUserToolTechnologyTranslation(ctx context.Context, toolTechID string, lang string, tr *profile.ToolTechnologyTranslation) error
	RemoveTransaltion(ctx context.Context, lang string) error

	ReportUser(ctx context.Context, report *userReport.Report) error

	SentEmailInvitation(ctx context.Context, email string, name string, companyID string) error
	GetInvitation(ctx context.Context) ([]invitation.Invitation, int32, error)
	GetInvitationForCompany(ctx context.Context, companyID string) ([]invitation.Invitation, int32, error)

	GetConectionsPrivacy(ctx context.Context, userID string) (account.PermissionType, error)

	ContactInvitationForWallet(ctx context.Context, name string, email string, message string, coins int32) error
	GetUserByInvitedID(ctx context.Context, userID string) (int32, error)
	AddGoldCoinsToWallet(ctx context.Context, userID string, coins int32) error
	CreateWalletAccount(ctx context.Context, userID string) error

	GetAllUsersForAdmin(ctx context.Context, first uint32, after string) (*profile.Users, error)
	ChangeUserStatus(ctx context.Context, userID string, status status.UserStatus) error

	GetUsersForAdvert(ctx context.Context, data account.UserForAdvert) ([]string, error)
}
