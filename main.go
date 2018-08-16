package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thegrandpackard/ctf-scoreboard/storage/mariadb"
)

func main() {
	log.Printf("CTF Scoreboard API v0.1")

	db, err := sql.Open("mysql", "ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard")
	if err != nil {
		log.Fatal("Error opening database: " + err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = mariadb.CreateCategoryTable(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Categories Table")

	category, err := mariadb.CreateCategory(db, "SMTP")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created category: %+v", category)

	categories, err := mariadb.GetCategories(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories %+v", categories)

	category, err = mariadb.CreateCategory(db, "SNMP")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created category: %+v", category)

	categories, err = mariadb.GetCategories(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories %+v", categories)

	category.Name = "SIP"
	err = mariadb.UpdateCategory(db, category)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Updated category: %+v", category)

	categories, err = mariadb.GetCategories(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories %+v", categories)

	for _, category := range categories {
		mariadb.DeleteCategory(db, category.ID)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Printf("Deleted All Categories")

	categories, err = mariadb.GetCategories(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories %+v", categories)

	err = mariadb.DeleteCategoryTable(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted Categories Table")

	log.Printf("CTF Scoreboard API Exited")
}
