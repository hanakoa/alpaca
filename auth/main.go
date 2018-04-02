package main

import (
	"database/sql"
	"fmt"
	"github.com/hanakoa/alpaca/auth/grpc"
	"github.com/hanakoa/alpaca/auth/models"
	"github.com/hanakoa/alpaca/auth/services"
	mfaGRPC "github.com/hanakoa/alpaca/mfa/grpc"
	"github.com/kevinmichaelchen/my-go-utils"
	"log"
	"sync"
	"time"
)

func main() {
	user := utils.EnvOrString("DB_USER", "alpaca")
	pass := utils.MustEnv("DB_PASSWORD")
	host := utils.MustEnv("DB_HOST")
	dbName := utils.EnvOrString("DB_DATABASE", "alpaca_auth")
	secret := utils.MustEnv("ALPACA_SECRET")
	origin := utils.MustEnv("ORIGIN_ALLOWED")
	port := utils.EnvOrInt("APP_PORT", 8080)
	maxWorkers := utils.EnvOrInt("MAX_WORKERS", 1)
	grpcPort := utils.EnvOrInt("GRPC_PORT", 50051)
	grpcMFAHost := utils.EnvOrString("GRPC_MFA_API_HOST", "localhost")
	grpcMFAPort := utils.EnvOrInt("GRPC_MFA_API_PORT", 50052)

	//snowFlakeNodeNumber := utils.StringToInt64(utils.MustEnv("SNOWFLAKE_NODE_NUMBER"))
	// TODO node number should come from env var
	var snowflakeNodeNumber int64 = 1

	db := InitDB(user, pass, host, dbName)

	var wg sync.WaitGroup

	snowflakeNode := utils.InitSnowflakeNode(snowflakeNodeNumber)
	// TODO configurable duration?
	iterationCount := models.CalibrateIterationCount(time.Millisecond * 1000)
	passwordService := services.PasswordService{
		DB:             db,
		SnowflakeNode:  snowflakeNode,
		IterationCount: iterationCount}

	a := App{}
	a.mfaClient = mfaGRPC.NewMFAClient(grpcMFAHost, grpcMFAPort)
	a.snowflakeNode = snowflakeNode
	a.Initialize(db, secret, maxWorkers)
	a.passwordService = passwordService
	log.Printf("Running on port %d...\n", port)
	wg.Add(1)
	go a.ServeRest(fmt.Sprintf(":%d", port), origin)

	wg.Add(1)
	grpcServer := &grpc.GrpcServer{Port: grpcPort, DB: db, PasswordService: passwordService}
	go grpcServer.Run()

	wg.Wait()
}

func InitDB(user, password, host, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	return utils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
