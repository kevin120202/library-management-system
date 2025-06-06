package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/kevin120202/library-management-system/internal/middleware"
	"github.com/kevin120202/library-management-system/internal/store"
	"github.com/kevin120202/library-management-system/internal/utils"
)

type registerUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	AccountType string `json:"account_type"`
}

type UserHandler struct {
	userStore  store.UserStore
	tokenStore store.TokenStore
	logger     *log.Logger
}

func NewUserHandler(userStore store.UserStore, tokenStore store.TokenStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if strings.ToLower(req.AccountType) != "user" && strings.ToLower(req.AccountType) != "admin" {
		return errors.New("enter valid account type")
	}
	if req.Address == "" {
		return errors.New("address is required")
	}

	return nil
}

// @desc    Create a user
// @route   POST /api/users
// @access  Public
func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username:    req.Username,
		Email:       req.Email,
		AccountType: req.AccountType,
		Address:     req.Address,
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: registering user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid request payload"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}

// @desc    Logout a user
// @route   POST /api/logout
// @access  Private
func (h *UserHandler) HandleLogoutUser(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	err := h.tokenStore.DeleteAllTokensForUser(currentUser.ID, currentUser.AccountType)
	if err != nil {
		h.logger.Printf("ERROR: HandleLogoutUser: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to lougout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"Logout": true})
}
