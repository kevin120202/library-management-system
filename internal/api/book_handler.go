package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kevin120202/library-management-system/internal/utils"
)

type BookHandler struct {
	Logger *log.Logger
}

func NewBookHandler(logger *log.Logger) *BookHandler {
	return &BookHandler{
		Logger: logger,
	}
}

func (bh *BookHandler) HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.Logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"book": fmt.Sprintf("book id is %d", bookID)})
}

func (bh *BookHandler) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"book": "this handler will create a book"})
}
