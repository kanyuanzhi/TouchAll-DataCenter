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

// 更新设备信息
func UpdateEquipmentBasicInformation(equipmentJson models.EquipmentBasicInformationAwareness) bool {
	equipmentMysql := models.TransformEquipmentFromJsonToMysql(&equipmentJson)
	result := mysqlDB.Model(&equipmentMysql).Omit("id", "data_type", "equipment_id", "network_mac_1", "network_mac_2", "authenticated").Updates(equipmentMysql)
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
