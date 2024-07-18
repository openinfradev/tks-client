package commands

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterImportCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		name             string
		clusterType      string
		organizationId   string
		description      string
		stackTemplateId  string
		kubeconfigPath   string
		kubeconfigString string
		policyIds        []string
		domains          []string
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

			var kubeconfig string
			if kubeconfigPath != "" {
				val, err := os.ReadFile(kubeconfigPath)
				if err != nil {
					log.Fatalf("Failed to read kubeconfig from [%s] path", err)
					log.Fatalf("Failed to read kubeconfig from [%s] path", kubeconfigPath)
				}
				kubeconfig = b64.StdEncoding.EncodeToString(val)
			} else if kubeconfigString != "" {
				kubeconfig = kubeconfigString
			} else {
				log.Fatalf("One of kubeconfigPath and kubeconfigString must be filled")
			}

			clusterDomains := make([]domain.ClusterDomain, len(domains))
			for i, domain := range domains {
				arrDomain := strings.Split(domain, "_")
				if len(arrDomain) > 0 {
					clusterDomains[i].DomainType = arrDomain[0]
					clusterDomains[i].Url = arrDomain[1]
				}
			}

			input := domain.ImportClusterRequest{
				OrganizationId:  organizationId,
				StackTemplateId: stackTemplateId,
				Name:            name,
				Description:     description,
				ClusterType:     clusterType,
				Kubeconfig:      kubeconfig,
				CloudService:    "BYOK",
				PolicyIds:       policyIds,
				Domains:         clusterDomains,
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

	command.Flags().StringVar(&kubeconfigPath, "kubeconfig-path", "", "the path of kubeconfig")
	command.Flags().StringVar(&kubeconfigString, "kubeconfig-string", "", "the contents of kubeconfig")

	command.Flags().StringSliceVar(&policyIds, "policy-ids", []string{}, "ex. policy_id1,policy_id1")

	command.Flags().StringSliceVar(&domains, "domains", []string{}, "ex. grafana_1.1.1.1:30001,thanos_1.1.1.1:30002")

	return command
}
