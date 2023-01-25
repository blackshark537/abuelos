package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/blackshark537/dataprod/src/app/core/entities"
)

const (
	TipoAbuelos       = "abuelos"
	TipoReproductoras = "reproductoras"
	TipoPollos        = "pollos"
)

// Abuelos
var (
	abu_sem_en_produccion     int     = 44
	abu_var_mort_recria       float32 = 4.0
	abu_var_mort_recria_ajust float32 = 5.0
	abu_var_mort_prod         float32 = 10.0
	abu_var_mort_prod_ajust   float32 = 12.0
	abu_var_prod_huevos_total float32 = 0
	abu_var_aprov_huevos      float32 = 0
	abu_var_nacimientos       float32 = 0
)

// Reproductoras
var (
	rep_sem_en_produccion       int     = 44
	rep_var_mort_recria         float32 = 5.0
	rep_var_mort_recria_ajust   float32 = 7.0
	rep_var_mort_prod           float32 = 10.0
	rep_var_mort_prod_ajust     float32 = 13.0
	rep_var_mort_pollitos       float32 = 10.72
	rep_var_mort_pollitos_ajust float32 = 13.7
	rep_var_prod_huevos_total   float32 = 11.75
	rep_var_aprov_huevos        float32 = 1
	rep_var_nacimientos         float32 = 3
)

// Pollos
var (
	var_venta_pollo_vivos      int     = 48
	var_pollos_vivos_mataderos int     = 52
	var_lb_pollo_procesados    int     = 70
	var_lb_pollo_matadero      int     = 52
	var_peso_promedio          float32 = 4.8
)

func LoadVariables(tipo string) {
	t := time.Now()
	defer bench("LoadVariables", t)
	vars := new(entities.Variable)
	filter := fmt.Sprintf(`{"tipo": {"$eq": %v}}`, tipo)
	err := vars.FindOne(filter)
	if err != nil {
		handleErr(err)
	}
	switch tipo {
	case "abuelos":
		abu_sem_en_produccion = vars.SemanasVida - sem_en_recria
		abu_var_mort_recria = vars.MortalidadRecria
		abu_var_mort_recria_ajust = vars.MortalidadRecria + getTotal(vars.MortalidadRecriaAjustado)
		abu_var_mort_prod = vars.MortalidadProduccion
		abu_var_mort_prod_ajust = vars.MortalidadProduccion + getTotal(vars.MortalidadProduccionAjustado)
		return
	case "reproductoras":
		rep_sem_en_produccion = vars.SemanasVida - sem_en_recria
		rep_var_mort_recria = vars.MortalidadRecria
		rep_var_mort_recria_ajust = vars.MortalidadRecria + getTotal(vars.MortalidadRecriaAjustado)
		rep_var_mort_prod = vars.MortalidadProduccion
		rep_var_mort_prod_ajust = vars.MortalidadProduccion + getTotal(vars.MortalidadProduccionAjustado)
		rep_var_prod_huevos_total = getTotal(vars.ProduccionHuevosAjustado)
		rep_var_aprov_huevos = getTotal(vars.AprovechamientoAjustado)
		rep_var_nacimientos = getTotal(vars.NacimientosAjustado)
		return
	case "pollos":
		rep_var_mort_pollitos = vars.MortalidadPollitos
		rep_var_mort_pollitos_ajust = vars.MortalidadPollitos + getTotal(vars.MortalidadPollitosAjustado)
		var_venta_pollo_vivos = vars.VentaPollosVivos
		var_pollos_vivos_mataderos = vars.PollosVivosMataderos
		var_lb_pollo_procesados = vars.LbPolloProcesados
		var_lb_pollo_matadero = vars.LbPolloMatadero
		var_peso_promedio = vars.PesoPromedio
		return
	default:
		handleErr(errors.New(fmt.Sprintf("Variable Tipe of %s Unknow", tipo)))
	}
	return
}

func getTotal(values []entities.Ajuste) float32 {
	var acc float32 = 0
	for _, el := range values {
		acc += el.Value
	}
	return acc
}

func CreateVariables(tipo string) {
	vars := new(entities.Variable)
	gumboro := entities.Ajuste{
		Name:  "gumboro",
		Prop:  "gumboro",
		Value: 0,
	}
	newcastle := entities.Ajuste{
		Name:  "newcastle",
		Prop:  "newcastle",
		Value: 0,
	}
	clostridium := entities.Ajuste{
		Name:  "Enteritis ó clostridium",
		Prop:  "clostridium",
		Value: 0,
	}
	influenza := entities.Ajuste{
		Name:  "Influenza Aviar",
		Prop:  "influenza",
		Value: 0,
	}
	anemia := entities.Ajuste{
		Name:  "Anemia Infecciosa",
		Prop:  "anemia",
		Value: 0,
	}
	viruela := entities.Ajuste{
		Name:  "Viruela",
		Prop:  "viruela",
		Value: 0,
	}
	nutricion := entities.Ajuste{
		Name:  "Deficiencia Nutricional",
		Prop:  "defi_nutricional",
		Value: 0,
	}
	clima := entities.Ajuste{
		Name:  "Fenómeno Natural ó Climatológico",
		Prop:  "clima",
		Value: 0,
	}
	otros := entities.Ajuste{
		Name:  "Otros Virus ú Enfermedades",
		Prop:  "otros",
		Value: 0,
	}

	baja_postura := entities.Ajuste{
		Name:  "Sindrome De Baja Postura",
		Prop:  "baja_postura",
		Value: 0,
	}
	coriza := entities.Ajuste{
		Name:  "Coriza Infecciosa Aviar",
		Prop:  "coriza",
		Value: 0,
	}
	colera := entities.Ajuste{
		Name:  "Cólera",
		Prop:  "colera",
		Value: 0,
	}
	pneumovirus := entities.Ajuste{
		Name:  "Pneumovirus",
		Prop:  "pneumovirus",
		Value: 0,
	}

	gallisepticum := entities.Ajuste{
		Name:  "Mycoplasma Gallisepticum",
		Prop:  "gallisepticum",
		Value: 0,
	}
	psynoviae := entities.Ajuste{
		Name:  "Mycoplasma Synoviae",
		Prop:  "synoviae",
		Value: 0,
	}

	huevos_piso := entities.Ajuste{
		Name:  "Huevos De Piso",
		Prop:  "huevos_piso",
		Value: 0,
	}

	cama := entities.Ajuste{
		Name:  "Condiciones De La Cama",
		Prop:  "cama",
		Value: 0,
	}

	nidal := entities.Ajuste{
		Name:  "Manejo Del Nidal",
		Prop:  "nidal",
		Value: 0,
	}

	recoleccion := entities.Ajuste{
		Name:  "Programa De Recolección De Huevos",
		Prop:  "recoleccion",
		Value: 0,
	}

	gallos := entities.Ajuste{
		Name:  "Falta De Gallos",
		Prop:  "gallos",
		Value: 0,
	}

	sobrepeso := entities.Ajuste{
		Name:  "Sobrepeso",
		Prop:  "sobrepeso",
		Value: 0,
	}
	huevos_sucios := entities.Ajuste{
		Name:  "Postura De Huevo y Suelo (Huevos Sucios)",
		Prop:  "huevos_sucios",
		Value: 0,
	}
	contaminacion := entities.Ajuste{
		Name:  "Contaminación",
		Prop:  "contaminacion",
		Value: 0,
	}
	encefalomielitis := entities.Ajuste{
		Name:  "Encefalomielitis Aviar",
		Prop:  "encefalomielitis",
		Value: 0,
	}

	hepatitis := entities.Ajuste{
		Name:  "Hepatitis",
		Prop:  "hepatitis",
		Value: 0,
	}

	coccidiosis := entities.Ajuste{
		Name:  "Coccidiosis Aviar",
		Prop:  "coccidiosis",
		Value: 0,
	}
	micoplasma := entities.Ajuste{
		Name:  "Micoplasma",
		Prop:  "micoplasma",
		Value: 0,
	}
	calefaccion := entities.Ajuste{
		Name:  "Calefacción",
		Prop:  "calefaccion",
		Value: 0,
	}

	switch tipo {
	case "abuelos":

		vars.SemanasVida = 24 + abu_sem_en_produccion
		vars.MortalidadRecria = abu_var_mort_recria
		vars.MortalidadRecriaAjustado = []entities.Ajuste{
			gumboro,
			newcastle,
			clostridium,
			influenza,
			anemia,
			viruela,
			nutricion,
			clima,
			otros,
		}
		vars.MortalidadProduccion = abu_var_mort_prod
		vars.MortalidadProduccionAjustado = []entities.Ajuste{
			baja_postura,
			newcastle,
			coriza,
			colera,
			influenza,
			pneumovirus,
			gallisepticum,
			psynoviae,
			anemia,
			viruela,
			nutricion,
			clima,
			otros,
		}
		break
	case "reproductoras":
		vars.SemanasVida = sem_en_recria + rep_sem_en_produccion
		vars.MortalidadRecria = rep_var_mort_recria
		vars.MortalidadRecriaAjustado = []entities.Ajuste{}
		vars.MortalidadProduccion = rep_var_mort_prod
		vars.MortalidadProduccionAjustado = []entities.Ajuste{
			baja_postura,
			newcastle,
			coriza,
			colera,
			influenza,
			pneumovirus,
			gallisepticum,
			psynoviae,
			anemia,
			viruela,
			nutricion,
			clima,
			otros,
		}

		vars.ProduccionHuevosAjustado = []entities.Ajuste{
			baja_postura,
			newcastle,
			coriza,
			colera,
			influenza,
			pneumovirus,
			gallisepticum,
			psynoviae,
			anemia,
			viruela,
			nutricion,
			clima,
			otros,
		}

		vars.AprovechamientoAjustado = []entities.Ajuste{
			huevos_piso,
			cama,
			nidal,
			baja_postura,
			recoleccion,
			nutricion,
			clima,
		}

		vars.NacimientosAjustado = []entities.Ajuste{
			gallos,
			sobrepeso,
			huevos_sucios,
			contaminacion,
			encefalomielitis,
			clima,
			nutricion,
		}
		break

	case "pollos":
		vars.MortalidadPollitos = rep_var_mort_pollitos
		vars.MortalidadPollitosAjustado = []entities.Ajuste{
			newcastle,
			gumboro,
			hepatitis,
			influenza,
			anemia,
			coccidiosis,
			micoplasma,
			calefaccion,
			clima,
			nutricion,
		}
		vars.VentaPollosVivos = var_venta_pollo_vivos
		vars.PollosVivosMataderos = var_pollos_vivos_mataderos
		vars.LbPolloProcesados = var_lb_pollo_procesados
		vars.LbPolloMatadero = var_lb_pollo_matadero
		vars.PesoPromedio = var_peso_promedio
		break
	default:
		handleErr(errors.New(fmt.Sprintf("Variable Tipe of %s Unknow", tipo)))
		return
	}
	vars.Tipo = tipo
	_, err := vars.Save()
	handleErr(err)
	return
}
