package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kevin120202/library-management-system/internal/middleware"
	"github.com/kevin120202/library-management-system/internal/store"
	"github.com/kevin120202/library-management-system/internal/utils"
)

type BookHandler struct {
	BookStore store.BookStore
	Logger    *log.Logger
}

func NewBookHandler(bookStore store.BookStore, logger *log.Logger) *BookHandler {
	return &BookHandler{
		BookStore: bookStore,
		Logger:    logger,
	}
}

// @desc    Get single book
// @route   Get /api/books/{id}
// @access  Public
func (bh *BookHandler) HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.Logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	book, err := bh.BookStore.GetBookByID(bookID)
	if err != nil {
		bh.Logger.Printf("ERROR: getBookByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"book": book})
}

// @desc    Get books
// @route   Get /api/books
// @access  Public
func (bh *BookHandler) HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := bh.BookStore.GetBooks()
	if err != nil {
		bh.Logger.Printf("ERROR: getBooks: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"books": books})
}

// @desc    Create a book
// @route   POST /api/books
// @access  Admin
func (bh *BookHandler) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book store.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		bh.Logger.Printf("ERROR: decodingCreateBook: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	if currentUser.AccountType != "admin" {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to create a book"})
		return
	}

	createdBook, err := bh.BookStore.CreateBook(&book)
	if err != nil {
		bh.Logger.Printf("ERROR: createBook: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create book"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"book": createdBook})
}

// @desc    Update a book
// @route   PUT /api/books/{id}
// @access  Admin
func (bh *BookHandler) HandleUpdateBookByID(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.Logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	existingBook, err := bh.BookStore.GetBookByID(bookID)
	if err != nil {
		bh.Logger.Printf("ERROR: getBookByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingBook == nil {
		http.NotFound(w, r)
		return
	}

	var updatedBookRequest struct {
		Title   *string `json:"title"`
		Author  *string `json:"author"`
		Summary *string `json:"summary"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedBookRequest)
	if err != nil {
		bh.Logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request paylaod"})
		return
	}

	if updatedBookRequest.Title != nil {
		existingBook.Title = *updatedBookRequest.Title
	}
	if updatedBookRequest.Author != nil {
		existingBook.Author = *updatedBookRequest.Author
	}
	if updatedBookRequest.Summary != nil {
		existingBook.Summary = *updatedBookRequest.Summary
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	if currentUser.AccountType != "admin" {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to create a book"})
		return
	}

	err = bh.BookStore.UpdateBook(existingBook)
	if err != nil {
		bh.Logger.Printf("ERROR: updatingBook: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"book": existingBook})
}

// @desc    Delete a book
// @route   DELETE /api/books/{id}
// @access  Admin
func (bh *BookHandler) HandleDeleteBookByID(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.Logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	if currentUser.AccountType != "admin" {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to create a book"})
		return
	}

	err = bh.BookStore.DeleteBook(bookID)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error deleting workout"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{"updated": true})
}

// @desc    Borrow a book
// @route   POST /api/books/{id}
// @access  Private
func (bh *BookHandler) HandleBorrowBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.Logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	err = bh.BookStore.BorrowBook(bookID, int64(currentUser.ID))
	if err != nil {
		bh.Logger.Printf(w, "ERROR: HandleBorrowBook: %v", err)
		return
	}

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error deleting workout"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{"updated": true})
}