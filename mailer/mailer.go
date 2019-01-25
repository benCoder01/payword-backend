package mailer

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/go-mail/mail"
)

func Send(mailaddress string, password string) error {

	if os.Getenv("APP_ENV") != "production" {
		fmt.Println(password)
		return nil
	}

	m := mail.NewMessage()
	m.SetHeader("From", "noreply@benediktricken.de")
	m.SetHeader("To", mailaddress)
	m.SetHeader("Subject", "Password Reset")

	var templatePath string

	if os.Getenv("APP_ENV") == "production" {
		templatePath = "/go/bin/template.html"
	} else {
		templatePath = "mailer/template.html"
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	m.SetBodyWriter("text/html", func(w io.Writer) error {
		return t.Execute(w, password)
	})

	emailUsername := os.Getenv("EMAIL_USERNAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	d := mail.NewDialer("smtp.udag.de", 465, emailUsername, emailPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// insecure
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
