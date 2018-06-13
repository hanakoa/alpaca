package services

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/hanakoa/alpaca/auth/models"
	"encoding/json"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"github.com/nbutton23/zxcvbn-go"
	"log"
	"strconv"
	"gopkg.in/guregu/null.v3"
	"errors"
)

type PasswordService struct {
	DB     *sql.DB
	SnowflakeNode *snowflake.Node
	IterationCount int
}

func (svc *PasswordService) UpdatePasswordHelper(tx *sql.Tx, p *models.Password, personId int64) (int, error) {
	person := &models.Person{Id: personId}
	if exists, err := person.Exists(tx); err != nil {
		return http.StatusNotFound, err
	} else if !exists {
		return http.StatusNotFound, errors.New("Person does not exist")
	}

	// Load the person's current password
	if err := person.GetPerson(tx); err != nil {
		return http.StatusInternalServerError, err
	}

	// Set the password's personId
	p.PersonID = personId

	if person.CurrentPasswordID.Valid {
		log.Println("Person already has a password... Updating it...")
		if err := p.UpdatePassword(tx); err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		log.Println("Person does not have a password... Creating one...")
		// TODO add more user info (e.g., name, email) to slice argument
		// The userInputs argument is an splice of strings that zxcvbn will add to its internal dictionary.
		// This can be whatever list of strings you like, but is meant for user inputs from other fields
		// of the form, like name and email. That way a password that includes the user's personal info
		// can be heavily penalized.
		minEntropyMatch := zxcvbn.PasswordStrength(p.PasswordText.String, []string {"alpaca"})
		if minEntropyMatch.Score < 2 {
			return http.StatusBadRequest, errors.New("Password not strong enough")
		}
		p.Id = snowflakeUtils.NewPrimaryKey(svc.SnowflakeNode)
		if err := p.CreatePassword(tx, svc.IterationCount); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Get password so we populate fields
	p.GetPasswordForPersonID(tx)

	if err := person.UpdateCurrentPassword(tx, p.Id); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (svc *PasswordService) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personId, ok := requestUtils.GetInt64(w, vars, "personId")
	if !ok {
		return
	}

	var p models.Password
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if !p.PasswordText.Valid {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Password text must be non-null")
		return
	}

	if p.PasswordText.String == "" {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Password text must be non-empty")
		return
	}

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	if statusCode, err := svc.UpdatePasswordHelper(tx, &p, personId); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, statusCode, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		// TODO we null out sensitive fields, which is tedious. Should we use a separate DTO struct?
		// TODO we should hide null fields from being returned in JSON
		p.Salt = nil
		p.PasswordHash = nil
		p.PasswordText = null.StringFrom("")
		p.IdStr = strconv.FormatInt(p.Id, 10)
		p.PersonIdStr = strconv.FormatInt(p.PersonID, 10)
		p.IterationCount = 0
		requestUtils.RespondWithJSON(w, http.StatusOK, p)
	}
}
