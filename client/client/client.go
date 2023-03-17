package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type MessageHandler func(conn *websocket.Conn, received []byte) error

type MsgType struct {
	From     string `json:"from"`
	To       string `json:"to"`
	SentTime string `json:"sent_time"`
	Content  string `json:"content"`
	// self id
	Name string `json:"name"`
}

var selfName string

func GenerateChatHandler(hdl func(inputMsg string) (string, error)) MessageHandler {
	if hdl == nil {
		return nil
	}
	return func(conn *websocket.Conn, received []byte) error {
		var msg MsgType
		err := json.Unmarshal(received, &msg)
		if err != nil {
			return err
		}
		// msg is id msg
		if msg.Name != "" {
			selfName = msg.Name
			return nil
		}
		if selfName == "" || selfName == msg.From || len(msg.Content) == 0 {
			return nil
		}
		// msg is content msg, and has id info
		sentTime, err := time.Parse(time.RFC3339Nano, msg.SentTime)
		if err != nil {
			return err
		}
		// skip too old msgs
		if time.Since(sentTime) > 10*time.Second {
			return nil
		}
		log.Println("got: ", msg.Content)
		out, err := hdl(msg.Content)
		if err != nil {
			out = fmt.Sprintf("some error happened: %v", err)
		}
		err = conn.WriteJSON(map[string]interface{}{
			"content": out,
			"to":      "", // leaves this empty
		})
		if err != nil {
			return err
		}
		log.Println("send: ", out)
		return nil
	}
}

func defaultHandler(inputMsg string) (string, error) {
	return "hello", nil
}

var handlerManager = map[string]func(inputMsg string) (string, error){
	"": defaultHandler,
}

func setHandler(name string, hdl func(inputMsg string) (string, error)) {
	if hdl == nil {
		return
	}
	if name == "" {
		return
	}
	handlerManager[name] = hdl
}

func GetHandler(name string) func(inputMsg string) (string, error) {
	ret, found := handlerManager[name]
	if !found {
		return defaultHandler
	}
	return ret
}

func StartClient(wsURL, handlerType string, interrupt chan os.Signal) {
	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatal("Invalid WS_URL:", err)
	}
	// Prepare handler
	handler := GenerateChatHandler(GetHandler(handlerType))
	for {
		log.Printf("connecting to %s", u.String())
		// Connect to websocket server
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("dial:", err)
			select {
			case <-interrupt:
				log.Println("interrupt")
				return
			case <-time.After(5 * time.Second): // Wait for 5 seconds before retrying
			}
			continue
		}
		defer c.Close()
		log.Printf("connected to %s", u.String())

		done := make(chan struct{})

		// Read messages from server
		go func() {
			defer close(done)
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				if mt != websocket.TextMessage {
					continue
				}
				// log.Printf("recv: %s", message)

				err = handler(c, message)
				if err != nil {
					log.Println("handler:", err)
					return
				}
			}
		}()

		select {
		case <-done:
			time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
