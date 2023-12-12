package ms_handler

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocket struct {
	OnConnect func(*WebSocketClient)
}

type WebSocketClient struct {
	conn *websocket.Conn
}

func (w WebSocket) Handle(args *Args) {
	// Устанавливаем соединение веб-сокета
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(args.Response.Response, args.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	w.OnConnect(&WebSocketClient{conn: conn})
}
