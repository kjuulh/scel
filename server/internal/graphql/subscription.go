package graphql

import (
	"context"
	"time"

	"github.com/kjuulh/scel/server/internal/app"
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

	go func(a *app.App) {
		goCtx := context.Background()
		goCtx, span := a.Tracer.Start(goCtx, "subscribeDownloads")
		defer span.End()

		ticker := time.NewTicker(time.Second * 5)
		done := make(chan bool)

		for {
			select {
			case <-ctx.Done():
				observability.Logger(goCtx).Error("downloads Done")
				close(downloadsChan)
				return
			case <-done:
				observability.Logger(goCtx).Error("downloads error")
				close(downloadsChan)
				return
			case <-ticker.C:
				recurringDownloads, err := a.DownloadsStore.GetDownloads(goCtx)
				observability.Logger(goCtx).Info("downloads get")
				if err != nil {
					observability.Logger(goCtx).Error("could not get downloads")
				}
				downloadsChan <- recurringDownloads
			}
		}
	}(r.App)

	return downloadsChan, nil
}
