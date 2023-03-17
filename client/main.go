package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/BeanLiu1994/tiny_chat/client/client"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	handlerType := os.Getenv("HDL_TYPE")

	// Use environment variable WS_URL as websocket server address
	wsURL, ok := os.LookupEnv("WS_URL")
	if !ok {
		log.Fatal("WS_URL environment variable not set")
	}

	client.StartClient(wsURL, handlerType, interrupt)
}
