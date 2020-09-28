package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	careercenter "gitlab.lan/Rightnao-site/microservices/jobs/internal/career-center"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/company"
	companyadmin "gitlab.lan/Rightnao-site/microservices/jobs/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/job"
	notmes "gitlab.lan/Rightnao-site/microservices/jobs/internal/notification_messages"
)

// GetCandidateProfile ...
func (s *Service) GetCandidateProfile(ctx context.Context) (*candidate.Profile, error) {
	span := s.tracer.MakeSpan(ctx, "GetCandidateProfile")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	profile, err := s.jobs.GetCandidateProfile(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	if ci := profile.CareerInterests; ci != nil {
		for i, city := range ci.Locations {
			cityName, subdivision, country, err := s.infoRPC.GetCityInformationByID(ctx, city.CityID, nil)
			if err != nil {
				s.tracer.LogError(span, err)
			}

			ci.Locations[i].City = cityName
			ci.Locations[i].Subdivision = subdivision
			ci.Locations[i].Country = country
		}
	}

	return profile, nil
}

// SetCareerInterests ...
func (s *Service) SetCareerInterests(ctx context.Context, data *candidate.CareerInterests) error {
	span := s.tracer.MakeSpan(ctx, "SetCareerInterests")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	data.SetUserID(userID)
	// data.NormalizedSalary = float32(data.SalaryMin) / float32(data.SalaryInterval.GetHours()) // TODO: add curency in calculation

	for i, city := range data.Locations {
		cityName, subdivision, country, err := s.infoRPC.GetCityInformationByID(ctx, city.CityID, nil)
		if err != nil {
			s.tracer.LogError(span, err)
		}

		data.Locations[i].City = cityName
		data.Locations[i].Subdivision = subdivision
		data.Locations[i].Country = country
	}

	profile := candidate.Profile{
		IsOpen:          true,
		CareerInterests: data,
	}
	profile.SetUserID(userID)

	err = s.jobs.UpsertCandidateProfile(ctx, &profile)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SetOpenFlag ...
func (s *Service) SetOpenFlag(ctx context.Context, flag bool) error {
	span := s.tracer.MakeSpan(ctx, "SetOpenFlag")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.jobs.SetOpenFlag(ctx, userID, flag)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// PostJob set current day as day of publish, set status `Active` and saves in DB
func (s *Service) PostJob(ctx context.Context, companyID string, post *job.Posting) (string, error) {
	span := s.tracer.MakeSpan(ctx, "PostJob")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	post.SetCompanyID(companyID)
	post.SetUserID(userID)
	id := post.GenerateID()

	// if !post.JobMetadata.Anonymous {
	post.CompanyDetails = &company.Details{
		// TODO: company avatar, URL, Industry, subindustry
	}
	post.CompanyDetails.SetCompanyID(companyID)
	// }

	now := time.Now()
	post.CreatedAt = now
	// post.ActivationDate = time.Date(
	// 	int(post.JobDetails.PublishYear),
	// 	time.Month(int(post.JobDetails.PublishMonth)),
	// 	int(post.JobDetails.PublishDay),
	// 	0, 0, 0, 0, time.Now().Location(),
	// )
	post.ActivationDate = now //time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// if post.ActivationDate.Before(post.CreatedAt) {
	// 	return "", errors.New("activation_date_can_not_be_in_past")
	// }

	post.Status = job.StatusActive //TODO: why it was Draft?
	// post.JobPriority = post.JobMetadata.JobPlan.GetPriority() // there is no more priority after removing plans?

	// if post.JobDetails.SalaryMin > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMin = float32(post.JobDetails.SalaryMin) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }

	// if post.JobDetails.SalaryMax > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMax = float32(post.JobDetails.SalaryMax) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }

	// calc expire date
	post.ExpirationDate = now.AddDate(0, 0, int(post.JobMetadata.AmountOfDays))

	err = s.jobs.PostJob(ctx, post)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ChangePost set current day as day of publish, set status `Active` and saves in DB
func (s *Service) ChangePost(ctx context.Context, draftID, companyID string, post *job.Posting) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangePost")
	defer span.Finish()

	// //  get userID
	// userID, err := s.authRPC.GetUserID(ctx)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return "", err
	// }

	err := post.SetID(draftID)
	if err != nil {
		return "", err
	}

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	// post.SetCompanyID(companyID)
	// post.SetUserID(userID)
	// // id := post.GenerateID()
	//
	// if !post.JobMetadata.Anonymous {
	// 	post.CompanyDetails = &company.Details{
	// 		// TODO: company avatar, URL, Industry, subindustry
	// 	}
	// 	post.CompanyDetails.SetCompanyID(companyID)
	// }
	//
	// now := time.Now()
	// post.CreatedAt = now
	// // post.ActivationDate = time.Date(
	// // 	int(post.JobDetails.PublishYear),
	// // 	time.Month(int(post.JobDetails.PublishMonth)),
	// // 	int(post.JobDetails.PublishDay),
	// // 	0, 0, 0, 0, time.Now().Location(),
	// // )
	// post.ActivationDate = now //time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	//
	// // if post.ActivationDate.Before(post.CreatedAt) {
	// // 	return "", errors.New("activation_date_can_not_be_in_past")
	// // }
	//
	// post.Status = job.StatusActive //TODO: why it was Draft?
	// post.JobPriority = post.JobMetadata.JobPlan.GetPriority()
	//
	// if post.JobDetails.SalaryMin > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMin = float32(post.JobDetails.SalaryMin) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }
	//
	// if post.JobDetails.SalaryMax > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMax = float32(post.JobDetails.SalaryMax) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }
	//
	// // calc expire date
	// post.ExpirationDate = now.AddDate(0, 0, post.JobMetadata.JobPlan.GetDays())

	err = s.jobs.UpdateJobDetails(ctx, post)
	if err != nil {
		return "", err
	}

	return post.GetID(), nil
}

// DeleteExpiredPost ...
func (s Service) DeleteExpiredPost(ctx context.Context, postID, companyID string) error {
	span := s.tracer.MakeSpan(ctx, "DeleteExpiredPost")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.jobs.DeleteExpiredPost(ctx, postID, companyID)
	if err != nil {
		return err
	}

	return nil
}

// SaveDraft set status `Draft` and saves in DB
func (s *Service) SaveDraft(ctx context.Context, companyID string, post *job.Posting) (string, error) {
	span := s.tracer.MakeSpan(ctx, "SaveDraft")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	post.SetCompanyID(companyID)
	post.SetUserID(userID)
	id := post.GenerateID()

	if !post.JobMetadata.Anonymous {
		post.CompanyDetails = &company.Details{
			// TODO: company avatar, URL, Industry, subindustry
		}
		post.CompanyDetails.SetCompanyID(companyID)
	}

	post.CreatedAt = time.Now()

	post.Status = job.StatusDraft
	// post.JobPriority = post.JobMetadata.JobPlan.GetPriority()

	// if post.JobDetails.SalaryMin > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMin = float32(post.JobDetails.SalaryMin) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }

	// if post.JobDetails.SalaryMax > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMax = float32(post.JobDetails.SalaryMax) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }

	err = s.jobs.PostJob(ctx, post)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ChangeDraft changes the fields in draft in DB
func (s *Service) ChangeDraft(ctx context.Context, draftID, companyID string, post *job.Posting) (string, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeDraft")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = post.SetCompanyID(companyID)
	if err != nil {
		return "", err
	}

	err = post.SetID(draftID)
	if err != nil {
		return "", err
	}

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		post.GetCompanyID(),
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return "", errors.New("not_allowed")
	}

	post.SetUserID(userID)
	// id := post.GenerateID()

	if !post.JobMetadata.Anonymous {
		post.CompanyDetails = &company.Details{
			// TODO: company avatar, URL, Industry, subindustry
		}
		post.CompanyDetails.SetCompanyID(post.GetCompanyID())
	}

	post.CreatedAt = time.Now()

	post.Status = job.StatusDraft
	// post.JobPriority = post.JobMetadata.JobPlan.GetPriority()

	// if post.JobDetails.SalaryMin > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMin = float32(post.JobDetails.SalaryMin) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }

	// if post.JobDetails.SalaryMax > 0 && post.JobDetails.SalaryInterval != "" {
	// 	post.NormalizedSalaryMax = float32(post.JobDetails.SalaryMax) / float32(post.JobDetails.SalaryInterval.GetHours()) // TODO also convert currency
	// }

	err = s.jobs.UpdateJobPosting(ctx, post)
	if err != nil {
		return "", err
	}

	return post.GetID(), nil
}

// ActivateJob checks if it is not `Draft` or `Paused`. Set status `Active`.
// If it was draft (never paused) just set activation date.
// If it was paused: recalculate `PausedDays`  (total paused hours divided by 24)
// and `ExpirationDate` (days(according to job plan) * (renewal + 1) + paused days)
func (s *Service) ActivateJob(ctx context.Context, companyID, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "ActivateJob")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	post, err := s.jobs.GetJobPosting(ctx, jobID, false)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if post.GetCompanyID() != companyID {
		return errors.New("not_allowed")
	}

	if post.Status != job.StatusDraft && post.Status != job.StatusPaused {
		return errors.New("not_allowed")
	}

	post.Status = job.StatusActive

	// if it's Draft
	if post.LastPauseDate.IsZero() {
		post.ActivationDate = time.Now()
	} else { // if it's Paused
		post.PausedDays += int(time.Now().Sub(post.LastPauseDate).Hours() / 24) // total paused hours divided by 24 (returns 0 if it was unpaused at the same day )
	}
	post.ExpirationDate = post.ActivationDate.AddDate(
		0,
		0,
		// BUG: calculation is sucks
		(int(post.JobMetadata.AmountOfDays)*(int(post.JobMetadata.Renewal)+1))+post.PausedDays, // days(according to job plan) * (renewal + 1) + paused days
	)

	err = s.jobs.UpdateJobPosting(ctx, post)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// PauseJob ...
func (s *Service) PauseJob(ctx context.Context, companyID, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "PauseJob")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	post, err := s.jobs.GetJobPosting(ctx, jobID, false)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if post.GetCompanyID() != companyID {
		return errors.New("not_allowed")
	}

	if post.Status != job.StatusActive {
		return errors.New("not_allowed")
	}

	post.Status = job.StatusPaused
	post.LastPauseDate = time.Now()

	err = s.jobs.UpdateJobPosting(ctx, post)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetListOfJobsWithSeenStat ...
func (s *Service) GetListOfJobsWithSeenStat(ctx context.Context, companyID string, first, after int32) ([]*job.ViewJobWithSeenStat, error) {
	span := s.tracer.MakeSpan(ctx, "GetListOfJobsWithSeenStat")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	if first <= 0 {
		first = 10
	}

	jobs, err := s.jobs.GetListOfJobsWithSeenStat(ctx, companyID, first, after)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// ApplyJob ...
func (s *Service) ApplyJob(ctx context.Context, jobID string, application *job.Application, fileIDs []string) error {
	span := s.tracer.MakeSpan(ctx, "ApplyJob")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	isApplied, err := s.jobs.IsJobApplied(ctx, userID, jobID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isApplied {
		return errors.New("already_applied")
	}

	application.SetUserID(userID)
	application.CreatedAt = time.Now()

	companyID, err := s.jobs.GetCompanyIDByJobID(ctx, jobID)
	if err != nil {
		return err
	}

	err = s.jobs.ApplyJob(ctx, userID, jobID, application, fileIDs)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// send notification
	err = s.mq.SendNewJobApplicant(companyID, &notmes.NewJobApplicant{
		CandidateID: userID,
		JobID:       jobID,
		CompanyID:   companyID,
	})
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// IgnoreInvitation ...
func (s *Service) IgnoreInvitation(ctx context.Context, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "IgnoreInvitation")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.jobs.IgnoreInvitation(ctx, userID, jobID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// AddJobView ...
func (s *Service) AddJobView(ctx context.Context, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "AddJobView")
	defer span.Finish()

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.jobs.AddJobView(ctx, userID, jobID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetRecommendedJobs ...
func (s *Service) GetRecommendedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error) {
	span := s.tracer.MakeSpan(ctx, "GetRecommendedJobs")
	defer span.Finish()

	country, err := s.infoRPC.GetUserCountry(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	if first <= 0 {
		first = 10
	}

	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		// s.tracer.LogError(span, err)
		// return nil, err
	}

	// log.Println("userID", userID)

	recJobs, err := s.jobs.GetRecommendedJobs(ctx, userID, country, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i, j := range recJobs {
		if j != nil {
			if recJobs[i].JobDetails.City != "" {
				lang := "en"
				cityID, err := strconv.Atoi(recJobs[i].JobDetails.City)
				if err != nil {
					s.tracer.LogError(span, err)
					return nil, err
				}

				cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
				if err != nil {
					s.tracer.LogError(span, err)
					return nil, err
				}

				recJobs[i].JobDetails.Location = job.Location{
					CityID:      recJobs[i].JobDetails.City,
					CityName:    cityName,
					Country:     countryID,
					Subdivision: subdivision,
				}
			}

			// remove all except job plan
			j.Metadata.AdvertisementCountries = []string{}
			j.Metadata.Anonymous = false
			j.Metadata.Currency = ""
			j.Metadata.NumOfLanguages = 0
			j.Metadata.Renewal = 0
			// check if it's applied
			j.IsApplied, err = s.jobs.IsJobApplied(ctx, userID, j.GetID())
			if err != nil {
				s.tracer.LogError(span, err)
			}

		}
	}

	return recJobs, nil
}

// GetJob ...
func (s *Service) GetJob(ctx context.Context, jobID string) (*job.ViewForUser, error) {
	span := s.tracer.MakeSpan(ctx, "GetJob")
	defer span.Finish()

	post, err := s.jobs.GetJobPosting(ctx, jobID, false)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	view := job.ViewForUser{
		JobDetails: post.JobDetails,
		Metadata:   post.JobMetadata,
	}
	if post.CompanyDetails != nil {
		view.CompanyDetails = *post.CompanyDetails
	}
	view.SetID(post.GetID())

	userID, err := s.authRPC.GetUserID(ctx)
	if err == nil && userID != "" {
		// isApplied
		view.IsApplied, err = s.jobs.IsJobApplied(ctx, userID, view.GetID())
		if err != nil {
			s.tracer.LogError(span, err)
		}
		// isSaved
		view.IsSaved, err = s.jobs.IsJobSaved(ctx, userID, view.GetID())
		if err != nil {
			s.tracer.LogError(span, err)
		}

		// log.Println(view.JobDetails.DeadlineYear, post.CreatedAt)

		if view.JobDetails.City != "" {
			lang := "en"
			cityID, err := strconv.Atoi(view.JobDetails.City)
			if err != nil {
				s.tracer.LogError(span, err)
				return nil, err
			}

			cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
			if err != nil {
				s.tracer.LogError(span, err)
				return nil, err
			}

			view.JobDetails.Location = job.Location{
				CityID:      view.JobDetails.City,
				CityName:    cityName,
				Country:     countryID,
				Subdivision: subdivision,
			}
		}

		if post.ExpirationDate.Year() != 1 {
			view.JobDetails.DeadlineDay = int32(post.ExpirationDate.Day())
			view.JobDetails.DeadlineYear = int32(post.ExpirationDate.Year())
			view.JobDetails.DeadlineMonth = int32(post.ExpirationDate.Month())
		}

		if post.CreatedAt.Year() != 1 {
			view.JobDetails.PublishDay = int32(post.CreatedAt.Day())
			view.JobDetails.PublishYear = int32(post.CreatedAt.Year())
			view.JobDetails.PublishMonth = int32(post.CreatedAt.Month())
		}

	} else {
		s.tracer.LogError(span, err)
	}

	return &view, nil
}

// GetDraft ...
func (s *Service) GetDraft(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetDraft")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	draft, err := s.jobs.GetDraft(ctx, draftID, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// view := job.ViewForCompany{
	// 	JobDetails: draft.JobDetails,
	// 	Metadata:   draft.Metadata,
	// }

	// view.SetID(draft.GetID())

	lang := "en"
	n, err := strconv.Atoi(draft.JobDetails.City)
	if err != nil {
		log.Println("wrong city id", err)
	} else {
		name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), &lang)
		if err != nil {
			log.Println("wrong city info", err)
		} else {
			draft.JobDetails.Location = job.Location{
				CityID:      draft.JobDetails.City,
				CityName:    name,
				Country:     country,
				Subdivision: sub,
			}
		}
	}

	return draft, nil
}

// GetPost ...
func (s *Service) GetPost(ctx context.Context, draftID, companyID string) (*job.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetPost")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	draft, err := s.jobs.GetPost(ctx, draftID, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// view := job.ViewForCompany{
	// 	JobDetails: draft.JobDetails,
	// }

	// view.SetID(draft.GetID())

	lang := "en"
	n, err := strconv.Atoi(draft.JobDetails.City)
	if err != nil {
		log.Println("wrong city id", err)
	} else {
		name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), &lang)
		if err != nil {
			log.Println("wrong city info", err)
		} else {
			draft.JobDetails.Location = job.Location{
				CityID:      draft.JobDetails.City,
				CityName:    name,
				Country:     country,
				Subdivision: sub,
			}
		}
	}

	return draft, nil
}

// GetJobApplicants ...
// TODO: add sort ?
func (s *Service) GetJobApplicants(ctx context.Context, companyID, jobID string, sort candidate.ApplicantSort, first, after int32) ([]*job.Applicant, error) {
	span := s.tracer.MakeSpan(ctx, "GetJobApplicants")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	// // if sort by date then just revert elements in array in `jobs` collection
	// var revert bool
	// if sort == candidate.AppicantPostedDate {
	// 	revert = true
	// }
	//

	postedJob, err := s.jobs.GetJobPosting(ctx, jobID, false /*revert*/)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// if postedJob.GetCompanyID() != companyID {
	// 	return nil, errors.New("not_allowed")
	// }

	// get ids of applicant
	userIDs := make([]string, 0, len(postedJob.Applications))
	for _, app := range postedJob.Applications {
		userIDs = append(userIDs, app.GetUserID())
	}

	// if sort != candidate.AppicantPostedDate {
	// 	ids, err := s.jobs.GetSortedListOfIDsOfApplicant(ctx, userIDs, sort, first, after)
	// 	if err != nil {
	// 		s.tracer.LogError(span, err)
	// 		return nil, err
	// 	}
	//
	// 	log.Println("ids GetSortedList:", ids)
	// } else {
	// 	after = 0
	// }

	careerInterests, err := s.GetCareerInterestsByIds(ctx, companyID, userIDs, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	candidates := make([]*job.Applicant, 0, len(postedJob.Applications))
	for _, app := range postedJob.Applications {

		// if ci, isExists := careerInterests[app.GetUserID()]; isExists {
		ci := careerInterests[app.GetUserID()]
		candidates = append(candidates, &job.Applicant{
			Application:     *app,
			CareerInterests: ci,
		})
		// }
	}
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range candidates {
		if candidates[i].CareerInterests != nil {
			candidates[i].CareerInterests.IsInvited, err = s.jobs.IsCandidateInvited(ctx, jobID, candidates[i].CareerInterests.GetUserID())
			if err != nil {
				s.tracer.LogError(span, err)
				return nil, err
			}
		}
	}

	return candidates, nil
}

// SaveJob ...
func (s *Service) SaveJob(ctx context.Context, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "SaveJob")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.jobs.SaveJob(ctx, userID, jobID)
	if err != nil {
		return err
	}

	return nil
}

// UnsaveJob ...
func (s *Service) UnsaveJob(ctx context.Context, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "UnsaveJob")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.jobs.UnsaveJob(ctx, userID, jobID)
	if err != nil {
		return err
	}

	return nil
}

// SkipJob ...
func (s *Service) SkipJob(ctx context.Context, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "SkipJob")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.jobs.SkipJob(ctx, userID, jobID)
	if err != nil {
		return err
	}

	return nil
}

// UnskipJob ...
func (s *Service) UnskipJob(ctx context.Context, jobID string) error {
	span := s.tracer.MakeSpan(ctx, "UnskipJob")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.jobs.UnskipJob(ctx, userID, jobID)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedJobs ...
func (s *Service) GetSavedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error) {
	span := s.tracer.MakeSpan(ctx, "GetSavedJobs")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	jobs, err := s.jobs.GetSavedJobs(ctx, userID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range jobs {
		lang := "en"
		cityID, err := strconv.Atoi(jobs[i].JobDetails.City)
		if err != nil {
			// s.tracer.LogError(span, err)
			// return nil, err
		}

		cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		jobs[i].JobDetails.Location = job.Location{
			CityID:      jobs[i].JobDetails.City,
			CityName:    cityName,
			Country:     countryID,
			Subdivision: subdivision,
		}

		if userID != "" {
			// isApplied
			jobs[i].IsApplied, err = s.jobs.IsJobApplied(ctx, userID, jobs[i].GetID())
			if err != nil {
				s.tracer.LogError(span, err)
			}
			// isSaved
			jobs[i].IsSaved, err = s.jobs.IsJobSaved(ctx, userID, jobs[i].GetID())
			if err != nil {
				s.tracer.LogError(span, err)
			}
		} else {
			s.tracer.LogError(span, err)
		}
	}

	return jobs, nil
}

// GetSkippedJobs ...
func (s *Service) GetSkippedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error) {
	span := s.tracer.MakeSpan(ctx, "GetSkippedJobs")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	jobs, err := s.jobs.GetSkippedJobs(ctx, userID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range jobs {
		lang := "en"
		cityID, err := strconv.Atoi(jobs[i].JobDetails.City)
		if err != nil {
			// s.tracer.LogError(span, err)
			// return nil, err
		}

		cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		jobs[i].JobDetails.Location = job.Location{
			CityID:      jobs[i].JobDetails.City,
			CityName:    cityName,
			Country:     countryID,
			Subdivision: subdivision,
		}

		if userID != "" {
			// isApplied
			jobs[i].IsApplied, err = s.jobs.IsJobApplied(ctx, userID, jobs[i].GetID())
			if err != nil {
				s.tracer.LogError(span, err)
			}
			// isSaved
			jobs[i].IsSaved, err = s.jobs.IsJobSaved(ctx, userID, jobs[i].GetID())
			if err != nil {
				s.tracer.LogError(span, err)
			}
		} else {
			s.tracer.LogError(span, err)
		}
	}

	return jobs, nil
}

// GetAppliedJobs ...
func (s *Service) GetAppliedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error) {
	span := s.tracer.MakeSpan(ctx, "GetAppliedJobs")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	jobs, err := s.jobs.GetAppliedJobs(ctx, userID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range jobs {
		lang := "en"
		cityID, err := strconv.Atoi(jobs[i].JobDetails.City)
		if err != nil {
			// s.tracer.LogError(span, err)
			// return nil, err
		}

		cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		jobs[i].JobDetails.Location = job.Location{
			CityID:      jobs[i].JobDetails.City,
			CityName:    cityName,
			Country:     countryID,
			Subdivision: subdivision,
		}
	}

	return jobs, nil
}

// SetJobApplicationSeen ...
func (s *Service) SetJobApplicationSeen(ctx context.Context, companyID, jobID, applicantID string, seen bool) error {
	span := s.tracer.MakeSpan(ctx, "SetJobApplicationSeen")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// post, err := s.jobs.GetJobPosting(ctx, jobID)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }
	//
	// if post.GetCompanyID() != companyID {
	// 	return errors.New("not_allowed")
	// }

	err := s.jobs.SetJobApplicationSeen(ctx, jobID, applicantID, seen)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SetJobApplicationCategory ...
func (s *Service) SetJobApplicationCategory(ctx context.Context, companyID, jobID, applicantID string, category job.ApplicantCategory) error {
	span := s.tracer.MakeSpan(ctx, "SetJobApplicationCategory")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// post, err := s.jobs.GetJobPosting(ctx, jobID)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }
	//
	// if post.GetCompanyID() != companyID {
	// 	return errors.New("not_allowed")
	// }

	err := s.jobs.SetJobApplicationCategory(ctx, jobID, applicantID, category)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetCandidates ...
func (s *Service) GetCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetCandidates")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	country, err := s.infoRPC.GetUserCountry(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	if first <= 0 {
		first = 10
	}

	profile, err := s.jobs.GetCandidates(ctx, companyID, country, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return profile, nil
}

// SaveCandidate ...
func (s *Service) SaveCandidate(ctx context.Context, companyID, candidateID string) error {
	span := s.tracer.MakeSpan(ctx, "SaveCandidate")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.jobs.SaveCandidate(ctx, companyID, candidateID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// UnsaveCandidate ...
func (s *Service) UnsaveCandidate(ctx context.Context, companyID, candidateID string) error {
	span := s.tracer.MakeSpan(ctx, "UnsaveCandidate")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.jobs.UnsaveCandidate(ctx, companyID, candidateID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// SkipCandidate ...
func (s *Service) SkipCandidate(ctx context.Context, companyID, candidateID string) error {
	span := s.tracer.MakeSpan(ctx, "SkipCandidate")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.jobs.SkipCandidate(ctx, companyID, candidateID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// UnskipCandidate ...
func (s *Service) UnskipCandidate(ctx context.Context, companyID, candidateID string) error {
	span := s.tracer.MakeSpan(ctx, "UnskipCandidate")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.jobs.UnskipCandidate(ctx, companyID, candidateID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetAmountsOfManageCandidates ...
func (s *Service) GetAmountsOfManageCandidates(ctx context.Context, companyID string) (int32, int32, int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetAmountsOfManageCandidates")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return 0, 0, 0, errors.New("not_allowed")
	}

	saved, skipped, err := s.jobs.GetAmountsOfManageCandidates(ctx, companyID)
	if err != nil {
		return 0, 0, 0, err
	}

	return saved, skipped, 0, nil
}

// GetAmountOfApplicantsPerCategory ...
func (s *Service) GetAmountOfApplicantsPerCategory(ctx context.Context, companyID string) (total, unseen, favorite, inReview, disqualified int32, err error) {
	span := s.tracer.MakeSpan(ctx, "GetAmountOfApplicantsPerCategory")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		err = errors.New("not_allowed")
		return
	}

	total, unseen, favorite, inReview, disqualified, err = s.jobs.GetAmountOfApplicantsPerCategory(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return
	}

	return
}

// GetSavedCandidates ...
func (s *Service) GetSavedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetSavedCandidates")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	if first <= 0 {
		first = 10
	}

	candidates, err := s.jobs.GetSavedCandidates(ctx, companyID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return candidates, nil
}

// GetSkippedCandidates ...
func (s *Service) GetSkippedCandidates(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetSkippedCandidates")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	if first <= 0 {
		first = 10
	}

	candidates, err := s.jobs.GetSkippedCandidates(ctx, companyID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return candidates, nil
}

// GetPostedJobs ...
func (s *Service) GetPostedJobs(ctx context.Context, companyID string) ([]*job.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetPostedJobs")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	jobs, err := s.jobs.GetPostedJobs(ctx, companyID)
	if err != nil {
		return nil, err
	}

	lang := "en"
	for i := range jobs {
		n, err := strconv.Atoi(jobs[i].JobDetails.City) // city id
		if err != nil {
			log.Println("wrong city id", err)
		} else {
			name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), &lang)
			if err != nil {
				log.Println("wrong city info", err)
			} else {
				jobs[i].JobDetails.Location = job.Location{
					CityID:      jobs[i].JobDetails.City,
					CityName:    name,
					Country:     country,
					Subdivision: sub,
				}
			}
		}

	}

	for i := range jobs {
		if jobs[i].ExpiredAt.Year() != 1 {
			jobs[i].JobDetails.PublishDay = int32(jobs[i].PublishedAt.Day())
			jobs[i].JobDetails.PublishYear = int32(jobs[i].PublishedAt.Year())
			jobs[i].JobDetails.PublishMonth = int32(jobs[i].PublishedAt.Month())
		}

		if jobs[i].ExpiredAt.Year() != 1 {
			jobs[i].JobDetails.DeadlineDay = int32(jobs[i].ExpiredAt.Day())
			jobs[i].JobDetails.DeadlineYear = int32(jobs[i].ExpiredAt.Year())
			jobs[i].JobDetails.DeadlineMonth = int32(jobs[i].ExpiredAt.Month())
		}
	}

	return jobs, nil
}

// GetJobForCompany ...
func (s *Service) GetJobForCompany(ctx context.Context, companyID string, jobID string) (*job.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetJobForCompany")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	jobs, err := s.jobs.GetJobForCompany(ctx, companyID, jobID)
	if err != nil {
		return nil, err
	}

	if jobs.JobDetails.City != "" {
		n, err := strconv.Atoi(jobs.JobDetails.City)
		if err != nil {
			log.Println("wrong city id")
		} else {
			lang := "en"
			name, sub, country, err := s.infoRPC.GetCityInformationByID(ctx, int32(n), &lang)
			if err != nil {
				return nil, err
			}

			jobs.JobDetails.Location = job.Location{
				CityID:      jobs.JobDetails.City,
				CityName:    name,
				Country:     country,
				Subdivision: sub,
			}
			log.Println("location:", jobs.JobDetails.Location)
		}
	}

	return jobs, nil
}

// InviteUserToApply ...
func (s *Service) InviteUserToApply(ctx context.Context, companyID string, jobID string, userID string, text string) error {
	span := s.tracer.MakeSpan(ctx, "InviteUserToApply")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	// check if it wasn't invited before
	isInvited, err := s.jobs.IsCandidateInvited(ctx, jobID, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isInvited {
		return errors.New("already_invited")
	}

	invitation := job.Invitation{
		Text: text,
	}
	invitation.SetUserID(userID)

	err = s.jobs.InviteUserToApply(ctx, jobID, invitation)
	if err != nil {
		return err
	}

	// send notification
	n := &notmes.NewInvitation{
		CompanyID: companyID,
		JobID:     jobID,
	}
	n.GenerateID()
	err = s.mq.SendNewInvitation(userID, n)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	return nil
}

// GetInvitedJobs ...
func (s *Service) GetInvitedJobs(ctx context.Context, first, after int32) ([]*job.ViewForUser, error) {
	span := s.tracer.MakeSpan(ctx, "GetInvitedJobs")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	jobs, err := s.jobs.GetInvitedJobs(ctx, userID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	for i := range jobs {
		lang := "en"
		cityID, err := strconv.Atoi(jobs[i].JobDetails.City)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		cityName, subdivision, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(cityID), &lang)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		jobs[i].JobDetails.Location = job.Location{
			CityID:      jobs[i].JobDetails.City,
			CityName:    cityName,
			Country:     countryID,
			Subdivision: subdivision,
		}

		if userID != "" {
			// isApplied
			jobs[i].IsApplied, err = s.jobs.IsJobApplied(ctx, userID, jobs[i].GetID())
			if err != nil {
				s.tracer.LogError(span, err)
			}
			// isSaved
			jobs[i].IsSaved, err = s.jobs.IsJobSaved(ctx, userID, jobs[i].GetID())
			if err != nil {
				s.tracer.LogError(span, err)
			}
		} else {
			s.tracer.LogError(span, err)
		}
	}

	return jobs, nil
}

// ReportJob ...
func (s *Service) ReportJob(ctx context.Context, jobID string, reportType job.ReportType, text string) error {
	span := s.tracer.MakeSpan(ctx, "ReportJob")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	report := job.Report{
		Text: text,
		Type: reportType,
	}
	report.GenerateID()
	report.SetJobID(jobID)
	report.SetReporterID(userID)

	err = s.jobs.ReportJob(ctx, &report)
	if err != nil {
		return err
	}

	return nil
}

// ReportCandidate ...
func (s *Service) ReportCandidate(ctx context.Context, companyID, candidateID string, text string) error {
	span := s.tracer.MakeSpan(ctx, "ReportCandidate")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	report := candidate.Report{
		Text: text,
	}
	report.GenerateID()
	report.SetReporterUserID(userID)
	report.SetReporterCompanyID(companyID)
	report.SetCandidateID(candidateID)

	err = s.jobs.ReportCandidate(ctx, &report)
	if err != nil {
		return err
	}

	return nil
}

// SaveJobSearchFilter ...
func (s *Service) SaveJobSearchFilter(ctx context.Context, filter *job.NamedSearchFilter) error {
	span := s.tracer.MakeSpan(ctx, "SaveJobSearchFilter")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	filter.GenerateID()
	filter.SetUserID(userID)

	err = s.jobs.SaveJobSearchFilter(ctx, filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetSavedJobSearchFilters ...
func (s *Service) GetSavedJobSearchFilters(ctx context.Context) ([]*job.NamedSearchFilter, error) {
	span := s.tracer.MakeSpan(ctx, "GetSavedJobSearchFilters")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	filters, err := s.jobs.GetSavedJobSearchFilters(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return filters, nil
}

// SaveCandidateSearchFilter ...
func (s *Service) SaveCandidateSearchFilter(ctx context.Context, companyID string, filter *candidate.NamedSearchFilter) error {
	span := s.tracer.MakeSpan(ctx, "SaveCandidateSearchFilter")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	filter.GenerateID()
	filter.SetCompanyID(companyID)

	err := s.jobs.SaveCandidateSearchFilter(ctx, companyID, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetSavedCandidateSearchFilters ...
func (s *Service) GetSavedCandidateSearchFilters(ctx context.Context, companyID string) ([]*candidate.NamedSearchFilter, error) {
	span := s.tracer.MakeSpan(ctx, "GetSavedCandidateSearchFilters")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	filters, err := s.jobs.GetSavedCandidateSearchFilters(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return filters, nil
}

// GetPlanPrices ...
func (s *Service) GetPlanPrices(ctx context.Context, companyID string, countriesID []string, currency string) ([]job.PlanPrices, error) {
	span := s.tracer.MakeSpan(ctx, "GetPlanPrices")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	plans, err := s.jobs.GetPlanPrices(ctx, countriesID, currency)
	if err != nil {
		return nil, err
	}

	return plans, nil
}

// GetPricingFor ...
func (s *Service) GetPricingFor(ctx context.Context, companyID string, meta *job.Meta) (*job.PricingResult, error) {
	span := s.tracer.MakeSpan(ctx, "GetPricingFor")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	filters, err := s.jobs.GetPricingFor(ctx, meta)
	if err != nil {
		return nil, err
	}

	return filters, nil
}

// GetAmountOfActiveJobsOfCompany ...
func (s *Service) GetAmountOfActiveJobsOfCompany(ctx context.Context, companyID string) (int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetAmountOfActiveJobsOfCompany")
	defer span.Finish()

	amount, err := s.jobs.GetAmountOfActiveJobsOfCompany(ctx, companyID)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// GetCareerInterestsByIds ...
func (s *Service) GetCareerInterestsByIds(ctx context.Context, companyID string, ids []string, first, after int32) (map[string]*candidate.CareerInterests, error) {
	span := s.tracer.MakeSpan(ctx, "GetCareerInterestsByIds")
	defer span.Finish()

	careers, err := s.jobs.GetCareerInterestsByIds(ctx, ids, first, after)
	if err != nil {
		return nil, err
	}
	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if allowed {
		for id := range careers {
			isSaved, err := s.jobs.IsCandidateSaved(ctx, companyID, id)
			if err != nil {
				s.tracer.LogError(span, err)
			}
			if careers[id] == nil {
				careers[id] = new(candidate.CareerInterests)
			}
			careers[id].IsSaved = isSaved
		}
	}

	return careers, nil
}

func (s *Service) UploadFileForJob(ctx context.Context, companyID, jobID string, file *job.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "UploadFileForJob")
	defer span.Finish()

	id := file.GenerateID()

	// check admin level of company
	if companyID != "" {
		allowed := s.checkAdminLevel(
			ctx,
			companyID,
			companyadmin.AdminLevelAdmin,
			companyadmin.AdminLevelVShop,
		)
		if !allowed {
			return "", errors.New("not_enought_authenticitation")
		}
	}

	err := s.jobs.UploadFileForJob(ctx, companyID, jobID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// UploadFileForApplication ...
func (s *Service) UploadFileForApplication(ctx context.Context, userID string, jobID string, file *job.File) (string, error) {
	span := s.tracer.MakeSpan(ctx, "UploadFileForApplication")
	defer span.Finish()

	id := file.GenerateID()

	err := s.jobs.UploadFileForApplication(ctx, userID, jobID, file)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// AddCVInCareerCenter ...
func (s *Service) AddCVInCareerCenter(ctx context.Context, companyID string, options careercenter.CVOptions) error {
	span := s.tracer.MakeSpan(ctx, "AddCVInCareerCenter")
	defer span.Finish()

	//  get userID
	userID, err := s.authRPC.GetUserID(ctx)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// TODO: check if not added before

	err = s.jobs.AddCVInCareerCenter(ctx, userID, companyID, options)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// GetSavedCVs ...
func (s *Service) GetSavedCVs(ctx context.Context, companyID string, first, after int32) ([]*candidate.ViewForCompany, error) {
	span := s.tracer.MakeSpan(ctx, "GetSavedCVs")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
	)
	if !allowed {
		return nil, errors.New("not_allowed")
	}

	if first <= 0 {
		first = 10
	}

	candidates, err := s.jobs.GetSavedCVs(ctx, companyID, first, after)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return candidates, nil
}

// RemoveCVs ...
func (s *Service) RemoveCVs(ctx context.Context, companyID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveCVs")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.jobs.RemoveCVs(ctx, companyID, ids)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// MakeFavoriteCVs ...
func (s *Service) MakeFavoriteCVs(ctx context.Context, companyID string, ids []string, isFavourite bool) error {
	span := s.tracer.MakeSpan(ctx, "MakeFavoriteCVs")
	defer span.Finish()

	// check admin level
	allowed := s.checkAdminLevel(
		ctx,
		companyID,
		companyadmin.AdminLevelAdmin,
		companyadmin.AdminLevelJob,
	)
	if !allowed {
		return errors.New("not_allowed")
	}

	err := s.jobs.MakeFavoriteCVs(ctx, companyID, ids, isFavourite)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// -----------

// return false if level doesn't much
func (s Service) checkAdminLevel(ctx context.Context, companyID string, requiredLevels ...companyadmin.AdminLevel) bool {
	span := s.tracer.MakeSpan(ctx, "checkAdminLevel")
	defer span.Finish()

	actualLevel, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		log.Println("Error: checkAdminLevel:", err)
		return false
	}

	for _, lvl := range requiredLevels {
		if lvl == actualLevel {
			return true
		}
	}

	return false
}
