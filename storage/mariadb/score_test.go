package mariadb

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func TestScore(t *testing.T) {

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

	err = storage.CreateUserTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created User Table")

	err = storage.CreateScoreTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Score Table")

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

	user := &model.User{Username: "1337 H4x0r", EmailAddress: "h4x0r@1337.com", Team: *redTeam}
	err = storage.CreateUser(user)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created User: %+v", user)

	users, err := storage.GetUsers()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Users: %+v", users)

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
		Text:     "SMTP 1",
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

	score := &model.Score{}
	score.User.ID = user.ID
	score.Question.ID = question1.ID
	err = storage.CreateScore(score)
	if err != nil {
		log.Fatal(err.Error())
	}

	scores, err := storage.GetScores()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Scores: %+v", scores)

	// Cleanup Tables

	err = storage.DeleteScoreTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted Score Table")

	err = storage.DeleteUserTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted User Table")

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
