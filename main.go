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

	socketClients := socket.NewSocketClients()
	go socketClients.Status()

	socketServer := socket.NewSocketServer(wsClients, socketClients)
	go socketServer.Start()

	for {
		time.Sleep(time.Second)
		//log.Println(time.Now().Unix())
	}
}
