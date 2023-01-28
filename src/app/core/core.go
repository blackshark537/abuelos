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
	config.DatabaseUri = "mongodb://localhost:27017"
	dbName := config.DatabaseName
	if dbName == "" {
		fmt.Printf("Invalid Enviroment DB_NAME %v\n", dbName)
		os.Exit(1)
	}
	portin.InjectApi(new(api.API))
	portout.InjectDatabase(&database.MongoDb{
		Uri:  config.DatabaseUri,
		Name: config.DatabaseName,
	})
}
