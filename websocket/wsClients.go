package websocket

import (
	"TouchAll-DataCenter/config"
	"TouchAll-DataCenter/models"
	"crypto/md5"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sort"
	"sync"
	"time"
)

type WsClients struct {
	members    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn

	requestForPeople chan *models.WsRequestForPeople
	requestForPerson chan *models.WsRequestForPerson
	//requestForEquipmentStatus    chan *models.WsRequestForEquipmentStatus
	requestForWsConnectionStatus   chan *models.WsRequestForWsConnectionStatus
	requestForEquipmentGroupStatus chan *models.WsRequestForEquipmentGroupStatus // 对设备组的请求

	peopleMembers             map[int]map[*websocket.Conn]bool
	personMembers             map[string]map[*websocket.Conn]bool
	wsConnectionStatusMembers map[*websocket.Conn]bool
	//equipmentStatusMembers map[int]map[*websocket.Conn]bool
	// 请求设备组的连接集合:
	// 每个请求独立，用以统计websocket连接数
	equipmentGroupStatusIndividualMembers map[*websocket.Conn]bool
	// 合并对相同设备状态发起的请求
	equipmentGroupStatusCombinedMembers map[[md5.Size]byte]map[*websocket.Conn]bool
	// 对相同设备请求[equipmentID1, equipmentID2, ...]的md5值到[id1, id2, ...]的映射
	equipmentGroupStatusMembersToIDs map[[md5.Size]byte][]int

	PeopleBroadcast          chan *models.PeopleAwareness
	PersonBroadcast          chan []*models.PersonAwareness
	EquipmentStatusBroadcast chan *models.EquipmentStatusAwareness

	// 设备状态流，将所有在线设备发送的设备状态信息汇总，通过websocket推送
	EquipmentStatusStream *models.EquipmentStatusStream

	//
	wsConnectionStatusStream *models.WsConnectionStatusStream
}

func NewWsClients() *WsClients {
	return &WsClients{
		members:    make(map[*websocket.Conn]bool),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),

		requestForPeople: make(chan *models.WsRequestForPeople),
		requestForPerson: make(chan *models.WsRequestForPerson),
		//requestForEquipmentStatus:    make(chan *models.WsRequestForEquipmentStatus),
		requestForWsConnectionStatus:   make(chan *models.WsRequestForWsConnectionStatus),
		requestForEquipmentGroupStatus: make(chan *models.WsRequestForEquipmentGroupStatus),

		peopleMembers:             make(map[int]map[*websocket.Conn]bool),
		personMembers:             make(map[string]map[*websocket.Conn]bool),
		wsConnectionStatusMembers: make(map[*websocket.Conn]bool),
		//equipmentStatusMembers: make(map[int]map[*websocket.Conn]bool),
		equipmentGroupStatusIndividualMembers: make(map[*websocket.Conn]bool),
		equipmentGroupStatusCombinedMembers:   make(map[[md5.Size]byte]map[*websocket.Conn]bool),
		equipmentGroupStatusMembersToIDs:      make(map[[md5.Size]byte][]int),

		PeopleBroadcast:          make(chan *models.PeopleAwareness),
		PersonBroadcast:          make(chan []*models.PersonAwareness),
		EquipmentStatusBroadcast: make(chan *models.EquipmentStatusAwareness),

		EquipmentStatusStream:    models.NewEquipmentStatusStream(),
		wsConnectionStatusStream: models.NewWsConnectionStatusStream(),
	}
}

func (wsClients *WsClients) Start() {
	config := config.NewConfig()
	_, pushInterval := config.GetWebSocketConfig()

	go func() {
		for {
			select {
			case member := <-wsClients.Register:
				wsClients.members[member] = true
			case member := <-wsClients.Unregister:
				if _, has := wsClients.members[member]; has {
					wsClients.members[member] = false
					_ = member.Close()
				}

			case wsRequest := <-wsClients.requestForPeople:
				if _, has := wsClients.peopleMembers[wsRequest.CameraID]; !has {
					wsClients.peopleMembers[wsRequest.CameraID] = make(map[*websocket.Conn]bool)
				}
				wsClients.peopleMembers[wsRequest.CameraID][wsRequest.Conn] = true
			case wsRequest := <-wsClients.requestForPerson:
				if _, has := wsClients.personMembers[wsRequest.Name]; !has {
					wsClients.personMembers[wsRequest.Name] = make(map[*websocket.Conn]bool)
				}
				wsClients.personMembers[wsRequest.Name][wsRequest.Conn] = true

			case wsRequest := <-wsClients.requestForWsConnectionStatus:
				wsClients.wsConnectionStatusMembers[wsRequest.Conn] = true

			//case wsRequest := <-wsClients.requestForEquipmentStatus:
			//	if _, has := wsClients.equipmentStatusMembers[wsRequest.EquipmentID]; !has {
			//		wsClients.equipmentStatusMembers[wsRequest.EquipmentID] = make(map[*websocket.Conn]bool)
			//	}
			//	wsClients.equipmentStatusMembers[wsRequest.EquipmentID][wsRequest.Conn] = true
			case wsRequest := <-wsClients.requestForEquipmentGroupStatus:
				wsClients.equipmentGroupStatusIndividualMembers[wsRequest.Conn] = true
				equipmentIDs := wsRequest.EquipmentIDs
				sort.Ints(equipmentIDs)
				bytes, _ := json.Marshal(equipmentIDs)
				key := md5.Sum(bytes)
				if _, has := wsClients.equipmentGroupStatusCombinedMembers[key]; !has {
					wsClients.equipmentGroupStatusCombinedMembers[key] = make(map[*websocket.Conn]bool)
				}
				wsClients.equipmentGroupStatusCombinedMembers[key][wsRequest.Conn] = true
				wsClients.equipmentGroupStatusMembersToIDs[key] = equipmentIDs
			}
		}
	}()
	go func() {
		for {
			select {
			case message := <-wsClients.PeopleBroadcast:
				if members, has := wsClients.peopleMembers[message.CameraID]; has {
					data, _ := json.Marshal(message)
					for member := range members {
						go func(member *websocket.Conn) {
							err := member.WriteMessage(websocket.TextMessage, data)
							if err != nil {
								log.Printf("write errro: %s", err)
								member.Close()
								delete(members, member)
								delete(wsClients.members, member)
							}
						}(member)
					}
				}
			case message := <-wsClients.PersonBroadcast:
				for _, personAwareness := range message {
					if members, has := wsClients.personMembers[personAwareness.Name]; has {
						data, _ := json.Marshal(personAwareness)
						for member := range members {
							go func(member *websocket.Conn) {
								err := member.WriteMessage(websocket.TextMessage, data)
								if err != nil {
									log.Printf("write errro: %s", err)
									member.Close()
									delete(members, member)
									delete(wsClients.members, member)
								}
							}(member)
						}
					}
				}
				//case message := <-wsClients.EquipmentStatusBroadcast:
				//	if members, has := wsClients.equipmentStatusMembers[message.EquipmentID]; has {
				//		data, _ := json.Marshal(message)
				//		for member := range members {
				//			go func(member *websocket.Conn) {
				//				err := member.WriteMessage(websocket.TextMessage, data)
				//				if err != nil {
				//					log.Printf("write errro: %s", err)
				//					member.Close()
				//					delete(members, member)
				//					delete(wsClients.members, member)
				//				}
				//			}(member)
				//		}
				//	}
			}
		}
	}()

	for {
		t := time.NewTimer(time.Duration(int(pushInterval.(float64))) * time.Second)
		<-t.C
		go wsClients.pushEquipmentStatus()
		go wsClients.pushWsConnectionStatus()
	}
}

func (wsClients *WsClients) pushEquipmentStatus() {
	var err error
	for key := range wsClients.equipmentGroupStatusCombinedMembers {
		go func(key [md5.Size]byte) {
			splitEquipmentStatusStream := models.NewEquipmentStatusStream()
			equipmentIDs := wsClients.equipmentGroupStatusMembersToIDs[key]
			wg := sync.WaitGroup{}
			for _, id := range equipmentIDs {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					if equipmentStatusAwareness, has := wsClients.EquipmentStatusStream.StatusStreamSyncMap.Load(id); has {
						splitEquipmentStatusStream.StatusStream[id] = equipmentStatusAwareness.(models.EquipmentStatusAwareness)
					}
				}(id)
			}
			wg.Wait()
			data, _ := json.Marshal(splitEquipmentStatusStream)
			for conn := range wsClients.equipmentGroupStatusCombinedMembers[key] {
				go func(conn *websocket.Conn, data []byte, key [md5.Size]byte) {
					err = conn.WriteMessage(websocket.TextMessage, data)
					if err != nil {
						log.Printf("write errro: %s", err)
						conn.Close()
						delete(wsClients.members, conn)
						delete(wsClients.equipmentGroupStatusCombinedMembers[key], conn)
						delete(wsClients.equipmentGroupStatusIndividualMembers, conn)
						if len(wsClients.equipmentGroupStatusCombinedMembers[key]) == 0 {
							delete(wsClients.equipmentGroupStatusCombinedMembers, key)
							delete(wsClients.equipmentGroupStatusMembersToIDs, key)
						}
					}
				}(conn, data, key)
			}
		}(key)
	}
}

func (wsClients *WsClients) pushWsConnectionStatus() {
	wsClients.wsConnectionStatusStream.AllConnections = len(wsClients.members)
	wsClients.wsConnectionStatusStream.EquipmentStatusConnections = len(wsClients.equipmentGroupStatusIndividualMembers)
	wsClients.wsConnectionStatusStream.WsConnectionStatusConnections = len(wsClients.wsConnectionStatusMembers)
	wsClients.wsConnectionStatusStream.UpdatedAt = time.Now().Unix()

	data, _ := json.Marshal(wsClients.wsConnectionStatusStream)
	for conn := range wsClients.wsConnectionStatusMembers {
		go func(conn *websocket.Conn) {
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("write errro: %s", err)
				conn.Close()
				delete(wsClients.members, conn)
				delete(wsClients.wsConnectionStatusMembers, conn)
			}
		}(conn)
	}
}

func (wsClients *WsClients) Status() {
	for {
		log.Printf("The number of socket clients: %d", len(wsClients.equipmentGroupStatusIndividualMembers))
		time.Sleep(time.Second * 1)
	}
}
