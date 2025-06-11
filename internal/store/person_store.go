package store

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Person struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birth_date"`
	Stack     *[]string `json:"stack"`
}

type PostgresPersonStore struct {
	db *sql.DB
}

func NewPostgresPersonStore(db *sql.DB) *PostgresPersonStore {
	return &PostgresPersonStore{
		db: db,
	}
}

type PersonStore interface {
	CountPersons() (int, error)
	CreatePerson(person *Person) (*uuid.UUID, error)
}

func (pg *PostgresPersonStore) CountPersons() (int, error) {
	var totalPersons int

	query := `
	SELECT COUNT(id)
	FROM persons
	`

	err := pg.db.QueryRow(query).Scan(&totalPersons)
	if err != nil {
		return 0, err
	}

	return totalPersons, nil
}

func (pg *PostgresPersonStore) CreatePerson(person *Person) (*uuid.UUID, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO persons (username, name, birth_date)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var id uuid.UUID
	err = tx.QueryRow(query, person.Username, person.Name, person.BirthDate).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "persons_username_key") {
			return nil, errors.New(ErrPersonUsernameAlreadyExists)
		}
		return nil, err
	}

	if person.Stack != nil && len(*person.Stack) > 0 {
		for _, stack := range *person.Stack {
			query = `
			INSERT INTO person_stacks (person_id, name)
			VALUES ($1, $2)
			`
			_, err = tx.Exec(query, id, stack)
			if err != nil {
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &id, nil
}
