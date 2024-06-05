package main

import (
	"flag"
	"strings"

	"github.com/downsized-devs/sdk-go/generator/services"
)

func main() {

	entityName := flag.String("entity_name", "example", "entity name as a camelCase string")
	fileLocation := flag.String("file_location", "generate", "aquahero-service file location as a string")
	api := flag.String("api", "create, edit, get, activate, delete", "choose generated api method \"create, edit, get, activate, delete\"")
	flag.Parse()
	apis := strings.Split(*api, ",")
	service := services.Init(services.Param{EntityName: *entityName, Location: *fileLocation, Api: apis})

	err := service.Execute()
	if err != nil {
		panic(err)
	}
}
