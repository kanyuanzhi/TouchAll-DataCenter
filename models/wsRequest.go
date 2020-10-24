package models

import "github.com/gorilla/websocket"

// RequestType:
// 10: PeopleAwareness
// 11: PersonAwareness
// 20: EnvironmentAwareness
// 30: EquipmentBasicInformationAwareness
// 31: EquipmentStatusAwareness

type WsRequest struct {
	RequestType int             `json:"request_type" bson:"request_type"`
	Conn        *websocket.Conn `json:"conn" bson:"conn"`
}

type WsRequestForPeople struct {
	WsRequest
	Camera int `json:"camera" bson:"camera"`
}

type WsRequestForPerson struct {
	WsRequest
	Name string `json:"name" bson:"name"`
}
