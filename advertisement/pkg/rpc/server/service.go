package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/advert"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/file"
)

// Service define functions inside Service
type Service interface {
	SaveBanner(ctx context.Context, ad *advert.Banner) (string, error)
	SaveBannerDraft(ctx context.Context, ad *advert.Banner) (string, error)
	ChangeBanner(ctx context.Context, ad *advert.Banner) error
	GetBannerDraftByID(ctx context.Context, companyID, bannerID string) (*advert.Banner, error)

	Publish(ctx context.Context, companyID string, bannerID string) error
	PutOnPause(ctx context.Context, companyID, bannerID string) error

	AddImageToGallery(ctx context.Context, id string, f *file.File) (string, error)
	GetGallery(ctx context.Context, companyID string, first uint32, after uint32) ([]*file.File, uint32, error)

	SaveJob(ctx context.Context, ad *advert.Job) (string, error)
	SaveCandidate(ctx context.Context, ad *advert.Candidate) (string, error)

	GetMyAdverts(ctx context.Context, companyID string, first uint32, after uint32) ([]*advert.Advert, uint32, error)
	GetBanners(ctx context.Context, countryID string /*place advert.Place,*/, format advert.Format, amount uint32) ([]*advert.Banner, error)
	GetCandidates(ctx context.Context, countryID string, format advert.Format, amount uint32) ([]string, error)
	GetJobs(ctx context.Context, countryID string, amount uint32) ([]string, error)

	RemoveAdvert(ctx context.Context, campaignID, advertID string) error
	RemoveAdvertCampaign(ctx context.Context, campaignID string) error

	CreateAdvertCampaign(ctx context.Context, companyID string, data advert.Campaing) (string, error)
	PauseAdvertCampaign(ctx context.Context, campaignID string) error
	ActiveAdvertCampaign(ctx context.Context, campaignID string) error

	CreateAdvertByCampaign(ctx context.Context, campaignID, typeID string, data advert.CampaingAdvert) (string, error)
	PauseAdvert(ctx context.Context, advertID string) error
	ActiveAdvert(ctx context.Context, advertID string) error

	GetAdvertCampaigns(ctx context.Context, compan string, first uint32, after string) (*advert.Campaings, error)
	GetAdvertsByCampaignID(ctx context.Context, campaignID, companyID string, first uint32, after string) (*advert.Adverts, error)
	GetAdvert(ctx context.Context, advertType advert.Type) (*advert.CampaingAdvert, error)
}
