package local_storage

import (
	"context"
	"fmt"
	"os"

	goerr "errors"

	"github.com/blevesearch/bleve"
	"github.com/downsized-devs/sdk-go/logger"
)

type Interface interface {
	Index(ctx context.Context, key string, data interface{}) error
	Search(ctx context.Context, query string) ([]string, error)
	DeleteIndex(ctx context.Context, indexDir string) error
	NewIndex(ctx context.Context, indexPath string) error
}

type Config struct {
	IndexPath string
}

type indexer struct {
	conf  Config
	log   logger.Interface
	index bleve.Index
}

func Init(cfg Config, log logger.Interface) Interface {

	i := &indexer{
		conf: cfg,
		log:  log,
	}

	return i
}

func (i *indexer) NewIndex(ctx context.Context, indexPath string) error {
	// Open or create the Bleve index at the specified indexPath
	index, err := bleve.Open(indexPath)
	if err != nil {
		if goerr.Is(err, bleve.ErrorIndexPathDoesNotExist) {
			// If the index path does not exist, create a new index
			indexMapping := bleve.NewIndexMapping()
			index, err = bleve.New(indexPath, indexMapping)
			if err != nil {
				i.log.Error(ctx, fmt.Sprintf("Error creating bleve index: %v", err))
				return err
			}
		} else {
			i.log.Error(ctx, fmt.Sprintf("Error opening bleve index: %v", err))
			return err
		}
	}

	i.index = index
	return nil

}

func (idx *indexer) Index(ctx context.Context, key string, data interface{}) error {
	if err := idx.index.Index(key, data); err != nil {
		return err
	}

	return nil
}

func (idx *indexer) Search(ctx context.Context, query string) ([]string, error) {
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchResult, err := idx.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var results []string

	for _, hit := range searchResult.Hits {
		// Append hit ID to results
		results = append(results, hit.ID)
	}

	return results, nil
}

// DeleteIndex deletes the specified index directory.
func (idx *indexer) DeleteIndex(ctx context.Context, indexDir string) error {
	// Remove the index directory
	err := os.RemoveAll(indexDir)
	if err != nil {
		return err
	}

	return nil
}
