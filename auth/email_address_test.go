package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func TestGetEmailAddresses(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch people with an empty database", func() {
			req, _ := http.NewRequest("GET", "/emailaddress", nil)
			response := ExecuteRequest(req)

			Convey("The response should be", func() {
				So(response.Code, ShouldEqual, http.StatusOK)
				body := response.Body.String()
				So(body, ShouldEqual, `{"data":[],"next_cursor":-1,"next_cursor_str":"-1","previous_cursor":-1,"previous_cursor_str":"-1"}`)
			})
		})

		Convey("When we fetch people", func() {
			ids := AddUsers(20)

			Convey("The response should be", func() {
				log.Println("First ID is", ids[0])
				endpoint := fmt.Sprintf("/emailaddress?count=3&cursor=%d", ids[0])
				req, _ := http.NewRequest("GET", endpoint, nil)
				response := ExecuteRequest(req)
				log.Println(response)
				So(response.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestGetNonExistentEmailAddress(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch a non-existent email address", func() {
			req, _ := http.NewRequest("GET", "/emailaddress/45", nil)
			response := ExecuteRequest(req)

			Convey("The response should be", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				var m map[string]string
				json.Unmarshal(response.Body.Bytes(), &m)
				So(m["error"], ShouldEqual, "EmailAddress not found")
			})
		})
	})
}

func TestCreateEmailAddress(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we try to create an email address with an id", func() {
			payload := []byte(`{"id": 5, "person_id": 1, "email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/emailaddress", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Do not provide an id."}`)
			})
		})

		Convey("When we try to create an email address without an email address", func() {
			payload := []byte(`{"person_id": 1}`)
			req, _ := http.NewRequest("POST", "/emailaddress", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Must supply email address."}`)
			})
		})

		Convey("When we try to create an email address with an empty email address", func() {
			payload := []byte(`{"person_id": 1, "email_address":" "}`)
			req, _ := http.NewRequest("POST", "/emailaddress", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Must supply email address."}`)
			})
		})

		Convey("When we try to create an email address with a malformed email address", func() {
			payload := []byte(`{"person_id":1,"email_address":"kevin"}`)
			req, _ := http.NewRequest("POST", "/emailaddress", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Email address has invalid format."}`)
			})
		})

		Convey("When we try to create an email address with an excessively long email address", func() {
			payload := []byte(`{"person_id":1,"email_address":"contact-admin-hello-webmaster-info-services-peter-crazy-but-oh-so-ubber-cool-english-alphabet-loverer-abcdefghijklmnopqrstuvwxyz@please-try-to.send-me-an-email-if-you-can-possibly-begin-to-remember-this-coz.this-is-the-longest-email-address-known-to-man-but-to-be-honest.this-is-such-a-stupidly-long-sub-domain-it-could-go-on-forever.pacraig.com"}`)
			req, _ := http.NewRequest("POST", "/emailaddress", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Email address cannot exceed 255 chars."}`)
			})
		})

		Convey("When we try to create an email address that is confirmed", func() {
			payload := []byte(`{"person_id":1,"email_address":"kevin.chen.bulk@gmail.com","confirmed":true}`)
			req, _ := http.NewRequest("POST", "/emailaddress", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Cannot create confirmed email address."}`)
			})
		})

		Convey("When we create an email address with the required fields", func() {
			ids := AddUsers(1)
			personIDString := strconv.FormatInt(ids[0], 10)
			payload := []byte(`{"person_id":` + personIDString + `,"email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/emailaddress", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate success", func() {
				So(response.Code, ShouldEqual, http.StatusCreated)
				var m map[string]interface{}

				json.Unmarshal(response.Body.Bytes(), &m)
				So(m["error"], ShouldBeNil)
				So(m["id"], ShouldNotEqual, 0.0)
				So(m["id_str"], ShouldNotEqual, "")
				So(m["confirmed"], ShouldEqual, false)
				So(m["email_address"], ShouldEqual, "kevin.chen.bulk@gmail.com")
				So(m["person_id_str"], ShouldEqual, personIDString)
				So(strconv.FormatInt(int64(m["person_id"].(float64)), 10), ShouldEqual, personIDString)
			})
		})

	})
}

func TestGetEmailAddress(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch an email address", func() {
			ids := AddUsers(1)

			Convey("The response should be successful", func() {
				req, _ := http.NewRequest("GET", fmt.Sprintf("/emailaddress/%d", ids[0]), nil)
				response := ExecuteRequest(req)
				So(response.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestUpdateEmailAddress(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we attempt to update a non-existent email address", func() {
			payload := []byte(`{"confirmed":true}`)
			req, _ := http.NewRequest("PUT", "/emailaddress/0", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				So(response.Body.String(), ShouldEqual, `{"error":"No email address found for id: 0"}`)
			})
		})

		Convey("When we update an email address", func() {
			ids := AddUsers(1)

			personID := ids[0]
			req, _ := http.NewRequest("GET", fmt.Sprintf("/person/%d", personID), nil)
			response := ExecuteRequest(req)
			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)

			log.Println("Hitting endpoint:", fmt.Sprintf("/person/%d", personID))
			log.Println("Got response:", response.Body)

			primaryEmailAddressID := GetInt64(m, "primary_email_address_id_str")
			log.Println("Primary email address is", primaryEmailAddressID)

			So(m["id_str"], ShouldNotBeNil)
			So(m["id_str"], ShouldNotBeEmpty)
			So(m["primary_email_address_id_str"], ShouldNotBeNil)
			So(m["primary_email_address_id_str"], ShouldNotBeEmpty)
			So(m["id_str"], ShouldNotEqual, m["primary_email_address_id_str"])
			So(GetInt64(m, "id_str"), ShouldEqual, personID)
			So(primaryEmailAddressID, ShouldNotEqual, 0)
			So(primaryEmailAddressID, ShouldNotEqual, personID)

			Convey("The response should indicate success", func() {
				payload := []byte(`{"confirmed":true}`)
				req, _ = http.NewRequest("PUT", fmt.Sprintf("/emailaddress/%d", primaryEmailAddressID), bytes.NewBuffer(payload))
				response = ExecuteRequest(req)
				log.Println(response.Body)

				So(response.Code, ShouldEqual, http.StatusOK)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(GetInt64(m, "id_str"), ShouldEqual, primaryEmailAddressID)
				So(GetInt64(m, "person_id_str"), ShouldEqual, personID)
				So(m["confirmed"], ShouldEqual, true)
			})
		})
	})
}

func TestDeleteEmailAddress(t *testing.T) {
	Convey("Given an email address", t, func() {
		ClearTable()
		ids := AddUsers(1)

		personID := ids[0]
		req, _ := http.NewRequest("GET", fmt.Sprintf("/person/%d", personID), nil)
		response := ExecuteRequest(req)
		log.Println("GOT RESPONSE =", response.Body)
		var m map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &m)

		primaryEmailAddressID := GetInt64(m, "primary_email_address_id_str")
		endpoint := fmt.Sprintf("/emailaddress/%d", primaryEmailAddressID)
		req, _ = http.NewRequest("GET", endpoint, nil)
		response = ExecuteRequest(req)
		log.Println(response.Body)
		So(response.Code, ShouldEqual, http.StatusOK)

		Convey("When we delete a primary email address, it should fail", func() {
			req, _ = http.NewRequest("DELETE", endpoint, nil)
			response = ExecuteRequest(req)
			log.Println(response.Body)
			So(response.Code, ShouldEqual, http.StatusBadRequest)
			body := response.Body.String()
			So(body, ShouldEqual, `{"error":"Email address is primary."}`)
		})
	})
}
