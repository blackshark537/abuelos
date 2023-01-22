package main

import (
	"log"
	"os"

	CommandLine "github.com/blackshark537/dataprod/src/app/cli"
)

func main() {
	app := CommandLine.New()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
