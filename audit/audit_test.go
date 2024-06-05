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
