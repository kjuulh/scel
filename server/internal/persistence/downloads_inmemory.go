package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/kjuulh/scel/server/internal/generated/graphql"
)

var _ DownloadsStore = &inMemoryDownloadsStore{}

func NewInMemoryDownloadsStore() DownloadsStore {
	return &inMemoryDownloadsStore{
		downloads: make(map[string]*graphql.Download),
	}
}

type inMemoryDownloadsStore struct {
	downloads map[string]*graphql.Download
}

// GetDownloads implements DownloadsStore
func (s *inMemoryDownloadsStore) GetDownloads(ctx context.Context) ([]*graphql.Download, error) {
	downloads := make([]*graphql.Download, 0)

	for _, d := range s.downloads {
		downloads = append(downloads, d)
	}

	return downloads, nil
}

// AddDownload implements DownloadsStore
func (s *inMemoryDownloadsStore) AddDownload(ctx context.Context, createDownload *graphql.CreateDownload) (*graphql.Download, error) {
	id, err := s.generateID(ctx)
	if err != nil {
		return nil, err
	}

	download := &graphql.Download{
		ID:     id,
		UserID: createDownload.UserID,
		Link:   createDownload.Link,
	}

	s.downloads[download.ID] = download

	return download, nil
}

func (s *inMemoryDownloadsStore) generateID(ctx context.Context) (string, error) {
	u := uuid.NewString()

	return u, nil
}
