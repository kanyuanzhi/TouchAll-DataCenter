package models

// DataType:
// 10: PeopleAwareness
// 11: PersonAwareness
// 20: EnvironmentAwareness
// 30: EquipmentBasicInformationAwareness
// 31: EquipmentStatusAwareness

type EquipmentBasicInformationAwareness struct {
	DataType            int                      `json:"data_type"`
	EquipmentID         int                      `json:"equipment_id"`
	NetworkMac          string                   `json:"network_mac"`
	OperateSystem       string                   `json:"operate_system"`
	NetworkName         string                   `json:"network_name"`
	Platform            string                   `json:"platform"`
	Architecture        string                   `json:"architecture"`
	BootTimeInTimestamp int64                    `json:"boot_time_in_timestamp"`
	BootTimeInString    string                   `json:"boot_time_in_string"`
	User                string                   `json:"user"`
	Host                string                   `json:"host"`
	UpdateTime          int64                    `json:"update_time"`
	CPU                 *CPUBasicInformation     `json:"cpu"`
	Memory              *DiskBasicInformation    `json:"memory"`
	Disk                *MemoryBasicInformation  `json:"disk"`
	Network             *NetworkBasicInformation `json:"network"`
}

type CPUBasicInformation struct {
	CPUCount int `json:"cpu_count"`
}

type DiskBasicInformation struct {
	DiskSize int `json:"disk_size"`
}

type MemoryBasicInformation struct {
	MemorySize int `json:"memory_size"`
}

type NetworkBasicInformation struct {
	NetworkMac int `json:"network_mac"`
	NetworkIP  int `json:"network_ip"`
}

type EquipmentsStatusAwareness struct {
	DataType    int            `json:"data_type"`
	EquipmentID int            `json:"equipment_id"`
	NetworkMac  string         `json:"network_mac"`
	RunningTime int64          `json:"running_time"`
	UpdateTime  int64          `json:"update_time"`
	CPU         *CPUStatus     `json:"cpu"`
	Memory      *MemoryStatus  `json:"memory"`
	Disk        *DiskStatus    `json:"disk"`
	Network     *NetworkStatus `json:"network"`
}

type CPUStatus struct {
	CPUPerUtilization     []float32 `json:"cpu_per_utilization"`
	CPUAverageUtilization float32   `json:"cpu_average_utilization"`
}

type DiskStatus struct {
	DiskUsed      float32 `json:"disk_used"`
	DiskAvailable float32 `json:"disk_available"`
}

type MemoryStatus struct {
	MemoryUtilization float32 `json:"memory_utilization"`
	MemoryAvailable   float32 `json:"memory_available"`
}

type NetworkStatus struct {
	NetworkDropIn           int     `json:"network_drop_in"`
	NetworkDropOut          int     `json:"network_drop_out"`
	NetworkSendGigabytes    float32 `json:"network_send_gigabytes"`
	NetworkReceiveGigabytes float32 `json:"network_receive_gigabytes"`
}
