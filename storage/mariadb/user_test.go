package mariadb

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func TestUser(t *testing.T) {

	storage, err := New("ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard")
	if err != nil {
		log.Fatal("Error opening database: " + err.Error())
	}

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

	blueTeam := &model.Team{Name: "Blue Team"}
	err = storage.CreateTeam(blueTeam)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Team: %+v", blueTeam)

	teams, err = storage.GetTeams()
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

	user2 := &model.User{Username: "The Terminator", EmailAddress: "terminator@sky.net", Team: *blueTeam}
	err = storage.CreateUser(user2)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created User: %+v", user)

	users, err = storage.GetUsers()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Users: %+v", users)

	teams, err = storage.GetTeams()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Teams: %+v", teams)

	for _, team := range teams {
		storage.DeleteTeam(team)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Printf("Deleted All Teams")

	users, err = storage.GetUsers()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Users: %+v", teams)

	for _, user := range users {
		storage.DeleteUser(user)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Printf("Deleted All Users")

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
}
