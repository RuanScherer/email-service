package kafkacmd

import (
	"github.com/RuanScherer/email-service/application/kafka"
	"github.com/RuanScherer/email-service/application/kafka/kafkaconsumers"
	"github.com/RuanScherer/email-service/application/usecase"
	"github.com/RuanScherer/email-service/config"
	"go.uber.org/fx"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const DefaultConsumersQuantity = 3

func StartConsumers(
	kafkaConsumerFactory kafka.ConsumerFactory,
	sendEmailUseCase usecase.SendEmailUseCase,
	shutdowner fx.Shutdowner,
) {
	var consumers []kafka.AppConsumer
	for i := range getDesiredConsumersQuantity() {
		c, err := subscribeEmailSendingRequested(kafkaConsumerFactory, sendEmailUseCase)
		if err != nil {
			slog.Error(
				"Error creating and subscribing consumer for 'email sending requested'.",
				slog.Uint64("consumerNumber", i),
				slog.Any("error", err),
			)
		} else {
			consumers = append(consumers, c)
		}
	}

	waitStopSignalAndUnsubscribe(consumers)
	err := shutdowner.Shutdown()
	if err != nil {
		slog.Error("Failed to shutdown FX.", slog.Any("error", err))
	}
}

func getDesiredConsumersQuantity() uint64 {
	consumersQuantity, err := strconv.ParseUint(os.Getenv(config.KafkaConsumerQuantity), 10, 8)
	if err != nil {
		slog.Error(
			"Failed to parse env. Fallbacking to default.",
			slog.String("envVar", config.KafkaConsumerQuantity),
			slog.Uint64("defaultValue", DefaultConsumersQuantity),
		)
		consumersQuantity = DefaultConsumersQuantity
	}
	return consumersQuantity
}

func subscribeEmailSendingRequested(
	kafkaConsumerFactory kafka.ConsumerFactory,
	sendEmailUseCase usecase.SendEmailUseCase,
) (
	*kafkaconsumers.EmailSendingRequestedConsumer,
	error,
) {
	c, err := kafkaconsumers.NewEmailSendingRequestedConsumer(kafkaConsumerFactory, sendEmailUseCase)
	if err != nil {
		return nil, err
	}

	err = c.Subscribe()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func waitStopSignalAndUnsubscribe(consumers []kafka.AppConsumer) {
	slog.Info("All consumers started. Keeping app running until a signal to stop be received.")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	slog.Info("Signal to stop has been received. Stopping app...")
	for _, c := range consumers {
		err := c.Unsubscribe()
		if err != nil {
			slog.Error("Error unsubscribing consumer.", slog.Any("error", err))
		}
	}
}
