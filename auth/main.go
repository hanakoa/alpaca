package main

import (
	"database/sql"
	"fmt"
	"github.com/TeslaGov/envy"
	"github.com/hanakoa/alpaca/auth/grpc"
	"github.com/hanakoa/alpaca/auth/models"
	mfaGRPC "github.com/hanakoa/alpaca/mfa/grpc"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"log"
	"sync"
	"time"
)

func main() {
	user := envy.StringOr("DB_USER", "alpaca")
	pass := envy.String("DB_PASSWORD")
	host := envy.String("DB_HOST")
	rabbitmqEnabled := envy.BoolOr("RABBITMQ_ENABLED", false)
	dbName := envy.StringOr("DB_DATABASE", "alpaca_auth")
	secret := envy.String("ALPACA_SECRET")
	origin := envy.String("ORIGIN_ALLOWED")
	port := envy.IntOr("APP_PORT", 8080)
	maxWorkers := envy.IntOr("MAX_WORKERS", 1)
	grpcPort := envy.IntOr("GRPC_PORT", 50051)
	grpcMFAHost := envy.StringOr("GRPC_MFA_API_HOST", "localhost")
	grpcMFAPort := envy.IntOr("GRPC_MFA_API_PORT", 50052)

	//snowFlakeNodeNumber := utils.StringToInt64(utils.String("SNOWFLAKE_NODE_NUMBER"))
	// TODO node number should come from env var
	var snowflakeNodeNumber int64 = 1

	db := InitDB(user, pass, host, dbName)

	var wg sync.WaitGroup

	snowflakeNode := snowflakeUtils.InitSnowflakeNode(snowflakeNodeNumber)
	iterationCount := envy.IntOr("ITERATION_COUNT", 0)
	if iterationCount == 0 {
		iterationCount = models.CalibrateIterationCount(time.Millisecond * 1000)
	}
	log.Println("Using iteration count:", iterationCount)

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
	return sqlUtils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
