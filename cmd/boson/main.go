package main

import (
	"context"
	"log"
	"net/http"

	"github.com/subtlepseudonym/boson"
	"github.com/subtlepseudonym/boson/email"

	"google.golang.org/api/gmail/v1"
)

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

	emailConfig := email.Config{
		FromUser:       fromUser,
		ReplyToAddress: replyToAddress,
	}

	emailHandler := boson.NewEmailHandler(emailConfig, emailService)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: emailHandler, // TODO: use a proper muxer
	}

	log.Fatal(srv.ListenAndServe())
}
