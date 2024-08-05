package usecase

import (
	"errors"
	"github.com/RuanScherer/email-service/application/email"
	"github.com/RuanScherer/email-service/application/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSendEmailUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockEmailService := email.NewMockService(ctrl)
	useCase := NewSendEmailUseCase(mockEmailService)

	err := useCase.Execute(model.SendEmailRequest{
		Subject: "fake email",
		Content: "This is a fake email",
	})
	require.Error(t, err)
	require.Equal(t, "Receiver email is required", err.Error())

	err = useCase.Execute(model.SendEmailRequest{
		To:      "john.doe@gmail.com",
		Content: "This is a fake email",
	})
	require.Error(t, err)
	require.Equal(t, "Subject is required", err.Error())

	err = useCase.Execute(model.SendEmailRequest{
		To:      "john.doe@gmail.com",
		Subject: "fake email",
	})
	require.Error(t, err)
	require.Equal(t, "Body is required", err.Error())

	mockEmailService.
		EXPECT().
		SendEmail(gomock.Any()).
		Return(errors.New("fake error"))

	err = useCase.Execute(model.SendEmailRequest{
		To:      "john.doe@gmail.com",
		Subject: "fake email",
		Content: "This is a fake email",
	})
	require.Error(t, err)

	mockEmailService.
		EXPECT().
		SendEmail(gomock.Any()).
		Return(nil)

	err = useCase.Execute(model.SendEmailRequest{
		To:      "john.doe@gmail.com",
		Subject: "fake email",
		Content: "This is a fake email",
	})
	require.Nil(t, err)
}
