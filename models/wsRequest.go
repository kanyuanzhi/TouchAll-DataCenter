package models

import "github.com/gorilla/websocket"

type WsRequest struct {
	RequestType int             `json:"request_type" bson:"request_type"`
	Conn        *websocket.Conn `json:"conn" bson:"conn"`
}

type WsRequestForPeople struct {
	WsRequest
	CameraID         int `json:"camera_id" bson:"camera_id"`
	PreviousCameraID int `json:"previous_camera_id"`
}

type WsRequestForPerson struct {
	WsRequest
	Name string `json:"name" bson:"name"`
}

type WsRequestForEquipmentStatus struct {
	WsRequest
	EquipmentID int    `json:"equipment_id" bson:"equipment_id"`
	NetworkMac  string `json:"network_mac" bson:"network_mac"`
}

// 用户端发出的对设备组状态的websocket请求
type WsRequestForEquipmentGroupStatus struct {
	WsRequest
	EquipmentIDs []int `json:"equipment_ids" bson:"equipment_ids"`
}

// 用户端发出的对数据中心websocket连接状态的websocket请求
type WsRequestForWsConnectionStatus struct {
	WsRequest
}

// 用户端发出的对环境状态的websocket请求
type WsRequestForEnvironment struct {
	WsRequest
	SensorIDs []int `json:"sensor_ids" bson:"sensor_ids"`
}
