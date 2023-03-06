package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/config"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func NewLoginCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		accountId string
		password  string
	)

	var command = &cobra.Command{
		Use:   "login SERVER",
		Short: "Log in to TKS",
		Long:  "Log in to TKS",
		Example: `# Login to TKS using a accountId and password
	tks login tks-api.tks.io`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println("You must specify server.")
				return fmt.Errorf("Usage: tks login [TKS_SERVER]")
			}
			server := args[0]

			accountId, password := PromptCredentials(accountId, password)
			input := domain.SignInRequest{
				AccountId: accountId,
				Password:  password,
			}

			var err error
			apiClient, err := _apiClient.New(server, "")
			helper.CheckError(err)

			body, err := apiClient.Post("auth/signin", input)
			if err != nil {
				return err
			}

			type DataInterface struct {
				User domain.User `json:"user"`
			}
			var out = DataInterface{}
			helper.Transcode(body, &out)

			fmt.Println(globalOpts.ConfigPath)
			localCfg, err := config.ReadLocalConfig(globalOpts.ConfigPath)
			helper.CheckError(err)
			if localCfg == nil {
				localCfg = &config.LocalConfig{}
			}

			localCfg.UpsertServer(config.Server{
				Server: server,
			})

			localCfg.UpsertUser(config.User{
				Name:         accountId,
				AuthToken:    out.User.Token,
				RefreshToken: "TODO",
			})

			err = config.WriteLocalConfig(*localCfg, globalOpts.ConfigPath)
			helper.CheckError(err)
			fmt.Printf("The user [%s] login successfully\n", accountId)

			return nil
		},
	}

	command.Flags().StringVar(&accountId, "accountId", "", "the accountId of an account to authenticate")
	command.Flags().StringVar(&password, "password", "", "the password of an account to authenticate")

	return command
}

func PromptCredentials(accountId string, password string) (string, string) {
	return PromptUsername(accountId), PromptPassword(password)
}

func PromptUsername(accountId string) string {
	return PromptMessage("AccountId", accountId)
}

func PromptMessage(message, value string) string {
	for value == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message + ": ")
		valueRaw, err := reader.ReadString('\n')
		helper.CheckError(err)
		value = strings.TrimSpace(valueRaw)
	}
	return value
}

func PromptPassword(password string) string {
	for password == "" {
		fmt.Print("Password: ")
		passwordRaw, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		helper.CheckError(err)
		password = string(passwordRaw)
		fmt.Print("\n")
	}
	return password
}
