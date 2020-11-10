package models

import (
	"gorm.io/gorm"
	"sync"
	"time"
)

type EquipmentBasicInformationAwarenessMysql struct {
	gorm.Model
	DataType       int       `json:"data_type" db:"data_type"`
	EquipmentID    int       `json:"equipment_id" db:"equipment_id" gorm:"index"`
	EquipmentType  int       `json:"equipment_type"`
	EquipmentGroup int       `json:"equipment_group"`
	OperateSystem  string    `json:"operate_system" db:"operate_system"`
	NetworkName    string    `json:"network_name" db:"network_name"`
	Platform       string    `json:"platform" db:"platform"`
	Architecture   string    `json:"architecture" db:"architecture"`
	Processor      string    `json:"processor"`
	BootTime       time.Time `json:"boot_time" db:"boot_time"`
	User           string    `json:"user" db:"user"`
	Host           string    `json:"host" db:"host"`
	CPUCount       int       `json:"cpu_count" db:"cpu_count"`
	DiskSize       float32   `json:"disk_size" db:"disk_size"`
	MemorySize     float32   `json:"memory_size" db:"memory_size"`
	NetworkMac1    string    `json:"network_mac_1" db:"network_mac_1" gorm:"column:network_mac_1"`
	NetworkIP1     string    `json:"network_ip_1" db:"network_ip_1" gorm:"column:network_ip_1"`
	NetworkMac2    string    `json:"network_mac_2" db:"network_mac_2" gorm:"column:network_mac_2"`
	NetworkIP2     string    `json:"network_ip_2" db:"network_ip_2" gorm:"column:network_ip_2"`
	NetworkCard1   string    `json:"network_card_1" db:"network_card_1" gorm:"column:network_card_1"`
	NetworkCard2   string    `json:"network_card_2" db:"network_card_2" gorm:"column:network_card_2"`
	Authenticated  int       `json:"authenticated" db:"authenticated"`
}

// EquipmentBasicInformationAwarenessMysql结构体对应mysql数据库中的equipment表
func (EquipmentBasicInformationAwarenessMysql) TableName() string {
	return "equipment"
}

type EquipmentBasicInformationAwareness struct {
	DataType      int                     `json:"data_type" db:"data_type"`
	EquipmentID   int                     `json:"equipment_id" db:"equipment_id"`
	OperateSystem string                  `json:"operate_system" db:"operate_system"`
	NetworkName   string                  `json:"network_name" db:"network_name"`
	Platform      string                  `json:"platform" db:"platform"`
	Architecture  string                  `json:"architecture" db:"architecture"`
	Processor     string                  `json:"processor"`
	BootTime      int64                   `json:"boot_time" db:"boot_time"`
	User          string                  `json:"user" db:"user"`
	Host          string                  `json:"host" db:"host"`
	UpdatedAt     int64                   `json:"updated_at"`
	CPU           CPUBasicInformation     `json:"cpu" db:"cpu"`
	Memory        MemoryBasicInformation  `json:"memory" db:"memory"`
	Disk          DiskBasicInformation    `json:"disk" db:"disk"`
	Network       NetworkBasicInformation `json:"network" db:"network"`
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
	NetworkMac  string `json:"network_mac" db:"network_mac"`
	NetworkIP   string `json:"network_ip" db:"network_ip"`
	NetworkCard string `json:"network_card" db:"network_card"`
}

func TransformEquipmentFromJsonToMysql(ebia *EquipmentBasicInformationAwareness) *EquipmentBasicInformationAwarenessMysql {
	return &EquipmentBasicInformationAwarenessMysql{
		DataType:       ebia.DataType,
		EquipmentID:    ebia.EquipmentID,
		EquipmentType:  0,
		EquipmentGroup: 0,
		OperateSystem:  ebia.OperateSystem,
		NetworkName:    ebia.NetworkName,
		Platform:       ebia.Platform,
		Architecture:   ebia.Architecture,
		Processor:      ebia.Processor,
		BootTime:       time.Unix(ebia.BootTime, 0),
		User:           ebia.User,
		Host:           ebia.Host,
		CPUCount:       ebia.CPU.CPUCount,
		DiskSize:       ebia.Disk.DiskSize,
		MemorySize:     ebia.Memory.MemorySize,
		NetworkMac1:    ebia.Network.NetworkMac,
		NetworkIP1:     ebia.Network.NetworkIP,
		NetworkMac2:    "",
		NetworkIP2:     "",
		NetworkCard1:   ebia.Network.NetworkCard,
		NetworkCard2:   "",
		Authenticated:  0,
	}
}

type EquipmentStatusAwareness struct {
	DataType    int           `json:"data_type"`
	EquipmentID int           `json:"equipment_id"`
	RunningTime int64         `json:"running_time"`
	UpdatedAt   int64         `json:"updated_at"`
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

// 响应设备端推送的设备基本信息
type ResponseForEquipmentBasicInformation struct {
	DataType      int  `json:"data_type"`
	EquipmentID   int  `json:"equipment_id"`
	Authenticated int  `json:"authenticated"`
	UseMysql      bool `json:"use_mysql"`
	UseMongodb    bool `json:"use_mongodb"`
}

func NewResponseForEquipmentBasicInformation() *ResponseForEquipmentBasicInformation {
	return &ResponseForEquipmentBasicInformation{
		DataType: 32,
	}
}

// websocket推送流
type EquipmentStatusStream struct {
	DataType int `json:"data_type"`

	// 使用普通字典记录终端请求的指定设备状态信息以返回终端
	StatusStream map[int]EquipmentStatusAwareness `json:"status_stream"`

	// 使用并发安全的字典类型sync.Map完成多个设备信息的并发合并
	StatusStreamSyncMap sync.Map `json:"status_stream_sync_map"`
}

func NewEquipmentStatusStream() *EquipmentStatusStream {
	return &EquipmentStatusStream{
		DataType:     33,
		StatusStream: make(map[int]EquipmentStatusAwareness),
	}
}
