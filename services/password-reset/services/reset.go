package services

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"encoding/json"
	"net/http"
	"strings"
	"github.com/badoux/checkmail"
	"fmt"
	"github.com/google/uuid"
	authGRPC "github.com/hanakoa/alpaca/services/auth/grpc"
	"github.com/hanakoa/alpaca/services/password-reset/models"
)

type PasswordResetSvc struct {
	DB            *sql.DB
	SnowflakeNode *snowflake.Node
	PassClient    authGRPC.PassClient
}

type PasswordResetRequest struct {
	Code        string `json:"code"`
	Account     string `json:"email_address"`
	NewPassword string `json:"password"`
}

func (svc *PasswordResetSvc) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var p PasswordResetRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB);
	if err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if strings.TrimSpace(p.NewPassword) == "" {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, "Password cannot be empty.")
		return
	}

	if strings.TrimSpace(p.Account) == "" {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, "Email address cannot be empty.")
		return
	}

	if err := checkmail.ValidateFormat(p.Account); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Email address has invalid format: %s", err.Error()))
		return
	}

	codeString := p.Code
	if u, err := uuid.Parse(p.Code); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		codeString = u.String()
	}

	c := &models.PasswordResetCode{Code: codeString}
	if valid, err := c.HasCode(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !valid {
		tx.Rollback()
		p.NewPassword = ""
		requestUtils.RespondWithJSON(w, http.StatusOK, p)
		return
	}

	var accountID int64
	if accountID, err = GetAccountIdForAccount(p.Account, tx); err != nil || accountID == 0 {
		tx.Rollback()
		p.NewPassword = ""
		// We deliberately do not leak if email is not found
		requestUtils.RespondWithJSON(w, http.StatusOK, p)
		return
	}

	c.AccountID = accountID

	if err := authGRPC.ResetPassword(svc.PassClient, accountID, p.NewPassword); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.MarkAsUsed(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.MarkAllAsUnusable(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		p.NewPassword = ""
		requestUtils.RespondWithJSON(w, http.StatusOK, p)
	}
}
