package email

import (
	"testing"
	"log"
	"os"
	"net/smtp"
)

func TestSendEmail(t *testing.T) {
	//send("this is a golang test")
	CreateEmail("Kevin Chen", "1a313ee8-111f-4ca7-bb94-5c22a130f71d")
}

func send(body string) {
	from := "kevin.chen.bulk@gmail.com"
	pass := os.Getenv("ALPACA_GMAIL_PASS")
	to := "kevin.chen.bulk@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}