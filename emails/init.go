package emails

import (
	"code-review/config"
	"code-review/logger"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/mailersend/mailersend-go"
)

var ms *mailersend.Mailersend

func init() {
	ms = mailersend.NewMailersend(config.GetEnv("MAILERSEND_API_KEY"))
}

func SendEmail(toEmail, toName, subject, text, html string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Name:  "Jeremiah Lena",
		Email: "jeremiah@trial-pr9084zz0jj4w63d.mlsender.net",
	}

	recipients := []mailersend.Recipient{
		{
			Name:  toName,
			Email: toEmail,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)

	res, error := ms.Email.Send(ctx, message)

	log.Println(*res, error)

	logger.Logger.Info("Email sent", zap.Any("response", res), zap.Error(error))

	return error
}

func SendEmailVerification(toEmail, toName, url string) (err error) {
	html := fmt.Sprintf("Please verify your email by clicking <a href=\"%s\">here</a>", url)
	text := fmt.Sprintf("Please verify your email by visiting the url: %s", url)

	err = SendEmail(toEmail, toName, "Email Verification", text, html)
	return
}

func GenerateVerificationToken(email string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
