package main

import (
	"byoj/shared/server"
	"byoj/shared/yamlconfig"
)

func main() {
	configuration, err := yamlconfig.ConfigLoad("config.yml")
	if err != nil {
		panic(err)
	}

	server.Run(configuration.Server)
}
