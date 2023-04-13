package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewOrganizationCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "organization",
		Short: "Operation for tks organization",
		Long:  `Operation for tks organization`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("organization called")
		},
	}

	command.AddCommand(NewOrganizationListCommand(globalOpts))
	command.AddCommand(NewOrganizationCreateCommand(globalOpts))

	return command
}
