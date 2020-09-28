package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gitlab.lan/Rightnao-site/microservices/statistics/internal/tracer"
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
	// user action
	s.router.HandleFunc(
		"/api/v1/statistics/user/{type}",
		func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			t := vars["type"]

			j := make(map[string]interface{})

			decoder := json.NewDecoder(request.Body)
			err := decoder.Decode(&j)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			ip, token, ua := getInfo(request)
			j["ip"] = ip
			j["token"] = token
			j["user_agent"] = ua
			j["created_at"] = time.Now()

			err = s.service.SaveUserEvent(ctx, j, t)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

		})

	// company activation
	s.router.HandleFunc(
		"/api/v1/statistics/company/{type}",
		func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			t := vars["type"]

			j := make(map[string]interface{})

			decoder := json.NewDecoder(request.Body)
			err := decoder.Decode(&j)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			ip, token, ua := getInfo(request)
			j["ip"] = ip
			j["token"] = token
			j["user_agent"] = ua
			j["created_at"] = time.Now()

			err = s.service.SaveCompanyEvent(ctx, j, t)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

		})
}

func getInfo(r *http.Request) (string, string, string) {
	var ip string
	if r.Header.Get("X-Forwarded-For") != "" {
		ip = r.Header.Get("X-Forwarded-For")
	} else {
		ip = r.RemoteAddr
	}
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		log.Printf("error: \nremote host:%s \nX-Forwarded-For: %s", r.RemoteAddr, r.Header.Get("X-Forwarded-For"))
		host = "" // TODO: what to do in this case?
	}

	var token string
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		splits := strings.Split(authorization, ":")
		if len(splits) > 0 {
			token = strings.TrimSpace(splits[0])
		}
	} else {
		cookie, err := r.Cookie("token_user")
		if err == nil {
			token = cookie.Value
		}
	}

	return host, token, r.UserAgent()
}
