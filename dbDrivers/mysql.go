package dbDrivers

import (
	"TouchAll-DataCenter/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var mysqlDB *gorm.DB

func initMysql() {
	config := config.NewConfig()
	user, pass, host, port, database := config.GetMysqlConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8&loc=Local", user, pass, host, port, database)
	var err error

	mysqlDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		panic("Mysql connection failed")
	}

	db, _ := mysqlDB.DB()
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(120 * time.Second)

	log.Printf("Mysql connection successed")
}

func GetMysqlDB() *gorm.DB {
	return mysqlDB
}
