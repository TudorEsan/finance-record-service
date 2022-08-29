package helpers

import (
	"App/database"
	"App/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var NetWorthCollection *mongo.Collection = database.OpenCollection(database.Client, "NetWorth")
var RecordCollection *mongo.Collection = database.OpenCollection(database.Client, "Records")

func AddRecord(userId primitive.ObjectID, record models.Record) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// add id to the record
	record.Id = primitive.NewObjectID()
	record.UserId = userId
	record.GenerateStatistics()

	_, err = RecordCollection.InsertOne(ctx, record)
	return
}

func GetRecords(userId primitive.ObjectID, page, limit int) (records []models.Record, err error) {
	records = make([]models.Record, 0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	l := int64(limit)
	skip := int64(page * limit)
	opt := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
		Sort: bson.M{
			"date": -1,
		},
	}
	curr, err := RecordCollection.Find(ctx, bson.M{"userId": userId}, &opt)
	if err != nil {
		return
	}
	for curr.Next(ctx) {
		var record models.Record
		err = curr.Decode(&record)
		if err != nil {
			return
		}
		records = append(records, record)
	}
	return
}

func DeleteRecord(userId, recordId primitive.ObjectID) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	RecordCollection.FindOneAndDelete(ctx, bson.M{"userId": userId, "_id": recordId})
	return
}

func GetRecord(userId primitive.ObjectID, recordId string) (record models.Record, err error) {
	id, err := primitive.ObjectIDFromHex(recordId)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = RecordCollection.FindOne(ctx, bson.M{"userId": userId, "_id": id}).Decode(&record)
	return
}

func UpdateRecord(userId primitive.ObjectID, record models.Record) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	record.GenerateStatistics()
	fmt.Println(userId)
	fmt.Print(record.Id)
	_, err = RecordCollection.UpdateOne(ctx, bson.M{"userId": userId, "_id": record.Id}, bson.M{"$set": record})
	fmt.Print("ERR ", err)
	return
}

func GetRecordCount(userId primitive.ObjectID) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	count, err = RecordCollection.CountDocuments(ctx, bson.M{"userId": userId})
	return
}
