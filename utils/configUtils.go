package utils

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	CViper *viper.Viper
}

func NewConfig() *Config {
	cviper := viper.New()
	cviper.SetConfigName("config")
	cviper.AddConfigPath("./")
	cviper.SetConfigType("json")
	err := cviper.ReadInConfig()
	if err != nil {
		log.Printf("config file error: %s\n", err)
	}
	return &Config{
		CViper: cviper,
	}
}

func (config *Config) GetValue(key string) interface{} {
	value := config.CViper.Get(key)
	return value
}

func (config *Config) GetMongodbConfig() (interface{}, interface{}, interface{}, interface{}) {
	username := config.GetValue("mongodb.username")
	password := config.GetValue("mongodb.password")
	host := config.GetValue("mongodb.host")
	port := config.GetValue("mongodb.port")
	return username, password, host, port
}
