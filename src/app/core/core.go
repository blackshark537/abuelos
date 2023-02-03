package core

import (
	portin "github.com/blackshark537/dataprod/src/app/core/port_in"
	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
)

type Core struct{}

func (c Core) InjectApi(api portin.ApiGateway) {
	portin.InjectApi(api)
}

func (c Core) InjectDatabase(database portout.Database) {
	portout.InjectDatabase(database)
}
