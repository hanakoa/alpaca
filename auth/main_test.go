package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hanakoa/alpaca/auth/models"
	"github.com/kevinmichaelchen/my-go-utils"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/guregu/null.v3"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/hanakoa/alpaca/auth/services"
)

var a App

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a = App{RabbitmqEnabled: false}
	db := InitDB("alpaca", "password", "localhost", "alpaca_auth_test")
	secret := "4FFFA6A10E744158464EB55133A475673264748804882A1B4F8106D545C584EF"
	a.snowflakeNode = utils.InitSnowflakeNode(1)
	a.Initialize(db, secret, 1)

	code := m.Run()

	os.Exit(code)
}

func TestGetUsers(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch people", func() {
			ids := AddUsers(20)

			Convey("The response should be", func() {
				log.Println("First ID is", ids[0])
				endpoint := fmt.Sprintf("/person?count=3&cursor=%d", ids[0])
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

		Convey("When we fetch people", func() {
			req, _ := http.NewRequest("GET", "/person", nil)
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

		Convey("When we fetch a non-existent person", func() {
			req, _ := http.NewRequest("GET", "/person/45", nil)
			response := ExecuteRequest(req)

			Convey("The response should be", func() {
				So(response.Code, ShouldEqual, http.StatusNotFound)
				var m map[string]string
				json.Unmarshal(response.Body.Bytes(), &m)
				So(m["error"], ShouldEqual, "Person not found")
			})
		})
	})
}

func TestCreateUser(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we create a person without an email address", func() {
			payload := []byte(`{"username":"kevinmchen"}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Must supply email address."}`)
			})
		})

		Convey("When we create a person with an empty email address", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":""}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Must supply email address."}`)
			})
		})

		Convey("When we create a person with a malformed email address", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":"kevin"}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Email address has invalid format."}`)
			})
		})

		Convey("When we create a person with an excessively long email address", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":"contact-admin-hello-webmaster-info-services-peter-crazy-but-oh-so-ubber-cool-english-alphabet-loverer-abcdefghijklmnopqrstuvwxyz@please-try-to.send-me-an-email-if-you-can-possibly-begin-to-remember-this-coz.this-is-the-longest-email-address-known-to-man-but-to-be-honest.this-is-such-a-stupidly-long-sub-domain-it-could-go-on-forever.pacraig.com"}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Email address cannot exceed 255 chars."}`)
			})
		})

		Convey("When we create a person with a non-null, empty username", func() {
			payload := []byte(`{"username":"","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Username must be non-empty."}`)
			})
		})

		Convey("The response should indicate failure", func() {
			Convey("When we submit a username with a dash", func() {
				payload := []byte(`{"username":"kevin-chen","email_address":"kevin.chen.bulk@gmail.com"}`)
				req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
				response := ExecuteRequest(req)
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"A username can only contain alphanumeric characters (letters A-Z, numbers 0-9) with the exception of underscores."}`)
			})
		})

		Convey("When we create a person with a username that is too short", func() {
			payload := []byte(`{"username":"kev","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Username length must be between 4 and 25."}`)
			})
		})

		Convey("When we create a person with a username that is too long", func() {
			payload := []byte(`{"username":"hanakoahanakoahanakoahanakoa","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate failure", func() {
				So(response.Code, ShouldEqual, http.StatusBadRequest)
				So(response.Body.String(), ShouldEqual, `{"error":"Username length must be between 4 and 25."}`)
			})
		})

		Convey("When we create a person with the required fields", func() {
			payload := []byte(`{"username":"kevinmchen","email_address":"kevin.chen.bulk@gmail.com"}`)
			req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(payload))
			response := ExecuteRequest(req)

			Convey("The response should indicate success", func() {
				So(response.Code, ShouldEqual, http.StatusCreated)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(m["error"], ShouldBeNil)
				So(m["id"], ShouldNotEqual, 0.0)
				So(m["id_str"], ShouldNotBeEmpty)
				So(GetInt64(m, "id_str"), ShouldNotEqual, "0")
				So(m["primaryEmailAddressId_str"], ShouldNotEqual, "0")
				So(m["primaryEmailAddressId_str"], ShouldNotEqual, 0)
				So(m["primaryEmailAddressId_str"], ShouldNotBeEmpty)
				So(m["username"], ShouldEqual, "kevinmchen")
				So(m["email_address"], ShouldEqual, "kevin.chen.bulk@gmail.com")
			})
		})
	})
}

func TestGetUser(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we fetch a person", func() {
			ids := AddUsers(1)

			Convey("The response should be", func() {
				req, _ := http.NewRequest("GET", fmt.Sprintf("/person/%d", ids[0]), nil)
				response := ExecuteRequest(req)
				So(response.Code, ShouldEqual, http.StatusOK)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(GetInt64(m, "id_str"), ShouldEqual, ids[0])
				So(m["primaryEmailAddressId_str"], ShouldNotEqual, "0")
				So(m["primaryEmailAddressId_str"], ShouldNotEqual, 0)
				So(m["primaryEmailAddressId_str"], ShouldNotBeEmpty)
			})
		})
	})
}

func TestUpdateUser(t *testing.T) {
	Convey("Given an empty database", t, func() {
		ClearTable()

		Convey("When we update a person", func() {
			ids := AddUsers(1)
			endpoint := fmt.Sprintf("/person/%d", ids[0])

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
				payload := []byte(`{"username":"hodor","primaryEmailAddressId":` + GetString(originalUser, "primaryEmailAddressId_str") + `}`)
				req, _ = http.NewRequest("PUT", endpoint, bytes.NewBuffer(payload))
				response = ExecuteRequest(req)
				log.Println(response.Body)

				So(response.Code, ShouldEqual, http.StatusOK)

				var m map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &m)

				So(m["id"], ShouldEqual, originalUser["id"])
				So(GetInt64(m, "id_str"), ShouldEqual, ids[0])
				So(m["username"], ShouldEqual, "hodor")
				So(m["username"], ShouldNotEqual, originalUser["username"])
				So(m["primaryEmailAddressId_str"], ShouldNotEqual, "0")
				So(m["primaryEmailAddressId_str"], ShouldNotEqual, 0)
				So(m["primaryEmailAddressId_str"], ShouldNotBeEmpty)
			})
		})
	})
}

func TestDeleteUser(t *testing.T) {
	Convey("Given a person", t, func() {
		ClearTable()
		ids := AddUsers(1)
		endpoint := fmt.Sprintf("/person/%d", ids[0])

		req, _ := http.NewRequest("GET", endpoint, nil)
		response := ExecuteRequest(req)
		So(response.Code, ShouldEqual, http.StatusOK)

		Convey("When we delete a person", func() {
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

func GetString(m map[string]interface{}, key string) string {
	return m[key].(string)
}

func GetInt64(m map[string]interface{}, key string) int64 {
	return utils.StringToInt64(GetString(m, key))
}

func ClearTable() {
	a.DB.Exec("UPDATE person SET primary_email_address_id = NULL")
	a.DB.Exec("UPDATE email_address SET person_id = NULL")
	a.DB.Exec("UPDATE phone_number SET person_id = NULL")
	a.DB.Exec("DELETE FROM email_address")
	a.DB.Exec("DELETE FROM login_attempt")
	a.DB.Exec("DELETE FROM password")
	a.DB.Exec("DELETE FROM person")
	a.DB.Exec("DELETE FROM phone_number")
}

func ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func AddUsers(count int) []int64 {
	if count < 1 {
		count = 1
	}

	ids := []int64{}
	for i := 1; i <= count; i++ {
		user := &services.CreatePersonRequest{
			Username:     null.StringFrom(fmt.Sprintf("user%d", i)),
			EmailAddress: fmt.Sprintf("user%d@gmail.com", i)}
		b, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", "/person", bytes.NewBuffer(b))
		if err != nil {
			panic(err)
		}
		res := ExecuteRequest(req)
		var p models.Person
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&p); err != nil {
			panic(err)
		}
		if p.Id == 0 {
			log.Println("RESPONSE =", res.Body)
			panic(fmt.Errorf("POST /person failed: %s", res.Body))
		}
		log.Printf("Created user with id: %d\n", p.Id)
		ids = append(ids, p.Id)
	}
	return ids
}
