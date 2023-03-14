package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAppGroupCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "appgroup",
		Short: "Operation for app group of cluster",
		Long:  `Operation for app group of cluster`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("appgroup called")
		},
	}

	command.AddCommand(NewAppGroupCreateCommand(globalOpts))
	command.AddCommand(NewAppGroupListCommand(globalOpts))
	command.AddCommand(NewAppGroupDeleteCommand(globalOpts))

	return command
}
