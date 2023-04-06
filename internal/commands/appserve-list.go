package commands

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewAppServeListCmd(globalOpts *GlobalOptions) *cobra.Command {
	var (
		organizationId string
		showAll        string
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of AppServing service.",
		Long: `Show list of AppServing service.
	
	Example:
	tks appserve list --organization-id <ORG_ID> --show-all <false|true>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if organizationId == "" {
				return errors.New("--organization-id is mandatory param")
			} else if showAll == "" {
				showAll = "true"
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			api := fmt.Sprintf("app-serve-apps?organizationId=%s&showAll=%s", organizationId, showAll)
			body, err := apiClient.Get(api)
			if err != nil {
				return err
			}

			type DataInterface struct {
				AppServeApps []domain.AppServeApp `json:"app_serve_apps"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			printAppServeApps(out.AppServeApps)

			return nil
		},
	}

	command.Flags().StringVar(&organizationId, "organization-id", "", "a organizationId")
	command.Flags().StringVar(&showAll, "show-all", "false", "a organizationId")

	return command
}

func printAppServeApps(d []domain.AppServeApp) {
	if len(d) == 0 {
		fmt.Println("No AppServeApp exists for specified organization!")
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
	t.AppendHeader(table.Row{"ORGANIZATION_ID", "NAME", "ID", "Type", "AppType", "ENDPOINT_URL", "CLUSTER_ID", "STATUS", "CREATED_AT", "UPDATED_AT"})
	for _, i := range d {
		tCreatedAt := helper.ParseTime(i.CreatedAt)
		var tUpdatedAt string
		if i.UpdatedAt != nil {
			tUpdatedAt = helper.ParseTime(*i.UpdatedAt)
		}
		t.AppendRow(table.Row{i.OrganizationId, i.Name, i.ID, i.Type, i.AppType, i.EndpointUrl, i.TargetClusterId, i.Status, tCreatedAt, tUpdatedAt})
	}
	fmt.Println(t.Render())
}
