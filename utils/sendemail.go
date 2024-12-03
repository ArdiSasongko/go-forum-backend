package utils

import (
	"fmt"

	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var TokenType = map[string]string{
	"email":           "Email Validation",
	"resend_email":    "Resend Email Validation",
	"password":        "Reset Password",
	"resend_password": "Resen Reset Password",
}

func SendToken(toEmail string, tokenType string, token int32) error {
	fromEmail := env.GetEnv("EMAIL_FROM", "")
	codeEmail := env.GetEnv("EMAIL_CODE", "")

	tokenTypeDescription, exists := TokenType[tokenType]
	if !exists {
		logrus.WithField("send email", "invalid type").Error("invalid type")
		return fmt.Errorf("invalid token type: %s", tokenType)
	}

	msg := fmt.Sprintf("This is your token for %s, token: %d", tokenTypeDescription, token)
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", tokenTypeDescription)
	m.SetBody("text/html", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, codeEmail)

	if err := d.DialAndSend(m); err != nil {
		logrus.WithField("send email", err.Error()).Error(err.Error())
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
