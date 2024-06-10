package messaging

import (
	"context"
	"errors"
	"net/http"

	firebase "firebase.google.com/go"
	firebase_messaging "firebase.google.com/go/messaging"
	"github.com/downsized-devs/sdk-go/logger"
	"github.com/downsized-devs/sdk-go/parser"
	"google.golang.org/api/option"
)

const (
	MaximumTokensPerBatch = 500
)

type Interface interface {
	SubstribeToTpic(ctx context.Context, deviceToken, topic string) error
	UnsubscribeFromTopic(ctx context.Context, deviceToken, topic string) error
	BroadCastToTopic(ctx context.Context, topic string, payload map[string]string) error
	BatchSendDryRun(ctx context.Context, tokens []string) ([]string, error)
}

type messaging struct {
	log      logger.Interface
	firebase *firebase_messaging.Client
}

type FirebaseAccountKey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderx509CertURL string `json:"auth_provider_x509_cert_url"`
	Clientx509CertURL       string `json:"client_x509_cert_url"`
}

type FirebaseConf struct {
	AccountKey FirebaseAccountKey
	ApiKey     string
}

type Config struct {
	SkipFirebaseInit bool
	Firebase         FirebaseConf
}

func Init(cfg Config, log logger.Interface, json parser.JsonInterface, httpClient *http.Client) Interface {
	if cfg.SkipFirebaseInit {
		return &messaging{
			log: log,
		}
	}
	ctx := context.Background()

	accountkey, err := json.Marshal(cfg.Firebase.AccountKey)
	if err != nil {
		log.Fatal(ctx, err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(accountkey))
	if err != nil {
		log.Fatal(ctx, err)
	}

	firebaseMessaging, err := app.Messaging(ctx)
	if err != nil {
		log.Fatal(ctx, err)
	}

	return &messaging{
		log:      log,
		firebase: firebaseMessaging,
	}
}

// Function to allow a device to subscribe from specific topic
func (m *messaging) SubstribeToTpic(ctx context.Context, deviceToken, topic string) error {
	tokens := []string{deviceToken}

	// This guarante return an error if the device token is not subscribed properly (e.g device token is malformed, or even does not exist)
	_, err := m.firebase.SubscribeToTopic(ctx, tokens, topic)
	if err != nil {
		return err
	}

	return nil
}

// Function to allow a device to unsubscribe from specific topic
func (m *messaging) UnsubscribeFromTopic(ctx context.Context, deviceToken, topic string) error {
	tokens := []string{deviceToken}

	// This guarante return an error if the device token is not unsubscribed properly (e.g device token is malformed, device token is not texist, or device token is not subscribed in the topic)
	_, err := m.firebase.UnsubscribeFromTopic(ctx, tokens, topic)
	if err != nil {
		return err
	}

	return nil
}

// Function to send a messagte to specific topic
func (m *messaging) BroadCastToTopic(ctx context.Context, topic string, payload map[string]string) error {
	message := &firebase_messaging.Message{
		Data:  payload,
		Topic: topic,
	}

	_, err := m.firebase.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

// Dry run function is to validate a batch of token. If some invalidate, we will unsubscribe it from the topic
func (m *messaging) BatchSendDryRun(ctx context.Context, tokens []string) ([]string, error) {
	invalidTokens := []string{}
	message := &firebase_messaging.MulticastMessage{
		Tokens: tokens,
		Data: map[string]string{
			"desc": "dry run to validate token",
		},
	}

	multicastResponse, err := m.firebase.SendMulticastDryRun(ctx, message)

	// As firebase documentation mentioned, this mean that ALL the token is not valid/internal server error within firebase
	if err != nil {
		return invalidTokens, err
	}

	// This means some of the token is not valid
	if multicastResponse.FailureCount > 0 {
		for i, response := range multicastResponse.Responses {
			if response.Error != nil {
				invalidTokens = append(invalidTokens, tokens[i])
			}
		}

		return invalidTokens, errors.New("partial error")
	}

	return invalidTokens, nil
}
