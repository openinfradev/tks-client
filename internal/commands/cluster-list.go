package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var organizationId string

	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of cluster.",
		Long: `Show list of clusterrganization.
	
	Example:
	tks cluster list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				organizationId = args[0]
			}

			if organizationId == "" {
				helper.PanicWithError("You must specify organizationId")
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Get(fmt.Sprintf("clusters?organizationId=%s", organizationId))
			if err != nil {
				return err
			}

			type DataInterface struct {
				Clusters []domain.Cluster `json:"clusters"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			printClusters(out.Clusters, true)

			return nil
		},
	}

	command.Flags().StringVar(&organizationId, "organization-id", "", "the organizationId with clusters")

	return command
}

func printClusters(r []domain.Cluster, long bool) {
	if len(r) == 0 {
		fmt.Println("No cluster exists for specified organization!")
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
	if long {
		t.AppendHeader(table.Row{"Name", "ID", "Status", "CREATED_AT", "UPDATED_AT", "CONTRACT_ID", "STATUS_DESC"})
		for _, s := range r {
			tCreatedAt := helper.ParseTime(s.CreatedAt)
			tUpdatedAt := helper.ParseTime(s.UpdatedAt)
			t.AppendRow(table.Row{s.Name, s.ID, s.Status, tCreatedAt, tUpdatedAt, s.OrganizationId})
		}
	} else {
		t.AppendHeader(table.Row{"Name", "ID", "Status", "CREATED_AT", "UPDATED_AT"})
		for _, s := range r {
			tCreatedAt := helper.ParseTime(s.CreatedAt)
			tUpdatedAt := helper.ParseTime(s.UpdatedAt)
			t.AppendRow(table.Row{s.Name, s.ID, s.Status, tCreatedAt, tUpdatedAt})
		}
	}
	fmt.Println(t.Render())
}
