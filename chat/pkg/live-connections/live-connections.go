package live_connections

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/tracer"
	"google.golang.org/grpc/metadata"
)

type ChatServiceInterface interface {
	SendMessage(ctx context.Context, message *model.Message) error
}

type OfflineRepo interface {
	SetOffline(id string, isOffline bool) error
	IsOffline(id string) (bool, error)
	DeleteOffline(id string) error
}

type LiveConnections struct {
	tracer *tracing.Tracer

	connections         map[string]*websocket.Conn
	connectionsMux      sync.Mutex
	chatService         ChatServiceInterface
	offlineParticipants map[string]bool
	offlineRepo         OfflineRepo
}

func NewLiveConnections(tracer *tracing.Tracer, repo OfflineRepo) *LiveConnections {
	return &LiveConnections{
		tracer:              tracer,
		connections:         make(map[string]*websocket.Conn),
		offlineParticipants: make(map[string]bool),
		offlineRepo:         repo,
	}
}

func (l *LiveConnections) SetChatService(c ChatServiceInterface) {
	l.chatService = c
}

func (l *LiveConnections) SetConnection(ctx context.Context, senderId string, conn *websocket.Conn) {
	span := l.tracer.MakeSpan(ctx, "SetConnection")
	defer span.Finish()

	l.connectionsMux.Lock()
	l.connections[senderId] = conn
	l.connectionsMux.Unlock()

	go func() {
		wsHandler := func() (connected bool) {
			connected = true
			defer func() {
				e := recover()
				if e != nil {
					log.Println("error in websocket handler: ", e)
					var err error

					switch er := e.(type) {
					case error:
						err = er
						if _, isClosed := err.(*websocket.CloseError); isClosed {
							connected = false
						}
					case string:
						err = errors.New(er)
					default:
						err = errors.New("Undefined error")
						log.Println("recoverd without error: ", e)
					}

					conn.WriteJSON(
						struct {
							Error string `json:"error"`
						}{
							err.Error(),
						},
					)
				}
			}()

			var message = model.Message{}
			err := conn.ReadJSON(&message)
			if err != nil {
				l.tracer.LogError(span, err)
				panic(err)
			}
			ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("user-id", senderId))
			err = l.chatService.SendMessage(ctx, &message)
			if err != nil {
				l.tracer.LogError(span, err)
				panic(err)
			}

			return
		}

		isConnected := true
		for isConnected {
			isConnected = wsHandler()
		}
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("on close: ", code, text)
		l.connectionsMux.Lock()
		defer l.connectionsMux.Unlock()
		delete(l.connections, senderId)
		return nil
	})
}

func (l *LiveConnections) SendTo(ctx context.Context, userId string, data interface{}) bool {
	span := l.tracer.MakeSpan(ctx, "SendTo")
	defer span.Finish()

	conn, ok := l.connections[userId]
	if ok {
		err := conn.WriteJSON(data)
		if err == nil {
			return true
		}
		l.tracer.LogError(span, err)
	}

	return false
}

func (l *LiveConnections) IsLive(userId string) bool {
	_, offline := l.offlineParticipants[userId]
	if offline {
		return false
	}

	_, isOnline := l.connections[userId]

	isSetOffline, err := l.offlineRepo.IsOffline(userId)
	if err != nil {
		log.Println("error: is offline:", err)
	}

	if isSetOffline {
		return false
	}

	return isOnline
}

func (l *LiveConnections) SetParticipantOffline(participantId string, status bool) {
	if status {
		err := l.offlineRepo.SetOffline(participantId, status)
		if err != nil {
			log.Println("error: set offline:", err)
		}
		l.offlineParticipants[participantId] = true
	} else {
		err := l.offlineRepo.DeleteOffline(participantId)
		if err != nil {
			log.Println("error: DeleteOffline:", err)
		}
		delete(l.offlineParticipants, participantId)
	}
}
