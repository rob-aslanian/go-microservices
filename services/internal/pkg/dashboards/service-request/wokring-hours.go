package servicerequest

// WorkingHours ...
type WorkingHours struct {
	IsAlwaysOpen bool          `bson:"is_always_open"`
	WorkingDate  []WorkingDate `bson:"working_date"`
}

// WorkingDate ...
type WorkingDate struct {
	HourFrom string    `bson:"hour_from"`
	HourTo   string    `bson:"hour_to"`
	WeekDays []WeekDay `bson:"week_days"`
}

// WeekDay ...
type WeekDay string

const (
	// WeekDayMonday ...
	WeekDayMonday WeekDay = "monday"
	// WeekDayTuesday ...
	WeekDayTuesday WeekDay = "tuesday"
	// WeekDayWednesday ...
	WeekDayWednesday WeekDay = "wednesday"
	// WeekDayThursday ...
	WeekDayThursday WeekDay = "thursday"
	// WeekDayFriday ...
	WeekDayFriday WeekDay = "friday"
	// WeekDaySaturday ...
	WeekDaySaturday WeekDay = "saturday"
	// WeekDaySunday ...
	WeekDaySunday WeekDay = "sunday"
)
