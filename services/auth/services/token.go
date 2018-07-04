package services

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"net/http"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	mfaGRPC "github.com/hanakoa/alpaca/services/mfa/grpc"
	"encoding/json"
	"strings"
	"time"
	"github.com/hanakoa/alpaca/services/auth/models"
	"log"
	"github.com/badoux/checkmail"
	"fmt"
	"github.com/ttacon/libphonenumber"
)

type LoginRequest struct {
	// Login is either an email address or username
	Login     string `json:"login"`
	Password  string `json:"password"`
}

type TokenService struct {
	DB            *sql.DB
	SnowflakeNode *snowflake.Node
	MFAClient     mfaGRPC.MFAClient
}

// Authenticate expects either an email address or username, and a password
func (svc *TokenService) Authenticate(w http.ResponseWriter, r *http.Request) {
	var resource LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&resource); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	login := strings.TrimSpace(resource.Login)

	if login == "" {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Must supply email address or username.")
		return
	}

	var account *models.Account
	var err error
	if isEmailAddress(login) {
		emailAddress := login
		if len(emailAddress) > 255 {
			requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address should not exceed 255 chars.")
			return
		}
		if err := checkmail.ValidateFormat(emailAddress); err != nil {
			requestUtils.RespondWithError(w, http.StatusBadRequest, "Malformed email address.")
			return
		}
		account, err = models.GetAccountByEmailAddress(svc.DB, emailAddress)
	} else if isPhoneNumber(login) {
		account, err = models.GetAccountByPhoneNumber(svc.DB, login)
	} else {
		username := login
		if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
			requestUtils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username length must be between %d and %d.", MinUsernameLength, MaxUsernameLength))
			return
		}
		account, err = models.GetAccountByUsername(svc.DB, username)
	}

	if err != nil {
		requestUtils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Authentication failed: %s", err.Error()))
		return
	}

	if !account.CurrentPasswordID.Valid {
		log.Printf("Account %d has no current password...\n", account.Id)
		requestUtils.RespondWithError(w, http.StatusForbidden, "Authentication failed.")
		return
	}

	password := &models.Password{Id: account.CurrentPasswordID.Int64}
	if err = password.GetPasswordForAccountID(svc.DB); err != nil {
		log.Printf("No Account exists for password %d...\n", account.CurrentPasswordID.Int64)
		requestUtils.RespondWithError(w, http.StatusForbidden, "Authentication failed.")
		return
	}

	passwordCorrect := models.MatchesHash(resource.Password, password)

	l := &models.LoginAttempt{Id: snowflakeUtils.NewPrimaryKey(svc.SnowflakeNode), Created: time.Now(), Success: passwordCorrect, AccountID: account.Id}
	if err := l.CreateLoginAttempt(svc.DB); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if passwordCorrect {
		if account.MultiFactorRequired {
			WriteMfaOptions(w, account)
		} else {
			requestUtils.RespondWithJSON(w, http.StatusOK, map[string]string{"msg": "Authenticated"})
		}
	} else {
		requestUtils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials.")
	}
}

func isEmailAddress(s string) bool {
	return strings.Contains(s, "@")
}

func isPhoneNumber(s string) bool {
	_, err := libphonenumber.Parse(s, "US")
	return err == nil
}
