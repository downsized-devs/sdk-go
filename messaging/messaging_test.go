package messaging

import (
	"context"
	"errors"
	"fmt"
	"testing"

	firebase_messaging "firebase.google.com/go/messaging"
	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeFirebase records calls and returns canned responses — implements firebaseMessenger.
type fakeFirebase struct {
	subTokens    []string
	subTopic     string
	subErr       error
	unsubTokens  []string
	unsubTopic   string
	unsubErr     error
	sendMessage  *firebase_messaging.Message
	sendErr      error
	multicastMsg *firebase_messaging.MulticastMessage
	multicastErr error
	multicastRsp *firebase_messaging.BatchResponse
	// multicastCalls records each batch passed to SendMulticastDryRun so that
	// chunking behavior can be asserted from tests.
	multicastCalls [][]string
	// multicastRspFor lets callers return different BatchResponses per batch
	// (indexed by call number). Falls back to multicastRsp when nil/short.
	multicastRspFor []*firebase_messaging.BatchResponse
}

func (f *fakeFirebase) SubscribeToTopic(_ context.Context, tokens []string, topic string) (*firebase_messaging.TopicManagementResponse, error) {
	f.subTokens = tokens
	f.subTopic = topic
	return nil, f.subErr
}

func (f *fakeFirebase) UnsubscribeFromTopic(_ context.Context, tokens []string, topic string) (*firebase_messaging.TopicManagementResponse, error) {
	f.unsubTokens = tokens
	f.unsubTopic = topic
	return nil, f.unsubErr
}

func (f *fakeFirebase) Send(_ context.Context, message *firebase_messaging.Message) (string, error) {
	f.sendMessage = message
	return "msg-id", f.sendErr
}

func (f *fakeFirebase) SendMulticastDryRun(_ context.Context, message *firebase_messaging.MulticastMessage) (*firebase_messaging.BatchResponse, error) {
	f.multicastMsg = message
	callIdx := len(f.multicastCalls)
	tokensCopy := append([]string(nil), message.Tokens...)
	f.multicastCalls = append(f.multicastCalls, tokensCopy)
	if f.multicastErr != nil {
		return nil, f.multicastErr
	}
	if callIdx < len(f.multicastRspFor) && f.multicastRspFor[callIdx] != nil {
		return f.multicastRspFor[callIdx], nil
	}
	return f.multicastRsp, nil
}

func newMessaging(fake *fakeFirebase) *messaging {
	return &messaging{
		log:      logger.Init(logger.Config{}),
		firebase: fake,
	}
}

// ------------------------- Init: SkipFirebaseInit ------------------------- //

func TestInit_SkipFirebaseInit(t *testing.T) {
	got := Init(Config{SkipFirebaseInit: true}, logger.Init(logger.Config{}), nil, nil)
	require.NotNil(t, got)
	m, ok := got.(*messaging)
	require.True(t, ok)
	assert.Nil(t, m.firebase)
}

// ------------------------- SubscribeToTopic ------------------------- //

func TestSubscribeToTopic_PassesTokenAndTopic(t *testing.T) {
	fake := &fakeFirebase{}
	m := newMessaging(fake)

	err := m.SubscribeToTopic(context.Background(), "device-1", "news")
	require.NoError(t, err)
	assert.Equal(t, []string{"device-1"}, fake.subTokens)
	assert.Equal(t, "news", fake.subTopic)
}

func TestSubscribeToTopic_PropagatesError(t *testing.T) {
	fake := &fakeFirebase{subErr: errors.New("bad token")}
	m := newMessaging(fake)

	err := m.SubscribeToTopic(context.Background(), "device-1", "news")
	require.Error(t, err)
}

// ------------------------- UnsubscribeFromTopic ------------------------- //

func TestUnsubscribeFromTopic_PassesTokenAndTopic(t *testing.T) {
	fake := &fakeFirebase{}
	m := newMessaging(fake)

	err := m.UnsubscribeFromTopic(context.Background(), "device-2", "alerts")
	require.NoError(t, err)
	assert.Equal(t, []string{"device-2"}, fake.unsubTokens)
	assert.Equal(t, "alerts", fake.unsubTopic)
}

func TestUnsubscribeFromTopic_PropagatesError(t *testing.T) {
	fake := &fakeFirebase{unsubErr: errors.New("not subscribed")}
	m := newMessaging(fake)

	err := m.UnsubscribeFromTopic(context.Background(), "device-2", "alerts")
	require.Error(t, err)
}

// ------------------------- BroadcastToTopic ------------------------- //

func TestBroadcastToTopic_SendsMessageWithDataAndTopic(t *testing.T) {
	fake := &fakeFirebase{}
	m := newMessaging(fake)

	payload := map[string]string{"k": "v"}
	err := m.BroadcastToTopic(context.Background(), "news", payload)
	require.NoError(t, err)
	require.NotNil(t, fake.sendMessage)
	assert.Equal(t, "news", fake.sendMessage.Topic)
	assert.Equal(t, payload, fake.sendMessage.Data)
}

func TestBroadcastToTopic_PropagatesError(t *testing.T) {
	fake := &fakeFirebase{sendErr: errors.New("send fail")}
	m := newMessaging(fake)

	err := m.BroadcastToTopic(context.Background(), "news", nil)
	require.Error(t, err)
}

// ------------------------- BatchSendDryRun ------------------------- //

func TestBatchSendDryRun_NoInvalidTokens(t *testing.T) {
	fake := &fakeFirebase{
		multicastRsp: &firebase_messaging.BatchResponse{
			SuccessCount: 2,
			FailureCount: 0,
			Responses: []*firebase_messaging.SendResponse{
				{Success: true},
				{Success: true},
			},
		},
	}
	m := newMessaging(fake)

	invalid, err := m.BatchSendDryRun(context.Background(), []string{"a", "b"})
	require.NoError(t, err)
	assert.Empty(t, invalid)
	require.NotNil(t, fake.multicastMsg)
	assert.Equal(t, []string{"a", "b"}, fake.multicastMsg.Tokens)
}

func TestBatchSendDryRun_PartialFailureReturnsInvalidTokens(t *testing.T) {
	fake := &fakeFirebase{
		multicastRsp: &firebase_messaging.BatchResponse{
			SuccessCount: 1,
			FailureCount: 2,
			Responses: []*firebase_messaging.SendResponse{
				{Success: true},
				{Error: errors.New("invalid 1")},
				{Error: errors.New("invalid 2")},
			},
		},
	}
	m := newMessaging(fake)

	invalid, err := m.BatchSendDryRun(context.Background(), []string{"a", "b", "c"})
	require.Error(t, err, "partial error expected")
	assert.Equal(t, []string{"b", "c"}, invalid)
}

func TestBatchSendDryRun_AllFailure(t *testing.T) {
	fake := &fakeFirebase{multicastErr: errors.New("network down")}
	m := newMessaging(fake)

	invalid, err := m.BatchSendDryRun(context.Background(), []string{"a", "b"})
	require.Error(t, err)
	assert.Empty(t, invalid)
}

// TestBatchSendDryRun_ChunksOverMax verifies that requests larger than
// MaximumTokensPerBatch are split into multiple SendMulticastDryRun calls
// and that invalid-token results merge across batches.
func TestBatchSendDryRun_ChunksOverMax(t *testing.T) {
	const total = MaximumTokensPerBatch + 1
	tokens := make([]string, total)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("tok-%d", i)
	}

	// First batch (500 tokens) reports one failure at index 0.
	first := &firebase_messaging.BatchResponse{
		SuccessCount: MaximumTokensPerBatch - 1,
		FailureCount: 1,
		Responses:    make([]*firebase_messaging.SendResponse, MaximumTokensPerBatch),
	}
	for i := range first.Responses {
		first.Responses[i] = &firebase_messaging.SendResponse{Success: true}
	}
	first.Responses[0] = &firebase_messaging.SendResponse{Error: errors.New("invalid")}

	// Second batch (1 token) all success.
	second := &firebase_messaging.BatchResponse{
		SuccessCount: 1,
		Responses:    []*firebase_messaging.SendResponse{{Success: true}},
	}

	fake := &fakeFirebase{multicastRspFor: []*firebase_messaging.BatchResponse{first, second}}
	m := newMessaging(fake)

	invalid, err := m.BatchSendDryRun(context.Background(), tokens)
	require.Error(t, err, "partial-error expected when any chunk reports failures")
	require.Len(t, fake.multicastCalls, 2, "expected two SendMulticastDryRun calls for 501 tokens")
	assert.Equal(t, MaximumTokensPerBatch, len(fake.multicastCalls[0]))
	assert.Equal(t, 1, len(fake.multicastCalls[1]))
	assert.Equal(t, []string{tokens[0]}, invalid)
}

func TestBatchSendDryRun_ChunksOverMaxAllSuccess(t *testing.T) {
	const total = MaximumTokensPerBatch + 1
	tokens := make([]string, total)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("tok-%d", i)
	}

	first := &firebase_messaging.BatchResponse{
		SuccessCount: MaximumTokensPerBatch,
		Responses:    make([]*firebase_messaging.SendResponse, MaximumTokensPerBatch),
	}
	for i := range first.Responses {
		first.Responses[i] = &firebase_messaging.SendResponse{Success: true}
	}
	second := &firebase_messaging.BatchResponse{
		SuccessCount: 1,
		Responses:    []*firebase_messaging.SendResponse{{Success: true}},
	}

	fake := &fakeFirebase{multicastRspFor: []*firebase_messaging.BatchResponse{first, second}}
	m := newMessaging(fake)

	invalid, err := m.BatchSendDryRun(context.Background(), tokens)
	require.NoError(t, err)
	assert.Empty(t, invalid)
	require.Len(t, fake.multicastCalls, 2)
}
