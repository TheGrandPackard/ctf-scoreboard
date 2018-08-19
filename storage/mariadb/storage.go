package mariadb

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// TODO Foreign Key Constraints

// Storage - storage
type Storage struct {
	db *sqlx.DB
}

// New - new
func New(config string) (db *Storage, err error) {

	db = &Storage{}

	if config == "" {
		err = fmt.Errorf("Invalid configuration passed (empty)")
		return
	}

	if db.db, err = sqlx.Open("mysql", config); err != nil {
		return
	}

	err = db.db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
