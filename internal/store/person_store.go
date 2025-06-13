package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db *pgxpool.Pool
}

func NewPostgresPersonStore(db *pgxpool.Pool) *PostgresPersonStore {
	return &PostgresPersonStore{
		db: db,
	}
}

type PersonStore interface {
	CountPersons() (int, error)
	CreatePerson(person *Person) (*uuid.UUID, error)
	GetPersonByID(id uuid.UUID) (*Person, error)
	GetPersonsByTerm(term string, limit int) ([]Person, error)
}

func (pg *PostgresPersonStore) CountPersons() (int, error) {
	var totalPersons int

	query := `
	SELECT COUNT(id)
	FROM persons
	`

	err := pg.db.QueryRow(context.Background(), query).Scan(&totalPersons)
	if err != nil {
		return 0, err
	}

	return totalPersons, nil
}

func (pg *PostgresPersonStore) CreatePerson(person *Person) (*uuid.UUID, error) {
	query := `
	INSERT INTO persons (username, name, birth_date, stack)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	var id uuid.UUID
	err := pg.db.QueryRow(context.Background(), query, person.Username, person.Name, person.BirthDate.Time, person.Stack).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "persons_username_key") {
			return nil, errors.New(ErrPersonUsernameAlreadyExists)
		}
		return nil, err
	}

	return &id, nil
}

func (pg *PostgresPersonStore) GetPersonByID(id uuid.UUID) (*Person, error) {
	query := `
	SELECT id, username, name, birth_date, stack
	FROM persons
	WHERE id = $1
	`

	var person Person
	err := pg.db.QueryRow(context.Background(), query, id).Scan(&person.ID, &person.Username, &person.Name, &person.BirthDate.Time, &person.Stack)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &person, nil
}

func (pg *PostgresPersonStore) GetPersonsByTerm(term string, limit int) ([]Person, error) {
	query := `
	SELECT p.id, p.username, p.name, p.birth_date, p.stack
	FROM persons p
	WHERE searchable ILIKE '%' || $1 || '%'
	LIMIT $2
	`
	rows, err := pg.db.Query(context.Background(), query, term, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := make([]Person, 0, limit)
	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.ID, &person.Username, &person.Name, &person.BirthDate.Time, &person.Stack); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}
