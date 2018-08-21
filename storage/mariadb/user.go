package mariadb

import (
	"time"

	"github.com/thegrandpackard/ctf-scoreboard/model"
)

// UserStorage - userStorage
type UserStorage interface {
	CreateUserTable() (err error)
	DeleteUserTable() (err error)

	CreateUser(user *model.User) (err error)
	GetUsers() (users []*model.User, err error)
	GetUser(user *model.User) (err error)
	UpdateUser(user *model.User) (err error)
	UpdateUserPassword(user *model.User) (err error)
	DeleteUser(user *model.User) (err error)
}

// CreateUserTable - createUserTable
func (s *Storage) CreateUserTable() (err error) {

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
		id INT(6) UNSIGNED AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email_address TEXT NOT NULL,
		admin TINYINT(1) NOT NULL DEFAULT 0,
		team_id INT(6) UNSIGNED NOT NULL,
		created TIMESTAMP,
		PRIMARY KEY (id)
		)`)

	return
}

// DeleteUserTable - deleteUserTable
func (s *Storage) DeleteUserTable() (err error) {

	_, err = s.db.Exec(`DROP TABLE IF EXISTS user`)

	return
}

// CreateUser - createUser
func (s *Storage) CreateUser(user *model.User) (err error) {

	result, err := s.db.Exec(`
		INSERT 
		INTO user
		SET username = ?,
			password = ?,
			email_address = ?,
			team_id = ?,
			created = NOW()`,
		user.Username,
		user.Password,
		user.EmailAddress,
		user.Team.ID)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	user.ID = uint64(id)
	user.Created = uint64(time.Now().Unix())
	user.Password = ""
	return
}

// GetUsers - getUsers
func (s *Storage) GetUsers() (users []*model.User, err error) {

	rows, err := s.db.Query(`
		SELECT id,
			   admin,
			   username,
			   email_address,
			   team_id,
			   UNIX_TIMESTAMP(created)
		FROM user
		ORDER BY username ASC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.ID, &user.Admin, &user.Username, &user.EmailAddress, &user.Team.ID, &user.Created)
		if err != nil {
			return
		}
		if user.Team.ID > 0 {
			err = s.GetTeam(&user.Team)
			if err != nil {
				return nil, err
			}
		}
		users = append(users, user)
	}

	return
}

// GetUser - getUser
func (s *Storage) GetUser(user *model.User) (err error) {

	err = s.db.QueryRow(`
		SELECT id,
			   admin,
			   username,
			   email_address,
			   team_id,
			   UNIX_TIMESTAMP(created)
		FROM user
		WHERE id = ?`, user.ID).Scan(&user.ID, &user.Admin, &user.Username, &user.EmailAddress, &user.Team.ID, &user.Created)
	if err != nil {
		return
	}

	if user.Team.ID > 0 {
		err = s.GetTeam(&user.Team)
		if err != nil {
			return err
		}
	}

	return
}

// UpdateUser - updateUser
func (s *Storage) UpdateUser(user *model.User) (err error) {

	_, err = s.db.Exec(`
	UPDATE user
	SET username = ?,
	email_address = ?,
	team_id = ?
	WHERE id = ?`,
		user.Username,
		user.EmailAddress,
		user.Team.ID)

	return
}

// UpdateUserPassword - updateUserPassword
func (s *Storage) UpdateUserPassword(user *model.User) (err error) {

	_, err = s.db.Exec(`
	UPDATE user
	SET password = ?
	WHERE id = ?`,
		user.Password,
		user.Team.ID)

	return
}

// DeleteUser - deleteUser
func (s *Storage) DeleteUser(user *model.User) (err error) {

	_, err = s.db.Exec(`
	DELETE 
	FROM user
	WHERE id = ?`,
		user.ID)

	return
}
