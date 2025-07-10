package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DerekBelloni/fem_project/internal/store"
	"github.com/DerekBelloni/fem_project/internal/tokens"
	"github.com/DerekBelloni/fem_project/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (t *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		t.logger.Printf("ERROR: createTokenRequest")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	fmt.Printf("handle create token, username: %v\n", req.Username)

	user, err := t.userStore.GetUserByUsername(req.Username)
	if err != nil {
		t.logger.Printf("ERROR: getUserByUsername")
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}

	passswordsMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		t.logger.Printf("ERROR: PasswordHash.Matches: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error:": "internal server error"})
		return
	}

	if !passswordsMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := t.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		t.logger.Printf("ERROR: Creating Token: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error:": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
