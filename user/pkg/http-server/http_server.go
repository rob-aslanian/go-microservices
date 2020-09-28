package serverHTTP

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	tracing "gitlab.lan/Rightnao-site/microservices/user/pkg/tracer"
)

// ServerHTTP ...
type ServerHTTP struct {
	address string
	router  *mux.Router
	service Service
	tracer  *tracing.Tracer
}

// Settings ...
type Settings struct {
	Address string
	Service Service
	Tracer  *tracing.Tracer
}

// NewHTTPServer creates new instance of http server
func NewHTTPServer(ctx context.Context, settings Settings) (*ServerHTTP, error) {
	server := &ServerHTTP{
		address: settings.Address,
		service: settings.Service,
		router:  mux.NewRouter(),
		tracer:  settings.Tracer,
	}
	server.registerHandlers(ctx)
	return server, nil
}

// Launch ...
func (s *ServerHTTP) Launch() error {
	server := http.Server{
		Addr:           s.address,
		Handler:        s,
		ReadTimeout:    15 * time.Minute,
		WriteTimeout:   15 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}

func (s *ServerHTTP) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *ServerHTTP) registerHandlers(ctx context.Context) {
	// user activation
	s.router.HandleFunc(
		"/api/activate/user",
		func(writer http.ResponseWriter, request *http.Request) {
			_, err := s.service.ActivateUser(
				ctx,
				request.URL.Query().Get("token"),
				request.URL.Query().Get("user_id"),
			)
			if err != nil {
				http.Redirect(writer, request /*request.Host+*/, "/user/login?error="+url.PathEscape(err.Error()), http.StatusTemporaryRedirect)
				return
			}

			http.Redirect(writer, request /*request.Host+*/, "/user/login", http.StatusTemporaryRedirect)
		})

	// email activation
	s.router.HandleFunc(
		"/api/activate/email",
		func(writer http.ResponseWriter, request *http.Request) {
			err := s.service.ActivateEmail(
				ctx,
				request.URL.Query().Get("code"),
				request.URL.Query().Get("user_id"),
			)
			if err != nil {
				http.Redirect(writer, request /*request.Host+*/, "/user/?error="+url.PathEscape(err.Error()), http.StatusTemporaryRedirect) // where to redirect?
				return
			}

			http.Redirect(writer, request /*request.Host+*/, "/user/login", http.StatusTemporaryRedirect)
		})
}
