package CommandLine

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

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
	path       string
	lote       string
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
				Destination: &data,
				Required:    true,
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
	service := new(services.AbuelosProjection)
	fmt.Printf("%s List Colection: %s\n", instance, collection)
	service.List(lote)
	return nil
}

func projectTable(ctx *cli.Context) error {
	switch collection {
	case "abuelos":
		service := new(services.AbuelosProjection)
		fmt.Printf("%s List Colection: %s\n", instance, collection)
		service.Table(2023, data, true)
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
	sub := make(map[string]entities.EntityList)
	sub["empresas"] = new(entities.Empresa)
	sub["lotes"] = new(entities.Lote)
	if sub[collection] == nil {
		return noMatch()
	}
	List(sub[collection], filter)
	return nil
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
	if collection == "lotes" {
		lotes := []entities.Lote{}
		json.Unmarshal([]byte(data), &lotes)
		for _, el := range lotes {
			el.Save()
		}
	}
	return nil
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
	sub := make(map[string]entities.EntityDelete)
	sub["empresas"] = new(entities.Empresa)
	sub["lotes"] = new(entities.Lote)
	if sub[collection] == nil {
		return noMatch()
	}

	fmt.Printf("%s Delete From Colection: %s\n", instance, collection)
	fmt.Println("_Id: ", objectId)
	return Delete(sub[collection], objectId)
}

func Delete(e entities.EntityDelete, id string) error {
	return e.Delete(id)
}

func noMatch() error {
	return errors.New("Sorry No collection matched!")
}
