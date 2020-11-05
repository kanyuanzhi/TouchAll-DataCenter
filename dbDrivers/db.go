package dbDrivers

import "dataCenter/config"

func init() {
	config := config.NewConfig()
	useMongodb := config.GetValue("mongodb.use").(bool)
	useMysql := config.GetValue("mysql.use").(bool)
	if useMongodb {
		initMongodb()
	}
	if useMysql {
		initMysql()
	}
}
