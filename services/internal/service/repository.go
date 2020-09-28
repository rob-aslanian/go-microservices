package service

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/qualifications"
	serviceorder "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/service-order"

	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"

	offer "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/offers"
	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	office "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/v-office"

	notmes "gitlab.lan/Rightnao-site/microservices/services/internal/notification_messages"

	companyadmin "gitlab.lan/Rightnao-site/microservices/services/internal/company-admin"
	review "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/review"
)

// AuthRPC represents auth service
type AuthRPC interface {
	GetUserID(ctx context.Context) (string, error)
}

// NetworkRPC represents network service
type NetworkRPC interface {
	GetAdminLevel(ctx context.Context, companyID string) (companyadmin.AdminLevel, error)
}

// InfoRPC represents info service
type InfoRPC interface {
	GetUserCountry(ctx context.Context) (string, error)
	GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error)
}

// JobsRPC ...
type JobsRPC interface {
	GetAmountOfActiveJobsOfCompany(ctx context.Context, companyID string) (int32, error)
}

// MailRPC represents a Mail gRPC client
type MailRPC interface {
	SendEmail(ctx context.Context, email string, message string) error
}

// UserRPC represents a Info gRPC client
type UserRPC interface {
	CheckPassword(ctx context.Context, password string) error
}

// ChatRPC ...
type ChatRPC interface {
	IsLive(ctx context.Context, id string) (bool, error)
}

// ReviewsRepository contains functions which have to be in reviews repository
type ReviewsRepository interface {
	AddReview(ctx context.Context, feedback *review.Review) error
	DeleteOfficeReview(ctx context.Context, officeID string, userID string, feedbackID string) error
	GetOfficeReviews(ctx context.Context, officeID string, first uint32, after uint32) ([]*review.Review, error)
	GetUsersRevies(ctx context.Context, userID string, first uint32, after uint32) ([]*review.Review, error)
	GetAvarageRateOfOffice(ctx context.Context, officeID string) (float32, uint32, error)
	GetAmountOfEachRate(ctx context.Context, officeID string) (map[uint32]uint32, error)
	AddOfficeReviewReport(ctx context.Context, feedbackReport *review.ReviewReport) error
}

// ServicesRepository is ll about offices. creating it, changing it, posting services from it etc
type ServicesRepository interface {
	GetVOffice(ctx context.Context, companyID, userID string) ([]*office.Office, error)
	GetVOfficeByID(ctx context.Context, officeID string) (*office.Office, error)

	RemoveVOffice(ctx context.Context, officeID string) error
	IsURLBusy(ctx context.Context, url string) (bool, error)
	CreateVOffice(ctx context.Context, office *office.Office) error
	ChangeVOffice(ctx context.Context, office *office.Office) error

	GetVOfficeServices(ctx context.Context, officeID string) ([]*servicerequest.Service, error)
	GetAllServices(ctx context.Context, profileID string, isComapny bool) ([]*servicerequest.Service, error)
	GetVOfficeService(ctx context.Context, officeID, serviceID string) (*servicerequest.Service, error)
	GetVOfficeServiceByID(ctx context.Context, serviceID string) (*servicerequest.Service, error)

	ChangeVOfficeServiceStatus(ctx context.Context, officeID, serviceID string, serviceStatus servicerequest.ServiceStatus) error

	AddVOfficeService(ctx context.Context, service *servicerequest.Service) error
	ChangeVOfficeService(ctx context.Context, serviceID, officeID string, service *servicerequest.Service) error
	RemoveVOfficeService(ctx context.Context, serviceID string) error
	RemoveFilesInVOfficeService(ctx context.Context, serviceID string, fileIDs []string) error
	IsOutOfOffice(ctx context.Context, officeID string, isOut bool, returnTime *time.Time) error
	ChangeVOfficeName(ctx context.Context, officeID, name string) error
	AddVOfficePortfolio(ctx context.Context, officeID string, portfolio *office.Portfolio, userID string, companyID string) error
	GetVOfficePortfolio(ctx context.Context, officeID string) ([]*servicerequest.Portfolio, error)
	ChangeVOfficePortfolio(ctx context.Context, officeID, portfolioID string, portfolio *office.Portfolio) error
	RemoveLinksInVOfficePortfolio(tx context.Context, officeID, portfolioID string, linkIDs []string) error
	RemoveVOfficePortfolio(ctx context.Context, officeID, portfolioID string) error

	AddChangeVOfficeDescription(ctx context.Context, officeID, description string) error
	RemoveFilesInVOfficePortfolio(ctx context.Context, officeID, portfolioID string, filesIDs []string) error

	AddVOfficeLanguages(ctx context.Context, userID, companyID, officeID string, langs []*qualifications.Language) error
	ChangeVOfficeLanguage(ctx context.Context, officeID string, langs []*qualifications.Language) error
	RemoveVOfficeLanguages(ctx context.Context, officeID string, languageIds []string) error

	AddServicesRequest(ctx context.Context, companyID string, request *servicerequest.Request) error
	RemoveFilesInServiceRequest(ctx context.Context, serviceID string, ids []string) error
	ChangeServicesRequest(ctx context.Context, request *servicerequest.Request) error
	HasLikedService(ctx context.Context, profileID string, serviceID string, requestID string) (bool, error)
	HasOrderedService(ctx context.Context, profileID string, serviceID string, requestID string) (bool, error)

	GetServicesRequest(ctx context.Context, userID, companyID string) ([]*servicerequest.Request, error)
	GetServiceRequest(ctx context.Context, serviceID string) (*servicerequest.Request, error)
	ChangeServicesRequestStatus(ctx context.Context, serviceID string, serviceStatus servicerequest.ServiceReqestStatus) error
	RemoveServicesRequest(ctx context.Context, requestID string) error

	ChangeVofficeCover(ctx context.Context, officeID string, companyID string, file string) error
	ChangeVofficeOriginCover(ctx context.Context, officeID string, companyID string, file string) error
	RemoveVofficeCover(ctx context.Context, officeID string, companyID string) error

	AddFileInVofficeService(ctx context.Context, officeID, serviceID, companyID string, file *file.File) error
	AddFileInServiceRequest(ctx context.Context, serviceID, companyID string, file *file.File) error
	AddFileInOrderService(ctx context.Context, orderID string, file *file.File) error

	AddFileInVOfficePortfolio(ctx context.Context, officeID, portfolioID, companyID string, file *file.File) error

	AcceptOrderService(ctx context.Context, data *serviceorder.Order, oldID string) error
	DeclineServiceOrder(ctx context.Context, profileID, orderID string) error
	CancelServiceOrder(ctx context.Context, profileID, orderID string) error
	DeliverServiceOrder(ctx context.Context, profileID, orderID string) error
	AcceptDeliverdServiceOrder(ctx context.Context, profileID, orderID string) error
	CancelDeliverdServiceOrder(ctx context.Context, profileID, orderID string) error

	GetOrderByReferalID(ctx context.Context, referalID string) (*serviceorder.Order, error)
	OrderService(ctx context.Context, data *serviceorder.Order) error
	GetVOfficerServiceOrders(ctx context.Context, ownerID string, officeID string, orderType serviceorder.OrderType, orderstatus serviceorder.OrderStatus, first int, after int) (*serviceorder.GetOrder, error)
	GetVOfficerServiceOrderByID(ctx context.Context, orderID string) (*serviceorder.Order, error)

	GetSavedVOfficeServices(ctx context.Context, profileID string, first int, after int) ([]string, int32, error)
	AddNoteForOrderService(ctx context.Context, orderID, profileID, text string) error

	SaveVOfficeService(ctx context.Context, profileID, serviceID string) error
	UnSaveVOfficeService(ctx context.Context, profileID, serviceID string) error

	GetProposalByID(ctx context.Context, prodileID, proposalID string) (*offer.Proposal, error)
	GetReceivedProposals(ctx context.Context, companyID, requestID string, first int, after int) (*offer.GetProposal, error)
	GetSendedProposals(ctx context.Context, profileID string, first int, after int) (*offer.GetProposal, error)
	GetProposalAmount(ctx context.Context, requestID string) (int32, error)

	SendProposalForServiceRequest(ctx context.Context, proposal *offer.Proposal) error
	IgnoreProposalForServiceRequest(ctx context.Context, profileID, proposalID string) error

	GetSavedServicesRequest(ctx context.Context, profileID string, first int, after int) ([]string, int32, error)
	SaveServiceRequest(ctx context.Context, profileID, serviceID string) error
	UnSaveServiceRequest(ctx context.Context, profileID, serviceID string) error

	WriteReviewForService(ctx context.Context, review review.Review) error
	WriteReviewForServiceRequest(ctx context.Context, review review.Review) error

	GetServicesReview(ctx context.Context, profileID, offficeID string, first int, after int) (*review.GetReview, error)
}

// RequestRepository is all about requests. posting it, changing it, adding files etc
type RequestRepository interface {
}

// CacheRepository contains functions which have to be in cache repository
type CacheRepository interface {
	CreateTemporaryCodeForEmailActivation(ctx context.Context, companyID string, email string) (string, error)
	CheckTemporaryCodeForEmailActivation(ctx context.Context, companyID string, code string) (bool, string, error)
	Remove(ctx context.Context, key string) error
}

// MQ ...
type MQ interface {
	OrderService(targetID string, not *notmes.NewOrder) error
	SendProposalForServiceRequest(targetID string, not *notmes.NewProposal) error
}
