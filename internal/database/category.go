package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (c *Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			return nil, err
		}

		categories = append(categories, Category{ID: id, Name: name, Description: description})
	}

	return categories, nil
}

func (c *Category) Find(id string) (Category, error) {
	stmt, err := c.db.Prepare("SELECT id, name, description FROM categories WHERE id = $1")
	if err != nil {
		fmt.Println("b", err)
		return Category{}, nil
	}
	defer stmt.Close()

	var category Category
	err = stmt.QueryRow(id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return Category{}, nil
	}

	return category, nil
}

func (c *Category) Create(name string, description string) (Category, error) {
	id := uuid.New().String()
	_, err := c.db.Exec(
		"INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)",
		id, name, description,
	)
	if err != nil {
		return Category{}, err
	}

	return Category{ID: id, Name: name, Description: description}, nil
}
