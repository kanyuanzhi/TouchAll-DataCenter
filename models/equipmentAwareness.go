package models

type EquipmentBasicInformationAwarenessMysql struct {
	DataType            int     `json:"data_type" db:"data_type"`
	EquipmentID         int     `json:"equipment_id" db:"equipment_id"`
	OperateSystem       string  `json:"operate_system" db:"operate_system"`
	NetworkName         string  `json:"network_name" db:"network_name"`
	Platform            string  `json:"platform" db:"platform"`
	Architecture        string  `json:"architecture" db:"architecture"`
	BootTimeInTimestamp int64   `json:"boot_time_in_timestamp" db:"boot_time_in_timestamp"`
	BootTimeInString    string  `json:"boot_time_in_string" db:"boot_time_in_string"`
	User                string  `json:"user" db:"user"`
	Host                string  `json:"host" db:"host"`
	UpdateTime          int64   `json:"update_time" db:"update_time"`
	CPUCount            int     `json:"cpu_count" db:"cpu_count"`
	DiskSize            float32 `json:"disk_size" db:"disk_size"`
	MemorySize          float32 `json:"memory_size" db:"memory_size"`
	NetworkMac1         string  `json:"network_mac_1" db:"network_mac_1"`
	NetworkIP1          string  `json:"network_ip_1" db:"network_ip_1"`
	NetworkMac2         string  `json:"network_mac_2" db:"network_mac_2"`
	NetworkIP2          string  `json:"network_ip_2" db:"network_ip_2"`
	NetworkCard1        string  `json:"network_card_1" db:"network_card_1"`
	NetworkCard2        string  `json:"network_card_2" db:"network_card_2"`
	Authenticated       int     `json:"authenticated" db:"authenticated"`
}

type EquipmentBasicInformationAwareness struct {
	DataType            int                     `json:"data_type" db:"data_type"`
	EquipmentID         int                     `json:"equipment_id" db:"equipment_id"`
	OperateSystem       string                  `json:"operate_system" db:"operate_system"`
	NetworkName         string                  `json:"network_name" db:"network_name"`
	Platform            string                  `json:"platform" db:"platform"`
	Architecture        string                  `json:"architecture" db:"architecture"`
	BootTimeInTimestamp int64                   `json:"boot_time_in_timestamp" db:"boot_time_in_timestamp"`
	BootTimeInString    string                  `json:"boot_time_in_string" db:"boot_time_in_string"`
	User                string                  `json:"user" db:"user"`
	Host                string                  `json:"host" db:"host"`
	UpdateTime          int64                   `json:"update_time" db:"update_time"`
	CPU                 CPUBasicInformation     `json:"cpu" db:"cpu"`
	Memory              MemoryBasicInformation  `json:"memory" db:"memory"`
	Disk                DiskBasicInformation    `json:"disk" db:"disk"`
	Network             NetworkBasicInformation `json:"network" db:"network"`
}

type CPUBasicInformation struct {
	CPUCount int `json:"cpu_count" db:"cpu_count"`
}

type DiskBasicInformation struct {
	DiskSize float32 `json:"disk_size" db:"disk_size"`
}

type MemoryBasicInformation struct {
	MemorySize float32 `json:"memory_size" db:"memory_size"`
}

type NetworkBasicInformation struct {
	NetworkMac string `json:"network_mac" db:"network_mac"`
	NetworkIP  string `json:"network_ip" db:"network_ip"`
}

func TransformEquipmentFromMongoToMysql(ebia *EquipmentBasicInformationAwareness) *EquipmentBasicInformationAwarenessMysql {
	return &EquipmentBasicInformationAwarenessMysql{
		DataType:            ebia.DataType,
		EquipmentID:         ebia.EquipmentID,
		OperateSystem:       ebia.OperateSystem,
		NetworkName:         ebia.NetworkName,
		Platform:            ebia.Platform,
		Architecture:        ebia.Architecture,
		BootTimeInTimestamp: ebia.BootTimeInTimestamp,
		BootTimeInString:    ebia.BootTimeInString,
		User:                ebia.User,
		Host:                ebia.Host,
		UpdateTime:          ebia.UpdateTime,
		CPUCount:            ebia.CPU.CPUCount,
		DiskSize:            ebia.Disk.DiskSize,
		MemorySize:          ebia.Memory.MemorySize,
		NetworkMac1:         ebia.Network.NetworkMac,
		NetworkIP1:          ebia.Network.NetworkIP,
		NetworkMac2:         "",
		NetworkIP2:          "",
		NetworkCard1:        "",
		NetworkCard2:        "",
		Authenticated:       0,
	}
}

type EquipmentStatusAwareness struct {
	DataType    int           `json:"data_type"`
	EquipmentID int           `json:"equipment_id"`
	RunningTime int64         `json:"running_time"`
	UpdateTime  int64         `json:"update_time"`
	CPU         CPUStatus     `json:"cpu"`
	Memory      MemoryStatus  `json:"memory"`
	Disk        DiskStatus    `json:"disk"`
	Network     NetworkStatus `json:"network"`
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

type ResponseForEquipmentBasicInformation struct {
	DataType      int `json:"data_type"`
	EquipmentID   int `json:"equipment_id"`
	Authenticated int `json:"authenticated"`
}

func NewResponseForEquipmentBasicInformation() *ResponseForEquipmentBasicInformation {
	return &ResponseForEquipmentBasicInformation{
		DataType: 32,
	}
}
