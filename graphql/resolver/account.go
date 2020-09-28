package resolver

import (
	"context"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/notificationsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func (_ *Resolver) IdentifyCountry(ctx context.Context) (*string, error) {
	countryID, err := user.IdentifyCountry(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}

	if countryID.GetID() != "" {
		id := countryID.GetID()
		return &id, nil
	}

	return nil, nil
}

func (_ *Resolver) GetAccount(ctx context.Context) (*AccountResolver, error) {
	account, err := user.GetAccount(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}

	a := AccountResolver{}
	a.R = &Account{
		Firstname: account.GetFirstName(),
		Lastname:  account.GetLastname(),
		Status: func(s userRPC.Status) string {
			switch s {
			case userRPC.Status_ACTIVATED:
				return "activated"
			case userRPC.Status_DISABLED:
				return "disabled"
			case userRPC.Status_BLOCKED:
				return "blocked"
			}
			return "not_activated"
		}(userRPC.Status_ACTIVATED), // TODO:
		// 	Wallet Wallet
		Firstname_native: &FirstnameNative{
			Firstname: account.GetNativeName().GetName(),
			Language:  account.GetNativeName().GetLanguageID(),
			// Permission: Permission{
			// 	Type: account.GetNativeName().GetPermission().GetType().String(),
			// },
			Permission: toPermission(account.GetNativeName().GetPermission()),
		},
		// Title *Title
		Patronymic: &Patronymic{
			Patronymic: account.GetPatronymic().GetPatronymic(),
			// Permission: Permission{
			// 	Type: account.GetPatronymic().GetPermission().GetType().String(),
			// },
			Permission: toPermission(account.GetPatronymic().GetPermission()),
		},
		Nickname: &Nickname{
			Nickname: account.GetNickName().GetNickname(),
			// Permission: Permission{
			// 	Type: account.GetNickName().GetPermission().GetType().String(),
			// },
			Permission: toPermission(account.GetNickName().GetPermission()),
		},
		Middlename: &Middlename{
			Middlename: account.GetMiddleName().GetMiddlename(),
			// Permission: Permission{
			// 	Type: account.GetMiddleName().GetPermission().GetType().String(),
			// },
			Permission: toPermission(account.GetMiddleName().GetPermission()),
		},

		Birthday: &Birthday{
			Birthday: account.GetBirthday().GetBirthday(),
			// Permission: Permission{
			// 	Type: account.GetBirthday().GetPermission().GetType().String(),
			// },
			Permission: toPermission(account.GetBirthday().GetPermission()),
		},

		Gender: &Gender{
			Gender: func(s string) string {
				if s == "FEMALE" {
					return "female"
				}
				return "male"
			}(account.GetGender().GetGender().String()),
			Permission: toPermission(account.GetGender().GetPermission()),
		},

		Ui_language: account.GetLanguageID(),

		Privacy: Privacy{
			// Sharing_edits:    account.GetPrivacies().GetShareEdits().String(),
			Sharing_edits: toPermissionType(account.GetPrivacies().GetShareEdits()),
			// Profile_pictures: account.GetPrivacies().GetProfilePicture().String(),
			Profile_pictures: toPermissionType(account.GetPrivacies().GetProfilePicture()),
			// My_connections:   account.GetPrivacies().GetMyConnections().String(),
			My_connections: toPermissionType(account.GetPrivacies().GetMyConnections()),
			// Find_by_email:  account.GetPrivacies().GetFindByEmail().String(),
			Find_by_email: toPermissionType(account.GetPrivacies().GetFindByEmail()),
			// Find_by_phone: account.GetPrivacies().GetFindByPhone().String(),
			Find_by_phone: toPermissionType(account.GetPrivacies().GetFindByPhone()),
			// Active_status: account.GetPrivacies().GetActiveStatus().String(),
			Active_status: toPermissionType(account.GetPrivacies().GetActiveStatus()),
		},

		Notifications: Notification{
			Birthdays:              account.GetNotifications().GetBirthdays(),
			Endorsements:           account.GetNotifications().GetEndorsements(),
			Job_changes_in_network: account.GetNotifications().GetJobChangesInNetwork(),
			Import_contacts_joined: account.GetNotifications().GetImportContactsJoined(),
			Job_recommendations:    account.GetNotifications().GetJobRecommendations(),
			Accept_invitation:      account.GetNotifications().GetAcceptInvitation(),
			New_followers:          account.GetNotifications().GetNewFollowers(),
			New_chat_message:       account.GetNotifications().GetNewChatMessage(),
			Email_updates:          account.GetNotifications().GetEmailUpdates(),
			Connection_request:     account.GetNotifications().GetConnectionRequest(),
		},
		Editable:             account.GetIsEditable(),
		Amount_of_sessions:   account.GetAmountOfSessions(),
		Last_change_password: account.GetLastChangePassword(),
	}

	a.R.Phone = make([]PhoneAccount, len(account.GetPhones()))
	for i := range account.GetPhones() {
		var phone PhoneAccount

		phone.ID = account.GetPhones()[i].GetId()
		phone.Country_iso = account.GetPhones()[i].GetCountryAbbreviation()
		phone.Country_code = account.GetPhones()[i].GetCountryCode().GetCode() // TODO: how about ID?
		phone.Number = account.GetPhones()[i].GetNumber()
		phone.Activated = account.GetPhones()[i].GetIsActivated()
		phone.Primary = account.GetPhones()[i].GetIsPrimary()
		// phone.Permission = Permission{
		// 	Type: account.GetPhones()[i].GetPermission().GetType().String(),
		// }
		phone.Permission = toPermission(account.GetPhones()[i].GetPermission())

		a.R.Phone[i] = phone
	}

	a.R.Email = make([]EmailAccount, len(account.GetEmails()))
	for i := range account.GetEmails() {
		var email EmailAccount

		email.ID = account.GetEmails()[i].GetId()
		email.Email = account.GetEmails()[i].GetEmail()
		email.Activated = account.GetEmails()[i].GetIsActivated()
		email.Primary = account.GetEmails()[i].GetIsPrimary()
		// email.Permission = Permission{
		// 	Type: account.GetEmails()[i].GetPermission().GetType().String(),
		// }
		email.Permission = toPermission(account.GetEmails()[i].GetPermission())

		a.R.Email[i] = email
	}

	a.R.My_address = make([]MyAddress, len(account.GetMyAddresses()))
	for i := range account.GetMyAddresses() {
		var address MyAddress

		address.ID = account.GetMyAddresses()[i].GetID()
		address.Name = account.GetMyAddresses()[i].GetName()
		address.Firstname = account.GetMyAddresses()[i].GetFirstname()
		address.Lastname = account.GetMyAddresses()[i].GetLastname()
		address.Street = account.GetMyAddresses()[i].GetStreet()
		address.Apartment = account.GetMyAddresses()[i].GetApartment()
		address.Zip = account.GetMyAddresses()[i].GetZIP()
		address.City = CityInAddress{
			ID:          Int32ToID(account.GetMyAddresses()[i].GetLocation().GetCity().GetId()),
			City:        account.GetMyAddresses()[i].GetLocation().GetCity().GetTitle(),
			Subdivision: account.GetMyAddresses()[i].GetLocation().GetCity().GetSubdivision(),
		}
		address.Country_id = account.GetMyAddresses()[i].GetLocation().GetCountry().GetId()
		address.Primary = account.GetMyAddresses()[i].GetIsPrimary()

		a.R.My_address[i] = address
	}

	a.R.OtherAddress = make([]OtherAddress, len(account.GetOtherAddresses()))
	for i := range account.GetOtherAddresses() {
		var address OtherAddress

		address.ID = account.GetOtherAddresses()[i].GetID()
		address.Fisrtname = account.GetOtherAddresses()[i].GetFirstname()
		address.Lastname = account.GetOtherAddresses()[i].GetLastname()
		address.Name = account.GetOtherAddresses()[i].GetName()
		address.Street = account.GetOtherAddresses()[i].GetStreet()
		address.City = CityInAddress{
			ID:          Int32ToID(account.GetOtherAddresses()[i].GetLocation().GetCity().GetId()),
			City:        account.GetOtherAddresses()[i].GetLocation().GetCity().GetTitle(),
			Subdivision: account.GetOtherAddresses()[i].GetLocation().GetCity().GetSubdivision(),
		}
		address.Appartment = account.GetOtherAddresses()[i].GetApartment()
		address.Zip = account.GetOtherAddresses()[i].GetZIP()
		address.Country_id = account.GetOtherAddresses()[i].GetLocation().GetCountry().GetId()

		a.R.OtherAddress[i] = address
	}

	//----
	// a.R.Sessions = make([]Sessions, len(account.GetSessions()))
	// for i := range account.GetSessions() {
	// 	var s Sessions
	//
	// 	s.ID = account.GetSessions()[i].GetID()
	// 	s.Browser = account.GetSessions()[i].GetBrowser()
	// 	s.Browser_version = account.GetSessions()[i].GetBrowserVersion()
	// 	s.Os = account.GetSessions()[i].GetOS()
	// 	s.Os_version = account.GetSessions()[i].GetOSVersion()
	// 	s.City = strconv.Itoa(int(account.GetSessions()[i].GetCity()))
	// 	s.Country_id = account.GetSessions()[i].GetCountryID()
	// 	s.Device_type = account.GetSessions()[i].GetDeviceType()
	// 	s.Last_activity_time = account.GetSessions()[i].GetLastActivityTime()
	// 	s.Current_session = account.GetSessions()[i].GetCurrentSession()
	//
	// 	a.R.Sessions[i] = s
	// }

	return &a, nil
}

func (_ *Resolver) GetNotificationSettings(ctx context.Context) (*NotificationSettingsResolver, error) {
	settings, err := notifications.GetSettings(ctx, &notificationsRPC.Empty{})
	if err != nil {
		return nil, err
	}

	return &NotificationSettingsResolver{
		R: &NotificationSettings{
			Approved_connection:    settings.GetApprovedConnection(),
			New_connection:         settings.GetNewConnection(),
			New_endorsement:        settings.GetNewConnection(),
			New_follow:             settings.GetNewFollow(),
			New_recommendation:     settings.GetNewRecommendation(),
			Recommendation_request: settings.GetRecommendationRequest(),
			Job_invitation:         settings.GetNewJobInvitation(),
		},
	}, nil
}

func (_ *Resolver) ChangeNotificationsSetting(ctx context.Context, input ChangeNotificationsSettingRequest) (*SuccessResolver, error) {
	_, err := notifications.ChangeSettings(ctx, &notificationsRPC.ChangeSettingsRequest{
		Property: toPropertyOption(input.Property),
		Value:    input.Value,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) GetNotifications(ctx context.Context, input GetNotificationsRequest) (*NotificationsListResolver, error) {
	var first string = "10"
	var after string = "0"

	if input.Pagination.First != nil {
		first = strconv.Itoa(int(*input.Pagination.First))
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	nots, err := notifications.GetNotifications(ctx, &notificationsRPC.Pagination{
		First: first,
		After: after,
	})
	if err != nil {
		return nil, err
	}

	// log.Printf("Got %#v\n", nots.GetAmount())
	// for i := range nots.GetNotifications() {
	// 	log.Printf("Recived Notification %d: %v\n", i, nots.GetNotifications()[i].GetNotification())
	// }

	notsJSON := make([]string, 0, len(nots.GetNotifications()))
	for _, n := range nots.GetNotifications() {
		// 	k, err := json.Marshal(nots.GetNotifications()[i].GetNotification())
		// 	if err != nil {
		// 		log.Println(err)
		// 	} else {
		// 		notsJSON = append(notsJSON, string(k))
		// }

		notsJSON = append(notsJSON, n.GetNotification())
	}

	return &NotificationsListResolver{R: &NotificationsList{
		Amount_not_seen:   nots.GetAmount(),
		Notification_json: notsJSON,
	}}, nil
}

func (_ *Resolver) GetUnseenNotifications(ctx context.Context, input GetNotificationsRequest) (*NotificationsListResolver, error) {
	var first string = "10"
	var after string = "0"

	if input.Pagination.First != nil {
		first = strconv.Itoa(int(*input.Pagination.First))
	}
	if input.Pagination.After != nil {
		after = *input.Pagination.After
	}

	nots, err := notifications.GetNotifications(ctx, &notificationsRPC.Pagination{
		First: first,
		After: after,
	})
	if err != nil {
		return nil, err
	}

	// log.Printf("Got %#v\n", nots.GetAmount())
	// for i := range nots.GetNotifications() {
	// 	log.Printf("Recived Notification %d: %v\n", i, nots.GetNotifications()[i].GetNotification())
	// }

	notsJSON := make([]string, 0, len(nots.GetNotifications()))
	for _, n := range nots.GetNotifications() {
		// 	k, err := json.Marshal(nots.GetNotifications()[i].GetNotification())
		// 	if err != nil {
		// 		log.Println(err)
		// 	} else {
		// 		notsJSON = append(notsJSON, string(k))
		// }

		notsJSON = append(notsJSON, n.GetNotification())
	}

	return &NotificationsListResolver{R: &NotificationsList{
		Amount_not_seen:   nots.GetAmount(),
		Notification_json: notsJSON,
	}}, nil
}

func (_ *Resolver) ChangeFirstName(ctx context.Context, input ChangeFirstNameRequest) (*SuccessResolver, error) {
	_, err := user.ChangeFirstName(ctx, &userRPC.FirstName{FirstName: input.Name})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeLastName(ctx context.Context, input ChangeLastNameRequest) (*SuccessResolver, error) {
	_, err := user.ChangeLastname(ctx, &userRPC.Lastname{Lastname: input.Lastname})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangePatronymic(ctx context.Context, input ChangePatronymicRequest) (*SuccessResolver, error) {
	var isPatronymicNull bool

	if input.Patronymic.Patronymic == nil {
		isPatronymicNull = true
	}

	_, err := user.ChangePatronymic(ctx, &userRPC.Patronymic{
		Patronymic:       NullToString(input.Patronymic.Patronymic),
		Permission:       NullPermissionInputToRPC(input.Patronymic.Permission),
		IsPatronymicNull: isPatronymicNull,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeNickname(ctx context.Context, input ChangeNicknameRequest) (*SuccessResolver, error) {
	var isNicknameNull bool

	if input.Nickname.Nickname == nil {
		isNicknameNull = true
	}

	_, err := user.ChangeNickname(ctx, &userRPC.Nickname{
		Nickname:       NullToString(input.Nickname.Nickname),
		Permission:     NullPermissionInputToRPC(input.Nickname.Permission),
		IsNicknameNull: isNicknameNull,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeMiddlename(ctx context.Context, input ChangeMiddlenameRequest) (*SuccessResolver, error) {
	var isMiddlenameNull bool

	if input.Middlename.Middlename == nil {
		isMiddlenameNull = true
	}
	_, err := user.ChangeMiddleName(ctx, &userRPC.Middlename{
		Middlename:       NullToString(input.Middlename.Middlename),
		Permission:       NullPermissionInputToRPC(input.Middlename.Permission),
		IsMiddlenameNull: isMiddlenameNull,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeNativeName(ctx context.Context, input ChangeNativeNameRequest) (*SuccessResolver, error) {
	var isNameNull bool

	if input.Name.Name == nil {
		isNameNull = true
	}

	_, err := user.ChangeNameOnNativeLanguage(ctx, &userRPC.NativeName{
		Name:       NullToString(input.Name.Name),
		Permission: NullPermissionInputToRPC(input.Name.Permission),
		LanguageID: NullToString(input.Name.Language),
		IsNameNull: isNameNull,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeBithday(ctx context.Context, input ChangeBithdayRequest) (*SuccessResolver, error) {
	_, err := user.ChangeBirthday(ctx, &userRPC.Birthday{
		Birthday:   NullToString(input.Birthday.Date),
		Permission: NullPermissionInputToRPC(input.Birthday.Permission),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeGender(ctx context.Context, input ChangeGenderRequest) (*SuccessResolver, error) {
	var isGenderNull bool

	if input.Gender.Gender == nil {
		isGenderNull = true
	}

	_, err := user.ChangeGender(ctx, &userRPC.Gender{
		Gender: func(g *string) userRPC.GenderValue {
			var val userRPC.GenderValue
			if g == nil {
				val = userRPC.GenderValue_FEMALE
				return val
			}
			if *g == "male" {
				return userRPC.GenderValue_MALE
			}
			val = userRPC.GenderValue_FEMALE
			return val
		}(input.Gender.Gender),
		Permission:   NullPermissionInputToRPC(input.Gender.Permission),
		IsGenderNull: isGenderNull,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddEmail(ctx context.Context, input AddEmailRequest) (*SuccessResolver, error) {
	id, err := user.AddEmail(ctx, &userRPC.Email{
		Email:      input.Email.Email,
		Permission: PermissionInputToRPC(input.Email.Permission),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		Success: true,
		ID:      id.GetID(),
	}}, nil
}

func (_ *Resolver) RemoveEmail(ctx context.Context, input RemoveEmailRequest) (*SuccessResolver, error) {
	_, err := user.RemoveEmail(ctx, &userRPC.Email{
		Id: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeEmail(ctx context.Context, input ChangeEmailRequest) (*SuccessResolver, error) {
	_, err := user.ChangeEmail(ctx, &userRPC.Email{
		Id:         input.Changes.ID,
		Permission: NullPermissionInputToRPC(input.Changes.Permission),
		IsPrimary:  NullBoolToBool(input.Changes.Primary),
	})

	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddPhone(ctx context.Context, input AddPhoneRequest) (*SuccessResolver, error) {
	id, err := user.AddPhone(ctx, &userRPC.Phone{
		CountryCode: &userRPC.CountryCode{
			ID: uint32(input.Phone.Country_code_id),
		},
		Number:     input.Phone.Number,
		Permission: PermissionInputToRPC(input.Phone.Permission),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) ChangePhone(ctx context.Context, input ChangePhoneRequest) (*SuccessResolver, error) {
	_, err := user.ChangePhone(ctx, &userRPC.Phone{
		Id:         input.Changes.ID,
		Permission: NullPermissionInputToRPC(input.Changes.Permission),
		IsPrimary:  NullBoolToBool(input.Changes.Primary),
	})

	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemovePhone(ctx context.Context, input RemovePhoneRequest) (*SuccessResolver, error) {
	_, err := user.RemovePhone(ctx, &userRPC.Phone{
		Id: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddMyAddress(ctx context.Context, input AddMyAddressRequest) (*SuccessResolver, error) {
	id, err := user.AddMyAddress(ctx, &userRPC.MyAddress{
		Name:      input.Address.Name,
		Firstname: input.Address.Firstname,
		Lastname:  input.Address.Lastname,
		Apartment: input.Address.Apartment,
		Street:    input.Address.Street,
		ZIP:       input.Address.Zip,
		Location: &userRPC.Location{
			Country: &userRPC.Country{
				Id: input.Address.Location.Country_id,
			},
			City: &userRPC.City{
				Id:          NullIDToInt32(input.Address.Location.City.ID),
				Title:       NullToString(input.Address.Location.City.City),
				Subdivision: NullToString(input.Address.Location.City.Subdivision),
			},
		},
		// IsPrimary: input.Address.,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) ChangeMyAddress(ctx context.Context, input ChangeMyAddressRequest) (*SuccessResolver, error) {
	var address userRPC.MyAddress

	address.ID = input.ID
	address.Name = NullToString(input.Address.Name)
	address.Firstname = NullToString(input.Address.Firstname)
	address.Lastname = NullToString(input.Address.Lastname)
	address.Apartment = NullToString(input.Address.Apartment)
	address.Street = NullToString(input.Address.Street)
	address.ZIP = NullToString(input.Address.Zip)

	if input.Address.Location != nil {
		address.Location = &userRPC.Location{}

		address.Location.City = &userRPC.City{
			Id:          NullIDToInt32(input.Address.Location.City.ID),
			Title:       NullToString(input.Address.Location.City.City),
			Subdivision: NullToString(input.Address.Location.City.Subdivision),
		}

		address.Location.Country = &userRPC.Country{
			Id: input.Address.Location.Country_id,
		}
	}

	address.IsPrimary = NullBoolToBool(input.Address.Primary)

	_, err := user.ChangeMyAddress(ctx, &address)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveMyAddress(ctx context.Context, input RemoveMyAddressRequest) (*SuccessResolver, error) {
	_, err := user.RemoveMyAddress(ctx, &userRPC.MyAddress{
		ID: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddOtherAddress(ctx context.Context, input AddMyAddressRequest) (*SuccessResolver, error) {
	id, err := user.AddOtherAddress(ctx, &userRPC.OtherAddress{
		Name:      input.Address.Name,
		Firstname: input.Address.Firstname,
		Lastname:  input.Address.Lastname,
		Apartment: input.Address.Apartment,
		Street:    input.Address.Street,
		ZIP:       input.Address.Zip,
		Location: &userRPC.Location{
			Country: &userRPC.Country{
				Id: input.Address.Location.Country_id,
			},
			City: &userRPC.City{
				Id:          NullIDToInt32(input.Address.Location.City.ID),
				Title:       NullToString(input.Address.Location.City.City),
				Subdivision: NullToString(input.Address.Location.City.Subdivision),
			},
		},
		// IsPrimary: input.Address.,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{
		ID:      id.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) ChangeOtherAddress(ctx context.Context, input ChangeMyAddressRequest) (*SuccessResolver, error) {
	var address userRPC.OtherAddress

	address.ID = input.ID
	address.Name = NullToString(input.Address.Name)
	address.Firstname = NullToString(input.Address.Firstname)
	address.Lastname = NullToString(input.Address.Lastname)
	address.Apartment = NullToString(input.Address.Apartment)
	address.Street = NullToString(input.Address.Street)
	address.ZIP = NullToString(input.Address.Zip)

	if input.Address.Location != nil {
		address.Location = &userRPC.Location{}

		address.Location.City = &userRPC.City{
			Id:          NullIDToInt32(input.Address.Location.City.ID),
			Title:       NullToString(input.Address.Location.City.City),
			Subdivision: NullToString(input.Address.Location.City.Subdivision),
		}

		address.Location.Country = &userRPC.Country{
			Id: input.Address.Location.Country_id,
		}
	}

	_, err := user.ChangeOtherAddress(ctx, &address)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveOtherAddress(ctx context.Context, input RemoveMyAddressRequest) (*SuccessResolver, error) {
	_, err := user.RemoveOtherAddress(ctx, &userRPC.OtherAddress{
		ID: input.ID,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeUILanguage(ctx context.Context, input ChangeUILanguageRequest) (*SuccessResolver, error) {
	_, err := user.ChangeUILanguage(ctx, &userRPC.Language{
		Language: input.Language,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangePrivacy(ctx context.Context, input ChangePrivacyRequest) (*SuccessResolver, error) {
	_, err := user.ChangePrivacy(ctx, &userRPC.ChangePrivacyRequest{
		Privacy: func(s string) userRPC.PrivacySettings {
			switch s {
			case "find_by_email":
				return userRPC.PrivacySettings_find_by_email
			case "find_by_phone":
				return userRPC.PrivacySettings_find_by_phone
			case "active_status":
				return userRPC.PrivacySettings_active_status
			case "sharing_edits":
				return userRPC.PrivacySettings_sharing_edits
			case "profile_pictures":
				return userRPC.PrivacySettings_profile_pictures
			case "my_connections":
				return userRPC.PrivacySettings_my_connections
			}
			return userRPC.PrivacySettings_my_connections
		}(input.Privacy),

		Permission: PermissionToRPC(input.Permission),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangePassword(ctx context.Context, input ChangePasswordRequest) (*SuccessResolver, error) {
	_, err := user.ChangePassword(ctx, &userRPC.ChangePasswordRequest{
		OldPassword: input.Old_password,
		NewPassword: input.New_password,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) Init2FA(ctx context.Context) (*TwoFAResponseResolver, error) {
	data, err := user.Init2FA(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}
	return &TwoFAResponseResolver{
		R: &TwoFAResponse{
			Qr_code: data.GetQR(),
			Key:     data.GetKey(),
		},
	}, nil
}

func (_ *Resolver) Enable2FA(ctx context.Context, input Enable2FARequest) (*SuccessResolver, error) {
	_, err := user.Enable2FA(ctx, &userRPC.TwoFACode{Code: input.Code})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) Disable2FA(ctx context.Context, input Disable2FARequest) (*SuccessResolver, error) {
	_, err := user.Disable2FA(ctx, &userRPC.TwoFACode{Code: input.Code})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) DeactivateUserAccount(ctx context.Context, input DeactivateUserAccountRequest) (*SuccessResolver, error) {
	_, err := user.DeactivateAccount(ctx, &userRPC.CheckPasswordRequest{
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) CheckPassword(ctx context.Context, input CheckPasswordRequest) (*SuccessResolver, error) {
	_, err := user.CheckPassword(ctx, &userRPC.CheckPasswordRequest{
		Password: input.Password,
	})
	if err != nil {
		// return nil, err
		return &SuccessResolver{R: &Success{Success: false}}, nil
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

// func (_ *Resolver) ChangeNotification(ctx context.Context, input ChangeNotificationRequest) (*SuccessResolver, error) {
// 	_, err := user.ChangeNotification(ctx, &userRPC.ChangeNotificationRequest{
// 		Notification: func(s string) userRPC.Notification {
// 			switch s {
// 			case "accept_invitation":
// 				return userRPC.Notification_accept_invitation
// 			case "new_followers":
// 				return userRPC.Notification_new_followers
// 			case "new_chat_message":
// 				return userRPC.Notification_new_chat_message
// 			case "birthdays":
// 				return userRPC.Notification_birthdays
// 			case "endorsements":
// 				return userRPC.Notification_endorsements
// 			case "email_updates":
// 				return userRPC.Notification_email_updates
// 			case "job_changes_in_network":
// 				return userRPC.Notification_job_changes_in_network
// 			case "import_contacts_joined":
// 				return userRPC.Notification_import_contacts_joined
// 			case "job_recommendations":
// 				return userRPC.Notification_job_recommendations
// 			}
// 			return userRPC.Notification_connection_request
// 		}(input.Notification),
//
// 		Value: input.Value,
// 	})
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &SuccessResolver{R: &Success{Success: true}}, nil
// }

func toPermission(p *userRPC.Permission) Permission {
	var per Permission
	switch p.GetType() {
	case userRPC.PermissionType_ME:
		per.Type = "me"
	case userRPC.PermissionType_MEMBERS:
		per.Type = "members"
	case userRPC.PermissionType_MY_CONNECTIONS:
		per.Type = "my_connections"
	default:
		per.Type = "members"
	}
	return per
}

func toPermissionType(t userRPC.PermissionType) string {
	switch t {
	case userRPC.PermissionType_ME:
		return "me"
	case userRPC.PermissionType_MEMBERS:
		return "members"
	case userRPC.PermissionType_MY_CONNECTIONS:
		return "my_connections"
	}
	return "members"
}

func toPropertyOption(s string) notificationsRPC.ChangeSettingsRequest_PropertyOption {
	switch s {
	case "new_endorsement":
		return notificationsRPC.ChangeSettingsRequest_NewEndorsement
	case "new_follow":
		return notificationsRPC.ChangeSettingsRequest_NewFollow
	case "new_connection":
		return notificationsRPC.ChangeSettingsRequest_NewConnection
	case "approved_connection":
		return notificationsRPC.ChangeSettingsRequest_ApprovedConnection
	case "recommendation_request":
		return notificationsRPC.ChangeSettingsRequest_RecommendationRequest
	case "new_recommendation":
		return notificationsRPC.ChangeSettingsRequest_NewRecommendation
	case "job_invitation":
		return notificationsRPC.ChangeSettingsRequest_NewJobInvitation
	}

	return notificationsRPC.ChangeSettingsRequest_UnknownProperty
}
