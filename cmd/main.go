package main

import (
	"fmt"
	"log"
	"ntsiris/product-microservice/api"
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var store = storage.MySQLStore{}

	err := store.InitStore(&config.EnvDBConfig)
	if err != nil {
		log.Fatalf("Failed to %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Fatalf("Abnormal Store Close: %v", err)
		}
	}()

	if err = store.VerifyStoreConnection(); err != nil {
		log.Fatalf("Failed to %v", err)
	}
	log.Printf("Successfully established connection to storage component")

	if config.EnvAPIServerConfig.MigrateUp {
		log.Printf("Running Up Migrations from %s", config.EnvAPIServerConfig.MigrationPath)
		if err := store.RunMigrationUp(config.EnvAPIServerConfig.MigrationPath); err != nil {
			log.Fatalf("Failed to %v", err)
		}
		log.Print("Up Migrations finished successfully!")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		log.Print("Shuting down server...")

		if config.EnvAPIServerConfig.MigrateDown {
			log.Printf("Running Down Migrations from %s", config.EnvAPIServerConfig.MigrationPath)
			if err := store.RunMigrationDown(config.EnvAPIServerConfig.MigrationPath); err != nil {
				log.Fatalf("Failed to %v", err)
			}
			log.Print("Down Migrations finished successfully!")
		}

		log.Print("Server Stopped!")
		os.Exit(0)
	}()

	apiServerAddress := fmt.Sprintf("%s:%s", config.EnvAPIServerConfig.PublicHost, config.EnvAPIServerConfig.Port)
	apiServer := api.NewAPIServer(apiServerAddress, &store)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("Failed to start API server: %v\n", err)
	}
}
