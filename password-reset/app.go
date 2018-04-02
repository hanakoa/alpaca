package main

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	authGRPC "github.com/hanakoa/alpaca/auth/grpc"
	"github.com/hanakoa/alpaca/password-reset/services"
	"github.com/kevinmichaelchen/my-go-utils"
	"github.com/kevinmichaelchen/my-go-utils/rabbitmq"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type App struct {
	Router        *mux.Router
	DB            *sql.DB
	snowflakeNode *snowflake.Node
	svc           services.PasswordResetSvc
}

// Initialize initializes the database connection and services,
// and creates routes for our REST API.
func (a *App) Initialize(
	db *sql.DB,
	snowflakeNodeNumber int64,
	numWorkers int,
	passClient authGRPC.PassClient) {
	rabbitmq.NewDispatcher(numWorkers, 10)

	a.snowflakeNode = utils.InitSnowflakeNode(snowflakeNodeNumber)
	a.DB = db
	a.svc = services.PasswordResetSvc{DB: a.DB, SnowflakeNode: a.snowflakeNode, PassClient: passClient}
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router = mux.NewRouter()

	// Send code to requestBody.emailAddress
	a.Router.HandleFunc("/password-reset", a.svc.SendCode).Methods("POST")
	a.Router.HandleFunc("/password-reset/{code:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}", a.svc.VerifyCode).Methods("GET")
	a.Router.HandleFunc("/password-reset", a.svc.ResetPassword).Methods("PUT")
}

// ServeRest runs the server
func (a *App) ServeRest(addr, origin string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{origin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "DELETE", "POST", "PUT", "OPTIONS"})
	log.Printf("Allowing origin: %s\n", origin)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)))
}

// ListenForRabbitMqEvents listens for events
func (a *App) ListenForRabbitMqEvents() {
	l := rabbitmq.NewRabbitListener("alpaca-auth-exchange", "person.#", "alpaca-password-reset-queue")
	l.Listen()
}
