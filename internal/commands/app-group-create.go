package commands

import (
	"fmt"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewAppGroupCreateCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		name         string
		description  string
		clusterId    string
		appGroupType string
	)

	var command = &cobra.Command{
		Use:   "create",
		Short: "Create a TKS AppGroup.",
		Long: `Create a TKS AppGroup.
	  
	Example:
	tks appgroup create <APP_GROUP_NAME> [--cluster-id CLUSTER_ID --type APPGROUP_TYPE]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}

			if name == "" {
				helper.PanicWithError("You must specify name")
			}

			input := domain.CreateAppGroupRequest{
				Name:         name,
				Description:  description,
				ClusterId:    domain.ClusterId(clusterId),
				AppGroupType: appGroupType,
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			body, err := apiClient.Post("app-groups", input)
			if err != nil {
				return err
			}

			var out domain.CreateAppGroupResponse
			helper.Transcode(body, &out)

			fmt.Println("appGroupId : ", out.ID)

			return nil
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "clusterId")
	helper.CheckError(command.MarkFlagRequired("cluster-id"))

	command.Flags().StringVarP(&name, "name", "n", "", "the name of appGroup")

	command.Flags().StringVarP(&description, "description", "d", "", "the description of appGroup")
	command.Flags().StringVarP(&appGroupType, "type", "t", "LMA", "the type of appgroup. ex) LMA, SERVICE_MESH")
	helper.CheckError(command.MarkFlagRequired("type"))

	return command
}
