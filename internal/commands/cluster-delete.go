package commands

import (
	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterDeleteCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		clusterId string
	)

	var command = &cobra.Command{
		Use:   "delete",
		Short: "Delete a TKS Cluster.",
		Long: `Delete a TKS Cluster.
	  
	Example:
	tks cluster delete <CLUSTERNAME>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				clusterId = args[0]
			}

			if clusterId == "" {
				helper.PanicWithError("You must specify clusterId")
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			_, err = apiClient.Delete("clusters/"+clusterId, nil)
			if err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "the Id for deleting cluster")
	return command
}
