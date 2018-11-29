package main

import (
	"context"
	"log"

	"github.com/subtlepseudonym/boson/email"

	"google.golang.org/api/gmail/v1"
)

const (
	fromUser       = "Machine Spirit"
	replyToAddress = "mechanicusdeus@gmail.com"
	toAddress      = "subtlepseudonym@gmail.com"
	tokFile        = "secrets/token.json"
	credsFile      = "secrets/credentials.json"
)

func main() {
	srv, err := email.NewService(context.Background(), credsFile, tokFile, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("create new email service failed: %s", err)
	}

	m := email.Message{
		From:    fromUser,
		ReplyTo: replyToAddress,
		To:      toAddress,
		Subject: "Test Mail",
		Body:    "This is a test email from boson",
	}

	err = srv.Send(m)
	if err != nil {
		log.Fatalf("send message failed: %s", err)
	}
}
