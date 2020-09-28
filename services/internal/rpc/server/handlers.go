package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
)

// CheckIfURLForVOfficeIsTaken checks if the URL of virtual office is taken
func (s Server) CheckIfURLForVOfficeIsTaken(ctx context.Context, data *servicesRPC.URL) (*servicesRPC.BooleanValue, error) {
	b, err := s.service.CheckIfURLForVOfficeIsTaken(ctx, data.GetURL())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.BooleanValue{
		Value: b,
	}, nil
}

// CreateVOffice creates VOffice
func (s Server) CreateVOffice(ctx context.Context, data *servicesRPC.CreateVOfficeRequest) (*servicesRPC.ID, error) {
	id, err := s.service.CreateVOffice(ctx, data.GetCompanyID(), createOfficeRequestRPCToOfficeAccount(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// ChangeVOffice ...
func (s Server) ChangeVOffice(ctx context.Context, data *servicesRPC.CreateVOfficeRequest) (*servicesRPC.Empty, error) {

	err := s.service.ChangeVOffice(ctx, data.GetCompanyID(), data.GetID(), createOfficeRequestRPCToOfficeAccount(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// RemoveVOffice ...
func (s Server) RemoveVOffice(ctx context.Context, data *servicesRPC.RemoveVOfficeRequest) (*servicesRPC.Empty, error) {

	err := s.service.RemoveVOffice(ctx, data.GetCompanyID(), data.GetOfficeID())

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// ChangeVofficeCover ...
func (s Server) ChangeVofficeCover(ctx context.Context, data *servicesRPC.File) (*servicesRPC.ID, error) {
	id, err := s.service.ChangeVofficeCover(ctx, data.GetTargetID(), data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// ChangeVofficeOriginCover ...
func (s Server) ChangeVofficeOriginCover(ctx context.Context, data *servicesRPC.File) (*servicesRPC.ID, error) {
	id, err := s.service.ChangeVofficeOriginCover(ctx, data.GetTargetID(), data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}
	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// RemoveVofficeCover ...
func (s Server) RemoveVofficeCover(ctx context.Context, data *servicesRPC.RemoveCover) (*servicesRPC.Empty, error) {

	err := s.service.RemoveVofficeCover(ctx, data.GetID(), data.GetCompanyID())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// AddFileInVofficeService ...
func (s Server) AddFileInVofficeService(ctx context.Context, data *servicesRPC.File) (*servicesRPC.ID, error) {
	id, err := s.service.AddFileInVofficeService(ctx, data.GetTargetID(), data.GetItemID(), data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// AddFileInServiceRequest ...
func (s Server) AddFileInServiceRequest(ctx context.Context, data *servicesRPC.File) (*servicesRPC.ID, error) {
	id, err := s.service.AddFileInServiceRequest(ctx, data.GetTargetID(), data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// AddFileInOrderService ...
func (s Server) AddFileInOrderService(ctx context.Context, data *servicesRPC.File) (*servicesRPC.ID, error) {
	id, err := s.service.AddFileInOrderService(ctx, data.GetTargetID(), data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// GetVOffice gets the voffice for user OR company
func (s Server) GetVOffice(ctx context.Context, data *servicesRPC.GetOfficeRequest) (*servicesRPC.VOffices, error) {
	office, err := s.service.GetVOffice(ctx, data.GetCompanyID(), data.GetUserID())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.VOffices{
		VOffice: officesToOfficesRPC(office),
	}, nil
}

// GetVOfficeByID ...
func (s Server) GetVOfficeByID(ctx context.Context, data *servicesRPC.GetVOfficeByIDRequest) (*servicesRPC.VOffice, error) {

	office, err := s.service.GetVOfficeByID(ctx, data.GetCompanyID(), data.GetOfficeID())
	if err != nil {
		return nil, err
	}

	return officeToOfficeRPC(office), nil
}

// GetServicesRequest ...
func (s Server) GetServicesRequest(ctx context.Context, data *servicesRPC.GerServiceRequst) (*servicesRPC.GetServicesResponse, error) {
	service, err := s.service.GetServicesRequest(ctx, data.GetOwnerID(), data.GetCompanyID())
	if err != nil {
		return nil, err
	}

	req := servicesRequestToRPC(service)

	return &servicesRPC.GetServicesResponse{
		Request: req,
	}, nil

}

// GetServiceRequest ...
func (s Server) GetServiceRequest(ctx context.Context, data *servicesRPC.GerServiceRequst) (*servicesRPC.GetServiceResponse, error) {
	service, err := s.service.GetServiceRequest(ctx, data.GetCompanyID(), data.GetServiceID())
	if err != nil {
		return nil, err
	}

	req := serviceRequestToRPC(service)

	return &servicesRPC.GetServiceResponse{
		Request: req,
	}, nil

}

// ChangeServicesRequestStatus ...
func (s Server) ChangeServicesRequestStatus(ctx context.Context, data *servicesRPC.GerServiceRequst) (*servicesRPC.Empty, error) {

	if data == nil {
		return nil, nil
	}

	err := s.service.ChangeServicesRequestStatus(ctx,
		data.GetCompanyID(),
		data.GetServiceID(),
		data.GetServiceStatus().String(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// IsOutOfOffice ...
func (s Server) IsOutOfOffice(ctx context.Context, data *servicesRPC.IsOutOfOfficeRequest) (*servicesRPC.Empty, error) {

	err := s.service.IsOutOfOffice(ctx,
		data.GetCompanyID(),
		data.GetOfficeID(),
		data.GetIsOut(),
		pointerDayMonthAndYearToTime(data.GetReturnDate()),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// ChangeVOfficeName changes the name of vOffice
func (s Server) ChangeVOfficeName(ctx context.Context, data *servicesRPC.ChangeVOfficeNameRequest) (*servicesRPC.Empty, error) {
	err := s.service.ChangeVOfficeName(ctx, data.GetCompanyID(), data.GetVOfficeID(), data.GetName())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// AddChangeVOfficeDescription adds or changes the description of vOffice depending what front-end is passing
func (s Server) AddChangeVOfficeDescription(ctx context.Context, data *servicesRPC.AddChangeDescriptionRequest) (*servicesRPC.Empty, error) {
	err := s.service.AddChangeVOfficeDescription(ctx, data.GetCompanyID(), data.GetVOfficeID(), data.GetDescription())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// AddVOfficePortfolio adds  the vOffice portfolio depending on what front-end is passing.
func (s Server) AddVOfficePortfolio(ctx context.Context, data *servicesRPC.AddChangeVOfficePortfolioPortfolioRequest) (*servicesRPC.VofficeAddPortfolioResponse, error) {
	id, links, err := s.service.AddVOfficePortfolio(ctx, data.GetCompanyID(), data.GetOfficeID(), serviceRPCPortfolioToPortfolio(data.GetPortfolio()))
	if err != nil {
		return nil, err
	}

	result := &servicesRPC.VofficeAddPortfolioResponse{}
	result.ID = id

	for _, link := range links {
		result.Links = append(result.Links, &servicesRPC.Link{
			ID:  link.GetID(),
			URL: link.URL,
		})
	}
	return result, nil
}

// ChangeVOfficePortfolio changes the vOffice portfolio depending on what front-end is passing.
func (s Server) ChangeVOfficePortfolio(ctx context.Context, data *servicesRPC.AddChangeVOfficePortfolioPortfolioRequest) (*servicesRPC.ID, error) {
	id, err := s.service.ChangeVOfficePortfolio(ctx, data.GetCompanyID(), data.GetOfficeID(), data.GetID(), serviceRPCPortfolioToPortfolio(data.GetPortfolio()))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// AddFileInVOfficePortfolio ...
func (s Server) AddFileInVOfficePortfolio(ctx context.Context, data *servicesRPC.File) (*servicesRPC.ID, error) {
	id, err := s.service.AddFileInVOfficePortfolio(ctx, data.GetTargetID(), data.GetItemID(), data.GetCompanyID(), fileRPCToProfileFile(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// RemoveLinksInVOfficePortfolio removes links which are in portfolio
func (s Server) RemoveLinksInVOfficePortfolio(ctx context.Context, data *servicesRPC.RemoveLinksRequest) (*servicesRPC.Empty, error) {
	ids := make([]string, 0, len(data.GetLinks()))
	for i := range data.GetLinks() {
		ids = append(ids, data.GetLinks()[i].GetID())
	}

	err := s.service.RemoveLinksInVOfficePortfolio(ctx, data.GetCompanyID(), data.GetFirstTargetID(), data.GetSecondTargetID(), ids)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// RemoveVOfficePortfolio removes the whole portfolio from vOffice
func (s Server) RemoveVOfficePortfolio(ctx context.Context, data *servicesRPC.RemoveVOfficePortfolioRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}
	err := s.service.RemoveVOfficePortfolio(ctx, data.GetCompanyID(), data.GetVOfficeID(), data.GetPortfolioID())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// RemoveFilesInVOfficePortfolio removes files in the desired Portfolio
func (s Server) RemoveFilesInVOfficePortfolio(ctx context.Context, data *servicesRPC.RemoveFilesRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInVOfficePortfolio(ctx, data.GetCompanyID(), data.GetFirstTargetID(), data.GetSecondTargetID(), ids)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// AddVOfficeLanguages ...
func (s Server) AddVOfficeLanguages(ctx context.Context, data *servicesRPC.ChangeVOfficeQualificationsRequest) (*servicesRPC.IDs, error) {
	if data == nil {
		return nil, nil
	}

	ids, err := s.service.AddVOfficeLanguages(ctx, data.GetCompanyID(), data.GetOfficeID(), serviceRPCLanguagesToLanguages(data.GetLanguages()))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.IDs{
		IDs: ids,
	}, nil
}

// ChangeVOfficeLanguage  changes qualifications of the office
func (s Server) ChangeVOfficeLanguage(ctx context.Context, data *servicesRPC.ChangeVOfficeQualificationsRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.ChangeVOfficeLanguage(ctx, data.GetCompanyID(), data.GetOfficeID(), serviceRPCLanguagesToLanguages(data.GetLanguages()))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// RemoveVOfficeLanguages removes languages in qualfiications
func (s Server) RemoveVOfficeLanguages(ctx context.Context, data *servicesRPC.RemoveLanguagesRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.RemoveVOfficeLanguages(ctx, data.GetCompanyID(), data.GetOfficeID(), data.GetID())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// GetVOfficeServices ....
func (s Server) GetVOfficeServices(ctx context.Context, data *servicesRPC.GetVOfficeServicesRequest) (*servicesRPC.GetVOfficeServicesResponse, error) {
	services, err := s.service.GetVOfficeServices(ctx, data.GetCompanyID(), data.GetOfficeID())

	if err != nil {
		return nil, err
	}

	se := serviceArrayToServiceArrayRPC(services)

	return &servicesRPC.GetVOfficeServicesResponse{
		Services: se,
	}, nil

}

// GetAllServices ....
func (s Server) GetAllServices(ctx context.Context, data *servicesRPC.GetVOfficeServicesRequest) (*servicesRPC.GetVOfficeServicesResponse, error) {
	services, err := s.service.GetAllServices(ctx, data.GetCompanyID())

	if err != nil {
		return nil, err
	}

	se := serviceArrayToServiceArrayRPC(services)

	return &servicesRPC.GetVOfficeServicesResponse{
		Services: se,
	}, nil

}

// GetVOfficeService ....
func (s Server) GetVOfficeService(ctx context.Context, data *servicesRPC.GetVOfficeServicesRequest) (*servicesRPC.Service, error) {
	service, err := s.service.GetVOfficeService(ctx, data.GetCompanyID(), data.GetOfficeID(), data.GetServiceID())

	if err != nil {
		return nil, err
	}

	se := serviceRPCToServiceRPC(service)

	return se, nil

}

// AddVOfficeService adds o the service that is posted by voffice
func (s Server) AddVOfficeService(ctx context.Context, data *servicesRPC.AddChangeVOfficeServiceRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}

	id, err := s.service.AddVOfficeService(ctx, data.GetCompanyID(), data.GetVOfficeID(), servicesRPCServiceToServiceStruct(data.GetService()))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// ChangeVOfficeService achanges  the service that is posted by voffice
func (s Server) ChangeVOfficeService(ctx context.Context, data *servicesRPC.AddChangeVOfficeServiceRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}

	id, err := s.service.ChangeVOfficeService(ctx, data.GetCompanyID(), data.GetID(), data.GetVOfficeID(), servicesRPCServiceToServiceStruct(data.GetService()))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// ChangeVOfficeServiceStatus ...
func (s Server) ChangeVOfficeServiceStatus(ctx context.Context, data *servicesRPC.ChangeVOfficeServiceStatusRequest) (*servicesRPC.Empty, error) {

	if data == nil {
		return nil, nil
	}

	err := s.service.ChangeVOfficeServiceStatus(ctx,
		data.GetCompanyID(),
		data.GetOfficeID(),
		data.GetServiceID(),
		data.GetServiceStatus().String(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// RemoveVOfficeService removes the service
func (s Server) RemoveVOfficeService(ctx context.Context, data *servicesRPC.RemoveRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	_, err := s.service.RemoveVOfficeService(ctx, data.GetCompanyID(), data.GetTargetID())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// RemoveFilesInVOfficeService removes files in the desired service
func (s Server) RemoveFilesInVOfficeService(ctx context.Context, data *servicesRPC.RemoveFilesRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInVOfficeService(ctx, data.GetCompanyID(), data.GetFirstTargetID(), ids)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// RemoveFilesInServiceRequest ...
func (s Server) RemoveFilesInServiceRequest(ctx context.Context, data *servicesRPC.RemoveFilesRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	ids := make([]string, 0, len(data.GetFiles()))
	for i := range data.GetFiles() {
		ids = append(ids, data.GetFiles()[i].GetID())
	}

	err := s.service.RemoveFilesInVOfficeService(ctx, data.GetCompanyID(), data.GetFirstTargetID(), ids)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// AddServicesRequest adds or changes the request that may be posted by user or company
func (s Server) AddServicesRequest(ctx context.Context, data *servicesRPC.AddServiceRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}

	id, err := s.service.AddServicesRequest(ctx, data.GetCompanyID(), requestRPCToRequest(data.GetRequest()))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// ChangeServicesRequest ...
func (s Server) ChangeServicesRequest(ctx context.Context, data *servicesRPC.ChangeServicesRequestRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}

	id, err := s.service.ChangeServicesRequest(ctx, data.GetCompanyID(), data.GetServiceID(), requestRPCToRequest(data.GetRequest()))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// RemoveServicesRequest removes the service
func (s Server) RemoveServicesRequest(ctx context.Context, data *servicesRPC.RemoveRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.RemoveServicesRequest(ctx, data.GetCompanyID(), data.GetTargetID())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// OrderService ...
func (s Server) OrderService(ctx context.Context, data *servicesRPC.OrderServiceRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}

	id, err := s.service.OrderService(ctx, orderRPCToOrder(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// AcceptOrderService ...
func (s Server) AcceptOrderService(ctx context.Context, data *servicesRPC.AcceptOrderServiceRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.AcceptOrderService(ctx, data.GetCompanyID(), data.GetServiceID(), data.GetOrderID())
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// CancelServiceOrder ...
func (s Server) CancelServiceOrder(ctx context.Context, data *servicesRPC.CancelServiceOrderRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.CancelServiceOrder(ctx,
		data.GetCompanyID(),
		data.GetOrderID(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// DeclineServiceOrder ...
func (s Server) DeclineServiceOrder(ctx context.Context, data *servicesRPC.CancelServiceOrderRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.DeclineServiceOrder(ctx,
		data.GetCompanyID(),
		data.GetOrderID(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// DeliverServiceOrder ...
func (s Server) DeliverServiceOrder(ctx context.Context, data *servicesRPC.CancelServiceOrderRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.DeliverServiceOrder(ctx,
		data.GetCompanyID(),
		data.GetOrderID(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// AcceptDeliverdServiceOrder ...
func (s Server) AcceptDeliverdServiceOrder(ctx context.Context, data *servicesRPC.CancelServiceOrderRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.AcceptDeliverdServiceOrder(ctx,
		data.GetCompanyID(),
		data.GetOrderID(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// CancelDeliverdServiceOrder ...
func (s Server) CancelDeliverdServiceOrder(ctx context.Context, data *servicesRPC.CancelServiceOrderRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.CancelDeliverdServiceOrder(ctx,
		data.GetCompanyID(),
		data.GetOrderID(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// SendProposalForServiceRequest ...
func (s Server) SendProposalForServiceRequest(ctx context.Context, data *servicesRPC.SendProposalRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}

	id, err := s.service.SendProposalForServiceRequest(ctx, proposalRPCToProposal(data))
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{
		ID: id,
	}, nil
}

// OrderProposalForServiceRequest ...
func (s Server) OrderProposalForServiceRequest(ctx context.Context, data *servicesRPC.IgnoreProposalRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	proposal, err := s.service.GetProposalByID(
		ctx,
		data.GetCompanyID(),
		data.GetProposalID(),
	)
	if err != nil {
		return nil, err
	}

	err = s.service.OrderProposalForServiceRequest(ctx, proposalToOrderType(proposal))

	s.service.IgnoreProposalForServiceRequest(ctx, data.GetCompanyID(), data.GetProposalID())

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// IgnoreProposalForServiceRequest ...
func (s Server) IgnoreProposalForServiceRequest(ctx context.Context, data *servicesRPC.IgnoreProposalRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.IgnoreProposalForServiceRequest(ctx,
		data.GetCompanyID(),
		data.GetProposalID(),
	)

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// GetVOfficerServiceOrders ...
func (s Server) GetVOfficerServiceOrders(ctx context.Context, data *servicesRPC.GetVOfficerServiceOrdersRequest) (*servicesRPC.OrderServices, error) {
	if data == nil {
		return nil, nil
	}

	res, err := s.service.GetVOfficerServiceOrders(ctx,
		data.GetOwnerID(),
		data.GetOfficeID(),
		orderTypeRPCToOrderType(data.GetOrderType()),
		orderStatusRPCToOrderStatus(data.GetOrderStatus()),
		data.GetFirst(),
		data.GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.OrderServices{
		OrdersAmount:  res.OrderAmount,
		OrderServices: ordersToRPC(res.Orders),
	}, nil
}

// GetReceivedProposals ...
func (s Server) GetReceivedProposals(ctx context.Context, data *servicesRPC.GetReceivedProposalsRequest) (*servicesRPC.ProposalsResponse, error) {
	if data == nil {
		return nil, nil
	}

	res, err := s.service.GetReceivedProposals(ctx,
		data.GetCompanyID(),
		data.GetRequestID(),
		data.GetFirst(),
		data.GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ProposalsResponse{
		ProposalsAmount: res.ProposalAmount,
		Proposals:       proposalsToRPC(res.Proposals),
	}, nil
}

// GetSendedProposals ...
func (s Server) GetSendedProposals(ctx context.Context, data *servicesRPC.GetReceivedProposalsRequest) (*servicesRPC.ProposalsResponse, error) {
	if data == nil {
		return nil, nil
	}

	res, err := s.service.GetSendedProposals(ctx,
		data.GetFirst(),
		data.GetAfter(),
		data.GetCompanyID(),
	)
	if err != nil {
		return nil, err
	}

	return &servicesRPC.ProposalsResponse{
		ProposalsAmount: res.ProposalAmount,
		Proposals:       proposalsToRPC(res.Proposals),
	}, nil
}

// SaveVOfficeService ...
func (s Server) SaveVOfficeService(ctx context.Context, data *servicesRPC.VOfficeServiceActionRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.SaveVOfficeService(ctx, data.GetCompanyID(), data.GetServiceID())

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// UnSaveVOfficeService ...
func (s Server) UnSaveVOfficeService(ctx context.Context, data *servicesRPC.VOfficeServiceActionRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.UnSaveVOfficeService(ctx, data.GetCompanyID(), data.GetServiceID())

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// SaveServiceRequest ...
func (s Server) SaveServiceRequest(ctx context.Context, data *servicesRPC.VOfficeServiceActionRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.SaveServiceRequest(ctx, data.GetCompanyID(), data.GetServiceID())

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// UnSaveServiceRequest ...
func (s Server) UnSaveServiceRequest(ctx context.Context, data *servicesRPC.VOfficeServiceActionRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.UnSaveServiceRequest(ctx, data.GetCompanyID(), data.GetServiceID())

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// GetSavedVOfficeServices ...
func (s Server) GetSavedVOfficeServices(ctx context.Context, data *servicesRPC.GetSavedVOfficeServicesRequest) (*servicesRPC.GetVOfficeServicesResponse, error) {

	res, err := s.service.GetSavedVOfficeServices(ctx, data.GetCompanyID(), data.GetFirst(), data.GetAfter())

	if err != nil {
		return nil, err
	}

	se := serviceArrayToServiceArrayRPC(res.Services)

	return &servicesRPC.GetVOfficeServicesResponse{
		Services:      se,
		ServiceAmount: res.ServiceAmount,
	}, nil
}

// GetSavedServicesRequest ...
func (s Server) GetSavedServicesRequest(ctx context.Context, data *servicesRPC.GetSavedVOfficeServicesRequest) (*servicesRPC.GetServicesResponse, error) {

	res, err := s.service.GetSavedServicesRequest(ctx, data.GetCompanyID(), data.GetFirst(), data.GetAfter())

	if err != nil {
		return nil, err
	}

	req := servicesRequestToRPC(res.Services)

	return &servicesRPC.GetServicesResponse{
		Request:       req,
		ServiceAmount: res.ServiceAmount,
	}, nil
}

// AddNoteForOrderService ...
func (s Server) AddNoteForOrderService(ctx context.Context, data *servicesRPC.AddNoteForOrderServiceRequest) (*servicesRPC.Empty, error) {
	if data == nil {
		return nil, nil
	}

	err := s.service.AddNoteForOrderService(ctx, data.GetOrderID(), data.GetCompanyID(), data.GetText())

	if err != nil {
		return nil, err
	}

	return &servicesRPC.Empty{}, nil
}

// WriteReviewForService ...
func (s Server) WriteReviewForService(ctx context.Context, data *servicesRPC.WriteReviewRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}

	id, err := s.service.WriteReviewForService(ctx, reviewServiceRPCToReviewService(data))

	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{ID: id}, nil
}

// WriteReviewForServiceRequest ...
func (s Server) WriteReviewForServiceRequest(ctx context.Context, data *servicesRPC.WriteReviewRequest) (*servicesRPC.ID, error) {
	if data == nil {
		return nil, nil
	}
	id, err := s.service.WriteReviewForServiceRequest(ctx, reviewServiceRequestRPCToReviewServiceRequest(data))

	if err != nil {
		return nil, err
	}

	return &servicesRPC.ID{ID: id}, nil
}

// GetServicesReview ...
func (s Server) GetServicesReview(ctx context.Context, data *servicesRPC.GetReviewRequest) (*servicesRPC.GetReviewResponse, error) {
	if data == nil {
		return nil, nil
	}
	res, err := s.service.GetServicesReview(ctx,
		data.GetOfficeID(),
		data.GetFirst(),
		data.GetAfter(),
	)

	if err != nil {
		return nil, err
	}

	return &servicesRPC.GetReviewResponse{
		ReviewAmount:     res.ReviewAmount,
		ServiceReviewAVG: reviewAVGToRPC(res),
		Reviews:          reviewsToRPC(res.Reviews),
	}, nil
}

// GetServicesRequestReview ...
func (s Server) GetServicesRequestReview(ctx context.Context, data *servicesRPC.GetReviewRequest) (*servicesRPC.GetReviewResponse, error) {
	if data == nil {
		return nil, nil
	}
	res, err := s.service.GetServicesRequestReview(ctx,
		data.GetOwnerID(),
		data.GetCompanyID(),
		data.GetFirst(),
		data.GetAfter(),
	)

	if err != nil {
		return nil, err
	}

	return &servicesRPC.GetReviewResponse{
		ReviewAmount:     res.ReviewAmount,
		ServiceReviewAVG: reviewAVGToRPC(res),
		Reviews:          reviewsToRPC(res.Reviews),
	}, nil
}
