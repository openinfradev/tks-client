package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/config"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewCloudAccountListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		organizationId string
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of cloud-account.",
		Long: `Show list of cloud-account.
	
	Example:
	tks cloud-account list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && organizationId == "" {
				localCfg, err := config.ReadLocalConfig(globalOpts.ConfigPath)
				if err == nil {
					organizationId = localCfg.GetUser().OrganizationId
				}
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Get("/organizations/" + organizationId + "/cloud-accounts")
			if err != nil {
				return err
			}

			var out = domain.GetCloudAccountsResponse{}
			helper.Transcode(body, &out)

			printCloudAccounts(out.CloudAccounts)

			return nil
		},
	}

	command.Flags().StringVarP(&organizationId, "organization-id", "o", "", "the organizationId with clusters")

	return command
}

func printCloudAccounts(r []domain.CloudAccountResponse) {
	if len(r) == 0 {
		fmt.Println("No cloudAccount exists for user organization!")
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
	t.AppendHeader(table.Row{"ORGANIZATION_ID", "ID", "NAME", "DESCRIPTION", "CLOUD_SERVICE", "RESOURCE", "CLUSTERS", "CREATED_AT", "UPDATED_AT"})
	for _, s := range r {
		tCreatedAt := helper.ParseTime(s.CreatedAt)
		tUpdatedAt := helper.ParseTime(s.UpdatedAt)
		t.AppendRow(table.Row{s.OrganizationId, s.ID, s.Name, s.Description, s.CloudService, s.Resource, s.Clusters, tCreatedAt, tUpdatedAt})
	}
	fmt.Println(t.Render())
}
