package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	careercenter "gitlab.lan/Rightnao-site/microservices/jobs/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job"
)

// Service define functions inside Service
type Service interface {
	GetCandidateProfile(ctx context.Context) (*candidate.Profile, error)
	SetCareerInterests(ctx context.Context, data *candidate.CareerInterests) error
	SetOpenFlag(ctx context.Context, flag bool) error
	PostJob(ctx context.Context, companyID string, post *job.Posting) (string, error)
	ChangePost(ctx context.Context, draftID, companyID string, post *job.Posting) (string, error)
	DeleteExpiredPost(ctx context.Context, postID, companyID string) error
	SaveDraft(ctx context.Context, companyID string, post *job.Posting) (string, error)
	ChangeDraft(ctx context.Context, draftID, companyID string, post *job.Posting) (string, error)
	ActivateJob(ctx context.Context, companyID, jobID string) error
	PauseJob(ctx context.Context, companyID, jobID string) error
	ApplyJob(ctx context.Context, jobID string, application *job.Application, fileIDs []string) error
	IgnoreInvitation(ctx context.Context, jobID string) error
	AddJobView(ctx context.Context, jobID string) error
	GetRecommendedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error)
	GetJob(ctx context.Context, jobID string) (*job.ViewForUser, error)
	GetDraft(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error)
	GetPost(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error)
	GetJobApplicants(ctx context.Context, companyID, jobID string, sort candidate.ApplicantSort, first, after int32) ([]*job.Applicant, error)
	SaveJob(ctx context.Context, jobID string) error
	UnsaveJob(ctx context.Context, jobID string) error
	SkipJob(ctx context.Context, jobID string) error
	UnskipJob(ctx context.Context, jobID string) error
	GetSavedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error)
	GetSkippedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error)
	GetListOfJobsWithSeenStat(ctx context.Context, companyID string, first, after int32) ([]*job.ViewJobWithSeenStat, error)
	GetAmountOfApplicantsPerCategory(ctx context.Context, companyID string) (total, unseen, favorite, inReview, disqualified int32, err error)
	GetAppliedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error)
	SetJobApplicationSeen(ctx context.Context, companyID, jobID, applicantID string, seen bool) error
	SetJobApplicationCategory(ctx context.Context, companyID, jobID, applicantID string, category job.ApplicantCategory) error
	GetCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error)
	SaveCandidate(ctx context.Context, companyID, candidateID string) error
	UnsaveCandidate(ctx context.Context, companyID, candidateID string) error
	SkipCandidate(ctx context.Context, companyID, candidateID string) error
	UnskipCandidate(ctx context.Context, companyID, candidateID string) error
	GetAmountsOfManageCandidates(ctx context.Context, companyID string) (int32, int32, int32, error)
	GetSavedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error)
	GetSkippedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error)
	GetPostedJobs(ctx context.Context, companyID string) ([]*job.ViewForCompany, error)
	GetJobForCompany(ctx context.Context, companyID string, jobID string) (*job.ViewForCompany, error)
	InviteUserToApply(ctx context.Context, companyID string, jobID string, userID string, text string) error
	GetInvitedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error)
	ReportJob(ctx context.Context, jobID string, reportType job.ReportType, text string) error
	ReportCandidate(ctx context.Context, companyID, candidateID string, text string) error
	SaveJobSearchFilter(ctx context.Context, filter *job.NamedSearchFilter) error
	GetSavedJobSearchFilters(ctx context.Context) ([]*job.NamedSearchFilter, error)
	SaveCandidateSearchFilter(ctx context.Context, companyID string, filter *candidate.NamedSearchFilter) error
	GetSavedCandidateSearchFilters(ctx context.Context, companyID string) ([]*candidate.NamedSearchFilter, error)
	GetPlanPrices(ctx context.Context, companyID string, countriesID []string, currency string) ([]job.PlanPrices, error)
	GetPricingFor(ctx context.Context, companyID string, meta *job.Meta) (*job.PricingResult, error)
	GetAmountOfActiveJobsOfCompany(ctx context.Context, companyID string) (int32, error)
	GetCareerInterestsByIds(ctx context.Context, companyID string, ids []string, first, after int32) (map[string]*candidate.CareerInterests, error)

	UploadFileForApplication(ctx context.Context, userID string, jobID string, file *job.File) (string, error)
	UploadFileForJob(ctx context.Context, companyID string, jobID string, file *job.File) (string, error)

	AddCVInCareerCenter(ctx context.Context, companyID string, options careercenter.CVOptions) error
	GetSavedCVs(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error)
	RemoveCVs(ctx context.Context, companyID string, ids []string) error
	MakeFavoriteCVs(ctx context.Context, companyID string, ids []string, isFavourite bool) error
}
