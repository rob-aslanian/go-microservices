package clientRPC

type (
	Configuration interface {
		GetAuthAddress() string
		GetUserAddress() string
		GetCompanyAddress() string
		GetJobAddress() string
		GetAdvertAddress() string
		GetServicesAddress() string
		GetNewsfeedAddress() string
		GetStuffAddress() string
	}
)
