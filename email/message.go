package email

import (
	"encoding/base64"
	"io"

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

// toGmailMessage converts the Message into RFC 822 compliant format and
// encapsulates it within a gmail.Message type
// TODO: prefer $NAME <$address> format, should research if that's worth
func (m Message) toGmailMessage() *gmail.Message {
	var msg gmail.Message
	s := "From: " + m.From + "\r\n" +
		"reply-to: " + m.ReplyTo + "\r\n" +
		"To: " + m.To + "\r\n" +
		"Subject: " + m.Subject + "\r\n" +
		"\r\n" + m.Body
	msg.Raw = base64.StdEncoding.EncodeToString([]byte(s))

	return &msg
}
