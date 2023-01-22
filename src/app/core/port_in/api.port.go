package portin

import (
	api "github.com/blackshark537/dataprod/src/app/Api"
)

type ApiAdapter struct{}

var _api *api.API = nil

func (ap *ApiAdapter) ForRoot(port string) error {
	if _api == nil {
		_api = new(api.API)
	}
	return _api.ForRoot(port)
}
