package main

import (
	"hello_server/ws"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// set server
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/chat", ws.WsChat)
	r.GET("/echo", ws.Echo)
	r.GET("/home", ws.Home)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
