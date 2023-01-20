package services

import (
	"fmt"
	"time"

	"github.com/blackshark537/dataprod/src/app/core/entities"
	"github.com/fatih/color"
)

type DataProjected struct {
	Mortalidad         float32
	MortalidadReal     float32
	Aves               int64
	AvesReal           int64
	HvosProducidos     int64
	HvosProducidosReal int64
	HvosIncubables     int64
	HvosIncubablesReal int64
	Nacimientos        int64
	NacimientosReal    int64
	Pollitos           int64
	PollitosReal       int64
}

type LoteProjection struct {
	Lote         int64
	Fecha        time.Time
	Edad         int8
	Day          int16
	Month        int8
	Year         int64
	EnProduccion bool
	Data         DataProjected
}

type TableProjection struct {
	Day   int16
	Month []DataProjected
	Year  int64
}

var _PROD = []float32{
	0, 8.40, 28.50, 48, 60.30, 69.70, 74.10, 77.50, 77.40, 77.10, 76.30, 74.80, 73.10, 71.40, 69.30, 67.20, 65.10, 63.30,
	61.10, 59, 57.10, 55.30, 53.50, 51.60, 49.70, 48, 46.20, 44.80, 43.60, 42.40, 41.40, 40.50, 39.50, 38.60, 37.70, 36.80,
	36, 35.10, 34.20, 33.30, 32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22,
}
var _APROV = []float32{
	0, 50, 78, 84, 88, 90, 92, 94, 94, 95, 95, 95, 95, 95, 95, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 95, 95, 95, 95, 95, 95,
	94, 94, 94, 94, 94, 93, 93, 93, 93, 93, 93, 92, 92, 92, 92, 91, 91, 90,
}
var _Nac = []float32{
	0, 0, 55, 65, 70, 75, 80, 81, 82, 83, 83, 82, 82, 82, 81, 81, 81, 80, 80, 80, 80, 79, 79, 78, 78, 77, 77, 77, 76, 76, 75, 75, 74, 74,
	73, 73, 72, 70, 69, 68, 67, 66, 65, 64, 63, 62, 61, 60, 59, 58, 57,
}

const (
	semanas_en_recria     = 24
	semanas_en_produccion = 41
	semanas_de_retardo    = 0
)

var (
	variable_mortalidad_recria          float32 = 5.0
	variable_mortalidad_recria_ajustado float32 = 7.0

	variable_mortalidad_produccion          float32 = 12.0
	variable_mortalidad_produccion_ajustado float32 = 17.0

	variable_produccion_huevos_totales float32 = 2
	variable_aprovechamiento_huevos    float32 = 1
	variable_nacimientos               float32 = 10
)

var instance string = color.MagentaString("[Projection]:")

type Projection struct{}

func (p *Projection) List(filters string) {
	results := p.Project(filters)
	fmt.Printf("%s %v Items\n", color.MagentaString("[Results]:"), len(results))
	for _, el := range results {

		fmt.Println("------------------------------------------------------------")
		fmt.Printf("Lote: %d\n", el.Lote)
		fmt.Printf("Fecha: %d-%d-%d\n", el.Day, el.Month, el.Year)
		fmt.Printf("Semana: %d\n", el.Edad)
		fmt.Printf("Mortalidad: %f\n", el.Data.Mortalidad)
		fmt.Printf("Mortalidad Real: %f\n", el.Data.MortalidadReal)
		fmt.Println("")
		fmt.Printf("Aves: %v\n", el.Data.Aves)
		fmt.Printf("Aves Real: %v\n", el.Data.AvesReal)
		fmt.Println("")
		fmt.Printf("Huevos Prod: %v\n", el.Data.HvosProducidos)
		fmt.Printf("Huevos Prod Real: %v\n", el.Data.HvosProducidosReal)
		fmt.Println("")
		fmt.Printf("Huevos Inc: %v\n", el.Data.HvosIncubables)
		fmt.Printf("Huevos Inc Real: %v\n", el.Data.HvosIncubablesReal)
		fmt.Println("")
		fmt.Printf("Nacimientos: %v\n", el.Data.Nacimientos)
		fmt.Printf("Nacimientos Real: %v\n", el.Data.NacimientosReal)

		//fmt.Printf("Pollitos: %v\n", el.Data.Pollitos)
	}
	fmt.Println("------------------------------------------------------------")
}

func (p *Projection) Project(filters string) []LoteProjection {
	operation("LoteProjection")
	loteEntity := new(entities.Lote)
	lotes, err := loteEntity.GetAll(filters)
	if err != nil {
		handleErr(err)
	}
	projection := []LoteProjection{}

	for _, lote := range lotes {
		recria := p.getRecria(lote)
		produccion := p.getProduccion(recria[len(recria)-1])
		projection = append(projection, recria...)
		projection = append(projection, produccion...)
	}
	return projection
}

func (p *Projection) getRecria(lote entities.Lote) []LoteProjection {
	recrias := []LoteProjection{}
	var idx int8 = 0
	fecha, err := time.Parse("2006-01-02", lote.Entrada)
	if err != nil {
		handleErr(err)
	}
	var percent float32 = 100.0
	var percentReal float32 = 100.0

	for i := 0; i < semanas_en_recria*7; i++ {
		if i%7 == 0 {
			idx += 1
		}
		if len(recrias) > 0 {
			percent = recrias[i-1].Data.Mortalidad
			percentReal = recrias[i-1].Data.MortalidadReal
		}

		mortality := percent - (variable_mortalidad_recria / float32(semanas_en_recria*7))
		mortalityReal := percentReal - (variable_mortalidad_recria_ajustado / float32((semanas_en_recria * 7)))

		data := DataProjected{
			Mortalidad:         mortality,
			MortalidadReal:     mortalityReal,
			Aves:               int64(lote.Hembras) - (int64(lote.Hembras) * int64(100.0-mortality) / 100),
			AvesReal:           int64(lote.Hembras) - (int64(lote.Hembras) * int64(100.0-mortalityReal) / 100),
			HvosProducidos:     0,
			HvosProducidosReal: 0,
			HvosIncubables:     0,
			HvosIncubablesReal: 0,
			Nacimientos:        0,
			NacimientosReal:    0,
			Pollitos:           0,
			PollitosReal:       0,
		}
		fecha = fecha.AddDate(0, 0, 1)

		recria := LoteProjection{
			Lote:         lote.Numero,
			Fecha:        fecha,
			Edad:         idx,
			Day:          int16(fecha.Day()),
			Month:        int8(fecha.Month()),
			Year:         int64(fecha.Year()),
			EnProduccion: false,
			Data:         data,
		}

		recrias = append(recrias, recria)
	}
	return recrias
}

func (p *Projection) getProduccion(lote LoteProjection) []LoteProjection {
	produccion := []LoteProjection{}
	var idx int8 = 0
	fecha, err := time.Parse("2006-1-2", fmt.Sprintf("%d-%d-%d", lote.Year, lote.Month, lote.Day))
	if err != nil {
		handleErr(err)
	}
	var percent float32 = 100.0
	var percentReal float32 = 100.0

	for i := 0; i < semanas_en_produccion*7; i++ {
		if i%7 == 0 {
			idx += 1
		}
		if len(produccion) > 0 {
			percent = produccion[i-1].Data.Mortalidad
			percentReal = produccion[i-1].Data.MortalidadReal
		}

		mortality := percent - (variable_mortalidad_produccion / float32(semanas_en_produccion*7))
		mortalityReal := percentReal - (variable_mortalidad_produccion_ajustado / float32((semanas_en_produccion * 7)))
		aves := int64(lote.Data.Aves) - (int64(lote.Data.Aves) * int64(100.0-mortality) / 100)
		aves_real := int64(lote.Data.Aves) - (int64(lote.Data.Aves) * int64(100.0-mortalityReal) / 100)

		// Standar de produccion
		std_prod := _PROD[idx] - (_PROD[idx] * variable_produccion_huevos_totales / 100.0)
		// Standar de Aprovechamiento
		std_aprov := _APROV[idx] - (_APROV[idx] * variable_aprovechamiento_huevos / 100.0)
		// Standar de Nacimientos
		std_nac := _Nac[idx] - (_Nac[idx] * variable_nacimientos / 100.0)

		h_totales_real := aves * int64(std_prod) / 100
		h_totales := aves_real * int64(_PROD[idx]) / 100

		h_incubables_real := h_totales * int64(std_aprov) / 100
		h_incubables := h_totales_real * int64(_APROV[idx]) / 100

		p_nacidos_real := h_incubables * int64(std_nac) / 100
		p_nacidos := h_incubables_real * int64(_Nac[idx]) / 100

		data := DataProjected{
			Mortalidad:         mortality,
			MortalidadReal:     mortalityReal,
			Aves:               aves,
			AvesReal:           aves_real,
			HvosProducidos:     h_totales,
			HvosProducidosReal: h_totales_real,
			HvosIncubables:     h_incubables,
			HvosIncubablesReal: h_incubables_real,
			Nacimientos:        p_nacidos / 2,
			NacimientosReal:    p_nacidos_real / 2,
			Pollitos:           0,
			PollitosReal:       0,
		}
		fecha = fecha.AddDate(0, 0, 1)

		project := LoteProjection{
			Lote:         lote.Lote,
			Fecha:        fecha,
			Edad:         idx + 24,
			Day:          int16(fecha.Day()),
			Month:        int8(fecha.Month()),
			Year:         int64(fecha.Year()),
			EnProduccion: true,
			Data:         data,
		}

		produccion = append(produccion, project)
	}
	return produccion
}

func operation(name string) {
	fmt.Printf("%s Operation: %s\n", instance, name)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
