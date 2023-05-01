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
	from      = "Machine Spirit"
	replyTo   = "mechanicusdeus@gmail.com"
	credsFile = "secrets/credentials.json"
)

func main() {
	emailService, err := email.NewService(credsFile)
	if err != nil {
		log.Fatalf("create new email service: %s", err)
	}

	emailConfig := api.EmailConfig{
		From:    fmt.Sprintf("%q <%s>", from, replyTo),
		ReplyTo: replyTo,
	}

	emailHandler := api.NewEmailHandler(emailConfig, emailService)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: emailHandler, // TODO: use a proper muxer
	}

	log.Printf("Listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
