package utils

import (
	"dataCenter/models"
	"github.com/jmoiron/sqlx"
	"log"
)

func DeleteEquipmentBasicInformation(ebia *models.EquipmentBasicInformationAwareness, db *sqlx.DB) {
	sqlStr := "DELETE FROM equipment WHERE network_mac='" + ebia.Network.NetworkMac + "'"
	res := db.MustExec(sqlStr)
	if res != nil {
		//log.Println(res)
	}
}

func CheckEquipmentBasicInformation(ebia *models.EquipmentBasicInformationAwareness, db *sqlx.DB) bool {
	sqlStr := "SELECT network_mac FROM equipment_information WHERE network_mac=?"
	var ebiaMysql models.EquipmentBasicInformationAwarenessMysql
	err := db.Get(&ebiaMysql, sqlStr, ebia.Network.NetworkMac)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func IsEquipmentIDExisted(id int, db *sqlx.DB) bool {
	sqlStr := "select equipment_id from equipment_information where equipment_id=?"
	var result models.EquipmentBasicInformationAwarenessMysql
	err := db.Get(&result, sqlStr, id)
	if err != nil {
		log.Println(err.Error() + " IsEquipmentIDExisted=false")
		return false
	}
	return true
}

func UpdateEquipmentBasicInformation(ebia *models.EquipmentBasicInformationAwareness, db *sqlx.DB) bool {
	ebiaMysql := models.TransformEquipmentFromMongoToMysql(ebia)
	sqlStr := "UPDATE equipment_information SET " +
		"network_name=:network_name, operate_system=:operate_system,architecture=:architecture," +
		"boot_time_in_timestamp=:boot_time_in_timestamp,platform=:platform," +
		"boot_time_in_string=:boot_time_in_string,user=:user,host=:host,update_time:=update_time," +
		"cpu_count=:cpu_count,disk_size=:disk_size,memory_size=:memory_size,data_type=:data_type " +
		"WHERE equipment_id=:equipment_id"
	_, err := db.NamedExec(sqlStr, ebiaMysql)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func InsertEquipmentBasicInformation(ebia *models.EquipmentBasicInformationAwareness, db *sqlx.DB) bool {
	ebiaMysql := models.TransformEquipmentFromMongoToMysql(ebia)
	sqlStr := "INSERT INTO equipment_information(" +
		"equipment_id ,network_name, operate_system,architecture,boot_time_in_timestamp,platform," +
		"boot_time_in_string,user,host,update_time," +
		"cpu_count,disk_size,memory_size,data_type,authenticated,network_mac_1) " +
		"VALUES(:equipment_id,:network_name,:operate_system,:architecture,:boot_time_in_timestamp,:platform," +
		":boot_time_in_string,:user,:host,:update_time," +
		":cpu_count,:disk_size,:memory_size,:data_type,:authenticated,:network_mac_1)"
	_, err := db.NamedExec(sqlStr, ebiaMysql)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func CheckNetworkMac(mac string, db *sqlx.DB) (bool, int, int) {
	sqlStr := "select equipment_id,authenticated from equipment_information where network_mac_1=? or network_mac_2=?"
	var ebiaMysql models.EquipmentBasicInformationAwarenessMysql
	err := db.Get(&ebiaMysql, sqlStr, mac, mac)
	if err != nil {
		log.Println(err.Error() + " isEquipmentVisited=false")
		return false, 0, 0
	}
	return true, ebiaMysql.EquipmentID, ebiaMysql.Authenticated
}
