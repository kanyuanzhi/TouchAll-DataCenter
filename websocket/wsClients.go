package websocket

import (
	"dataCenter/config"
	"dataCenter/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WsClients struct {
	members    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn

	requestForPeople          chan *models.WsRequestForPeople
	requestForPerson          chan *models.WsRequestForPerson
	requestForEquipmentStatus chan *models.WsRequestForEquipmentStatus

	// 对设备组的请求
	requestForEquipmentGroupStatus chan *models.WsRequestForEquipmentGroupStatus
	// 请求设备组的连接集合
	EquipmentGroupStatusMembers map[*models.WsRequestForEquipmentGroupStatus]bool

	peopleMembers          map[int]map[*websocket.Conn]bool
	personMembers          map[string]map[*websocket.Conn]bool
	equipmentStatusMembers map[int]map[*websocket.Conn]bool

	PeopleBroadcast          chan *models.PeopleAwareness
	PersonBroadcast          chan []*models.PersonAwareness
	EquipmentStatusBroadcast chan *models.EquipmentStatusAwareness

	// 设备状态流，将所有在线设备发送的设备状态信息汇总，通过websocket推送
	EquipmentStatusStream *models.EquipmentStatusStream
}

func NewWsClients() *WsClients {
	return &WsClients{
		members:    make(map[*websocket.Conn]bool),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),

		requestForPeople:          make(chan *models.WsRequestForPeople),
		requestForPerson:          make(chan *models.WsRequestForPerson),
		requestForEquipmentStatus: make(chan *models.WsRequestForEquipmentStatus),

		requestForEquipmentGroupStatus: make(chan *models.WsRequestForEquipmentGroupStatus),
		EquipmentGroupStatusMembers:    make(map[*models.WsRequestForEquipmentGroupStatus]bool),

		peopleMembers:          make(map[int]map[*websocket.Conn]bool),
		personMembers:          make(map[string]map[*websocket.Conn]bool),
		equipmentStatusMembers: make(map[int]map[*websocket.Conn]bool),

		PeopleBroadcast:          make(chan *models.PeopleAwareness),
		PersonBroadcast:          make(chan []*models.PersonAwareness),
		EquipmentStatusBroadcast: make(chan *models.EquipmentStatusAwareness),

		EquipmentStatusStream: models.NewEquipmentStatusStream(),
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
					member.Close()
				}

			case wsRequest := <-wsClients.requestForPeople:
				if _, has := wsClients.peopleMembers[wsRequest.Camera]; !has {
					wsClients.peopleMembers[wsRequest.Camera] = make(map[*websocket.Conn]bool)
				}
				wsClients.peopleMembers[wsRequest.Camera][wsRequest.Conn] = true
			case wsRequest := <-wsClients.requestForPerson:
				if _, has := wsClients.personMembers[wsRequest.Name]; !has {
					wsClients.personMembers[wsRequest.Name] = make(map[*websocket.Conn]bool)
				}
				wsClients.personMembers[wsRequest.Name][wsRequest.Conn] = true
			case wsRequest := <-wsClients.requestForEquipmentStatus:
				if _, has := wsClients.equipmentStatusMembers[wsRequest.EquipmentID]; !has {
					wsClients.equipmentStatusMembers[wsRequest.EquipmentID] = make(map[*websocket.Conn]bool)
				}
				wsClients.equipmentStatusMembers[wsRequest.EquipmentID][wsRequest.Conn] = true
			case wsRequest := <-wsClients.requestForEquipmentGroupStatus:
				wsClients.EquipmentGroupStatusMembers[wsRequest] = true
			}
		}
	}()
	go func() {
		for {
			select {
			case message := <-wsClients.PeopleBroadcast:
				if members, has := wsClients.peopleMembers[message.Camera]; has {
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
			case message := <-wsClients.EquipmentStatusBroadcast:
				wsClients.EquipmentStatusStream.StatusStream[message.EquipmentID] = *message
				if members, has := wsClients.equipmentStatusMembers[message.EquipmentID]; has {
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
			}
		}
	}()

	for {
		t := time.NewTimer(time.Duration(int(pushInterval.(float64))) * time.Second)
		<-t.C
		for wsRequestForEquipmentGroupStatus := range wsClients.EquipmentGroupStatusMembers {
			go func(wsRequestForEquipmentGroupStatus *models.WsRequestForEquipmentGroupStatus) {
				splitEquipmentStatusStream := models.NewEquipmentStatusStream()
				for _, equipmentID := range wsRequestForEquipmentGroupStatus.EquipmentIDs {
					if equipmentStatusAwareness, has := wsClients.EquipmentStatusStream.StatusStream[equipmentID]; has {
						splitEquipmentStatusStream.StatusStream[equipmentID] = equipmentStatusAwareness
					}
				}
				data, _ := json.Marshal(splitEquipmentStatusStream)
				err := wsRequestForEquipmentGroupStatus.Conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Printf("write errro: %s", err)
					wsRequestForEquipmentGroupStatus.Conn.Close()
					delete(wsClients.EquipmentGroupStatusMembers, wsRequestForEquipmentGroupStatus)
				}
			}(wsRequestForEquipmentGroupStatus)
		}
	}
}

func (wsClients *WsClients) Status() {
	for {
		log.Printf("The number of socket clients: %d", len(wsClients.EquipmentGroupStatusMembers))
		time.Sleep(time.Second * 1)
	}
}
