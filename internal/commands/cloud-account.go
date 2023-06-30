package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCloudAccountCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "cloud-account",
		Short: "Operation for cloud setting",
		Long:  `Operation for cloud setting`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cloud-account called")
		},
	}

	command.AddCommand(NewCloudAccountCreateCommand(globalOpts))
	command.AddCommand(NewCloudAccountListCommand(globalOpts))

	return command
}
