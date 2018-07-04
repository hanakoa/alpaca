package main

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hanakoa/alpaca/services/auth/services"
	mfaGRPC "github.com/hanakoa/alpaca/services/mfa/grpc"
	"github.com/kevinmichaelchen/my-go-utils/rabbitmq"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type App struct {
	Router          *mux.Router
	DB              *sql.DB
	RabbitmqEnabled bool
	accountService  services.AccountService
	passwordService services.PasswordService
	emailAddressSvc services.EmailAddressService
	tokenService    services.TokenService
	snowflakeNode   *snowflake.Node
	mfaClient       mfaGRPC.MFAClient
	iterationCount  int
}

// Initialize initializes the database connection and services,
// and creates routes for our REST API.
func (a *App) Initialize(db *sql.DB, secret string, numWorkers int) {
	if a.RabbitmqEnabled {
		log.Println("Creating RabbitMQ dispatcher...")
		rabbitmq.NewDispatcher(numWorkers, 10)
	}

	a.DB = db
	a.initializeServices()
	a.initializeRoutes()
}

// initializeServices initializes the Service layer.
func (a *App) initializeServices() {
	a.accountService = services.NewAccountService(a.DB, a.snowflakeNode, a.RabbitmqEnabled)
	a.emailAddressSvc = services.NewEmailAddressService(a.DB, a.snowflakeNode, a.RabbitmqEnabled)
	a.passwordService = services.PasswordService{
		DB:             a.DB,
		SnowflakeNode:  a.snowflakeNode,
		IterationCount: a.iterationCount}
	a.tokenService = services.TokenService{
		DB:            a.DB,
		SnowflakeNode: a.snowflakeNode,
		MFAClient:     a.mfaClient}
}

func (a *App) initializeRoutes() {
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/token", a.tokenService.Authenticate).Methods("POST")

	a.Router.HandleFunc("/account/{accountId:[0-9]+}/password", a.passwordService.UpdatePassword).Methods("PUT")

	a.Router.HandleFunc("/account", a.accountService.GetAccounts).Methods("GET")
	a.Router.HandleFunc("/account/{accountId:[0-9]+}", a.accountService.GetAccount).Methods("GET")
	a.Router.HandleFunc("/account", a.accountService.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/account/{accountId:[0-9]+}", a.accountService.UpdateAccount).Methods("PUT")
	a.Router.HandleFunc("/account/{accountId:[0-9]+}", a.accountService.DeleteAccount).Methods("DELETE")

	// TODO get account by primary email address
	// TODO get account by email address
	// TODO get account by username

	a.Router.HandleFunc("/emailaddress", a.emailAddressSvc.GetEmailAddresses).Methods("GET")
	a.Router.HandleFunc("/emailaddress/{id:[0-9]+}", a.emailAddressSvc.GetEmailAddress).Methods("GET")
	a.Router.HandleFunc("/emailaddress", a.emailAddressSvc.CreateEmailAddress).Methods("POST")
	a.Router.HandleFunc("/emailaddress/{id:[0-9]+}", a.emailAddressSvc.UpdateEmailAddress).Methods("PUT")
	a.Router.HandleFunc("/emailaddress/{id:[0-9]+}", a.emailAddressSvc.DeleteEmailAddress).Methods("DELETE")
}

// ServeRest runs the server
func (a *App) ServeRest(addr, origin string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{origin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "DELETE", "POST", "PUT", "OPTIONS"})
	log.Printf("Allowing origin: %s\n", origin)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(originsOk, handlers.AllowCredentials(), headersOk, methodsOk)(a.Router)))
}
