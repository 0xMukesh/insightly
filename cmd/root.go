package cmd

import (
	"context"

	"github.com/0xmukesh/insightly/internal/commands"
	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Version: "0.0.1",
		Use:     "insightly",
		Long:    "`insightly` is a CLI that analyzes website's a11y by combining a11y auditing tools with AI powered insights",
		Example: "insightly",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	genUxCmd := commands.GenerateUxReportCmd{}
	setupCmd := commands.SetupCmd{}
	configCmd := commands.ConfigCmd{}

	rootCmd.AddCommand(genUxCmd.New())
	rootCmd.AddCommand(setupCmd.New())
	rootCmd.AddCommand(configCmd.New())

	return rootCmd.ExecuteContext(context.Background())
}
