package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewClusterNodeCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "node",
		Short: "Operation for tks cluster node",
		Long:  `Operation for tks cluster node`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cluster called")
		},
	}

	command.AddCommand(NewClusterNodeGetCommand(globalOpts))
	command.AddCommand(NewClusterNodeListCommand(globalOpts))

	return command

}
