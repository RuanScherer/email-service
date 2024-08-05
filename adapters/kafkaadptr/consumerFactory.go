package kafkaadptr

import (
	"errors"
	"log/slog"

	appkafka "github.com/RuanScherer/email-service/application/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ConsumerFactory struct{}

func NewConsumerFactory() *ConsumerFactory {
	return &ConsumerFactory{}
}

func (cf *ConsumerFactory) NewConsumer(config map[string]any) (appkafka.Consumer, error) {
	consumer := &consumer{}
	if config == nil {
		return nil, errors.New("config is required")
	}

	cfg := kafka.ConfigMap{}
	for k, v := range config {
		cfg[k] = v
	}

	kafkaConsumer, err := kafka.NewConsumer(&cfg)
	if err != nil {
		slog.Error("Error creating kafka consumer.", slog.Any("error", err))
		return nil, err
	}

	consumer.consumer = kafkaConsumer
	return consumer, nil
}

type consumer struct {
	consumer *kafka.Consumer
	polling  bool
}

func (c *consumer) Subscribe(topic string, messagesChan chan appkafka.Message) error {
	if c.consumer == nil {
		return errors.New("kafka consumer is not initialized")
	}

	err := c.consumer.Subscribe(topic, nil)
	if err != nil {
		slog.Error(
			"Error subscribing to topic.",
			slog.String("topic", topic),
			slog.Any("error", err),
		)
		return err
	}

	go c.poll(messagesChan)
	return nil
}

func (c *consumer) poll(messagesChan chan appkafka.Message) {
	c.polling = true
	for c.polling == true {
		evt := c.consumer.Poll(100)
		switch evt.(type) {
		case *kafka.Message:
			message := evt.(*kafka.Message)
			headers := make(map[string][]byte)
			for _, h := range message.Headers {
				headers[h.Key] = h.Value
			}

			messagesChan <- appkafka.Message{
				Key:     string(message.Key),
				Headers: headers,
				Value:   message.Value,
			}
		case kafka.Error:
			slog.Error(
				"Error consuming event.",
				slog.String("error", evt.(kafka.Error).Error()),
			)
		default:
		}
	}
}

func (c *consumer) Close() error {
	c.polling = false
	return c.consumer.Close()
	// esperar fechar
}
