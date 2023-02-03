package config

import (
	"os"
	"strings"
)

var VERSION string = "0.0.7"
var DatabaseName string = os.Getenv("DB_NAME")
var DatabaseUri string = strings.Join(strings.Split(os.Getenv("DB_URL"), `\`), "")

var IsBench bool = true
