package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"testing"
)

func TestGetUsers(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch accounts", func() {
			ids := AddUsers(20)

			Convey("The response should be", func() {
				log.Println("First ID is", ids[0].Id)
				endpoint := fmt.Sprintf("/account?count=3&cursor=%d", ids[0].Id)
				req, _ := http.NewRequest("GET", endpoint, nil)
				response := ExecuteRequest(req)
				log.Println(response)
				So(response.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestEmptyTable(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch accounts", func() {
			req, _ := http.NewRequest("GET", "/account", nil)
			response := ExecuteRequest(req)

			Convey("The response should be", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
				body := response.Body.String()
				So(body, ShouldEqual, `{"data":[],"next_cursor":-1,"next_cursor_str":"-1","previous_cursor":-1,"previous_cursor_str":"-1"}`)
			})
		})
	})
}

func TestGetNonExistentUser(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch a non-existent account", func() {
			req, _ := http.NewRequest("GET", "/account/45", nil)
			response := ExecuteRequest(req)

			Convey("The response should be", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				var m map[string]string
				json.Unmarshal(response.Body.Bytes(), &m)
				So(m["error"], ShouldEqual, "Account not found")
			})
		})
	})
}

func TestCreateUser(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we create a account without an email address", func() {
			payload := []byte(`{"username":"kevinmchen"}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Must supply email address."}`)
			})
		})

		Convey("When we create a account with an empty email address", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":""}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Must supply email address."}`)
			})
		})

		Convey("When we create a account with a malformed email address", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":"kevin"}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Email address has invalid format."}`)
			})
		})

		Convey("When we create a account with an excessively long email address", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":"contact-admin-hello-webmaster-info-services-peter-crazy-but-oh-so-ubber-cool-english-alphabet-loverer-abcdefghijklmnopqrstuvwxyz@please-try-to.send-me-an-email-if-you-can-possibly-begin-to-remember-this-coz.this-is-the-longest-email-address-known-to-man-but-to-be-honest.this-is-such-a-stupidly-long-sub-domain-it-could-go-on-forever.pacraig.com"}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Email address cannot exceed 255 chars."}`)
			})
		})

		Convey("When we create a account with a non-null, empty username", func() {
			payload := []byte(`{"username":"","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Username must be non-empty."}`)
			})
		})

		Convey("The response should indicate failure", func() {
			Convey("When we submit a username with a dash", func() {
				payload := []byte(`{"username":"kevin-chen","email_address":"kevin.chen.bulk@gmail.com"}`)
				req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
				response := ExecuteRequest(req)
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"A username can only contain alphanumeric characters (letters A-Z, numbers 0-9) with the exception of underscores."}`)
			})
		})

		Convey("When we create a account with a username that is too short", func() {
			payload := []byte(`{"username":"kev","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Username length must be between 4 and 25."}`)
			})
		})

		Convey("When we create a account with a username that is too long", func() {
			payload := []byte(`{"username":"hanakoahanakoahanakoahanakoa","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Username length must be between 4 and 25."}`)
			})
		})

		Convey("When we create a account with the required fields", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate success", func() {
				So(response.Code, ShouldEqual, http.StatusCreated)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(m["error"], ShouldBeNil)
				So(m["id"], ShouldNotEqual, 0.0)
				So(m["id_str"], ShouldNotBeEmpty)
				So(GetInt64(m, "id_str"), ShouldNotEqual, "0")
				So(m["primary_email_address_id_str"], ShouldNotEqual, "0")
				So(m["primary_email_address_id_str"], ShouldNotEqual, 0)
				So(m["primary_email_address_id_str"], ShouldNotBeEmpty)
				So(m["username"], ShouldEqual, "kevinmchen")
			})
		})
	})
}

func TestGetUser(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch a account", func() {
			ids := AddUsers(1)

			Convey("The response should be", func() {
				req, _ := http.NewRequest("GET", fmt.Sprintf("/account/%d", ids[0].Id), nil)
				response := ExecuteRequest(req)
				So(response.Code, ShouldEqual, http.StatusOK)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(GetInt64(m, "id_str"), ShouldEqual, ids[0].Id)
				So(m["primary_email_address_id_str"], ShouldNotEqual, "0")
				So(m["primary_email_address_id_str"], ShouldNotEqual, 0)
				So(m["primary_email_address_id_str"], ShouldNotBeEmpty)
			})
		})
	})
}

func TestUpdateUser(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we update a account", func() {
			ids := AddUsers(1)
			endpoint := fmt.Sprintf("/account/%d", ids[0].Id)

			req, _ := http.NewRequest("GET", endpoint, nil)
			response := ExecuteRequest(req)

			var originalUser map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &originalUser)

			Convey("without a primary email address ID the response should indicate failure", func() {
				payload := []byte(`{"username":"hodor"}`)
				req, _ = http.NewRequest("PUT", endpoint, bytes.NewBuffer(payload))
				response = ExecuteRequest(req)
				log.Println(response.Body)

				So(response.Code, ShouldEqual, http.StatusBadRequest)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(m["error"], ShouldEqual, "Must provide primary email address ID")
			})

			Convey("with a primary email address ID the response should indicate success", func() {
				payload := []byte(`{"username":"hodor","primary_email_address_id":` + GetString(originalUser, "primary_email_address_id_str") + `}`)
				req, _ = http.NewRequest("PUT", endpoint, bytes.NewBuffer(payload))
				response = ExecuteRequest(req)
				log.Println(response.Body)

				So(response.Code, ShouldEqual, http.StatusOK)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(m["id"], ShouldEqual, originalUser["id"])
				So(GetInt64(m, "id_str"), ShouldEqual, ids[0].Id)
				So(m["username"], ShouldEqual, "hodor")
				So(m["username"], ShouldNotEqual, originalUser["username"])
				So(m["primary_email_address_id_str"], ShouldNotEqual, "0")
				So(m["primary_email_address_id_str"], ShouldNotEqual, 0)
				So(m["primary_email_address_id_str"], ShouldNotBeEmpty)
			})
		})
	})
}

func TestDeleteUser(t *testing.T) {
	Convey("Given a account", t, func() {
		ClearTable()
		ids := AddUsers(1)
		endpoint := fmt.Sprintf("/account/%d", ids[0].Id)

		req, _ := http.NewRequest("GET", endpoint, nil)
		response := ExecuteRequest(req)
		So(response.Code, ShouldEqual, http.StatusOK)

		Convey("When we delete a account", func() {
			req, _ = http.NewRequest("DELETE", endpoint, nil)
			response = ExecuteRequest(req)
			log.Println(response.Body)
			So(response.Code, ShouldEqual, http.StatusOK)

			Convey("GET should return Not Found", func() {
				req, _ = http.NewRequest("GET", endpoint, nil)
				response = ExecuteRequest(req)
				So(response.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}
