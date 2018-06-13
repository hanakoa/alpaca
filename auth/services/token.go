package services

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"net/http"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	mfaGRPC "github.com/hanakoa/alpaca/mfa/grpc"
	"encoding/json"
	"gopkg.in/guregu/null.v3"
	"strings"
	"time"
	"github.com/hanakoa/alpaca/auth/models"
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

	var person *models.Person
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
		person, err = models.GetPersonByEmailAddress(svc.DB, emailAddress)
	} else if isPhoneNumber(login) {
		person, err = models.GetPersonByPhoneNumber(svc.DB, login)
	} else {
		username := login
		if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
			requestUtils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username length must be between %d and %d.", MinUsernameLength, MaxUsernameLength))
			return
		}
		person = &models.Person{Username: null.StringFrom(username)}
		err = person.GetPersonByUsername(svc.DB)
	}

	if err != nil {
		requestUtils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Authentication failed: %s", err.Error()))
		return
	}

	if !person.CurrentPasswordID.Valid {
		log.Printf("Person %d has no current password...\n", person.Id)
		requestUtils.RespondWithError(w, http.StatusForbidden, "Authentication failed.")
		return
	}

	password := &models.Password{Id: person.CurrentPasswordID.Int64}
	if err = password.GetPasswordForPersonID(svc.DB); err != nil {
		log.Printf("No Person exists for password %d...\n", person.CurrentPasswordID.Int64)
		requestUtils.RespondWithError(w, http.StatusForbidden, "Authentication failed.")
		return
	}

	passwordCorrect := models.MatchesHash(resource.Password, password)

	l := &models.LoginAttempt{Id: snowflakeUtils.NewPrimaryKey(svc.SnowflakeNode), Created: time.Now(), Success: passwordCorrect, PersonID: person.Id}
	if err := l.CreateLoginAttempt(svc.DB); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if passwordCorrect {
		if person.MultiFactorRequired {
			WriteMfaOptions(w, person)
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
