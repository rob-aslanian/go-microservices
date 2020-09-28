package http_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func (s *HttpServer) registerHandlers() {
	s.router.HandleFunc("/ws/messenger", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		span := s.tracer.MakeSpan(ctx, "ws")
		defer span.Finish()

		var senderId string
		var err error
		companyId := request.URL.Query().Get("company_id")
		if companyId == "" {
			senderId, err = s.authenticateUser(request)
			if err != nil {
				writer.WriteHeader(401)
				return
			}
		} else {
			_, err = s.requireAdminLevelForCompany(request, companyId, "Admin")
			if err != nil {
				writer.WriteHeader(401)
				return
			}
			senderId = companyId
		}

		conn, err := s.upgrader.Upgrade(writer, request, nil)
		if err != nil {
			s.tracer.LogError(span, err)
			writer.WriteHeader(400)
			return
		}

		s.connectionsHolder.SetConnection(ctx, senderId, conn)
	})

	s.router.HandleFunc("/fs", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("panic: ", e)
				writer.WriteHeader(500)
			}
		}()
		token := getToken(request)
		ctx := context.Background()
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token", token))

		err := request.ParseMultipartForm(1 << 32)
		if err != nil {
			log.Println("parsing mulitpart form error:", err)
			writer.WriteHeader(400)
			return
		}

		//conversationId := request.FormValue("conversation_id")
		//fmt.Println("conversation id: ", conversationId)
		//if conversationId == "" {
		//	writer.WriteHeader(400)
		//	writer.Write([]byte("conversation_id should not be empty"))
		//}

		files := make([]interface{}, 0, len(request.MultipartForm.File))

		// _, header, err := request.FormFile("file")
		// if err != nil {
		// 	writer.WriteHeader(400)
		// 	return
		// }

		log.Println("request.MultipartForm.File:", request.MultipartForm.File)

		for _, f := range request.MultipartForm.File {
			for _, header := range f {
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}

				if header.Size >= (10485760) { // 10 MB
					writer.WriteHeader(http.StatusBadRequest)
					return
				}

				log.Println("Uploading:", header.Filename)

				id, err := s.service.UploadFile(ctx, header)
				if err != nil {
					fmt.Println(err)
					writer.WriteHeader(500)
					return
				}

				files = append(files, map[string]interface{}{
					"id":   id,
					"name": header.Filename,
					"size": header.Size,
				})
			}
		}
		writer.WriteHeader(201)
		writer.Header().Set("Content-Type", "application/json")

		// response := map[string]interface{}{
		// 	"id":   id,
		// 	"name": header.Filename,
		// 	"size": header.Size,
		// }

		bytes, err := json.Marshal(map[string]interface{}{
			"files": files,
		})
		if err != nil {
			log.Println("error marshaling json:", err)
		}
		log.Println("files:", files)

		writer.Write(bytes)
	}).Methods("POST")

	s.router.HandleFunc("/fs/{fileId}", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("panic: ", e)
				writer.WriteHeader(500)
			}
		}()
		token := getToken(request)
		ctx := context.Background()
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token", token))

		vars := mux.Vars(request)

		if bson.IsObjectIdHex(vars["fileId"]) {
			file, err := s.service.ReadFile(ctx, vars["fileId"])
			if err != nil {
				writer.WriteHeader(404)
				return
			}
			defer file.Close()

			if t := getContentType(file.Name()); t != "" {
				writer.Header().Set("content-type", t)
			}
			if isAttachment(file.Name()) {
				writer.Header().Set("content-disposition", "attachment; filename="+file.Name())
			}
			http.ServeContent(writer, request, file.Name(), file.UploadDate(), file)
		} else {
			http.Redirect(writer, request, "/file/"+vars["fileId"], http.StatusMovedPermanently)
		}
	})
}
