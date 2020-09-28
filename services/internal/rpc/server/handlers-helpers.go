package serverRPC

import (
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/review"

	offer "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/offers"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
	additionaldetails "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/additional-details"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/category"
	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/location"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/qualifications"
	serviceorder "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/service-order"
	office "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/v-office"
)

func createOfficeRequestRPCToOfficeAccount(data *servicesRPC.CreateVOfficeRequest) *office.Office {
	if data == nil {
		return nil
	}

	off := office.Office{
		Name:        data.GetName(),
		Description: data.GetDescription(),
		Category:    data.GetCategory(),
		Languages:   make([]*qualifications.Language, 0, len(data.GetLanguages())),
	}

	off.Location = location.Location{
		Country: &location.Country{},
		City: &location.City{
			ID: data.GetLocation().GetCity().GetId(),
		},
	}

	for _, lang := range data.GetLanguages() {
		off.Languages = append(off.Languages, serviceRPCLanguageToQualificationsLanguage(lang))
	}

	return &off
}

func pointerServicesSKillsRPCToQualificaitonsSkills(data *servicesRPC.Skill) *qualifications.Skill {
	if data == nil {
		return nil
	}
	skill := qualifications.Skill{
		Skill: data.GetSkill(),
	}
	_ = skill.SetID(data.GetID())

	return &skill
}

func pointerServiceToolsRPCToQualificationsTools(data *servicesRPC.ToolTechnology) *qualifications.ToolTechnology {
	if data == nil {
		return nil
	}

	tool := qualifications.ToolTechnology{
		ToolTechnology: data.GetToolTechnology(),
		Rank:           servicesLevelRPCToQualificationsLevel(data.GetRank()),
	}

	return &tool
}

func servicesLevelRPCToQualificationsLevel(data servicesRPC.Level) *qualifications.Level {
	level := qualifications.LevelBeginner

	switch data {
	case servicesRPC.Level_Intermediate:
		level = qualifications.LevelIntermediate
	case servicesRPC.Level_Advanced:
		level = qualifications.LevelAdvanced
	case servicesRPC.Level_Master:
		level = qualifications.LevelMaster
	}

	return &level
}

func qualificationsLevelToServicesRPCLevel(data *qualifications.Level) servicesRPC.Level {
	level := servicesRPC.Level_Beginner

	switch *data {
	case qualifications.LevelIntermediate:
		level = servicesRPC.Level_Intermediate
	case qualifications.LevelAdvanced:
		level = servicesRPC.Level_Advanced
	case qualifications.LevelMaster:
		level = servicesRPC.Level_Master
	}

	return level
}

func servicesRPCServiceToServiceStruct(data *servicesRPC.Service) *servicerequest.Service {
	if data == nil {
		return nil
	}
	objOfficeID, err := primitive.ObjectIDFromHex(data.OfficeID)
	if err != nil {
		return nil
	}

	serv := servicerequest.Service{
		OfficeID:          objOfficeID,
		Title:             data.Tittle,
		Description:       data.Description,
		Price:             servicesRPCPriceToPriceStruct(data.Price),
		FixedPriceAmmount: data.FixedPriceAmmount,
		MinPriceAmmout:    data.MinPriceAmmount,
		MaxPriceAmmout:    data.MaxPriceAmmount,
		DeliveryTime:      serviceRPCDeliveryTimeToDeliveryTime(data.DeliveryTime),
		LocationType:      servicesRPCLocationTypeToLocationType(data.LocationType),
		Location:          serviceRPCLocationToLocation(data.Location),
		Currency:          data.Currency,
		CreatedAt:         stringDayMonthAndYearToTime(data.CreatedAt),
		IsSaved:           data.IsSaved,
		IsRemote:          data.IsRemote,
		IsDraft:           data.IsDraft,
		WorkingHours:      serviceWorkinHourRPCToWorkingHour(data.GetWorkingDate()),
	}

	serv.Category = &category.Category{
		Main: data.GetCategory().GetMain(),
		Sub:  data.GetCategory().GetSub(),
	}

	serv.AdditionalDetails = &additionaldetails.AdditionalDetails{
		Qualifications: serviceRPCQualificationsToQualifications(data.AdditionalDetails.Qualifications),
	}

	return &serv
}

func requestRPCToRequest(data *servicesRPC.Request) *servicerequest.Request {
	if data == nil {
		return nil
	}

	request := servicerequest.Request{
		Tittle:            data.Tittle,
		CustomDate:        data.GetCustomDate(),
		Price:             servicesRPCPriceToPriceStruct(data.Price),
		Description:       data.Description,
		FixedPriceAmmount: data.FixedPriceAmmount,
		MinPriceAmmout:    data.MinPriceAmmount,
		MaxPriceAmmout:    data.MaxPriceAmmount,
		ProjectType:       servicesRPCProjectTypeToProjectType(data.GetProjectType()),
		DeliveryTime:      serviceRPCDeliveryTimeToDeliveryTime(data.DeliveryTime),
		LocationType:      servicesRPCLocationTypeToLocationType(data.LocationType),
		Location:          serviceRPCLocationToLocation(data.Location),
		Currency:          data.GetCurrency(),
		CreatedAt:         stringDayMonthAndYearToTime(data.CreatedAt),
		IsRemote:          data.IsRemote,
		IsDraft:           data.IsDraft,
		Status:            serviceRPCStatusToStatus(data.GetStatus()),
		AdditionalDetails: serviceRequAdditionalDetailsToAdditinalDetails(data.GetAdditionalDetails()),
	}

	request.Category = category.Category{
		Main: data.GetCategory().GetMain(),
		Sub:  data.GetCategory().GetSub(),
	}

	return &request
}

func serviceWorkinHourToRPC(data *servicerequest.WorkingHours) *servicesRPC.WorkingHour {
	if data == nil {
		return nil
	}

	return &servicesRPC.WorkingHour{
		IsAlwaysOpen: data.IsAlwaysOpen,
		WorkingHours: serviceWorkingDatesToRPC(data.WorkingDate),
	}
}

func serviceWorkingDatesToRPC(data []servicerequest.WorkingDate) []*servicesRPC.WorkingDate {
	if data == nil {
		return nil
	}

	days := make([]*servicesRPC.WorkingDate, 0, len(data))

	for _, d := range data {
		days = append(days, serviceWorkingDateToRPC(d))
	}

	return days
}

func serviceWorkingDateToRPC(data servicerequest.WorkingDate) *servicesRPC.WorkingDate {

	return &servicesRPC.WorkingDate{
		FromDate: data.HourFrom,
		ToDate:   data.HourTo,
		WeekDays: serviceWeekDaysToRPC(data.WeekDays),
	}

}

func serviceWeekDaysToRPC(data []servicerequest.WeekDay) []servicesRPC.WeekDays {
	if len(data) <= 0 {
		return nil
	}

	weekDays := make([]servicesRPC.WeekDays, 0, len(data))

	for _, wd := range data {
		weekDays = append(weekDays, serviceWeekDayToRPC(wd))
	}

	return weekDays

}
func serviceWeekDayToRPC(data servicerequest.WeekDay) servicesRPC.WeekDays {
	switch data {
	case servicerequest.WeekDayMonday:
		return servicesRPC.WeekDays_MONDAY
	case servicerequest.WeekDayTuesday:
		return servicesRPC.WeekDays_TUESDAY
	case servicerequest.WeekDayWednesday:
		return servicesRPC.WeekDays_WEDNESDAY
	case servicerequest.WeekDayThursday:
		return servicesRPC.WeekDays_THURSDAY
	case servicerequest.WeekDayFriday:
		return servicesRPC.WeekDays_FRIDAY
	case servicerequest.WeekDaySaturday:
		return servicesRPC.WeekDays_SATURDAY
	}
	return servicesRPC.WeekDays_SUNDAY
}

func serviceWorkinHourRPCToWorkingHour(data *servicesRPC.WorkingHour) *servicerequest.WorkingHours {
	if data == nil {
		return nil
	}

	return &servicerequest.WorkingHours{
		IsAlwaysOpen: data.GetIsAlwaysOpen(),
		WorkingDate:  serviceWorkingDatesRPCToWorkingDates(data.GetWorkingHours()),
	}
}

func serviceWorkingDatesRPCToWorkingDates(data []*servicesRPC.WorkingDate) []servicerequest.WorkingDate {
	if data == nil {
		return nil
	}

	days := make([]servicerequest.WorkingDate, 0, len(data))

	for _, d := range data {
		days = append(days, serviceWorkingDateRPCToWorkingDate(d))
	}

	return days
}

func serviceWorkingDateRPCToWorkingDate(data *servicesRPC.WorkingDate) servicerequest.WorkingDate {

	if data == nil {
		return servicerequest.WorkingDate{}
	}

	return servicerequest.WorkingDate{
		HourFrom: data.GetFromDate(),
		HourTo:   data.GetToDate(),
		WeekDays: serviceWeekDaysRPCToWeekDays(data.GetWeekDays()),
	}

}

func serviceWeekDaysRPCToWeekDays(data []servicesRPC.WeekDays) []servicerequest.WeekDay {
	if len(data) <= 0 {
		return nil
	}

	weekDays := make([]servicerequest.WeekDay, 0, len(data))

	for _, wd := range data {
		weekDays = append(weekDays, serviceWeekDayRPCToWeekDay(wd))
	}

	return weekDays

}
func serviceWeekDayRPCToWeekDay(data servicesRPC.WeekDays) servicerequest.WeekDay {
	switch data {
	case servicesRPC.WeekDays_MONDAY:
		return servicerequest.WeekDayMonday
	case servicesRPC.WeekDays_TUESDAY:
		return servicerequest.WeekDayTuesday
	case servicesRPC.WeekDays_WEDNESDAY:
		return servicerequest.WeekDayWednesday
	case servicesRPC.WeekDays_THURSDAY:
		return servicerequest.WeekDayThursday
	case servicesRPC.WeekDays_FRIDAY:
		return servicerequest.WeekDayFriday
	case servicesRPC.WeekDays_SATURDAY:
		return servicerequest.WeekDaySaturday
	}
	return servicerequest.WeekDaySunday
}

func serviceRequAdditionalDetailsToRPC(data *additionaldetails.AdditionalDetails) *servicesRPC.RequestAdditionalDetails {
	if data == nil {
		return nil
	}

	return &servicesRPC.RequestAdditionalDetails{
		Languages:       languageArrayToLanguageArrayRPC(data.Qualifications.Languages),
		Tools:           toolsArrayToToolsArrayRPC(data.Qualifications.ToolsAndTechnologies),
		Skills:          skillsArrayToSkillsArrayRPC(data.Qualifications.Skills),
		ServiceProvider: serviceProviderToRPC(data.ServiceProvider),
	}
}

func serviceRequAdditionalDetailsToAdditinalDetails(data *servicesRPC.RequestAdditionalDetails) *additionaldetails.AdditionalDetails {
	if data == nil {
		return nil
	}

	return &additionaldetails.AdditionalDetails{
		Qualifications: &qualifications.Qualifications{
			Languages:            serviceRPCLanguagesToLanguages(data.GetLanguages()),
			Skills:               skillsArrayRPCToArray(data.GetSkills()),
			ToolsAndTechnologies: toolsArrayRPCToArray(data.GetTools()),
		},
		ServiceProvider: serviceProviderRPCToProvider(data.GetServiceProvider()),
	}
}

func serviceProviderToRPC(data *additionaldetails.ServiceProvider) servicesRPC.ServiceProviderEnum {

	switch *data {
	case additionaldetails.ServiceProviderCompany:
		return servicesRPC.ServiceProviderEnum_Company
	case additionaldetails.ServiceProviderFreelancer:
		return servicesRPC.ServiceProviderEnum_Freelancer
	case additionaldetails.ServiceProviderProffesional:
		return servicesRPC.ServiceProviderEnum_Professional

	}
	return servicesRPC.ServiceProviderEnum_Professional
}

func serviceProviderRPCToProvider(data servicesRPC.ServiceProviderEnum) *additionaldetails.ServiceProvider {
	provider := additionaldetails.ServiceProviderNoPreference

	switch data {
	case servicesRPC.ServiceProviderEnum_Company:
		provider = additionaldetails.ServiceProviderCompany
	case servicesRPC.ServiceProviderEnum_Freelancer:
		provider = additionaldetails.ServiceProviderFreelancer
	case servicesRPC.ServiceProviderEnum_Professional:
		provider = additionaldetails.ServiceProviderProffesional

	}
	return &provider
}

func servicesRPCProjectTypeToProjectType(data servicesRPC.RequestProjectTypeEnum) servicerequest.ProjectType {
	status := servicerequest.ProjectTypeNotSure

	switch data {
	case servicesRPC.RequestProjectTypeEnum_On_Going_Project:
		return servicerequest.ProjectTypeOnGoing
	case servicesRPC.RequestProjectTypeEnum_One_Time_Project:
		return servicerequest.ProjectTypeOneTime
	}

	return status
}

func serviceRPCStatusToStatus(data servicesRPC.StatusEnum) servicerequest.Status {
	status := servicerequest.StatusActive

	switch data {
	case servicesRPC.StatusEnum_Status_Closed:
		return servicerequest.StatusClosed
	case servicesRPC.StatusEnum_Status_Draft:
		return servicerequest.StatusDraft
	case servicesRPC.StatusEnum_Status_Paused:
		return servicerequest.StatusPaused
	case servicesRPC.StatusEnum_Status_Pending:
		return servicerequest.StatusPending
	case servicesRPC.StatusEnum_Status_Rejected:
		return servicerequest.StatusRejected
	}

	return status
}

func serviceRPCToVisibilityType(data servicesRPC.VisibilityEnum) servicerequest.Visibility {
	visibility := servicerequest.VisibilityAnyone

	switch data {
	case servicesRPC.VisibilityEnum_Invited_Only:
		return servicerequest.VisibilityInvitedOnly
	case servicesRPC.VisibilityEnum_Only_RightNao_User:
		return servicerequest.VisibilityOnlyRigntNaoUsers
	}

	return visibility
}
func serviceRPCLocationToLocation(data *servicesRPC.Location) *location.Location {
	if data == nil {
		return nil
	}

	return &location.Location{
		Country: &location.Country{
			ID: data.GetCountry().GetId(),
		},
		City: &location.City{
			ID: data.GetCity().GetId(),
		},
	}
}

func serviceRPCDeliveryTimeToDeliveryTime(data servicesRPC.DeliveryTimeEnum) servicerequest.DeliveryTime {
	deliveryTime := servicerequest.DeliveryUpTo24Hours

	switch data {
	case servicesRPC.DeliveryTimeEnum_Up_To_3_Days:
		return servicerequest.DeliveryUpTo3Days
	case servicesRPC.DeliveryTimeEnum_Up_To_7_Days:
		return servicerequest.DeliveryUpTo7Days
	case servicesRPC.DeliveryTimeEnum_Weeks_1_2:
		return servicerequest.Delivery12Weeks
	case servicesRPC.DeliveryTimeEnum_Weeks_2_4:
		return servicerequest.Delivery2Weeks
	case servicesRPC.DeliveryTimeEnum_Month_And_More:
		return servicerequest.DeliveryMonthAndMore
	case servicesRPC.DeliveryTimeEnum_Custom:
		return servicerequest.DeliveryCustom
	}

	return deliveryTime
}

func servicesRPCLocationTypeToLocationType(data servicesRPC.LocationEnum) location.LocationType {

	if data == servicesRPC.LocationEnum_Remote_only {
		return location.LocationTypeRemoTeOnly
	}

	return location.LocationTypeOnSiteWork
}

func serviceRPCLanguagesToLanguages(data []*servicesRPC.Language) []*qualifications.Language {
	if data == nil {
		return nil
	}

	langs := make([]*qualifications.Language, 0, len(data))

	for _, lang := range data {
		langs = append(langs, serviceRPCLanguageToQualificationsLanguage(lang))
	}

	return langs

}

func serviceRPCQualificationsToQualifications(data *servicesRPC.Qualifications) *qualifications.Qualifications {
	if data == nil {
		return nil
	}
	qua := &qualifications.Qualifications{
		// Languages:            make([]*qualifications.Language, 0, len(data.GetLanguages())),
		Skills:               make([]*qualifications.Skill, 0, len(data.GetSkills())),
		ToolsAndTechnologies: make([]*qualifications.ToolTechnology, 0, len(data.GetTools())),
	}

	// for i := range data.GetLanguages() {
	// 	qua.Languages = append(qua.Languages, serviceRPCLanguageToQualificationsLanguage(data.GetLanguages()[i]))
	// }

	for i := range data.GetSkills() {
		qua.Skills = append(qua.Skills, serviceRPCSkillToQualificationsSkill(data.GetSkills()[i]))
	}

	for i := range data.GetTools() {
		qua.ToolsAndTechnologies = append(qua.ToolsAndTechnologies, serviceRPCToolsToQualificationsTools(data.GetTools()[i]))
	}

	return qua
}

func qualificationsToserviceRPCQualifications(data qualifications.Qualifications) *servicesRPC.Qualifications {
	qua := &servicesRPC.Qualifications{
		Languages: make([]*servicesRPC.Language, 0, len(data.Languages)),
		Skills:    make([]*servicesRPC.Skill, 0, len(data.Skills)),
		Tools:     make([]*servicesRPC.ToolTechnology, 0, len(data.ToolsAndTechnologies)),
	}

	for i := range data.Languages {
		qua.Languages = append(qua.Languages, qualificationsLanguageToServicesRPCLanguage(data.Languages[i]))
	}

	for i := range data.Skills {
		qua.Skills = append(qua.Skills, qualificationsSkillToServicesRPCSkill(data.Skills[i]))
	}

	for i := range data.ToolsAndTechnologies {
		qua.Tools = append(qua.Tools, qualificationsToolsTechnologiesToServicesRPCToolsTechnologies(data.ToolsAndTechnologies[i]))
	}

	return qua
}

func serviceRPCServiceIncludesToServiceIncludes(data servicesRPC.ServiceIncludesEnum) *additionaldetails.ServiceIncludes {
	serviceIncludes := additionaldetails.ServiceIncludesSourceFile

	switch data {
	case servicesRPC.ServiceIncludesEnum_Service_Include_Print_Ready:
		serviceIncludes = additionaldetails.ServiceIncludesPrintReady
	case servicesRPC.ServiceIncludesEnum_Service_Include_Photo_Editing:
		serviceIncludes = additionaldetails.ServiceIncludesPhotoEditing
	case servicesRPC.ServiceIncludesEnum_Service_Include_Custom_Graphics:
		serviceIncludes = additionaldetails.ServiceIncludesCustomGraphics
	case servicesRPC.ServiceIncludesEnum_Service_Include_Stock_Photos:
		serviceIncludes = additionaldetails.ServiceIncludesStockPhotos
	}

	return &serviceIncludes
}

func serviceIncludesRPCServiceIncludes(data *additionaldetails.ServiceIncludes) servicesRPC.ServiceIncludesEnum {
	serviceIncludes := servicesRPC.ServiceIncludesEnum_Service_Includes_Source_File

	switch *data {
	case additionaldetails.ServiceIncludesPrintReady:
		serviceIncludes = servicesRPC.ServiceIncludesEnum_Service_Include_Print_Ready
	case additionaldetails.ServiceIncludesPhotoEditing:
		serviceIncludes = servicesRPC.ServiceIncludesEnum_Service_Include_Photo_Editing
	case additionaldetails.ServiceIncludesCustomGraphics:
		serviceIncludes = servicesRPC.ServiceIncludesEnum_Service_Include_Custom_Graphics
	case additionaldetails.ServiceIncludesStockPhotos:
		serviceIncludes = servicesRPC.ServiceIncludesEnum_Service_Include_Stock_Photos
	}

	return serviceIncludes
}

func servicesRPCPriceToPriceStruct(data servicesRPC.PriceEnum) servicerequest.Price {
	price := servicerequest.PriceNegotiable

	switch data {
	case servicesRPC.PriceEnum_Price_Fixed:
		price = servicerequest.PriceFixed
	case servicesRPC.PriceEnum_Price_Hourly:
		price = servicerequest.PriceHourly
	}

	return price
}

func serviceRPCPurposeToPurpose(data servicesRPC.PurposeEnum) *additionaldetails.Purpose {
	purpose := additionaldetails.PurposeAcademic

	switch data {
	case servicesRPC.PurposeEnum_Purpose_Corporate:
		purpose = additionaldetails.PurposeCorporate
	case servicesRPC.PurposeEnum_Purpose_Event:
		purpose = additionaldetails.PurposeEvent
	case servicesRPC.PurposeEnum_Purpose_Promotional:
		purpose = additionaldetails.PurposePromotional
	case servicesRPC.PurposeEnum_Purpose_Medical:
		purpose = additionaldetails.PurposeMedical
	case servicesRPC.PurposeEnum_Purpose_Real_Estate:
		purpose = additionaldetails.PurposeRealEstate
	case servicesRPC.PurposeEnum_Purpose_Travel:
		purpose = additionaldetails.PurposeTravel
	}

	return &purpose
}

func timeToStringDayMonthAndYear(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	y, m, d := t.UTC().Date()
	return strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
}

func stringDayMonthAndYearToTime(s string) time.Time {
	if date, err := time.Parse("2-1-2006", s); err == nil {
		return date
	}
	return time.Time{}
}

func pointerDayMonthAndYearToTime(s string) *time.Time {
	if date, err := time.Parse("2-1-2006", s); err == nil {
		return &date
	}
	return nil
}

func serviceRPCPortfolioToPortfolio(data *servicesRPC.Portfolio) *office.Portfolio {
	if data == nil {
		return nil
	}

	port := office.Portfolio{
		ContentType: serviceRPCContentTypeEnumToContentTypeEnum(data.ContentType),
		CreatedAt:   stringDayMonthAndYearToTime(data.CreatedAt),
		Description: data.Description,
		Files:       make([]*file.File, 0, len(data.Files)),
		Tittle:      data.Tittle,
		Link:        make([]*file.Link, 0, len(data.Links)),
	}

	port.SetID(data.GetId())
	port.SetUserID(data.GetUserID())
	port.SetCompanyID(data.GetCompanyID())

	for _, link := range data.GetLinks() {
		var ids primitive.ObjectID

		if link.ID != "" {
			id, _ := primitive.ObjectIDFromHex(link.GetID())
			ids = id
		}
		port.Link = append(port.Link, &file.Link{
			ID:  ids,
			URL: link.URL,
		})
	}

	return &port
}

func officePortfolioToPortfolioRPC(data *office.Portfolio) *servicesRPC.Portfolio {
	if data == nil {
		return nil
	}

	port := servicesRPC.Portfolio{
		Id:          data.GetID(),
		UserID:      data.GetUserID(),
		CompanyID:   data.GetCompanyID(),
		ContentType: officeContentTypeToContentTypeEnum(data.ContentType),
		Tittle:      data.Tittle,
		CreatedAt:   timeToStringDayMonthAndYear(data.CreatedAt),
		Description: data.Description,
		Files:       fileToserviceRPCFile(data.Files),
		Links:       make([]*servicesRPC.Link, 0, len(data.Link)),
	}

	/// Links
	for _, link := range data.Link {
		port.Links = append(port.Links, &servicesRPC.Link{
			ID:  link.ID.Hex(),
			URL: link.URL,
		})
	}

	return &port
}

func portfolioRPCArrayToPortfolioProfileArray(data []*office.Portfolio) []*servicesRPC.Portfolio {
	if data == nil {
		return nil
	}

	portfolio := make([]*servicesRPC.Portfolio, 0, len(data))
	for _, p := range data {
		portfolio = append(portfolio, officePortfolioToPortfolioRPC(p))
	}

	return portfolio
}

func serviceRPCContentTypeEnumToContentTypeEnum(data servicesRPC.Portfolio_ContentTypeEnum) office.ContentType {
	contenttype := office.ContentTypeOther

	switch data {
	case servicesRPC.Portfolio_Content_Type_Audio:
		return office.ContentTypeAudio
	case servicesRPC.Portfolio_Content_Type_Code:
		return office.ContentTypeCode
	case servicesRPC.Portfolio_Content_Type_Image:
		return office.ContentTypeImage
	case servicesRPC.Portfolio_Content_Type_Article:
		return office.ContentTypeArticle
	case servicesRPC.Portfolio_Content_Type_Video:
		return office.ContentTypeVideo
	}

	return contenttype
}

func officeContentTypeToContentTypeEnum(data office.ContentType) servicesRPC.Portfolio_ContentTypeEnum {
	contenttype := servicesRPC.Portfolio_Content_Type_Other

	switch data {
	case office.ContentTypeAudio:
		return servicesRPC.Portfolio_Content_Type_Audio
	case office.ContentTypeCode:
		return servicesRPC.Portfolio_Content_Type_Code
	case office.ContentTypeImage:
		return servicesRPC.Portfolio_Content_Type_Image
	case office.ContentTypeArticle:
		return servicesRPC.Portfolio_Content_Type_Article
	case office.ContentTypeVideo:
		return servicesRPC.Portfolio_Content_Type_Video
	}

	return contenttype
}

func officesToOfficesRPC(data []*office.Office) []*servicesRPC.VOffice {
	if data == nil {
		return nil
	}

	offices := make([]*servicesRPC.VOffice, 0, len(data))

	for _, office := range data {
		offices = append(offices, officeToOfficeRPC(office))
	}

	return offices

}

func officeToOfficeRPC(data *office.Office) *servicesRPC.VOffice {
	if data == nil {
		return nil
	}
	office := &servicesRPC.VOffice{
		ID:          data.GetID(),
		UserID:      data.GetUserID(),
		CompanyID:   data.GetCompanyID(),
		IsMe:        data.IsMe,
		CreatedAt:   timeToStringDayMonthAndYear(data.CreatedAt),
		Description: data.Description,
		Location:    locationRPC(&data.Location),
		Cover:       data.CoverImage,
		CoverOrigin: data.CoverOriginImage,
		Name:        data.Name,
		Category:    data.Category,
		Languages:   make([]*servicesRPC.Language, 0, len(data.Languages)),
		IsOut:       data.IsOut,
	}

	for _, lang := range data.Languages {
		office.Languages = append(office.Languages, qualificationsLanguageToServicesRPCLanguage(lang))
	}

	if data.ReturnDate != nil {
		office.ReturnDate = timeToStringDayMonthAndYear(*data.ReturnDate)
	}

	return office
}

// @TODO
func servicesRequestToRPC(data []*servicerequest.Request) []*servicesRPC.Request {
	if data == nil {
		return nil
	}

	service := make([]*servicesRPC.Request, 0, len(data))
	for _, p := range data {
		service = append(service, serviceRequestToRPC(p))
	}

	return service
}

func serviceRequestToRPC(p *servicerequest.Request) *servicesRPC.Request {
	if p == nil {
		return nil
	}

	return &servicesRPC.Request{
		ID:                p.GetID(),
		UserID:            p.GetUserID(),
		CompanyID:         p.GetCompanyID(),
		Category:          categoryToCategoryRPC(&p.Category),
		CreatedAt:         timeToStringDayMonthAndYear(p.CreatedAt),
		Currency:          p.Currency,
		CustomDate:        p.CustomDate,
		DeliveryTime:      deliveryTimeToDeliveryTimeRPC(p.DeliveryTime),
		Files:             fileToserviceRPCFile(p.Files),
		FixedPriceAmmount: p.FixedPriceAmmount,
		IsDraft:           p.IsDraft,
		IsRemote:          p.IsRemote,
		IsClosed:          p.IsClosed,
		IsPaused:          p.IsPaused,
		ProposalAmount:    p.ProposalAmount,
		Location:          locationRPC(p.Location),
		LocationType:      locationTypeToLocationRPC(p.LocationType),
		MaxPriceAmmount:   p.MaxPriceAmmout,
		MinPriceAmmount:   p.MinPriceAmmout,
		Price:             priceRPC(p.Price),
		ProjectType:       projcetTypeRPC(p.ProjectType),
		Status:            statusRPC(p.GetStatus()),
		Tittle:            p.Tittle,
		AdditionalDetails: serviceRequAdditionalDetailsToRPC(p.AdditionalDetails),
		Description:       p.Description,
		HasLiked:          p.HasLiked,
	}
}

func projcetTypeRPC(data servicerequest.ProjectType) servicesRPC.RequestProjectTypeEnum {
	projectEnum := servicesRPC.RequestProjectTypeEnum_Not_Sure

	switch data {
	case servicerequest.ProjectTypeOnGoing:
		projectEnum = servicesRPC.RequestProjectTypeEnum_On_Going_Project
	case servicerequest.ProjectTypeOneTime:
		projectEnum = servicesRPC.RequestProjectTypeEnum_One_Time_Project
	}

	return projectEnum
}

func visibilityRPC(data servicerequest.Visibility) servicesRPC.VisibilityEnum {
	visibilityEnum := servicesRPC.VisibilityEnum_Anyone

	switch data {
	case servicerequest.VisibilityInvitedOnly:
		visibilityEnum = servicesRPC.VisibilityEnum_Invited_Only
	case servicerequest.VisibilityOnlyRigntNaoUsers:
		visibilityEnum = servicesRPC.VisibilityEnum_Only_RightNao_User
	}

	return visibilityEnum
}

func statusRPC(data servicerequest.Status) servicesRPC.StatusEnum {
	statusEnum := servicesRPC.StatusEnum_Status_Active

	switch data {
	case servicerequest.StatusClosed:
		statusEnum = servicesRPC.StatusEnum_Status_Closed
	case servicerequest.StatusDraft:
		statusEnum = servicesRPC.StatusEnum_Status_Draft
	case servicerequest.StatusPaused:
		statusEnum = servicesRPC.StatusEnum_Status_Paused
	case servicerequest.StatusPending:
		statusEnum = servicesRPC.StatusEnum_Status_Pending
	case servicerequest.StatusRejected:
		statusEnum = servicesRPC.StatusEnum_Status_Rejected

	}

	return statusEnum
}

func servicesToServiceRPCService(data servicerequest.Service) *servicesRPC.Service {

	ser := &servicesRPC.Service{
		ID:                data.GetID(),
		UserID:            data.GetUserID(),
		CompanyID:         data.GetCompanyID(),
		OfficeID:          data.GetOfficeID(),
		Files:             fileToserviceRPCFile(data.Files),
		Tittle:            data.Title,
		Description:       data.Description,
		DeliveryTime:      deliveryTimeToDeliveryTimeRPC(data.DeliveryTime),
		Currency:          data.Currency,
		Price:             priceRPC(data.Price),
		Category:          categoryToCategoryRPC(data.Category),
		FixedPriceAmmount: data.FixedPriceAmmount,
		LocationType:      locationTypeToLocationRPC(data.LocationType),
		MaxPriceAmmount:   data.MaxPriceAmmout,
		MinPriceAmmount:   data.MinPriceAmmout,
		IsDraft:           data.IsDraft,
		IsRemote:          data.IsRemote,
		IsSaved:           data.IsSaved,
		Location:          locationRPC(data.Location),
	}

	return ser
}

func fileToserviceRPCFile(data []*file.File) []*servicesRPC.File {
	if data == nil {
		return nil
	}

	files := make([]*servicesRPC.File, 0, len(data))
	for _, file := range data {
		files = append(files, &servicesRPC.File{
			ID:       file.GetID(),
			Name:     file.Name,
			MimeType: file.MimeType,
			URL:      file.URL,
		})
	}

	return files
}

func priceRPC(data servicerequest.Price) servicesRPC.PriceEnum {
	priceEnum := servicesRPC.PriceEnum_Price_Fixed

	switch data {
	case servicerequest.PriceHourly:
		priceEnum = servicesRPC.PriceEnum_Price_Hourly
	case servicerequest.PriceNegotiable:
		priceEnum = servicesRPC.PriceEnum_Price_Negotiable
	}

	return priceEnum
}
func locationRPC(data *location.Location) *servicesRPC.Location {
	if data == nil {
		return nil
	}
	return &servicesRPC.Location{
		Country: locationCountryToRPC(data.Country),
		City:    locationCityToRPC(data.City),
	}
}

func locationCountryToRPC(data *location.Country) *servicesRPC.Country {
	return &servicesRPC.Country{
		Id: data.ID,
	}
}

func locationCityToRPC(data *location.City) *servicesRPC.City {
	return &servicesRPC.City{
		Id:          data.ID,
		Title:       data.Name,
		Subdivision: data.Subdivision,
	}
}

func vOfficeCategoryToCategoryRPC(data category.VOfficeCategory) *servicesRPC.Category {
	return &servicesRPC.Category{
		Main: data.Main,
	}
}

func fileRPCToProfileFile(data *servicesRPC.File) *file.File {
	if data == nil {
		return nil
	}

	f := file.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		URL:      data.GetURL(),
	}

	_ = f.SetID(data.GetID())

	return &f
}

func qualificationsLanguageToServicesRPCLanguage(data *qualifications.Language) *servicesRPC.Language {
	if data == nil {
		return nil
	}

	qual := &servicesRPC.Language{
		ID:       data.GetID(),
		Language: data.Language,
		Rank:     qualificationsLevelToServicesRPCLevel(data.Rank),
	}

	return qual
}

func qualificationsSkillToServicesRPCSkill(data *qualifications.Skill) *servicesRPC.Skill {
	if data == nil {
		return nil
	}

	qual := &servicesRPC.Skill{
		ID:    data.GetID(),
		Skill: data.Skill,
	}

	return qual
}

func qualificationsToolsTechnologiesToServicesRPCToolsTechnologies(data *qualifications.ToolTechnology) *servicesRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	qual := &servicesRPC.ToolTechnology{
		ID:             data.GetID(),
		ToolTechnology: data.ToolTechnology,
		Rank:           qualificationsLevelToServicesRPCLevel(data.Rank),
	}

	return qual
}

func serviceRPCLanguageToQualificationsLanguage(data *servicesRPC.Language) *qualifications.Language {
	if data == nil {
		return nil
	}

	qual := &qualifications.Language{
		Language: data.GetLanguage(),
		Rank:     servicesLevelRPCToQualificationsLevel(data.GetRank()),
	}

	if data.GetID() == "" {
		qual.GenerateID()
	} else {
		_ = qual.SetID(data.GetID())
	}

	return qual
}

func serviceRPCSkillToQualificationsSkill(data *servicesRPC.Skill) *qualifications.Skill {
	if data == nil {
		return nil
	}

	qual := &qualifications.Skill{
		Skill: data.GetSkill(),
	}

	if data.GetID() == "" {
		qual.GenerateID()
	} else {
		_ = qual.SetID(data.GetID())
	}

	return qual
}

func serviceRPCToolsToQualificationsTools(data *servicesRPC.ToolTechnology) *qualifications.ToolTechnology {
	if data == nil {
		return nil
	}

	qual := &qualifications.ToolTechnology{
		ToolTechnology: data.GetToolTechnology(),
		Rank:           servicesLevelRPCToQualificationsLevel(data.GetRank()),
	}

	if data.GetID() == "" {
		qual.GenerateID()
	} else {
		_ = qual.SetID(data.GetID())
	}

	return qual
}

func serviceArrayToGetVOfficeResponse(data []*servicerequest.Service) *servicesRPC.GetVOfficeServicesResponse {
	if data == nil {
		return nil
	}

	return &servicesRPC.GetVOfficeServicesResponse{
		Services: serviceArrayToServiceArrayRPC(data),
	}
}

func serviceArrayToServiceArrayRPC(data []*servicerequest.Service) []*servicesRPC.Service {
	if data == nil || len(data) <= 0 {
		return nil
	}

	serviceArray := make([]*servicesRPC.Service, 0, len(data))
	for _, s := range data {
		serviceArray = append(serviceArray, serviceRPCToServiceRPC(s))
	}

	return serviceArray
}

func serviceRPCToServiceRPC(data *servicerequest.Service) *servicesRPC.Service {
	if data == nil {
		return nil
	}
	return &servicesRPC.Service{
		ID:                data.GetID(),
		AdditionalDetails: additionalDetailsToAdditionalDetailsRPC(data.AdditionalDetails),
		Action:            actionToActionRPC(data.Action),
		Cancellations:     data.Cancellations,
		Category:          categoryToCategoryRPC(data.Category),
		UserID:            data.GetUserID(),
		CompanyID:         data.GetCompanyID(),
		Currency:          data.Currency,
		OfficeID:          data.GetOfficeID(),
		Description:       data.Description,
		Tittle:            data.Title,
		Files:             fileToserviceRPCFile(data.Files),
		Price:             priceRPC(data.Price),
		FixedPriceAmmount: data.FixedPriceAmmount,
		MinPriceAmmount:   data.MinPriceAmmout,
		MaxPriceAmmount:   data.MaxPriceAmmout,
		DeliveryTime:      deliveryTimeToDeliveryTimeRPC(data.DeliveryTime),
		LocationType:      locationTypeToLocationRPC(data.LocationType),
		Location:          locationRPC(data.Location),
		IsDraft:           data.IsDraft,
		IsRemote:          data.IsRemote,
		IsSaved:           data.IsSaved,
		IsPaused:          data.IsPaused,
		HasLiked:          data.HasLiked,
		WorkingDate:       serviceWorkinHourToRPC(data.WorkingHours),
	}

}

func locationTypeToLocationRPC(data location.LocationType) servicesRPC.LocationEnum {
	if data == location.LocationTypeRemoTeOnly {
		return servicesRPC.LocationEnum_Remote_only
	}
	return servicesRPC.LocationEnum_On_Site_Work
}

func deliveryTimeToDeliveryTimeRPC(data servicerequest.DeliveryTime) servicesRPC.DeliveryTimeEnum {
	switch data {
	case servicerequest.DeliveryUpTo24Hours:
		return servicesRPC.DeliveryTimeEnum_Up_To_24_Hours
	case servicerequest.DeliveryUpTo3Days:
		return servicesRPC.DeliveryTimeEnum_Up_To_3_Days
	case servicerequest.DeliveryUpTo7Days:
		return servicesRPC.DeliveryTimeEnum_Up_To_7_Days
	case servicerequest.Delivery12Weeks:
		return servicesRPC.DeliveryTimeEnum_Weeks_1_2
	case servicerequest.Delivery2Weeks:
		return servicesRPC.DeliveryTimeEnum_Weeks_2_4
	case servicerequest.DeliveryMonthAndMore:
		return servicesRPC.DeliveryTimeEnum_Month_And_More
	}

	return servicesRPC.DeliveryTimeEnum_Custom

}

func additionalDetailsToAdditionalDetailsRPC(data *additionaldetails.AdditionalDetails) *servicesRPC.AdditionalDetails {
	if data == nil {
		return nil
	}

	return &servicesRPC.AdditionalDetails{
		// Purpose:         additionalPurposeToAdditionaPurposeRPC(data.Purpose),
		Qualifications: qualificationsToQualificationsRPC(data.Qualifications),
		// ServiceIncludes: serviceIncludesRPCServiceIncludes(data.ServiceIncludes),
	}
}

func additionalPurposeToAdditionaPurposeRPC(data *additionaldetails.Purpose) servicesRPC.PurposeEnum {
	if data == nil {
		return servicesRPC.PurposeEnum_Purpose_Academic
	}

	switch *data {
	case additionaldetails.PurposeCorporate:
		return servicesRPC.PurposeEnum_Purpose_Corporate
	case additionaldetails.PurposeEvent:
		return servicesRPC.PurposeEnum_Purpose_Event
	case additionaldetails.PurposeMedical:
		return servicesRPC.PurposeEnum_Purpose_Medical
	case additionaldetails.PurposePromotional:
		return servicesRPC.PurposeEnum_Purpose_Promotional
	case additionaldetails.PurposeRealEstate:
		return servicesRPC.PurposeEnum_Purpose_Real_Estate
	case additionaldetails.PurposeTravel:
		return servicesRPC.PurposeEnum_Purpose_Travel
	}

	return servicesRPC.PurposeEnum_Purpose_Travel
}

func qualificationsToQualificationsRPC(data *qualifications.Qualifications) *servicesRPC.Qualifications {
	if data == nil {
		return nil
	}
	return &servicesRPC.Qualifications{
		Languages: languageArrayToLanguageArrayRPC(data.Languages),
		Skills:    skillsArrayToSkillsArrayRPC(data.Skills),
		Tools:     toolsArrayToToolsArrayRPC(data.ToolsAndTechnologies),
	}
}

func languageArrayToLanguageArrayRPC(data []*qualifications.Language) []*servicesRPC.Language {
	if data == nil {
		return nil
	}

	langArray := make([]*servicesRPC.Language, 0, len(data))

	for _, lang := range data {
		langArray = append(langArray, languageToLanguageRPC(lang))
	}

	return langArray
}

func languageToLanguageRPC(data *qualifications.Language) *servicesRPC.Language {
	if data == nil {
		return nil
	}

	return &servicesRPC.Language{
		ID:       data.GetID(),
		Language: data.Language,
		Rank:     qualificationsLevelToServicesRPCLevel(data.Rank),
	}
}

func skillsArrayRPCToArray(data []*servicesRPC.Skill) []*qualifications.Skill {
	if data == nil {
		return nil
	}

	skillArray := make([]*qualifications.Skill, 0, len(data))

	for _, sk := range data {
		skillArray = append(skillArray, skillRPCToSkill(sk))
	}

	return skillArray
}

func skillsArrayToSkillsArrayRPC(data []*qualifications.Skill) []*servicesRPC.Skill {
	if data == nil {
		return nil
	}

	skillArray := make([]*servicesRPC.Skill, 0, len(data))

	for _, sk := range data {
		skillArray = append(skillArray, skillToSkillRPC(sk))
	}

	return skillArray
}

func skillRPCToSkill(data *servicesRPC.Skill) (skill *qualifications.Skill) {
	if data == nil {
		return nil
	}

	res := &qualifications.Skill{
		Skill: data.Skill,
	}

	if data.GetID() == "" {
		res.GenerateID()
	} else {
		_ = res.SetID(data.GetID())
	}

	return res

}

func skillToSkillRPC(data *qualifications.Skill) *servicesRPC.Skill {
	if data == nil {
		return nil
	}

	return &servicesRPC.Skill{
		ID:    data.GetID(),
		Skill: data.Skill,
	}
}

func toolsArrayRPCToArray(data []*servicesRPC.ToolTechnology) []*qualifications.ToolTechnology {
	if data == nil {
		return nil
	}

	toolArray := make([]*qualifications.ToolTechnology, 0, len(data))

	for _, tool := range data {
		toolArray = append(toolArray, toolRPCToTool(tool))
	}

	return toolArray
}

func toolRPCToTool(data *servicesRPC.ToolTechnology) *qualifications.ToolTechnology {
	if data == nil {
		return nil
	}

	res := &qualifications.ToolTechnology{
		ToolTechnology: data.ToolTechnology,
		Rank:           servicesLevelRPCToQualificationsLevel(data.Rank),
	}

	if data.GetID() == "" {
		res.GenerateID()
	} else {
		_ = res.SetID(data.GetID())
	}

	return res
}

func toolsArrayToToolsArrayRPC(data []*qualifications.ToolTechnology) []*servicesRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	toolArray := make([]*servicesRPC.ToolTechnology, 0, len(data))

	for _, tool := range data {
		toolArray = append(toolArray, toolToToolRPC(tool))
	}

	return toolArray
}

func toolToToolRPC(data *qualifications.ToolTechnology) *servicesRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	return &servicesRPC.ToolTechnology{
		ID:             data.GetID(),
		ToolTechnology: data.ToolTechnology,
		Rank:           qualificationsLevelToServicesRPCLevel(data.Rank),
	}
}

func actionToActionRPC(data servicerequest.Action) servicesRPC.ActionEnum {
	switch data {
	case servicerequest.ActionDelete:
		return servicesRPC.ActionEnum_Action_Delete
	case servicerequest.ActionEdit:
		return servicesRPC.ActionEnum_Action_Edit
	case servicerequest.ActionPause:
		return servicesRPC.ActionEnum_Action_Pause
	}

	return servicesRPC.ActionEnum_Action_Share
}

func categoryToCategoryRPC(data *category.Category) *servicesRPC.ServiceCategory {
	if data == nil {
		return nil
	}

	return &servicesRPC.ServiceCategory{
		Main: data.Main,
		Sub:  data.Sub,
	}
}

func orderTypeRPCToOrderType(data servicesRPC.OrderType) serviceorder.OrderType {

	if data == servicesRPC.OrderType_SELLER {
		return serviceorder.OrderTypeSeller
	}
	return serviceorder.OrderTypeBuyer
}

func orderStatusRPCToOrderStatus(data servicesRPC.OrderStatusEnum) serviceorder.OrderStatus {

	switch data {
	case servicesRPC.OrderStatusEnum_Status_Canceled:
		return serviceorder.OrderStatusCanceled
	case servicesRPC.OrderStatusEnum_Status_Completed:
		return serviceorder.OrderStatusCompleted
	case servicesRPC.OrderStatusEnum_Status_Delivered:
		return serviceorder.OrderStatusDelivered
	case servicesRPC.OrderStatusEnum_Status_Disputed:
		return serviceorder.OrderStatusDisputed
	case servicesRPC.OrderStatusEnum_Status_In_Progress:
		return serviceorder.OrderStatusInProgress
	case servicesRPC.OrderStatusEnum_Status_Out_Of_Schedule:
		return serviceorder.OrderStatusOutOfSchedule
	case servicesRPC.OrderStatusEnum_Status_New:
		return serviceorder.OrderStatusNew

	}
	return serviceorder.OrderStatusAny
}

func orderStatusToRPC(data serviceorder.OrderStatus) servicesRPC.OrderStatusEnum {

	switch data {
	case serviceorder.OrderStatusCanceled:
		return servicesRPC.OrderStatusEnum_Status_Canceled
	case serviceorder.OrderStatusCompleted:
		return servicesRPC.OrderStatusEnum_Status_Completed
	case serviceorder.OrderStatusDelivered:
		return servicesRPC.OrderStatusEnum_Status_Delivered
	case serviceorder.OrderStatusDisputed:
		return servicesRPC.OrderStatusEnum_Status_Disputed
	case serviceorder.OrderStatusInProgress:
		return servicesRPC.OrderStatusEnum_Status_In_Progress
	case serviceorder.OrderStatusOutOfSchedule:
		return servicesRPC.OrderStatusEnum_Status_Out_Of_Schedule
	case serviceorder.OrderStatusNew:
		return servicesRPC.OrderStatusEnum_Status_New

	}
	return servicesRPC.OrderStatusEnum_Status_Any
}

func ordersToRPC(data []serviceorder.Order) []*servicesRPC.OrderService {
	if len(data) <= 0 {
		return nil
	}

	orders := make([]*servicesRPC.OrderService, 0, len(data))

	for _, o := range data {
		orders = append(orders, orderToRPC(o))
	}

	return orders
}

func orderToRPC(data serviceorder.Order) *servicesRPC.OrderService {

	orderDetaill := data.OrderDetail

	res := &servicesRPC.OrderService{
		ID:              data.GetID(),
		IsCompany:       orderDetaill.IsCompany,
		Currency:        orderDetaill.Currency,
		Description:     orderDetaill.Description,
		PriceAmount:     orderDetaill.PriceAmount,
		DeliveryTime:    deliveryTimeToDeliveryTimeRPC(orderDetaill.DeliveryTime),
		Files:           fileToserviceRPCFile(orderDetaill.Files),
		Service:         servicesToServiceRPCService(data.Service),
		Request:         serviceRequestToRPC(&data.Request),
		OrderStatus:     orderStatusToRPC(orderDetaill.Status),
		PriceType:       priceRPC(orderDetaill.PriceType),
		MinPriceAmmount: orderDetaill.MinPriceAmount,
		MaxPriceAmmount: orderDetaill.MaxPriceAmount,
		CustomDate:      orderDetaill.CustomDate,
	}

	if !orderDetaill.ProfileID.IsZero() {
		res.ProfileID = orderDetaill.GetProfileID()
	}

	if orderDetaill.Note != nil {
		res.Note = *orderDetaill.Note
	}

	return res
}

func proposalsToRPC(data []offer.Proposal) []*servicesRPC.ProposalDetail {
	if len(data) <= 0 {
		return nil
	}

	proposals := make([]*servicesRPC.ProposalDetail, 0, len(data))

	for _, p := range data {
		proposals = append(proposals, proposalToRPC(p))
	}

	return proposals
}

func proposalToRPC(data offer.Proposal) *servicesRPC.ProposalDetail {

	proposalDetaill := data.ProposalDetail

	res := &servicesRPC.ProposalDetail{
		ID:              data.GetID(),
		ProfileID:       proposalDetaill.GetProfileID(),
		Currency:        proposalDetaill.Currency,
		Message:         proposalDetaill.Message,
		PriceAmount:     proposalDetaill.PriceAmount,
		MinPriceAmmount: proposalDetaill.MinPriceAmount,
		MaxPriceAmmount: proposalDetaill.MaxPriceAmount,
		ExperationTime:  proposalDetaill.ExperationTime,
		CustomDate:      proposalDetaill.CustomDate,
		DeliveryTime:    deliveryTimeToDeliveryTimeRPC(proposalDetaill.DeliveryTime),
		Service:         servicesToServiceRPCService(data.Service),
		Request:         serviceRequestToRPC(data.Request),
		OrderStatus:     orderStatusToRPC(proposalDetaill.Status),
		PriceType:       priceRPC(proposalDetaill.PriceType),
	}

	if !proposalDetaill.ProfileID.IsZero() {
		res.ProfileID = proposalDetaill.GetProfileID()
	}

	return res
}

func proposalRPCToProposal(data *servicesRPC.SendProposalRequest) *offer.Proposal {
	if data == nil {
		return nil
	}

	// Owner id
	objOwnerID, err := primitive.ObjectIDFromHex(data.GetOwnerID())
	if err != nil {
		return nil
	}

	// Request id
	objRequestID, err := primitive.ObjectIDFromHex(data.GetRequestID())
	if err != nil {
		return nil
	}

	return &offer.Proposal{
		OwnerID:        objOwnerID,
		RequestID:      objRequestID,
		IsOwnerCompany: data.GetIsOwnerCompany(),
		ProposalDetail: proposalDetailRPCToProposalDetail(data.GetProposalDetail()),
	}
}

func proposalDetailRPCToProposalDetail(data *servicesRPC.ProposalDetail) (od offer.ProposalDetail) {
	if data == nil {
		return
	}

	// Profile id
	objProfileID, err := primitive.ObjectIDFromHex(data.GetProfileID())
	if err != nil {
		return
	}

	// Service id
	objServiceID, err := primitive.ObjectIDFromHex(data.GetServiceID())
	if err != nil {
		return
	}

	// Office id
	objOfficeID, err := primitive.ObjectIDFromHex(data.GetOfficeID())
	if err != nil {
		return
	}

	return offer.ProposalDetail{
		ProfileID:      objProfileID,
		IsCompany:      data.GetIsCompany(),
		ServiceID:      objServiceID,
		OfficeID:       objOfficeID,
		Message:        data.GetMessage(),
		CustomDate:     data.GetCustomDate(),
		ExperationTime: data.GetExperationTime(),
		PriceType:      servicesRPCPriceToPriceStruct(data.GetPriceType()),
		PriceAmount:    data.GetPriceAmount(),
		MinPriceAmount: data.GetMinPriceAmmount(),
		MaxPriceAmount: data.GetMaxPriceAmmount(),
		Currency:       data.GetCurrency(),
		DeliveryTime:   serviceRPCDeliveryTimeToDeliveryTime(data.GetDeliveryTime()),
	}
}

func orderRPCToOrder(data *servicesRPC.OrderServiceRequest) *serviceorder.Order {
	if data == nil {
		return nil
	}

	// Owner id
	objOwnerID, err := primitive.ObjectIDFromHex(data.GetOwnerID())
	if err != nil {
		return nil
	}

	// Service id
	objServiceID, err := primitive.ObjectIDFromHex(data.GetServiceID())
	if err != nil {
		return nil
	}

	// Office id
	objOfficeID, err := primitive.ObjectIDFromHex(data.GetOfficeID())
	if err != nil {
		return nil
	}

	return &serviceorder.Order{
		IsOwnerCompany: data.GetIsOwnerCompany(),
		OwnerID:        objOwnerID,
		ServiceID:      objServiceID,
		OfficeID:       objOfficeID,
		OrderDetail:    orderDetailRPCToOrderDetail(data.GetOrderDetail()),
	}
}

func orderDetailRPCToOrderDetail(data *servicesRPC.OrderService) (od serviceorder.OrderDetail) {
	if data == nil {
		return
	}

	// Buyer id
	objBuyerID, err := primitive.ObjectIDFromHex(data.GetProfileID())
	if err != nil {
		return
	}

	return serviceorder.OrderDetail{
		ProfileID:      objBuyerID,
		IsCompany:      data.GetIsCompany(),
		Description:    data.GetDescription(),
		PriceType:      servicesRPCPriceToPriceStruct(data.GetPriceType()),
		PriceAmount:    data.GetPriceAmount(),
		Currency:       data.GetCurrency(),
		DeliveryTime:   serviceRPCDeliveryTimeToDeliveryTime(data.GetDeliveryTime()),
		MinPriceAmount: data.GetMinPriceAmmount(),
		MaxPriceAmount: data.GetMaxPriceAmmount(),
		CustomDate:     data.GetCustomDate(),
	}
}

func reviewServiceRequestRPCToReviewServiceRequest(data *servicesRPC.WriteReviewRequest) (rev review.Review) {
	if data == nil {
		return
	}

	// owner id
	objOwnerID, err := primitive.ObjectIDFromHex(data.GetOwnerID())
	if err != nil {
		return
	}

	return review.Review{
		IsOwnerCompnay: data.GetIsOwnerCompany(),
		OwnerID:        objOwnerID,
		ReviewDetail:   reviewDetailRPCToReviewDetail(data.GetReviewDetail()),
	}
}

func reviewServiceRPCToReviewService(data *servicesRPC.WriteReviewRequest) (rev review.Review) {
	if data == nil {
		return
	}

	// office id
	objOfficeID, err := primitive.ObjectIDFromHex(data.GetOfficeID())
	if err != nil {
		return
	}

	// service id
	objServiceID, err := primitive.ObjectIDFromHex(data.GetServiceID())
	if err != nil {
		return
	}

	// owner id
	objOwnerID, err := primitive.ObjectIDFromHex(data.GetOwnerID())
	if err != nil {
		return
	}

	return review.Review{
		OfficeID:     &objOfficeID,
		ServiceID:    objServiceID,
		OwnerID:      objOwnerID,
		ReviewDetail: reviewDetailRPCToReviewDetail(data.GetReviewDetail()),
	}
}

func reviewHireRPCToReviewHire(data servicesRPC.HireEnum) review.Hire {

	switch data {
	case servicesRPC.HireEnum_Not_Hire:
		return review.NotHire
	case servicesRPC.HireEnum_Will_Hire:
		return review.WillHire
	}

	return review.NotAnswer
}

func reviewHireToRPC(data review.Hire) servicesRPC.HireEnum {

	switch data {
	case review.NotHire:
		return servicesRPC.HireEnum_Not_Hire
	case review.WillHire:
		return servicesRPC.HireEnum_Will_Hire
	}

	return servicesRPC.HireEnum_Not_Answer
}

func reviewDetailRPCToReviewDetail(data *servicesRPC.ReviewDetail) (r review.ReviewDetail) {
	if data == nil {
		return
	}

	// profile id
	objProfileID, err := primitive.ObjectIDFromHex(data.GetProfileID())
	if err != nil {
		return
	}

	return review.ReviewDetail{
		ProfileID:     objProfileID,
		IsCompany:     data.GetIsCompany(),
		Clarity:       data.GetClarity(),
		Communication: data.GetCommunication(),
		Payment:       data.GetPayment(),
		Hire:          reviewHireRPCToReviewHire(data.GetHire()),
		Description:   data.GetDescription(),
	}
}

func reviewAVGToRPC(data *review.GetReview) *servicesRPC.ServiceReviewAVG {
	if data == nil {
		return nil
	}

	return &servicesRPC.ServiceReviewAVG{
		ClarityAVG:       data.ClartityAVG,
		CommunicationAVG: data.CommunicationAVG,
		PaymentAVG:       data.PaymentAVG,
	}
}

func reviewsToRPC(data []review.Review) []*servicesRPC.ReviewDetail {
	if len(data) <= 0 {
		return nil
	}

	reviews := make([]*servicesRPC.ReviewDetail, 0, len(data))

	for _, r := range data {
		reviews = append(reviews, reviewToRPC(r))
	}

	return reviews
}

func reviewToRPC(data review.Review) *servicesRPC.ReviewDetail {
	detail := data.ReviewDetail

	res := &servicesRPC.ReviewDetail{
		ID:            data.GetID(),
		ReviewAVG:     data.ReviewAVG,
		ProfileID:     detail.GetProfileID(),
		IsCompany:     detail.IsCompany,
		Clarity:       detail.Clarity,
		Communication: detail.Communication,
		CreatedAt:     detail.CreatedAt.String(),
		Description:   detail.Description,
		Hire:          reviewHireToRPC(detail.Hire),
		Payment:       detail.Payment,
	}

	// Get Service
	if data.GetServiceID() != "" {
		res.Service = servicesToServiceRPCService(data.Service)
	}

	// Get Reqeust
	if data.GetRequestID() != "" {
		res.Request = serviceRequestToRPC(&data.Request)
	}

	return res

}

func nullToString(data *string) string {
	if data == nil {
		return ""
	}

	return *data
}

func proposalToOrderType(data *offer.Proposal) *serviceorder.Order {
	if data == nil {
		return nil
	}

	detail := data.ProposalDetail

	res := &serviceorder.Order{
		IsOwnerCompany: detail.IsCompany,
		OrderDetail: serviceorder.OrderDetail{
			IsCompany:      data.IsOwnerCompany,
			Status:         "in_progress",
			Description:    detail.Message,
			PriceType:      detail.PriceType,
			PriceAmount:    detail.PriceAmount,
			MinPriceAmount: detail.MinPriceAmount,
			MaxPriceAmount: detail.MaxPriceAmount,
			Currency:       detail.Currency,
			DeliveryTime:   detail.DeliveryTime,
			CreatedAt:      time.Now(),
		},
		OrderType: "seller",
	}

	// Set Owner Profile id , switch from offer detail to order owner
	if !detail.ProfileID.IsZero() {
		res.SetOwnerID(detail.GetProfileID())
	}

	// Set Office ID
	if !detail.OfficeID.IsZero() {
		res.SetOfficeID(detail.GetOfficeID())
	}

	// Set Service Id
	if !detail.ServiceID.IsZero() {
		res.SetServiceID(detail.GetServiceID())
	}

	// Set Request Id
	if !data.RequestID.IsZero() {
		res.SetRequestID(data.GetRequestID())
	}

	// Set Profile id
	if !data.OwnerID.IsZero() {
		res.OrderDetail.SetProfileID(data.GetOwnerID())
	}

	res.GenerateID()

	return res
}
