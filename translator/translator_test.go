package translator

import (
	"context"
	"testing"

	"github.com/downsized-devs/sdk-go/appcontext"
	"github.com/downsized-devs/sdk-go/language"
	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_translator_Translate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	conf := Config{
		FallbackLanguageID:   language.English,
		SupportedLanguageIDs: []string{language.Indonesian},
		TranslationDir:       "translations",
	}

	mockCtxEnglish := appcontext.SetAcceptLanguage(context.Background(), language.English)
	mockCtxIndo := appcontext.SetAcceptLanguage(context.Background(), language.Indonesian)

	type args struct {
		ctx    context.Context
		key    interface{}
		params []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test fallback EN",
			args: args{
				ctx: mockCtxEnglish,
				key: "test",
			},
			want: "This is a Test",
		},
		{
			name: "invalid key",
			args: args{
				ctx: mockCtxEnglish,
				key: "invalid",
			},
			wantErr: true,
		},
		{
			name: "success EN",
			args: args{
				ctx: mockCtxEnglish,
				key: "test",
			},
			want: "This is a Test",
		},
		{
			name: "success ID",
			args: args{
				ctx: mockCtxIndo,
				key: "test",
			},
			want: "Ini adalah sebuah Test",
		},
		{
			name: "success: key is empty string",
			args: args{
				ctx: mockCtxIndo,
				key: "",
			},
			want: "",
		},
		{
			name: "success: key is nil",
			args: args{
				ctx: mockCtxIndo,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := Init(conf, logger)
			got, err := ut.Translate(tt.args.ctx, tt.args.key, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("translator.Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
