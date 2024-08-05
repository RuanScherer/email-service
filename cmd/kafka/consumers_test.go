package kafkacmd

import (
	"github.com/RuanScherer/email-service/application/kafka"
	"github.com/RuanScherer/email-service/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGetDesiredConsumersQuantity(t *testing.T) {
	err := os.Setenv(config.KafkaConsumerQuantity, "4")
	if err != nil {
		require.Failf(t, err.Error(), "Not able to set env %v", config.KafkaConsumerQuantity)
	}

	quantity := getDesiredConsumersQuantity()
	if quantity != 4 {
		require.Equal(t, 4, quantity)
	}
}

func TestGetDesiredConsumersQuantity_DefaultValue(t *testing.T) {
	var defaultQuantity uint64 = 3
	quantity := getDesiredConsumersQuantity()
	if quantity != defaultQuantity {
		require.Equal(t, defaultQuantity, quantity)
	}
}

func TestWaitStopSignalAndUnsubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	cm1 := kafka.NewMockAppConsumer(ctrl)
	cm2 := kafka.NewMockAppConsumer(ctrl)
	consumers := []kafka.AppConsumer{cm1, cm2}

	cm1.
		EXPECT().
		Unsubscribe().
		Times(1)
	cm2.
		EXPECT().
		Unsubscribe().
		Times(1)

	time.AfterFunc(1*time.Second, func() {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if err != nil {
			require.Errorf(t, err, "Failed to emit interrupt signal.")
		}
	})
	waitStopSignalAndUnsubscribe(consumers)
}
