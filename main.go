package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/thegrandpackard/ctf-scoreboard/rest"
	"github.com/thegrandpackard/ctf-scoreboard/storage/mariadb"
)

func main() {
	log.Printf("CTF Scoreboard API v0.1")

	// Initialize storage
	// TODO Get config from text file
	storage, err := mariadb.New("ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard")
	if err != nil {
		log.Fatal("Error opening database: " + err.Error())
	}

	// REST Interface
	err = rest.InitializeRoutes(storage)
	if err != nil {
		log.Fatal("Error initializeing routes: " + err.Error())
	}

	// TODO: Websocket Interface

	// TODO: Static Webserver Interface for Angular frontend
}
