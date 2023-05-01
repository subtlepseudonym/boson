package email

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

type Service struct {
	From     string
	password string
}

func NewService(credentialsPath string) (*Service, error) {
	var secrets struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	f, err := os.Open(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("open credentials: %w", err)
	}
	err = json.NewDecoder(f).Decode(&secrets)
	if err != nil {
		return nil, fmt.Errorf("decode credentials: %w", err)
	}

	return &Service{
		From:     secrets.Email,
		password: secrets.Password,
	}, nil
}

func (s *Service) Send(msg Message) error {
	client, err := smtp.Dial(fmt.Sprintf("%s:%s", smtpHost, smtpPort))
	if err != nil {
		return fmt.Errorf("dial smtp addr: %w", err)
	}
	defer client.Quit()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	err = client.StartTLS(tlsConfig)
	if err != nil {
		return fmt.Errorf("smtp start tls: %w", err)
	}
	err = client.Auth(smtp.PlainAuth("", s.From, s.password, smtpHost))
	if err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	err = client.Mail(s.From)
	if err != nil {
		return fmt.Errorf("smtp mail: %q: %w", s.From, err)
	}
	for _, to := range msg.To {
		err = client.Rcpt(to)
		if err != nil {
			return fmt.Errorf("smtp rcpt: %q: %w", to, err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	_, err = w.Write([]byte(msg.String()))
	w.Close()
	if err != nil {
		return fmt.Errorf("smtp write data: %w", err)
	}

	return nil
}
