package commands

import (
	"fmt"
	"strconv"
	"strings"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterCreateCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		name             string
		clusterType      string
		organizationId   string
		description      string
		stackTemplateId  string
		cloudService     string
		cloudAccountId   string
		stack            int
		tksCpNode        int
		tksCpNodeMax     int
		tksCpNodeType    string
		tksInfraNode     int
		tksInfraNodeMax  int
		tksInfraNodeType string
		tksUserNode      int
		tksUserNodeMax   int
		tksUserNodeType  string
		clusterEndpoint  string
		policyIds        []string
		domains          []string
	)

	var command = &cobra.Command{
		Use:   "create",
		Short: "Create a TKS Cluster.",
		Long: `Create a TKS Cluster.
	  
	Example:
	tks cluster create <CLUSTERNAME> [--cloud-service AWS] [--template TEMPLATE_NAME]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}

			if name == "" {
				helper.PanicWithError("You must specify name")
			}

			isStack := false
			if stack > 0 {
				isStack = true
			}

			byoClusterEndpointHost := ""
			byoClusterEndpointPort := 0
			if cloudService == domain.CloudService_BYOH {
				if clusterEndpoint == "" {
					return fmt.Errorf("invalid clusterEndpoint")
				}

				arr := strings.Split(clusterEndpoint, ":")
				if len(arr) != 2 {
					return fmt.Errorf("invalid clusterEndpoint")
				}
				byoClusterEndpointHost = arr[0]
				byoClusterEndpointPort, _ = strconv.Atoi(arr[1])
			}

			clusterDomains := make([]domain.ClusterDomain, len(domains))
			for i, domain := range domains {
				arrDomain := strings.Split(domain, "|")
				if len(arrDomain) > 0 {
					clusterDomains[i].DomainType = arrDomain[0]
					clusterDomains[i].Url = arrDomain[1]
				}
			}

			input := domain.CreateClusterRequest{
				Name:                   name,
				Description:            description,
				ClusterType:            clusterType,
				CloudService:           cloudService,
				OrganizationId:         organizationId,
				StackTemplateId:        stackTemplateId,
				CloudAccountId:         cloudAccountId,
				ByoClusterEndpointHost: byoClusterEndpointHost,
				ByoClusterEndpointPort: byoClusterEndpointPort,
				IsStack:                isStack,
				TksCpNode:              tksCpNode,
				TksCpNodeMax:           tksCpNodeMax,
				TksCpNodeType:          tksCpNodeType,
				TksInfraNode:           tksInfraNode,
				TksInfraNodeMax:        tksInfraNodeMax,
				TksInfraNodeType:       tksInfraNodeType,
				TksUserNode:            tksUserNode,
				TksUserNodeMax:         tksUserNodeMax,
				TksUserNodeType:        tksUserNodeType,
				PolicyIds:              policyIds,
				Domains:                clusterDomains,
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Post("clusters", input)
			if err != nil {
				return err
			}

			var out domain.CreateClusterResponse
			helper.Transcode(body, &out)

			fmt.Println("clusterId : ", out.ID)

			return nil
		},
	}

	command.Flags().StringVarP(&organizationId, "organization-id", "o", "", "the organizationId with clusters")
	helper.CheckError(command.MarkFlagRequired("organization-id"))

	command.Flags().StringVar(&cloudService, "cloud-service", "AWS", "the cloud service for cluster (AWS | BYOH)")
	command.Flags().StringVar(&clusterType, "cluster-type", "USER", "the cluster type (USER | ADMIN)")

	command.Flags().StringVarP(&cloudAccountId, "cloud-account-id", "s", "", "the cloudAccountId for cluster")

	command.Flags().StringVarP(&stackTemplateId, "stack-template-id", "t", "", "the template for installation")
	helper.CheckError(command.MarkFlagRequired("stack-template-id"))

	command.Flags().StringVarP(&name, "name", "n", "", "the name of organization")
	command.Flags().StringVarP(&description, "description", "d", "", "the description of organization")

	command.Flags().IntVar(&tksCpNode, "tks-cp-node", 3, "number of control-plane nodes")
	command.Flags().IntVar(&tksCpNodeMax, "tks-cp-node-max", 3, "max number of control-plane nodes")
	command.Flags().StringVar(&tksCpNodeType, "tks-cp-node-type", "t3.large", "machine type for tks cp node")

	command.Flags().IntVar(&tksInfraNode, "tks-infra-node", 3, "number of tks infra nodes")
	command.Flags().IntVar(&tksInfraNodeMax, "tks-infra-node-max", 3, "max number of tks infra nodes")
	command.Flags().StringVar(&tksInfraNodeType, "tks-infra-node-type", "t3.2xlarge", "machine type for tks infra node")

	command.Flags().IntVar(&tksUserNode, "tks-user-node", 1, "number of user nodes")
	command.Flags().IntVar(&tksUserNodeMax, "tks-user-node-max", 1, "max number of user nodes")
	command.Flags().StringVar(&tksUserNodeType, "tks-user-node-type", "t3.large", "machine type for user node")

	command.Flags().IntVar(&stack, "stack", 0, "enable creating stack")

	command.Flags().StringVar(&clusterEndpoint, "cluster-endpoint", "", "cluster endpoint host for byoh")

	command.Flags().StringSliceVar(&policyIds, "policy-ids", []string{}, "ex. policy_id1,policy_id1")

	command.Flags().StringSliceVar(&domains, "domains", []string{}, "ex. grafana_1.1.1.1:30001,thanos_1.1.1.1:30002")

	return command
}
