package models

// DataType:
// 10: PeopleAwareness
// 11: PersonAwareness
// 20: EnvironmentAwareness
// 30: EquipmentBasicInformationAwareness
// 31: EquipmentStatusAwareness

type EquipmentAwareness struct {
	DataType    int     `json:"data_type"`
	EquipmentID int     `json:"equipment_id"`
	RunningTime int64   `json:"running_time"`
	CPU         float32 `json:"cpu"`
	Memory      float32 `json:"memory"`
	Disk        float32 `json:"disk"`
	UpdateTime  int64   `json:"update_time"`
}

func NewEquipmentAwareness(equipmentID int, runningTime int64, cpu float32, memory float32, disk float32, updateTime int64) *EquipmentAwareness {
	return &EquipmentAwareness{
		DataType:    4,
		EquipmentID: equipmentID,
		RunningTime: runningTime,
		CPU:         cpu,
		Memory:      memory,
		Disk:        disk,
		UpdateTime:  updateTime,
	}
}
