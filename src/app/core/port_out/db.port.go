package portout

import (
	"fmt"
)

type DbPort[t interface{}] struct {
	Name     string
	Entity   *t
	Entities *[]t
}

type Database interface {
	GenerateId() interface{}
	SelectTable(table string)
	SetFilters(f string)
	Where(prop string, cond string, value any)
	Count() (int64, error)
	Create(object interface{}) (interface{}, error)
	Find(entity interface{}) error
	FindOne(entity interface{}) error
	UpdateById(id string, entity interface{}) error
	DeleteById(id string) error
	DeleteMany() interface{}
	InsertMany(documents []interface{}) interface{}
}

var database Database = nil

func InjectDatabase(db Database) {
	if database == nil {
		database = db
	}
}

func (port *DbPort[t]) GenerateId() interface{} {
	return database.GenerateId()
}

func (port *DbPort[t]) Count(filters string) (int64, error) {
	database.SetFilters(filters)
	count, err := database.Count()
	return count, err
}

func (port *DbPort[t]) Save() (interface{}, error) {
	database.SelectTable(port.Name)
	return database.Create(port.Entity)
}

func (port *DbPort[t]) InsetMany(documents []interface{}) interface{} {
	database.SelectTable(port.Name)
	return database.InsertMany(documents)
}

func (port *DbPort[t]) Update(id string) error {
	database.SelectTable(port.Name)
	result := database.UpdateById(id, port.Entity)
	return result
}

func (port *DbPort[t]) GetAll(filters string) ([]t, error) {
	database.SelectTable(port.Name)
	database.SetFilters(filters)
	var results []t
	err := database.Find(&results)
	return results, err
}

func (port *DbPort[t]) FindOne(filters string) error {
	database.SelectTable(port.Name)
	database.SetFilters(filters)
	result := database.FindOne(&port.Entity)
	return result
}

func (port *DbPort[t]) Delete(id string) error {
	database.SelectTable(port.Name)
	result := database.DeleteById(id)
	return result
}

func (port *DbPort[t]) DeleteMany(ids []string) interface{} {
	database.SelectTable(port.Name)
	database.SetFilters(fmt.Sprintf("{'id': { '$eq': [%s] } }", ids))
	result := database.DeleteMany()
	return result
}
