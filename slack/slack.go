package slack

import (
	"context"

	slackgo "github.com/slack-go/slack"
)

type Attachment slackgo.Attachment

type AttachmentField slackgo.AttachmentField

type Config struct {
	Enabled bool
	Token   string
}

type Interface interface {
	SendMessage(ctx context.Context, channelID string, attachment Attachment, fields []AttachmentField) error
}

type slack struct {
	client *slackgo.Client
	conf   Config
}

func Init(conf Config) Interface {
	client := slackgo.New(conf.Token)
	return &slack{
		client: client,
		conf:   conf,
	}
}

func (sb *slack) SendMessage(ctx context.Context, channelID string, attachment Attachment, fields []AttachmentField) error {
	if !sb.conf.Enabled {
		return nil
	}

	for _, f := range fields {
		attachment.Fields = append(attachment.Fields, slackgo.AttachmentField(f))
	}

	_, _, err := sb.client.PostMessageContext(ctx, channelID, slackgo.MsgOptionAttachments(slackgo.Attachment(attachment)))
	if err != nil {
		return err
	}

	return nil
}
