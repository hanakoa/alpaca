package main

import (
	"database/sql"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hanakoa/alpaca/mfa/services"
	"github.com/kevinmichaelchen/my-go-utils/rabbitmq"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type App struct {
	Router        *mux.Router
	DB            *sql.DB
	snowflakeNode *snowflake.Node
	TwilioService services.TwilioSvc
}

// Initialize initializes the database connection and services,
// and creates routes for our REST API.
func (a *App) Initialize(db *sql.DB, snowflakeNodeNumber int64, numWorkers int) {
	rabbitmq.NewDispatcher(numWorkers, 10)
	a.snowflakeNode = snowflakeUtils.InitSnowflakeNode(snowflakeNodeNumber)
	a.DB = db
	a.initializeServices()
	a.initializeRoutes()
}

// initializeServices initializes the Service layer.
func (a *App) initializeServices() {
	// TODO does MFA have any services?
}

func (a *App) initializeRoutes() {
	a.Router = mux.NewRouter()
}

// ServeRest runs the server
func (a *App) ServeRest(addr, origin string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{origin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "DELETE", "POST", "PUT", "OPTIONS"})
	log.Printf("Allowing origin: %s\n", origin)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(originsOk, handlers.AllowCredentials(), headersOk, methodsOk)(a.Router)))
}
