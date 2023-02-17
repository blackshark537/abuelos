package main

import (
	"log"
	"os"

	api "github.com/blackshark537/abuelos/src/app/Api"
	CommandLine "github.com/blackshark537/abuelos/src/app/cli"
	core "github.com/blackshark537/abuelos/src/app/core"
	"github.com/blackshark537/abuelos/src/app/infraestructure/database"
	"github.com/blackshark537/abuelos/src/env"
)

func main() {
	core := new(core.Core)
	dbName, dbUri := env.GetEnv()

	core.InjectDatabase(
		&database.MongoDb{
			Uri:  dbUri,
			Name: dbName,
		},
	)
	core.InjectApi(
		new(api.API),
	)
	app := CommandLine.New()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
