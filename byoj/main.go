package main

import "byoj/shared/server"

func main() {
	server.Run(server.Server{
		Port: 3435,
	})
}
