package commands

import (
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterInstallCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		clusterId string
	)

	var command = &cobra.Command{
		Use:   "install",
		Short: "Install a TKS Cluster.",
		Long: `Install a TKS Cluster.
	  
	Example:
	tks cluster install <CLUSTER_ID>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				clusterId = args[0]
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			_, err = apiClient.Post("clusters/"+clusterId+"/install", nil)
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
