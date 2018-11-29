package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// saveToken saves a token to a file path
// FIXME: decide if this is useful
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "open file failed")
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return errors.Wrap(err, "encode token failed")
	}
	return nil
}

// requestTokenFromUser prompts the user to visit a web page, authorize this application,
// generate a token, and copy it to the command line
// FIXME: decide if this is useful
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
