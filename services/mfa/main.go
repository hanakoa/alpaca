package main

import (
	"database/sql"
	"fmt"
	"github.com/TeslaGov/envy"
	"github.com/hanakoa/alpaca/services/mfa/grpc"
	"github.com/hanakoa/alpaca/services/mfa/services"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"github.com/sfreiberg/gotwilio"
	"log"
	"sync"
	"time"
)

func main() {
	a := App{}

	twilioAccountSid := envy.String("TWILIO_ACCOUNT_SID")
	twilioAuthToken := envy.String("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := envy.String("TWILIO_PHONE_NUMBER")
	user := envy.StringOr("DB_USER", "alpaca")
	pass := envy.String("DB_PASSWORD")
	host := envy.String("DB_HOST")
	dbName := envy.StringOr("DB_DATABASE", "alpaca_mfa")
	origin := envy.String("ORIGIN_ALLOWED")
	port := envy.IntOr("APP_PORT", 8082)
	maxWorkers := envy.IntOr("MAX_WORKERS", 1)
	grpcPort := envy.IntOr("GRPC_PORT", 50052)

	twilio := gotwilio.NewTwilioClient(twilioAccountSid, twilioAuthToken)
	twilioService := services.TwilioSvc{
		Twilio:            twilio,
		TwilioAccountSid:  twilioAccountSid,
		TwilioAuthToken:   twilioAuthToken,
		TwilioPhoneNumber: twilioPhoneNumber}

	var snowflakeNodeNumber int64 = 1

	db := initDB(user, pass, host, dbName)

	var wg sync.WaitGroup

	a.TwilioService = twilioService
	a.Initialize(db, snowflakeNodeNumber, maxWorkers)
	log.Printf("Running on port %d...\n", port)
	wg.Add(1)
	go a.ServeRest(fmt.Sprintf(":%d", port), origin)

	wg.Add(1)
	grpcServer := &grpc.GrpcServer{
		Port:          grpcPort,
		DB:            db,
		TwilioService: twilioService}
	go grpcServer.Run()

	wg.Wait()
}

func initDB(user, password, host, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	return sqlUtils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
