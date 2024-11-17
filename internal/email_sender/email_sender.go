package email_sender

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

var (
	senderEmail, emailExists   = os.LookupEnv("EMAIL")
	senderPassword, passExists = os.LookupEnv("PROGRAM_CODE")
)

func EmailSender(to, subject, body string) error {
	if !emailExists || !passExists {
		return fmt.Errorf("environment variables EMAIL or PROGRAM_CODE not set")
	}

	msg := buildMessage(senderEmail, to, subject, body)
	auth := smtp.PlainAuth("", senderEmail, senderPassword, "smtp.gmail.com")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.gmail.com",
	}

	conn, err := tls.Dial("tcp", "smtp.gmail.com:465", tlsConfig)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, "smtp.gmail.com")
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(senderEmail); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

func buildMessage(from, to, subject, body string) string {
	return fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s", from, to, subject, body)
}
