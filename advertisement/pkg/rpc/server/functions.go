package serverRPC

import (
	"context"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/advert"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/file"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/location"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/advertRPC"
)

// SaveBanner ...
func (s Server) SaveBanner(ctx context.Context, data *advertRPC.Banner) (*advertRPC.ID, error) {
	id, err := s.service.SaveBanner(ctx, advertRPCBannerToAdvertBanner(data))
	if err != nil {
		return nil, err
	}

	return &advertRPC.ID{
		ID: id,
	}, nil
}

// SaveBannerDraft ...
func (s Server) SaveBannerDraft(ctx context.Context, data *advertRPC.Banner) (*advertRPC.ID, error) {
	id, err := s.service.SaveBannerDraft(ctx, advertRPCBannerToAdvertBanner(data))
	if err != nil {
		return nil, err
	}

	return &advertRPC.ID{
		ID: id,
	}, nil
}

// ChangeBanner ...
func (s Server) ChangeBanner(ctx context.Context, data *advertRPC.Banner) (*advertRPC.Empty, error) {
	err := s.service.ChangeBanner(ctx, advertRPCBannerToAdvertBanner(data))
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// GetBannerDraft ...
func (s Server) GetBannerDraft(ctx context.Context, data *advertRPC.IDWithCompanyID) (*advertRPC.Banner, error) {
	banner, err := s.service.GetBannerDraftByID(
		ctx,
		data.GetCompanyID(),
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}
	return advertBannerToAdvertRPCBanner(banner), nil
}

// Publish ...
func (s Server) Publish(ctx context.Context, data *advertRPC.IDWithCompanyID) (*advertRPC.Empty, error) {
	err := s.service.Publish(ctx, data.GetCompanyID(), data.GetID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// PutOnPause ...
func (s Server) PutOnPause(ctx context.Context, data *advertRPC.IDWithCompanyID) (*advertRPC.Empty, error) {
	err := s.service.PutOnPause(ctx, data.GetCompanyID(), data.GetID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// AddImageToGallery ...
func (s Server) AddImageToGallery(ctx context.Context, data *advertRPC.File) (*advertRPC.ID, error) {
	id, err := s.service.AddImageToGallery(ctx, data.GetTargetID(), advertRPCFileToFile(data))
	if err != nil {
		return nil, err
	}

	return &advertRPC.ID{
		ID: id,
	}, nil
}

// GetGallery ...
func (s Server) GetGallery(ctx context.Context, data *advertRPC.GetGalleryRequest) (*advertRPC.Files, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, _ := strconv.Atoi(data.GetPagination().GetFirst())
		first = uint32(f)
		a, _ := strconv.Atoi(data.GetPagination().GetAfter())
		after = uint32(a)
	}

	files, amount, err := s.service.GetGallery(
		ctx,
		data.GetCompanyID(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	filesRPC := make([]*advertRPC.File, 0, len(files))

	for _, f := range files {
		filesRPC = append(filesRPC, fileToadvertRPCFile(f))
	}

	return &advertRPC.Files{
		Files:  filesRPC,
		Amount: amount,
	}, nil
}

// GetMyAdverts ...
func (s Server) GetMyAdverts(ctx context.Context, data *advertRPC.GetMyAdvertsRequest) (*advertRPC.Adverts, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, _ := strconv.Atoi(data.GetPagination().GetFirst())
		first = uint32(f)
		a, _ := strconv.Atoi(data.GetPagination().GetAfter())
		after = uint32(a)
	}

	ads, _, err := s.service.GetMyAdverts(
		ctx,
		data.GetCompanyID(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	adsRPC := make([]*advertRPC.Advert, 0, len(ads))

	for _, a := range ads {
		adsRPC = append(adsRPC, advertAdvertToAdvertRPCAdvert(a))
	}

	return &advertRPC.Adverts{
		Adverts: adsRPC,
	}, nil
}

// SaveJob ...
func (s Server) SaveJob(ctx context.Context, data *advertRPC.Job) (*advertRPC.ID, error) {
	ad := advertRPCAdvertToAdvertAdvert(data.GetAdvert())

	job := advert.Job{}
	if ad != nil {
		job.Advert = *ad
	}

	id, err := s.service.SaveJob(
		ctx,
		&job,
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.ID{
		ID: id,
	}, nil
}

// SaveCandidate ...
func (s Server) SaveCandidate(ctx context.Context, data *advertRPC.Candidate) (*advertRPC.ID, error) {
	ad := advertRPCAdvertToAdvertAdvert(data.GetAdvert())

	cand := advert.Candidate{}
	if ad != nil {
		cand.Advert = *ad
	}

	id, err := s.service.SaveCandidate(
		ctx,
		&cand,
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.ID{
		ID: id,
	}, nil
}

// GetBanners ...
func (s Server) GetBanners(ctx context.Context, data *advertRPC.GetBannersRequest) (*advertRPC.Banners, error) {
	ads, err := s.service.GetBanners(
		ctx,
		data.GetCountryID(),
		// advertRPCPlaceToAdvertPlace(data.GetPlace()),
		advertRPCFromatToFormat(data.GetFormat()),
		uint32(data.GetAmount()),
	)
	if err != nil {
		return nil, err
	}

	adsRPC := make([]*advertRPC.Banner, 0, len(ads))

	for _, a := range ads {
		adsRPC = append(adsRPC, advertBannerToAdvertRPCBanner(a))
	}

	return &advertRPC.Banners{
		Banners: adsRPC,
	}, nil
}

// GetCandidates ...
func (s Server) GetCandidates(ctx context.Context, data *advertRPC.GetCandidatesRequest) (*advertRPC.IDs, error) {
	ids, err := s.service.GetCandidates(
		ctx,
		data.GetCountryID(),
		advertRPCFromatToFormat(data.GetFormat()),
		uint32(data.GetAmount()),
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.IDs{
		IDs: ids,
	}, nil
}

// GetJobs ...
func (s Server) GetJobs(ctx context.Context, data *advertRPC.GetJobsRequest) (*advertRPC.IDs, error) {
	ids, err := s.service.GetJobs(
		ctx,
		data.GetCountryID(),
		uint32(data.GetAmount()),
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.IDs{
		IDs: ids,
	}, nil
}

// RemoveAdvert ...
func (s Server) RemoveAdvert(ctx context.Context, data *advertRPC.PauseAdvertRequest) (*advertRPC.Empty, error) {
	err := s.service.RemoveAdvert(ctx, data.GetCampaignID(), data.GetAdvertID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// RemoveAdvertCampaign ...
func (s Server) RemoveAdvertCampaign(ctx context.Context, data *advertRPC.PauseAdvertRequest) (*advertRPC.Empty, error) {
	err := s.service.RemoveAdvertCampaign(ctx, data.GetCampaignID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// CreateAdvertCampaign ...
func (s Server) CreateAdvertCampaign(ctx context.Context, data *advertRPC.AdvertCampaign) (*advertRPC.ID, error) {
	res, err := s.service.CreateAdvertCampaign(
		ctx,
		data.GetCompanyID(),
		advertCampaignRPCToStruct(data),
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.ID{
		ID: res,
	}, nil
}

// PauseAdvertCampaign ...
func (s Server) PauseAdvertCampaign(ctx context.Context, data *advertRPC.PauseAdvertRequest) (*advertRPC.Empty, error) {
	err := s.service.PauseAdvertCampaign(ctx, data.GetCampaignID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// ActiveAdvertCampaign ...
func (s Server) ActiveAdvertCampaign(ctx context.Context, data *advertRPC.PauseAdvertRequest) (*advertRPC.Empty, error) {
	err := s.service.ActiveAdvertCampaign(ctx, data.GetCampaignID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// PauseAdvert ...
func (s Server) PauseAdvert(ctx context.Context, data *advertRPC.PauseAdvertRequest) (*advertRPC.Empty, error) {
	err := s.service.PauseAdvert(ctx, data.GetAdvertID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// ActiveAdvert ...
func (s Server) ActiveAdvert(ctx context.Context, data *advertRPC.PauseAdvertRequest) (*advertRPC.Empty, error) {
	err := s.service.ActiveAdvert(ctx, data.GetAdvertID())
	if err != nil {
		return nil, err
	}

	return &advertRPC.Empty{}, nil
}

// CreateAdvertByCampaign ...
func (s Server) CreateAdvertByCampaign(ctx context.Context, data *advertRPC.CampaignAdvert) (*advertRPC.ID, error) {
	res, err := s.service.CreateAdvertByCampaign(
		ctx,
		data.GetCampaignID(),
		data.GetTypeID(),
		advertCampaignAdvertRPCToStruct(data),
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.ID{
		ID: res,
	}, nil
}

// GetAdvertCampaigns ...
func (s Server) GetAdvertCampaigns(ctx context.Context, data *advertRPC.AdvertGetRequest) (*advertRPC.Campaigns, error) {
	res, err := s.service.GetAdvertCampaigns(
		ctx,
		data.GetCompanyID(),
		data.GetFirst(),
		data.GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.Campaigns{
		TotalAmount:     res.Amount,
		AdvertCampaigns: advertCampaignsToRPC(res.Campaings),
	}, nil
}

// GetAdvertsByCampaignID ...
func (s Server) GetAdvertsByCampaignID(ctx context.Context, data *advertRPC.AdvertGetRequest) (*advertRPC.Adverts, error) {
	res, err := s.service.GetAdvertsByCampaignID(
		ctx,
		data.GetCampaignID(),
		data.GetCompanyID(),
		data.GetFirst(),
		data.GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	return &advertRPC.Adverts{
		Amount:   res.Amount,
		Campaign: advertCampaignToRPC(res.Campaing),
		Adverts:  advertAdvertsToRPC(res.Adverts),
	}, nil
}

// GetAdvert ...
func (s Server) GetAdvert(ctx context.Context, data *advertRPC.GetAdvertRequest) (*advertRPC.Advert, error) {
	res, err := s.service.GetAdvert(ctx, advertRPCTypeToType(data.GetAdvertType()))

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return advertAdvertToRPC(*res), nil
}

func advertAdvertsToRPC(data []advert.CampaingAdvert) []*advertRPC.Advert {
	if len(data) <= 0 {
		return nil
	}
	ads := make([]*advertRPC.Advert, 0, len(data))

	for _, a := range data {
		ads = append(ads, advertAdvertToRPC(a))
	}

	return ads
}

func advertAdvertToRPC(data advert.CampaingAdvert) *advertRPC.Advert {
	res := &advertRPC.Advert{
		ID:          data.GetAdvertID(),
		TypeID:      data.GetTypeID(),
		Name:        data.Name,
		URL:         data.URL,
		Impressions: data.Impressions,
		Clicks:      data.Clicks,
		Forwarding:  data.Forwarding,
		StartDate:   data.CreatedAt.String(),
		AdType:      advertTypeToAdvertRPCAdvertType(data.AdvertType),
		AdStatus:    advertStatusToAdvertRPCAdvertStatus(data.Status),
		Contents:    advertContentsToRPC(data.Contents),
		Formats:     advertFormatsToRPC(data.Formats),
	}

	filesRPC := make([]*advertRPC.File, 0, len(data.Files))

	for _, f := range data.Files {
		filesRPC = append(filesRPC, fileToadvertRPCFile(&f))
		res.Files = filesRPC
	}

	return res

}

func advertCampaignsToRPC(data []advert.Campaing) []*advertRPC.AdvertCampaign {
	if len(data) <= 0 {
		return nil
	}

	adverts := make([]*advertRPC.AdvertCampaign, 0, len(data))

	for _, a := range data {
		adverts = append(adverts, advertCampaignToRPC(a))
	}

	return adverts
}

func advertCampaignToRPC(data advert.Campaing) *advertRPC.AdvertCampaign {
	return &advertRPC.AdvertCampaign{
		ID:          data.GetID(),
		Clicks:      data.Clicks,
		Impressions: data.Impressions,
		Name:        data.Name,
		Referals:    data.Referals,
		Forwarding:  data.Forwarding,
		CrtAVG:      data.CtrAVG,
		Languages:   data.Languages,
		StartDate:   data.StartDate.String(),
		AgeFrom:     nullToInt32(data.AgeFrom),
		AgeTo:       nullToInt32(data.AgeTo),
		Gender:      advertGenderToRPC(data.Gender),
		Status:      advertStatusToAdvertRPCAdvertStatus(data.Status),
		AdvertType:  advertTypeToAdvertRPCAdvertType(data.Type),
		Formats:     advertFormatsToRPC(data.Formats),
		Locations:   advertLoctionsToRPC(data.Locations),
	}
}

func advertFormatsToRPC(data []advert.Format) []advertRPC.Format {
	if len(data) <= 0 {
		return nil
	}

	formats := make([]advertRPC.Format, 0, len(data))

	for _, f := range data {
		formats = append(formats, advertFormatToRPC(f))
	}

	return formats
}

func advertGenderToRPC(data *string) advertRPC.GenderEnum {
	if data == nil {
		return advertRPC.GenderEnum_WITHOUT_GENDER
	}

	if *data == "MALE" {
		return advertRPC.GenderEnum_MALE
	}

	return advertRPC.GenderEnum_FEMALE
}

func advertCampaignAdvertRPCToStruct(data *advertRPC.CampaignAdvert) (r advert.CampaingAdvert) {
	if data == nil {
		return
	}

	cmAdvert := advert.CampaingAdvert{
		Name:       data.GetName(),
		URL:        data.GetURL(),
		Status:     advert.StatusActive,
		AdvertType: advertRPCTypeToType(data.GetAdvertType()),
		Contents:   advertContentsToAdvertRPCContents(data.GetContents()),
	}

	return cmAdvert

}

func advertContentsToAdvertRPCContents(data []*advertRPC.Content) *[]advert.Content {
	if data == nil {
		return nil
	}

	contents := make([]advert.Content, 0, len(data))

	for _, c := range data {
		contents = append(contents, *advertRPCContentToAdvertContent(c))
	}

	return &contents
}

func advertCampaignRPCToStruct(data *advertRPC.AdvertCampaign) (r advert.Campaing) {
	if data == nil {
		return
	}

	ad := advert.Campaing{
		Clicks:      data.GetClicks(),
		Currency:    data.GetCurrency(),
		Forwarding:  data.GetForwarding(),
		Impressions: data.GetImpressions(),
		Languages:   data.GetLanguages(),
		Name:        data.GetName(),
		Referals:    data.GetReferals(),
		StartDate:   stringDateToTime(data.GetStartDate()),
		Locations:   make([]location.Location, 0, len(data.GetLocations())),
		Type:        advertRPCTypeToType(data.GetAdvertType()),
		Formats:     advertRPCFromatsToFormats(data.GetFormats()),
	}

	if data.GetAgeFrom() != 0 {
		ageFrom := data.GetAgeFrom()
		ad.AgeFrom = &ageFrom
	}

	if data.GetAgeTo() != 0 {
		ageTo := data.GetAgeTo()
		ad.AgeTo = &ageTo
	}

	if data.GetGender().String() != "WITHOUT_GENDER" {
		gender := data.GetGender().String()
		ad.Gender = &gender
	}

	for _, locs := range data.GetLocations() {
		if loc := advertRPCLocationToLocation(locs); loc != nil {
			ad.Locations = append(ad.Locations, *loc)
		}
	}

	return ad
}

// ---------------------------

func advertRPCFromatsToFormats(data []advertRPC.Format) []advert.Format {
	if len(data) <= 0 {
		return nil
	}

	formats := make([]advert.Format, 0, len(data))

	for _, f := range data {
		formats = append(formats, advertRPCFromatToFormat(f))
	}

	return formats
}

func advertRPCBannerToAdvertBanner(data *advertRPC.Banner) *advert.Banner {
	if data == nil {
		return nil
	}

	banner := advert.Banner{
		// DestinationURL: data.GetDestinationURL(),
		IsResponsive: data.GetIsResponsive(),
		// Places:         make([]advert.Place, 0, len(data.GetPlaces())),
		Contents: make([]advert.Content, 0, len(data.GetContents())),
	}

	if ad := advertRPCAdvertToAdvertAdvert(data.GetAdvert()); ad != nil {
		banner.Advert = *ad
	}

	// for _, place := range data.GetPlaces() {
	// 	banner.Places = append(banner.Places, advertRPCPlaceToAdvertPlace(place))
	// }

	// TODO: contents
	for _, content := range data.GetContents() {
		if c := advertRPCContentToAdvertContent(content); c != nil {
			banner.Contents = append(banner.Contents, *c)
		}
	}

	return &banner
}

// func advertRPCPlaceToAdvertPlace(data advertRPC.Place) advert.Place {
// 	switch data {
// 	case advertRPC.Place_User:
// 		return advert.PlaceUser
// 	case advertRPC.Place_Company:
// 		return advert.PlaceCompany
// 	case advertRPC.Place_LocalBusiness:
// 		return advert.PlaceLocalBusiness
// 	case advertRPC.Place_Brands:
// 		return advert.PlaceBrands
// 	case advertRPC.Place_Groups:
// 		return advert.PlaceGroups
// 	}
//
// 	return advert.PlaceUser
// }

func advertRPCAdvertToAdvertAdvert(data *advertRPC.Advert) *advert.Advert {
	if data == nil {
		return nil
	}

	ad := advert.Advert{
		Name:      data.GetName(),
		Currency:  data.GetCurrency(),
		Status:    advertRPCAdvertStatusToAdvertStatus(data.GetAdStatus()),
		StartDate: stringDateToTime(data.GetStartDate()),
		Button:    data.GetButtonTitle(),
	}

	_ = ad.SetID(data.GetID())
	_ = ad.SetCompanyID(data.GetCompanyID())

	finishedDate := stringDateToTime(data.GetFinishDate())
	if !finishedDate.IsZero() {
		ad.FinishDate = &finishedDate
	}

	return &ad
}

func advertRPCAdvertStatusToAdvertStatus(data advertRPC.Status) advert.Status {
	switch data {
	case advertRPC.Status_Paused:
		return advert.StatusPaused
	case advertRPC.Status_Completed:
		return advert.StatusCompleted
	}

	return advert.StatusDraft
}

func advertRPCLocationToLocation(data *advertRPC.Location) *location.Location {
	if data == nil {
		return nil
	}

	loc := location.Location{
		CountryID: data.GetCountryID(),
	}

	if city := data.GetCity(); city != "" {
		loc.City = &city
	}

	if cityID := data.GetCityID(); cityID != "" {
		loc.CityID = &cityID
	}

	if subdivision := data.GetSubdivision(); subdivision != "" {
		loc.Subdivision = &subdivision
	}

	return &loc
}

func advertLoctionsToRPC(data []location.Location) []*advertRPC.Location {
	if len(data) <= 0 {
		return nil
	}

	locations := make([]*advertRPC.Location, 0, len(data))

	for _, l := range data {
		locations = append(locations, &advertRPC.Location{
			City:      nullToString(l.City),
			CountryID: l.CountryID,
		})
	}
	return locations
}
func advertRPCContentToAdvertContent(data *advertRPC.Content) (r *advert.Content) {
	if data == nil {
		return nil
	}

	content := advert.Content{
		Content:        data.GetContent(),
		Title:          data.GetTitle(),
		DestinationURL: data.GetDestinationURL(),
		CustomButton:   data.GetCustomButton(),
	}

	id := content.GenerateID()
	content.SetID(id)

	return &content
}

func advertRPCFileToFile(data *advertRPC.File) *file.File {
	if data == nil {
		return nil
	}

	f := file.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		URL:      data.GetURL(),
	}

	_ = f.SetID(data.GetID())
	_ = f.SetUserID(data.GetUserID())

	if data.GetCompanyID() != "" {
		_ = f.SetCompanyID(data.GetCompanyID())
	}

	return &f
}

func advertRPCFromatToFormat(data advertRPC.Format) advert.Format {
	switch data {
	case advertRPC.Format_CAROUSEL:
		return advert.FormatCarousel
	case advertRPC.Format_SLIDE:
		return advert.FormatSlide
	case advertRPC.Format_VIDEO:
		return advert.FormatVideo
	case advertRPC.Format_IMAGE:
		return advert.FormatImage
	case advertRPC.Format_RESPONSIVE:
		return advert.FormatResponsive
	case advertRPC.Format_SPOTLIGHT:
		return advert.FormatSpotlight
	case advertRPC.Format_SIDE_PIN:
		return advert.FormatSidePin
	case advertRPC.Format_HEAD_PIN:
		return advert.FormatHeadPin
	case advertRPC.Format_BUSINESS_SEARCH:
		return advert.FormatBusinessSearch
	case advertRPC.Format_PROFESSIONAL_SEARCH:
		return advert.FormatProfessionalSearch
	case advertRPC.Format_CANDIDATE_SEARCH:
		return advert.FormatCandidatesearch
	case advertRPC.Format_JOB_SEARCH:
		return advert.FormatJobsearch
	case advertRPC.Format_SERVICE_SEARCH:
		return advert.FormatServiceSearch
	}

	return advert.FormatSingleImage
}

func advertFormatToRPC(data advert.Format) advertRPC.Format {
	switch data {
	case advert.FormatCarousel:
		return advertRPC.Format_CAROUSEL
	case advert.FormatSlide:
		return advertRPC.Format_SLIDE
	case advert.FormatVideo:
		return advertRPC.Format_VIDEO
	case advert.FormatImage:
		return advertRPC.Format_IMAGE
	case advert.FormatResponsive:
		return advertRPC.Format_RESPONSIVE
	case advert.FormatSpotlight:
		return advertRPC.Format_SPOTLIGHT
	case advert.FormatSidePin:
		return advertRPC.Format_SIDE_PIN
	case advert.FormatHeadPin:
		return advertRPC.Format_HEAD_PIN
	case advert.FormatBusinessSearch:
		return advertRPC.Format_BUSINESS_SEARCH
	case advert.FormatProfessionalSearch:
		return advertRPC.Format_PROFESSIONAL_SEARCH
	case advert.FormatCandidatesearch:
		return advertRPC.Format_CANDIDATE_SEARCH
	case advert.FormatJobsearch:
		return advertRPC.Format_JOB_SEARCH
	case advert.FormatServiceSearch:
		return advertRPC.Format_SERVICE_SEARCH
	}

	return advertRPC.Format_SINGLE_IMAGE
}

// ---------------------------

func fileToadvertRPCFile(data *file.File) *advertRPC.File {
	if data == nil {
		return nil
	}

	f := advertRPC.File{
		ID:        data.GetID(),
		UserID:    data.GetUserID(),
		CompanyID: data.GetCompanyID(),
		MimeType:  data.MimeType,
		Name:      data.Name,
		URL:       data.URL,
	}

	return &f
}

func advertAdvertToAdvertRPCAdvert(data *advert.Advert) *advertRPC.Advert {
	if data == nil {
		return nil
	}

	ad := advertRPC.Advert{
		Name:        data.Name,
		ButtonTitle: data.Button,
		Currency:    data.Currency,
		AdStatus:    advertStatusToAdvertRPCAdvertStatus(data.Status),
		AdType:      advertTypeToAdvertRPCAdvertType(data.Type),
		StartDate:   timeToStringDayMonthAndYear(data.StartDate),
	}

	ad.ID = data.GetID()

	ad.CompanyID = data.GetCompanyID()
	ad.CreatorID = data.GetCreatorID()

	if data.FinishDate != nil {
		finishedDate := timeToStringDayMonthAndYear(*data.FinishDate)
		if finishedDate != "" {
			ad.FinishDate = finishedDate
		}
	}

	return &ad
}

func advertStatusToAdvertRPCAdvertStatus(data advert.Status) advertRPC.Status {
	switch data {
	case advert.StatusActive:
		return advertRPC.Status_Active
	case advert.StatusPaused:
		return advertRPC.Status_Paused
	case advert.StatusCompleted:
		return advertRPC.Status_Completed
	}

	return advertRPC.Status_In_Active
}

func advertTypeToAdvertRPCAdvertType(data advert.Type) advertRPC.Type {
	switch data {
	case advert.TypeJob:
		return advertRPC.Type_Job_Type
	case advert.TypeCandidate:
		return advertRPC.Type_Candidate_Type
	case advert.TypeOffice:
		return advertRPC.Type_Office_Type
	case advert.TypeShop:
		return advertRPC.Type_Shop_Type
	case advert.TypeBrand:
		return advertRPC.Type_Brand_Type
	case advert.TypeAuto:
		return advertRPC.Type_Auto_Type
	case advert.TypeRealEstate:
		return advertRPC.Type_Real_Estate_Type
	case advert.TypeProduct:
		return advertRPC.Type_Product_Type
	case advert.TypeService:
		return advertRPC.Type_Service_Type
	case advert.TypeOrganization:
		return advertRPC.Type_Organization_Type
	case advert.TypeCompany:
		return advertRPC.Type_Company_Type
	case advert.TypeProffesional:
		return advertRPC.Type_Professional_Type
	}

	return advertRPC.Type_Banner_Type
}

func advertRPCTypeToType(data advertRPC.Type) advert.Type {

	switch data {
	case advertRPC.Type_Job_Type:
		return advert.TypeJob
	case advertRPC.Type_Candidate_Type:
		return advert.TypeCandidate
	case advertRPC.Type_Office_Type:
		return advert.TypeOffice
	case advertRPC.Type_Shop_Type:
		return advert.TypeShop
	case advertRPC.Type_Brand_Type:
		return advert.TypeBrand
	case advertRPC.Type_Auto_Type:
		return advert.TypeAuto
	case advertRPC.Type_Real_Estate_Type:
		return advert.TypeRealEstate
	case advertRPC.Type_Product_Type:
		return advert.TypeProduct
	case advertRPC.Type_Service_Type:
		return advert.TypeService
	case advertRPC.Type_Organization_Type:
		return advert.TypeOrganization
	case advertRPC.Type_Company_Type:
		return advert.TypeCompany
	case advertRPC.Type_Professional_Type:
		return advert.TypeProffesional
	case advertRPC.Type_Banner_Type:
		return advert.TypeBanner
	}

	return advert.TypeAny
}

func locationToAdvertRPCLocation(data *location.Location) *advertRPC.Location {
	if data == nil {
		return nil
	}

	loc := advertRPC.Location{
		CountryID: data.CountryID,
	}

	if city := data.City; city != nil {
		loc.City = *city
	}

	if cityID := data.CityID; cityID != nil {
		loc.CityID = *cityID
	}

	if subdivision := data.Subdivision; subdivision != nil {
		loc.Subdivision = *subdivision
	}

	return &loc
}

func advertBannerToAdvertRPCBanner(data *advert.Banner) *advertRPC.Banner {
	if data == nil {
		return nil
	}

	banner := advertRPC.Banner{
		Advert: advertAdvertToAdvertRPCAdvert(&data.Advert),
		// DestinationURL: data.DestinationURL,
		IsResponsive: data.IsResponsive,
		Contents:     make([]*advertRPC.Content, 0, len(data.Contents)),
		// Places
	}

	for _, c := range data.Contents {
		banner.Contents = append(banner.Contents, advertContentToAdvertRPCContent(&c))
	}

	return &banner
}

func advertContentsToRPC(data *[]advert.Content) []*advertRPC.Content {
	if data == nil {
		return nil
	}

	contents := make([]*advertRPC.Content, 0, len(*data))

	for _, c := range *data {
		contents = append(contents, advertContentToAdvertRPCContent(&c))
	}

	return contents
}

func advertContentToAdvertRPCContent(data *advert.Content) *advertRPC.Content {
	if data == nil {
		return nil
	}

	content := advertRPC.Content{
		DestinationURL: data.DestinationURL,
		ImageURL:       data.ImageURL,
		Title:          data.Title,
		Content:        data.Content,
		CustomButton:   data.CustomButton,
	}

	return &content
}

// ---------------------------

func stringDateToTime(s string) time.Time {
	if date, err := time.Parse("2-1-2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func stringDayMonthAndYearToTime(s string) time.Time {
	if date, err := time.Parse("1-2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func stringYearToDate(s string) time.Time {
	if date, err := time.Parse("2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func timeToStringMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, _ := t.UTC().Date()
	return strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func timeToStringDayMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, d := t.UTC().Date()
	return strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func nullToInt32(data *int32) int32 {
	if data == nil {
		return 0
	}

	return *data
}

func nullToString(data *string) string {
	if data == nil {
		return ""
	}

	return *data
}
