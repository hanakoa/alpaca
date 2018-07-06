package main

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
	"log"
)

func TestLogin(t *testing.T) {
	ClearTable()

	Convey("Given a user with an email, username, and password", t, func() {

		users := AddUsers(1)

		payload := []byte(`{"password":"potato-tomato-cherry-gun"}`)
		req, err := http.NewRequest("PUT", fmt.Sprintf("/account/%d/password", users[0].Id), bytes.NewBuffer(payload))
		So(err, ShouldBeNil)
		response := ExecuteRequest(req)
		So(response.Code, ShouldEqual, http.StatusOK)

		Convey("Should return authenticated", func() {
			payload := []byte(`{"login":"user1","password":"potato-tomato-cherry-gun"}`)
			req, err := http.NewRequest("POST", "/token", bytes.NewBuffer(payload))
			So(err, ShouldBeNil)
			response := ExecuteRequest(req)
			log.Println("Set-Cookie", response.HeaderMap["Set-Cookie"])
			log.Println("response", response.Body.String())
			So(response.Code, ShouldEqual, http.StatusOK)
		})
	})
}
