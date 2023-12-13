package ms_handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocket struct {
	OnConnect func(*WebSocketClient)
	OnMessage func(*WebSocketClient, int, []byte)
	OnClose   func(*WebSocketClient)
}

type WebSocketClient struct {
	conn *websocket.Conn
}

func (c *WebSocketClient) Send(message any) error {
	var err error

	switch message.(type) {
	case []byte:
		err = c.conn.WriteMessage(websocket.BinaryMessage, (message).([]byte))
		break
	case string:
		err = c.conn.WriteMessage(websocket.TextMessage, []byte((message).(string)))
		break
	default:
		x, err2 := json.Marshal(message)
		if err2 != nil {
			return err2
		}
		err = c.conn.WriteMessage(websocket.TextMessage, x)
		break
	}

	return err
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

	client := &WebSocketClient{conn: conn}
	if w.OnConnect != nil {
		w.OnConnect(client)
	}

	go func() {
		for {
			if conn == nil {
				break
			}
			kind, msg, err2 := conn.ReadMessage()
			if err2 != nil {
				fmt.Printf("WS ERR: %v\n", err2)
				// log.Println(err2)
				break
			} else {
				log.Printf("Received message: %v %v", kind, msg)
				if w.OnMessage != nil {
					w.OnMessage(client, kind, msg)
				}
			}
		}

		if w.OnClose != nil {
			w.OnClose(client)
		}
	}()
}
