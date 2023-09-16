package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the semantic version set during go build.
var Version string

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the jo version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
}
