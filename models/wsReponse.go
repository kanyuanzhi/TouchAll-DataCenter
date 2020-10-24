package models

// DataType:
// 10: PeopleAwareness
// 11: PersonAwareness
// 20: EnvironmentAwareness
// 30: EquipmentBasicInformationAwareness
// 31: EquipmentStatusAwareness

type WsResponse struct {
	DataType int  `json:"data_type"`
	Success  bool `json:"success"`
}

func NewWsResponse(success bool) *WsResponse {
	return &WsResponse{
		DataType: 0,
		Success:  success,
	}
}
