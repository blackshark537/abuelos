package config

import "os"

var VERSION string = "0.0.7"
var DatabaseName string = os.Getenv("DB_NAME")
var DatabaseUri string = "mongodb://localhost:27017"

var IsBench bool = true
