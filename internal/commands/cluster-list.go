package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		organizationId string
	)

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

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			api := fmt.Sprintf("clusters?organizationId=%s", organizationId)
			body, err := apiClient.Get(api)
			if err != nil {
				return err
			}

			var out domain.GetClustersResponse
			helper.Transcode(body, &out)

			printClusters(out.Clusters)

			return nil
		},
	}

	command.Flags().StringVarP(&organizationId, "organization-id", "o", "", "the organizationId with clusters")

	return command
}

func printClusters(r []domain.ClusterResponse) {
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
	t.AppendHeader(table.Row{"ORGANIZATION_ID", "NAME", "ID", "STATUS", "CREATED_AT", "UPDATED_AT"})
	for _, s := range r {
		tCreatedAt := helper.ParseTime(s.CreatedAt)
		tUpdatedAt := helper.ParseTime(s.UpdatedAt)
		t.AppendRow(table.Row{s.OrganizationId, s.Name, s.ID, s.Status, tCreatedAt, tUpdatedAt})
	}
	fmt.Println(t.Render())
}

func ModelToJson(in any) string {
	a, _ := json.Marshal(in)
	n := len(a)        //Find the length of the byte array
	s := string(a[:n]) //convert to string
	return s
}
