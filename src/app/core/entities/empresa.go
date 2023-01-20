package entities

import (
	"fmt"
	"log"
	"time"

	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Empresa struct {
	Id         primitive.ObjectID `json:"id" xml:"id" form:"id"`
	CreatedAt  time.Time          `json:"CreatedAt" xml:"CreatedAt" form:"CreatedAt"`
	UpdatedAt  time.Time          `json:"UpdatedAt" xml:"UpdatedAt" form:"UpdatedAt"`
	Nombre     string             `json:"nombre" xml:"nombre" form:"nombre" validate:"required,min=10,max=100"`
	Direccion  string             `json:"direccion" xml:"direccion" form:"direccion"  validate:"required,min=10,max=100"`
	Provincia  string             `json:"provincia" xml:"provincia" form:"provincia"`
	Correo     string             `json:"correo" xml:"correo" form:"correo" `
	Telefonos  []string           `json:"telefonos" xml:"telefonos" form:"telefonos"  validate:"required"`
	Principal  int8               `json:"principal" xml:"principal" form:"principal"  validate:"number"`
	Asignacion int64              `json:"asignacion" xml:"asignacion" form:"asignacion"  validate:"required,number,min=0"`
	Tipo       string             `json:"tipo" xml:"tipo" form:"tipo"  validate:"required,min=1,max=25"`
	Geopoint   Geopoint           `json:"geopoint" xml:"geopoint" form:"geopoint"`
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

func (e *Empresa) Count(filters string) int64 {
	count, err := e.GetDbPort().Count(filters)
	if err != nil {
		count = 0
	}
	return count
}

func (e *Empresa) Save(geo *Geopoint) (interface{}, error) {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	e.Id = primitive.NewObjectID()
	if geo != nil {
		if geo.Lat != 0 && geo.Lon != 0 {
			e.Geopoint = *geo
		}
	}
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
	fmt.Printf("%s %v Items\n", color.MagentaString("[Results]:"), len(results))
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
		log.Fatal(err)
	}
}
