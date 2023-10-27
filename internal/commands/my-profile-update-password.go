package commands

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
)

func NewMyProfileUpdatePasswordCommand(globalOpts *GlobalOptions) *cobra.Command {
	var (
		originPassword string
		newPassword    string
	)

	var command = &cobra.Command{
		Use:   "update-password",
		Short: "update my password.",
		Long: `Update my password.
	
	Example:
	tks my-profile update-password --origin-password <ORIGIN_PASSWORD> --new-password <NEW_PASSWORD>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if globalOpts.CurrentOrganizationId == "" {
				return errors.New("current organization is not set")
			}
			originPassword, newPassword := prompt("origin-password", "new-password")

			input := domain.UpdatePasswordRequest{
				NewPassword:    newPassword,
				OriginPassword: originPassword,
			}

			apiClient, err := _apiClient.NewWithToken(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("organizations/%v/my-profile/password", globalOpts.CurrentOrganizationId)
			body, err := apiClient.Put(url, input)
			if err != nil {
				return err
			}

			helper.Transcode(body, &struct{}{})

			return nil
		},
	}

	command.Flags().StringVar(&originPassword, "origin-password", "", "a current Password")
	command.Flags().StringVar(&newPassword, "new-password", "", "a new Password")

	return command
}

func prompt(names ...string) (string, string) {
	var values []string
	for _, name := range names {
		if strings.Contains(name, "password") {
			fmt.Printf("%s: ", name)
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			helper.CheckError(err)
			value := string(bytePassword)
			values = append(values, strings.TrimSpace(value))
			fmt.Print("\n")
			continue
		}
		fmt.Printf("%s: ", name)
		reader := bufio.NewReader(os.Stdin)
		value, _ := reader.ReadString('\n')
		values = append(values, strings.TrimSpace(value))
	}
	return values[0], values[1]
}
