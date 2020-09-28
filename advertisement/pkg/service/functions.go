package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/advert"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/file"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
)

// SaveBanner ...
func (s Service) SaveBanner(ctx context.Context, ad *advert.Banner) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveBanner")
	defer span.Finish()

	id, err := s.saveBanner(ctx, advert.StatusActive, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// SaveBannerDraft ...
func (s Service) SaveBannerDraft(ctx context.Context, ad *advert.Banner) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveBanner")
	defer span.Finish()

	id, err := s.saveBanner(ctx, advert.StatusDraft, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

func (s Service) saveBanner(ctx context.Context, t advert.Status, ad *advert.Banner) (string, error) {
	span := s.tracer.MakeSpan(ctx, "saveBanner")
	defer span.Finish()

	err := ad.Validate()
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	ad.Trim()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}
	ad.SetCreatorID(userID)
	id := ad.GenerateID()
	ad.Type = advert.TypeBanner

	if t == advert.StatusActive {
		ad.CalculateDates()
	} else {
		ad.Status = advert.StatusDraft
	}

	err = s.repository.SaveBanner(ctx, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	return id, nil
}

// ChangeBanner ...
func (s Service) ChangeBanner(ctx context.Context, ad *advert.Banner) error {
	span := s.tracer.MakeSpan(ctx, "ChangeBanner")
	defer span.Finish()

	err := ad.Validate()
	if err != nil {
		return err
	}
	ad.Trim()

	banner, err := s.repository.GetBannerByID(ctx, ad.GetID())
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if banner.Status != advert.StatusDraft {
		return errors.New(`not_allowed`)
	}

	err = s.repository.ChangeBanner(ctx, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetBannerDraftByID ...
func (s Service) GetBannerDraftByID(ctx context.Context, companyID, bannerID string) (*advert.Banner, error) {
	span := s.tracer.MakeSpan(ctx, "GetBannerDraftByID")
	defer span.Finish()

	banner, err := s.repository.GetBannerByID(ctx, bannerID)
	if err != nil {
		return nil, err
	}

	if banner.Status != advert.StatusDraft {
		return nil, errors.New(`not_found`)
	}

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}

	// TODO: check if admin

	if companyID != "" &&
		(banner.GetCompanyID() != companyID) {
		return nil, errors.New(`not_found`)
	}
	if banner.GetCreatorID() != userID {
		return nil, errors.New(`not_found`)
	}

	return banner, nil
}

// Publish ...
func (s Service) Publish(ctx context.Context, companyID string, advertID string) error {
	span := s.tracer.MakeSpan(ctx, "Publish")
	defer span.Finish()

	// TODO: check if admin of company

	ad, err := s.repository.GetAdvertByID(ctx, advertID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if !(ad.Type == advert.TypeJob ||
		ad.Type == advert.TypeBanner ||
		ad.Type == advert.TypeCandidate) {
		return errors.New(`not_allowed`)
	}

	if ad.Status != advert.StatusDraft &&
		ad.Status != advert.StatusPaused &&
		ad.Status != advert.StatusNotRunning {
		return errors.New(`not_allowed`)
	}

	switch ad.Status {

	case advert.StatusPaused:
		var pausedTime time.Duration
		// add previous paused time
		if ad.PausedTime != nil {
			pausedTime = *ad.PausedTime
		}
		// add current paused time
		if ad.LastPausedTime != nil {
			pausedTime = pausedTime + time.Duration(ad.LastPausedTime.Nanosecond())
		}
		ad.PausedTime = &pausedTime
		finishDate := ad.StartDate.Add(pausedTime)
		ad.FinishDate = &finishDate

	case advert.StatusDraft:
		ad.StartDate = time.Now()
		ad.CalculateDates()
	}

	ad.Status = advert.StatusActive

	err = s.repository.ChangeAdvert(ctx, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// PutOnPause ...
func (s Service) PutOnPause(ctx context.Context, companyID, advertID string) error {
	span := s.tracer.MakeSpan(ctx, "PutOnPause")
	defer span.Finish()

	// TODO: check if admin

	ad, err := s.repository.GetAdvertByID(ctx, advertID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if ad.Status != advert.StatusActive {
		return errors.New(`not_allowed`)
	}

	if !(ad.Type == advert.TypeJob ||
		ad.Type == advert.TypeBanner ||
		ad.Type == advert.TypeCandidate) {
		return errors.New(`not_allowed`)
	}

	now := time.Now()
	ad.LastPausedTime = &now
	ad.Status = advert.StatusPaused

	err = s.repository.ChangeAdvert(ctx, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddImageToGallery ...
func (s Service) AddImageToGallery(ctx context.Context, campaignID string, f *file.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "AddImageToGallery")
	defer span.Finish()

	id := f.GenerateID()

	err := s.repository.AddImageToGallery(ctx, campaignID, f)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

func (s Service) getImageURLByID(ctx context.Context, id string) (string, error) {
	span := s.tracer.MakeSpan(ctx, "GetImageURLByID")
	defer span.Finish()

	url, err := s.repository.GetImageURLByID(ctx, id)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", nil
	}

	return url, nil
}

// GetGallery ...
func (s Service) GetGallery(ctx context.Context, companyID string, first uint32, after uint32) ([]*file.File, uint32, error) {
	span := s.tracer.MakeSpan(ctx, "GetGallery")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, errors.New("token_is_empty")
	}
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, errors.New("token_is_empty")
	}

	// TODO: check company permissions

	files, amount, err := s.repository.GetGallery(ctx, userID, companyID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return files, amount, nil
}

// GetMyAdverts ...
func (s Service) GetMyAdverts(ctx context.Context, companyID string, first uint32, after uint32) ([]*advert.Advert, uint32, error) {
	span := s.tracer.MakeSpan(ctx, "GetMyAdverts")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, errors.New("token_is_empty")
	}
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, errors.New("token_is_empty")
	}

	// TODO: check company permissions

	ads, amount, err := s.repository.GetMyAdverts(ctx, userID, companyID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return ads, amount, nil
}

// SaveJob ...
func (s Service) SaveJob(ctx context.Context, ad *advert.Job) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveJob")
	defer span.Finish()

	err := ad.Validate()
	if err != nil {
		return "", err
	}
	ad.Trim()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}

	// TODO: check if admin

	ad.SetCreatorID(userID)
	id := ad.GenerateID()
	ad.Type = advert.TypeJob

	ad.CalculateDates()

	err = s.repository.SaveJob(ctx, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	return id, nil
}

// SaveCandidate ...
func (s Service) SaveCandidate(ctx context.Context, ad *advert.Candidate) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveJob")
	defer span.Finish()

	err := ad.Validate()
	if err != nil {
		return "", err
	}
	ad.Trim()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}

	ad.SetCreatorID(userID)
	id := ad.GenerateID()
	ad.Type = advert.TypeCandidate

	ad.CalculateDates()

	err = s.repository.SaveCandidate(ctx, ad)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	return id, nil
}

// GetBanners ...
func (s Service) GetBanners(ctx context.Context, countryID string /*, place advert.Place*/, format advert.Format, amount uint32) ([]*advert.Banner, error) {
	span := s.tracer.MakeSpan(ctx, "GetBanners")
	defer span.Finish()

	banners, err := s.repository.GetBanners(ctx, countryID /*place,*/, format, amount)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return banners, nil
}

// GetCandidates ...
func (s Service) GetCandidates(ctx context.Context, countryID string, format advert.Format, amount uint32) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "GetCandidates")
	defer span.Finish()

	ids, err := s.repository.GetCandidates(ctx, countryID, format, amount)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}
	return ids, nil
}

// GetJobs ...
func (s Service) GetJobs(ctx context.Context, countryID string, amount uint32) ([]string, error) {
	span := s.tracer.MakeSpan(ctx, "GetJobs")
	defer span.Finish()

	ids, err := s.repository.GetJobs(ctx, countryID, amount)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return ids, nil
}

// RemoveAdvert ...
func (s Service) RemoveAdvert(ctx context.Context, campaignID, advertID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveAdvert")
	defer span.Finish()

	// TODO: check if admin

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}
	_, err = s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	err = s.repository.RemoveAdvert(ctx, campaignID, advertID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// RemoveAdvert ...
func (s Service) RemoveAdvertCampaign(ctx context.Context, campaignID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveAdvertCampaign")
	defer span.Finish()

	// TODO: check if admin

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}
	_, err = s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	err = s.repository.RemoveAdvertCampaign(ctx, campaignID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// CreateAdvertCampaign ...
func (s Service) CreateAdvertCampaign(ctx context.Context, companyID string, data advert.Campaing) (string, error) {
	span := s.tracer.MakeSpan(ctx, "CreateAdvertCampaign")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}

	profileID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}

	if companyID != "" {
		profileID = companyID
		data.IsComapny = true
	}

	id := data.GenerateID()
	data.SetID(id)
	data.SetCreatorID(profileID)
	data.Adverts = make([]advert.CampaingAdvert, 0)

	ids, getUserErr := s.userRPC.GetUsersForAdvert(ctx, toUserAdvert(data))

	if getUserErr == nil {
		users := make([]primitive.ObjectID, 0, len(ids))
		for i := range ids {
			id, idErr := primitive.ObjectIDFromHex(ids[i])
			if idErr == nil {
				users = append(users, id)
			}
		}

		data.RelevantUsers = users
	}

	if time.Now().After(data.StartDate) {
		data.Status = advert.StatusActive
	} else {
		data.Status = advert.StatusInActive
	}

	err = s.repository.CreateAdvertCampaign(ctx, data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// PauseAdvertCampaign ...
func (s Service) PauseAdvertCampaign(ctx context.Context, campaignID string) error {
	span := s.tracer.MakeSpan(ctx, "PauseAdvertCampaign")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	_, err = s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	err = s.repository.ChangeStatus(ctx, campaignID, "", advert.StatusPaused)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ActiveAdvertCampaign ...
func (s Service) ActiveAdvertCampaign(ctx context.Context, campaignID string) error {
	span := s.tracer.MakeSpan(ctx, "ActiveAdvertCampaign")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	_, err = s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	err = s.repository.ChangeStatus(ctx, campaignID, "", advert.StatusActive)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// PauseAdvert ...
func (s Service) PauseAdvert(ctx context.Context, advertID string) error {
	span := s.tracer.MakeSpan(ctx, "PauseAdvert")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	_, err = s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	err = s.repository.ChangeStatus(ctx, "", advertID, advert.StatusPaused)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ActiveAdvert ...
func (s Service) ActiveAdvert(ctx context.Context, advertID string) error {
	span := s.tracer.MakeSpan(ctx, "ActiveAdvert")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	_, err = s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return errors.New("token_is_empty")
	}

	err = s.repository.ChangeStatus(ctx, "", advertID, advert.StatusActive)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// CreateAdvertByCampaign ...
func (s Service) CreateAdvertByCampaign(ctx context.Context, campaignID, typeID string, data advert.CampaingAdvert) (string, error) {
	span := s.tracer.MakeSpan(ctx, "CreateAdvertByCampaign")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}

	_, err = s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", errors.New("token_is_empty")
	}

	if typeID != "" {
		data.SetTypeID(typeID)
	}

	id := data.GenerateAdvertID()
	data.SetAdvertID(id)
	data.CreatedAt = time.Now()

	err = s.repository.CreateAdvertByCampaign(ctx, campaignID, data)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// GetAdvertCampaigns ...
func (s Service) GetAdvertCampaigns(ctx context.Context, companyID string, first uint32, after string) (*advert.Campaings, error) {
	span := s.tracer.MakeSpan(ctx, "GetAdvertCampaigns")
	defer span.Finish()

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}

	profileID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}

	if companyID != "" {
		profileID = companyID
	}

	res, err := s.repository.GetAdvertCampaigns(ctx, profileID, int(first), int(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return res, nil
}

// GetAdvertsByCampaignID ...
func (s Service) GetAdvertsByCampaignID(ctx context.Context, campaignID, companyID string, first uint32, after string) (*advert.Adverts, error) {
	span := s.tracer.MakeSpan(ctx, "GetAdvertsByCampaignID")
	defer span.Finish()

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}

	profileID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}

	if companyID != "" {
		profileID = companyID
	}

	res, err := s.repository.GetAdvertsByCampaignID(ctx, campaignID, profileID, int(first), int(afterNumber))
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return res, nil
}

// GetAdvert ...
func (s Service) GetAdvert(ctx context.Context, advertType advert.Type) (*advert.CampaingAdvert, error) {
	span := s.tracer.MakeSpan(ctx, "GetAdvert")
	defer span.Finish()

	token, err := s.extractToken(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}

	profileID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, errors.New("token_is_empty")
	}

	res, err := s.repository.GetAdvert(ctx, profileID, advertType)

	if err != nil {
		return nil, err
	}

	if res.Impressions <= 0 || res.Clicks <= 0 {
		err = s.ChangeStatus(ctx, "", res.Adverts.GetAdvertID(), advert.StatusCompleted)
	} else {
		err = s.ChangeAdvert(ctx, res.Adverts.GetAdvertID(), advert.ActionTypeImpressions)
	}

	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return &res.Adverts, nil

}

// ChangeAdvert ...
func (s Service) ChangeAdvert(ctx context.Context, advertID string, data advert.ActionType) error {
	span := s.tracer.MakeSpan(ctx, "ChangeAdvert")
	defer span.Finish()

	err := s.repository.ChangeAdvertActions(ctx, advertID, data)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeStatus ...
func (s Service) ChangeStatus(ctx context.Context, campaignID, advertID string, data advert.Status) error {
	span := s.tracer.MakeSpan(ctx, "ChangeStatus")
	defer span.Finish()

	err := s.repository.ChangeStatus(ctx, campaignID, advertID, data)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// --------------------------

func toUserAdvert(data advert.Campaing) account.UserForAdvert {
	// Locations
	locations := make([]string, 0, len(data.Locations))
	for _, l := range data.Locations {
		locations = append(locations, strings.ToUpper(l.CountryID))
	}

	res := account.UserForAdvert{
		OwnerID:   data.GetCreatorID(),
		Locations: locations,
		Languages: data.Languages,
	}

	// Gender
	if data.Gender != nil {
		res.Gender = *data.Gender
	}
	// Age From
	if data.AgeFrom != nil {
		res.AgeFrom = int32(time.Now().Year()) - *data.AgeFrom
	}
	// Age to
	if data.AgeTo != nil {

		res.AgeTo = int32(time.Now().Year()) - *data.AgeTo
	}

	return res

}

func (s *Service) extractToken(ctx context.Context) (string, error) {
	span := opentracing.SpanFromContext(ctx)
	span = span.Tracer().StartSpan("extractToken", opentracing.ChildOf(span.Context()))
	defer span.Finish()

	res := make(map[string]string, 1)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			res["token"] = arr[0]
		}
	}

	token, ok := res["token"]
	if !ok {
		return "", errors.New("token is empty")
	}
	return token, nil
}
