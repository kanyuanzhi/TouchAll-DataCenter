package dbDrivers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func GetConn(dbname string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	opt := options.Client().ApplyURI("mongodb://root:root@10.211.55.9:27017")
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
	return client.Database(dbname), nil
}
