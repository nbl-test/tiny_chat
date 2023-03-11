package ws

import (
	_ "embed"
	"encoding/json"
	"log"
	"time"

	"github.com/BeanLiu1994/tiny_chat/chat"
	"github.com/BeanLiu1994/tiny_chat/session"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsSession struct {
	c           *gin.Context
	ws          *websocket.Conn
	sessName    string
	created     time.Time
	sendChannel chan string
}

func (s *WsSession) GetID() string {
	return s.sessName
}
func (s *WsSession) Created() time.Time {
	return s.created
}
func (s *WsSession) Send(b []byte) {
	ch := s.sendChannel
	if ch == nil {
		return
	}
	ch <- string(b)
}
func (s *WsSession) SendString(str string) {
	ch := s.sendChannel
	if ch == nil {
		return
	}
	ch <- str
}
func (s *WsSession) stopChannel() {
	s.sendChannel = nil
}

func WsChat(c *gin.Context) {
	// controls ws connectivity
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade err:", err)
		return
	}
	defer ws.Close()

	// create session name
	sessName := session.DefaultSessionManager.CreateID()

	// prepare send method
	sendChannel := make(chan string, 8)
	defer close(sendChannel)
	go func() {
		for v := range sendChannel {
			log.Printf("says to %v: %v", sessName, v)
			err = ws.WriteMessage(websocket.TextMessage, []byte(v))
			if err != nil {
				log.Println("err send text:", err)
			}
		}
	}()

	// record session
	sess := &WsSession{
		c:           c,
		ws:          ws,
		created:     time.Now(),
		sessName:    sessName,
		sendChannel: sendChannel,
	}
	defer sess.stopChannel()
	defer session.DefaultSessionManager.Set(sessName, nil)
	session.DefaultSessionManager.Set(sessName, sess)
	log.Println(sessName, "connected")
	defer log.Println(sessName, "disconnected")

	// init send global msg
	chat.DefaultChatManager.OnConnect(sessName)

	// read loop
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt != websocket.TextMessage {
			continue
		}
		type tmpStruct struct {
			To      string `json:"to"`
			Content string `json:"content"`
		}
		tmp := tmpStruct{}
		err = json.Unmarshal(message, &tmp)
		if err != nil {
			log.Println("parse:", err)
			session.SendJsonErr(sess, "unmarshal input message failed")
			continue
		}
		// send init messages
		chat.DefaultChatManager.Say(sessName, tmp.To, tmp.Content)
	}
}
