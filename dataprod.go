package main

import (
	"log"
	"os"

	CommandLine "github.com/blackshark537/dataprod/src/app/cli"
)

func main() {
	app := CommandLine.New()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
