package resolver

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/advertRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func (_ *Resolver) GetAdvertGallery(ctx context.Context, input GetAdvertGalleryRequest) (AdvertGalleryResolver, error) {

	companyID := ""

	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	resp, err := advert.GetGallery(ctx, &advertRPC.GetGalleryRequest{
		CompanyID: companyID,
		Pagination: &advertRPC.Pagination{
			After: input.After,
			First: input.First,
		},
	})
	if err != nil {
		return AdvertGalleryResolver{}, err
	}

	files := make([]File, 0, len(resp.GetFiles()))

	for _, f := range resp.GetFiles() {
		files = append(files, File{
			ID:        f.GetID(),
			Address:   f.GetURL(),
			Mime_type: f.GetMimeType(),
			Name:      f.GetName(),
		})
	}

	return AdvertGalleryResolver{
		R: &AdvertGallery{
			Files:  files,
			Amount: int32(resp.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) RemoveAdvert(ctx context.Context, input RemoveAdvertRequest) (*SuccessResolver, error) {

	_, err := advert.RemoveAdvert(ctx, &advertRPC.PauseAdvertRequest{
		CampaignID: input.Campaign_id,
		AdvertID:   input.Advert_id,
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

func (_ *Resolver) RemoveAdvertCampaign(ctx context.Context, input RemoveAdvertCampaignRequest) (*SuccessResolver, error) {

	_, err := advert.RemoveAdvertCampaign(ctx, &advertRPC.PauseAdvertRequest{
		CampaignID: input.Campaign_id,
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

func (_ *Resolver) GetMyAdvert(ctx context.Context, input GetMyAdvertRequest) (AdvertRecordsResolver, error) {
	companyID := ""

	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	resp, err := advert.GetMyAdverts(ctx, &advertRPC.GetMyAdvertsRequest{
		CompanyID: companyID,
		Pagination: &advertRPC.Pagination{
			After: input.After,
			First: input.First,
		},
	})
	if err != nil {
		return AdvertRecordsResolver{}, err
	}

	ads := make([]Advert, 0, len(resp.GetAdverts()))

	for _, a := range resp.GetAdverts() {
		ads = append(ads, Advert{
			ID:         a.GetID(),
			Name:       a.GetName(),
			Start_date: a.GetStartDate(),
			Type:       advertRPCAdvertTypeToString(a.GetAdType()),
			End_date:   a.GetFinishDate(),
			Status:     advertRPCAdvertStatusToString(a.GetAdStatus()),
		})
	}

	return AdvertRecordsResolver{
		R: &AdvertRecords{
			Ads:    ads,
			Amount: int32(resp.GetAmount()),
		},
	}, nil
}

func (_ *Resolver) GetAdvertBanners(ctx context.Context, input GetAdvertBannersRequest) ([]AdvertBannerResolver, error) {
	resp, err := advert.GetBanners(ctx, &advertRPC.GetBannersRequest{
		CountryID: input.CountryID,
		Amount:    input.Amount,
		Format:    stringToAdvertFormat(input.Format),
	})
	if err != nil {
		return nil, err
	}

	ads := make([]AdvertBannerResolver, 0, len(resp.GetBanners()))

	for _, a := range resp.GetBanners() {
		contents := make([]AdvertContent, 0, len(a.GetContents()))

		for _, c := range a.GetContents() {
			contents = append(contents, AdvertContent{
				Description:     c.GetContent(),
				Title:           c.GetTitle(),
				Destination_url: c.GetDestinationURL(),
			})
		}

		ads = append(ads, AdvertBannerResolver{
			R: &AdvertBanner{
				// Is_responsive: a.GetIsResponsive(),
				// Url:      a.GetDestinationURL(),
				Contents: contents,
			},
		},
		)
	}

	for i := range ads {
		ads[i].R.Button_title = resp.GetBanners()[i].GetAdvert().GetButtonTitle()
	}

	return ads, nil
}

func (_ *Resolver) GetAdvertCandidates(ctx context.Context, input GetAdvertCandidatesRequest) ([]ProfileResolver, error) {
	resp, err := advert.GetCandidates(ctx, &advertRPC.GetCandidatesRequest{
		CountryID: input.CountryID,
		Amount:    input.Amount,
		Format:    stringToAdvertFormat(input.Format),
	})
	if err != nil {
		return nil, err
	}

	candidates := make([]ProfileResolver, 0, len(resp.GetIDs()))

	profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: resp.GetIDs(),
	})
	if err != nil {
		return nil, err
	}

	for _, p := range profiles.GetProfiles() {
		profile := ToProfile(ctx, p)
		candidates = append(candidates, ProfileResolver{
			R: &profile,
		},
		)
	}

	return candidates, nil
}

func (_ *Resolver) GetAdvertJobs(ctx context.Context, input GetAdvertJobsRequest) ([]CompanyProfileResolver, error) {
	resp, err := advert.GetJobs(ctx, &advertRPC.GetJobsRequest{
		CountryID: input.CountryID,
		Amount:    input.Amount,
	})
	if err != nil {
		return nil, err
	}

	candidates := make([]CompanyProfileResolver, 0, len(resp.GetIDs()))

	profiles, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: resp.GetIDs(),
	})
	if err != nil {
		return nil, err
	}

	for _, p := range profiles.GetProfiles() {
		if p != nil {
			profile := toCompanyProfile(ctx, *p)
			candidates = append(candidates, CompanyProfileResolver{
				R: &profile,
			},
			)
		}
	}

	return candidates, nil
}

func (_ *Resolver) GetAdvertBannerDraft(ctx context.Context, input GetAdvertBannerDraftRequest) (AdvertBannerDraftResolver, error) {
	companyID := ""

	if input.Company_id != nil {
		companyID = *input.Company_id
	}

	ad, err := advert.GetBannerDraft(ctx, &advertRPC.IDWithCompanyID{
		CompanyID: companyID,
		ID:        input.Banner_id,
	})
	if err != nil {
		return AdvertBannerDraftResolver{}, err
	}

	adv := advertRPCBannertoAdvertBannerDraft(ad)

	return AdvertBannerDraftResolver{
		R: adv,
	}, nil
}

// CreateAdvertCampaign ...
func (*Resolver) CreateAdvertCampaign(ctx context.Context, in CreateAdvertCampaignRequest) (*SuccessResolver, error) {
	res, err := advert.CreateAdvertCampaign(ctx, &advertRPC.AdvertCampaign{
		CompanyID:   NullToString(in.Input.Company_id),
		Name:        in.Input.Name,
		Currency:    in.Input.Currency,
		Clicks:      in.Input.Clicks,
		Forwarding:  in.Input.Forwarding,
		Impressions: in.Input.Impressions,
		Languages:   in.Input.Languages,
		Referals:    in.Input.Referals,
		StartDate:   in.Input.Start_date,
		AdvertType:  advertTypeToRPC(in.Input.Type),
		Locations:   advertLocationsToRPC(in.Input.Locations),
		Formats:     advertFormatsToRPC(in.Input.Formats),
		AgeFrom:     NullToInt32(in.Input.Age_from),
		AgeTo:       NullToInt32(in.Input.Age_to),
		Gender:      advertGenderToRPC(in.Input.Gender),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID: res.GetID(),
		},
	}, nil
}

func (*Resolver) PauseAdvertCampaign(ctx context.Context, in PauseAdvertCampaignRequest) (*SuccessResolver, error) {

	_, err := advert.PauseAdvertCampaign(ctx, &advertRPC.PauseAdvertRequest{
		CampaignID: in.Campaign_id,
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

func (*Resolver) ActiveAdvertCampaign(ctx context.Context, in ActiveAdvertCampaignRequest) (*SuccessResolver, error) {

	_, err := advert.ActiveAdvertCampaign(ctx, &advertRPC.PauseAdvertRequest{
		CampaignID: in.Campaign_id,
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

func (*Resolver) PauseAdvert(ctx context.Context, in PauseAdvertRequest) (*SuccessResolver, error) {

	_, err := advert.PauseAdvert(ctx, &advertRPC.PauseAdvertRequest{
		AdvertID: in.Advert_id,
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

func (*Resolver) ActiveAdvert(ctx context.Context, in ActiveAdvertRequest) (*SuccessResolver, error) {

	_, err := advert.ActiveAdvert(ctx, &advertRPC.PauseAdvertRequest{
		AdvertID: in.Advert_id,
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

// CreateAdvertByCampaign ...
func (*Resolver) CreateAdvertByCampaign(ctx context.Context, in CreateAdvertByCampaignRequest) (*SuccessResolver, error) {

	data := &advertRPC.CampaignAdvert{
		CampaignID: in.Campaign_id,
		Name:       in.Input.Name,
		TypeID:     NullToString(in.Input.ID),
		AdvertType: advertTypeToRPC(in.Input.Type),
		URL:        NullToString(in.Input.Url),
	}

	if content := in.Input.Content; content != nil {
		data.Contents = advertContentsToRPC(content)
	}

	res, err := advert.CreateAdvertByCampaign(ctx, data)

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID: res.GetID(),
		},
	}, nil
}

// GetAdvertCampaigns ...
func (*Resolver) GetAdvertCampaigns(ctx context.Context, in GetAdvertCampaignsRequest) (*CampaingsResolver, error) {

	ad, err := advert.GetAdvertCampaigns(ctx, &advertRPC.AdvertGetRequest{
		CompanyID: NullToString(in.Company_id),
		First:     Nullint32ToUint32(in.Pagination.First),
		After:     NullToString(in.Pagination.After),
	})

	if err != nil {
		return nil, err
	}

	// adv := advertRPCBannertoAdvertBannerDraft(ad)

	return &CampaingsResolver{
		R: &Campaings{
			Total_amount: ad.GetTotalAmount(),
			Campaings:    advertCampaignsRPCToCampaigns(ad.GetAdvertCampaigns()),
		},
	}, nil
}

// GetAdvertsByCampaignID ...
func (*Resolver) GetAdvertsByCampaignID(ctx context.Context, in GetAdvertsByCampaignIDRequest) (*AdvertsResolver, error) {

	ad, err := advert.GetAdvertsByCampaignID(ctx, &advertRPC.AdvertGetRequest{
		CampaignID: in.Campaign_id,
		CompanyID:  NullToString(in.Company_id),
		First:      Nullint32ToUint32(in.Pagination.First),
		After:      NullToString(in.Pagination.After),
	})

	if err != nil {
		return nil, err
	}

	return &AdvertsResolver{
		R: &Adverts{
			Total_amount: ad.GetAmount(),
			Capmaign:     advertCampaignRPCToCampaign(ad.GetCampaign()),
			Adverts:      advertAdvertsRPCToAdverts(ad.GetAdverts()),
		},
	}, nil
}

// GetAdvert ...
func (*Resolver) GetAdvert(ctx context.Context, in GetAdvertRequest) (*AdvertResolver, error) {
	ad, err := advert.GetAdvert(ctx, &advertRPC.GetAdvertRequest{
		AdvertType: advertTypeToRPC(in.Type),
	})
	if err != nil {
		return nil, err
	}

	if ad == nil {
		return &AdvertResolver{}, nil
	}

	res := advertAdvertRPCToAdvert(ad)

	return &AdvertResolver{
		R: &res,
	}, nil
}

func advertAdvertsRPCToAdverts(data []*advertRPC.Advert) []Advert {
	if len(data) <= 0 {
		return nil
	}

	adverts := make([]Advert, 0, len(data))

	for _, a := range data {
		adverts = append(adverts, advertAdvertRPCToAdvert(a))
	}

	return adverts
}

func advertAdvertRPCToAdvert(data *advertRPC.Advert) Advert {
	if data == nil {
		return Advert{}
	}

	res := Advert{
		ID:          data.GetID(),
		Type_id:     data.GetTypeID(),
		Name:        data.GetName(),
		Status:      advertRPCAdvertStatusToString(data.GetAdStatus()),
		Type:        advertRPCAdvertTypeToString(data.GetAdType()),
		Start_date:  data.GetStartDate(),
		Impressions: data.GetImpressions(),
		Clicks:      data.GetClicks(),
		Forwarding:  data.GetForwarding(),
		Referals:    data.GetReferals(),
		Ctr_avg:     data.GetCrtAVG(),
		Url:         data.GetURL(),
		Contents:    advertContentsRPCToContents(data.GetContents()),
		Formats:     advertFormatsToString(data.GetFormats()),
	}

	files := make([]File, 0, len(data.GetFiles()))

	for _, f := range data.GetFiles() {
		files = append(files, File{
			ID:        f.GetID(),
			Address:   f.GetURL(),
			Mime_type: f.GetMimeType(),
			Name:      f.GetName(),
		})

		res.Files = files
	}

	return res
}

func advertGenderRPCToGender(data advertRPC.GenderEnum) string {

	switch data {
	case advertRPC.GenderEnum_MALE:
		return "male"
	case advertRPC.GenderEnum_FEMALE:
		return "female"
	}
	return "both"
}

func advertContentsRPCToContents(data []*advertRPC.Content) []AdvertContent {
	if len(data) <= 0 {
		return nil
	}

	contents := make([]AdvertContent, 0, len(data))

	for _, c := range data {
		contents = append(contents, advertContentRPCToContent(c))
	}

	return contents
}

func advertContentRPCToContent(data *advertRPC.Content) AdvertContent {
	if data == nil {
		return AdvertContent{}
	}

	return AdvertContent{
		Title:           data.GetTitle(),
		Destination_url: data.GetDestinationURL(),
		Description:     data.GetContent(),
		Custom_button:   data.GetCustomButton(),
	}
}
func advertCampaignsRPCToCampaigns(data []*advertRPC.AdvertCampaign) []AdvertCampaign {
	if len(data) <= 0 {
		return nil
	}

	adverts := make([]AdvertCampaign, 0, len(data))

	for _, a := range data {
		adverts = append(adverts, advertCampaignRPCToCampaign(a))
	}

	return adverts
}

func advertCampaignRPCToCampaign(data *advertRPC.AdvertCampaign) AdvertCampaign {
	if data == nil {
		return AdvertCampaign{}
	}

	return AdvertCampaign{
		ID:          data.GetID(),
		Clicks:      data.GetClicks(),
		Ctr_avg:     data.GetCrtAVG(),
		Name:        data.GetName(),
		Forwarding:  data.GetForwarding(),
		Impressions: data.GetImpressions(),
		Referals:    data.GetReferals(),
		Start_date:  data.GetStartDate(),
		Age_from:    data.GetAgeFrom(),
		Age_to:      data.GetAgeTo(),
		Languages:   data.GetLanguages(),
		Status:      advertRPCAdvertStatusToString(data.GetStatus()),
		Type:        advertRPCAdvertTypeToString(data.GetAdvertType()),
		Gender:      advertGenderRPCToGender(data.GetGender()),
		Format:      advertFormatsToString(data.GetFormats()),
		Locations:   advertLocationsToRPCToLocations(data.GetLocations()),
	}
}

func advertGenderToRPC(data *string) advertRPC.GenderEnum {
	if data == nil {
		return advertRPC.GenderEnum_WITHOUT_GENDER
	}

	if *data == "female" {
		return advertRPC.GenderEnum_FEMALE
	}

	if *data == "both" {
		return advertRPC.GenderEnum_WITHOUT_GENDER
	}

	return advertRPC.GenderEnum_MALE
}

func advertContentsToRPC(data *[]AdvertCampaignContnetInput) []*advertRPC.Content {
	if data == nil {
		return nil
	}

	contnets := make([]*advertRPC.Content, 0, len(*data))

	for _, c := range *data {
		contnets = append(contnets, &advertRPC.Content{
			Title:          c.Headline,
			DestinationURL: NullToString(c.Url),
			CustomButton:   NullToString(c.Custom_button),
			Content:        NullToString(c.Description),
		})
	}

	return contnets
}
func advertFormatsToRPC(data []string) []advertRPC.Format {
	if len(data) <= 0 {
		return nil
	}

	formats := make([]advertRPC.Format, 0, len(data))

	for _, f := range data {
		formats = append(formats, stringToAdvertFormat(f))
	}

	return formats
}

func advertLocationsToRPC(data []LocationInput) []*advertRPC.Location {
	if len(data) <= 0 {
		return nil
	}

	locations := make([]*advertRPC.Location, 0, len(data))

	for _, l := range data {
		locations = append(locations, locationInputToAdvertLocation(&l))
	}

	return locations
}

func advertLocationsToRPCToLocations(data []*advertRPC.Location) []AdvertLocation {
	if len(data) <= 0 {
		return nil
	}

	locations := make([]AdvertLocation, 0, len(data))

	for _, l := range data {
		locations = append(locations, AdvertLocation{
			Country: l.GetCountryID(),
			City:    l.GetCity(),
		})
	}

	return locations
}

func bannerInputToAdvertRPCBanner(data BannerInput) *advertRPC.Banner {
	advert := advertRPC.Advert{
		StartDate:   data.Start_date,
		Name:        data.Name,
		Currency:    data.Currency,
		ButtonTitle: data.Button_title,
	}

	banner := advertRPC.Banner{
		Advert: &advert,
		// IsResponsive:   data.Is_responsive,
		// DestinationURL: data.Destination_url,
		Contents: make([]*advertRPC.Content, 0, len(data.Contents)),
		// Places:   make([]advertRPC.Place, 0, len(data.Places)),
	}

	for _, content := range data.Contents {
		banner.Contents = append(banner.Contents, contentInputToAdvertContentInput(&content))
	}

	// for _, place := range data.Places {
	// 	banner.Places = append(banner.Places, stringToAdvertPlace(place))
	// }

	return &banner
}

func locationInputToAdvertLocation(data *LocationInput) *advertRPC.Location {
	if data == nil {
		return nil
	}

	loc := advertRPC.Location{
		CountryID: data.Country_id,
	}

	if data.City.City != nil {
		loc.City = *data.City.City
	}

	if data.City.ID != nil {
		loc.CityID = *data.City.ID
	}

	if data.City.Subdivision != nil {
		loc.Subdivision = *data.City.Subdivision
	}

	return &loc
}

func contentInputToAdvertContentInput(data *AdvertContentInput) *advertRPC.Content {
	if data == nil {
		return nil
	}

	content := advertRPC.Content{
		FileID:         data.File_id,
		Title:          data.Title,
		Content:        data.Description,
		DestinationURL: data.Destination_url,
	}

	return &content
}

func stringToAdvertPlace(data string) advertRPC.Place {
	switch data {
	case "user":
		return advertRPC.Place_User
	case "company":
		return advertRPC.Place_Company
	case "local_business":
		return advertRPC.Place_LocalBusiness
	case "brands":
		return advertRPC.Place_Brands
	case "groups":
		return advertRPC.Place_Groups
	}

	return advertRPC.Place_User
}

func stringToAdvertFormat(data string) advertRPC.Format {
	switch data {
	case "CAROUSEL":
		return advertRPC.Format_CAROUSEL
	case "SLIDE":
		return advertRPC.Format_SLIDE
	case "VIDEO":
		return advertRPC.Format_VIDEO
	case "IMAGE":
		return advertRPC.Format_IMAGE
	case "RESPONSIVE":
		return advertRPC.Format_RESPONSIVE
	case "SPOTLIGHT":
		return advertRPC.Format_SPOTLIGHT
	case "SIDE_PIN":
		return advertRPC.Format_SIDE_PIN
	case "HEAD_PIN":
		return advertRPC.Format_HEAD_PIN
	case "BUSINESS_SEARCH":
		return advertRPC.Format_BUSINESS_SEARCH
	case "PROFESSIONAL_SEARCH":
		return advertRPC.Format_PROFESSIONAL_SEARCH
	case "CANDIDATE_SEARCH":
		return advertRPC.Format_CANDIDATE_SEARCH
	case "JOB_SEARCH":
		return advertRPC.Format_JOB_SEARCH
	case "SERVICE_SEARCH":
		return advertRPC.Format_SERVICE_SEARCH
	}

	return advertRPC.Format_SINGLE_IMAGE
}

func advertRPCAdvertStatusToString(data advertRPC.Status) string {
	switch data {
	case advertRPC.Status_Active:
		return "active"
	case advertRPC.Status_In_Active:
		return "in_active"
	case advertRPC.Status_Completed:
		return "complited"
	}

	return "paused"
}

func jobAdvertInputToAdvertRPCJob(data AdvertJobInput) *advertRPC.Advert {
	advert := advertRPC.Advert{
		StartDate: data.Start_date,
		Name:      data.Name,
		Currency:  data.Currency,
	}

	return &advert
}

func candidateAdvertInputToAdvertRPCCandidate(data AdvertCandidateInput) *advertRPC.Advert {
	advert := advertRPC.Advert{
		StartDate: data.Start_date,
		Name:      data.Name,
		Currency:  data.Currency,
	}

	return &advert
}

func advertRPCAdvertTypeToString(data advertRPC.Type) string {
	switch data {
	case advertRPC.Type_Candidate_Type:
		return "candidate"
	case advertRPC.Type_Job_Type:
		return "job"
	case advertRPC.Type_Shop_Type:
		return "shop"
	case advertRPC.Type_Office_Type:
		return "office"
	case advertRPC.Type_Auto_Type:
		return "auto"
	case advertRPC.Type_Real_Estate_Type:
		return "real_estate"
	case advertRPC.Type_Product_Type:
		return "product"
	case advertRPC.Type_Service_Type:
		return "service"
	case advertRPC.Type_Organization_Type:
		return "organization"
	case advertRPC.Type_Company_Type:
		return "company"
	case advertRPC.Type_Professional_Type:
		return "professional"

	}

	return "banner"
}

func advertTypeToRPC(data string) advertRPC.Type {
	switch data {
	case "candidate":
		return advertRPC.Type_Candidate_Type
	case "job":
		return advertRPC.Type_Job_Type
	case "shop":
		return advertRPC.Type_Shop_Type
	case "office":
		return advertRPC.Type_Office_Type
	case "brand":
		return advertRPC.Type_Brand_Type
	case "auto":
		return advertRPC.Type_Auto_Type
	case "real_estate":
		return advertRPC.Type_Real_Estate_Type
	case "product":
		return advertRPC.Type_Product_Type
	case "service":
		return advertRPC.Type_Service_Type
	case "organization":
		return advertRPC.Type_Organization_Type
	case "company":
		return advertRPC.Type_Company_Type
	case "professional":
		return advertRPC.Type_Professional_Type
	case "banner":
		return advertRPC.Type_Banner_Type
	}

	return advertRPC.Type_Advert_Any
}

func advertRPCBannertoAdvertBannerDraft(data *advertRPC.Banner) *AdvertBannerDraft {
	if data == nil {
		return nil
	}

	ad := AdvertBannerDraft{
		// Is_responsive: data.GetIsResponsive(),
		// Destination_url: data.GetDestinationURL(),
		// Places:          make([]string, 0, len(data.GetPlaces())),
		Contents: make([]AdvertBannerContent, 0, len(data.GetContents())),
	}

	// for _, place := range data.GetPlaces() {
	// 	ad.Places = append(ad.Places, advertPlaceToString(place))
	// }

	for _, c := range data.GetContents() {
		if content := advertContentRPCToAdvertBannerContent(c); content != nil {
			ad.Contents = append(ad.Contents, *content)
		}
	}

	if a := data.GetAdvert(); a != nil {
		ad.Currency = a.GetCurrency()
		ad.Name = a.GetName()
		ad.Start_date = a.GetStartDate()
	}

	return &ad
}

func advertPlaceToString(data advertRPC.Place) string {
	switch data {
	case advertRPC.Place_Company:
		return "company"
	case advertRPC.Place_LocalBusiness:
		return "local_business"
	case advertRPC.Place_Brands:
		return "brands"
	case advertRPC.Place_Groups:
		return "groups"
	}

	return "user"
}

func advertContentRPCToAdvertBannerContent(data *advertRPC.Content) *AdvertBannerContent {
	if data == nil {
		return nil
	}

	content := AdvertBannerContent{
		File_id:     data.GetFileID(),
		Title:       data.GetTitle(),
		Description: data.GetContent(),
		Url:         data.GetImageURL(),
	}

	return &content
}

func advertRPCLocationsToCity(data []*advertRPC.Location) []*City {
	if len(data) <= 0 {
		return nil
	}

	cities := make([]*City, 0, len(data))

	for _, c := range data {
		cities = append(cities, advertRPCLocationToCity(c))
	}

	return cities
}
func advertRPCLocationToCity(data *advertRPC.Location) *City {
	if data == nil {
		return nil
	}

	city := City{
		City:        data.GetCity(),
		Country:     data.GetCountryID(),
		ID:          data.GetCityID(),
		Subdivision: data.GetSubdivision(),
	}

	return &city
}

func advertFormatsToString(data []advertRPC.Format) []string {

	if len(data) <= 0 {
		return nil
	}

	formats := make([]string, 0, len(data))

	for _, f := range data {
		formats = append(formats, advertFormatToString(f))
	}

	return formats
}

func advertFormatToString(data advertRPC.Format) string {
	switch data {
	case advertRPC.Format_CAROUSEL:
		return "CAROUSEL"
	case advertRPC.Format_SLIDE:
		return "SLIDE"
	case advertRPC.Format_VIDEO:
		return "VIDEO"
	case advertRPC.Format_IMAGE:
		return "IMAGE"
	case advertRPC.Format_RESPONSIVE:
		return "RESPONSIVE"
	case advertRPC.Format_SPOTLIGHT:
		return "SPOTLIGHT"
	case advertRPC.Format_SIDE_PIN:
		return "SIDE_PIN"
	case advertRPC.Format_HEAD_PIN:
		return "HEAD_PIN"
	case advertRPC.Format_BUSINESS_SEARCH:
		return "BUSINESS_SEARCH"
	case advertRPC.Format_PROFESSIONAL_SEARCH:
		return "PROFESSIONAL_SEARCH"
	case advertRPC.Format_CANDIDATE_SEARCH:
		return "CANDIDATE_SEARCH"
	case advertRPC.Format_JOB_SEARCH:
		return "JOB_SEARCH"
	case advertRPC.Format_SERVICE_SEARCH:
		return "SERVICE_SEARCH"
	}

	return "SINGLE_IMAGE"
}
