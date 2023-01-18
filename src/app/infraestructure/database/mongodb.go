package database

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDb struct {
	Uri        string
	Name       string
	Collection *mongo.Collection
	Filters    bson.M
}

var ctx = context.TODO()

func (db *MongoDb) createClient() *mongo.Database {
	clientOptions := options.Client().ApplyURI(db.Uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	handleErr(err)

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB Successfully connected")
	return client.Database(db.Name)
}

func (db *MongoDb) SelectTable(table string) {
	db.Filters = bson.M{}
	client := db.createClient()
	db.Collection = client.Collection(table)
}

func (db *MongoDb) SetFilters(f string) {
	json.Unmarshal([]byte(f), &db.Filters)
}

func (db MongoDb) Where(prop string, cond string, value any) {
	db.Filters = bson.M{
		prop: bson.M{cond: value},
	}
}

func (db *MongoDb) Create(object interface{}) (*mongo.InsertOneResult, error) {
	db.isCollection()
	result, err := db.Collection.InsertOne(ctx, object)
	defer db.close()
	return result, err
}

func (db *MongoDb) ReadAll() *mongo.Cursor {
	db.isCollection()
	cursor, err := db.Collection.Find(ctx, db.Filters)
	handleErr(err)
	defer db.close()
	return cursor
}

func (db *MongoDb) Read() *mongo.SingleResult {
	db.isCollection()
	defer db.close()
	return db.Collection.FindOne(ctx, db.Filters)
}

func (db *MongoDb) UpdateById(id string, entity interface{}) *mongo.SingleResult {
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()
	defer db.close()
	return db.Collection.FindOneAndUpdate(ctx, bson.M{"id": objectId}, bson.M{"$set": entity})
}

func (db *MongoDb) DeleteById(id string) *mongo.SingleResult {
	objectId, err := primitive.ObjectIDFromHex(id)
	handleErr(err)
	db.isCollection()
	defer db.close()
	return db.Collection.FindOneAndDelete(ctx, bson.M{"id": objectId})
}

func (db *MongoDb) isCollection() {
	if db.Collection == nil {
		panic("MongoDb Collection instance null")
	}
}

func (db *MongoDb) close() {
	err := db.Collection.Database().Client().Disconnect(ctx)
	handleErr(err)
	fmt.Println("MongoDB Successfully disconnected")
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
