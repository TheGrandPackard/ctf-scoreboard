package mariadb

import (
	"database/sql"
	"time"

	"github.com/thegrandpackard/ctf-scoreboard/model"
)

// CreateCategoryTable - createCategoryTable
func CreateCategoryTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
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
func DeleteCategoryTable(db *sql.DB) (err error) {

	_, err = db.Exec(`DROP TABLE IF EXISTS category`)

	return
}

// CreateCategory - createCategory
func CreateCategory(db *sql.DB, name string) (category model.Category, err error) {

	result, err := db.Exec(`
		INSERT 
		INTO category
		SET name = ?,
			created = NOW(),
			display_order = IFNULL((SELECT MAX(display_order) FROM category AS cat2) + 1, 0)`,
		name)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	category = model.Category{ID: uint64(id), Name: name, Created: uint64(time.Now().Unix())}

	return
}

// GetCategories - getCategories
func GetCategories(db *sql.DB) (categories []model.Category, err error) {

	rows, err := db.Query(`
		SELECT id,
			   name,
			   UNIX_TIMESTAMP(created)
		FROM category
		ORDER BY display_order ASC, name ASC`)

	defer rows.Close()

	for rows.Next() {
		category := model.Category{}
		err = rows.Scan(&category.ID, &category.Name, &category.Created)
		if err != nil {
			return
		}
		categories = append(categories, category)
	}

	return
}

// UpdateCategory - updateCategory
func UpdateCategory(db *sql.DB, category model.Category) (err error) {

	_, err = db.Exec(`
	UPDATE category
	SET name = ?
	WHERE id = ?`,
		category.Name,
		category.ID)

	return
}

// DeleteCategory - deleteCategory
func DeleteCategory(db *sql.DB, id uint64) (err error) {

	_, err = db.Exec(`
	DELETE 
	FROM category
	WHERE id = ?`,
		id)

	return
}
