package commands

import (
	"errors"
	"fmt"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

type rollback struct {
	TaskId string `yaml:"task_id" json:"taskId"`
}

func NewAppserveRollbackCmd(globalOpts *GlobalOptions) *cobra.Command {
	var (
		organizationId string
		appId          string
	)
	var command = &cobra.Command{
		Use:   "rollback",
		Short: "Rollback an app to revision",
		Long: `Rollback an app to revision.
  
	Example:
	tks appserve rollback <TASK_ID> --organization-id <ORG_ID> --app-id <APP_ID>`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			taskId := args[0]
			if len(taskId) < 1 {
				return errors.New("TASK_ID is mandatory! Run 'tks appserve rollback --help'")
			}

			if organizationId == "" {
				return errors.New("--organization-id is mandatory param")
			}
			if appId == "" {
				return errors.New("--app-id is mandatory param")
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/app-serve-apps/%v/rollback", organizationId, appId)
			params := rollback{TaskId: taskId}

			body, err := apiClient.Post(url, params)
			if err != nil {
				return err
			}

			fmt.Printf("Response: %T Type\n %v\n", body, fmt.Sprintf("%v", body))

			return nil
		},
	}

	command.Flags().StringVar(&organizationId, "organization-id", "", "a organizationId")
	command.Flags().StringVar(&appId, "app-id", "", "a appId")

	return command
}
