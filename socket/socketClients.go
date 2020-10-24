package socket

import (
	"log"
	"net"
	"time"
)

type SocketClients struct {
	members map[*net.Conn]bool
}

func NewSocketClients() *SocketClients {
	return &SocketClients{
		members: make(map[*net.Conn]bool),
	}
}

func (socketClients *SocketClients) Status() {
	for {
		log.Printf("The number of socket clients: %d", len(socketClients.members))
		time.Sleep(time.Second * 5)
	}
}
