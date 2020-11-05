package dbDrivers

import (
	"context"
	"dataCenter/config"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var MongoDB *mongo.Database

func initMongodb() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	config := config.NewConfig()
	username, password, host, port, database := config.GetMongodbConfig()
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	opt := options.Client().ApplyURI(url)
	opt.SetMaxPoolSize(50)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatal("Mongodb connection failed" + err.Error())
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal("Mongodb connection failed" + err.Error())
	}

	log.Printf("Mongodb connection successed")

	MongoDB = client.Database(database.(string))
}

func GetMongoDB() *mongo.Database {
	return MongoDB
}
