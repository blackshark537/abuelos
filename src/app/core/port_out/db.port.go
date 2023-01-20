package portout

import (
	"context"

	"github.com/blackshark537/dataprod/src/app/core/config"
	"github.com/blackshark537/dataprod/src/app/infraestructure/database"
)

type DbPort[t interface{}] struct {
	Name   string
	Entity *t
}

type DbAdapter struct{}

var mongodb *database.MongoDb = nil

func (a *DbAdapter) ForRoot() {
	mongodb = &database.MongoDb{
		Uri:  config.DatabaseUri,
		Name: config.DatabaseName,
	}
}

func (port *DbPort[t]) Count(filters string) (int64, error) {
	mongodb.SetFilters(filters)
	count, err := mongodb.Count()
	return count, err
}

func (port *DbPort[t]) Save() (interface{}, error) {
	mongodb.SelectTable(port.Name)
	return mongodb.Create(port.Entity)
}

func (port *DbPort[t]) Update(id string) error {
	mongodb.SelectTable(port.Name)
	result := mongodb.UpdateById(id, port.Entity)
	return result.Err()
}

func (port *DbPort[t]) GetAll(filters string) ([]t, error) {
	mongodb.SelectTable(port.Name)
	mongodb.SetFilters(filters)
	var results []t
	cursor := mongodb.Find()
	err := cursor.All(context.TODO(), &results)
	return results, err
}

func (port *DbPort[t]) FindOne(filters string) error {
	mongodb.SelectTable(port.Name)
	mongodb.SetFilters(filters)
	result := mongodb.FindOne()
	return result.Decode(&port.Entity)
}

func (port *DbPort[t]) Delete(id string) error {
	mongodb.SelectTable(port.Name)
	result := mongodb.DeleteById(id)
	return result.Err()
}
