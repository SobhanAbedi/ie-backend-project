package mailer

import (
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"ie-backend-project/common"
	"ie-backend-project/model"
	"strings"
)

const GmailSMTPServer = "smtp.gmail.com:587"

type Mailer struct {
	Username string
	Auth     sasl.Client
}

func NewMailer(username, password string) Mailer {
	mailer := Mailer{Username: username}
	mailer.Auth = sasl.NewPlainClient("", username, password)
	return mailer
}

func composeEmail(sender, recipient, subject, body string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		sender, recipient, subject, body,
	))
}

func (m Mailer) SendMails(recipients []model.Student, r []interface{}, ch chan int) {
	for i, recipient := range recipients {
		to := []string{recipient.Email}
		msg := strings.NewReader(
			string(composeEmail(m.Username, recipient.Email, "Results from: "+recipient.Course.String(), recipient.String())))
		err := smtp.SendMail(GmailSMTPServer, m.Auth, m.Username, to, msg)
		if err != nil {
			r[i] = common.Error{Note: err.Error()}
		} else {
			fmt.Println("Mailer: Mail Sent:", recipient.String())
			r[i] = common.Success{Note: "Mail Sent"}
		}
	}
	fmt.Println("Mailer: Sending results through channel")
	ch <- 0
}
