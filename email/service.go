package email

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Config holds values used to mutate the way Service behaves
type Config struct {
	FromUser       string
	ReplyToAddress string
	// TODO: oauth2 scopes?
}

type Service struct {
	gmailService *gmail.Service
}

func NewService(ctx context.Context, credentialsPath, tokenPath string, scope ...string) (*Service, error) {
	cb, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, errors.Wrap(err, "read credentials failed")
	}

	config, err := google.ConfigFromJSON(cb, scope...)
	if err != nil {
		return nil, errors.Wrap(err, "create config from credentials failed")
	}

	tb, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "read token failed")
	}

	var token oauth2.Token
	err = json.Unmarshal(tb, &token)
	if err != nil {
		return nil, errors.Wrap(err, "decode token failed")
	}

	client := config.Client(ctx, &token)
	gmailService, err := gmail.New(client)
	if err != nil {
		return nil, errors.Wrap(err, "retrieve gmail client failed")
	}

	service := Service{
		gmailService: gmailService,
	}
	return &service, nil
}

func (s *Service) Send(msg Message) (*gmail.Message, error) {
	// FIXME: why is user "me" ?
	sendCall := s.gmailService.Users.Messages.Send("me", msg.toGmailMessage())
	if msg.Attachment != nil {
		sendCall.Media(msg.Attachment) // FIXME: use options field
	}
	gmailMessage, err := sendCall.Do()
	if err != nil {
		return nil, errors.Wrap(err, "send call failed")
	}

	return gmailMessage, nil
}
