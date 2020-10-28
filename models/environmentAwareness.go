package models

type EnvironmentAwareness struct {
	DataType    int     `json:"data_type"`
	SensorID    int     `json:"sensor_id"`
	Location    int     `json:"location"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	UpdateTime  int64   `json:"update_time"`
}

func NewEnvironmentAwareness(sensorID int, location int, temperature float32, humidity float32, updateTime int64) *EnvironmentAwareness {
	return &EnvironmentAwareness{
		DataType:    3,
		SensorID:    sensorID,
		Location:    location,
		Temperature: temperature,
		Humidity:    humidity,
		UpdateTime:  updateTime,
	}
}
