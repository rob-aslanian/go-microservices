package serverHTTP

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcuga/golongpoll"
	"google.golang.org/grpc/metadata"
)

func subscribe(lpManager *golongpoll.LongpollManager, auth AuthRPC, network NetworkRPC) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		token, _ := request.Cookie("token_user")

		// don't allow unauthorized access
		if token == nil {
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := metadata.AppendToOutgoingContext(request.Context(), "token", token.Value)

		userIDFromAuth, err := auth.GetUserID(ctx, token.Value)
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
			log.Println(err)
			return
		}

		vars := mux.Vars(request)
		id := vars["id"]
		targetType := vars["target_type"]

		// golongpoll requer category parameter in URL
		values := request.URL.Query()
		values.Add("category", id)
		request.URL.RawQuery = values.Encode()

		// autharization
		switch targetType {
		case "user":
			if id != userIDFromAuth {
				writer.WriteHeader(http.StatusForbidden)
				return
			}

		case "company":
			isAdmin, err := network.IsAdmin(ctx, id)
			if !isAdmin || err != nil {
				writer.WriteHeader(http.StatusForbidden)
				return
			}

		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		lpManager.SubscriptionHandler(writer, request)
	}
}
