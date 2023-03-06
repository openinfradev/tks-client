package commands

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewOrganizationListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of organization.",
		Long: `Show list of organization.
	
	Example:
	tks organization list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Get("organizations")
			if err != nil {
				return err
			}

			type DataInterface struct {
				Organizations []domain.Organization `json:"organizations"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			printOrganizations(out.Organizations, true)

			fmt.Println("Success")
			return nil
		},
	}

	return command
}

/*
var organizationListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list of organization.",
	Long: `Show list of organization.

Example:
tks organization list`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("server ")

			apiClient, err := _apiClient.New(config.GetServer().Server, config.GetUser().AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Get("organizations")
			if err != nil {
				return err
			}

			type DataInterface struct {
				Organizations []domain.Organization `json:"organizations"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			printOrganizations(out.Organizations, true)

			fmt.Println("Success")
		return nil
	},
}
*/

func printOrganizations(r []domain.Organization, long bool) {
	t := table.NewWriter()
	tTemp := table.Table{}
	tTemp.Render()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	if long {
		t.AppendHeader(table.Row{"ORGANIZATION_ID", "NAME", "DESCRIPTION", "CREATED_AT", "UPDATED_AT"})
		for _, s := range r {
			tCreatedAt := parseTime(s.CreatedAt)
			tUpdatedAt := parseTime(s.UpdatedAt)

			t.AppendRow(table.Row{s.ID, s.Name, s.Description, tCreatedAt, tUpdatedAt})
		}
	} else {
		t.AppendHeader(table.Row{"TYPE", "SERVICE_ID", "STATUS", "CREATED_AT", "UPDATED_AT"})
		for _, s := range r {
			tCreatedAt := parseTime(s.CreatedAt)
			tUpdatedAt := parseTime(s.UpdatedAt)

			t.AppendRow(table.Row{s.ID, s.Name, s.Description, tCreatedAt, tUpdatedAt})
		}
	}

	if len(r) > 0 {
		fmt.Println(t.Render())
	} else {
		fmt.Println("No organization found.")
	}
}
func parseTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
