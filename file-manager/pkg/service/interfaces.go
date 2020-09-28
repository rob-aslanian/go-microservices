package service

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/advertRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/newsfeedRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

type RPCClient = interface {
	Auth() authRPC.AuthServiceClient
	User() userRPC.UserServiceClient
	Company() companyRPC.CompanyServiceClient
	Jobs() jobsRPC.JobsServiceClient
	Advert() advertRPC.AdvertServiceClient
	Services() servicesRPC.ServicesServiceClient
	Newsfeed() newsfeedRPC.NewsfeedServiceClient
	Stuff() stuffRPC.StuffServiceClient
}

type FileListRepository interface {
	SaveFileInfo(multipart.FileHeader, string, string, string) error
	GetFileInfo(context.Context, string) (FileInfo, error)
}

type FileStorageRepository interface {
	Upload(ctx context.Context, request *http.Request) ([]multipart.FileHeader, []string, error)
}

type FileInfo = interface {
	GetID() string
	GetURL() string
	GetName() string
	GetInternalName() string
	GetMimeType() string
	GetSize() int64
	GetOwnerID() string
	GetCreatedAt() time.Time
}
