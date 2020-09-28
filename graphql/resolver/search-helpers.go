package resolver

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/searchRPC"
)

func serviceOwnerRPCToServiceOwner(data *string) searchRPC.ServiceOwnerEnum {
	if data == nil {
		return searchRPC.ServiceOwnerEnum_Any_Owner
	}

	switch *data {
	case "Owner_User":
		return searchRPC.ServiceOwnerEnum_Owner_User
	case "Owner_Company":
		return searchRPC.ServiceOwnerEnum_Owner_Company
	}

	return searchRPC.ServiceOwnerEnum_Any_Owner

}

func serviceOwnerRPCToString(data searchRPC.ServiceOwnerEnum) string {

	switch data {
	case searchRPC.ServiceOwnerEnum_Owner_User:
		return "Owner_User"
	case searchRPC.ServiceOwnerEnum_Owner_Company:
		return "Owner_Company"
	}

	return "Any_Owner"

}

func stringTotExperienceEnumRPC(s string) searchRPC.ExperienceEnum {
	switch s {
	case "without_experience":
		return searchRPC.ExperienceEnum_WithoutExperience
	case "less_then_one_year":
		return searchRPC.ExperienceEnum_LessThenOneYear
	case "one_two_years":
		return searchRPC.ExperienceEnum_OneTwoYears
	case "two_three_years":
		return searchRPC.ExperienceEnum_TwoThreeYears
	case "three_five_years":
		return searchRPC.ExperienceEnum_ThreeFiveYears
	case "five_seven_years":
		return searchRPC.ExperienceEnum_FiveSevenyears
	case "seven_ten_years":
		return searchRPC.ExperienceEnum_SevenTenYears
	case "ten_years_and_more":
		return searchRPC.ExperienceEnum_TenYearsAndMore
	}

	return searchRPC.ExperienceEnum_UnknownExperience
}

func searchExperienceEnumToString(s searchRPC.ExperienceEnum) string {
	switch s {
	case searchRPC.ExperienceEnum_WithoutExperience:
		return "without_experience"
	case searchRPC.ExperienceEnum_LessThenOneYear:
		return "less_then_one_year"
	case searchRPC.ExperienceEnum_OneTwoYears:
		return "one_two_years"
	case searchRPC.ExperienceEnum_TwoThreeYears:
		return "two_three_years"
	case searchRPC.ExperienceEnum_ThreeFiveYears:
		return "three_five_years"
	case searchRPC.ExperienceEnum_FiveSevenyears:
		return "five_seven_years"
	case searchRPC.ExperienceEnum_SevenTenYears:
		return "seven_ten_years"
	case searchRPC.ExperienceEnum_TenYearsAndMore:
		return "ten_years_and_more"
	}

	return "experience_unknown"
}

func stringToDateEnumRPC(s string) searchRPC.DatePostedEnum {
	switch s {
	case "past_24_hours":
		return searchRPC.DatePostedEnum_Past24Hours
	case "past_week":
		return searchRPC.DatePostedEnum_PastWeek
	case "past_month":
		return searchRPC.DatePostedEnum_PastMonth
	}

	return searchRPC.DatePostedEnum_Anytime
}

func searchRPCDateEnumToString(s searchRPC.DatePostedEnum) string {
	switch s {
	case searchRPC.DatePostedEnum_Past24Hours:
		return "past_24_hours"
	case searchRPC.DatePostedEnum_PastWeek:
		return "past_week"
	case searchRPC.DatePostedEnum_PastMonth:
		return "past_month"
	}

	return "anytime"
}

func stringToSearchRPC(data string) searchRPC.CompanySizeEnum {
	size := searchRPC.CompanySizeEnum_SIZE_UNDEFINED

	switch data {
	case "size_self_employed":
		size = searchRPC.CompanySizeEnum_SIZE_SELF_EMPLOYED
	case "size_1_10_employees":
		size = searchRPC.CompanySizeEnum_SIZE_1_10_EMPLOYEES
	case "size_11_50_employees":
		size = searchRPC.CompanySizeEnum_SIZE_11_50_EMPLOYEES
	case "size_51_200_employees":
		size = searchRPC.CompanySizeEnum_SIZE_51_200_EMPLOYEES
	case "size_201_500_employees":
		size = searchRPC.CompanySizeEnum_SIZE_201_500_EMPLOYEES
	case "size_501_1000_employees":
		size = searchRPC.CompanySizeEnum_SIZE_501_1000_EMPLOYEES
	case "size_1001_5000_employees":
		size = searchRPC.CompanySizeEnum_SIZE_1001_5000_EMPLOYEES
	case "size_5001_10000_employees":
		size = searchRPC.CompanySizeEnum_SIZE_5001_10000_EMPLOYEES
	case "size_10001_plus_employees":
		size = searchRPC.CompanySizeEnum_SIZE_10001_PLUS_EMPLOYEES
	}

	return size
}

func searchRPCToString(data searchRPC.CompanySizeEnum) string {
	size := "size_unknown"

	switch data {
	case searchRPC.CompanySizeEnum_SIZE_SELF_EMPLOYED:
		size = "size_self_employed"
	case searchRPC.CompanySizeEnum_SIZE_1_10_EMPLOYEES:
		size = "size_1_10_employees"
	case searchRPC.CompanySizeEnum_SIZE_11_50_EMPLOYEES:
		size = "size_11_50_employees"
	case searchRPC.CompanySizeEnum_SIZE_51_200_EMPLOYEES:
		size = "size_51_200_employees"
	case searchRPC.CompanySizeEnum_SIZE_201_500_EMPLOYEES:
		size = "size_201_500_employees"
	case searchRPC.CompanySizeEnum_SIZE_501_1000_EMPLOYEES:
		size = "size_501_1000_employees"
	case searchRPC.CompanySizeEnum_SIZE_1001_5000_EMPLOYEES:
		size = "size_1001_5000_employees"
	case searchRPC.CompanySizeEnum_SIZE_5001_10000_EMPLOYEES:
		size = "size_5001_10000_employees"
	case searchRPC.CompanySizeEnum_SIZE_10001_PLUS_EMPLOYEES:
		size = "size_10001_plus_employees"
	}

	return size
}

func searchRPCCityArrayToCityArray(c []*searchRPC.City) []City {
	if c == nil {
		return nil
	}

	srcArray := make([]City, 0, len(c))

	for _, city := range c {
		srcArray = append(srcArray, City{
			City:        city.GetCity(),
			Country:     city.GetCountry(),
			ID:          city.GetID(),
			Subdivision: city.GetSubdivision(),
		})
	}

	return srcArray
}
