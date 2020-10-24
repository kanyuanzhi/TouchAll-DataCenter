package main

import (
	"dataCenter/socket"
	"dataCenter/websocket"
	"time"
)

func main() {
	wsClients := websocket.NewWsClients()
	go wsClients.Start()

	wsServer := websocket.NewWsServer(wsClients)
	go wsServer.Start()

	socketServer := socket.NewSocketServer(wsClients)
	go socketServer.Start()

	for {
		time.Sleep(time.Second)
		//log.Println(time.Now().Unix())
	}
}
