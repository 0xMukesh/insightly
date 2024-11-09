package commands

import "github.com/spf13/cobra"

type BaseCmd struct {
	Cmd  *cobra.Command
	Args []string
}
