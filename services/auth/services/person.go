package services

import (
	"database/sql"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/hanakoa/alpaca/services/auth/models"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"time"
	"gopkg.in/guregu/null.v3"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"strings"
	"fmt"
	"log"
	"github.com/kevinmichaelchen/my-go-utils/rabbitmq"
	"github.com/badoux/checkmail"
	"regexp"
	"github.com/TeslaGov/cursor"
)

const (
	MinUsernameLength = 4
	MaxUsernameLength = 25
)

type AccountService struct {
	DB     *sql.DB
	SnowflakeNode *snowflake.Node
	AccountSender rabbitmq.Sender
}

type CreateAccountRequest struct {
	EmailAddress string      `json:"email_address"`
	Username     null.String `json:"username"`
}

type LogSender struct {
}

func (l LogSender) Send(i interface{}) {
	log.Println("Sending message: " + i.(string))
}

func NewAccountService(db *sql.DB, snowflakeNode *snowflake.Node, rabbitmqEnabled bool) AccountService {
	svc := AccountService{DB: db, SnowflakeNode: snowflakeNode, AccountSender: nil}
	if rabbitmqEnabled {
		svc.AccountSender = rabbitmq.NewRabbitSender("alpaca-auth-exchange", "account.#")
	} else {
		svc.AccountSender = LogSender{}
	}
	return svc
}

// TODO only admins can call this endpoint
func (svc *AccountService) GetAccounts(w http.ResponseWriter, r *http.Request) {
	count := cursor.GetCount(r, DefaultPageSize, MaxPageSize)
	c := cursor.GetCursor(r)
	sort := cursor.GetSort(r)

	accounts, err := models.GetAccounts(svc.DB, int64(c), sort, count)
	if err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response interface{}
	if len(accounts) != 0 {
		var data = make([]interface{}, len(accounts))
		for i, p := range accounts {
			data[i] = p
		}

		lastId := accounts[len(accounts)-1].Id
		response = cursor.MakePage(count, data, c, int(lastId))
	} else {
		response = cursor.EmptyPage()
	}
	requestUtils.RespondWithJSON(w, http.StatusOK, response)
	rabbitmq.Send(svc.AccountSender, "Getting accounts...")
}

func (svc *AccountService) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := requestUtils.GetInt64(w, vars, "accountId")
	if !ok {
		return
	}

	log.Printf("Looking up account: %d\n", id)
	p := models.Account{Id: id}
	if err := p.GetAccount(svc.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			requestUtils.RespondWithError(w, http.StatusNotFound, "Account not found")
		default:
			requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	rabbitmq.Send(svc.AccountSender, "getting account...")
	setStringsForAccount(&p)
	requestUtils.RespondWithJSON(w, http.StatusOK, p)
}

// TODO only admins can create
func (svc *AccountService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	p := &models.Account{}
	var req CreateAccountRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Println("Invalid request payload")
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
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
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Must supply email address.")
		return
	}

	if len(req.EmailAddress) > 255 {
		log.Println("Email address cannot exceed 255 chars.")
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address cannot exceed 255 chars.")
		return
	}

	if err := checkmail.ValidateFormat(req.EmailAddress); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Email address has invalid format.")
		return
	}

	if p.Username.Valid {
		username := strings.TrimSpace(p.Username.String)
		p.Username = null.StringFrom(username)
		if username == "" {
			requestUtils.RespondWithError(w, http.StatusBadRequest, "Username must be non-empty.")
			return
		}
		usernameLen := len(username)
		if usernameLen > MaxUsernameLength || usernameLen < MinUsernameLength {
			requestUtils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username length must be between %d and %d.", MinUsernameLength, MaxUsernameLength))
			return
		}
		if !isValidUsername(username) {
			requestUtils.RespondWithError(w, http.StatusBadRequest, "A username can only contain alphanumeric characters (letters A-Z, numbers 0-9) with the exception of underscores.")
			return
		}
	}

	// TODO email address cannot already exist and be confirmed

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	accountId := snowflakeUtils.NewPrimaryKey(svc.SnowflakeNode)
	p.Id = accountId
	if err := p.CreateAccount(tx); err != nil {
		tx.Rollback()
		log.Printf("Could not create account: %s", err.Error())
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	emailAddressId := snowflakeUtils.NewPrimaryKey(svc.SnowflakeNode)
	emailAddress := &models.EmailAddress{ID: emailAddressId, Primary: true, EmailAddress: req.EmailAddress, AccountID: p.Id}
	if err := emailAddress.CreateEmailAddress(tx); err != nil {
		tx.Rollback()
		log.Println("Could not create email address")
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created email address %d for account %d\n", emailAddress.ID, p.Id)

	p.PrimaryEmailAddressID = null.IntFrom(emailAddressId)
	if err := p.UpdateAccount(tx); err != nil {
		tx.Rollback()
		log.Println("Could not set primary email address for account")
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		setStringsForAccount(p)
		rabbitmq.Send(svc.AccountSender, "created account")
		requestUtils.RespondWithJSON(w, http.StatusCreated, p)
	}
}

func (svc *AccountService) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := requestUtils.GetInt64(w, vars, "accountId")
	if !ok {
		return
	}

	var p models.Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.Id = id

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	// TODO do we need 2 calls?
	if exists, err := p.Exists(tx); err != nil {
		requestUtils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	} else if !exists {
		requestUtils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("No account found for id: %d", id))
		return
	}
	// TODO update disabled
	// TODO username must not be taken
	if !p.PrimaryEmailAddressID.Valid {
		requestUtils.RespondWithError(w, http.StatusBadRequest, "Must provide primary email address ID")
		return
	}

	if err := p.UpdateAccount(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := p.GetAccount(tx); err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		setStringsForAccount(&p)
		rabbitmq.Send(svc.AccountSender, "updated account")
		requestUtils.RespondWithJSON(w, http.StatusOK, p)
	}
}

func (svc *AccountService) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := requestUtils.GetInt64(w, vars, "accountId")
	if !ok {
		return
	}

	var tx *sql.Tx
	tx, err := sqlUtils.StartTransaction(w, svc.DB); if err != nil {
		return
	}

	// TODO you can only delete yourself, unless you're an admin
	p := models.Account{Id: id, Deleted: null.TimeFrom(time.Now())}
	if err := p.DeleteAccount(tx); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO delete email addresses

	// Load new fields, like deleted_timestamp
	if err := p.GetDeletedAccount(tx); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		setStringsForAccount(&p)
		rabbitmq.Send(svc.AccountSender, "deleted account")
		requestUtils.RespondWithJSON(w, http.StatusOK, p)
	}
}

func setStringsForAccount(p *models.Account) {
	p.IdStr = strconv.FormatInt(p.Id, 10)
	// TODO PrimaryEmailAddressID should not be nullable (because email#accountId is not-nullable)
	p.PrimaryEmailAddressIdStr = strconv.FormatInt(p.PrimaryEmailAddressID.Int64, 10)
	p.CurrentPasswordIdStr = strconv.FormatInt(p.CurrentPasswordID.Int64, 10)
}

func isValidUsername(username string) bool {
	match, _ := regexp.MatchString("^[A-Za-z0-9_]*$", username)
	return match
}