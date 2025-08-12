package main

import (
	"log"
	"os"

	"www.github.com/kharljhon14/kwento-ko/cmd/api"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Fatalln("failed to create new server:", err)
	}

	err = server.Start(os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		log.Fatalln("failed to start server:", err)
	}
}
