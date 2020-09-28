package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	careercenter "gitlab.lan/Rightnao-site/microservices/company/pkg/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/profile"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/status"
	notmes "gitlab.lan/Rightnao-site/microservices/company/pkg/notification_messages"
	arangorepo "gitlab.lan/Rightnao-site/microservices/company/pkg/repository/arango"
	"google.golang.org/grpc/metadata"
)

// CheckIfURLForCompanyIsTaken ...
func (s Service) CheckIfURLForCompanyIsTaken(ctx context.Context, url string) (bool, error) {
	span := s.tracer.MakeSpan(ctx, "CreateNewAccount")
	defer span.Finish()

	// checks if url is not busy
	isURLBusy, err := s.repository.Company.IsURLBusy(ctx, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return false, err
	}
	if isURLBusy {
		s.tracer.LogError(span, errors.New("url_already_taken"))
		return true, nil
	}

	return false, nil
}

// CreateNewAccount ...
func (s Service) CreateNewAccount(ctx context.Context, acc *account.Account) (string, string, error) {
	span := s.tracer.MakeSpan(ctx, "CreateNewAccount")
	defer span.Finish()

	// retrive id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", "", err
	}

	// checks if url is not busy
	isURLBusy, err := s.repository.Company.IsURLBusy(ctx, acc.URL)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", "", err
	}
	if isURLBusy {
		s.tracer.LogError(span, errors.New("url_already_taken"))
		return "", "", errors.New("url_already_taken")
	}

	// check if email is not busy
	isExists, err := s.repository.Company.IsEmailExists(ctx, acc.Emails[0].Email)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", "", err
	}
	if isExists {
		s.tracer.LogError(span, errors.New("email_already_in_use"))
		return "", "", errors.New("email_already_in_use")
	}

	// generating ID for company
	companyID := acc.GenerateID()
	acc.Status = status.CompanyStatusActivated
	acc.Size = account.SizeUnknown
	acc.CompanyType = account.Company
	acc.CreatedAt = time.Now()
	// generating id for phone
	acc.Phones[0].GenerateID()
	// generating id for email
	acc.Emails[0].GenerateID()
	acc.Emails[0].Primary = true
	// generating id for address
	acc.Addresses[0].GenerateID()
	acc.Addresses[0].IsPrimary = true
	acc.SetOwnerID(userID)

	// define country code info
	cd, countryID, err := s.infoRPC.GetCountryIDAndCountryCode(ctx, int32(acc.Phones[0].CountryCode.ID))
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return "", "", err
	}

	if countryID == "" {
		return "", "", errors.New("bad_country_code_id")
	}

	acc.Phones[0].CountryCode.Code = cd
	acc.Phones[0].CountryCode.CountryID = countryID

	//checks that fields are not empty and/or is over 120 characters
	//also checks that url and phone is valid
	// log.Printf("company register %#v\n", acc.Phones[0])
	err = registerValidator(acc)
	if err != nil {
		return "", "", err
	}

	// get info about location
	cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, acc.Addresses[0].Location.City.ID, nil)
	if errInfo != nil {
		s.tracer.LogError(span, errInfo)
		// internal_error
		return "", "", errInfo
	}

	acc.Addresses[0].Location.City.Name = cityName
	acc.Addresses[0].Location.City.Subdivision = subdivision
	acc.Addresses[0].Location.Country.ID = countryID

	// save in DB
	err = s.repository.Company.SaveNewCompanyAccount(ctx, acc)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", "", err
	}

	err = s.repository.arrangoRepo.SaveCompany(ctx, &arangorepo.Company{
		ID:        acc.GetID(),
		CreatedAt: time.Now(),
		Industry: arangorepo.Industry{
			Main: acc.Industry.Main,
			Sub:  acc.Industry.Sub,
		},
		Name:   acc.Name,
		Status: "ACTIVATED",
		Type:   string(acc.Type),
		URL:    acc.URL,
	})
	if err != nil {
		log.Println("arrangoRepo.SaveCompany:", err)
	}

	// set owner of company
	err = s.networkRPC.MakeCompanyOwner(ctx, userID, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", "", err // TODO: remove company
	}

	// set admin of company
	err = s.networkRPC.MakeCompanyAdmin(ctx, userID, companyID, account.AdminLevelAdmin)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", "", err // TODO: remove company
	}

	// generate tmp code for activation
	tmpCode, err := s.repository.Cache.CreateTemporaryCodeForEmailActivation(ctx, companyID, acc.Emails[0].Email)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// send email
	// err = s.mailRPC.SendEmail(
	// 	ctx,
	// 	acc.Emails[0].Email,
	// 	fmt.Sprint("<html><body><a target='_blank' href='https://"+s.host+"/api/activate/company?token=", tmpCode, "&company_id=", companyID, "'>Activate</a></body></html>")) // TODO: write template for message
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// }
	// fmt.Println(fmt.Sprint("https://"+s.host+"/api/activate/company?token=", tmpCode, "&company_id=", companyID)) // TODO: delete later

	emailMessage := fmt.Sprint("<html><body><a target='_blank' href='https://"+s.host+"/api/activate/company?token=", tmpCode, "&company_id=", companyID, "'>Activate</a></body></html>") // TODO: write template for message
	log.Println(emailMessage)

	err = s.mq.SendEmail(acc.Emails[0].Email, emailMessage)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	if acc.GetInvitedByID() != "" {
		s.AddGoldCoinsToWallet(ctx, acc.GetInvitedByID(), 1)
	}

	// TODO: send SMS

	return companyID, acc.URL, nil
}

// ActivateCompany checks temporary code and activates company and email
func (s Service) ActivateCompany(ctx context.Context, companyID string, code string) error {
	span := s.tracer.MakeSpan(ctx, "ActivateCompany")
	defer span.Finish()

	// check tmp code
	matched, email, err := s.repository.Cache.CheckTemporaryCodeForEmailActivation(ctx, companyID, code)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	if !matched {
		return errors.New("wrong_activation_code")
	}

	// change status of user
	err = s.repository.Company.ChangeStatusOfCompany(ctx, companyID, status.CompanyStatusActivated)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return err
	}

	// change status of email
	err = s.repository.Company.ActivateEmail(ctx, companyID, email)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return err
	}

	err = s.repository.Cache.Remove(ctx, email)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// DeactivateCompany ...
func (s Service) DeactivateCompany(ctx context.Context, companyID string, password string) error {
	span := s.tracer.MakeSpan(ctx, "DeactivateCompany")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	followersAmount, err := s.networkRPC.GetCompanyFollowersNumber(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check user's password
	err = s.userRPC.CheckPassword(ctx, password)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if followersAmount > account.MaximumAmountOfFollowersForDeactivation {
		return errors.New("too_much_followers")
	}

	err = s.repository.Company.ChangeStatusOfCompany(ctx, companyID, status.CompanyStatusDeactivated)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetCompanyAccount ...
func (s Service) GetCompanyAccount(ctx context.Context, companyID string) (*account.Account, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyAccount")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	acc, err := s.repository.Company.GetCompanyAccount(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// get city
	for _, ac := range acc.Addresses {
		if ac.Location.City == nil {
			continue
		}

		ac.Location = location.Location{
			City: &location.City{
				ID: ac.Location.City.ID,
			},
			Country: &location.Country{},
		}

		// if city have id
		if ac.Location.City != nil {
			cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, ac.Location.City.ID, nil)
			if errInfo != nil {
				s.tracer.LogError(span, errInfo)
				// internal_error
				// return "", errInfo
				continue
			}
			// ac.Location.City.ID = ac.Location.City.ID
			ac.Location.City.Name = cityName
			ac.Location.City.Subdivision = subdivision
			ac.Location.Country.ID = countryID
		}
	}

	return acc, nil
}

// ChangeCompanyName ...
func (s Service) ChangeCompanyName(ctx context.Context, companyID string, name string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyName")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//check if name is empty
	err := emptyValidator(name)
	if err != nil {
		return err
	}

	//check if name is over 120 characters
	err = length128Validator(name)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyName(ctx, companyID, name)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyURL ...
func (s Service) ChangeCompanyURL(ctx context.Context, companyID string, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyURL")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}
	// checks if url is not busy
	isURLBusy, err := s.repository.Company.IsURLBusy(ctx, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isURLBusy {
		return errors.New("url_already_taken")
	}
	err = emptyValidator(url)
	if err != nil {
		return err
	}

	//checks if url is valid
	err = urlValidator(url)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyURL(ctx, companyID, url)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyFoundationDate ...
func (s Service) ChangeCompanyFoundationDate(ctx context.Context, companyID string, foundationDate time.Time) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyFoundationDate")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//checks if foundation date is in the future
	err := foundationDateValidator(foundationDate)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyFoundationDate(ctx, companyID, foundationDate)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyIndustry ...
func (s Service) ChangeCompanyIndustry(ctx context.Context, companyID string, industry *account.Industry) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyIndustry")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeCompanyIndustry(ctx, companyID, industry)
	if err != nil {
		return err
	}

	//checks if industry is returned correctly
	err = isNumber(industry.Main)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyType ...
func (s Service) ChangeCompanyType(ctx context.Context, companyID string, companyType account.Type) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyType")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeCompanyType(ctx, companyID, companyType)
	if err != nil {
		return err
	}

	if companyType == account.Type("") {
		return errors.New("Invalid Type")
	}

	return nil
}

// ChangeCompanySize ...
func (s Service) ChangeCompanySize(ctx context.Context, companyID string, size account.Size) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanySize")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeCompanySize(ctx, companyID, size)
	if err != nil {
		return err
	}

	//checks if input is empty
	if size == account.Size("") {
		return errors.New("Invalid Size")
	}

	return nil
}

// AddCompanyEmail ...
func (s Service) AddCompanyEmail(ctx context.Context, companyID string, email *account.Email) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyEmail")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	// TODO: check if email wasn't added before
	isAdded, err := s.repository.Company.IsEmailAdded(ctx, companyID, email.Email)
	if err != nil {
		return "", err
	}
	if isAdded {
		return "", errors.New("email_already_added")
	}

	id := email.GenerateID()

	//checks if email is valid
	err = emailValidator(email.Email)
	if err != nil {
		return "", err
	}

	err = s.repository.Company.AddCompanyEmail(ctx, companyID, email)
	if err != nil {
		return "", err
	}

	// send activation link
	tmpCode, err := s.repository.Cache.CreateTemporaryCodeForEmailActivation(ctx, companyID, email.Email)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return "", err
	}
	// err = s.mailRPC.SendEmail(ctx, email.Email, fmt.Sprint("<html><body><a target='_blank' href='https://"+s.host+"/api/activate/company_email?code=", tmpCode, "&company_id=", companyID, "'>Activate</a></body></html>"))
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	// return "", err
	// }

	emailMessage := fmt.Sprint("<html><body><a target='_blank' href='https://"+s.host+"/api/activate/company_email?code=", tmpCode, "&company_id=", companyID, "'>Activate</a></body></html>")
	log.Println(emailMessage)

	err = s.mq.SendEmail(email.Email, emailMessage)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return id, nil
}

// ActivateEmail ...
func (s Service) ActivateEmail(ctx context.Context, companyID string, code string) error {
	span := s.tracer.MakeSpan(ctx, "ActivateEmail")
	defer span.Finish()

	matched, email, err := s.repository.Cache.CheckTemporaryCodeForEmailActivation(ctx, companyID, code)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}
	if !matched {
		return errors.New("wrong_activation_code")
	}

	err = s.repository.Company.ActivateEmail(ctx, companyID, email)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCompanyEmail ...
func (s Service) DeleteCompanyEmail(ctx context.Context, companyID string, emailID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyEmail")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// TODO:
	// check if it is not primary
	isPrimary, err := s.repository.Company.IsEmailPrimary(ctx, companyID, emailID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isPrimary {
		return errors.New("email_is_primary")
	}

	err = s.repository.Company.DeleteCompanyEmail(ctx, companyID, emailID)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyEmail ...
func (s Service) ChangeCompanyEmail(ctx context.Context, companyID string, emailID string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyEmail")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// TODO:
	// check if email exists by ID

	// check if it activated
	isActivated, err := s.repository.Company.IsEmailActivated(ctx, companyID, emailID)
	if err != nil {
		return err
	}
	if !isActivated {
		return errors.New("email_is_not_activated")
	}

	// make primary
	err = s.repository.Company.MakeEmailPrimary(ctx, companyID, emailID)
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyPhone ...
func (s Service) AddCompanyPhone(ctx context.Context, companyID string, phone *account.Phone) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyPhone")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	// TODO: phone.CountryCode shouldn't be nil

	// check if phone wasn't added before
	isAdded, err := s.repository.Company.IsPhoneAdded(ctx, companyID, phone)
	if err != nil {
		return "", err
	}
	if isAdded {
		return "", errors.New("phone_already_added")
	}

	id := phone.GenerateID()

	// define country code info
	cd, countryID, err := s.infoRPC.GetCountryIDAndCountryCode(ctx, int32(phone.CountryCode.ID))
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return "", err
	}

	if countryID == "" {
		return "", errors.New("bad_country_code_id")
	}

	phone.CountryCode.Code = cd
	phone.CountryCode.CountryID = countryID

	//checks if phone number is valid according to country id and code
	err = phoneValidator(phone.Number, phone.CountryCode.Code, phone.CountryCode.CountryID)
	if err != nil {
		return "", err
	}

	err = s.repository.Company.AddCompanyPhone(ctx, companyID, phone)
	if err != nil {
		return "", err
	}

	// TODO: generate code for phone activation
	// TODO: send SMS

	return id, nil
}

// DeleteCompanyPhone ...
func (s Service) DeleteCompanyPhone(ctx context.Context, companyID string, phoneID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyPhone")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// check if it is not primary

	err := s.repository.Company.DeleteCompanyPhone(ctx, companyID, phoneID)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyPhone ...
func (s Service) ChangeCompanyPhone(ctx context.Context, companyID string, phoneID string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyEmail")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// TODO:
	// check if phone exists
	// isExists, err := s.repository.Company.IsPhoneExists(ctx, companyID, phoneID)
	// if err != nil {
	// 	return err
	// }
	// if !isExists {
	// 	return errors.New("phone_not_found")
	// }

	// TODO: uncomment when we will have SMS gateway
	// check if it activated
	// isActivated, err := s.repository.Company.IsPhoneActivated(ctx, companyID, phoneID)
	// if err != nil {
	// 	return err
	// }
	// if !isActivated {
	// 	return errors.New("phone_is_not_activated")
	// }

	// make primary
	err := s.repository.Company.MakePhonePrimary(ctx, companyID, phoneID)
	if err != nil {
		return err
	}

	return nil
}

// AddCompanyAddress ...
func (s Service) AddCompanyAddress(ctx context.Context, companyID string, address *account.Address) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyAddress")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := address.GenerateID()

	// get info about city
	cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, address.Location.City.ID, nil)
	if errInfo != nil {
		s.tracer.LogError(span, errInfo)
		// internal_error
		return "", errInfo
	}
	address.Location.City.Name = cityName
	address.Location.City.Subdivision = subdivision
	address.Location.Country.ID = countryID

	for i := range address.Phones {
		_ = address.Phones[i].GenerateID()

		// define country code info
		cd, countryID, err := s.infoRPC.GetCountryIDAndCountryCode(ctx, int32(address.Phones[i].CountryCode.ID))
		if err != nil {
			s.tracer.LogError(span, err)
			// internal_error
			return "", err
		}

		if countryID == "" {
			return "", errors.New("bad_country_code_id")
		}

		address.Phones[i].CountryCode.Code = cd
		address.Phones[i].CountryCode.CountryID = countryID
	}

	for i := range address.BusinessHours {
		_ = address.BusinessHours[i].GenerateID()
	}

	//checks if address is valid.(empty/ is over 120 characters apartment,zipcode,street )
	err := addressValidator(address)
	if err != nil {
		return "", err
	}

	err = s.repository.Company.AddCompanyAddress(ctx, companyID, address)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// ChangeCompanyAddress ...
func (s Service) ChangeCompanyAddress(ctx context.Context, companyID string, address *account.Address) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyAddress")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//checks if address is valid.(empty/ is over 120 characters apartment,zipcode,street )

	err := addressValidator(address)
	if err != nil {
		return err
	}

	for i := range address.BusinessHours {
		if address.BusinessHours[i].GetID() == "" {
			address.BusinessHours[i].GenerateID()
		}
	}

	for i := range address.Websites {
		if address.Websites[i].GetID() == "" {
			address.Websites[i].GenerateID()
		}
	}

	for i := range address.Phones {
		if address.Phones[i].GetID() == "" {
			address.Phones[i].GenerateID()
		}
	}

	for i := range address.Phones {
		_ = address.Phones[i].GenerateID()

		// define country code info
		cd, countryID, err := s.infoRPC.GetCountryIDAndCountryCode(ctx, int32(address.Phones[i].CountryCode.ID))
		if err != nil {
			s.tracer.LogError(span, err)
			// internal_error
			return err
		}

		if countryID == "" {
			return errors.New("bad_country_code_id")
		}

		address.Phones[i].CountryCode.Code = cd
		address.Phones[i].CountryCode.CountryID = countryID
	}

	err = s.repository.Company.ChangeCompanyAddress(ctx, companyID, address)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// DeleteCompanyAddress ...
func (s Service) DeleteCompanyAddress(ctx context.Context, companyID string, addressID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyAddress")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// check if it is not primary

	err := s.repository.Company.DeleteCompanyAddress(ctx, companyID, addressID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddCompanyWebsite ...
func (s Service) AddCompanyWebsite(ctx context.Context, companyID string, website string) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyWebsite")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	//checks if website is valid.
	err := emptyValidator(website)
	if err != nil {
		return "", err
	}
	err = urlValidator(website)
	if err != nil {
		return "", err
	}

	//checks if it is over 120 characters
	err = length128Validator(website)
	if err != nil {
		return "", err
	}

	ws := account.Website{
		Site: website,
	}

	id := ws.GenerateID()

	err = s.repository.Company.AddCompanyWebsite(ctx, companyID, &ws)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// DeleteCompanyWebsite ...
func (s Service) DeleteCompanyWebsite(ctx context.Context, companyID string, websiteID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyWebsite")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.DeleteCompanyWebsite(ctx, companyID, websiteID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCompanyWebsite ...
func (s Service) ChangeCompanyWebsite(ctx context.Context, companyID string, websiteID string, website string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyWebsite")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeCompanyWebsite(ctx, companyID, websiteID, website)
	if err != nil {
		return err
	}

	return nil
}

// ChangeCompanyParking ...
func (s Service) ChangeCompanyParking(ctx context.Context, companyID string, parking account.Parking) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyParking")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeCompanyParking(ctx, companyID, parking)
	if err != nil {
		return err
	}

	//cheks that it doesn't return empty input
	if parking == account.Parking("") {
		return errors.New("Invalid Parking")
	}

	return nil
}

// ChangeCompanyBenefits ...
func (s Service) ChangeCompanyBenefits(ctx context.Context, companyID string, benefits []profile.Benefit) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyBenefits")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeCompanyBenefits(ctx, companyID, benefits)
	if err != nil {
		return err
	}

	//cheks that it doesn't return empty input
	for _, t := range benefits {
		if t == profile.Benefit("") {
			return errors.New("Empty Benefits Enter")
		}
	}

	return nil
}

// AddCompanyAdmin ...
func (s Service) AddCompanyAdmin(ctx context.Context, companyID string, userID string, level account.AdminLevel, password string) error {
	span := s.tracer.MakeSpan(ctx, "AddCompanyAdmin")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// check user's password
	err := s.userRPC.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	err = s.networkRPC.AddCompanyAdmin(ctx, companyID, userID, level)
	if err != nil {
		return err
	}

	// TODO: send notification

	return nil
}

// DeleteCompanyAdmin ...
func (s Service) DeleteCompanyAdmin(ctx context.Context, companyID string, userID string, password string) error {
	span := s.tracer.MakeSpan(ctx, "AddCompanyAdmin")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// check user's password
	err := s.userRPC.CheckPassword(ctx, password)
	if err != nil {
		return err
	}

	err = s.networkRPC.DeleteCompanyAdmin(ctx, companyID, userID)
	if err != nil {
		return err
	}

	return nil
}

// profile

// GetCompanyProfile ...
func (s Service) GetCompanyProfile(ctx context.Context, url string, lang string) (*profile.Profile, account.AdminLevel, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyProfile")
	defer span.Finish()

	prof, err := s.repository.Company.GetCompanyProfile(ctx, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, account.AdminLevelUnknown, err
	}

	if lang == "" {
		// take from cookie
		lang = s.retriveUILang(ctx)
	}

	err = s.processCompanyProfile(ctx, prof, lang)
	if err != nil {
		return nil, account.AdminLevelUnknown, err
	}

	//getting admin level
	lvl, err := s.networkRPC.GetAdminLevel(ctx, prof.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, account.AdminLevelUnknown, err
	}

	return prof, lvl, nil
}

// GetCompanyProfileByID ...
func (s Service) GetCompanyProfileByID(ctx context.Context, companyID string, lang string) (*profile.Profile, account.AdminLevel, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyProfileByID")
	defer span.Finish()

	// maybe apply some permissions?

	prof, err := s.repository.Company.GetCompanyProfileByID(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, account.AdminLevelUnknown, err
	}

	// // get city
	// for _, pr := range prof.Addresses {
	// 	if pr.Location.City == nil {
	// 		continue
	// 	}
	//
	// 	pr.Location = location.Location{
	// 		City: &location.City{
	// 			ID: pr.Location.City.ID,
	// 		},
	// 		Country: &location.Country{},
	// 	}
	//
	// 	// if city have id
	// 	if pr.Location.City != nil {
	// 		cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, pr.Location.City.ID, nil)
	// 		if errInfo != nil {
	// 			s.tracer.LogError(span, errInfo)
	// 			// internal_error
	// 			// return "", errInfo
	// 			continue
	// 		}
	// 		// pr.Location.City.ID = pr.Location.City.ID
	// 		pr.Location.City.Name = cityName
	// 		pr.Location.City.Subdivision = subdivision
	// 		pr.Location.Country.ID = countryID
	// 	}
	// }
	//
	// prof.Translate(ctx, lang)
	//
	// followings, followers, employees, err := s.networkRPC.GetCompanyCountings(ctx, prof.GetID())
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// }
	//
	// prof.AmountOfFollowers = followers
	// prof.AmountOfFollowings = followings
	// prof.AmountOfEmployees = employees
	// // prof.AmountOfJobs
	//
	// // is favourite
	// prof.IsFavorite, err = s.networkRPC.IsFavourite(ctx, prof.GetID())
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// }
	//
	// // is follow
	// prof.IsFollow, err = s.networkRPC.IsFollow(ctx, prof.GetID())
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// }
	//
	// // is blocked
	// prof.IsBlocked, err = s.networkRPC.IsBlockedCompany(ctx, prof.GetID())
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// }
	//
	// // get avarage rating
	// avg, _, err := s.repository.Reviews.GetAvarageRateOfCompany(ctx, prof.GetID())
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// }
	// prof.AvarageRating = avg

	err = s.processCompanyProfile(ctx, prof, lang)
	if err != nil {
		return nil, account.AdminLevelUnknown, err
	}

	//getting admin level
	lvl, err := s.networkRPC.GetAdminLevel(ctx, prof.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, account.AdminLevelUnknown, err
	}

	// // checks if company is online
	// prof.IsOnline, err = s.chatRPC.IsLive(ctx, prof.GetID())
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return nil, account.AdminLevelUnknown, err
	// }

	return prof, lvl, nil
}

// GetCompanyProfiles ...
func (s Service) GetCompanyProfiles(ctx context.Context, ids []string) ([]*profile.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyProfiles")
	defer span.Finish()

	// maybe apply some permissions?

	profs, err := s.repository.Company.GetCompanyProfiles(ctx, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for _, prof := range profs {
		err = s.processCompanyProfile(ctx, prof, "en") // TODO: lang
		if err != nil {
			return nil, err
		}
		// // get city
		// for _, pr := range prof.Addresses {
		// 	if pr.Location.City == nil {
		// 		continue
		// 	}
		//
		// 	pr.Location = location.Location{
		// 		City: &location.City{
		// 			ID: pr.Location.City.ID,
		// 		},
		// 		Country: &location.Country{},
		// 	}
		//
		// 	// if city have id
		// 	if pr.Location.City != nil {
		// 		cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, pr.Location.City.ID, nil)
		// 		if errInfo != nil {
		// 			s.tracer.LogError(span, errInfo)
		// 			// internal_error
		// 			// return "", errInfo
		// 			continue
		// 		}
		// 		// pr.Location.City.ID = pr.Location.City.ID
		// 		pr.Location.City.Name = cityName
		// 		pr.Location.City.Subdivision = subdivision
		// 		pr.Location.Country.ID = countryID
		// 	}
		// }
		//
		// // checks if company is online
		// prof.IsOnline, err = s.chatRPC.IsLive(ctx, prof.GetID())
		// if err != nil {
		// 	s.tracer.LogError(span, err)
		// 	return nil, err
		// }
		//
		// // is follow
		// prof.IsFollow, err = s.networkRPC.IsFollow(ctx, prof.GetID())
		// if err != nil {
		// 	s.tracer.LogError(span, err)
		// }
	}

	return profs, nil
}

// GetCompanyProfilesMap ...
func (s Service) GetCompanyProfilesMap(ctx context.Context, ids []string) (map[string]*profile.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyProfiles")
	defer span.Finish()

	// maybe apply some permissions?

	profs, err := s.repository.Company.GetCompanyProfiles(ctx, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// checks if company is online
	for _, prof := range profs {
		err = s.processCompanyProfile(ctx, prof, "en") // TODO: lang
		if err != nil {
			return nil, err
		}

		// prof.IsOnline, err = s.chatRPC.IsLive(ctx, prof.GetID())
		// if err != nil {
		// 	s.tracer.LogError(span, err)
		// 	return nil, err
		// }
	}

	companyProfilesMap := make(map[string]*profile.Profile, len(profs))

	for _, pr := range profs {
		companyProfilesMap[pr.GetID()] = pr
	}

	return companyProfilesMap, nil
}

// ChangeCompanyAboutUs ...
func (s Service) ChangeCompanyAboutUs(ctx context.Context, companyID string, aboutUs *profile.AboutUs) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyAboutUs")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	for i := range aboutUs.BusinessHours {
		aboutUs.BusinessHours[i].GenerateID()
	}

	//checks length of fields, emptiness, checks foundation is not in future
	//checks parking, type and size not to return empty
	err := aboutUsValidator(aboutUs)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyAboutUs(ctx, companyID, aboutUs)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetFounders ...
func (s Service) GetFounders(ctx context.Context, companyID string, first int32, after string) ([]*profile.Founder, error) {
	span := s.tracer.MakeSpan(ctx, "GetFounders")
	defer span.Finish()

	var afterNumber int
	if after != "" {
		afterNumber, err := strconv.Atoi(after)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, errors.New("bad_after_value")
		}
		if afterNumber < 0 {
			return nil, errors.New("bad_after_value")
		}
	}

	// retrive id of user
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// Check if user is admin or owner
	lvl, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		// return nil, err
	}
	// if lvl != account.AdminLevelAdmin {
	// return nil, errors.New("you_don't_have_permission")
	// }

	// get founders
	founders, err := s.repository.Company.GetFounders(ctx, companyID, first, afterNumber)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	filteredFounders := make([]*profile.Founder, 0, len(founders))

	// filter if user is not admin
	if lvl != account.AdminLevelAdmin {
		for i := range founders {
			if (userID != "" && founders[i].GetUserID() == userID) || founders[i].Approved == true || founders[i].GetUserID() == "" {
				filteredFounders = append(filteredFounders, founders[i])
			}
		}

		return filteredFounders, nil
	}

	return founders, nil
}

// AddCompanyFounder ...
func (s Service) AddCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyFounder")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := founder.GenerateID()
	founder.CreatedAt = time.Now()

	//checks fields are not empty
	err := emptyValidator(founder.Name, founder.Position)
	if err != nil {
		return "", err
	}

	//checks field's length is not over 120 characters
	err = length128Validator(founder.Name, founder.Position)
	if err != nil {
		return "", err
	}

	err = dashAndSpace(founder.Name)
	if err != nil {
		return "", err
	}

	if founder.GetUserID() == "" {
		founder.Approved = true
	}

	err = s.repository.Company.AddCompanyFounder(ctx, companyID, founder)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// send notification
	if founder.GetUserID() != "" {
		err = s.mq.SendNewFounderRequest(companyID, &notmes.NewFounderRequest{
			Founder:   founder.GetUserID(),
			RequestID: id,
		})
		if err != nil {
			s.tracer.LogError(span, err)
		}
	}

	return id, nil
}

// DeleteCompanyFounder ...
func (s Service) DeleteCompanyFounder(ctx context.Context, companyID string, founderID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyFounder")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.DeleteCompanyFounder(ctx, companyID, founderID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ApproveFounderRequest ...
func (s Service) ApproveFounderRequest(ctx context.Context, companyID, requestID string) error {
	span := s.tracer.MakeSpan(ctx, "ApproveFounderRequest")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Company.ApproveFounderRequest(ctx, companyID, requestID, userID, true) // TODO: check if admin of company or user sent this request
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveFounderRequest ...
func (s Service) RemoveFounderRequest(ctx context.Context, companyID, requestID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFounderRequest")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Company.RemoveFounderRequest(ctx, companyID, requestID, userID, true) // TODO: check if admin of company or user sent this request
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCompanyFounder ...
func (s Service) ChangeCompanyFounder(ctx context.Context, companyID string, founder *profile.Founder) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyFounder")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//checks field's length is not over 120 characters
	err := length128Validator(founder.Name, founder.Position)
	if err != nil {
		return err
	}
	//checks fields are not empty
	err = emptyValidator(founder.Name, founder.Position)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyFounder(ctx, companyID, founder)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCompanyFounderAvatar ...
func (s Service) ChangeCompanyFounderAvatar(ctx context.Context, companyID string, founderID string, file *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyFounderAvatar")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := file.GenerateID()

	err := s.repository.Company.ChangeCompanyFounderAvatar(ctx, companyID, founderID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveCompanyFounderAvatar ...
func (s Service) RemoveCompanyFounderAvatar(ctx context.Context, companyID string, founderID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveCompanyFounderAvatar")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.RemoveCompanyFounderAvatar(ctx, companyID, founderID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddCompanyAward ...
func (s Service) AddCompanyAward(ctx context.Context, companyID string, award *profile.Award) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyAward")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := award.GenerateID()
	award.CreatedAt = time.Now()

	for i := range award.Links {
		award.Links[i].GenerateID()
	}

	err := awardValidator(award)
	if err != nil {
		return "", err
	}

	err = s.repository.Company.AddCompanyAward(ctx, companyID, award)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// DeleteCompanyAward ...
func (s Service) DeleteCompanyAward(ctx context.Context, companyID string, awardID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyAward")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.DeleteCompanyAward(ctx, companyID, awardID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCompanyAward ...
func (s Service) ChangeCompanyAward(ctx context.Context, companyID string, award *profile.Award) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyParking")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//checks date is not in future, checks title and issuer length(120 characters)/emptiness
	err := awardValidator(award)
	if err != nil {
		return err
	}

	for i := range award.Links {
		award.Links[i].GenerateID()
	}

	err = s.repository.Company.ChangeCompanyAward(ctx, companyID, award)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddLinksInCompanyAward ...
func (s Service) AddLinksInCompanyAward(ctx context.Context, companyID, awardID string, links []*profile.Link) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "AddLinksInCompanyAward")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	// userID, err := s.authRPC.GetUserID(ctx, token)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return nil, err
	// }

	ids := make([]string, 0, len(links))

	for i := range links {
		ids = append(ids, links[i].GenerateID())
	}

	err := s.repository.Company.AddLinksInCompanyAward(ctx, companyID, awardID, links)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return ids, nil
}

// AddFileInCompanyAward add information about file to user's CompanyAward. Called in file-manager service.
// It doesn't verify who tries to add!
func (s Service) AddFileInCompanyAward(ctx context.Context, companyID, awardID string, file *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInCompanyAward")
	defer span.Finish()

	id := file.GenerateID()

	err := s.repository.Company.AddFileInCompanyAward(ctx, companyID, awardID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveFilesInCompanyAward ...
func (s Service) RemoveFilesInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInCompanyAward")
	defer span.Finish()

	err := s.repository.Company.RemoveFilesInCompanyAward(ctx, companyID, awardID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeLinkInCompanyAward ...
func (s Service) ChangeLinkInCompanyAward(ctx context.Context, companyID, awardID, linkID, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeLinkInCompanyAward")
	defer span.Finish()

	err := s.repository.Company.ChangeLinkInCompanyAward(ctx, companyID, awardID, linkID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveLinksInCompanyAward ...
func (s Service) RemoveLinksInCompanyAward(ctx context.Context, companyID, awardID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveLinksInCompanyAward")
	defer span.Finish()

	err := s.repository.Company.RemoveLinksInCompanyAward(ctx, companyID, awardID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetUploadedFilesInCompanyAward ...
func (s Service) GetUploadedFilesInCompanyAward(ctx context.Context, companyID string) ([]*profile.File, error) {
	span := s.tracer.MakeSpan(ctx, "GetUploadedFilesInCompanyAward")
	defer span.Finish()

	files, err := s.repository.Company.GetUploadedFilesInCompanyAward(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return files, nil
}

// GetCompanyGallery ...
func (s Service) GetCompanyGallery(ctx context.Context, companyID string, first, afterNumber uint32) ([]*profile.File, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyGallery")
	defer span.Finish()

	// var afterNumber int
	// if after != "" {
	// 	afterNumber, err := strconv.Atoi(after)
	// 	if err != nil {
	// 		s.tracer.LogError(span, err)
	// 		return nil, errors.New("bad_after_value")
	// 	}
	// 	if afterNumber < 0 {
	// 		return nil, errors.New("bad_after_value")
	// 	}
	// }

	// // retrive id of user
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return nil, err
	// }
	//
	// // Check if user is admin or owner
	// lvl, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	// return nil, err
	// }
	// // if lvl != account.AdminLevelAdmin {
	// // return nil, errors.New("you_don't_have_permission")
	// // }

	// get galleries
	galleries, err := s.repository.Company.GetCompanyGallery(ctx, companyID, first, afterNumber)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// filteredGalleries := make([]*profile.File, 0, len(galleries))

	// // filter if user is not admin
	// if lvl != account.AdminLevelAdmin {
	// 	for i := range galleries {
	// 		if userID != "" && galleries[i].GetID() == userID {
	// 			filteredGalleries = append(filteredGalleries, galleries[i])
	// 		}
	// 	}
	//
	// 	return filteredGalleries, nil
	// }

	return galleries, nil
}

// AddFileInCompanyGallery add information about file to user's CompanyGallery. Called in file-manager service.
func (s Service) AddFileInCompanyGallery(ctx context.Context, companyID string, file *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInCompanyGallery")
	defer span.Finish()

	id := file.GenerateID()

	err := s.repository.Company.AddFileInCompanyGallery(ctx, companyID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveFilesInCompanyGallery ...
func (s Service) RemoveFilesInCompanyGallery(ctx context.Context, companyID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInCompanyGallery")
	defer span.Finish()

	err := s.repository.Company.RemoveFilesInCompanyGallery(ctx, companyID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetUploadedFilesInCompanyGallery ...
func (s Service) GetUploadedFilesInCompanyGallery(ctx context.Context, companyID string) ([]*profile.File, error) {
	span := s.tracer.MakeSpan(ctx, "GetUploadedFilesInCompanyGallery")
	defer span.Finish()

	files, err := s.repository.Company.GetUploadedFilesInCompanyGallery(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return files, nil
}

// AddCompanyMilestone ...
func (s Service) AddCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyMilestone")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := milestone.GenerateID()
	milestone.CreatedAt = time.Now()

	//checks date is not in future, checks title and description length(120 characters)/emptiness
	err := milestoneValidator(milestone)
	if err != nil {
		return "", err
	}

	err = s.repository.Company.AddCompanyMilestone(ctx, companyID, milestone)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// DeleteCompanyMilestone ...
func (s Service) DeleteCompanyMilestone(ctx context.Context, companyID string, milestoneID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyMilestone")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.DeleteCompanyMilestone(ctx, companyID, milestoneID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCompanyMilestone ...
func (s Service) ChangeCompanyMilestone(ctx context.Context, companyID string, milestone *profile.Milestone) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyMilestone")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//checks date is not in future, checks title and description length(120 characters)/emptiness
	err := milestoneValidator(milestone)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyMilestone(ctx, companyID, milestone)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeImageMilestone ...
func (s Service) ChangeImageMilestone(ctx context.Context, companyID string, milestoneID string, image *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeImageMilestone")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := image.GenerateID()

	err := s.repository.Company.ChangeImageMilestone(ctx, companyID, milestoneID, image)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveImageInMilestone ...
func (s Service) RemoveImageInMilestone(ctx context.Context, companyID string, milestoneID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveImageInMilestone")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.RemoveImageInMilestone(ctx, companyID, milestoneID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddCompanyProduct ...
func (s Service) AddCompanyProduct(ctx context.Context, companyID string, product *profile.Product) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyProduct")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVShop,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := product.GenerateID()
	product.CreatedAt = time.Now()

	//checks for emptiness, checks for name length(120 characters).
	err := productValidator(product)
	if err != nil {
		return "", err
	}

	err = s.repository.Company.AddCompanyProduct(ctx, companyID, product)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// DeleteCompanyProduct ...
func (s Service) DeleteCompanyProduct(ctx context.Context, companyID string, productID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyProduct")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVShop,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.DeleteCompanyProduct(ctx, companyID, productID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCompanyProduct ...
func (s Service) ChangeCompanyProduct(ctx context.Context, companyID string, product *profile.Product) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyProduct")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVShop,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//checks for emptiness, checks for name length(120 characters).
	err := productValidator(product)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyProduct(ctx, companyID, product)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeImageProduct ...
func (s Service) ChangeImageProduct(ctx context.Context, companyID string, productID string, image *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeImageProduct")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVShop,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := image.GenerateID()

	err := s.repository.Company.ChangeImageProduct(ctx, companyID, productID, image)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveImageInProduct ...
func (s Service) RemoveImageInProduct(ctx context.Context, companyID string, productID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveImageInProduct")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVShop,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.RemoveImageInProduct(ctx, companyID, productID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddCompanyService ...
func (s Service) AddCompanyService(ctx context.Context, companyID string, service *profile.Service) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyService")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVService,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := service.GenerateID()
	service.CreatedAt = time.Now()

	//checks for emptiness, checks for name length(120 characters).
	err := serviceValidator(service)
	if err != nil {
		return "", err
	}

	err = s.repository.Company.AddCompanyService(ctx, companyID, service)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// DeleteCompanyService ...
func (s Service) DeleteCompanyService(ctx context.Context, companyID string, serviceID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyService")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVService,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.DeleteCompanyService(ctx, companyID, serviceID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeCompanyService ...
func (s Service) ChangeCompanyService(ctx context.Context, companyID string, service *profile.Service) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanyService")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVService,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//checks for emptiness, checks for name length(120 characters).
	err := serviceValidator(service)
	if err != nil {
		return err
	}

	err = s.repository.Company.ChangeCompanyService(ctx, companyID, service)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeImageService ...
func (s Service) ChangeImageService(ctx context.Context, companyID string, serviceID string, image *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeImageService")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVService,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	id := image.GenerateID()

	err := s.repository.Company.ChangeImageService(ctx, companyID, serviceID, image)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveImageInService ...
func (s Service) RemoveImageInService(ctx context.Context, companyID string, serviceID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveImageInService")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
		account.AdminLevelVService,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.RemoveImageInService(ctx, companyID, serviceID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddCompanyReport ...
func (s Service) AddCompanyReport(ctx context.Context, companyID string, report *profile.Report) error {
	span := s.tracer.MakeSpan(ctx, "AddCompanyReport")
	defer span.Finish()

	report.GenerateID()
	report.CreatedAt = time.Now()

	//check only for description length
	err := length500Validator(report.Description)
	if err != nil {
		return err
	}

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	report.SetCreatorID(userID)

	err = s.repository.Company.AddCompanyReport(ctx, companyID, report)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddCompanyReview ...
func (s Service) AddCompanyReview(ctx context.Context, companyID string, review *profile.Review) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddCompanyReview")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// TODO: check if user not added review for this company before

	id := review.GenerateID()
	review.CreatedAt = time.Now()

	err = review.SetAuthorID(userID)
	if err != nil {
		return "", err
	}

	err = review.SetCompanyID(companyID)
	if err != nil {
		return "", err
	}
	//checks for length(120 characters) and emptiness of headline and description
	err = reviewValidator(review)
	if err != nil {
		return "", err
	}

	err = s.repository.Reviews.AddReview(ctx, review)
	if err != nil {
		return "", err
	}

	// send notification
	err = s.mq.SendNewCompanyReview(companyID, &notmes.NewCompanyReview{
		ReviewID:       id,
		ReviewerUserID: userID,
	})
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return id, nil
}

// DeleteCompanyReview ...
func (s Service) DeleteCompanyReview(ctx context.Context, companyID string, reviewID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteCompanyReview")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Reviews.DeleteCompanyReview(ctx, companyID, userID, reviewID)
	if err != nil {
		return err
	}

	return nil
}

// GetCompanyReviews ...
func (s Service) GetCompanyReviews(ctx context.Context, companyID string, first uint32, after string) ([]*profile.Review, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyReviews")
	defer span.Finish()

	var afterNumber int
	if after != "" {
		afterNumber, err := strconv.Atoi(after)
		if err != nil {
			s.tracer.LogError(span, err)
			log.Println("error: after has a bad value:", after)
			return nil, errors.New("bad_after_value")
		}
		if afterNumber < 0 {
			log.Println("error: after has a bad value:", after)
			return nil, errors.New("bad_after_value")
		}
	}

	reviews, err := s.repository.Reviews.GetCompanyReviews(ctx, companyID, first, uint32(afterNumber))
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// GetUsersRevies ...
func (s Service) GetUsersRevies(ctx context.Context, userID string, first uint32, after string) ([]*profile.Review, error) {
	span := s.tracer.MakeSpan(ctx, "GetUsersRevies")
	defer span.Finish()

	var afterNumber int
	if after != "" {
		afterNumber, err := strconv.Atoi(after)
		if err != nil {
			s.tracer.LogError(span, err)
			log.Println("error: after has a bad value:", after)
			return nil, errors.New("bad_after_value")
		}
		if afterNumber < 0 {
			log.Println("error: after has a bad value:", after)
			return nil, errors.New("bad_after_value")
		}
	}

	reviews, err := s.repository.Reviews.GetUsersRevies(ctx, userID, first, uint32(afterNumber))
	if err != nil {
		return nil, err
	}

	companyIDs := make([]string, 0, len(reviews))

	for _, r := range reviews {
		if r != nil {
			companyIDs = append(companyIDs, r.GetCompanyID())
		}
	}

	profiles, err := s.GetCompanyProfilesMap(ctx, companyIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for _, r := range reviews {
		prof := profiles[r.GetCompanyID()]
		if prof != nil {
			r.Company = *prof
		}
	}

	return reviews, nil
}

// GetAvarageRateOfCompany ...
func (s Service) GetAvarageRateOfCompany(ctx context.Context, companyID string) (float32, uint32, error) {
	span := s.tracer.MakeSpan(ctx, "GetAvarageRateOfCompany")
	defer span.Finish()

	avg, amount, err := s.repository.Reviews.GetAvarageRateOfCompany(ctx, companyID)
	if err != nil {
		return 0, 0, err
	}

	return avg, amount, nil
}

// GetAmountOfEachRate ...
func (s Service) GetAmountOfEachRate(ctx context.Context, companyID string) (map[uint32]uint32, error) {
	span := s.tracer.MakeSpan(ctx, "GetAmountOfEachRate")
	defer span.Finish()

	rates, err := s.repository.Reviews.GetAmountOfEachRate(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

// AddCompanyReviewReport ...
func (s Service) AddCompanyReviewReport(ctx context.Context, reviewReport *profile.ReviewReport) error {
	span := s.tracer.MakeSpan(ctx, "AddCompanyReviewReport")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	reviewReport.SetAuthorID(userID)

	reviewReport.GenerateID()
	reviewReport.CreatedAt = time.Now()

	err = s.repository.Reviews.AddCompanyReviewReport(ctx, reviewReport)
	if err != nil {
		return err
	}

	return nil
}

// ChangeAvatar ...
func (s Service) ChangeAvatar(ctx context.Context, companyID string, file *profile.File) error {
	span := s.tracer.MakeSpan(ctx, "ChangeAvatar")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeAvatar(ctx, companyID, file.URL)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeOriginAvatar ...
func (s Service) ChangeOriginAvatar(ctx context.Context, companyID string, file *profile.File) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOriginAvatar")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeOriginAvatar(ctx, companyID, file.URL)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveAvatar ...
func (s Service) RemoveAvatar(ctx context.Context, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveAvatar")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.RemoveAvatar(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetOriginAvatar ...
func (s Service) GetOriginAvatar(ctx context.Context, companyID string) (string, error) {
	span := s.tracer.MakeSpan(ctx, "GetOriginAvatar")
	defer span.Finish()

	// Check if user is admin or owner
	// lvl, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	// if err != nil {
	// 	return err
	// }
	// if lvl != account.AdminLevelAdmin {
	// 	return errors.New("you_don't_have_permission")
	// }

	url, err := s.repository.Company.GetOriginAvatar(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return url, nil
}

// ChangeCover ...
func (s Service) ChangeCover(ctx context.Context, companyID string, file *profile.File) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCover")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeCover(ctx, companyID, file.URL)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeOriginCover ...
func (s Service) ChangeOriginCover(ctx context.Context, companyID string, file *profile.File) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOriginCover")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.ChangeOriginCover(ctx, companyID, file.URL)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveCover ...
func (s Service) RemoveCover(ctx context.Context, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveCover")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.RemoveCover(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetOriginCover ...
func (s Service) GetOriginCover(ctx context.Context, companyID string) (string, error) {
	span := s.tracer.MakeSpan(ctx, "GetOriginCover")
	defer span.Finish()

	// Check if user is admin or owner
	// lvl, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	// if err != nil {
	// 	return err
	// }
	// if lvl != account.AdminLevelAdmin {
	// 	return errors.New("you_don't_have_permission")
	// }

	url, err := s.repository.Company.GetOriginCover(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return url, nil
}

// SaveCompanyProfileTranslation ...
func (s Service) SaveCompanyProfileTranslation(ctx context.Context, companyID string, lang string, tr *profile.Translation) error {
	span := s.tracer.MakeSpan(ctx, "SaveCompanyProfileTranslation")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.SaveCompanyProfileTranslation(ctx, companyID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveCompanyMilestoneTranslation ...
func (s Service) SaveCompanyMilestoneTranslation(ctx context.Context, companyID string, milestoneID string, language string, translation *profile.MilestoneTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveCompanyMilestoneTranslation")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.SaveCompanyMilestoneTranslation(ctx, companyID, milestoneID, language, translation)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveCompanyAwardTranslation ...
func (s Service) SaveCompanyAwardTranslation(ctx context.Context, companyID string, awardID string, language string, translation *profile.AwardTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveCompanyAwardTranslation")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.repository.Company.SaveCompanyAwardTranslation(ctx, companyID, awardID, language, translation)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetAmountOfReviewsOfUser ...
func (s Service) GetAmountOfReviewsOfUser(ctx context.Context, userID string) (int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetAmountOfReviewsOfUser")
	defer span.Finish()

	amount, err := s.repository.Reviews.GetAmountOfReviewsOfUser(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return 0, err
	}

	return amount, nil
}

// AddGoldCoinsToWallet ...
func (s Service) AddGoldCoinsToWallet(ctx context.Context, userID string, coins int32) error {
	span := s.tracer.MakeSpan(ctx, "AddGoldCoinsToWallet")
	defer span.Finish()

	err := s.stuffRPC.AddGoldCoinsToWallet(ctx, userID, coins)

	if err != nil {
		s.tracer.LogError(span, err)

		return err
	}

	return nil
}

// OpenCareerCenter ...
func (s Service) OpenCareerCenter(ctx context.Context, companyID string, cc *careercenter.CareerCenter) error {
	span := s.tracer.MakeSpan(ctx, "OpenCareerCenter")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		account.AdminLevelAdmin,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// TODO: check if not opened

	cc.IsOpened = true
	cc.CreatedAt = time.Now()

	err := s.repository.Company.OpenCareerCenter(ctx, companyID, cc)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// -------------

// return false if level doesn't much
func (s Service) checkAdminLevel(ctx context.Context, companyID string, requiredLevels ...account.AdminLevel) bool {
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

	// return true
}

func (s Service) processCompanyProfile(ctx context.Context, prof *profile.Profile, lang string) error {
	span := s.tracer.MakeSpan(ctx, "processCompanyProfile")
	defer span.Finish()

	// get city
	for _, pr := range prof.Addresses {
		if pr.Location.City == nil {
			continue
		}

		pr.Location = location.Location{
			City: &location.City{
				ID: pr.Location.City.ID,
			},
			Country: &location.Country{},
		}

		// if city have id
		if pr.Location.City != nil {
			cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, pr.Location.City.ID, &lang)
			if errInfo != nil {
				s.tracer.LogError(span, errInfo)
				// internal_error
				// return "", errInfo
				continue
			}
			// pr.Location.City.ID = pr.Location.City.ID
			pr.Location.City.Name = cityName
			pr.Location.City.Subdivision = subdivision
			pr.Location.Country.ID = countryID
		}
	}

	prof.CurrentTranslation = prof.Translate(ctx, lang)
	prof.AvailableTranslations = make([]string, 0, len(prof.Translation))
	prof.AvailableTranslations = append(prof.AvailableTranslations, "en") // english is default
	for key := range prof.Translation {
		prof.AvailableTranslations = append(prof.AvailableTranslations, key)
	}

	followings, followers, employees, err := s.networkRPC.GetCompanyCountings(ctx, prof.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
	}

	prof.AmountOfFollowers = followers
	prof.AmountOfFollowings = followings
	prof.AmountOfEmployees = employees
	prof.AmountOfJobs, err = s.jobsRPC.GetAmountOfActiveJobsOfCompany(ctx, prof.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
	}

	companyID := s.retriveCompanyID(ctx)

	if companyID != "" {
		isBlockedByCompany, err := s.networkRPC.IsBlockedCompanyByCompany(ctx, prof.GetID(), companyID)
		if err != nil {
			return err
		}
		if isBlockedByCompany {
			return errors.New("you_are_blocked")
		}

		// is follow
		prof.IsFollow, err = s.networkRPC.IsFollowForCompany(ctx, prof.GetID(), companyID)
		if err != nil {
			s.tracer.LogError(span, err)
		}

		// is blocked
		prof.IsBlocked, err = s.networkRPC.IsBlockedCompanyForCompany(ctx, prof.GetID(), companyID)
		if err != nil {
			s.tracer.LogError(span, err)
		}

		// // is favourite
		// prof.IsFavorite, err = s.networkRPC.IsFavouriteCompany(ctx, prof.GetID(), companyID)
		// if err != nil {
		// 	s.tracer.LogError(span, err)
		// }

	} else {
		isBlockedByCompany, err := s.networkRPC.IsBlockedCompanyByUser(ctx, prof.GetID())
		if err != nil {
			return err
		}
		if isBlockedByCompany {
			return errors.New("you_are_blocked")
		}

		// is favourite
		prof.IsFavorite, err = s.networkRPC.IsFavourite(ctx, prof.GetID())
		if err != nil {
			s.tracer.LogError(span, err)
		}

		// is follow
		prof.IsFollow, err = s.networkRPC.IsFollow(ctx, prof.GetID())
		if err != nil {
			s.tracer.LogError(span, err)
		}

		// is blocked
		prof.IsBlocked, err = s.networkRPC.IsBlockedCompany(ctx, prof.GetID())
		if err != nil {
			s.tracer.LogError(span, err)
		}
	}

	// get avarage rating
	avg, _, err := s.repository.Reviews.GetAvarageRateOfCompany(ctx, prof.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
	}
	prof.AvarageRating = avg

	// checks if company is online
	prof.IsOnline, err = s.chatRPC.IsLive(ctx, prof.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

func (s Service) retriveCompanyID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("company_id")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}

func (s Service) retriveUILang(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("ui_lang")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}
