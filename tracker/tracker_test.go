//go:build integration
// +build integration

package tracker

import (
	"context"
	"fmt"
	"testing"

	"github.com/downsized-devs/sdk-go/log"
)

func Test_tracker_Push(t *testing.T) {
	type fields struct {
		opt Options
		log log.Interface
	}
	type args struct {
		ctx          context.Context
		trackingName string
		labels       map[string]string
	}

	mocklabels1 := map[string]string{
		"user_id":    fmt.Sprintf("%d", 1),
		"url":        "/v1/user/profile",
		"error_code": fmt.Sprintf("%d", 400),
	}

	mocklabels2 := map[string]string{
		"user_id":    fmt.Sprintf("%d", 2),
		"url":        "/v1/user/profile",
		"error_code": fmt.Sprintf("%d", 400),
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "succes tracking error count - user 1",
			fields: fields{
				opt: Options{
					Enabled: true,
					URL:     defaultURL,
					Port:    defaultPort,
					JobName: defaultJobName,
				},
			},
			args: args{
				ctx:          context.Background(),
				trackingName: "error_count",
				labels:       mocklabels1,
			},
			wantErr: false,
		},
		{
			name: "succes tracking error count - user 2",
			fields: fields{
				opt: Options{
					Enabled: true,
					URL:     defaultURL,
					Port:    defaultPort,
					JobName: defaultJobName,
				},
			},
			args: args{
				ctx:          context.Background(),
				trackingName: "error_count",
				labels:       mocklabels2,
			},
			wantErr: false,
		},
		{
			name: "success disable push tracker",
			fields: fields{
				opt: Options{
					Enabled: false,
					URL:     defaultURL,
					Port:    defaultPort,
					JobName: defaultJobName,
				},
			},
			args: args{
				ctx:          context.Background(),
				trackingName: "error_count",
				labels:       mocklabels1,
			},
			wantErr: false,
		},
		{
			name: "error tracking error count",
			fields: fields{
				opt: Options{
					Enabled: true,
					URL:     defaultURL,
					Port:    ":8080",
					JobName: defaultJobName,
				},
			},
			args: args{
				ctx:          context.Background(),
				trackingName: "error_count",
				labels:       mocklabels1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Init(tt.fields.opt, tt.fields.log)
			if err := tr.Push(tt.args.ctx, tt.args.trackingName, tt.args.labels); (err != nil) != tt.wantErr {
				t.Errorf("tracker.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
