package entities

import (
	"fmt"
	"time"

	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ajuste struct {
	Name  string
	Prop  string
	Value float32
}

type Variable struct {
	Id                           primitive.ObjectID `json:"id" xml:"id" form:"id"`
	CreatedAt                    time.Time          `json:"CreatedAt" xml:"CreatedAt" form:"CreatedAt"`
	UpdatedAt                    time.Time          `json:"UpdatedAt" xml:"UpdatedAt" form:"UpdatedAt"`
	SemanasVida                  int                `json:"semanas_vida" xml:"semanas_vida" form:"semanas_vida"`
	MortalidadRecria             float32            `json:"mortalidad_recria" xml:"mortalidad_recria" form:"mortalidad_recria"`
	MortalidadRecriaAjustado     []Ajuste           `json:"mortalidad_recria_ajustado" xml:"mortalidad_recria_ajustado" form:"mortalidad_recria_ajustado"`
	MortalidadProduccion         float32            `json:"mortalidad_produccion" xml:"mortalidad_produccion" form:"mortalidad_produccion"`
	MortalidadProduccionAjustado []Ajuste           `json:"mortalidad_produccion_ajustado" xml:"mortalidad_produccion_ajustado" form:"mortalidad_produccion_ajustado"`
	MortalidadPollitos           float32            `json:"mortalidad_pollitos" xml:"mortalidad_pollitos" form:"mortalidad_pollitos"`
	MortalidadPollitosAjustado   []Ajuste           `json:"mortalidad_pollitos_ajustado" xml:"mortalidad_pollitos_ajustado" form:"mortalidad_pollitos_ajustado"`
	ProduccionHuevosAjustado     []Ajuste           `json:"produccion_huevos_ajustado" xml:"produccion_huevos_ajustado" form:"produccion_huevos_ajustado"`
	AprovechamientoAjustado      []Ajuste           `json:"aprovechamiento_ajustado" xml:"aprovechamiento_ajustado" form:"aprovechamiento_ajustado"`
	NacimientosAjustado          []Ajuste           `json:"nacimientos_ajustado" xml:"nacimientos_ajustado" form:"nacimientos_ajustado"`

	Tipo                 string  `json:"tipo" xml:"tipo" form:"tipo"`
	VentaPollosVivos     int     `json:"venta_pollos_vivos" xml:"venta_pollos_vivos" form:"venta_pollos_vivos"`
	PollosVivosMataderos int     `json:"pollos_vivos_mataderos" xml:"pollos_vivos_mataderos" form:"pollos_vivos_mataderos"`
	LbPolloProcesados    int     `json:"lb_pollo_procesados" xml:"lb_pollo_procesados" form:"lb_pollo_procesados"`
	LbPolloMatadero      int     `json:"lb_pollo_matadero" xml:"lb_pollo_matadero" form:"lb_pollo_matadero"`
	PesoPromedio         float32 `json:"peso_promedio" xml:"peso_promedio" form:"peso_promedio"`
}

func (v *Variable) GetDbPort() *portout.DbPort[Variable] {
	return &portout.DbPort[Variable]{
		Name:   "variables",
		Entity: v,
	}
}

func (v *Variable) Save() (interface{}, error) {
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	v.Id = primitive.NewObjectID()
	return v.GetDbPort().Save()
}

func (v *Variable) Update(id string) error {
	return v.GetDbPort().Update(id)
}

func (v *Variable) Delete(id string) error {
	return v.GetDbPort().Delete(id)
}

func (v *Variable) GetAll(filters string) ([]Variable, error) {
	return v.GetDbPort().GetAll(filters)
}

func (v *Variable) FindOne(filters string) error {
	return v.GetDbPort().FindOne(filters)
}

func (v *Variable) List(filter string) {
	results, err := v.GetAll(filter)
	handleErr(err)
	fmt.Printf("%s %v Items\n", color.MagentaString("[Results]:"), len(results))
	for _, el := range results {

		fmt.Println("------------------------------------------------------------")
		fmt.Printf("Id: %v\n", el.Id)
		fmt.Printf("Tipo: %v\n", el.Tipo)
		fmt.Printf("Mortalidad Recria: %v\n", el.MortalidadRecria)
		fmt.Printf("Mortalidad Produccion: %v\n", el.MortalidadProduccion)

	}
	fmt.Println("------------------------------------------------------------")
}
