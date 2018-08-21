package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/thegrandpackard/ctf-scoreboard/rest"
	"github.com/thegrandpackard/ctf-scoreboard/storage/mariadb"
)

func main() {
	log.Printf("CTF Scoreboard API v0.1")

	// TODO: Get config from text file
	storageConfig := "ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard"

	// Initialize Storage
	storage, err := mariadb.New(storageConfig)
	if err != nil {
		log.Fatal("Error opening database: " + err.Error())
	}

	// Initialize REST Interface
	err = rest.InitializeRoutes(storage)
	if err != nil {
		log.Fatal("Error initializeing routes: " + err.Error())
	}

	// TODO: Initialize Websocket Interface

	// TODO: Initialize Static Webserver Interface for Angular frontend
}
