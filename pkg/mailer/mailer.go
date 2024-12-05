package mailer

import "context"

// MailService represents the interface for our mail service.
type MailService interface {
	SendMail(ctx context.Context, mailReq Mail) error
}

const (
	SEND_VERIFICATION_CODE_EMAIL         = "SEND_VERIFICATION_CODE_EMAIL"
	SEND_VERIFICATION_CODE_EMAIL_SUBJECT = "MeowHasiswa - Verification Code"
	SEND_RESET_PASSWORD_EMAIL            = "SEND_RESET_PASSWORD_EMAIL"
	SEND_RESET_PASSWORD_EMAIL_SUBJECT    = "MeowHasiswa - Reset Password"
)

// Mail represents a email request
type Mail struct {
	From        string
	To          []string
	Name        string
	Subject     string
	HtmlContent string
	TextContent string
	TemplateID  int64
	Data        map[string]interface{}
}
