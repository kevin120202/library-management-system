package store

import (
	"database/sql"
	"time"
)

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresBookStore struct {
	db *sql.DB
}

func NewPostgresBookStore(db *sql.DB) *PostgresBookStore {
	return &PostgresBookStore{db: db}
}

type BookStore interface {
	CreateBook(*Book) (*Book, error)
	GetBookByID(id int64) (*Book, error)
	// UpdateBook(*Book) error
	// DeleteBook(id int64) error
}

func (pg *PostgresBookStore) CreateBook(book *Book) (*Book, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO books (title, author, summary)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = tx.QueryRow(query, book.Title, book.Author, book.Summary).Scan(&book.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (pg *PostgresBookStore) GetBookByID(id int64) (*Book, error) {
	book := &Book{}

	query := `
		SELECT id, title, author, summary, created_at, updated_at FROM books
		WHERE id = $1
	`

	err := pg.db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Summary, &book.CreatedAt, &book.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return book, nil
}
