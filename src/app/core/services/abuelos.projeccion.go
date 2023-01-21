package services

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/blackshark537/dataprod/src/app/core/entities"
	"github.com/fatih/color"
)

type DataProjected struct {
	Mortalidad         float64
	MortalidadReal     float64
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
	Semana       int8
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

var _PROD = []float64{
	0, 8.40, 28.50, 48, 60.30, 69.70, 74.10, 77.50, 77.40, 77.10, 76.30, 74.80, 73.10, 71.40, 69.30, 67.20, 65.10, 63.30,
	61.10, 59, 57.10, 55.30, 53.50, 51.60, 49.70, 48, 46.20, 44.80, 43.60, 42.40, 41.40, 40.50, 39.50, 38.60, 37.70, 36.80,
	36, 35.10, 34.20, 33.30, 32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22,
}
var _APROV = []float64{
	0, 50, 78, 84, 88, 90, 92, 94, 94, 95, 95, 95, 95, 95, 95, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 95, 95, 95, 95, 95, 95,
	94, 94, 94, 94, 94, 93, 93, 93, 93, 93, 93, 92, 92, 92, 92, 91, 91, 90,
}
var _Nac = []float64{
	0, 0, 55, 65, 70, 75, 80, 81, 82, 83, 83, 82, 82, 82, 81, 81, 81, 80, 80, 80, 80, 79, 79, 78, 78, 77, 77, 77, 76, 76, 75, 75, 74, 74,
	73, 73, 72, 70, 69, 68, 67, 66, 65, 64, 63, 62, 61, 60, 59, 58, 57,
}

const (
	semanas_en_recria     = 24
	semanas_en_produccion = 41
	semanas_de_retardo    = 0
)

var (
	variable_mortalidad_recria          float64 = 4.0
	variable_mortalidad_recria_ajustado float64 = 5.0

	variable_mortalidad_produccion          float64 = 10.0
	variable_mortalidad_produccion_ajustado float64 = 12.0

	variable_produccion_huevos_totales float64 = 1
	variable_aprovechamiento_huevos    float64 = 0.6
	variable_nacimientos               float64 = 3
)

var instance string = color.MagentaString("[Projection]:")

type AbuelosProjection struct{}

// [Warning] For CLI Use Only
func ListAbuelos(lote string) {
	filters := fmt.Sprintf(`{"numero" : {"$eq": %v}}`, lote)
	results := ProjectAbuelos(filters)
	fmt.Printf("%s %v Items\n", color.MagentaString("[Results]:"), len(results))
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("| Lote\t | Sem\t | mortalidad\t | Aves\t | Produc | Incub | Nac\t |")
	fmt.Println("-------------------------------------------------------------------")
	for _, el := range results {
		fmt.Printf("| %d\t | %d\t | %f\t | %d\t | %d\t  | %d\t | %d\t |\n", el.Lote, el.Semana, el.Data.Mortalidad, el.Data.Aves, el.Data.HvosProducidos, el.Data.HvosIncubables, el.Data.Nacimientos)
	}
	fmt.Println("-------------------------------------------------------------------")
}

// [Warning] For CLI Use Only
func AbuelosTable(year string, dataType string, isProduccion bool) {
	cols, rows := AbuelosProjectionTable(year, dataType, isProduccion)

	fmt.Printf("Data: %s Year: %v\n", dataType, year)
	fmt.Println("----------------------------------------------------------------------------------------------------------")
	fmt.Printf("| %v\t | %v\t | %v\t | %v\t | %v\t | %v\t | %v\t | %v\t | %v\t | %v\t | %v\t | %v\t | %v\t |\n", cols[0], cols[1], cols[2], cols[3], cols[4], cols[5], cols[6], cols[7], cols[8], cols[9], cols[10], cols[11], cols[12])
	for d := 0; d < len(rows[0]); d++ {
		fmt.Printf("| %v\t | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v |\n", rows[0][d], rows[1][d], rows[2][d], rows[3][d], rows[4][d], rows[5][d], rows[6][d], rows[7][d], rows[8][d], rows[9][d], rows[10][d], rows[11][d], rows[12][d])
	}
	fmt.Println("----------------------------------------------------------------------------------------------------------")
}

func ProjectAbuelos(filters string) []LoteProjection {
	operation("LoteProjection")
	loteEntity := new(entities.Lote)
	lotes, err := loteEntity.GetAll(filters)
	if err != nil {
		handleErr(err)
	}
	projection := []LoteProjection{}

	for _, lote := range lotes {
		recria := getAbuelosRecria(lote)
		produccion := getAbuelosProduccion(recria[len(recria)-1])
		projection = append(projection, recria...)
		projection = append(projection, produccion...)
	}
	return projection
}

func getAbuelosRecria(lote entities.Lote) []LoteProjection {
	recrias := []LoteProjection{}
	data := DataProjected{
		Mortalidad:         100,
		MortalidadReal:     100,
		Aves:               int64(lote.Hembras),
		AvesReal:           int64(lote.Hembras),
		HvosProducidos:     0,
		HvosProducidosReal: 0,
		HvosIncubables:     0,
		HvosIncubablesReal: 0,
		Nacimientos:        0,
		NacimientosReal:    0,
		Pollitos:           0,
		PollitosReal:       0,
	}
	fecha, err := time.Parse("2006-01-02", lote.Entrada)
	if err != nil {
		handleErr(err)
	}

	recria := LoteProjection{
		Lote:         lote.Numero,
		Fecha:        fecha,
		Semana:       0,
		Day:          int16(fecha.Day()),
		Month:        int8(fecha.Month()),
		Year:         int64(fecha.Year()),
		EnProduccion: false,
		Data:         data,
	}

	recrias = append(recrias, recria)
	var idx int8 = 0
	var percent float64 = 0
	var percentReal float64 = 0

	for i := 0; i < semanas_en_recria*7; i++ {
		if i%7 == 0 {
			idx += 1
		}
		if len(recrias) > 0 {
			percent = recrias[i].Data.Mortalidad
			percentReal = recrias[i].Data.MortalidadReal
		}

		mortality := percent - (variable_mortalidad_recria / float64(semanas_en_recria*7))
		mortalityReal := percentReal - (variable_mortalidad_recria_ajustado / float64((semanas_en_recria * 7)))

		aves := float64(lote.Hembras) * mortality / 100
		aves_real := float64(lote.Hembras) * mortalityReal / 100

		data := DataProjected{
			Mortalidad:         mortality,
			MortalidadReal:     mortalityReal,
			Aves:               int64(math.Floor(aves)),
			AvesReal:           int64(math.Floor(aves_real)),
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
			Semana:       idx,
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

func getAbuelosProduccion(lote LoteProjection) []LoteProjection {
	produccion := []LoteProjection{}
	var idx int8 = 0
	fecha, err := time.Parse("2006-1-2", fmt.Sprintf("%d-%d-%d", lote.Year, lote.Month, lote.Day))
	if err != nil {
		handleErr(err)
	}
	var percent float64 = 100.0
	var percentReal float64 = 100.0

	for i := 0; i < semanas_en_produccion*7; i++ {
		if i%7 == 0 {
			idx += 1
		}
		if len(produccion) > 0 {
			percent = produccion[i-1].Data.Mortalidad
			percentReal = produccion[i-1].Data.MortalidadReal
		}

		mortality := percent - (variable_mortalidad_produccion / float64(semanas_en_produccion*7))
		mortalityReal := percentReal - (variable_mortalidad_produccion_ajustado / float64((semanas_en_produccion * 7)))
		aves := (float64(lote.Data.Aves) * mortality / 100)
		aves_real := (float64(lote.Data.Aves) * mortalityReal / 100)

		// Standar de produccion Real
		std_prod := _PROD[idx] - (_PROD[idx] * variable_produccion_huevos_totales / 100.0)
		// Standar de Aprovechamiento Real
		std_aprov := _APROV[idx] - (_APROV[idx] * variable_aprovechamiento_huevos / 100.0)
		// Standar de Nacimientos Real
		std_nac := _Nac[idx] - (_Nac[idx] * variable_nacimientos / 100.0)

		h_totales_real := aves * float64(std_prod) / 100
		h_totales := aves_real * float64(_PROD[idx]) / 100

		h_incubables_real := h_totales * float64(std_aprov) / 100
		h_incubables := h_totales_real * float64(_APROV[idx]) / 100

		p_nacidos_real := h_incubables * float64(std_nac) / 100
		p_nacidos := h_incubables_real * float64(_Nac[idx]) / 100

		data := DataProjected{
			Mortalidad:         mortality,
			MortalidadReal:     mortalityReal,
			Aves:               int64(math.Floor(aves)),
			AvesReal:           int64(math.Floor(aves_real)),
			HvosProducidos:     int64(math.Floor(h_totales)),
			HvosProducidosReal: int64(math.Floor(h_totales_real)),
			HvosIncubables:     int64(math.Floor(h_incubables)),
			HvosIncubablesReal: int64(math.Floor(h_incubables_real)),
			Nacimientos:        int64(math.Floor(p_nacidos / 2)),
			NacimientosReal:    int64(math.Floor(p_nacidos_real / 2)),
			Pollitos:           0,
			PollitosReal:       0,
		}
		fecha = fecha.AddDate(0, 0, 1)

		project := LoteProjection{
			Lote:         lote.Lote,
			Fecha:        fecha,
			Semana:       idx + 24, //+24 weeks in recria
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

func AbuelosProjectionTable(year string, dataType string, isProduccion bool) ([]string, [13][32]int64) {
	operation("ProjectionTable")
	if year == "" {
		year = fmt.Sprintf("%d", time.Now().Year())
	}
	if dataType == "" {
		dataType = "aves"
	}
	y, err := strconv.ParseInt(year, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	cols := []string{"Day", "jan", "feb", "mar", "apr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dec"}
	rows := [13][32]int64{}
	filter := fmt.Sprintf("{'year':{'$gte': %d}}", y-2)
	projections := ProjectAbuelos(filter)

	for _, el := range projections {
		if el.Year == y && el.EnProduccion == isProduccion {
			rows[0][el.Day-1] = int64(el.Day)
			switch dataType {
			case "aves":
				rows[el.Month][el.Day-1] += el.Data.Aves
				break
			case "hprod":
				rows[el.Month][el.Day-1] += el.Data.HvosProducidos
				break
			case "hinc":
				rows[el.Month][el.Day-1] += el.Data.HvosIncubables
				break
			case "nac":
				rows[el.Month][el.Day-1] += el.Data.Nacimientos
				break
			default:
				break
			}

		}
	}
	return cols, rows
}

func operation(name string) {
	fmt.Printf("%s Operation: %s\n", instance, name)
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
