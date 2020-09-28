package additionaldetails

import (
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/qualifications"
)

// AdditionalDetails ...
type AdditionalDetails struct {
	Qualifications  *qualifications.Qualifications `bson:"qualifications"`
	// Purpose         *Purpose                       `bson:"purpose"`
	// ServiceIncludes *ServiceIncludes               `bson:"service_includes"`
	ServiceProvider *ServiceProvider               `bson:"service_provider"`

}
