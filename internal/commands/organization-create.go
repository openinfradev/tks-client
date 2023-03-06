package commands

import (
	"fmt"

	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

var organizationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a tks organization",
	Long: `Create a tks organization

Example:
tks organization create <ORGANIZATION NAME>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("You must specify organization name.")
			return fmt.Errorf("Usage: tks organization create <ORGANIZATION NAME>")
		}
		fmt.Println("Organization Name: ", args[0])
		name := args[0]
		description, _ := cmd.Flags().GetString("description")
		creator, _ := cmd.Flags().GetString("creator")

		input := domain.CreateOrganizationRequest{
			Name:        name,
			Description: description,
			Creator:     creator,
		}

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
