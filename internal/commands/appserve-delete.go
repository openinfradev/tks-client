package commands

import (
	"errors"
	"fmt"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewAppserveDeleteCmd(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "delete",
		Short: "delete app deployed by AppServing service",
		Long: `delete app deployed by AppServing service.
  
	Example:
	tks appserve delete <APP_UUID>`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			appId := args[0]

			if len(appId) < 1 {
				return errors.New("APP_ID is mandatory! Run 'tks appserve delete --help'")
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Delete("app-serve-apps/"+appId, nil)
			if err != nil {
				return err
			}

			fmt.Printf("Response: %T Type\n %v", body, fmt.Sprintf("%v", body))

			return nil
		},
	}

	return command
}
