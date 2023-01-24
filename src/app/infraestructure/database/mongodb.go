package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/blackshark537/dataprod/src/app/core/config"
	"github.com/fatih/color"
)

type MongoDb struct {
	Uri        string
	Name       string
	Collection *mongo.Collection
	Filters    bson.M
}

var ctx = context.TODO()
var instance = color.MagentaString("[MongoDB]:")
var client *mongo.Client = nil

func (db *MongoDb) createClient() *mongo.Database {
	if client == nil {
		clientOptions := options.Client().ApplyURI(db.Uri)
		_client, err := mongo.Connect(context.TODO(), clientOptions)
		client = _client
		handleErr(err)

		if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
			handleErr(err)
		}
		fmt.Printf("%s Successfully connected\n", instance)
	}
	return client.Database(db.Name)
}

func (db *MongoDb) SelectTable(table string) {
	db.Filters = bson.M{}
	client := db.createClient()
	db.Collection = client.Collection(table)
	fmt.Printf("%s Select Colection: %s\n", instance, color.GreenString(table))
}

func (db *MongoDb) SetFilters(f string) {
	json.Unmarshal([]byte(f), &db.Filters)
}

func (db MongoDb) Where(prop string, cond string, value any) {
	db.Filters = bson.M{
		prop: bson.M{cond: value},
	}
}

func (db *MongoDb) Count() (int64, error) {
	t := time.Now()
	defer bench("Count", t)
	db.isCollection()
	count, err := db.Collection.CountDocuments(ctx, db.Filters)

	return count, err
}

func (db *MongoDb) Create(object interface{}) (*mongo.InsertOneResult, error) {
	t := time.Now()
	defer bench("Create", t)
	db.isCollection()
	result, err := db.Collection.InsertOne(ctx, object)

	return result, err
}

func (db *MongoDb) Find() *mongo.Cursor {
	t := time.Now()
	defer bench("FindAll", t)
	db.isCollection()
	cursor, err := db.Collection.Find(ctx, db.Filters)
	handleErr(err)

	return cursor
}

func (db *MongoDb) FindOne() *mongo.SingleResult {
	t := time.Now()
	defer bench("FindOne", t)
	db.isCollection()

	return db.Collection.FindOne(ctx, db.Filters)
}

func (db *MongoDb) UpdateById(id string, entity interface{}) *mongo.SingleResult {
	t := time.Now()
	defer bench("UpdateById", t)
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()

	return db.Collection.FindOneAndUpdate(ctx, bson.M{"id": objectId}, bson.M{"$set": entity})
}

func (db *MongoDb) DeleteById(id string) *mongo.SingleResult {
	t := time.Now()
	defer bench("DeleteById", t)
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()

	return db.Collection.FindOneAndDelete(ctx, bson.M{"id": objectId})
}

func (db *MongoDb) DeleteMany() *mongo.DeleteResult {
	t := time.Now()
	defer bench("DeleteMany", t)
	db.isCollection()

	result, err := db.Collection.DeleteMany(ctx, db.Filters)
	handleErr(err)
	return result
}

func (db *MongoDb) InsertMany(documents []interface{}) *mongo.InsertManyResult {
	t := time.Now()
	defer bench("InsertMany", t)
	db.isCollection()

	result, err := db.Collection.InsertMany(ctx, documents)
	handleErr(err)
	return result
}

func (db *MongoDb) isCollection() {
	if db.Collection == nil {
		handleErr(errors.New(color.YellowString("[MongoDB]: Collection instance null")))
	}
}

func (db *MongoDb) close() {
	err := db.Collection.Database().Client().Disconnect(ctx)
	handleErr(err)
	fmt.Printf("%s Successfully disconnected\n", instance)
}

func bench(name string, t time.Time) {
	if config.IsBench {
		fmt.Printf("%s Operation: %s - %v mili secs\n", instance, name, time.Since(t).Milliseconds())
	}
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
