package job

import (
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"
	jobShared "gitlab.lan/Rightnao-site/microservices/jobs/internal/job-functions"
)

// Details ...
type Details struct {
	Title           string
	Country         string
	Region          string
	City            string
	Location        Location     `bson:"-"`
	LocationType    LocationType `bson:"location_type"`
	JobFunctions    []jobShared.JobFunction
	EmploymentTypes []candidate.JobType `bson:"employment_types"`
	Descriptions    []*Description

	Required  ApplicantQualification `bson:"required"`
	Preferred ApplicantQualification `bson:"preferred"`

	Files []File `bson:"files,omitempty"`

	SalaryCurrency        string
	SalaryMin             int32
	SalaryMax             int32
	SalaryInterval        candidate.SalaryInterval
	AddtionalCompensation []AddtionalCompensation `bson:"additional_compensation"`
	AdditionalInfo        AdditionalInfo          `bson:"additional_info"`

	Benefits []Benefit

	NumberOfPositions int32
	PublishDay        int32
	PublishMonth      int32
	PublishYear       int32
	DeadlineDay       int32 `bson:"deadlineday"`
	DeadlineMonth     int32 `bson:"deadlinemonth"`
	DeadlineYear      int32 `bson:"deadlineyear"`
	HiringDay         int32
	HiringMonth       int32
	HiringYear        int32
	CoverLetter       bool   `bson:"cover_letter"`
	WorkRemotly       bool   `bson:"work_remote"`
	HeaderURL         string // what for?
}

// Location ...
type Location struct {
	CityID      string
	CityName    string
	Country     string
	Subdivision string
}
