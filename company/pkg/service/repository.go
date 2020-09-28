package service

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	careercenter "gitlab.lan/Rightnao-site/microservices/company/pkg/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/status"
	arangorepo "gitlab.lan/Rightnao-site/microservices/company/pkg/repository/arango"
)

// CompanyRepository contains functions which have to be in company repository
type CompanyRepository interface {
	SaveNewCompanyAccount(ctx context.Context, data *account.Account) error
	ChangeStatusOfCompany(ctx context.Context, companyID string, stat status.CompanyStatus) error
	GetCompanyAccount(ctx context.Context, companyID string) (*account.Account, error)
	ChangeCompanyName(ctx context.Context, companyID string, name string) error
	IsURLBusy(ctx context.Context, url string) (bool, error)
	ChangeCompanyURL(ctx context.Context, companyID string, url string) error
	ChangeCompanyFoundationDate(ctx context.Context, companyID string, foundationDate time.Time) error
	ChangeCompanyIndustry(ctx context.Context, companyID string, industry *account.Industry) error
	ChangeCompanyType(ctx context.Context, companyID string, companyType account.Type) error
	ChangeCompanySize(ctx context.Context, companyID string, size account.Size) error

	AddCompanyEmail(ctx context.Context, companyID string, email *account.Email) error
	DeleteCompanyEmail(ctx context.Context, companyID string, emailID string) error
	IsEmailExists(ctx context.Context, emailID string) (bool, error)
	IsEmailActivated(ctx context.Context, companyID string, emailID string) (bool, error)
	IsEmailPrimary(ctx context.Context, companyID string, emailID string) (bool, error)
	IsEmailAdded(ctx context.Context, companyID string, email string) (bool, error)
	MakeEmailPrimary(ctx context.Context, companyID string, emailID string) error
	ActivateEmail(ctx context.Context, companyID string, email string) error

	AddCompanyPhone(ctx context.Context, companyID string, phone *account.Phone) error
	DeleteCompanyPhone(ctx context.Context, companyID string, phoneID string) error
	IsPhoneExists(ctx context.Context, companyID string, phone *account.Phone) (bool, error)
	IsPhoneActivated(ctx context.Context, companyID string, phoneID string) (bool, error)
	IsPhonePrimary(ctx context.Context, companyID string, phoneID string) (bool, error)
	IsPhoneAdded(ctx context.Context, companyID string, phone *account.Phone) (bool, error)
	MakePhonePrimary(ctx context.Context, companyID string, phoneID string) error

	AddCompanyAddress(ctx context.Context, companyID string, address *account.Address) error
	DeleteCompanyAddress(ctx context.Context, companyID string, addressID string) error
	ChangeCompanyAddress(ctx context.Context, companyID string, address *account.Address) error

	AddCompanyWebsite(ctx context.Context, companyID string, website *account.Website) error
	DeleteCompanyWebsite(ctx context.Context, companyID string, websiteID string) error
	ChangeCompanyWebsite(ctx context.Context, companyID string, websiteID string, website string) error

	ChangeCompanyParking(ctx context.Context, companyID string, parking account.Parking) error
	ChangeCompanyBenefits(ctx context.Context, companyID string, benefits []profile.Benefit) error

	GetCompanyProfile(ctx context.Context, url string) (*profile.Profile, error)
	GetCompanyProfileByID(ctx context.Context, companyID string) (*profile.Profile, error)
	GetCompanyProfiles(ctx context.Context, ids []string) ([]*profile.Profile, error)

	ChangeCompanyAboutUs(ctx context.Context, companyID string, aboutUs *profile.AboutUs) error

	GetFounders(ctx context.Context, companyID string, first int32, afterNumber int) ([]*profile.Founder, error)
	AddCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) error
	DeleteCompanyFounder(ctx context.Context, companyID string, founderID string) error
	ChangeCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) error
	ChangeCompanyFounderAvatar(ctx context.Context, companyID string, milestoneID string, image *profile.File) error
	RemoveCompanyFounderAvatar(ctx context.Context, companyID string, milestoneID string) error
	ApproveFounderRequest(ctx context.Context, companyID, requestID, userID string, byCompanyAdmin bool) error
	RemoveFounderRequest(ctx context.Context, companyID, requestID, userID string, byCompanyAdmin bool) error

	AddCompanyAward(ctx context.Context, companyID string, award *profile.Award) error
	DeleteCompanyAward(ctx context.Context, companyID string, awardID string) error
	ChangeCompanyAward(ctx context.Context, companyID string, award *profile.Award) error
	AddLinksInCompanyAward(ctx context.Context, companyID, expID string, links []*profile.Link) error
	RemoveLinksInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error
	RemoveCompanyAward(ctx context.Context, companyID, awardID string) error
	AddFileInCompanyAward(ctx context.Context, companyID, awardID string, file *profile.File) error
	RemoveFilesInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error
	ChangeLinkInCompanyAward(ctx context.Context, companyID, awardID string, linkID string, url string) error
	GetUploadedFilesInCompanyAward(ctx context.Context, companyID string) ([]*profile.File, error)

	AddFileInCompanyGallery(ctx context.Context, companyID string, file *profile.File) error
	RemoveFilesInCompanyGallery(ctx context.Context, companyID string, ids []string) error
	GetUploadedFilesInCompanyGallery(ctx context.Context, companyID string) ([]*profile.File, error)
	GetCompanyGallery(ctx context.Context, companyID string, first, afterNumber uint32) ([]*profile.File, error)

	AddCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) error
	DeleteCompanyMilestone(ctx context.Context, companyID string, milestoneID string) error
	ChangeCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) error
	ChangeImageMilestone(ctx context.Context, companyID string, milestoneID string, image *profile.File) error
	RemoveImageInMilestone(ctx context.Context, companyID string, milestoneID string) error

	AddCompanyProduct(ctx context.Context, companyID string, product *profile.Product) error
	DeleteCompanyProduct(ctx context.Context, companyID string, productID string) error
	ChangeCompanyProduct(ctx context.Context, companyID string, product *profile.Product) error
	ChangeImageProduct(ctx context.Context, companyID string, productID string, image *profile.File) error
	RemoveImageInProduct(ctx context.Context, companyID string, productID string) error

	AddCompanyService(ctx context.Context, companyID string, service *profile.Service) error
	DeleteCompanyService(ctx context.Context, companyID string, serviceID string) error
	ChangeCompanyService(ctx context.Context, companyID string, service *profile.Service) error
	ChangeImageService(ctx context.Context, companyID string, serviceID string, image *profile.File) error
	RemoveImageInService(ctx context.Context, companyID string, serviceID string) error

	AddCompanyReport(ctx context.Context, companyID string, report *profile.Report) error

	ChangeAvatar(ctx context.Context, companyID string, image string) error
	ChangeOriginAvatar(ctx context.Context, companyID string, image string) error
	RemoveAvatar(ctx context.Context, companyID string) error
	GetOriginAvatar(ctx context.Context, companyID string) (string, error)

	ChangeCover(ctx context.Context, companyID string, image string) error
	ChangeOriginCover(ctx context.Context, companyID string, image string) error
	RemoveCover(ctx context.Context, companyID string) error
	GetOriginCover(ctx context.Context, companyID string) (string, error)

	SaveCompanyProfileTranslation(ctx context.Context, companyID string, lang string, tr *profile.Translation) error
	SaveCompanyMilestoneTranslation(ctx context.Context, companyID, milestoneID, language string, translation *profile.MilestoneTranslation) error
	SaveCompanyAwardTranslation(ctx context.Context, companyID, awardID, language string, translation *profile.AwardTranslation) error
	OpenCareerCenter(ctx context.Context, companyID string, cc *careercenter.CareerCenter) error
}

// CacheRepository contains functions which have to be in cache repository
type CacheRepository interface {
	CreateTemporaryCodeForEmailActivation(ctx context.Context, companyID string, email string) (string, error)
	CheckTemporaryCodeForEmailActivation(ctx context.Context, companyID string, code string) (bool, string, error)
	Remove(ctx context.Context, key string) error
}

// ReviewsRepository contains functions which have to be in reviews repository
type ReviewsRepository interface {
	AddReview(ctx context.Context, review *profile.Review) error
	DeleteCompanyReview(ctx context.Context, companyID string, userID string, reviewID string) error
	GetCompanyReviews(ctx context.Context, companyID string, first uint32, after uint32) ([]*profile.Review, error)
	GetUsersRevies(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Review, error)
	GetAvarageRateOfCompany(ctx context.Context, companyID string) (float32, uint32, error)
	GetAmountOfEachRate(ctx context.Context, companyID string) (map[uint32]uint32, error)
	AddCompanyReviewReport(ctx context.Context, reviewReport *profile.ReviewReport) error
	GetAmountOfReviewsOfUser(ctx context.Context, userID string) (int32, error)
}

// ArangoRepo ...
type ArangoRepo interface {
	SaveCompany(ctx context.Context, c *arangorepo.Company) error
}
