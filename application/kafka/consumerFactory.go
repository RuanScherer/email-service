package kafka

type ConsumerFactory interface {
	NewConsumer(config map[string]any) (Consumer, error)
}

type Consumer interface {
	Subscribe(topic string, messagesChan chan Message) error
	Close() error
}

type Message struct {
	Key     string
	Headers map[string][]byte
	Value   []byte
}
