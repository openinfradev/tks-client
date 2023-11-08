package commands

import (
	"fmt"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewCloudAccountCreateCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		organizationId  string
		name            string
		description     string
		cloudService    string
		awsAccountId    string
		accessKeyId     string
		secretAccessKey string
	)

	var command = &cobra.Command{
		Use:   "create",
		Short: "Create a cloud account.",
		Long: `Create a cloud account.
	  
	Example:
	tks cloud-account create <NAME> [--awsAccountId AWS_ACCOUNT_ID --accessKeyID ACCESS_KEY_ID --secretAccesssKey SECRET_ACCESS_KEY]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}

			if name == "" {
				helper.PanicWithError("You must specify name")
			}

			input := domain.CreateCloudAccountRequest{
				Name:            name,
				Description:     description,
				CloudService:    cloudService,
				AwsAccountId:    awsAccountId,
				AccessKeyId:     accessKeyId,
				SecretAccessKey: secretAccessKey,
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/cloud-accounts", organizationId)
			body, err := apiClient.Post(url, input)
			if err != nil {
				return err
			}

			var out domain.CreateCloudAccountResponse
			helper.Transcode(body, &out)

			fmt.Println("cloudAccountId : ", out.ID)

			return nil
		},
	}

	command.Flags().StringVarP(&organizationId, "organization-id", "o", "", "the organizationId with cloud accounts")
	helper.CheckError(command.MarkFlagRequired("organization-id"))

	command.Flags().StringVarP(&name, "name", "n", "", "the name of cloud account")
	command.Flags().StringVarP(&description, "description", "d", "", "the description of cloud account")
	command.Flags().StringVar(&cloudService, "cloud-service", "AWS", "the type of cloud account")

	command.Flags().StringVar(&awsAccountId, "aws-account-id", "", "The accountId of aws")
	helper.CheckError(command.MarkFlagRequired("aws-account-id"))

	command.Flags().StringVar(&accessKeyId, "access-key-id", "", "The accessKeyId of aws")
	helper.CheckError(command.MarkFlagRequired("access-key-id"))

	command.Flags().StringVar(&secretAccessKey, "secret-access-key", "", "The secret access key of accessKey Id")
	helper.CheckError(command.MarkFlagRequired("secret-access-key"))

	return command
}
