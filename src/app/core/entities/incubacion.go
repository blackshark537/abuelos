package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TpMaquina []string = []string{"incubadora", "ponedora"}

type Maquina struct {
	Id        primitive.ObjectID `bson: "_id"`
	CreatedAt time.Time          `bson: "created_at"`
	UpdatedAt time.Time          `bson: "updated_at"`
	Numero    int16              `bson: "numero"`
	Marca     string             `bson: "marca"`
	Modelo    string             `bson: "modelo"`
	Capacidad int64              `bson: "capacidad"`
	Tipo      string             `bson: "tipo"`
}

type Incubacion struct {
	Id         primitive.ObjectID `bson: "_id"`
	CreatedAt  time.Time          `bson: "created_at"`
	UpdatedAt  time.Time          `bson: "updated_at"`
	Lote       string             `bson: "lote"`
	Numero     int64              `bson: "numero"`
	Fecha      time.Time          `bson: "fecha"`
	Huevos     int64              `bson: "huevos"`
	Incubadora Maquina            `bson: "incubadora"`
}
