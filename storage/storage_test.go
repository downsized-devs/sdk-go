package storage

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	mock_log "github.com/downsized-devs/sdk-go/tests/mock/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_storage_Upload_Download_Delete(t *testing.T) {
	t.SkipNow() // remove this if you want to run the test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()

	storageConfig := Config{
		AWSS3: AWSS3Config{
			AccessKeyID:     "", // fill these in if you want to run the test
			BucketName:      "", // fill these in if you want to run the test
			Region:          "", // fill these in if you want to run the test
			SecretAccessKey: "", // fill these in if you want to run the test
		},
	}
	data := []byte(`{"name":"test","desc":"ok"}`)

	type args struct {
		ctx          context.Context
		key          string
		filename     string
		filemimetype string
		data         []byte
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantByte []byte
		wantErr  bool
	}{
		{
			name: "ok",
			args: args{
				ctx:          context.Background(),
				key:          "220426/11/testfile.json",
				filename:     "testfile",
				filemimetype: "json",
				data:         data,
			},
			want:     "https://aquahero-storage-staging.s3.ap-southeast-1.amazonaws.com/220426/11/testfile.json",
			wantByte: data,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Init(storageConfig, logger)
			got, err := s.Upload(tt.args.ctx, tt.args.key, tt.args.filename, tt.args.filemimetype, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("storage.Upload() = %v, want %v", got, tt.want)
				return
			}

			u, _ := url.Parse(got)
			gotByte, err := s.Download(tt.args.ctx, u.Path[1:])
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotByte) != string(tt.wantByte) {
				t.Errorf("storage.Download() = %v, want %v", string(gotByte), string(tt.wantByte))
				return
			}

			err = s.Delete(tt.args.ctx, u.Path[1:])
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_storage_CreateUrlByKey(t *testing.T) {
	storageConfig := Config{
		AWSS3: AWSS3Config{
			AccessKeyID:     "access-key",
			BucketName:      "bucket-name",
			Region:          "region",
			SecretAccessKey: "secret-key",
		},
	}

	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "success",
			key:  "testKey.txt",
			want: fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", storageConfig.AWSS3.BucketName, storageConfig.AWSS3.Region, "testKey.txt"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Init(storageConfig, nil)

			url := s.CreateUrlByKey(tt.key)
			assert.Equal(t, tt.want, url)
		})
	}
}
