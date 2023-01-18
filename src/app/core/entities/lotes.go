package entities

import (
	"fmt"
	"strings"
	"time"

	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lote struct {
	Id        primitive.ObjectID `bson: "_id"`
	CreatedAt time.Time          `bson: "created_at"`
	UpdatedAt time.Time          `bson: "updated_at"`
	Numero    int64              `bson: "numero"`
	Entrada   string             `bson: "entrada"`
	Month     string             `bson: "month"`
	Year      string             `bson: "year"`
	Empresa   string             `bson: "empresa"`
	Hembras   int32              `bson: "hembras"`
	Machos    int32              `bson: "machos"`
}

/****************************************************************************
*	Lote Application Servces
****************************************************************************/
func (l *Lote) GetDbPort() *portout.DbPort[Lote] {
	return &portout.DbPort[Lote]{
		Name:   "lotes",
		Entity: l,
	}
}

func (l *Lote) Save() (interface{}, error) {
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	l.Id = primitive.NewObjectID()
	l.Month = strings.Split(l.Entrada, "/")[1]
	l.Year = strings.Split(l.Entrada, "/")[2]
	return l.GetDbPort().Save()
}

func (l *Lote) Update(id string) error {
	return l.GetDbPort().Update(id)
}

func (l *Lote) Delete(id string) error {
	return l.GetDbPort().Delete(id)
}

func (l *Lote) GetAll(filters string) ([]Lote, error) {
	return l.GetDbPort().GetAll(filters)
}

func (l *Lote) FindOne(filters string) error {
	return l.GetDbPort().FindOne(filters)
}

func (l *Lote) List(filter string) {
	results, err := l.GetAll(filter)
	handleErr(err)
	fmt.Printf("results: %v\n", len(results))
	for _, el := range results {

		fmt.Println("------------------------------------------------------------")
		fmt.Printf("Id: %v\n", el.Id)
		fmt.Printf("Lote: %v\n", el.Numero)
		fmt.Printf("Empresa: %v\n", el.Empresa)
		fmt.Printf("Fecha De Entrada: %v\n", el.Entrada)
		fmt.Printf("Hembras: %v\n", el.Hembras)
		fmt.Printf("Machos: %v\n", el.Machos)
	}
	fmt.Println("------------------------------------------------------------")
}
