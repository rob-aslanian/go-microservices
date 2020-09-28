package graphqlws

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"

	"gitlab.lan/Rightnao-site/microservices/graphql/middleware"
	"gitlab.lan/Rightnao-site/microservices/graphql/server_ws/internal/connection"
)

const protocolGraphQLWS = "graphql-ws"

var upgrader = websocket.Upgrader{
	CheckOrigin:  func(r *http.Request) bool { return true },
	Subprotocols: []string{protocolGraphQLWS},
}

// NewHandlerFunc returns an http.HandlerFunc that supports GraphQL over websockets
func NewHandlerFunc(svc connection.GraphQLService, httpHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//---
		ctx := middleware.GetInfo( /*r.Context()*/ context.Background(), r)
		//---
		for _, subprotocol := range websocket.Subprotocols(r) {
			if subprotocol == "graphql-ws" {
				ws, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					return
				}

				if ws.Subprotocol() != protocolGraphQLWS {
					ws.Close()
					return
				}

				go connection.Connect(ctx, ws, svc)
				return
			}
		}

		// Fallback to HTTP
		httpHandler.ServeHTTP(w, r)
	}
}
