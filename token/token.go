package token

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// GetClient retrieves a token, saves the token, then returns the generated client
func GetClient(config *oauth2.Config, tokFile string) (*http.Client, error) {
	tok, err := GetTokenFromFile(tokFile)
	if err != nil {
		// FIXME: figure out if I'm even going to use this fn. If yes, just return an error here
		tok = requestTokenFromUser(config)
		SaveToken(tokFile, tok)
	}
	//log.Printf("token valid: %t\n", tok.Valid())
	return config.Client(context.Background(), tok), nil
}

// requestTokenFromUser prompts the user to visit a web page, authorize this application,
// generate a token, and copy it to the command line
func requestTokenFromUser(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

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
