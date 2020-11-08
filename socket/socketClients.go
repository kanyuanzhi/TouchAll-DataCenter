package socket

import (
	"log"
	"net"
	"time"
)

type SocketClients struct {
	members          map[net.Conn]bool
	equipmentMembers map[net.Conn]int // 设备与数据中心的连接到设备ID的映射
}

func NewSocketClients() *SocketClients {
	return &SocketClients{
		members:          make(map[net.Conn]bool),
		equipmentMembers: make(map[net.Conn]int),
	}
}

func (socketClients *SocketClients) Status() {
	for {
		log.Printf("The number of socket clients: %d", len(socketClients.equipmentMembers))
		time.Sleep(time.Second * 1)
	}
}
