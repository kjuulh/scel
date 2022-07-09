package graphql

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/kjuulh/scel/server/internal/app"
	graphql1 "github.com/kjuulh/scel/server/internal/generated/graphql"
	"github.com/ravilushqa/otelgqlgen"
)

type resolver struct {
	App *app.App
}

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

		router := chi.NewRouter()

		srv := handler.NewDefaultServer(
			graphql1.NewExecutableSchema(
				graphql1.Config{
					Resolvers:  &resolver{App: a},
					Directives: graphql1.DirectiveRoot{},
					Complexity: graphql1.ComplexityRoot{},
				},
			),
		)
		srv.AddTransport(transport.POST{})
		srv.AddTransport(transport.Websocket{
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
				return websocketInit(ctx, initPayload)
			},
		})
		srv.Use(extension.Introspection{})
		srv.Use(otelgqlgen.Middleware())

		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
		router.Handle("/query", srv)
		log.Printf("connect to http://0.0.0.0:%d/ for graphql playground", port)

		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router)

		doneChan <- true
	}(a, graphQlPort, doneChan)

	return nil
}

func websocketInit(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
	// Get the token from payload
	//any := initPayload["authToken"]
	//token, ok := any.(string)
	//if !ok || token == "" {
	//	return nil, errors.New("authToken not found in transport payload")
	//}

	//// Perform token verification and authentication...
	//userId := "john.doe" // e.g. userId, err := GetUserFromAuthentication(token)

	//// put it in context
	//ctxNew := context.WithValue(ctx, "username", userId)

	//return ctxNew, nil
	return ctx, nil
}
