package services

import (
	"database/sql"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/hanakoa/alpaca/auth/models"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"time"
	"gopkg.in/guregu/null.v3"
	"github.com/kevinmichaelchen/my-go-utils"
	"strings"
	"fmt"
	"log"
	"github.com/kevinmichaelchen/my-go-utils/rabbitmq"
	"github.com/badoux/checkmail"
	"regexp"
)

const (
	MinUsernameLength = 4
	MaxUsernameLength = 25
)

type PersonService struct {
	DB     *sql.DB
	SnowflakeNode *snowflake.Node
	PersonSender rabbitmq.Sender
}

type CreatePersonRequest struct {
	EmailAddress string      `json:"email_address"`
	Username     null.String `json:"username"`
}

type LogSender struct {
}

func (l LogSender) Send(i interface{}) {
	log.Println("Sending message: " + i.(string))
}

func NewPersonService(db *sql.DB, snowflakeNode *snowflake.Node, rabbitmqEnabled bool) PersonService {
	svc := PersonService{DB: db, SnowflakeNode: snowflakeNode, PersonSender: nil}
	if rabbitmqEnabled {
		svc.PersonSender = rabbitmq.NewRabbitSender("alpaca-auth-exchange", "person.#")
	} else {
		svc.PersonSender = LogSender{}
	}
	return svc
}

// TODO only admins can call this endpoint
func (svc *PersonService) GetPersons(w http.ResponseWriter, r *http.Request) {
	count := getCount(r)
	cursor := getCursor(r)
	sort := getSort(r)

	people, err := models.GetPersons(svc.DB, cursor, sort, count)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response interface{}
	if len(people) != 0 {
		var data = make([]interface{}, len(people))
		for i, p := range people {
			data[i] = p
		}

		lastId := people[len(people)-1].Id
		response = makePage(count, data, cursor, lastId)
	} else {
		response = emptyPage()
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
	svc.PersonSender.Send("Getting people")
}

func (svc *PersonService) GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := utils.GetInt64(w, vars, "personId")
	if !ok {
		return
	}

	log.Printf("Looking up person: %d\n", id)
	p := models.Person{Id: id}
	if err := p.GetPerson(svc.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "Person not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	svc.PersonSender.Send("getting person...")
	setStringsForPerson(&p)
	utils.RespondWithJSON(w, http.StatusOK, p)
}

// TODO only admins can create
func (svc *PersonService) CreatePerson(w http.ResponseWriter, r *http.Request) {
	p := &models.Person{}
	var req CreatePersonRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Println("Invalid request payload")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	now := time.Now()
	p.Created = null.TimeFrom(now)
	p.LastModified = null.TimeFrom(now)
	if req.Username.Valid {
		req.Username = null.StringFrom(strings.TrimSpace(req.Username.String))
	}
	p.Username = req.Username

	req.EmailAddress = strings.TrimSpace(req.EmailAddress)
	if req.EmailAddress == "" {
		log.Println("Must supply email address.")
		utils.RespondWithError(w, http.StatusBadRequest, "Must supply email address.")
		return
	}

	if len(req.EmailAddress) > 255 {
		log.Println("Email address cannot exceed 255 chars.")
		utils.RespondWithError(w, http.StatusBadRequest, "Email address cannot exceed 255 chars.")
		return
	}

	if err := checkmail.ValidateFormat(req.EmailAddress); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Email address has invalid format.")
		return
	}

	if p.Username.Valid {
		username := strings.TrimSpace(p.Username.String)
		p.Username = null.StringFrom(username)
		if username == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Username must be non-empty.")
			return
		}
		usernameLen := len(username)
		if usernameLen > MaxUsernameLength || usernameLen < MinUsernameLength {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username length must be between %d and %d.", MinUsernameLength, MaxUsernameLength))
			return
		}
		if !isValidUsername(username) {
			utils.RespondWithError(w, http.StatusBadRequest, "A username can only contain alphanumeric characters (letters A-Z, numbers 0-9) with the exception of underscores.")
			return
		}
	}

	// TODO email address cannot already exist and be confirmed

	var tx *sql.Tx
	tx, err := utils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	personId := utils.NewPrimaryKey(svc.SnowflakeNode)
	p.Id = personId
	if err := p.CreatePerson(tx); err != nil {
		tx.Rollback()
		log.Printf("Could not create person: %s", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	emailAddressId := utils.NewPrimaryKey(svc.SnowflakeNode)
	emailAddress := &models.EmailAddress{ID: emailAddressId, Primary: true, EmailAddress: req.EmailAddress, PersonID: p.Id}
	if err := emailAddress.CreateEmailAddress(tx); err != nil {
		tx.Rollback()
		log.Println("Could not create email address")
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created email address %d for person %d\n", emailAddress.ID, p.Id)

	p.PrimaryEmailAddressID = null.IntFrom(emailAddressId)
	if err := p.UpdatePerson(tx); err != nil {
		tx.Rollback()
		log.Println("Could not set primary email address for person")
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("PERSON CREATE - COMMIT FAILED")
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		setStringsForPerson(p)
		rabbitmq.Send(svc.PersonSender, "created person")
		utils.RespondWithJSON(w, http.StatusCreated, p)
	}
}

func (svc *PersonService) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := utils.GetInt64(w, vars, "personId")
	if !ok {
		return
	}

	var p models.Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.Id = id

	var tx *sql.Tx
	tx, err := utils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	// TODO do we need 2 calls?
	if exists, err := p.Exists(tx); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	} else if !exists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("No person found for id: %d", id))
		return
	}
	// TODO update disabled
	// TODO username must not be taken
	if !p.PrimaryEmailAddressID.Valid {
		utils.RespondWithError(w, http.StatusBadRequest, "Must provide primary email address ID")
		return
	}

	if err := p.UpdatePerson(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := p.GetPerson(tx); err != nil {
		tx.Rollback()
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		setStringsForPerson(&p)
		svc.PersonSender.Send("updated person")
		utils.RespondWithJSON(w, http.StatusOK, p)
	}
}

func (svc *PersonService) DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := utils.GetInt64(w, vars, "personId")
	if !ok {
		return
	}

	var tx *sql.Tx
	tx, err := utils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	// TODO you can only delete yourself, unless you're an admin
	p := models.Person{Id: id, Deleted: null.TimeFrom(time.Now())}
	if err := p.DeletePerson(tx); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO delete email addresses

	// Load new fields, like deleted_timestamp
	if err := p.GetDeletedPerson(tx); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		setStringsForPerson(&p)
		svc.PersonSender.Send("deleted person")
		utils.RespondWithJSON(w, http.StatusOK, p)
	}
}

func setStringsForPerson(p *models.Person) {
	p.IdStr = strconv.FormatInt(p.Id, 10)
	// TODO PrimaryEmailAddressID should not be nullable (because email#personId is not-nullable)
	p.PrimaryEmailAddressIdStr = strconv.FormatInt(p.PrimaryEmailAddressID.Int64, 10)
	p.CurrentPasswordIdStr = strconv.FormatInt(p.CurrentPasswordID.Int64, 10)
}

func isValidUsername(username string) bool {
	match, _ := regexp.MatchString("^[A-Za-z0-9_]*$", username)
	return match
}