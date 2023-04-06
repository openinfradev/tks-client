package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewStackTemplateCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "stack-template",
		Short: "Operation for stack templates",
		Long:  `Operation for stack templates`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("stack-template called")
		},
	}

	command.AddCommand(NewStackTemplateListCommand(globalOpts))

	return command
}
