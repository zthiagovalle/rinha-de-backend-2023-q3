package store

import "database/sql"

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
