package mailer

import (
	"fmt"
	"html/template"
	"io"

	"github.com/go-mail/mail"
)

func Send(mailaddress string, password string) error {

	m := mail.NewMessage()
	m.SetHeader("From", "noreply@benediktricken.de")
	m.SetHeader("To", mailaddress)
	m.SetHeader("Subject", "Password Reset")

	t, err := template.ParseFiles("mailer/template.html")
	if err != nil {
		return err
	}

	m.SetBodyWriter("text/html", func(w io.Writer) error {
		return t.Execute(w, password)
	})

	d := mail.NewDialer("smtp.udag.de", 587, "benediktricken-de-0001", "8ztwtwtw")
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Password: ", password)
	return nil
}
