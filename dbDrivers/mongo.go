package dbDrivers

import (
	"context"
	"dataCenter/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func GetConn(dbname string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	config := utils.NewConfig()
	username, password, host, port := config.GetMongodbConfig()
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	opt := options.Client().ApplyURI(url)
	opt.SetMaxPoolSize(50)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Printf("Mongodb connected")

	return client.Database(dbname), nil
}
