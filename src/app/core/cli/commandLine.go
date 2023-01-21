package CommandLine

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/blackshark537/dataprod/src/app/core/config"
	"github.com/blackshark537/dataprod/src/app/core/entities"
	portin "github.com/blackshark537/dataprod/src/app/core/port_in"
	"github.com/blackshark537/dataprod/src/app/core/services"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	port         string
	collection   string
	objectId     string
	data         string
	filter       string
	path         string
	lote         string
	isProduccion string
	year         string
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
				Usage:       `Apply Filters Ex: '{"property": {"$regex": "myValue"}}'`,
				Aliases:     []string{"w"},
				Destination: &filter,
				Value:       "",
			},
		},
		Action: listTable,
	},
	{
		Name:    "project",
		Aliases: []string{"prj"},
		Usage:   "Project a lote",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lote",
				Usage:       "Lote number to project",
				Aliases:     []string{"l"},
				Destination: &lote,
				Required:    true,
			},
		},
		Action: projectLote,
	},
	{
		Name:    "table",
		Aliases: []string{"tb"},
		Usage:   "Show a data projection table",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "type",
				Usage:       "Type of chicken ('abuelos', 'reproductoras', 'pollos')",
				Aliases:     []string{"t"},
				Destination: &collection,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "select",
				Usage:       "Type of data to show ('aves', 'hprod', 'hinc', 'nac')",
				Aliases:     []string{"s"},
				Value:       "aves",
				Destination: &data,
			},
			&cli.StringFlag{
				Name:        "year",
				Usage:       "Data Year",
				Aliases:     []string{"y"},
				Value:       fmt.Sprintf("%d", time.Now().Year()),
				Destination: &year,
			},
			&cli.StringFlag{
				Name:        "produccion",
				Usage:       "If True shows birds in Produccion else in Recria",
				Aliases:     []string{"p"},
				Value:       "yes",
				Destination: &isProduccion,
			},
		},
		Action: projectTable,
	},
	{
		Name:    "insert",
		Aliases: []string{"i"},
		Usage:   "Insert an object into a given collection",
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
		Name:    "bulk",
		Aliases: []string{"b"},
		Usage:   "Insert many objects into a given collection from a .json file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "into",
				Usage:       "Database Collection",
				Aliases:     []string{"to"},
				Destination: &collection,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "path",
				Usage:       "File path of data to bulk",
				Aliases:     []string{"p"},
				Destination: &path,
				Required:    true,
			},
		},
		Action: InsertMany,
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
	{
		Name:    "clear",
		Aliases: []string{"clr"},
		Usage:   "Delete all items from a given collection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "from",
				Usage:       "Database Collection",
				Aliases:     []string{"f"},
				Destination: &collection,
				Required:    true,
			},
		},
		Action: DeleteAll,
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

func projectLote(ctx *cli.Context) error {
	collection = "lotes"
	fmt.Printf("%s List Colection: %s\n", instance, collection)
	services.ListAbuelos(lote)
	return nil
}

func projectTable(ctx *cli.Context) error {
	switch collection {
	case "abuelos":

		fmt.Printf("%s List Colection: %s\n", instance, collection)
		prod := false
		if isProduccion == "y" || isProduccion == "yes" {
			prod = true
		}
		services.AbuelosTable(year, data, prod)
		return nil
	case "reproductoras":
		return noMatch()
	case "pollos":
		return noMatch()
	default:
		return noMatch()
	}
}

func listTable(ctx *cli.Context) error {
	switch collection {
	case "empresas":
		e := new(entities.Empresa)
		e.List(filter)
		return nil
	case "lotes":
		e := new(entities.Lote)
		e.List(filter)
		return nil
	default:
		return noMatch()
	}
}

func List(e entities.EntityList, filter string) {
	e.List(filter)
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
		return noMatch()
	}
}

func InsertMany(ctx *cli.Context) error {
	file, err := os.ReadFile(fmt.Sprintf("./%s", path))
	if err != nil {
		return err
	}
	data = string(file)
	switch collection {
	case "lotes":
		lotes := []entities.Lote{}
		err := json.Unmarshal([]byte(data), &lotes)
		if err != nil {
			return err
		}
		for _, el := range lotes {
			el.Save()
		}
		return nil
	case "empresas":
		empresas := []entities.Empresa{}
		err := json.Unmarshal([]byte(data), &empresas)
		if err != nil {
			return err
		}
		for _, el := range empresas {
			geo := new(entities.Geopoint)
			json.Unmarshal([]byte(data), &geo)
			el.Save(geo)
		}
		return nil
	default:
		return noMatch()
	}
}

func DeleteAll(ctx *cli.Context) error {
	switch collection {
	case "lotes":
		l := new(entities.Lote)
		return l.DeleteMany(filter)
	case "empresas":
		l := new(entities.Lote)
		return l.DeleteMany(filter)
	default:
		return noMatch()
	}
}

func deleteFromTable(ctx *cli.Context) error {
	fmt.Printf("%s Delete From Colection: %s\n", instance, collection)
	fmt.Println("_Id: ", objectId)
	switch collection {
	case "empresas":
		e := new(entities.Empresa)
		return e.Delete(objectId)
	case "lotes":
		e := new(entities.Lote)
		return e.Delete(objectId)
	default:
		return noMatch()
	}
}

func noMatch() error {
	return errors.New("Sorry No collection matched!")
}

func handleErr(e error) error {
	if e != nil {
		return e
	}
	return nil
}
