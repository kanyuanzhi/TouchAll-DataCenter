package models

import "github.com/gorilla/websocket"

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

type WsRequestForEquipmentStatus struct {
	WsRequest
	EquipmentID int `json:"equipment_id" bson:"equipment_id"`
}

type WsRequestForEquipmentGroupStatus struct {
	WsRequest
	EquipmentIDs []int `json:"equipment_ids" bson:"equipment_ids"`
}
