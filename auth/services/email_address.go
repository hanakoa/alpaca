package services

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"github.com/kevinmichaelchen/my-go-utils"
	"github.com/gorilla/mux"
	"github.com/hanakoa/alpaca/auth/models"
	"gopkg.in/guregu/null.v3"
	"time"
	"strconv"
	"encoding/json"
	"log"
	"strings"
	"fmt"
	"github.com/kevinmichaelchen/my-go-utils/rabbitmq"
	"net/http"
	"github.com/badoux/checkmail"
)

type EmailAddressService struct {
	DB     *sql.DB
	SnowflakeNode *snowflake.Node
	EmailAddressSender rabbitmq.Sender
}

func NewEmailAddressService(db *sql.DB, snowflakeNode *snowflake.Node, rabbitmqEnabled bool) EmailAddressService {
	svc := EmailAddressService{DB: db, SnowflakeNode: snowflakeNode, EmailAddressSender: nil}
	if rabbitmqEnabled {
		svc.EmailAddressSender = rabbitmq.NewRabbitSender("alpaca-auth-exchange", "emailAddress.#")
	}
	return svc
}

func setStringsForEmailAddress(e *models.EmailAddress) {
	e.IdStr = strconv.FormatInt(e.ID, 10)
	// TODO PrimaryEmailAddressID should not be nullable
	e.PersonIdStr = strconv.FormatInt(e.PersonID, 10)
}

// TODO only admins can call this endpoint
func (svc *EmailAddressService) GetEmailAddresses(w http.ResponseWriter, r *http.Request) {
	count := getCount(r)
	cursor := getCursor(r)
	sort := getSort(r)

	emailAddresses, err := models.GetEmailAddresses(svc.DB, cursor, sort, count)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response interface{}
	if len(emailAddresses) != 0 {
		var data = make([]interface{}, len(emailAddresses))
		for i, e := range emailAddresses {
			data[i] = e
		}

		lastId := emailAddresses[len(emailAddresses) - 1].ID
		response = makePage(count, data, cursor, lastId)
	} else {
		response = emptyPage()
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (svc *EmailAddressService) GetEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := utils.GetInt64(w, vars, "id")
	if !ok {
		return
	}

	svc.EmailAddressSender.Send("getting email address...")
	log.Printf("Looking up email address: %d\n", id)
	e := models.EmailAddress{ID: id}
	if err := e.GetEmailAddress(svc.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "EmailAddress not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, e)
}

// TODO only admins can create, unless personId is you
func (svc *EmailAddressService) CreateEmailAddress(w http.ResponseWriter, r *http.Request) {
	var e models.EmailAddress
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	now := time.Now()
	e.Created = null.TimeFrom(now)
	e.LastModified = null.TimeFrom(now)

	if e.ID != 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Do not provide an id.")
		return
	}

	e.EmailAddress = strings.TrimSpace(e.EmailAddress)
	if e.EmailAddress == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Must supply email address.")
		return
	}

	if len(e.EmailAddress) > 255 {
		utils.RespondWithError(w, http.StatusBadRequest, "Email address cannot exceed 255 chars.")
		return
	}

	if e.Confirmed {
		utils.RespondWithError(w, http.StatusBadRequest, "Cannot create confirmed email address.")
		return
	}

	if confirmed, err := e.IsConfirmed(svc.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if confirmed {
		utils.RespondWithError(w, http.StatusBadRequest, "Email address already exists and is confirmed.")
		return
	}

	if e.PersonID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Email address must have personId.")
		return
	}

	if err := checkmail.ValidateFormat(e.EmailAddress); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Email address has invalid format.")
		return
	}

	var tx *sql.Tx
	tx, err := utils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	p := &models.Person{Id: e.PersonID}
	if exists, err := p.Exists(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if !exists {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("No person found for personId: %d", e.PersonID))
		return
	}

	e.ID = utils.NewPrimaryKey(svc.SnowflakeNode)
	e.IdStr = strconv.FormatInt(e.ID, 10)
	e.PersonIdStr = strconv.FormatInt(e.PersonID, 10)
	if err := e.CreateEmailAddress(tx); err != nil {
		tx.Rollback()
		log.Println("Could not create email address")
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		// TODO publish RabbitMQ message for email confirmation code
		svc.EmailAddressSender.Send("created email address")
		utils.RespondWithJSON(w, http.StatusCreated, e)
	}
}

func (svc *EmailAddressService) UpdateEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := utils.GetInt64(w, vars, "id")
	if !ok {
		return
	}

	var e models.EmailAddress
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	e.ID = id

	var tx *sql.Tx
	tx, err := utils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	if exists, err := e.Exists(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !exists {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("No email address found for id: %d", id))
		return
	}

	if err := e.UpdateEmailAddress(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := e.GetEmailAddress(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	e.IdStr = strconv.FormatInt(e.ID, 10)
	e.PersonIdStr = strconv.FormatInt(e.PersonID, 10)

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		utils.RespondWithJSON(w, http.StatusOK, e)
	}
}

// TODO when emails are deleted
// deleter must own that email unless they're an admin
// nobody (not even admins) can delete primary emails
func (svc *EmailAddressService) DeleteEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := utils.GetInt64(w, vars, "id")
	if !ok {
		return
	}

	var tx *sql.Tx
	tx, err := utils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	// TODO you can only delete your own email address, unless you're an admin
	e := models.EmailAddress{ID: id, Deleted: null.TimeFrom(time.Now())}

	if exists, err := e.Exists(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !exists {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("No email address found for id: %d", id))
		return
	}

	if err := e.GetEmailAddress(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if e.Primary {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusBadRequest, "Email address is primary.")
		return
	}

	if err := e.DeleteEmailAddress(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := e.GetDeletedEmailAddress(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	e.IdStr = strconv.FormatInt(e.ID, 10)
	e.PersonIdStr = strconv.FormatInt(e.PersonID, 10)

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		utils.RespondWithJSON(w, http.StatusOK, e)
	}
}