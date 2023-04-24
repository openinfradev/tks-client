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
		name                string
		organizationId      string
		description         string
		stackTemplateId     string
		cloudAccountId      string
		cpNodeCnt           int
		cpNodeMachineType   string
		tksNodeCnt          int
		tksNodeMachineType  string
		userNodeCnt         int
		userNodeMachineType string
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
				OrganizationId:      organizationId,
				StackTemplateId:     stackTemplateId,
				Name:                name,
				Description:         description,
				CloudAccountId:      cloudAccountId,
				CpNodeCnt:           cpNodeCnt,
				CpNodeMachineType:   cpNodeMachineType,
				TksNodeCnt:          tksNodeCnt,
				TksNodeMachineType:  tksNodeMachineType,
				UserNodeCnt:         userNodeCnt,
				UserNodeMachineType: userNodeMachineType,
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

	command.Flags().StringVarP(&cloudAccountId, "cloud-account-id", "s", "", "the cloudAccountId for cluster")
	helper.CheckError(command.MarkFlagRequired("cloud-account-id"))

	command.Flags().StringVarP(&stackTemplateId, "stack-template-id", "t", "", "the template for installation")
	helper.CheckError(command.MarkFlagRequired("stack-template-id"))

	command.Flags().StringVarP(&name, "name", "n", "", "the name of organization")
	command.Flags().StringVarP(&description, "description", "d", "", "the description of organization")

	command.Flags().IntVar(&cpNodeCnt, "cp-node-cnt", 3, "number of control-plane nodes")
	command.Flags().StringVar(&cpNodeMachineType, "cp-node-machine-type", "t3.large", "machine type for tks cp node")
	command.Flags().IntVar(&tksNodeCnt, "tks-node-cnt", 3, "number of tks nodes")
	command.Flags().StringVar(&tksNodeMachineType, "tks-node-machine-type", "t3.2xlarge", "machine type for tks node")
	command.Flags().IntVar(&userNodeCnt, "user-node-cnt", 1, "number of control-plane nodes")
	command.Flags().StringVar(&userNodeMachineType, "user-node-machine-type", "t3.large", "machine type for user node")

	return command
}
