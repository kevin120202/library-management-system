package store

import "database/sql"

type PostgresBorrowReturnStore struct {
	db *sql.DB
}

func NewPostgresBorrowReturnStore(db *sql.DB) *PostgresBorrowReturnStore {
	return &PostgresBorrowReturnStore{db: db}
}

type BorrowBookStore interface {
	BorrowBook(id int64) (*Book, error)
	ReturnBook(id int64)
	RenewBook(id int64) (*Book, error)
}

func (pg *PostgresBorrowReturnStore) BorrowBook(bookID int64, userID int64) (error) {
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
