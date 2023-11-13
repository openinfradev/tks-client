package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewMyProfileCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "my-profile",
		Short: "Operation for my profile",
		Long:  `Operation for my profile`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("my-profile called")
		},
	}

	command.AddCommand(NewMyProfileUpdatePasswordCommand(globalOpts))

	return command
}
