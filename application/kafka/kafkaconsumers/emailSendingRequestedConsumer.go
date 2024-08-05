package kafkaconsumers

import (
	"encoding/json"
	"github.com/RuanScherer/email-service/application/kafka"
	"github.com/RuanScherer/email-service/application/model"
	"github.com/RuanScherer/email-service/application/usecase"
	"github.com/RuanScherer/email-service/config"
	"log/slog"
	"os"
)

const EmailSendingRequestedTopic = "email-sending-requested"

type EmailSendingRequestedConsumer struct {
	consumer         kafka.Consumer
	messagesChan     chan kafka.Message
	sendEmailUseCase usecase.SendEmailUseCase
}

func NewEmailSendingRequestedConsumer(
	consumerFactory kafka.ConsumerFactory,
	sendEmailUseCase usecase.SendEmailUseCase,
) (*EmailSendingRequestedConsumer, error) {
	consumer, err := consumerFactory.NewConsumer(map[string]any{
		"bootstrap.servers":        os.Getenv(config.KafkaBootstrapServers),
		"group.id":                 os.Getenv(config.KafkaConsumerGroupId),
		"auto.offset.reset":        "beginning",
		"allow.auto.create.topics": "true",
	})
	if err != nil {
		return nil, err
	}

	return &EmailSendingRequestedConsumer{
		consumer:         consumer,
		messagesChan:     nil,
		sendEmailUseCase: sendEmailUseCase,
	}, nil
}

func (c *EmailSendingRequestedConsumer) Subscribe() error {
	c.messagesChan = make(chan kafka.Message)
	err := c.consumer.Subscribe(EmailSendingRequestedTopic, c.messagesChan)
	if err != nil {
		close(c.messagesChan)
		return err
	}

	go func() {
		for m := range c.messagesChan {
			message := &model.SendEmailRequest{}
			err := json.Unmarshal(m.Value, message)
			if err != nil {
				slog.Error(
					"Failed to parse message.",
					slog.String("message", string(m.Value)),
					slog.Any("error", err),
				)
				continue
			}

			// TODO: implement error handling
			_ = c.sendEmailUseCase.Execute(*message)
		}
	}()
	return nil
}

func (c *EmailSendingRequestedConsumer) Unsubscribe() error {
	err := c.consumer.Close()
	if err != nil {
		return err
	}

	close(c.messagesChan)
	return nil
}
