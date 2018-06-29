package main

import (
	"database/sql"
	"fmt"
	"github.com/TeslaGov/envy"
	authGRPC "github.com/hanakoa/alpaca/auth/grpc"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"log"
	"sync"
	"time"
)

func main() {
	a := App{}

	user := envy.StringOr("DB_USER", "alpaca")
	pass := envy.String("DB_PASSWORD")
	host := envy.String("DB_HOST")
	dbName := envy.StringOr("DB_DATABASE", "alpaca_password_reset")
	origin := envy.String("ORIGIN_ALLOWED")
	port := envy.IntOr("APP_PORT", 8081)
	maxWorkers := envy.IntOr("MAX_WORKERS", 1)

	grpcAuthHost := envy.StringOr("GRPC_AUTH_API_HOST", "localhost")
	grpcAuthPort := envy.IntOr("GRPC_AUTH_API_PORT", 50051)

	//snowFlakeNodeNumber := utils.StringToInt64(utils.String("SNOWFLAKE_NODE_NUMBER"))
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
	return sqlUtils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
