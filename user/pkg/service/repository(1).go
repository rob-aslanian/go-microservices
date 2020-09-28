package service

import (
	"context"
	"net"
	"time"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/invitation"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
	arangorepo "gitlab.lan/Rightnao-site/microservices/user/pkg/repository/arango"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/status"
	userReport "gitlab.lan/Rightnao-site/microservices/user/pkg/user_report"
)

// UsersRepository contains functions which have to be in users repository
type UsersRepository interface {
	SaveNewAccount(ctx context.Context, acc *account.Account, password string) error
	IsUsernameBusy(ctx context.Context, username string) (bool, error)
	IsEmailAlreadyInUse(ctx context.Context, email string) (bool, error)
	GetUserIDAndUsernameAndPrimaryEmailByLogin(ctx context.Context, login string) (id string, username, email string, err error)
	GetUserIDAndPrimaryEmailByLogin(ctx context.Context, login string) (id string, email string, err error)
	// GetPrimaryEmailByUserID(ctx context.Context, userID string) (email string, err error)
	ChangeStatusOfUser(ctx context.Context, userID string, st status.UserStatus) error
	ActivateEmail(ctx context.Context, userID string, email string) error
	ChangePassword(ctx context.Context, userID string, password string) error
	GetCredentialsAndStatus(ctx context.Context, login string) (res account.LoginResponse, err error)
	GetCredentialsByUserID(ctx context.Context, userID string) (res account.LoginResponse, err error)

	GetAccount(ctx context.Context, userID string) (*account.Account, error)
	ChangeFirstName(ctx context.Context, userID string, firstname string) error
	ChangeLastName(ctx context.Context, userID string, lastname string) error
	GetDateOfRegistration(ctx context.Context, userID string) (time.Time, error)
	GetDateOfActivation(ctx context.Context, userID string) (time.Time, error)
	SetDateOfActivation(ctx context.Context, userID string, date time.Time) error
	ChangePatronymic(ctx context.Context, userID string, patronymic *string, permission *account.Permission) error
	ChangeNickname(ctx context.Context, userID string, nickname *string, permission *account.Permission) error
	ChangeMiddleName(ctx context.Context, userID string, middlename *string, permission *account.Permission) error
	ChangeNameOnNativeLanguage(ctx context.Context, userID string, name *string, lang *string, permission *account.Permission) error
	ChangeBirthday(ctx context.Context, userID string, birthday *time.Time, permission *account.Permission) error
	ChangeGender(ctx context.Context, userID string, gender *string, permission *account.Permission) error

	IsEmailAdded(ctx context.Context, userID string, email string) (bool, error)
	AddEmail(ctx context.Context, userID string, email *account.Email) error
	IsPrimaryEmail(ctx context.Context, userID string, emailID string) (bool, error)
	RemoveEmail(ctx context.Context, userID string, emailID string) error
	ChangeEmailPermission(ctx context.Context, userID string, emailID string, permission *account.Permission) error
	IsEmailActivated(ctx context.Context, userID string, emailID string) (bool, error)
	MakeEmailPrimary(ctx context.Context, userID string, emailID string) error

	AddPhone(ctx context.Context, userID string, phone *account.Phone) error
	IsPhoneAdded(ctx context.Context, userID string, phone *account.Phone) (bool, error)
	IsPhoneAlreadyInUse(ctx context.Context, phone *account.Phone) (bool, error)
	RemovePhone(ctx context.Context, userID string, phoneID string) error
	IsPrimaryPhone(ctx context.Context, userID string, phoneID string) (bool, error)
	ChangePhonePermission(ctx context.Context, userID string, phoneID string, permission *account.Permission) error
	IsPhoneActivated(ctx context.Context, userID string, phoneID string) (bool, error)
	MakePhonePrimary(ctx context.Context, userID string, phoneID string) error

	AddMyAddress(ctx context.Context, userID string, address *account.MyAddress) error
	RemoveMyAddress(ctx context.Context, userID string, addressID string) error
	ChangeMyAddress(ctx context.Context, userID string, address *account.MyAddress) error

	AddOtherAddress(ctx context.Context, userID string, address *account.OtherAddress) error
	RemoveOtherAddress(ctx context.Context, userID string, addressID string) error
	ChangeOtherAddress(ctx context.Context, userID string, address *account.OtherAddress) error

	ChangeUILanguage(ctx context.Context, userID string, languageID string) error
	ChangePrivacy(ctx context.Context, userID string, priv account.PrivacyItem, value account.PermissionType) error

	Save2FASecret(ctx context.Context, userID string, secret []byte) error
	Get2FAInfo(ctx context.Context, userID string) (is2FAEnabled bool, secret []byte, err error)
	Enable2FA(ctx context.Context, userID string) error
	Disable2FA(ctx context.Context, userID string) error

	GetProfileByURL(ctx context.Context, url string) (*profile.Profile, error)
	GetProfileByID(ctx context.Context, userID string) (*profile.Profile, error)
	GetProfilesByID(ctx context.Context, ids []string) ([]*profile.Profile, error)

	GetPortfolios(ctx context.Context, userID string, first uint32, after uint32, contentType string) (*profile.Portfolios, error)
	GetPortfolioByID(ctx context.Context, userID string, portfolioID string) (*profile.Portfolio, error)
	GetUserPortfolioInfo(ctx context.Context, userID string) (*profile.PortfolioInfo, error)
	GetUserPortfolioCommentAmount(ctx context.Context, userID string) (int32, error)
	GetUserPortfolioLikeAmount(ctx context.Context, userID string) (int32, error)
	GetUserPortfolioViewAmount(ctx context.Context, userID string) (int32, error)

	AddPortfolio(ctx context.Context, userID string, port *profile.Portfolio) error
	AddSavedCountToPortfolio(ctx context.Context, ownerID string, portfolioID string) error
	AddViewCountToPortfolio(ctx context.Context, port *profile.PortfolioAction) error

	GetPortfolioComments(ctx context.Context, porfolioID string, first uint32, after uint32) (*profile.GetPortfolioComments, error)
	AddCommentToPortfolio(ctx context.Context, comment *profile.PortfolioComment) error
	RemoveCommentInPortfolio(ctx context.Context, userID, portfolioID, commentID string) error

	LikeUserPortfolio(ctx context.Context, port *profile.PortfolioAction) error
	UnLikeUserPortfolio(ctx context.Context, port *profile.PortfolioAction) error

	GetPortfolioViewCount(ctx context.Context, portfolioID string) (int32, error)
	GetPortfolioLikes(ctx context.Context, profileID, portfolioID string) (*profile.PortfolioLikes, error)

	GetPositionOfLastFile(ctx context.Context, userID, portfolioID string) (uint32, error)
	ChangeOrderFilesInPortfolio(ctx context.Context, userID, portfolioID, fileID string, position uint32) error
	ChangePortfolio(ctx context.Context, userID string, port *profile.Portfolio) error
	RemovePortfolio(ctx context.Context, userID string, portID string) error
	AddLinksInPortfolio(ctx context.Context, userID string, expID string, links []*profile.Link) error
	ChangeLinkInPortfolio(ctx context.Context, userID string, eduID string, linkID string, url string) error
	RemoveLinksInPortfolio(ctx context.Context, userID string, expID string, ids []string) error
	AddFileInPortfolio(ctx context.Context, userID string, portID string, file *profile.File) error
	RemoveFilesInPortfolio(ctx context.Context, userID string, portID string, ids []string) error

	GetToolTechnology(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.ToolTechnology, error)
	AddToolTechnology(ctx context.Context, userID string, tools []*profile.ToolTechnology) error
	ChangeToolTechnology(ctx context.Context, userID string, tools []*profile.ToolTechnology) error
	RemoveToolTechnology(ctx context.Context, userID string, toolsID []string) error

	GetExperiences(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Experience, error)
	AddExperience(ctx context.Context, userID string, exp *profile.Experience) error
	ChangeExperience(ctx context.Context, userID string, exp *profile.Experience, changeIsCurrentlyWorking bool) error
	AddLinksInExperience(ctx context.Context, userID string, expID string, links []*profile.Link) error
	RemoveLinksInExperience(ctx context.Context, userID string, expID string, ids []string) error
	RemoveExperience(ctx context.Context, userID string, expID string) error
	AddFileInExperience(ctx context.Context, userID string, expID string, file *profile.File) error
	RemoveFilesInExperience(ctx context.Context, userID string, expID string, ids []string) error
	ChangeLinkInExperience(ctx context.Context, userID string, eduID string, linkID string, url string) error
	GetUploadedFilesInExperience(ctx context.Context, userID string) ([]*profile.File, error)

	GetEducations(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Education, error)
	AddEducation(ctx context.Context, userID string, edu *profile.Education) error
	ChangeEducation(ctx context.Context, userID string, edu *profile.Education) error
	AddLinksInEducation(ctx context.Context, userID string, eduID string, links []*profile.Link) error
	RemoveEducation(ctx context.Context, userID string, eduID string) error
	AddFileInEducation(ctx context.Context, userID string, eduID string, file *profile.File) error
	RemoveFilesInEducation(ctx context.Context, userID string, eduID string, ids []string) error
	ChangeLinkInEducation(ctx context.Context, userID string, eduID string, linkID string, url string) error
	RemoveLinksInEducation(ctx context.Context, userID string, eduID string, ids []string) error
	GetUploadedFilesInEducation(ctx context.Context, userID string) ([]*profile.File, error)

	GetSkills(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Skill, error)
	GetEndorsements(ctx context.Context, skillID string, first uint32, after uint32) ([]string, error)
	AddSkills(ctx context.Context, userID string, skills []*profile.Skill) error
	GetPositionOfLastSkill(ctx context.Context, userID string) (uint32, error)
	ChangeOrderOfSkill(ctx context.Context, userID string, skillID string, position uint32) error
	RemoveSkills(ctx context.Context, userID string, ids []string) error
	VerifySkill(ctx context.Context, userID string, targetID string, skillID string) error
	IsSkillVerified(ctx context.Context, userID string, targetID string, skillID string) (bool, error)
	UnverifySkill(ctx context.Context, userID string, targetID string, skillID string) error

	GetInterests(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Interest, error)
	AddInterest(ctx context.Context, userID string, interest *profile.Interest) error
	ChangeInterest(ctx context.Context, userID string, interest *profile.Interest) error
	RemoveInterest(ctx context.Context, userID string, interestID string) error
	ChangeImageInterest(ctx context.Context, userID string, interestID string, image *profile.File) error
	RemoveImageInInterest(ctx context.Context, userID string, interestID string) error
	GetUnuploadImageInInterest(ctx context.Context, userID string) (*profile.File, error)
	GetOriginImageInInterest(ctx context.Context, userID, interestID string) (string, error)
	ChangeOriginImageInInterest(ctx context.Context, userID, interestID, url string) error

	GetAccomplishments(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Accomplishment, error)
	AddAccomplishment(ctx context.Context, userID string, accoplishment *profile.Accomplishment) error
	ChangeAccomplishment(ctx context.Context, userID string, accoplishment *profile.Accomplishment) error
	RemoveAccomplishment(ctx context.Context, userID string, accoplishmentID string) error

	GetKnownLanguages(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.KnownLanguage, error)
	AddKnownLanguage(ctx context.Context, userID string, lang *profile.KnownLanguage) error
	ChangeKnownLanguage(ctx context.Context, userID string, lang *profile.KnownLanguage) error
	RemoveKnownLanguage(ctx context.Context, userID string, langID string) error

	ChangeHeadline(ctx context.Context, userID string, headline string) error
	ChangeStory(ctx context.Context, userID string, headline string) error

	GetOriginAvatar(ctx context.Context, userID string) (avatarPath string, err error)
	ChangeOriginAvatar(ctx context.Context, userID string, url string) error
	ChangeAvatar(ctx context.Context, userID string, url string) error
	RemoveAvatar(ctx context.Context, userID string) error

	GetInfoAboutCompletionProfile(ctx context.Context, userID string) (exp, edu, skills, langs, interests, tools bool, err error)

	AddFileInAccomplishment(ctx context.Context, userID string, expID string, file *profile.File) error
	RemoveFilesInAccomplishment(ctx context.Context, userID string, expID string, ids []string) error
	AddLinksInAccomplishment(ctx context.Context, userID string, expID string, links []*profile.Link) error
	RemoveLinksInAccomplishment(ctx context.Context, userID string, expID string, ids []string) error
	GetUploadedFilesInAccomplishment(ctx context.Context, userID string) ([]*profile.File, error)

	SaveUserProfileTranslation(ctx context.Context, userID string, lang string, tr *profile.Translation) error
	SaveUserExperienceTranslation(ctx context.Context, userID string, expID string, lang string, tr *profile.ExperienceTranslation) error
	SaveUserEducationTranslation(ctx context.Context, userID string, educationID string, lang string, tr *profile.EducationTranslation) error
	SaveUserInterestTranslation(ctx context.Context, userID string, interestID string, lang string, tr *profile.InterestTranslation) error
	SaveUserSkillTranslation(ctx context.Context, userID string, skillID string, lang string, tr *profile.SkillTranslation) error
	SaveUserAccomplishmentTranslation(ctx context.Context, userID string, accomplishmentID string, lang string, tr *profile.AccomplishmentTranslation) error
	SaveUserPortfolioTranslation(ctx context.Context, userID string, portfolioID string, lang string, tr *profile.PortfolioTranslation) error
	SaveUserToolTechnologyTranslation(ctx context.Context, userID string, toolTechID string, lang string, tr *profile.ToolTechnologyTranslation) error
	RemoveTransaltion(ctx context.Context, userID string, lang string) error

	SaveReport(ctx context.Context, userID string, report *userReport.Report) error

	IsInvitationSend(ctx context.Context, email string) (bool, error)
	SaveInvitation(ctx context.Context, inv invitation.Invitation) error
	GetInvitation(ctx context.Context, userID string) ([]invitation.Invitation, int32, error)
	GetInvitationForCompany(ctx context.Context, companyID string) ([]invitation.Invitation, int32, error)

	GetPrivacyMyConnections(ctx context.Context, userID string) (account.PermissionType, error)
	GetAllUsersForAdmin(ctx context.Context, first uint32, after uint32) (*profile.Users, error)
	ChangeUserStatus(ctx context.Context, userID string, status status.UserStatus) error
	GetUserByInvitedID(ctx context.Context, userID string) (int32, error)

	GetUsersForAdvert(ctx context.Context, data account.UserForAdvert) ([]string, error)
}

// CacheRepository contains functions which have to be in cache repository
type CacheRepository interface {
	CreateTemporaryCodeForEmailActivation(ctx context.Context, userID string, email string) (string, error)
	CreateTemporaryCodeForNotActivatedUser(id string) (string, error)
	CreateTemporaryCodeForRecoveryByEmail(id string) (string, error)
	CheckTemporaryCodeForEmailActivation(ctx context.Context, userID string, code string) (bool, string, error)
	CheckTemporaryCodeForNotActivatedUser(userID string, code string) (bool, error)
	CheckTemporaryCodeForRecoveryByEmail(userID string, code string) (bool, error)
	Remove(key string) error
}

// GeoipRepository contains functions which have to be in geoip repository
type GeoipRepository interface {
	GetCountryISOCode(ipAddress net.IP) (string, error)
}

type ArangoRepo interface {
	SaveUser(ctx context.Context, u *arangorepo.User) error
}
