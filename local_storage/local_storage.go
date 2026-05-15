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

// NewIndex opens or creates a bleve index at indexPath.
//
// SECURITY: indexPath is passed straight to bleve.Open / bleve.New, which
// creates a directory on disk. Callers MUST validate indexPath against a
// known-good root (e.g. Config.IndexPath) before invoking this method —
// passing user-controlled input enables arbitrary filesystem writes
// wherever this process has permission.
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
	if idx == nil {
		return fmt.Errorf("indexer is nil")
	}

	if idx.index == nil {
		return fmt.Errorf("bleve index is not initialized: call NewIndex before Index")
	}

	if err := idx.index.Index(key, data); err != nil {
		return err
	}

	return nil
}

func (idx *indexer) Search(ctx context.Context, query string) ([]string, error) {
	if idx.index == nil {
		return nil, fmt.Errorf("index is not initialized")
	}

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

// DeleteIndex recursively removes the directory at indexDir.
//
// SECURITY: indexDir is passed straight to os.RemoveAll with no validation.
// Callers MUST ensure indexDir is constrained to a known-good root (e.g.
// Config.IndexPath) before calling — passing user-controlled input enables
// arbitrary recursive deletion anywhere this process can write. This
// method does not enforce path containment; that is the caller's
// responsibility.
func (idx *indexer) DeleteIndex(ctx context.Context, indexDir string) error {
	// Remove the index directory
	err := os.RemoveAll(indexDir)
	if err != nil {
		return err
	}

	return nil
}
