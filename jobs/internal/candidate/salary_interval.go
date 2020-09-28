package candidate

// SalaryInterval ...
type SalaryInterval string

const (
	// SalaryIntervalUnknown unknwon salary interval
	SalaryIntervalUnknown SalaryInterval = "Unknown"
	// SalaryIntervalHour salary per hour
	SalaryIntervalHour SalaryInterval = "Hour"
	// SalaryIntervalDay salary per day
	// SalaryIntervalMonth salary per month
	SalaryIntervalMonth SalaryInterval = "Month"
	// SalaryIntervalYear salary per year
	SalaryIntervalYear SalaryInterval = "Year"
)

var salaryIntervalHours = map[string]int{
	"Hour":  1,
	"Day":   8,
	"Week":  8 * 5,
	"Month": 8 * 22,
	"Year":  8 * 261,
}

// GetHours returns amount of hours of given interval
func (s SalaryInterval) GetHours() int {
	return salaryIntervalHours[string(s)]
}
