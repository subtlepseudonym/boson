package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/subtlepseudonym/boson/email"
)

type SMS struct {
	To   []string `json:"to"`
	Body string   `json:"body"`
}

type SMSConfig struct {
	From string
}

type smsHandler struct {
	config  SMSConfig
	service *email.Service
}

func NewSMSHandler(cfg SMSConfig, srv *email.Service) http.Handler {
	return smsHandler{
		config:  cfg,
		service: srv,
	}
}

func (h smsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var msg SMS
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonMessage("unable to decode request body"))
		return
	}

	message := email.Message{
		From: h.config.From,
		To:   msg.To,
		Body: msg.Body,
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
