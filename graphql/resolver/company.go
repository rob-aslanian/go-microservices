package resolver

import (
	"context"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/notificationsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func (_ *Resolver) GetOriginCompanyAvatar(ctx context.Context, input GetOriginCompanyAvatarRequest) (string, error) {
	url, err := company.GetOriginAvatar(ctx, &companyRPC.ID{
		ID: input.Company_id,
	})
	if err != nil {
		return "", err
	}

	return url.GetURL(), nil
}

func (_ *Resolver) GetOriginCompanyCover(ctx context.Context, input GetOriginCompanyCoverRequest) (string, error) {
	url, err := company.GetOriginCover(ctx, &companyRPC.ID{
		ID: input.Company_id,
	})
	if err != nil {
		return "", err
	}

	return url.GetURL(), nil
}

func (_ *Resolver) CheckIfURLForCompanyIsTaken(ctx context.Context, input CheckIfURLForCompanyIsTakenRequest) (bool, error) {
	res, err := company.CheckIfURLForCompanyIsTaken(ctx, &companyRPC.URL{
		URL: input.Url,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

// DONE (ALMOST)
func (_ *Resolver) RegisterCompany(ctx context.Context, in RegisterCompanyRequest) (*RegisterCompanyResponseResolver, error) {

	cityID, _ := strconv.Atoi(in.Input.City_id)

	response, err := company.RegisterCompany(ctx, &companyRPC.RegisterCompanyRequest{
		Apartment: NullToString(in.Input.Apartment),
		CityId:    int32(cityID),
		Email:     in.Input.Email,
		Industry: &companyRPC.Industry{
			Main: in.Input.Industry,
			// Subs: in.Input.Industry.Subindustries,
		},
		Websites: websitesToRPC(in.Input.Websites),
		Name:     in.Input.Name,
		Phone: &companyRPC.Phone{
			CountryCode: &companyRPC.CountryCode{
				Id: in.Input.Phone.Country_code_id,
			},
			Number: in.Input.Phone.Number,
		},
		// FoundationDate: in.Input.Foundation_date,
		// State:         in.Input.State,
		StreetAddress: in.Input.Address,
		Type:          stringToCompanyTypeRPC(in.Input.Type),
		URL:           in.Input.Url,
		ZipCode:       in.Input.Zip,
		VAT:           NullToString(in.Input.Vat),
		InvitedBy:     NullToString(in.Input.Invited_by),
	})

	if err != nil {
		return nil, err
	}

	return &RegisterCompanyResponseResolver{
		R: &RegisterCompanyResponse{
			ID:      response.GetId(),
			Success: true,
			Url:     response.GetURL(),
		},
	}, nil
}

func websitesToRPC(data *[]string) []string {
	if data == nil {
		return nil
	}

	websites := make([]string, 0, len(*data))

	for _, webiste := range *data {
		websites = append(websites, webiste)
	}

	return websites
}

func (_ *Resolver) DeactivateCompany(ctx context.Context, in DeactivateCompanyRequest) (*SuccessResolver, error) {
	_, err := company.DeactivateCompany(ctx, &companyRPC.DeactivateCompanyRequest{
		Id:       in.Company_id,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) GetCompanyAccount(ctx context.Context, in GetCompanyAccountRequest) (*CompanyAccountResolver, error) {
	response, err := company.GetCompanyAccount(ctx, &companyRPC.ID{
		ID: in.Company_id,
	})
	if err != nil {
		return nil, err
	}

	return &CompanyAccountResolver{
		R: companyAccountRPCToCompanyAccount(response),
	}, nil
}

func (_ *Resolver) GetCompanyNotificationSettings(ctx context.Context, in GetCompanyNotificationSettingsRequest) (*CompanyNotificationSettingsResolver, error) {
	settings, err := notifications.GetCompanySettings(ctx, &notificationsRPC.ID{
		ID: in.Company_id,
	})
	if err != nil {
		return nil, err
	}

	return &CompanyNotificationSettingsResolver{
		R: &CompanyNotificationSettings{
			New_follow:    settings.GetNewFollow(),
			New_review:    settings.GetNewReview(),
			New_applicant: settings.GetNewApplicant(),
		},
	}, nil
}

func (_ *Resolver) GetCompanyNotifications(ctx context.Context, in GetCompanyNotificationsRequest) (*NotificationsListResolver, error) {
	var first string = "10"
	var after string = "0"

	if in.Pagination.First != nil {
		first = strconv.Itoa(int(*in.Pagination.First))
	}
	if in.Pagination.After != nil {
		after = *in.Pagination.After
	}

	nots, err := notifications.GetCompanyNotifications(ctx, &notificationsRPC.PaginationWithID{
		ID:    in.Company_id,
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

func (_ *Resolver) ChangeCompanyNotificationsSetting(ctx context.Context, in ChangeCompanyNotificationsSettingRequest) (*SuccessResolver, error) {
	_, err := notifications.ChangeCompanySettings(ctx, &notificationsRPC.ChangeCompanySettingsRequest{
		CompanyID: in.Company_id,
		Property:  toCompanyPropertyOption(in.Property),
		Value:     in.Value,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) MarkNotificationAsSeenForCompany(ctx context.Context, in MarkNotificationAsSeenForCompanyRequest) (*SuccessResolver, error) {
	_, err := notifications.MarkAsSeenForCompany(ctx, &notificationsRPC.IDWithIDs{
		ID:  in.Company_id,
		IDs: in.Ids,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveNotificationForCompany(ctx context.Context, in MarkNotificationAsSeenForCompanyRequest) (*SuccessResolver, error) {
	_, err := notifications.RemoveNotificationForCompany(ctx, &notificationsRPC.IDWithIDs{
		ID:  in.Company_id,
		IDs: in.Ids,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyName(ctx context.Context, in ChangeCompanyNameRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyName(ctx, &companyRPC.ChangeCompanyNameRequest{
		Id:   in.Company_id,
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyUrl(ctx context.Context, in ChangeCompanyUrlRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyURL(ctx, &companyRPC.ChangeCompanyUrlRequest{
		Id:  in.Company_id,
		Url: in.Url,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyFoundationDate(ctx context.Context, in ChangeCompanyFoundationDateRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyFoundationDate(ctx, &companyRPC.ChangeCompanyFoundationDateRequest{
		Id:             in.Company_id,
		FoundationDate: in.Foundation_date,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

// TODO:
func (_ *Resolver) ChangeCompanyIndustry(ctx context.Context, in ChangeCompanyIndustryRequest) (*SuccessResolver, error) {

	_, err := company.ChangeCompanyIndustry(ctx, &companyRPC.ChangeCompanyIndustryRequest{
		ID: in.Company_id,
		Industry: &companyRPC.Industry{
			Main: in.Input.ID,
			Subs: NullStringArrayToStringArray(in.Input.Subindustries),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil

	// 	comp := companyRPC.ChangeCompanyIndustryRequest{
	// 		Id: in.Company_id,
	// 		Industry: &companyRPC.Industry{
	//
	// 			Id: func() companyRPC.IndustryId {
	// 				if in.Input.ID == "industry_information_technology" {
	// 					return companyRPC.IndustryId_INDUSTRY_INFORMATION_TECHNOLOGY
	// 				}
	// 				return companyRPC.IndustryId_INDUSTRY_UNKNOWN
	// 			}(),
	//
	// 			Subindustries: make([]*companyRPC.Industry_Subindustry, 0),
	// 		},
	// 	}
	//
	// 	if in.Input.Subindustries != nil {
	// 		var res = []*companyRPC.Industry_Subindustry{}
	// 		for _, v := range *in.Input.Subindustries {
	// 			if v == "subindustry_web_development" {
	// 				res = append(res, &companyRPC.Industry_Subindustry{
	// 					Id: companyRPC.SubindustryId_SUBINDUSTRY_WEB_DEVEVELOPMENT,
	// 				})
	// 			} else if v == "subindustry_it_management" {
	// 				res = append(res, &companyRPC.Industry_Subindustry{
	// 					Id: companyRPC.SubindustryId_SUBINDUSTRY_IT_MANAGEMENT,
	// 				})
	// 			} else {
	// 				res = append(res, &companyRPC.Industry_Subindustry{
	// 					Id: companyRPC.SubindustryId_SUBINDUSTRY_UNKNOWN,
	// 				})
	// 			}
	// 		}
	//
	// 		comp.GetIndustry().Subindustries = res
	// 	}
	//
	// 	_, err := company.ChangeCompanyIndustry(ctx, &comp)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyType(ctx context.Context, in ChangeCompanyTypeRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyType(ctx, &companyRPC.ChangeCompanyTypeRequest{
		Id:   in.Company_id,
		Type: stringToCompanyTypeRPC(in.Type),
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

func (_ *Resolver) ChangeCompanySize(ctx context.Context, in ChangeCompanySizeRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanySize(ctx, &companyRPC.ChangeCompanySizeRequest{
		Id:   in.Company_id,
		Size: stringToCompanySizeRPC(in.Size),
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

func (_ *Resolver) AddCompanyEmail(ctx context.Context, in AddCompanyEmailRequest) (*SuccessResolver, error) {
	response, err := company.AddCompanyEmail(ctx, &companyRPC.AddCompanyEmailRequest{
		ID: in.Company_id,
		Email: &companyRPC.Email{
			Email: in.Input.Email,
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

	// 	resp, err := company.AddCompanyEmail(ctx, &companyRPC.AddCompanyEmailRequest{
	// 		Id: in.Company_id,
	// 		Email: &companyRPC.Email{
	// 			Email: in.Input.Email,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &SuccessResolver{R: &Success{Success: true, ID: strconv.Itoa(int(resp.EmailId))}}, nil
}

func (_ *Resolver) DeleteCompanyEmail(ctx context.Context, in DeleteCompanyEmailRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyEmail(ctx, &companyRPC.DeleteCompanyEmailRequest{
		Id:      in.Company_id,
		EmailId: in.ID,
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

func (_ *Resolver) ChangeCompanyEmail(ctx context.Context, in ChangeCompanyEmailRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyEmail(ctx, &companyRPC.ChangeCompanyEmailRequest{
		Id:             in.Company_id,
		EmailId:        in.Changes.ID,
		EmailIsPrimary: NullBoolToBool(in.Changes.Is_primary),
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddCompanyPhone(ctx context.Context, in AddCompanyPhoneRequest) (*SuccessResolver, error) {
	response, err := company.AddCompanyPhone(ctx, &companyRPC.AddCompanyPhoneRequest{
		Id: in.Company_id,
		Phone: &companyRPC.Phone{
			Number: in.Input.Number,
			CountryCode: &companyRPC.CountryCode{
				Id: in.Input.Country_code_id,
			},
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

func (_ *Resolver) DeleteCompanyPhone(ctx context.Context, in DeleteCompanyPhoneRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyPhone(ctx, &companyRPC.DeleteCompanyPhoneRequest{
		Id:      in.Company_id,
		PhoneId: in.ID,
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

func (_ *Resolver) ChangeCompanyPhone(ctx context.Context, in ChangeCompanyPhoneRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyPhone(ctx, &companyRPC.ChangeCompanyPhoneRequest{
		Id:             in.Company_id,
		PhoneId:        in.Changes.ID,
		PhoneIsPrimary: NullBoolToBool(in.Changes.Is_primary),
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

// TODO:
func (_ *Resolver) AddCompanyAddress(ctx context.Context, in AddCompanyAddressRequest) (*SuccessResolver, error) {
	address := companyRPC.AddCompanyAddressRequest{
		Id: in.Company_id,
		Address: &companyRPC.Address{
			Name:          in.Input.Name,
			Apartment:     in.Input.Apartment,
			ZipCode:       in.Input.Zip_code,
			StreetAddress: in.Input.Street_address,
			Phones:        make([]*companyRPC.Phone, 0, len(in.Input.Phones)),
			Location: &companyRPC.Location{
				City: &companyRPC.City{
					Id: in.Input.City_id,
				},
				Country: &companyRPC.Country{},
			},
			BusinessHours: make([]*companyRPC.BusinessHoursItem, 0, len(in.Input.Business_hours)),
			IsPrimary:     in.Input.Primary,
		},
	}

	if gp := in.Input.Geo_pos; gp != nil {
		address.Address.GeoPos = &companyRPC.GeoPos{
			Lantitude: in.Input.Geo_pos.Lantitude,
			Longitude: in.Input.Geo_pos.Longitude,
		}
	}

	for i := range in.Input.Phones {
		address.Address.Phones = append(address.Address.Phones, &companyRPC.Phone{
			CountryCode: &companyRPC.CountryCode{
				Id: in.Input.Phones[i].Country_code_id,
			},
			Number: in.Input.Phones[i].Number,
		})
	}

	for i := range in.Input.Business_hours {
		bh := in.Input.Business_hours[i]
		address.Address.BusinessHours = append(address.Address.BusinessHours, businessHourItemToCompanyBusinessHourRPC(&bh))
	}

	response, err := company.AddCompanyAddress(ctx, &address)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      response.GetID(),
			Success: true,
		},
	}, nil

	// 	_, err := company.AddCompanyAddress(ctx, &companyRPC.AddCompanyAddressRequest{
	// 		Id: in.Company_id,
	// 		Address: &companyRPC.Address{
	// 			Id:            0, //not used
	// 			Name:          in.Input.Name,
	// 			ZipCode:       in.Input.Zip_code,
	// 			Apartment:     in.Input.Apartment,
	// 			StreetAddress: in.Input.Street_address,
	// 			CityId:        in.Input.City_id,
	// 			CountryId:     in.Input.Country_id,
	// 			State:         in.Input.State,
	// 			Phones: func() []*companyRPC.Phone {
	// 				res := []*companyRPC.Phone{}
	// 				for _, v := range in.Input.Phones {
	// 					res = append(res, &companyRPC.Phone{
	// 						//CountryAbbreviation: "",							//TODO
	// 						CountryCode: &companyRPC.CountryCode{
	// 							Id: v.Country_code_id,
	// 							//Code: "",										//TODO
	// 						},
	// 						Number: v.Number,
	// 					})
	// 				}
	// 				return res
	// 			}(),
	// 			BusinessHours: func() []*companyRPC.BusinessHoursItem {
	// 				res := []*companyRPC.BusinessHoursItem{}
	// 				for _, v := range in.Input.Business_hours {
	// 					res = append(res, &companyRPC.BusinessHoursItem{
	// 						WeekDays: v.Week_days,
	// 						HourFrom: v.Hour_from,
	// 						HourTo:   v.Hour_to,
	// 					})
	// 				}
	// 				return res
	// 			}(),
	// 			GeoPos: func() *companyRPC.GeoPos {
	// 				res := &companyRPC.GeoPos{}
	// 				if in.Input.Geo_pos != nil {
	// 					res.Lantitude = float64(in.Input.Geo_pos.Lantitude)
	// 					res.Longitude = float64(in.Input.Geo_pos.Longitude)
	// 				}
	// 				return res
	// 			}(),
	// 			IsPrimary: false, //not used
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) DeleteCompanyAddress(ctx context.Context, in DeleteCompanyAddressRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyAddress(ctx, &companyRPC.DeleteCompanyAddressRequest{
		Id:        in.Company_id,
		AddressId: in.ID,
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

//  TODO:
func (_ *Resolver) ChangeCompanyAddress(ctx context.Context, in ChangeCompanyAddressRequest) (*SuccessResolver, error) {
	address := companyRPC.ChangeCompanyAddressRequest{
		Id: in.Company_id,
		Address: &companyRPC.Address{
			ID:            in.Changes.ID,
			Name:          NullToString(in.Changes.Name),
			Apartment:     NullToString(in.Changes.Apartment),
			ZipCode:       NullToString(in.Changes.Zip_code),
			StreetAddress: NullToString(in.Changes.Street_address),
			// Phones:        make([]*companyRPC.Phone, 0, len(in.Changes.Phones)),
			Location: &companyRPC.Location{
				City: &companyRPC.City{
					Id: NullToString(in.Changes.City_id),
				},
				Country: &companyRPC.Country{},
			},
			IsPrimary: NullBoolToBool(in.Changes.Primary),
			// BusinessHours
		},
	}

	if gp := in.Changes.Geo_pos; gp != nil {
		address.Address.GeoPos = &companyRPC.GeoPos{
			Lantitude: in.Changes.Geo_pos.Lantitude,
			Longitude: in.Changes.Geo_pos.Longitude,
		}
	}

	if phones := in.Changes.Phones; phones != nil {

		address.Address.Phones = make([]*companyRPC.Phone, 0, len(*phones))

		for i := range *phones {
			address.Address.Phones = append(address.Address.Phones, &companyRPC.Phone{
				CountryCode: &companyRPC.CountryCode{
					Id: (*phones)[i].Country_code_id,
				},
				Number: (*phones)[i].Number,
			})
		}

		if in.Changes.Business_hours != nil {
			address.Address.BusinessHours = make([]*companyRPC.BusinessHoursItem, 0, len(*in.Changes.Business_hours))

			for i := range *in.Changes.Business_hours {
				address.Address.BusinessHours = append(address.Address.BusinessHours, &companyRPC.BusinessHoursItem{
					// ID:       (*in.Changes.Business_hours)[i].ID,
					WeekDays: (*in.Changes.Business_hours)[i].Week_days,
					HourFrom: (*in.Changes.Business_hours)[i].Hour_from,
					HourTo:   (*in.Changes.Business_hours)[i].Hour_to,
				})
			}
		}
	}

	_, err := company.ChangeCompanyAddress(ctx, &address)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			// ID:      response.GetID(),
			Success: true,
		},
	}, nil

	// OLD
	// 	_, err := company.ChangeCompanyAddress(ctx, &companyRPC.ChangeCompanyAddressRequest{
	// 		Id: in.Company_id,
	// 		Address: &companyRPC.Address{
	// 			Id:            in.Changes.ID,
	// 			Name:          NullToString(in.Changes.Name),
	// 			ZipCode:       NullToString(in.Changes.Zip_code),
	// 			Apartment:     NullToString(in.Changes.Apartment),
	// 			StreetAddress: NullToString(in.Changes.Street_address),
	// 			CityId:        NullToInt32(in.Changes.City_id),
	// 			CountryId:     NullToString(in.Changes.Country_id),
	// 			State:         NullToString(in.Changes.State),
	// 			Phones: func() []*companyRPC.Phone {
	// 				res := []*companyRPC.Phone{}
	// 				if in.Changes.Phones != nil {
	// 					for _, v := range *in.Changes.Phones {
	// 						res = append(res, &companyRPC.Phone{
	// 							//CountryAbbreviation: "",							//TODO
	// 							CountryCode: &companyRPC.CountryCode{
	// 								Id: v.Country_code_id,
	// 								//Code: "",										//TODO
	// 							},
	// 							Number: v.Number,
	// 						})
	// 					}
	// 				}
	// 				return res
	// 			}(),
	// 			BusinessHours: func() []*companyRPC.BusinessHoursItem {
	// 				res := []*companyRPC.BusinessHoursItem{}
	// 				if in.Changes.Business_hours != nil {
	// 					for _, v := range *in.Changes.Business_hours {
	// 						res = append(res, &companyRPC.BusinessHoursItem{
	// 							WeekDays: v.Week_days,
	// 							HourFrom: v.Hour_from,
	// 							HourTo:   v.Hour_to,
	// 						})
	// 					}
	// 				}
	// 				return res
	// 			}(),
	// 			GeoPos: func() *companyRPC.GeoPos {
	// 				res := &companyRPC.GeoPos{}
	// 				if in.Changes.Geo_pos != nil {
	// 					res.Lantitude = in.Changes.Geo_pos.Lantitude
	// 					res.Longitude = in.Changes.Geo_pos.Longitude
	// 				}
	// 				return res
	// 			}(),
	// 			IsPrimary: NullBoolToBool(in.Changes.Is_primary),
	// 		},
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &SuccessResolver{R: &Success{Success: true}}, nil
	// return nil, nil
}

func (_ *Resolver) AddCompanyWebsite(ctx context.Context, in AddCompanyWebsiteRequest) (*SuccessResolver, error) {
	response, err := company.AddCompanyWebsite(ctx, &companyRPC.AddCompanyWebsiteRequest{
		Id:      in.Company_id,
		Website: in.Website,
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

func (_ *Resolver) DeleteCompanyWebsite(ctx context.Context, in DeleteCompanyWebsiteRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyWebsite(ctx, &companyRPC.DeleteCompanyWebsiteRequest{
		Id:        in.Company_id,
		WebsiteId: in.ID,
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

func (_ *Resolver) ChangeCompanyWebsite(ctx context.Context, in ChangeCompanyWebsiteRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyWebsite(ctx, &companyRPC.ChangeCompanyWebsiteRequest{
		Id:        in.Company_id,
		WebsiteId: in.Changes.ID,
		Website:   NullToString(in.Changes.Website),
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

func (_ *Resolver) ChangeCompanyParking(ctx context.Context, in ChangeCompanyParkingRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyParking(ctx, &companyRPC.ChangeCompanyParkingRequest{
		Id:      in.Company_id,
		Parking: stringToCompanyParkingRPC(in.Parking),
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyBenefits(ctx context.Context, in ChangeCompanyBenefitsRequest) (*SuccessResolver, error) {
	benefits := make([]companyRPC.BenefitEnum, len(*in.Benefits))
	for i, t := range *in.Benefits {
		benefits[i] = companyRPC.BenefitEnum(companyRPC.BenefitEnum_value[t])
	}
	_, err := company.ChangeCompanyBenefits(ctx, &companyRPC.ChangeCompanyBenefitsRequest{
		ID:             in.Company_id,
		CompanyBenefit: benefits,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

// does id should be returned?
func (_ *Resolver) AddCompanyAdmin(ctx context.Context, in AddCompanyAdminRequest) (*SuccessResolver, error) {
	_, err := company.AddCompanyAdmin(ctx, &companyRPC.AddCompanyAdminRequest{
		Id:       in.Company_id,
		Password: in.Password,
		UserId:   in.User_id,
		Role:     stringToCompanyAdminRoleRPC(in.Role),
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

func (_ *Resolver) DeleteCompanyAdmin(ctx context.Context, in AddCompanyAdminRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyAdmin(ctx, &companyRPC.DeleteCompanyAdminRequest{
		Id:       in.Company_id,
		Password: in.Password,
		UserId:   in.User_id,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) GetCompanyAdmins(ctx context.Context, in GetCompanyAdminsRequest) ([]CompanyAdminResolver, error) {
	result, err := network.GetCompanyAdmins(ctx, &networkRPC.Company{
		Id: in.Company_id,
	})
	if err != nil {
		return nil, err
	}

	admins := make([]CompanyAdminResolver, 0, len(result.GetList()))

	ids := make([]string, 0, len(result.GetList()))
	for _, r := range result.GetList() {
		ids = append(ids, string(r.GetUser().GetId()))
	}

	if len(ids) > 0 {
		profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
			ID: ids,
		})
		if err != nil {
			return nil, err
		}

		for _, res := range result.GetList() {
			admin := CompanyAdminResolver{
				R: &CompanyAdmin{
					User: ToProfile(ctx, profiles.GetProfiles()[string(res.GetUser().GetId())]),
					Role: networkAdminRoleRPCToString(res.GetLevel()),
				},
			}

			admins = append(admins, admin)
		}
	}

	return admins, nil
}

//Profile management

func (_ *Resolver) GetCompanyProfile(ctx context.Context, in GetCompanyProfileRequest) (*CompanyProfileResolver, error) {
	var lang string

	if in.Lang != nil {
		lang = *in.Lang
	}

	resp, err := company.GetCompanyProfile(ctx, &companyRPC.GetCompanyProfileRequest{
		Url:  in.Url,
		Lang: lang,
	},
	)
	if err != nil {
		return nil, err
	}

	var r CompanyProfile

	if resp != nil {
		r = toCompanyProfile(ctx, *resp.GetProfile())
	}

	return &CompanyProfileResolver{
		R: &r,
	}, nil
}

func (_ *Resolver) GetCompanyProfileByID(ctx context.Context, in GetCompanyProfileByIDRequest) (*CompanyProfileResolver, error) {
	resp, err := company.GetCompanyProfileByID(ctx, &companyRPC.ID{
		ID: in.Company_id,
	},
	)
	if err != nil {
		return nil, err
	}

	var r CompanyProfile

	if resp != nil {
		r = toCompanyProfile(ctx, *resp.GetProfile())
	}

	return &CompanyProfileResolver{
		R: &r,
	}, nil
}

// TODO:
func toCompanyProfile(ctx context.Context, profile companyRPC.Profile) CompanyProfile {
	prof := CompanyProfile{
		ID:              profile.GetId(),
		Avatar:          profile.GetAvatar(),
		Url:             profile.GetURL(),
		Cover:           profile.GetCover(),
		Name:            profile.GetName(),
		Description:     profile.GetDescription(),
		Mission:         profile.GetMission(),
		Type:            companyTypeRPCToString(profile.GetType()),
		Foundation_date: profile.GetFoundationDate(),
		Size:            companySizeRPCToString(profile.GetSize()),
		Parking:         companyParkingRPCToString(profile.GetParking()),
		Business_hours:  make([]BusinessHoursItem, 0, len(profile.GetBusinessHours())),
		// Founders:        make([]CompanyFounder, 0, len(profile.GetFounders())),
		Awards: make([]CompanyAward, 0, len(profile.GetAwards())),
		// Galleries:              make([]GalleryProfile, 0, len(profile.GetGalleries())),
		Milstones:              make([]CompanyMilestone, 0, len(profile.GetMilestones())),
		Products:               make([]CompanyProduct, 0, len(profile.GetProducts())),
		Services:               make([]CompanyService, 0, len(profile.GetServices())),
		Addresses:              make([]CompanyAddress, 0, len(profile.GetAddresses())),
		Websites:               profile.GetWebsites(),
		Email:                  profile.GetEmail().GetEmail(),
		Emails:                 profile.GetEmails(),
		Phones:                 profile.GetPhones(),
		My_role:                companyAdminRoleRPCToString(profile.GetRole()),
		Amount_jobs:            profile.GetAmountOfJobs(),
		Avarage_rating:         float64(profile.GetAvarageRating()),
		Is_about_us_set:        profile.GetWasAboutUsSet(),
		Favorite:               profile.GetIsFavorite(),
		Follow:                 profile.GetIsFollow(),
		Online:                 profile.GetIsOnline(),
		Blocked:                profile.GetIsBlocked(),
		Benefits:               companyBenefitsRPCToStringArray(profile.GetBenefits()),
		Available_translations: profile.GetAvailableTranslations(),
		Current_translation:    profile.GetCurrentTranslation(),
		Career_center:          careerCenterRPCToCareerCenter(profile.GetCareerCenter()),
	}

	if profile.GetIndustry() != nil {
		prof.Industry = CompanyIndustry{
			ID:            profile.GetIndustry().GetMain(),
			Subindustries: profile.GetIndustry().GetSubs(),
		}
	} else {
		prof.Industry = CompanyIndustry{}
	}

	for i := range profile.GetBusinessHours() {

		bh := BusinessHoursItem{
			ID:        profile.GetBusinessHours()[i].GetID(),
			Hour_from: profile.GetBusinessHours()[i].GetHourFrom(),
			Hour_to:   profile.GetBusinessHours()[i].GetHourTo(),
			Week_days: profile.GetBusinessHours()[i].GetWeekDays(), // TODO: change
		}

		prof.Business_hours = append(prof.Business_hours, bh)
	}

	for i := range profile.GetAwards() {
		award := CompanyAward{
			ID:     profile.GetAwards()[i].GetID(),
			Issuer: profile.GetAwards()[i].GetIssuer(),
			Title:  profile.GetAwards()[i].GetTitle(),
			Year:   profile.GetAwards()[i].GetYear(),
			Link:   make([]Link, 0, len(profile.GetAwards()[i].GetLinks())),
			File:   make([]File, 0, len(profile.GetAwards()[i].GetFiles())),
		}

		for _, l := range profile.GetAwards()[i].GetLinks() {
			award.Link = append(award.Link, Link{
				Address: l.GetURL(),
				ID:      l.GetID(),
			})
		}

		for _, f := range profile.GetAwards()[i].GetFiles() {
			award.File = append(award.File, File{
				Address:   f.GetURL(),
				ID:        f.GetID(),
				Mime_type: f.GetMimeType(),
				Name:      f.GetName(),
			})
		}

		prof.Awards = append(prof.Awards, award)
	}

	// for i := range profile.GetGalleries() {
	// 	gallery := GalleryProfile{
	// 		ID:    profile.GetGalleries()[i].GetID(),
	// 		Files: make([]File, 0, len(profile.GetGalleries()[i].GetFile())),
	// 	}

	// 	for _, f := range profile.GetGalleries()[i].GetFile() {
	// 		gallery.Files = append(gallery.Files, File{
	// 			Address:   f.GetURL(),
	// 			ID:        f.GetID(),
	// 			Mime_type: f.GetMimeType(),
	// 			Name:      f.GetName(),
	// 		})
	// 	}
	// }

	for i := range profile.GetMilestones() {
		prof.Milstones = append(prof.Milstones, CompanyMilestone{
			ID:          profile.GetMilestones()[i].GetId(),
			Description: profile.GetMilestones()[i].GetDescription(),
			Image:       profile.GetMilestones()[i].GetImage(),
			Title:       profile.GetMilestones()[i].GetTitle(),
			Year:        profile.GetMilestones()[i].GetYear(),
		})
	}

	for i := range profile.GetProducts() {
		prof.Products = append(prof.Products, CompanyProduct{
			ID:      profile.GetProducts()[i].GetID(),
			Image:   profile.GetProducts()[i].GetImage(),
			Website: profile.GetProducts()[i].GetWebsite(),
			Name:    profile.GetProducts()[i].GetName(),
		})
	}

	for i := range profile.GetServices() {
		prof.Services = append(prof.Services, CompanyService{
			ID:      profile.GetServices()[i].GetID(),
			Image:   profile.GetServices()[i].GetImage(),
			Name:    profile.GetServices()[i].GetName(),
			Website: profile.GetServices()[i].GetWebsite(),
		})
	}

	for i := range profile.GetAddresses() {
		if address := companyAddressRPCToCompanyAddress(profile.GetAddresses()[i]); address != nil {
			prof.Addresses = append(prof.Addresses, *address)
		}
	}

	if profile.GetPhone() != nil {
		prof.Phone = &PhoneProfile{
			Number: profile.GetPhone().GetNumber(),
		}
		if profile.GetPhone().GetCountryCode() != nil {
			prof.Phone.Country_code = profile.GetPhone().GetCountryCode().GetCode()
			// prof.Phone.Country_code = profile.GetPhone().GetNumber()
		}
	} else {
		prof.Phone = &PhoneProfile{}
	}

	prof.Network_info = &NetworkInfoInCompanyProfile{
		Employees:  profile.GetAmountOfEmployees(),
		Followers:  profile.GetAmountOfFollowers(),
		Followings: profile.GetAmountOfFollowings(),
	}

	return prof

	// return CompanyProfile{
	// 	ID:          profile.GetId(),
	// 	Url:         profile.GetURL(),
	// 	Name:        profile.GetName(),
	// 	Description: profile.GetDescription(),
	// 	Mission:     profile.GetMission(),
	//
	// 	Industry: CompanyIndustry{
	// 		// ID: companyRPC.IndustryId_name[int32(profile.GetIndustry().GetID())],
	// 		// Subindustries: func() []string {
	// 		// 	res := []string{}
	// 		// 	for _, v := range profile.GetIndustry().GetSubindustries() {
	// 		// 		res = append(res, companyRPC.SubindustryId_name[int32(v.GetId())])
	// 		// 	}
	// 		// 	return res
	// 		// }(),
	// 	},
	//
	// 	Type:            companyRPC.Type_name[int32(profile.GetType())],
	// 	Foundation_date: profile.GetFoundationDate(),
	// 	Size:            companyRPC.Size_name[int32(profile.GetSize())],
	// 	Parking:         companyRPC.Parking_name[int32(profile.GetParking())],
	// 	Business_hours: func() []BusinessHoursItem {
	// 		res := []BusinessHoursItem{}
	// 		for _, v := range profile.GetBusinessHours() {
	// 			res = append(res, BusinessHoursItem{
	// 				// ID:        v.GetID(),
	// 				Week_days: v.GetWeekDays(),
	// 				Hour_from: v.GetHourFrom(),
	// 				Hour_to:   v.GetHourTo(),
	// 			})
	// 		}
	// 		return res
	// 	}(),
	//
	// 	Founders: func() []CompanyFounder {
	// 		res := []CompanyFounder{}
	// 		for _, v := range profile.GetFounders() {
	// 			res = append(res, CompanyFounder{
	// 				// ID:             v.GetId(),
	// 				Name:           v.GetName(),
	// 				Position_title: v.GetPositionTitle(),
	// 				Avatar:         v.GetAvatar(),
	// 			})
	// 		}
	// 		return res
	// 	}(),
	//
	// 	Awards: func() []CompanyAward {
	// 		res := []CompanyAward{}
	// 		for _, v := range profile.GetAwards() {
	// 			res = append(res, CompanyAward{
	// 				// ID:     v.GetId(),
	// 				Title:  v.GetTitle(),
	// 				Issuer: v.GetIssuer(),
	// 				Year:   v.GetYear(),
	// 			})
	// 		}
	// 		return res
	// 	}(),
	// 	Milstones: func() []CompanyMilestone {
	// 		res := []CompanyMilestone{}
	// 		for _, v := range profile.GetMilestones() {
	// 			res = append(res, CompanyMilestone{
	// 				// ID:          v.GetId(),
	// 				Image:       v.GetImage(),
	// 				Year:        v.GetYear(),
	// 				Title:       v.GetTitle(),
	// 				Description: v.GetDescription(),
	// 			})
	// 		}
	// 		return res
	// 	}(),
	//
	// 	Products: func() []CompanyProduct {
	// 		res := []CompanyProduct{}
	// 		for _, v := range profile.GetProducts() {
	// 			res = append(res, CompanyProduct{
	// 				// ID:      v.GetId(),
	// 				Image:   v.GetImage(),
	// 				Name:    v.GetName(),
	// 				Website: v.GetWebsite(),
	// 			})
	// 		}
	// 		return res
	// 	}(),
	// 	Services: func() []CompanyService {
	// 		res := []CompanyService{}
	// 		for _, v := range profile.GetServices() {
	// 			res = append(res, CompanyService{
	// 				// ID:      v.GetId(),
	// 				Image:   v.GetImage(),
	// 				Name:    v.GetName(),
	// 				Website: v.GetWebsite(),
	// 			})
	// 		}
	// 		return res
	// 	}(),
	//
	// 	Addresses: func() []CompanyAddress {
	// 		res := []CompanyAddress{}
	// 		for _, v := range profile.GetAddresses() {
	// 			res = append(res, CompanyAddress{
	// 				// ID:             v.GetId(),
	// 				Name:      v.GetName(),
	// 				Zip_code:  v.GetZipCode(),
	// 				Apartment: v.GetApartment(),
	// 				// State:          v.GetState(),
	// 				Street_address: v.GetStreetAddress(),
	// 				// City_id:        v.GetCityId(),
	// 				// Country_id:     v.GetCountryId(),
	// 				Phones: func() []CompanyPhone {
	// 					res := []CompanyPhone{}
	// 					for _, v := range v.GetPhones() {
	// 						res = append(res, CompanyPhone{
	// 							// ID:              v.GetId(),
	// 							Country_code_id: v.GetCountryCode().GetId(),
	// 							Number:          v.GetNumber(),
	// 							Is_primary:      v.GetIsPrimary(),
	// 						})
	// 					}
	// 					return res
	// 				}(),
	// 				Business_hours: func() []BusinessHoursItem {
	// 					res := []BusinessHoursItem{}
	// 					for _, v := range v.GetBusinessHours() {
	// 						res = append(res, BusinessHoursItem{
	// 							// ID:        v.GetId(),
	// 							Week_days: v.GetWeekDays(),
	// 							Hour_from: v.GetHourFrom(),
	// 							Hour_to:   v.GetHourTo(),
	// 						})
	// 					}
	// 					return res
	// 				}(),
	// 				Geo_pos: GeoPos{
	// 					Lantitude: v.GetGeoPos().GetLantitude(),
	// 					Longitude: v.GetGeoPos().GetLongitude(),
	// 				},
	// 			})
	// 		}
	// 		return res
	// 	}(),
	//
	// 	Email: profile.GetEmail().GetEmail(),
	// 	Phone: func() *PhoneProfile {
	// 		if profile.GetPhone() != nil {
	// 			return &PhoneProfile{
	// 				Primary:      true,
	// 				Number:       profile.GetPhone().GetNumber(),
	// 				Country_code: profile.GetPhone().GetCountryCode().GetCode(),
	// 			}
	// 		}
	// 		return nil
	// 	}(),
	// }
}

// About us

// TODO:
func (_ *Resolver) ChangeCompanyAboutUs(ctx context.Context, in ChangeCompanyAboutUsRequest) (*SuccessResolver, error) {
	changes := companyRPC.ChangeCompanyAboutUsRequest{
		Id:             in.Company_id,
		Description:    NullToString(in.Changes.Description),
		Mission:        NullToString(in.Changes.Mission),
		Type:           stringToCompanyTypeRPC(NullToString(in.Changes.Type)),
		Size:           stringToCompanySizeRPC(NullToString(in.Changes.Size)),
		FoundationDate: NullToString(in.Changes.Foundation_date),
		Parking:        stringToCompanyParkingRPC(NullToString(in.Changes.Parking)),
	}

	if in.Changes.Description == nil {
		changes.IsDescriptionNull = true
	}

	if in.Changes.Mission == nil {
		changes.IsMissionNull = true
	}

	if in.Changes.Type == nil {
		changes.IsTypeNull = true
	}

	if in.Changes.Size == nil {
		changes.IsSizeNull = true
	}

	if in.Changes.Parking == nil {
		changes.IsParkingNull = true
	}

	if in.Changes.Industry == nil {
		changes.IsSubindustryNull = true
	}

	if in.Changes.Industry != nil {
		changes.Industry = &companyRPC.Industry{
			Main: in.Changes.Industry.ID,
			Subs: in.Changes.Industry.Subindustries,
		}
	}

	if in.Changes.Business_hours != nil {
		changes.BusinessHours = make([]*companyRPC.BusinessHoursItem, 0, len(*in.Changes.Business_hours))

		for i := range *in.Changes.Business_hours {
			changes.BusinessHours = append(changes.BusinessHours, &companyRPC.BusinessHoursItem{
				// ID:       (*in.Changes.Business_hours)[i].ID,
				WeekDays: (*in.Changes.Business_hours)[i].Week_days,
				HourFrom: (*in.Changes.Business_hours)[i].Hour_from,
				HourTo:   (*in.Changes.Business_hours)[i].Hour_to,
			})
		}
	}
	_, err := company.ChangeCompanyAboutUs(ctx, &changes)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil

	// 	_, err := company.ChangeCompanyAboutUs(ctx, &companyRPC.ChangeCompanyAboutUsRequest{
	// 		Id:          in.Company_id,
	// 		Description: NullToString(in.Changes.Description),
	// 		Mission:     NullToString(in.Changes.Mission),
	// 		Industry: &companyRPC.Industry{
	// 			Id: func() companyRPC.IndustryId {
	// 				if in.Changes.Industry.ID == "industry_information_technology" {
	// 					return companyRPC.IndustryId_INDUSTRY_INFORMATION_TECHNOLOGY
	// 				}
	// 				return companyRPC.IndustryId_INDUSTRY_UNKNOWN
	// 			}(),
	// 			Subindustries: func() []*companyRPC.Industry_Subindustry {
	// 				var res = []*companyRPC.Industry_Subindustry{}
	// 				for _, v := range in.Changes.Industry.Subindustries {
	// 					if v == "subindustry_web_development" {
	// 						res = append(res, &companyRPC.Industry_Subindustry{
	// 							Id: companyRPC.SubindustryId_SUBINDUSTRY_WEB_DEVEVELOPMENT,
	// 						})
	// 					} else if v == "subindustry_it_management" {
	// 						res = append(res, &companyRPC.Industry_Subindustry{
	// 							Id: companyRPC.SubindustryId_SUBINDUSTRY_IT_MANAGEMENT,
	// 						})
	// 					} else {
	// 						res = append(res, &companyRPC.Industry_Subindustry{
	// 							Id: companyRPC.SubindustryId_SUBINDUSTRY_UNKNOWN,
	// 						})
	// 					}
	// 				}
	// 				return res
	// 			}(),
	// 		},
	// 		Type: func() companyRPC.Type {
	// 			if in.Changes.Type != nil {
	// 				switch *in.Changes.Type {
	// 				case "type_self_employed":
	// 					return companyRPC.Type_TYPE_SELF_EMPLOYED
	// 				case "type_educational_institution":
	// 					return companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION
	// 				case "type_government_agency":
	// 					return companyRPC.Type_TYPE_GOVERNMENT_AGENSY
	// 				case "type_sole_proprietorship":
	// 					return companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP
	// 				case "type_privately_held":
	// 					return companyRPC.Type_TYPE_PRIVATELY_HELD
	// 				case "type_partnership":
	// 					return companyRPC.Type_TYPE_PARTNERSHIP
	// 				case "type_public_company":
	// 					return companyRPC.Type_TYPE_PUBLIC_COMPANY
	// 				}
	// 			}
	// 			return companyRPC.Type_TYPE_UNKNOWN
	// 		}(),
	// 		Size: func() companyRPC.Size {
	// 			if in.Changes.Size != nil {
	// 				switch *in.Changes.Size {
	// 				case "size_self_employed":
	// 					return companyRPC.Size_SIZE_SELF_EMPLOYED
	// 				case "size_1_10_employees":
	// 					return companyRPC.Size_SIZE_1_10_EMPLOYEES
	// 				case "size_11_50_employees":
	// 					return companyRPC.Size_SIZE_11_50_EMPLOYEES
	// 				case "size_51_200_employees":
	// 					return companyRPC.Size_SIZE_51_200_EMPLOYEES
	// 				case "size_201_500_employees":
	// 					return companyRPC.Size_SIZE_201_500_EMPLOYEES
	// 				case "size_501_1000_employees":
	// 					return companyRPC.Size_SIZE_501_1000_EMPLOYEES
	// 				case "size_1001_5000_employees":
	// 					return companyRPC.Size_SIZE_1001_5000_EMPLOYEES
	// 				case "size_5001_10000_employees":
	// 					return companyRPC.Size_SIZE_5001_10000_EMPLOYEES
	// 				case "size_10001_plus_employees":
	// 					return companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES
	// 				}
	// 			}
	// 			return companyRPC.Size_SIZE_UNKNOWN
	// 		}(),
	// 		FoundationDate: NullToString(in.Changes.Foundation_date),
	// 		Parking: func() companyRPC.Parking {
	// 			if in.Changes.Parking != nil {
	// 				switch *in.Changes.Parking {
	// 				case "parking_no_parking":
	// 					return companyRPC.Parking_PARKING_NO_PARKING
	// 				case "parking_parking_lot":
	// 					return companyRPC.Parking_PARKING_PARKING_LOT
	// 				case "parking_street_parking":
	// 					return companyRPC.Parking_PARKING_STREET_PARKING
	// 				}
	// 			}
	// 			return companyRPC.Parking_PARKING_UNKNOWN
	// 		}(),
	// 		BusinessHours: func() []*companyRPC.BusinessHoursItem {
	// 			res := []*companyRPC.BusinessHoursItem{}
	// 			if in.Changes.Business_hours != nil {
	// 				for _, v := range *in.Changes.Business_hours {
	// 					res = append(res, &companyRPC.BusinessHoursItem{
	// 						WeekDays: v.Week_days,
	// 						HourFrom: v.Hour_from,
	// 						HourTo:   v.Hour_to,
	// 					})
	// 				}
	// 			}
	// 			return res
	// 		}(),
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &SuccessResolver{R: &Success{Success: true}}, nil
}

//Founders

// TODO:
func (_ *Resolver) AddCompanyFounder(ctx context.Context, in AddCompanyFounderRequest) (*SuccessResolver, error) {
	result, err := company.AddCompanyFounder(ctx, &companyRPC.AddCompanyFounderRequest{
		ID: in.Company_id,
		Founder: &companyRPC.Founder{
			Name:          NullToString(in.Input.Name),
			PositionTitle: in.Input.Position_title,
			UserID:        NullToString(in.Input.User_id),
			// Avatar:
		},
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{
		ID:      result.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) DeleteCompanyFounder(ctx context.Context, in DeleteCompanyFounderRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyFounder(ctx, &companyRPC.DeleteCompanyFounderRequest{
		Id:        in.Company_id,
		FounderId: in.ID,
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyFounder(ctx context.Context, in ChangeCompanyFounderRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyFounder(ctx, &companyRPC.ChangeCompanyFounderRequest{
		Id: in.Company_id,
		Founder: &companyRPC.Founder{
			ID:            in.Changes.ID,
			Name:          NullToString(in.Changes.Name),
			PositionTitle: NullToString(in.Changes.Position_title),
			// Avatar:        NullToString(in.Changes.Avatar),
		},
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ApproveFounderRequest(ctx context.Context, in ApproveFounderRequestRequest) (*SuccessResolver, error) {
	_, err := company.ApproveFounderRequest(ctx, &companyRPC.ApproveFounderRequestRequest{
		CompanyID: in.Company_id,
		RequestID: in.Request_id,
	},
	)
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveFounderRequest(ctx context.Context, in RemoveFounderRequestRequest) (*SuccessResolver, error) {
	_, err := company.RemoveFounderRequest(ctx, &companyRPC.RemoveFounderRequestRequest{
		CompanyID: in.Company_id,
		RequestID: in.Request_id,
	},
	)
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

//AddCompanyAward ...
func (_ *Resolver) AddCompanyAward(ctx context.Context, in AddCompanyAwardRequest) (*SuccessResolver, error) {
	result, err := company.AddCompanyAward(ctx, &companyRPC.AddCompanyAwardRequest{
		ID: in.Company_id,
		Award: &companyRPC.Award{
			Title:  in.Input.Title,
			Issuer: in.Input.Issuer,
			Year:   in.Input.Year,
			Files: func(ids *[]string) []*companyRPC.File {
				if ids == nil {
					return nil
				}
				files := make([]*companyRPC.File, 0, len(*ids))
				for i := range *ids {
					files = append(files, &companyRPC.File{ID: (*ids)[i]})
				}

				return files
			}(in.Input.Files_id),

			Links: func(in *[]LinkInput) []*companyRPC.Link {
				if in != nil {
					ar := make([]*companyRPC.Link, len(*in))
					for r := range *in {
						ar[r] = &companyRPC.Link{
							URL: (*in)[r].Url,
						}
					}

					return ar
				}
				return []*companyRPC.Link{}
			}(in.Input.Link),
		},
	})
	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{
		ID:      result.GetID(),
		Success: true,
	}}, nil
}

func (_ *Resolver) DeleteCompanyAward(ctx context.Context, in DeleteCompanyAwardRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyAward(ctx, &companyRPC.DeleteCompanyAwardRequest{
		Id:      in.Company_id,
		AwardId: in.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyAward(ctx context.Context, in ChangeCompanyAwardRequest) (*SuccessResolver, error) {
	award := &companyRPC.ChangeCompanyAwardRequest{
		Id: in.Company_id,
		Award: &companyRPC.Award{
			ID:     in.Changes.ID,
			Title:  NullToString(in.Changes.Title),
			Issuer: NullToString(in.Changes.Issuer),
			Year:   NullToInt32(in.Changes.Year),
		},
	}

	if in.Changes.Link != nil {
		award.Award.Links = make([]*companyRPC.Link, 0, len(*in.Changes.Link))
		for _, link := range *in.Changes.Link {
			award.Award.Links = append(award.Award.Links, &companyRPC.Link{
				URL: link.Url,
				ID:  NullToString(link.ID),
			})
		}
	}

	_, err := company.ChangeCompanyAward(ctx, award)
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddLinksInCompanyAward(ctx context.Context, input AddLinksInCompanyAwardRequest) (*SuccessResolver, error) {
	_, err := company.AddLinksInCompanyAward(ctx, &companyRPC.AddLinksRequest{
		ID:      input.Company_id,
		AwardID: input.ID,
		Links: func(in []LinkInput) []*companyRPC.Link {
			ar := make([]*companyRPC.Link, len(in))
			for r := range in {
				ar[r] = &companyRPC.Link{
					URL: in[r].Url,
				}
			}
			return ar
		}(input.Input),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeLinkInCompanyAward(ctx context.Context, input ChangeLinkInCompanyAwardRequest) (*SuccessResolver, error) {
	_, err := company.ChangeLinkInCompanyAward(ctx, &companyRPC.ChangeLinkRequest{
		ID:      input.Company_id,
		AwardID: input.ID,
		Link: &companyRPC.Link{
			ID:  input.Link_id,
			URL: input.Url,
		},
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveLinksInCompanyAward(ctx context.Context, input RemoveLinksInCompanyAwardRequest) (*SuccessResolver, error) {
	_, err := company.RemoveLinksInCompanyAward(ctx, &companyRPC.RemoveLinksRequest{
		ID:      input.Company_id,
		AwardID: input.ID,
		Links: func(in []string) []*companyRPC.Link {
			ar := make([]*companyRPC.Link, len(in))
			for r := range in {
				ar[r] = &companyRPC.Link{
					ID: in[r],
				}
			}
			return ar
		}(input.Links_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveFilesInCompanyAward(ctx context.Context, input RemoveFilesInCompanyAwardRequest) (*SuccessResolver, error) {
	_, err := company.RemoveFilesInCompanyAward(ctx, &companyRPC.RemoveFilesRequest{
		ID:      input.Company_id,
		AwardID: input.ID,
		Files: func(in []string) []*companyRPC.File {
			ar := make([]*companyRPC.File, len(in))
			for r := range in {
				ar[r] = &companyRPC.File{
					ID: in[r],
				}
			}

			return ar
		}(input.Files_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveFilesInCompanyGallery(ctx context.Context, input RemoveFilesInCompanyGalleryRequest) (*SuccessResolver, error) {
	_, err := company.RemoveFilesInCompanyGallery(ctx, &companyRPC.RemoveGalleryFileRequest{
		ID: input.Company_id,
		Files: func(in []string) []*companyRPC.GalleryFile {
			ar := make([]*companyRPC.GalleryFile, len(in))
			for r := range in {
				ar[r] = &companyRPC.GalleryFile{
					ID: in[r],
				}
			}

			return ar
		}(input.Files_id),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

//Milestones
func (_ *Resolver) AddCompanyMilestone(ctx context.Context, in AddCompanyMilestoneRequest) (*SuccessResolver, error) {
	// 	var img string
	// 	if in.Input.Image != nil {
	// 		img = *in.Input.Image
	// 	}

	response, err := company.AddCompanyMilestone(ctx, &companyRPC.AddCompanyMilestoneRequest{
		ID: in.Company_id,
		Milestone: &companyRPC.Milestone{
			// 		Image:       img,
			Year:        in.Input.Year,
			Title:       in.Input.Title,
			Description: in.Input.Description,
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

func (_ *Resolver) DeleteCompanyMilestone(ctx context.Context, in DeleteCompanyMilestoneRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyMilestone(ctx, &companyRPC.DeleteCompanyMilestoneRequest{
		Id:          in.Company_id,
		MilestoneId: in.ID,
	})

	if err != nil {
		return nil, err
	}
	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyMilestone(ctx context.Context, in ChangeCompanyMilestoneRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyMilestone(ctx, &companyRPC.ChangeCompanyMilestoneRequest{
		Id: in.Company_id,
		Milestone: &companyRPC.Milestone{
			Id: in.Changes.ID,
			// Image:       NullToString(in.Changes.Image),
			Year:        NullToInt32(in.Changes.Year),
			Title:       NullToString(in.Changes.Title),
			Description: NullToString(in.Changes.Description),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveImageInMilestone(ctx context.Context, in RemoveImageInMilestoneRequest) (*SuccessResolver, error) {
	_, err := company.RemoveImageInMilestone(ctx, &companyRPC.RemoveImageInMilestoneRequest{
		ID:        in.ID,
		CompanyID: in.Company_id,
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

//Products
func (_ *Resolver) AddCompanyProduct(ctx context.Context, in AddCompanyProductRequest) (*SuccessResolver, error) {
	// 	var img string
	//
	// 	if in.Input.Image != nil {
	// 		*in.Input.Image = img
	// 	}

	response, err := company.AddCompanyProduct(ctx, &companyRPC.AddCompanyProductRequest{
		Id: in.Company_id,
		Product: &companyRPC.Product{
			// Image:   img,
			Name:    in.Input.Name,
			Website: in.Input.Website,
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

func (_ *Resolver) DeleteCompanyProduct(ctx context.Context, in DeleteCompanyProductRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyProduct(ctx, &companyRPC.DeleteCompanyProductRequest{
		Id:        in.Company_id,
		ProductId: in.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyProduct(ctx context.Context, in ChangeCompanyProductRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyProduct(ctx, &companyRPC.ChangeCompanyProductRequest{
		Id: in.Company_id,
		Product: &companyRPC.Product{
			ID: in.Changes.ID,
			// Image:   NullToString(in.Changes.Image),
			Name:    NullToString(in.Changes.Name),
			Website: NullToString(in.Changes.Website),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveImageInProduct(ctx context.Context, in RemoveImageInProductRequest) (*SuccessResolver, error) {
	_, err := company.RemoveImageInProduct(ctx, &companyRPC.RemoveImageInProductRequest{
		ID:        in.ID,
		CompanyID: in.Company_id,
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

// Services
func (_ *Resolver) AddCompanyService(ctx context.Context, in AddCompanyServiceRequest) (*SuccessResolver, error) {
	// 	var img string
	// 	if in.Input.Image != nil {
	// 		img = *in.Input.Image
	// 	}

	response, err := company.AddCompanyService(ctx, &companyRPC.AddCompanyServiceRequest{
		ID: in.Company_id,
		Service: &companyRPC.Service{
			// Image:   img,
			Name:    in.Input.Name,
			Website: in.Input.Website,
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

func (_ *Resolver) DeleteCompanyService(ctx context.Context, in DeleteCompanyServiceRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyService(ctx, &companyRPC.DeleteCompanyServiceRequest{
		Id:        in.Company_id,
		ServiceId: in.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) ChangeCompanyService(ctx context.Context, in ChangeCompanyServiceRequest) (*SuccessResolver, error) {
	_, err := company.ChangeCompanyService(ctx, &companyRPC.ChangeCompanyServiceRequest{
		Id: in.Company_id,
		Service: &companyRPC.Service{
			ID: in.Changes.ID,
			// Image:   NullToString(in.Changes.Image),
			Name:    NullToString(in.Changes.Name),
			Website: NullToString(in.Changes.Website),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RemoveImageInService(ctx context.Context, in RemoveImageInServiceRequest) (*SuccessResolver, error) {
	_, err := company.RemoveImageInService(ctx, &companyRPC.RemoveImageInServiceRequest{
		ID:        in.ID,
		CompanyID: in.Company_id,
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

// Review
func (_ *Resolver) AddCompanyReview(ctx context.Context, in AddCompanyReviewRequest) (*SuccessResolver, error) {
	response, err := company.AddCompanyReview(ctx, &companyRPC.AddCompanyReviewRequest{
		Id: in.Company_id,
		Review: &companyRPC.Review{
			Description: NullToString(in.Input.Description),
			Headline:    in.Input.Headline,
			Rate:        stringScroreToUint32(in.Input.Score),
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

// TODO:
func (_ *Resolver) DeleteCompanyReview(ctx context.Context, in DeleteCompanyReviewRequest) (*SuccessResolver, error) {
	_, err := company.DeleteCompanyReview(ctx, &companyRPC.DeleteCompanyReviewRequest{
		ID:       in.Company_id,
		ReviewID: in.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) AddCompanyReviewReport(ctx context.Context, in AddCompanyReviewReportRequest) (*SuccessResolver, error) {
	_, err := company.AddCompanyReviewReport(ctx, &companyRPC.AddCompanyReviewReportRequest{
		ReviewReport: &companyRPC.ReviewReport{
			CompanyID:   in.Company_id,
			ReviewId:    in.ID,
			Explanation: NullToString(in.Input.Explanation),
			Report:      stringToReviewReport(in.Input.Report),
		},
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

// TODO: change schema
func (_ *Resolver) GetCompanyReviews(ctx context.Context, in GetCompanyReviewsRequest) (*[]CompanyReviewResolver, error) {
	response, err := company.GetCompanyReviews(ctx, &companyRPC.GetCompanyReviewsRequest{
		ID:    in.Company_id,
		First: uint32(NullToInt32(in.Pagination.First)),
		After: NullToString((in.Pagination.After)),
	})
	if err != nil {
		return nil, err
	}

	reviews := make([]CompanyReviewResolver, 0, len(response.GetReviews()))
	for i := range response.GetReviews() {
		rev := companyReviewRPCToReview(response.GetReviews()[i])

		if response.GetReviews()[i].GetAuthorID() != "" {
			prof, err := user.GetProfileByID(ctx, &userRPC.ID{
				ID: response.GetReviews()[i].GetAuthorID(),
			})
			if err != nil {
				return nil, err
			}

			rev.User = ToProfile(ctx, prof)
		}

		reviews = append(reviews, CompanyReviewResolver{
			R: rev,
		})
	}

	return &reviews, nil
}

func (_ *Resolver) GetCompanyReviewsOfUser(ctx context.Context, in GetCompanyReviewsOfUserRequest) (*[]CompanyReviewForUserResolver, error) {
	response, err := company.GetUsersReviews(ctx, &companyRPC.GetCompanyReviewsRequest{
		ID:    in.User_id,
		First: uint32(NullToInt32(in.Pagination.First)),
		After: NullToString((in.Pagination.After)),
	})
	if err != nil {
		return nil, err
	}

	reviews := make([]CompanyReviewForUserResolver, 0, len(response.GetReviews()))
	for i := range response.GetReviews() {
		reviews = append(reviews, CompanyReviewForUserResolver{
			R: companyReviewRPCToReviewForUser(response.GetReviews()[i]),
		})
	}

	return &reviews, nil
}

// func (_ *Resolver) ChangeCompanyCover(ctx context.Context, in ChangeCompanyCoverRequest) (*SuccessResolver, error) {
// 	// 	_, err := company.ChangeCompanyCover(ctx, &companyRPC.ChangeCompanyCoverRequest{
// 	// 		Id:    in.Company_id,
// 	// 		Cover: NullToString(in.Cover),
// 	// 	})
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	return &SuccessResolver{R: &Success{Success: true}}, nil
// 	return nil, nil
// }

func (_ *Resolver) AddCompanyReport(ctx context.Context, in AddCompanyReportRequest) (*SuccessResolver, error) {
	// 	_, err := company.AddCompanyReport(ctx, &companyRPC.AddCompanyReportRequest{
	// 		Id: in.Company_id,
	// 		Report: func() companyRPC.ReportEnum {
	// 			switch in.Input.Report {
	// 			case "report_violates_terms_of_use":
	// 				return companyRPC.ReportEnum_REPORT_VIOLATES_TERMS_OF_USE
	// 			case "report_not_real_organization":
	// 				return companyRPC.ReportEnum_REPORT_NOT_REAL_ORGANIZATION
	// 			case "report_may_have_been_hacked":
	// 				return companyRPC.ReportEnum_REPORT_MAY_HAVE_BEEN_HACKED
	// 			case "report_picture_is_not_logo":
	// 				return companyRPC.ReportEnum_REPORT_PICTURE_IS_NOT_LOGO
	// 			case "report_duplicate":
	// 				return companyRPC.ReportEnum_REPORT_DUPLICATE
	// 			case "report_somthing_else":
	// 				return companyRPC.ReportEnum_REPORT_SOMTHING_ELSE
	// 			}
	// 			return companyRPC.ReportEnum_REPORT_UNKNOWN
	// 		}(),
	// 		Explanation: NullToString(in.Input.Explanation),
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &SuccessResolver{R: &Success{Success: true}}, nil

	return nil, nil
}

// TODO: remove
func (_ *Resolver) DeleteCompanyReport(ctx context.Context, in DeleteCompanyReportRequest) (*SuccessResolver, error) {
	// 	_, err := company.DeleteCompanyReport(ctx, &companyRPC.DeleteCompanyReportRequest{
	// 		Id:       in.Company_id,
	// 		ReportId: in.ID,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &SuccessResolver{R: &Success{Success: true}}, nil
	return nil, nil
}

// TODO: remove
func (_ *Resolver) GetCompanyReports(ctx context.Context, in GetCompanyReportsRequest) (*[]CompanyReportResolver, error) {
	// 	resp, err := company.GetCompanyReports(ctx, &companyRPC.GetCompanyReportsRequest{
	// 		Id: in.Company_id,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return func() *[]CompanyReportResolver {
	// 		if len(resp.Reports) == 0 {
	// 			return nil
	// 		}
	// 		res := []CompanyReportResolver{}
	// 		for _, v := range resp.Reports {
	// 			res = append(res, CompanyReportResolver{
	// 				R: &CompanyReport{
	// 					ID:          v.Id,
	// 					Report:      companyRPC.ReportEnum_name[int32(v.GetReport())],
	// 					Explanation: v.Explanation,
	// 				},
	// 			})
	// 		}
	// 		return &res
	// 	}(), nil

	return nil, nil
}

func (_ *Resolver) GetCompanyRate(ctx context.Context, in GetCompanyRateRequest) (*CompanyAvarageRateResolver, error) {
	result, err := company.GetAvarageRateOfCompany(ctx, &companyRPC.ID{
		ID: in.Company_id,
	})
	if err != nil {
		return nil, err
	}

	return &CompanyAvarageRateResolver{
		R: &CompanyAvarageRate{
			Amount_reviews: int32(result.GetAmountReviews()),
			Avarage_rate:   float64(result.GetAvarageRate()),
		},
	}, nil
}

func (_ *Resolver) GetAmountOfEachRate(ctx context.Context, in GetAmountOfEachRateRequest) (*CompanyScoreResolver, error) {
	result, err := company.GetAmountOfEachRate(ctx, &companyRPC.ID{
		ID: in.Company_id,
	})
	if err != nil {
		return nil, err
	}

	scores := result.GetRate()

	return &CompanyScoreResolver{
		R: &CompanyScore{
			Score_unknown:   int32(scores[0]),
			Score_poor:      int32(scores[1]),
			Score_fair:      int32(scores[2]),
			Score_good:      int32(scores[3]),
			Score_very_good: int32(scores[4]),
			Score_excellent: int32(scores[5]),
		},
	}, nil
}

func (_ *Resolver) RemoveCompanyAvatar(ctx context.Context, in RemoveCompanyAvatarRequest) (*SuccessResolver, error) {
	_, err := company.RemoveAvatar(ctx, &companyRPC.ID{
		ID: in.Company_id,
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

func (_ *Resolver) RemoveCompanyCover(ctx context.Context, in RemoveCompanyCoverRequest) (*SuccessResolver, error) {
	_, err := company.RemoveCover(ctx, &companyRPC.ID{
		ID: in.Company_id,
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

func (_ *Resolver) SaveCompanyProfileTranslation(ctx context.Context, in SaveCompanyProfileTranslationRequest) (*SuccessResolver, error) {
	_, err := company.SaveCompanyProfileTranslation(ctx, &companyRPC.ProfileTranslation{
		CompanyID:   in.Company_id,
		Description: in.Translations.Description,
		Mission:     in.Translations.Mission,
		Name:        in.Translations.Name,
		Language:    in.Translations.Language,
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

func (_ *Resolver) SaveCompanyMilestoneTranslation(ctx context.Context, in SaveCompanyMilestoneTranslationRequest) (*SuccessResolver, error) {
	_, err := company.SaveCompanyMilestoneTranslation(ctx, &companyRPC.MilestoneTranslation{
		CompanyID:   in.Company_id,
		Language:    in.LanguageID,
		MilestoneID: in.Translations.Milestone_id,
		Title:       in.Translations.Title,
		Desciption:  in.Translations.Description,
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

func (_ *Resolver) SaveCompanyAwardTranslation(ctx context.Context, in SaveCompanyAwardTranslationRequest) (*SuccessResolver, error) {
	_, err := company.SaveCompanyAwardTranslation(ctx, &companyRPC.AwardTranslation{
		CompanyID: in.Company_id,
		Language:  in.LanguageID,
		AwardID:   in.Translations.Award_id,
		Title:     in.Translations.Title,
		Issuer:    in.Translations.Issuer,
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

// OpenCareerCenter ...
func (_ *Resolver) OpenCareerCenter(ctx context.Context, in OpenCareerCenterRequest) (*bool, error) {
	_, err := company.OpenCareerCenter(ctx, &companyRPC.OpenCareerCenterRequest{
		CompanyID:           in.Company_id,
		CVButtonEnabled:     in.Input.Cv_button_enabled,
		CustomButtonEnabled: in.Input.Custom_button.Enabled,
		CustomButtonURL:     in.Input.Custom_button.Url,
		CustomButtontitle:   in.Input.Custom_button.Title,
		Description:         in.Input.Description,
		Title:               in.Input.Title,
	})
	if err != nil {
		return nil, err
	}

	t := true
	return &t, nil
}

// To RPC
func stringToCompanyTypeRPC(data string) companyRPC.Type {
	switch data {
	case "type_self_employed":
		return companyRPC.Type_TYPE_SELF_EMPLOYED
	case "type_educational_institution":
		return companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION
	case "type_government_agency":
		return companyRPC.Type_TYPE_GOVERNMENT_AGENCY
	case "type_sole_proprietorship":
		return companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP
	case "type_privately_held":
		return companyRPC.Type_TYPE_PRIVATELY_HELD
	case "type_partnership":
		return companyRPC.Type_TYPE_PARTNERSHIP
	case "type_public_company":
		return companyRPC.Type_TYPE_PUBLIC_COMPANY

	}

	return companyRPC.Type_TYPE_UNKNOWN
}

func stringToCompanySizeRPC(data string) companyRPC.Size {

	switch data {
	case "size_self_employed":
		return companyRPC.Size_SIZE_SELF_EMPLOYED
	case "size_1_10_employees":
		return companyRPC.Size_SIZE_1_10_EMPLOYEES
	case "size_11_50_employees":
		return companyRPC.Size_SIZE_11_50_EMPLOYEES
	case "size_51_200_employees":
		return companyRPC.Size_SIZE_51_200_EMPLOYEES
	case "size_201_500_employees":
		return companyRPC.Size_SIZE_201_500_EMPLOYEES
	case "size_501_1000_employees":
		return companyRPC.Size_SIZE_501_1000_EMPLOYEES
	case "size_1001_5000_employees":
		return companyRPC.Size_SIZE_1001_5000_EMPLOYEES
	case "size_5001_10000_employees":
		return companyRPC.Size_SIZE_5001_10000_EMPLOYEES
	case "size_10001_plus_employees":
		return companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES
	}

	return companyRPC.Size_SIZE_UNKNOWN
}

func stringToCompanyParkingRPC(data string) companyRPC.Parking {
	parking := companyRPC.Parking_PARKING_UNKNOWN

	switch data {
	case "parking_no_parking":
		parking = companyRPC.Parking_PARKING_NO_PARKING
	case "parking_parking_lot":
		parking = companyRPC.Parking_PARKING_PARKING_LOT
	case "parking_street_parking":
		parking = companyRPC.Parking_PARKING_STREET_PARKING
	}

	return parking
}

func stringToCompanyAdminRoleRPC(data string) companyRPC.AdminRole {
	role := companyRPC.AdminRole_ROLE_UNKNOWN

	switch data {
	case "admin":
		role = companyRPC.AdminRole_ROLE_ADMIN
	case "job_editor":
		role = companyRPC.AdminRole_ROLE_JOB_EDITOR
	case "commercial_admin":
		role = companyRPC.AdminRole_ROLE_COMMERCIAL_ADMIN
	case "v_shop_admin":
		role = companyRPC.AdminRole_ROLE_V_SHOP_ADMIN
	case "v_service_admin":
		role = companyRPC.AdminRole_ROLE_V_SERVICE_ADMIN
	}

	return role
}

// TODO:
func companyEmailToCompanyEmail(data *CompanyEmail) *companyRPC.Email {
	if data == nil {
		return nil
	}

	email := companyRPC.Email{
		// ID:    data.ID,
		Email: data.Email,
	}

	return &email
}

func stringScroreToUint32(data string) uint32 {
	var score uint32

	switch data {
	case "score_poor":
		score = 1
	case "score_fair":
		score = 2
	case "score_good":
		score = 3
	case "score_very_good":
		score = 4
	case "score_excellent":
		score = 5
	}

	return score
}

func stringToReviewReport(data string) companyRPC.ReviewReportEnum {
	rp := companyRPC.ReviewReportEnum_REVIEW_REPORT_UNKNOWN

	switch data {
	case "review_report_spam":
		rp = companyRPC.ReviewReportEnum_REVIEW_REPORT_SPAM
	case "review_report_scam":
		rp = companyRPC.ReviewReportEnum_REVIEW_REPORT_SCAM
	case "review_report_inappropriate_offensive":
		rp = companyRPC.ReviewReportEnum_REVIEW_REPORT_INAPPROPRIATE_OFFENSIVE
	case "review_report_false_fake":
		rp = companyRPC.ReviewReportEnum_REVIEW_REPORT_FALSE_FAKE
	case "review_report_off_topic":
		rp = companyRPC.ReviewReportEnum_REVIEW_REPORT_OFF_TOPIC
	case "review_report_somthing_else":
		rp = companyRPC.ReviewReportEnum_REVIEW_REPORT_SOMTHING_ELSE
	}

	return rp
}

func companyPhoneToCompanyPhoneRPC(data *CompanyPhone) *companyRPC.Phone {
	if data == nil {
		return nil
	}

	phone := companyRPC.Phone{
		ID:        data.ID,
		IsPrimary: data.Is_primary,
		Number:    data.Number,
		// CountryCode: CountryCode{
		// 	ID: data.Country_code_id,
		// },
	}

	return &phone
}

func businessHourItemToCompanyBusinessHourRPC(data *CompanyBusinessHoursInput) *companyRPC.BusinessHoursItem {
	if data == nil {
		return nil
	}

	bh := companyRPC.BusinessHoursItem{
		// ID:       data.,
		HourFrom: data.Hour_from,
		HourTo:   data.Hour_to,
		WeekDays: data.Week_days,
	}
	return &bh
}

// From RPC

func companyAccountRPCToCompanyAccount(data *companyRPC.Account) *CompanyAccount {
	if data == nil {
		return nil
	}

	account := CompanyAccount{
		ID:   data.GetID(),
		Name: data.GetName(),
		Url:  data.GetURL(),
		Industry: CompanyIndustry{
			ID:            data.GetIndustry().GetMain(),
			Subindustries: data.GetIndustry().GetSubs(),
		},
		Type:            companyTypeRPCToString(data.GetType()),
		Size:            companySizeRPCToString(data.GetSize()),
		Owner_id:        data.GetOwnerID(),
		Foundation_date: data.GetFoundationDate(),
		Addresses:       make([]CompanyAddress, 0, len(data.GetAddresses())),
		Emails:          make([]CompanyEmail, 0, len(data.GetEmails())),
		Phones:          make([]CompanyPhone, 0, len(data.GetPhones())),
		Websites:        make([]CompanyWebsite, 0, len(data.GetWebsites())),
		Parking:         companyParkingRPCToString(data.GetParking()),
		Status:          companyStatusRPCToString(data.GetStatus()),
		Business_hours:  make([]BusinessHoursItem, 0, len(data.GetBusinessHours())),
	}

	for i := range data.GetAddresses() {
		if address := companyAddressRPCToCompanyAddress(data.GetAddresses()[i]); address != nil {
			account.Addresses = append(account.Addresses, *address)
		}
	}

	for i := range data.GetEmails() {
		if email := companyEmailRPCToCompanyEmail(data.GetEmails()[i]); email != nil {
			account.Emails = append(account.Emails, *email)
		}
	}

	for i := range data.GetPhones() {
		if phone := companyPhoneRPCToCompanyPhone(data.GetPhones()[i]); phone != nil {
			account.Phones = append(account.Phones, *phone)
		}
	}

	for i := range data.GetWebsites() {
		if website := companyWebsiteRPCToCompanyWebsite(data.GetWebsites()[i]); website != nil {
			account.Websites = append(account.Websites, *website)
		}
	}

	for i := range data.GetBusinessHours() {
		if bh := companyBusinessHourRPCToBusinessHourItem(data.GetBusinessHours()[i]); bh != nil {
			account.Business_hours = append(account.Business_hours, *bh)
		}
	}

	return &account
}

// TODO:
func companyAddressRPCToCompanyAddress(data *companyRPC.Address) *CompanyAddress {
	if data == nil {
		return nil
	}

	address := CompanyAddress{
		ID:             data.GetID(),
		Apartment:      data.GetApartment(),
		Zip_code:       data.GetZipCode(),
		Name:           data.GetName(),
		Street_address: data.GetStreetAddress(),
		Business_hours: make([]BusinessHoursItem, 0, len(data.GetBusinessHours())),
		Phones:         make([]CompanyPhone, 0, len(data.GetPhones())),
		Country_id:     data.GetLocation().GetCountry().GetId(),
		Primary:        data.GetIsPrimary(),
		City: &CityInAddress{
			City:        data.GetLocation().GetCity().GetTitle(),
			ID:          data.GetLocation().GetCity().GetId(),
			Subdivision: data.GetLocation().GetCity().GetSubdivision(),
		},
		// State
		// Country_id
		// City_id
	}

	if gp := companyGeoPosRPCToGeoPos(data.GetGeoPos()); gp != nil {
		address.Geo_pos = *gp
	}

	for i := range data.GetBusinessHours() {
		if bh := companyBusinessHourRPCToBusinessHourItem(data.GetBusinessHours()[i]); bh != nil {
			address.Business_hours = append(address.Business_hours, *bh)
		}
	}

	for i := range data.GetPhones() {
		if phone := companyPhoneRPCToCompanyPhone(data.GetPhones()[i]); phone != nil {
			address.Phones = append(address.Phones, *phone)
		}
	}

	return &address
}

func companyBusinessHourRPCToBusinessHourItem(data *companyRPC.BusinessHoursItem) *BusinessHoursItem {
	if data == nil {
		return nil
	}

	bh := BusinessHoursItem{
		ID:        data.GetID(),
		Hour_from: data.GetHourFrom(),
		Hour_to:   data.GetHourTo(),
		Week_days: data.GetWeekDays(),
	}
	return &bh
}

func companyPhoneRPCToCompanyPhone(data *companyRPC.Phone) *CompanyPhone {
	if data == nil {
		return nil
	}

	phone := CompanyPhone{
		ID:         data.GetID(),
		Is_primary: data.GetIsPrimary(),
		Country_id: data.GetCountryAbbreviation(),
		Number:     data.GetNumber(),
	}

	if data.GetCountryCode() != nil {
		phone.Country_code = data.GetCountryCode().GetCode()
		phone.Country_code_id = data.GetCountryCode().GetId()
	}

	return &phone
}

func companyGeoPosRPCToGeoPos(data *companyRPC.GeoPos) *GeoPos {
	if data == nil {
		return nil
	}

	gp := GeoPos{
		Lantitude: data.GetLantitude(),
		Longitude: data.GetLongitude(),
	}

	return &gp
}

func companyTypeRPCToString(data companyRPC.Type) string {
	switch data {
	case companyRPC.Type_TYPE_SELF_EMPLOYED:
		return "type_self_employed"
	case companyRPC.Type_TYPE_EDUCATIONAL_INSTITUTION:
		return "type_educational_institution"
	case companyRPC.Type_TYPE_GOVERNMENT_AGENCY:
		return "type_government_agency"
	case companyRPC.Type_TYPE_SOLE_PROPRIETORSHIP:
		return "type_sole_proprietorship"
	case companyRPC.Type_TYPE_PRIVATELY_HELD:
		return "type_privately_held"
	case companyRPC.Type_TYPE_PARTNERSHIP:
		return "type_partnership"
	case companyRPC.Type_TYPE_PUBLIC_COMPANY:
		return "type_public_company"
	}

	return "type_unknown"
}

func companyTypeRPCToStringArray(data []companyRPC.Type) []string {
	ctype := make([]string, 0, len(data))

	for _, t := range data {
		ctype = append(ctype, companyTypeRPCToString(t))
	}

	return ctype
}

func companySizeRPCToString(data companyRPC.Size) string {
	switch data {
	case companyRPC.Size_SIZE_SELF_EMPLOYED:
		return "size_self_employed"
	case companyRPC.Size_SIZE_1_10_EMPLOYEES:
		return "size_1_10_employees"
	case companyRPC.Size_SIZE_11_50_EMPLOYEES:
		return "size_11_50_employees"
	case companyRPC.Size_SIZE_51_200_EMPLOYEES:
		return "size_51_200_employees"
	case companyRPC.Size_SIZE_201_500_EMPLOYEES:
		return "size_201_500_employees"
	case companyRPC.Size_SIZE_501_1000_EMPLOYEES:
		return "size_501_1000_employees"
	case companyRPC.Size_SIZE_1001_5000_EMPLOYEES:
		return "size_1001_5000_employees"
	case companyRPC.Size_SIZE_5001_10000_EMPLOYEES:
		return "size_5001_10000_employees"
	case companyRPC.Size_SIZE_10001_PLUS_EMPLOYEES:
		return "size_10001_plus_employees"
	}

	return "size_unknown"
}

func companySizeRPCToStringArray(data []companyRPC.Size) []string {
	size := make([]string, 0, len(data))

	for _, t := range data {
		size = append(size, companySizeRPCToString(t))
	}

	return size
}

func companyEmailRPCToCompanyEmail(data *companyRPC.Email) *CompanyEmail {
	email := CompanyEmail{
		ID:    data.GetID(),
		Email: data.GetEmail(),
	}

	return &email
}

func companyWebsiteRPCToCompanyWebsite(data *companyRPC.Website) *CompanyWebsite {
	if data == nil {
		return nil
	}

	website := CompanyWebsite{
		ID:      data.GetId(),
		Website: data.GetWebsite(),
	}

	return &website
}

func companyParkingRPCToString(data companyRPC.Parking) string {
	parking := "parking_unknown"

	switch data {
	case companyRPC.Parking_PARKING_NO_PARKING:
		parking = "parking_no_parking"
	case companyRPC.Parking_PARKING_PARKING_LOT:
		parking = "parking_parking_lot"
	case companyRPC.Parking_PARKING_STREET_PARKING:
		parking = "parking_street_parking"
	}

	return parking
}

func companyStatusRPCToString(data companyRPC.Status) string {
	s := "status_not_activated"

	switch data {
	case companyRPC.Status_STATUS_ACTIVATED:
		s = "status_activated"
	case companyRPC.Status_STATUS_DEACTIVATED:
		s = "status_deactivated"
	case companyRPC.Status_STATUS_NOT_ACTIVATED:
		s = "status_not_activated"
	}

	return s
}

func companyReviewRPCToReview(data *companyRPC.Review) *CompanyReview {
	if data == nil {
		return nil
	}

	review := CompanyReview{
		ID:          data.GetId(),
		Score:       scoreToString(data.GetRate()),
		Headline:    data.GetHeadline(),
		Description: data.GetDescription(),
		Created_at:  data.GetDate(),
	}

	return &review
}

func companyReviewRPCToReviewForUser(data *companyRPC.Review) (review *CompanyReviewForUser) {
	if data == nil {
		return nil
	}

	review = &CompanyReviewForUser{
		ID:          data.GetId(),
		Score:       scoreToString(data.GetRate()),
		Headline:    data.GetHeadline(),
		Description: data.GetDescription(),
		Created_at:  data.GetDate(),
	}

	if data.GetCompany() != nil {
		co := data.GetCompany()
		review.Company = toCompanyProfile(context.TODO(), *co)
	}

	return review
}

func companyAdminRoleRPCToString(data companyRPC.AdminRole) string {
	role := "role_unknown"

	switch data {
	case companyRPC.AdminRole_ROLE_ADMIN:
		role = "role_admin"
	case companyRPC.AdminRole_ROLE_JOB_EDITOR:
		role = "role_job_editor"
	case companyRPC.AdminRole_ROLE_COMMERCIAL_ADMIN:
		role = "role_commercial_admin"
	case companyRPC.AdminRole_ROLE_V_SHOP_ADMIN:
		role = "role_v_shop_admin"
	case companyRPC.AdminRole_ROLE_V_SERVICE_ADMIN:
		role = "role_v_service_admin"
	}

	return role
}

func networkAdminRoleRPCToString(data networkRPC.AdminLevel) string {
	role := "role_unknown"

	switch data {
	case networkRPC.AdminLevel_Admin:
		role = "role_admin"
	case networkRPC.AdminLevel_JobAdmin:
		role = "role_job_editor"
	case networkRPC.AdminLevel_CommercialAdmin:
		role = "role_commercial_admin"
	case networkRPC.AdminLevel_VShopAdmin:
		role = "role_v_shop_admin"
	case networkRPC.AdminLevel_VServiceAdmin:
		role = "role_v_service_admin"
	}

	return role
}

func scoreToString(s uint32) string {
	r := "score_unknown"

	switch s {
	case 1:
		return "score_poor"
	case 2:
		return "score_fair"
	case 3:
		return "score_good"
	case 4:
		return "score_very_good"
	case 5:
		return "score_excellent"
	}

	return r
}

func toCompanyPropertyOption(s string) notificationsRPC.ChangeCompanySettingsRequest_PropertyOption {
	switch s {
	case "new_follow":
		return notificationsRPC.ChangeCompanySettingsRequest_NewFollow
	case "new_review":
		return notificationsRPC.ChangeCompanySettingsRequest_NewReview
	case "new_applicant":
		return notificationsRPC.ChangeCompanySettingsRequest_NewApplicant
	}

	return notificationsRPC.ChangeCompanySettingsRequest_UnknownProperty
}

func companyBenefitToString(s string) companyRPC.BenefitEnum {
	switch s {
	case "childcare":
		return companyRPC.BenefitEnum_childcare
	case "labor_agreement":
		return companyRPC.BenefitEnum_labor_agreement
	case "remote_working":
		return companyRPC.BenefitEnum_remote_working
	case "floater":
		return companyRPC.BenefitEnum_floater
	case "paid_timeoff":
		return companyRPC.BenefitEnum_paid_timeoff
	case "flexible_working_hours":
		return companyRPC.BenefitEnum_flexible_working_hours
	case "additional_timeoff":
		return companyRPC.BenefitEnum_additional_timeoff
	case "additional_parental_leave":
		return companyRPC.BenefitEnum_parental_bonus
	case "sick_leave_for_family_members":
		return companyRPC.BenefitEnum_sick_leave_for_family_members
	case "company_daycare":
		return companyRPC.BenefitEnum_company_daycare
	case "company_canteen":
		return companyRPC.BenefitEnum_company_canteen
	case "sport_facilities":
		return companyRPC.BenefitEnum_sport_facilities
	case "access_for_handicapped_persons":
		return companyRPC.BenefitEnum_access_for_handicapped_persons
	case "employee_parking":
		return companyRPC.BenefitEnum_employee_parking
	case "shuttle_service":
		return companyRPC.BenefitEnum_shuttle_service
	case "multiple_work_spaces":
		return companyRPC.BenefitEnum_multiple_work_spaces
	case "corporate_events":
		return companyRPC.BenefitEnum_corporate_events
	case "trainig_and_development":
		return companyRPC.BenefitEnum_trainig_and_development
	case "pets_allowed":
		return companyRPC.BenefitEnum_pets_allowed
	case "corporate_medical_staff":
		return companyRPC.BenefitEnum_corporate_medical_staff
	case "game_consoles":
		return companyRPC.BenefitEnum_game_consoles
	case "snack_and_drink_selfservice":
		return companyRPC.BenefitEnum_snack_and_drink_selfservice
	case "private_pension_scheme":
		return companyRPC.BenefitEnum_private_pension_scheme
	case "health_insurance":
		return companyRPC.BenefitEnum_health_insurance
	case "dental_care":
		return companyRPC.BenefitEnum_dental_care
	case "car_insurance":
		return companyRPC.BenefitEnum_car_insurance
	case "tution_fees":
		return companyRPC.BenefitEnum_tution_fees
	case "permfomance_related_bonus":
		return companyRPC.BenefitEnum_permfomance_related_bonus
	case "stock_options":
		return companyRPC.BenefitEnum_stock_options
	case "profit_earning_bonus":
		return companyRPC.BenefitEnum_profit_earning_bonus
	case "additional_months_salary":
		return companyRPC.BenefitEnum_additional_months_salary
	case "employers_matching_contributions":
		return companyRPC.BenefitEnum_employers_matching_contributions
	case "parental_bonus":
		return companyRPC.BenefitEnum_parental_bonus
	case "tax_deductions":
		return companyRPC.BenefitEnum_tax_deductions
	case "language_courses":
		return companyRPC.BenefitEnum_language_courses
	case "company_car":
		return companyRPC.BenefitEnum_company_car
	case "laptop":
		return companyRPC.BenefitEnum_laptop
	case "discounts_on_company_products_and_services":
		return companyRPC.BenefitEnum_discounts_on_company_products_and_services
	case "holiday_vouchers":
		return companyRPC.BenefitEnum_holiday_vouchers
	case "restraunt_vouchers":
		return companyRPC.BenefitEnum_restraunt_vouchers
	case "corporate_housing":
		return companyRPC.BenefitEnum_corporate_housing
	case "mobile_phone":
		return companyRPC.BenefitEnum_mobile_phone
	case "gift_vouchers":
		return companyRPC.BenefitEnum_gift_vouchers
	case "cultural_or_sporting_activites":
		return companyRPC.BenefitEnum_cultural_or_sporting_activites
	case "employee_service_vouchers":
		return companyRPC.BenefitEnum_employee_service_vouchers
	case "corporate_credit_card":
		return companyRPC.BenefitEnum_corporate_credit_card
	case "transportation":
		return companyRPC.BenefitEnum_transportation
	}

	return companyRPC.BenefitEnum_other
}

func companyBenefitsRPCToStringArray(data []companyRPC.BenefitEnum) []string {
	cBenefits := make([]string, 0, len(data))

	for _, t := range data {
		cBenefits = append(cBenefits, stringToCompanyRPCBenefits(t))
	}

	return cBenefits
}

func stringToCompanyRPCBenefits(data companyRPC.BenefitEnum) string {
	switch data {
	case companyRPC.BenefitEnum_childcare:
		return "childcare"
	case companyRPC.BenefitEnum_labor_agreement:
		return "labor_agreement"
	case companyRPC.BenefitEnum_remote_working:
		return "remote_working"
	case companyRPC.BenefitEnum_floater:
		return "floater"
	case companyRPC.BenefitEnum_paid_timeoff:
		return "paid_timeoff"
	case companyRPC.BenefitEnum_flexible_working_hours:
		return "flexible_working_hours"
	case companyRPC.BenefitEnum_additional_timeoff:
		return "additional_timeoff"
	case companyRPC.BenefitEnum_additional_parental_leave:
		return "additional_parental_leave"
	case companyRPC.BenefitEnum_sick_leave_for_family_members:
		return "sick_leave_for_family_members"
	case companyRPC.BenefitEnum_company_daycare:
		return "company_daycare"
	case companyRPC.BenefitEnum_company_canteen:
		return "company_canteen"
	case companyRPC.BenefitEnum_sport_facilities:
		return "sport_facilities"
	case companyRPC.BenefitEnum_access_for_handicapped_persons:
		return "access_for_handicapped_persons"
	case companyRPC.BenefitEnum_employee_parking:
		return "employee_parking"
	case companyRPC.BenefitEnum_shuttle_service:
		return "shuttle_service"
	case companyRPC.BenefitEnum_multiple_work_spaces:
		return "multiple_work_spaces"
	case companyRPC.BenefitEnum_corporate_events:
		return "corporate_events"
	case companyRPC.BenefitEnum_trainig_and_development:
		return "trainig_and_development"
	case companyRPC.BenefitEnum_pets_allowed:
		return "pets_allowed"
	case companyRPC.BenefitEnum_corporate_medical_staff:
		return "corporate_medical_staff"
	case companyRPC.BenefitEnum_game_consoles:
		return "game_consoles"
	case companyRPC.BenefitEnum_snack_and_drink_selfservice:
		return "snack_and_drink_selfservice"
	case companyRPC.BenefitEnum_private_pension_scheme:
		return "private_pension_scheme"
	case companyRPC.BenefitEnum_health_insurance:
		return "health_insurance"
	case companyRPC.BenefitEnum_dental_care:
		return "dental_care"
	case companyRPC.BenefitEnum_car_insurance:
		return "car_insurance"
	case companyRPC.BenefitEnum_tution_fees:
		return "tution_fees"
	case companyRPC.BenefitEnum_permfomance_related_bonus:
		return "permfomance_related_bonus"
	case companyRPC.BenefitEnum_stock_options:
		return "stock_options"
	case companyRPC.BenefitEnum_profit_earning_bonus:
		return "profit_earning_bonus"
	case companyRPC.BenefitEnum_additional_months_salary:
		return "additional_months_salary"
	case companyRPC.BenefitEnum_employers_matching_contributions:
		return "employers_matching_contributions"
	case companyRPC.BenefitEnum_tax_deductions:
		return "tax_deductions"
	case companyRPC.BenefitEnum_language_courses:
		return "language_courses"
	case companyRPC.BenefitEnum_company_car:
		return "company_car"
	case companyRPC.BenefitEnum_laptop:
		return "laptop"
	case companyRPC.BenefitEnum_discounts_on_company_products_and_services:
		return "discounts_on_company_products_and_services"
	case companyRPC.BenefitEnum_holiday_vouchers:
		return "holiday_vouchers"
	case companyRPC.BenefitEnum_restraunt_vouchers:
		return "restraunt_vouchers"
	case companyRPC.BenefitEnum_corporate_housing:
		return "corporate_housing"
	case companyRPC.BenefitEnum_mobile_phone:
		return "mobile_phone"
	case companyRPC.BenefitEnum_gift_vouchers:
		return "gift_vouchers"
	case companyRPC.BenefitEnum_cultural_or_sporting_activites:
		return "cultural_or_sporting_activites"
	case companyRPC.BenefitEnum_employee_service_vouchers:
		return "employee_service_vouchers"
	case companyRPC.BenefitEnum_corporate_credit_card:
		return "corporate_credit_card"
	case companyRPC.BenefitEnum_transportation:
		return "transportation"
	case companyRPC.BenefitEnum_relocation_package:
		return "relocation_package"
	}

	return "other"
}

func careerCenterRPCToCareerCenter(data *companyRPC.CareerCenter) *CareerCenter {
	if data == nil {
		return &CareerCenter{}
		// return nil
	}

	return &CareerCenter{
		Custom_button_enabled: data.GetCustomButtonEnabled(),
		Custom_button_title:   data.GetCustomButtontitle(),
		Custom_button_url:     data.GetCustomButtonURL(),
		Cv_button_enabled:     data.GetCVButtonEnabled(),
		Description:           data.GetDescription(),
		Title:                 data.GetTitle(),
	}
}
