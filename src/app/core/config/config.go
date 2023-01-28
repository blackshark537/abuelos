package config

import "os"

var VERSION string = "0.0.6"
var DatabaseName string = os.Getenv("DB_NAME")
var DatabaseUri string = ""
var IsBench bool = true
