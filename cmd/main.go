package main

import (
	"log"
	"ntsiris/product-microservice/api"
)

func main() {
	apiServer := api.NewAPIServer(":8080", nil)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("starting API server: %v\n", err)
	}
}
