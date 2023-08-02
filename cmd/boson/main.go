package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/subtlepseudonym/boson/api"
	"github.com/subtlepseudonym/boson/email"

	"github.com/emersion/go-smtp"
)

// TODO: read these from file
const (
	smtpHost = "smtp.gmail.com"
	smtpPort = 587

	from                   = "Machine Spirit"
	replyTo                = "mechanicusdeus@gmail.com"
	defaultCredentialsPath = "secrets/credentials.json"
)

var credentialsPath string

func main() {
	flag.StringVar(&credentialsPath, "credentials", defaultCredentialsPath, "Path to JSON file with email credentials")
	flag.Parse()

	emailService, err := email.NewService(smtpHost, smtpPort, credentialsPath)
	if err != nil {
		log.Fatalf("create new email service: %s", err)
	}

	emailConfig := api.EmailConfig{
		From:    fmt.Sprintf("%q <%s>", from, replyTo),
		ReplyTo: replyTo,
	}
	emailHandler := api.NewEmailHandler(emailConfig, emailService)

	smsConfig := api.SMSConfig{
		From: fmt.Sprintf("%q <%s>", from, replyTo),
	}

	// exit gracefully
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	stop := make(chan struct{})

	go func() {
		<-interrupt
		close(stop)
	}()

	// set up smtp passthrough
	smtpSrv := smtp.NewServer(emailHandler)
	smtpSrv.Addr = ":40589"
	smtpSrv.Domain = "boson.home.arpa" // TODO: use a real value, from file
	smtpSrv.WriteTimeout = 60 * time.Second
	smtpSrv.ReadTimeout = 60 * time.Second
	smtpSrv.MaxMessageBytes = 1024 * 1024
	smtpSrv.MaxRecipients = 8
	smtpSrv.AllowInsecureAuth = true

	// set up push-notification rest service
	mux := http.NewServeMux()
	mux.Handle("/email", emailHandler)
	mux.Handle("/sms", api.NewSMSHandler(smsConfig, emailService))
	restSrv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	go func() {
		err := restSrv.ListenAndServe()
		log.Println("REST:", err)
	}()
	log.Printf("Listening on %s\n", restSrv.Addr)
	go func() {
		err := smtpSrv.ListenAndServe()
		log.Println("SMTP:", err)
	}()
	log.Printf("SMTP passthrough on %s\n", smtpSrv.Addr)

	<-stop
	smtpSrv.Shutdown(context.Background())
	restSrv.Shutdown(context.Background())
}
