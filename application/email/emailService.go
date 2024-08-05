package email

import "github.com/asaskevich/govalidator"

type Service interface {
	SendEmail(email SendingConfig) error
}

type SendingConfig struct {
	To      string `valid:"required~Receiver email is required,email~Receiver email is invalid"`
	Subject string `valid:"required~Subject is required"`
	Content string `valid:"required~Body is required"`
}

func NewSendingConfig(to, subject, content string) (*SendingConfig, error) {
	sc := &SendingConfig{
		To:      to,
		Subject: subject,
		Content: content,
	}

	_, err := govalidator.ValidateStruct(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}
