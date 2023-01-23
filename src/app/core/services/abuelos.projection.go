package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/blackshark537/dataprod/src/app/core/config"
	"github.com/blackshark537/dataprod/src/app/core/entities"
	"github.com/fatih/color"
)

type DataProjected struct {
	Mortalidad         float32
	MortalidadReal     float32
	Aves               int
	AvesReal           int
	HvosProducidos     int
	HvosProducidosReal int
	HvosIncubables     int
	HvosIncubablesReal int
	Nacimientos        int
	NacimientosReal    int
	Pollitos           int
	PollitosReal       int
}

type LoteProjection struct {
	Fecha        time.Time
	Semana       int8
	Day          int16
	Month        int8
	EnProduccion bool
	Year         int64
	Lote         int64
	StdProd      float32
	StdApro      float32
	StdNac       float32
	StdProdReal  float32
	StdAproReal  float32
	StdNacReal   float32
	Data         DataProjected
}

type TableProjection struct {
	Day   int16
	Month []DataProjected
	Year  int64
}

var cache map[string][]LoteProjection

var ABU_PROD = []float32{
	0, 8.40, 28.50, 48, 60.30, 69.70, 74.10, 77.50, 77.40, 77.10, 76.30, 74.80, 73.10, 71.40, 69.30, 67.20, 65.10, 63.30,
	61.10, 59, 57.10, 55.30, 53.50, 51.60, 49.70, 48, 46.20, 44.80, 43.60, 42.40, 41.40, 40.50, 39.50, 38.60, 37.70, 36.80,
	36, 35.10, 34.20, 33.30, 32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22,
}
var ABU_APROV = []float32{
	0, 50, 78, 84, 88, 90, 92, 94, 94, 95, 95, 95, 95, 95, 95, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 95, 95, 95, 95, 95, 95,
	94, 94, 94, 94, 94, 93, 93, 93, 93, 93, 93, 92, 92, 92, 92, 91, 91, 90,
}
var ABU_NAC = []float32{
	0, 0, 55, 65, 70, 75, 80, 81, 82, 83, 83, 82, 82, 82, 81, 81, 81, 80, 80, 80, 80, 79, 79, 78, 78, 77, 77, 77, 76, 76, 75, 75, 74, 74,
	73, 73, 72, 70, 69, 68, 67, 66, 65, 64, 63, 62, 61, 60, 59, 58, 57,
}

const (
	sem_en_recria = 24
)

var (
	abu_sem_en_produccion     int     = 41
	abu_var_mort_recria       float32 = 4.0
	abu_var_mort_recria_ajust float32 = 5.0
	abu_var_mort_prod         float32 = 10.0
	abu_var_mort_prod_ajust   float32 = 12.0
	abu_var_prod_huevos_total float32 = 0
	abu_var_aprov_huevos      float32 = 0
	abu_var_nacimientos       float32 = 0
)

var instance string = color.MagentaString("[Projection]:")

// [Warning] For CLI Use Only
func ListAbuelos(lote string) {
	filters := fmt.Sprintf(`{"numero" : {"$eq": %v}}`, lote)
	results := ProjectAbuelos(filters)
	fmt.Printf("%s %v Items\n", color.MagentaString("[Results]:"), len(results))
	fmt.Printf("%s %v\n", color.MagentaString("[Lote]:"), lote)
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("| Lote\t | Sem\t | mortalidad\t | Aves\t | STD\t | Produc\t | STD\t | Incub | STD\t | Nac\t |")
	fmt.Println("-------------------------------------------------------------------")
	for _, el := range results {
		fmt.Printf("| %d\t | %f\t | %d\t | %f\t | %d\t  | %f\t | %d\t | %f\t | %d\t |\n", el.Semana, el.Data.Mortalidad, el.Data.Aves, el.StdProd, el.Data.HvosProducidos, el.StdApro, el.Data.HvosIncubables, el.StdNac, el.Data.Nacimientos)
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
	t := time.Now()
	defer bench("LoteProjection", t)
	if cache[filters] != nil {
		return cache[filters]
	}
	loteEntity := new(entities.Lote)
	lotes, err := loteEntity.GetAll(filters)
	if err != nil {
		handleErr(err)
	}
	projection := []LoteProjection{}

	for _, lote := range lotes {
		recriaCh := make(chan []LoteProjection, 500)
		prodCh := make(chan []LoteProjection, 500)
		go func() {
			recriaCh <- getAbuelosRecria(lote)
		}()
		recria := <-recriaCh
		go func() {
			prodCh <- getAbuelosProduccion(recria[len(recria)-1])
		}()
		produccion := <-prodCh
		close(prodCh)
		close(recriaCh)
		projection = append(projection, recria...)
		projection = append(projection, produccion...)
	}
	cache = make(map[string][]LoteProjection, len(projection))
	cache[filters] = projection
	return projection
}

func getAbuelosRecria(lote entities.Lote) []LoteProjection {
	recrias := []LoteProjection{}
	data := DataProjected{
		Mortalidad:     100,
		MortalidadReal: 100,
		Aves:           int(lote.Hembras),
		AvesReal:       int(lote.Hembras),
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
	var percent float32 = 0
	var percentReal float32 = 0

	for i := 0; i < sem_en_recria*7; i++ {
		if i%7 == 0 {
			idx += 1
		}
		if len(recrias) > 0 {
			percent = recrias[i].Data.Mortalidad
			percentReal = recrias[i].Data.MortalidadReal
		}

		mortality := percent - (abu_var_mort_recria / float32(sem_en_recria*7))
		mortalityReal := percentReal - (abu_var_mort_recria_ajust / float32((sem_en_recria * 7)))

		aves := float32(lote.Hembras) * mortality / 100.0
		aves_real := float32(lote.Hembras) * mortalityReal / 100.0

		data := DataProjected{
			Mortalidad:     mortality,
			MortalidadReal: mortalityReal,
			Aves:           int(aves),
			AvesReal:       int(aves_real),
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
	fecha := lote.Fecha
	var percent float32 = 100
	var percentReal float32 = 100

	for i := 0; i < abu_sem_en_produccion*7; i++ {
		if i%7 == 0 {
			idx += 1
		}
		if len(produccion) > 0 {
			percent = produccion[i-1].Data.Mortalidad
			percentReal = produccion[i-1].Data.MortalidadReal
		}

		mortality := percent - (abu_var_mort_prod / float32(abu_sem_en_produccion*7))
		mortalityReal := percentReal - (abu_var_mort_prod_ajust / float32((abu_sem_en_produccion * 7)))
		aves := (float32(lote.Data.Aves) * mortality / 100)
		aves_real := (float32(lote.Data.Aves) * mortalityReal / 100)

		// Standar de produccion Real
		std_prod := ABU_PROD[idx] - (ABU_PROD[idx] * abu_var_prod_huevos_total / 100)
		// Standar de Aprovechamiento Real
		std_aprov := ABU_APROV[idx] - (ABU_APROV[idx] * abu_var_aprov_huevos / 100)
		// Standar de Nacimientos Real
		std_nac := ABU_NAC[idx] - (ABU_NAC[idx] * abu_var_nacimientos / 100)

		h_totales := aves * ABU_PROD[idx] / 100
		h_totales_real := aves_real * std_prod / 100

		h_incubables := h_totales * ABU_APROV[idx] / 100
		h_incubables_real := h_totales_real * std_aprov / 100

		p_NACidos := h_incubables * ABU_NAC[idx] / 100
		p_NACidos_real := h_incubables_real * std_nac / 100

		data := DataProjected{
			Mortalidad:         mortality,
			MortalidadReal:     mortalityReal,
			Aves:               int(aves),
			AvesReal:           int(aves_real),
			HvosProducidos:     int(h_totales),
			HvosProducidosReal: int(h_totales_real),
			HvosIncubables:     int(h_incubables),
			HvosIncubablesReal: int(h_incubables_real),
			Nacimientos:        int(p_NACidos / 2),
			NacimientosReal:    int(p_NACidos_real / 2),
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
			StdProd:      ABU_PROD[idx],
			StdApro:      ABU_APROV[idx],
			StdNac:       ABU_NAC[idx],
			StdProdReal:  std_prod,
			StdAproReal:  std_aprov,
			StdNacReal:   std_nac,
			Data:         data,
		}

		produccion = append(produccion, project)
	}
	return produccion
}

func AbuelosProjectionTable(year string, dataType string, isProduccion bool) ([]string, [13][32]int) {
	t := time.Now()
	defer bench("ProjectionTable", t)
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
	rows := [13][32]int{}
	filter := fmt.Sprintf(`{"year":{"$gte": %d}}`, y-2)
	projections := ProjectAbuelos(filter)

	for _, el := range projections {
		if el.Year == y && el.EnProduccion == isProduccion {
			rows[0][el.Day-1] = int(el.Day)
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

func bench(name string, t time.Time) {
	if config.IsBench {
		fmt.Printf("%s Operation: %s - %v mili secs\n", instance, name, time.Since(t).Milliseconds())
	}
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
