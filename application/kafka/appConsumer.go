package kafka

type AppConsumer interface {
	Subscribe() error
	Unsubscribe() error
}
