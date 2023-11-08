package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterNodeGetCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		clusterId string
		long      bool
	)

	var command = &cobra.Command{
		Use:   "get",
		Short: "Show information of cluster node.",
		Long: `Show information of cluster node.
	
	Example:
	tks cluster node get`,
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			api := fmt.Sprintf("clusters/%s/nodes", clusterId)
			body, err := apiClient.Get(api)
			if err != nil {
				return err
			}

			var out domain.GetClusterNodesResponse
			helper.Transcode(body, &out)

			printClusterNodes(out.Nodes, long)

			return nil
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "the clusterId for nodes")
	helper.CheckError(command.MarkFlagRequired("cluster-id"))

	command.Flags().BoolVarP(&long, "long", "l", false, "enabled detail information")

	return command
}

func printClusterNodes(r []domain.ClusterNode, long bool) {
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

	t.AppendHeader(table.Row{"TYPE", "TARGETED", "REGISTERED", "REGISTERING", "STATUS", "COMMAND", "VALIDITY"})
	for _, s := range r {
		if !long {
			if s.Status == "DELETED" {
				continue
			}
		}
		t.AppendRow(table.Row{s.Type, s.Targeted, s.Registered, s.Registering, s.Status, s.Command, s.Validity})
	}

	if len(r) > 0 {
		fmt.Println(t.Render())
	} else {
		fmt.Println("No node found.")
	}
}
