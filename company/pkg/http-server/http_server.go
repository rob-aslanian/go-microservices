package serverHTTP

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/tracer"
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
		"/api/activate/company",
		func(writer http.ResponseWriter, request *http.Request) {
			err := s.service.ActivateCompany(
				ctx,
				request.URL.Query().Get("company_id"),
				request.URL.Query().Get("token"),
			)
			if err != nil {
				http.Redirect(writer, request /*request.Host+*/, "/user/login?error="+url.PathEscape(err.Error()), http.StatusTemporaryRedirect) // TODO: change path
				return
			}

			http.Redirect(writer, request /*request.Host+*/, "/user/login", http.StatusTemporaryRedirect) // TODO: change path
		})

	// email activation
	s.router.HandleFunc(
		"/api/activate/company_email",
		func(writer http.ResponseWriter, request *http.Request) {
			err := s.service.ActivateEmail(
				ctx,
				request.URL.Query().Get("company_id"),
				request.URL.Query().Get("code"),
			)
			if err != nil {
				http.Redirect(writer, request /*request.Host+*/, "/user/?error="+url.PathEscape(err.Error()), http.StatusTemporaryRedirect) // where to redirect?
				return
			}

			http.Redirect(writer, request /*request.Host+*/, "/user/login", http.StatusTemporaryRedirect)
		})
}