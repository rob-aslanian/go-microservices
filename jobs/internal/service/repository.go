package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	careercenter "gitlab.lan/Rightnao-site/microservices/jobs/internal/career-center"
	companyadmin "gitlab.lan/Rightnao-site/microservices/jobs/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job"
	notmes "gitlab.lan/Rightnao-site/microservices/jobs/internal/notification_messages"
)

// AuthRPC represents auth service
type AuthRPC interface {
	GetUserID(ctx context.Context) (string, error)
}

// NetworkRPC represents network service
type NetworkRPC interface {
	GetAdminLevel(ctx context.Context, companyID string) (companyadmin.AdminLevel, error)
}

// InfoRPC represents info service
type InfoRPC interface {
	GetUserCountry(ctx context.Context) (string, error)
	GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error)
}

// JobsRepository represents jobs repository
type JobsRepository interface {
	GetCandidateProfile(ctx context.Context, userID string) (*candidate.Profile, error)
	UpsertCandidateProfile(ctx context.Context, profile *candidate.Profile) error
	SetOpenFlag(ctx context.Context, userID string, flag bool) error
	PostJob(ctx context.Context, details *job.Posting) error
	UpdateJobPosting(ctx context.Context, post *job.Posting) error
	UpdateJobDetails(ctx context.Context, post *job.Posting) error
	DeleteExpiredPost(ctx context.Context, postID, companyID string) error
	GetJobPosting(ctx context.Context, jobID string, revertApplicant bool) (*job.Posting, error)
	GetDraft(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error)
	GetPost(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error)
	GetIDsOfApplicants(ctx context.Context, jobID string, reverse bool) ([]string, error)
	GetListOfJobsWithSeenStat(ctx context.Context, companyID string, first, after int32) ([]*job.ViewJobWithSeenStat, error)
	IsJobApplied(ctx context.Context, userID, jobID string) (bool, error)
	IgnoreInvitation(ctx context.Context, userID, jobID string) error
	IsJobSaved(ctx context.Context, userID, jobID string) (bool, error)
	GetCompanyIDByJobID(ctx context.Context, jobID string) (string, error)
	ApplyJob(ctx context.Context, userID, jobID string, application *job.Application, fileIDs []string) error
	AddJobView(ctx context.Context, userID, jobID string) error
	GetRecommendedJobs(ctx context.Context, userID, countryID string, first, after int32) ([]*job.ViewForUser, error)
	SaveJob(ctx context.Context, userID, jobID string) error
	UnsaveJob(ctx context.Context, userID, jobID string) error
	SkipJob(ctx context.Context, userID, jobID string) error
	UnskipJob(ctx context.Context, userID, jobID string) error
	GetSavedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error)
	GetSkippedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error)
	GetAppliedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error)
	SetJobApplicationSeen(ctx context.Context, jobID, applicantID string, seen bool) error
	SetJobApplicationCategory(ctx context.Context, jobID, applicantID string, category job.ApplicantCategory) error
	GetCandidates(ctx context.Context, companyID string, country string, first, after int32) ([]*candidate.ViewForCompany, error)
	SaveCandidate(ctx context.Context, companyID, candidateID string) error
	UnsaveCandidate(ctx context.Context, companyID, candidateID string) error
	SkipCandidate(ctx context.Context, companyID, candidateID string) error
	UnskipCandidate(ctx context.Context, companyID, candidateID string) error
	GetAmountsOfManageCandidates(ctx context.Context, companyID string) (int32, int32, error)
	GetAmountOfApplicantsPerCategory(ctx context.Context, companyID string) (total, unseen, favorite, inReview, disqualified int32, err error)
	GetSavedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error)
	GetSkippedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error)
	GetPostedJobs(ctx context.Context, companyID string) ([]*job.ViewForCompany, error)
	InviteUserToApply(ctx context.Context, jobID string, invitation job.Invitation) error
	GetInvitedJobs(ctx context.Context, userID string, first, after int32) ([]*job.ViewForUser, error)
	IsCandidateInvited(ctx context.Context, jobID, candidateID string) (bool, error)
	GetJobForCompany(ctx context.Context, companyID string, jobID string) (*job.ViewForCompany, error)
	ReportJob(ctx context.Context, report *job.Report) error
	ReportCandidate(ctx context.Context, report *candidate.Report) error
	SaveJobSearchFilter(ctx context.Context, filter *job.NamedSearchFilter) error
	GetSavedJobSearchFilters(ctx context.Context, userID string) ([]*job.NamedSearchFilter, error)
	SaveCandidateSearchFilter(ctx context.Context, companyID string, filter *candidate.NamedSearchFilter) error
	GetSavedCandidateSearchFilters(ctx context.Context, companyID string) ([]*candidate.NamedSearchFilter, error)
	GetPlanPrices(ctx context.Context, countriesID []string, currency string) ([]job.PlanPrices, error)
	GetPricingFor(ctx context.Context, meta *job.Meta) (*job.PricingResult, error)
	GetAmountOfActiveJobsOfCompany(ctx context.Context, companyID string) (int32, error)
	GetCareerInterestsByIds(ctx context.Context, ids []string, first, after int32) (map[string]*candidate.CareerInterests, error)
	GetSortedListOfIDsOfApplicant(ctx context.Context, userIDs []string, sort candidate.ApplicantSort, first, after int32) ([]string, error)
	UploadFileForApplication(ctx context.Context, userID string, jobID string, file *job.File) error
	UploadFileForJob(ctx context.Context, companyID string, jobID string, file *job.File) error

	IsCandidateSaved(ctx context.Context, companyID, candidateID string) (bool, error)

	AddCVInCareerCenter(ctx context.Context, userID, companyID string, options careercenter.CVOptions) error
	GetSavedCVs(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error)
	RemoveCVs(ctx context.Context, companyID string, ids []string) error
	MakeFavoriteCVs(ctx context.Context, companyID string, ids []string, isFavourite bool) error
}

// MQ ...
type MQ interface {
	SendNewInvitation(targetID string, not *notmes.NewInvitation) error
	SendNewJobApplicant(targetID string, not *notmes.NewJobApplicant) error
}
