package mariadb

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

func TestTeam(t *testing.T) {

	storage, err := New("ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard")
	if err != nil {
		log.Fatal("Error opening database: " + err.Error())
	}

	err = storage.CreateTeamTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Team Table")

	blueTeam := &model.Team{Name: "Blue Team"}
	err = storage.CreateTeam(blueTeam)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Team: %+v", blueTeam)

	teams, err := storage.GetTeams()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Teams: %+v", teams)

	redTeam := &model.Team{Name: "Red Team"}
	err = storage.CreateTeam(redTeam)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created Team: %+v", redTeam)

	teams, err = storage.GetTeams()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Teams: %+v", teams)

	redTeam.Name = "R3d Team"
	err = storage.UpdateTeam(redTeam)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Updated Team: %+v", redTeam)

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

	teams, err = storage.GetTeams()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Teams: %+v", teams)

	err = storage.DeleteTeamTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Deleted Team Table")
}
