package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewClusterCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "cluster",
		Short: "Operation for tks cluster",
		Long:  `Operation for tks cluster`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cluster called")
		},
	}

	command.AddCommand(NewClusterListCommand(globalOpts))
	command.AddCommand(NewClusterCreateCommand(globalOpts))
	command.AddCommand(NewClusterDeleteCommand(globalOpts))

	return command
}
