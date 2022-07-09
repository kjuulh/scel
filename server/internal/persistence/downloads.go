package persistence

import (
	"context"

	"github.com/kjuulh/scel/server/internal/generated/graphql"
)

type DownloadsStore interface {
	AddDownload(ctx context.Context, createDownload *graphql.CreateDownload) (*graphql.Download, error)
	GetDownloads(ctx context.Context) ([]*graphql.Download, error)
}
