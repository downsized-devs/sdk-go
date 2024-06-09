package local_storage

import (
	"context"
	"os"
	"testing"

	"github.com/blevesearch/bleve"
	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"go.uber.org/mock/gomock"
)

func TestNewIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	type args struct {
		ctx       context.Context
		indexPath string
	}

	tests := []struct {
		name          string
		args          args
		prepIndexMock func() bleve.Index
		wantErr       bool
	}{
		{
			name: "New index creation",
			args: args{
				ctx:       context.Background(),
				indexPath: "non_existing_index_path",
			},
			prepIndexMock: func() bleve.Index {
				indexMapping := bleve.NewIndexMapping()
				index, _ := bleve.NewMemOnly(indexMapping)
				return index
			},
			wantErr: false,
		},
		{
			name: "Open existing index",
			args: args{
				ctx:       context.Background(),
				indexPath: "existing_index_path",
			},
			prepIndexMock: func() bleve.Index {
				indexMapping := bleve.NewIndexMapping()
				index, _ := bleve.NewMemOnly(indexMapping)
				return index
			},
			wantErr: false,
		},
		{
			name: "Invalid index path",
			args: args{
				ctx:       context.Background(),
				indexPath: "/invalid/index/path",
			},
			prepIndexMock: func() bleve.Index {
				return nil // Simulate failure to create index
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := tt.prepIndexMock()
			idx := &indexer{
				index: index,
				log:   logger,
			}
			err := idx.NewIndex(tt.args.ctx, tt.args.indexPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			e := os.RemoveAll("non_existing_index_path")
			if e != nil {
				t.Fatal(e)
			}

			f := os.RemoveAll("existing_index_path")
			if f != nil {
				t.Fatal(f)
			}
		})
	}
}

func Test_indexer_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx  context.Context
		key  string
		data interface{}
	}

	testKey := "test_key"
	testData := "test_data"

	tests := []struct {
		name          string
		args          args
		prepIndexMock func() bleve.Index
		wantErr       bool
	}{
		{
			name: "Index success",
			args: args{
				ctx:  context.Background(),
				key:  testKey,
				data: testData,
			},
			prepIndexMock: func() bleve.Index {
				indexMapping := bleve.NewIndexMapping()
				index, _ := bleve.NewMemOnly(indexMapping)

				return index
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := tt.prepIndexMock()
			b := indexer{
				index: im,
			}
			err := b.Index(tt.args.ctx, tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexer.Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
