package graphql

import (
	"context"

	graphql1 "github.com/kjuulh/scel/server/internal/generated/graphql"
	"github.com/kjuulh/scel/server/internal/observability"
)

func (r *resolver) Subscription() graphql1.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *resolver }

// SubscribeDownloads implements graphql.SubscriptionResolver
func (r *subscriptionResolver) SubscribeDownloads(ctx context.Context, userID string) (<-chan []*graphql1.Download, error) {
	downloadsChan := make(chan []*graphql1.Download, 1)

	initialDownloads, err := r.App.DownloadsStore.GetDownloads(ctx)
	observability.Logger(ctx).Error("downloads get")
	if err != nil {
		observability.Logger(ctx).Error("could not get downloads")
	}
	downloadsChan <- initialDownloads

	return downloadsChan, nil
}
