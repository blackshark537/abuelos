package portin

import (
	"errors"

	"github.com/fatih/color"
)

type ApiAdapter struct{}

func (a *ApiAdapter) ForRoot(port string) error {
	msg := color.RedString("Api Sever not available")
	return errors.New(msg) //api.ForRoot(port)
}
