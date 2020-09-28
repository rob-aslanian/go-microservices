package http_server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func getToken(r *http.Request) string {
	var token string
	token = r.URL.Query().Get("token")
	if token != "" {
		return token
	}
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		splits := strings.Split(authorization, " ")
		if len(splits) > 1 {
			token = splits[1]
		}
	} else {
		cookie, err := r.Cookie("token_user")
		if err == nil {
			token = cookie.Value
		}
	}
	return token
}

func (w *HttpServer) authenticateUser(request *http.Request) (userId string, err error) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println("panic: ", e)
			err = errors.New("You are not authenticated")
		}
	}()

	token := getToken(request)
	log.Println("token: ", token)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("token", token))
	userId = w.service.AuthenticateUser(ctx)

	return
}

func (w *HttpServer) requireAdminLevelForCompany(request *http.Request, companyKey string, levels ...string) (userId string, err error) {
	defer func() {
		e := recover()
		if e != nil {
			switch er := e.(type) {
			case string:
				err = errors.New(er)
			case error:
				err = er
			default:
				err = errors.New("Something wrong in server")
			}
			fmt.Println("panic: ", e)
			fmt.Println("err: ", err)
		}
	}()

	token := getToken(request)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("token", token))
	userId = w.service.RequireAdminLevelForCompany(ctx, companyKey, levels...)

	return
}

func findParticipant(id string, list []*model.Participant) *model.Participant {
	for _, p := range list {
		if p.Id == id {
			return p
		}
	}
	return nil
}

var nonAttachmentRegex = regexp.MustCompile("(?i)\\.(png)|(jpg)$")

func isAttachment(name string) bool {
	return !nonAttachmentRegex.Match([]byte(name))
}

var types map[string]*regexp.Regexp

func getContentType(name string) string {
	if types == nil {
		types = map[string]*regexp.Regexp{
			"application/msword": compileRegexForExtension("doc", "dot"),
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document": compileRegexForExtension("docx", "dotx"),
			"application/vnd.ms-excel": compileRegexForExtension("xls", "xlt", "xla"),
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": compileRegexForExtension("xlsx"),
		}
	}

	bytes := []byte(name)
	for t, reg := range types {
		if reg.Match(bytes) {
			return t
		}
	}
	return ""
}

func compileRegexForExtension(extentions ...string) *regexp.Regexp {
	inBrackets := make([]string, len(extentions))
	for i, e := range extentions {
		inBrackets[i] = "(" + e + ")"
	}
	e := strings.Join(inBrackets, "|")
	return regexp.MustCompile("(?i)\\." + e + "$")
}
