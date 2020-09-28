package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/advert"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/file"
)

// Repository ...
type Repository interface {
	SaveBanner(context.Context, *advert.Banner) error
	GetBannerByID(ctx context.Context, bannerID string) (*advert.Banner, error)
	ChangeBanner(ctx context.Context, ad *advert.Banner) error
	SaveJob(context.Context, *advert.Job) error
	SaveCandidate(ctx context.Context, ad *advert.Candidate) error

	ChangeAdvert(ctx context.Context, ad *advert.Advert) error
	GetAdvertByID(ctx context.Context, advertID string) (*advert.Advert, error)

	AddImageToGallery(context.Context, string, *file.File) error
	GetImageURLByID(ctx context.Context, fileID string) (string, error)
	GetGallery(ctx context.Context, userID, companyID string, first, after uint32) ([]*file.File, uint32, error)
	GetMyAdverts(ctx context.Context, userID, companyID string, first, after uint32) ([]*advert.Advert, uint32, error)
	RemoveAdvert(ctx context.Context, campaignID, advertID string) error
	RemoveAdvertCampaign(ctx context.Context, campaignID string) error

	GetBanners(ctx context.Context, countryID string /*, place advert.Place*/, format advert.Format, amount uint32) ([]*advert.Banner, error)
	GetCandidates(ctx context.Context, countryID string, format advert.Format, amount uint32) ([]string, error)
	GetJobs(ctx context.Context, countryID string, amount uint32) ([]string, error)

	CreateAdvertCampaign(ctx context.Context, data advert.Campaing) error

	CreateAdvertByCampaign(ctx context.Context, campaignID string, data advert.CampaingAdvert) error
	GetAdvertCampaigns(ctx context.Context, profileID string, first int, after int) (*advert.Campaings, error)
	GetAdvertsByCampaignID(ctx context.Context, campaignID, profileID string, first int, after int) (*advert.Adverts, error)
	GetAdvert(ctx context.Context, profileID string, advertType advert.Type) (*advert.GetAdvert, error)
	ChangeAdvertActions(ctx context.Context, advertID string, data advert.ActionType) error
	ChangeStatus(ctx context.Context, campaignID, advertID string, data advert.Status) error
}
