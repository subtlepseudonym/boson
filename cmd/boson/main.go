package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/subtlepseudonym/boson/api"
	"github.com/subtlepseudonym/boson/email"
)

// TODO: read these from file
const (
	smtpHost = "smtp.gmail.com"
	smtpPort = 587

	from      = "Machine Spirit"
	replyTo   = "mechanicusdeus@gmail.com"
	credsFile = "secrets/credentials.json"
)

func main() {
	emailService, err := email.NewService(smtpHost, smtpPort, credsFile)
	if err != nil {
		log.Fatalf("create new email service: %s", err)
	}

	emailConfig := api.EmailConfig{
		From:    fmt.Sprintf("%q <%s>", from, replyTo),
		ReplyTo: replyTo,
	}
	smsConfig := api.SMSConfig{
		From: fmt.Sprintf("%q <%s>", from, replyTo),
	}

	mux := http.NewServeMux()
	mux.Handle("/email", api.NewEmailHandler(emailConfig, emailService))
	mux.Handle("/sms", api.NewSMSHandler(smsConfig, emailService))

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	log.Printf("Listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
