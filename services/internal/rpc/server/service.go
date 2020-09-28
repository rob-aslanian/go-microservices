package serverRPC

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/review"

	offer "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/offers"

	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	serviceorder "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/service-order"

	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/qualifications"
	office "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/v-office"
)

// Service define functions inside service
type Service interface {
	CheckIfURLForVOfficeIsTaken(ctx context.Context, url string) (bool, error)
	IsOutOfOffice(ctx context.Context, companyID, officeID string, isOut bool, returnDate *time.Time) error
	CreateVOffice(ctx context.Context, companyID string, vOffice *office.Office) (string, error)
	ChangeVOffice(ctx context.Context, companyID, officeID string, vOffice *office.Office) error

	RemoveVOffice(ctx context.Context, companyID, officeID string) error
	GetVOffice(ctx context.Context, companyID, userID string) ([]*office.Office, error)
	GetVOfficeByID(ctx context.Context, companyID, officeID string) (*office.Office, error)
	ChangeVOfficeName(ctx context.Context, companyID, voffficeID string, name string) error
	// ChangeVOfficeLocation(ctx context.Context, location location.Location) error
	AddChangeVOfficeDescription(ctx context.Context, companyID, officeID string, description string) error
	AddVOfficePortfolio(ctx context.Context, companyID, officeID string, portfolio *office.Portfolio) (string, []*file.Link, error)
	ChangeVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string, portfolio *office.Portfolio) (string, error)
	RemoveVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string) error

	AddFileInVofficeService(ctx context.Context, officeID, serviceID, companyID string, file *file.File) (string, error)
	AddFileInServiceRequest(ctx context.Context, serviceID, companyID string, file *file.File) (string, error)
	AddFileInOrderService(ctx context.Context, orderID, companyID string, file *file.File) (string, error)

	AddFileInVOfficePortfolio(ctx context.Context, officeID, portfolioID, companyID string, file *file.File) (string, error)

	RemoveFilesInVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string, ids []string) error
	AddVOfficeLanguages(ctx context.Context, companyID, officeID string, lang []*qualifications.Language) ([]string, error)
	ChangeVOfficeLanguage(ctx context.Context, companyID, officeID string, lang []*qualifications.Language) error
	RemoveVOfficeLanguages(ctx context.Context, companyID, officeID string, languageIDs []string) error
	GetVOfficeServices(ctx context.Context, companyID, officeID string) ([]*servicerequest.Service, error)
	GetAllServices(ctx context.Context, companyID string) ([]*servicerequest.Service, error)

	GetVOfficeService(ctx context.Context, companyID, officeID, serviceID string) (*servicerequest.Service, error)
	ChangeVOfficeServiceStatus(ctx context.Context, companyID, officeID, serviceID string, serviceStatus string) error
	AddVOfficeService(ctx context.Context, companyID, officeID string, service *servicerequest.Service) (string, error)
	ChangeVOfficeService(ctx context.Context, companyID, serviceID, officeID string, service *servicerequest.Service) (string, error)
	RemoveVOfficeService(ctx context.Context, companyID, serviceID string) (string, error)
	RemoveFilesInVOfficeService(ctx context.Context, companyID, serviceID string, ids []string) error
	RemoveFilesInServiceRequest(ctx context.Context, companyID, serviceID string, ids []string) error

	RemoveLinksInVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string, linksIDs []string) error

	ChangeVofficeCover(ctx context.Context, officeID string, companyID string, file *file.File) (string, error)
	ChangeVofficeOriginCover(ctx context.Context, officeID string, companyID string, file *file.File) (string, error)

	RemoveVofficeCover(ctx context.Context, officeID string, companyID string) error

	/// Service Request
	AddServicesRequest(ctx context.Context, companyID string, request *servicerequest.Request) (string, error)
	ChangeServicesRequest(ctx context.Context, companyID, serviceID string, request *servicerequest.Request) (string, error)
	GetServicesRequest(ctx context.Context, ownerID, companyID string) ([]*servicerequest.Request, error)
	GetServiceRequest(ctx context.Context, companyID string, serviceID string) (*servicerequest.Request, error)
	ChangeServicesRequestStatus(ctx context.Context, companyID, serviceID string, serviceStatus string) error
	RemoveServicesRequest(ctx context.Context, companyID, requestID string) error

	OrderService(ctx context.Context, data *serviceorder.Order) (string, error)
	OrderProposalForServiceRequest(ctx context.Context, data *serviceorder.Order) error
	AcceptOrderService(ctx context.Context, companyID, serviceID, orderID string) error
	DeclineServiceOrder(ctx context.Context, companyID, orderID string) error
	CancelServiceOrder(ctx context.Context, companyID, orderID string) error
	DeliverServiceOrder(ctx context.Context, companyID, orderID string) error
	AcceptDeliverdServiceOrder(ctx context.Context, companyID, orderID string) error
	CancelDeliverdServiceOrder(ctx context.Context, companyID, orderID string) error

	GetVOfficerServiceOrders(ctx context.Context, ownerID string, officeID string, orderType serviceorder.OrderType, orderstatus serviceorder.OrderStatus, first uint32, after string) (*serviceorder.GetOrder, error)
	AddNoteForOrderService(ctx context.Context, orderID, companyID, text string) error

	GetSavedVOfficeServices(ctx context.Context, companyID string, first uint32, after string) (*servicerequest.GetServices, error)
	GetProposalByID(ctx context.Context, companyID, proposalID string) (*offer.Proposal, error)

	SaveVOfficeService(ctx context.Context, companyID, serviceID string) error
	UnSaveVOfficeService(ctx context.Context, companyID, serviceID string) error

	GetReceivedProposals(ctx context.Context, companyID, requestID string, first uint32, after string) (*offer.GetProposal, error)
	GetSendedProposals(ctx context.Context, first uint32, after, companyID string) (*offer.GetProposal, error)

	SendProposalForServiceRequest(ctx context.Context, proposal *offer.Proposal) (string, error)
	IgnoreProposalForServiceRequest(ctx context.Context, companyID, proposalID string) error

	GetSavedServicesRequest(ctx context.Context, companyID string, first uint32, after string) (*servicerequest.GetServicesRequest, error)
	SaveServiceRequest(ctx context.Context, companyID, serviceID string) error
	UnSaveServiceRequest(ctx context.Context, companyID, serviceID string) error

	WriteReviewForService(ctx context.Context, review review.Review) (string, error)
	WriteReviewForServiceRequest(ctx context.Context, review review.Review) (string, error)
	GetServicesReview(ctx context.Context, offficeID string, first uint32, after string) (*review.GetReview, error)
	GetServicesRequestReview(ctx context.Context, ownerID, companyID string, first uint32, after string) (*review.GetReview, error)
}
