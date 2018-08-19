package mariadb

import (
	"time"

	"github.com/thegrandpackard/ctf-scoreboard/model"
)

// CategoryStorage - categoryStorage
type CategoryStorage interface {
	CreateCategoryTable() (err error)
	DeleteCategoryTable() (err error)

	CreateCategory(category *model.Category) (err error)
	GetCategories() (categories []*model.Category, err error)
	GetCategory(category *model.Category) (err error)
	UpdateCategory(category *model.Category) (err error)
	DeleteCategory(category *model.Category) (err error)
}

// CreateCategoryTable - createCategoryTable
func (s *Storage) CreateCategoryTable() (err error) {

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS category (
		id INT(6) UNSIGNED AUTO_INCREMENT,
		name TEXT NOT NULL,
		display_order INT(6) NOT NULL DEFAULT 0,
		created TIMESTAMP,
		PRIMARY KEY (id)
		)`)

	return
}

// DeleteCategoryTable - deleteCategoryTable
func (s *Storage) DeleteCategoryTable() (err error) {

	_, err = s.db.Exec(`DROP TABLE IF EXISTS category`)

	return
}

// CreateCategory - createCategory
func (s *Storage) CreateCategory(category *model.Category) (err error) {

	result, err := s.db.Exec(`
		INSERT 
		INTO category
		SET name = ?,
			created = NOW(),
			display_order = IFNULL((SELECT MAX(display_order) FROM category AS cat2) + 1, 0)`,
		category.Name)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	category.ID = uint64(id)
	category.Created = uint64(time.Now().Unix())
	return
}

// GetCategories - getCategories
func (s *Storage) GetCategories() (categories []*model.Category, err error) {

	rows, err := s.db.Query(`
		SELECT id,
			   name,
			   UNIX_TIMESTAMP(created)
		FROM category
		ORDER BY display_order ASC, name ASC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		category := &model.Category{}
		err = rows.Scan(&category.ID, &category.Name, &category.Created)
		if err != nil {
			return
		}
		categories = append(categories, category)
	}

	return
}

// GetCategory - getCategory
func (s *Storage) GetCategory(category *model.Category) (err error) {

	err = s.db.QueryRow(`
		SELECT id,
			   name,
			   UNIX_TIMESTAMP(created)
		FROM category
		WHERE id = ?`, category.ID).Scan(&category.ID, &category.Name, &category.Created)
	return
}

// UpdateCategory - updateCategory
func (s *Storage) UpdateCategory(category *model.Category) (err error) {

	_, err = s.db.Exec(`
	UPDATE category
	SET name = ?
	WHERE id = ?`,
		category.Name,
		category.ID)

	return
}

// DeleteCategory - deleteCategory
func (s *Storage) DeleteCategory(category *model.Category) (err error) {

	_, err = s.db.Exec(`
	DELETE 
	FROM category
	WHERE id = ?`,
		category.ID)

	return
}
