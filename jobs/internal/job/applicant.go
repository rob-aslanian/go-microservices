package job

import "gitlab.lan/Rightnao-site/microservices/jobs/internal/candidate"

// Applicant ...
type Applicant struct {
	Application     Application
	CareerInterests *candidate.CareerInterests
}
