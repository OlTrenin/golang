package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"number-service/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func InitDB(host, port, user, password, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS numbers (
			id SERIAL PRIMARY KEY,
			value INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_numbers_value ON numbers(value);
	`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

func (r *PostgresRepository) Save(value int) error {
	query := "INSERT INTO numbers (value) VALUES ($1)"
	_, err := r.db.Exec(query, value)
	if err != nil {
		return fmt.Errorf("failed to save number: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetAllSorted() ([]int, error) {
	query := "SELECT value FROM numbers ORDER BY value ASC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query numbers: %w", err)
	}
	defer rows.Close()

	var numbers []int
	for rows.Next() {
		var value int
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		numbers = append(numbers, value)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if numbers == nil {
		numbers = []int{}
	}

	return numbers, nil
}

var _ domain.NumberRepository = (*PostgresRepository)(nil)
