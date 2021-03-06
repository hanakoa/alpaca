package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/TeslaGov/envy"
	"github.com/hanakoa/alpaca/services/password-reset/services"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App
var db *sql.DB

func TestMain(m *testing.M) {
	a = App{}
	dbHost := envy.StringOr("PASSWORD_RESET_DB_HOST", "localhost")
	db = InitDB("alpaca", "password", dbHost, "alpaca_password_reset_test")
	ClearTable(db)
	a.Initialize(db, 1, 1, nil)

	code := m.Run()

	os.Exit(code)
}

func TestSendPasswordReset(t *testing.T) {
	ClearTable(db)
	Convey("Given a user with an email and username", t, WithTransaction(db, func(tx *sql.Tx) {
		tx, err := db.Begin()
		So(err, ShouldBeNil)
		tx.Exec(`INSERT INTO account (id, username) VALUES (1, 'kevin_chen')`)
		tx.Exec(`INSERT INTO email_address (id, account_id, email_address, confirmed, is_primary) VALUES (1, 1, 'kevin.chen.bulk@gmail.com', TRUE, TRUE)`)
		err = tx.Commit()
		So(err, ShouldBeNil)

		Convey("Should return no options and state that an email has been sent", func() {
			So(sendCodeRequest("kevin_chen"), ShouldEqual, `"Reset request submitted."`)
			So(sendCodeRequest("kevin.chen.bulk@gmail.com"), ShouldEqual, `"Reset request submitted."`)
		})
	}))
}

func TestSendPasswordResetPhoneNumber(t *testing.T) {
	ClearTable(db)
	Convey("Given a user with an email, username, and phone number", t, WithTransaction(db, func(tx *sql.Tx) {
		tx, err := db.Begin()
		So(err, ShouldBeNil)
		tx.Exec(`INSERT INTO account (id, username) VALUES (1, 'kevin_chen')`)
		tx.Exec(`INSERT INTO email_address (id, account_id, email_address, confirmed, is_primary) VALUES (1, 1, 'kevin.chen.bulk@gmail.com', TRUE, TRUE)`)
		tx.Exec(`INSERT INTO phone_number (id, account_id, phone_number, confirmed) VALUES (1, 1, '5555555555', TRUE)`)
		err = tx.Commit()
		So(err, ShouldBeNil)

		Convey("Options should include emails and phone", func() {
			res := sendCodeRequest("kevin_chen")
			So(res, ShouldEqual, `{"phone_numbers":[{"id":1,"phone_number":"55","account_id":1}],"email_addresses":[{"id":1,"email_address":"ke*************@g****.com","account_id":1}]}`)

			res = sendCodeRequest("kevin.chen.bulk@gmail.com")
			So(res, ShouldEqual, `{"phone_numbers":[{"id":1,"phone_number":"55","account_id":1}],"email_addresses":[{"id":1,"email_address":"ke*************@g****.com","account_id":1}]}`)

			res = sendCodeRequest("5555555555")
			So(res, ShouldEqual, `{"phone_numbers":[{"id":1,"phone_number":"55","account_id":1}],"email_addresses":[{"id":1,"email_address":"ke*************@g****.com","account_id":1}]}`)

			res = sendCodeRequest("5554444444")
			So(res, ShouldEqual, `"Reset request submitted."`)
		})
	}))
}

func sendCodeRequest(account string) string {
	sendCodeRequest := &services.CodeOptionsRequest{Account: account}
	b, err := json.Marshal(sendCodeRequest)
	So(err, ShouldBeNil)

	req, err := http.NewRequest("POST", "/password-reset", bytes.NewBuffer(b))
	So(err, ShouldBeNil)
	res := ExecuteRequest(req)
	return res.Body.String()
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

func ClearTable(db *sql.DB) {
	db.Exec("DELETE FROM password_reset_code")
	db.Exec("DELETE FROM email_address")
	db.Exec("DELETE FROM phone_number")
	db.Exec("DELETE FROM account")
}

func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
