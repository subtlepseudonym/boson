package main

import (
	"context"
	"log"
	"net/http"

	"github.com/subtlepseudonym/boson/api"
	"github.com/subtlepseudonym/boson/email"

	"google.golang.org/api/gmail/v1"
)

// TODO: read these from file
const (
	fromUser       = "Machine Spirit"
	replyToAddress = "mechanicusdeus@gmail.com"
	tokFile        = "secrets/token.json"
	credsFile      = "secrets/credentials.json"
)

func main() {
	emailService, err := email.NewService(context.Background(), credsFile, tokFile, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("create new email service failed: %s", err)
	}

	emailConfig := api.EmailConfig{
		FromUser:       fromUser,
		ReplyToAddress: replyToAddress,
	}

	emailHandler := api.NewEmailHandler(emailConfig, emailService)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: emailHandler, // TODO: use a proper muxer
	}

	log.Printf("Listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
