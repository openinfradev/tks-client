package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewUserCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "user",
		Short: "Operation for user of organization",
		Long:  `Operation for user of organization`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("user called")
		},
	}

	command.AddCommand(NewUserCreateCommand(globalOpts))
	command.AddCommand(NewUserListCommand(globalOpts))
	command.AddCommand(NewUserDeleteCommand(globalOpts))

	return command
}
