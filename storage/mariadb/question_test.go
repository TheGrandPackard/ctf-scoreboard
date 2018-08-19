package mariadb

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func TestQuestion(t *testing.T) {

	storage, err := New("ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard")
	if err != nil {
		log.Fatal("Error opening database: " + err.Error())
	}

	// Create Tables

	err = storage.CreateCategoryTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Category Table")

	err = storage.CreateQuestionTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Question Table")

	err = storage.CreateTeamTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Team Table")

	// Testing

	redTeam := &model.Team{Name: "Red Team"}
	err = storage.CreateTeam(redTeam)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Team: %+v", redTeam)

	teams, err := storage.GetTeams()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Teams: %+v", teams)

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

	question1 := &model.Question{
		Category: *smtp,
		Team:     *redTeam,
		Text:     "Question 1",
		Answer:   "flag{question1}",
		Hint:     "No hint",
		File:     "http://www.google.com",
		Points:   25,
	}
	err = storage.CreateQuestion(question1)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Question: %+v", question1)

	questions, err := storage.GetQuestions()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Questions: %+v", questions)

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

	// Cleanup Tables

	err = storage.DeleteTeamTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted Team Table")

	err = storage.DeleteQuestionTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted Question Table")

	err = storage.DeleteCategoryTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted Categories Table")
}
