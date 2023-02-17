package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/blackshark537/abuelos/src/app/core/entities"
	portin "github.com/blackshark537/abuelos/src/app/core/port_in"
	"github.com/blackshark537/abuelos/src/app/core/services"
	"github.com/fatih/color"
)

type CliCommand struct {
	Collection string
}

var emp *entities.Empresa = new(entities.Empresa)
var lot *entities.Lote = new(entities.Lote)
var vars *entities.Variable = new(entities.Variable)
var insumos *entities.Insumo = new(entities.Insumo)

var instance = color.MagentaString("[CLI]:")

func (cmd *CliCommand) ServerStart(port string) error {
	fmt.Printf("%s Server Starting on port: %s\n", instance, port)
	return new(portin.ApiPort).ForRoot(port)
}

func (cmd *CliCommand) SetVariables() {
	services.CreateVariables(services.TipoAbuelos)
	services.CreateVariables(services.TipoReproductoras)
	services.CreateVariables(services.TipoPollos)
}

func (cmd *CliCommand) ListAbuelos(loteId string) {
	services.ListAbuelos(loteId)
}

func (cmd *CliCommand) ListReproductoras(loteId string) {
	services.ListReproductoras(loteId)
}

func (cmd *CliCommand) ListTable(filter string) error {
	switch cmd.Collection {
	case "empresas":
		emp.List(filter)
		return nil
	case "lotes":
		lot.List(filter)
		return nil
	case "insumos":
		insumos.List(filter)
		return nil
	case "variables":
		vars.List(filter)
		return nil
	default:
		return noMatch()
	}
}

func (cmd *CliCommand) ProjectTable(year string, data string, prod bool) error {
	switch cmd.Collection {
	case "abuelos":
		fmt.Printf("%s List Colection: %s\n", instance, cmd.Collection)
		services.AbuelosTable(year, data, prod)
		return nil
	case "incubations":
		services.LisIncubations(year)
		return nil
	case "reproductoras":
		services.ReprodTable(year, data, prod)
		return nil
	case "pollos":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	case "vent-pollos-vivos":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	case "pollos-vivos-mat":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	case "libras-pollo":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	case "lb-pollo-proc":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	case "kl-pollo-proc":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	case "tn-pollo-proc":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	case "lb-pollo-mat":
		services.ReprodTable(year, cmd.Collection, prod)
		return nil
	default:
		return noMatch()
	}
}

func (cmd *CliCommand) InsertIntoCollection(data string) error {
	fmt.Printf("%s Insert Into Colection: %s\n", instance, cmd.Collection)
	fmt.Printf("%s Data: %s\n", instance, data)
	switch cmd.Collection {
	case "empresas":
		empresa := new(entities.Empresa)
		gpoin := new(entities.Geopoint)
		json.Unmarshal([]byte(data), empresa)
		json.Unmarshal([]byte(data), gpoin)
		res, err := empresa.Save(gpoin)
		if err != nil {
			return err
		}
		fmt.Printf("Object %v successfully created!\n", res)
		return nil
	case "lotes":
		json.Unmarshal([]byte(data), lot)
		res, err := lot.Save()
		if err != nil {
			return err
		}
		fmt.Printf("Object %v successfully created!\n", res)
		return nil
	case "insumos":
		json.Unmarshal([]byte(data), insumos)
		res, err := insumos.Save()
		if err != nil {
			return err
		}
		fmt.Printf("Object %v successfully created!\n", res)
		return nil
	default:
		return noMatch()
	}
}

func (cmd *CliCommand) InsertFromFile(path string) error {
	file, err := os.ReadFile(fmt.Sprintf("./%s", path))
	if err != nil {
		return err
	}
	data := string(file)

	switch cmd.Collection {
	case "lotes":
		lotes := []entities.Lote{}
		if err := json.Unmarshal([]byte(data), &lotes); err != nil {
			return err
		}
		for _, el := range lotes {
			if _, err = el.Save(); err != nil {
				break
			}
		}
		return err
	case "empresas":
		empresas := []entities.Empresa{}
		if err := json.Unmarshal([]byte(data), &empresas); err != nil {
			return err
		}
		for _, el := range empresas {
			geo := new(entities.Geopoint)
			json.Unmarshal([]byte(data), &geo)
			if _, err = el.Save(geo); err != nil {
				break
			}
		}
		return err
	case "isumos":
		insumos := []entities.Insumo{}
		if err := json.Unmarshal([]byte(data), &insumos); err != nil {
			return err
		}
		for _, insumo := range insumos {
			if _, err = insumo.Save(); err != nil {
				break
			}
		}
		return err
	default:
		return noMatch()
	}
}

func (cmd *CliCommand) ClearCollection(filter string) error {
	switch cmd.Collection {
	case "lotes":
		return lot.DeleteMany(filter)
	case "empresas":
		return emp.DeleteMany(filter)
	case "insumos":
		return insumos.DeleteMany(filter)
	default:
		return noMatch()
	}
}

func (cmd *CliCommand) RemoveFromCollection(objectId string) error {
	fmt.Printf("%s Delete From Colection: %s\n", instance, cmd.Collection)
	fmt.Println("_Id: ", objectId)
	switch cmd.Collection {
	case "empresas":
		return emp.Delete(objectId)
	case "lotes":
		return lot.Delete(objectId)
	case "insumos":
		return insumos.Delete(objectId)
	default:
		return noMatch()
	}
}

func noMatch() error {
	return errors.New("Sorry No collection matched!")
}
