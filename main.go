package main

import (
	"github.com/RuanScherer/email-service/adapters/kafkaadptr"
	"github.com/RuanScherer/email-service/adapters/smtpemailadptr"
	"github.com/RuanScherer/email-service/application/email"
	"github.com/RuanScherer/email-service/application/kafka"
	"github.com/RuanScherer/email-service/application/usecase"
	kafkacmd "github.com/RuanScherer/email-service/cmd/kafka"
	"go.uber.org/fx"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file", err)
		panic(err)
	}

	fx.New(
		fx.Provide(
			fx.Annotate(
				smtpemailadptr.NewService,
				fx.As(new(email.Service)),
			),
			fx.Annotate(
				usecase.NewSendEmailUseCase,
				fx.As(new(usecase.SendEmailUseCase)),
			),
			fx.Annotate(
				kafkaadptr.NewConsumerFactory,
				fx.As(new(kafka.ConsumerFactory)),
			),
		),

		fx.Invoke(kafkacmd.StartConsumers),
	).Run()
}
