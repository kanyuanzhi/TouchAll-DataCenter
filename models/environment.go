package models

import "sync"

type Environment struct {
	DataType    int     `json:"data_type"`
	SensorID    int     `json:"sensor_id"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	UpdatedAt   int64   `json:"updated_at"`
}

// websocket推送流
type EnvironmentStream struct {
	DataType           int                 `json:"data_type"`
	Environment        map[int]Environment `json:"environment"`
	EnvironmentSyncMap sync.Map            `json:"environment_sync_map"`
}

func NewEnvironmentStream() *EnvironmentStream {
	return &EnvironmentStream{
		DataType:    20,
		Environment: make(map[int]Environment),
	}
}
