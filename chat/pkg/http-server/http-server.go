package http_server

import (
	"mime/multipart"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/tracer"
	"golang.org/x/net/context"
)

type AuthClient interface {
	GetUserId(context.Context, string) (string, error)
}

type ConnectionsHolder interface {
	SetConnection(ctx context.Context, userId string, conn *websocket.Conn)
}

type ChatService interface {
	UploadFile(ctx context.Context, header *multipart.FileHeader) (string, error)
	ReadFile(ctx context.Context, fileId string) (*mgo.GridFile, error)

	AuthenticateUser(ctx context.Context) string
	RequireAdminLevelForCompany(ctx context.Context, companyKey string, levels ...string) string
}

type HttpServer struct {
	tracer            *tracing.Tracer
	router            *mux.Router
	authRepo          AuthClient
	service           ChatService
	connectionsHolder ConnectionsHolder

	upgrader *websocket.Upgrader
}

func NewHttpServer(auth AuthClient, chatService ChatService, connections ConnectionsHolder, tracer *tracing.Tracer) *HttpServer {
	server := &HttpServer{
		tracer:            tracer,
		router:            mux.NewRouter(),
		authRepo:          auth,
		service:           chatService,
		connectionsHolder: connections,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	server.registerHandlers()

	return server
}

func (s *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}
