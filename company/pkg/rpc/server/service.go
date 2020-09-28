package serverRPC

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	careercenter "gitlab.lan/Rightnao-site/microservices/company/pkg/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
)

// Service define functions inside Service
type Service interface {
	CheckIfURLForCompanyIsTaken(ctx context.Context, url string) (bool, error)
	CreateNewAccount(ctx context.Context, acc *account.Account) (string, string, error)
	DeactivateCompany(ctx context.Context, companyID string, password string) error

	GetCompanyAccount(ctx context.Context, companyID string) (*account.Account, error)

	ChangeCompanyName(ctx context.Context, companyID string, name string) error
	ChangeCompanyURL(ctx context.Context, companyID string, url string) error
	ChangeCompanyFoundationDate(ctx context.Context, companyID string, foundationDate time.Time) error
	ChangeCompanyIndustry(ctx context.Context, companyID string, industry *account.Industry) error
	ChangeCompanyType(ctx context.Context, companyID string, companyType account.Type) error
	ChangeCompanySize(ctx context.Context, companyID string, size account.Size) error

	AddCompanyEmail(ctx context.Context, companyID string, email *account.Email) (string, error)
	DeleteCompanyEmail(ctx context.Context, companyID string, emailID string) error
	ChangeCompanyEmail(ctx context.Context, companyID string, emailID string) error

	AddCompanyPhone(ctx context.Context, companyID string, phone *account.Phone) (string, error)
	DeleteCompanyPhone(ctx context.Context, companyID string, phoneID string) error
	ChangeCompanyPhone(ctx context.Context, companyID string, phoneID string) error

	AddCompanyAddress(ctx context.Context, companyID string, address *account.Address) (string, error)
	DeleteCompanyAddress(ctx context.Context, companyID string, addressID string) error
	ChangeCompanyAddress(ctx context.Context, companyID string, address *account.Address) error

	AddCompanyWebsite(ctx context.Context, companyID string, website string) (string, error)
	DeleteCompanyWebsite(ctx context.Context, companyID string, websiteID string) error
	ChangeCompanyWebsite(ctx context.Context, companyID string, websiteID string, website string) error

	ChangeCompanyParking(ctx context.Context, companyID string, parking account.Parking) error
	ChangeCompanyBenefits(ctx context.Context, companyID string, benefits []profile.Benefit) error

	AddCompanyAdmin(ctx context.Context, companyID string, userID string, level account.AdminLevel, password string) error
	DeleteCompanyAdmin(ctx context.Context, companyID string, userID string, password string) error

	// profile

	GetCompanyProfile(ctx context.Context, url string, lang string) (*profile.Profile, account.AdminLevel, error)
	GetCompanyProfileByID(ctx context.Context, url string, lang string) (*profile.Profile, account.AdminLevel, error)
	GetCompanyProfiles(ctx context.Context, ids []string) ([]*profile.Profile, error)

	ChangeCompanyAboutUs(ctx context.Context, companyID string, aboutUs *profile.AboutUs) error

	GetFounders(ctx context.Context, companyID string, first int32, after string) ([]*profile.Founder, error)
	AddCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) (string, error)
	DeleteCompanyFounder(ctx context.Context, companyID string, founderID string) error
	ChangeCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) error
	ChangeCompanyFounderAvatar(ctx context.Context, companyID string, founderID string, file *profile.File) (string, error)
	RemoveCompanyFounderAvatar(ctx context.Context, companyID string, founderID string) error
	ApproveFounderRequest(ctx context.Context, companyID, requestID string) error
	RemoveFounderRequest(ctx context.Context, companyID, requestID string) error

	AddCompanyAward(ctx context.Context, companyID string, award *profile.Award) (string, error)
	DeleteCompanyAward(ctx context.Context, companyID string, awardID string) error
	ChangeCompanyAward(ctx context.Context, companyID string, award *profile.Award) error
	AddLinksInCompanyAward(ctx context.Context, companyID, awardID string, links []*profile.Link) ([]string, error)
	AddFileInCompanyAward(ctx context.Context, companyID, awardID string, file *profile.File) (string, error)
	RemoveFilesInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error
	ChangeLinkInCompanyAward(ctx context.Context, companyID, awardID, linkID, url string) error
	RemoveLinksInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error
	GetUploadedFilesInCompanyAward(ctx context.Context, companyID string) ([]*profile.File, error)

	AddCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) (string, error)
	DeleteCompanyMilestone(ctx context.Context, companyID string, milestoneID string) error
	ChangeCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) error
	ChangeImageMilestone(ctx context.Context, companyID string, milestoneID string, file *profile.File) (string, error)
	RemoveImageInMilestone(ctx context.Context, companyID string, milestoneID string) error

	AddCompanyProduct(ctx context.Context, companyID string, product *profile.Product) (string, error)
	DeleteCompanyProduct(ctx context.Context, companyID string, productID string) error
	ChangeCompanyProduct(ctx context.Context, companyID string, product *profile.Product) error
	ChangeImageProduct(ctx context.Context, companyID string, productID string, file *profile.File) (string, error)
	RemoveImageInProduct(ctx context.Context, companyID string, productID string) error

	AddCompanyService(ctx context.Context, companyID string, service *profile.Service) (string, error)
	DeleteCompanyService(ctx context.Context, companyID string, serviceID string) error
	ChangeCompanyService(ctx context.Context, companyID string, service *profile.Service) error
	ChangeImageService(ctx context.Context, companyID string, serviceID string, file *profile.File) (string, error)
	RemoveImageInService(ctx context.Context, companyID string, serviceID string) error

	AddCompanyReport(ctx context.Context, companyID string, report *profile.Report) error

	AddCompanyReview(ctx context.Context, companyID string, review *profile.Review) (string, error)
	DeleteCompanyReview(ctx context.Context, companyID string, reviewID string) error
	GetCompanyReviews(ctx context.Context, companyID string, first uint32, after string) ([]*profile.Review, error)
	GetUsersRevies(ctx context.Context, userID string, first uint32, after string) ([]*profile.Review, error)

	GetAvarageRateOfCompany(ctx context.Context, companyID string) (float32, uint32, error)
	GetAmountOfEachRate(ctx context.Context, companyID string) (map[uint32]uint32, error)

	AddCompanyReviewReport(ctx context.Context, reviewReport *profile.ReviewReport) error

	ChangeAvatar(ctx context.Context, companyID string, file *profile.File) error
	GetOriginAvatar(ctx context.Context, companyID string) (string, error)
	ChangeOriginAvatar(ctx context.Context, companyID string, file *profile.File) error
	RemoveAvatar(ctx context.Context, companyID string) error

	ChangeCover(ctx context.Context, companyID string, file *profile.File) error
	GetOriginCover(ctx context.Context, companyID string) (string, error)
	ChangeOriginCover(ctx context.Context, companyID string, file *profile.File) error
	RemoveCover(ctx context.Context, companyID string) error

	SaveCompanyProfileTranslation(ctx context.Context, companyID string, lang string, tr *profile.Translation) error
	SaveCompanyMilestoneTranslation(ctx context.Context, companyID string, milestoneID string, language string, translation *profile.MilestoneTranslation) error
	SaveCompanyAwardTranslation(ctx context.Context, companyID string, awardID string, language string, translation *profile.AwardTranslation) error

	GetAmountOfReviewsOfUser(ctx context.Context, userID string) (int32, error)

	GetCompanyGallery(ctx context.Context, companyID string, first, after uint32) ([]*profile.File, error)
	AddFileInCompanyGallery(ctx context.Context, companyID string, file *profile.File) (string, error)
	RemoveFilesInCompanyGallery(ctx context.Context, companyID string, ids []string) error
	GetUploadedFilesInCompanyGallery(ctx context.Context, companyID string) ([]*profile.File, error)

	AddGoldCoinsToWallet(ctx context.Context, userID string, coins int32) error

	OpenCareerCenter(ctx context.Context, companyID string, cc *careercenter.CareerCenter) error
}
