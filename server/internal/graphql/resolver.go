package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kjuulh/scel/server/internal/app"
	graphql1 "github.com/kjuulh/scel/server/internal/generated/graphql"
	"github.com/kjuulh/scel/server/internal/observability"
	"github.com/ravilushqa/otelgqlgen"
	"go.uber.org/zap"
)

type resolver struct {
	App *app.App
}

// // foo
func (r *queryResolver) Downloads(ctx context.Context, userID string) ([]*graphql1.Download, error) {
	ctx, tracer := r.App.Tracer.Start(ctx, "downloads")
	defer tracer.End()

	observability.Logger(ctx).Info("GetDownloads", zap.String("request", "some-request"))

	downloads := []*graphql1.Download{
		{
			ID:     "some-id",
			UserID: "some-user-id",
		},
	}

	return downloads, nil
}

// Query returns graphql1.QueryResolver implementation.
func (r *resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *resolver }

func ServeGraphQL(a *app.App, doneChan chan bool) error {
	port := os.Getenv("SCEL_GRAPHQL_PORT")
	if port == "" {
		port = "15000"
	}

	graphQlPort, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("SCEL_GRAPHQL_PORT is not an integer, error: %w", err)
	}

	go func(a *app.App, port int, doneChan chan bool) {
		srv := handler.NewDefaultServer(graphql1.NewExecutableSchema(graphql1.Config{Resolvers: &resolver{App: a}}))
		srv.Use(otelgqlgen.Middleware())

		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", srv)
		log.Printf("connect to http://0.0.0.0:%d/ for graphql playground", port)

		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)

		doneChan <- true
	}(a, graphQlPort, doneChan)

	return nil
}
