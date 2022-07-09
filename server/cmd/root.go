package cmd

import (
	"fmt"
	"os"

	"github.com/kjuulh/scel/server/cmd/server"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "scel",
	}

	server.RegisterCommand(cmd)

	return cmd
}

func Execute() {
	rootCmd := NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
