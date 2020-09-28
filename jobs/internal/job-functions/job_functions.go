package jobShared

type JobFunction string

// none
//   accounting
//   administrative
//   arts_design
//   business_development
//   community_social_services
//   consulting
//   education
//   engineering
//   entrepreneurship
//   finance
//   healthcare_services
//   human_resources
//   information_technology
//   legal
//   marketing
//   media_communications
//   military_protective_services
//   operations
//   product_management
//   program_product_management
//   purchasing
//   quality_assurance
//   real_estate
//   rersearch
//   sales
//   support
const (
	// JobFunctionNone ...
	JobFunctionNone JobFunction = "none"

	JobFunctionAccounting JobFunction = "accounting"

	JobFunctionAdministrative JobFunction = "administrative"

	JobFunctionArts_Design JobFunction = "arts_design"

	JobFunctionBusiness_Development JobFunction = "business_development"

	JobFunctionCommunity_Social_Services JobFunction = "community_social_services"

	JobFunctionConsulting JobFunction = "consulting"

	JobFunctionEducation JobFunction = "education"

	JobFunctionEngineering JobFunction = "engineering"

	JobFunctionEntrepreneurship JobFunction = "entrepreneurship"

	JobFunctionFinance JobFunction = "finance"

	JobFunctionHealthcare_Services JobFunction = "healthcare_services"

	JobFunctionHuman_Resources JobFunction = "human_resources"

	JobFunctionInformation_Technology JobFunction = "information_technology"

	JobFunctionLegal JobFunction = "legal"

	JobFunctionMarketing JobFunction = "marketing"

	JobFunctionMedia_Communications JobFunction = "media_communications"

	JobFunctionMilitary_Protective_Services JobFunction = "military_protective_services"

	JobFunctionOperator JobFunction = "operations"

	JobFunctionProduct_Management JobFunction = "product_management"

	JobFunctionProgram_Product_Management JobFunction = "program_product_management"

	JobFunctionPurchasing JobFunction = "purchasing"

	JobFunctionQuality_Assurance JobFunction = "quality_assurance"

	JobFunctionReal_Estate JobFunction = "real_estate"

	JobFunctionRersearch JobFunction = "rersearch"

	JobFunctionSales JobFunction = "sales"

	JobFunctionSupport JobFunction = "support"
)

func String(e JobFunction) string {
	switch e {
	case JobFunctionAccounting:
		return "accounting"
	case JobFunctionAdministrative:
		return "administrative"
	case JobFunctionArts_Design:
		return "arts_design"
	case JobFunctionBusiness_Development:
		return "business_development"
	case JobFunctionCommunity_Social_Services:
		return "community_social_services"
	case JobFunctionConsulting:
		return "consulting"
	case JobFunctionEducation:
		return "education"
	case JobFunctionEngineering:
		return "engineering"
	case JobFunctionEntrepreneurship:
		return "entrepreneurship"
	case JobFunctionFinance:
		return "finance"
	case JobFunctionHealthcare_Services:
		return "healthcare_services"
	case JobFunctionHuman_Resources:
		return "human_resources"
	case JobFunctionInformation_Technology:
		return "information_technology"
	case JobFunctionLegal:
		return "legal"
	case JobFunctionMarketing:
		return "marketing"
	case JobFunctionMedia_Communications:
		return "media_communications"
	case JobFunctionMilitary_Protective_Services:
		return "military_protective_services"
	case JobFunctionOperator:
		return "operations"
	case JobFunctionProduct_Management:
		return "product_management"
	case JobFunctionProgram_Product_Management:
		return "program_product_management"
	case JobFunctionPurchasing:
		return "purchasing"
	case JobFunctionQuality_Assurance:
		return "quality_assurance"
	case JobFunctionReal_Estate:
		return "real_estate"
	case JobFunctionRersearch:
		return "rersearch"
	case JobFunctionSales:
		return "sales"
	case JobFunctionSupport:
		return "support"
	}

	return "none"
}
