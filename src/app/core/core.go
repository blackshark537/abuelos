package core

import (
	"log"
	"os"

	CommandLine "github.com/blackshark537/dataprod/src/app/core/cli"
	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
)

func ForRoot() {
	dbPort := portout.DbAdapter{}
	dbPort.ForRoot()
	app := CommandLine.New()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
