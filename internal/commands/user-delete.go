package commands

import (
	"errors"
	"fmt"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewUserDeleteCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		accountId string
	)

	var command = &cobra.Command{
		Use:   "delete",
		Short: "Delete users.",
		Long: `Delete users.
	
	Example:
	tks user delete --account-id <ACCOUNT_ID>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if globalOpts.CurrentOrganizationId == "" {
				return errors.New("current organization is not set")
			}
			if accountId == "" {
				return errors.New("account-id is not set")
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/users/%v", globalOpts.CurrentOrganizationId, accountId)
			_, err = apiClient.Delete(url, nil)
			if err != nil {
				return err
			}
			return nil
		},
	}

	command.Flags().StringVar(&accountId, "account-id", "", "user accountId")
	helper.CheckError(command.MarkFlagRequired("account-id"))

	return command
}
