//go:build integration
// +build integration

package email

import (
	"context"
	"testing"

	"github.com/downsized-devs/sdk-go/log"
)

func initTestEmail() Interface {
	conf := Config{
		SMTP: SMTPConfig{
			// fill this config to run the test
			Host:     "",
			Port:     0,
			Username: "",
			Password: "",
		},
	}

	return Init(conf, log.Init(log.Config{Level: "debug"}))
}

func Test_email_SendEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		param SendEmailParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success send test email",
			args: args{
				ctx: context.Background(),
				param: SendEmailParams{
					Body:        "<strong>This is a test email</strong> and <a href=\"https://www.google.com\">This is a link</a>",
					BodyType:    BodyContentTypeHTML,
					Subject:     "Test Email",
					SenderName:  "Delos Test",
					SenderEmail: "no-reply@delosaqua.com",
					Recipients: Recipient{
						ToEmails: []string{"bambang.widyatmoko@delosaqua.com"},
					},
					Headers: map[string]string{
						"X-PM-Message-Stream": "outbound",
						"X-PM-Tag":            "test",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initTestEmail().SendEmail(tt.args.ctx, tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("email.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
