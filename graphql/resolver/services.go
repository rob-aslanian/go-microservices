package resolver

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func (_ *Resolver) CreateVOffice(ctx context.Context, in CreateVOfficeRequest) (*SuccessResolver, error) {

	response, err := services.CreateVOffice(ctx, &servicesRPC.CreateVOfficeRequest{
		CompanyID:   NullToString(in.Company_id),
		Name:        in.Input.Name,
		Description: NullToString(in.Input.Description),
		Category:    in.Input.Category,
		Languages:   languagesInputArrayToServicesRPCLanguages(in.Input.Languages),
		Location:    locationToServicesRPCLocation(&in.Input.Location),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      response.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) ChangeVOffice(ctx context.Context, in ChangeVOfficeRequest) (*SuccessResolver, error) {

	_, err := services.ChangeVOffice(ctx, &servicesRPC.CreateVOfficeRequest{
		ID:          in.Office_id,
		CompanyID:   NullToString(in.Company_id),
		Name:        in.Input.Name,
		Description: NullToString(in.Input.Description),
		Category:    in.Input.Category,
		Languages:   languagesInputArrayToServicesRPCLanguages(in.Input.Languages),
		Location:    locationToServicesRPCLocation(&in.Input.Location),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveVOffice(ctx context.Context, in RemoveVOfficeRequest) (*SuccessResolver, error) {

	_, err := services.RemoveVOffice(ctx, &servicesRPC.RemoveVOfficeRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  in.Office_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) ChangeVOfficeName(ctx context.Context, in ChangeVOfficeNameRequest) (*SuccessResolver, error) {
	_, err := services.ChangeVOfficeName(ctx, &servicesRPC.ChangeVOfficeNameRequest{
		CompanyID: NullToString(in.Company_id),
		VOfficeID: in.Office_id,
		Name:      in.Name,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) GetVOffice(ctx context.Context, in GetVOfficeRequest) (*VOfficesResolver, error) {
	office, err := services.GetVOffice(ctx, &servicesRPC.GetOfficeRequest{
		CompanyID: NullToString(in.Company_id),
		UserID:    NullToString(in.User_id),
	})

	if err != nil {
		return nil, err
	}

	return &VOfficesResolver{
		R: &VOffices{
			V_offices: vOfficesRPCToVOffices(office.GetVOffice()),
		},
	}, nil
}

func (_ *Resolver) GetVOfficeByID(ctx context.Context, in GetVOfficeByIDRequest) (*VOfficeResolver, error) {
	office, err := services.GetVOfficeByID(ctx, &servicesRPC.GetVOfficeByIDRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  in.Office_id,
	})

	if err != nil {
		return nil, err
	}

	result := vOfficeRPCToVOffice(office)

	return &VOfficeResolver{
		R: &result,
	}, nil
}

func (_ *Resolver) IsOutOfOffice(ctx context.Context, in IsOutOfOfficeRequest) (*SuccessResolver, error) {
	_, err := services.IsOutOfOffice(ctx, &servicesRPC.IsOutOfOfficeRequest{
		CompanyID:  NullToString(in.Company_id),
		OfficeID:   in.Office_id,
		IsOut:      in.Is_Out,
		ReturnDate: NullToString(in.Return_Date),
	})
	if err != nil {
		return nil, err
	}
	sucess := &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}

	return sucess, nil
}

func (_ *Resolver) AddChangeVOfficeDescription(ctx context.Context, in AddChangeVOfficeDescriptionRequest) (*SuccessResolver, error) {
	_, err := services.AddChangeVOfficeDescription(ctx, &servicesRPC.AddChangeDescriptionRequest{
		VOfficeID:   in.Office_id,
		CompanyID:   NullToString(in.Company_id),
		Description: in.Description,
	})
	if err != nil {
		return nil, err
	}

	success := &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}

	return success, nil
}

func (_ *Resolver) RemoveVofficeCover(ctx context.Context, in RemoveVofficeCoverRequest) (*SuccessResolver, error) {
	_, err := services.RemoveVofficeCover(ctx, &servicesRPC.RemoveCover{
		ID:        in.Office_id,
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	success := &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}

	return success, nil
}

func (_ *Resolver) AddVOfficePortfolio(ctx context.Context, in AddVOfficePortfolioRequest) (*AddVofficeResponseResolver, error) {
	response, err := services.AddVOfficePortfolio(ctx, &servicesRPC.AddChangeVOfficePortfolioPortfolioRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  in.Office_id,
		Portfolio: &servicesRPC.Portfolio{
			Tittle:      in.Portfolio.Tittle,
			Description: in.Portfolio.Description,
			ContentType: stringToServicesRPCContentType(in.Portfolio.ContentType),
			Files:       files(in.Portfolio.Files_id),
			Links:       links(in.Portfolio.Links),
		},
	})
	if err != nil {
		return nil, err
	}

	success := &AddVofficeResponseResolver{
		R: &AddVofficeResponse{
			Success: true,
		},
	}

	if response.GetID() != "" {
		success.R.ID = response.GetID()
	}

	for _, link := range response.GetLinks() {
		success.R.Links = append(success.R.Links, &Link{
			ID:      link.GetID(),
			Address: link.GetURL(),
		})
	}

	return success, nil
}

func (_ *Resolver) ChangeVOfficePortfolio(ctx context.Context, in ChangeVOfficePortfolioRequest) (*SuccessResolver, error) {
	response, err := services.ChangeVOfficePortfolio(ctx, &servicesRPC.AddChangeVOfficePortfolioPortfolioRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  in.Office_id,
		ID:        NullToString(in.Portfolio_id),
		Portfolio: &servicesRPC.Portfolio{
			Tittle:      in.Portfolio.Tittle,
			Description: in.Portfolio.Description,
			ContentType: stringToServicesRPCContentType(in.Portfolio.ContentType),
			Files:       files(in.Portfolio.Files_id),
			Links:       links(in.Portfolio.Links),
		},
	})
	if err != nil {
		return nil, err
	}

	success := &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}

	if response.GetID() != "" {
		success.R.ID = response.GetID()
	}

	return success, nil
}

func (_ *Resolver) RemoveLinksInVOfficePortfolio(ctx context.Context, in RemoveLinksInVOfficePortfolioRequest) (*SuccessResolver, error) {
	_, err := services.RemoveLinksInVOfficePortfolio(ctx, &servicesRPC.RemoveLinksRequest{
		CompanyID:      NullToString(in.Company_id),
		FirstTargetID:  in.Office_id,
		SecondTargetID: in.Portfolio_id,
		Links:          linksString(in.Links_ids),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveFilesInVOfficePortfolio(ctx context.Context, in RemoveFilesInVOfficePortfolioRequest) (*SuccessResolver, error) {
	_, err := services.RemoveFilesInVOfficePortfolio(ctx, &servicesRPC.RemoveFilesRequest{
		CompanyID:      NullToString(in.Company_id),
		FirstTargetID:  in.Office_id,
		SecondTargetID: in.Portfolio_id,
		Files:          files(&in.Files_ids),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveVOfficePortfolio(ctx context.Context, in RemoveVOfficePortfolioRequest) (*SuccessResolver, error) {
	_, err := services.RemoveVOfficePortfolio(ctx, &servicesRPC.RemoveVOfficePortfolioRequest{
		CompanyID:   NullToString(in.Company_id),
		PortfolioID: in.Portfolio_id,
		VOfficeID:   in.Office_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) AddVOfficeLanguages(ctx context.Context, in ChangeVOfficeLanguageRequest) (*IDsResolver, error) {
	res, err := services.AddVOfficeLanguages(ctx, &servicesRPC.ChangeVOfficeQualificationsRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  in.Office_id,
		Languages: changeLanguagesInputArrayToServicesRPCLanguages(in.Languages),
	})
	if err != nil {
		return nil, err
	}

	return &IDsResolver{
		R: &IDs{
			Ids: res.GetIDs(),
		},
	}, nil
}

func (_ *Resolver) ChangeVOfficeLanguage(ctx context.Context, in ChangeVOfficeLanguageRequest) (*SuccessResolver, error) {
	_, err := services.ChangeVOfficeLanguage(ctx, &servicesRPC.ChangeVOfficeQualificationsRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  in.Office_id,
		Languages: changeLanguagesInputArrayToServicesRPCLanguages(in.Languages),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveVOfficeLanguages(ctx context.Context, in RemoveVOfficeLanguagesRequest) (*SuccessResolver, error) {
	_, err := services.RemoveVOfficeLanguages(ctx, &servicesRPC.RemoveLanguagesRequest{
		ID:        in.Language_ids,
		OfficeID:  in.Office_id,
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) GetVOfficeServices(ctx context.Context, in GetVOfficeServicesRequest) (*ServicesResolver, error) {

	vOfficeServices, err := services.GetVOfficeServices(ctx, &servicesRPC.GetVOfficeServicesRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  NullToString(in.Office_id),
	})
	if err != nil {
		return nil, err
	}

	return &ServicesResolver{
		R: servicesRPCServicesToServicesResolver(vOfficeServices),
	}, nil
}

func (_ *Resolver) GetAllServices(ctx context.Context, in GetAllServicesRequest) (*ServicesResolver, error) {

	vOfficeServices, err := services.GetAllServices(ctx, &servicesRPC.GetVOfficeServicesRequest{
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &ServicesResolver{
		R: servicesRPCServicesToServicesResolver(vOfficeServices),
	}, nil
}

func (_ *Resolver) GetVOfficeService(ctx context.Context, in GetVOfficeServiceRequest) (*ServiceResolver, error) {
	vOfficeService, err := services.GetVOfficeService(ctx, &servicesRPC.GetVOfficeServicesRequest{
		CompanyID: NullToString(in.Company_id),
		OfficeID:  NullToString(in.Office_id),
		ServiceID: in.Service_id,
	})
	if err != nil {
		return nil, err
	}

	res := serviceRPCToService(vOfficeService)

	return &ServiceResolver{
		R: &res,
	}, nil

}

func (_ *Resolver) GetServicesRequest(ctx context.Context, in GetServicesRequestRequest) (*[]ServiceRequestResolver, error) {

	serRequest, err := services.GetServicesRequest(ctx, &servicesRPC.GerServiceRequst{
		CompanyID: NullToString(in.Company_id),
		OwnerID:   NullToString(in.Owner_id),
	})
	if err != nil {
		return nil, err
	}

	reqs := make([]ServiceRequestResolver, 0, len(serRequest.GetRequest()))

	for i := range serRequest.GetRequest() {
		reqs = append(reqs, ServiceRequestResolver{
			R: serviceRequestRPCToServiceRequest(serRequest.GetRequest()[i]),
		})
	}

	return &reqs, nil

}

func (_ *Resolver) GetServiceRequest(ctx context.Context, in GetServiceRequestRequest) (*ServiceRequestResolver, error) {

	serRequest, err := services.GetServiceRequest(ctx, &servicesRPC.GerServiceRequst{
		CompanyID: NullToString(in.Company_id),
		ServiceID: in.Service_id,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceRequestResolver{
		R: serviceRequestRPCToServiceRequest(serRequest.GetRequest()),
	}, nil

}

func (_ *Resolver) ChangeServicesRequestStatus(ctx context.Context, in ChangeServicesRequestStatusRequest) (*SuccessResolver, error) {
	_, err := services.ChangeServicesRequestStatus(ctx, &servicesRPC.GerServiceRequst{
		CompanyID:     NullToString(in.Company_id),
		ServiceID:     in.Service_id,
		ServiceStatus: serviceServiceStatusToRPC(in.Status),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) AddVOfficeService(ctx context.Context, in AddVOfficeServiceRequest) (*SuccessResolver, error) {

	response, err := services.AddVOfficeService(ctx, &servicesRPC.AddChangeVOfficeServiceRequest{
		CompanyID: NullToString(in.Company_id),
		VOfficeID: NullToString(&in.Office_id),
		Service: &servicesRPC.Service{
			Tittle:            in.Service.Title,
			Currency:          in.Service.Currency,
			Category:          category(in.Service.Category),
			DeliveryTime:      stringToServiceRPCDeliveryTimeEnum(in.Service.Delivery_time),
			AdditionalDetails: additionalDetails(in.Service.Additional_details),
			Files:             files(in.Service.Files_id),
			OfficeID:          in.Service.Office_id,
			Description:       in.Service.Description,
			Price:             stringToServiceRPCPriceEnum(in.Service.Price),
			MinPriceAmmount:   NullToInt32(in.Service.Min_price_amount),
			MaxPriceAmmount:   NullToInt32(in.Service.Max_price_amount),
			FixedPriceAmmount: NullToInt32(in.Service.Fixed_price_amount),
			Location:          locationToServicesRPCLocation(in.Service.Location),
			LocationType:      locationTypeToServiceRPCLocationTypeEnum(in.Service.Location_type),
			IsDraft:           in.Service.Is_Draft,
			IsRemote:          in.Service.Is_Remote,
			WorkingDate:       workingDateToRPC(in.Service.Wokring_hour),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      response.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) ChangeVOfficeService(ctx context.Context, in ChangeVOfficeServiceRequest) (*SuccessResolver, error) {
	response, err := services.ChangeVOfficeService(ctx, &servicesRPC.AddChangeVOfficeServiceRequest{
		ID:        in.Service_id,
		CompanyID: NullToString(in.Company_id),
		VOfficeID: in.Office_id,
		Service: &servicesRPC.Service{
			Tittle:            in.Service.Title,
			Category:          category(in.Service.Category),
			DeliveryTime:      stringToServiceRPCDeliveryTimeEnum(in.Service.Delivery_time),
			AdditionalDetails: additionalDetails(in.Service.Additional_details),
			Files:             files(in.Service.Files_id),
			OfficeID:          in.Service.Office_id,
			Description:       in.Service.Description,
			Price:             stringToServiceRPCPriceEnum(in.Service.Price),
			MinPriceAmmount:   NullToInt32(in.Service.Min_price_amount),
			MaxPriceAmmount:   NullToInt32(in.Service.Max_price_amount),
			FixedPriceAmmount: NullToInt32(in.Service.Fixed_price_amount),
			Location:          locationToServicesRPCLocation(in.Service.Location),
			LocationType:      locationTypeToServiceRPCLocationTypeEnum(in.Service.Location_type),
			WorkingDate:       workingDateToRPC(in.Service.Wokring_hour),
			IsDraft:           in.Service.Is_Draft,
			Currency:          in.Service.Currency,
			IsRemote:          in.Service.Is_Remote,
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      response.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) ChangeVOfficeServiceStatus(ctx context.Context, in ChangeVOfficeServiceStatusRequest) (*SuccessResolver, error) {

	_, err := services.ChangeVOfficeServiceStatus(ctx, &servicesRPC.ChangeVOfficeServiceStatusRequest{
		CompanyID:     NullToString(in.Company_id),
		OfficeID:      in.Office_id,
		ServiceID:     in.Service_id,
		ServiceStatus: serviceServiceStatusToRPC(in.Status),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveVOfficeService(ctx context.Context, in RemoveVOfficeServiceRequest) (*SuccessResolver, error) {
	_, err := services.RemoveVOfficeService(ctx, &servicesRPC.RemoveRequest{
		CompanyID: NullToString(in.Company_id),
		TargetID:  in.Service_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveFilesInVOfficeService(ctx context.Context, in RemoveFilesInVOfficeServiceRequest) (*SuccessResolver, error) {
	_, err := services.RemoveFilesInVOfficeService(ctx, &servicesRPC.RemoveFilesRequest{
		CompanyID:     NullToString(in.Company_id),
		FirstTargetID: in.Service_id,
		Files:         files(&in.Files_ids),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveServicesRequest(ctx context.Context, in RemoveServicesRequestRequest) (*SuccessResolver, error) {
	_, err := services.RemoveServicesRequest(ctx, &servicesRPC.RemoveRequest{
		CompanyID: NullToString(in.Company_id),
		TargetID:  in.Request_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) OrderService(ctx context.Context, in OrderServiceRequest) (*SuccessResolver, error) {
	res, err := services.OrderService(ctx, &servicesRPC.OrderServiceRequest{
		OwnerID:        in.Input.Owner_id,
		IsOwnerCompany: in.Input.Is_owner_company,
		ServiceID:      in.Input.Service_id,
		OfficeID:       in.Input.Office_id,
		OrderDetail:    orderDetailToRPC(in.Input.Order_details),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) AcceptOrderService(ctx context.Context, in AcceptOrderServiceRequest) (*SuccessResolver, error) {
	_, err := services.AcceptOrderService(ctx, &servicesRPC.AcceptOrderServiceRequest{
		CompanyID: NullToString(in.Company_id),
		ServiceID: in.Service_id,
		OrderID:   in.Order_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) CancelServiceOrder(ctx context.Context, in CancelServiceOrderRequest) (*SuccessResolver, error) {
	_, err := services.CancelServiceOrder(ctx, &servicesRPC.CancelServiceOrderRequest{
		OrderID:   in.Order_id,
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) DeclineServiceOrder(ctx context.Context, in DeclineServiceOrderRequest) (*SuccessResolver, error) {
	_, err := services.DeclineServiceOrder(ctx, &servicesRPC.CancelServiceOrderRequest{
		OrderID:   in.Order_id,
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) DeliverServiceOrder(ctx context.Context, in DeliverServiceOrderRequest) (*SuccessResolver, error) {
	_, err := services.DeliverServiceOrder(ctx, &servicesRPC.CancelServiceOrderRequest{
		OrderID:   in.Order_id,
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) AcceptDeliverdServiceOrder(ctx context.Context, in AcceptDeliverdServiceOrderRequest) (*SuccessResolver, error) {
	_, err := services.AcceptDeliverdServiceOrder(ctx, &servicesRPC.CancelServiceOrderRequest{
		OrderID:   in.Order_id,
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) CancelDeliverdServiceOrder(ctx context.Context, in CancelDeliverdServiceOrderRequest) (*SuccessResolver, error) {
	_, err := services.CancelDeliverdServiceOrder(ctx, &servicesRPC.CancelServiceOrderRequest{
		OrderID:   in.Order_id,
		CompanyID: NullToString(in.Company_id),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) SendProposalForServiceRequest(ctx context.Context, in SendProposalForServiceRequestRequest) (*SuccessResolver, error) {
	res, err := services.SendProposalForServiceRequest(ctx, &servicesRPC.SendProposalRequest{
		OwnerID:        in.Input.Owner_id,
		RequestID:      in.Input.Request_id,
		IsOwnerCompany: in.Input.Is_owner_company,
		ProposalDetail: proposalDetailToRPC(in.Input.Proposal_detail),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) OrderProposalForServiceRequest(ctx context.Context, in OrderProposalForServiceRequestRequest) (*SuccessResolver, error) {
	_, err := services.OrderProposalForServiceRequest(ctx, &servicesRPC.IgnoreProposalRequest{
		CompanyID:  NullToString(in.Company_id),
		ProposalID: in.Proposal_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) IgnoreProposalForServiceRequest(ctx context.Context, in IgnoreProposalForServiceRequestRequest) (*SuccessResolver, error) {
	_, err := services.IgnoreProposalForServiceRequest(ctx, &servicesRPC.IgnoreProposalRequest{
		CompanyID:  NullToString(in.Company_id),
		ProposalID: in.Proposal_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) RemoveFilesInServiceRequest(ctx context.Context, in RemoveFilesInServiceRequestRequest) (*SuccessResolver, error) {
	_, err := services.RemoveFilesInServiceRequest(ctx, &servicesRPC.RemoveFilesRequest{
		CompanyID:     NullToString(in.Company_id),
		FirstTargetID: in.Service_id,
		Files:         files(&in.Files_ids),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) AddServicesRequest(ctx context.Context, in AddServicesRequestRequest) (*SuccessResolver, error) {
	response, err := services.AddServicesRequest(ctx, &servicesRPC.AddServiceRequest{
		CompanyID: NullToString(in.Company_id),
		Request: &servicesRPC.Request{
			Tittle:            in.Request.Title,
			Category:          category(in.Request.Category),
			Currency:          in.Request.Currency,
			DeliveryTime:      stringToServiceRPCDeliveryTimeEnum(in.Request.Delivery_time),
			Files:             files(in.Request.Files_id),
			Description:       in.Request.Description,
			Price:             stringToServiceRPCPriceEnum(in.Request.Price),
			MinPriceAmmount:   NullToInt32(in.Request.Min_price_amount),
			MaxPriceAmmount:   NullToInt32(in.Request.Max_price_amount),
			FixedPriceAmmount: NullToInt32(in.Request.Fixed_price_amount),
			Location:          locationToServicesRPCLocation(in.Request.Location),
			LocationType:      locationTypeToServiceRPCLocationTypeEnum(in.Request.Location_type),
			AdditionalDetails: requestAdditionalDetails(in.Request.Additional_details),
			IsDraft:           in.Request.Is_Draft,
			ProjectType:       stringToServiceRPCRequestProjectType(in.Request.Project_type),
			CustomDate:        NullToString(in.Request.Custom_date),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      response.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) ChangeServicesRequest(ctx context.Context, in ChangeServicesRequestRequest) (*SuccessResolver, error) {
	response, err := services.ChangeServicesRequest(ctx, &servicesRPC.ChangeServicesRequestRequest{
		CompanyID: NullToString(in.Company_id),
		ServiceID: in.Service_id,
		Request: &servicesRPC.Request{
			Tittle:            in.Request.Title,
			Category:          category(in.Request.Category),
			Currency:          in.Request.Currency,
			DeliveryTime:      stringToServiceRPCDeliveryTimeEnum(in.Request.Delivery_time),
			Files:             files(in.Request.Files_id),
			Description:       in.Request.Description,
			Price:             stringToServiceRPCPriceEnum(in.Request.Price),
			MinPriceAmmount:   NullToInt32(in.Request.Min_price_amount),
			MaxPriceAmmount:   NullToInt32(in.Request.Max_price_amount),
			FixedPriceAmmount: NullToInt32(in.Request.Fixed_price_amount),
			Location:          locationToServicesRPCLocation(in.Request.Location),
			LocationType:      locationTypeToServiceRPCLocationTypeEnum(in.Request.Location_type),
			AdditionalDetails: requestAdditionalDetails(in.Request.Additional_details),
			IsDraft:           in.Request.Is_Draft,
			ProjectType:       stringToServiceRPCRequestProjectType(in.Request.Project_type),
			CustomDate:        NullToString(in.Request.Custom_date),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      response.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) GetServiceOrders(ctx context.Context, in GetServiceOrdersRequest) (*ServiceOrdersResolver, error) {
	response, err := services.GetVOfficerServiceOrders(ctx, &servicesRPC.GetVOfficerServiceOrdersRequest{
		OwnerID:     in.Owner_id,
		OfficeID:    NullToString(in.Office_id),
		OrderType:   orderTypeToRPC(in.Order_type),
		First:       Nullint32ToUint32(in.Pagination.First),
		After:       NullToString(in.Pagination.After),
		OrderStatus: orderStatusToRPC(in.Order_status),
	})

	if err != nil {
		return nil, err
	}

	orders := make([]ServiceOrder, 0, len(response.GetOrderServices()))

	for _, order := range response.GetOrderServices() {

		if c := orderRPCToOrder(order); c != nil {

			c.User_profile = &Profile{}
			c.Company_profile = &CompanyProfile{}

			orders = append(orders, *c)
		}
	}

	userIDs := make([]string, 0, len(response.GetOrderServices()))
	companiesIDs := make([]string, 0, len(response.GetOrderServices()))

	for _, c := range response.GetOrderServices() {
		if c.GetIsCompany() {
			companiesIDs = append(companiesIDs, c.GetProfileID())
		} else {
			userIDs = append(userIDs, c.GetProfileID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range response.GetOrderServices() {
		for i := range orders {

			if orders[i].ID == p.GetID() {

				// company profile
				if p.GetIsCompany() {
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == p.GetProfileID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						orders[i].Company_profile = &pr
						orders[i].User_profile = &Profile{
							Location: &LocationProfile{},
						}
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[p.GetProfileID()])
					orders[i].User_profile = &profile

				}

			}

		}
	}

	return &ServiceOrdersResolver{
		R: &ServiceOrders{
			Order_amount: response.GetOrdersAmount(),
			Orders:       orders,
		},
	}, nil
}

func (_ *Resolver) GetReceivedProposals(ctx context.Context, in GetReceivedProposalsRequest) (*ProposalsResolver, error) {
	response, err := services.GetReceivedProposals(ctx, &servicesRPC.GetReceivedProposalsRequest{
		CompanyID: NullToString(in.Company_id),
		RequestID: NullToString(in.Request_id),
		First:     Nullint32ToUint32(in.Pagination.First),
		After:     NullToString(in.Pagination.After),
	})

	if err != nil {
		return nil, err
	}

	proposals := make([]ProposalDetail, 0, len(response.GetProposals()))

	for _, proposal := range response.GetProposals() {

		if c := proposalRPCToproposal(proposal); c != nil {

			c.User_profile = &Profile{}
			proposals = append(proposals, *c)
		}
	}

	userIDs := make([]string, 0, len(response.GetProposals()))

	for _, c := range response.GetProposals() {
		userIDs = append(userIDs, c.GetProfileID())

	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range response.GetProposals() {
		for i := range proposals {

			if proposals[i].ID == p.GetID() {
				// user profile
				profile := ToProfile(ctx, userResp.GetProfiles()[p.GetProfileID()])
				proposals[i].User_profile = &profile

			}

		}
	}

	return &ProposalsResolver{
		R: &Proposals{
			Proposal_amount: response.GetProposalsAmount(),
			Proposals:       proposals,
		},
	}, nil
}

func (_ *Resolver) GetSendedProposals(ctx context.Context, in GetSendedProposalsRequest) (*ProposalsResolver, error) {
	response, err := services.GetSendedProposals(ctx, &servicesRPC.GetReceivedProposalsRequest{
		First:     Nullint32ToUint32(in.Pagination.First),
		After:     NullToString(in.Pagination.After),
		CompanyID: NullToString(in.Company_id),
	})

	if err != nil {
		return nil, err
	}

	proposals := make([]ProposalDetail, 0, len(response.GetProposals()))

	for _, proposal := range response.GetProposals() {

		if c := proposalRPCToproposal(proposal); c != nil {

			c.User_profile = &Profile{}
			proposals = append(proposals, *c)
		}
	}

	userIDs := make([]string, 0, len(response.GetProposals()))

	for _, c := range response.GetProposals() {
		userIDs = append(userIDs, c.GetProfileID())

	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range response.GetProposals() {
		for i := range proposals {

			if proposals[i].ID == p.GetID() {
				// user profile
				profile := ToProfile(ctx, userResp.GetProfiles()[p.GetProfileID()])
				proposals[i].User_profile = &profile
			}
		}
	}

	return &ProposalsResolver{
		R: &Proposals{
			Proposal_amount: response.GetProposalsAmount(),
			Proposals:       proposals,
		},
	}, nil
}

func (_ *Resolver) AddNoteForOrderService(ctx context.Context, in AddNoteForOrderServiceRequest) (*SuccessResolver, error) {
	_, err := services.AddNoteForOrderService(ctx, &servicesRPC.AddNoteForOrderServiceRequest{
		CompanyID: NullToString(in.Company_id),
		OrderID:   in.Order_id,
		Text:      in.Text,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) SaveVOfficeService(ctx context.Context, in SaveVOfficeServiceRequest) (*SuccessResolver, error) {
	_, err := services.SaveVOfficeService(ctx, &servicesRPC.VOfficeServiceActionRequest{
		CompanyID: NullToString(in.Company_id),
		ServiceID: in.Service_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) UnSaveVOfficeService(ctx context.Context, in UnSaveVOfficeServiceRequest) (*SuccessResolver, error) {
	_, err := services.UnSaveVOfficeService(ctx, &servicesRPC.VOfficeServiceActionRequest{
		CompanyID: NullToString(in.Company_id),
		ServiceID: in.Service_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) SaveServiceRequest(ctx context.Context, in SaveServiceRequestRequest) (*SuccessResolver, error) {
	_, err := services.SaveServiceRequest(ctx, &servicesRPC.VOfficeServiceActionRequest{
		CompanyID: NullToString(in.Company_id),
		ServiceID: in.Service_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) UnSaveServiceRequest(ctx context.Context, in UnSaveServiceRequestRequest) (*SuccessResolver, error) {
	_, err := services.UnSaveServiceRequest(ctx, &servicesRPC.VOfficeServiceActionRequest{
		CompanyID: NullToString(in.Company_id),
		ServiceID: in.Service_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func (_ *Resolver) GetSavedVOfficeServices(ctx context.Context, in GetSavedVOfficeServicesRequest) (*ServicesResolver, error) {

	savedServices, err := services.GetSavedVOfficeServices(ctx, &servicesRPC.GetSavedVOfficeServicesRequest{
		First:     Nullint32ToUint32(in.Pagination.First),
		After:     NullToString(in.Pagination.After),
		CompanyID: NullToString(in.Company_id),
	})

	if err != nil {
		return nil, err
	}

	return &ServicesResolver{
		R: servicesRPCServicesToServicesResolver(savedServices),
	}, nil
}

func (_ *Resolver) GetSavedServicesRequest(ctx context.Context, in GetSavedServicesRequestRequest) (*ServicesRequestResolver, error) {

	savedServices, err := services.GetSavedServicesRequest(ctx, &servicesRPC.GetSavedVOfficeServicesRequest{
		First:     Nullint32ToUint32(in.Pagination.First),
		After:     NullToString(in.Pagination.After),
		CompanyID: NullToString(in.Company_id),
	})

	if err != nil {
		return nil, err
	}

	return &ServicesRequestResolver{
		R: &ServicesRequest{
			Service_amount: savedServices.GetServiceAmount(),
			Services:       servicesRequestRPCToServicesRequest(savedServices.GetRequest()),
		},
	}, nil
}

func (_ *Resolver) WriteReviewForService(ctx context.Context, in WriteReviewForServiceRequest) (*SuccessResolver, error) {
	res, err := services.WriteReviewForService(ctx, &servicesRPC.WriteReviewRequest{
		OfficeID:     in.Office_id,
		ServiceID:    in.Service_id,
		OwnerID:      in.Owner_id,
		ReviewDetail: reviewDetailToRPC(in.Review_detail),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) WriteReviewForServiceRequest(ctx context.Context, in WriteReviewForServiceRequestRequest) (*SuccessResolver, error) {
	res, err := services.WriteReviewForServiceRequest(ctx, &servicesRPC.WriteReviewRequest{
		OwnerID:        in.Owner_id,
		IsOwnerCompany: in.Is_owner_company,
		ReviewDetail:   reviewDetailToRPC(in.Review_detail),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      res.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) GetServicesReview(ctx context.Context, in GetServicesReviewRequest) (*ServicesReviewResolver, error) {
	res, err := services.GetServicesReview(ctx, &servicesRPC.GetReviewRequest{
		First:     Nullint32ToUint32(in.Pagination.First),
		After:     NullToString(in.Pagination.After),
		CompanyID: NullToString(in.Company_id),
		OfficeID:  in.Office_id,
	})
	if err != nil {
		return nil, err
	}

	reviews := make([]Review, 0, len(res.GetReviews()))

	for _, review := range res.GetReviews() {

		if c := reviewRPCToReview(review); c != nil {

			c.User_profile = &Profile{}
			c.Company_profile = &CompanyProfile{}

			reviews = append(reviews, *c)
		}
	}

	userIDs := make([]string, 0, len(res.GetReviews()))
	companiesIDs := make([]string, 0, len(res.GetReviews()))

	for _, c := range res.GetReviews() {
		if c.GetIsCompany() {
			companiesIDs = append(companiesIDs, c.GetProfileID())
		} else {
			userIDs = append(userIDs, c.GetProfileID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range res.GetReviews() {
		for i := range reviews {

			if reviews[i].ID == p.GetID() {

				// company profile
				if p.GetIsCompany() {
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == p.GetProfileID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						reviews[i].Company_profile = &pr
						reviews[i].User_profile = &Profile{
							Location: &LocationProfile{},
						}
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[p.GetProfileID()])
					reviews[i].User_profile = &profile

				}

			}

		}
	}

	return &ServicesReviewResolver{
		R: &ServicesReview{
			Review_amount: res.GetReviewAmount(),
			Reviews_avg:   reviewAVGRPCToReviewAVG(res.GetServiceReviewAVG()),
			Reviews:       reviews,
		},
	}, nil
}

func (_ *Resolver) GetServicesRequestReview(ctx context.Context, in GetServicesRequestReviewRequest) (*ServicesReviewResolver, error) {
	res, err := services.GetServicesRequestReview(ctx, &servicesRPC.GetReviewRequest{
		First:     Nullint32ToUint32(in.Pagination.First),
		After:     NullToString(in.Pagination.After),
		CompanyID: NullToString(in.Company_id),
		OwnerID:   NullToString(in.Owner_id),
	})
	if err != nil {
		return nil, err
	}

	reviews := make([]Review, 0, len(res.GetReviews()))

	for _, review := range res.GetReviews() {

		if c := reviewRPCToReview(review); c != nil {

			c.User_profile = &Profile{}
			c.Company_profile = &CompanyProfile{}

			reviews = append(reviews, *c)
		}
	}

	userIDs := make([]string, 0, len(res.GetReviews()))
	companiesIDs := make([]string, 0, len(res.GetReviews()))

	for _, c := range res.GetReviews() {
		if c.GetIsCompany() {
			companiesIDs = append(companiesIDs, c.GetProfileID())
		} else {
			userIDs = append(userIDs, c.GetProfileID())
		}
	}

	companyResp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: companiesIDs,
	})
	if err != nil {
		return nil, err
	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range res.GetReviews() {
		for i := range reviews {

			if reviews[i].ID == p.GetID() {

				// company profile
				if p.GetIsCompany() {
					var profile *companyRPC.Profile

					for _, company := range companyResp.GetProfiles() {
						if company.GetId() == p.GetProfileID() {
							profile = company
							break
						}
					}

					if profile != nil {
						pr := toCompanyProfile(ctx, *profile)
						reviews[i].Company_profile = &pr
						reviews[i].User_profile = &Profile{
							Location: &LocationProfile{},
						}
					}
				} else {
					// user profile
					profile := ToProfile(ctx, userResp.GetProfiles()[p.GetProfileID()])
					reviews[i].User_profile = &profile

				}

			}

		}
	}

	return &ServicesReviewResolver{
		R: &ServicesReview{
			Review_amount: res.GetReviewAmount(),
			Reviews_avg:   reviewAVGRPCToReviewAVG(res.GetServiceReviewAVG()),
			Reviews:       reviews,
		},
	}, nil
}
