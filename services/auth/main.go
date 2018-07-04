package main

import (
	"database/sql"
	"fmt"
	"github.com/TeslaGov/envy"
	"github.com/bwmarrin/snowflake"
	"github.com/hanakoa/alpaca/services/auth/grpc"
	"github.com/hanakoa/alpaca/services/auth/models"
	mfaGRPC "github.com/hanakoa/alpaca/services/mfa/grpc"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	sqlUtils "github.com/kevinmichaelchen/my-go-utils/sql"
	"gopkg.in/guregu/null.v3"
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

// TODO will also need to emit a message event so other microservices can have this new Account and EmailAddress
func seedAdminAccount(db *sql.DB, snowflakeNode *snowflake.Node) {
	username := envy.StringOr("SEED_ACCOUNT_USERNAME", "")
	email := envy.StringOr("SEED_ACCOUNT_EMAIL", "")

	if count, err := models.TotalAccountCount(db); err != nil {
		log.Fatal("Could not reach DB while doing a count(*)")
	} else if count > 0 {
		log.Println("Account table is not empty. Skipping...")
		return
	} else if count == 0 && username == "" && email == "" {
		log.Fatal("Account table is empty, so you must supply env vars to seed admin account")
	}

	var acc1, acc2 *models.Account
	var err error
	acc1, err = models.GetAccountByEmailAddress(db, email)
	if err != nil {
		log.Fatal("Could not contact DB during account seeding, while trying " +
			"to check for an existing account by email address")
	}
	acc2, err = models.GetAccountByUsername(db, username)
	if err != nil {
		log.Fatal("Could not contact DB during account seeding, while trying " +
			"to check for an existing account by username")
	}

	if acc1 != nil && acc2 == nil {
		log.Fatal("An account already exists with that email, but not that username")
	}
	if acc2 != nil && acc1 == nil {
		log.Fatal("An account already exists with that username, but not that email")
	}
	if acc1 != nil && acc2 != nil {
		if acc1.Id == acc2.Id {
			log.Println("A seed account already exists with that email and username. Skipping...")
			return
		} else {
			log.Fatal("Two existing accounts already exist, one has the username, the other has the email")
		}
	}

	// TODO validate email and username
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Could not start transaction.")
	}

	accountID := snowflakeUtils.NewPrimaryKey(snowflakeNode)
	a := &models.Account{
		Id:       accountID,
		Username: null.StringFrom(username),
	}
	if err := a.CreateAccount(tx); err != nil {
		tx.Rollback()
		log.Fatal("Could not create account")
	}

	// TODO add admin role
	emailAddressID := snowflakeUtils.NewPrimaryKey(snowflakeNode)
	emailAddress := &models.EmailAddress{
		ID:           emailAddressID,
		Primary:      true,
		EmailAddress: email,
		AccountID:    accountID,
	}
	if err := emailAddress.CreateEmailAddress(tx); err != nil {
		tx.Rollback()
		log.Fatal("Could not create email address")
	}
	log.Printf("Seeded user %s - %s", username, email)
}

func InitDB(user, password, host, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	return sqlUtils.InitDatabase("postgres", connectionString, 3, time.Second*5)
}
