package utils

import (
	"TouchAll-DataCenter/dbDrivers"
	"TouchAll-DataCenter/models"
	"errors"
	"gorm.io/gorm"
)

var mysqlDB = dbDrivers.GetMysqlDB()

// 判断设备ID是否已注册
func IsEquipmentIDExisted(id int) bool {
	var equipment models.EquipmentBasicInformationAwarenessMysql
	result := mysqlDB.Select("id").Where("equipment_id=?", id).First(&equipment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		} else {
			panic(result.Error.Error())
			return false
		}
	}
	return true
}

// 更新设备信息（设备初始连接数据中心时发送设备基本信息，自动）
func UpdateEquipmentBasicInformation(equipmentJson models.EquipmentBasicInformationAwareness) bool {
	equipmentMysql := models.TransformEquipmentFromJsonToMysql(&equipmentJson)
	result := mysqlDB.Model(&equipmentMysql).Where("equipment_id=?", equipmentMysql.EquipmentID).Omit("id", "data_type", "equipment_id", "network_mac_1",
		"network_mac_2", "network_ip_2", "network_card_2", "authenticated").Updates(equipmentMysql)
	if result.Error != nil {
		panic(result.Error.Error())
		return false
	}
	return true
}

// 新建设备信息
func InsertEquipmentBasicInformation(equipmentJson models.EquipmentBasicInformationAwareness) bool {
	equipmentMysql := models.TransformEquipmentFromJsonToMysql(&equipmentJson)
	result := mysqlDB.Create(&equipmentMysql)
	if result.Error != nil {
		panic(result.Error.Error())
		return false
	}
	return true
}

// 判断设备网卡是否已注册
func IsEquipmentNetworkMacExisted(mac string) (bool, int, int) {
	var equipment models.EquipmentBasicInformationAwarenessMysql
	result := mysqlDB.Select("equipment_id", "authenticated").Where("network_mac_1=?", mac).Or("network_mac_2=?", mac).First(&equipment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, 0, 0
		} else {
			panic(result.Error.Error())
			return false, 0, 0
		}
	}
	return true, equipment.EquipmentID, equipment.Authenticated
}

// 判断监控摄像机地址是否已注册
func IsCameraHostExisted(host string) (bool, int, int) {
	var camera models.Camera
	result := mysqlDB.Select("camera_id").Where("camera_host=?", host).First(&camera)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, 0, 0
		} else {
			panic(result.Error.Error())
			return false, 0, 0
		}
	}
	return true, camera.CameraID, camera.Authenticated
}

// 判断监控摄像机ID是否已注册
func IsCameraIDExisted(id int) bool {
	var camera models.Camera
	result := mysqlDB.Select("id").Where("camera_id=?", id).First(&camera)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		} else {
			panic(result.Error.Error())
			return false
		}
	}
	return true
}

// 新建监控摄像机信息
func InsertCamera(camera models.Camera) bool {
	result := mysqlDB.Create(&camera)
	if result.Error != nil {
		panic(result.Error.Error())
		return false
	}
	return true
}

// 判断AI监控摄像机地址是否已注册
func IsAICameraHostExisted(host string) (bool, int, int) {
	var aiCamera models.AICamera
	result := mysqlDB.Select("camera_id").Where("camera_host=?", host).First(&aiCamera)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, 0, 0
		} else {
			panic(result.Error.Error())
			return false, 0, 0
		}
	}
	return true, aiCamera.CameraID, aiCamera.Authenticated
}

// 判断AI监控摄像机ID是否已注册
func IsAICameraIDExisted(id int) bool {
	var aiCamera models.AICamera
	result := mysqlDB.Select("id").Where("camera_id=?", id).First(&aiCamera)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		} else {
			panic(result.Error.Error())
			return false
		}
	}
	return true
}

// 新建AI监控摄像机信息
func InsertAICamera(aiCamera models.AICamera) bool {
	result := mysqlDB.Create(&aiCamera)
	if result.Error != nil {
		panic(result.Error.Error())
		return false
	}
	return true
}
