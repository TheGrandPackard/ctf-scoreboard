package mariadb

import (
	"time"

	"github.com/thegrandpackard/ctf-scoreboard/model"
)

// ScoreStorage - scoreStorage
type ScoreStorage interface {
	CreateScoreTable() (err error)
	DeleteScoreTable() (err error)

	CreateScore(score *model.Score) (err error)
	GetScores() (scores []*model.Score, err error)
	GetUserScores(user model.User) (scores []*model.Score, err error)
	GetQuestionScores(question model.Question) (scores []*model.Score, err error)
	GetScore(user model.User, question model.Question) (score *model.Score, err error)
	DeleteScore(score *model.Score) (err error)
}

// CreateScoreTable - createScoreTable
func (s *Storage) CreateScoreTable() (err error) {

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS score (
		user_id INT(6) UNSIGNED NOT NULL,
		question_id INT(6) UNSIGNED NOT NULL,
		created TIMESTAMP,
		PRIMARY KEY (user_id, question_id)
		)`)

	return
}

// DeleteScoreTable - deleteScoreTable
func (s *Storage) DeleteScoreTable() (err error) {

	_, err = s.db.Exec(`DROP TABLE IF EXISTS score`)

	return
}

// CreateScore - createScore
func (s *Storage) CreateScore(score *model.Score) (err error) {

	_, err = s.db.Exec(`
		INSERT 
		INTO score
		SET user_id = ?,
			question_id = ?,
			created = NOW()`,
		score.User.ID,
		score.Question.ID)
	if err != nil {
		return
	}

	score.Created = uint64(time.Now().Unix())
	return
}

// GetScores - getScores
func (s *Storage) GetScores() (scores []*model.Score, err error) {

	rows, err := s.db.Query(`
		SELECT user_id,
			   question_id,
			   UNIX_TIMESTAMP(created)
		FROM score
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		score := &model.Score{}
		err = rows.Scan(&score.User.ID, &score.Question.ID, &score.Created)
		if err != nil {
			return
		}
		scores = append(scores, score)
	}

	return
}

// GetUserScores - getUserScores
func (s *Storage) GetUserScores(user model.User) (scores []*model.Score, err error) {

	rows, err := s.db.Query(`
		SELECT user_id,
			   question_id,
			   UNIX_TIMESTAMP(created)
		FROM score
		WHERE user_id = ?
		ORDER BY name ASC`, user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		score := &model.Score{}
		err = rows.Scan(&score.User.ID, &score.Question.ID, &score.Created)
		if err != nil {
			return
		}
		scores = append(scores, score)
	}

	return
}

// GetQuestionScores - getQuestionScores
func (s *Storage) GetQuestionScores(question model.Question) (scores []*model.Score, err error) {

	rows, err := s.db.Query(`
		SELECT user_id,
			   question_id,
			   UNIX_TIMESTAMP(created)
		FROM score
		WHERE question_id = ?
		ORDER BY name ASC`, question.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		score := &model.Score{}
		err = rows.Scan(&score.User.ID, &score.Question.ID, &score.Created)
		if err != nil {
			return
		}
		scores = append(scores, score)
	}

	return
}

// GetScore - getScore
func (s *Storage) GetScore(user model.User, question model.Question) (score *model.Score, err error) {

	err = s.db.QueryRow(`
		SELECT UNIX_TIMESTAMP(created)
		FROM score
		WHERE user_id = ? AND 
			  question_id = ?`, user.ID, question.ID).Scan(&score.Created)

	return
}

// DeleteScore - deleteScore
func (s *Storage) DeleteScore(score *model.Score) (err error) {

	_, err = s.db.Exec(`
	DELETE 
	FROM score
	WHERE user_id = ? AND
		  question_id = ?`,
		score.User.ID,
		score.Question.ID)

	return
}
