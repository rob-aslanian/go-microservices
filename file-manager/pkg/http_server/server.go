package http_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

type HttpServer struct {
	router  *mux.Router
	config  Configuration
	service Service
}

func NewHttpServer(ctx context.Context, config Configuration, service Service) *HttpServer {
	r := mux.NewRouter()

	server := &HttpServer{
		service: service,
		config:  config,
		router:  r,
	}
	server.registerHandlers(ctx)

	return server
}

func (s *HttpServer) Start() error {
	server := http.Server{
		Addr:           s.config.GetAddress(),
		Handler:        s,
		ReadTimeout:    15 * time.Minute,
		WriteTimeout:   15 * time.Minute,
		MaxHeaderBytes: 1 << 20,
		// ErrorLog:       log.New(os.Stdout, "http:", log.LstdFlags),
	}
	return server.ListenAndServe()
}

func (s *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *HttpServer) registerHandlers(ctx context.Context) {
	s.router.HandleFunc("/api/v1/uploading/company/{company_id}/{target}/", s.Upload)                              //.Methods("POST")
	s.router.HandleFunc("/api/v1/uploading/company/{company_id}/{target}", s.Upload).Queries("base64", "{base64}") //.Methods("POST")
	s.router.HandleFunc("/api/v1/uploading/company/{company_id}/{target}/{target_id}", s.Upload)                   //.Methods("POST")
	s.router.HandleFunc("/api/v1/uploading/company/{company_id}/{target}/{target_id}/{item_id}", s.Upload)         //.Methods("POST")

	s.router.HandleFunc("/api/v1/uploading/{target}/", s.Upload)                                          //.Methods("POST")
	s.router.HandleFunc("/api/v1/uploading/{target}/{target_id}", s.Upload)                               //.Methods("POST")
	s.router.HandleFunc("/api/v1/uploading/{target}/{target_id}/{item_id}", s.Upload)                     //.Methods("POST")
	s.router.HandleFunc("/api/v1/uploading/{target}/{target_id}", s.Upload).Queries("base64", "{base64}") //.Methods("POST")

	s.router.PathPrefix("/file/{file_url}").HandlerFunc(s.Download)
}

// ------- CORS

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// ------- CORS

func (s *HttpServer) Upload(response http.ResponseWriter, request *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(request.Context(), "Uploading file")
	defer span.Finish()

	setupResponse(&response, request)
	if (*request).Method == "OPTIONS" {
		fmt.Println("request is options")
		return
	}

	// log.Println("request:", request)

	var companyID, token, target, targetID, itemID string
	tokenUser, err := request.Cookie("token_user")
	if err == nil {
		token = tokenUser.Value
	}

	// log.Println("Token:", token)

	vars := mux.Vars(request)
	target = vars["target"]
	targetID = vars["target_id"]
	itemID = vars["item_id"]
	companyID = vars["company_id"]

	// base64 := vars["base64"]
	base64 := request.FormValue("base64")

	// log.Println("values:", base64)
	// log.Println("vars:", vars)

	var isBase64 bool
	if base64 == "true" {
		// log.Println("base64 uploading")
		isBase64 = true
	}

	ctx = getInfo(ctx, request)

	info, err := s.service.Upload(ctx, token, companyID, target, targetID, itemID, request, isBase64)

	// BUG: In case error server just reset connection. Why?
	if err != nil {

		log.Println("error:", err)

		switch err {

		case errors.New("You are not authenticated"):
			fallthrough
		case errors.New("token is empty"):
			http.Error(response, err.Error(), http.StatusUnauthorized)

		default:
			http.Error(response, err.Error(), http.StatusInternalServerError)

		}
		return
	}

	response.WriteHeader(http.StatusOK)
	j, _ := json.Marshal(
		map[string]interface{}{
			"success": true,
			"info":    info,
		},
	)
	response.Write(j)
}

func (s *HttpServer) Download(response http.ResponseWriter, request *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(request.Context(), "Downloading file")
	defer span.Finish()

	vars := mux.Vars(request)
	fileURL := vars["file_url"]

	if fileURL == "" {
		http.Error(response, "file not specified", http.StatusBadRequest)
		return
	}

	err := s.service.GetFile(ctx, response, request, fileURL)
	if err != nil {
		return
	}
}

func getInfo(ctx context.Context, r *http.Request) context.Context {
	var host string

	ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	if len(ips) != 0 {
		host = ips[0]
	} else {
		ip := r.RemoteAddr
		var err error
		host, _, err = net.SplitHostPort(ip)
		if err != nil {
			host = "" // TODO: what to do in this case?
		}
		if host == "" {
			host = "37.232.15.11" // TODO needs to be deleted
		}
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

	companyID := ""
	cookie, err := r.Cookie("company_id")
	if err == nil {
		companyID = cookie.Value
	}

	uiLang := ""
	cookie, err = r.Cookie("selected_lang")
	if err == nil {
		uiLang = cookie.Value
	}

	return metadata.AppendToOutgoingContext(
		ctx,
		"ip", host,
		"http_user_agent", r.UserAgent(),
		"token", token,
		"company_id", companyID,
		"ui_lang", uiLang,
	)
}
