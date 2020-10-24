package websocket

import (
	"dataCenter/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

//1: PeopleAwareness
//2: PersonAwareness
//3: EnvironmentAwareness
//4: EquipmentAwareness

type WsClients struct {
	members    map[*websocket.Conn]interface{}
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn

	requestForPeople chan *models.WsRequestForPeople
	requestForPerson chan *models.WsRequestForPerson

	peopleMembers      map[int]map[*websocket.Conn]bool
	personMembers      map[string]map[*websocket.Conn]bool
	environmentMembers map[int]map[*websocket.Conn]bool
	equipmentMembers   map[int]map[*websocket.Conn]bool

	PeopleBroadcast      chan *models.PeopleAwareness
	PersonBroadcast      chan []*models.PersonAwareness
	environmentBroadcast chan map[int][]byte
	equipmentBroadcast   chan map[int][]byte
}

func NewWsClients() *WsClients {
	return &WsClients{
		members:    make(map[*websocket.Conn]interface{}),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),

		requestForPeople: make(chan *models.WsRequestForPeople),
		requestForPerson: make(chan *models.WsRequestForPerson),

		peopleMembers:      make(map[int]map[*websocket.Conn]bool),
		personMembers:      make(map[string]map[*websocket.Conn]bool),
		environmentMembers: make(map[int]map[*websocket.Conn]bool),
		equipmentMembers:   make(map[int]map[*websocket.Conn]bool),

		PeopleBroadcast:      make(chan *models.PeopleAwareness),
		PersonBroadcast:      make(chan []*models.PersonAwareness),
		environmentBroadcast: make(chan map[int][]byte),
		equipmentBroadcast:   make(chan map[int][]byte),
	}
}

func (wsClients *WsClients) Start() {
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
		}
	}
}
