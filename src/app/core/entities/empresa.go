package entities

import (
	"fmt"
	"time"

	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Empresa struct {
	Id         primitive.ObjectID `bson: "_id"`
	CreatedAt  time.Time          `bson: "created_at"`
	UpdatedAt  time.Time          `bson: "updated_at"`
	Nombre     string             `bson: "nombre", unique: true`
	Direccion  string             `bson: "direccion"`
	Provincia  string             `bson: "direccion"`
	Correo     string             `bson: "correos"`
	Telefonos  []string           `bson: "telefonos"`
	Principal  int8               `bson: "principal"`
	Asignacion int64              `bson: "asignacion"`
	Tipo       string             `bson: "tipo"`
	Geopoint   Geopoint           `bson: "geopoint"`
}

/****************************************************************************
*	Empresa Application Servces
****************************************************************************/

func (e *Empresa) GetDbPort() *portout.DbPort[Empresa] {
	return &portout.DbPort[Empresa]{
		Name:   "empresas",
		Entity: e,
	}
}

func (e *Empresa) Save() (interface{}, error) {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	e.Id = primitive.NewObjectID()
	geo := Geopoint{}
	geo.Lat = 0
	geo.Lon = 0
	e.Geopoint = geo
	return e.GetDbPort().Save()
}

func (e *Empresa) Update(id string) error {
	return e.GetDbPort().Update(id)
}

func (e *Empresa) Delete(id string) error {
	return e.GetDbPort().Delete(id)
}

func (e *Empresa) GetAll(filters string) ([]Empresa, error) {
	return e.GetDbPort().GetAll(filters)
}

func (e *Empresa) FindOne(filters string) error {
	return e.GetDbPort().FindOne(filters)
}

func (e *Empresa) List(filter string) {
	results, err := e.GetAll(filter)
	handleErr(err)
	fmt.Printf("results: %v\n", len(results))
	for _, el := range results {

		fmt.Println("------------------------------------------------------------")
		fmt.Printf("Id: %v\n", el.Id)
		fmt.Printf("Nombre: %v\n", el.Nombre)
		fmt.Printf("Dirección: %v\n", el.Direccion)
		if len(el.Telefonos) > 0 {
			fmt.Printf("Telefono: %v\n", el.Telefonos[el.Principal])
		}
		fmt.Printf("Asignación: %v\n", el.Asignacion)
		fmt.Printf("Tipo: %v\n", el.Tipo)
	}
	fmt.Println("------------------------------------------------------------")
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
