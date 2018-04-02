package main

import (
	"database/sql"
	"fmt"
	"github.com/hanakoa/alpaca/mfa/grpc"
	"github.com/hanakoa/alpaca/mfa/services"
	"github.com/kevinmichaelchen/my-go-utils"
	"github.com/sfreiberg/gotwilio"
	"log"
	"sync"
	"time"
)

func main() {
	a := App{}

	twilioAccountSid := utils.MustEnv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := utils.MustEnv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := utils.MustEnv("TWILIO_PHONE_NUMBER")
	user := utils.EnvOrString("DB_USER", "alpaca")
	pass := utils.MustEnv("DB_PASSWORD")
	host := utils.MustEnv("DB_HOST")
	dbName := utils.EnvOrString("DB_DATABASE", "alpaca_mfa")
	origin := utils.MustEnv("ORIGIN_ALLOWED")
	port := utils.EnvOrInt("APP_PORT", 8082)
	maxWorkers := utils.EnvOrInt("MAX_WORKERS", 1)
	grpcPort := utils.EnvOrInt("GRPC_PORT", 50052)

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
	return utils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
