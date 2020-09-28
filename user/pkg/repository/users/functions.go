package usersRepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/invitation"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/status"
	userReport "gitlab.lan/Rightnao-site/microservices/user/pkg/user_report"
	// "github.com/mongodb/mongo-go-driver/bson"
)

const (
	usersCollection            = "users"
	portfolioViews             = "portfolio_views"
	portfolioLikes             = "portfolio_likes"
	portfolioComments          = "portfolio_comments"
	sessionsCollection         = "sessions"
	reportsCollection          = "reports"
	emailInvitationsCollection = "email_invitation"
)

// // TestFunctionInRepository ...
// func (r Repository) TestFunctionInRepository(ctx context.Context) error {
// 	// mongo driver
// 	// result := r.collections[usersCollection].FindOne(ctx, bson.D{})
// 	// if result.Err() != nil {
// 	// 	return result.Err()
// 	// }
//
// 	// mgo driver
// 	result := r.collections[usersCollection].Find(bson.M{})
// 	n, err := result.Count()
// 	doc := bson.D{}
// 	result.One(&doc)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(n)
// 	fmt.Println(&doc)
//
// 	return nil
// }

// SaveNewAccount saves new user account in database
func (r Repository) SaveNewAccount(ctx context.Context, acc *account.Account, password string) error {
	m := bson.D{}
	bs, _ := bson.Marshal(acc)
	bson.Unmarshal(bs, &m)

	m = append(m, bson.DocElem{
		Name:  "password",
		Value: password,
	})
	m = append(m, bson.DocElem{
		Name:  "primary_email",
		Value: acc.Emails[0].Email,
	})

	if err := r.collections[usersCollection].Insert(m); err != nil {
		if mgo.IsDup(err) {
			//  already exists
			// TODO: return error
		}
	}

	return nil
}

// IsEmailAlreadyInUse checks if email is used by anyone.
func (r Repository) IsEmailAlreadyInUse(ctx context.Context, email string) (bool, error) {
	// look for it in account doc
	c, err := r.collections[usersCollection].Find(bson.M{
		"primary_email": email,
	}).Count()

	if err != nil {
		return false, err
	}

	if c > 0 {
		return true, nil
	}

	// look for it in array of emails
	count, err := r.collections[usersCollection].Find(bson.M{
		"emails": bson.M{
			"$elemMatch": bson.M{
				"email":     email,
				"activated": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// IsUsernameBusy checks if email is used by anyone.
func (r Repository) IsUsernameBusy(ctx context.Context, username string) (bool, error) {
	// look for it in account doc
	c, err := r.collections[usersCollection].Find(bson.M{
		"username": username,
	}).Count()

	if err != nil {
		return false, err
	}

	if c > 0 {
		return true, nil
	}
	//
	// // look for it in array of emails
	// count, err := r.collections[usersCollection].Find(bson.M{
	// 	"emails": bson.M{
	// 		"$elemMatch": bson.M{
	// 			"email":     email,
	// 			"activated": true,
	// 		},
	// 	},
	// }).Count()
	// if err != nil {
	// 	return false, err
	// }
	//
	// if count > 0 {
	// 	return true, nil
	// }

	return false, nil
}

// GetUserIDAndUsernameAndPrimaryEmailByLogin returns id, username and primary email of user by login
func (r Repository) GetUserIDAndUsernameAndPrimaryEmailByLogin(ctx context.Context, login string) (id string, username, email string, err error) {
	m := struct {
		ID       bson.ObjectId `bson:"_id"`
		Email    string        `bson:"primary_email"`
		Username string        `bson:"username"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{"primary_email": login})

	if err := result.One(&m); err != nil {
		return "", "", "", err
	}

	return m.ID.Hex(), m.Username, m.Email, nil
}

// GetUserIDAndPrimaryEmailByLogin returns id and primary email of user by login
func (r Repository) GetUserIDAndPrimaryEmailByLogin(ctx context.Context, login string) (id string, email string, err error) {
	m := struct {
		ID    bson.ObjectId `bson:"_id"`
		Email string        `bson:"primary_email"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{"primary_email": login})

	if err := result.One(&m); err != nil {
		return "", "", err
	}

	return m.ID.Hex(), m.Email, nil
}

// // GetPrimaryEmailByUserID returns email by id of user
// func (r Repository) GetPrimaryEmailByUserID(ctx context.Context, userID string) (email string, err error) {
// 	if !bson.IsObjectIdHex(userID) {
// 		return "", errors.New("wrong_id")
// 	}
//
// 	m := struct {
// 		Email string `bson:"primary_email"`
// 	}{}
//
// 	result := r.collections[usersCollection].Find(bson.M{
// 		"_id": bson.ObjectIdHex(userID),
// 	}).Sort("primary_email")
//
// 	result.One(&m)
//
// 	return m.Email, nil
// }

// ChangeStatusOfUser changes status of user
func (r Repository) ChangeStatusOfUser(ctx context.Context, userID string, st status.UserStatus) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"status": st,
			},
		},
	)

	return err
}

// ActivateEmail ...
func (r Repository) ActivateEmail(ctx context.Context, userID string, email string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":          bson.ObjectIdHex(userID),
			"emails.email": email,
		},
		bson.M{
			"$set": bson.M{
				"emails.$.activated": true,
			},
		},
	)

	return err
}

// ChangePassword changes password. It Doesn"t encrypt password!
func (r Repository) ChangePassword(ctx context.Context, userID string, password string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		}, bson.M{
			"$set": bson.M{
				"password":             password,
				"last_change_password": time.Now(),
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetCredentialsAndStatus ...
func (r Repository) GetCredentialsAndStatus(ctx context.Context, login string) (response account.LoginResponse, err error) {
	m := struct {
		ID        bson.ObjectId     `bson:"_id"`
		Status    status.UserStatus `bson:"status"`
		URL       string            `bson:"url"`
		Avatar    string            `bson:"avatar"`
		FirstName string            `bson:"first_name"`
		LastName  string            `bson:"last_name"`
		Password  string            `bson:"password"`
		Email     string            `bson:"primary_email"`
		Gender    struct {
			Gender string `bson:"gender"`
		} `bson:"gender"`

		TwoFA struct {
			Secret   []byte `bson:"secret"`
			IsEnable bool   `bson:"enabled"`
		} `bson:"two_fa"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{
		"username": login,
	})

	err = result.One(&m)
	if err != nil {
		return
	}

	return account.LoginResponse{
		ID:           m.ID.Hex(),
		Password:     m.Password,
		Status:       m.Status,
		Is2FAEnabled: m.TwoFA.IsEnable,
		TwoFASecret:  m.TwoFA.Secret,
		URL:          m.URL,
		Avatar:       m.Avatar,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		Gender:       m.Gender.Gender,
		Email:        m.Email,
	}, nil

}

// GetCredentialsByUserID ...
func (r Repository) GetCredentialsByUserID(ctx context.Context, userID string) (res account.LoginResponse, err error) {
	if !bson.IsObjectIdHex(userID) {
		err = errors.New("wrong_id")
		return
	}

	m := struct {
		URL       string `bson:"url"`
		FirstName string `bson:"first_name"`
		LastName  string `bson:"last_name"`
		Password  string `bson:"password"`
		twoFA     struct {
			Secret   []byte `bson:"secret"`
			IsEnable bool   `bson:"enabled"`
		} `bson:"two_fa"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
	})

	err = result.
		Select(bson.M{
			"url":        1,
			"first_name": 1,
			"last_name":  1,
			"password":   1,
			"two_fa":     1,
		}).
		One(&m)
	if err != nil {
		return
	}

	// m.Password, m.twoFA.IsEnable, m.twoFA.Secret, nil

	return account.LoginResponse{
		Password:     m.Password,
		Is2FAEnabled: m.twoFA.IsEnable,
		TwoFASecret:  m.twoFA.Secret,
		URL:          m.URL,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
	}, nil
}

// GetAccount returns account of user
func (r Repository) GetAccount(ctx context.Context, userID string) (*account.Account, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	var acc account.Account
	result := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
	})
	err := result.One(&acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

// ChangeFirstName changes first name of user.
func (r Repository) ChangeFirstName(ctx context.Context, userID string, firstname string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"first_name": firstname,
			},
		},
	)

	return err
}

// ChangeLastName changes first name of user.
func (r Repository) ChangeLastName(ctx context.Context, userID string, lastname string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"last_name": lastname,
			},
		},
	)

	return err
}

// GetDateOfRegistration returns time of registration of user
func (r Repository) GetDateOfRegistration(ctx context.Context, userID string) (time.Time, error) {
	if !bson.IsObjectIdHex(userID) {
		return time.Time{}, errors.New("wrong_id")
	}

	m := struct {
		CreatedAt time.Time `bson:"created_at"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
	})

	err := result.One(&m)
	if err != nil {
		return time.Time{}, err
	}

	return m.CreatedAt, nil
}

// GetDateOfActivation returns time of activation of user
func (r Repository) GetDateOfActivation(ctx context.Context, userID string) (time.Time, error) {
	if !bson.IsObjectIdHex(userID) {
		return time.Time{}, errors.New("wrong_id")
	}

	m := struct {
		CreatedAt time.Time `bson:"activated_at"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
	})

	err := result.One(&m)
	if err != nil {
		return time.Time{}, err
	}

	return m.CreatedAt, nil
}

// SetDateOfActivation ...
func (r Repository) SetDateOfActivation(ctx context.Context, userID string, date time.Time) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"activated_at": date,
			},
		},
	)

	return err
}

// ChangePatronymic changes patronymic and its permission if permission is not nil.
func (r Repository) ChangePatronymic(ctx context.Context, userID string, patronymic *string, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	changeQuery := bson.M{}
	if patronymic != nil {
		changeQuery["patronymic.name"] = patronymic
	}

	if permission != nil && permission.Type != account.PermissionTypeNone {
		changeQuery["patronymic.permission"] = permission
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": changeQuery,
		},
	)

	return err
}

// ChangeNickname changes nickname and its permission if permission is not nil.
func (r Repository) ChangeNickname(ctx context.Context, userID string, nickname *string, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	changeQuery := bson.M{}
	if nickname != nil {
		changeQuery["nickname.name"] = nickname
	}

	if permission != nil && permission.Type != account.PermissionTypeNone {
		changeQuery["nickname.permission"] = permission
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": changeQuery,
		},
	)

	return err
}

// ChangeMiddleName changes nickname and its permission if permission is not nil.
func (r Repository) ChangeMiddleName(ctx context.Context, userID string, middlename *string, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	changeQuery := bson.M{}
	if middlename != nil {
		changeQuery["middlename.name"] = middlename
	}

	if permission != nil && permission.Type != account.PermissionTypeNone {
		changeQuery["middlename.permission"] = permission
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": changeQuery,
		},
	)

	return err
}

// ChangeNameOnNativeLanguage changes native name, language of name and its permission if permission is not nil.
func (r Repository) ChangeNameOnNativeLanguage(ctx context.Context, userID string, name *string, lang *string, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	changeQuery := bson.M{}

	if name != nil {
		changeQuery["native_name.name"] = name
	}

	if lang != nil {
		changeQuery["native_name.language"] = lang
	}

	if permission != nil && permission.Type != account.PermissionTypeNone {
		changeQuery["native_name.permission"] = permission
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": changeQuery,
		},
	)

	return err
}

// ChangeBirthday ...
func (r Repository) ChangeBirthday(ctx context.Context, userID string, birthday *time.Time, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	changeQuery := bson.M{}
	if birthday != nil {
		changeQuery["birthday.birthday"] = birthday

		year, month, day := birthday.Date()
		changeQuery["birthday_date.day"] = day
		changeQuery["birthday_date.month"] = int(month)
		changeQuery["birthday_date.year"] = year
	}

	if permission != nil && permission.Type != account.PermissionTypeNone {
		changeQuery["birthday.permission"] = permission
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": changeQuery,
		},
	)

	return err
}

// ChangeGender ...
func (r Repository) ChangeGender(ctx context.Context, userID string, gender *string, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	changeQuery := bson.M{}
	if gender != nil {
		changeQuery["gender.gender"] = gender
	}

	if permission != nil && permission.Type != account.PermissionTypeNone {
		changeQuery["gender.permission"] = permission
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": changeQuery,
		},
	)

	return err
}

// IsEmailAdded checks if user already added this email before
func (r Repository) IsEmailAdded(ctx context.Context, userID string, email string) (bool, error) {
	if !bson.IsObjectIdHex(userID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[usersCollection].Find(bson.M{
		"_id":          bson.ObjectIdHex(userID),
		"emails.email": email,
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// AddEmail ...
func (r Repository) AddEmail(ctx context.Context, userID string, email *account.Email) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"emails": email,
			},
		},
	)

	return err
}

// IsPrimaryEmail checks if email with given id is primary for user
func (r Repository) IsPrimaryEmail(ctx context.Context, userID string, emailID string) (bool, error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(emailID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
		"emails": bson.M{
			"$elemMatch": bson.M{
				"id":      bson.ObjectIdHex(emailID),
				"primary": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// RemoveEmail ...
func (r Repository) RemoveEmail(ctx context.Context, userID string, emailID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(emailID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"emails": bson.M{
					"id": bson.ObjectIdHex(emailID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeEmailPermission changes permission of email if permission is not nil
func (r Repository) ChangeEmailPermission(ctx context.Context, userID string, emailID string, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(emailID) {
		return errors.New("wrong_id")
	}

	if permission != nil {
		err := r.collections[usersCollection].Update(
			bson.M{
				"_id":       bson.ObjectIdHex(userID),
				"emails.id": bson.ObjectIdHex(emailID),
			},
			bson.M{
				"$set": bson.M{
					"emails.$.permission": permission,
				},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsEmailActivated checks if email is activated
func (r Repository) IsEmailActivated(ctx context.Context, userID string, emailID string) (bool, error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(emailID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[usersCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
			"emails": bson.M{
				"$elemMatch": bson.M{
					"id":        bson.ObjectIdHex(emailID),
					"activated": true,
				},
			},
		}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// MakeEmailPrimary set for all emails primary false, and makes true only for given email
func (r Repository) MakeEmailPrimary(ctx context.Context, userID string, emailID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(emailID) {
		return errors.New("wrong_id")
	}

	// set primary false for all emails
	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(userID),
			"emails.id": bson.ObjectIdHex(emailID),
		},
		bson.M{
			"$set": bson.M{
				"emails.$[].primary": false,
			},
		},
	)
	if err != nil {
		return err
	}

	// set primary true given email
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(userID),
			"emails.id": bson.ObjectIdHex(emailID),
		},
		bson.M{
			"$set": bson.M{
				"emails.$.primary": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddPhone ...
func (r Repository) AddPhone(ctx context.Context, userID string, phone *account.Phone) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"phones": phone,
			},
		},
	)

	return err
}

// IsPhoneAdded returns true if user already added this phone
func (r Repository) IsPhoneAdded(ctx context.Context, userID string, phone *account.Phone) (bool, error) {
	if !bson.IsObjectIdHex(userID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[usersCollection].Find(bson.M{
		"_id":                    bson.ObjectIdHex(userID),
		"phones.number":          phone.Number,
		"phones.country_code.id": phone.CountryCode.ID,
		// "primary":                true,
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// IsPhoneAlreadyInUse returns true if this phone already added by someone
func (r Repository) IsPhoneAlreadyInUse(ctx context.Context, phone *account.Phone) (bool, error) {
	count, err := r.collections[usersCollection].Find(bson.M{
		"phones.number":          phone.Number,
		"phones.country_code.id": phone.CountryCode.ID,
		"primary":                true,
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// RemovePhone ...
func (r Repository) RemovePhone(ctx context.Context, userID string, phoneID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(phoneID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"phones": bson.M{"id": bson.ObjectIdHex(phoneID)},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// IsPrimaryPhone checks if phone with given id is primary for user
func (r Repository) IsPrimaryPhone(ctx context.Context, userID string, phoneID string) (bool, error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(phoneID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
		"phones": bson.M{
			"$elemMatch": bson.M{
				"id":      bson.ObjectIdHex(phoneID),
				"primary": true,
			},
		},
	}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// ChangePhonePermission changes permission of phone if permission is not nil
func (r Repository) ChangePhonePermission(ctx context.Context, userID string, phoneID string, permission *account.Permission) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(phoneID) {
		return errors.New("wrong_id")
	}

	if permission != nil {
		err := r.collections[usersCollection].Update(
			bson.M{
				"_id":       bson.ObjectIdHex(userID),
				"phones.id": bson.ObjectIdHex(phoneID),
			},
			bson.M{
				"$set": bson.M{
					"phones.$.permission": permission,
				},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsPhoneActivated checks if phone is activated
func (r Repository) IsPhoneActivated(ctx context.Context, userID string, phoneID string) (bool, error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(phoneID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[usersCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
			"phones": bson.M{
				"$elemMatch": bson.M{
					"id":        bson.ObjectIdHex(phoneID),
					"activated": true,
				},
			},
		}).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

// MakePhonePrimary ...
func (r Repository) MakePhonePrimary(ctx context.Context, userID string, phoneID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(phoneID) {
		return errors.New("wrong_id")
	}

	// set primary false for all emails
	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(userID),
			"phones.id": bson.ObjectIdHex(phoneID),
		},
		bson.M{
			"$set": bson.M{
				"phones.$[].primary": false,
			},
		},
	)
	if err != nil {
		return err
	}

	err = r.collections[usersCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(userID),
			"phones.id": bson.ObjectIdHex(phoneID),
		},
		bson.M{
			"$set": bson.M{
				"phones.$.primary": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddMyAddress ...
func (r Repository) AddMyAddress(ctx context.Context, userID string, address *account.MyAddress) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"my_addresses": address,
			},
		},
	)
	return err
}

// RemoveMyAddress ...
func (r Repository) RemoveMyAddress(ctx context.Context, userID string, addressID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(addressID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"my_addresses": bson.M{"id": bson.ObjectIdHex(addressID)},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeMyAddress ...
func (r Repository) ChangeMyAddress(ctx context.Context, userID string, address *account.MyAddress) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(address.GetID()) {
		return errors.New("wrong_id")
	}

	updateQuery := bson.M{}
	if address.Apartment != "" {
		updateQuery["my_addresses.$.apartment"] = address.Apartment
	}
	if address.Firstname != "" {
		updateQuery["my_addresses.$.firstname"] = address.Firstname
	}
	if address.Lastname != "" {
		updateQuery["my_addresses.$.lastname"] = address.Lastname
	}
	if address.Name != "" {
		updateQuery["my_addresses.$.name"] = address.Name
	}
	if address.Street != "" {
		updateQuery["my_addresses.$.street"] = address.Street
	}
	if address.ZIP != "" {
		updateQuery["my_addresses.$.zip"] = address.ZIP
	}
	if address.Location.City != nil &&
		(address.Location.City.ID != 0 ||
			address.Location.City.Name != "" ||
			address.Location.City.Subdivision != "") {
		updateQuery["my_addresses.$.location.city"] = address.Location.City
	}
	if address.Location.Country != nil && address.Location.Country.ID != "" {
		updateQuery["my_addresses.$.location.country"] = address.Location.Country
	}

	if len(updateQuery) > 0 {
		err := r.collections[usersCollection].Update(
			bson.M{
				"_id":             bson.ObjectIdHex(userID),
				"my_addresses.id": bson.ObjectIdHex(address.GetID()),
			},
			bson.M{
				"$set": updateQuery,
			},
		)
		if err != nil {
			return err
		}
	}

	if address.IsPrimary {
		// make every not primary
		err := r.collections[usersCollection].Update(
			bson.M{
				"_id":             bson.ObjectIdHex(userID),
				"my_addresses.id": bson.ObjectIdHex(address.GetID()),
			},
			bson.M{
				"$set": bson.M{
					"my_addresses.$[].primary": false,
				},
			},
		)
		if err != nil {
			return err
		}

		// make primary
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id":             bson.ObjectIdHex(userID),
				"my_addresses.id": bson.ObjectIdHex(address.GetID()),
			},
			bson.M{
				"$set": bson.M{
					"my_addresses.$.primary": true,
				},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddOtherAddress ...
func (r Repository) AddOtherAddress(ctx context.Context, userID string, address *account.OtherAddress) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"other_addresses": address,
			},
		},
	)
	return err
}

// RemoveOtherAddress ...
func (r Repository) RemoveOtherAddress(ctx context.Context, userID string, addressID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(addressID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"other_addresses": bson.M{"id": bson.ObjectIdHex(addressID)},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeOtherAddress ...
func (r Repository) ChangeOtherAddress(ctx context.Context, userID string, address *account.OtherAddress) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(address.GetID()) {
		return errors.New("wrong_id")
	}

	updateQuery := bson.M{}
	if address.Apartment != "" {
		updateQuery["other_addresses.$.apartment"] = address.Apartment
	}
	if address.Firstname != "" {
		updateQuery["other_addresses.$.firstname"] = address.Firstname
	}
	if address.Lastname != "" {
		updateQuery["other_addresses.$.lastname"] = address.Lastname
	}
	if address.Name != "" {
		updateQuery["other_addresses.$.name"] = address.Name
	}
	if address.Street != "" {
		updateQuery["other_addresses.$.street"] = address.Street
	}
	if address.ZIP != "" {
		updateQuery["other_addresses.$.zip"] = address.ZIP
	}

	if address.Location.City != nil &&
		(address.Location.City.ID != 0 ||
			address.Location.City.Name != "" ||
			address.Location.City.Subdivision != "") {
		updateQuery["other_addresses.$.location.city"] = address.Location.City
	}
	if address.Location.Country != nil && address.Location.Country.ID != "" {
		updateQuery["other_addresses.$.location.country"] = address.Location.Country
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"other_addresses.id": bson.ObjectIdHex(address.GetID()),
		},
		bson.M{
			"$set": updateQuery,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeUILanguage changes UI language  for user
func (r Repository) ChangeUILanguage(ctx context.Context, userID string, languageID string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"language_id": languageID,
			},
		})
	if err != nil {
		return err
	}
	return nil
}

// ChangePrivacy ...
func (r Repository) ChangePrivacy(ctx context.Context, userID string, priv account.PrivacyItem, value account.PermissionType) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"privacy." + string(priv): value,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// Save2FASecret ...
func (r Repository) Save2FASecret(ctx context.Context, userID string, secret []byte) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"two_fa": bson.M{
					"secret": secret,
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// Get2FAInfo return 2fa secret and boolean value if it is enabled
func (r Repository) Get2FAInfo(ctx context.Context, userID string) (is2FAEnabled bool, secret []byte, err error) {
	if !bson.IsObjectIdHex(userID) {
		err = errors.New("wrong_id")
		return
	}

	m := struct {
		TwoFA struct {
			Secret    []byte `bson:"secret"`
			IsEnabled bool   `bson:"enabled"`
		} `bson:"two_fa"`
	}{}

	result := r.collections[usersCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
	)

	err = result.One(&m)
	if err != nil {
		return
	}

	return m.TwoFA.IsEnabled, m.TwoFA.Secret, nil
}

// Enable2FA ...
func (r Repository) Enable2FA(ctx context.Context, userID string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"two_fa.enabled": true,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// Disable2FA ...
func (r Repository) Disable2FA(ctx context.Context, userID string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"two_fa.secret":  []byte{},
				"two_fa.enabled": false,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// Profile

// GetProfileByURL ...
func (r Repository) GetProfileByURL(ctx context.Context, url string) (*profile.Profile, error) {

	var prof profile.Profile

	result := r.collections[usersCollection].Find(
		bson.M{
			"url":    url,
			"status": "ACTIVATED",
		},
	)
	err := result.One(&prof)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}

// GetProfileByID ...
func (r Repository) GetProfileByID(ctx context.Context, userID string) (*profile.Profile, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	var prof profile.Profile

	result := r.collections[usersCollection].Find(
		bson.M{
			"_id":    bson.ObjectIdHex(userID),
			"status": "ACTIVATED",
		},
	)
	err := result.One(&prof)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}

// GetProfilesByID ...
func (r Repository) GetProfilesByID(ctx context.Context, ids []string) ([]*profile.Profile, error) {

	idsObject := make([]bson.ObjectId, 0, len(ids))

	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return nil, errors.New("wrong_id")
		}
	}

	var result []*profile.Profile

	err := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.M{
						"$in": idsObject,
					},
				},
			},
		},
	).All(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetPortfolios ...
func (r Repository) GetPortfolios(ctx context.Context, userID string, first uint32, after uint32, contetType string) (*profile.Portfolios, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := new(profile.Portfolios)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":        0,
					"portfolios": 1,
				},
			},

			{
				"$unwind": bson.M{
					"path": "$portfolios",
				},
			},

			{
				"$match": bson.M{
					"portfolios.content_type": contetType,
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"portfolios": bson.M{
						"$push": "$$ROOT.portfolios",
					},
				},
			},
			{
				"$addFields": bson.M{
					"portfolio_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$portfolios"},
							bson.M{"$size": "$portfolios"},
							0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"portfolio_amount": 1,
					"portfolios": bson.M{
						"$slice": []interface{}{
							"$portfolios",
							after,
							first,
						},
					},
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetPortfolioByID ...
func (r Repository) GetPortfolioByID(tx context.Context, userID string, portfolioID string) (*profile.Portfolio, error) {

	if !bson.IsObjectIdHex(userID) && !bson.IsObjectIdHex(portfolioID) {
		return nil, errors.New("wrong_id")
	}

	m := new(profile.Portfolio)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},

			{
				"$unwind": bson.M{
					"path": "$portfolios",
				},
			},
			{
				"$match": bson.M{
					"portfolios.id": bson.ObjectIdHex(portfolioID),
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$portfolios",
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetPortfolioLikes ...
func (r Repository) GetPortfolioLikes(ctx context.Context, profileID, portfolioID string) (*profile.PortfolioLikes, error) {

	if !bson.IsObjectIdHex(profileID) && !bson.IsObjectIdHex(portfolioID) {
		return nil, errors.New("wrong_id")
	}

	m := new(profile.PortfolioLikes)

	res := r.collections[portfolioLikes].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(portfolioID),
				},
			},
			{
				"$addFields": bson.M{
					"likes": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$likes"},
							bson.M{"$size": "$likes"},
							0,
						},
					},
					"has_liked": bson.M{
						"$in": []interface{}{
							bson.ObjectIdHex(profileID), "$likes.profile_id",
						},
					},
				},
			},
			{
				"$project": bson.M{
					"_id":       0,
					"likes":     1,
					"has_liked": 1,
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

//@@ New Portfolio @@//

func (r Repository) AddPortfolio(ctx context.Context, userID string, port *profile.Portfolio) error {

	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	m := []struct {
		File *profile.File `bson:"files"`
	}{}

	port.Files = make([]*profile.File, len(m))
	for i := range m {
		port.Files[i] = m[i].File
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"portfolios": port,
			},
		})

	if err != nil {
		return err
	}

	return nil
}

// AddSavedCountToPortfolio ...
func (r Repository) AddSavedCountToPortfolio(ctx context.Context, ownerID string, portfolioID string) error {

	if !bson.IsObjectIdHex(ownerID) && !bson.IsObjectIdHex(portfolioID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(ownerID),
			"portfolios.id": bson.ObjectIdHex(portfolioID),
		},
		bson.M{
			"$inc": bson.M{
				"portfolios.$.saved_count": 1,
			},
		})
	if err != nil {
		return err
	}

	return nil

}

// AddViewCountToPortfolio ...
func (r Repository) AddViewCountToPortfolio(ctx context.Context, port *profile.PortfolioAction) error {

	if !bson.IsObjectIdHex(port.OwnerID) &&
		!bson.IsObjectIdHex(port.ProfileID) &&
		!bson.IsObjectIdHex(port.PortfolioID) {
		return errors.New("wrong_id")
	}

	type portfolios struct {
		OwnerID   bson.ObjectId `bson:"owner_id"`
		ProfileID bson.ObjectId `bson:"profile_id"`
		IsCompany bool          `bson:"is_company"`
		CreatedAt time.Time     `bson:"created_at"`
	}

	selector := bson.M{"_id": bson.ObjectIdHex(port.PortfolioID)}
	update := bson.M{
		"$push": bson.M{
			"views": &portfolios{
				ProfileID: bson.ObjectIdHex(port.ProfileID),
				OwnerID:   bson.ObjectIdHex(port.OwnerID),
				IsCompany: port.IsCompany,
				CreatedAt: port.CreatedAt,
			},
		},
	}

	// adding portfolio view
	_, err := r.collections[portfolioViews].Upsert(selector, update)

	if err != nil {
		return err
	}

	return nil

}

// AddCommentToPortfolio ...
func (r Repository) AddCommentToPortfolio(ctx context.Context, comment *profile.PortfolioComment) error {

	if !bson.IsObjectIdHex(comment.OwnerID) &&
		!bson.IsObjectIdHex(comment.ProfileID) &&
		!bson.IsObjectIdHex(comment.PortfolioID) {
		return errors.New("wrong_id")
	}

	type portfolios struct {
		ID        bson.ObjectId `bson:"id"`
		Comment   string        `bson:"comment"`
		OwnerID   bson.ObjectId `bson:"owner_id"`
		ProfileID bson.ObjectId `bson:"profile_id"`
		IsCompany bool          `bson:"is_company"`
		CreatedAt time.Time     `bson:"created_at"`
	}

	selector := bson.M{"_id": bson.ObjectIdHex(comment.PortfolioID)}
	update := bson.M{
		"$push": bson.M{
			"comments": &portfolios{
				ID:        comment.ID,
				Comment:   comment.Comment,
				ProfileID: bson.ObjectIdHex(comment.ProfileID),
				OwnerID:   bson.ObjectIdHex(comment.OwnerID),
				IsCompany: comment.IsCompany,
				CreatedAt: comment.CreatedAt,
			},
		},
	}

	// adding portfolio comments
	_, err := r.collections[portfolioComments].Upsert(selector, update)

	if err != nil {
		return err
	}

	return nil

}

// GetAllUsersForAdmin ...
func (r Repository) GetAllUsersForAdmin(ctx context.Context, first uint32, after uint32) (*profile.Users, error) {

	m := new(profile.Users)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$group": bson.M{
					"_id": nil,
					"users": bson.M{
						"$push": bson.M{
							"id":            "$_id",
							"first_name":    "$first_name",
							"last_name":     "$last_name",
							"avatar":        "$avatar",
							"url":           "$url",
							"status":        "$status",
							"created_at":    "$created_at",
							"birthday":      "$birthday",
							"birthday_date": "$birthday_date",
							"location":      "$location",
							"gender":        "$gender.gender",
							"email":         "$primary_email",
							"phones":        "$phones",
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"users_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$users"},
							bson.M{"$size": "$users"},
							0,
						},
					},
				},
			},
			{
				"$facet": bson.M{
					"users": []interface{}{
						bson.M{"$unwind": "$users"},
						bson.M{"$skip": after},
						bson.M{"$limit": first},
					},
				},
			},
			{
				"$project": bson.M{
					"users": "$users.users",
					"users_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$users"},
							bson.M{"$arrayElemAt": []interface{}{"$users.users_amount", 0}}, 0,
						},
					},
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

// ChangeUserStatus ...
func (r Repository) ChangeUserStatus(ctx context.Context, userID string, status status.UserStatus) error {

	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"status": status,
			},
		})

	if err != nil {
		return err
	}

	return nil
}

// GetUserPortfolioInfo ...
func (r Repository) GetUserPortfolioInfo(ctx context.Context, userID string) (*profile.PortfolioInfo, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := new(profile.PortfolioInfo)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id": 0,
					"has_video": bson.M{
						"$cond": bson.M{
							"if": bson.M{"$isArray": "$portfolios"},
							"then": bson.M{
								"$in": []interface{}{
									"video", "$portfolios.content_type",
								},
							},
							"else": false,
						},
					},
					"has_article": bson.M{
						"$cond": bson.M{
							"if": bson.M{"$isArray": "$portfolios"},
							"then": bson.M{
								"$in": []interface{}{
									"article", "$portfolios.content_type",
								},
							},
							"else": false,
						},
					},
					"has_audio": bson.M{
						"$cond": bson.M{
							"if": bson.M{"$isArray": "$portfolios"},
							"then": bson.M{
								"$in": []interface{}{
									"audio", "$portfolios.content_type",
								},
							},
							"else": false,
						},
					},
					"has_photo": bson.M{
						"$cond": bson.M{
							"if": bson.M{"$isArray": "$portfolios"},
							"then": bson.M{
								"$in": []interface{}{
									"photo", "$portfolios.content_type",
								},
							},
							"else": false,
						},
					},
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return nil, err
	}

	// Get Comments amount
	commentAmount, err := r.GetUserPortfolioCommentAmount(ctx, userID)

	if err != nil {
		commentAmount = 0
	}
	m.CommentCount = commentAmount

	// Get view amount
	viewAmount, err := r.GetUserPortfolioViewAmount(ctx, userID)

	if err != nil {
		viewAmount = 0
	}
	m.ViewsCount = viewAmount

	// Get like amount
	likeAmount, err := r.GetUserPortfolioLikeAmount(ctx, userID)

	if err != nil {
		likeAmount = 0
	}
	m.LikesCount = likeAmount

	return m, nil

}

// GetUserPortfolioCommentAmount ...
func (r Repository) GetUserPortfolioCommentAmount(ctx context.Context, userID string) (int32, error) {

	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	m := struct {
		Amount int32 `bson:"comments_amount"`
	}{}

	res := r.collections[portfolioComments].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"comments.owner_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"count": bson.M{
						"$sum": bson.M{
							"$size": "$comments",
						},
					},
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"comments_amount": bson.M{
						"$sum": "$count",
					},
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return 0, err
	}

	return m.Amount, nil

}

// GetUserPortfolioLikeAmount ...
func (r Repository) GetUserPortfolioLikeAmount(ctx context.Context, userID string) (int32, error) {

	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	m := struct {
		Amount int32 `bson:"likes_amount"`
	}{}

	res := r.collections[portfolioLikes].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"likes.owner_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"count": bson.M{
						"$sum": bson.M{
							"$size": "$likes",
						},
					},
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"likes_amount": bson.M{
						"$sum": "$count",
					},
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return 0, err
	}

	return m.Amount, nil

}

// GetUserPortfolioViewAmount ...
func (r Repository) GetUserPortfolioViewAmount(ctx context.Context, userID string) (int32, error) {

	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	m := struct {
		Amount int32 `bson:"views_amount"`
	}{}

	res := r.collections[portfolioViews].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"views.owner_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$unwind": bson.M{
					"path": "$views",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$views",
				},
			},
			{
				"$project": bson.M{
					"year": bson.M{
						"$year": "$created_at",
					},
					"month": bson.M{
						"$month": "$created_at",
					},
					"day": bson.M{
						"$dayOfMonth": "$created_at",
					},
					"profile_id": 1,
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"year":       "$year",
						"month":      "$month",
						"day":        "$day",
						"profile_id": "$profile_id",
					},
					"views_amount": bson.M{
						"$sum": 1,
					},
				},
			},
			{
				"$count": "views_amount",
			},
		})

	err := res.One(&m)

	if err != nil {
		return 0, err
	}

	return m.Amount, nil

}

// GetPortfolioComments ...
func (r Repository) GetPortfolioComments(ctx context.Context, porfolioID string, first uint32, after uint32) (*profile.GetPortfolioComments, error) {

	if !bson.IsObjectIdHex(porfolioID) {
		return nil, errors.New("wrong_id")
	}

	m := new(profile.GetPortfolioComments)

	res := r.collections[portfolioComments].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(porfolioID),
				},
			},
			{
				"$addFields": bson.M{
					"comments_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$comments"},
							bson.M{"$size": "$comments"},
							0,
						},
					},
				},
			},
			{
				"$facet": bson.M{
					"comments": []bson.M{
						bson.M{"$unwind": "$comments"},
						bson.M{
							"$sort": bson.M{
								"created_at": 1,
							},
						},
						bson.M{"$skip": after},
						bson.M{"$limit": first},
					},
				},
			},
			{
				"$project": bson.M{
					"comments": "$comments.comments",
					"comments_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$comments.comments_amount"},
							bson.M{"$arrayElemAt": []interface{}{"$comments.comments_amount", 0}}, 0,
						},
					},
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetUserByInvitedID ...
func (r Repository) GetUserByInvitedID(ctx context.Context, userID string) (int32, error) {

	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	m := struct {
		Count int32 `bson:"count"`
	}{}

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"invited_by": bson.ObjectIdHex(userID),
				},
			},
			{
				"$count": "count",
			},
		})

	err := res.One(&m)

	if err != nil {
		return 0, err
	}

	return m.Count, nil
}

// RemoveCommentInPortfolio ...
func (r Repository) RemoveCommentInPortfolio(ctx context.Context, profileID, portfolioID, commentID string) error {

	if !bson.IsObjectIdHex(profileID) &&
		!bson.IsObjectIdHex(portfolioID) &&
		!bson.IsObjectIdHex(commentID) {
		return errors.New("wrong_id")
	}

	err := r.collections[portfolioComments].Update(
		bson.M{
			"_id": bson.ObjectIdHex(portfolioID),
		},
		bson.M{
			"$pull": bson.M{
				"comments": bson.M{
					"id":         bson.ObjectIdHex(commentID),
					"profile_id": bson.ObjectIdHex(profileID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// LikeUserPortfolio ...
func (r Repository) LikeUserPortfolio(ctx context.Context, port *profile.PortfolioAction) error {

	if !bson.IsObjectIdHex(port.OwnerID) &&
		!bson.IsObjectIdHex(port.ProfileID) &&
		!bson.IsObjectIdHex(port.PortfolioID) {
		return errors.New("wrong_id")
	}

	type portfolios struct {
		OwnerID   bson.ObjectId `bson:"owner_id"`
		ProfileID bson.ObjectId `bson:"profile_id"`
		IsCompany bool          `bson:"is_company"`
		CreatedAt time.Time     `bson:"created_at"`
	}

	selector := bson.M{"_id": bson.ObjectIdHex(port.PortfolioID)}
	update := bson.M{
		"$push": bson.M{
			"likes": &portfolios{
				ProfileID: bson.ObjectIdHex(port.ProfileID),
				OwnerID:   bson.ObjectIdHex(port.OwnerID),
				IsCompany: port.IsCompany,
				CreatedAt: port.CreatedAt,
			},
		},
	}

	// adding portfolio view
	_, err := r.collections[portfolioLikes].Upsert(selector, update)

	if err != nil {
		return err
	}

	return nil

}

// UnLikeUserPortfolio ...
func (r Repository) UnLikeUserPortfolio(ctx context.Context, port *profile.PortfolioAction) error {

	if !bson.IsObjectIdHex(port.OwnerID) &&
		!bson.IsObjectIdHex(port.ProfileID) &&
		!bson.IsObjectIdHex(port.PortfolioID) {
		return errors.New("wrong_id")
	}

	err := r.collections[portfolioLikes].Update(
		bson.M{
			"_id": bson.ObjectIdHex(port.PortfolioID),
		},
		bson.M{
			"$pull": bson.M{
				"likes": bson.M{
					"profile_id": bson.ObjectIdHex(port.ProfileID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil

}

// GetPortfolioViewCount ...
func (r Repository) GetPortfolioViewCount(ctx context.Context, portfolioID string) (int32, error) {

	if !bson.IsObjectIdHex(portfolioID) {
		return 0, errors.New("wrong_id")
	}

	m := struct {
		Sum int32 `bson:"sum"`
	}{}

	res := r.collections[portfolioViews].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(portfolioID),
				},
			},
			{
				"$unwind": bson.M{
					"path": "$views",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$views",
				},
			},
			{
				"$project": bson.M{
					"year": bson.M{
						"$year": "$created_at",
					},
					"month": bson.M{
						"$month": "$created_at",
					},
					"day": bson.M{
						"$dayOfMonth": "$created_at",
					},
					"profile_id": 1,
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"year":       "$year",
						"month":      "$month",
						"day":        "$day",
						"profile_id": "$profile_id",
					},
					"sum": bson.M{
						"$sum": 1,
					},
				},
			},
			{
				"$count": "sum",
			},
		},
	)

	err := res.One(&m)
	if err != nil {
		log.Println("error:", err)
		return 0, err
	}
	return m.Sum, nil

}

// GetPositionOfLastFile returns the position of last file
func (r Repository) GetPositionOfLastFile(ctx context.Context, userID, portfolioID string) (uint32, error) {
	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	m := make([]struct {
		Files *profile.File `bson:"files"`
	},
		0,
	)

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id":           bson.ObjectIdHex(userID),
					"portfolios.id": bson.ObjectIdHex(portfolioID),
				},
			},
			{
				"$project": bson.M{
					"_id":   0,
					"files": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$files",
				},
			},
			{
				"$sort": bson.M{
					"files.position": -1,
				},
			},
			{
				"$limit": 1,
			},
		})

	err := result.All(&m)
	if err != nil {
		return 0, err
	}

	if len(m) > 0 {
		return m[0].Files.Position, nil
	}

	return 0, nil
}

// ChangeOrderFilesInPortfolio ...
// BUG: it incerement all elements instead ignoring condition
func (r Repository) ChangeOrderFilesInPortfolio(ctx context.Context, userID, portfolioID, fileID string, position uint32) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(fileID) {
		return errors.New("wrong_id")
	}

	m := struct {
		Portfolio struct {
			Files []*profile.File `bson:"files"`
		} `bson:"portfolios"`
	}{}

	// getting skill
	err := r.collections[usersCollection].Find(
		bson.M{
			"_id":                 bson.ObjectIdHex(userID),
			"portfolios.id":       bson.ObjectIdHex(portfolioID),
			"portfolios.files.id": bson.ObjectIdHex(fileID),
		},
	).Select(
		bson.M{
			"portfolios.files.$.id": 1,
		},
	).One(&m)

	if err != nil {
		return err
	}

	if len(m.Portfolio.Files) < 1 {
		return errors.New("files not found")
	}

	// set position
	m.Portfolio.Files[0].Position = position

	// removing skill from array
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"portfolios": bson.M{
					"files.id": bson.ObjectIdHex(fileID),
				},
			},
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	// increment position
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
			"portfolios.files.position": bson.M{
				"$gt": position,
			},
		},
		bson.M{
			"$inc": bson.M{
				"portfolios.files.$[].position": 1,
			},
		},
	)
	if err != nil {
		return err
	}

	// inserting skill in certain position
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"portfolios": bson.M{
					"$each": m.Portfolio.Files,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangePortfolio ...
func (r Repository) ChangePortfolio(ctx context.Context, userID string, port *profile.Portfolio) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	bs := make(bson.M)

	if port.Description != "" {
		bs["portfolios.$.description"] = port.Description
	}
	if port.Title != "" {
		bs["portfolios.$.title"] = port.Title
	}

	if len(port.Tools) > 0 {
		bs["portfolios.$.tools"] = port.Tools
	}

	bs["portfolios.$.is_comment_closed"] = port.IsCommentClosed

	if len(bs) == 0 {
		return errors.New("nothing_to_change")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"portfolios.id": port.ID,
		},
		bson.M{
			"$set": bs,
		})
	if err != nil {
		return err
	}

	return nil
}

// RemovePortfolio ...
func (r Repository) RemovePortfolio(ctx context.Context, userID string, portID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(portID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"portfolios": bson.M{
					"id": bson.ObjectIdHex(portID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeLinkInPortfolio ...
func (r Repository) ChangeLinkInPortfolio(ctx context.Context, userID string, portID string, linkID string, url string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(portID) || !bson.IsObjectIdHex(linkID) {
		return errors.New("wrong_id")
	}

	result := struct {
		Index int `bson:"index"`
	}{}

	err := r.collections[usersCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(userID),
			},
		},
		{
			"$project": bson.M{
				"portfolios": 1,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$portfolios",
				"includeArrayIndex":          "portfolios.index",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"portfolios.id": bson.ObjectIdHex(portID),
			},
		},
		{
			"$project": bson.M{
				"index": "$portfolios.index",
			},
		},
	}).One(&result)
	if err != nil {
		return err
	}

	fmt.Println("Index:", result.Index)

	err = r.collections[usersCollection].Update(
		bson.M{
			"_id":                 bson.ObjectIdHex(userID),
			"portfolios.links.id": bson.ObjectIdHex(linkID),
		},
		bson.M{
			"$set": bson.M{
				"portfolios." + strconv.Itoa(result.Index) + ".links.$.url": url,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveLinksInPortfolio ...
func (r Repository) RemoveLinksInPortfolio(ctx context.Context, userID string, portID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(portID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"portfolios.id": bson.ObjectIdHex(portID),
		},
		bson.M{
			"$pull": bson.M{
				"portfolios.$.links": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddLinksInPortfolio ...
func (r Repository) AddLinksInPortfolio(ctx context.Context, userID string, portID string, links []*profile.Link) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(portID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"portfolios.id": bson.ObjectIdHex(portID),
		},
		bson.M{
			"$push": bson.M{
				"portfolios.$.links": bson.M{
					"$each": links,
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddFileInPortfolio ...
func (r Repository) AddFileInPortfolio(ctx context.Context, userID string, portID string, file *profile.File) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}
	if portID != "" && !bson.IsObjectIdHex(portID) {
		return errors.New("wrong_id")
	}

	var err error

	f := struct {
		File *profile.File `bson:",inline"`
		Type string
	}{
		File: file,
		Type: "Portfolio",
	}

	if portID != "" {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id":           bson.ObjectIdHex(userID),
				"portfolios.id": bson.ObjectIdHex(portID),
			},
			bson.M{
				"$push": bson.M{
					"portfolios.$.files": f,
				},
			},
		)
	} else {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(userID),
			},
			bson.M{
				"$push": bson.M{
					"unused_files": f,
				},
			},
		)
	}

	if err != nil {
		return err
	}
	return nil
}

// RemoveFilesInPortfolio ...
func (r Repository) RemoveFilesInPortfolio(ctx context.Context, userID string, portID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(portID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"portfolios.id": bson.ObjectIdHex(portID),
		},
		bson.M{
			"$pull": bson.M{
				"portfolios.$.files": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetToolTechnology ...
func (r Repository) GetToolTechnology(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.ToolTechnology, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]struct {
		ToolTechnology *profile.ToolTechnology `bson:"tools_technologies"`
	},
		0,
	)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":                0,
					"tools_technologies": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$tools_technologies",
				},
			},
			{
				"$sort": bson.M{
					"tools_technologies.created_at": 1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := res.All(&m)
	if err != nil {
		return nil, err
	}

	tools := make([]*profile.ToolTechnology, 0, len(m))
	for i := range m {
		tools = append(tools, m[i].ToolTechnology)
	}

	return tools, nil
}

// AddToolTechnology ...
func (r Repository) AddToolTechnology(ctx context.Context, userID string, tools []*profile.ToolTechnology) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	// adding ToolTechnology
	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"tools_technologies": bson.M{
					"$each": tools,
				},
			},
		})

	if err != nil {
		return err
	}

	return nil
}

// ChangeToolTechnology ...
func (r Repository) ChangeToolTechnology(ctx context.Context, userID string, tools []*profile.ToolTechnology) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	// adding ToolTechnology
	bs := make(bson.M)
	for _, tool := range tools {
		if tool.ToolTechnology != "" {
			bs["tools_technologies.$.tool-technology"] = tool.ToolTechnology
		}
		if tool.Rank != "" {
			bs["tools_technologies.$.rank"] = tool.Rank
		}

		err := r.collections[usersCollection].Update(
			bson.M{
				"_id":                   bson.ObjectIdHex(userID),
				"tools_technologies.id": tool.ID,
			},
			bson.M{
				"$set": bs,
			})

		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveToolTechnology ...
func (r Repository) RemoveToolTechnology(ctx context.Context, userID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"tools_technologies": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetExperiences ...
func (r Repository) GetExperiences(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Experience, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]struct {
		Experience *profile.Experience `bson:"experiences"`
	},
		0,
	)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":         0,
					"experiences": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$experiences",
				},
			},
			{
				"$sort": bson.M{
					"experiences.start_date":  -1,
					"experiences.finish_date": 1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := res.All(&m)
	if err != nil {
		return nil, err
	}

	exps := make([]*profile.Experience, 0, len(m))
	for i := range m {
		exps = append(exps, m[i].Experience)
	}

	return exps, nil
}

// AddExperience ...
func (r Repository) AddExperience(ctx context.Context, userID string, exp *profile.Experience) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	fileIDs := make([]bson.ObjectId, 0, len(exp.Files))
	for i := range exp.Files {
		fileIDs = append(fileIDs, exp.Files[i].ID)
	}

	// append files
	m := []struct {
		File *profile.File `bson:"unused_files"`
	}{}
	err := r.collections[usersCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "Experience",
					"unused_files.id": bson.M{
						"$in": fileIDs,
					},
				},
			},
		},
	).All(&m)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	exp.Files = make([]*profile.File, len(m))
	for i := range m {
		exp.Files[i] = m[i].File
	}

	// adding experience
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"experiences": exp,
			},
		})

	if err != nil {
		return err
	}

	// delete old files
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"unused_files": bson.M{
					"type": "Experience",
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeExperience ...
func (r Repository) ChangeExperience(ctx context.Context, userID string, exp *profile.Experience, changeIsCurrentlyWorking bool) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	bs := make(bson.M)

	if exp.Position != "" {
		bs["experiences.$.position"] = exp.Position
	}
	if exp.Company != "" {
		bs["experiences.$.company"] = exp.Company
	}
	if exp.CityID != nil {
		bs["experiences.$.city_id"] = exp.CityID
	}

	if changeIsCurrentlyWorking {
		if exp.CurrentlyWork {
			bs["experiences.$.currently_work"] = true
			bs["experiences.$.finish_date"] = nil
		} else {
			if !exp.FinishDate.IsZero() {
				bs["experiences.$.currently_work"] = false
				bs["experiences.$.finish_date"] = exp.FinishDate
			} else {
				// return errors.New("wrong_date")
			}
		}
	}

	// BUG:  possible to set start date after finish date
	if !exp.StartDate.IsZero() {
		bs["experiences.$.start_date"] = exp.StartDate
	}
	if exp.Description != nil {
		bs["experiences.$.description"] = exp.Description
	}

	if len(exp.Links) > 0 {
		bs["experiences.$.links"] = exp.Links
	}

	if len(bs) == 0 {
		return errors.New("nothing_to_change")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":            bson.ObjectIdHex(userID),
			"experiences.id": exp.ID,
		},
		bson.M{
			"$set": bs,
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveExperience ...
func (r Repository) RemoveExperience(ctx context.Context, userID string, expID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"experiences": bson.M{
					"id": bson.ObjectIdHex(expID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddLinksInExperience ...
func (r Repository) AddLinksInExperience(ctx context.Context, userID string, expID string, links []*profile.Link) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":            bson.ObjectIdHex(userID),
			"experiences.id": bson.ObjectIdHex(expID),
		},
		bson.M{
			"$push": bson.M{
				"experiences.$.links": bson.M{
					"$each": links,
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddFileInExperience ...
func (r Repository) AddFileInExperience(ctx context.Context, userID string, expID string, file *profile.File) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}
	if expID != "" && !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	var err error

	f := struct {
		File *profile.File `bson:",inline"`
		Type string
	}{
		File: file,
		Type: "Experience",
	}

	if expID != "" {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id":            bson.ObjectIdHex(userID),
				"experiences.id": bson.ObjectIdHex(expID),
			},
			bson.M{
				"$push": bson.M{
					"experiences.$.files": f,
				},
			},
		)
	} else {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(userID),
			},
			bson.M{
				"$push": bson.M{
					"unused_files": f,
				},
			},
		)
	}

	if err != nil {
		return err
	}
	return nil
}

// RemoveFilesInExperience ...
func (r Repository) RemoveFilesInExperience(ctx context.Context, userID string, expID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":            bson.ObjectIdHex(userID),
			"experiences.id": bson.ObjectIdHex(expID),
		},
		bson.M{
			"$pull": bson.M{
				"experiences.$.files": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeLinkInExperience ...
func (r Repository) ChangeLinkInExperience(ctx context.Context, userID string, expID string, linkID string, url string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) || !bson.IsObjectIdHex(linkID) {
		return errors.New("wrong_id")
	}

	result := struct {
		Index int `bson:"index"`
	}{}

	err := r.collections[usersCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(userID),
			},
		},
		{
			"$project": bson.M{
				"experiences": 1,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$experiences",
				"includeArrayIndex":          "experiences.index",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"experiences.id": bson.ObjectIdHex(expID),
			},
		},
		{
			"$project": bson.M{
				"index": "$experiences.index",
			},
		},
	}).One(&result)
	if err != nil {
		return err
	}

	fmt.Println("Index:", result.Index)

	err = r.collections[usersCollection].Update(
		bson.M{
			"_id":                  bson.ObjectIdHex(userID),
			"experiences.links.id": bson.ObjectIdHex(linkID),
		},
		bson.M{
			"$set": bson.M{
				"experiences." + strconv.Itoa(result.Index) + ".links.$.url": url,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveLinksInExperience ...
func (r Repository) RemoveLinksInExperience(ctx context.Context, userID string, expID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":            bson.ObjectIdHex(userID),
			"experiences.id": bson.ObjectIdHex(expID),
		},
		bson.M{
			"$pull": bson.M{
				"experiences.$.links": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetUploadedFilesInExperience ...
func (r Repository) GetUploadedFilesInExperience(ctx context.Context, userID string) ([]*profile.File, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	f := []struct {
		File *profile.File `bson:"unused_files"`
	}{}

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "Experience",
				},
			},
		},
	)
	err := result.All(&f)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		return nil, err
	}

	files := make([]*profile.File, 0, len(f))

	for i := range f {
		files = append(files, f[i].File)
	}

	return files, nil
}

// GetEducations ...
func (r Repository) GetEducations(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Education, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]struct {
		Education *profile.Education `bson:"educations"`
	},
		0,
	)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":        0,
					"educations": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$educations",
				},
			},
			{
				"$sort": bson.M{
					"educations.start_date":  -1,
					"educations.finish_date": 1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := res.All(&m)
	if err != nil {
		return nil, err
	}

	edus := make([]*profile.Education, 0, len(m))
	for i := range m {
		edus = append(edus, m[i].Education)
	}

	return edus, nil
}

// AddAccomplishment ...
func (r Repository) AddAccomplishment(ctx context.Context, userID string, edu *profile.Accomplishment) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	fileIDs := make([]bson.ObjectId, 0, len(edu.Files))
	for i := range edu.Files {
		fileIDs = append(fileIDs, edu.Files[i].ID)
	}

	// append files
	m := []struct {
		File *profile.File `bson:"unused_files"`
	}{}
	err := r.collections[usersCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "Accomplishment",
					"unused_files.id": bson.M{
						"$in": fileIDs,
					},
				},
			},
		},
	).All(&m)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	edu.Files = make([]*profile.File, len(m))
	for i := range m {
		edu.Files[i] = m[i].File
	}

	// adding accomplishment
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"accomplishments": edu,
			},
		})

	if err != nil {
		return err
	}

	// delete old files
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"unused_files": bson.M{
					"type": "Accomplishment",
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddFileInAccomplishment ...
func (r Repository) AddFileInAccomplishment(ctx context.Context, userID string, expID string, file *profile.File) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}
	if expID != "" && !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	var err error

	f := struct {
		File *profile.File `bson:",inline"`
		Type string
	}{
		File: file,
		Type: "Accomplishment",
	}

	if expID != "" {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id":                bson.ObjectIdHex(userID),
				"accomplishments.id": bson.ObjectIdHex(expID),
			},
			bson.M{
				"$push": bson.M{
					"accomplishments.$.files": f,
				},
			},
		)
	} else {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(userID),
			},
			bson.M{
				"$push": bson.M{
					"unused_files": f,
				},
			},
		)
	}

	if err != nil {
		return err
	}
	return nil
}

// RemoveFilesInAccomplishment ...
func (r Repository) RemoveFilesInAccomplishment(ctx context.Context, userID string, expID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"accomplishments.id": bson.ObjectIdHex(expID),
		},
		bson.M{
			"$pull": bson.M{
				"accomplishments.$.files": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddLinksInAccomplishment ...
func (r Repository) AddLinksInAccomplishment(ctx context.Context, userID string, accID string, links []*profile.Link) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(accID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"accomplishments.id": bson.ObjectIdHex(accID),
		},
		bson.M{
			"$push": bson.M{
				"accomplishments.$.links": bson.M{
					"$each": links,
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveLinksInAccomplishment ...
func (r Repository) RemoveLinksInAccomplishment(ctx context.Context, userID string, accID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(accID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"accomplishments.id": bson.ObjectIdHex(accID),
		},
		bson.M{
			"$pull": bson.M{
				"accomplishments.$.links": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetUploadedFilesInAccomplishment ...
func (r Repository) GetUploadedFilesInAccomplishment(ctx context.Context, userID string) ([]*profile.File, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	f := []struct {
		File *profile.File `bson:"unused_files"`
	}{}

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "Accomplishment",
				},
			},
		},
	)
	err := result.All(&f)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		return nil, err
	}

	files := make([]*profile.File, 0, len(f))

	for i := range f {
		files = append(files, f[i].File)
	}

	return files, nil
}

// ChangeEducation ...
func (r Repository) ChangeEducation(ctx context.Context, userID string, edu *profile.Education) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	bs := make(bson.M)

	bs["educations.$.is_currenlty_study"] = edu.IsCurrentlyStudy

	if edu.School != "" {
		bs["educations.$.school"] = edu.School
	}
	if edu.Degree != nil {
		bs["educations.$.degree"] = edu.Degree
	}
	if edu.FieldStudy != "" {
		bs["educations.$.field_study"] = edu.FieldStudy
	}
	if edu.Grade != nil {
		bs["educations.$.grade"] = edu.Grade
	}
	if edu.CityID != nil {
		bs["educations.$.city_id"] = edu.CityID
	}

	if !edu.FinishDate.IsZero() {
		bs["educations.$.finish_date"] = edu.FinishDate
	}
	if !edu.StartDate.IsZero() {
		bs["educations.$.start_date"] = edu.StartDate
	}
	if edu.Description != nil {
		bs["educations.$.description"] = edu.Description
	}

	if edu.IsCurrentlyStudy {
		bs["educations.$.finish_date"] = time.Time{}
	}

	if len(edu.Links) > 0 {
		bs["educations.$.links"] = edu.Links
	}

	if len(bs) == 0 {
		return errors.New("nothing_to_change")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"educations.id": edu.ID,
		},
		bson.M{
			"$set": bs,
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveEducation ...
func (r Repository) RemoveEducation(ctx context.Context, userID string, eduID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(eduID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"educations": bson.M{
					"id": bson.ObjectIdHex(eduID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddEducation ...
func (r Repository) AddEducation(ctx context.Context, userID string, edu *profile.Education) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	fileIDs := make([]bson.ObjectId, 0, len(edu.Files))
	for i := range edu.Files {
		fileIDs = append(fileIDs, edu.Files[i].ID)
	}

	// append files
	m := []struct {
		File *profile.File `bson:"unused_files"`
	}{}
	err := r.collections[usersCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "Education",
					"unused_files.id": bson.M{
						"$in": fileIDs,
					},
				},
			},
		},
	).All(&m)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	edu.Files = make([]*profile.File, len(m))
	for i := range m {
		edu.Files[i] = m[i].File
	}

	// adding experience
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"educations": edu,
			},
		})

	if err != nil {
		return err
	}

	// delete old files
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"unused_files": bson.M{
					"type": "Education",
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddLinksInEducation ...
func (r Repository) AddLinksInEducation(ctx context.Context, userID string, eduID string, links []*profile.Link) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(eduID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"educations.id": bson.ObjectIdHex(eduID),
		},
		bson.M{
			"$push": bson.M{
				"educations.$.links": bson.M{
					"$each": links,
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// AddFileInEducation ...
func (r Repository) AddFileInEducation(ctx context.Context, userID string, eduID string, file *profile.File) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}
	if eduID != "" && !bson.IsObjectIdHex(eduID) {
		return errors.New("wrong_id")
	}

	var err error

	f := struct {
		File *profile.File `bson:",inline"`
		Type string
	}{
		File: file,
		Type: "Education",
	}

	if eduID != "" {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id":           bson.ObjectIdHex(userID),
				"educations.id": bson.ObjectIdHex(eduID),
			},
			bson.M{
				"$push": bson.M{
					"educations.$.files": f,
				},
			},
		)
	} else {
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id": bson.ObjectIdHex(userID),
			},
			bson.M{
				"$push": bson.M{
					"unused_files": f,
				},
			},
		)
	}

	if err != nil {
		return err
	}
	return nil
}

// RemoveFilesInEducation ...
func (r Repository) RemoveFilesInEducation(ctx context.Context, userID string, eduID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(eduID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"educations.id": bson.ObjectIdHex(eduID),
		},
		bson.M{
			"$pull": bson.M{
				"educations.$.files": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeLinkInEducation ...
func (r Repository) ChangeLinkInEducation(ctx context.Context, userID string, eduID string, linkID string, url string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(eduID) || !bson.IsObjectIdHex(linkID) {
		return errors.New("wrong_id")
	}

	result := struct {
		Index int `bson:"index"`
	}{}

	err := r.collections[usersCollection].Pipe([]bson.M{
		{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(userID),
			},
		},
		{
			"$project": bson.M{
				"educations": 1,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$educations",
				"includeArrayIndex":          "educations.index",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"educations.id": bson.ObjectIdHex(eduID),
			},
		},
		{
			"$project": bson.M{
				"index": "$educations.index",
			},
		},
	}).One(&result)
	if err != nil {
		return err
	}

	fmt.Println("Index:", result.Index)

	err = r.collections[usersCollection].Update(
		bson.M{
			"_id":                 bson.ObjectIdHex(userID),
			"educations.links.id": bson.ObjectIdHex(linkID),
		},
		bson.M{
			"$set": bson.M{
				"educations." + strconv.Itoa(result.Index) + ".links.$.url": url,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveLinksInEducation ...
func (r Repository) RemoveLinksInEducation(ctx context.Context, userID string, eduID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(eduID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"educations.id": bson.ObjectIdHex(eduID),
		},
		bson.M{
			"$pull": bson.M{
				"educations.$.links": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetUploadedFilesInEducation ...
func (r Repository) GetUploadedFilesInEducation(ctx context.Context, userID string) ([]*profile.File, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	f := []struct {
		File *profile.File `bson:"unused_files"`
	}{}

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			bson.M{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			bson.M{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			bson.M{
				"$match": bson.M{
					"unused_files.type": "Education",
				},
			},
		},
	)
	err := result.All(&f)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		return nil, err
	}

	files := make([]*profile.File, 0, len(f))

	for i := range f {
		files = append(files, f[i].File)
	}

	return files, nil
}

// GetSkills ...
func (r Repository) GetSkills(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Skill, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]struct {
		Skills *profile.Skill `bson:"skills"`
	},
		0,
	)

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":    0,
					"skills": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$skills",
				},
			},
			{
				"$sort": bson.M{
					"skills.position": 1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := result.All(&m)
	if err != nil {
		return nil, err
	}

	sks := make([]*profile.Skill, 0, len(m))
	for i := range m {
		sks = append(sks, m[i].Skills)
	}

	return sks, nil
}

// GetEndorsements ...
func (r Repository) GetEndorsements(ctx context.Context, skillID string, first uint32, after uint32) ([]string, error) {
	if !bson.IsObjectIdHex(skillID) {
		return nil, errors.New("wrong_id")
	}

	e := make([]struct {
		Skill struct {
			IDs []bson.ObjectId `bson:"endorsements"`
		} `bson:"skills"`
	}, 0)

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"skills.id": bson.ObjectIdHex(skillID),
				},
			},
			{
				"$project": bson.M{
					"_id":                 0,
					"skills.id":           1,
					"skills.endorsements": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$skills",
				},
			},
			{
				"$match": bson.M{
					"skills.id": bson.ObjectIdHex(skillID),
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	)

	err := result.All(&e)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ids := make([]string, 0, len(e))
	for i := range e {
		for j := range e[i].Skill.IDs {
			ids = append(ids, e[i].Skill.IDs[j].Hex())
		}
	}

	return ids, nil
}

// AddSkills ...
func (r Repository) AddSkills(ctx context.Context, userID string, skills []*profile.Skill) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"skills": bson.M{
					"$each": skills,
				},
			},
		})

	if err != nil {
		return err
	}

	return nil
}

// GetPositionOfLastSkill returns the position of last skill
func (r Repository) GetPositionOfLastSkill(ctx context.Context, userID string) (uint32, error) {
	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	m := make([]struct {
		Skills *profile.Skill `bson:"skills"`
	},
		0,
	)

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":    0,
					"skills": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$skills",
				},
			},
			{
				"$sort": bson.M{
					"skills.position": -1,
				},
			},
			{
				"$limit": 1,
			},
		})

	err := result.All(&m)
	if err != nil {
		return 0, err
	}

	if len(m) > 0 {
		return m[0].Skills.Position, nil
	}

	return 0, nil
}

// ChangeOrderOfSkill ...
// BUG: it incerement all elements instead ignoring condition
func (r Repository) ChangeOrderOfSkill(ctx context.Context, userID string, skillID string, position uint32) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(skillID) {
		return errors.New("wrong_id")
	}

	m := struct {
		Skills []*profile.Skill `bson:"skills"`
	}{}

	// getting skill
	err := r.collections[usersCollection].Find(
		bson.M{
			"_id":       bson.ObjectIdHex(userID),
			"skills.id": bson.ObjectIdHex(skillID),
		},
	).Select(
		bson.M{
			"skills.$.id": 1,
		},
	).One(&m)

	if err != nil {
		return err
	}

	if len(m.Skills) < 1 {
		return errors.New("skill not found")
	}

	// set position
	m.Skills[0].Position = position

	// removing skill from array
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"skills": bson.M{
					"id": bson.ObjectIdHex(skillID),
				},
			},
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	// increment position
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
			"skills.position": bson.M{
				"$gt": position,
			},
		},
		bson.M{
			"$inc": bson.M{
				"skills.$[].position": 1,
			},
		},
	)
	if err != nil {
		return err
	}

	// inserting skill in certain position
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"skills": bson.M{
					"$each": m.Skills,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveSkills ...
func (r Repository) RemoveSkills(ctx context.Context, userID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))
	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
	for i := range idsObject {
		removeIDs = append(removeIDs, idsObject[i])
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"skills": bson.M{
					"id": bson.M{
						"$in": removeIDs,
					},
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// VerifySkill ...
func (r Repository) VerifySkill(ctx context.Context, userID string, targetID string, skillID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(targetID) || !bson.IsObjectIdHex(skillID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(targetID),
			"skills.id": bson.ObjectIdHex(skillID),
		},
		bson.M{
			"$push": bson.M{
				"skills.$.endorsements": bson.ObjectIdHex(userID),
			},
		})
	if err != nil {
		return err
	}
	return nil
}

// IsSkillVerified checks if user varify skill for target
func (r Repository) IsSkillVerified(ctx context.Context, userID string, targetID string, skillID string) (bool, error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(targetID) || !bson.IsObjectIdHex(skillID) {
		return false, errors.New("wrong_id")
	}

	count, err := r.collections[usersCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
			"skills": bson.M{
				"$elemMatch": bson.M{
					"id": bson.ObjectIdHex(skillID),
					"endorsements": bson.M{
						"$in": []bson.ObjectId{bson.ObjectIdHex(targetID)},
					},
				},
			},
		}).Count()

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// UnverifySkill ...
// TODO: make it work
func (r Repository) UnverifySkill(ctx context.Context, userID string, targetID string, skillID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(targetID) || !bson.IsObjectIdHex(skillID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(targetID),
			"skills.id": bson.ObjectIdHex(skillID),
		},
		bson.M{
			"$pull": bson.M{
				"skills.$.endorsements": bson.ObjectIdHex(userID),
			},
		})
	if err != nil {
		return err
	}
	return nil
}

// GetInterests ...
func (r Repository) GetInterests(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Interest, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]struct {
		Interests *profile.Interest `bson:"interests"`
	},
		0,
	)

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":       0,
					"interests": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$interests",
				},
			},
			// {
			// 	"$sort": bson.M{
			// 		"interests.position": 1,
			// 	},
			// },
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := result.All(&m)
	if err != nil {
		return nil, err
	}

	ints := make([]*profile.Interest, 0, len(m))
	for i := range m {
		ints = append(ints, m[i].Interests)
	}

	return ints, nil
}

// AddInterest ...
func (r Repository) AddInterest(ctx context.Context, userID string, interest *profile.Interest) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	var err error

	// append files
	if interest.Image != nil {
		m := struct {
			File *profile.File `bson:"unused_files"`
		}{}
		err = r.collections[usersCollection].Pipe(
			[]bson.M{
				bson.M{
					"$match": bson.M{
						"_id": bson.ObjectIdHex(userID),
					},
				},
				bson.M{
					"$project": bson.M{
						"_id":          0,
						"unused_files": 1,
					},
				},
				bson.M{
					"$unwind": bson.M{
						"path": "$unused_files",
					},
				},
				bson.M{
					"$match": bson.M{
						"unused_files.type": "Interest",
						"unused_files.id":   interest.Image.ID,
					},
				},
			},
		).One(&m)
		if err != nil && err != mgo.ErrNotFound {
			return err
		}

		interest.Image = m.File
	}

	// adding interest
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"interests": interest,
			},
		})
	if err != nil {
		return err
	}

	// delete old files
	err = r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"unused_files": bson.M{
					"type": "Interest",
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeInterest ...
func (r Repository) ChangeInterest(ctx context.Context, userID string, interest *profile.Interest) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	bs := make(bson.M)

	if interest.Description != nil {
		bs["interests.$.description"] = interest.Description
	}
	if interest.Interest != "" {
		bs["interests.$.interest"] = interest.Interest
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":          bson.ObjectIdHex(userID),
			"interests.id": bson.ObjectIdHex(interest.GetID()),
		},
		bson.M{
			"$set": bs,
		})

	if err != nil {
		return err
	}

	return nil
}

// RemoveInterest ...
func (r Repository) RemoveInterest(ctx context.Context, userID string, interestID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(interestID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"interests": bson.M{
					"id": bson.ObjectIdHex(interestID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeImageInterest ...
// TODO:
func (r Repository) ChangeImageInterest(ctx context.Context, userID string, interestID string, image *profile.File) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}
	if interestID != "" && !bson.IsObjectIdHex(interestID) {
		return errors.New("wrong_id")
	}

	var err error

	if interestID != "" {
		// add to the existing interest
		err = r.collections[usersCollection].Update(
			bson.M{
				"_id":          bson.ObjectIdHex(userID),
				"interests.id": bson.ObjectIdHex(interestID),
			},
			bson.M{
				"$set": bson.M{
					// "interests.$.image": data.GetURL(),
					"interests.$.image": image,
				},
			},
		)
	} else {
		log.Println("unused image for interest") // delete

		// add for new interest
		f := struct {
			File *profile.File `bson:",inline"`
			Type string
		}{
			File: image,
			Type: "Interest",
		}

		// retrive uploaded file
		err = r.collections[usersCollection].Pipe(
			[]bson.M{
				bson.M{
					"$match": bson.M{
						"_id": bson.ObjectIdHex(userID),
					},
				},
				bson.M{
					"$project": bson.M{
						"_id":          0,
						"unused_files": 1,
					},
				},
				bson.M{
					"$unwind": bson.M{
						"path": "$unused_files",
					},
				},
				bson.M{
					"$match": bson.M{
						"unused_files.type": "Interest",
					},
				},
			},
		).One(&f)
		if err != nil && err != mgo.ErrNotFound {
			return err
		}

		// if file has been uploaded
		if f.File != nil {
			err = r.collections[usersCollection].Update(
				bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
				bson.M{
					"$push": bson.M{
						"unused_files": image,
					},
				},
			)
			if err != nil {
				return err
			}
		}
		// else {
		// 	err = r.collections[usersCollection].Update(
		// 		bson.M{
		// 			"_id":             bson.ObjectIdHex(userID),
		// 			"unused_files.id": f.File.ID,
		// 		},
		// 		bson.M{
		// 			"$set": bson.M{
		// 				"unused_files.$.url": data.GetURL(),
		// 			},
		// 		},
		// 	)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	if err != nil {
		return err
	}

	return nil
}

// RemoveImageInInterest ...
func (r Repository) RemoveImageInInterest(ctx context.Context, userID string, interestID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(interestID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":          bson.ObjectIdHex(userID),
			"interests.id": bson.ObjectIdHex(interestID),
		},
		bson.M{
			"$set": bson.M{
				"interests.$.image": nil,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetUnuploadImageInInterest ...
func (r Repository) GetUnuploadImageInInterest(ctx context.Context, userID string) (*profile.File, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	f := struct {
		File *profile.File `bson:"unused_files"`
	}{}

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":          0,
					"unused_files": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$unused_files",
				},
			},
			{
				"$match": bson.M{
					"unused_files.type": "Interest",
				},
			},
		},
	)
	err := result.One(&f)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		return nil, err
	}

	return f.File, nil
}

// GetOriginImageInInterest ...
func (r Repository) GetOriginImageInInterest(ctx context.Context, userID, interestID string) (string, error) {
	if !bson.IsObjectIdHex(userID) {
		return "", errors.New("wrong_id")
	}

	if !bson.IsObjectIdHex(interestID) {
		return "", errors.New("wrong_id")
	}

	m := struct {
		Image string `bson:"image_origin"`
	}{}

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"interests": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$interests",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$interests",
				},
			},
			{
				"$match": bson.M{
					"id": bson.ObjectIdHex(interestID),
				},
			},
		},
	)

	err := result.One(&m)
	if err != nil {
		return "", err
	}

	return m.Image, nil
}

// ChangeOriginImageInInterest ...
func (r Repository) ChangeOriginImageInInterest(ctx context.Context, userID, interestID, url string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	if !bson.IsObjectIdHex(interestID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":          bson.ObjectIdHex(userID),
			"interests.id": bson.ObjectIdHex(interestID),
		},
		bson.M{
			"$set": bson.M{
				"interests.$.image_origin": url,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetAccomplishments ...
func (r Repository) GetAccomplishments(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.Accomplishment, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]struct {
		Accomplishment *profile.Accomplishment `bson:"accomplishments"`
	},
		0,
	)

	res := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":             0,
					"accomplishments": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$accomplishments",
				},
			},
			{
				"$sort": bson.M{
					"accomplishments.start_date":  -1,
					"accomplishments.finish_date": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := res.All(&m)
	if err != nil {
		return nil, err
	}

	acc := make([]*profile.Accomplishment, 0, len(m))
	for i := range m {
		acc = append(acc, m[i].Accomplishment)
	}

	return acc, nil
}

// // AddFileInExperience ...
// func (r Repository) AddFileInExperience(ctx context.Context, userID string, expID string, file *profile.File) error {
// 	if !bson.IsObjectIdHex(userID) {
// 		return errors.New("wrong_id")
// 	}
// 	if expID != "" && !bson.IsObjectIdHex(expID) {
// 		return errors.New("wrong_id")
// 	}

// 	var err error

// 	f := struct {
// 		File *profile.File `bson:",inline"`
// 		Type string
// 	}{
// 		File: file,
// 		Type: "Experience",
// 	}

// 	if expID != "" {
// 		err = r.collections[usersCollection].Update(
// 			bson.M{
// 				"_id":            bson.ObjectIdHex(userID),
// 				"experiences.id": bson.ObjectIdHex(expID),
// 			},
// 			bson.M{
// 				"$push": bson.M{
// 					"experiences.$.files": f,
// 				},
// 			},
// 		)
// 	} else {
// 		err = r.collections[usersCollection].Update(
// 			bson.M{
// 				"_id": bson.ObjectIdHex(userID),
// 			},
// 			bson.M{
// 				"$push": bson.M{
// 					"unused_files": f,
// 				},
// 			},
// 		)
// 	}

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // RemoveFilesInExperience ...
// func (r Repository) RemoveFilesInExperience(ctx context.Context, userID string, expID string, ids []string) error {
// 	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) {
// 		return errors.New("wrong_id")
// 	}

// 	idsObject := make([]bson.ObjectId, 0, len(ids))
// 	for i := range ids {
// 		if bson.IsObjectIdHex(ids[i]) {
// 			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
// 		} else {
// 			return errors.New("wrong_id")
// 		}
// 	}

// 	removeIDs := make([]bson.ObjectId, 0, len(idsObject))
// 	for i := range idsObject {
// 		removeIDs = append(removeIDs, idsObject[i])
// 	}

// 	err := r.collections[usersCollection].Update(
// 		bson.M{
// 			"_id":            bson.ObjectIdHex(userID),
// 			"experiences.id": bson.ObjectIdHex(expID),
// 		},
// 		bson.M{
// 			"$pull": bson.M{
// 				"experiences.$.files": bson.M{
// 					"id": bson.M{
// 						"$in": removeIDs,
// 					},
// 				},
// 			},
// 		})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // AddAccomplishment ...
// func (r Repository) AddAccomplishment(ctx context.Context, userID string, accomplishment *profile.Accomplishment) error {
// 	if !bson.IsObjectIdHex(userID) {
// 		return errors.New("wrong_id")
// 	}

// 	// bs := make(bson.M)
// 	//
// 	// if accomplishment.FinishDate != nil {
// 	// 	bs["finish_date"] = accomplishment.FinishDate
// 	// }
// 	//
// 	// switch accomplishment.Type {
// 	// // adding certificate
// 	// case profile.AccomplishmentTypeCertificate:
// 	// 	if accomplishment.Issuer != nil {
// 	// 		bs["issuer"] = accomplishment.Issuer
// 	// 	}
// 	// 	if data.GetIsLicenseNumberNull() == false {
// 	// 		bs["license_number"] = data.GetLicenseNumber()
// 	// 	}
// 	// 	if data.GetIsIsExpireNull() == false {
// 	// 		bs["is_expire"] = data.GetIsExpire()
// 	// 	}
// 	// 	if data.GetIsURLNull() == false {
// 	// 		bs["url"] = data.GetURL()
// 	// 	}
// 	//
// 	// case Accomplishment_AccomplishmentType_License:
// 	// 	if data.GetIsIssuerNull() == false {
// 	// 		bs["issuer"] = data.GetIssuer()
// 	// 	}
// 	// 	if data.GetIsLicenseNumberNull() == false {
// 	// 		bs["license_number"] = data.GetLicenseNumber()
// 	// 	}
// 	// 	if data.GetIsIsExpireNull() == false {
// 	// 		bs["is_expire"] = data.GetIsExpire()
// 	// 	}
// 	//
// 	// case Accomplishment_AccomplishmentType_Award:
// 	// 	if data.GetIsIssuerNull() == false {
// 	// 		bs["issuer"] = data.GetIssuer()
// 	// 	}
// 	// 	if data.GetIsDescriptionNull() == false {
// 	// 		bs["description"] = data.GetDescription()
// 	// 	}
// 	//
// 	// case Accomplishment_AccomplishmentType_Project:
// 	// 	if data.GetIsURLNull() == false {
// 	// 		bs["url"] = data.GetURL()
// 	// 	}
// 	// 	if data.GetIsDescriptionNull() == false {
// 	// 		bs["description"] = data.GetDescription()
// 	// 	}
// 	// 	if data.GetIsStartDateNull() == false {
// 	// 		bs["start_date"], _ = time.Parse("1-2006", data.GetStartDate())
// 	// 	}
// 	//
// 	// case Accomplishment_AccomplishmentType_Publication:
// 	// 	if data.GetIsIssuerNull() == false {
// 	// 		bs["issuer"] = data.GetIssuer()
// 	// 	}
// 	// 	if data.GetIsURLNull() == false {
// 	// 		bs["url"] = data.GetURL()
// 	// 	}
// 	// 	if data.GetIsDescriptionNull() == false {
// 	// 		bs["description"] = data.GetDescription()
// 	// 	}
// 	//
// 	// case Accomplishment_AccomplishmentType_Test:
// 	// 	if data.GetIsDescriptionNull() == false {
// 	// 		bs["description"] = data.GetDescription()
// 	// 	}
// 	// 	bs["score"] = data.GetScore()
// 	// }

// 	err := r.collections[usersCollection].Update(
// 		bson.M{
// 			"_id": bson.ObjectIdHex(userID),
// 		},
// 		bson.M{
// 			"$push": bson.M{
// 				"accomplishments": accomplishment,
// 			},
// 		})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// ChangeAccomplishment ...
func (r Repository) ChangeAccomplishment(ctx context.Context, userID string, accomplishment *profile.Accomplishment) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(accomplishment.GetID()) {
		return errors.New("wrong_id")
	}

	a := struct {
		Accomplishments []*profile.Accomplishment `bson:"accomplishments"`
	}{}

	// getting accomplishment
	err := r.collections[usersCollection].Find(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"accomplishments.id": bson.ObjectIdHex(accomplishment.GetID()),
		},
	).Select(
		bson.M{
			"accomplishments.$.id": bson.ObjectIdHex(accomplishment.GetID()),
		},
	).One(&a)

	if err != nil {
		return err
	}

	if len(a.Accomplishments) < 1 {
		return errors.New("not found")
	}

	bs := make(bson.M)

	if accomplishment.Name != "" {
		bs["accomplishments.$.name"] = accomplishment.Name
	}

	switch a.Accomplishments[0].Type {
	// certificate
	case profile.AccomplishmentTypeCertificate:

		if len(accomplishment.Links) > 0 {
			bs["accomplishments.$.links"] = accomplishment.Links
		}
		if accomplishment.Issuer != nil {
			bs["accomplishments.$.issuer"] = accomplishment.Issuer
		}
		if accomplishment.LicenseNumber != nil {
			bs["accomplishments.$.license_number"] = accomplishment.LicenseNumber
		}
		if accomplishment.StartDate != nil {
			bs["accomplishments.$.start_date"] = accomplishment.StartDate
		}
		if accomplishment.IsExpire != nil {
			// bs["accomplishments.$.is_expire"] = accomplishment.IsExpire
			if *accomplishment.IsExpire == false {
				bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
				bs["accomplishments.$.is_expire"] = nil
			} else {
				bs["accomplishments.$.is_expire"] = true
				bs["accomplishments.$.finish_date"] = nil
			}
		}
		// if accomplishment.IsExpire != nil {
		// 	bs["accomplishments.$.is_expire"] = accomplishment.IsExpire
		// }
		// if accomplishment.FinishDate != nil {
		// 	bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
		// } else {
		// 	bs["accomplishments.$.is_expire"] = false
		// }
		if accomplishment.URL != nil {
			bs["accomplishments.$.url"] = accomplishment.URL
		}
	// license
	case profile.AccomplishmentTypeLicense:

		if len(accomplishment.Links) > 0 {
			bs["accomplishments.$.links"] = accomplishment.Links
		}
		if accomplishment.Issuer != nil {
			bs["accomplishments.$.issuer"] = accomplishment.Issuer
		}
		if accomplishment.LicenseNumber != nil {
			bs["accomplishments.$.license_number"] = accomplishment.LicenseNumber
		}
		if accomplishment.StartDate != nil {
			bs["accomplishments.$.start_date"] = accomplishment.StartDate
		}
		if accomplishment.IsExpire != nil {
			bs["accomplishments.$.is_expire"] = accomplishment.IsExpire
			if *accomplishment.IsExpire == false {
				bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
				bs["accomplishments.$.is_expire"] = nil
			} else {
				bs["accomplishments.$.is_expire"] = true
				bs["accomplishments.$.finish_date"] = nil
			}
		}
		// if accomplishment.IsExpire != nil {
		// 	bs["accomplishments.$.is_expire"] = accomplishment.IsExpire
		// }
		// if accomplishment.FinishDate != nil {
		// 	bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
		// } else {
		// 	bs["accomplishments.$.is_expire"] = false
		// }
		//  award
	case profile.AccomplishmentTypeAward:

		if len(accomplishment.Links) > 0 {
			bs["accomplishments.$.links"] = accomplishment.Links
		}
		if accomplishment.Issuer != nil {
			bs["accomplishments.$.issuer"] = accomplishment.Issuer
		}
		if accomplishment.Description != nil {
			bs["accomplishments.$.description"] = accomplishment.Description
		}
		if accomplishment.FinishDate != nil {
			bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
		}
		// project
	case profile.AccomplishmentTypeProject:

		if len(accomplishment.Links) > 0 {
			bs["accomplishments.$.links"] = accomplishment.Links
		}
		if accomplishment.URL != nil {
			bs["accomplishments.$.url"] = accomplishment.URL
		}
		if accomplishment.Description != nil {
			bs["accomplishments.$.description"] = accomplishment.Description
		}
		if accomplishment.StartDate != nil {
			bs["accomplishments.$.start_date"] = accomplishment.StartDate
		}
		if accomplishment.FinishDate != nil {
			bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
		}
		// publication
	case profile.AccomplishmentTypePublication:

		if len(accomplishment.Links) > 0 {
			bs["accomplishments.$.links"] = accomplishment.Links
		}
		if accomplishment.Issuer != nil {
			bs["accomplishments.$.issuer"] = accomplishment.Issuer
		}
		if accomplishment.URL != nil {
			bs["accomplishments.$.url"] = accomplishment.URL
		}
		if accomplishment.Description != nil {
			bs["accomplishments.$.description"] = accomplishment.Description
		}
		if accomplishment.FinishDate != nil {
			bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
		}
		// test
	case profile.AccomplishmentTypeTest:

		if len(accomplishment.Links) > 0 {
			bs["accomplishments.$.links"] = accomplishment.Links
		}
		if accomplishment.Description != nil {
			bs["accomplishments.$.description"] = accomplishment.Description
		}
		if accomplishment.Score != nil {
			bs["accomplishments.$.score"] = accomplishment.Score
		}
		if accomplishment.FinishDate != nil {
			bs["accomplishments.$.finish_date"] = accomplishment.FinishDate
		}
	}

	if len(bs) == 0 {
		return errors.New("nothing_to_change")
	}

	err = r.collections[usersCollection].Update(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"accomplishments.id": bson.ObjectIdHex(accomplishment.GetID()),
		},
		bson.M{
			"$set": bs,
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveAccomplishment ...
func (r Repository) RemoveAccomplishment(ctx context.Context, userID string, accomplishmentID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(accomplishmentID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"accomplishments": bson.M{
					"id": bson.ObjectIdHex(accomplishmentID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetKnownLanguages ...
func (r Repository) GetKnownLanguages(ctx context.Context, userID string, first uint32, after uint32) ([]*profile.KnownLanguage, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := make([]struct {
		KnownLanguages *profile.KnownLanguage `bson:"known_languages"`
	},
		0,
	)

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": bson.M{
					"_id":             0,
					"known_languages": 1,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$known_languages",
				},
			},
			{
				"$sort": bson.M{
					"known_languages.rank": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		})

	err := result.All(&m)
	if err != nil {
		return nil, err
	}

	langs := make([]*profile.KnownLanguage, 0, len(m))
	for i := range m {
		langs = append(langs, m[i].KnownLanguages)
	}

	return langs, nil
}

// AddKnownLanguage ...
func (r Repository) AddKnownLanguage(ctx context.Context, userID string, lang *profile.KnownLanguage) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	// adding known language
	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$push": bson.M{
				"known_languages": lang,
			},
		})

	if err != nil {
		return err
	}

	return nil
}

// ChangeKnownLanguage ...
func (r Repository) ChangeKnownLanguage(ctx context.Context, userID string, lang *profile.KnownLanguage) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(lang.GetID()) {
		return errors.New("wrong_id")
	}

	bs := make(bson.M)

	if lang.Rank != 0 {
		bs["known_languages.$.rank"] = lang.Rank
	}
	if lang.Language != "" {
		bs["known_languages.$.language"] = lang.Language
	}

	if len(bs) == 0 {
		return errors.New("nothing_to_change")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"known_languages.id": bson.ObjectIdHex(lang.GetID()),
		},
		bson.M{
			"$set": bs,
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveKnownLanguage ...
func (r Repository) RemoveKnownLanguage(ctx context.Context, userID string, langID string) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(langID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$pull": bson.M{
				"known_languages": bson.M{
					"id": bson.ObjectIdHex(langID),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeHeadline ...
func (r Repository) ChangeHeadline(ctx context.Context, userID string, headline string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"headline": headline,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeStory ...
func (r Repository) ChangeStory(ctx context.Context, userID string, story string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"story": story,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// GetOriginAvatar ...
func (r Repository) GetOriginAvatar(ctx context.Context, userID string) (avatarPath string, err error) {
	if !bson.IsObjectIdHex(userID) {
		return "", errors.New("wrong_id")
	}

	m := struct {
		AvatarOrigin string `bson:"avatar_origin"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
	})

	err = result.One(&m)
	if err != nil {
		return
	}

	return m.AvatarOrigin, nil
}

// ChangeOriginAvatar ...
func (r Repository) ChangeOriginAvatar(ctx context.Context, userID string, url string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"avatar_origin": url,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// ChangeAvatar ...
func (r Repository) ChangeAvatar(ctx context.Context, userID string, url string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"avatar": url,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// RemoveAvatar ...
func (r Repository) RemoveAvatar(ctx context.Context, userID string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"avatar":        nil,
				"avatar_origin": nil,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// SaveReport ...
func (r Repository) SaveReport(ctx context.Context, userID string, report *userReport.Report) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(report.GetUserID()) {
		return errors.New("wrong_id")
	}

	if err := r.collections[reportsCollection].Insert(report); err != nil {
		return err
	}

	return nil
}

// GetInfoAboutCompletionProfile returns true for fileds if it is not empty
func (r Repository) GetInfoAboutCompletionProfile(ctx context.Context, userID string) (exp, edu, skills, langs, interests, tools bool, err error) {
	if !bson.IsObjectIdHex(userID) {
		err = errors.New("wrong_id")
		return
	}

	m := struct {
		Exp    bool `bson:"experiences"`
		Edu    bool `bson:"educations"`
		Tols   bool `bson:"tools_technologies"`
		Skills bool `bson:"skills"`
		Langs  bool `bson:"langs"`
		Inter  bool `bson:"interests"`
	}{}

	result := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID)},
			},
			{
				"$project": bson.M{
					"experiences": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$experiences"},
							bson.M{"$gt": []interface{}{bson.M{"$size": "$experiences"}, 0}},
							false,
						},
					},
					"educations": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$educations"},
							bson.M{"$gt": []interface{}{bson.M{"$size": "$educations"}, 0}},
							false,
						},
					},
					"tools_technologies": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$tools_technologies"},
							bson.M{"$gt": []interface{}{bson.M{"$size": "$tools_technologies"}, 0}},
							false,
						},
					},
					"skills": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$skills"},
							bson.M{"$gt": []interface{}{bson.M{"$size": "$skills"}, 0}},
							false,
						},
					},
					"langs": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$known_languages"},
							bson.M{"$gt": []interface{}{bson.M{"$size": "$known_languages"}, 0}},
							false,
						},
					},
					"interests": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$interests"},
							bson.M{"$gt": []interface{}{bson.M{"$size": "$interests"}, 0}},
							false,
						},
					},
				},
			},
		},
	)

	err = result.One(&m)
	if err != nil {
		return
	}

	return m.Exp, m.Edu, m.Skills, m.Langs, m.Inter, m.Tols, nil
}

// Translations

// SaveUserProfileTranslation ...
func (r Repository) SaveUserProfileTranslation(ctx context.Context, userID string, lang string, tr *profile.Translation) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// if nickname != nil {
	// 	bs["nickname"] = *nickname
	// }
	//
	// // bs["lang"] = lang
	// bs["firstname"] = firstname
	// bs["lastname"] = lastname
	// bs["headline"] = lastname

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				"translation." + lang: tr,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserExperienceTranslation ...
func (r Repository) SaveUserExperienceTranslation(ctx context.Context, userID string, expID string, lang string, tr *profile.ExperienceTranslation /*position, company, description string*/) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(expID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// // bs["experiences.$.lang"] = lang
	// bs["position"] = position
	// bs["company"] = company
	// bs["description"] = description

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":            bson.ObjectIdHex(userID),
			"experiences.id": bson.ObjectIdHex(expID),
		},
		bson.M{
			"$set": bson.M{
				"experiences.$.translations." + lang: tr, //bs,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserEducationTranslation ...
func (r Repository) SaveUserEducationTranslation(ctx context.Context, userID string, educationID string, lang string, tr *profile.EducationTranslation /*school, degree, fieldStudy, grade, description string*/) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(educationID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// // bs["experiences.$.lang"] = lang
	// bs["school"] = school
	// bs["degree"] = degree
	// bs["field_study"] = fieldStudy
	// bs["grade"] = grade
	// bs["description"] = description

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"educations.id": bson.ObjectIdHex(educationID),
		},
		bson.M{
			"$set": bson.M{
				"educations.$.translations." + lang: tr, //bs,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserInterestTranslation ...
func (r Repository) SaveUserInterestTranslation(ctx context.Context, userID string, interestID string, lang string, tr *profile.InterestTranslation /*interest, description string*/) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(interestID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// // bs["experiences.$.lang"] = lang
	// bs["interest"] = interest
	// bs["description"] = description

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":          bson.ObjectIdHex(userID),
			"interests.id": bson.ObjectIdHex(interestID),
		},
		bson.M{
			"$set": bson.M{
				"interests.$.translations." + lang: tr, //bs,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserSkillTranslation ...
func (r Repository) SaveUserSkillTranslation(ctx context.Context, userID string, skillID string, lang string, tr *profile.SkillTranslation /*skill string*/) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(skillID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// // bs["experiences.$.lang"] = lang
	// bs["skill"] = skill

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":       bson.ObjectIdHex(userID),
			"skills.id": bson.ObjectIdHex(skillID),
		},
		bson.M{
			"$set": bson.M{
				"skills.$.translations." + lang: tr, //bs,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserAccomplishmentTranslation ...
func (r Repository) SaveUserAccomplishmentTranslation(ctx context.Context, userID string, accomplishmentID string, lang string, tr *profile.AccomplishmentTranslation) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(accomplishmentID) {
		return errors.New("wrong_id")
	}

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":                bson.ObjectIdHex(userID),
			"accomplishments.id": bson.ObjectIdHex(accomplishmentID),
		},
		bson.M{
			"$set": bson.M{
				"accomplishments.$.translations." + lang: tr, //bs,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserPortfolioTranslation ...
func (r Repository) SaveUserPortfolioTranslation(ctx context.Context, userID string, portfolioID string, lang string, tr *profile.PortfolioTranslation /*interest, description string*/) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(portfolioID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// // bs["experiences.$.lang"] = lang
	// bs["interest"] = interest
	// bs["description"] = description

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":           bson.ObjectIdHex(userID),
			"portfolios.id": bson.ObjectIdHex(portfolioID),
		},
		bson.M{
			"$set": bson.M{
				"portfolios.$.translations." + lang: tr, //bs,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserToolTechnologyTranslation ...
func (r Repository) SaveUserToolTechnologyTranslation(ctx context.Context, userID string, tooltechID string, lang string, tr *profile.ToolTechnologyTranslation /*interest, description string*/) error {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(tooltechID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// // bs["experiences.$.lang"] = lang
	// bs["interest"] = interest
	// bs["description"] = description

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id":                   bson.ObjectIdHex(userID),
			"tools_technologies.id": bson.ObjectIdHex(tooltechID),
		},
		bson.M{
			"$set": bson.M{
				"tools_technologies.$.translations." + lang: tr, //bs,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveTransaltion ...
func (r Repository) RemoveTransaltion(ctx context.Context, userID string, lang string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	// bs := make(bson.M)
	//
	// if nickname != nil {
	// 	bs["nickname"] = *nickname
	// }
	//
	// // bs["lang"] = lang
	// bs["firstname"] = firstname
	// bs["lastname"] = lastname
	// bs["headline"] = lastname

	err := r.collections[usersCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$unset": bson.M{
				"translation." + lang: "",
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsInvitationSend ...
func (r Repository) IsInvitationSend(ctx context.Context, email string) (bool, error) {
	count, err := r.collections[emailInvitationsCollection].Find(
		bson.M{
			"email": email,
		},
	).Count()
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// SaveInvitation ...
func (r Repository) SaveInvitation(ctx context.Context, inv invitation.Invitation) error {
	err := r.collections[emailInvitationsCollection].Insert(inv)
	if err != nil {
		return err
	}

	return nil
}

// GetInvitation ...
func (r Repository) GetInvitation(ctx context.Context, userID string) ([]invitation.Invitation, int32, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, 0, errors.New("wrong_id")
	}

	result := struct {
		Invs   []invitation.Invitation `bson:"invitations"`
		Amount int32                   `bson:"amount"`
	}{}

	err := r.collections[emailInvitationsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"user_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$group": bson.M{
					"_id": "$user_id",
					"invitations": bson.M{
						"$addToSet": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"amount": bson.M{
						"$sum": bson.M{"$size": "$invitations"},
					},
				},
			},
		},
	).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	return result.Invs, result.Amount, nil
}

// GetInvitationForCompany ...
func (r Repository) GetInvitationForCompany(ctx context.Context, companyID string) ([]invitation.Invitation, int32, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, 0, errors.New("wrong_id")
	}

	result := struct {
		Invs   []invitation.Invitation `bson:"invitations"`
		Amount int32                   `bson:"amount"`
	}{}
	err := r.collections[emailInvitationsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"company_id": bson.ObjectIdHex(companyID),
				},
			},
			{
				"$group": bson.M{
					"_id": "$group_id",
					"invitations": bson.M{
						"$addToSet": "$$ROOT",
					},
				},
			},
			{
				"$addFields": bson.M{
					"amount": bson.M{
						"$sum": bson.M{"$size": "$invitations"},
					},
				},
			},
		},
	).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	return result.Invs, result.Amount, nil
}

// GetPrivacyMyConnections ...
func (r Repository) GetPrivacyMyConnections(ctx context.Context, userID string) (account.PermissionType, error) {
	if !bson.IsObjectIdHex(userID) {
		return account.PermissionTypeNone, errors.New("wrong_id")
	}

	pr := struct {
		Privacy account.Privacy `bson:"privacy"`
	}{}

	result := r.collections[usersCollection].Find(bson.M{
		"_id": bson.ObjectIdHex(userID),
	})
	err := result.One(&pr)
	if err != nil {
		return account.PermissionTypeNone, err
	}

	return pr.Privacy.MyConnections, nil
}

// GetUsersForAdvert ...
func (r Repository) GetUsersForAdvert(ctx context.Context, data account.UserForAdvert) ([]string, error) {
	if !bson.IsObjectIdHex(data.OwnerID) {
		return nil, errors.New("wrong_id")
	}

	// "_id":{
	// 	$ne: ObjectId('5d8b1b6cd356540001682dca')
	//   }
	// 	"birthday.birthday":{
	// 	 $gt: ISODate('2020')
	// 	 $lt
	// 	},
	// 	"gender.gender":"MALE"
	//   }
	// "location.country.id": {
	//     $in:['AS' , 'GE']
	//   },

	match := bson.M{
		"_id": bson.M{
			"$ne": bson.ObjectIdHex(data.OwnerID),
		},
	}

	if data.Gender != "" {
		match["gender.gender"] = data.Gender
	}

	if len(data.Locations) > 0 {
		match["location.country.id"] = bson.M{
			"$in": data.Locations,
		}
	}

	if data.AgeFrom != 0 || data.AgeTo != 0 {

		match["birthday.birthday"] = bson.M{
			"$lt": time.Date(int(data.AgeFrom), time.January, 1, 0, 0, 0, 0, &time.Location{}),
			"$gt": time.Date(int(data.AgeTo), time.January, 1, 0, 0, 0, 0, &time.Location{}),
		}
	}

	log.Printf("ADVERT DATA %+v", match)
	ids := struct {
		ID []bson.ObjectId `bson:"ids"`
	}{}

	err := r.collections[usersCollection].Pipe(
		[]bson.M{
			{
				"$match": match,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"ids": bson.M{
						"$push": "$$ROOT._id",
					},
				},
			},
		},
	).One(&ids)

	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(ids.ID))

	for _, i := range ids.ID {
		res = append(res, i.Hex())
	}

	return res, nil
}
