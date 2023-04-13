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
		name            string
		organizationId  string
		description     string
		template        string
		region          string
		cloudSettingId  string
		machineType     string
		numOfAz         int
		machineReplicas int
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
				OrganizationId:  organizationId,
				TemplateId:      template,
				Name:            name,
				Description:     description,
				CloudSettingId:  cloudSettingId,
				NumOfAz:         numOfAz,
				MachineType:     machineType,
				Region:          region,
				MachineReplicas: machineReplicas,
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

	command.Flags().StringVarP(&cloudSettingId, "cloud-setting-id", "s", "", "the cloudSettingId for cluster")
	helper.CheckError(command.MarkFlagRequired("cloud-setting-id"))

	command.Flags().StringVarP(&name, "name", "n", "", "the name of organization")
	command.Flags().StringVarP(&description, "description", "d", "", "the description of organization")
	command.Flags().StringVar(&template, "template", "aws-reference", "the template for installation")

	command.Flags().StringVar(&region, "region", "ap-northeast-2", "AWS region")
	command.Flags().StringVar(&machineType, "machine-type", "", "machine type for user node")
	command.Flags().IntVar(&numOfAz, "num-of-az", 1, "number of available zone")
	command.Flags().IntVar(&machineReplicas, "machine-replicas", 1, "the number of machine replica")

	return command
}
