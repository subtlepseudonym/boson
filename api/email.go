package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/subtlepseudonym/boson/email"

	"github.com/emersion/go-smtp"
)

// Email defines the email.Message fields that are exposed to REST clients
// TODO: add support for multiple from addresses
type Email struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// EmailConfig holds values used to mutate the way Service behaves
type EmailConfig struct {
	From    string
	ReplyTo string
}

type EmailHandler struct {
	config  EmailConfig
	service *email.Service
}

func NewEmailHandler(cfg EmailConfig, srv *email.Service) *EmailHandler {
	return &EmailHandler{
		config:  cfg,
		service: srv,
	}
}

// ServeHTTP accepts a POST request with a JSON encoded request body containing
// the email contents, as defined by the Email struct
// FIXME: accept attachments
func (h *EmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(jsonMessage("only method POST allowed"))
		return
	}
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonMessage("request body cannot be empty"))
		return
	}

	var msg Email
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonMessage("unable to decode request body"))
		return
	}

	// TODO: allow client to pick which address they'd like to send from
	// 		 This will require sending a password as part of the request
	message := email.Message{
		From:    h.config.From,
		ReplyTo: h.config.ReplyTo,
		To:      msg.To,
		Headers: []string{fmt.Sprint("Subject:", msg.Subject)},
		Body:    msg.Body,
	}
	err = h.service.Send(message)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonMessage("unable to send email"))
		return
	}
	w.Write(jsonMessage("sent"))
}

// NewSession is a stub that returns a canned session
func (h *EmailHandler) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &session{
		service: h.service,
		Message: email.Message{
			From:    h.config.From,
			ReplyTo: h.config.ReplyTo,
		},
	}, nil
}

type session struct {
	service *email.Service
	Message email.Message
}

func (s *session) AuthPlain(username, password string) error {
	// TODO: check password against file
	// TODO: add auth checking to rcpt() and data()
	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	s.Message.Headers = append(s.Message.Headers, fmt.Sprintf("Forwarded-For: %s", from))
	return nil
}

func (s *session) Rcpt(to string) error {
	s.Message.To = append(s.Message.To, to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}
	body := string(b)

	// FIXME: there's got to be a better way than using regex
	subjectRegex := regexp.MustCompile(`(?i)(subject:.*)\r\n`)
	sidx := subjectRegex.FindStringSubmatchIndex(body)
	if len(sidx) >= 4 {
		subject := body[sidx[2]:sidx[3]]
		s.Message.Headers = append(s.Message.Headers, subject)
		body = body[:sidx[2]] + body[sidx[3]:] // remove subject line
	}

	lineBreakRegex := regexp.MustCompile(`\r\n\r\n`) // indicates end of headers
	idx := lineBreakRegex.FindStringIndex(body)
	if len(idx) >= 2 {
		headers := body[:idx[0]]
		for _, header := range strings.Split(headers, "\r\n") {
			s.Message.Headers = append(s.Message.Headers, fmt.Sprintf("Proxied-%s", header))
		}
		s.Message.Body = body[idx[1]:]
	} else {
		s.Message.Body = body
	}

	err = s.service.Send(s.Message)
	s.Reset()
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (s *session) Reset() {
	s.Message.To = nil
	s.Message.Headers = nil
	s.Message.Body = ""
}

func (s *session) Logout() error {
	return nil
}
