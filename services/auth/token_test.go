package main

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hanakoa/alpaca/services/auth/services"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"strings"
	"testing"
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
			cookies := response.HeaderMap["Set-Cookie"]
			So(response.Code, ShouldEqual, http.StatusOK)
			So(cookies, ShouldNotBeNil)
			So(len(cookies), ShouldEqual, 1)
			jwtString := parseJWT(cookies)["alpacajwt"]
			log.Println("Parsing jwt", jwtString)
			alpacaClaims, err := parseAlpacaClaims(jwtString)
			So(err, ShouldBeNil)
			So(alpacaClaims, ShouldNotBeNil)
			So(alpacaClaims.Issuer, ShouldEqual, "alpaca")
			So(alpacaClaims.Subject, ShouldNotBeNil)
		})
	})
}

// parseAlpacaClaims parses claims from a JWT string
func parseAlpacaClaims(jwtString string) (*services.AlpacaClaims, error) {
	resource := &services.AlpacaClaims{}
	// JWT validation happens at the API Gateway.
	if _, _, err := new(jwt.Parser).ParseUnverified(jwtString, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func parseJWT(cookieStrings []string) map[string]string {
	m := make(map[string]string)
	for _, c := range cookieStrings {
		eq := strings.Split(c, "=")
		key := eq[0]
		val := strings.Split(eq[1], ";")[0]
		m[key] = val
	}
	return m
}
