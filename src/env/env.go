package env

import (
	"fmt"
	"os"
	"strings"
)

var databaseName string = os.Getenv("DB_NAME")
var databaseUri string = strings.Join(strings.Split(os.Getenv("DB_URL"), `\`), "")

// Return DB_NAME, DB_URI
func GetEnv() (string, string) {
	if databaseName == "" {
		fmt.Printf("Invalid Enviroment DB_NAME %v\n, please set this variable", databaseName)
		os.Exit(1)
	}
	if databaseUri == "" {
		fmt.Printf("Invalid Enviroment DB_URL %v, please set this variable\n", databaseUri)
		os.Exit(1)
	}
	return databaseName, databaseUri
}
