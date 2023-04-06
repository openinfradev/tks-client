package commands

import (
	"fmt"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewOrganizationCreateCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		name        string
		description string
	)

	var command = &cobra.Command{
		Use:   "create",
		Short: "Create a tks organization",
		Long: `Create a tks organization
	
	Example:
	tks organization create <ORGANIZATION NAME>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}
			if name == "" {
				helper.PanicWithError("You must specify name")
			}

			input := domain.CreateOrganizationRequest{
				Name:        name,
				Description: description,
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)
			body, err := apiClient.Post("organizations", input)
			if err != nil {
				return err
			}

			type DataInterface struct {
				OrganizationId string `json:"organizationId"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			fmt.Println("Success: The request to create organization ", body, " was accepted.")

			return nil
		},
	}
	command.Flags().StringVar(&name, "name", "", "the name of organization")
	command.Flags().StringVar(&description, "description", "", "the description of organization")

	return command
}
