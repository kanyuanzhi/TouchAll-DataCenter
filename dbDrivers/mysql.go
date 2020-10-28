package dbDrivers

import (
	"dataCenter/utils"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func GetMysqlConn() *sqlx.DB {
	config := utils.NewConfig()
	user, pass, host, port, database := config.GetMysqlConfig()
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, database)
	db, err := sqlx.Connect("mysql", url)
	if err != nil {
		log.Println("url: " + url)
		log.Fatal("Mysql url wrong: " + err.Error())
	}
	// 最大连接数
	db.SetMaxOpenConns(100)
	// 闲置连接数
	db.SetMaxIdleConns(20)
	// 最大连接周期
	db.SetConnMaxLifetime(120 * time.Second)

	if err = db.Ping(); nil != err {
		log.Fatal("Mysql connection failed: " + err.Error())
	}
	log.Printf("Mysql connection successed")

	return db
}
