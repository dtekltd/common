package ws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dtekltd/common/system"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Map map[string]any

type Message struct {
	To      string `json:"to,omitempty"`
	From    string `json:"from,omitempty"`
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	Data    *Map   `json:"data,omitempty"`
}

var (
	initialized = false
	userConns   = make(map[uint64]*websocket.Conn) // [userID]conn
	clientConns = make(map[uint64]*websocket.Conn) // [clientID]conn
	broadcast   = make(chan []byte)
	register    = make(chan *websocket.Conn)
	unregister  = make(chan *websocket.Conn)
)

func NewHandler() fiber.Handler {
	initialized = true
	return websocket.New(func(conn *websocket.Conn) {
		// When the function returns,
		// unregister the client and close the connection
		defer func() {
			unregister <- conn
			conn.Close()
		}()

		// Register the client
		register <- conn

		for {
			t, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					system.Logger.Error("WS read error:", err.Error())
				}
				return // defer
			}

			if t == websocket.TextMessage {
				// broadcast <- message
				system.Logger.Info("WS message from client:", message)
			} else {
				system.Logger.Error("WS message-type not supported:", t)
			}
		}
	})
}

func RunHub() {
	if !initialized {
		return
	}
	for {
		select {
		case message := <-broadcast:
			handleMessage(message)

		case conn := <-register:
			id := toUint64(conn.Params("id"))
			clientConns[id] = conn
			if uID := conn.Query("id"); uID != "" {
				userConns[toUint64(uID)] = conn
				system.Logger.Infof("WS new connection - cid: %d, uid: %s, name: %s", id, uID, conn.Query("name"))
			} else {
				system.Logger.Info("WS new connection - cid: %d", id)
			}

		case conn := <-unregister:
			id := toUint64(conn.Params("id"))
			delete(clientConns, id)
			if uID := conn.Query("id"); uID != "" {
				delete(userConns, toUint64(uID))
				system.Logger.Infof("WS close connection - cid: %d, uid: %s, name: %s", id, uID, conn.Query("name"))
			} else {
				system.Logger.Infof("WS close connection - cid: %d", id)
			}

			// ...
		}
	}
}

func SendMessage(message *Message) error {
	if !initialized {
		return nil
	}
	if message.To != "" {
		if conn, ok := clientConns[toUint64(message.To)]; !ok {
			return fmt.Errorf("ws connection not found for client #%s", message.To)
		} else {
			message.To = "" // skip to from json
			return conn.WriteMessage(websocket.TextMessage, messageToString(message))
		}
	} else {
		if message.Type == "" {
			// broadcast to all clients
			message.Type = "announce"
		}
		msg, _ := json.Marshal(message)
		broadcast <- msg
	}
	return nil
}

func handleMessage(msg []byte) error {
	message := &Message{}
	if err := json.Unmarshal(msg, message); err != nil {
		system.Logger.Error("WS Unmarshal message failed:", err.Error())
		return err
	}
	switch message.Type {
	case "p2p":
		// send to a specific user
		system.Logger.Error("WS not implemented!")
	case "announce":
		for _, conn := range clientConns {
			go func(conn *websocket.Conn) {
				// send to each client in parallel
				// so we don't block on a slow client
				if err := conn.WriteMessage(websocket.TextMessage, messageToString(message)); err != nil {
					system.Logger.Error("WS write error:", err.Error())
				}
			}(conn)
		}
	default:
		system.Logger.Errorf("WS message TYPE '%s' not supported", message.Type)
	}
	return nil
}

func toUint64(str string) uint64 {
	num, _ := strconv.ParseFloat(str, 64)
	return uint64(num)
}

func messageToString(msg *Message) []byte {
	if text, err := json.Marshal(msg); err != nil {
		system.Logger.Errorf("WS marshal failed msg: %v, err: %s", msg, err)
		return []byte{}
	} else {
		return text
	}
}
