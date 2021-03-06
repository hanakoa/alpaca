package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/TeslaGov/envy"
	"github.com/hanakoa/alpaca/services/auth/models"
	"github.com/hanakoa/alpaca/services/auth/services"
	snowflakeUtils "github.com/kevinmichaelchen/my-go-utils/snowflake"
	stringUtils "github.com/kevinmichaelchen/my-go-utils/string"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/guregu/null.v3"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var MyApp App
var DB *sql.DB

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	MyApp = App{
		RabbitmqEnabled: false,
		iterationCount:  10,
		CookieConfig: &services.CookieConfiguration{
			Name:     "alpacajwt",
			Domain:   "",
			HttpOnly: true,
		},
		JWTSecret: "4FFFA6A10E744158464EB55133A475673264748804882A1B4F8106D545C584EF",
	}
	dbHost := envy.StringOr("AUTH_DB_HOST", "localhost")
	DB = InitDB("alpaca", "password", dbHost, "alpaca_auth_test")
	MyApp.snowflakeNode = snowflakeUtils.InitSnowflakeNode(1)
	MyApp.Initialize(DB, 1)

	code := m.Run()

	os.Exit(code)
}

func GetString(m map[string]interface{}, key string) string {
	return m[key].(string)
}

func GetInt64(m map[string]interface{}, key string) int64 {
	return stringUtils.StringToInt64(GetString(m, key))
}

func ClearTable() {
	MyApp.DB.Exec("UPDATE account SET primary_email_address_id = NULL")
	MyApp.DB.Exec("UPDATE account SET current_password_id = NULL")
	MyApp.DB.Exec("UPDATE email_address SET account_id = NULL")
	MyApp.DB.Exec("UPDATE phone_number SET account_id = NULL")
	MyApp.DB.Exec("DELETE FROM email_address")
	MyApp.DB.Exec("DELETE FROM login_attempt")
	MyApp.DB.Exec("DELETE FROM password")
	MyApp.DB.Exec("DELETE FROM account")
	MyApp.DB.Exec("DELETE FROM phone_number")
}

func WithTransaction(db *sql.DB, f func(tx *sql.Tx)) func() {
	return func() {
		tx, err := db.Begin()
		So(err, ShouldBeNil)

		Reset(func() {
			/* Verify that the transaction is alive by executing a command */
			_, err := tx.Exec("SELECT 1")
			So(err, ShouldBeNil)

			tx.Rollback()
		})

		/* Here we invoke the actual test-closure and provide the transaction */
		f(tx)
	}
}

func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	MyApp.Router.ServeHTTP(rr, req)

	return rr
}

func AddUsers(count int) []models.Account {
	if count < 1 {
		count = 1
	}

	users := []models.Account{}
	for i := 1; i <= count; i++ {
		user := &services.CreateAccountRequest{
			Username:     null.StringFrom(fmt.Sprintf("user%d", i)),
			EmailAddress: fmt.Sprintf("user%d@gmail.com", i)}
		b, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", "/account", bytes.NewBuffer(b))
		if err != nil {
			panic(err)
		}
		res := ExecuteRequest(req)
		var p models.Account
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&p); err != nil {
			panic(err)
		}
		if p.Id == 0 {
			panic(fmt.Errorf("POST /account failed: %s", res.Body))
		}
		users = append(users, p)
	}
	return users
}
