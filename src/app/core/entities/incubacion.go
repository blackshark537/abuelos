package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TpMaquina []string = []string{"incubadora", "ponedora"}

type Maquina struct {
	Id        primitive.ObjectID `json:"id" xml:"id" form:"id"`
	CreatedAt time.Time          `json:"createdAt" xml:"createdAt" form:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" xml:"updatedAt" form:"updatedAt"`
	Numero    int16              `json:"numero" xml:"numero" form:"numero"`
	Marca     string             `json:"marca" xml:"marca" form:"marca"`
	Modelo    string             `json:"modelo" xml:"modelo" form:"modelo"`
	Capacidad int64              `json:"capacidad" xml:"capacidad" form:"capacidad"`
	Tipo      string             `json:"tipo" xml:"tipo" form:"tipo"`
}

type Incubacion struct {
	Id         primitive.ObjectID `json:"id" xml:"id" form:"id"`
	CreatedAt  time.Time          `json:"createdAt" xml:"createdAt" form:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" xml:"updatedAt" form:"updatedAt"`
	Lote       string             `json:"lote" xml:"lote" form:"lote"`
	Numero     int64              `json:"numero" xml:"numero" form:"numero"`
	Fecha      time.Time          `json:"fecha" xml:"fecha" form:"fecha"`
	Huevos     int64              `json:"huevos" xml:"huevos" form:"huevos"`
	Incubadora Maquina            `json:"incubadora" xml:"incubadora" form:"incubadora"`
}
