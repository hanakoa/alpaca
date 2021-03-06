package services

import (
	"database/sql"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"strings"
	"github.com/badoux/checkmail"
	"fmt"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"github.com/hanakoa/alpaca/services/password-reset/models"
	"github.com/ttacon/libphonenumber"
)

type CodeOptionsRequest struct {
	// Account can be an email address, phone number, or username
	Account string `json:"account"`
}

// CodeOptionsResponse contains all possible ways a user can login with.
type CodeOptionsResponse struct {
	PhoneNumbers   []models.PhoneNumber  `json:"phone_numbers"`
	EmailAddresses []models.EmailAddress `json:"email_addresses"`
	// TODO eventually we'll add Yubikey devices and backup recovery codes
}

func (s *CodeOptionsResponse) NumOptions() int {
	num := 0
	if s.PhoneNumbers != nil {
		num += len(s.PhoneNumbers)
	}
	if s.EmailAddresses != nil {
		num += len(s.EmailAddresses)
	}
	return num
}

func (svc *PasswordResetSvc) SendCodeOptions(w http.ResponseWriter, r *http.Request) {
	var p CodeOptionsRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB);
	if err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if strings.TrimSpace(p.Account) == "" {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, "Account field cannot be empty.")
		return
	}

	var accountID int64
	if accountID, err = GetAccountIdForAccount(p.Account, tx); err != nil || accountID == 0 {
		tx.Rollback()
		// We deliberately do not leak if email is not found
		// TODO RespondWithJSON doesn't actually return JSON
		requestUtils.RespondWithJSON(w, http.StatusOK, "Reset request submitted.")
		return
	}

	sendCodeOptions, err := getSendOptions(accountID, tx)
	if err != nil {
		tx.Rollback()
		// We deliberately do not leak if email is not found
		requestUtils.RespondWithJSON(w, http.StatusOK, "Reset request submitted.")
		return
	}
	if sendCodeOptions.NumOptions() == 1 {
		// TODO actually send an email
		log.Println("Fake sending an email")

		expiration := time.Now().Add(time.Minute * 30)
		if resetCode, err := models.NewPasswordResetCode(accountID, expiration); err != nil {
			tx.Rollback()
			requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			if err := resetCode.CreatePasswordResetCode(tx); err != nil {
				tx.Rollback()
				requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}

		if err := tx.Commit(); err != nil {
			requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			requestUtils.RespondWithJSON(w, http.StatusOK, "Reset request submitted.")
		}
		return
	}

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		requestUtils.RespondWithJSON(w, http.StatusOK, sendCodeOptions)
	}
}

func getSendOptions(accountID int64, tx *sql.Tx) (*CodeOptionsResponse, error) {
	options := &CodeOptionsResponse{}
	if phoneNumbers, err := models.GetPhoneNumbersForAccount(accountID, tx); err != nil {
		return nil, err
	} else {
		options.PhoneNumbers = phoneNumbers
	}
	if emailAddresses, err := models.GetEmailAddressesForAccount(accountID, tx); err != nil {
		return nil, err
	} else {
		options.EmailAddresses = emailAddresses
	}
	return options, nil
}

func GetAccountIdForAccount(account string, tx *sql.Tx) (int64, error) {
	if isEmailAddress(account) {
		if err := checkmail.ValidateFormat(account); err != nil {
			log.Println(err.Error())
			return 0, fmt.Errorf("email address has invalid format: %s", account)
		}
		emailAddress := &models.EmailAddress{EmailAddress: account}
		// TODO should be case insensitive
		emailAddress.GetConfirmedEmailAddress(tx)
		return emailAddress.AccountID, nil
	} else if isPhoneNumber(account) {
		phoneNumber := &models.PhoneNumber{PhoneNumber: account}
		phoneNumber.GetPhoneNumberByPhoneNumber(tx)
		return phoneNumber.AccountID, nil
	} else {
		account := &models.Account{Username: account}
		account.GetAccountByUsername(tx)
		return account.ID, nil
	}
}

func isEmailAddress(s string) bool {
	return strings.Contains(s, "@")
}

func isPhoneNumber(s string) bool {
	_, err := libphonenumber.Parse(s, "US")
	return err == nil
}
