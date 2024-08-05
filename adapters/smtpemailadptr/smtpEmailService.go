package smtpemailadptr

import (
	"github.com/RuanScherer/email-service/application/email"
	"github.com/RuanScherer/email-service/config"
	"gopkg.in/gomail.v2"
	"log/slog"
	"os"
	"strconv"
)

type Service struct {
	dialer *gomail.Dialer
}

func NewService() (*Service, error) {
	port, err := strconv.Atoi(os.Getenv(config.EmailSmtpPort))
	if err != nil {
		return nil, err
	}

	return &Service{
		dialer: gomail.NewDialer(
			os.Getenv(config.EmailSmtpHost),
			port,
			os.Getenv(config.EmailSmtpUsername),
			os.Getenv(config.EmailSmtpPassword),
		),
	}, nil
}

func (service *Service) SendEmail(email email.SendingConfig) error {
	message := gomail.NewMessage()
	message.SetHeader("From", service.dialer.Username)
	message.SetHeader("To", email.To)
	message.SetHeader("Subject", email.Subject)
	message.SetBody("text/html", email.Content)

	err := service.dialer.DialAndSend(message)
	if err != nil {
		slog.Error("Error sending email.", slog.Any("error", err))
		return err
	}

	slog.Info(
		"Email sent successfully!",
		slog.String("subject", email.Subject),
		slog.String("to", email.To),
	)
	return nil
}
