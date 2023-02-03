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

func (db *MongoDb) Where(prop string, cond string, value any) {
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

func (db *MongoDb) GenerateId() interface{} {
	return primitive.NewObjectID()
}

func (db *MongoDb) Create(object interface{}) (interface{}, error) {
	t := time.Now()
	defer bench("Create", t)
	db.isCollection()
	result, err := db.Collection.InsertOne(ctx, object)

	return result, err
}

func (db *MongoDb) Find(entity interface{}) error {
	t := time.Now()
	defer bench("FindAll", t)
	db.isCollection()
	cursor, err := db.Collection.Find(ctx, db.Filters)
	handleErr(err)
	return cursor.All(ctx, entity)
}

func (db *MongoDb) FindOne(entity interface{}) error {
	t := time.Now()
	defer bench("FindOne", t)
	db.isCollection()

	res := db.Collection.FindOne(ctx, db.Filters)
	return res.Decode(entity)
}

func (db *MongoDb) UpdateById(id string, entity interface{}) error {
	t := time.Now()
	defer bench("UpdateById", t)
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()

	result := db.Collection.FindOneAndUpdate(ctx, bson.M{"id": objectId}, bson.M{"$set": entity})
	return result.Err()
}

func (db *MongoDb) DeleteById(id string) error {
	t := time.Now()
	defer bench("DeleteById", t)
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()

	result := db.Collection.FindOneAndDelete(ctx, bson.M{"id": objectId})
	return result.Err()
}

func (db *MongoDb) DeleteMany() interface{} {
	t := time.Now()
	defer bench("DeleteMany", t)
	db.isCollection()

	result, err := db.Collection.DeleteMany(ctx, db.Filters)
	handleErr(err)
	return result
}

func (db *MongoDb) InsertMany(documents []interface{}) interface{} {
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
