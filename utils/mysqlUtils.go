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
	sqlStr := "SELECT network_mac FROM equipment WHERE network_mac=?"
	var ebiaMysql models.EquipmentBasicInformationAwarenessMysql
	err := db.Get(&ebiaMysql, sqlStr, ebia.Network.NetworkMac)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func InsertEquipmentBasicInformation(ebia *models.EquipmentBasicInformationAwareness, db *sqlx.DB) bool {
	if CheckEquipmentBasicInformation(ebia, db) {
		DeleteEquipmentBasicInformation(ebia, db)
	}
	ebiaMysql := models.TransformEquipmentFromMongoToMysql(ebia)
	sqlStr := "INSERT INTO equipment(network_mac,equipment_id,network_ip,network_name," +
		"operate_system,architecture,boot_time_in_timestamp,platform," +
		"boot_time_in_string,user,host,update_time," +
		"cpu_count,disk_size,memory_size,data_type) " +
		"VALUES(:network_mac,:equipment_id,:network_ip,:network_name," +
		":operate_system,:architecture,:boot_time_in_timestamp,:platform," +
		":boot_time_in_string,:user,:host,:update_time," +
		":cpu_count,:disk_size,:memory_size,:data_type)"
	_, err := db.NamedExec(sqlStr, ebiaMysql)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
