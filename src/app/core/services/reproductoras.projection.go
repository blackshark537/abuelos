package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var REP_PROD = []float32{
	3, 20, 50, 82, 85, 86, 85.8, 84.8, 83.8, 82.8, 81.8, 80.8, 79.8, 78.8, 77.8, 76.7, 75.6, 74.5, 73.3, 72.1, 70.9,
	69.7, 68.5, 67.3, 66.1, 64.9, 63.6, 62.3, 61, 59.6, 58.2, 56.8, 55.4, 54, 52.5, 51, 49.5, 47.9, 46.3, 44.7, 43.1,
	42, 41, 40, 39.5, 38.5, 37.5, 36.5, 35.5, 34.5, 33.6, 33.6, 33.6, 33.6, 33.6, 33.6, 32.6, 32.6, 32.6, 32.6, 32.6, 32.6,
	32.6, 32.6, 32.6, 32.6, 32.6, 32.6, 32.6, 32.6, 32.6, 31.6, 31.6, 31.6, 31.6, 31.6,
}

var REP_APROV = []float32{
	0, 50, 75, 80, 92, 95, 96, 96.5, 97.5, 97.5, 97.5, 97.8, 97.8, 97.8, 97.8, 97.8, 97.8, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97,
	97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97,
	97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 96, 96, 96, 96, 96,
}

var REP_NAC = []float32{
	0, 72, 77.2, 80, 82.1, 83.8, 85.2, 86.4, 87.5, 88.5, 89.4, 90, 90, 89.9, 89.7, 89.5, 89.3, 89, 88.7, 88.4, 88.1, 87.7, 87.3,
	86.9, 86.5, 86.1, 85.6, 85.1, 84.6, 84.1, 83.6, 83.1, 82.6, 82, 81.5, 81, 80.6, 80.1, 79.7, 79.2, 78.7, 78.3, 77.8, 76.8,
	75.8, 74.8, 73.8, 72.8, 71.8, 70.8, 69.8, 68.8, 67.8, 66.8, 65.8, 64.8, 63.8, 62.8, 61.8, 60.8, 59, 58, 57, 56, 55, 54, 53, 52,
	51, 50, 49, 48, 47, 46, 45, 44,
}

var cache2 map[string][]LoteProjection = make(map[string][]LoteProjection)

func ProjectReproductoras(year string) []LoteProjection {
	t := time.Now()
	defer bench("LoteProjection", t)
	if cache2[year] != nil {
		return cache2[year]
	}
	y, err := strconv.ParseInt(year, 10, 64)
	handleErr(err)

	LoadVariables(TipoReproductoras)
	LoadVariables(TipoPollos)

	lotes := GetIncubations(fmt.Sprint(y - 2))
	lotes = append(lotes, GetIncubations(fmt.Sprint(y-1))...)
	lotes = append(lotes, GetIncubations(year)...)

	projection := []LoteProjection{}

	for _, lote := range lotes {
		recriaCh := make(chan []LoteProjection, 5)
		prodCh := make(chan []LoteProjection, 5)
		go func() {
			recriaCh <- getReprodRecria(lote)
		}()
		recria := <-recriaCh
		go func() {
			prodCh <- getReprodProduccion(recria[len(recria)-1])
		}()
		produccion := <-prodCh
		close(prodCh)
		close(recriaCh)
		projection = append(projection, recria...)
		projection = append(projection, produccion...)
	}
	cache2[year] = projection
	return projection
}

func getReprodRecria(lote Incubation) []LoteProjection {
	recrias := []LoteProjection{}
	data := DataProjected{
		Mortalidad:     100,
		MortalidadReal: 100,
		Aves:           int(lote.Nacimientos),
		AvesReal:       int(lote.Nacimientos),
	}
	fecha := lote.Fecha

	recria := LoteProjection{
		Lote:         int64(lote.Id),
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

		mortality := percent - (rep_var_mort_recria / float32(sem_en_recria*7))
		mortalityReal := percentReal - (rep_var_mort_recria_ajust / float32((sem_en_recria * 7)))

		aves := float32(lote.Nacimientos) * mortality / 100.0
		aves_real := float32(lote.Nacimientos) * mortalityReal / 100.0

		data := DataProjected{
			Mortalidad:     mortality,
			MortalidadReal: mortalityReal,
			Aves:           int(aves),
			AvesReal:       int(aves_real),
		}
		fecha = fecha.AddDate(0, 0, 1)

		recria := LoteProjection{
			Lote:         int64(lote.Id),
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

func getReprodProduccion(lote LoteProjection) []LoteProjection {
	produccion := []LoteProjection{}
	var idx int8 = 0
	fecha := lote.Fecha
	var percent float32 = 100
	var percentReal float32 = 100

	for i := 0; i < rep_sem_en_produccion*7; i++ {
		if i%7 == 0 {
			idx += 1
		}
		if len(produccion) > 0 {
			percent = produccion[i-1].Data.Mortalidad
			percentReal = produccion[i-1].Data.MortalidadReal
		}

		mortality := percent - (rep_var_mort_prod / float32(rep_sem_en_produccion*7))
		mortalityReal := percentReal - (rep_var_mort_prod_ajust / float32((rep_sem_en_produccion * 7)))
		aves := (float32(lote.Data.Aves) * mortality / 100)
		aves_real := (float32(lote.Data.Aves) * mortalityReal / 100)

		// Standar de produccion Real
		std_prod := REP_PROD[idx] - (REP_PROD[idx] * rep_var_prod_huevos_total / 100)
		// Standar de Aprovechamiento Real
		std_aprov := REP_APROV[idx] - (REP_APROV[idx] * rep_var_aprov_huevos / 100)
		// Standar de Nacimientos Real
		std_nac := REP_NAC[idx] - (REP_NAC[idx] * rep_var_nacimientos / 100)

		h_totales := aves * REP_PROD[idx] / 100
		h_totales_real := aves_real * std_prod / 100

		h_incubables := h_totales * REP_APROV[idx] / 100
		h_incubables_real := h_totales_real * std_aprov / 100

		p_nacidos := h_incubables * REP_NAC[idx] / 100
		p_nacidos_real := h_incubables_real * std_nac / 100

		data := DataProjected{
			Mortalidad:         mortality,
			MortalidadReal:     mortalityReal,
			Aves:               int(aves),
			AvesReal:           int(aves_real),
			HvosProducidos:     int(h_totales),
			HvosProducidosReal: int(h_totales_real),
			HvosIncubables:     int(h_incubables),
			HvosIncubablesReal: int(h_incubables_real),
			Nacimientos:        int(p_nacidos),
			NacimientosReal:    int(p_nacidos_real),
			Pollitos:           int(p_nacidos - (p_nacidos * rep_var_mort_pollitos / 100)),
			PollitosReal:       int(p_nacidos_real - (p_nacidos_real * rep_var_mort_pollitos_ajust / 100)),
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
			StdProd:      REP_PROD[idx],
			StdApro:      REP_APROV[idx],
			StdNac:       REP_NAC[idx],
			StdProdReal:  std_prod,
			StdAproReal:  std_aprov,
			StdNacReal:   std_nac,
			Data:         data,
		}

		produccion = append(produccion, project)
	}
	return produccion
}

func ReprodProjectionTable(year string, dataType string, isProduccion bool) ([]string, [13][32]int) {
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
	cols := []string{"Day", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dec"}
	rows := [13][32]int{}

	projections := ProjectReproductoras(fmt.Sprintf("%d", y))

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
			case "pollos":
				rows[el.Month][el.Day-1] += el.Data.Pollitos
				break
			case "vent-pollos-vivos":
				rows[el.Month][el.Day-1] += el.Data.Pollitos * var_venta_pollo_vivos / 100
				break
			case "pollos-vivos-mat":
				rows[el.Month][el.Day-1] += el.Data.Pollitos * var_pollos_vivos_mataderos / 100
				break
			case "libras-pollo":
				rows[el.Month][el.Day-1] += int(float32(el.Data.Pollitos) * var_peso_promedio)
				break
			case "lb-pollo-proc":
				rows[el.Month][el.Day-1] += int(float32(el.Data.Pollitos)*var_peso_promedio) * var_lb_pollo_procesados / 100
				break
			case "kl-pollo-proc":
				lb := int(float32(el.Data.Pollitos)*var_peso_promedio) * var_lb_pollo_procesados / 100
				rows[el.Month][el.Day-1] += int(float32(lb) / 2.20462)
				break
			case "tn-pollo-proc":
				lb := int(float32(el.Data.Pollitos)*var_peso_promedio) * var_lb_pollo_procesados / 100
				rows[el.Month][el.Day-1] += int(float32(lb) / (2.20462 * 1000))
				break
			case "lb-pollo-mat":
				rows[el.Month][el.Day-1] += int(float32(el.Data.Pollitos)*var_peso_promedio) * var_lb_pollo_matadero / 100
				break
			default:
				break
			}

		}
	}
	return cols, rows
}

// [Warning] For CLI Use Only
func ListReproductoras(lote string) {
	filters := fmt.Sprintf(`{"numero" : {"$eq": %v}}`, lote)
	results := ProjectReproductoras(filters)
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
func ReprodTable(year string, dataType string, isProduccion bool) {
	cols, rows := ReprodProjectionTable(year, dataType, isProduccion)
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New(cols[0], cols[1], cols[2], cols[3], cols[4], cols[5], cols[6], cols[7], cols[8], cols[9], cols[10], cols[11], cols[12])
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	fmt.Printf("Data: %s Year: %v\n", dataType, year)
	for d := 0; d < len(rows[0]); d++ {
		tbl.AddRow(rows[0][d], rows[1][d], rows[2][d], rows[3][d], rows[4][d], rows[5][d], rows[6][d], rows[7][d], rows[8][d], rows[9][d], rows[10][d], rows[11][d], rows[12][d])
	}

	tbl.Print()
}
