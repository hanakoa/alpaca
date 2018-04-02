package services

import (
	"github.com/gorilla/mux"
	"log"
	"github.com/kevinmichaelchen/my-go-utils"
	"fmt"
	"net/http"
	"github.com/hanakoa/alpaca/password-reset/models"
	"github.com/google/uuid"
)

func (svc *PasswordResetSvc) VerifyCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codeString := vars["code"]
	log.Println("Parsing:", codeString)
	if _, err := uuid.Parse(codeString); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c := &models.PasswordResetCode{Code: codeString}
	if valid, err := c.HasCode(svc.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else if valid {
		utils.RespondWithJSON(w, http.StatusOK, "Found valid code")
		return
	} else {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("No password reset code for: %s", codeString))
		return
	}
}