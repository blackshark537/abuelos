package entities

import (
	"time"

	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TipoMicro = 0 << iota
	TipoMacro
)

type Insumo struct {
	Id        primitive.ObjectID `json:"id" xml:"id" form:"id"`
	CreatedAt time.Time          `json:"CreatedAt" xml:"CreatedAt" form:"CreatedAt"`
	UpdatedAt time.Time          `json:"UpdatedAt" xml:"UpdatedAt" form:"UpdatedAt"`
	Nombre    string
	Unidad    string
	Precio    int64
	Cantidad  int64
	Tipo      uint8
}

/****************************************************************************
*	Insumo Application Servces
****************************************************************************/

func (i *Insumo) GetDbPort() *portout.DbPort[Insumo] {
	return &portout.DbPort[Insumo]{
		Name:   "insumo",
		Entity: i,
	}
}

func (i *Insumo) Save() (interface{}, error) {
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	i.Id = primitive.NewObjectID()
	return i.GetDbPort().Save()
}

func (i *Insumo) Update(id string) error {
	return i.GetDbPort().Update(id)
}

func (i *Insumo) Delete(id string) error {
	return i.GetDbPort().Delete(id)
}

func (i *Insumo) DeleteMany(filters string) error {
	var err error = nil
	results, err := i.GetAll(filters)
	if err != nil {
		return err
	}
	for _, el := range results {
		err = el.Delete(el.Id.Hex())
		if err != nil {
			break
		}
	}
	return err
}

// Queries
func (i *Insumo) GetAll(filters string) ([]Insumo, error) {
	return i.GetDbPort().GetAll(filters)
}

func (i *Insumo) FindOne(filters string) error {
	return i.GetDbPort().FindOne(filters)
}

func (i *Insumo) List(filter string) {
	insumos, _ := i.GetAll("")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Id", "Nombre", "Cantidad", "Unidad", "Precio", "Tipo", "Fecha")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, el := range insumos {
		Tipo := "Micro"
		if el.Tipo > 0 {
			Tipo = "Macro"
		}
		tbl.AddRow(el.Id.Hex(), el.Nombre, el.Cantidad/100, el.Unidad, el.Precio/100, Tipo, el.CreatedAt.String())
	}
	tbl.Print()
}
