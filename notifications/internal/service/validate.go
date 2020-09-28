package service

//
// import (
// 	"errors"
// 	"strings"
// 	"unicode/utf8"
//
// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
//
// 	"github.com/asaskevich/govalidator"
// )
//
// const (
// 	namePattern    = "^[a-zA-Z\\-'\\s]+(?:[a-zA-Z]+)$"
// 	specailPattern = "[\\[\\]\\{\\}\\|\\<\\>\\(\\)\\'~!@#$%^&*?.,/\\-+]"
// )
//
// // check fields if any empty for registration form
// func registrationFieldsValidator(acc *account.Account, password string) error {
// 	if acc.FirstName == "" || acc.Lastname == "" || acc.Emails[0].Email == "" || password == "" ||
// 		acc.Birthday.Birthday.IsZero() || acc.Gender.Gender == "" {
// 		return usersErrors.AllRequerd
// 	}
//
// 	return nil
// }
//
// // check name latin alphabet
// func nameValidator(name string) error {
// 	if govalidator.StringMatches(name, namePattern) == false {
// 		return usersErrors.ValidName
// 	}
//
// 	return nil
// }
//
// // check email
// func emailValidator(email string) error {
// 	if govalidator.IsEmail(email) == false {
// 		return usersErrors.ValidEmail
// 	}
//
// 	return nil
// }
//
// // check password
// func passwordValidator(password string) error {
// 	// todo return error from users-errors pkg
// 	var errMsg string
// 	if utf8.RuneCountInString(password) < 9 {
// 		errMsg = "Password must be at least 9 characters\n"
// 	}
//
// 	if govalidator.StringMatches(password, "[A-Z]") == false {
// 		errMsg += "Password must be at least one Uppercase Letter (A-Z)\n"
// 	}
//
// 	if govalidator.StringMatches(password, "[0-9]") == false {
// 		errMsg += "Password must be at least one number (0-9)\n"
// 	}
//
// 	if govalidator.StringMatches(password, specailPattern) == false {
// 		errMsg += "Password must be at least one Special characters"
// 	}
//
// 	if strings.TrimSpace(errMsg) != "" {
// 		return errors.New(strings.TrimSpace(errMsg))
// 	}
//
// 	return nil
// }
