package model

type SendEmailRequest struct {
	To      string `json:"to" valid:"required~Receiver email is required"`
	Subject string `json:"subject" valid:"required~Subject is required"`
	Content string `json:"content" valid:"required~Body is required"`
}
