package email

import (
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/subtlepseudonym/boson/email/token"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

// Message holds the fields for a basic email message
type Message struct {
	From       string // from name
	ReplyTo    string // from address
	To         string
	Subject    string
	Body       string
	Attachment io.Reader
}

// String outputs the Message fields in standard email format, excluding
// the attachment
func (m Message) String() string {
	return "From: " + m.From + "\r\n" +
		"reply-to: " + m.ReplyTo + "\r\n" +
		"To: " + m.To + "\r\n" +
		"Subject: " + m.Subject + "\r\n" +
		"\r\n" + m.Body
}

func (m Message) ToGmailMessage() gmail.Message {
	var msg gmail.Message
	msg.Raw = m.String()

	if m.Attachment != nil {
		// TODO: might want to do resumable upload instead
	}
	return msg
}

// GetClient builds a new oauth2 client with the provided settings and the token
// located at the provided file path
func GetClient(config *oauth2.Config, tokenPath string) (*http.Client, error) {
	tok, err := token.GetTokenFromFile(tokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "get token from file failed")
	}

	return config.Client(context.Background(), tok), nil
}
