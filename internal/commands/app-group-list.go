package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewAppGroupListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		clusterId string
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of cluster.",
		Long: `Show list of clusterrganization.
	
	Example:
	tks cluster list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				clusterId = args[0]
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			api := fmt.Sprintf("app-groups?clusterId=%s", clusterId)
			body, err := apiClient.Get(api)
			if err != nil {
				return err
			}

			var out domain.GetAppGroupsResponse
			helper.Transcode(body, &out)

			printAppGroups(out.AppGroups)

			return nil
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "the clusterId")

	return command
}

func printAppGroups(r []domain.AppGroupResponse) {
	if len(r) == 0 {
		fmt.Println("No appGroup exists for specified cluster!")
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
	t.AppendHeader(table.Row{"CLUSTER_ID", "NAME", "ID", "TYPE", "STATUS", "CREATED_AT", "UPDATED_AT"})
	for _, s := range r {
		tCreatedAt := helper.ParseTime(s.CreatedAt)
		tUpdatedAt := helper.ParseTime(s.UpdatedAt)
		t.AppendRow(table.Row{s.ClusterId, s.Name, s.ID, s.AppGroupType, s.Status, tCreatedAt, tUpdatedAt})
	}
	fmt.Println(t.Render())
}
