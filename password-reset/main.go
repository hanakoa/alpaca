package main

import (
	"database/sql"
	"fmt"
	authGRPC "github.com/hanakoa/alpaca/auth/grpc"
	"github.com/kevinmichaelchen/my-go-utils"
	"log"
	"sync"
	"time"
	"github.com/TeslaGov/envy"
)

func main() {
	a := App{}

	user := envy.EnvOrString("DB_USER", "alpaca")
	pass := envy.MustEnv("DB_PASSWORD")
	host := envy.MustEnv("DB_HOST")
	dbName := envy.EnvOrString("DB_DATABASE", "alpaca_password_reset")
	origin := envy.MustEnv("ORIGIN_ALLOWED")
	port := envy.EnvOrInt("APP_PORT", 8081)
	maxWorkers := envy.EnvOrInt("MAX_WORKERS", 1)

	grpcAuthHost := envy.EnvOrString("GRPC_AUTH_API_HOST", "localhost")
	grpcAuthPort := envy.EnvOrInt("GRPC_AUTH_API_PORT", 50051)

	//snowFlakeNodeNumber := utils.StringToInt64(utils.MustEnv("SNOWFLAKE_NODE_NUMBER"))
	// TODO node number should come from env var
	var snowflakeNodeNumber int64 = 1

	var wg sync.WaitGroup

	db := InitDB(user, pass, host, dbName)
	a.Initialize(db, snowflakeNodeNumber, maxWorkers, authGRPC.NewPassClient(grpcAuthHost, grpcAuthPort))
	log.Printf("Running on port %d...\n", port)
	wg.Add(1)
	go a.ServeRest(fmt.Sprintf(":%d", port), origin)

	wg.Add(1)
	go a.ListenForRabbitMqEvents()

	wg.Wait()
}

func InitDB(user, password, host, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	return utils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
