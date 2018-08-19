package mariadb

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func TestCategory(t *testing.T) {

	storage, err := New("ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard")
	if err != nil {
		log.Fatal("Error opening database: " + err.Error())
	}

	err = storage.CreateCategoryTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Category Table")

	smtp := &model.Category{Name: "SMTP"}
	err = storage.CreateCategory(smtp)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Category: %+v", smtp)

	categories, err := storage.GetCategories()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories: %+v", categories)

	snmp := &model.Category{Name: "SNMP"}
	err = storage.CreateCategory(snmp)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Category: %+v", snmp)

	categories, err = storage.GetCategories()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories: %+v", categories)

	snmp.Name = "SIP"
	err = storage.UpdateCategory(snmp)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Updated Category: %+v", snmp)

	categories, err = storage.GetCategories()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories: %+v", categories)

	for _, category := range categories {
		storage.DeleteCategory(category)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Printf("Deleted All Categories")

	categories, err = storage.GetCategories()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Categories %+v", categories)

	err = storage.DeleteCategoryTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted Categories Table")
}
