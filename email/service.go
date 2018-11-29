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

type Service struct {
	credentialsPath string
	tokenPath       string
	gmailService    *gmail.Service
	// TODO: oauth2 scopes?
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
		credentialsPath: credentialsPath,
		tokenPath:       tokenPath,
		gmailService:    gmailService,
	}
	return &service, nil
}

func (s *Service) Send(msg Message) error {
	// FIXME: why is user "me" ?
	sendCall := s.gmailService.Users.Messages.Send("me", msg.toGmailMessage())
	if msg.Attachment != nil {
		sendCall.Media(msg.Attachment) // FIXME: use options field
	}
	_, err := sendCall.Do()
	return err
}
