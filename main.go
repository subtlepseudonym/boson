package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/subtlepseudonym/boson/token"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const (
	fromUser       = "Machine Spirit"
	replyToAddress = "mechanicusdeus@gmail.com"
	toAddress      = "subtlepseudonym@gmail.com"
	tokFile        = "token.json"
)

type message struct {
	From       string // from name
	ReplyTo    string // reply-to address
	To         string // to address
	Subject    string
	Body       string
	Attachment io.Reader // optional attachment
}

func (m message) String() string {
	return "From: " + m.From + "\r\n" +
		"reply-to: " + m.ReplyTo + "\r\n" +
		"To: " + m.To + "\r\n" +
		"Subject: " + m.Subject + "\r\n" +
		"\r\n" + m.Body
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tok, err := token.GetTokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		token.SaveToken(tokFile, tok)
	}
	log.Printf("token valid: %t\n", tok.Valid())
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	m := message{
		From:    fromUser,
		ReplyTo: replyToAddress,
		To:      toAddress,
		Subject: "Test Mail",
		Body:    "This is a test email from boson",
	}

	var msg gmail.Message
	msg.Raw = base64.StdEncoding.EncodeToString([]byte(m.String()))

	user := "me"
	r, err := srv.Users.Messages.Send(user, &msg).Do()
	if err != nil {
		log.Fatalf("Send message failed: %v", err)
	}
	log.Printf("%#v\n", r)
	log.Printf("%d %s\n", r.ServerResponse.HTTPStatusCode, http.StatusText(r.ServerResponse.HTTPStatusCode))
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
