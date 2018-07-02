package services

import "github.com/sfreiberg/gotwilio"

type TwilioSvc struct {
	Twilio            *gotwilio.Twilio
	TwilioAccountSid  string
	TwilioAuthToken   string
	TwilioPhoneNumber string
}

func (s *TwilioSvc) SendSms(to, message string) {
	s.Twilio.SendSMS(s.TwilioPhoneNumber, to, message, "", "")
}
