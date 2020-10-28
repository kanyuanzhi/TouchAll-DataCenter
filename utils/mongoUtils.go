package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func InsertOneRecord(document interface{}, connToCollection *mongo.Collection) (success bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = connToCollection.InsertOne(ctx, document)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	//log.Println(res.InsertedID)
	return true, nil
}

func InsertManyRecords(documents []interface{}, connToCollection *mongo.Collection) (success bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = connToCollection.InsertMany(ctx, documents)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	//log.Println(res.InsertedIDs)
	return true, nil
}
