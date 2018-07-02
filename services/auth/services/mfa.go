package services

import (
	"net/http"
	"github.com/google/uuid"
	"log"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	"github.com/hanakoa/alpaca/services/auth/models"
)

type LoginResponse struct {
	// MfaCode is a UUID for the person's password reset.
	// We send it back so the UI can re-trigger re-sends.
	MfaCode uuid.UUID `json:"mfa_code"`
}

func WriteMfaOptions(w http.ResponseWriter, person *models.Person) {
	log.Printf("2FA required for person %d", person.Id)
	if resetCodeID, err := uuid.NewRandom(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		log.Printf("Reset code is %s", resetCodeID)
		mfaResponse := LoginResponse{MfaCode: resetCodeID}
		requestUtils.RespondWithJSON(w, http.StatusOK, mfaResponse)
	}
}
