package portin

import api "github.com/blackshark537/dataprod/src/app/Api"

type ApiAdapter struct{}

func (a *ApiAdapter) ForRoot(port string) error {
	//msg := color.RedString("Api Sever not available")
	//return errors.New(msg)
	return api.ForRoot(port)
}
