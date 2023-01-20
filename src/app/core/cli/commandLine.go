package CommandLine

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blackshark537/dataprod/src/app/core/config"
	"github.com/blackshark537/dataprod/src/app/core/entities"
	portin "github.com/blackshark537/dataprod/src/app/core/port_in"
	"github.com/blackshark537/dataprod/src/app/core/services"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	port       string
	collection string
	objectId   string
	data       string
	filter     string
)

var instance = color.MagentaString("[CLI]:")

var Commands []*cli.Command = []*cli.Command{
	{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Start the server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       "3000",
				Usage:       "Server port listener",
				Destination: &port,
			},
		},
		Action: serverStart,
	},
	{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "Show a collection list",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "from",
				Usage:       "Database Collection",
				Aliases:     []string{"f"},
				Destination: &collection,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "where",
				Usage:       "Filter",
				Aliases:     []string{"w"},
				Destination: &filter,
				Value:       "",
			},
		},
		Action: listTable,
	},
	{
		Name:    "insert",
		Aliases: []string{"i"},
		Usage:   "Insert object into a given collection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "into",
				Usage:       "Database Collection",
				Aliases:     []string{"to"},
				Destination: &collection,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "data",
				Usage:       "Data to insert",
				Aliases:     []string{"d"},
				Destination: &data,
				Required:    true,
			},
		},
		Action: insertIntoTable,
	},
	{
		Name:    "delete",
		Aliases: []string{"del"},
		Usage:   "Delete an item from a given collection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "from",
				Usage:       "Database Collection",
				Aliases:     []string{"f"},
				Destination: &collection,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "ObjectID",
				Usage:       "Database Id",
				Aliases:     []string{"id"},
				Destination: &objectId,
				Required:    true,
			},
		},
		Action: deleteFromTable,
	},
}

// Create and returns the cli app instance
func New() *cli.App {
	website := color.CyanString("WEBSITE: https://dataprod.cloud")
	support := color.GreenString("SUPPORT: support@dataprod.cloud")
	cli.AppHelpTemplate = fmt.Sprintf(`%s

	%s
	
	%s
	
	`, cli.AppHelpTemplate, website, support)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"v"},
		Usage:   "print only the version",
	}

	AppName := color.MagentaString("Dataprod Core")
	AppDesc := color.MagentaString("A simple CLI program to manage your dataprod system")
	return &cli.App{
		Name:     AppName,
		Version:  "v" + config.VERSION,
		Usage:    AppDesc,
		Commands: Commands,
	}
}

func serverStart(ctx *cli.Context) error {
	fmt.Printf("%s Server Starting on port: %s\n", instance, port)
	api := portin.ApiAdapter{}
	return api.ForRoot(port)
}

func listTable(ctx *cli.Context) error {
	fmt.Printf("%s List Colection: %s\n", instance, collection)
	switch collection {
	case "empresas":
		empresa := entities.Empresa{}
		empresa.List(filter)
		return nil
	case "lotes":
		lotes := entities.Lote{}
		lotes.List(filter)
		return nil
	case "projection":
		projection := services.Projection{}
		projection.List(filter)
		return nil
	default:
		return errors.New("Not a valid collection")
	}
}

func insertIntoTable(ctx *cli.Context) error {
	fmt.Printf("%s Insert Into Colection: %s\n", instance, collection)
	fmt.Printf("%s Data: %s\n", instance, data)
	switch collection {
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
		return nil
	case "alimentos":
		return nil
	default:
		return errors.New("No valid collection")
	}
}

func deleteFromTable(ctx *cli.Context) error {
	fmt.Printf("%s Delete From Colection: %s\n", instance, collection)
	fmt.Println("_Id: ", objectId)
	switch collection {
	case "empresas":
		empresa := entities.Empresa{}
		return empresa.Delete(objectId)
	default:
		return errors.New("No collection matched")
	}

}
