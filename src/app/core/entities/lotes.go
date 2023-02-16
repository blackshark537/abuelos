package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	portout "github.com/blackshark537/dataprod/src/app/core/port_out"

	"github.com/fatih/color"
)

type Lote struct {
	Id        interface{} `json:"id" xml:"id" form:"id"`
	CreatedAt time.Time   `json:"createdAt" xml:"createdAt" form:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt" xml:"updatedAt" form:"updatedAt"`
	Numero    int64       `json:"numero" xml:"numero" form:"numero" validate:"required, number"`
	Entrada   string      `json:"entrada" xml:"entrada" form:"entrada" validate:"required"`
	Month     int64       `json:"month" xml:"month" form:"month"`
	Year      int64       `json:"year" xml:"year" form:"year"`
	Empresa   string      `json:"empresa" xml:"empresa" form:"empresa" validate:"required, min=10, max=100"`
	Hembras   int32       `json:"hembras" xml:"hembras" form:"hembras" validate:"required, number, min=0"`
	Machos    int32       `json:"machos" xml:"machos" form:"machos" validate:"required, number, min=0"`
	Tipo      string      `json:"tipo" xml:"tipo" form:"tipo"  validate:"required,min=1,max=25"`
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

// Not Tested
func (l *Lote) Count(filters string) int64 {
	count, err := l.GetDbPort().Count(filters)
	if err != nil {
		count = 0
	}
	return count
}

func (l *Lote) Save() (interface{}, error) {
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	l.Id = l.GetDbPort().NewId()
	month, err := strconv.ParseInt(strings.Split(l.Entrada, "-")[1], 10, 64)
	l.Month = month
	year, err := strconv.ParseInt(strings.Split(l.Entrada, "-")[0], 10, 64)
	if err != nil {
		return nil, err
	}
	l.Year = year
	return l.GetDbPort().Save()
}

func (l *Lote) Update(id string) error {
	return l.GetDbPort().Update(id)
}

func (l *Lote) Delete(id string) error {
	return l.GetDbPort().Delete(id)
}

func (l *Lote) DeleteMany(filters string) error {
	return errors.New("Not implemented")
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
	fmt.Printf("%s %v Items\n", color.MagentaString("[Results]:"), len(results))
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
