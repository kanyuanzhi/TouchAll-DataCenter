package models

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
