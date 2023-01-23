package entities

import (
	"time"

	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ajuste struct {
	Name  string
	Value float32
}

type Variable struct {
	Id                           primitive.ObjectID `json:"id" xml:"id" form:"id"`
	CreatedAt                    time.Time          `json:"CreatedAt" xml:"CreatedAt" form:"CreatedAt"`
	UpdatedAt                    time.Time          `json:"UpdatedAt" xml:"UpdatedAt" form:"UpdatedAt"`
	SemanasProduccion            int                `json:"semanas_produccion" xml:"semanas_produccion" form:"semanas_produccion"`
	MortalidadRecria             float32            `json:"mortalidad_recria" xml:"mortalidad_recria" form:"mortalidad_recria"`
	MortalidadRecriaAjustado     []Ajuste           `json:"mortalidad_recria_ajustado" xml:"mortalidad_recria_ajustado" form:"mortalidad_recria_ajustado"`
	MortalidadProduccion         float32            `json:"mortalidad_produccion" xml:"mortalidad_produccion" form:"mortalidad_produccion"`
	MortalidadProduccionAjustado []Ajuste           `json:"mortalidad_produccion_ajustado" xml:"mortalidad_produccion_ajustado" form:"mortalidad_produccion_ajustado"`
	ProduccionHuevos             float32            `json:"produccion_huevos" xml:"produccion_huevos" form:"produccion_huevos"`
	ProduccionHuevosAjustado     []Ajuste           `json:"produccion_huevos_ajustado" xml:"produccion_huevos_ajustado" form:"produccion_huevos_ajustado"`
	Aprovechamiento              float32            `json:"aprovechamiento" xml:"aprovechamiento" form:"aprovechamiento"`
	AprovechamientoAjustado      []Ajuste           `json:"aprovechamiento_ajustado" xml:"aprovechamiento_ajustado" form:"aprovechamiento_ajustado"`
	Nacimientos                  float32            `json:"nacimientos" xml:"nacimientos" form:"nacimientos"`
	NacimientosAjustado          []Ajuste           `json:"nacimientos_ajustado" xml:"nacimientos_ajustado" form:"nacimientos_ajustado"`
	Tipo                         string             `json:"tipo" xml:"tipo" form:"tipo"`
}

func (v *Variable) GetDbPort() *portout.DbPort[Variable] {
	return &portout.DbPort[Variable]{
		Name:   "empresas",
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
