package commands

import (
	"errors"
	"fmt"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewAppserveAbortCmd(globalOpts *GlobalOptions) *cobra.Command {
	var organizationId string
	var command = &cobra.Command{
		Use:   "abort",
		Short: "Abort an app paused in blue-green state",
		Long: `Abort an app paused in blue-green state.
  
	Example:
	tks appserve abort <APP_ID> --organization-id <ORG_ID>`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			appId := args[0]
			if len(appId) < 1 {
				return errors.New("APP_ID is mandatory! Run 'tks appserve abort --help'")
			}

			if organizationId == "" {
				return errors.New("--organization-id is mandatory param")
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/app-serve-apps/%v", organizationId, appId)
			params := conf{Strategy: "blue-green", Abort: true}

			body, err := apiClient.Put(url, params)
			if err != nil {
				return err
			}

			fmt.Printf("Response: %T Type\n %v\n", body, fmt.Sprintf("%v", body))

			return nil
		},
	}

	command.Flags().StringVar(&organizationId, "organization-id", "", "a organizationId")

	return command
}
