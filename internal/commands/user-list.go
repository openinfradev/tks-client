package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewUserListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of users.",
		Long: `Show list of Tks users.
	
	Example:
	tks user list `,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/users", globalOpts.CurrentOrganizationId)
			body, err := apiClient.Get(url)
			if err != nil {
				return err
			}

			type DataInterface struct {
				Users []domain.UserResponse `json:"users"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			printUsers(out.Users)

			return nil
		},
	}

	return command
}

func printUsers(d []domain.UserResponse) {
	if len(d) == 0 {
		fmt.Println("No User exists for specified organization!")
		return
	}

	t := table.NewWriter()
	tTemp := table.Table{}
	tTemp.Render()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"ORGANIZATION_ID", "ACCOUNT_ID", "NAME", "ID", "email", "Department", "Description", "CREATED_AT", "UPDATED_AT"})
	for _, i := range d {
		tCreatedAt := helper.ParseTime(i.CreatedAt)
		tUpdatedAt := helper.ParseTime(i.UpdatedAt)
		t.AppendRow(table.Row{i.Organization.ID, i.AccountId, i.Name, i.ID, i.Email, i.Department, i.Description, tCreatedAt, tUpdatedAt})
	}
	fmt.Println(t.Render())
}
