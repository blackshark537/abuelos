package core

import (
	"fmt"
	"os"

	api "github.com/blackshark537/dataprod/src/app/Api"
	"github.com/blackshark537/dataprod/src/app/core/config"
	portin "github.com/blackshark537/dataprod/src/app/core/port_in"
	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
	"github.com/blackshark537/dataprod/src/app/infraestructure/database"
)

func ForRoot() {
	if config.DatabaseName == "" {
		fmt.Printf("Invalid Enviroment DB_NAME %v\n", config.DatabaseName)
		os.Exit(1)
	}
	portout.InjectDatabase(
		&database.MongoDb{
			Uri:  config.DatabaseUri,
			Name: config.DatabaseName,
		},
	)

	portin.InjectApi(
		new(api.API),
	)
}
