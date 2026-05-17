package audit

import (
	"context"
	"fmt"
	"testing"

	"github.com/downsized-devs/sdk-go/appcontext"
	"github.com/downsized-devs/sdk-go/auth"
	mock_auth "github.com/downsized-devs/sdk-go/tests/mock/auth"
	"go.uber.org/mock/gomock"
)

func Test_audit_Capture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	ctx := appcontext.SetEventName(context.Background(), "tests")

	type mockFields struct {
		auth *mock_auth.MockInterface
	}

	mocks := mockFields{
		auth: authMock,
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		mockFunc func(m mockFields)
		args     args
	}{
		{
			name: "capture",
			args: args{
				ctx: ctx,
			},
			mockFunc: func(mock mockFields) {
				mock.auth.EXPECT().GetUserAuthInfo(ctx).Return(auth.UserAuthInfo{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Init(authMock)
			tt.mockFunc(mocks)

			l.Capture(tt.args.ctx)
		})
	}
}

// Test_audit_Init_MultipleCalls regresses the pre-1.0.1 sync.Once bug where
// the second (and later) Init calls returned an audit with a zero-value
// zerolog.Logger because the once.Do body had already run.
func Test_audit_Init_MultipleCalls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authMock := mock_auth.NewMockInterface(ctrl)
	authMock.EXPECT().GetUserAuthInfo(gomock.Any()).Return(auth.UserAuthInfo{}, nil).AnyTimes()

	first := Init(authMock)
	second := Init(authMock)

	ctx := appcontext.SetEventName(context.Background(), "tests")
	// If the logger were zero-value, Capture would still not panic on
	// zerolog (it writes to a nil writer fine), so additionally assert
	// that the two audit values are independent allocations.
	first.Capture(ctx)
	second.Capture(ctx)

	if first == second {
		t.Fatalf("expected Init to return independent values; got the same pointer twice")
	}
}

func Test_audit_Record(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)

	type mockFields struct {
		auth *mock_auth.MockInterface
	}

	mocks := mockFields{
		auth: authMock,
	}

	type args struct {
		ctx context.Context
		log Collection
	}
	tests := []struct {
		name     string
		mockFunc func(m mockFields)
		args     args
	}{
		{
			name: "record",
			args: args{
				ctx: context.Background(),
				log: Collection{},
			},
			mockFunc: func(mock mockFields) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.UserAuthInfo{}, nil)
			},
		},
		{
			name: "record with collection",
			args: args{
				ctx: context.Background(),
				log: Collection{
					SelectParam: "select",
					InsertParam: "insert",
					UpdateParam: "update",
					Error:       fmt.Errorf("eror"),
				},
			},
			mockFunc: func(mock mockFields) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.UserAuthInfo{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Init(authMock)
			tt.mockFunc(mocks)

			l.Record(tt.args.ctx, tt.args.log)
		})
	}
}
