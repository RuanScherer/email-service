package kafkaconsumers

import (
	"github.com/RuanScherer/email-service/application/kafka"
	"github.com/RuanScherer/email-service/application/usecase"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestEmailSendingRequestedConsumer_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockKafkaConsumerFactory := kafka.NewMockConsumerFactory(ctrl)
	mockKafkaConsumer := kafka.NewMockConsumer(ctrl)
	mockUseCase := usecase.NewMockSendEmailUseCase(ctrl)

	mockKafkaConsumerFactory.
		EXPECT().
		NewConsumer(gomock.Any()).
		Return(mockKafkaConsumer, nil)
	consumer, err := NewEmailSendingRequestedConsumer(mockKafkaConsumerFactory, mockUseCase)
	require.Nil(t, err)
	require.NotNil(t, consumer)

	mockKafkaConsumer.
		EXPECT().
		Subscribe("email-sending-requested", gomock.Any()).
		Return(nil)
	err = consumer.Subscribe()
	require.Nil(t, err)

	mockUseCase.
		EXPECT().
		Execute(gomock.Any()).
		Return(nil)
	consumer.messagesChan <- kafka.Message{
		Value: []byte("{\"to\": \"john.doe@gmail.com\",\"subject\": \"test\",\"content\": \"hello, test!\"}"),
	}
	time.Sleep(time.Second)
}

func TestEmailSendingRequestedConsumer_Unsubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockKafkaConsumerFactory := kafka.NewMockConsumerFactory(ctrl)
	mockKafkaConsumer := kafka.NewMockConsumer(ctrl)
	mockUseCase := usecase.NewMockSendEmailUseCase(ctrl)

	mockKafkaConsumerFactory.
		EXPECT().
		NewConsumer(gomock.Any()).
		Return(mockKafkaConsumer, nil)
	consumer, err := NewEmailSendingRequestedConsumer(mockKafkaConsumerFactory, mockUseCase)
	require.Nil(t, err)
	require.NotNil(t, consumer)

	mockKafkaConsumer.
		EXPECT().
		Subscribe("email-sending-requested", gomock.Any()).
		Return(nil)
	err = consumer.Subscribe()
	require.Nil(t, err)

	mockKafkaConsumer.
		EXPECT().
		Close()
	err = consumer.Unsubscribe()
	require.Nil(t, err)
}
