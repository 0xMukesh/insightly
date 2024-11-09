package cmd

import (
	"context"

	"github.com/0xmukesh/ratemywebsite/internal/commands"
	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Version: "0.0.1",
		Use:     "something",
		Long:    "something is something",
		Example: "something",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	genUxCmd := commands.GenerateUxReportCmd{}
	setupCmd := commands.SetupCmd{}

	rootCmd.AddCommand(genUxCmd.New())
	rootCmd.AddCommand(setupCmd.New())

	return rootCmd.ExecuteContext(context.Background())
}
