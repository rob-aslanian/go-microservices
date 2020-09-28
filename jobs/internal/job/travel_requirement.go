package job

type TravelRequirement string

const (
	// TravelRequirementNone ...
	TravelRequirementNone TravelRequirement = "none"

	// TravelRequirementAll ...
	TravelRequirementAll TravelRequirement = "all_time"

	// TravelRequirementWeek ...
	TravelRequirementWeek TravelRequirement = "once_week"

	// TravelRequirementMonth ...
	TravelRequirementMonth TravelRequirement = "once_month"

	// TravelRequirementFew ...
	TravelRequirementFew TravelRequirement = "few_times"

	// TravelRequirementYear ...
	TravelRequirementYear TravelRequirement = "once_year"
)
