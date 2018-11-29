package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"

	"github.com/subtlepseudonym/boson/email"

	"golang.org/x/oauth2/google"
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
	b, err := ioutil.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, err := email.GetClient(config, tokFile)
	if err != nil {
		log.Fatalf("Unable to get oauth client: %v", err)
	}

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	m := email.Message{
		From:    fromUser,
		ReplyTo: replyToAddress,
		To:      toAddress,
		Subject: "Test Mail",
		Body:    "This is a test email from boson",
	}

	var msg gmail.Message
	msg.Raw = base64.StdEncoding.EncodeToString([]byte(m.String()))

	user := "me"
	_, err = srv.Users.Messages.Send(user, &msg).Do()
	if err != nil {
		log.Fatalf("Send message failed: %v", err)
	}
	//        r, err := srv.Users.Labels.List(user).Do()
	//        if err != nil {
	//                log.Fatalf("Unable to retrieve labels: %v", err)
	//        }
	//        if len(r.Labels) == 0 {
	//                fmt.Println("No labels found.")
	//                return
	//        }
	//        fmt.Println("Labels:")
	//        for _, l := range r.Labels {
	//                fmt.Printf("- %s\n", l.Name)
	//        }
}
