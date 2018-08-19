package mariadb

import (
	"github.com/thegrandpackard/ctf-scoreboard/model"
)

// QuestionStorage - questionStorage
type QuestionStorage interface {
	CreateQuestionTable() (err error)
	DeleteQuestionTable() (err error)

	CreateQuestion(question *model.Question) (err error)
	GetQuestions() (questions []*model.Question, err error)
	GetQuestion(question *model.Question) (err error)
	UpdateQuestion(question *model.Question) (err error)
	DeleteQuestion(question *model.Question) (err error)
}

// CreateQuestionTable - createQuestionTable
func (s *Storage) CreateQuestionTable() (err error) {

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS question (
		id INT(6) UNSIGNED AUTO_INCREMENT,
		category_id INT(6) UNSIGNED NOT NULL,
		team_id INT(6) UNSIGNED NOT NULL DEFAULT 0,
		text TEXT NOT NULL,
		answer TEXT NOT NULL,
		hint TEXT NOT NULL,
		file TEXT NOT NULL,
		points INT(6) UNSIGNED NOT NULL DEFAULT 0,
		created TIMESTAMP,
		PRIMARY KEY (id)
		)`)

	return
}

// DeleteQuestionTable - deleteQuestionTable
func (s *Storage) DeleteQuestionTable() (err error) {

	_, err = s.db.Exec(`DROP TABLE IF EXISTS question`)

	return
}

// CreateQuestion - createQuestion
func (s *Storage) CreateQuestion(question *model.Question) (err error) {

	result, err := s.db.Exec(`
		INSERT 
		INTO question
		SET category_id = ?,
			team_id = ?,
			text = ?,
			answer = ?,
			hint = ?,
			file = ?,
			points = ?,
			created = NOW()`,
		question.Category.ID,
		question.Team.ID,
		question.Text,
		question.Answer,
		question.Hint,
		question.File,
		question.Points)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	question.ID = uint64(id)
	return
}

// GetQuestions - getQuestions
func (s *Storage) GetQuestions() (questions []*model.Question, err error) {

	rows, err := s.db.Query(`
		SELECT id,
			   category_id,
			   team_id,
			   text,
			   answer,
			   hint,
			   file,
			   points,
			   UNIX_TIMESTAMP(created)
		FROM question
		ORDER BY id ASC`)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		question := &model.Question{}
		err = rows.Scan(
			&question.ID,
			&question.Category.ID,
			&question.Team.ID,
			&question.Text,
			&question.Answer,
			&question.Hint,
			&question.File,
			&question.Points,
			&question.Created)
		if err != nil {
			return
		}

		err = s.GetCategory(&question.Category)
		if err != nil {
			return
		}

		if question.Team.ID > 0 {
			err = s.GetTeam(&question.Team)
			if err != nil {
				return
			}
		}

		questions = append(questions, question)
	}

	return
}

// GetQuestion - getQuestion
func (s *Storage) GetQuestion(question *model.Question) (err error) {

	err = s.db.QueryRow(`
		SELECT id,
			   category_id,
			   team_id,
			   text,
			   answer,
			   hint,
			   file,
			   points,
			   UNIX_TIMESTAMP(created)
		FROM question
		ORDER BY name ASC`).Scan(
		&question.ID,
		&question.Category.ID,
		&question.Team.ID,
		&question.Text,
		&question.Answer,
		&question.Hint,
		&question.File,
		&question.Points,
		&question.Created)
	if err != nil {
		return
	}

	err = s.GetCategory(&question.Category)
	if err != nil {
		return
	}

	if question.Team.ID > 0 {
		err = s.GetTeam(&question.Team)
		if err != nil {
			return
		}
	}

	return
}

// UpdateQuestion - updateQuestion
func (s *Storage) UpdateQuestion(Question *model.Question) (err error) {

	_, err = s.db.Exec(`
	UPDATE question
	SET category_id = ?,
		team_id = ?,
		text = ?,
		answer = ?,
		hint = ?,
		file = ?,
		points = ?
	WHERE id = ?`,
		Question.Category.ID,
		Question.Team.ID,
		Question.Text,
		Question.Answer,
		Question.Hint,
		Question.File,
		Question.Points)

	return
}

// DeleteQuestion - deleteQuestion
func (s *Storage) DeleteQuestion(Question *model.Question) (err error) {

	_, err = s.db.Exec(`
	DELETE 
	FROM question
	WHERE id = ?`,
		Question.ID)

	return
}
