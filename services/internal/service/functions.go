package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	notmes "gitlab.lan/Rightnao-site/microservices/services/internal/notification_messages"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/location"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/qualifications"
	review "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/review"
	"go.mongodb.org/mongo-driver/bson/primitive"

	companyadmin "gitlab.lan/Rightnao-site/microservices/services/internal/company-admin"
	offer "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/offers"
	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	serviceorder "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/service-order"

	office "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/v-office"
)

// CheckIfURLForVOfficeIsTaken checks if url is already taken or not
func (s Service) CheckIfURLForVOfficeIsTaken(ctx context.Context, url string) (bool, error) {
	span := s.tracer.MakeSpan(ctx, "CreateNewAccount")
	defer span.Finish()

	// checks if url is busy
	isURLBusy, err := s.repository.Services.IsURLBusy(ctx, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return false, err
	}
	if isURLBusy {
		s.tracer.LogError(span, errors.New("url_already_taken"))
		return false, err
	}

	return false, nil
}

// ChangeVofficeCover ...
func (s Service) ChangeVofficeCover(ctx context.Context, officeID string, companyID string, file *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeOriginCover")
	defer span.Finish()

	// check admin level
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_allowed")
		}
	}

	err := s.repository.Services.ChangeVofficeCover(ctx, officeID, companyID, file.URL)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return "", nil
}

/// ChangeVofficeOriginCover ...
func (s Service) ChangeVofficeOriginCover(ctx context.Context, officeID string, companyID string, file *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeOriginCover")
	defer span.Finish()

	id := file.GenerateID()

	// check admin level
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_allowed")
		}
	}

	err := s.repository.Services.ChangeVofficeOriginCover(ctx, officeID, companyID, file.URL)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// AddFileInVofficeService ...
func (s Service) AddFileInVofficeService(ctx context.Context, officeID, serviceID, companyID string, file *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInVofficeService")
	defer span.Finish()

	id := file.GenerateID()

	// check admin level
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_allowed")
		}
	}

	err := s.repository.Services.AddFileInVofficeService(ctx, officeID, serviceID, companyID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveVofficeCover ...
func (s Service) RemoveVofficeCover(ctx context.Context, officeID, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveVofficeCover")
	defer span.Finish()

	// check admin level
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_allowed")
		}
	}

	err := s.repository.Services.RemoveVofficeCover(ctx, officeID, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// CreateVOffice creates office
func (s Service) CreateVOffice(ctx context.Context, companyID string, vOffice *office.Office) (string, error) {
	span := s.tracer.MakeSpan(ctx, "CreateVOffice")
	defer span.Finish()

	// retrieve id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	if companyID != "" {
		// check admin level of company
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
	} else {
		vOffice.SetUserID(userID)
	}

	// checks if url is busy
	// isURLBusy, err := s.CheckIfURLForVOfficeIsTaken(ctx, vOffice.URL)
	// if err != nil {
	// 	return "", err
	// }
	// if isURLBusy {
	// 	s.tracer.LogError(span, errors.New("url_already_taken"))
	// 	return "url_already_taken", err
	// }

	vOfficeID := vOffice.GenerateID()
	vOffice.CreatedAt = time.Now()
	vOffice.SetCompanyID(companyID)

	err = s.repository.Services.CreateVOffice(
		ctx,
		vOffice,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return vOfficeID, nil
}

// ChangeVOffice ...
func (s Service) ChangeVOffice(ctx context.Context, companyID, officeID string, vOffice *office.Office) error {
	span := s.tracer.MakeSpan(ctx, "ChangeVOffice")
	defer span.Finish()

	// retrieve id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if companyID != "" {
		// check admin level of company
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enough_authenticitation")
		}
	} else {
		vOffice.SetUserID(userID)
	}

	vOffice.SetID(officeID)
	vOffice.SetCompanyID(companyID)

	err = s.repository.Services.ChangeVOffice(
		ctx,
		vOffice,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveVOffice ...
func (s Service) RemoveVOffice(ctx context.Context, companyID, officeID string) error {
	span := s.tracer.MakeSpan(ctx, "CreateVOffice")
	defer span.Finish()

	// retrieve id of user
	_, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if companyID != "" {
		// check admin level of company
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enough_authenticitation")
		}
	}

	err = s.repository.Services.RemoveVOffice(
		ctx,
		officeID,
	)

	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetVOffice gets the vOffice for user OR company
func (s Service) GetVOffice(ctx context.Context, companyID string, userID string) ([]*office.Office, error) {
	span := s.tracer.MakeSpan(ctx, "GetVOffice")
	defer span.Finish()

	if companyID == "" && userID == "" {
		return nil, errors.New(`id_is_empty`)
	}

	vOffices, err := s.repository.Services.GetVOffice(ctx, companyID, userID)
	if err != nil {
		return nil, err
	}

	if len(vOffices) > 0 {
		for _, vOffice := range vOffices {
			// Check if user is admin
			currentUserID, err := s.authRPC.GetUserID(ctx)
			if err == nil && userID != "" {
				vOffice.IsMe = currentUserID == userID
			}

			/// Check if is comapny admin
			if companyID != "" {
				allowed := s.checkAdminLevel(
					ctx,
					companyID,
					companyadmin.AdminLevelAdmin,
					companyadmin.AdminLevelVShop,
				)
				vOffice.IsMe = allowed
			}

			/// Location
			cityID := vOffice.Location.City.ID
			lang := "en"
			cityIDStr, _ := strconv.Atoi(cityID)
			cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityIDStr), &lang)
			if err != nil {
				s.tracer.LogError(span, err)
				return nil, err

			}
			vOffice.Location = location.Location{
				City: &location.City{
					ID:          cityID,
					Name:        cityName,
					Subdivision: subdivision,
				},
				Country: &location.Country{
					ID: countryID,
				},
			}

			if vOffice.ReturnDate != nil {
				if time.Now().After(*vOffice.ReturnDate) {
					err := s.repository.Services.IsOutOfOffice(ctx, vOffice.GetID(), false, nil)
					if err != nil {
						return nil, err
					}

					vOffice.IsOut = false
					vOffice.ReturnDate = nil
				}
			}
		}
	}

	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return vOffices, err
}

// GetVOfficeByID ...
func (s Service) GetVOfficeByID(ctx context.Context, companyID, officeID string) (*office.Office, error) {
	span := s.tracer.MakeSpan(ctx, "GetVOfficeByID")
	defer span.Finish()

	// retrieve id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	if companyID == "" && userID == "" {
		return nil, errors.New(`id_is_empty`)
	}

	vOffice, err := s.repository.Services.GetVOfficeByID(ctx, officeID)

	if err != nil {
		return nil, err
	}

	/// Check if user is admin
	if err == nil && vOffice.UserID != nil {
		vOffice.IsMe = userID == vOffice.UserID.Hex()
	}

	/// Check if is comapny admin
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		vOffice.IsMe = allowed
	}

	/// Location
	cityID := vOffice.Location.City.ID
	lang := "en"
	cityIDStr, _ := strconv.Atoi(cityID)
	cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityIDStr), &lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err

	}
	vOffice.Location = location.Location{
		City: &location.City{
			ID:          cityID,
			Name:        cityName,
			Subdivision: subdivision,
		},
		Country: &location.Country{
			ID: countryID,
		},
	}

	if vOffice.ReturnDate != nil {
		if time.Now().After(*vOffice.ReturnDate) {
			err := s.repository.Services.IsOutOfOffice(ctx, vOffice.GetID(), false, nil)
			if err != nil {
				return nil, err
			}

			vOffice.IsOut = false
			vOffice.ReturnDate = nil
		}
	}

	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return vOffice, err
}

// GetServicesRequest ...
func (s Service) GetServicesRequest(ctx context.Context, ownerID, companyID string) ([]*servicerequest.Request, error) {
	span := s.tracer.MakeSpan(ctx, "GetServicesRequest")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	if ownerID != "" {
		profileID = ownerID
	}

	services, err := s.repository.Services.GetServicesRequest(ctx, profileID, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for _, vos := range services {

		proposalAmount, err := s.repository.Services.GetProposalAmount(ctx, vos.GetID())

		hasLiked, _ := s.repository.Services.HasLikedService(ctx, profileID, "", vos.GetID())
		vos.HasLiked = hasLiked

		if proposalAmount > 0 && err == nil {
			vos.ProposalAmount = proposalAmount
		}

		cityID, _ := strconv.Atoi(vos.Location.City.ID)
		lang := "en"
		/// Location
		cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err

		}
		vos.Location = &location.Location{
			City: &location.City{
				ID:          strconv.Itoa(cityID),
				Name:        cityName,
				Subdivision: subdivision,
			},
			Country: &location.Country{
				ID: countryID,
			},
		}
	}

	return services, nil
}

// GetServiceRequest ...
func (s Service) GetServiceRequest(ctx context.Context, companyID string, serviceID string) (*servicerequest.Request, error) {
	span := s.tracer.MakeSpan(ctx, "GetServiceRequest")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	service, err := s.repository.Services.GetServiceRequest(ctx, serviceID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	proposalAmount, err := s.repository.Services.GetProposalAmount(ctx, service.GetID())

	if proposalAmount > 0 && err == nil {
		service.ProposalAmount = proposalAmount
	}

	hasLiked, _ := s.repository.Services.HasLikedService(ctx, profileID, "", service.GetID())
	service.HasLiked = hasLiked

	cityID, _ := strconv.Atoi(service.Location.City.ID)
	lang := "en"
	/// Location
	cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err

	}
	service.Location = &location.Location{
		City: &location.City{
			ID:          strconv.Itoa(cityID),
			Name:        cityName,
			Subdivision: subdivision,
		},
		Country: &location.Country{
			ID: countryID,
		},
	}

	return service, nil
}

// ChangeServicesRequestStatus ...
func (s Service) ChangeServicesRequestStatus(ctx context.Context, companyID, serviceID string, serviceStatus string) error {

	span := s.tracer.MakeSpan(ctx, "ChangeServicesRequestStatus")
	defer span.Finish()

	// retrieve id of user
	_, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enough_authenticitation")
		}
	}

	data := servicerequest.ServiceReqestStatus{}

	switch serviceStatus {
	case "SERVICE_ACTIVE":
		data.IsDraft = false
		data.IsPaused = false
	case "SERVICE_DRAFT":
		data.IsDraft = true
	case "SERVICE_PAUSED":
		data.IsPaused = true
	case "SERVICE_DEACTIVATE":
		data.IsDraft = true
		data.IsPaused = true
	case "SERVICE_CLOSED":
		data.IsClosed = true

	}

	err = s.repository.Services.ChangeServicesRequestStatus(ctx, serviceID, data)

	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// IsOutOfOffice tells the office to be out of service for a period of time(returnTime)
func (s Service) IsOutOfOffice(ctx context.Context, companyID, officeID string, isOut bool, returnTime *time.Time) error {
	span := s.tracer.MakeSpan(ctx, "IsOutOfOffice")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.IsOutOfOffice(ctx, officeID, isOut, returnTime)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil

}

// ChangeVOfficeName changes office name
func (s Service) ChangeVOfficeName(ctx context.Context, companyID, officeID, name string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeVOfficeName")
	defer span.Finish()

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	rErr := s.repository.Services.ChangeVOfficeName(ctx, officeID, name)

	if rErr != nil {
		s.tracer.LogError(span, rErr)
		return rErr
	}

	return nil
}

// GetVOfficeServices is the logic why which it returns services according to officeID
func (s Service) GetVOfficeServices(ctx context.Context, companyID, officeID string) ([]*servicerequest.Service, error) {
	span := s.tracer.MakeSpan(ctx, "GetVOfficeServices")
	defer span.Finish()

	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enough_authenticitation")
		}

		profileID = companyID
	}
	vOfficeService, err := s.repository.Services.GetVOfficeServices(ctx, officeID)

	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for _, vos := range vOfficeService {

		// Liked
		hasLiked, _ := s.repository.Services.HasLikedService(ctx, profileID, vos.GetID(), "")
		vos.HasLiked = hasLiked

		cityID, _ := strconv.Atoi(vos.Location.City.ID)
		lang := "en"
		/// Location
		cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err

		}
		vos.Location = &location.Location{
			City: &location.City{
				ID:          strconv.Itoa(cityID),
				Name:        cityName,
				Subdivision: subdivision,
			},
			Country: &location.Country{
				ID: countryID,
			},
		}
	}

	return vOfficeService, nil
}

// GetAllServices is the logic why which it returns services according to officeID
func (s Service) GetAllServices(ctx context.Context, companyID string) ([]*servicerequest.Service, error) {
	span := s.tracer.MakeSpan(ctx, "GetAllServices")
	defer span.Finish()

	isCompany := false
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enough_authenticitation")
		}
		isCompany = true
		profileID = companyID
	}
	vOfficeService, err := s.repository.Services.GetAllServices(ctx, profileID, isCompany)

	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for _, vos := range vOfficeService {

		cityID, _ := strconv.Atoi(vos.Location.City.ID)
		lang := "en"
		/// Location
		cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err

		}
		vos.Location = &location.Location{
			City: &location.City{
				ID:          strconv.Itoa(cityID),
				Name:        cityName,
				Subdivision: subdivision,
			},
			Country: &location.Country{
				ID: countryID,
			},
		}
	}

	return vOfficeService, nil
}

// GetVOfficeService ...
func (s Service) GetVOfficeService(ctx context.Context, companyID, officeID, serviceID string) (*servicerequest.Service, error) {
	span := s.tracer.MakeSpan(ctx, "GetVOfficeService")
	defer span.Finish()

	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enough_authenticitation")
		}
		profileID = companyID
	}
	vOfficeService, err := s.repository.Services.GetVOfficeService(ctx, officeID, serviceID)

	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// Liked
	hasLiked, _ := s.repository.Services.HasLikedService(ctx, profileID, vOfficeService.GetID(), "")
	vOfficeService.HasLiked = hasLiked

	cityID, _ := strconv.Atoi(vOfficeService.Location.City.ID)
	lang := "en"
	/// Location
	cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err

	}
	vOfficeService.Location = &location.Location{
		City: &location.City{
			ID:          strconv.Itoa(cityID),
			Name:        cityName,
			Subdivision: subdivision,
		},
		Country: &location.Country{
			ID: countryID,
		},
	}

	return vOfficeService, nil
}

// ChangeVOfficeServiceStatus ...
func (s Service) ChangeVOfficeServiceStatus(ctx context.Context, companyID, officeID, serviceID string, serviceStatus string) error {

	span := s.tracer.MakeSpan(ctx, "ChangeVOfficeServiceStatus")
	defer span.Finish()

	// retrieve id of user
	_, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enough_authenticitation")
		}
	}

	data := servicerequest.ServiceStatus{}

	switch serviceStatus {
	case "SERVICE_ACTIVE":
		data.IsDraft = false
		data.IsPaused = false
	case "SERVICE_DRAFT":
		data.IsDraft = true
	case "SERVICE_PAUSED":
		data.IsPaused = true
	case "SERVICE_DEACTIVATE":
		data.IsDraft = true
		data.IsPaused = true
	}

	err = s.repository.Services.ChangeVOfficeServiceStatus(ctx, officeID, serviceID, data)

	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddVOfficeService adds or changes the service which was definetly posted by voffice, which belongs to user OR company
func (s Service) AddVOfficeService(ctx context.Context, companyID string, officeID string, service *servicerequest.Service) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddVOfficeService")
	defer span.Finish()

	// retrieve id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	if service != nil {
		// check admin level of company
		if companyID != "" {
			allowed := s.checkAdminLevel(
				ctx,
				companyID,
				companyadmin.AdminLevelAdmin,
				companyadmin.AdminLevelVShop,
			)
			if !allowed {
				return "", errors.New("not_enough_authenticitation")
			}

			service.SetCompanyID(companyID)
		} else {
			service.SetUserID(userID)
		}

		servID := service.GenerateID()
		service.CreatedAt = time.Now()
		service.SetID(servID)
		service.SetOfficeID(officeID)

		err := s.repository.Services.AddVOfficeService(ctx, service)
		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}

		return servID, nil
	}

	return "", nil
}

// ChangeVOfficeService adds or changes the service which was definetly posted by voffice, which belongs to user OR company
func (s Service) ChangeVOfficeService(ctx context.Context, companyID, serviceID, officeID string, service *servicerequest.Service) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeVOfficeService")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return "", err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
	}

	err := s.repository.Services.ChangeVOfficeService(ctx, serviceID, officeID, service)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return serviceID, nil
}

// RemoveVOfficeService removes service
func (s Service) RemoveVOfficeService(ctx context.Context, companyID, serviceID string) (string, error) {
	span := s.tracer.MakeSpan(ctx, "RemoveVOfficeService")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.RemoveVOfficeService(ctx, serviceID)
	if err != nil {
		return "", err
	}

	return "", nil
}

// RemoveFilesInVOfficeService removes files from service
func (s Service) RemoveFilesInVOfficeService(ctx context.Context, companyID, serviceID string, fileIDs []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInVOfficeService")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.RemoveFilesInVOfficeService(ctx, serviceID, fileIDs)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFilesInServiceRequest ...
func (s Service) RemoveFilesInServiceRequest(ctx context.Context, companyID, serviceID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInServiceRequest")
	defer span.Finish()

	// retrieve id of user
	_, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err = s.repository.Services.RemoveFilesInServiceRequest(ctx, serviceID, ids)
	if err != nil {
		return err
	}

	return nil
}

// AddServicesRequest adds the request which was posted by user od ID
func (s Service) AddServicesRequest(ctx context.Context, companyID string, request *servicerequest.Request) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddServicesRequest")
	defer span.Finish()

	// retrieve id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enought_authenticitation")
		}
	}

	request.CreatedAt = time.Now()
	requestID := request.GenerateID()
	request.SetCompanyID(companyID)
	request.SetUserID(userID)

	err = s.repository.Services.AddServicesRequest(ctx, companyID, request)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return requestID, nil
}

// ChangeServicesRequest ...
func (s Service) ChangeServicesRequest(ctx context.Context, companyID, serviceID string, request *servicerequest.Request) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeServicesRequest")
	defer span.Finish()

	// retrieve id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enought_authenticitation")
		}
	}

	request.SetID(serviceID)
	request.SetCompanyID(companyID)
	request.SetUserID(userID)

	request.CreatedAt = time.Now()

	err = s.repository.Services.ChangeServicesRequest(ctx, request)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return serviceID, nil
}

// RemoveServicesRequest removes service
func (s Service) RemoveServicesRequest(ctx context.Context, companyID, requestID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveServicesRequest")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.RemoveServicesRequest(ctx, requestID)
	if err != nil {
		return err
	}

	return nil
}

// AddChangeVOfficeDescription adds or changes the description of vOffice
func (s Service) AddChangeVOfficeDescription(ctx context.Context, companyID, officeID, description string) error {
	span := s.tracer.MakeSpan(ctx, "AddChangeVOfficeDescription")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.AddChangeVOfficeDescription(ctx, officeID, description)
	if err != nil {
		return err
	}

	return nil
}

// AddVOfficePortfolio adds  portfolio in vOffice
func (s Service) AddVOfficePortfolio(ctx context.Context, companyID, officeID string, portfolio *office.Portfolio) (string, []*file.Link, error) {
	span := s.tracer.MakeSpan(ctx, "AddVOfficePortfolio")
	defer span.Finish()

	if portfolio == nil {
		return "", nil, errors.New("huy")
	}

	id := portfolio.GenerateID()

	// retrieve id of user
	userID, authErr := s.authRPC.GetUserID(ctx)
	if authErr != nil {
		s.tracer.LogError(span, authErr)
		return "", nil, authErr
	}

	// // check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", nil, errors.New("not_enought_authenticitation")
		}
	}

	portfolio.CreatedAt = time.Now()

	for i := range portfolio.Link {
		portfolio.Link[i].ID = primitive.NewObjectID()
	}

	err := s.repository.Services.AddVOfficePortfolio(ctx, officeID, portfolio, userID, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", nil, err
	}

	return id, portfolio.Link, nil
}

// ChangeVOfficePortfolio Changes  portfolio in vOffice
func (s Service) ChangeVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string, portfolio *office.Portfolio) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeVOfficePortfolio")
	defer span.Finish()

	id := portfolio.GenerateID()

	// retrieve id of user
	// userID, err := s.authRPC. UserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return "", err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enough_authenticitation")
		}
	}

	portfolio.CreatedAt = time.Now()

	err := s.repository.Services.ChangeVOfficePortfolio(ctx, officeID, portfolioID, portfolio)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveVOfficePortfolio removes portfolio from the vOffice
func (s Service) RemoveVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveVOfficePortfolio")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.RemoveVOfficePortfolio(ctx, officeID, portfolioID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddFileInVOfficePortfolio adds file in vOffice portfolio
func (s Service) AddFileInVOfficePortfolio(ctx context.Context, officeID, portfolioID, companyID string, file *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInVOfficePortfolio")
	defer span.Finish()

	id := file.GenerateID()

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.AddFileInVOfficePortfolio(ctx, officeID, portfolioID, companyID, file)
	if err != nil {
		return "", err
	}

	return id, nil
}

// AddFileInServiceRequest ...
func (s Service) AddFileInServiceRequest(ctx context.Context, serviceID, companyID string, file *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInServiceRequest")
	defer span.Finish()

	id := file.GenerateID()

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.AddFileInServiceRequest(ctx, serviceID, companyID, file)
	if err != nil {
		return "", err
	}

	return id, nil
}

// AddFileInOrderService ...
func (s Service) AddFileInOrderService(ctx context.Context, orderID, companyID string, file *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInOrderService")
	defer span.Finish()

	id := file.GenerateID()

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.AddFileInOrderService(ctx, orderID, file)
	if err != nil {
		return "", err
	}

	return id, nil
}

// RemoveFilesInVOfficePortfolio removes files from vOffice portfolio
func (s Service) RemoveFilesInVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string, fileIDs []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInVOfficePortfolio")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.RemoveFilesInVOfficePortfolio(ctx, officeID, portfolioID, fileIDs)
	if err != nil {
		return err
	}

	return nil

}

// return false if level doesn't much
func (s Service) checkAdminLevel(ctx context.Context, companyID string, requiredLevels ...companyadmin.AdminLevel) bool {
	span := s.tracer.MakeSpan(ctx, "checkAdminLevel")
	defer span.Finish()

	actualLevel, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		log.Println("Error: checkAdminLevel:", err)
		return false
	}

	for _, lvl := range requiredLevels {
		if lvl == actualLevel {
			return true
		}
	}

	return false
}

// AddVOfficeLanguages adds or changes office qualifications
func (s Service) AddVOfficeLanguages(ctx context.Context, companyID, officeID string, langs []*qualifications.Language) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "AddVOfficeLanguages")
	defer span.Finish()

	// retrieve id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}
	}

	ids := make([]string, 0, len(langs))

	for i := range langs {
		id := langs[i].GetID()
		ids = append(ids, id)
	}

	err = s.repository.Services.AddVOfficeLanguages(ctx, userID, companyID, officeID, langs)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// ChangeVOfficeLanguage adds or changes office qualifications
func (s Service) ChangeVOfficeLanguage(ctx context.Context, companyID, officeID string, langs []*qualifications.Language) error {
	span := s.tracer.MakeSpan(ctx, "ChangeVOfficeLanguage")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return "", err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	for i := range langs {
		if langs[i].GetID() == "" {
			_ = langs[i].GenerateID()
		}
	}

	err := s.repository.Services.ChangeVOfficeLanguage(ctx, officeID, langs)
	if err != nil {
		return err
	}

	return nil
}

// RemoveVOfficeLanguages removes Language in qualifications
func (s Service) RemoveVOfficeLanguages(ctx context.Context, companyID, officeID string, languageIds []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveVOfficeLanguages")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enough_authenticitation")
		}
	}

	err := s.repository.Services.RemoveVOfficeLanguages(ctx, officeID, languageIds)
	if err != nil {
		return err
	}

	return nil
}

// RemoveLinksInVOfficePortfolio removes links
func (s Service) RemoveLinksInVOfficePortfolio(ctx context.Context, companyID, officeID, portfolioID string, linkIDs []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInVOfficeService")
	defer span.Finish()

	// retrieve id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}
	}

	err := s.repository.Services.RemoveLinksInVOfficePortfolio(ctx, officeID, portfolioID, linkIDs)
	if err != nil {
		return err
	}

	return nil
}

// SendProposalForServiceRequest ...
func (s Service) SendProposalForServiceRequest(ctx context.Context, data *offer.Proposal) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SendProposalForServiceRequest")
	defer span.Finish()

	if data == nil {
		return "", errors.New("proposal data can`t be empty")
	}

	id := data.GenerateID()
	data.ProposalDetail.CreatedAt = time.Now()
	data.ProposalDetail.Status = "new"

	err := s.repository.Services.SendProposalForServiceRequest(ctx, data)
	if err != nil {
		return "", err
	}

	// Nottification
	profileID := data.ProposalDetail.GetProfileID()
	note := &notmes.NewProposal{
		RequestID: data.GetRequestID(),
	}

	// Sender is company
	if data.ProposalDetail.IsCompany {
		note.CompanyID = profileID
	} else {
		note.UserSenderID = profileID
	}
	note.GenerateID()
	err = s.mq.SendProposalForServiceRequest(data.GetOwnerID(), note)

	if err != nil {
		s.tracer.LogError(span, err)
	}

	return id, nil
}

// GetProposalByID ...
func (s Service) GetProposalByID(ctx context.Context, companyID, proposalID string) (*offer.Proposal, error) {
	span := s.tracer.MakeSpan(ctx, "GetProposalByID")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	proposal, err := s.repository.Services.GetProposalByID(ctx, profileID, proposalID)
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

// OrderProposalForServiceRequest ...
func (s Service) OrderProposalForServiceRequest(ctx context.Context, data *serviceorder.Order) error {
	span := s.tracer.MakeSpan(ctx, "OrderProposalForServiceRequest")
	defer span.Finish()

	if data == nil {
		return errors.New("order data can`t be empty")
	}

	err := s.repository.Services.OrderService(ctx, data)

	if err != nil {
		return err
	}

	oldID := data.GetID()

	detail := data.OrderDetail

	res := &serviceorder.Order{
		IsOwnerCompany: detail.IsCompany,
		RequestID:      data.RequestID,
		ServiceID:      data.ServiceID,
		OrderType:      "buyer",
		OrderDetail: serviceorder.OrderDetail{
			IsCompany:      data.IsOwnerCompany,
			Currency:       detail.Currency,
			DeliveryTime:   detail.DeliveryTime,
			Description:    detail.Description,
			MaxPriceAmount: detail.MaxPriceAmount,
			MinPriceAmount: detail.MinPriceAmount,
			PriceAmount:    detail.PriceAmount,
			PriceType:      detail.PriceType,
			Status:         "in_progress",
			CustomDate:     detail.CustomDate,
			CreatedAt:      time.Now(),
		},
	}

	res.GenerateID()
	res.SetReferalID(oldID)

	// Set owner id , switch from profile id to owner id
	if !detail.ProfileID.IsZero() {
		res.SetOwnerID(detail.GetProfileID())
	}

	// Set profile id , switch from owner id to profile id
	if !data.OwnerID.IsZero() {
		res.OrderDetail.SetProfileID(data.GetOwnerID())
	}

	err = s.repository.Services.AcceptOrderService(ctx, res, oldID)
	if err != nil {
		return err
	}

	return nil
}

// IgnoreProposalForServiceRequest ...
func (s Service) IgnoreProposalForServiceRequest(ctx context.Context, companyID, proposalID string) error {
	span := s.tracer.MakeSpan(ctx, "IgnoreProposalForServiceRequest")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.IgnoreProposalForServiceRequest(ctx, profileID, proposalID)
	if err != nil {
		return err
	}

	return nil
}

// OrderService ...
func (s Service) OrderService(ctx context.Context, data *serviceorder.Order) (string, error) {
	span := s.tracer.MakeSpan(ctx, "OrderService")
	defer span.Finish()

	if data == nil {
		return "", errors.New("order data can`t be empty")
	}

	id := data.GenerateID()
	data.OrderDetail.CreatedAt = time.Now()
	data.OrderDetail.Status = "new"
	data.OrderType = "seller"

	err := s.repository.Services.OrderService(ctx, data)

	if err != nil {
		return "", err
	}

	// Nottification
	profileID := data.OrderDetail.GetProfileID()
	note := &notmes.NewOrder{
		ServiceID: data.GetServiceID(),
	}

	// Sender is company
	if data.OrderDetail.IsCompany {
		note.CompanyID = profileID
	} else {
		note.UserSenderID = profileID
	}
	note.GenerateID()
	err = s.mq.OrderService(data.GetOwnerID(), note)

	if err != nil {
		s.tracer.LogError(span, err)
	}

	return id, nil
}

// AcceptOrderService ...
func (s Service) AcceptOrderService(ctx context.Context, companyID, serviceID, orderID string) error {
	span := s.tracer.MakeSpan(ctx, "AcceptOrderService")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	res, err := s.repository.Services.GetVOfficerServiceOrderByID(ctx, orderID)

	if err != nil {
		return errors.New("Can`t find order")
	}

	oldID := res.GetID()

	newID := res.GenerateID()
	res.SetID(newID)
	res.SetReferalID(oldID)
	res.SetOwnerID(res.OrderDetail.GetProfileID())
	res.OrderType = "buyer"
	res.OrderDetail.SetProfileID(profileID)
	res.OrderDetail.CreatedAt = time.Now()
	res.OrderDetail.IsCompany = res.IsOwnerCompany
	res.OrderDetail.Status = "in_progress"
	res.SetServiceID(serviceID)

	err = s.repository.Services.AcceptOrderService(ctx, res, oldID)
	if err != nil {
		return err
	}

	return nil

}

// CancelServiceOrder ...
func (s Service) CancelServiceOrder(ctx context.Context, companyID, orderID string) error {
	span := s.tracer.MakeSpan(ctx, "CancelServiceOrder")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.CancelServiceOrder(ctx, profileID, orderID)
	if err != nil {
		return err
	}

	return nil

}

// DeclineServiceOrder ...
func (s Service) DeclineServiceOrder(ctx context.Context, companyID, orderID string) error {
	span := s.tracer.MakeSpan(ctx, "DeclineServiceOrder")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.DeclineServiceOrder(ctx, profileID, orderID)
	if err != nil {
		return err
	}

	return nil

}

// DeliverServiceOrder ...
func (s Service) DeliverServiceOrder(ctx context.Context, companyID, orderID string) error {
	span := s.tracer.MakeSpan(ctx, "DeliverServiceOrder")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.DeliverServiceOrder(ctx, profileID, orderID)
	if err != nil {
		return err
	}

	return nil

}

// AcceptDeliverdServiceOrder ...
func (s Service) AcceptDeliverdServiceOrder(ctx context.Context, companyID, orderID string) error {
	span := s.tracer.MakeSpan(ctx, "AcceptDeliverdServiceOrder")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.AcceptDeliverdServiceOrder(ctx, profileID, orderID)
	if err != nil {
		return err
	}

	return nil

}

// CancelDeliverdServiceOrder ...
func (s Service) CancelDeliverdServiceOrder(ctx context.Context, companyID, orderID string) error {
	span := s.tracer.MakeSpan(ctx, "CancelDeliverdServiceOrder")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.CancelDeliverdServiceOrder(ctx, profileID, orderID)
	if err != nil {
		return err
	}

	return nil

}

// GetVOfficerServiceOrders ...
func (s Service) GetVOfficerServiceOrders(ctx context.Context, ownerID, officeID string,
	orderType serviceorder.OrderType,
	orderstatus serviceorder.OrderStatus, first uint32, after string) (*serviceorder.GetOrder, error) {

	span := s.tracer.MakeSpan(ctx, "GetVOfficerServiceOrders")
	defer span.Finish()

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	res, err := s.repository.Services.GetVOfficerServiceOrders(ctx, ownerID, officeID, orderType, orderstatus, int(first), int(afterNumber))

	if err != nil {
		return nil, nil
	}

	for i, r := range res.Orders {
		// Get Services
		serv, err := s.repository.Services.GetVOfficeServiceByID(ctx, r.GetServiceID())

		if serv != nil && err == nil {
			res.Orders[i].Service = *serv
		}
		// Get Request
		req, reqErr := s.repository.Services.GetServiceRequest(ctx, r.GetRequestID())
		if req != nil && reqErr == nil {
			res.Orders[i].Request = *req
		}
	}

	return res, nil
}

// GetReceivedProposals ...
func (s Service) GetReceivedProposals(ctx context.Context, companyID, requestID string, first uint32, after string) (*offer.GetProposal, error) {
	span := s.tracer.MakeSpan(ctx, "GetReceivedProposals")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	res, err := s.repository.Services.GetReceivedProposals(ctx, profileID, requestID, int(first), int(afterNumber))

	if err != nil {
		return nil, err
	}

	for i, r := range res.Proposals {
		serv, err := s.repository.Services.GetVOfficeServiceByID(ctx, r.ProposalDetail.GetServiceID())

		if serv != nil && err == nil {
			res.Proposals[i].Service = *serv
		}

		req, reqErr := s.repository.Services.GetServiceRequest(ctx, r.GetRequestID())

		if req != nil && reqErr == nil {
			res.Proposals[i].Request = req
		} else {
			res.Proposals[i].Request = &servicerequest.Request{}
		}
	}

	return res, nil
}

// GetSendedProposals ...
func (s Service) GetSendedProposals(ctx context.Context, first uint32, after, companyID string) (*offer.GetProposal, error) {
	span := s.tracer.MakeSpan(ctx, "GetSendedProposals")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	res, err := s.repository.Services.GetSendedProposals(ctx, profileID, int(first), int(afterNumber))

	if err != nil {
		return nil, nil
	}

	for i, r := range res.Proposals {
		// Get Services
		serv, err := s.repository.Services.GetVOfficeServiceByID(ctx, r.ProposalDetail.GetServiceID())
		if serv != nil && err == nil {
			res.Proposals[i].Service = *serv
		}

		// Get Request
		req, reqErr := s.repository.Services.GetServiceRequest(ctx, r.GetRequestID())
		if req != nil && reqErr == nil {
			res.Proposals[i].Request = req
		} else {
			res.Proposals[i].Request = &servicerequest.Request{}
		}
	}

	return res, nil
}

// SaveVOfficeService ...
func (s Service) SaveVOfficeService(ctx context.Context, companyID, serviceID string) error {
	span := s.tracer.MakeSpan(ctx, "SaveVOfficeService")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.SaveVOfficeService(ctx, profileID, serviceID)
	if err != nil {
		return err
	}

	return nil
}

// UnSaveVOfficeService ...
func (s Service) UnSaveVOfficeService(ctx context.Context, companyID, serviceID string) error {
	span := s.tracer.MakeSpan(ctx, "UnSaveVOfficeService")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.UnSaveVOfficeService(ctx, profileID, serviceID)
	if err != nil {
		return err
	}

	return nil
}

// SaveServiceRequest ...
func (s Service) SaveServiceRequest(ctx context.Context, companyID, serviceID string) error {
	span := s.tracer.MakeSpan(ctx, "SaveServiceRequest")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.SaveServiceRequest(ctx, profileID, serviceID)
	if err != nil {
		return err
	}

	return nil
}

// UnSaveServiceRequest ...
func (s Service) UnSaveServiceRequest(ctx context.Context, companyID, serviceID string) error {
	span := s.tracer.MakeSpan(ctx, "UnSaveServiceRequest")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.UnSaveServiceRequest(ctx, profileID, serviceID)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedVOfficeServices ...
func (s Service) GetSavedVOfficeServices(ctx context.Context, companyID string, first uint32, after string) (*servicerequest.GetServices, error) {
	span := s.tracer.MakeSpan(ctx, "GetSavedVOfficeServices")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	res, amount, err := s.repository.Services.GetSavedVOfficeServices(ctx, profileID, int(first), int(afterNumber))

	if err != nil {
		return nil, nil
	}

	data := &servicerequest.GetServices{
		ServiceAmount: amount,
	}

	for _, r := range res {
		serv, _ := s.repository.Services.GetVOfficeServiceByID(ctx, r)

		if serv != nil && err == nil {
			data.Services = append(data.Services, serv)
		}
	}

	return data, nil
}

// GetSavedServicesRequest ...
func (s Service) GetSavedServicesRequest(ctx context.Context, companyID string, first uint32, after string) (*servicerequest.GetServicesRequest, error) {
	span := s.tracer.MakeSpan(ctx, "GetSavedServicesRequest")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	res, amount, err := s.repository.Services.GetSavedServicesRequest(ctx, profileID, int(first), int(afterNumber))

	if err != nil {
		return nil, nil
	}

	data := &servicerequest.GetServicesRequest{
		ServiceAmount: amount,
	}

	for _, r := range res {
		serv, serErr := s.GetServiceRequest(ctx, companyID, r)

		if serv != nil && serErr == nil {
			data.Services = append(data.Services, serv)
		}
	}

	return data, nil
}

// AddNoteForOrderService ...
func (s Service) AddNoteForOrderService(ctx context.Context, orderID, companyID, text string) error {
	span := s.tracer.MakeSpan(ctx, "AddNoteForOrderService")
	defer span.Finish()

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	err = s.repository.Services.AddNoteForOrderService(ctx, orderID, profileID, text)
	if err != nil {
		return err
	}

	return nil
}

// WriteReviewForService ...
func (s Service) WriteReviewForService(ctx context.Context, review review.Review) (string, error) {
	span := s.tracer.MakeSpan(ctx, "WriteReviewForService")
	defer span.Finish()

	id := review.GenerateID()
	review.ReviewDetail.CreatedAt = time.Now()

	err := s.repository.Services.WriteReviewForService(ctx, review)
	if err != nil {
		return "", err
	}

	return id, nil
}

// WriteReviewForServiceRequest ...
func (s Service) WriteReviewForServiceRequest(ctx context.Context, review review.Review) (string, error) {
	span := s.tracer.MakeSpan(ctx, "WriteReviewForServiceRequest")
	defer span.Finish()

	id := review.GenerateID()
	review.ReviewDetail.CreatedAt = time.Now()

	err := s.repository.Services.WriteReviewForServiceRequest(ctx, review)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetServicesReview ...
func (s Service) GetServicesReview(ctx context.Context, officeID string, first uint32, after string) (*review.GetReview, error) {
	span := s.tracer.MakeSpan(ctx, "GetServicesReview")
	defer span.Finish()

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	res, err := s.repository.Services.GetServicesReview(ctx, "", officeID, int(first), int(afterNumber))

	if err != nil {
		return nil, err
	}

	for i, r := range res.Reviews {
		// Get Service
		if r.GetServiceID() != "" {
			serv, err := s.repository.Services.GetVOfficeServiceByID(ctx, r.GetServiceID())

			if serv != nil && err == nil {
				res.Reviews[i].Service = *serv
			}
		}
	}

	return res, nil
}

// GetServicesRequestReview ...
func (s Service) GetServicesRequestReview(ctx context.Context, ownerID, companyID string, first uint32, after string) (*review.GetReview, error) {
	span := s.tracer.MakeSpan(ctx, "GetServicesRequestReview")
	defer span.Finish()

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	// retrieve id of user
	profileID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return nil, errors.New("not_enought_authenticitation")
		}

		profileID = companyID
	}

	if ownerID != "" {
		profileID = ownerID
	}

	res, err := s.repository.Services.GetServicesReview(ctx, profileID, "", int(first), int(afterNumber))

	if err != nil {
		return nil, err
	}

	for i, r := range res.Reviews {
		// Get Request
		if r.GetRequestID() != "" {
			req, err := s.repository.Services.GetServiceRequest(ctx, r.GetRequestID())

			if req != nil && err == nil {
				res.Reviews[i].Request = *req
			}
		}
	}

	return res, nil
}
