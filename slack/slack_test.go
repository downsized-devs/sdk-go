//go:build integration
// +build integration

package slack

import (
	"context"
	"testing"
)

func Test_slack_SendMessage(t *testing.T) {

	configEnabled := Config{
		Enabled: true,
		Token:   "", // fill these in if you want to run the test
	}

	sampleField := []AttachmentField{}
	sampleAttachment := Attachment{}

	type args struct {
		ctx        context.Context
		channelID  string
		attachment Attachment
		fields     []AttachmentField
	}
	tests := []struct {
		name    string
		conf    Config
		args    args
		wantErr bool
	}{
		{
			name: "disable",
			args: args{
				ctx:        context.Background(),
				channelID:  "",
				attachment: sampleAttachment,
				fields:     sampleField,
			},
			wantErr: false,
		},
		{
			name: "error",
			conf: configEnabled,
			args: args{
				ctx:        context.Background(),
				channelID:  "",
				attachment: sampleAttachment,
				fields:     sampleField,
			},
			wantErr: true,
		},
		{
			name: "ok",
			conf: configEnabled,
			args: args{
				ctx:       context.Background(),
				channelID: "C03F81AQYG4",
				attachment: Attachment{
					Pretext: "testing",
				},
				fields: []AttachmentField{
					{
						Title: "titel",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := Init(tt.conf)
			if err := sb.SendMessage(tt.args.ctx, tt.args.channelID, tt.args.attachment, tt.args.fields); (err != nil) != tt.wantErr {
				t.Errorf("slack.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
