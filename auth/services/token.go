package services

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"net/http"
	"github.com/kevinmichaelchen/my-go-utils"
	mfaGRPC "github.com/hanakoa/alpaca/mfa/grpc"
	"encoding/json"
	"gopkg.in/guregu/null.v3"
	"strings"
	"time"
	"github.com/hanakoa/alpaca/auth/models"
	"log"
	"github.com/badoux/checkmail"
	"fmt"
	"github.com/google/uuid"
)

type LoginRequest struct {
	// Login is either an email address or username
	Login     string `json:"login"`
	Password  string `json:"password"`
}

type LoginResponse struct {
	// MfaCode is a UUID for the person's password reset.
	// We send it back so the UI can re-trigger re-sends.
	MfaCode uuid.UUID `json:"mfaCode"`
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	login := strings.TrimSpace(resource.Login)

	if login == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Must supply email address or username.")
		return
	}

	var person *models.Person
	var err error
	isEmail := strings.Contains(login, "@")
	if isEmail {
		emailAddress := login
		if len(emailAddress) > 255 {
			utils.RespondWithError(w, http.StatusBadRequest, "Email address should not exceed 255 chars.")
			return
		}
		if err := checkmail.ValidateFormat(emailAddress); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Malformed email address.")
			return
		}
		person = &models.Person{EmailAddress: emailAddress}
		err = person.GetPersonByEmailAddress(svc.DB)
	} else {
		username := login
		if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username length must be between %d and %d.", MinUsernameLength, MaxUsernameLength))
			return
		}
		person = &models.Person{Username: null.StringFrom(username)}
		err = person.GetPersonByUsername(svc.DB)
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("Authentication failed: %s", err.Error()))
		return
	}

	if !person.CurrentPasswordID.Valid {
		log.Printf("Person %d has no current password...\n", person.Id)
		utils.RespondWithError(w, http.StatusForbidden, "Authentication failed.")
		return
	}

	password := &models.Password{Id: person.CurrentPasswordID.Int64}
	if err = password.GetPasswordForPersonID(svc.DB); err != nil {
		log.Printf("No Person exists for password %d...\n", person.CurrentPasswordID.Int64)
		utils.RespondWithError(w, http.StatusForbidden, "Authentication failed.")
		return
	}

	success := models.MatchesHash(resource.Password, password)

	l := &models.LoginAttempt{Id: utils.NewPrimaryKey(svc.SnowflakeNode), Created: time.Now(), Success: success, PersonID: person.Id}
	if err := l.CreateLoginAttempt(svc.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if success {
		if person.MultiFactorRequired {
			log.Printf("2FA required for person %d", person.Id)
			if resetCodeID, err := uuid.NewRandom(); err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			} else {
				log.Printf("Reset code is %s", resetCodeID)
				// TODO check RABBITMQ_ENABLED
				if err := mfaGRPC.Send2FACode(svc.MFAClient, person.Id, resetCodeID); err != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
				mfaResponse := LoginResponse{MfaCode: resetCodeID}
				utils.RespondWithJSON(w, http.StatusOK, mfaResponse)
			}
		} else {
			utils.RespondWithJSON(w, http.StatusOK, map[string]string{"msg": "Authenticated"})
		}
	} else {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials.")
	}
}
