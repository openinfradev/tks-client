package commands

import (
	"fmt"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewCloudSettingListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		all bool
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of cloud-setting.",
		Long: `Show list of cloud-setting.
	
	Example:
	tks cloud-setting list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Get("cloud-settings?all=" + strconv.FormatBool(all))
			if err != nil {
				return err
			}

			var out = domain.GetCloudSettingsResponse{}
			helper.Transcode(body, &out)

			printCloudSettings(out.CloudSettings)

			return nil
		},
	}

	command.Flags().BoolVarP(&all, "all", "A", false, "show all organizations")

	return command
}

func printCloudSettings(r []domain.CloudSettingResponse) {
	if len(r) == 0 {
		fmt.Println("No cloudSetting exists for user organization!")
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
	t.AppendHeader(table.Row{"ORGANIZATION_ID", "ID", "NAME", "DESCRIPTION", "TYPE", "RESOURCE", "CLUSTERS", "CREATED_AT", "UPDATED_AT"})
	for _, s := range r {
		tCreatedAt := helper.ParseTime(s.CreatedAt)
		tUpdatedAt := helper.ParseTime(s.UpdatedAt)
		t.AppendRow(table.Row{s.OrganizationId, s.ID, s.Name, s.Description, s.Type, s.Resource, s.Clusters, tCreatedAt, tUpdatedAt})
	}
	fmt.Println(t.Render())
}
