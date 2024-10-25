package main

import (
	"fmt"
	"log"
	"ntsiris/product-microservice/api"
	"ntsiris/product-microservice/config"
	"ntsiris/product-microservice/internal/database"
)

func main() {
	db, err := database.NewSQLDatabase(config.EnvDBConfig)
	if err != nil {
		log.Fatalf("Failed to %v", err)
	}

	if err = database.VerifyDatabaseConnection(db); err != nil {
		log.Fatalf("Failed to %v", err)
	}
	log.Printf("Successfully established connection to database")

	apiServerAddress := fmt.Sprintf("%s:%s", config.EnvAPIServerConfig.PublicHost, config.EnvAPIServerConfig.Port)
	apiServer := api.NewAPIServer(apiServerAddress, db)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("Failed to start API server: %v\n", err)
	}
	log.Printf("Product API Server running on address: %s\n", apiServer.Address)
}
