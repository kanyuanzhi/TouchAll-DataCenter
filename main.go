package main

import (
	"time"
)

func main() {
	wsClients := NewWsClients()
	go wsClients.Start()

	wsServer := NewWsServer(wsClients)
	go wsServer.Start()

	socketServer := NewSocketServer(wsClients)
	go socketServer.Start()

	for {
		time.Sleep(time.Second)
		//log.Println(time.Now().Unix())
	}
}
