package token

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// GetTokenFromFile retrieves a token from a local file
func GetTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "open token file failed")
	}
	defer f.Close()

	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	if err != nil {
		return nil, errors.Wrap(err, "decode token failed")
	}

	return &tok, nil
}

// SaveToken saves a token to a file path
func SaveToken(path string, token *oauth2.Token) error {
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
