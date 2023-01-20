package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

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

func (db *MongoDb) createClient() *mongo.Database {
	clientOptions := options.Client().ApplyURI(db.Uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	handleErr(err)

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		handleErr(err)
	}

	fmt.Printf("%s Successfully connected\n", instance)
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
	operation("Count")
	db.isCollection()
	count, err := db.Collection.CountDocuments(ctx, db.Filters)
	defer db.close()
	return count, err
}

func (db *MongoDb) Create(object interface{}) (*mongo.InsertOneResult, error) {
	operation("Create")
	db.isCollection()
	result, err := db.Collection.InsertOne(ctx, object)
	defer db.close()
	return result, err
}

func (db *MongoDb) Find() *mongo.Cursor {
	operation("FindAll")
	db.isCollection()
	cursor, err := db.Collection.Find(ctx, db.Filters)
	handleErr(err)
	defer db.close()
	return cursor
}

func (db *MongoDb) FindOne() *mongo.SingleResult {
	operation("FindOne")
	db.isCollection()
	defer db.close()
	return db.Collection.FindOne(ctx, db.Filters)
}

func (db *MongoDb) UpdateById(id string, entity interface{}) *mongo.SingleResult {
	operation("UpdateById")
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()
	defer db.close()
	return db.Collection.FindOneAndUpdate(ctx, bson.M{"id": objectId}, bson.M{"$set": entity})
}

func (db *MongoDb) DeleteById(id string) *mongo.SingleResult {
	operation("DeleteById")
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()
	defer db.close()
	return db.Collection.FindOneAndDelete(ctx, bson.M{"id": objectId})
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

func operation(name string) {
	fmt.Printf("%s Operation: %s\n", instance, name)
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
