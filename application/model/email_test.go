package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewEmail(t *testing.T) {
	req := SendEmailRequest{
		Subject: "fake email",
		Content: "This is a fake email",
	}
	_, err := govalidator.ValidateStruct(req)
	require.Error(t, err)
	require.Equal(t, "Receiver email is required", err.Error())

	req = SendEmailRequest{
		To:      "john.doe@gmail.com",
		Content: "This is a fake email",
	}
	_, err = govalidator.ValidateStruct(req)
	require.Error(t, err)
	require.Equal(t, "Subject is required", err.Error())

	req = SendEmailRequest{
		To:      "john.doe@gmail.com",
		Subject: "fake email",
	}
	_, err = govalidator.ValidateStruct(req)
	require.Error(t, err)
	require.Equal(t, "Body is required", err.Error())

	req = SendEmailRequest{
		To:      "john.doe@gmail.com",
		Subject: "fake email",
		Content: "This is a fake email",
	}
	_, err = govalidator.ValidateStruct(req)
	require.Nil(t, err)
}
