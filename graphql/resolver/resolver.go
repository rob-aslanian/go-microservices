package resolver

import (
	"context"
	"log"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/rentalRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/advertRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/newsfeedRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/statisticsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
	"google.golang.org/grpc/metadata"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"

	"github.com/graph-gophers/graphql-go"

	rpc "gitlab.lan/Rightnao-site/microservices/graphql/grpc"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/searchRPC"

	// "gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/groupsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/notificationsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/shopRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

var (
	user          userRPC.UserServiceClient
	auth          authRPC.AuthServiceClient
	info          infoRPC.InfoServiceClient
	company       companyRPC.CompanyServiceClient
	services      servicesRPC.ServicesServiceClient
	network       networkRPC.NetworkServiceClient
	search        searchRPC.SearchServiceClient
	jobs          jobsRPC.JobsServiceClient
	chat          chatRPC.ChatServiceClient
	advert        advertRPC.AdvertServiceClient
	stuff         stuffRPC.StuffServiceClient
	statistics    statisticsRPC.StatisticsClient
	groups        groupsRPC.GroupsServiceClient
	notifications notificationsRPC.NotificationsServiceClient
	newsfeed      newsfeedRPC.NewsfeedServiceClient
	shop          shopRPC.ShopServiceClient
	rental        rentalRPC.RentalServiceClient
)

type Resolver struct {
	AddedPostEvents     chan *NewsfeedPostResolverCustom
	addedPostSubscriber chan *addedPostSubscriber
	// -----
	AddedCommentPostEvents     chan *NewsfeedPostCommentResolverCustom
	addedCommentPostSubscriber chan *addedCommentPostSubscriber
	// -----
	AddedLikePostEvents     chan *LikeResolverCustom
	addedLikePostSubscriber chan *addedLikePostSubscriber
}

func (r *Resolver) Init() {
	r.AddedPostEvents = make(chan *NewsfeedPostResolverCustom)
	r.addedPostSubscriber = make(chan *addedPostSubscriber)
	go r.broadcastAddedPost()
	r.AddedCommentPostEvents = make(chan *NewsfeedPostCommentResolverCustom)
	r.addedCommentPostSubscriber = make(chan *addedCommentPostSubscriber)
	go r.broadcastAddedCommentPost()
	r.AddedLikePostEvents = make(chan *LikeResolverCustom)
	r.addedLikePostSubscriber = make(chan *addedLikePostSubscriber)
	go r.broadcastAddedLikePost()

	user = rpc.GetGrpcClient().UserClient
	auth = rpc.GetGrpcClient().AuthClient
	info = rpc.GetGrpcClient().InfoClient
	company = rpc.GetGrpcClient().CompanyClient
	network = rpc.GetGrpcClient().NetworkClient
	search = rpc.GetGrpcClient().SearchClient
	jobs = rpc.GetGrpcClient().JobsClient
	chat = rpc.GetGrpcClient().ChatClient
	advert = rpc.GetGrpcClient().AdvertClient
	stuff = rpc.GetGrpcClient().StuffClient
	statistics = rpc.GetGrpcClient().StatisticsClient
	notifications = rpc.GetGrpcClient().NotificationsClient
	services = rpc.GetGrpcClient().ServicesClient
	newsfeed = rpc.GetGrpcClient().NewsfeedClient
	groups = rpc.GetGrpcClient().GroupsClient
	shop = rpc.GetGrpcClient().ShopClient
	rental = rpc.GetGrpcClient().RentalClient
}

// extra part, that wasn't generated automatically

func (r AccountResolver) Sessions(ctx context.Context, args SessionsRequest) ([]SessionsResolver, error) {
	var first int32 = 2
	var after int32 = 0

	if args.First != nil {
		first = int32(*args.First)
	}
	if args.After != nil {
		after = int32(*args.After)
	}

	sessions, err := auth.GetListOfSessions(ctx, &authRPC.ListOfSessionsQuery{
		After: after,
		First: first,
	})
	if err != nil {
		return nil, err
	}

	ses := make([]SessionsResolver, 0, len(sessions.GetSessions()))

	for _, s := range sessions.GetSessions() {
		se := Sessions{
			Browser_version:    s.GetBrowserVersion(),
			Os:                 s.GetOS(),
			Os_version:         s.GetOSVersion(),
			Browser:            s.GetBrowser(),
			Last_activity_time: s.GetLastActivityTime(),
			ID:                 s.GetID(),
			Device_type:        s.GetDeviceType(),
			// City:               s.GetCity(),
			// Country_id:      s.GetCountryID(),
			Current_session: s.GetCurrentSession(),
		}

		se.Location = &City{
			Country: s.GetCountryID(),
		}

		n := strconv.Itoa(int(s.GetCity()))
		city, err := info.GetCityInfoByID(ctx, &infoRPC.IDWithLang{
			ID: n,
		})
		if err != nil {
			log.Println("error: get city info")
		} else {
			se.Location.City = city.GetTitle()
			se.Location.ID = strconv.Itoa(int(city.GetId()))
			se.Location.Subdivision = city.GetSubdivision()
		}

		ses = append(ses, SessionsResolver{
			R: &se,
		})
	}

	return ses, nil
}

func (r ProfileResolver) Experiences(ctx context.Context, args ExperiencesRequest) (*[]*ExperienceProfileResolver, error) {

	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	exp, err := user.GetExperiences(ctx, &userRPC.RequestExperiences{
		UserID:   string(r.ID()),
		First:    first,
		After:    after,
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	experiences := make([]*ExperienceProfileResolver, len(exp.GetExperiences()))

	for i := range experiences {
		e := ExperienceProfile{
			ID:          exp.GetExperiences()[i].GetID(),
			Title:       exp.GetExperiences()[i].GetPosition(),
			Company:     exp.GetExperiences()[i].GetCompany(),
			Description: exp.GetExperiences()[i].GetDescription(),
			Start_date:  exp.GetExperiences()[i].GetStartDate(),
			Finish_date: exp.GetExperiences()[i].GetFinishDate(),
			Currently:   exp.GetExperiences()[i].GetCurrentlyWork(),
			File:        make([]File, len(exp.GetExperiences()[i].GetFiles())),
		}

		for i, f := range exp.GetExperiences()[i].GetFiles() {
			e.File[i] = File{
				Name:      f.GetName(),
				Address:   f.GetURL(),
				ID:        f.GetID(),
				Mime_type: f.GetMimeType(),
			}
		}

		// e.Link = exp.GetExperiences()[i].GetLinks()
		e.Link = func(f []*userRPC.Link) []Link {
			links := make([]Link, len(f))

			for i := range f {
				var link Link
				link.ID = f[i].GetID()
				link.Address = f[i].GetURL()
				// TODO:

				links[i] = link
			}

			return links
		}(exp.GetExperiences()[i].GetLinks())

		// c := exp.GetExperiences()[i].GetCityID()
		// e.Location = &Location{
		// 	City:    &City{},
		// 	Country: &Country{},
		// }
		// e.Location.City.ID = strconv.Itoa(int(c))
		// e.Location.City.City = exp.GetExperiences()[i].GetLocation()

		loc := exp.GetExperiences()[i].GetLocation()

		e.Location = &Location{
			City:    &City{},
			Country: &Country{},
		}

		if loc != nil {
			if loc.GetCity() != nil {
				e.Location.City.ID = strconv.Itoa(int(loc.GetCity().GetId()))
				e.Location.City.City = loc.GetCity().GetTitle()
				e.Location.City.Subdivision = loc.GetCity().GetSubdivision()
			}
			if loc.GetCountry() != nil {
				e.Location.Country.ID = loc.GetCountry().GetId()
			}
		}

		experiences[i] = &ExperienceProfileResolver{R: &e}
	}

	return &experiences, nil
}

func (r ProfileResolver) Educations(ctx context.Context, args EducationsRequest) (*[]*EducationProfileResolver, error) {
	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	exp, err := user.GetEducations(ctx, &userRPC.RequestEducations{
		UserID:   string(r.ID()),
		First:    first,
		After:    after,
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	experiences := make([]*EducationProfileResolver, len(exp.GetEducations()))

	for i := range experiences {
		var e EducationProfile

		e.ID = exp.GetEducations()[i].GetID()
		e.School = exp.GetEducations()[i].GetSchool()
		e.Degree = exp.GetEducations()[i].GetDegree()
		e.Field_study = exp.GetEducations()[i].GetFieldStudy()
		e.Grade = exp.GetEducations()[i].GetGrade()

		e.Description = exp.GetEducations()[i].GetDescription()

		e.Start_date = exp.GetEducations()[i].GetStartDate()
		e.Finish_date = exp.GetEducations()[i].GetFinishDate()
		e.Currently_study = exp.GetEducations()[i].GetIsCurrentlyStudy()

		e.File = func(f []*userRPC.File) []File {
			files := make([]File, len(f))

			for i := range f {
				var file File
				file.Name = f[i].GetName()
				file.Address = f[i].GetURL()
				file.ID = f[i].GetID()
				file.Mime_type = f[i].GetMimeType()

				files[i] = file
			}

			return files
		}(exp.GetEducations()[i].GetFiles())

		// e.Link = exp.GetExperiences()[i].GetLinks()
		e.Link = func(f []*userRPC.Link) []Link {
			links := make([]Link, len(f))

			for i := range f {
				var link Link
				link.ID = f[i].GetID()
				link.Address = f[i].GetURL()
				links[i] = link
			}

			return links
		}(exp.GetEducations()[i].GetLinks())

		loc := exp.GetEducations()[i].GetLocation()
		e.Location = &Location{
			City:    &City{},
			Country: &Country{},
		}
		if loc != nil {
			if loc.GetCity() != nil {
				e.Location.City.ID = strconv.Itoa(int(loc.GetCity().GetId()))
				e.Location.City.City = loc.GetCity().GetTitle()
				e.Location.City.Subdivision = loc.GetCity().GetSubdivision()
			}
			if loc.GetCountry() != nil {
				e.Location.Country.ID = loc.GetCountry().GetId()
			}
		}

		experiences[i] = &EducationProfileResolver{R: &e}
	}

	return &experiences, nil
}

func (r ProfileResolver) Skills(ctx context.Context, args EducationsRequest) (*[]*SkillProfileResolver, error) {
	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	exp, err := user.GetSkills(ctx, &userRPC.RequestSkills{
		UserID:   string(r.ID()),
		First:    first,
		After:    after,
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	experiences := make([]*SkillProfileResolver, len(exp.GetSkills()))

	for i := range experiences {
		var e SkillProfile

		e.ID = exp.GetSkills()[i].GetID()
		e.Name = exp.GetSkills()[i].GetSkill()
		// e.Amount_endorsements = exp.GetSkills()[i].Get // TODO:

		experiences[i] = &SkillProfileResolver{R: &e, language: r.language}
	}

	return &experiences, nil
}

func (r SkillProfileResolver) Endorsements(ctx context.Context, args EndorsementsRequest) (*[]*ProfileResolver, error) {
	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	profile, err := user.GetEndorsements(ctx, &userRPC.RequestEndorsements{
		SkillID:  r.R.ID,
		First:    first,
		After:    after,
		Language: r.language,
	},
	)
	if err != nil {
		return nil, err
	}

	profiles := make([]*ProfileResolver, len(profile.GetProfiles()))

	for i := range profiles {
		pr := ToProfile(ctx, profile.GetProfiles()[i])
		profiles[i] = &ProfileResolver{R: &pr}
	}

	return &profiles, nil
}

func (r ProfileResolver) Interests(ctx context.Context, args EducationsRequest) (*[]*InterestProfileResolver, error) {
	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	exp, err := user.GetInterests(ctx, &userRPC.RequestInterests{
		UserID:   string(r.ID()),
		First:    first,
		After:    after,
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	experiences := make([]*InterestProfileResolver, len(exp.GetInterests()))

	for i := range experiences {
		var e InterestProfile

		e.ID = exp.GetInterests()[i].GetID()

		if exp.GetInterests()[i].GetIsImageNull() == false {
			e.Image = exp.GetInterests()[i].GetImage()
		}
		if exp.GetInterests()[i].GetIsInterestNull() == false {
			e.Interest = exp.GetInterests()[i].GetInterest()
		}
		if exp.GetInterests()[i].GetIsDescriptionNull() == false {
			e.Description = exp.GetInterests()[i].GetDescription()
		}

		experiences[i] = &InterestProfileResolver{R: &e}
	}

	return &experiences, nil
}

// Accomplishments ...
func (r ProfileResolver) Accomplishments(ctx context.Context, args AccomplishmentsRequest) (*[]*AccomplishmentResolver, error) {

	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	exp, err := user.GetAccomplishments(ctx, &userRPC.RequestAccomplshments{
		UserID:   string(r.ID()),
		First:    first,
		After:    after,
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	accomplishments := make([]*AccomplishmentResolver, len(exp.GetAccomplishments()))

	for i := range accomplishments {

		var e Accomplishment

		switch exp.GetAccomplishments()[i].GetType() {
		case userRPC.Accomplishment_Award:
			res := AwardResolver{
				R: &Award{
					ID:          exp.GetAccomplishments()[i].GetID(),
					Title:       exp.GetAccomplishments()[i].GetName(),
					Date:        exp.GetAccomplishments()[i].GetFinishDate(),
					Description: exp.GetAccomplishments()[i].GetDescription(),
					Issuer:      exp.GetAccomplishments()[i].GetIssuer(),
					File:        make([]File, len(exp.GetAccomplishments()[i].GetFiles())),
				},
			}
			if exp.GetAccomplishments()[i].IsIssuerNull == false {
				res.R.Issuer = exp.GetAccomplishments()[i].GetIssuer()
			}
			if exp.GetAccomplishments()[i].IsFinishDateNull == false {
				res.R.Date = exp.GetAccomplishments()[i].GetFinishDate()
			}
			if exp.GetAccomplishments()[i].IsDescriptionNull == false {
				res.R.Description = exp.GetAccomplishments()[i].GetDescription()
			}
			res.R.Link = func(f []*userRPC.Link) []Link {
				links := make([]Link, len(f))

				for i := range f {
					var link Link
					link.ID = f[i].GetID()
					link.Address = f[i].GetURL()
					// TODO:

					links[i] = link
				}

				return links
			}(exp.GetAccomplishments()[i].GetLinks())

			res.R.File = func(f []*userRPC.File) []File {
				files := make([]File, len(f))

				for i := range f {
					var file File
					file.Name = f[i].GetName()
					file.Address = f[i].GetURL()
					file.ID = f[i].GetID()
					file.Mime_type = f[i].GetMimeType()

					files[i] = file
				}

				return files
			}(exp.GetAccomplishments()[i].GetFiles())

			e = res

		case userRPC.Accomplishment_Test:
			res := TestResolver{
				R: &Test{
					ID:          exp.GetAccomplishments()[i].GetID(),
					Title:       exp.GetAccomplishments()[i].GetName(),
					Score:       int32(exp.GetAccomplishments()[i].GetScore()),
					Description: exp.GetAccomplishments()[i].GetDescription(),
					Date:        exp.GetAccomplishments()[i].GetFinishDate(),
					File:        make([]File, len(exp.GetAccomplishments()[i].GetFiles())),
				},
			}
			res.R.Link = func(f []*userRPC.Link) []Link {
				links := make([]Link, len(f))

				for i := range f {
					var link Link
					link.ID = f[i].GetID()
					link.Address = f[i].GetURL()
					// TODO:

					links[i] = link
				}

				return links
			}(exp.GetAccomplishments()[i].GetLinks())

			res.R.File = func(f []*userRPC.File) []File {
				files := make([]File, len(f))

				for i := range f {
					var file File
					file.Name = f[i].GetName()
					file.Address = f[i].GetURL()
					file.ID = f[i].GetID()
					file.Mime_type = f[i].GetMimeType()

					files[i] = file
				}

				return files
			}(exp.GetAccomplishments()[i].GetFiles())

			if exp.GetAccomplishments()[i].IsFinishDateNull == false {
				res.R.Date = exp.GetAccomplishments()[i].GetFinishDate()
			}
			if exp.GetAccomplishments()[i].IsDescriptionNull == false {
				res.R.Description = exp.GetAccomplishments()[i].GetDescription()
			}

			e = res

		case userRPC.Accomplishment_Project:
			res := ProjectResolver{
				R: &Project{
					ID:          exp.GetAccomplishments()[i].GetID(),
					Name:        exp.GetAccomplishments()[i].GetName(),
					Description: exp.GetAccomplishments()[i].GetDescription(),
					Start_date:  exp.GetAccomplishments()[i].GetStartDate(),
					Finish_date: exp.GetAccomplishments()[i].GetFinishDate(),
					Url:         exp.GetAccomplishments()[i].GetURL(),
					File:        make([]File, len(exp.GetAccomplishments()[i].GetFiles())),
				},
			}
			if exp.GetAccomplishments()[i].IsURLNull == false {
				res.R.Url = exp.GetAccomplishments()[i].GetURL()
			}
			if exp.GetAccomplishments()[i].IsDescriptionNull == false {
				res.R.Description = exp.GetAccomplishments()[i].GetDescription()
			}
			if exp.GetAccomplishments()[i].IsFinishDateNull == false {
				res.R.Finish_date = exp.GetAccomplishments()[i].GetFinishDate()
			}
			if exp.GetAccomplishments()[i].IsStartDateNull == false {
				res.R.Start_date = exp.GetAccomplishments()[i].GetStartDate()
			}
			if exp.GetAccomplishments()[i].IsIsExpireNull == false {
				res.R.Is_expire = exp.GetAccomplishments()[i].GetIsExpire()
			}
			res.R.Link = func(f []*userRPC.Link) []Link {
				links := make([]Link, len(f))

				for i := range f {
					var link Link
					link.ID = f[i].GetID()
					link.Address = f[i].GetURL()
					// TODO:

					links[i] = link
				}

				return links
			}(exp.GetAccomplishments()[i].GetLinks())

			res.R.File = func(f []*userRPC.File) []File {
				files := make([]File, len(f))

				for i := range f {
					var file File
					file.Name = f[i].GetName()
					file.Address = f[i].GetURL()
					file.ID = f[i].GetID()
					file.Mime_type = f[i].GetMimeType()

					files[i] = file
				}

				return files
			}(exp.GetAccomplishments()[i].GetFiles())

			e = res

		case userRPC.Accomplishment_License:
			res := LicenseResolver{
				R: &License{
					ID:             exp.GetAccomplishments()[i].GetID(),
					Name:           exp.GetAccomplishments()[i].GetName(),
					File:           make([]File, len(exp.GetAccomplishments()[i].GetFiles())),
					License_number: exp.GetAccomplishments()[i].GetLicenseNumber(),
					Start_date:     exp.GetAccomplishments()[i].GetStartDate(),
					Finish_date:    exp.GetAccomplishments()[i].GetFinishDate(),
					Is_expire:      exp.GetAccomplishments()[i].GetIsExpire(),
					Issuer:         exp.GetAccomplishments()[i].GetIssuer(),
				},
			}
			if exp.GetAccomplishments()[i].IsIssuerNull == false {
				res.R.Issuer = exp.GetAccomplishments()[i].GetIssuer()
			}
			if exp.GetAccomplishments()[i].IsLicenseNumberNull == false {
				res.R.License_number = exp.GetAccomplishments()[i].GetLicenseNumber()
			}
			if exp.GetAccomplishments()[i].IsFinishDateNull == false {
				res.R.Finish_date = exp.GetAccomplishments()[i].GetFinishDate()
			}
			if exp.GetAccomplishments()[i].IsStartDateNull == false {
				res.R.Start_date = exp.GetAccomplishments()[i].GetStartDate()
			}
			if exp.GetAccomplishments()[i].IsIsExpireNull == false {
				res.R.Is_expire = exp.GetAccomplishments()[i].GetIsExpire()
			}
			res.R.Link = func(f []*userRPC.Link) []Link {
				links := make([]Link, len(f))

				for i := range f {
					var link Link
					link.ID = f[i].GetID()
					link.Address = f[i].GetURL()
					// TODO:

					links[i] = link
				}

				return links
			}(exp.GetAccomplishments()[i].GetLinks())

			res.R.File = func(f []*userRPC.File) []File {
				files := make([]File, len(f))

				for i := range f {
					var file File
					file.Name = f[i].GetName()
					file.Address = f[i].GetURL()
					file.ID = f[i].GetID()
					file.Mime_type = f[i].GetMimeType()

					files[i] = file
				}

				return files
			}(exp.GetAccomplishments()[i].GetFiles())

			e = res

		case userRPC.Accomplishment_Certificate:
			res := CertificationResolver{
				R: &Certification{
					ID:             exp.GetAccomplishments()[i].GetID(),
					Name:           exp.GetAccomplishments()[i].GetName(),
					Start_date:     exp.GetAccomplishments()[i].GetStartDate(),
					Finish_date:    exp.GetAccomplishments()[i].GetFinishDate(),
					Is_expire:      exp.GetAccomplishments()[i].GetIsExpire(),
					Url:            exp.GetAccomplishments()[i].GetURL(),
					License_number: exp.GetAccomplishments()[i].GetLicenseNumber(),
					File:           make([]File, len(exp.GetAccomplishments()[i].GetFiles())),
				},
			}
			if exp.GetAccomplishments()[i].IsIssuerNull == false {
				res.R.Certification_authority = exp.GetAccomplishments()[i].GetIssuer()
			}
			if exp.GetAccomplishments()[i].IsLicenseNumberNull == false {
				res.R.License_number = exp.GetAccomplishments()[i].GetLicenseNumber()
			}
			if exp.GetAccomplishments()[i].IsFinishDateNull == false {
				res.R.Finish_date = exp.GetAccomplishments()[i].GetFinishDate()
			}
			if exp.GetAccomplishments()[i].IsStartDateNull == false {
				res.R.Start_date = exp.GetAccomplishments()[i].GetStartDate()
			}
			if exp.GetAccomplishments()[i].IsIsExpireNull == false {
				res.R.Is_expire = exp.GetAccomplishments()[i].GetIsExpire()
			}
			if exp.GetAccomplishments()[i].IsURLNull == false {
				res.R.Url = exp.GetAccomplishments()[i].GetURL()
			}
			res.R.Link = func(f []*userRPC.Link) []Link {
				links := make([]Link, len(f))

				for i := range f {
					var link Link
					link.ID = f[i].GetID()
					link.Address = f[i].GetURL()
					// TODO:

					links[i] = link
				}

				return links
			}(exp.GetAccomplishments()[i].GetLinks())

			res.R.File = func(f []*userRPC.File) []File {
				files := make([]File, len(f))

				for i := range f {
					var file File
					file.Name = f[i].GetName()
					file.Address = f[i].GetURL()
					file.ID = f[i].GetID()
					file.Mime_type = f[i].GetMimeType()

					files[i] = file
				}

				return files
			}(exp.GetAccomplishments()[i].GetFiles())

			e = res

		case userRPC.Accomplishment_Publication:
			res := PublicationResolver{
				R: &Publication{
					ID:          exp.GetAccomplishments()[i].GetID(),
					Title:       exp.GetAccomplishments()[i].GetName(),
					Publisher:   exp.GetAccomplishments()[i].GetIssuer(),
					Description: exp.GetAccomplishments()[i].GetDescription(),
					Url:         exp.GetAccomplishments()[i].GetURL(),
					File:        make([]File, len(exp.GetAccomplishments()[i].GetFiles())),
				},
			}

			if exp.GetAccomplishments()[i].IsIssuerNull == false {
				res.R.Publisher = exp.GetAccomplishments()[i].GetIssuer()
			}
			if exp.GetAccomplishments()[i].IsFinishDateNull == false {
				res.R.Date = exp.GetAccomplishments()[i].GetFinishDate()
			}
			if exp.GetAccomplishments()[i].IsURLNull == false {
				res.R.Url = exp.GetAccomplishments()[i].GetURL()
			}
			if exp.GetAccomplishments()[i].IsDescriptionNull == false {
				res.R.Description = exp.GetAccomplishments()[i].GetDescription()
			}
			res.R.Link = func(f []*userRPC.Link) []Link {
				links := make([]Link, len(f))

				for i := range f {
					var link Link
					link.ID = f[i].GetID()
					link.Address = f[i].GetURL()
					// TODO:

					links[i] = link
				}

				return links
			}(exp.GetAccomplishments()[i].GetLinks())

			res.R.File = func(f []*userRPC.File) []File {
				files := make([]File, len(f))

				for i := range f {
					var file File
					file.Name = f[i].GetName()
					file.Address = f[i].GetURL()
					file.ID = f[i].GetID()
					file.Mime_type = f[i].GetMimeType()

					files[i] = file
				}

				return files
			}(exp.GetAccomplishments()[i].GetFiles())

			e = res
		}
		accomplishments[i] = &AccomplishmentResolver{r: &e}
	}

	return &accomplishments, nil
}

// func (r ProfileResolver) Portfolios(ctx context.Context, args PortfoliosRequest) (*[]*PortfolioProfileResolver, error) {

// 	var first uint32 = 2
// 	var after string = "0"

// 	if args.First != nil {
// 		first = uint32(*args.First)
// 	}
// 	if args.After != nil {
// 		after = *args.After
// 	}

// 	port, err := user.GetPortfolios(ctx, &userRPC.RequestPortfolios{
// 		UserID: string(r.ID()),
// 		First:  first,
// 		After:  after,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	portfolios := make([]*PortfolioProfileResolver, len(port.GetPortfolios()))

// 	for i := range portfolios {
// 		p := PortfolioProfile{
// 			ID:          port.GetPortfolios()[i].GetID(),
// 			Title:       port.GetPortfolios()[i].GetTittle(),
// 			Description: port.GetPortfolios()[i].GetDescription(),
// 			ContentType: userRPCContentTypeEnumToString(port.GetPortfolios()[i].GetContentType()),
// 			Link:        make([]Link, len(port.GetPortfolios()[i].GetLinks())),
// 			Files:       make([]File, len(port.GetPortfolios()[i].GetFiles())),
// 		}

// 		for i, f := range port.GetPortfolios()[i].GetFiles() {
// 			p.Files[i] = File{
// 				Name:      f.GetName(),
// 				Address:   f.GetURL(),
// 				ID:        f.GetID(),
// 				Mime_type: f.GetMimeType(),
// 			}
// 		}

// 		p.Link = func(f []*userRPC.Link) []Link {
// 			links := make([]Link, len(f))

// 			for i := range f {
// 				var link Link
// 				link.ID = f[i].GetID()
// 				link.Address = f[i].GetURL()
// 				// TODO:

// 				links[i] = link
// 			}

// 			return links
// 		}(port.GetPortfolios()[i].GetLinks())

// 		portfolios[i] = &PortfolioProfileResolver{R: &p}
// 	}

// 	return &portfolios, nil
// }

func (r CompanyProfileResolver) Gallery(ctx context.Context, args GalleryRequest) (GalleryProfileResolver, error) {

	var first uint32 = 2
	var after uint32 = 0

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		first = uint32(*args.First)
	}

	gallery, err := company.GetCompanyGallery(ctx, &companyRPC.RequestGallery{
		CompanyID: string(r.ID()),
		First:     first,
		After:     after,
	})
	if err != nil {
		return GalleryProfileResolver{}, err
	}

	files := make([]File, 0, len(gallery.GetGalleryFiles()))

	for _, f := range gallery.GetGalleryFiles() {
		files = append(files, File{
			Address:   f.GetURL(),
			ID:        f.GetID(),
			Mime_type: f.GetMimeType(),
			Name:      f.GetName(),
		})
	}

	gal := GalleryProfile{
		Files: files,
	}

	return GalleryProfileResolver{
		R: &gal,
	}, nil
}

func (r ProfileResolver) ToolsTechnologies(ctx context.Context, args ToolsTechnologiesRequest) (*[]*ToolsTechnologiesProfileResolver, error) {

	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	tool, err := user.GetToolsTechnologies(ctx, &userRPC.RequestToolTechnology{
		UserID:   string(r.ID()),
		First:    first,
		After:    after,
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	tools := make([]*ToolsTechnologiesProfileResolver, len(tool.GetToolsTechnologies()))

	for i := range tools {
		p := ToolsTechnologiesProfile{
			ID:              tool.GetToolsTechnologies()[i].GetID(),
			Tool_Technology: tool.GetToolsTechnologies()[i].GetToolTechnology(),
			Rank:            userRPCToolsLevelToToolsLevel(tool.GetToolsTechnologies()[i].GetRank()),
		}

		tools[i] = &ToolsTechnologiesProfileResolver{R: &p}
	}

	return &tools, nil
}

func (r AccomplishmentResolver) ID() graphql.ID {
	return (*r.r).ID()
}

func (r AccomplishmentResolver) ToLicense() (*LicenseResolver, bool) {
	res, ok := (*r.r).(LicenseResolver)
	return &res, ok
}

func (r AccomplishmentResolver) ToCertification() (*CertificationResolver, bool) {
	res, ok := (*r.r).(CertificationResolver)
	return &res, ok
}

func (r AccomplishmentResolver) ToAward() (*AwardResolver, bool) {
	res, ok := (*r.r).(AwardResolver)
	return &res, ok
}

func (r AccomplishmentResolver) ToProject() (*ProjectResolver, bool) {
	res, ok := (*r.r).(ProjectResolver)
	return &res, ok
}

func (r AccomplishmentResolver) ToPublication() (*PublicationResolver, bool) {
	res, ok := (*r.r).(PublicationResolver)
	return &res, ok
}

func (r AccomplishmentResolver) ToTest() (*TestResolver, bool) {
	res, ok := (*r.r).(TestResolver)
	return &res, ok
}

// languages
func (r ProfileResolver) Languages(ctx context.Context, args EducationsRequest) (*[]*KnownLanguageProfileResolver, error) {

	var first uint32 = 2
	var after string = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	exp, err := user.GetKnownLanguages(ctx, &userRPC.RequestKnownLanguages{
		UserID: string(r.ID()),
		First:  first,
		After:  after,
	})
	if err != nil {
		return nil, err
	}

	languages := make([]*KnownLanguageProfileResolver, len(exp.GetKnownLanguages()))

	for i := range languages {
		var e KnownLanguageProfile

		e.ID = exp.GetKnownLanguages()[i].GetID()
		e.Rate = int32(exp.GetKnownLanguages()[i].GetRank())
		e.Language = exp.GetKnownLanguages()[i].GetLanguage()

		languages[i] = &KnownLanguageProfileResolver{R: &e}
	}

	return &languages, nil
}

// company
func (r CompanyProfileResolver) Founders(ctx context.Context, args FoundersRequest) (*[]*CompanyFounderResolver, error) {

	var first int32 = 2
	var after string = "0"

	if args.First != nil {
		first = int32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	founders, err := company.GetFounders(ctx, &companyRPC.GetFoundersRequest{
		CompanyID: string(r.ID()),
		First:     first,
		After:     after,
	})
	if err != nil {
		return nil, err
	}

	foun := make([]*CompanyFounderResolver, 0, len(founders.GetFounders()))

	for i := range founders.GetFounders() {
		var e CompanyFounder

		e.Approved = founders.GetFounders()[i].GetIsApproved()
		e.ID = founders.GetFounders()[i].GetID()
		e.Position_title = founders.GetFounders()[i].GetPositionTitle()
		e.UserID = founders.GetFounders()[i].GetUserID()

		if founders.GetFounders()[i].GetUserID() != "" {
			prof, err := user.GetProfileByID(ctx, &userRPC.ID{
				ID: e.UserID,
			})
			if err != nil {
				continue
				// return nil, err
			}

			e.Avatar = prof.GetAvatar()
			e.Name = prof.GetFirstname() + " " + prof.GetLastname()

		} else {
			e.Avatar = founders.GetFounders()[i].GetAvatar()
			e.Name = founders.GetFounders()[i].GetName()

		}
		foun = append(foun, &CompanyFounderResolver{R: &e})
	}

	return &foun, nil
}

// recommendations
func (r ProfileResolver) RecievedRecommendation(ctx context.Context, args PaginationInput) (*[]RecommendationResolver, error) {
	var first uint32 = 2
	var after = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	recommendations, err := user.GetReceivedRecommendations(ctx, &userRPC.IDWithPagination{
		ID: string(r.ID()),
		Pagination: &userRPC.Pagination{
			First: first,
			After: after,
		},
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	rec := make([]RecommendationResolver, 0, len(recommendations.GetRecommendations()))

	for i := range recommendations.GetRecommendations() {
		var r Recommendation

		log.Printf("Reccomendation %+v", r)

		r.ID = recommendations.GetRecommendations()[i].GetID()
		r.Text = recommendations.GetRecommendations()[i].GetText()
		r.Receiver = ToProfile(ctx, recommendations.GetRecommendations()[i].GetReceiver())
		r.Recommendator = ToProfile(ctx, recommendations.GetRecommendations()[i].GetRecommendator())
		r.Created_at = time.Unix(recommendations.GetRecommendations()[i].GetCreatedAt(), 0)
		r.Title = recommendations.GetRecommendations()[i].GetTitle()
		r.Relation = relationRPCToString(recommendations.GetRecommendations()[i].GetRelation())

		if recommendations.GetRecommendations()[i].GetIsIsHiddenNull() {
			r.Is_hidden = nil
		} else {
			b := recommendations.GetRecommendations()[i].GetIsHidden()
			r.Is_hidden = &b
		}

		rec = append(rec, RecommendationResolver{R: &r})
	}

	return &rec, nil
}

func (r ProfileResolver) GivenRecommendations(ctx context.Context, args PaginationInput) (*[]RecommendationResolver, error) {
	var first uint32 = 2
	var after = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	recommendations, err := user.GetGivenRecommendations(ctx, &userRPC.IDWithPagination{
		ID: string(r.ID()),
		Pagination: &userRPC.Pagination{
			First: first,
			After: after,
		},
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	rec := make([]RecommendationResolver, 0, len(recommendations.GetRecommendations()))

	for i := range recommendations.GetRecommendations() {
		var r Recommendation

		r.ID = recommendations.GetRecommendations()[i].GetID()
		r.Text = recommendations.GetRecommendations()[i].GetText()
		r.Receiver = ToProfile(ctx, recommendations.GetRecommendations()[i].GetReceiver())
		r.Recommendator = ToProfile(ctx, recommendations.GetRecommendations()[i].GetRecommendator())
		r.Created_at = time.Unix(recommendations.GetRecommendations()[i].GetCreatedAt(), 0)
		r.Title = recommendations.GetRecommendations()[i].GetTitle()
		r.Relation = relationRPCToString(recommendations.GetRecommendations()[i].GetRelation())
		// r.Is_hidden

		rec = append(rec, RecommendationResolver{R: &r})
	}

	return &rec, nil
}

func (r ProfileResolver) ReceivedRecommendationRequests(ctx context.Context, args PaginationInput) (*[]RecommendationRequestResolver, error) {
	var first uint32 = 2
	var after = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	recommendations, err := user.GetReceivedRecommendationRequests(ctx, &userRPC.IDWithPagination{
		ID: string(r.ID()),
		Pagination: &userRPC.Pagination{
			First: first,
			After: after,
		},
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	rec := make([]RecommendationRequestResolver, 0, len(recommendations.GetRecommendationRequests()))

	for i := range recommendations.GetRecommendationRequests() {
		var r RecommendationRequest

		r.ID = recommendations.GetRecommendationRequests()[i].GetID()
		r.Text = recommendations.GetRecommendationRequests()[i].GetText()
		r.Requested = ToProfile(ctx, recommendations.GetRecommendationRequests()[i].GetRequested())
		r.Requestor = ToProfile(ctx, recommendations.GetRecommendationRequests()[i].GetRequestor())
		r.Created_at = time.Unix(recommendations.GetRecommendationRequests()[i].GetCreatedAt(), 0)
		r.Title = recommendations.GetRecommendationRequests()[i].GetTitle()
		r.Relation = relationRPCToString(recommendations.GetRecommendationRequests()[i].GetRelation())

		rec = append(rec, RecommendationRequestResolver{R: &r})
	}

	return &rec, nil
}

// RequestedRecommendationRequests ....
func (r ProfileResolver) RequestedRecommendationRequests(ctx context.Context, args PaginationInput) (*[]RecommendationRequestResolver, error) {
	var first uint32 = 2
	var after = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	recommendations, err := user.GetRequestedRecommendationRequests(ctx, &userRPC.IDWithPagination{
		ID: string(r.ID()),
		Pagination: &userRPC.Pagination{
			First: first,
			After: after,
		},
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	rec := make([]RecommendationRequestResolver, 0, len(recommendations.GetRecommendationRequests()))

	for i := range recommendations.GetRecommendationRequests() {
		var r RecommendationRequest

		r.ID = recommendations.GetRecommendationRequests()[i].GetID()
		r.Text = recommendations.GetRecommendationRequests()[i].GetText()
		r.Requested = ToProfile(ctx, recommendations.GetRecommendationRequests()[i].GetRequested())
		r.Requestor = ToProfile(ctx, recommendations.GetRecommendationRequests()[i].GetRequestor())
		r.Created_at = time.Unix(recommendations.GetRecommendationRequests()[i].GetCreatedAt(), 0)
		r.Title = recommendations.GetRecommendationRequests()[i].GetTitle()
		r.Relation = relationRPCToString(recommendations.GetRecommendationRequests()[i].GetRelation())

		rec = append(rec, RecommendationRequestResolver{R: &r})
	}

	return &rec, nil
}

func (r ProfileResolver) HiddenRecommendation(ctx context.Context, args PaginationInput) (*[]RecommendationResolver, error) {
	var first uint32 = 2
	var after = "0"

	if args.First != nil {
		first = uint32(*args.First)
	}
	if args.After != nil {
		after = *args.After
	}

	recommendations, err := user.GetHiddenRecommendations(ctx, &userRPC.IDWithPagination{
		ID: string(r.ID()),
		Pagination: &userRPC.Pagination{
			First: first,
			After: after,
		},
		Language: r.language,
	})
	if err != nil {
		return nil, err
	}

	rec := make([]RecommendationResolver, 0, len(recommendations.GetRecommendations()))

	for i := range recommendations.GetRecommendations() {
		var r Recommendation

		r.ID = recommendations.GetRecommendations()[i].GetID()
		r.Text = recommendations.GetRecommendations()[i].GetText()
		r.Receiver = ToProfile(ctx, recommendations.GetRecommendations()[i].GetReceiver())
		r.Recommendator = ToProfile(ctx, recommendations.GetRecommendations()[i].GetRecommendator())
		r.Title = recommendations.GetRecommendations()[i].GetTitle()
		r.Relation = relationRPCToString(recommendations.GetRecommendations()[i].GetRelation())

		if recommendations.GetRecommendations()[i].GetIsIsHiddenNull() {
			r.Is_hidden = nil
		} else {
			b := recommendations.GetRecommendations()[i].GetIsHidden()
			r.Is_hidden = &b
		}

		rec = append(rec, RecommendationResolver{R: &r})
	}

	return &rec, nil
}

// func userRPCContentTypeEnumToString(data userRPC.ContentTypeEnum) string {
// 	content := "Photo"

// 	switch data {
// 	case userRPC.ContentTypeEnum_Content_Type_Article:
// 		content = "Article"
// 	case userRPC.ContentTypeEnum_Content_Type_Video:
// 		content = "Video"
// 	case userRPC.ContentTypeEnum_Content_Type_Audio:
// 		content = "Audio"
// 	}

// 	return content
// }

func stringToUserRPCContentType(data string) userRPC.ContentTypeEnum {
	var content = userRPC.ContentTypeEnum_Content_Type_Photo

	switch data {
	case "Article":
		content = userRPC.ContentTypeEnum_Content_Type_Article
	case "Video":
		content = userRPC.ContentTypeEnum_Content_Type_Video
	case "Audio":
		content = userRPC.ContentTypeEnum_Content_Type_Audio
	}

	return content
}

func stringToUserRPCToolTechnologyLevel(data string) userRPC.ToolTechnology_Level {
	var level userRPC.ToolTechnology_Level

	switch data {
	case "Level_Begginer":
		level = userRPC.ToolTechnology_Level_Beginner
	case "Level_Intermediate":
		level = userRPC.ToolTechnology_Level_Intermediate
	case "Level_Advanced":
		level = userRPC.ToolTechnology_Level_Advanced
	case "Level_Master":
		level = userRPC.ToolTechnology_Level_Master
	}

	return level
}

func userRPCToolsLevelToToolsLevel(data userRPC.ToolTechnology_Level) string {
	content := "Level_Begginer"

	switch data {
	case userRPC.ToolTechnology_Level_Intermediate:
		content = "Level_Intermediate"
	case userRPC.ToolTechnology_Level_Advanced:
		content = "Level_Advanced"
	case userRPC.ToolTechnology_Level_Master:
		content = "Level_Master"
	}

	return content
}

func passContext(ctx *context.Context) {

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		log.Println("error while passing context")
	}
}
