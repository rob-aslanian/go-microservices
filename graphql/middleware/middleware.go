package middleware

import (
	"context"
	"net"
	"net/http"

	"strings"

	"google.golang.org/grpc/metadata"
)

// GetInfo returns context ready to passing through gRPC with metadata:
// ip, ui_lang and etc
func GetInfo(ctx context.Context, r *http.Request) context.Context {
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
