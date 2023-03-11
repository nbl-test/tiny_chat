package main

import (
	"log"
	"net/http"

	"github.com/BeanLiu1994/tiny_chat/chat"
	"github.com/BeanLiu1994/tiny_chat/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	// log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// keepalive
	done := chat.DefaultChatManager.KeepAlive()
	defer close(done)

	// set server
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/chat", ws.WsChat)
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/static/", http.StatusTemporaryRedirect)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
