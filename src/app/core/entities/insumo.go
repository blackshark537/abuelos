package entities

import (
	"errors"
	"time"

	portout "github.com/blackshark537/abuelos/src/app/core/port_out"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const (
	TipoMicro = 0 << iota
	TipoMacro
)

type Insumo struct {
	Id        interface{} `json:"id" xml:"id" form:"id"`
	CreatedAt time.Time   `json:"createdAt" xml:"createdAt" form:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt" xml:"updatedAt" form:"updatedAt"`
	Nombre    string      `json:"nombre" xml:"nombre" form:"nombre"`
	Unidad    string      `json:"unidad" xml:"unidad" form:"unidad"`
	Precio    int64       `json:"precio" xml:"precio" form:"precio"`
	Cantidad  int64       `json:"cantidad" xml:"cantidad" form:"cantidad"`
	Tipo      uint8       `json:"tipo" xml:"tipo" form:"tipo"`
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
	i.Id = i.GetDbPort().NewId()
	return i.GetDbPort().Save()
}

func (i *Insumo) Update(id string) error {
	return i.GetDbPort().Update(id)
}

func (i *Insumo) Delete(id string) error {
	return i.GetDbPort().Delete(id)
}

func (i *Insumo) DeleteMany(filters string) error {
	return errors.New("Not Implemented")
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
		tbl.AddRow(el.Id, el.Nombre, el.Cantidad/100, el.Unidad, el.Precio/100, Tipo, el.CreatedAt.String())
	}
	tbl.Print()
}
