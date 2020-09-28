package job

// Plan ...
type Plan string

const (
	// PlanUnknown ...
	PlanUnknown Plan = "unknown"
	// PlanBasic ...
	PlanBasic Plan = "basic"
	// PlanStart ...
	PlanStart Plan = "start"
	// PlanStandard ...
	PlanStandard Plan = "standard"
	// PlanProfessional ...
	PlanProfessional Plan = "professional"
	// PlanProfessionalPlus ...
	PlanProfessionalPlus Plan = "professionalPlus"
	// PlanExclusive ...
	PlanExclusive Plan = "exclusive"
	// PlanPremium ...
	PlanPremium Plan = "premium"
)

var jobPlanDays = map[string]int{
	"basic":            5,
	"start":            15,
	"standard":         30,
	"professional":     30,
	"professionalPlus": 30,
	"exclusive":        30,
	"premium":          30,
}
var jobPlanPriorities = map[string]int{
	"basic":            7,
	"start":            6,
	"standard":         5,
	"professional":     4,
	"professionalPlus": 3,
	"exclusive":        2,
	"premium":          1,
}

// GetDays ...
func (p Plan) GetDays() int {
	return jobPlanDays[string(p)]
}

// GetPriority ...
func (p Plan) GetPriority() int {
	return jobPlanPriorities[string(p)]
}
