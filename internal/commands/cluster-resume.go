package commands

import (
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterResumeCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		clusterId string
	)

	var command = &cobra.Command{
		Use:   "resume",
		Short: "Resume a TKS Cluster.",
		Long: `Resume a TKS Cluster.
	  
	Example:
	tks cluster resume <CLUSTER_ID>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				clusterId = args[0]
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			_, err = apiClient.Put("clusters/"+clusterId+"/resume", nil)
			if err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "the clusterId with clusters")
	helper.CheckError(command.MarkFlagRequired("cluster-id"))

	return command
}
