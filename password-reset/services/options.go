package services

import (
	"database/sql"
	"github.com/kevinmichaelchen/my-go-utils"
	"strings"
	"github.com/badoux/checkmail"
	"fmt"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"github.com/hanakoa/alpaca/password-reset/models"
	"strconv"
)

type SendCodeRequest struct {
	// Account can be an email address, phone number, or username
	Account string `json:"account"`
}

type SendCodeOptions struct {
	PhoneNumbers   []models.PhoneNumber  `json:"phone_numbers"`
	EmailAddresses []models.EmailAddress `json:"email_addresses"`
	// TODO eventually we'll add Yubikey devices and backup recovery codes
}

func (s *SendCodeOptions) NumOptions() int {
	num := 0
	if s.PhoneNumbers != nil {
		num += len(s.PhoneNumbers)
	}
	if s.EmailAddresses != nil {
		num += len(s.EmailAddresses)
	}
	return num
}

func (svc *PasswordResetSvc) SendCode(w http.ResponseWriter, r *http.Request) {
	var p SendCodeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var tx *sql.Tx
	tx, err := utils.StartTransaction(w, svc.DB);
	if err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if strings.TrimSpace(p.Account) == "" {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, "Account field cannot be empty.")
		return
	}

	var personID int64
	if personID, err = GetPersonIdForAccount(p.Account, tx); err != nil || personID == 0 {
		tx.Rollback()
		// We deliberately do not leak if email is not found
		// TODO RespondWithJSON doesn't actually return JSON
		utils.RespondWithJSON(w, http.StatusOK, "Reset request submitted.")
		return
	}

	sendCodeOptions, err := getSendOptions(personID, tx)
	if err != nil {
		tx.Rollback()
		// We deliberately do not leak if email is not found
		utils.RespondWithJSON(w, http.StatusOK, "Reset request submitted.")
		return
	}
	if sendCodeOptions.NumOptions() == 1 {
		// TODO actually send an email
		log.Println("Fake sending an email")

		expiration := time.Now().Add(time.Minute * 30)
		if resetCode, err := models.NewPasswordResetCode(personID, expiration); err != nil {
			tx.Rollback()
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			if err := resetCode.CreatePasswordResetCode(tx); err != nil {
				tx.Rollback()
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}

		if err := tx.Commit(); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		} else {
			utils.RespondWithJSON(w, http.StatusOK, "Reset request submitted.")
		}
		return
	}

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		utils.RespondWithJSON(w, http.StatusOK, sendCodeOptions)
	}
}

func getSendOptions(personID int64, tx *sql.Tx) (*SendCodeOptions, error) {
	options := &SendCodeOptions{}
	if phoneNumbers, err := models.GetPhoneNumbersForPerson(personID, tx); err != nil {
		return nil, err
	} else {
		options.PhoneNumbers = phoneNumbers
	}
	if emailAddresses, err := models.GetEmailAddressesForPerson(personID, tx); err != nil {
		return nil, err
	} else {
		options.EmailAddresses = emailAddresses
	}
	return options, nil
}

func GetPersonIdForAccount(account string, tx *sql.Tx) (int64, error) {
	if isEmailAddress(account) {
		if err := checkmail.ValidateFormat(account); err != nil {
			log.Println(err.Error())
			return 0, fmt.Errorf("email address has invalid format: %s", account)
		}
		emailAddress := &models.EmailAddress{EmailAddress: account}
		// TODO should be case insensitive
		emailAddress.GetConfirmedEmailAddress(tx)
		return emailAddress.PersonID, nil
	} else if isPhoneNumber(account) {
		phoneNumber := &models.PhoneNumber{PhoneNumber: account}
		phoneNumber.GetPhoneNumberByPhoneNumber(tx)
		return phoneNumber.PersonID, nil
	} else {
		person := &models.Person{Username: account}
		person.GetPersonByUsername(tx)
		return person.ID, nil
	}
}

func isEmailAddress(s string) bool {
	return strings.Contains(s, "@")
}

func isPhoneNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
