package graphql

import (
	"context"

	graphql1 "github.com/kjuulh/scel/server/internal/generated/graphql"
)

// Query returns graphql1.QueryResolver implementation.
func (r *resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *resolver }

func (r *queryResolver) Downloads(ctx context.Context, userID string) ([]*graphql1.Download, error) {
	downloads, err := r.App.DownloadsStore.GetDownloads(ctx)
	if err != nil {
		return nil, err
	}

	return downloads, nil
}
