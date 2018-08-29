package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/thegrandpackard/ctf-scoreboard/config"
	"github.com/thegrandpackard/ctf-scoreboard/rest"
	"github.com/thegrandpackard/ctf-scoreboard/storage/mariadb"
)

var err error
var storage *mariadb.Storage

func main() {
	log.Printf("CTF Scoreboard API v0.1")

	configuration := config.LoadConfig()

	// Initialize Storage
	// FEATURE: Add support for additional storage types
	if configuration.StorageType == "mysql" {
		storage, err = mariadb.New(configuration.StorageConfig)
		if err != nil {
			log.Fatal("Error opening database: " + err.Error())
		}
	} else {
		log.Fatalf("Unsupported storage type: %s", configuration.StorageType)
	}

	// Initialize REST Interface
	err = rest.InitializeRoutes(storage)
	if err != nil {
		log.Fatal("Error initializeing routes: " + err.Error())
	}

	// TODO: Initialize Websocket Interface

	// TODO: Initialize Static Webserver Interface for Angular frontend
}
