package mailer

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	brevo "github.com/getbrevo/brevo-go/lib"
	"go.uber.org/zap"
	"time"
)

type BrevoService struct {
	cfg       config.Config
	brevoRepo *brevo.APIClient
}

func NewBrevoService(cfg config.Config, brevoRepo *brevo.APIClient) *BrevoService {
	return &BrevoService{
		cfg:       cfg,
		brevoRepo: brevoRepo,
	}
}

func (b *BrevoService) SendMail(ctx context.Context, mailReq Mail) error {
	//ctx, endFunc := trace.Start(ctx, "BrevoService.SendMail", "external_service")
	//defer endFunc()

	var recipients []brevo.SendSmtpEmailTo

	for _, toEmail := range mailReq.To {
		recipients = append(recipients, brevo.SendSmtpEmailTo{
			Email: toEmail,
		})
	}

	mailBody := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  b.cfg.GetMailerCfg().DefaultSenderName,
			Email: b.cfg.GetMailerCfg().DefaultSenderEmail,
		},
		To:          recipients,
		Subject:     mailReq.Subject,
		ScheduledAt: time.Now().Add(5 * time.Second),
	}

	if mailReq.TextContent != "" {
		mailBody.TextContent = mailReq.TextContent
	}

	if mailReq.HtmlContent != "" {
		mailBody.TextContent = mailReq.HtmlContent
	}

	if mailReq.TemplateID > 0 {
		mailBody.TemplateId = mailReq.TemplateID
	}

	if len(mailReq.Data) > 0 {
		mailBody.Params = mailReq.Data
	}

	obj, _, err := b.brevoRepo.TransactionalEmailsApi.SendTransacEmail(ctx, mailBody)
	if err != nil {
		b.cfg.Logger().ErrorWithContext(ctx, "[BrevoSvc.SendMail] Error in TransactionalEmailsApi->SendTransacEmail", zap.String("error", err.Error()))
		return err
	}

	b.cfg.Logger().InfoWithContext(ctx, "[BrevoSvc.SendMail] Success send mail", zap.String("mail_type", mailReq.Name), zap.Any("data", obj))
	return nil
}
