package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_apiClient "github.com/openinfradev/tks-api/pkg/api-client"
	"github.com/openinfradev/tks-client/internal/helper"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func NewAppserveUpdateCmd(globalOpts *GlobalOptions) *cobra.Command {
	var (
		appserveCfgFile string
		appCfgFile      string
		appSecretFile   string
	)

	var command = &cobra.Command{
		Use:   "update",
		Short: "Update an app deployed by AppServing service",
		Long: `Update an app deployed by AppServing service.
  
	Example:
	tks appserve update <APP_ID> --appserve-config <CONFIGFILE>`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			appId := args[0]

			if len(appId) < 1 {
				return errors.New("APP_ID is mandatory! Run 'tks appserve update --help'")
			}

			var c conf
			if appserveCfgFile == "" {
				return errors.New("--appserve-config is mandatory param")
			}

			// Get Appserving request params from config file
			yamlData, err := os.ReadFile(appserveCfgFile)
			if err != nil {
				return fmt.Errorf("error: %s", err)
			}

			fmt.Printf("*******\nConfig:\n%s\n*******\n", yamlData)

			// Get application config from file
			if appCfgFile != "" {
				appCfgBytes, err := os.ReadFile(appCfgFile)
				if err != nil {
					return fmt.Errorf("error: %s", err)
				}
				// Add appCfg to existing struct
				c.AppConfig = string(appCfgBytes)
			}

			// Get application secret from file
			if appSecretFile != "" {
				appSecretBytes, err := os.ReadFile(appSecretFile)
				if err != nil {
					return fmt.Errorf("error: %s", err)
				}
				c.AppSecret = string(appSecretBytes)
			}

			// Unmarshal yaml content into struct
			if err = yaml.Unmarshal(yamlData, &c); err != nil {
				return fmt.Errorf("error: %s", err)
			}

			// Exclude immutable params
			if c.Name != "" || c.Type != "" || c.AppType != "" {
				fmt.Println("Following params are immutable in update action, " +
					"so they're ignored.\n\t- name\n\t- type\n\t- app_type\n\t- target_cluster")
				c.Name = ""
				c.Type = ""
				c.AppType = ""
				c.TargetClusterId = ""
			}

			// Convert map to Json
			cBytes, err := json.Marshal(&c)
			if err != nil {
				return fmt.Errorf("unable to marshal config to JSON: %s", err)
			}

			fmt.Println("========== ")
			fmt.Println("Json data: ")
			fmt.Println(string(cBytes))
			fmt.Println("========== ")

			apiClient, err := _apiClient.New(globalOpts.ServerAddr, globalOpts.AuthToken)
			helper.CheckError(err)

			url := fmt.Sprintf("app-serve-apps/%v", appId)

			body, err := apiClient.Put(url, c)
			if err != nil {
				return err
			}

			fmt.Println("Response:\n", fmt.Sprintf("%v", body))

			return nil
		},
	}

	command.Flags().StringVar(&appserveCfgFile, "appserve-config", "", "config file for AppServing service")
	command.Flags().StringVar(&appCfgFile, "app-config", "", "custom config file for user application")
	command.Flags().StringVar(&appSecretFile, "app-secret", "", "custom secret file for user application")

	return command
}
