package commands

import (
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
		creator         string
		template        string
		region          string
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

			type CreateClusterRequest struct {
				OrganizationId  string `json:"organizationId"`
				TemplateId      string `json:"templateId"`
				Name            string `json:"name"`
				Description     string `json:"description"`
				NumberOfAz      string `json:"numberOfAz"`
				MachineType     string `json:"machineType"`
				Region          string `json:"region"`
				MachineReplicas int    `json:"machineReplicas"`
			}

			input := domain.CreateClusterRequest{
				OrganizationId:  organizationId,
				TemplateId:      template,
				Name:            name,
				Description:     description,
				NumberOfAz:      numOfAz,
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

			type DataInterface struct {
				ClusterId string `json:"clusterId"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			return nil
		},
	}

	command.Flags().StringVar(&organizationId, "organization-id", "", "the organizationId with clusters")
	helper.CheckError(command.MarkFlagRequired("organization-id"))

	command.Flags().StringVar(&name, "name", "", "the name of organization")
	command.Flags().StringVar(&description, "description", "", "the description of organization")
	command.Flags().StringVar(&creator, "creator", "", "the user's uuid for creating organization")
	command.Flags().StringVar(&template, "template", "aws-reference", "the template for installation")

	command.Flags().StringVar(&region, "region", "ap-northeast-2", "AWS region")
	command.Flags().StringVar(&machineType, "machine-type", "", "machine type for user node")
	command.Flags().IntVar(&numOfAz, "num-of-az", 1, "number of available zone")
	command.Flags().IntVar(&machineReplicas, "machine-replicas", 1, "the number of machine replica")

	return command
}
