package http_server

import (
	"context"
	"net/http"

	"gitlab.lan/Rightnao-site/microservices/file-manager/pkg/file_response"
)

type Configuration interface {
	GetAddress() string
}

type Service interface {
	Upload(context.Context, string, string, string, string, string, *http.Request, bool) ([]fileResponse.FileResponse, error)
	GetFile(context.Context, http.ResponseWriter, *http.Request, string) error
}
