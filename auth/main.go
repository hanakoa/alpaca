package main

import (
	"database/sql"
	"fmt"
	"github.com/TeslaGov/envy"
	"github.com/hanakoa/alpaca/auth/grpc"
	"github.com/hanakoa/alpaca/auth/models"
	mfaGRPC "github.com/hanakoa/alpaca/mfa/grpc"
	"github.com/kevinmichaelchen/my-go-utils"
	"log"
	"sync"
	"time"
)

func main() {
	user := envy.EnvOrString("DB_USER", "alpaca")
	pass := envy.MustEnv("DB_PASSWORD")
	host := envy.MustEnv("DB_HOST")
	rabbitmqEnabled := envy.EnvOrBool("RABBITMQ_ENABLED", false)
	dbName := envy.EnvOrString("DB_DATABASE", "alpaca_auth")
	secret := envy.MustEnv("ALPACA_SECRET")
	origin := envy.MustEnv("ORIGIN_ALLOWED")
	port := envy.EnvOrInt("APP_PORT", 8080)
	maxWorkers := envy.EnvOrInt("MAX_WORKERS", 1)
	grpcPort := envy.EnvOrInt("GRPC_PORT", 50051)
	grpcMFAHost := envy.EnvOrString("GRPC_MFA_API_HOST", "localhost")
	grpcMFAPort := envy.EnvOrInt("GRPC_MFA_API_PORT", 50052)

	//snowFlakeNodeNumber := utils.StringToInt64(utils.MustEnv("SNOWFLAKE_NODE_NUMBER"))
	// TODO node number should come from env var
	var snowflakeNodeNumber int64 = 1

	db := InitDB(user, pass, host, dbName)

	var wg sync.WaitGroup

	snowflakeNode := utils.InitSnowflakeNode(snowflakeNodeNumber)
	// TODO configurable duration?
	iterationCount := models.CalibrateIterationCount(time.Millisecond * 1000)

	a := App{RabbitmqEnabled: rabbitmqEnabled, iterationCount: iterationCount}
	a.mfaClient = mfaGRPC.NewMFAClient(grpcMFAHost, grpcMFAPort)
	a.snowflakeNode = snowflakeNode
	a.Initialize(db, secret, maxWorkers)
	log.Printf("Running on port %d...\n", port)
	wg.Add(1)
	go a.ServeRest(fmt.Sprintf(":%d", port), origin)

	wg.Add(1)
	grpcServer := &grpc.GrpcServer{Port: grpcPort, DB: db, PasswordService: a.passwordService}
	go grpcServer.Run()

	wg.Wait()
}

func InitDB(user, password, host, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	return utils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
