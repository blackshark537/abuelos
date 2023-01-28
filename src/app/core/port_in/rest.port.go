package portin

import (
	api "github.com/blackshark537/dataprod/src/app/Api"
)

type ApiPort struct{}

var rest *api.API = nil

func InjectApi(a *api.API) {
	if rest == nil {
		rest = a
	}
}

func (ap *ApiPort) ForRoot(port string) error {
	return rest.ForRoot(port)
}
