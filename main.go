package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BeanLiu1994/tiny_chat/chat"
	"github.com/BeanLiu1994/tiny_chat/client/client"
	"github.com/BeanLiu1994/tiny_chat/ws"

	"github.com/gin-gonic/gin"

	hideMyAssParsing "github.com/veksa/hide-my-ass-parsing"
)

func main() {
	// log flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// keepalive
	done := chat.DefaultChatManager.KeepAlive()
	defer close(done)

	// self activator
	serveUrl := os.Getenv("SERVE_URL")
	if serveUrl != "" {
		done2 := make(chan bool)
		go func(done chan bool) {
			ticker := time.NewTicker(10 * time.Minute)
			go func() {
				for {
					select {
					case <-done:
						return
					case <-ticker.C:
						http.Get(serveUrl)
					}
				}
			}()
		}(done2)
		defer close(done2)
	}

	// port env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// prepare client
	_, found := os.LookupEnv("ENABLE_BOT")
	if found {
		done3 := startDefaultClient(port)
		defer close(done3)
	}

	// set server
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/hma", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"proxies": hideMyAssParsing.GetProxies(),
		})
	})
	r.GET("/chat", ws.WsChat)
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/static/", http.StatusTemporaryRedirect)
	})
	r.Run(":" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func startDefaultClient(port string) chan os.Signal {
	done := make(chan os.Signal)
	go func() {
		time.Sleep(time.Second * 5)
		// add client
		client.StartClient("ws://127.0.0.1:"+port+"/chat", "gpt_proxy", done)
	}()
	return done
}
