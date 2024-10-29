package main

import (
	"fmt"
	"io"
	"log"
	"ntsiris/product-microservice/api"
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logFile := setUpFileLog()
	defer resourceCleanUp(logFile)

	var store = storage.MySQLStore{}
	err := store.InitStore(&config.EnvDBConfig)
	if err != nil {
		log.Fatalf("Store Initialization %v", err)
	}
	defer resourceCleanUp(logFile)

	if err = store.VerifyStoreConnection(); err != nil {
		log.Fatalf("Verify Storage Connection %v", err)
	}
	log.Printf("Successfully established connection to storage component")

	if config.EnvAPIServerConfig.MigrateUp {
		log.Printf("Running Up Migrations from %s", config.EnvAPIServerConfig.MigrationPath)
		if err := store.RunMigrationUp(config.EnvAPIServerConfig.MigrationPath); err != nil {
			log.Fatalf("Run Up Migration %v", err)
		}
		log.Print("Up Migrations finished successfully!")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go serverCleanUp(&store, sigs)

	apiServerAddress := fmt.Sprintf("%s:%s", config.EnvAPIServerConfig.PublicHost, config.EnvAPIServerConfig.Port)
	apiServer := api.NewAPIServer(apiServerAddress, &store)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("Start API server %v\n", err)
	}
}

func setUpFileLog() *os.File {
	logFile, err := os.OpenFile(config.EnvAPIServerConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Open log file \"%s\": %v", config.EnvAPIServerConfig.LogFile, err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stderr)

	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile
}

func serverCleanUp(store storage.ProductStore, sigs chan os.Signal) {
	<-sigs

	log.Print("Shuting down server...")

	if config.EnvAPIServerConfig.MigrateDown {
		log.Printf("Running Down Migrations from %s", config.EnvAPIServerConfig.MigrationPath)
		if err := store.RunMigrationDown(config.EnvAPIServerConfig.MigrationPath); err != nil {
			log.Fatalf("Run Down Migration %v", err)
		}
		log.Print("Down Migrations finished successfully!")
	}

	log.Print("Server Stopped!")
	os.Exit(0)
}

func resourceCleanUp(resource io.Closer) {
	if err := resource.Close(); err != nil {
		log.Fatal(err)
	}
}
