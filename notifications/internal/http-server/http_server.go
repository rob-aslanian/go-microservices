package serverHTTP

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jcuga/golongpoll"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/mq"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/tracer"
)

// ServerHTTP ...
type ServerHTTP struct {
	address    string
	router     *mux.Router
	service    Service
	tracer     *tracing.Tracer
	longpoll   *golongpoll.LongpollManager
	mq         *mq.RabbitmqConnection
	authRPC    AuthRPC
	networkRPC NetworkRPC
}

// Settings ...
type Settings struct {
	Address         string
	Service         Service
	Tracer          *tracing.Tracer
	LongPollManager *golongpoll.LongpollManager
	AuthRPC         AuthRPC
	NetworkRPC      NetworkRPC
}

// NewHTTPServer creates new instance of http server
func NewHTTPServer(ctx context.Context, settings Settings) (*ServerHTTP, error) {
	server := &ServerHTTP{
		address:    settings.Address,
		service:    settings.Service,
		router:     mux.NewRouter(),
		tracer:     settings.Tracer,
		longpoll:   settings.LongPollManager,
		authRPC:    settings.AuthRPC,
		networkRPC: settings.NetworkRPC,
	}
	server.registerHandlers(ctx)
	return server, nil
}

// Launch ...
func (s *ServerHTTP) Launch() error {
	server := http.Server{
		Addr:           s.address,
		Handler:        s,
		ReadTimeout:    3 * time.Minute,
		WriteTimeout:   3 * time.Minute,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	return server.ListenAndServe()
}

func (s *ServerHTTP) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *ServerHTTP) registerHandlers(ctx context.Context) {
	s.router.HandleFunc(
		"/api/v1/notifications/{target_type}/{id}",
		subscribe(s.longpoll, s.authRPC, s.networkRPC),
	)
}
