package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCloudSettingCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "cloud-setting",
		Short: "Operation for cloud setting",
		Long:  `Operation for cloud setting`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cloud-setting called")
		},
	}

	command.AddCommand(NewCloudSettingListCommand(globalOpts))

	return command
}
