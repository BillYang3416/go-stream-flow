package email

import (
	"context"

	"github.com/bgg/go-file-gate/internal/entity"
	"github.com/bgg/go-file-gate/pkg/logger"
	mail "github.com/xhit/go-simple-mail/v2"
)

type UserUploadedFileEmailSender struct {
	smtp *mail.SMTPClient
	l    logger.Logger
}

func NewUserUploadedFileEmailSender(smtp *mail.SMTPClient, l logger.Logger) *UserUploadedFileEmailSender {
	return &UserUploadedFileEmailSender{smtp: smtp, l: l}
}

func (s *UserUploadedFileEmailSender) Send(ctx context.Context, uuf entity.UserUploadedFile) error {
	s.l.Info("Sending email to %s", uuf.EmailRecipient)
	email := mail.NewMSG()
	email.SetFrom("bgg@mail.com").
		AddTo(uuf.EmailRecipient).
		SetSubject("File uploaded").
		SetBody(mail.TextHTML, "File uploaded")

	attachment := mail.File{
		Name: uuf.Name,
		Data: uuf.Content,
	}
	email.Attach(&attachment)

	err := email.Send(s.smtp)
	if err != nil {
		s.l.Error(err)
		return err
	}

	return nil
}
