package mariadb

import (
	"time"

	"github.com/thegrandpackard/ctf-scoreboard/model"
)

// TeamStorage - teamStorage
type TeamStorage interface {
	CreateTeamTable() (err error)
	DeleteTeamTable() (err error)

	CreateTeam(team *model.Team) (err error)
	GetTeams() (teams []*model.Team, err error)
	GetTeam(team *model.Team) (err error)
	UpdateTeam(team *model.Team) (err error)
	DeleteTeam(team *model.Team) (err error)
}

// CreateTeamTable - createTeamTable
func (s *Storage) CreateTeamTable() (err error) {

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS team (
		id INT(6) UNSIGNED AUTO_INCREMENT,
		name TEXT NOT NULL,
		created TIMESTAMP,
		PRIMARY KEY (id)
		)`)

	return
}

// DeleteTeamTable - deleteTeamTable
func (s *Storage) DeleteTeamTable() (err error) {

	_, err = s.db.Exec(`DROP TABLE IF EXISTS team`)

	return
}

// CreateTeam - createTeam
func (s *Storage) CreateTeam(team *model.Team) (err error) {

	result, err := s.db.Exec(`
		INSERT 
		INTO team
		SET name = ?,
			created = NOW()`,
		team.Name)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	team.ID = uint64(id)
	team.Created = uint64(time.Now().Unix())
	return
}

// GetTeams - getTeams
func (s *Storage) GetTeams() (teams []*model.Team, err error) {

	rows, err := s.db.Query(`
		SELECT id,
			   name,
			   UNIX_TIMESTAMP(created)
		FROM team
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		Team := &model.Team{}
		err = rows.Scan(&Team.ID, &Team.Name, &Team.Created)
		if err != nil {
			return
		}
		teams = append(teams, Team)
	}

	return
}

// GetTeam - getTeam
func (s *Storage) GetTeam(team *model.Team) (err error) {

	err = s.db.QueryRow(`
		SELECT id,
			   name,
			   UNIX_TIMESTAMP(created)
		FROM team
		WHERE id = ?`, team.ID).Scan(&team.ID, &team.Name, &team.Created)

	return
}

// UpdateTeam - updateTeam
func (s *Storage) UpdateTeam(Team *model.Team) (err error) {

	_, err = s.db.Exec(`
	UPDATE team
	SET name = ?
	WHERE id = ?`,
		Team.Name,
		Team.ID)

	return
}

// DeleteTeam - deleteTeam
func (s *Storage) DeleteTeam(Team *model.Team) (err error) {

	_, err = s.db.Exec(`
	DELETE 
	FROM team
	WHERE id = ?`,
		Team.ID)

	return
}
