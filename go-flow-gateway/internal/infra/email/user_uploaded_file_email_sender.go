package email

import (
	"context"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	mail "github.com/xhit/go-simple-mail/v2"
)

type UserUploadedFileEmailSender struct {
	smtp   *mail.SMTPClient
	logger logger.Logger
}

func NewUserUploadedFileEmailSender(smtp *mail.SMTPClient, l logger.Logger) *UserUploadedFileEmailSender {
	return &UserUploadedFileEmailSender{smtp: smtp, logger: l}
}

func (s *UserUploadedFileEmailSender) Send(ctx context.Context, uuf entity.UserUploadedFile) error {
	s.logger.Info("UserUploadedFileEmailSender - Send: sending email", "userUploadedFileID", uuf.ID)

	email := mail.NewMSG()
	email.SetFrom("bgg@mail.com").
		AddTo(uuf.EmailRecipient).
		SetSubject("Your File Upload Confirmation").
		SetBody(mail.TextHTML, "<h1>File Upload Successful</h1>"+
			"<p>Hello,</p>"+
			"<p>We have successfully received your file upload.</p>"+
			"<p><b>File Name:</b> "+uuf.Name+"</p>"+
			"<p>If you have any questions or need further assistance, please do not hesitate to contact us.</p>"+
			"<p>Best Regards,<br>Your Support Team</p>")

	attachment := mail.File{
		Name: uuf.Name,
		Data: uuf.Content,
	}
	email.Attach(&attachment)

	err := email.Send(s.smtp)
	if err != nil {
		s.logger.Error("UserUploadedFileEmailSender - Send: failed to send email", "error", err)
		return err
	}

	s.logger.Info("UserUploadedFileEmailSender - Send: successfully sent email", "userUploadedFileID", uuf.ID)
	return nil
}
