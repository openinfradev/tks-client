package commands

import (
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewAppGroupDeleteCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		appGroupId string
	)

	var command = &cobra.Command{
		Use:   "delete",
		Short: "Delete a AppGroup in cluster.",
		Long: `Delete a AppGroup in cluster.
	  
	Example:
	tks appgroup delete <APP_GROUP_ID>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				appGroupId = args[0]
			}

			if appGroupId == "" {
				helper.PanicWithError("You must specify appGroupId")
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			_, err = apiClient.Delete("app-groups/"+appGroupId, nil)
			if err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().StringVarP(&appGroupId, "appgroup-id", "a", "", "the Id of appGroup")
	return command
}
