package localstorage

import (
	"context"
	"os"
	"slices"
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
		{
			name: "Index with uninitialized index field",
			args: args{
				ctx:  context.Background(),
				key:  testKey,
				data: testData,
			},
			prepIndexMock: func() bleve.Index {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := tt.prepIndexMock()
			b := &indexer{
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

func Test_indexer_Index_NilReceiver(t *testing.T) {
	var idx *indexer
	err := idx.Index(context.Background(), "k", "v")
	if err == nil {
		t.Fatal("expected error when receiver is nil, got nil")
	}
	if err.Error() != "indexer is nil" {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

func Test_indexer_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx   context.Context
		query string
	}

	tests := []struct {
		name          string
		args          args
		prepIndexMock func() bleve.Index
		seed          func(bleve.Index)
		wantErr       bool
		wantHitID     string
	}{
		{
			name: "Search success returns hit IDs",
			args: args{
				ctx:   context.Background(),
				query: "alpha",
			},
			prepIndexMock: func() bleve.Index {
				indexMapping := bleve.NewIndexMapping()
				index, _ := bleve.NewMemOnly(indexMapping)
				return index
			},
			seed: func(idx bleve.Index) {
				_ = idx.Index("doc-1", map[string]string{"name": "alpha"})
			},
			wantErr:   false,
			wantHitID: "doc-1",
		},
		{
			name: "Search with uninitialized index field",
			args: args{
				ctx:   context.Background(),
				query: "alpha",
			},
			prepIndexMock: func() bleve.Index {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := tt.prepIndexMock()
			if tt.seed != nil && im != nil {
				tt.seed(im)
			}
			b := &indexer{index: im}
			results, err := b.Search(tt.args.ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexer.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantHitID != "" && !slices.Contains(results, tt.wantHitID) {
				t.Errorf("indexer.Search() results = %v, want hit %q", results, tt.wantHitID)
			}
		})
	}
}
