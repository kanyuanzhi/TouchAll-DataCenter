package utils

import (
	"context"
	"dataCenter/dbDrivers"
	"log"
	"time"
)

var mongoDB = dbDrivers.GetMongoDB()

func InsertOneRecord(document interface{}, collection string) (success bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = mongoDB.Collection(collection).InsertOne(ctx, document)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	//log.Println(res.InsertedID)
	return true, nil
}

func InsertManyRecords(documents []interface{}, collection string) (success bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = mongoDB.Collection(collection).InsertMany(ctx, documents)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	//log.Println(res.InsertedIDs)
	return true, nil
}
