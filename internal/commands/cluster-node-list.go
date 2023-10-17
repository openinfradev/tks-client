package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterNodeListCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		clusterId string
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "Show list of cluster node.",
		Long: `Show list of cluster node.
	
	Example:
	tks cluster node list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			api := fmt.Sprintf("clusters/%s/nodes", clusterId)
			body, err := apiClient.Get(api)
			if err != nil {
				return err
			}

			var out domain.GetClusterNodesResponse
			helper.Transcode(body, &out)

			printClusterHosts(out.Nodes)

			return nil
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "the clusterId for nodes")
	helper.CheckError(command.MarkFlagRequired("cluster-id"))

	return command
}

func printClusterHosts(r []domain.ClusterNode) {
	if len(r) == 0 {
		fmt.Println("No cluster nodes exists for specified cluster!")
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

	t.AppendHeader(table.Row{"TYPE", "NAME", "STATUS"})
	for _, s := range r {

		for _, host := range s.Hosts {
			t.AppendRow(table.Row{s.Type, host.Name, host.Status})
		}
	}

	if len(r) > 0 {
		fmt.Println(t.Render())
	} else {
		fmt.Println("No host found.")
	}
}
