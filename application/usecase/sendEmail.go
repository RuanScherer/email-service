package usecase

import (
	"github.com/asaskevich/govalidator"
	"log/slog"

	"github.com/RuanScherer/email-service/application/email"
	"github.com/RuanScherer/email-service/application/model"
)

type SendEmailUseCase interface {
	Execute(req model.SendEmailRequest) error
}

type SendEmailUseCaseImpl struct {
	emailService email.Service
}

func NewSendEmailUseCase(emailService email.Service) *SendEmailUseCaseImpl {
	return &SendEmailUseCaseImpl{
		emailService,
	}
}

func (useCase *SendEmailUseCaseImpl) Execute(req model.SendEmailRequest) error {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		slog.Error("Invalid request to send email.", slog.Any("error", err))
		return err
	}

	emailConfig, err := email.NewSendingConfig(req.To, req.Subject, req.Content)
	if err != nil {
		slog.Error("Error creating config to send email.", slog.Any("error", err))
		return err
	}

	err = useCase.emailService.SendEmail(*emailConfig)
	if err != nil {
		return err
	}
	return nil
}
