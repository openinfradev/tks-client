package commands

import (
	"errors"
	"fmt"
	"log"
	"os"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterImportCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		name            string
		clusterType     string
		organizationId  string
		description     string
		stackTemplateId string
		kubeconfigPath  string
	)

	var command = &cobra.Command{
		Use:   "import",
		Short: "Import a TKS Cluster.",
		Long: `Import a TKS Cluster.
	  
	Example:
	tks cluster import <CLUSTERNAME> [--cloud-service AWS] [--template TEMPLATE_NAME]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println("You must specify cluster name.")
				return errors.New("Usage: tks cluster import <CLUSTERNAME> --contract-id <CONTRACTID>")
			}

			if len(args) == 1 {
				name = args[0]
			}

			kubeconfig, err := os.ReadFile(kubeconfigPath)
			if err != nil {
				log.Fatalf("Failed to read kubeconfig from [%s] path", err)
				log.Fatalf("Failed to read kubeconfig from [%s] path", kubeconfigPath)
			}
			input := domain.ImportClusterRequest{
				OrganizationId:  organizationId,
				StackTemplateId: stackTemplateId,
				Name:            name,
				Description:     description,
				ClusterType:     clusterType,
				Kubeconfig:      kubeconfig,
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Post("clusters/import", input)
			if err != nil {
				return err
			}

			var out domain.ImportClusterResponse
			helper.Transcode(body, &out)

			fmt.Println("clusterId : ", out.ID)

			return nil
		},
	}
	command.Flags().StringVarP(&organizationId, "organization-id", "o", "", "the organizationId with clusters")
	helper.CheckError(command.MarkFlagRequired("organization-id"))

	command.Flags().StringVar(&clusterType, "cluster-type", "USER", "the cluster type (USER | ADMIN)")

	command.Flags().StringVarP(&stackTemplateId, "stack-template-id", "t", "", "the template for installation")
	helper.CheckError(command.MarkFlagRequired("stack-template-id"))

	command.Flags().StringVarP(&name, "name", "n", "", "the name of organization")
	command.Flags().StringVarP(&description, "description", "d", "", "the description of organization")

	command.Flags().StringVar(&kubeconfigPath, "kubeconfig-path", "~/.kube/config", "the path of kubeconfig")

	return command
}
