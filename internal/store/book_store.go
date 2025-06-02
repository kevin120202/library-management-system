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
	GetBooks() ([]Book, error)
	GetBookByID(id int64) (*Book, error)
	UpdateBook(*Book) error
	DeleteBook(id int64) error
	BorrowBook(bookID int64, userID int64) (error)
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

func (pg *PostgresBookStore) GetBooks() ([]Book, error) {
	var books []Book

	query := `
		SELECT id, title, author, summary, created_at, updated_at FROM books
	`

	rows, err := pg.db.Query(query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.Summary, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
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

func (pg *PostgresBookStore) UpdateBook(book *Book) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE books
		SET title = $1, author = $2, summary = $3
		WHERE id = $4
	`

	result, err := tx.Exec(query, book.Title, book.Author, book.Summary, book.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (pg *PostgresBookStore) DeleteBook(id int64) error {
	query := `
		DELETE from books
		WHERE id = $1
	`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pg *PostgresBookStore) BorrowBook(bookID int64, userID int64) (error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO borrows_returns (book_id, user_id)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(query, bookID, userID)
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE books
		SET availability_status = 'borrowed'
		WHERE id = $1
	`

	_, err = tx.Exec(updateQuery, bookID, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresBookStore) CheckIsBookAvailable(bookID int64) (*string, error) {
	return nil,nil
	// tx, err := pg.db.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer tx.Rollback()

	// var availabilityStatus string;

	// query := `SELECT availability_status FROM books WHERE $1`
	// err = tx.QueryRow(query, bookID).Scan(&availabilityStatus)
	// if err != nil {
	// 	return nil, err
	// }

	// if availabilityStatus == "borrowed" {
	// 	return 
	// }


	// err = tx.Commit()
	// if err != nil {
	// 	return err
	// }

	// return nil
}