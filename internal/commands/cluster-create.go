package commands

import (
	"fmt"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewClusterCreateCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		name             string
		organizationId   string
		stackId          string
		description      string
		stackTemplateId  string
		cloudAccountId   string
		tksCpNode        int
		tksCpNodeMax     int
		tksCpNodeType    string
		tksInfraNode     int
		tksInfraNodeMax  int
		tksInfraNodeType string
		tksUserNode      int
		tksUserNodeMax   int
		tksUserNodeType  string
	)

	var command = &cobra.Command{
		Use:   "create",
		Short: "Create a TKS Cluster.",
		Long: `Create a TKS Cluster.
	  
	Example:
	tks cluster create <CLUSTERNAME> [--template TEMPLATE_NAME]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}

			if name == "" {
				helper.PanicWithError("You must specify name")
			}

			input := domain.CreateClusterRequest{
				OrganizationId:   organizationId,
				StackId:          domain.StackId(stackId),
				StackTemplateId:  stackTemplateId,
				Name:             name,
				Description:      description,
				CloudAccountId:   cloudAccountId,
				TksCpNode:        tksCpNode,
				TksCpNodeMax:     tksCpNodeMax,
				TksCpNodeType:    tksCpNodeType,
				TksInfraNode:     tksInfraNode,
				TksInfraNodeMax:  tksInfraNodeMax,
				TksInfraNodeType: tksInfraNodeType,
				TksUserNode:      tksUserNode,
				TksUserNodeMax:   tksUserNodeMax,
				TksUserNodeType:  tksUserNodeType,
			}

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
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

	command.Flags().StringVarP(&stackId, "stack-id", "s", "", "the stackId")

	command.Flags().StringVarP(&cloudAccountId, "cloud-account-id", "s", "", "the cloudAccountId for cluster")
	helper.CheckError(command.MarkFlagRequired("cloud-account-id"))

	command.Flags().StringVarP(&stackTemplateId, "stack-template-id", "t", "", "the template for installation")
	helper.CheckError(command.MarkFlagRequired("stack-template-id"))

	command.Flags().StringVarP(&name, "name", "n", "", "the name of organization")
	command.Flags().StringVarP(&description, "description", "d", "", "the description of organization")

	command.Flags().IntVar(&tksCpNode, "tks-cp-node", 0, "number of control-plane nodes")
	command.Flags().IntVar(&tksCpNodeMax, "tks-cp-node-max", 0, "max number of control-plane nodes")
	command.Flags().StringVar(&tksCpNodeType, "tks-cp-node-type", "t3.large", "machine type for tks cp node")

	command.Flags().IntVar(&tksInfraNode, "tks-infra-node", 1, "number of tks infra nodes")
	command.Flags().IntVar(&tksInfraNodeMax, "tks-infra-node-max", 1, "max number of tks infra nodes")
	command.Flags().StringVar(&tksInfraNodeType, "tks-infra-node-type", "t3.2xlarge", "machine type for tks infra node")

	command.Flags().IntVar(&tksUserNode, "tks-user-node", 1, "number of user nodes")
	command.Flags().IntVar(&tksUserNodeMax, "tks-user-node-max", 1, "max number of user nodes")
	command.Flags().StringVar(&tksUserNodeType, "tks-user-node-type", "t3.large", "machine type for user node")

	return command
}
