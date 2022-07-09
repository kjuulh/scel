package graphql

import (
	"context"
	graphql1 "github.com/kjuulh/scel/server/internal/generated/graphql"
)

func (r *resolver) Mutation() graphql1.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *resolver }

// AddDownload implements graphql.MutationResolver
func (r *mutationResolver) AddDownload(ctx context.Context, download graphql1.CreateDownload) (*graphql1.Download, error) {
	return r.App.DownloadsStore.AddDownload(ctx, &download)
}
