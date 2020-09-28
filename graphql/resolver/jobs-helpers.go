package resolver

import (
	"log"
	"strings"

	"github.com/pkg/errors"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	hc_errors "gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
	"golang.org/x/net/context"
)

func handleError(err error) (error, bool) {
	if err != nil {
		log.Println(err)
		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		if !b {
			return err, true
		}
		return errors.New(e.Description), true
	}
	return nil, false
}

func jobs_candidateProfileToGql(ctx context.Context, arg interface{}) *CandidateProfile {
	switch profile := arg.(type) {
	case *jobsRPC.CandidateProfile:
		return &CandidateProfile{
			Is_open:          profile.IsOpen,
			Career_interests: *jobs_careerInterestsToGql(profile.CareerInterests),
		}
	case *jobsRPC.CandidateViewForCompany:
		return &CandidateProfile{
			User:             getUserProfile(ctx, profile.UserId),
			User_id:          profile.UserId,
			Career_interests: *jobs_careerInterestsToGql(profile.CareerInterests),
			Is_saved:         profile.IsSaved,
		}
	default:
		return nil
	}
}

func stringToHighlightEnum(data *string) jobsRPC.JobHighlight {
	if data == nil {
		return jobsRPC.JobHighlight_None
	}
	var highlight = *data

	switch highlight {
	case "blue":
		return jobsRPC.JobHighlight_Blue
	case "white":
		return jobsRPC.JobHighlight_White
	}

	return jobsRPC.JobHighlight_None
}

func AddtionCompensationToRPC(data []string) []jobsRPC.AdditionalCompensation {

	comps := make([]jobsRPC.AdditionalCompensation, 0, len(data))

	for _, comp := range data {
		comps = append(comps, stringTpAddtionCompensationEnum(comp))
	}

	return comps

}

func stringTpAddtionCompensationEnum(data string) jobsRPC.AdditionalCompensation {

	switch data {
	case "sales_commission":
		return jobsRPC.AdditionalCompensation_Sales_Commission
	case "tips_gratuities":
		return jobsRPC.AdditionalCompensation_Tips_Gratuities

	case "profit_sharing":
		return jobsRPC.AdditionalCompensation_Profit_Sharing
	case "bonus":
		return jobsRPC.AdditionalCompensation_Bonus
	}

	return jobsRPC.AdditionalCompensation_Any

}

func getUserProfile(ctx context.Context, userId string) Profile {
	passContext(&ctx)

	var gqlProfile Profile
	userProfile, err := user.GetProfileByID(ctx, &userRPC.ID{ID: userId})
	if err != nil {
		log.Println("error: getUserProfile", err)
	}
	// if err == nil && len(userProfile.Profiles) > 0 {
	// 	gqlProfile = ToProfile(userProfile.Profiles[0])
	// }
	gqlProfile = ToProfile(ctx, userProfile)
	return gqlProfile
}

func jobs_careerInterestsToGql(ci *jobsRPC.CareerInterests) *CareerInterests {
	if ci == nil {
		return &CareerInterests{
			Company_size:    jobCompanyRPCToString(jobsRPC.CompanySize_SIZE_UNDEFINED),
			Salary_interval: "Any",
			Experience:      jobExperienceEnumToString(jobsRPC.ExperienceEnum_UnknownExperience),
		}
		// return nil
	}

	jobTypes := make([]string, len(ci.JobTypes))
	for i, t := range ci.JobTypes {
		jobTypes[i] = t.String()
	}

	c := CareerInterests{
		Jobs:            ci.Jobs,
		Industry:        ci.Industry,
		Subindustry:     ci.Subindustry,
		Company_size:    jobCompanyRPCToString(ci.CompanySize),
		Job_types:       jobTypes,
		Salary_currency: ci.SalaryCurrency,
		Salary_min:      ci.SalaryMin,
		Salary_max:      ci.SalaryMax,
		Salary_interval: jobsRPCSalaryIntervalToString(ci.SalaryInterval),
		Relocate:        ci.Relocate,
		Remote:          ci.Remote,
		Travel:          ci.Travel,
		Experience:      jobExperienceEnumToString(ci.Experience),
		Cities:          make([]City, 0, len(ci.GetLocations())),
		Suitable_for:    suitableForRPCToArray(ci.GetSuitableFor()),
		Is_invited:      ci.IsInvited,
		Is_saved:        ci.IsSaved,
	}

	for _, loc := range ci.GetLocations() {
		c.Cities = append(c.Cities, City{
			City:        loc.GetCityName(),
			Country:     loc.GetCountry(),
			ID:          loc.GetCityID(),
			Subdivision: loc.GetSubdivision(),
		})

	}

	return &c
}

func jobCompanyRPCToString(data jobsRPC.CompanySize) string {
	switch data {
	case jobsRPC.CompanySize_SIZE_10001_PLUS_EMPLOYEES:
		return "size_10001_plus_employees"
	case jobsRPC.CompanySize_SIZE_1_10_EMPLOYEES:
		return "size_1_10_employees"
	case jobsRPC.CompanySize_SIZE_11_50_EMPLOYEES:
		return "size_11_50_employees"
	case jobsRPC.CompanySize_SIZE_51_200_EMPLOYEES:
		return "size_51_200_employees"
	case jobsRPC.CompanySize_SIZE_201_500_EMPLOYEES:
		return "size_201_500_employees"
	case jobsRPC.CompanySize_SIZE_501_1000_EMPLOYEES:
		return "size_501_1000_employees"
	case jobsRPC.CompanySize_SIZE_1001_5000_EMPLOYEES:
		return "size_1001_5000_employees"
	case jobsRPC.CompanySize_SIZE_5001_10000_EMPLOYEES:
		return "size_5001_10000_employees"

	}

	return "size_unknown"
}

func jobs_careerInterestsInputToRPC(ci *CareerInterestsInput) *jobsRPC.CareerInterests {
	jobTypes := make([]jobsRPC.JobType, len(ci.Job_types))
	for i, t := range ci.Job_types {
		jobTypes[i] = jobsRPC.JobType(jobsRPC.JobType_value[t])
	}
	ce := jobsRPC.CareerInterests{
		Jobs:           ci.Jobs,
		Industry:       ci.Industry,
		Subindustry:    ci.Subindustry,
		CompanySize:    jobsRPC.CompanySize(jobsRPC.CompanySize_value[strings.ToUpper(ci.Company_size)]),
		JobTypes:       jobTypes,
		SalaryCurrency: ci.Salary_currency,
		SalaryMin:      ci.Salary_min,
		SalaryMax:      ci.Salary_max,
		SalaryInterval: jobsRPC.SalaryInterval(jobsRPC.SalaryInterval_value[ci.Salary_interval]),
		Relocate:       ci.Relocate,
		Remote:         ci.Remote,
		Travel:         ci.Travel,
		Experience:     stringToCareerInterestExperienceEnumRPC(ci.Experience),
		Locations:      make([]*jobsRPC.Location, 0, len(ci.Cities)),
		SuitableFor:    suitableForToRPCArray(ci.Suitable_for),
	}

	for _, c := range ci.Cities {
		ce.Locations = append(ce.Locations, &jobsRPC.Location{
			CityID: c,
		})
	}

	return &ce
}

func stringToCareerInterestExperienceEnumRPC(s string) jobsRPC.ExperienceEnum {
	switch s {
	case "without_experience":
		return jobsRPC.ExperienceEnum_WithoutExperience
	case "less_then_one_year":
		return jobsRPC.ExperienceEnum_LessThenOneYear
	case "one_two_years":
		return jobsRPC.ExperienceEnum_OneTwoYears
	case "two_three_years":
		return jobsRPC.ExperienceEnum_TwoThreeYears
	case "three_five_years":
		return jobsRPC.ExperienceEnum_ThreeFiveYears
	case "five_seven_years":
		return jobsRPC.ExperienceEnum_FiveSevenyears
	case "seven_ten_years":
		return jobsRPC.ExperienceEnum_SevenTenYears
	case "ten_years_and_more":
		return jobsRPC.ExperienceEnum_TenYearsAndMore
	}

	return jobsRPC.ExperienceEnum_UnknownExperience
}

// TODO: split it
func jobs_jobPostingToResolver(ctx context.Context, arg interface{}) *JobPostingResolver {
	switch job := arg.(type) {
	case *jobsRPC.JobViewForUser:
		comp, err := company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
			Ids: []string{job.GetCompanyInfo().GetCompanyId()},
		})

		if err != nil {
			log.Println(err)
		}

		post := JobPostingResolver{
			R: &JobPosting{
				ID: job.Id,

				// Company_id:  job.CompanyInfo.CompanyId,
				Job_details:     *jobs_jobDetailsToGql(job.JobDetails),
				Job_meta:        *jobs_jobMetaToGql(job.Metadata),
				Application:     jobs_applicationToGql(job.Application),
				Is_saved:        job.GetIsSaved(),
				Is_applied:      job.GetIsApplied(),
				Text_invitation: job.GetInvitationText(),
			},
		}

		if len(comp.GetProfiles()) > 0 {
			post.R.Company = toCompanyProfile(ctx, *(comp.GetProfiles()[0]))
		}

		return &post
	case *jobsRPC.JobViewForCompany:
		return &JobPostingResolver{
			R: &JobPosting{
				ID:                     job.Id,
				User_id:                job.GetUserId(),
				Job_details:            *jobs_jobDetailsToGql(job.JobDetails),
				Status:                 job.Status,
				Job_meta:               *jobs_jobMetaToGql(job.Metadata),
				Number_of_views:        job.NumberOfViews,
				Number_of_applications: job.NumberOfApplications,
			},
		}
	default:
		return nil
	}
}

func jobs_jobMetaToGql(meta *jobsRPC.JobMeta) *JobMeta {
	return &JobMeta{
		Anonymous:               meta.Anonymous,
		Num_of_languages:        meta.GetNumOfLanguages(),
		Highlight:               highlightRPCToString(meta.GetHighlight()),
		Renewal:                 meta.Renewal,
		Currency:                meta.GetCurrency(),
		Advertisement_countries: meta.AdvertisementCountries,
		// Job_plan:                jobPlanRPCToString(meta.JobPlan), // TODO:
		Amount_of_days: meta.GetAmountOfDays(),
	}
}

func highlightRPCToString(data jobsRPC.JobHighlight) string {
	switch data {
	case jobsRPC.JobHighlight_Blue:
		return "blue"
	case jobsRPC.JobHighlight_White:
		return "white"
	}
	return "none"
}
func jobPlanRPCToString(data jobsRPC.JobPlan) string {
	switch data {
	case jobsRPC.JobPlan_Basic:
		return "basic"
	case jobsRPC.JobPlan_Start:
		return "start"
	case jobsRPC.JobPlan_Professional:
		return "professional"
	case jobsRPC.JobPlan_Standard:
		return "standard"
	case jobsRPC.JobPlan_ProfessionalPlus:
		return "professionalPlus"
	case jobsRPC.JobPlan_Exclusive:
		return "exclusive"
	case jobsRPC.JobPlan_Premium:
		return "premium"
	}

	return "unknown"
}

func jobs_jobApplicantToResolver(ctx context.Context, applicant *jobsRPC.JobApplicant) *JobApplicantResolver {
	app := JobApplicant{
		UserId:      applicant.UserId,
		User:        getUserProfile(ctx, applicant.UserId),
		Application: *jobs_applicationToGql(applicant.Application),
		// Career_interests: *jobs_careerInterestsToGql(applicant.CareerInterests),
	}
	if ca := jobs_careerInterestsToGql(applicant.CareerInterests); ca != nil {
		app.Career_interests = ca
	}
	return &JobApplicantResolver{
		R: &app,
	}
}

func jobs_applicationToGql(app *jobsRPC.Application) *Application {
	if app == nil {
		return &Application{}
	}
	a := Application{
		Job_id:       app.JobId,
		Email:        app.Email,
		Phone:        app.Phone,
		Cover_letter: app.CoverLetter,
		// Documents:    app.Documents,
		Documents:  make([]File, 0, len(app.Documents)),
		Created_at: app.GetCreatedAt(),
		Metadata: ApplicationMeta{
			Category: jobsRPCApplicantCategoryEnumToString(app.Metadata.Category),
			Seen:     app.Metadata.Seen,
		},
	}

	for _, f := range app.Documents {
		if f != nil {
			a.Documents = append(a.Documents, *jobfileRPCToFile(f))
		}
	}

	return &a
}

func jobs_jobDetailsToGql(job *jobsRPC.JobDetails) *JobDetails {
	if job == nil {
		return nil
	}
	jobTypes := make([]string, len(job.EmploymentTypes))
	for i, t := range job.EmploymentTypes {
		jobTypes[i] = t.String()
	}
	descriptions := make([]JobDescription, len(job.Descriptions))
	for i, d := range job.Descriptions {
		descriptions[i] = *jobs_jobDescriptionToGql(d)
	}
	benefits := make([]string, 0, len(job.Benefits))
	for _, b := range job.Benefits {
		benefits = append(benefits, jobBenefitRPCToJobBenefit(b))
	}

	jd := JobDetails{
		Title:                   job.Title,
		Country:                 job.Country,
		Location_type:           jobLocationTypeRPCToString(job.GetLocationType()),
		Job_functions:           jobFunctionRPCToArray(job.JobFunctions),
		Employment_types:        jobTypes,
		Descriptions:            descriptions,
		Salary_currency:         job.SalaryCurrency,
		Salary_min:              job.SalaryMin,
		Salary_max:              job.SalaryMax,
		Salary_interval:         job.SalaryInterval.String(),
		Benefits:                benefits,
		Number_of_positions:     job.NumberOfPositions,
		Publish_day:             job.PublishDay,
		Publish_month:           job.PublishMonth,
		Publish_year:            job.PublishYear,
		Deadline_day:            job.DeadlineDay,
		Deadline_month:          job.DeadlineMonth,
		Deadline_year:           job.DeadlineYear,
		Hiring_day:              job.HiringDay,
		Hiring_month:            job.HiringMonth,
		Hiring_year:             job.HiringYear,
		Cover_letter:            job.CoverLetter,
		Header_url:              job.HeaderUrl,
		Additional_compensation: additionalCompensationRPCToArray(job.GetAdditionalCompensation()),
		Additional_info:         additionalInfoRPCToAdditionalInfo(job.GetAdditionalInfo()),
		Required:                jobQuailificationRPCToQuialification(job.GetRequired()),
		Preterred:               jobQuailificationRPCToQuialification(job.GetPreterred()),
		Files:                   filesRPCToFile(job.GetFiles()),
	}

	if job.GetLocation() != nil {
		jd.Location = City{
			City:        job.GetLocation().GetCityName(),
			Country:     job.GetLocation().GetCountry(),
			ID:          job.GetLocation().GetCityID(),
			Subdivision: job.GetLocation().GetSubdivision(),
		}
	}

	return &jd
}

func jobFunctionRPCToArray(data []jobsRPC.JobFunction) []string {

	jbFns := make([]string, 0, len(data))

	for _, jbFn := range data {
		jbFns = append(jbFns, jobFunctionRPCToEnum(jbFn))
	}

	return jbFns
}

func jobFunctionRPCToEnum(data jobsRPC.JobFunction) string {

	log.Printf("Job Functions %+v \n", data)
	switch data {
	case jobsRPC.JobFunction_Accounting:
		return "accounting"
	case jobsRPC.JobFunction_Administrative:
		return "administrative"
	case jobsRPC.JobFunction_Arts_Design:
		return "arts_design"
	case jobsRPC.JobFunction_Business_Development:
		return "business_development"
	case jobsRPC.JobFunction_Community_Social_Services:
		return "community_social_services"
	case jobsRPC.JobFunction_Consulting:
		return "consulting"
	case jobsRPC.JobFunction_Education:
		return "education"
	case jobsRPC.JobFunction_Engineering:
		return "engineering"
	case jobsRPC.JobFunction_Entrepreneurship:
		return "entrepreneurship"
	case jobsRPC.JobFunction_Finance:
		return "finance"
	case jobsRPC.JobFunction_Healthcare_Services:
		return "healthcare_services"
	case jobsRPC.JobFunction_Human_Resources:
		return "human_resources"
	case jobsRPC.JobFunction_Information_Technology:
		return "information_technology"
	case jobsRPC.JobFunction_Legal:
		return "legal"
	case jobsRPC.JobFunction_Marketing:
		return "marketing"
	case jobsRPC.JobFunction_Media_Communications:
		return "media_communications"
	case jobsRPC.JobFunction_Military_Protective_Services:
		return "military_protective_services"
	case jobsRPC.JobFunction_Operations:
		return "operations"
	case jobsRPC.JobFunction_Product_Management:
		return "product_management"
	case jobsRPC.JobFunction_Program_Product_Management:
		return "program_product_management"
	case jobsRPC.JobFunction_Purchasing:
		return "purchasing"
	case jobsRPC.JobFunction_Quality_Assurance:
		return "quality_assurance"
	case jobsRPC.JobFunction_Real_Estate:
		return "real_estate"
	case jobsRPC.JobFunction_Rersearch:
		return "rersearch"
	case jobsRPC.JobFunction_Sales:
		return "sales"
	case jobsRPC.JobFunction_Support:
		return "support"

	}

	return "none"
}

func jobFunctionArrayToRPC(data []string) []jobsRPC.JobFunction {

	jbFns := make([]jobsRPC.JobFunction, 0, len(data))

	for _, jbFn := range data {
		jbFns = append(jbFns, jobFunctionEnumToRPC(jbFn))
	}

	return jbFns

}

func jobFunctionEnumToRPC(data string) jobsRPC.JobFunction {

	switch data {
	case "accounting":
		return jobsRPC.JobFunction_Accounting
	case "administrative":
		return jobsRPC.JobFunction_Administrative
	case "arts_design":
		return jobsRPC.JobFunction_Arts_Design
	case "business_development":
		return jobsRPC.JobFunction_Business_Development
	case "community_social_services":
		return jobsRPC.JobFunction_Community_Social_Services
	case "consulting":
		return jobsRPC.JobFunction_Consulting
	case "education":
		return jobsRPC.JobFunction_Education
	case "engineering":
		return jobsRPC.JobFunction_Engineering
	case "entrepreneurship":
		return jobsRPC.JobFunction_Entrepreneurship
	case "finance":
		return jobsRPC.JobFunction_Finance
	case "healthcare_services":
		return jobsRPC.JobFunction_Healthcare_Services
	case "human_resources":
		return jobsRPC.JobFunction_Human_Resources
	case "information_technology":
		return jobsRPC.JobFunction_Information_Technology
	case "legal":
		return jobsRPC.JobFunction_Legal
	case "marketing":
		return jobsRPC.JobFunction_Marketing
	case "media_communications":
		return jobsRPC.JobFunction_Media_Communications
	case "military_protective_services":
		return jobsRPC.JobFunction_Military_Protective_Services
	case "operations":
		return jobsRPC.JobFunction_Operations
	case "product_management":
		return jobsRPC.JobFunction_Product_Management
	case "program_product_management":
		return jobsRPC.JobFunction_Program_Product_Management
	case "purchasing":
		return jobsRPC.JobFunction_Purchasing
	case "quality_assurance":
		return jobsRPC.JobFunction_Quality_Assurance
	case "real_estate":
		return jobsRPC.JobFunction_Real_Estate
	case "rersearch":
		return jobsRPC.JobFunction_Rersearch
	case "sales":
		return jobsRPC.JobFunction_Sales
	case "support":
		return jobsRPC.JobFunction_Support

	}

	return jobsRPC.JobFunction_None_Job_Func
}

func filesRPCToFile(data []*jobsRPC.File) []File {
	if data == nil {
		return nil
	}

	files := make([]File, 0, len(data))

	for _, file := range data {
		files = append(files, fileRPCToFile(file))
	}

	return files
}

func fileRPCToFile(data *jobsRPC.File) File {
	if data == nil {
		return File{}
	}

	return File{
		ID:        data.GetID(),
		Name:      data.GetName(),
		Mime_type: data.GetMimeType(),
		Address:   data.GetURL(),
	}
}

func jobQuailificationRPCToQuialification(data *jobsRPC.ApplicationQuailification) *ApplicantQuailification {
	if data == nil {
		return &ApplicantQuailification{}
	}

	return &ApplicantQuailification{
		Experience: jobExperienceEnumToString(data.GetExperience()),
		Skills:     data.GetSkills(),
		Education:  data.GetEducations(),
		Languages:  jobLanguageRPCArrayToLanguageArray(data.GetLanguages()),
		Tools:      jobToolsRPCArrayToToolsArray(data.GetTools()),
		License:    data.GetLicense(),
		Work:       data.GetWork(),
	}

}

func jobToolsRPCArrayToToolsArray(data []*jobsRPC.ApplcantToolsAndTechnology) []ToolsTechnologies {
	if data == nil {
		return nil
	}

	tools := make([]ToolsTechnologies, 0, len(data))

	for _, tool := range data {
		tools = append(tools, jobToolRPCToTool(tool))
	}

	return tools
}

func jobToolRPCToTool(data *jobsRPC.ApplcantToolsAndTechnology) ToolsTechnologies {
	if data == nil {
		return ToolsTechnologies{}
	}

	return ToolsTechnologies{
		ID:   data.GetID(),
		Rank: jobApplicationLevelRPCTostring(data.GetRank()),
		Tool: data.GetTool(),
	}
}

func jobLanguageRPCArrayToLanguageArray(data []*jobsRPC.ApplicantLanguage) []Language {
	if data == nil {
		return nil
	}

	languages := make([]Language, 0, len(data))

	for _, language := range data {
		languages = append(languages, jobLanguageRPCToLanguage(language))
	}

	return languages
}

func jobLanguageRPCToLanguage(data *jobsRPC.ApplicantLanguage) Language {
	if data == nil {
		return Language{}
	}

	return Language{
		ID:       data.GetID(),
		Language: data.GetLanguage(),
		Rank:     jobApplicationLevelRPCTostring(data.GetRank()),
	}
}

func qualificationInputToRPC(data *ApplicantQuailificationInput) *jobsRPC.ApplicationQuailification {
	if data == nil {
		return &jobsRPC.ApplicationQuailification{}
	}

	return &jobsRPC.ApplicationQuailification{
		Experience: stringToJobExperienceEnum(data.Experience),
		Educations: NullStringArrayToStringArray(data.Education),
		Skills:     NullStringArrayToStringArray(data.Skills),
		Tools:      toolsArrayToRPC(data.Tools),
		Languages:  languagesArrayToRPC(data.Languages),
		License:    data.License,
		Work:       data.Work,
	}
}

func languagesArrayToRPC(data *[]LanguageInput) []*jobsRPC.ApplicantLanguage {
	if data == nil {
		return nil
	}

	languages := make([]*jobsRPC.ApplicantLanguage, 0, len(*data))

	for _, language := range *data {
		languages = append(languages, quialifationLanguageToRPC(language))
	}

	return languages
}

func quialifationLanguageToRPC(data LanguageInput) *jobsRPC.ApplicantLanguage {
	return &jobsRPC.ApplicantLanguage{
		Language: data.Language,
		Rank:     stringToApplicationLevel(data.Rank),
	}
}

func toolsArrayToRPC(data *[]ToolsTechnologiesInput) []*jobsRPC.ApplcantToolsAndTechnology {
	if data == nil {
		return nil
	}

	tools := make([]*jobsRPC.ApplcantToolsAndTechnology, 0, len(*data))

	for _, tool := range *data {
		tools = append(tools, quialifationToolToRPC(tool))
	}

	return tools
}

func quialifationToolToRPC(data ToolsTechnologiesInput) *jobsRPC.ApplcantToolsAndTechnology {

	return &jobsRPC.ApplcantToolsAndTechnology{
		Tool: data.Tool,
		Rank: stringToApplicationLevel(data.Rank),
	}
}

func jobApplicationLevelRPCTostring(data jobsRPC.ApplicationLevel) string {

	switch data {
	case jobsRPC.ApplicationLevel_Level_Intermediate:
		return "Level_Intermediate"
	case jobsRPC.ApplicationLevel_Level_Advanced:
		return "Level_Advanced"
	case jobsRPC.ApplicationLevel_Level_Master:
		return "Level_Master"
	}

	return "Level_Begginer"
}

func stringToApplicationLevel(data string) jobsRPC.ApplicationLevel {
	switch data {
	case "Level_Begginer":
		return jobsRPC.ApplicationLevel_Level_Begginer
	case "Level_Intermediate":
		return jobsRPC.ApplicationLevel_Level_Intermediate
	case "Level_Advanced":
		return jobsRPC.ApplicationLevel_Level_Advanced
	case "Level_Master":
		return jobsRPC.ApplicationLevel_Level_Master
	}

	return jobsRPC.ApplicationLevel_Level_Unknown
}

func additionalInfoRPCToAdditionalInfo(data *jobsRPC.AdditionalInfo) *AdditionalInfo {
	if data == nil {
		return &AdditionalInfo{}
	}

	return &AdditionalInfo{
		Suitable_for:       suitableForRPCToArray(data.GetSuitableFor()),
		Travel_requirement: travelRequirementRPCToEnum(data.GetTravelRequirement()),
	}

}

func travelRequirementRPCToEnum(data jobsRPC.TravelRequirement) string {

	switch data {
	case jobsRPC.TravelRequirement_All_time:
		return "all_time"
	case jobsRPC.TravelRequirement_Few_times:
		return "few_times"
	case jobsRPC.TravelRequirement_Once_month:
		return "once_month"
	case jobsRPC.TravelRequirement_Once_week:
		return "once_week"
	case jobsRPC.TravelRequirement_Once_year:
		return "once_year"
	}

	return "none"
}

func suitableForRPCToArray(data []jobsRPC.SuitableFor) []string {

	suitables := make([]string, 0, len(data))

	for _, suitable := range data {
		suitables = append(suitables, suitableForRPCToSuitableFor(suitable))
	}

	return suitables
}

func suitableForRPCToSuitableFor(data jobsRPC.SuitableFor) string {
	switch data {
	case jobsRPC.SuitableFor_Student:
		return "student"
	case jobsRPC.SuitableFor_Person_With_Disability:
		return "person_with_a_disability"
	case jobsRPC.SuitableFor_Single_Parent:
		return "single_parent"
	case jobsRPC.SuitableFor_Veterans:
		return "veterans"
	}
	return "none"
}

func additionalInfoToRPC(data *AdditionalInfoInput) *jobsRPC.AdditionalInfo {
	if data == nil {
		return nil
	}

	return &jobsRPC.AdditionalInfo{
		SuitableFor:       suitableForToRPCArray(data.Suitable_for),
		TravelRequirement: travelRequirementForRPC(data.Travel_requirement),
	}

}

func travelRequirementForRPC(data string) jobsRPC.TravelRequirement {

	switch data {
	case "all_time":
		return jobsRPC.TravelRequirement_All_time
	case "once_week":
		return jobsRPC.TravelRequirement_Once_week
	case "once_month":
		return jobsRPC.TravelRequirement_Once_month
	case "few_times":
		return jobsRPC.TravelRequirement_Few_times
	case "once_year":
		return jobsRPC.TravelRequirement_Once_year
	}

	return jobsRPC.TravelRequirement_Travel_req_none
}

func suitableForToRPCArray(data []string) []jobsRPC.SuitableFor {

	suitables := make([]jobsRPC.SuitableFor, 0, len(data))

	for _, suitable := range data {
		suitables = append(suitables, suitableForToRPC(suitable))
	}

	return suitables
}

func suitableForToRPC(data string) jobsRPC.SuitableFor {

	switch data {
	case "student":
		return jobsRPC.SuitableFor_Student
	case "person_with_a_disability":
		return jobsRPC.SuitableFor_Person_With_Disability
	case "single_parent":
		return jobsRPC.SuitableFor_Single_Parent
	case "veterans":
		return jobsRPC.SuitableFor_Veterans
	}

	return jobsRPC.SuitableFor_None_Suitable
}

func jobLocationTypeRPCToString(data jobsRPC.LocationType) string {
	if data == jobsRPC.LocationType_On_Site {
		return "On_Site_Work"
	}
	return "Remote_only"
}

func stringLocationTypeToEnum(data string) jobsRPC.LocationType {

	if data == "Remote_only" {
		return jobsRPC.LocationType_Remote
	}

	return jobsRPC.LocationType_On_Site
}

func additionalCompensationRPCToArray(data []jobsRPC.AdditionalCompensation) []string {

	comps := make([]string, 0, len(data))

	for i := range data {
		comps = append(comps, additionalCompensationEnumToString(data[i]))
	}

	return comps
}

func additionalCompensationEnumToString(data jobsRPC.AdditionalCompensation) string {

	switch data {
	case jobsRPC.AdditionalCompensation_Profit_Sharing:
		return "profit_sharing"
	case jobsRPC.AdditionalCompensation_Sales_Commission:
		return "sales_commission"
	case jobsRPC.AdditionalCompensation_Tips_Gratuities:
		return "tips_gratuities"
	}

	return "bonus"
}
func jobs_jobDescriptionToGql(desc *jobsRPC.JobDescription) *JobDescription {
	return &JobDescription{
		Language:    desc.Language,
		Description: desc.Description,
		Why_us:      desc.WhyUs,
	}
}
func jobsRPCDescriptionToJBDSCR(desc *jobsRPC.JobDescription) JobDescription {
	return JobDescription{
		Language:    desc.Language,
		Description: desc.Description,
		Why_us:      desc.WhyUs,
	}
}

// func jobs_jobAlertToResolver(alert *jobsRPC.JobAlert) *JobAlertResolver {
// 	return &JobAlertResolver{
// 		R: &JobAlert{
// 			ID:                  alert.Id,
// 			Name:                alert.Name,
// 			Interval:            alert.Interval,
// 			Notify_email:        alert.NotifyEmail,
// 			Notify_notification: alert.NotifyNotification,
// 			Filter:              *jobs_jobSearchFilterToGql(alert.Filter),
// 		},
// 	}
// }

// func jobs_candidateAlertToResolver(alert *jobsRPC.CandidateAlert) *CandidateAlertResolver {
// 	return &CandidateAlertResolver{
// 		R: &CandidateAlert{
// 			ID:                  alert.Id,
// 			Name:                alert.Name,
// 			Interval:            alert.Interval,
// 			Notify_email:        alert.NotifyEmail,
// 			Notify_notification: alert.NotifyNotification,
// 			Filter:              *jobs_candidateSearchFilterToGql(alert.Filter),
// 		},
// 	}
// }

// func jobs_NamedJobSearchFilterToResolver(filter *jobsRPC.NamedJobSearchFilter) *NamedJobSearchFilterResolver {
// 	return &NamedJobSearchFilterResolver{
// 		R: &NamedJobSearchFilter{
// 			ID:     filter.Id,
// 			Name:   filter.Name,
// 			Filter: *jobs_jobSearchFilterToGql(filter.Filter),
// 		},
// 	}
// }

// func jobs_NamedCandidateSearchFilterToResolver(filter *jobsRPC.NamedCandidateSearchFilter) *NamedCandidateSearchFilterResolver {
// 	return &NamedCandidateSearchFilterResolver{
// 		R: &NamedCandidateSearchFilter{
// 			ID:     filter.Id,
// 			Name:   filter.Name,
// 			Filter: *jobs_candidateSearchFilterToGql(filter.Filter),
// 		},
// 	}
// }

// func jobs_jobSearchFilterToGql(filter *jobsRPC.JobSearchFilter) *JobSearchFilter {
// 	dates := make([]JobsDate, len(filter.DatePosted))
// 	for i, d := range filter.DatePosted {
// 		dates[i] = JobsDate{
// 			Day:   int32(d.Day),
// 			Month: int32(d.Month),
// 			Year:  int32(d.Year),
// 		}
// 	}
//
// 	sizes := make([]int32, len(filter.CompanySize))
// 	for i, s := range filter.CompanySize {
// 		sizes[i] = int32(s)
// 	}
// 	return &JobSearchFilter{
// 		Keyword:     filter.Keyword,
// 		Date_posted: dates,
//
// 		Experience_level: int32(filter.ExperienceLevel),
//
// 		Degree: filter.Degree,
//
// 		Country: filter.Country,
// 		City:    filter.City,
//
// 		Job_type: filter.JobType,
//
// 		Language: filter.Language,
//
// 		Industry:     filter.Industry,
// 		Subindustry:  filter.Subindustry,
// 		Company_name: filter.CompanyName,
// 		Company_size: sizes,
//
// 		Currency:   filter.Currency,
// 		Period:     filter.Period,
// 		Min_salary: int32(filter.MinSalary),
// 		Max_salary: int32(filter.MaxSalary),
//
// 		Skill: filter.Skill,
//
// 		Is_following:         filter.IsFollowing,
// 		Without_cover_letter: filter.WithoutCoverLetter,
// 		With_salary:          filter.WithSalary,
// 	}
// }

// func jobs_jobSearchFilterInputToRPC(filter *JobSearchFilterInput) *jobsRPC.JobSearchFilter {
// 	dates := make([]*jobsRPC.Date, len(filter.Date_posted))
// 	for i, d := range filter.Date_posted {
// 		dates[i] = &jobsRPC.Date{
// 			Day:   uint32(d.Day),
// 			Month: uint32(d.Month),
// 			Year:  uint32(d.Year),
// 		}
// 	}
//
// 	sizes := make([]uint32, len(filter.Company_size))
// 	for i, s := range filter.Company_size {
// 		sizes[i] = uint32(s)
// 	}
// 	f := &jobsRPC.JobSearchFilter{
// 		Keyword:    filter.Keyword,
// 		DatePosted: dates,
//
// 		Degree: filter.Degree,
//
// 		Country: filter.Country,
// 		City:    filter.City,
//
// 		JobType: filter.Job_type,
//
// 		Language: filter.Language,
//
// 		Industry:    filter.Industry,
// 		Subindustry: filter.Subindustry,
// 		CompanyName: filter.Company_name,
// 		CompanySize: sizes,
//
// 		Currency:  filter.Currency,
// 		Period:    filter.Period,
// 		MinSalary: uint32(filter.Min_salary),
// 		MaxSalary: uint32(filter.Max_salary),
//
// 		Skill: filter.Skill,
//
// 		IsFollowing:        filter.Is_following,
// 		WithoutCoverLetter: filter.Without_cover_letter,
// 		WithSalary:         filter.With_salary,
// 	}
// 	if filter.Experience_level != nil {
// 		f.ExperienceLevel = uint32(*filter.Experience_level)
// 		f.IsExperienceLevelNull = false
// 	} else {
// 		f.IsExperienceLevelNull = true
// 	}
//
// 	return f
// }

//func jobs_jobSearchFilterToRPC(filter *JobSearchFilter) *jobsRPC.JobSearchFilter {
//	return &jobsRPC.JobSearchFilter{
//		Query:              filter.Query,
//		Country:            filter.Country,
//		City:               filter.City,
//		Industry:           filter.Industry,
//		Subindustry:        filter.Subindustry,
//		CompanyName:        filter.Company_name,
//		CompanySize:        jobsRPC.CompanySize(jobsRPC.CompanySize_value[filter.Company_size]),
//		Experience:         filter.Experience,
//		JobType:            jobsRPC.JobType(jobsRPC.JobType_value[filter.Job_type]),
//		SalaryCurrency:     filter.Salary_currency,
//		SalaryMin:          filter.Salary_min,
//		SalaryMax:          filter.Salary_max,
//		SalaryInterval:     jobsRPC.SalaryInterval(jobsRPC.SalaryInterval_value[filter.Salary_interval]),
//		Skill:              filter.Skill,
//		Degree:             filter.Degree,
//		Language:           filter.Language,
//		FollowingCompanies: filter.Following_companies,
//		WithoutCoverLetter: filter.Without_coverLetter,
//		WithSalary:         filter.With_salary,
//	}
//}

func jobs_candidateSearchFilterInputToRPC(filter *CandidateSearchFilterInput) *jobsRPC.CandidateSearchFilter {
	experiences := make([]uint32, len(filter.Experience_level))
	for i, e := range filter.Experience_level {
		experiences[i] = uint32(e)
	}
	return &jobsRPC.CandidateSearchFilter{
		Keywords:        NullStringArrayToStringArray(filter.Keywords),
		Country:         NullStringArrayToStringArray(filter.Country),
		City:            NullStringArrayToStringArray(filter.City),
		CurrentCompany:  NullStringArrayToStringArray(filter.Current_company),
		PastCompany:     NullStringArrayToStringArray(filter.Past_company),
		Industry:        NullStringArrayToStringArray(filter.Industry),
		SubIndustry:     NullStringArrayToStringArray(filter.Sub_industry),
		ExperienceLevel: stringToCareerInterestExperienceEnumRPC(filter.Experience_level),
		JobType:         NullStringArrayToStringArray(filter.Job_type),
		Skill:           NullStringArrayToStringArray(filter.Skill),
		Language:        NullStringArrayToStringArray(filter.Language),
		School:          NullStringArrayToStringArray(filter.School),
		Degree:          NullStringArrayToStringArray(filter.Degree),
		FieldOfStudy:    NullStringArrayToStringArray(filter.Field_of_study),
		IsStudent:       filter.Is_student,
		Currency:        NullToString(filter.Currency),
		Period:          NullToString(filter.Period),
		MinSalary:       Nullint32ToUint32(filter.Min_salary),
		MaxSalary:       Nullint32ToUint32(filter.Max_salary),

		IsWillingToTravel:      filter.Is_willing_to_travel,
		IsWillingToWorkRemotly: filter.Is_willing_to_work_remotly,
		IsPossibleToRelocate:   filter.Is_possible_to_relocate,
	}
}

// func jobs_candidateSearchFilterToGql(filter *jobsRPC.CandidateSearchFilter) *CandidateSearchFilter {
// 	// experiences := make([]int32, len(filter.ExperienceLevel))
// 	// for i, e := range filter.ExperienceLevel {
// 	// 	experiences[i] = int32(e)
// 	// }
// 	return &CandidateSearchFilter{
// 		Keywords: filter.Keywords,
//
// 		Country: filter.Country,
// 		City:    filter.City,
//
// 		Current_company: filter.CurrentCompany,
// 		Past_company:    filter.PastCompany,
//
// 		Industry:     filter.Industry,
// 		Sub_industry: filter.SubIndustry,
//
// 		Experience_level: CareerInterestExperienceEnumRPC(filter.ExperienceLevel),
// 		Job_type:         filter.JobType,
//
// 		Skill:    filter.Skill,
// 		Language: filter.Language,
//
// 		School:         filter.School,
// 		Degree:         filter.Degree,
// 		Field_of_study: filter.FieldOfStudy,
// 		Is_student:     filter.IsStudent,
//
// 		Currency:   filter.Currency,
// 		Period:     filter.Period,
// 		Min_salary: int32(filter.MinSalary),
// 		Max_salary: int32(filter.MaxSalary),
//
// 		Is_willing_to_travel:       filter.IsWillingToTravel,
// 		Is_willing_to_work_remotly: filter.IsWillingToWorkRemotly,
// 		Is_possible_to_relocate:    filter.IsPossibleToRelocate,
// 	}
// }

func jobBenefitToString(s string) jobsRPC.JobDetails_JobBenefit {
	switch s {
	case "labor_agreement":
		return jobsRPC.JobDetails_labor_agreement
	case "remote_working":
		return jobsRPC.JobDetails_remote_working
	case "floater":
		return jobsRPC.JobDetails_floater
	case "paid_timeoff":
		return jobsRPC.JobDetails_paid_timeoff
	case "flexible_working_hours":
		return jobsRPC.JobDetails_flexible_working_hours
	case "additional_timeoff":
		return jobsRPC.JobDetails_additional_timeoff
	case "additional_parental_leave":
		return jobsRPC.JobDetails_additional_parental_leave
	case "sick_leave_for_family_members":
		return jobsRPC.JobDetails_sick_leave_for_family_members
	case "company_daycare":
		return jobsRPC.JobDetails_company_daycare
	case "company_canteen":
		return jobsRPC.JobDetails_company_canteen
	case "sport_facilities":
		return jobsRPC.JobDetails_sport_facilities
	case "access_for_handicapped_persons":
		return jobsRPC.JobDetails_access_for_handicapped_persons
	case "employee_parking":
		return jobsRPC.JobDetails_employee_parking
	case "shuttle_service":
		return jobsRPC.JobDetails_shuttle_service
	case "multiple_work_spaces":
		return jobsRPC.JobDetails_multiple_work_spaces
	case "corporate_events":
		return jobsRPC.JobDetails_corporate_events
	case "trainig_and_development":
		return jobsRPC.JobDetails_trainig_and_development
	case "pets_allowed":
		return jobsRPC.JobDetails_pets_allowed
	case "corporate_medical_staff":
		return jobsRPC.JobDetails_corporate_medical_staff
	case "game_consoles":
		return jobsRPC.JobDetails_game_consoles
	case "snack_and_drink_selfservice":
		return jobsRPC.JobDetails_snack_and_drink_selfservice
	case "private_pension_scheme":
		return jobsRPC.JobDetails_private_pension_scheme
	case "health_insurance":
		return jobsRPC.JobDetails_health_insurance
	case "dental_care":
		return jobsRPC.JobDetails_dental_care
	case "car_insurance":
		return jobsRPC.JobDetails_car_insurance
	case "tution_fees":
		return jobsRPC.JobDetails_tution_fees
	case "permfomance_related_bonus":
		return jobsRPC.JobDetails_permfomance_related_bonus
	case "stock_options":
		return jobsRPC.JobDetails_stock_options
	case "profit_earning_bonus":
		return jobsRPC.JobDetails_profit_earning_bonus
	case "additional_months_salary":
		return jobsRPC.JobDetails_additional_months_salary
	case "employers_matching_contributions":
		return jobsRPC.JobDetails_employers_matching_contributions
	case "parental_bonus":
		return jobsRPC.JobDetails_parental_bonus
	case "tax_deductions":
		return jobsRPC.JobDetails_tax_deductions
	case "language_courses":
		return jobsRPC.JobDetails_language_courses
	case "company_car":
		return jobsRPC.JobDetails_company_car
	case "laptop":
		return jobsRPC.JobDetails_laptop
	case "discounts_on_company_products_and_services":
		return jobsRPC.JobDetails_discounts_on_company_products_and_services
	case "holiday_vouchers":
		return jobsRPC.JobDetails_holiday_vouchers
	case "restraunt_vouchers":
		return jobsRPC.JobDetails_restraunt_vouchers
	case "corporate_housing":
		return jobsRPC.JobDetails_corporate_housing
	case "mobile_phone":
		return jobsRPC.JobDetails_mobile_phone
	case "gift_vouchers":
		return jobsRPC.JobDetails_gift_vouchers
	case "cultural_or_sporting_activites":
		return jobsRPC.JobDetails_cultural_or_sporting_activites
	case "employee_service_vouchers":
		return jobsRPC.JobDetails_employee_service_vouchers
	case "corporate_credit_card":
		return jobsRPC.JobDetails_corporate_credit_card
	}

	return jobsRPC.JobDetails_other
}

func jobBenefitRPCToJobBenefit(s jobsRPC.JobDetails_JobBenefit) string {
	switch s {
	case jobsRPC.JobDetails_labor_agreement:
		return "labor_agreement"
	case jobsRPC.JobDetails_remote_working:
		return "remote_working"
	case jobsRPC.JobDetails_floater:
		return "floater"
	case jobsRPC.JobDetails_paid_timeoff:
		return "paid_timeoff"
	case jobsRPC.JobDetails_flexible_working_hours:
		return "flexible_working_hours"
	case jobsRPC.JobDetails_additional_timeoff:
		return "additional_timeoff"
	case jobsRPC.JobDetails_additional_parental_leave:
		return "additional_parental_leave"
	case jobsRPC.JobDetails_sick_leave_for_family_members:
		return "sick_leave_for_family_members"
	case jobsRPC.JobDetails_company_daycare:
		return "company_daycare"
	case jobsRPC.JobDetails_company_canteen:
		return "company_canteen"
	case jobsRPC.JobDetails_sport_facilities:
		return "sport_facilities"
	case jobsRPC.JobDetails_access_for_handicapped_persons:
		return "access_for_handicapped_persons"
	case jobsRPC.JobDetails_employee_parking:
		return "employee_parking"
	case jobsRPC.JobDetails_shuttle_service:
		return "shuttle_service"
	case jobsRPC.JobDetails_multiple_work_spaces:
		return "multiple_work_spaces"
	case jobsRPC.JobDetails_corporate_events:
		return "corporate_events"
	case jobsRPC.JobDetails_trainig_and_development:
		return "trainig_and_development"
	case jobsRPC.JobDetails_pets_allowed:
		return "pets_allowed"
	case jobsRPC.JobDetails_corporate_medical_staff:
		return "corporate_medical_staff"
	case jobsRPC.JobDetails_game_consoles:
		return "game_consoles"
	case jobsRPC.JobDetails_snack_and_drink_selfservice:
		return "snack_and_drink_selfservice"
	case jobsRPC.JobDetails_private_pension_scheme:
		return "private_pension_scheme"
	case jobsRPC.JobDetails_health_insurance:
		return "health_insurance"
	case jobsRPC.JobDetails_dental_care:
		return "dental_care"
	case jobsRPC.JobDetails_car_insurance:
		return "car_insurance"
	case jobsRPC.JobDetails_tution_fees:
		return "tution_fees"
	case jobsRPC.JobDetails_permfomance_related_bonus:
		return "permfomance_related_bonus"
	case jobsRPC.JobDetails_stock_options:
		return "stock_options"
	case jobsRPC.JobDetails_profit_earning_bonus:
		return "profit_earning_bonus"
	case jobsRPC.JobDetails_additional_months_salary:
		return "additional_months_salary"
	case jobsRPC.JobDetails_employers_matching_contributions:
		return "employers_matching_contributions"
	case jobsRPC.JobDetails_parental_bonus:
		return "parental_bonus"
	case jobsRPC.JobDetails_tax_deductions:
		return "tax_deductions"
	case jobsRPC.JobDetails_language_courses:
		return "language_courses"
	case jobsRPC.JobDetails_company_car:
		return "company_car"
	case jobsRPC.JobDetails_laptop:
		return "laptop"
	case jobsRPC.JobDetails_discounts_on_company_products_and_services:
		return "discounts_on_company_products_and_services"
	case jobsRPC.JobDetails_holiday_vouchers:
		return "holiday_vouchers"
	case jobsRPC.JobDetails_restraunt_vouchers:
		return "restraunt_vouchers"
	case jobsRPC.JobDetails_corporate_housing:
		return "corporate_housing"
	case jobsRPC.JobDetails_mobile_phone:
		return "mobile_phone"
	case jobsRPC.JobDetails_gift_vouchers:
		return "gift_vouchers"
	case jobsRPC.JobDetails_cultural_or_sporting_activites:
		return "cultural_or_sporting_activites"
	case jobsRPC.JobDetails_employee_service_vouchers:
		return "employee_service_vouchers"
	case jobsRPC.JobDetails_corporate_credit_card:
		return "corporate_credit_card"
	}

	return "other"
}

func jobsJobWithSeenStatToResolver(data *jobsRPC.ViewJobWithSeenStat) *JobWithSeenStatResolver {
	return &JobWithSeenStatResolver{
		R: &JobWithSeenStat{
			ID:            data.GetID(),
			Title:         data.GetTitle(),
			Total_amount:  data.GetTotalAmount(),
			Unseen_amount: data.GetUnseenAmount(),
			Status:        data.GetStatus(),
		},
	}
}

func stringToJobsRPCApplicantCategoryEnum(data string) jobsRPC.ApplicantCategoryEnum {
	switch data {
	case "Favorite":
		return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryFavorite
	case "In_review":
		return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryInReview
	case "Disqualified":
		return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryDisqualified
	}
	return jobsRPC.ApplicantCategoryEnum_ApplicantCategoryNone
}

func jobsRPCApplicantCategoryEnumToString(data jobsRPC.ApplicantCategoryEnum) string {
	switch data {
	case jobsRPC.ApplicantCategoryEnum_ApplicantCategoryFavorite:
		return "Favorite"
	case jobsRPC.ApplicantCategoryEnum_ApplicantCategoryInReview:
		return "In_review"
	case jobsRPC.ApplicantCategoryEnum_ApplicantCategoryDisqualified:
		return "Disqualified"
	}

	return "None"
}

func jobfileRPCToFile(data *jobsRPC.File) *File {
	if data == nil {
		return nil
	}

	f := File{
		ID:        data.GetID(),
		Address:   data.GetURL(),
		Mime_type: data.GetMimeType(),
		Name:      data.GetName(),
	}

	return &f
}

func jobExperienceEnumToString(s jobsRPC.ExperienceEnum) string {
	switch s {
	case jobsRPC.ExperienceEnum_WithoutExperience:
		return "without_experience"
	case jobsRPC.ExperienceEnum_LessThenOneYear:
		return "less_then_one_year"
	case jobsRPC.ExperienceEnum_OneTwoYears:
		return "one_two_years"
	case jobsRPC.ExperienceEnum_TwoThreeYears:
		return "two_three_years"
	case jobsRPC.ExperienceEnum_ThreeFiveYears:
		return "three_five_years"
	case jobsRPC.ExperienceEnum_FiveSevenyears:
		return "five_seven_years"
	case jobsRPC.ExperienceEnum_SevenTenYears:
		return "seven_ten_years"
	case jobsRPC.ExperienceEnum_TenYearsAndMore:
		return "ten_years_and_more"
	}

	return "experience_unknown"
}

func jobsRPCSalaryIntervalToString(s jobsRPC.SalaryInterval) string {
	switch s {
	case jobsRPC.SalaryInterval_Hour:
		return "Hour"
	case jobsRPC.SalaryInterval_Month:
		return "Month"
	case jobsRPC.SalaryInterval_Year:
		return "Year"
	}

	return "Unknown"
}

func jobsRPCJobTypeToString(s jobsRPC.JobType) string {
	switch s {
	case jobsRPC.JobType_Consultancy:
		return "Consultancy"
	case jobsRPC.JobType_Contractual:
		return "Contractual"
	case jobsRPC.JobType_FullTime:
		return "FullTime"
	case jobsRPC.JobType_Internship:
		return "Internship"
	case jobsRPC.JobType_PartTime:
		return "PartTime"
	case jobsRPC.JobType_Partner:
		return "Partner"
	case jobsRPC.JobType_Temporary:
		return "Temporary"
	case jobsRPC.JobType_Volunteer:
		return "Volunteer"
	}

	return "UnknownJobType"
}
