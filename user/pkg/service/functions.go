package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
	twoFA "gitlab.lan/Rightnao-site/microservices/user/pkg/internal/two-fa"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/invitation"
	notmes "gitlab.lan/Rightnao-site/microservices/user/pkg/notification_messages"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
	arangorepo "gitlab.lan/Rightnao-site/microservices/user/pkg/repository/arango"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/status"
	userReport "gitlab.lan/Rightnao-site/microservices/user/pkg/user_report"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

// IdentifyCountry ...
func (s Service) IdentifyCountry(ctx context.Context) (string, error) {
	span := s.tracer.MakeSpan(ctx, "IdentifyCountry")
	defer span.Finish()

	// define location by IP address
	var ip string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.tracer.LogError(span, errors.New("coudn't resolve ip address"))
	} else {
		strArr := md.Get("ip")
		if len(strArr) > 0 {
			ip = strArr[0]
		}
	}
	country, err := s.repository.GeoIP.GetCountryISOCode(net.ParseIP(ip))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return country, nil
}

// IsUsernameBusy ...
func (s Service) IsUsernameBusy(ctx context.Context, username string) (bool, error) {
	span := s.tracer.MakeSpan(ctx, "IsUsernameBusy")
	defer span.Finish()

	// check if usernmae is not busy
	usernameInUse, err := s.repository.Users.IsUsernameBusy(ctx, username)
	if err != nil {
		s.tracer.LogError(span, err)
		return false, err
	}
	if usernameInUse {
		return true, nil
	}

	return false, nil
}

// CreateNewAccount creates new account for user
func (s Service) CreateNewAccount(ctx context.Context, acc *account.Account, password string) (id, url, tmpToken string, err error) {
	span := s.tracer.MakeSpan(ctx, "CreateNewAccount")
	defer span.Finish()

	// pass data in context
	s.passContext(&ctx)

	acc.FirstName = strings.TrimSpace(acc.FirstName)
	acc.Lastname = strings.TrimSpace(acc.Lastname)
	acc.Emails[0].Email = strings.TrimSpace(acc.Emails[0].Email)
	acc.Emails[0].Email = strings.ToLower(acc.Emails[0].Email)
	acc.Username = strings.ToLower(acc.Username)
	acc.Username = strings.TrimSpace(acc.Username)

	year, month, day := acc.Birthday.Birthday.Date()
	acc.BirthdayDate = account.Date{
		Day:   day,
		Month: int(month),
		Year:  year,
	}

	err = emptyValidator(acc.FirstName, acc.Lastname, acc.Username)
	if err != nil {
		return "", "", "", err
	}
	err = fromTwoToSixtyFour(acc.FirstName)
	if err != nil {
		return "", "", "", err
	}
	err = fromTwoToHundredTwentyEight(acc.Lastname)
	if err != nil {
		return "", "", "", err
	}

	err = userNameValidator(acc.Username)
	if err != nil {
		return "", "", "", err
	}
	if len(acc.Emails) > 0 {
		err = emailValidator(acc.Emails[0].Email)
		if err != nil {
			return "", "", "", err
		}
	} else {
		return "", "", "", errors.New("Please Enter Email")
	}
	err = validPassword(password)
	if err != nil {
		return "", "", "", err
	}

	// TODO: trim data!
	// TODO: make first letter capital in some fields!

	// check if email is not busy
	inUse, err := s.repository.Users.IsEmailAlreadyInUse(ctx, acc.Emails[0].Email)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	if inUse {
		err = errors.New("this_email_already_in_use") // TODO: how it return as gRPC status?
		return
	}

	// check if usernmae is not busy
	usernameInUse, err := s.repository.Users.IsUsernameBusy(ctx, acc.Username)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	if usernameInUse {
		err = errors.New("this_username_already_in_use") // TODO: how it return as gRPC status?
		return
	}

	// TODO: check phone is not busy yet (in future)

	// define location by IP address
	var ip string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.tracer.LogError(span, errors.New("coudn't resolve ip address"))
	} else {
		strArr := md.Get("ip")
		if len(strArr) > 0 {
			ip = strArr[0]
		}
	}
	country, err := s.repository.GeoIP.GetCountryISOCode(net.ParseIP(ip))
	if err != nil {
		s.tracer.LogError(span, err)
	}
	if country != "" {
		acc.Location = &account.UserLocation{
			Location: location.Location{
				Country: &location.Country{
					ID: country,
				},
			},
		}
	}

	id = acc.GenerateID()
	url = acc.GenerateURL()
	acc.Status = status.UserStatusNotActivated // set not_activated status
	acc.CreatedAt = time.Now()                 // set date of registration
	acc.Emails[0].Primary = true               // set email as primary
	acc.Emails[0].GenerateID()

	// encode password
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	err = s.repository.Users.SaveNewAccount(ctx, acc, string(encryptedPass))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	emptyString := ""

	err = s.repository.arrangoRepo.SaveUser(ctx, &arangorepo.User{
		ID:           acc.GetID(),
		CreatedAt:    time.Now(),
		Firstname:    acc.FirstName,
		Lastname:     acc.Lastname,
		Status:       "ACTIVATED",
		URL:          acc.URL,
		PrimaryEmail: acc.Emails[0].Email,
		Gender: arangorepo.Gender{
			Gender: acc.Gender.Gender,
			Type:   &emptyString,
		},
	})
	if err != nil {
		log.Println("arrangoRepo.SaveUser:", err)
	}

	if acc.GetInvitedByID() != "" {
		s.AddGoldCoinsToWallet(ctx, acc.GetInvitedByID(), 1)
	}

	// generate tmp code for activation
	tmpCode, err := s.repository.Cache.CreateTemporaryCodeForEmailActivation(ctx, id, acc.Emails[0].Email)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// log.Println("activation code:", tmpCode) // TODO: delete later
	// log.Println("user id:", id)              // TODO: delete later

	// send email
	// err = s.mailRPC.SendEmail(
	// 	ctx,
	// 	acc.Emails[0].Email,
	// 	fmt.Sprint("<html><body><a target='_blank' href='https://"+s.Host+"/api/activate/user?token=", tmpCode, "&user_id=", id, "'>Activate</a></body></html>")) // TODO: write template for message
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// }
	// // fmt.Println(fmt.Sprint("https://"+s.Host+"/api/activate/user?token=", tmpCode, "&user_id=", id)) // TODO: delete later

	// emailMessage := fmt.Sprint("<html><body><a target='_blank' href='https://"+s.Host+"/api/activate/user?token=", tmpCode, "&user_id=", id, "'>Activate</a></body></html>")
	emailMessage := s.tpl.GetActivationMessage(fmt.Sprint( /*"https://"+s.Host+"/api/activate/user?token=", tmpCode, "&user_id=", id)*/ tmpCode))
	// log.Println(acc.Emails[0].Email, emailMessage)

	err = s.mq.SendEmail(acc.Emails[0].Email, "Activation", emailMessage)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// generate tmp token for not activated user
	tmpToken, err = s.repository.Cache.CreateTemporaryCodeForNotActivatedUser(id)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	err = s.CreateWalletAccount(ctx, id)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return id, url, tmpToken, nil
}

// Recover ...
func (s Service) Recover(ctx context.Context, email string, remindUsername, resetPassword bool) error {
	span := s.tracer.MakeSpan(ctx, "Recover")
	span.Finish()

	emailMessage := "<html><body>"

	id, username, email, err := s.repository.Users.GetUserIDAndUsernameAndPrimaryEmailByLogin(ctx, email)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return err
	}

	if remindUsername {
		emailMessage += fmt.Sprintf("<p>You username is: %s</p>", username)
	}

	err = emailValidator(email)
	if err != nil {
		return err
	}

	if remindUsername == false && resetPassword == false {
		return errors.New("Invalid Request")
	}

	if resetPassword {
		code, err := s.repository.Cache.CreateTemporaryCodeForRecoveryByEmail(id)
		if err != nil {
			s.tracer.LogError(span, err)
			// internal error
			return err
		}
		emailMessage += fmt.Sprint("<a target='_blank' href='https://"+s.Host+"/registration/PasswordRecovery?token=", code, "&user_id=", id, "'>Recover</a>")
	}

	emailMessage += "</body></html>"
	fmt.Println(emailMessage)

	err = s.mq.SendEmail(email, "Recover", emailMessage)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// err = s.mailRPC.SendEmail(ctx, email, emailMessage)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	// internal error
	// 	return err
	// }

	return nil
}

// GenerateRecoveryCode generates code for recovery
// should it allow change pass for only activated users?
func (s Service) GenerateRecoveryCode(ctx context.Context, login string) error {
	span := s.tracer.MakeSpan(ctx, "GenerateRecoveryCode")
	span.Finish()

	// TODO: validate data!

	id, email, err := s.repository.Users.GetUserIDAndPrimaryEmailByLogin(ctx, login)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return err
	}

	code, err := s.repository.Cache.CreateTemporaryCodeForRecoveryByEmail(id)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return err
	}

	message := fmt.Sprint("<html><body><a target='_blank' href='http://"+s.Host+"/user/PasswordRecovery?token=", code, "&user_id=", id, "'>Recover</a></body></html>")
	log.Println(message)

	// err = s.mailRPC.SendEmail(ctx, email, fmt.Sprint("<html><body><a target='_blank' href='http://"+s.Host+"/user/PasswordRecovery?token=", code, "&user_id=", id, "'>Recover</a></body></html>"))
	err = s.mq.SendEmail(email, "Recover", message)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return err
	}

	return nil
}

// ActivateUser looks for code and return id. Called From HTTP server
func (s Service) ActivateUser(ctx context.Context, code string, userID string) (*account.LoginResponse, error) {
	span := s.tracer.MakeSpan(ctx, "ActivateUser")
	defer span.Finish()

	// check tmp code
	matched, email, err := s.repository.Cache.CheckTemporaryCodeForEmailActivation(ctx, userID, code)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	if !matched {
		return nil, errors.New("wrong_activation_code")
	}

	// change status of user
	err = s.repository.Users.ChangeStatusOfUser(ctx, userID, status.UserStatusActivated)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return nil, err
	}

	// change status of email
	err = s.repository.Users.ActivateEmail(ctx, userID, email)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return nil, err
	}

	res, err := s.repository.Users.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return nil, err
	}

	s.passContext(&ctx)

	result := &account.LoginResponse{}

	result.Token, err = s.authRPC.LoginUser(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	// TODO:
	// SetDateOfActivation

	err = s.repository.Cache.Remove(email)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return &account.LoginResponse{
		ID:        userID,
		URL:       res.URL,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Avatar:    res.Avatar,
		Token:     result.Token,
		Gender:    res.Gender,
	}, nil
}

// ActivateEmail ...
func (s Service) ActivateEmail(ctx context.Context, code string, userID string) error {
	span := s.tracer.MakeSpan(ctx, "ActivateEmail")
	defer span.Finish()

	matched, email, err := s.repository.Cache.CheckTemporaryCodeForEmailActivation(ctx, userID, code)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}
	if !matched {
		return errors.New("wrong_activation_code")
	}

	s.repository.Users.ActivateEmail(ctx, userID, email)

	return nil
}

// RecoverPassword if tmp code is matched changes password.
func (s Service) RecoverPassword(ctx context.Context, code, userID, password string) error {
	span := s.tracer.MakeSpan(ctx, "RecoverPassword")
	defer span.Finish()

	matched, err := s.repository.Cache.CheckTemporaryCodeForRecoveryByEmail(userID, code)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	if !matched {
		return errors.New("wrong_activation_code")
	}

	// encrypting password
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	err = s.repository.Users.ChangePassword(ctx, userID, string(pass))
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	err = s.repository.Cache.Remove(userID)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// CheckToken ...
func (s Service) CheckToken(ctx context.Context) (bool, error) {
	span := s.tracer.MakeSpan(ctx, "CheckToken")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return false, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return false, err
	}

	if userID != "" {
		return true, nil
	}

	return false, nil
}

// Login ...
func (s Service) Login(ctx context.Context, login, password, twoFACode string) (result account.LoginResponse, err error) {
	span := s.tracer.MakeSpan(ctx, "Login")
	defer span.Finish()

	// check user status
	res, err := s.repository.Users.GetCredentialsAndStatus(ctx, login)
	if err != nil {
		return
	}

	is2FAEnabled := res.Is2FAEnabled
	twoFASecret := res.TwoFASecret
	userID := res.ID

	result = account.LoginResponse{
		ID:           res.ID,
		Is2FAEnabled: res.Is2FAEnabled,
		Status:       res.Status,
		Avatar:       res.Avatar,
		TwoFASecret:  res.TwoFASecret,
		URL:          res.URL,
		Password:     res.Password,
		FirstName:    res.FirstName,
		LastName:     res.LastName,
		Gender:       res.Gender,
		Email:        res.Email,
	}

	//trimming spaces
	strings.TrimSpace(login)
	strings.TrimSpace(password)
	strings.TrimSpace(twoFACode)

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))
	if err != nil {
		err = errors.New("wrong_credentials")
		return
	}

	// if not activated
	// generate token in cache repository
	if res.Status == status.UserStatusNotActivated {
		tmpToken, errCache := s.repository.Cache.CreateTemporaryCodeForNotActivatedUser(userID)
		if err != nil {
			s.tracer.LogError(span, errCache)
			// internal error
		}
		result.Token = tmpToken
		fmt.Println("unimplemented: temp Token:", tmpToken)
		return
	}

	if res.Status == status.UserStatusDeactivated {
		err = s.repository.Users.ChangeStatusOfUser(ctx, userID, status.UserStatusActivated)
		res.Status = status.UserStatusActivated
		if err != nil {
			s.tracer.LogError(span, err)
		}
	}

	// if activated
	if is2FAEnabled {
		isMatch := twoFA.CheckCode(twoFASecret, twoFACode)
		if !isMatch {
			err = errors.New("wrong_2fa_code")
			return
		}
	}

	s.passContext(&ctx)

	result.Token, err = s.authRPC.LoginUser(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	// save session
	// it is saved in auth service

	// fmt.Println(stat, userID, password, is2FAEnabled, twoFASecret, url, token)
	return
}

// SignOut ...
func (s Service) SignOut(ctx context.Context) error {
	span := s.tracer.MakeSpan(ctx, "SignOut")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("wrong_token")
	}

	s.passContext(&ctx)

	err := s.authRPC.SignOut(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	return err
}

// SignOutSession ...
func (s Service) SignOutSession(ctx context.Context, sessionID string) error {
	span := s.tracer.MakeSpan(ctx, "SignOutSession")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("wrong_token")
	}

	s.passContext(&ctx)

	err := s.authRPC.SignOutSession(ctx, sessionID)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	return err
}

// SignOutFromAll ...
func (s Service) SignOutFromAll(ctx context.Context) error {
	span := s.tracer.MakeSpan(ctx, "SignOutFromAll")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("wrong_token")
	}

	s.passContext(&ctx)

	err := s.authRPC.SignOutFromAll(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	return err
}

// GetAccount returns account of user.
func (s Service) GetAccount(ctx context.Context) (*account.Account, error) {
	span := s.tracer.MakeSpan(ctx, "GetAccount")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	acc, err := s.repository.Users.GetAccount(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	amount, err := s.authRPC.GetAmountOfSessions(ctx)
	if err != nil {
		return nil, err
	}

	acc.AmountOfSessions = amount

	return acc, nil
}

// ChangeFirstName changes first name of user,
// this action can be done only within 5 days after registration.
func (s Service) ChangeFirstName(ctx context.Context, firstname string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeFirstName")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for firstname not to be empty or over 32 or contain !alphabets
	err = fromTwoToSixtyFour(firstname)
	if err != nil {
		return err
	}
	err = dashAndSpace(firstname)
	if err != nil {
		return err
	}

	dateReg, err := s.repository.Users.GetDateOfRegistration(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		//internal error
		return err
	}

	if time.Since(dateReg) > 5*(24*time.Hour) {
		return errors.New("time_for_this_action_is_passed")
	}

	err = s.repository.Users.ChangeFirstName(ctx, userID, firstname)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangeLastName changes last name of user,
// this action can be done only within 5 days after registration.
func (s Service) ChangeLastName(ctx context.Context, lastname string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeLastName")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for lastname not to be empty or over 32 or contain !alphabets
	err = fromTwoToHundredTwentyEight(lastname)
	if err != nil {
		return err
	}
	err = dashAndSpace(lastname)
	if err != nil {
		return err
	}

	dateReg, err := s.repository.Users.GetDateOfRegistration(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		//internal error
		return err
	}

	if time.Since(dateReg) > 5*(24*time.Hour) {
		return errors.New("time_for_this_action_is_passed")
	}

	err = s.repository.Users.ChangeLastName(ctx, userID, lastname)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangePatronymic changes patronomic and its permission
func (s Service) ChangePatronymic(ctx context.Context, patronymic *string, permission *account.Permission) error {
	span := s.tracer.MakeSpan(ctx, "ChangePatronymic")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for patronymic. only contain alphabets and not be over 120 characters
	err = middlenicknameValidator(patronymic)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangePatronymic(ctx, userID, patronymic, permission)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangeNickname changes nickname and its permission
func (s Service) ChangeNickname(ctx context.Context, nickname *string, permission *account.Permission) error {
	span := s.tracer.MakeSpan(ctx, "ChangeNickname")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for nickname. only contain alphabets and not be over 120 characters
	err = middlenicknameValidator(nickname)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangeNickname(ctx, userID, nickname, permission)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangeMiddleName changes middle name and its permission
func (s Service) ChangeMiddleName(ctx context.Context, middlename *string, permission *account.Permission) error {
	span := s.tracer.MakeSpan(ctx, "ChangeMiddleName")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for nickname. only contain alphabets and not be over 120 characters
	err = middlenicknameValidator(middlename)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangeMiddleName(ctx, userID, middlename, permission)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangeNameOnNativeLanguage changes middle name and its permission
func (s Service) ChangeNameOnNativeLanguage(ctx context.Context, name *string, lang *string, permission *account.Permission) error {
	span := s.tracer.MakeSpan(ctx, "ChangeNameOnNativeLanguage")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	// TODO: check if it is a valid language

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check language and name. lang must be 2 characters long.name !="" !>32 && !alphabets
	err = changeNameOnNativeLanguageValidator(name, lang)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangeNameOnNativeLanguage(ctx, userID, name, lang, permission)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangeBirthday changes birthday of user and its permission,
// user can change date only within 5 days after registration.
func (s Service) ChangeBirthday(ctx context.Context, birthday *time.Time, permission *account.Permission) error {
	span := s.tracer.MakeSpan(ctx, "ChangeBirthday")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for input. not allow below 1960
	err = birthdayValidator(birthday)
	if err != nil {
		return err
	}

	dateReg, err := s.repository.Users.GetDateOfRegistration(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		//internal error
		return err
	}

	if birthday != nil {
		if time.Since(dateReg) > 5*(24*time.Hour) {
			return errors.New("time_for_this_action_is_passed")
		}
	}

	err = s.repository.Users.ChangeBirthday(ctx, userID, birthday, permission)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangeGender changes gender of user
// user can change date only within 5 days after registration.
func (s Service) ChangeGender(ctx context.Context, gender *string, permission *account.Permission) error {
	span := s.tracer.MakeSpan(ctx, "ChangeGender")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	dateReg, err := s.repository.Users.GetDateOfRegistration(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		//internal error
		return err
	}

	if gender != nil {
		if time.Since(dateReg) > 5*(24*time.Hour) {
			return errors.New("time_for_this_action_is_passed")
		}
	}

	err = s.repository.Users.ChangeGender(ctx, userID, gender, permission)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// AddEmail ...
func (s Service) AddEmail(ctx context.Context, email string, permission *account.Permission) (id string, err error) {
	span := s.tracer.MakeSpan(ctx, "AddEmail")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	//check if user input email is valid. is email format
	err = emailValidator(email)
	if err != nil {
		return "", err
	}

	// check if the user already added this email before
	inUse, err := s.repository.Users.IsEmailAdded(ctx, userID, email)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return "", err
	}
	if inUse {
		return "", errors.New("email_already_added")
	}

	// check if this activated email already exits among all users
	alreadyExists, err := s.repository.Users.IsEmailAlreadyInUse(ctx, email)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return "", err
	}
	if alreadyExists {
		return "", errors.New("email_already_in_use")
	}

	// saving email in db
	em := account.Email{
		Email: email,
	}
	if permission != nil {
		em.Permission = *permission
	}

	id = em.GenerateID()
	err = s.repository.Users.AddEmail(ctx, userID, &em)
	if err != nil {
		// internal error
		s.tracer.LogError(span, err)
	}

	// send activation link
	tmpCode, err := s.repository.Cache.CreateTemporaryCodeForEmailActivation(ctx, userID, email)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return
	}

	// // TODO: send to queue
	// // Now it has LEAK of goroutine
	// go func() {
	// 	_ = s.mailRPC.SendEmail(ctx, email, fmt.Sprint("<html><body><a target='_blank' href='https://"+s.Host+"/api/activate/email?code=", tmpCode, "&user_id=", userID, "'>Activate</a></body></html>"))
	// }()

	emailMessage := fmt.Sprint("https://"+s.Host+"/api/activate/email?code=", tmpCode, "&user_id=", userID)
	err = s.mq.SendEmail(email, "Activation", emailMessage)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return
}

// RemoveEmail ...
func (s Service) RemoveEmail(ctx context.Context, emailID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveEmail")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if it is primary
	isPrimary, err := s.repository.Users.IsPrimaryEmail(ctx, userID, emailID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	if isPrimary {
		return errors.New("cant_delete_primary_email")
	}

	err = s.repository.Users.RemoveEmail(ctx, userID, emailID)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangeEmail ...
func (s Service) ChangeEmail(ctx context.Context, emailID string, permission *account.Permission, isPrimary bool) error {
	span := s.tracer.MakeSpan(ctx, "ChangeEmail")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if permission != nil && permission.Type != account.PermissionTypeNone {
		// change permission
		err := s.repository.Users.ChangeEmailPermission(ctx, userID, emailID, permission)
		if err != nil {
			s.tracer.LogError(span, err)
			// internal_error
			return err
		}
	}

	// make email primary if isPrimary == true
	if isPrimary {
		// check if it is activated
		isActivated, err := s.repository.Users.IsEmailActivated(ctx, userID, emailID)
		if err != nil {
			s.tracer.LogError(span, err)
			// internal_error
			return err
		}
		if !isActivated {
			return errors.New("email_is_not_activated")
		}

		err = s.repository.Users.MakeEmailPrimary(ctx, userID, emailID)
		if err != nil {
			s.tracer.LogError(span, err)
			// internal_error
			return err
		}
	}

	return nil
}

// AddPhone ...
func (s Service) AddPhone(ctx context.Context, countryCode *account.CountryCode, number string, permission *account.Permission) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddPhone")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// check if user didn't add it before
	inUse, err := s.repository.Users.IsPhoneAdded(ctx, userID, &account.Phone{
		Number:      number,
		CountryCode: *countryCode,
	})
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return "", err
	}
	if inUse {
		return "", errors.New("phone_already_added")
	}

	// check if phone is not added before by someone else
	inUse, err = s.repository.Users.IsPhoneAlreadyInUse(ctx, &account.Phone{
		Number:      number,
		CountryCode: *countryCode,
	})
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
		return "", err
	}
	if inUse {
		return "", errors.New("phone_already_in_use")
	}

	// define country code and country id by country code id
	cd, countryID, err := s.infoRPC.GetCountryIDAndCountryCode(ctx, int32(countryCode.ID))
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return "", err
	}

	if countryID == "" {
		return "", errors.New("bad_country_code_id")
	}

	countryCode.Code = cd
	countryCode.CountryID = countryID

	// checks given numbers for it's countries format.
	err = phoneValidator(countryCode.Code, number, countryCode.CountryID)
	if err != nil {
		return "", err
	}

	// saving in db
	phone := account.Phone{
		Number:      number,
		CountryCode: *countryCode,
	}
	id := phone.GenerateID()

	if permission != nil {
		phone.Permission = *permission
	}

	err = s.repository.Users.AddPhone(
		ctx,
		userID,
		&phone,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// TODO: send SMS for activation

	return id, nil
}

// TODO: func (s Service) ActivatePhone(ctx context.Context, phoneID string, tmpCode string) error{}

// RemovePhone ...
func (s Service) RemovePhone(ctx context.Context, phoneID string) error {
	span := s.tracer.MakeSpan(ctx, "RemovePhone")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if it is primary
	isPrimary, err := s.repository.Users.IsPrimaryPhone(ctx, userID, phoneID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	if isPrimary {
		return errors.New("cant_delete_primary_phone")
	}

	err = s.repository.Users.RemovePhone(ctx, userID, phoneID)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// ChangePhone ...
func (s Service) ChangePhone(ctx context.Context, phoneID string, permission *account.Permission, isPrimary bool) error {
	span := s.tracer.MakeSpan(ctx, "ChangePhone")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// change permission
	if permission != nil && permission.Type != account.PermissionTypeNone {
		err := s.repository.Users.ChangePhonePermission(ctx, userID, phoneID, permission)
		if err != nil {
			s.tracer.LogError(span, err)
			// internal_error
			return err
		}
	}

	// make phone primary if isPrimary == true
	if isPrimary {
		// TODO: uncheck when we will have SMS gateway
		// // check if it is activated
		// isActivated, err := s.repository.Users.IsPhoneActivated(ctx, userID, phoneID)
		// if err != nil {
		// 	s.tracer.LogError(span, err)
		// 	// internal_error
		// 	return err
		// }
		// if !isActivated {
		// 	return errors.New("phone_is_not_activated")
		// }

		err = s.repository.Users.MakePhonePrimary(ctx, userID, phoneID)
		if err != nil {
			s.tracer.LogError(span, err)
			// internal_error
			return err
		}
	}

	return nil
}

// AddMyAddress ...
func (s Service) AddMyAddress(ctx context.Context, address *account.MyAddress) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddMyAddress")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	//check for address inputs not to be over 32 characters or empty
	err = addressValidator(address)
	if err != nil {
		return "", err
	}

	// TODO: check if it is the first address then make it primary

	// if city have id then save get info
	if address.Location.City.ID != 0 {
		cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, address.Location.City.ID, nil)
		if errInfo != nil {
			s.tracer.LogError(span, errInfo)
			// internal_error
			return "", errInfo
		}

		address.Location.City.Name = cityName
		address.Location.City.Subdivision = subdivision
		address.Location.Country.ID = countryID
	}

	id := address.GenerateID()

	err = s.repository.Users.AddMyAddress(ctx, userID, address)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return "", err
	}
	return id, nil
}

// RemoveMyAddress ...
func (s Service) RemoveMyAddress(ctx context.Context, addressID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveMyAddress")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// it's possible to remove primary address

	err = s.repository.Users.RemoveMyAddress(ctx, userID, addressID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	return nil
}

// ChangeMyAddress ...
func (s Service) ChangeMyAddress(ctx context.Context, address *account.MyAddress) error {
	span := s.tracer.MakeSpan(ctx, "ChangeMyAddress")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for address inputs not to be over 32 characters or empty
	err = addressValidator(address)
	if err != nil {
		return err
	}

	// get info about city
	if address.Location.City != nil && address.Location.City.ID != 0 {
		cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, address.Location.City.ID, nil)
		if errInfo != nil {
			s.tracer.LogError(span, errInfo)
			// internal_error
			return errInfo
		}

		address.Location.City.Name = cityName
		address.Location.City.Subdivision = subdivision

		if address.Location.Country == nil {
			address.Location.Country = &location.Country{}
		}

		address.Location.Country.ID = countryID
	}

	err = s.repository.Users.ChangeMyAddress(ctx, userID, address)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}
	return nil
}

// AddOtherAddress ....
func (s Service) AddOtherAddress(ctx context.Context, address *account.OtherAddress) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddOtherAddress")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// if city have id then save get info
	if address.Location.City.ID != 0 {
		cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, address.Location.City.ID, nil)
		if errInfo != nil {
			s.tracer.LogError(span, errInfo)
			// internal_error
			return "", errInfo
		}

		address.Location.City.Name = cityName
		address.Location.City.Subdivision = subdivision
		address.Location.Country.ID = countryID
	}

	id := address.GenerateID()

	err = s.repository.Users.AddOtherAddress(ctx, userID, address)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return "", err
	}
	return id, nil
}

// RemoveOtherAddress ...
func (s Service) RemoveOtherAddress(ctx context.Context, addressID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveOtherAddress")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveOtherAddress(ctx, userID, addressID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	return nil
}

// ChangeOtherAddress ...
func (s Service) ChangeOtherAddress(ctx context.Context, address *account.OtherAddress) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOtherAddress")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// get info about city
	if address.Location.City != nil && address.Location.City.ID != 0 {
		cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, address.Location.City.ID, nil)
		if errInfo != nil {
			s.tracer.LogError(span, errInfo)
			// internal_error
			return errInfo
		}

		address.Location.City.Name = cityName
		address.Location.City.Subdivision = subdivision

		if address.Location.Country == nil {
			address.Location.Country = &location.Country{}
		}

		address.Location.Country.ID = countryID
	}

	err = s.repository.Users.ChangeOtherAddress(ctx, userID, address)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}
	return nil
}

// ChangeUILanguage ...
func (s Service) ChangeUILanguage(ctx context.Context, lang string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeUILanguage")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// TODO: check if it's allowed language
	if lang != "en" {
		return errors.New("unsupported_language")
	}

	err = s.repository.Users.ChangeUILanguage(ctx, userID, lang)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	return nil
}

// ChangePrivacy ...
func (s Service) ChangePrivacy(ctx context.Context, priv *account.PrivacyItem, value *account.PermissionType) error {
	span := s.tracer.MakeSpan(ctx, "ChangePrivacy")
	defer span.Finish()

	if priv == nil || value == nil {
		return errors.New("got null")
	}

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.ChangePrivacy(ctx, userID, *priv, *value)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	return nil
}

// ChangePassword ...
func (s Service) ChangePassword(ctx context.Context, oldPass string, newPass string) error {
	span := s.tracer.MakeSpan(ctx, "ChangePassword")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for new password to be valid
	err = validPassword(newPass)
	if err != nil {
		return err
	}

	res, err := s.repository.Users.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return errors.New("wrong_credentials")
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(oldPass))
	if err != nil {
		return errors.New("wrong_credentials")
	}

	// encrypting password
	newPassHashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	err = s.repository.Users.ChangePassword(ctx, userID, string(newPassHashed))
	if err != nil {
		s.tracer.LogError(span, err)
		// internal error
	}

	return nil
}

// Init2FA resturns qr code encoded as base64 and url. Saves 2fa secret in db.
func (s Service) Init2FA(ctx context.Context) (string, string, string, error) {
	span := s.tracer.MakeSpan(ctx, "Init2FA")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", "", "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", "", "", err
	}

	secret := twoFA.CreateSecret()

	url := twoFA.GenerateLink(secret, "") // TODO: retrive login from db instead of ""
	qr, err := twoFA.GenerateQRCode(url)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return "", "", "", err
	}

	err = s.repository.Users.Save2FASecret(ctx, userID, secret)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return "", "", "", err
	}

	key := twoFA.GetKey(secret)

	return qr, url, key, nil
}

// Enable2FA checks if code is valid and enable 2fa for user
func (s Service) Enable2FA(ctx context.Context, code string) error {
	span := s.tracer.MakeSpan(ctx, "Enable2FA")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// get 2fa secret
	is2FAEnabled, secret, err := s.repository.Users.Get2FAInfo(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	if is2FAEnabled {
		return errors.New("2fa_already_enabled")
	}

	isMatch := twoFA.CheckCode(secret, code)
	if !isMatch {
		return errors.New("wrong_2fa_code")
	}

	err = s.repository.Users.Enable2FA(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	return nil
}

// Disable2FA disables 2fa and remove 2fa secret
func (s Service) Disable2FA(ctx context.Context, code string) error {
	span := s.tracer.MakeSpan(ctx, "Disable2FA")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// get 2fa secret
	is2FAEnabled, secret, err := s.repository.Users.Get2FAInfo(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	if !is2FAEnabled {
		return errors.New("2fa_already_disabled")
	}

	isMatch := twoFA.CheckCode(secret, code)
	if !isMatch {
		return errors.New("wrong_2fa_code")
	}

	err = s.repository.Users.Disable2FA(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return err
	}

	return nil
}

// DeactivateAccount ...
func (s Service) DeactivateAccount(ctx context.Context, password string) error {
	span := s.tracer.MakeSpan(ctx, "DeactivateAccount")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	res, err := s.repository.Users.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))
	if err != nil {
		return errors.New("wrong_password")
	}

	s.passContext(&ctx)

	err = s.authRPC.SignOut(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	err = s.repository.Users.ChangeStatusOfUser(ctx, userID, status.UserStatusDeactivated)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// Profile

// GetProfile returns profile of user with applied permissions and privacy settings
func (s Service) GetProfile(ctx context.Context, url string, language string) (*profile.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetProfile")
	defer span.Finish()

	// retrive profile of target
	prof, err := s.repository.Users.GetProfileByURL(ctx, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	err = s.processProfile(ctx, language, prof)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return prof, nil
}

// GetProfileByID ...
func (s Service) GetProfileByID(ctx context.Context, targetUserID string) (*profile.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetProfileByID")
	defer span.Finish()

	// retrive profile of target
	prof, err := s.repository.Users.GetProfileByID(ctx, targetUserID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	lang := s.retriveUILang(ctx)

	err = s.processProfile(ctx, lang, prof)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return prof, nil
}

// GetProfilesByID ...
func (s Service) GetProfilesByID(ctx context.Context, ids []string, lang string) ([]*profile.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetProfilesByID")
	defer span.Finish()

	// retrive profile of target
	profiles, err := s.repository.Users.GetProfilesByID(ctx, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	if lang == "" {
		lang = s.retriveUILang(ctx)
	}

	var wg sync.WaitGroup
	wg.Add(len(profiles))

	guard := make(chan struct{}, 10)

	for i := range profiles {
		guard <- struct{}{}
		go func(n int) {
			defer wg.Done()

			err = s.processProfile(ctx, lang, profiles[n])
			if err != nil {
				s.tracer.LogError(span, err)
				return
			}

			<-guard
		}(i)
	}

	wg.Wait()

	return profiles, nil
}

// GetMapProfilesByID ...
func (s Service) GetMapProfilesByID(ctx context.Context, ids []string, lang string) (map[string]*profile.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetMapProfilesByID")
	defer span.Finish()

	// retrive profile of target
	profiles, err := s.repository.Users.GetProfilesByID(ctx, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	if lang == "" {
		lang = s.retriveUILang(ctx)
	}
	var wg sync.WaitGroup
	wg.Add(len(profiles))

	guard := make(chan struct{}, 10)

	for i := range profiles {
		guard <- struct{}{}
		go func(n int) {
			defer wg.Done()

			err = s.processProfile(ctx, lang, profiles[n])
			if err != nil {
				s.tracer.LogError(span, err)
				return
			}

			<-guard
		}(i)
	}

	wg.Wait()

	mapProfiles := make(map[string]*profile.Profile, len(profiles))

	for _, pr := range profiles {
		mapProfiles[pr.GetID()] = pr
	}

	return mapProfiles, nil
}

// GetMyCompanies ...
func (s Service) GetMyCompanies(ctx context.Context) (interface{}, error) {
	span := s.tracer.MakeSpan(ctx, "GetMyCompanies")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	companyIDs, err := s.networkRPC.GetUserCompanies(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return nil, err
	}

	companies, err := s.companyRPC.GetCompanies(ctx, companyIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		// internal_error
		return nil, err
	}

	return companies, nil
}

// Wallet

// ContactInvitationForWallet ...
func (s Service) ContactInvitationForWallet(ctx context.Context, name string, email string, message string, coins int32) error {
	span := s.tracer.MakeSpan(ctx, "ContactInvitationForWallet")
	defer span.Finish()

	s.passContext(&ctx)

	err := s.mq.SendEmail(email, "Invitation", message)
	if err != nil {
		s.tracer.LogError(span, err)

		return err
	}

	err = s.stuffRPC.ContactInvitationForWallet(ctx, name, email, message, coins)

	if err != nil {
		s.tracer.LogError(span, err)

		return err
	}

	return nil
}

// GetUserByInvitedID ...
func (s Service) GetUserByInvitedID(ctx context.Context, userID string) (int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetUserByInvitedID")
	defer span.Finish()

	s.passContext(&ctx)

	count, err := s.repository.Users.GetUserByInvitedID(ctx, userID)

	if err != nil {
		s.tracer.LogError(span, err)

		return 0, err
	}

	return count, nil
}

// GetAllUsersForAdmin ...
func (s Service) GetAllUsersForAdmin(ctx context.Context, first uint32, after string) (*profile.Users, error) {
	span := s.tracer.MakeSpan(ctx, "GetAllUsersForAdmin")
	defer span.Finish()

	s.passContext(&ctx)

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

	profile, err := s.repository.Users.GetAllUsersForAdmin(ctx, first, uint32(afterNumber))

	if err != nil {
		s.tracer.LogError(span, err)

		return nil, err
	}

	return profile, nil
}

// ChangeUserStatus ...
func (s Service) ChangeUserStatus(ctx context.Context, userID string, status status.UserStatus) error {
	span := s.tracer.MakeSpan(ctx, "ChangeUserStatus")
	defer span.Finish()

	s.passContext(&ctx)

	err := s.repository.Users.ChangeUserStatus(ctx, userID, status)

	if err != nil {
		s.tracer.LogError(span, err)

		return err
	}

	return nil
}

// AddGoldCoinsToWallet ...
func (s Service) AddGoldCoinsToWallet(ctx context.Context, userID string, coins int32) error {
	span := s.tracer.MakeSpan(ctx, "AddGoldCoinsToWallet")
	defer span.Finish()

	s.passContext(&ctx)

	err := s.stuffRPC.AddGoldCoinsToWallet(ctx, userID, coins)

	if err != nil {
		s.tracer.LogError(span, err)

		return err
	}

	return nil
}

// CreateWalletAccount ...
func (s Service) CreateWalletAccount(ctx context.Context, userID string) error {

	span := s.tracer.MakeSpan(ctx, "CreateWalletAccount")
	defer span.Finish()

	s.passContext(&ctx)

	err := s.stuffRPC.CreateWalletAccount(ctx, userID)

	if err != nil {
		s.tracer.LogError(span, err)

		return err
	}

	return nil
}

//@@@  NEW_PORTOFLIO @@@//

// GetUserPortfolioInfo ...
func (s Service) GetUserPortfolioInfo(ctx context.Context, userID string) (*profile.PortfolioInfo, error) {
	span := s.tracer.MakeSpan(ctx, "GetUserPortfolioInfo")
	defer span.Finish()

	res, err := s.repository.Users.GetUserPortfolioInfo(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return res, nil
}

// AddPortfolio ...
func (s Service) AddPortfolio(ctx context.Context, port *profile.Portfolio) (id string, err error) {

	span := s.tracer.MakeSpan(ctx, "AddPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	id = port.GenerateID()

	port.CreatedAt = time.Now()

	err = s.repository.Users.AddPortfolio(ctx, userID, port)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// AddSavedCountToPortfolio ...
func (s Service) AddSavedCountToPortfolio(ctx context.Context, ownerID string, portfolioID string) error {

	span := s.tracer.MakeSpan(ctx, "AddSavedCountToPortfolio")
	defer span.Finish()

	err := s.repository.Users.AddSavedCountToPortfolio(ctx, ownerID, portfolioID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// LikeUserPortfolio ...
func (s Service) LikeUserPortfolio(ctx context.Context, ownerID string, portfolioID string, companyID string) error {

	var isCompany bool = false
	var profileID string

	span := s.tracer.MakeSpan(ctx, "LikeUserPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	profileID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	_ = s.UnLikeUserPortfolio(ctx, ownerID, portfolioID, companyID)

	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	if companyID != "" {
		profileID = companyID
		isCompany = true
	}

	err = s.repository.Users.LikeUserPortfolio(ctx, &profile.PortfolioAction{
		PortfolioID: portfolioID,
		OwnerID:     ownerID,
		ProfileID:   profileID,
		IsCompany:   isCompany,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// UnLikeUserPortfolio ...
func (s Service) UnLikeUserPortfolio(ctx context.Context, ownerID string, portfolioID string, companyID string) error {

	var isCompany bool = false
	var profileID string

	span := s.tracer.MakeSpan(ctx, "UnLikeUserPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	profileID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if companyID != "" {
		profileID = companyID
		isCompany = true
	}

	err = s.repository.Users.UnLikeUserPortfolio(ctx, &profile.PortfolioAction{
		PortfolioID: portfolioID,
		OwnerID:     ownerID,
		ProfileID:   profileID,
		IsCompany:   isCompany,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddViewCountToPortfolio ...
func (s Service) AddViewCountToPortfolio(ctx context.Context, ownerID string, portfolioID string, companyID string) error {

	var isCompany bool = false
	var profileID string

	span := s.tracer.MakeSpan(ctx, "AddViewCountToPortfolio")
	defer span.Finish()
	token := s.retriveToken(ctx)

	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	profileID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if companyID != "" {
		profileID = companyID
		isCompany = true
	}

	err = s.repository.Users.AddViewCountToPortfolio(ctx, &profile.PortfolioAction{
		PortfolioID: portfolioID,
		OwnerID:     ownerID,
		ProfileID:   profileID,
		IsCompany:   isCompany,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetPortfolioComments ...
func (s Service) GetPortfolioComments(ctx context.Context, porfolioID string, first uint32, after string) (*profile.GetPortfolioComments, error) {

	span := s.tracer.MakeSpan(ctx, "GetPortfolioComments")
	defer span.Finish()

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

	result, err := s.repository.Users.GetPortfolioComments(ctx, porfolioID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return result, nil
}

// AddCommentToPortfolio ...
func (s Service) AddCommentToPortfolio(ctx context.Context, comment *profile.PortfolioComment) (string, error) {

	span := s.tracer.MakeSpan(ctx, "AddCommentToPortfolio")
	defer span.Finish()

	if !comment.IsCompany {
		token := s.retriveToken(ctx)

		if token == "" {
			return "", errors.New("token_is_empty")
		}

		s.passContext(&ctx)

		profileID, err := s.authRPC.GetUserID(ctx, token)

		if err != nil {
			s.tracer.LogError(span, err)
			return "", err
		}

		comment.ProfileID = profileID
	}

	id := comment.GenerateCommentID()

	err := s.repository.Users.AddCommentToPortfolio(ctx, &profile.PortfolioComment{
		ID:          comment.ID,
		Comment:     comment.Comment,
		PortfolioID: comment.PortfolioID,
		ProfileID:   comment.ProfileID,
		OwnerID:     comment.OwnerID,
		IsCompany:   comment.IsCompany,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveCommentInPortfolio ...
func (s Service) RemoveCommentInPortfolio(ctx context.Context, portfolioID, commentID, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveCommentInPortfolio")
	defer span.Finish()

	var profileID string
	var err error

	// If Company
	if companyID != "" {
		profileID = companyID

	} else {
		token := s.retriveToken(ctx)
		if token == "" {
			return errors.New("token_is_empty")
		}

		s.passContext(&ctx)

		profileID, err = s.authRPC.GetUserID(ctx, token)

		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	err = s.repository.Users.RemoveCommentInPortfolio(ctx,
		profileID,
		portfolioID,
		commentID,
	)

	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetPortfolioViewCount ...
func (s Service) GetPortfolioViewCount(ctx context.Context, portfolioID string) (int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetPortfolioViewCount")
	defer span.Finish()

	result, err := s.repository.Users.GetPortfolioViewCount(ctx, portfolioID)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return result, nil
}

// GetPortfolioLikes ...
func (s Service) GetPortfolioLikes(ctx context.Context, profileID, portfolioID string) (*profile.PortfolioLikes, error) {
	span := s.tracer.MakeSpan(ctx, "GetPortfolioLikes")
	defer span.Finish()

	result, err := s.repository.Users.GetPortfolioLikes(ctx, profileID, portfolioID)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return result, nil
}

// GetPortfolios ...
func (s Service) GetPortfolios(ctx context.Context, comapnyID string, userID string, first uint32, after string, contentType string) (*profile.Portfolios, error) {
	span := s.tracer.MakeSpan(ctx, "GetPortfolios")
	defer span.Finish()

	token := s.retriveToken(ctx)

	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	profileID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	if comapnyID != "" {
		profileID = comapnyID
	}

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

	log.Printf("AFTER FOR PORT %+v", afterNumber)
	result, err := s.repository.Users.GetPortfolios(ctx, userID, first, uint32(afterNumber), contentType)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	if result != nil {
		for _, r := range result.Portfolios {
			count, err := s.GetPortfolioViewCount(ctx, r.GetID())
			likes, lErr := s.GetPortfolioLikes(ctx, profileID, r.GetID())

			if lErr != nil {
				s.tracer.LogError(span, err)
			}

			if err != nil {
				s.tracer.LogError(span, err)
			}

			r.ViewsCount = count

			if likes != nil {
				r.LikesCount = likes.Likes
				r.HasLiked = likes.HasLiked
			}

		}
	}

	return result, nil
}

// GetPortfolioByID ...
func (s Service) GetPortfolioByID(ctx context.Context, companyID, userID, portfolioID string) (*profile.Portfolio, error) {

	span := s.tracer.MakeSpan(ctx, "GetPortfolioByID")
	defer span.Finish()

	token := s.retriveToken(ctx)

	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	profileID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	if companyID != "" {
		profileID = companyID
	}

	result, err := s.repository.Users.GetPortfolioByID(ctx, userID, portfolioID)

	if err != nil {
		s.tracer.LogError(span, err)
	}

	count, err := s.GetPortfolioViewCount(ctx, result.GetID())
	likes, lErr := s.GetPortfolioLikes(ctx, profileID, result.GetID())

	if lErr != nil {
		s.tracer.LogError(span, err)
	}

	if err != nil {
		s.tracer.LogError(span, err)
	}

	if likes != nil {
		result.LikesCount = likes.Likes
		result.HasLiked = likes.HasLiked
	}

	result.ViewsCount = count

	return result, nil
}

// ChangeOrderFilesInPortfolio ...
func (s Service) ChangeOrderFilesInPortfolio(ctx context.Context, portfolioID, fileID string, position uint32) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOrderFilesInPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.ChangeOrderFilesInPortfolio(ctx, userID, portfolioID, fileID, position)
	if err != nil {
		return err
	}

	return nil
}

// ChangePortfolio ...
func (s Service) ChangePortfolio(ctx context.Context, port *profile.Portfolio) error {
	span := s.tracer.MakeSpan(ctx, "ChangePortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.ChangePortfolio(ctx, userID, port)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemovePortfolio ...
func (s Service) RemovePortfolio(ctx context.Context, portID string) error {
	span := s.tracer.MakeSpan(ctx, "RemovePortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemovePortfolio(ctx, userID, portID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// AddFileInPortfolio add information about file to user's portfolio. Called in file-manager service.
// It doesn't verify who tries to add!
func (s Service) AddFileInPortfolio(ctx context.Context, userID string, portID string, file *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInPortfolio")
	defer span.Finish()

	id := file.GenerateID()

	err := s.repository.Users.AddFileInPortfolio(ctx, userID, portID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveFilesInPortfolio ...
func (s Service) RemoveFilesInPortfolio(ctx context.Context, portID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveFilesInPortfolio(ctx, userID, portID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddLinksInPortfolio ...
func (s Service) AddLinksInPortfolio(ctx context.Context, expID string, links []*profile.Link) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "AddLinksInPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids := make([]string, 0, len(links))

	for i := range links {
		ids = append(ids, links[i].GenerateID())
	}

	err = s.repository.Users.AddLinksInPortfolio(ctx, userID, expID, links)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return ids, nil
}

// ChangeLinkInPortfolio ...
func (s Service) ChangeLinkInPortfolio(ctx context.Context, expID string, linkID, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeLinkInPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.ChangeLinkInPortfolio(ctx, userID, expID, linkID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveLinksInPortfolio ...
func (s Service) RemoveLinksInPortfolio(ctx context.Context, expID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveLinksInPortfolio")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveLinksInPortfolio(ctx, userID, expID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetToolsTechnologies ...
func (s Service) GetToolsTechnologies(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.ToolTechnology, error) {
	span := s.tracer.MakeSpan(ctx, "GetToolTechnology")
	defer span.Finish()

	// TODO: apply permission ?

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

	result, err := s.repository.Users.GetToolTechnology(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return result, nil
}

// AddToolTechnology adds ToolTechnology to user profile
func (s Service) AddToolTechnology(ctx context.Context, tools []*profile.ToolTechnology) (ids []string, err error) {
	span := s.tracer.MakeSpan(ctx, "AddToolTechnology")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	toolsIDs := make([]string, len(tools))

	for i := range tools {
		toolsIDs[i] = tools[i].GenerateID()
	}

	err = s.repository.Users.AddToolTechnology(ctx, userID, tools)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return toolsIDs, nil
}

// ChangeToolTechnology ...
func (s Service) ChangeToolTechnology(ctx context.Context, tools []*profile.ToolTechnology) error {
	span := s.tracer.MakeSpan(ctx, "ChangeToolTechnology")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for empty inputs. check input characters not be over 32. finishdate not be before start
	// err = experienceValidator(exp)
	// if err != nil {
	// 	return err
	// }
	err = s.repository.Users.ChangeToolTechnology(ctx, userID, tools)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemoveToolTechnology ...
func (s Service) RemoveToolTechnology(ctx context.Context, toolsIDs []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveSkills")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveToolTechnology(ctx, userID, toolsIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetExperiences ...
func (s Service) GetExperiences(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Experience, error) {
	span := s.tracer.MakeSpan(ctx, "GetExperiences")
	defer span.Finish()

	// TODO: apply permission ?

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

	result, err := s.repository.Users.GetExperiences(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	for _, exp := range result {
		exp.Location = location.Location{
			City:    &location.City{},
			Country: &location.Country{},
		}

		// if city have id
		if exp.CityID != nil {
			cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, int32(*exp.CityID), nil)
			if errInfo != nil {
				s.tracer.LogError(span, errInfo)
				// internal_error
				// return "", errInfo
				continue
			}
			exp.Location.City.ID = int32(*exp.CityID)
			exp.Location.City.Name = cityName
			exp.Location.City.Subdivision = subdivision
			exp.Location.Country.ID = countryID
		}
	}

	// translate
	if lang != "" {
		for i := range result {
			result[i].Translate(lang)
		}
	}

	return result, nil
}

// AddExperience ...
func (s Service) AddExperience(ctx context.Context, exp *profile.Experience) (id string, err error) {
	span := s.tracer.MakeSpan(ctx, "AddExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	//experience validator.
	//check for empty inputs
	// err = experienceValidator(exp)
	// if err != nil {
	// 	return "", err
	// }

	if exp.CurrentlyWork {
		exp.FinishDate = nil
	} else {
		if exp.FinishDate.IsZero() {
			return "", err
		}
	}

	id = exp.GenerateID()

	for i := range exp.Links {
		exp.Links[i].GenerateID()
	}

	if city := exp.Location.City; city != nil {
		id := uint32(city.ID)
		exp.CityID = &id
	}

	err = s.repository.Users.AddExperience(ctx, userID, exp)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ChangeExperience ...
func (s Service) ChangeExperience(ctx context.Context, exp *profile.Experience, changeIsCurrentlyWorking bool) error {
	span := s.tracer.MakeSpan(ctx, "ChangeExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for empty inputs. check input characters not be over 32. finishdate not be before start
	// err = experienceValidator(exp)
	// if err != nil {
	// 	return err
	// }

	if exp.CurrentlyWork {
		exp.FinishDate = nil
	} else {
		if exp.FinishDate.IsZero() {
			return err
		}
	}

	if city := exp.Location.City; city != nil {
		id := uint32(city.ID)
		exp.CityID = &id
	}

	for i := range exp.Links {
		if exp.Links[i].GetID() == "" {
			exp.Links[i].GenerateID()
		}
	}

	err = s.repository.Users.ChangeExperience(ctx, userID, exp, changeIsCurrentlyWorking)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemoveExperience ...
func (s Service) RemoveExperience(ctx context.Context, expID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveExperience(ctx, userID, expID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// AddLinksInExperience ...
func (s Service) AddLinksInExperience(ctx context.Context, expID string, links []*profile.Link) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "AddLinksInExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids := make([]string, 0, len(links))

	for i := range links {
		ids = append(ids, links[i].GenerateID())
	}

	err = s.repository.Users.AddLinksInExperience(ctx, userID, expID, links)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return ids, nil
}

// AddFileInExperience add information about file to user's experience. Called in file-manager service.
// It doesn't verify who tries to add!
func (s Service) AddFileInExperience(ctx context.Context, userID string, expID string, file *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInExperience")
	defer span.Finish()

	id := file.GenerateID()

	err := s.repository.Users.AddFileInExperience(ctx, userID, expID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveFilesInExperience ...
func (s Service) RemoveFilesInExperience(ctx context.Context, expID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveFilesInExperience(ctx, userID, expID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeLinkInExperience ...
func (s Service) ChangeLinkInExperience(ctx context.Context, expID string, linkID, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeLinkInExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.ChangeLinkInExperience(ctx, userID, expID, linkID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveLinksInExperience ...
func (s Service) RemoveLinksInExperience(ctx context.Context, expID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveLinksInExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveLinksInExperience(ctx, userID, expID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetUploadedFilesInExperience ...
func (s Service) GetUploadedFilesInExperience(ctx context.Context) ([]*profile.File, error) {
	span := s.tracer.MakeSpan(ctx, "GetUploadedFilesInExperience")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	files, err := s.repository.Users.GetUploadedFilesInExperience(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return files, nil
}

// GetEducations ...
func (s Service) GetEducations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Education, error) {
	span := s.tracer.MakeSpan(ctx, "GetEducations")
	defer span.Finish()

	// TODO: apply permission ?

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

	result, err := s.repository.Users.GetEducations(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	for _, edu := range result {
		edu.Location = location.Location{
			City:    &location.City{},
			Country: &location.Country{},
		}

		// if city have id
		if edu.CityID != nil {
			cityName, subdivision, countryID, errInfo := s.infoRPC.GetCityInformationByID(ctx, int32(*edu.CityID), nil)
			if errInfo != nil {
				s.tracer.LogError(span, errInfo)
				// internal_error
				// return "", errInfo
				continue
			}
			edu.Location.City.ID = int32(*edu.CityID)
			edu.Location.City.Name = cityName
			edu.Location.City.Subdivision = subdivision
			edu.Location.Country.ID = countryID
		}
	}

	// translate
	if lang != "" {
		for i := range result {
			result[i].Translate(lang)
		}
	}

	return result, nil
}

// AddEducation ...
func (s Service) AddEducation(ctx context.Context, edu *profile.Education) (id string, err error) {
	span := s.tracer.MakeSpan(ctx, "AddEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// Education Validator.
	// Check for empty inputs
	err = educationValidator(edu)
	if err != nil {
		return "", err
	}

	if !edu.IsCurrentlyStudy {
		if edu.FinishDate.IsZero() {
			return "", err
		}
	}

	id = edu.GenerateID()

	if city := edu.Location.City; city != nil {
		id := uint32(city.ID)
		edu.CityID = &id
	}

	for i := range edu.Links {
		edu.Links[i].GenerateID()
	}

	if edu.IsCurrentlyStudy {
		edu.FinishDate = time.Time{}
	}

	err = s.repository.Users.AddEducation(ctx, userID, edu)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ChangeEducation ...
func (s Service) ChangeEducation(ctx context.Context, edu *profile.Education) error {
	span := s.tracer.MakeSpan(ctx, "ChangeEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//Education Validator.
	//Check for empty inputs
	// err = educationValidator(edu)
	// if err != nil {
	// 	return err
	// }

	if edu.IsCurrentlyStudy {
		edu.FinishDate = time.Time{}
	} else {
		if edu.FinishDate.IsZero() {
			return err
		}
	}

	if city := edu.Location.City; city != nil {
		id := uint32(city.ID)
		edu.CityID = &id
	}

	for i := range edu.Links {
		if edu.Links[i].GetID() == "" {
			edu.Links[i].GenerateID()
		}
	}

	err = s.repository.Users.ChangeEducation(ctx, userID, edu)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemoveEducation ...
func (s Service) RemoveEducation(ctx context.Context, eduID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveEducation(ctx, userID, eduID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// AddLinksInEducation ...
func (s Service) AddLinksInEducation(ctx context.Context, eduID string, links []*profile.Link) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "AddLinksInEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids := make([]string, 0, len(links))

	for i := range links {
		ids = append(ids, links[i].GenerateID())
	}

	err = s.repository.Users.AddLinksInEducation(ctx, userID, eduID, links)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return ids, nil
}

// AddFileInEducation add information about file to user's education. Called in file-manager service.
func (s Service) AddFileInEducation(ctx context.Context, userID string, eduID string, file *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInEducation")
	defer span.Finish()

	id := file.GenerateID()

	err := s.repository.Users.AddFileInEducation(ctx, userID, eduID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveFilesInEducation ...
func (s Service) RemoveFilesInEducation(ctx context.Context, eduID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveFilesInEducation(ctx, userID, eduID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeLinkInEducation ...
func (s Service) ChangeLinkInEducation(ctx context.Context, expID string, linkID, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeLinkInEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.ChangeLinkInEducation(ctx, userID, expID, linkID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveLinksInEducation ...
func (s Service) RemoveLinksInEducation(ctx context.Context, eduID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveLinksInEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveLinksInEducation(ctx, userID, eduID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetUploadedFilesInEducation ...
func (s Service) GetUploadedFilesInEducation(ctx context.Context) ([]*profile.File, error) {
	span := s.tracer.MakeSpan(ctx, "GetUploadedFilesInEducation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	files, err := s.repository.Users.GetUploadedFilesInEducation(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return files, nil
}

// GetSkills ...
func (s Service) GetSkills(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Skill, error) {
	span := s.tracer.MakeSpan(ctx, "GetSkills")
	defer span.Finish()

	// TODO: apply permission ?

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

	result, err := s.repository.Users.GetSkills(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// translate
	if lang != "" {
		for i := range result {
			result[i].Translate(lang)
		}
	}

	return result, nil
}

// GetEndorsements ...
func (s Service) GetEndorsements(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetEndorsements")
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

	ids, err := s.repository.Users.GetEndorsements(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// getting profiles with applyed permissions
	profiles, err := s.GetProfilesByID(ctx, ids, lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return profiles, nil
}

// AddSkills ...
func (s Service) AddSkills(ctx context.Context, skills []*profile.Skill) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "AddSkills")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	//check for skills not to exceed certain number
	err = skillsValidator(skills)
	if err != nil {
		return nil, err
	}

	lastPosition, err := s.repository.Users.GetPositionOfLastSkill(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	skillsIDs := make([]string, len(skills))

	for i := range skills {
		skillsIDs[i] = skills[i].GenerateID()
		skills[i].Position = lastPosition + 1 + uint32(i)
	}

	err = s.repository.Users.AddSkills(ctx, userID, skills)
	if err != nil {
		return nil, err
	}

	return skillsIDs, nil
}

// ChangeOrderOfSkill ...
func (s Service) ChangeOrderOfSkill(ctx context.Context, skillID string, position uint32) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOrderOfSkill")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.ChangeOrderOfSkill(ctx, userID, skillID, position)
	if err != nil {
		return err
	}

	return nil
}

// RemoveSkills ...
func (s Service) RemoveSkills(ctx context.Context, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveSkills")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveSkills(ctx, userID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// VerifySkill ...
func (s Service) VerifySkill(ctx context.Context, targetID string, skillID string) error {
	span := s.tracer.MakeSpan(ctx, "VerifySkill")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if user tries verify his own skill
	if userID == targetID {
		return errors.New("not_allowed")
	}

	// check if it wasn't verified before
	isVerified, err := s.repository.Users.IsSkillVerified(ctx, userID, targetID, skillID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isVerified {
		return errors.New("already_verified")
	}

	err = s.repository.Users.VerifySkill(ctx, userID, targetID, skillID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// send notification
	n := &notmes.NewEndorsement{
		UserSenderID:  userID,
		EndorsementID: targetID,
		SkillID:       skillID,
	}
	n.GenerateID()
	err = s.mq.SendNewEndorsement(targetID, n)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// UnverifySkill ...
func (s Service) UnverifySkill(ctx context.Context, targetID string, skillID string) error {
	span := s.tracer.MakeSpan(ctx, "UnverifySkill")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if user tries unverify his own skill
	if userID == targetID {
		return errors.New("not_allowed")
	}

	// // check if it wasn't verified before
	// isVerified, err := s.repository.Users.IsSkillVerified(ctx, userID, targetID, skillID)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }
	// if isVerified {
	// 	return errors.New("already_unverified")
	// }

	err = s.repository.Users.UnverifySkill(ctx, userID, targetID, skillID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetInterests ...
func (s Service) GetInterests(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Interest, error) {
	span := s.tracer.MakeSpan(ctx, "GetInterests")
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

	interests, err := s.repository.Users.GetInterests(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// translate
	if lang != "" {
		for i := range interests {
			interests[i].Translate(lang)
		}
	}

	return interests, nil
}

// AddInterest ...
func (s Service) AddInterest(ctx context.Context, interest *profile.Interest) (id string, err error) {
	span := s.tracer.MakeSpan(ctx, "AddInterest")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	//check if interest is empty or exceeds certain character number
	strings.TrimSpace(interest.Interest)
	strings.TrimSpace(*interest.Description)
	err = emptyValidator(interest.Interest)
	if err != nil {
		return "", err
	}
	err = length64Validator(interest.Interest)
	if err != nil {
		return "", err
	}
	if interest.Description != nil {
		err = length128Validator(*interest.Description)
		if err != nil {
			return "", err
		}
	}
	err = pointerLettersNumbersSpecialValidator(interest.Description)
	if err != nil {
		return "", err
	}
	err = lettersNumbersSpecialValidator(interest.Interest)
	if err != nil {
		return "", err
	}

	id = interest.GenerateID()

	err = s.repository.Users.AddInterest(ctx, userID, interest)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ChangeInterest ...
func (s Service) ChangeInterest(ctx context.Context, interest *profile.Interest) error {
	span := s.tracer.MakeSpan(ctx, "ChangeInterest")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	strings.TrimSpace(interest.Interest)
	strings.TrimSpace(*interest.Description)

	//check if interest is empty or exceeds certain character number
	err = length64Validator(interest.Interest)
	if err != nil {
		return err
	}
	err = length128Validator(*interest.Description)
	if err != nil {
		return err
	}
	err = pointerLettersNumbersSpecialValidator(interest.Description)
	if err != nil {
		return err
	}
	err = lettersNumbersSpecialValidator(interest.Interest)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangeInterest(ctx, userID, interest)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemoveInterest ...
func (s Service) RemoveInterest(ctx context.Context, interestID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveInterest")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveInterest(ctx, userID, interestID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// ChangeImageInterest ...
// called from file-manager
func (s Service) ChangeImageInterest(ctx context.Context, userID string, interestID string, image *profile.File) (id string, err error) {
	span := s.tracer.MakeSpan(ctx, "ChangeImageInterest")
	defer span.Finish()

	id = image.GenerateID()

	err = s.repository.Users.ChangeImageInterest(ctx, userID, interestID, image)
	if err != nil {
		s.tracer.LogError(span, err)
		return
	}

	return id, nil
}

// RemoveImageInInterest ...
func (s Service) RemoveImageInInterest(ctx context.Context, interestID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveImageInInterest")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveImageInInterest(ctx, userID, interestID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetUnuploadImageInInterest ...
func (s Service) GetUnuploadImageInInterest(ctx context.Context) (*profile.File, error) {
	span := s.tracer.MakeSpan(ctx, "GetUnuploadImageInInterest")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	file, err := s.repository.Users.GetUnuploadImageInInterest(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return file, nil
}

// GetOriginImageInInterest ...
func (s Service) GetOriginImageInInterest(ctx context.Context, interestID string) (string, error) {
	span := s.tracer.MakeSpan(ctx, "GetOriginImageInInterest")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	image, err := s.repository.Users.GetOriginImageInInterest(ctx, userID, interestID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	return image, nil
}

// ChangeOriginImageInInterest ...
func (s Service) ChangeOriginImageInInterest(ctx context.Context, userID string, interestID string, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOriginImageInInterest")
	defer span.Finish()

	// fmt.Println(userID)

	err := s.repository.Users.ChangeOriginImageInInterest(ctx, userID, interestID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// AddFileInAccomplishment add information about file to user's education. Called in file-manager service.
func (s Service) AddFileInAccomplishment(ctx context.Context, userID string, eduID string, file *profile.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddFileInAccomplishment")
	defer span.Finish()

	id := file.GenerateID()

	err := s.repository.Users.AddFileInAccomplishment(ctx, userID, eduID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// RemoveFilesInAccomplishment ...
func (s Service) RemoveFilesInAccomplishment(ctx context.Context, eduID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFilesInAccomplishment")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveFilesInAccomplishment(ctx, userID, eduID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddLinksInAccomplishment ...
func (s Service) AddLinksInAccomplishment(ctx context.Context, accID string, links []*profile.Link) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "AddLinksInAccomplishment")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids := make([]string, 0, len(links))

	for i := range links {
		ids = append(ids, links[i].GenerateID())
	}

	err = s.repository.Users.AddLinksInAccomplishment(ctx, userID, accID, links)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return ids, nil
}

// RemoveLinksInAccomplishment ...
func (s Service) RemoveLinksInAccomplishment(ctx context.Context, accID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveLinksInAccomplishment")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveLinksInAccomplishment(ctx, userID, accID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetUploadedFilesInAccomplishment ...
func (s Service) GetUploadedFilesInAccomplishment(ctx context.Context) ([]*profile.File, error) {
	span := s.tracer.MakeSpan(ctx, "GetUploadedFilesInAccomplishment")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	files, err := s.repository.Users.GetUploadedFilesInAccomplishment(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return files, nil
}

// GetAccomplishments ...
func (s Service) GetAccomplishments(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Accomplishment, error) {
	span := s.tracer.MakeSpan(ctx, "GetAccomplishments")
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

	accomplishments, err := s.repository.Users.GetAccomplishments(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	// translate
	if lang != "" {
		for i := range accomplishments {
			accomplishments[i].Translate(lang)
		}
	}

	return accomplishments, nil
}

// AddAccomplishment ...
func (s Service) AddAccomplishment(ctx context.Context, accoplishment *profile.Accomplishment) (id string, err error) {
	span := s.tracer.MakeSpan(ctx, "AddAccomplishment")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = accomplishmentValidator(accoplishment)
	if err != nil {
		return "", err
	}

	id = accoplishment.GenerateID()

	for i := range accoplishment.Links {
		accoplishment.Links[i].GenerateID()
	}

	err = s.repository.Users.AddAccomplishment(ctx, userID, accoplishment)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ChangeAccomplishment ...
func (s Service) ChangeAccomplishment(ctx context.Context, accoplishment *profile.Accomplishment) error {
	span := s.tracer.MakeSpan(ctx, "ChangeAccomplishment")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = accomplishmentValidator(accoplishment)
	if err != nil {
		return err
	}

	for i := range accoplishment.Links {
		if accoplishment.Links[i].GetID() == "" {
			accoplishment.Links[i].GenerateID()
		}
	}

	err = s.repository.Users.ChangeAccomplishment(ctx, userID, accoplishment)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemoveAccomplishment ...
func (s Service) RemoveAccomplishment(ctx context.Context, accoplishmentID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveAccomplishment")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveAccomplishment(ctx, userID, accoplishmentID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetReceivedRecommendations ...
func (s Service) GetReceivedRecommendations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Recommendation, error) {
	span := s.tracer.MakeSpan(ctx, "GetReceivedRecommendations")
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

	recommendations, err := s.networkRPC.GetReceivedRecommendationByID(ctx, targetID, int32(first), int32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids := make([]string, 0, len(recommendations)*2)

	for i := range recommendations {
		if recommendations[i].Receiver != nil {
			ids = append(ids, recommendations[i].Receiver.GetID())
		}
		if recommendations[i].Recommendator != nil {
			ids = append(ids, recommendations[i].Recommendator.GetID())
		}
	}

	profileMap, err := s.GetMapProfilesByID(ctx, ids, lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range recommendations {
		if recommendations[i].Receiver != nil {
			recommendations[i].Receiver = profileMap[recommendations[i].Receiver.GetID()]
		}
		if recommendations[i].Recommendator != nil {
			recommendations[i].Recommendator = profileMap[recommendations[i].Recommendator.GetID()]
		}
	}

	return recommendations, nil
}

// GetGivenRecommendations ...
func (s Service) GetGivenRecommendations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Recommendation, error) {
	span := s.tracer.MakeSpan(ctx, "GetGivenRecommendations")
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

	recommendations, err := s.networkRPC.GetGivenRecommendationByID(ctx, targetID, int32(first), int32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids := make([]string, 0, len(recommendations)*2)

	for i := range recommendations {
		if recommendations[i].Receiver != nil {
			ids = append(ids, recommendations[i].Receiver.GetID())
		}
		if recommendations[i].Recommendator != nil {
			ids = append(ids, recommendations[i].Recommendator.GetID())
		}
	}

	profileMap, err := s.GetMapProfilesByID(ctx, ids, lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range recommendations {
		if recommendations[i].Receiver != nil {
			recommendations[i].Receiver = profileMap[recommendations[i].Receiver.GetID()]
		}
		if recommendations[i].Recommendator != nil {
			recommendations[i].Recommendator = profileMap[recommendations[i].Recommendator.GetID()]
		}
	}

	return recommendations, nil
}

// GetReceivedRecommendationRequests ...
func (s Service) GetReceivedRecommendationRequests(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.RecommendationRequest, error) {
	span := s.tracer.MakeSpan(ctx, "GetReceivedRecommendationRequests")
	defer span.Finish()

	s.passContext(&ctx)

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

	recommendations, err := s.networkRPC.GetReceivedRecommendationRequests(ctx, targetID, int32(first), int32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}
	ids := make([]string, 0, len(recommendations))
	for _, recom := range recommendations {
		// if i == 0 {
		// 	ids = append(ids, recom.Requestor.GetID())
		// }
		ids = append(ids, recom.Requestor.GetID())
	}

	// log.Println("ids:", ids)

	profiles, errProf := s.GetMapProfilesByID(ctx, ids, lang)
	if errProf != nil {
		s.tracer.LogError(span, err)
		return nil, errProf
	}

	for i := range recommendations {
		// recommendations[i].Requested = profiles[recommendations[i].Requested.GetID()]
		// recommendations[i].Requestor = profiles[recommendations[i].Requestor.GetID()]
		recommendations[i].Requestor = profiles[recommendations[i].Requestor.GetID()]
		// log.Println("profiles:", profiles)
	}
	// for i := range recommendations {
	// 	profiles, errProf := s.GetProfilesByID(
	// 		ctx,
	// 		[]string{recommendations[i].Requested.GetID(), recommendations[i].Requestor.GetID()},
	// 		lang,
	// 	)
	// 	if err != nil {
	// 		s.tracer.LogError(span, errProf)
	// 		return nil, errProf
	// 	}
	//
	// 	if len(profiles) != 2 {
	// 		err = errors.New("internal_error")
	// 		s.tracer.LogError(span, err)
	// 		return nil, err
	// 	}
	//
	// 	recommendations[i].Requested = profiles[0]
	// 	recommendations[i].Requestor = profiles[1]
	// }

	return recommendations, nil
}

// GetRequestedRecommendationRequests ...
func (s Service) GetRequestedRecommendationRequests(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.RecommendationRequest, error) {
	span := s.tracer.MakeSpan(ctx, "GetRequestedRecommendationRequests")
	defer span.Finish()

	s.passContext(&ctx)

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

	recommendations, err := s.networkRPC.GetRequestedRecommendationRequests(ctx, targetID, int32(first), int32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	ids := make([]string, 0, len(recommendations))
	for _, recom := range recommendations {
		ids = append(ids, recom.Requested.GetID())
	}

	profiles, errProf := s.GetMapProfilesByID(ctx, ids, lang)
	if errProf != nil {
		s.tracer.LogError(span, err)
		return nil, errProf
	}

	for i := range recommendations {
		// recommendations[i].Requestor = profiles[recommendations[i].Requestor.GetID()]
		recommendations[i].Requested = profiles[recommendations[i].Requested.GetID()]
	}

	return recommendations, nil
}

// GetHiddenRecommendations ...
func (s Service) GetHiddenRecommendations(ctx context.Context, targetID string, first uint32, after string, lang string) ([]*profile.Recommendation, error) {
	span := s.tracer.MakeSpan(ctx, "GetHiddenRecommendations")
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

	recommendations, err := s.networkRPC.GetHiddenRecommendationByID(ctx, targetID, int32(first), int32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range recommendations {
		profiles, errProf := s.GetProfilesByID(
			ctx,
			[]string{recommendations[i].Receiver.GetID(), recommendations[i].Recommendator.GetID()},
			lang,
		)
		if err != nil {
			s.tracer.LogError(span, errProf)
			return nil, errProf
		}

		if len(profiles) != 2 {
			err = errors.New("internal_error")
			s.tracer.LogError(span, err)
			return nil, err
		}

		recommendations[i].Receiver = profiles[0]
		recommendations[i].Recommendator = profiles[1]
	}

	return recommendations, nil
}

// GetKnownLanguages ...
func (s Service) GetKnownLanguages(ctx context.Context, targetID string, first uint32, after string) ([]*profile.KnownLanguage, error) {
	span := s.tracer.MakeSpan(ctx, "GetKnownLanguages")
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

	result, err := s.repository.Users.GetKnownLanguages(ctx, targetID, first, uint32(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return result, nil
}

// AddKnownLanguage ...
func (s Service) AddKnownLanguage(ctx context.Context, lang *profile.KnownLanguage) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddKnownLanguage")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	//check for empty language
	err = languageValidator(lang)
	if err != nil {
		return "", err
	}

	id := lang.GenerateID()

	err = s.repository.Users.AddKnownLanguage(ctx, userID, lang)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ChangeKnownLanguage ...
func (s Service) ChangeKnownLanguage(ctx context.Context, lang *profile.KnownLanguage) error {
	span := s.tracer.MakeSpan(ctx, "ChangeKnownLanguage")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for empty language
	err = languageValidator(lang)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangeKnownLanguage(ctx, userID, lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemoveKnownLanguage ...
func (s Service) RemoveKnownLanguage(ctx context.Context, langID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveKnownLanguage")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveKnownLanguage(ctx, userID, langID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// ChangeHeadline ...
func (s Service) ChangeHeadline(ctx context.Context, headline string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeHeadline")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check headline input no to be empty && not > 120 characters
	err = headlineValidator(headline)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangeHeadline(ctx, userID, headline)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// ChangeStory ...
func (s Service) ChangeStory(ctx context.Context, story string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeStory")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check for story not be over 400 characters
	err = storyValidator(story)
	if err != nil {
		return err
	}

	err = s.repository.Users.ChangeStory(ctx, userID, story)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// GetOriginAvatar ...
func (s Service) GetOriginAvatar(ctx context.Context) (string, error) {
	span := s.tracer.MakeSpan(ctx, "GetOriginAvatar")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return "", errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	avatar, err := s.repository.Users.GetOriginAvatar(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	return avatar, nil
}

// ChangeOriginAvatar ...
func (s Service) ChangeOriginAvatar(ctx context.Context, userID string, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeOriginAvatar")
	defer span.Finish()

	// fmt.Println(userID)

	err := s.repository.Users.ChangeOriginAvatar(ctx, userID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// ChangeAvatar ...
func (s Service) ChangeAvatar(ctx context.Context, userID string, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeAvatar")
	defer span.Finish()

	err := s.repository.Users.ChangeAvatar(ctx, userID, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// RemoveAvatar ...
func (s Service) RemoveAvatar(ctx context.Context) error {
	span := s.tracer.MakeSpan(ctx, "RemoveAvatar")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveAvatar(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	return nil
}

// CheckPassword ...
func (s Service) CheckPassword(ctx context.Context, password string) error {
	span := s.tracer.MakeSpan(ctx, "CheckPassword")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	res, err := s.repository.Users.GetCredentialsByUserID(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))
	if err != nil {
		return errors.New("wrong_password")
	}

	return nil
}

// ReportUser ...
func (s Service) ReportUser(ctx context.Context, report *userReport.Report) error {
	span := s.tracer.MakeSpan(ctx, "ReportUser")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	//check if description characters are over 400
	err = reportValidator(report)
	if err != nil {
		return err
	}

	report.SetCreatorID(userID)
	report.CreatedAt = time.Now()

	s.repository.Users.SaveReport(ctx, userID, report)

	return nil
}

// Translations

// SaveUserProfileTranslation ...
func (s Service) SaveUserProfileTranslation(ctx context.Context, lang string, tr *profile.Translation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserProfileTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserProfileTranslation(ctx, userID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveUserExperienceTranslation ...
func (s Service) SaveUserExperienceTranslation(ctx context.Context, expID string, lang string, tr *profile.ExperienceTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserExperienceTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserExperienceTranslation(ctx, userID, expID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveUserEducationTranslation ...
func (s Service) SaveUserEducationTranslation(ctx context.Context, educationID string, lang string, tr *profile.EducationTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserEducationTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserEducationTranslation(ctx, userID, educationID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveUserInterestTranslation ...
func (s Service) SaveUserInterestTranslation(ctx context.Context, interestID string, lang string, tr *profile.InterestTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserInterestTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserInterestTranslation(ctx, userID, interestID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveUserPortfolioTranslation ...
func (s Service) SaveUserPortfolioTranslation(ctx context.Context, portfolioID string, lang string, tr *profile.PortfolioTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserPortfolioTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserPortfolioTranslation(ctx, userID, portfolioID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveUserToolTechnologyTranslation ...
func (s Service) SaveUserToolTechnologyTranslation(ctx context.Context, toolTechID string, lang string, tr *profile.ToolTechnologyTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserToolTechnologyTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserToolTechnologyTranslation(ctx, userID, toolTechID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveUserSkillTranslation ...
func (s Service) SaveUserSkillTranslation(ctx context.Context, skillID string, lang string, tr *profile.SkillTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserSkillTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserSkillTranslation(ctx, userID, skillID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SaveUserAccomplishmentTranslation ...
func (s Service) SaveUserAccomplishmentTranslation(ctx context.Context, accomplishmentID string, lang string, tr *profile.AccomplishmentTranslation) error {
	span := s.tracer.MakeSpan(ctx, "SaveUserAccomplishmentTranslation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.SaveUserAccomplishmentTranslation(ctx, userID, accomplishmentID, lang, tr)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveTransaltion ...
func (s Service) RemoveTransaltion(ctx context.Context, lang string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveTransaltion")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Users.RemoveTransaltion(ctx, userID, lang)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SentEmailInvitation ...
func (s Service) SentEmailInvitation(ctx context.Context, email string, name string, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "SentEmailInvitation")
	defer span.Finish()

	err := emailValidator(email)
	if err != nil {
		return err
	}

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if user not registered already
	inUse, err := s.repository.Users.IsEmailAlreadyInUse(ctx, email)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if inUse {
		return errors.New("user_already_registered")
	}

	// TODO: check if this invitaion on this email has not been send before
	isSend, err := s.repository.Users.IsInvitationSend(ctx, email)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isSend {
		return errors.New("invitation_already_send")
	}

	// sent email
	err = s.mq.SendEmail(email, "Invitation", name+`, you has been invited on https://rightnao.com/`)
	if err != nil {
		return err
	}

	// save invitation in DB
	inv := invitation.Invitation{
		Email:    email,
		Name:     name,
		CratedAt: time.Now(),
	}
	inv.GenerateID()

	if companyID != "" {
		inv.SetCompanyID(companyID)
	} else {
		inv.SetUserID(userID)
	}

	err = s.repository.Users.SaveInvitation(ctx, inv)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetInvitation ...
func (s Service) GetInvitation(ctx context.Context) ([]invitation.Invitation, int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetInvitation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, 0, errors.New("token_is_empty")
	}

	s.passContext(&ctx)

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	invs, amount, err := s.repository.Users.GetInvitation(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return invs, amount, err
}

// GetInvitationForCompany ...
func (s Service) GetInvitationForCompany(ctx context.Context, companyID string) ([]invitation.Invitation, int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetInvitationForCompany")
	defer span.Finish()

	// token := s.retriveToken(ctx)
	// if token == "" {
	// 	return nil, 0, errors.New("token_is_empty")
	// }
	//
	// s.passContext(&ctx)
	//
	// userID, err := s.authRPC.GetUserID(ctx, token)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return nil, 0, err
	// }

	invs, amount, err := s.repository.Users.GetInvitationForCompany(ctx, companyID)
	if err != nil {
		return nil, 0, err
	}

	return invs, amount, err
}

// GetConectionsPrivacy ...
func (s Service) GetConectionsPrivacy(ctx context.Context, userID string) (account.PermissionType, error) {
	span := s.tracer.MakeSpan(ctx, "GetConectionsPrivacy")
	defer span.Finish()

	// token := s.retriveToken(ctx)
	// if token == "" {
	// 	return account.PermissionTypeNone, errors.New("token_is_empty")
	// }
	//
	// s.passContext(&ctx)
	//
	// userID, err := s.authRPC.GetUserID(ctx, token)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return account.PermissionTypeNone, err
	// }

	pr, err := s.repository.Users.GetPrivacyMyConnections(ctx, userID)
	if err != nil {
		return account.PermissionTypeNone, err
	}

	return pr, nil
}

// GetUsersForAdvert ...
func (s Service) GetUsersForAdvert(ctx context.Context, data account.UserForAdvert) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "GetUsersForAdvert")
	defer span.Finish()

	ids, err := s.repository.Users.GetUsersForAdvert(ctx, data)
	if err != nil {
		return nil, err
	}

	return ids, nil

}

// ------------------------------------

// TODO: return error
func (s Service) getInfoFromNetworkForProfile(ctx context.Context, userID string) (isFriend, isBlocked, isBlockedByUser, isFavorite, isFollowing, isFriendRequestSend, isFriendRequestRecieved bool, friendshipID string, mutualConnectionsAmount int32, err error) {
	span := s.tracer.MakeSpan(ctx, "getInfoFromNetworkForProfile")
	defer span.Finish()

	var friend, blocked, blockedByUser, favorite, following bool
	var friendErr, blockedErr, blockedByUserErr, favoriteErr, followingErr, friendRequestErr, friendRequestRecievedErr, friendshiIDdErr, mutualConnectionsAmountErr error

	wg := sync.WaitGroup{}
	wg.Add(9)

	// isFriend?
	go func() {
		defer wg.Done()
		friend, friendErr = s.networkRPC.IsFriend(ctx, userID)
		if friendErr != nil {
			s.tracer.LogError(span, friendErr)
			// return nil, err
		}
	}()

	// is Blocked?
	go func() {
		defer wg.Done()
		blocked, blockedErr = s.networkRPC.IsBlocked(ctx, userID)
		if blockedErr != nil {
			s.tracer.LogError(span, blockedErr)
			// return nil, err
		}
	}()

	// are you blocked?
	go func() {
		defer wg.Done()
		blockedByUser, blockedByUserErr = s.networkRPC.IsBlockedByUser(ctx, userID)
		if blockedByUserErr != nil {
			s.tracer.LogError(span, blockedByUserErr)
			// return nil, err
		}
	}()

	// is Favorite?
	go func() {
		defer wg.Done()
		favorite, favoriteErr = s.networkRPC.IsFavourite(ctx, userID)
		if favoriteErr != nil {
			s.tracer.LogError(span, favoriteErr)
			// return nil, err
		}
	}()

	// is Following?
	go func() {
		defer wg.Done()
		following, followingErr = s.networkRPC.IsFollowing(ctx, userID)
		if followingErr != nil {
			s.tracer.LogError(span, followingErr)
			// return nil, err
		}
	}()

	// is Friend Request Send?
	go func() {
		defer wg.Done()
		isFriendRequestSend, friendRequestErr = s.networkRPC.IsFriendRequestSend(ctx, userID)
		if friendRequestErr != nil {
			s.tracer.LogError(span, friendRequestErr)
			// return nil, err
		}
	}()

	// is Friend Request recivied?
	go func() {
		defer wg.Done()
		isFriendRequestRecieved, _, friendRequestRecievedErr = s.networkRPC.IsFriendRequestRecieved(ctx, userID)

		// log.Printf("isFriendRequestRecieved: %v\nfriendshipID: %v\nfriendRequestRecievedErr: %v\n", isFriendRequestRecieved, friendshipID, friendRequestRecievedErr)

		if friendRequestRecievedErr != nil {
			s.tracer.LogError(span, friendRequestRecievedErr)
			// return nil, err
		}
	}()

	// is Friend Request recivied?
	go func() {
		defer wg.Done()
		friendshipID, friendshiIDdErr = s.networkRPC.GetFriendshipID(ctx, userID)

		// log.Printf("isFriendRequestRecieved: %v\nfriendshipID: %v\nfriendRequestRecievedErr: %v\n", isFriendRequestRecieved, friendshipID, friendRequestRecievedErr)

		if friendshiIDdErr != nil {
			s.tracer.LogError(span, friendRequestRecievedErr)
			// return nil, err
		}
	}()

	// is Friend Request recivied?
	go func() {
		defer wg.Done()
		mutualConnectionsAmount, mutualConnectionsAmountErr = s.networkRPC.GetAmountOfMutualConnections(ctx, userID)

		// log.Printf("isFriendRequestRecieved: %v\nfriendshipID: %v\nfriendRequestRecievedErr: %v\n", isFriendRequestRecieved, friendshipID, friendRequestRecievedErr)

		if mutualConnectionsAmountErr != nil {
			s.tracer.LogError(span, mutualConnectionsAmountErr)
			// return nil, err
		}
	}()

	wg.Wait()
	return friend, blocked, blockedByUser, favorite, following, isFriendRequestSend, isFriendRequestRecieved, friendshipID, mutualConnectionsAmount, nil
}

func (s Service) getInfoFromNetworkForProfileForCompany(ctx context.Context, userID string, companyID string) (isFriend, isBlocked, isBlockedByUser, isFollowing bool, err error) {
	span := s.tracer.MakeSpan(ctx, "getInfoFromNetworkForProfileForCompany")
	defer span.Finish()

	var friend, blocked, blockedByUser, following bool
	var friendErr, blockedErr, blockedByUserErr, followingErr error

	wg := sync.WaitGroup{}
	wg.Add(4)

	// isFriend?
	go func() {
		defer wg.Done()
		friend, friendErr = s.networkRPC.IsFriend(ctx, userID)
		if friendErr != nil {
			s.tracer.LogError(span, friendErr)
			// return nil, err
		}
	}()

	// is Blocked?
	go func() {
		defer wg.Done()
		blocked, blockedErr = s.networkRPC.IsBlockedForCompany(ctx, userID, companyID)
		if blockedErr != nil {
			s.tracer.LogError(span, blockedErr)
			// return nil, err
		}
	}()

	// are you blocked?
	go func() {
		defer wg.Done()
		blockedByUser, blockedByUserErr = s.networkRPC.IsBlockedByCompany(ctx, userID, companyID)
		if blockedByUserErr != nil {
			s.tracer.LogError(span, blockedByUserErr)
			// return nil, err
		}
	}()

	// is Following?
	go func() {
		defer wg.Done()
		following, followingErr = s.networkRPC.IsFollowingForCompany(ctx, userID, companyID)
		if followingErr != nil {
			s.tracer.LogError(span, followingErr)
			// return nil, err
		}
	}()

	wg.Wait()

	return friend, blocked, blockedByUser, following, nil
}

// ------------------------------------

func (s Service) processProfile(ctx context.Context, language string, prof *profile.Profile) error {
	span := s.tracer.MakeSpan(ctx, "processProfile")
	defer span.Finish()

	var hasToken bool
	token := s.retriveToken(ctx)
	if token != "" {
		hasToken = true
	}

	s.passContext(&ctx)

	var userID string

	// retrive id of requestor
	if hasToken {
		var err error
		userID, err = s.authRPC.GetUserID(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
	}

	// get primary addr, put into location
	for _, addr := range prof.MyAddresses {
		if addr.IsPrimary {
			prof.Location = &account.UserLocation{
				Location: addr.Location,
			}
		}
	}

	var err error
	prof.IsOnline, err = s.chatRPC.IsOnline(ctx, prof.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
	}

	var isFriend, isBlocked, isBlockedByUser, isFavorite, isFollowing, isFriendRequestSend, isFriendRequestRecieved bool
	var mutualConnectionsAmount int32
	var friendshipID string

	if hasToken {
		companyID := s.retriveCompanyID(ctx)
		if companyID != "" {
			// as company
			isFriend, isBlocked, isBlockedByUser, isFollowing, err := s.getInfoFromNetworkForProfileForCompany(ctx, prof.GetID(), companyID)
			if err != nil {
				s.tracer.LogError(span, err)
				return err
			}
			if isBlockedByUser {
				return errors.New("you_are_blocked")
			}

			prof.IsFriend = isFriend

			prof.IsBlocked = isBlocked
			prof.IsFollow = isFollowing
			// ---
			if userID == prof.GetID() {
				prof.CompletePercent = s.evaluateComplete(ctx, prof)
				prof.IsMe = true
			}
			// ---
		} else {
			// as user
			isFriend, isBlocked, isBlockedByUser, isFavorite, isFollowing, isFriendRequestSend, isFriendRequestRecieved, friendshipID, mutualConnectionsAmount, err = s.getInfoFromNetworkForProfile(ctx, prof.GetID())
			if err != nil {
				s.tracer.LogError(span, err)
				// internal_error
				return err
			}

			if isBlockedByUser {
				return errors.New("you_are_blocked")
			}

			prof.IsFriend = isFriend
			prof.IsBlocked = isBlocked
			prof.IsFavorite = isFavorite
			prof.IsFollow = isFollowing
			prof.IsFriendRequestSend = isFriendRequestSend
			prof.IsFriendRequestRecieved = isFriendRequestRecieved // better call if isMe
			prof.FriendshipID = friendshipID                       // better call if isMe
			prof.IsMember = true
			prof.MutualConectionsAmount = mutualConnectionsAmount

			if userID == prof.GetID() {
				prof.CompletePercent = s.evaluateComplete(ctx, prof)
				prof.IsMe = true
			}
		}
	}

	prof.ApplyPermissions(ctx)
	prof.ApplyPrivacies(ctx)

	if language == "" {
		language = s.retriveUILang(ctx)
	}
	prof.CurrentTranslation = prof.Translate(ctx, language)
	prof.AvailableTranslations = make([]string, 0, len(prof.Translation))
	prof.AvailableTranslations = append(prof.AvailableTranslations, "en") // english is default
	for key := range prof.Translation {
		prof.AvailableTranslations = append(prof.AvailableTranslations, key)
	}

	return nil
}

func (s Service) passContext(ctx *context.Context) {
	// span := s.tracer.MakeSpan(*ctx, "passContext")
	// defer span.Finish()

	// fmt.Printf("Before: %+v\n", *ctx)

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		log.Println("can't pass context")
	}

	// fmt.Printf("After: %+v\n", *ctx)
}

func (s Service) retriveToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
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

func (s Service) evaluateComplete(ctx context.Context, prof *profile.Profile) int8 {
	var percent int8

	exp, edu, skills, langs, interests, tools, err := s.repository.Users.GetInfoAboutCompletionProfile(ctx, prof.GetID())
	if err != nil {
		log.Println("error: evaluateComplete:", err)
		return 0
	}

	if prof.Avatar != "" {
		percent += 20
	}

	if prof.Story != "" {
		percent += 5
	}

	if prof.Headline != "" {
		percent += 10
	}

	if exp {
		percent += 20
	}

	if edu {
		percent += 20
	}

	if skills {
		percent += 10
	}

	if tools {
		percent += 5
	}

	if langs {
		percent += 5
	}

	if interests {
		percent += 5
	}

	return percent
}
