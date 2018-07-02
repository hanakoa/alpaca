package email

import (
	"github.com/matcornic/hermes"
	"io/ioutil"
)

//var (
//	Sender = utils.MustEnv("SENDER")
//)
//
//func SendEmail() {
//	log.Println(Sender)
//}

func CreateEmail(name, resetCode string) {
	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Alpaca",
			Link: "https://example-hermes.com/",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				"You recently requested a password reset.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started, please click here:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Reset my password",
						Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"If you cannot view links, navigate to the password reset screen in the login page. Then manually enter the code:",
				resetCode,
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		// TODO remove panic
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	// Optionally, preview the generated HTML e-mail by writing it to a local file
	err = ioutil.WriteFile("preview.html", []byte(emailBody), 0644)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}
}