package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"database/sql"
	"encoding/json"
	"bytes"
	"github.com/hanakoa/alpaca/password-reset/services"
)

var a App
var db *sql.DB
func TestMain(m *testing.M) {
	a = App{}

	db = InitDB("alpaca", "password", "localhost", "alpaca_password_reset_test")
	ClearTable(db)
	a.Initialize(db, 1, 1, nil)

	code := m.Run()

	os.Exit(code)
}

func TestSendPasswordReset(t *testing.T) {
	if tx, err := db.Begin(); err != nil {
		panic(err)
	} else {
		tx.Exec(`INSERT INTO person (id, username) VALUES (1, 'kevin_chen')`)
		tx.Exec(`INSERT INTO email_address (id, person_id, email_address, confirmed, is_primary) VALUES (1, 1, 'kevin.chen.bulk@gmail.com', TRUE, TRUE)`)
		if err := tx.Commit(); err != nil {
			panic(err)
		}
	}

	Convey("Given a user with an email and username", t, WithTransaction(db, func(tx *sql.Tx) {
		Convey("Should return no options and state that an email has been sent", func() {
			So(sendCodeRequest("kevin_chen"), ShouldEqual, `"Reset request submitted."`)
			So(sendCodeRequest("kevin.chen.bulk@gmail.com"), ShouldEqual, `"Reset request submitted."`)
		})
	}))
}

func TestSendPasswordResetPhoneNumber(t *testing.T) {
	if tx, err := db.Begin(); err != nil {
		panic(err)
	} else {
		tx.Exec(`INSERT INTO person (id, username) VALUES (1, 'kevin_chen')`)
		tx.Exec(`INSERT INTO email_address (id, person_id, email_address, confirmed, is_primary) VALUES (1, 1, 'kevin.chen.bulk@gmail.com', TRUE, TRUE)`)
		tx.Exec(`INSERT INTO phone_number (id, person_id, phone_number, confirmed) VALUES (1, 1, '5555555555', TRUE)`)
		if err := tx.Commit(); err != nil {
			panic(err)
		}
	}

	Convey("Given a user with an email, username, and phone number", t, WithTransaction(db, func(tx *sql.Tx) {
		Convey("Options should include emails and phone", func() {
			res := sendCodeRequest("kevin_chen")
			So(res, ShouldEqual, `{"phone_numbers":[{"id":1,"phone_number":"55","person_id":1}],"email_addresses":[{"id":1,"email_address":"ke*************@g****.com","person_id":1}]}`)

			res = sendCodeRequest("kevin.chen.bulk@gmail.com")
			So(res, ShouldEqual, `{"phone_numbers":[{"id":1,"phone_number":"55","person_id":1}],"email_addresses":[{"id":1,"email_address":"ke*************@g****.com","person_id":1}]}`)

			res = sendCodeRequest("5555555555")
			So(res, ShouldEqual, `{"phone_numbers":[{"id":1,"phone_number":"55","person_id":1}],"email_addresses":[{"id":1,"email_address":"ke*************@g****.com","person_id":1}]}`)

			res = sendCodeRequest("5554444444")
			So(res, ShouldEqual, `"Reset request submitted."`)
		})
	}))
}

func sendCodeRequest(account string) string {
	sendCodeRequest := &services.SendCodeRequest{Account: account}
	b, err := json.Marshal(sendCodeRequest)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "/password-reset", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
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
	db.Exec("DELETE FROM person")
}

func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
