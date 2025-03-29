package email

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
	"gopkg.in/gomail.v2"
)

const (
	emailRawHeaderFrom = "From"
	emailRawHeaderTo   = "To"
	emailRawHeaderCc   = "Cc"
	emailRawHeaderBcc  = "Bcc"
	emailRawSubject    = "Subject"
)

type Interface interface {
	SendEmail(ctx context.Context, params SendEmailParams) error
	GenerateBody() TemplateInterface
}

type Config struct {
	SMTP     SMTPConfig
	Template TemplateConfig
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	TLSConfig struct {
		InsecureSkipVerify bool
	}
}

type AWSSESConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type email struct {
	dialer   *gomail.Dialer
	config   Config
	log      logger.Interface
	template TemplateInterface
}

func Init(cfg Config, log logger.Interface) Interface {
	dialer := gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)
	if cfg.SMTP.TLSConfig.InsecureSkipVerify {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	return &email{
		dialer:   dialer,
		config:   cfg,
		log:      log,
		template: initTemplate(cfg.Template, log),
	}
}

func (e *email) SendEmail(ctx context.Context, param SendEmailParams) error {
	if param.BodyType == "" {
		param.BodyType = BodyContentTypePlain
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader(emailRawHeaderFrom, fmt.Sprintf("%s <%s>", param.SenderName, param.SenderEmail))
	mailer.SetHeader(emailRawHeaderTo, param.Recipients.ToEmails...)
	mailer.SetHeader(emailRawHeaderCc, param.Recipients.CCEmails...)
	mailer.SetHeader(emailRawHeaderBcc, param.Recipients.BCCEmails...)
	mailer.SetHeader(emailRawSubject, param.Subject)
	mailer.SetBody(param.BodyType, param.Body)
	for hk, hv := range param.Headers {
		mailer.SetHeader(hk, hv)
	}
	for i := range param.Attachments {
		mailer.Attach(param.Attachments[i])
	}

	if err := e.dialer.DialAndSend(mailer); err != nil {
		return errors.NewWithCode(codes.CodeSendEmailFailed, "failed to send email, with err: %v", err)
	}

	return nil
}

func (e *email) GenerateBody() TemplateInterface {
	return e.template
}
