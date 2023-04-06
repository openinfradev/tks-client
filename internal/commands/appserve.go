package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewAppserveCommand(globalOpts *GlobalOptions) *cobra.Command {
	// appserveCmd represents the appserve command
	var command = &cobra.Command{
		Use:   "appserve",
		Short: "Operation for TKS Appserve",
		Long:  `Operation for TKS Appserve`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("appserve called")
		},
	}

	command.AddCommand(NewAppserveCreateCmd(globalOpts))
	command.AddCommand(NewAppServeListCmd(globalOpts))
	command.AddCommand(NewAppServeShowCmd(globalOpts))
	command.AddCommand(NewAppserveUpdateCmd(globalOpts))
	command.AddCommand(NewAppserveDeleteCmd(globalOpts))
	command.AddCommand(NewAppservePromoteCmd(globalOpts))
	command.AddCommand(NewAppserveAbortCmd(globalOpts))

	return command
}
