package email

import (
	"fmt"
	"io"
	"strings"
)

// Message holds the fields for a basic email message
type Message struct {
	From       string // from name
	ReplyTo    string // from address
	To         []string
	Subject    string
	Body       string
	Attachment io.Reader
}

func (m Message) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"From: %s\r\nReply-To: %s\r\n",
		m.From,
		m.ReplyTo,
	))
	for _, to := range m.To {
		sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	}
	sb.WriteString(fmt.Sprintf(
		"Subject: %s\r\n\r\n%s",
		m.Subject,
		m.Body,
	))

	return sb.String()
}
