package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

type Person struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"apelido"`
	Name      string    `json:"nome"`
	BirthDate DateOnly  `json:"nascimento"`
	Stack     *[]string `json:"stack"`
}

type DateOnly struct {
	time.Time
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%s\"", d.Time.Format(dateLayout))
	return []byte(s), nil
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
	GetPersonByID(id uuid.UUID) (*Person, error)
	GetPersonsByTerm(term string, limit int, offset int) ([]Person, error)
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
	err = tx.QueryRow(query, person.Username, person.Name, person.BirthDate.Time).Scan(&id)
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

func (pg *PostgresPersonStore) GetPersonByID(id uuid.UUID) (*Person, error) {
	query := `
	SELECT id, username, name, birth_date
	FROM persons
	WHERE id = $1
	`

	var person Person
	err := pg.db.QueryRow(query, id).Scan(&person.ID, &person.Username, &person.Name, &person.BirthDate.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	query = `
	SELECT name
	FROM person_stacks
	WHERE person_id = $1
	`
	rows, err := pg.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stacks []string
	for rows.Next() {
		var stack string
		if err := rows.Scan(&stack); err != nil {
			return nil, err
		}
		stacks = append(stacks, stack)
	}
	person.Stack = &stacks

	return &person, nil
}

func (pg *PostgresPersonStore) GetPersonsByTerm(term string, limit int, offset int) ([]Person, error) {
	query := `
	SELECT p.id, p.username, p.name, p.birth_date
	FROM persons p
	LEFT JOIN person_stacks ps ON p.id = ps.person_id
	WHERE p.username ILIKE '%' || $1 || '%'
	OR p.name ILIKE '%' || $1 || '%'
	OR ps.name ILIKE '%' || $1 || '%'
	GROUP BY p.id, p.username, p.name, p.birth_date
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := pg.db.Query(query, term, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []Person
	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.ID, &person.Username, &person.Name, &person.BirthDate.Time); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	if len(persons) == 0 {
		return persons, nil
	}

	for i, person := range persons {
		query = `
		SELECT name
		FROM person_stacks
		WHERE person_id = $1
		`
		rows, err := pg.db.Query(query, person.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var stacks []string
		for rows.Next() {
			var stack string
			if err := rows.Scan(&stack); err != nil {
				return nil, err
			}
			stacks = append(stacks, stack)
		}
		persons[i].Stack = &stacks
	}

	return persons, nil
}
