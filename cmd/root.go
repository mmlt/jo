package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd returns the parent of all jo commands.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jo",
		Short: "json operators brings operations like add, subtract, do to json. jo is a friend of jq.",
	}
	cmd.AddCommand(NewDoCmd())
	cmd.AddCommand(NewVersionCmd())
	return cmd
}
