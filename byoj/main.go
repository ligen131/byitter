package main

import (
	"byoj/model"
	"byoj/shared/database"
	"byoj/shared/server"
	"byoj/shared/yamlconfig"
)

func main() {
	configuration, err := yamlconfig.ConfigLoad("config.yml")
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(configuration.Database)
	if err != nil {
		panic(err)
	}

	err = model.InitModel(db)
	if err != nil {
		panic(err)
	}

	err = server.Run(configuration.Server)
	if err != nil {
		panic(err)
	}
}
