package core

import (
	portout "github.com/blackshark537/dataprod/src/app/core/port_out"
)

func ForRoot() {
	dbPort := portout.DbAdapter{}
	dbPort.ForRoot()
}
