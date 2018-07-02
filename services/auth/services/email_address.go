package services

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	"github.com/gorilla/mux"
	"github.com/hanakoa/alpaca/services/auth/models"
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
	"github.com/TeslaGov/cursor"
)

const (
	//200
	DefaultPageSize = 5
	MaxPageSize     = 1000
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
	} else {
		svc.EmailAddressSender = LogSender{}
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
	count := cursor.GetCount(r, DefaultPageSize, MaxPageSize)
	c := cursor.GetCursor(r)
	sort := cursor.GetSort(r)

	emailAddresses, err := models.GetEmailAddresses(svc.DB, int64(c), sort, count)
	if err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response interface{}
	if len(emailAddresses) != 0 {
		var data = make([]interface{}, len(emailAddresses))
		for i, e := range emailAddresses {
			data[i] = e
		}

		lastId := emailAddresses[len(emailAddresses) - 1].ID
		response = cursor.MakePage(count, data, int(c), int(lastId))
	} else {
		response = cursor.EmptyPage()
	}
	requestUtils.RespondWithJSON(w, http.StatusOK, response)
}

func (svc *EmailAddressService) GetEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := requestUtils.GetInt64(w, vars, "id")
	if !ok {
		return
	}

	rabbitmq.Send(svc.EmailAddressSender, "Getting email address...")
	log.Printf("Looking up email address: %d\n", id)
	e := models.EmailAddress{ID: id}
	if err := e.GetEmailAddress(svc.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			requestUtils.RespondWithError(w, http.StatusNotFound, "EmailAddress not found")
		default:
			requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	requestUtils.RespondWithJSON(w, http.StatusOK, e)
}

// TODO only admins can create, unless person is you
func (svc *EmailAddressService) CreateEmailAddress(w http.ResponseWriter, r *http.Request) {
	var e models.EmailAddress
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	now := time.Now()
	e.Created = null.TimeFrom(now)
	e.LastModified = null.TimeFrom(now)

	if e.ID != 0 {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Do not provide an id.")
		return
	}

	e.EmailAddress = strings.TrimSpace(e.EmailAddress)
	if e.EmailAddress == "" {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Must supply email address.")
		return
	}

	if len(e.EmailAddress) > 255 {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address cannot exceed 255 chars.")
		return
	}

	if e.Confirmed {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Cannot create confirmed email address.")
		return
	}

	if confirmed, err := e.IsConfirmed(svc.DB); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if confirmed {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address already exists and is confirmed.")
		return
	}

	if e.PersonID == 0 {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address must have person ID.")
		return
	}

	if err := checkmail.ValidateFormat(e.EmailAddress); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address has invalid format.")
		return
	}

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	p := &models.Person{Id: e.PersonID}
	if exists, err := p.Exists(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if !exists {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("No person found for person ID: %d", e.PersonID))
		return
	}

	e.ID = snowflakeUtils.NewPrimaryKey(svc.SnowflakeNode)
	e.IdStr = strconv.FormatInt(e.ID, 10)
	e.PersonIdStr = strconv.FormatInt(e.PersonID, 10)
	if err := e.CreateEmailAddress(tx); err != nil {
		tx.Rollback()
		log.Println("Could not create email address")
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		// TODO publish RabbitMQ message for email confirmation code
		rabbitmq.Send(svc.EmailAddressSender, "Created email address...")
		requestUtils.RespondWithJSON(w, http.StatusCreated, e)
	}
}

func (svc *EmailAddressService) UpdateEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := requestUtils.GetInt64(w, vars, "id")
	if !ok {
		return
	}

	var e models.EmailAddress
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	e.ID = id

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	if exists, err := e.Exists(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !exists {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("No email address found for id: %d", id))
		return
	}

	if err := e.UpdateEmailAddress(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := e.GetEmailAddress(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	e.IdStr = strconv.FormatInt(e.ID, 10)
	e.PersonIdStr = strconv.FormatInt(e.PersonID, 10)

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		requestUtils.RespondWithJSON(w, http.StatusOK, e)
	}
}

// TODO when emails are deleted
// deleter must own that email unless they're an admin
// nobody (not even admins) can delete primary emails
func (svc *EmailAddressService) DeleteEmailAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := requestUtils.GetInt64(w, vars, "id")
	if !ok {
		return
	}

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	// TODO you can only delete your own email address, unless you're an admin
	e := models.EmailAddress{ID: id, Deleted: null.TimeFrom(time.Now())}

	if exists, err := e.Exists(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else if !exists {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("No email address found for id: %d", id))
		return
	}

	if err := e.GetEmailAddress(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if e.Primary {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address is primary.")
		return
	}

	if err := e.DeleteEmailAddress(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := e.GetDeletedEmailAddress(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	e.IdStr = strconv.FormatInt(e.ID, 10)
	e.PersonIdStr = strconv.FormatInt(e.PersonID, 10)

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		requestUtils.RespondWithJSON(w, http.StatusOK, e)
	}
}