package resolver

import (
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
)

const (
	ContentTypeImage                = "ContentTypeImage"
	ContentTypeArticle              = "ContentTypeArticle"
	ContentTypeCode                 = "ContentTypeCode"
	ContentTypeVideo                = "ContentTypeVideo"
	ContentTypeAudio                = "ContentTypeAudio"
	ContentTypeOther                = "ContentTypeOther"
	up_To_24_Hours                  = "Up_To_24_Hours"
	up_To_3_Days                    = "Up_To_3_Days"
	up_To_7_Days                    = "Up_To_7_Days"
	up_To_1_2_Weeks                 = "Weeks_1_2"
	Weeks_2_4                       = "Weeks_2_4"
	Purpose_Promotional             = "Purpose_Promotional"
	Purpose_Real_Estate             = "Purpose_Real_Estate"
	Purpose_Corporate               = "Purpose_Corporate"
	Purpose_Travel                  = "Purpose_Travel"
	Purpose_Medical                 = "Purpose_Medical"
	Purpose_Academic                = "Purpose_Academic"
	Purpose_Event                   = "Purpose_Event"
	Price_Fixed                     = "Price_Fixed"
	Price_Hourly                    = "Price_Hourly"
	Price_Negotiable                = "Price_Negotiable"
	Remote_only                     = "Remote_only"
	On_Site_Work                    = "On_Site_Work"
	Service_Includes_Source_File    = "Service_Includes_Source_File"
	Service_Include_Print_Ready     = "Service_Include_Print_Ready"
	Service_Include_Photo_Editing   = "Service_Include_Photo_Editing"
	Service_Include_Custom_Graphics = "Service_Include_Custom_Graphics"
	Service_Include_Stock_Photos    = "Service_Include_Stock_Photos"
	No_Preference                   = "No_Preference"
	Company                         = "Company"
	Freelancer                      = "Freelancer"
	Responsive_Design               = "Responsive_Design"
	Fix_Documentation               = "Fix_Documentation"
	Content_Upload                  = "Content_Upload"
	Design_Customization            = "Design_Customization"
	Browser_Compatibility           = "Browser_Compatibility"
	Include_Source_Code             = "Include_Source_Code"
	Bug_Investigation               = "Bug_Investigation"
	Online_Consultants              = "Online_Consultants"
	Hosting                         = "Hosting"
	Detailed_Code_Comments          = "Detailed_Code_Comments"
	Server_Upload                   = "Server_Upload"
	Search_Engine_Optimization      = "Search_Engine_Optimization"
	Other                           = "Other"
	New_Website                     = "New_Website"
	Website_Refine                  = "Website_Refine"
	One_Time_Project                = "One_Time_Project"
	On_Going_Project                = "On_Going_Project"
	Level_Begginer                  = "Level_Begginer"
	Level_Intermediate              = "Level_Intermediate"
	Level_Advanced                  = "Level_Advanced"
	Level_Master                    = "Level_Master"
)

func files(data *[]string) []*servicesRPC.File {
	if data == nil {
		return nil
	}
	filesArray := make([]*servicesRPC.File, 0, len(*data))
	for i := range *data {
		filesArray = append(filesArray, &servicesRPC.File{
			ID: (*data)[i],
		})
	}
	return filesArray
}

func links(in *[]LinkWithIDInput) []*servicesRPC.Link {
	if in != nil {
		linksArray := make([]*servicesRPC.Link, len(*in))
		for i, link := range *in {
			var id string = ""

			if link.ID != nil {
				id = *link.ID
			}

			linksArray[i] = &servicesRPC.Link{
				ID:  id,
				URL: link.Url,
			}
		}
		return linksArray
	}
	return []*servicesRPC.Link{}
}

func linksString(in []string) []*servicesRPC.Link {
	if in != nil {
		linksArray := make([]*servicesRPC.Link, len(in))
		for i := range in {
			linksArray[i] = &servicesRPC.Link{
				ID: in[i],
			}
		}
		return linksArray
	}
	return []*servicesRPC.Link{}
}

func workingDateToRPC(data *WorkingHour) *servicesRPC.WorkingHour {
	if data == nil {
		return nil
	}

	return &servicesRPC.WorkingHour{
		IsAlwaysOpen: data.Is_always_open,
		WorkingHours: workingHoursToRPC(data.Working_date),
	}
}

func workingHoursToRPC(data []*ServiceWorkingHourInput) []*servicesRPC.WorkingDate {
	if data == nil {
		return nil
	}

	workingDates := make([]*servicesRPC.WorkingDate, 0, len(data))

	for _, wd := range data {
		workingDates = append(workingDates, workingHourToRPC(wd))
	}

	return workingDates
}

func workingHourToRPC(data *ServiceWorkingHourInput) *servicesRPC.WorkingDate {
	if data == nil {
		return nil
	}

	return &servicesRPC.WorkingDate{
		FromDate: data.Hour_from,
		ToDate:   data.Hour_to,
		WeekDays: weekDaysToRPC(data.Week_days),
	}

}

func weekDaysToRPC(data []string) []servicesRPC.WeekDays {
	if len(data) <= 0 {
		return nil
	}

	days := make([]servicesRPC.WeekDays, 0, len(data))

	for _, d := range data {
		days = append(days, weekDayToRPC(d))
	}

	return days
}

func weekDayToRPC(data string) servicesRPC.WeekDays {
	switch data {
	case "monday":
		return servicesRPC.WeekDays_MONDAY
	case "tuesday":
		return servicesRPC.WeekDays_TUESDAY
	case "wednesday":
		return servicesRPC.WeekDays_WEDNESDAY
	case "thursday":
		return servicesRPC.WeekDays_THURSDAY
	case "friday":
		return servicesRPC.WeekDays_FRIDAY
	case "saturday":
		return servicesRPC.WeekDays_SATURDAY
	}

	return servicesRPC.WeekDays_SUNDAY
}

func qualifications(data *Qualifications) *servicesRPC.Qualifications {
	if data != nil {
		return &servicesRPC.Qualifications{
			Languages: languagesArrayToServicesRPCLanguages(data.Languages),
			Skills:    skillsArrayToServicesRPCSkills(data.Skills),
			Tools:     ToolsArrayToServicesRPCTools(data.ToolTechnology),
		}
	}
	return nil
}

func qualificationsInput(data *QualificationsInput) *servicesRPC.Qualifications {
	if data != nil {
		return &servicesRPC.Qualifications{
			// Languages: languagesInputArrayToServicesRPCLanguages(data.Languages),
			Tools:  ToolsInputArrayToServicesRPCTools(data.ToolTechnology),
			Skills: skillsInputArrayToServicesRPCSkills(data.Skills),
		}
	}
	return nil
}

func changeQualificationsInput(data *ChangeQualificationsInput) *servicesRPC.Qualifications {
	if data != nil {
		return &servicesRPC.Qualifications{
			// Languages: changeLanguagesInputArrayToServicesRPCLanguages(data.Languages),
			Skills: changeSkillsInputArrayToServicesRPCSkills(data.Skills),
			Tools:  changeToolsInputArrayToServicesRPCTools(data.ToolTechnology),
		}
	}
	return nil
}

func additionalDetails(data *AdditionalDetailsInput) *servicesRPC.AdditionalDetails {
	if data != nil {
		return &servicesRPC.AdditionalDetails{
			Qualifications: qualificationsInput(data.Qualifications),
			// Purpose:         purpose(data.Purpose),
			// ServiceIncludes: serviceIncludes(data.Service_includes),
		}
	}

	return nil
}

func purpose(data *string) servicesRPC.PurposeEnum {
	if data != nil {
		return stringToservicesRPCPurposeEnum(*data)
	}
	return servicesRPC.PurposeEnum_Purpose_Academic
}

func serviceIncludes(data *string) servicesRPC.ServiceIncludesEnum {
	if data != nil {
		return stringToServicesRPCServiceIncludes(*data)
	}
	return servicesRPC.ServiceIncludesEnum_Service_Includes_Source_File
}

func serviceIncludesRPCToString(data servicesRPC.ServiceIncludesEnum) string {
	serviceIncludes := Service_Include_Custom_Graphics

	switch data {
	case servicesRPC.ServiceIncludesEnum_Service_Include_Photo_Editing:
		return Service_Include_Photo_Editing
	case servicesRPC.ServiceIncludesEnum_Service_Include_Print_Ready:
		return Service_Include_Print_Ready
	case servicesRPC.ServiceIncludesEnum_Service_Includes_Source_File:
		return Service_Includes_Source_File
	case servicesRPC.ServiceIncludesEnum_Service_Include_Stock_Photos:
		return Service_Include_Stock_Photos
	}
	return serviceIncludes
}

func stringToServiceRPCDeliveryTimeEnum(data string) servicesRPC.DeliveryTimeEnum {
	delivery_time := servicesRPC.DeliveryTimeEnum_Delivery_Time_Any

	switch data {
	case up_To_24_Hours:
		delivery_time = servicesRPC.DeliveryTimeEnum_Up_To_24_Hours
	case up_To_3_Days:
		delivery_time = servicesRPC.DeliveryTimeEnum_Up_To_3_Days
	case up_To_7_Days:
		delivery_time = servicesRPC.DeliveryTimeEnum_Up_To_7_Days
	case up_To_1_2_Weeks:
		delivery_time = servicesRPC.DeliveryTimeEnum_Weeks_1_2
	case Weeks_2_4:
		delivery_time = servicesRPC.DeliveryTimeEnum_Weeks_2_4
	case "Custom":
		delivery_time = servicesRPC.DeliveryTimeEnum_Custom
	case "Month_And_More":
		delivery_time = servicesRPC.DeliveryTimeEnum_Month_And_More
	}

	return delivery_time
}

func stringToservicesRPCPurposeEnum(data string) servicesRPC.PurposeEnum {
	purpose := servicesRPC.PurposeEnum_Purpose_Academic

	switch data {
	case Purpose_Promotional:
		return servicesRPC.PurposeEnum_Purpose_Promotional
	case Purpose_Corporate:
		purpose = servicesRPC.PurposeEnum_Purpose_Corporate
	case Purpose_Real_Estate:
		purpose = servicesRPC.PurposeEnum_Purpose_Real_Estate
	case Purpose_Travel:
		purpose = servicesRPC.PurposeEnum_Purpose_Travel
	case Purpose_Medical:
		purpose = servicesRPC.PurposeEnum_Purpose_Medical
	case Purpose_Academic:
		purpose = servicesRPC.PurposeEnum_Purpose_Academic
	case Purpose_Event:
		purpose = servicesRPC.PurposeEnum_Purpose_Event
	}

	return purpose
}

func stringToServicesRPCServiceIncludes(data string) servicesRPC.ServiceIncludesEnum {
	serviceIncludes := servicesRPC.ServiceIncludesEnum_Service_Includes_Source_File

	switch data {
	case Service_Include_Custom_Graphics:
		return servicesRPC.ServiceIncludesEnum_Service_Include_Custom_Graphics
	case Service_Include_Photo_Editing:
		return servicesRPC.ServiceIncludesEnum_Service_Include_Photo_Editing
	case Service_Include_Print_Ready:
		return servicesRPC.ServiceIncludesEnum_Service_Include_Print_Ready
	case Service_Include_Stock_Photos:
		return servicesRPC.ServiceIncludesEnum_Service_Include_Stock_Photos
	}

	return serviceIncludes
}

func stringToServicesRPCContentType(data string) servicesRPC.Portfolio_ContentTypeEnum {
	var content = servicesRPC.Portfolio_Content_Type_Other

	switch data {
	case ContentTypeImage:
		content = servicesRPC.Portfolio_Content_Type_Image
	case ContentTypeArticle:
		content = servicesRPC.Portfolio_Content_Type_Article
	case ContentTypeCode:
		content = servicesRPC.Portfolio_Content_Type_Code
	case ContentTypeVideo:
		content = servicesRPC.Portfolio_Content_Type_Video
	case ContentTypeAudio:
		content = servicesRPC.Portfolio_Content_Type_Audio
	}

	return content
}

func servicesRPCContentTypeToString(data servicesRPC.Portfolio_ContentTypeEnum) string {
	var content = ContentTypeOther

	switch data {
	case servicesRPC.Portfolio_Content_Type_Image:
		content = ContentTypeImage
	case servicesRPC.Portfolio_Content_Type_Article:
		content = ContentTypeArticle
	case servicesRPC.Portfolio_Content_Type_Code:
		content = ContentTypeCode
	case servicesRPC.Portfolio_Content_Type_Video:
		content = ContentTypeVideo
	case servicesRPC.Portfolio_Content_Type_Audio:
		content = ContentTypeAudio
	}

	return content
}

func skillsArrayToServicesRPCSkills(data []Skill) []*servicesRPC.Skill {
	if data == nil {
		return []*servicesRPC.Skill{}
	}

	serviceArray := make([]*servicesRPC.Skill, 0, len(data))

	for _, ser := range serviceArray {
		serviceArray = append(serviceArray, ser)
	}

	return serviceArray
}

func skillsArrayRPCToServicesSkills(data []*servicesRPC.Skill) []Skill {
	if data == nil {
		return []Skill{}
	}

	serviceArray := make([]Skill, 0, len(data))

	for _, ser := range data {
		serviceArray = append(serviceArray, skillRPCToSKill(ser))
	}

	return serviceArray
}

func skillRPCToSKill(data *servicesRPC.Skill) Skill {
	if data == nil {
		return Skill{}
	}

	return Skill{
		ID:    data.ID,
		Skill: data.Skill,
	}
}

func languagesArrayToServicesRPCLanguages(data []Language) []*servicesRPC.Language {

	serviceArray := make([]*servicesRPC.Language, 0, len(data))

	for _, ser := range serviceArray {
		serviceArray = append(serviceArray, ser)
	}

	return serviceArray
}

func languageArrayRPCToLanguages(data []*servicesRPC.Language) []Language {
	// if data == []Language{} {
	// 	return []*servicesRPC.Language{}
	// }

	serviceArray := make([]Language, 0, len(data))

	for _, ser := range data {
		serviceArray = append(serviceArray, languageRPCToLanguage(ser))
	}

	return serviceArray
}

func languageRPCToLanguage(data *servicesRPC.Language) Language {
	if data == nil {
		return Language{}
	}

	return Language{
		ID:       data.GetID(),
		Language: data.Language,
		Rank:     ServicesRPCLevelToString(data.Rank),
	}
}

func skillsInputArrayToServicesRPCSkills(data []*SkillInput) []*servicesRPC.Skill {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.Skill, 0, len(data))

	for _, ser := range data {
		serviceArray = append(serviceArray, skillInputToNullServicesSkill(*ser))
	}

	return serviceArray
}

func skillInputToNullServicesSkill(data SkillInput) *servicesRPC.Skill {
	skill := &servicesRPC.Skill{
		Skill: data.Skill,
	}
	return skill
}

func changeSkillsInputArrayToServicesRPCSkills(data *[]ChangeSkillInput) []*servicesRPC.Skill {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.Skill, 0, len(*data))

	for _, ser := range *data {
		serviceArray = append(serviceArray, changeSkillInputToNullServicesSkill(ser))
	}

	return serviceArray
}

func changeSkillInputToNullServicesSkill(data ChangeSkillInput) *servicesRPC.Skill {
	skill := &servicesRPC.Skill{
		ID:    data.ID,
		Skill: data.Skill,
	}
	return skill
}

func languagesInputArrayToServicesRPCLanguages(data *[]QualificationLanguageInput) []*servicesRPC.Language {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.Language, 0, len(*data))

	for _, ser := range *data {
		serviceArray = append(serviceArray, languageInputToNullServiceslanguage(ser))
	}

	return serviceArray
}

func languageInputToNullServiceslanguage(data QualificationLanguageInput) *servicesRPC.Language {
	language := &servicesRPC.Language{
		Language: data.Language,
		Rank:     stringToServicesRPCLevel(data.Rank),
	}

	return language
}

func changeLanguagesInputArrayToServicesRPCLanguages(data *[]ChangeQualificationLanguageInput) []*servicesRPC.Language {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.Language, 0, len(*data))

	for _, ser := range *data {
		serviceArray = append(serviceArray, changeLanguageInputToNullServiceslanguage(ser))
	}

	return serviceArray
}

func changeLanguageInputToNullServiceslanguage(data ChangeQualificationLanguageInput) *servicesRPC.Language {
	language := &servicesRPC.Language{
		ID:       NullToString(data.ID),
		Language: data.Language,
		Rank:     stringToServicesRPCLevel(data.Rank),
	}
	return language
}

func ToolsInputArrayToServicesRPCTools(data *[]VOfficeToolTechnologyInput) []*servicesRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.ToolTechnology, 0, len(*data))

	for _, ser := range *data {
		serviceArray = append(serviceArray, toolTechnologyInputToNullServicestoolTechnology(ser))
	}

	return serviceArray
}

func toolTechnologyInputToNullServicestoolTechnology(data VOfficeToolTechnologyInput) *servicesRPC.ToolTechnology {
	toolTechnology := &servicesRPC.ToolTechnology{
		ToolTechnology: data.Tool_Technology,
		Rank:           stringToServicesRPCLevel(data.Rank),
	}
	return toolTechnology
}

func changeToolsInputArrayToServicesRPCTools(data *[]ChangeVOfficeToolTechnologyInput) []*servicesRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.ToolTechnology, 0, len(*data))

	for _, ser := range *data {
		serviceArray = append(serviceArray, changeToolTechnologyInputToNullServicestoolTechnology(ser))
	}

	return serviceArray
}

func changeToolTechnologyInputToNullServicestoolTechnology(data ChangeVOfficeToolTechnologyInput) *servicesRPC.ToolTechnology {
	toolTechnology := &servicesRPC.ToolTechnology{
		ID:             data.ID,
		ToolTechnology: data.Tool_Technology,
		Rank:           stringToServicesRPCLevel(data.Rank),
	}
	return toolTechnology
}

// TODO check why putting it directly into array without tranfroming, doesn't throw eerors
func ToolsArrayToServicesRPCTools(data []VOfficeToolTechnology) []*servicesRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.ToolTechnology, 0, len(data))

	for _, ser := range data {
		serviceArray = append(serviceArray, toolsToToolsRPC(ser))
	}

	return serviceArray
}

func toolsToToolsRPC(data VOfficeToolTechnology) *servicesRPC.ToolTechnology {

	return &servicesRPC.ToolTechnology{
		ID:             data.ID,
		ToolTechnology: data.Tool_Technology,
		Rank:           stringToServicesRPCLevel(data.Rank),
	}
}

// TODO check why putting it directly into array without tranfroming, doesn't throw eerors
func ToolsArrayRPCToServicesTools(data []*servicesRPC.ToolTechnology) []VOfficeToolTechnology {
	if data == nil {
		return nil
	}

	serviceArray := make([]VOfficeToolTechnology, 0, len(data))

	for _, ser := range data {
		serviceArray = append(serviceArray, toolsRPCToTools(ser))
	}

	return serviceArray
}

func toolsRPCToTools(data *servicesRPC.ToolTechnology) VOfficeToolTechnology {
	if data == nil {
		return VOfficeToolTechnology{}
	}

	return VOfficeToolTechnology{
		ID:              data.ID,
		Tool_Technology: data.ToolTechnology,
		Rank:            ServicesRPCLevelToString(data.Rank),
	}
}

func stringToServicesRPCLevel(data string) servicesRPC.Level {
	level := servicesRPC.Level_Beginner

	switch data {
	case Level_Intermediate:
		return servicesRPC.Level_Intermediate
	case Level_Advanced:
		return servicesRPC.Level_Advanced
	case Level_Master:
		return servicesRPC.Level_Master
	}

	return level
}

func ServicesRPCLevelToString(data servicesRPC.Level) string {
	level := Level_Begginer

	switch data {
	case servicesRPC.Level_Intermediate:
		return Level_Intermediate
	case servicesRPC.Level_Advanced:
		return Level_Advanced
	case servicesRPC.Level_Master:
		return Level_Master
	}

	return level
}

func ToolsArrayNullToServicesRPCTools(data *[]VOfficeToolTechnologyInput) []*servicesRPC.ToolTechnology {
	if data == nil {
		return nil
	}

	serviceArray := make([]*servicesRPC.ToolTechnology, 0, len(*data))

	for _, ser := range *data {
		serviceArray = append(serviceArray, toolArrayNullToServicesRPCtool(ser))
	}

	return serviceArray
}

func toolArrayNullToServicesRPCtool(data VOfficeToolTechnologyInput) *servicesRPC.ToolTechnology {
	return &servicesRPC.ToolTechnology{
		ToolTechnology: data.Tool_Technology,
		Rank:           stringToServicesRPCLevel(data.Rank),
	}
}

func stringToServiceRPCPriceEnum(data string) servicesRPC.PriceEnum {
	priceEnum := servicesRPC.PriceEnum_Price_Any

	switch data {
	case Price_Hourly:
		priceEnum = servicesRPC.PriceEnum_Price_Hourly
	case Price_Negotiable:
		priceEnum = servicesRPC.PriceEnum_Price_Negotiable
	case Price_Fixed:
		priceEnum = servicesRPC.PriceEnum_Price_Fixed
	}

	return priceEnum
}

func locationToServicesRPCLocation(data *LocationInput) *servicesRPC.Location {
	if data == nil {
		return &servicesRPC.Location{}
	}

	location := servicesRPC.Location{
		City: serviceRPCLocationToCity(data.City),
		Country: &servicesRPC.Country{
			Id: data.Country_id,
		},
	}

	return &location
}

func serviceRPCLocationToCity(data CityInput) *servicesRPC.City {
	// if data == nil {
	// 	return nil
	// }

	city := servicesRPC.City{}

	if data.ID != nil {
		city.Id = *data.ID
	}

	if data.Subdivision != nil {
		city.Subdivision = *data.Subdivision
	}

	if data.City != nil {
		city.Title = *data.City
	}

	return &city
}

func locationTypeToServiceRPCLocationTypeEnum(data *string) servicesRPC.LocationEnum {

	if data == nil {
		return servicesRPC.LocationEnum_Remote_only
	}

	switch *data {
	case "Remote_only":
		return servicesRPC.LocationEnum_Remote_only
	case "On_Site_Work":
		return servicesRPC.LocationEnum_On_Site_Work
	}

	return servicesRPC.LocationEnum_Location_Any
}

func category(data CategoryInput) *servicesRPC.ServiceCategory {
	return &servicesRPC.ServiceCategory{
		Main: data.Main,
		Sub:  data.Sub_category,
	}
}

func serviceRPCCategoryToCategory(data *servicesRPC.ServiceCategory) *Category {
	if data == nil {
		return nil
	}

	c := Category{
		Main:         data.GetMain(),
		Sub_Category: data.GetSub(),
	}

	return &c
}

func serviceRPCCategoryToCity(city *servicesRPC.City) *City {
	if city == nil {
		return nil
	}

	data := City{
		ID:          city.GetId(),
		City:        city.GetTitle(),
		Subdivision: city.GetSubdivision(),
	}

	return &data
}

func serviceRPCCategoryToCountry(country *servicesRPC.Country) *Country {
	if country == nil {
		return nil
	}

	data := Country{
		Country: country.GetId(),
	}

	return &data
}

func serviceRPCCategoryToLocation(data *servicesRPC.Location) *Location {
	if data == nil {
		return &Location{
			City:    &City{},
			Country: &Country{},
		}
	}

	loc := Location{
		City:    serviceRPCCategoryToCity(data.GetCity()),
		Country: serviceRPCCategoryToCountry(data.GetCountry()),
	}

	return &loc
}

func visibiltyTypeToServiceRPCVisibilityType(data string) servicesRPC.VisibilityEnum {
	visibility := servicesRPC.VisibilityEnum_Anyone

	switch data {
	case "Only_RightNao_User":
		visibility = servicesRPC.VisibilityEnum_Only_RightNao_User
	case "Invited_Only":
		visibility = servicesRPC.VisibilityEnum_Invited_Only
	}

	return visibility
}

func requestAdditionalDetails(data *RequestAdditionalDetailsInput) *servicesRPC.RequestAdditionalDetails {
	if data != nil {
		return &servicesRPC.RequestAdditionalDetails{
			Languages:       languagesInputArrayToServicesRPCLanguages(data.Languages),
			Tools:           ToolsArrayNullToServicesRPCTools(data.Tools_Technologies),
			Skills:          skillsInputArrayToServicesRPCSkills(data.Skills),
			ServiceProvider: stringToServicesRPCRequestServiceProvider(data.Service_provider),
			// ServiceIncludes: stringToServicesRPCRequestServiceIncludes(data.Service_Includes),
			// ServiceType:     stringToServicesRPCRequestServiceType(data.Service_Type),
		}
	}

	return nil
}

func stringToServicesRPCRequestServiceProvider(data string) servicesRPC.ServiceProviderEnum {
	service_provider := servicesRPC.ServiceProviderEnum_No_Preference

	switch data {
	case "Company":
		service_provider = servicesRPC.ServiceProviderEnum_Company
	case "Freelancer":
		service_provider = servicesRPC.ServiceProviderEnum_Freelancer
	case "Proffesional":
		service_provider = servicesRPC.ServiceProviderEnum_Professional
	}

	return service_provider
}

func stringToServicesRPCRequestServiceIncludes(data string) servicesRPC.RequestServiceIncludesEnum {
	service_includes := servicesRPC.RequestServiceIncludesEnum_Other

	switch data {
	case Responsive_Design:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Responsive_Design
	case Fix_Documentation:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Fix_Documentation
	case Content_Upload:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Content_Upload
	case Design_Customization:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Design_Customization
	case Browser_Compatibility:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Browser_Compatibility
	case Include_Source_Code:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Include_Source_Code
	case Bug_Investigation:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Bug_Investigation
	case Online_Consultants:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Online_Consultants
	case Hosting:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Hosting
	case Detailed_Code_Comments:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Detailed_Code_Comments
	case Server_Upload:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Server_Upload
	case Search_Engine_Optimization:
		service_includes = servicesRPC.RequestServiceIncludesEnum_Search_Engine_Optimization
	}

	return service_includes
}

func stringToServicesRPCRequestServiceType(data string) servicesRPC.RequestAdditionalDetails_ServiceTypeEnum {
	service_type := servicesRPC.RequestAdditionalDetails_New_Website

	if data == Website_Refine {
		service_type = servicesRPC.RequestAdditionalDetails_Website_Refine
	}

	return service_type
}

func stringToServiceRPCRequestProjectType(data string) servicesRPC.RequestProjectTypeEnum {
	project_type := servicesRPC.RequestProjectTypeEnum_Not_Sure

	switch data {
	case On_Going_Project:
		project_type = servicesRPC.RequestProjectTypeEnum_On_Going_Project
	case One_Time_Project:
		project_type = servicesRPC.RequestProjectTypeEnum_One_Time_Project
	}

	return project_type
}

func stringToServiceRPCRequestProjectTypes(data *[]string) []servicesRPC.RequestProjectTypeEnum {
	if data == nil {
		return nil
	}

	projectTypes := make([]servicesRPC.RequestProjectTypeEnum, 0, len(*data))

	for _, pt := range *data {
		projectTypes = append(projectTypes, stringToServiceRPCRequestProjectType(pt))
	}

	return projectTypes
}

func projectTypesToRPC(data []servicesRPC.RequestProjectTypeEnum) []string {
	if data == nil {
		return nil
	}

	projectTypes := make([]string, 0, len(data))

	for _, pt := range data {
		projectTypes = append(projectTypes, servicesRPCProjectTypeToProjectType(pt))
	}

	return projectTypes
}

func vOfficesRPCToVOffices(data []*servicesRPC.VOffice) []VOffice {
	if data == nil {
		return nil
	}

	offices := make([]VOffice, 0, len(data))

	for _, office := range data {
		offices = append(offices, vOfficeRPCToVOffice(office))
	}

	return offices
}

func vOfficeRPCToVOffice(data *servicesRPC.VOffice) VOffice {
	if data == nil {
		return VOffice{}
	}

	office := VOffice{
		ID:           data.GetID(),
		CompanyID:    data.GetCompanyID(),
		UserID:       data.GetUserID(),
		Name:         data.GetName(),
		Description:  data.GetDescription(),
		Created_at:   data.GetCreatedAt(),
		Location:     locationRPCToLocationProfile(data.GetLocation()),
		Languages:    languageArrayRPCToLanguages(data.GetLanguages()),
		Category:     data.GetCategory(),
		IsMe:         data.GetIsMe(),
		Cover:        data.GetCover(),
		Cover_origin: data.GetCoverOrigin(),
		IsOut:        data.GetIsOut(),
		Return_date:  data.GetReturnDate(),
	}

	return office
}

func locationRPCToLocationProfile(data *servicesRPC.Location) *Location {

	if data == nil {
		return nil
	}

	city := City{
		ID:   data.GetCity().GetId(),
		City: data.GetCity().GetTitle(),
	}

	country := Country{
		Country: data.GetCountry().GetId(),
	}

	location := Location{
		City:    &city,
		Country: &country,
	}

	return &location
}

func portfolioRPCArrayToPortfolioProfileArray(data []*servicesRPC.Portfolio) []PortfolioProfile {
	if data == nil {
		return nil
	}

	portfolio := make([]PortfolioProfile, 0, len(data))
	for _, p := range data {
		portfolio = append(portfolio, portfolioRPCToPortfolioProfile(p))
	}

	return portfolio
}

func portfolioRPCToPortfolioProfile(data *servicesRPC.Portfolio) PortfolioProfile {
	if data == nil {
		return PortfolioProfile{}
	}

	portfolio := PortfolioProfile{
		ID:          data.GetId(),
		UserID:      data.GetUserID(),
		ContentType: servicesRPCContentTypeToString(data.GetContentType()),
		Created_at:  data.GetCreatedAt(),
		Description: data.GetDescription(),
		Files:       make([]File, 0, len(data.GetFiles())),
		Link:        make([]Link, 0, len(data.GetLinks())),
		Title:       data.GetTittle(),
	}

	/// Links
	for _, link := range data.Links {
		portfolio.Link = append(portfolio.Link, Link{
			ID:      link.ID,
			Address: link.URL,
		})
	}

	/// Files
	for _, file := range data.GetFiles() {
		portfolio.Files = append(portfolio.Files, File{
			ID:        file.ID,
			Address:   file.URL,
			Name:      file.Name,
			Mime_type: file.MimeType,
		})
	}

	return portfolio
}

func servicesRPCServicesToServicesResolver(data *servicesRPC.GetVOfficeServicesResponse) *Services {
	if data == nil {
		return nil
	}

	services := &Services{
		Service_amount: data.GetServiceAmount(),
		Services:       serviceRPCservicessArrayToservices(data.GetServices()),
	}

	return services
}

func serviceRPCservicessArrayToservices(data []*servicesRPC.Service) []Service {
	if data == nil {
		return nil
	}

	service := make([]Service, 0, len(data))
	for _, s := range data {
		service = append(service, serviceRPCToService(s))
	}

	return service
}

func serviceRPCToService(data *servicesRPC.Service) Service {
	if data == nil {
		return Service{}
	}

	service := Service{
		ID:                 data.GetID(),
		Additional_details: additionalDetailsRPCToAdditionalDetails(data.AdditionalDetails),
		Description:        data.GetDescription(),
		Title:              data.GetTittle(),
		CompanyID:          data.GetCompanyID(),
		UserID:             data.GetUserID(),
		OfficeID:           data.GetOfficeID(),
		Price:              servicesRPCPriceToPrice(data.GetPrice()),
		Currency:           data.GetCurrency(),
		Files:              servicesRPCFileToFile(data.GetFiles()),
		Delivery_time:      servicesRPCDeliveryTimeToDeliveryTime(data.GetDeliveryTime()),
		Location:           serviceRPCCategoryToLocation(data.GetLocation()),
		Location_type:      servicesRPCLocationTypeToLocationType(data.GetLocationType()),
		Fixed_price_amount: data.GetFixedPriceAmmount(),
		Min_price_amount:   data.GetMinPriceAmmount(),
		Max_price_amount:   data.GetMaxPriceAmmount(),
		Is_Draft:           data.GetIsDraft(),
		Is_Remote:          data.GetIsRemote(),
		Is_Paused:          data.GetIsPaused(),
		Has_liked:          data.GetHasLiked(),
		Wokring_hour:       serviceWorkingHourRPCToWorkingHour(data.GetWorkingDate()),
	}

	if s := serviceRPCCategoryToCategory(data.GetCategory()); s != nil {
		service.Category = *s
	}

	return service
}

func serviceWorkingHourRPCToWorkingHour(data *servicesRPC.WorkingHour) *WorkingHourType {
	if data == nil {
		return nil
	}

	return &WorkingHourType{
		Is_always_open: data.GetIsAlwaysOpen(),
		Working_date:   serviceWorkingDatesRPCToWorkingDates(data.GetWorkingHours()),
	}
}

func serviceWorkingDatesRPCToWorkingDates(data []*servicesRPC.WorkingDate) []*ServiceWorkingHour {
	if data == nil {
		return nil
	}

	dates := make([]*ServiceWorkingHour, 0, len(data))

	for _, d := range data {
		dates = append(dates, serviceWorkingDateRPCToWorkingDate(d))
	}

	return dates
}

func serviceWorkingDateRPCToWorkingDate(data *servicesRPC.WorkingDate) *ServiceWorkingHour {
	if data == nil {
		return nil
	}

	return &ServiceWorkingHour{
		Hour_from: data.GetFromDate(),
		Hour_to:   data.GetToDate(),
		Week_days: serviceWeekDaysRPCToWeekDays(data.GetWeekDays()),
	}
}

func serviceWeekDaysRPCToWeekDays(data []servicesRPC.WeekDays) []string {
	if len(data) <= 0 {
		return nil
	}

	weekDays := make([]string, 0, len(data))

	for _, wd := range data {
		weekDays = append(weekDays, serviceWeekDayRPCToWeekDay(wd))
	}

	return weekDays

}

func serviceWeekDayRPCToWeekDay(data servicesRPC.WeekDays) string {

	switch data {
	case servicesRPC.WeekDays_MONDAY:
		return "monday"
	case servicesRPC.WeekDays_TUESDAY:
		return "tuesday"
	case servicesRPC.WeekDays_WEDNESDAY:
		return "wednesday"
	case servicesRPC.WeekDays_THURSDAY:
		return "thursday"
	case servicesRPC.WeekDays_FRIDAY:
		return "friday"
	case servicesRPC.WeekDays_SATURDAY:
		return "saturday"
	}

	return "sunday"

}

func servicesRequestRPCToServicesRequest(data []*servicesRPC.Request) []ServiceRequest {
	if data == nil {
		return nil
	}

	services := make([]ServiceRequest, 0, len(data))

	for _, s := range data {
		res := serviceRequestRPCToServiceRequest(s)
		services = append(services, *res)
	}

	return services

}

func serviceRequestRPCToServiceRequest(data *servicesRPC.Request) *ServiceRequest {
	if data == nil {
		return nil
	}

	var services = &ServiceRequest{}

	services = &ServiceRequest{
		ID:                 data.GetID(),
		UserID:             data.GetUserID(),
		CompanyID:          data.GetCompanyID(),
		Created_at:         data.GetCreatedAt(),
		Delivery_time:      servicesRPCDeliveryTimeToDeliveryTime(data.GetDeliveryTime()),
		Description:        data.GetDescription(),
		Files:              servicesRPCFileToFile(data.GetFiles()),
		Fixed_price_amount: data.GetFixedPriceAmmount(),
		Location:           serviceRPCCategoryToLocation(data.GetLocation()),
		Location_type:      servicesRPCLocationTypeToLocationType(data.GetLocationType()),
		Max_price_amount:   data.GetMaxPriceAmmount(),
		Min_price_amount:   data.GetMinPriceAmmount(),
		Price:              servicesRPCPriceToPrice(data.GetPrice()),
		Currency:           data.GetCurrency(),
		Title:              data.GetTittle(),
		Status:             servicesRPCStatusToStatus(data.GetStatus()),
		Project_type:       servicesRPCProjectTypeToProjectType(data.GetProjectType()),
		Is_Draft:           data.GetIsDraft(),
		Is_Remote:          data.GetIsRemote(),
		Is_Closed:          data.GetIsClosed(),
		Is_Paued:           data.GetIsPaused(),
		Proposal_amount:    data.GetProposalAmount(),
		Additional_details: serviceRequestAdditionalDetailsRPCToAdditionalDetails(data.GetAdditionalDetails()),
		Custom_date:        data.GetCustomDate(),
		Has_liked:          data.GetHasLiked(),
	}

	if data.GetCategory() != nil {
		var category = serviceRPCCategoryToCategory(data.GetCategory())

		services.Category = *category
	}

	return services
}

func servicesRPCProjectTypeToProjectType(data servicesRPC.RequestProjectTypeEnum) string {
	switch data {
	case servicesRPC.RequestProjectTypeEnum_One_Time_Project:
		return "One_Time_Project"
	case servicesRPC.RequestProjectTypeEnum_On_Going_Project:
		return "On_Going_Project"
	}

	return "Not_Sure"
}

func serviceServiceStatusToRPC(data string) servicesRPC.ServiceStatusEnum {

	switch data {
	case "status_activate":
		return servicesRPC.ServiceStatusEnum_SERVICE_ACTIVE
	case "status_deactivate":
		return servicesRPC.ServiceStatusEnum_SERVICE_DEACTIVATE
	case "status_draft":
		return servicesRPC.ServiceStatusEnum_SERVICE_DRAFT
	case "status_paused":
		return servicesRPC.ServiceStatusEnum_SERVICE_PAUSED
	case "status_closed":
		return servicesRPC.ServiceStatusEnum_SERVICE_CLOSED
	}

	return servicesRPC.ServiceStatusEnum_UNKNOWN_SERVICE_STATUS
}

func serviceServiceStatusRPCToString(data servicesRPC.ServiceStatusEnum) string {
	switch data {
	case servicesRPC.ServiceStatusEnum_SERVICE_ACTIVE:
		return "status_activate"
	case servicesRPC.ServiceStatusEnum_SERVICE_DEACTIVATE:
		return "status_deactivate"
	case servicesRPC.ServiceStatusEnum_SERVICE_DRAFT:
		return "status_draft"
	case servicesRPC.ServiceStatusEnum_SERVICE_PAUSED:
		return "status_paused"
	case servicesRPC.ServiceStatusEnum_SERVICE_CLOSED:
		return "status_closed"
	}

	return "status_unknown"
}

func servicesRPCStatusToStatus(data servicesRPC.StatusEnum) string {
	switch data {
	case servicesRPC.StatusEnum_Status_Closed:
		return "status_closed"
	case servicesRPC.StatusEnum_Status_Draft:
		return "status_draft"
	case servicesRPC.StatusEnum_Status_Paused:
		return "status_paused"
	case servicesRPC.StatusEnum_Status_Pending:
		return "status_pending"
	case servicesRPC.StatusEnum_Status_Rejected:
		return "status_rejected"
	}

	return "status_activated"
}

func servicesRPCVisibilityToVisibility(data servicesRPC.VisibilityEnum) string {
	switch data {
	case servicesRPC.VisibilityEnum_Anyone:
		return "Anyone"
	case servicesRPC.VisibilityEnum_Only_RightNao_User:
		return "Only_RightNao_User"
	}

	return "Invited_Only"
}
func servicesRPCFileToFile(data []*servicesRPC.File) []File {
	if data == nil {
		return nil
	}

	files := make([]File, 0, len(data))
	for _, file := range data {
		files = append(files, File{
			ID:        file.ID,
			Mime_type: file.MimeType,
			Name:      file.Name,
			Address:   file.URL,
		})
	}
	return files
}
func servicesRPCLocationTypeToLocationType(data servicesRPC.LocationEnum) string {

	switch data {
	case servicesRPC.LocationEnum_On_Site_Work:
		return "On_Site_Work"
	case servicesRPC.LocationEnum_Remote_only:
		return "Remote_only"
	}

	return "Location_Any"
}

func servicesRPCDeliveryTimeToDeliveryTime(data servicesRPC.DeliveryTimeEnum) string {
	switch data {
	case servicesRPC.DeliveryTimeEnum_Month_And_More:
		return "Month_And_More"
	case servicesRPC.DeliveryTimeEnum_Up_To_3_Days:
		return "Up_To_3_Days"
	case servicesRPC.DeliveryTimeEnum_Up_To_7_Days:
		return "Up_To_7_Days"
	case servicesRPC.DeliveryTimeEnum_Weeks_1_2:
		return "Weeks_1_2"
	case servicesRPC.DeliveryTimeEnum_Weeks_2_4:
		return "Weeks_2_4"
	case servicesRPC.DeliveryTimeEnum_Up_To_24_Hours:
		return "Up_To_24_Hours"
	case servicesRPC.DeliveryTimeEnum_Custom:
		return "Custom"
	}

	return "Any"
}

func servicesRPCPriceToPrice(data servicesRPC.PriceEnum) string {
	switch data {
	case servicesRPC.PriceEnum_Price_Fixed:
		return "Price_Fixed"
	case servicesRPC.PriceEnum_Price_Hourly:
		return "Price_Hourly"
	case servicesRPC.PriceEnum_Price_Negotiable:
		return "Price_Negotiable"
	}

	return "Any"
}

func serviceRequestAdditionalDetailsRPCToAdditionalDetails(data *servicesRPC.RequestAdditionalDetails) RequestAdditionalDetails {

	if data == nil {
		return RequestAdditionalDetails{}
	}
	return RequestAdditionalDetails{
		Skills:           skillsArrayRPCToServicesSkills(data.GetSkills()),
		Languages:        languageArrayRPCToLanguages(data.GetLanguages()),
		ToolTechnology:   ToolsArrayRPCToServicesTools(data.GetTools()),
		Service_provider: serviceProviderRPCToServiceProvider(data.GetServiceProvider()),
	}
}

func serviceProviderRPCToServiceProvider(data servicesRPC.ServiceProviderEnum) string {
	switch data {
	case servicesRPC.ServiceProviderEnum_Company:
		return "Company"
	case servicesRPC.ServiceProviderEnum_Freelancer:
		return "Freelancer"
	case servicesRPC.ServiceProviderEnum_Professional:
		return "Proffesional"

	}

	return "No_Preference"
}

func additionalDetailsRPCToAdditionalDetails(data *servicesRPC.AdditionalDetails) *AdditionalDetails {
	if data == nil {
		return nil
	}

	return &AdditionalDetails{
		Purpose:          purposeString(data.Purpose),
		Qualifications:   qualificationsRPCToQualifications(data.Qualifications),
		Service_includes: serviceIncludesRPCToString(data.ServiceIncludes),
	}
}

func purposeString(data servicesRPC.PurposeEnum) string {
	return purposeEnumRPCToString(data)
}

func purposeEnumRPCToString(data servicesRPC.PurposeEnum) string {
	purpose := Purpose_Academic

	switch data {
	case servicesRPC.PurposeEnum_Purpose_Promotional:
		return Purpose_Promotional
	case servicesRPC.PurposeEnum_Purpose_Corporate:
		return Purpose_Corporate
	case servicesRPC.PurposeEnum_Purpose_Real_Estate:
		return Purpose_Real_Estate
	case servicesRPC.PurposeEnum_Purpose_Travel:
		return Purpose_Travel
	case servicesRPC.PurposeEnum_Purpose_Medical:
		return Purpose_Medical
	case servicesRPC.PurposeEnum_Purpose_Academic:
		return Purpose_Academic
	case servicesRPC.PurposeEnum_Purpose_Event:
		return Purpose_Event
	}

	return purpose
}

func qualificationsRPCToQualifications(data *servicesRPC.Qualifications) *Qualifications {
	if data != nil {
		return &Qualifications{
			Languages:      languageArrayRPCToLanguages(data.Languages),
			Skills:         skillsArrayRPCToServicesSkills(data.Skills),
			ToolTechnology: ToolsArrayRPCToServicesTools(data.Tools),
		}
	}
	return nil
}

func proposalDetailToRPC(data ProposalDetailInput) *servicesRPC.ProposalDetail {
	return &servicesRPC.ProposalDetail{
		ProfileID:       data.Profile_id,
		IsCompany:       data.Is_company,
		ServiceID:       data.Service_id,
		OfficeID:        data.Office_id,
		Message:         data.Message,
		ExperationTime:  data.Expertaion_time,
		CustomDate:      NullToString(data.Custom_date),
		Currency:        NullToString(data.Currency),
		PriceAmount:     NullToInt32(data.Price_amount),
		MinPriceAmmount: NullToInt32(data.Min_price),
		MaxPriceAmmount: NullToInt32(data.Max_price),
		PriceType:       stringToServiceRPCPriceEnum(data.Price_type),
		DeliveryTime:    stringToServiceRPCDeliveryTimeEnum(data.Delivery_time),
	}
}

func proposalRPCToproposal(data *servicesRPC.ProposalDetail) (so *ProposalDetail) {
	if data == nil {
		return nil
	}

	return &ProposalDetail{
		ID:              data.GetID(),
		Currency:        data.GetCurrency(),
		Delivery_time:   servicesRPCDeliveryTimeToDeliveryTime(data.GetDeliveryTime()),
		Custom_date:     data.GetCustomDate(),
		Expertaion_time: data.GetExperationTime(),
		Message:         data.GetMessage(),
		Price_amount:    data.GetPriceAmount(),
		Min_price:       data.GetMinPriceAmmount(),
		Max_price:       data.GetMaxPriceAmmount(),
		Price_type:      servicesRPCPriceToPrice(data.GetPriceType()),
		Service:         serviceRPCToService(data.GetService()),
		Request:         serviceRequestRPCToServiceRequest(data.GetRequest()),
		Status:          orderStatusRPCToOrderStatus(data.GetOrderStatus()),
	}
}

func orderDetailToRPC(data OrderServiceDetailInput) *servicesRPC.OrderService {
	return &servicesRPC.OrderService{
		ProfileID:       data.Profile_id,
		IsCompany:       data.Is_company,
		MinPriceAmmount: NullToInt32(data.Min_price),
		MaxPriceAmmount: NullToInt32(data.Max_price),
		CustomDate:      NullToString(data.Custom_date),
		Description:     data.Description,
		Currency:        NullToString(data.Currency),
		PriceAmount:     NullToInt32(data.Price_amount),
		PriceType:       stringToServiceRPCPriceEnum(data.Price_type),
		DeliveryTime:    stringToServiceRPCDeliveryTimeEnum(data.Delivery_time),
	}
}

func orderTypeToRPC(data string) servicesRPC.OrderType {
	if data == "seller" {
		return servicesRPC.OrderType_SELLER
	}
	return servicesRPC.OrderType_BUYER
}

func orderStatusToRPC(data string) servicesRPC.OrderStatusEnum {
	switch data {
	case "new":
		return servicesRPC.OrderStatusEnum_Status_New
	case "canceled":
		return servicesRPC.OrderStatusEnum_Status_Canceled
	case "in_progress":
		return servicesRPC.OrderStatusEnum_Status_In_Progress
	case "out_of_schedule":
		return servicesRPC.OrderStatusEnum_Status_Out_Of_Schedule
	case "delivered":
		return servicesRPC.OrderStatusEnum_Status_Delivered
	case "completed":
		return servicesRPC.OrderStatusEnum_Status_Completed
	case "disputed":
		return servicesRPC.OrderStatusEnum_Status_Disputed
	}
	return servicesRPC.OrderStatusEnum_Status_Any
}

func orderStatusRPCToOrderStatus(data servicesRPC.OrderStatusEnum) string {
	switch data {
	case servicesRPC.OrderStatusEnum_Status_Canceled:
		return "canceled"
	case servicesRPC.OrderStatusEnum_Status_In_Progress:
		return "in_progress"
	case servicesRPC.OrderStatusEnum_Status_Out_Of_Schedule:
		return "out_of_schedule"
	case servicesRPC.OrderStatusEnum_Status_Delivered:
		return "delivered"
	case servicesRPC.OrderStatusEnum_Status_Completed:
		return "completed"
	case servicesRPC.OrderStatusEnum_Status_Disputed:
		return "disputed"
	case servicesRPC.OrderStatusEnum_Status_New:
		return "new"
	}
	return "any"
}

func orderRPCToOrder(data *servicesRPC.OrderService) (so *ServiceOrder) {
	if data == nil {
		return nil
	}

	log.Printf("Check request in GRAPHQL %+v", data.GetRequest())

	return &ServiceOrder{
		ID:            data.GetID(),
		Currency:      data.GetCurrency(),
		Delivery_time: servicesRPCDeliveryTimeToDeliveryTime(data.GetDeliveryTime()),
		Description:   data.GetDescription(),
		Files:         servicesRPCFileToFile(data.GetFiles()),
		Note:          data.GetNote(),
		Price_amount:  data.GetPriceAmount(),
		Min_price:     data.GetMinPriceAmmount(),
		Max_price:     data.GetMaxPriceAmmount(),
		Custom_date:   data.GetCustomDate(),
		Price_type:    servicesRPCPriceToPrice(data.GetPriceType()),
		Service:       serviceRPCToService(data.GetService()),
		Request:       serviceRequestRPCToServiceRequest(data.GetRequest()),
		Status:        orderStatusRPCToOrderStatus(data.GetOrderStatus()),
	}
}

func reviewHireRPCToReviewHire(data servicesRPC.HireEnum) string {

	switch data {
	case servicesRPC.HireEnum_Will_Hire:
		return "will_hire"
	case servicesRPC.HireEnum_Not_Hire:
		return "not_hire"
	}

	return "not_answer"
}

func reviewHireToRPC(data string) servicesRPC.HireEnum {

	switch data {
	case "will_hire":
		return servicesRPC.HireEnum_Will_Hire
	case "not_hire":
		return servicesRPC.HireEnum_Not_Hire
	}

	return servicesRPC.HireEnum_Not_Answer
}

func reviewDetailToRPC(data ReviewDetailInput) *servicesRPC.ReviewDetail {
	return &servicesRPC.ReviewDetail{
		ProfileID:     data.Profile_id,
		IsCompany:     data.Is_company,
		Clarity:       int32ToUint32(data.Clarity),
		Communication: int32ToUint32(data.Communication),
		Payment:       int32ToUint32(data.Payment),
		Hire:          reviewHireToRPC(data.Hire),
		Description:   NullToString(data.Description),
	}
}

func reviewAVGRPCToReviewAVG(data *servicesRPC.ServiceReviewAVG) (sa ServiceReviewAvg) {
	if data == nil {
		return
	}

	return ServiceReviewAvg{
		Clarity_avg:       float64(data.GetClarityAVG()),
		Communication_avg: float64(data.GetCommunicationAVG()),
		Payment_avg:       float64(data.GetPaymentAVG()),
	}
}

func reviewsRPCToReviews(data []*servicesRPC.ReviewDetail) []Review {
	if data == nil {
		return nil
	}

	return nil
}

func reviewRPCToReview(data *servicesRPC.ReviewDetail) *Review {
	if data == nil {
		return nil
	}

	res := &Review{
		ID:            data.GetID(),
		Description:   data.GetDescription(),
		Hire:          reviewHireRPCToReviewHire(data.GetHire()),
		Review_at:     data.GetCreatedAt(),
		Clarity:       int32(data.GetClarity()),
		Communication: int32(data.GetCommunication()),
		Payment:       int32(data.GetPayment()),
		Review_avg:    data.GetReviewAVG(),
		Request:       &ServiceRequest{},
		Service:       serviceRPCToService(data.GetService()),
	}

	// Get Request
	if data.GetRequest() != nil {
		res.Request = serviceRequestRPCToServiceRequest(data.GetRequest())
	}

	return res

}

func int32ToUint32(data int32) uint32 {
	if data < 0 {
		return 0
	}

	return uint32(data)
}
