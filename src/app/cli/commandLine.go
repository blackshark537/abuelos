package CommandLine

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	core_command "github.com/blackshark537/dataprod/src/app/core/command"
	"github.com/blackshark537/dataprod/src/app/core/config"
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
var cmds *core_command.CliCommand = new(core_command.CliCommand)
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
	{
		Name:    "SetVars",
		Aliases: []string{"sv"},
		Usage:   "Set and initialize all variables",
		Action:  SetVariables,
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

	AppName := color.MagentaString(fmt.Sprintf("Dataprod Core Arch: %v", runtime.GOARCH))
	AppDesc := color.MagentaString("A simple CLI program to manage your dataprod system")
	return &cli.App{
		Name:     AppName,
		Version:  "v" + config.VERSION,
		Usage:    AppDesc,
		Commands: Commands,
	}
}

func SetVariables(ctx *cli.Context) error {
	cmds.SetVariables()
	return nil
}

func serverStart(ctx *cli.Context) error {
	return cmds.ServerStart(port)
}

func projectLote(ctx *cli.Context) error {
	cmds.Collection = "lotes"
	cmds.ListAbuelos(lote)
	return nil
}

func projectTable(ctx *cli.Context) error {
	cmds.Collection = collection
	return cmds.ProjectTable(year, data, prodToBool())
}

func prodToBool() bool {
	prod := false
	if isProduccion == "yes" || isProduccion == "y" {
		prod = true
	}
	return prod
}

func listTable(ctx *cli.Context) error {
	cmds.Collection = collection
	return cmds.ListTable(filter)
}

func insertIntoTable(ctx *cli.Context) error {
	cmds.Collection = collection
	return cmds.InsertIntoCollection(data)
}

func InsertMany(ctx *cli.Context) error {
	cmds.Collection = collection
	return cmds.InsertFromFile(path)
}

func DeleteAll(ctx *cli.Context) error {
	cmds.Collection = collection
	return cmds.ClearCollection(filter)
}

func deleteFromTable(ctx *cli.Context) error {
	cmds.Collection = collection
	return cmds.RemoveFromCollection(objectId)
}

func noMatch() error {
	return errors.New(fmt.Sprintf("%s Sorry No collection matched!", instance))
}

func handleErr(e error) error {
	if e != nil {
		return e
	}
	return nil
}
