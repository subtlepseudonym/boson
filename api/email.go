package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/subtlepseudonym/boson/email"
)

// Email defines the email.Message fields that are exposed to REST clients
// TODO: add support for multiple from addresses
type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// EmailConfig holds values used to mutate the way Service behaves
type EmailConfig struct {
	From    string
	ReplyTo string
}

type emailHandler struct {
	config  EmailConfig
	service *email.Service
}

func NewEmailHandler(cfg EmailConfig, srv *email.Service) emailHandler {
	return emailHandler{
		config:  cfg,
		service: srv,
	}
}

// ServeHTTP accepts a POST request with a JSON encoded request body containing
// the email contents, as defined by the Email struct
// FIXME: accept attachments
func (h emailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var msg email.Message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonMessage("unable to decode request body"))
		return
	}

	// TODO: allow client to pick which address they'd like to send from
	// 		 This will require sending a password as part of the request
	msg.From = h.config.From
	msg.ReplyTo = h.config.ReplyTo

	err = h.service.Send(msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonMessage("unable to send email"))
		return
	}
	w.Write(jsonMessage("sent"))
}
