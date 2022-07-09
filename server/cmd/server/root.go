package server

import (
	"context"
	"log"

	"github.com/kjuulh/scel/server/internal/app"
	"github.com/kjuulh/scel/server/internal/graphql"
	"github.com/kjuulh/scel/server/internal/observability"
	"github.com/spf13/cobra"
)

func NewServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the scel server",
		Run: func(cmd *cobra.Command, _ []string) {
			println("serving scel")

			a := app.NewApp()

			doneChan := make(chan bool, 1)

			shutdown, err := observability.NewOtlp(context.Background())
			if err != nil {
				panic(err)
			}
			defer shutdown()

			err = graphql.ServeGraphQL(a, doneChan)
			if err != nil {
				log.Fatalf("could not serve graphql: %w+", err)
			}

			<-doneChan
		},
	}
}

func RegisterCommand(cmd *cobra.Command) {
	cmd.AddCommand(NewServerCmd())
}
