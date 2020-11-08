package models

type WsConnectionStatusStream struct {
	DataType                      int   `json:"data_type"`
	AllConnections                int   `json:"all_connections"`
	EquipmentStatusConnections    int   `json:"equipment_status_connections"`
	WsConnectionStatusConnections int   `json:"ws_connection_status_connections"`
	UpdatedAt                     int64 `json:"updated_at"`
}

func NewWsConnectionStatusStream() *WsConnectionStatusStream {
	return &WsConnectionStatusStream{
		DataType:                      40,
		AllConnections:                0,
		EquipmentStatusConnections:    0,
		WsConnectionStatusConnections: 0,
		UpdatedAt:                     0,
	}
}
